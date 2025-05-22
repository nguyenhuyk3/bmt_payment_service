package services

import (
	"bmt_payment_service/dto/request"
	"context"
)

type IPayment interface {
	CreatePaymentURL(ctx context.Context, arg request.CreatePaymentURLReq) (string, int, error)
	CreatePaymentRecord(ctx context.Context, arg request.CreatePaymentRecordReq) (interface{}, int, error)
	HandleWebhook()
	Refund()
}
