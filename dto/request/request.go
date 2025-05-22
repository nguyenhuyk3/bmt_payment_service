package request

import "time"

type PaymentReq struct {
	Amount int64 `json:"amount" binding:"required"`
	// OrderInfo string `json:"order_info" binding:"required"`
	OrderId   uint64
	ExtraData string `json:"extra_data" binding:"required"`
}

type PaymentResult struct {
	OrderID     string
	Paid        bool
	Amount      int64
	PaymentTime time.Time
	RawData     map[string]string
}

type MoMoPayload struct {
	PartnerCode  string `json:"partnerCode"`
	AccessKey    string `json:"accessKey"`
	RequestID    string `json:"requestId"`
	Amount       string `json:"amount"`
	OrderID      string `json:"orderId"`
	OrderInfo    string `json:"orderInfo"`
	PartnerName  string `json:"partnerName"`
	StoreId      string `json:"storeId"`
	OrderGroupId string `json:"orderGroupId"`
	Lang         string `json:"lang"`
	AutoCapture  bool   `json:"autoCapture"`
	RedirectUrl  string `json:"redirectUrl"`
	IpnUrl       string `json:"ipnUrl"`
	ExtraData    string `json:"extraData"`
	RequestType  string `json:"requestType"`
	Signature    string `json:"signature"`
}
