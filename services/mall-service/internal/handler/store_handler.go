package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/mall-service/internal/service"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeSvc *service.StoreService
}

func NewStoreHandler(storeSvc *service.StoreService) *StoreHandler {
	return &StoreHandler{storeSvc: storeSvc}
}

// GET /api/mall/stores (MALL-006)
func (h *StoreHandler) List(c *gin.Context) {
	page, size := response.ParsePage(c)
	areaID, _ := strconv.ParseInt(c.Query("area_id"), 10, 64)

	stores, total, err := h.storeSvc.List(page, size, areaID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, stores, total, page, size)
}

// GET /api/mall/stores/:id
func (h *StoreHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	store, err := h.storeSvc.GetDetail(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, "门店不存在", nil)
		return
	}
	response.Success(c, store)
}
