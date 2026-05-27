package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/mall-service/internal/service"

	"github.com/gin-gonic/gin"
)

type PromotionHandler struct {
	promoSvc *service.PromotionService
}

func NewPromotionHandler(promoSvc *service.PromotionService) *PromotionHandler {
	return &PromotionHandler{promoSvc: promoSvc}
}

// GET /api/mall/promotions
func (h *PromotionHandler) List(c *gin.Context) {
	page, size := response.ParsePage(c)

	promos, total, err := h.promoSvc.List(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, promos, total, page, size)
}

// GET /api/mall/promotions/:id
func (h *PromotionHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	promo, err := h.promoSvc.GetDetail(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, "促销不存在", nil)
		return
	}
	response.Success(c, promo)
}
