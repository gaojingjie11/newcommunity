package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/mall-service/internal/service"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartSvc *service.CartService
}

func NewCartHandler(cartSvc *service.CartService) *CartHandler {
	return &CartHandler{cartSvc: cartSvc}
}

// POST /api/mall/cart/items (MALL-007)
func (h *CartHandler) Add(c *gin.Context) {
	userID := c.GetInt64("userID")
	var req struct {
		ProductID int64 `json:"product_id" binding:"required"`
		Quantity  int   `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.cartSvc.Add(userID, req.ProductID, req.Quantity); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// DELETE /api/mall/cart/items/:id (MALL-010)
func (h *CartHandler) Remove(c *gin.Context) {
	userID := c.GetInt64("userID")
	cartID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	if err := h.cartSvc.Remove(cartID, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// PUT /api/mall/cart/items/:id (MALL-011)
func (h *CartHandler) UpdateQuantity(c *gin.Context) {
	userID := c.GetInt64("userID")
	cartID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}
	var req struct {
		Quantity int64 `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.cartSvc.UpdateQuantity(cartID, userID, req.Quantity); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// GET /api/mall/cart/items
func (h *CartHandler) List(c *gin.Context) {
	userID := c.GetInt64("userID")

	items, err := h.cartSvc.List(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, toCartResponses(items), int64(len(items)), 1, len(items))
}
