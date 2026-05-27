package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/mall-service/internal/service"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentSvc *service.CommentService
}

func NewCommentHandler(commentSvc *service.CommentService) *CommentHandler {
	return &CommentHandler{commentSvc: commentSvc}
}

// GET /api/mall/comments
func (h *CommentHandler) List(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Query("product_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "商品参数错误", nil)
		return
	}
	page, size := response.ParsePage(c)

	comments, total, err := h.commentSvc.List(productID, page, size)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, comments, total, page, size)
}

// POST /api/mall/comments
func (h *CommentHandler) Create(c *gin.Context) {
	userID := c.GetInt64("userID")
	var req struct {
		ProductID int64  `json:"product_id" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Rating    int    `json:"rating" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.commentSvc.Create(userID, req.ProductID, req.Content, req.Rating); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}
