package services

import (
	"bmt_payment_service/dto/request"
	"context"
)

type IPayment interface {
	CreatePaymentURL(ctx context.Context, arg request.PaymentReq) (string, int, error)
	VerifyPaymentCallback(data map[string]string) (bool, error)
	HandleWebhook(requestBody []byte) (request.PaymentResult, error)
	Refund(orderID string, amount int64) error
}
