package momo

import (
	"bmt_payment_service/dto/request"
	"bmt_payment_service/dto/response"
	"bmt_payment_service/internal/services"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/sony/sonyflake"
)

type moMoPayment struct {
	Endpoint    string
	PartnerCode string
	AccessKey   string
	SecretKey   string
	RedirectURL string
	IpnURL      string
}

// CreatePaymentURL implements services.IPayment.
func (m *moMoPayment) CreatePaymentURL(ctx context.Context, arg request.PaymentReq) (string, int, error) {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	requestIDInt, _ := flake.NextID()

	orderID := strconv.FormatUint(arg.OrderId, 10)
	orderInfo := "Thanh toán vé xem phim"
	requestID := strconv.FormatUint(requestIDInt, 10)
	amountStr := strconv.FormatInt(arg.Amount, 10)

	rawSignature := fmt.Sprintf("accessKey=%s&amount=%s&extraData=%s&ipnUrl=%s&orderId=%s&orderInfo=%s&partnerCode=%s&redirectUrl=%s&requestId=%s&requestType=payWithMethod",
		m.AccessKey, amountStr, arg.ExtraData, m.IpnURL, orderID, orderInfo, m.PartnerCode, m.RedirectURL, requestID)

	h := hmac.New(sha256.New, []byte(m.SecretKey))
	h.Write([]byte(rawSignature))

	signature := hex.EncodeToString(h.Sum(nil))

	payload := request.MoMoPayload{
		PartnerCode:  m.PartnerCode,
		AccessKey:    m.AccessKey,
		RequestID:    requestID,
		Amount:       amountStr,
		OrderID:      orderID,
		OrderInfo:    orderInfo,
		PartnerName:  "MoMo Payment",
		StoreId:      "Test Store",
		OrderGroupId: "",
		Lang:         "vi",
		AutoCapture:  true,
		RedirectUrl:  m.RedirectURL,
		IpnUrl:       m.IpnURL,
		ExtraData:    arg.ExtraData,
		RequestType:  "payWithMethod",
		Signature:    signature,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	resp, err := http.Post(m.Endpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	var momoResp response.MoMoRes
	if err := json.Unmarshal(body, &momoResp); err != nil {
		return "", http.StatusInternalServerError, err
	}

	return momoResp.PayURL, http.StatusOK, nil
}

// HandleWebhook implements services.IPayment.
func (m *moMoPayment) HandleWebhook(requestBody []byte) (request.PaymentResult, error) {
	panic("unimplemented")
}

// Refund implements services.IPayment.
func (m *moMoPayment) Refund(orderID string, amount int64) error {
	panic("unimplemented")
}

// VerifyPaymentCallback implements services.IPayment.
func (m *moMoPayment) VerifyPaymentCallback(data map[string]string) (bool, error) {
	/**
	resultCode = 0: giao dịch thành công.
	resultCode = 9000: giao dịch được cấp quyền (authorization) thành công .
	resultCode <> 0: giao dịch thất bại.
	*/
	/**
	  * Dựa vào kết quả này để update trạng thái đơn hàng
	  * Kết quả log:
	  * {
			partnerCode: 'MOMO',
			orderId: 'MOMO1712108682648',
			requestId: 'MOMO1712108682648',
			amount: 10000,
			orderInfo: 'pay with MoMo',
			orderType: 'momo_wallet',
			transId: 4014083433,
			resultCode: 0,
			message: 'Thành công.',
			payType: 'qr',
			responseTime: 1712108811069,
			extraData: '',
			signature: '10398fbe70cd3052f443da99f7c4befbf49ab0d0c6cd7dc14efffd6e09a526c0'
		}
	*/
	panic("unimplemented")
}

func NewMomoPayment(
	endPoint, partnerCode, accessKey, secretKey, redirectUrl, ipnUrl string,
) services.IPayment {
	return &moMoPayment{
		Endpoint:    endPoint,
		PartnerCode: partnerCode,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		RedirectURL: redirectUrl,
		IpnURL:      ipnUrl,
	}
}
