package sqlc

import (
	"bmt_payment_service/dto/request"
	"context"
)

type IStore interface {
	CreatePaymentTran(ctx context.Context, arg request.CreatePaymentRecordReq) (Payment, error)
}
