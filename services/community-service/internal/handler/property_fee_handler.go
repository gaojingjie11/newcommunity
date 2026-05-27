package handler

import (
	"errors"
	"net/http"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/community-service/internal/repository"
	"smartcommunity-microservices/services/community-service/internal/service"

	"github.com/gin-gonic/gin"
)

type PropertyFeeHandler struct {
	svc *service.PropertyFeeService
}

func NewPropertyFeeHandler(svc *service.PropertyFeeService) *PropertyFeeHandler {
	return &PropertyFeeHandler{svc: svc}
}

func (h *PropertyFeeHandler) MyFees(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	page, size := response.ParsePage(c)
	items, total, err := h.svc.ListByUser(userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询物业费失败", nil)
		return
	}
	response.SuccessPaged(c, items, total, page, size)
}

func (h *PropertyFeeHandler) Pay(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req service.PayPropertyFeeRequest
	_ = c.ShouldBindJSON(&req)
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = c.GetHeader("Idempotency-Key")
	}
	result, err := h.svc.Pay(userID, id, req)
	if err != nil {
		if errors.Is(err, repository.ErrPropertyFeePaid) {
			response.Error(c, http.StatusConflict, 409, "物业费已缴纳", nil)
			return
		}
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{"payment_result": result})
}

func (h *PropertyFeeHandler) MyPayments(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	page, size := response.ParsePage(c)
	items, total, err := h.svc.ListPaymentsByUser(userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询缴费记录失败", nil)
		return
	}
	response.SuccessPaged(c, items, total, page, size)
}

func (h *PropertyFeeHandler) AdminCreate(c *gin.Context) {
	var req service.CreatePropertyFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}
	item, err := h.svc.Create(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "创建物业费失败: "+err.Error(), nil)
		return
	}
	response.Success(c, item)
}

func (h *PropertyFeeHandler) AdminList(c *gin.Context) {
	page, size := response.ParsePage(c)
	items, total, err := h.svc.ListAll(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询物业费失败", nil)
		return
	}
	response.SuccessPaged(c, items, total, page, size)
}

func (h *PropertyFeeHandler) AdminPayments(c *gin.Context) {
	page, size := response.ParsePage(c)
	items, total, err := h.svc.ListPayments(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询缴费记录失败", nil)
		return
	}
	response.SuccessPaged(c, items, total, page, size)
}
