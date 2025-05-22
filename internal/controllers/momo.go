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
	"github.com/sony/sonyflake"
)

type MoMoController struct {
	MoMoService services.IPayment
	Flake       *sonyflake.Sonyflake
}

func (m *MoMoController) CreatePaymentURL(c *gin.Context) {
	var req request.PaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.FailureResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	orderIDInt, _ := m.Flake.NextID()
	req.OrderId = orderIDInt
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
	// orderId := c.Query("orderId")
	// requestID := c.Query("requestId")
	// amount := c.Query("amount")
	// orderInfo := c.Query("orderInfo")
	// orderType := c.Query("orderType")
	// transId := c.Query("transId")
	resultCode := c.Query("resultCode")
	// message := c.Query("message")
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

	if resultCode == "0" {
		responses.SuccessResponse(c, http.StatusOK, "pay successfully", paymentInfo)
	} else {
		responses.FailureResponse(c, http.StatusInternalServerError, "pay failed")
	}
}

func NewMoMoController(
	moMoService services.IPayment,
) *MoMoController {
	return &MoMoController{
		MoMoService: moMoService,
		Flake:       sonyflake.NewSonyflake(sonyflake.Settings{}),
	}
}
