package momo

import (
	"bmt_payment_service/db/sqlc"
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
	SqlStore    sqlc.IStore
	Endpoint    string
	PartnerCode string
	AccessKey   string
	SecretKey   string
	RedirectURL string
	IpnURL      string
}

// CreatePaymentRecord implements services.IPayment.
func (m *moMoPayment) CreatePaymentRecord(ctx context.Context, arg request.CreatePaymentRecordReq) (interface{}, int, error) {
	payment, err := m.SqlStore.CreatePaymentTran(ctx, arg)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return payment, http.StatusOK, nil
}

// CreatePaymentURL implements services.IPayment.
func (m *moMoPayment) CreatePaymentURL(ctx context.Context, arg request.CreatePaymentURLReq) (string, int, error) {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	requestIDInt, _ := flake.NextID()
	orderIDInt, _ := flake.NextID()

	orderID := strconv.FormatUint(orderIDInt, 10)
	orderInfo := "Thanh toán vé xem phim"
	requestID := strconv.FormatUint(requestIDInt, 10)
	amountStr := strconv.FormatInt(arg.Amount, 10)

	rawSignature := fmt.Sprintf("accessKey=%s&amount=%s&extraData=%s&ipnUrl=%s&orderId=%s&orderInfo=%s&partnerCode=%s&redirectUrl=%s&requestId=%s&requestType=payWithMethod",
		m.AccessKey, amountStr, arg.ExtraData, m.IpnURL, orderID, orderInfo, m.PartnerCode, m.RedirectURL, requestID)

	h := hmac.New(sha256.New, []byte(m.SecretKey))
	h.Write([]byte(rawSignature))

	signature := hex.EncodeToString(h.Sum(nil))

	payload := moMoPayload{
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
func (m *moMoPayment) HandleWebhook() {
	panic("unimplemented")
}

// Refund implements services.IPayment.
func (m *moMoPayment) Refund() {
	panic("unimplemented")
}

func NewMomoPayment(
	sqlStore sqlc.IStore,
	endPoint, partnerCode, accessKey, secretKey, redirectUrl, ipnUrl string,
) services.IPayment {
	return &moMoPayment{
		SqlStore:    sqlStore,
		Endpoint:    endPoint,
		PartnerCode: partnerCode,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		RedirectURL: redirectUrl,
		IpnURL:      ipnUrl,
	}
}
