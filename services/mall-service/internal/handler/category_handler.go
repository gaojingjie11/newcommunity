package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/mall-service/internal/model"
	"smartcommunity-microservices/services/mall-service/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categorySvc *service.CategoryService
}

func NewCategoryHandler(categorySvc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categorySvc: categorySvc}
}

// GET /api/mall/categories
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.categorySvc.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, categories, int64(len(categories)), 1, len(categories))
}

// GET /api/mall/categories/:id
func (h *CategoryHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	cat, err := h.categorySvc.GetDetail(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, "分类不存在", nil)
		return
	}
	response.Success(c, cat)
}

// POST /api/admin/mall/categories (ADMIN-MALL-005)
func (h *CategoryHandler) Create(c *gin.Context) {
	var cat model.ProductCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.categorySvc.Create(&cat); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{"id": cat.ID})
}

// PUT /api/admin/mall/categories/:id (ADMIN-MALL-005)
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	var cat model.ProductCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	cat.ID = id

	if err := h.categorySvc.Update(&cat); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// DELETE /api/admin/mall/categories/:id (ADMIN-MALL-005)
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	if err := h.categorySvc.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}
