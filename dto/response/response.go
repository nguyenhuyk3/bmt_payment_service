package response

type MoMoRes struct {
	PayURL string `json:"payUrl"`
	// other fields like message, errorCode... can be added if needed
}

type MoMoPaymentRes struct {
	PartnerCode  string `json:"partnerCode"`
	OrderID      string `json:"orderId"`
	RequestID    string `json:"requestId"`
	Amount       int64  `json:"amount"`
	OrderInfo    string `json:"orderInfo"`
	OrderType    string `json:"orderType"`
	TransID      int64  `json:"transId"`
	ResultCode   int    `json:"resultCode"`
	Message      string `json:"message"`
	PayType      string `json:"payType"`
	ResponseTime int64  `json:"responseTime"`
	ExtraData    string `json:"extraData"`
	Signature    string `json:"signature"`
}
