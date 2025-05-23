package controllers

import (
	"bmt_payment_service/dto/request"
	"bmt_payment_service/internal/responses"
	"bmt_payment_service/internal/services"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MoMoController struct {
	MoMoService services.IPayment
}

func (m *MoMoController) CreatePaymentURL(c *gin.Context) {
	var req request.CreatePaymentURLReq
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	paymentURL, status, err := m.MoMoService.CreatePaymentURL(ctx, req)
	if err != nil {
		responses.FailureResponse(c, status, err.Error())
		return
	}

	responses.SuccessResponse(c, status, "create momo payment url successfully", gin.H{
		"payment_url": paymentURL,
	})
}

func (m *MoMoController) VerifyPaymentCallback(c *gin.Context) {
	/**
	resultCode = 0: transaction successful.
	resultCode = 9000: transaction authorized successfully.
	resultCode ≠ 0: transaction failed.
	*/

	// orderId := c.Query("orderId")
	// requestID := c.Query("requestId")
	amount := c.Query("amount")
	// orderInfo := c.Query("orderInfo")
	// orderType := c.Query("orderType")
	transId := c.Query("transId")
	resultCode := c.Query("resultCode")
	message := c.Query("message")
	// payType := c.Query("payType")
	// responseTime := c.Query("responseTime")
	extraData := c.Query("extraData")
	// signature := c.Query("signature")
	// partnerCode := c.Query("partnerCode")

	var paymentInfo map[string]interface{}
	err := json.Unmarshal([]byte(extraData), &paymentInfo)
	if err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid json format: %s", extraData))
		return
	}

	orderIdFloat, ok := paymentInfo["order_id"].(float64)
	if !ok {
		responses.FailureResponse(c, http.StatusBadRequest, "invalid order id format in extraData")
		return
	}

	orderIdInt32 := int32(orderIdFloat)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if resultCode == "0" || resultCode == "9000" {
		req := request.CreatePaymentRecordReq{
			Amount:        amount,
			OrderId:       orderIdInt32,
			Status:        "success",
			Method:        "momo",
			Message:       message,
			TransactionId: transId,
		}
		data, status, err := m.MoMoService.CreatePaymentRecord(ctx, req)
		if err != nil {
			responses.FailureResponse(c, status, fmt.Sprintf("failed to pay with resultCode = 0: %v", err))
			return
		}

		responses.SuccessResponse(c, status, "pay successfully", data)
	} else {
		req := request.CreatePaymentRecordReq{
			Amount:        amount,
			OrderId:       orderIdInt32,
			Status:        "failed",
			Method:        "momo",
			Message:       message,
			TransactionId: transId,
		}
		_, status, err := m.MoMoService.CreatePaymentRecord(ctx, req)
		if err != nil {
			responses.FailureResponse(c, status, fmt.Sprintf("failed to pay with resultCode != 0: %v", err))
			return
		}

		responses.FailureResponse(c, http.StatusInternalServerError, "pay failed")
	}
}

func NewMoMoController(
	moMoService services.IPayment,
) *MoMoController {
	return &MoMoController{
		MoMoService: moMoService,
	}
}
