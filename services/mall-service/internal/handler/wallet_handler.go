package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/mall-service/internal/service"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletSvc *service.WalletService
}

func NewWalletHandler(walletSvc *service.WalletService) *WalletHandler {
	return &WalletHandler{walletSvc: walletSvc}
}

// POST /api/mall/wallet/recharge (MALL-019)
func (h *WalletHandler) Recharge(c *gin.Context) {
	userID := c.GetInt64("userID")
	var req struct {
		Amount         float64 `json:"amount" binding:"required"`
		IdempotencyKey string  `json:"idempotency_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.walletSvc.Recharge(userID, centsFromYuan(req.Amount), req.IdempotencyKey); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// POST /api/mall/wallet/transfer (MALL-020)
func (h *WalletHandler) Transfer(c *gin.Context) {
	userID := c.GetInt64("userID")
	var req struct {
		ToUserID       int64   `json:"to_user_id" binding:"required"`
		Amount         float64 `json:"amount" binding:"required"`
		IdempotencyKey string  `json:"idempotency_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.walletSvc.Transfer(userID, req.ToUserID, centsFromYuan(req.Amount), req.IdempotencyKey); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// GET /api/mall/wallet/balance (MALL-021)
func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID := c.GetInt64("userID")

	balance, err := h.walletSvc.GetBalance(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{"balance": yuanFromCents(balance)})
}

// GET /api/mall/wallet/transactions (MALL-021)
func (h *WalletHandler) ListTransactions(c *gin.Context) {
	userID := c.GetInt64("userID")
	page, size := response.ParsePage(c)

	var txType *int
	if t := c.Query("type"); t != "" {
		v, err := strconv.Atoi(t)
		if err == nil {
			txType = &v
		}
	}

	txs, total, err := h.walletSvc.ListTransactions(userID, page, size, txType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, toWalletTransactionResponses(txs), total, page, size)
}
