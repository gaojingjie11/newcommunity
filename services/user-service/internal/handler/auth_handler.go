package handler

import (
	"net/http"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/user-service/internal/model"
	"smartcommunity-microservices/services/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// POST /api/users/register (AUTH-001)
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}

	response.Success(c, gin.H{"uid": user.ID})
}

// POST /api/users/login (AUTH-002)
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Mobile   string `json:"mobile" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	token, user, err := h.authService.Login(req.Mobile, req.Password, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		response.Error(c, http.StatusUnauthorized, 401, err.Error(), nil)
		return
	}

	response.Success(c, gin.H{
		"token":     token,
		"user_info": userInfoPayload(user, true),
	})
}

// POST /api/users/sms-code (AUTH-008a)
func (h *AuthHandler) SendLoginCode(c *gin.Context) {
	var req struct {
		Mobile string `json:"mobile" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	_, err := h.authService.SendLoginCode(req.Mobile)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "验证码发送失败", nil)
		return
	}

	response.Success(c, gin.H{
		"expires_in": 300,
	})
}

// POST /api/users/login-code (AUTH-008b)
func (h *AuthHandler) LoginByCode(c *gin.Context) {
	var req struct {
		Mobile string `json:"mobile" binding:"required"`
		Code   string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	result, err := h.authService.LoginByCode(req.Mobile, req.Code, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		response.Error(c, http.StatusUnauthorized, 401, err.Error(), nil)
		return
	}

	response.Success(c, gin.H{
		"token":             result.Token,
		"user_info":         userInfoPayload(result.User, result.ProfileCompleted),
		"is_new_user":       result.IsNewUser,
		"profile_completed": result.ProfileCompleted,
	})
}

// POST /api/users/logout (AUTH-007)
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetInt64("userID")
	if err := h.authService.Logout(userID); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "退出失败", nil)
		return
	}
	response.Success(c, nil)
}

// POST /api/users/password-reset/code (AUTH-003a)
func (h *AuthHandler) SendResetCode(c *gin.Context) {
	var req struct {
		Mobile string `json:"mobile" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.authService.SendResetCode(req.Mobile); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// POST /api/users/password-reset (AUTH-003b)
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Mobile      string `json:"mobile" binding:"required"`
		Code        string `json:"code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.authService.ResetPassword(req.Mobile, req.Code, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// PUT /api/users/me/password (AUTH-004)
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetInt64("userID")
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

func userInfoPayload(user *model.SysUser, profileCompleted bool) gin.H {
	return gin.H{
		"id":                user.ID,
		"username":          user.Username,
		"real_name":         user.RealName,
		"mobile":            user.Mobile,
		"avatar":            user.Avatar,
		"green_points":      user.GreenPoints,
		"role":              user.Role,
		"status":            user.Status,
		"face_registered":   user.FaceRegistered,
		"face_image_url":    user.FaceImageURL,
		"profile_completed": profileCompleted,
	}
}
