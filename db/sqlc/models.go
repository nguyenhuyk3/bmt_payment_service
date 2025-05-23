// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type PaymentMethods string

const (
	PaymentMethodsMomo         PaymentMethods = "momo"
	PaymentMethodsVnpay        PaymentMethods = "vnpay"
	PaymentMethodsZalopay      PaymentMethods = "zalopay"
	PaymentMethodsCreditCard   PaymentMethods = "credit_card"
	PaymentMethodsBankTransfer PaymentMethods = "bank_transfer"
)

func (e *PaymentMethods) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentMethods(s)
	case string:
		*e = PaymentMethods(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentMethods: %T", src)
	}
	return nil
}

type NullPaymentMethods struct {
	PaymentMethods PaymentMethods `json:"payment_methods"`
	Valid          bool           `json:"valid"` // Valid is true if PaymentMethods is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentMethods) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentMethods, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentMethods.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentMethods) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentMethods), nil
}

type PaymentStatuses string

const (
	PaymentStatusesCreated  PaymentStatuses = "created"
	PaymentStatusesSuccess  PaymentStatuses = "success"
	PaymentStatusesFailed   PaymentStatuses = "failed"
	PaymentStatusesCanceled PaymentStatuses = "canceled"
	PaymentStatusesExpired  PaymentStatuses = "expired"
)

func (e *PaymentStatuses) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = PaymentStatuses(s)
	case string:
		*e = PaymentStatuses(s)
	default:
		return fmt.Errorf("unsupported scan type for PaymentStatuses: %T", src)
	}
	return nil
}

type NullPaymentStatuses struct {
	PaymentStatuses PaymentStatuses `json:"payment_statuses"`
	Valid           bool            `json:"valid"` // Valid is true if PaymentStatuses is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullPaymentStatuses) Scan(value interface{}) error {
	if value == nil {
		ns.PaymentStatuses, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.PaymentStatuses.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullPaymentStatuses) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.PaymentStatuses), nil
}

type Outbox struct {
	ID             pgtype.UUID      `json:"id"`
	AggregatedType string           `json:"aggregated_type"`
	AggregatedID   int32            `json:"aggregated_id"`
	EventType      string           `json:"event_type"`
	Payload        []byte           `json:"payload"`
	CreatedAt      pgtype.Timestamp `json:"created_at"`
}

type Payment struct {
	ID            int32            `json:"id"`
	OrderID       int32            `json:"order_id"`
	Amount        string           `json:"amount"`
	Status        PaymentStatuses  `json:"status"`
	Method        PaymentMethods   `json:"method"`
	TransactionID string           `json:"transaction_id"`
	ErrorMessage  pgtype.Text      `json:"error_message"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
}
