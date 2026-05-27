package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/mall-service/internal/service"

	"github.com/gin-gonic/gin"
)

// InternalHandler serves service-to-service endpoints with internal-only semantics.
// These are NOT the same as admin endpoints — they use InternalToken auth,
// have no RBAC checks, and accept minimal request bodies.
type InternalHandler struct {
	orderSvc   *service.OrderService
	paymentSvc *service.PaymentService
	walletSvc  *service.WalletService
}

func NewInternalHandler(orderSvc *service.OrderService, paymentSvc *service.PaymentService, walletSvc *service.WalletService) *InternalHandler {
	return &InternalHandler{orderSvc: orderSvc, paymentSvc: paymentSvc, walletSvc: walletSvc}
}

// CancelExpiredOrder is called by the scheduler/gateway to cancel an expired order.
// It verifies the order is still pending before cancelling.
func (h *InternalHandler) CancelExpiredOrder(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.Reason == "" {
		req.Reason = "系统超时自动取消"
	}

	if err := h.orderSvc.CancelOrder(orderID, req.Reason); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// GetPaymentStatus is a lightweight query for other services to check payment state.
func (h *InternalHandler) GetPaymentStatus(c *gin.Context) {
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

// DebitWalletRequest is the request body for the internal wallet debit endpoint.
type DebitWalletRequest struct {
	UserID         int64  `json:"user_id" binding:"required"`
	Amount         int64  `json:"amount" binding:"required"`
	BizType        string `json:"biz_type"`
	BizID          string `json:"biz_id"`
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
	Remark         string `json:"remark"`
	PayType        string `json:"pay_type"`
	Password       string `json:"password"`
	FaceImageURL   string `json:"face_image_url"`
}

// DebitWallet debits amount from a user's wallet. Used by other services (e.g. community-service for property fees).
func (h *InternalHandler) DebitWallet(c *gin.Context) {
	var req DebitWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	if req.PayType != "" || req.Password != "" || req.FaceImageURL != "" {
		if err := h.paymentSvc.ValidatePayAuth(req.UserID, service.PayOrderRequest{
			PayType:      req.PayType,
			Password:     req.Password,
			FaceImageURL: req.FaceImageURL,
		}); err != nil {
			response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
			return
		}
	}

	walletTxID, balBefore, balAfter, err := h.walletSvc.DebitForExternal(
		req.UserID, req.Amount, req.BizType, req.BizID, req.IdempotencyKey, req.Remark)
	if err != nil {
		if err.Error() == "余额不足" {
			response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		} else {
			response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		}
		return
	}

	response.Success(c, gin.H{
		"wallet_transaction_id": walletTxID,
		"balance_before":        balBefore,
		"balance_after":         balAfter,
	})
}
