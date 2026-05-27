package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/mall-service/internal/service"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentSvc *service.PaymentService
}

func NewPaymentHandler(paymentSvc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentSvc: paymentSvc}
}

// POST /api/mall/orders/:id/pay (MALL-013)
func (h *PaymentHandler) Pay(c *gin.Context) {
	userID := c.GetInt64("userID")
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	var req service.PayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "支付参数不完整，请重试", nil)
		return
	}

	result, err := h.paymentSvc.PayOrder(orderID, userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{"payment_result": result})
}

// GET /api/mall/orders/:id/payment-status
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	status, err := h.paymentSvc.GetPaymentStatus(orderID)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, err.Error(), nil)
		return
	}
	response.Success(c, status)
}
