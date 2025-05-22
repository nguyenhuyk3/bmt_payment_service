package sqlc

import (
	"bmt_payment_service/dto/request"
	"bmt_payment_service/global"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SqlStore struct {
	connPool *pgxpool.Pool
}

// CreatePaymentTran implements IStore.
func (s *SqlStore) CreatePaymentTran(ctx context.Context, arg request.CreatePaymentRecordReq) (Payment, error) {
	var payment Payment
	err := s.execTran(ctx, func(q *Queries) error {
		var status PaymentStatuses

		err := status.Scan(arg.Status)
		if err != nil {
			return fmt.Errorf("failed to scan status: %w", err)
		}

		var method PaymentMethods
		err = method.Scan(arg.Method)
		if err != nil {
			return fmt.Errorf("failed to scan method: %w", err)
		}

		paymentData, err := q.CreatePayment(ctx,
			CreatePaymentParams{
				OrderID:       arg.OrderId,
				Amount:        arg.Amount,
				Status:        status,
				Method:        method,
				TransactionID: arg.TransactionId,
				ErrorMessage: pgtype.Text{
					String: arg.Message,
					Valid:  true,
				},
			})
		if err != nil {
			return fmt.Errorf("failed to create payment: %w", err)
		}

		var eventType string = global.PAYMENT_FAILED

		if arg.Status == "success" {
			eventType = global.PAYMENT_SUCCESS
		}

		payloadBytes, err := json.Marshal(gin.H{
			"order_id": arg.OrderId,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
		err = q.CreateOutbox(ctx,
			CreateOutboxParams{
				AggregatedType: "PAYMENT",
				AggregatedID:   arg.OrderId,
				EventType:      eventType,
				Payload:        payloadBytes,
			})
		if err != nil {
			return fmt.Errorf("failed to create outbox: %w", err)
		}

		payment = paymentData

		return nil
	})

	return payment, err
}

func (s *SqlStore) execTran(ctx context.Context, fn func(*Queries) error) error {
	// Start transaction
	tran, err := s.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tran)
	// fn performs a series of operations down the db
	err = fn(q)
	if err != nil {
		// If an error occurs, rollback the transaction
		if rbErr := tran.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tran err: %v, rollback err: %v", err, rbErr)
		}

		return err
	}

	return tran.Commit(ctx)
}

func NewStore(connPool *pgxpool.Pool) IStore {
	return &SqlStore{
		connPool: connPool,
	}
}
