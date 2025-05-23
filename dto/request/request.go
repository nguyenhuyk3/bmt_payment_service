package request

import "time"

type CreatePaymentURLReq struct {
	OrderId int32 `json:"order_id" binding:"required"`
	Amount  int64 `json:"amount" binding:"required"`
	// OrderInfo string `json:"order_info" binding:"required"`
	// ExtraData string `json:"extra_data" binding:"required"`
}

type CreatePaymentRecordReq struct {
	Amount        string
	OrderId       int32
	Status        string
	Method        string
	Message       string
	TransactionId string
}

type PaymentResult struct {
	OrderID     string
	Paid        bool
	Amount      int64
	PaymentTime time.Time
	RawData     map[string]string
}
