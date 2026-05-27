package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/mall-service/internal/service"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	favSvc *service.FavoriteService
}

func NewFavoriteHandler(favSvc *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{favSvc: favSvc}
}

// POST /api/mall/favorites (MALL-008)
func (h *FavoriteHandler) Add(c *gin.Context) {
	userID := c.GetInt64("userID")
	var req struct {
		ProductID int64 `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.favSvc.Add(userID, req.ProductID); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// DELETE /api/mall/favorites/:product_id (MALL-009)
func (h *FavoriteHandler) Remove(c *gin.Context) {
	userID := c.GetInt64("userID")
	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	if err := h.favSvc.Remove(userID, productID); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// GET /api/mall/favorites (MALL-018)
func (h *FavoriteHandler) List(c *gin.Context) {
	userID := c.GetInt64("userID")
	page, size := response.ParsePage(c)

	favorites, total, err := h.favSvc.List(userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, toFavoriteResponses(favorites), total, page, size)
}

// GET /api/mall/favorites/check/:product_id
func (h *FavoriteHandler) Check(c *gin.Context) {
	userID := c.GetInt64("userID")
	productID, err := strconv.ParseInt(c.Param("product_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	exists, err := h.favSvc.Check(userID, productID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{"is_favorite": exists})
}
