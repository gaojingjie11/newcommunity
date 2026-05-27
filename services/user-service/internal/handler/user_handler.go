package handler

import (
	"net/http"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GET /api/users/me (AUTH-005, AUTH-006)
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt64("userID")
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, "用户不存在", nil)
		return
	}
	response.Success(c, gin.H{
		"id":              user.ID,
		"username":        user.Username,
		"real_name":       user.RealName,
		"mobile":          user.Mobile,
		"age":             user.Age,
		"gender":          user.Gender,
		"email":           user.Email,
		"avatar":          user.Avatar,
		"green_points":    user.GreenPoints,
		"role":            user.Role,
		"status":          user.Status,
		"balance":         user.Balance,
		"face_registered": user.FaceRegistered,
		"face_image_url":  user.FaceImageURL,
		"created_at":      user.CreatedAt,
	})
}

// POST /api/users/me/face
func (h *UserHandler) RegisterFace(c *gin.Context) {
	userID := c.GetInt64("userID")
	var req struct {
		FaceImageURL string `json:"face_image_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	if err := h.userService.RegisterFace(userID, req.FaceImageURL); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{
		"face_registered": true,
		"face_image_url":  req.FaceImageURL,
	})
}

// PUT /api/users/me (AUTH-005)
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt64("userID")
	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.userService.UpdateProfile(userID, req); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "更新失败", nil)
		return
	}
	response.Success(c, nil)
}
