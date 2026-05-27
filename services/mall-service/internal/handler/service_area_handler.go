package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/mall-service/internal/model"
	"smartcommunity-microservices/services/mall-service/internal/service"

	"github.com/gin-gonic/gin"
)

type ServiceAreaHandler struct {
	areaSvc *service.ServiceAreaService
}

func NewServiceAreaHandler(areaSvc *service.ServiceAreaService) *ServiceAreaHandler {
	return &ServiceAreaHandler{areaSvc: areaSvc}
}

// GET /api/mall/service-areas
func (h *ServiceAreaHandler) List(c *gin.Context) {
	areas, err := h.areaSvc.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, areas, int64(len(areas)), 1, len(areas))
}

// POST /api/admin/mall/service-areas (ADMIN-MALL-008)
func (h *ServiceAreaHandler) Create(c *gin.Context) {
	var area model.ServiceArea
	if err := c.ShouldBindJSON(&area); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.areaSvc.Create(&area); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{"id": area.ID})
}

// PUT /api/admin/mall/service-areas/:id (ADMIN-MALL-008)
func (h *ServiceAreaHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	var area model.ServiceArea
	if err := c.ShouldBindJSON(&area); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	area.ID = id

	if err := h.areaSvc.Update(&area); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// DELETE /api/admin/mall/service-areas/:id (ADMIN-MALL-008)
func (h *ServiceAreaHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	if err := h.areaSvc.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}
