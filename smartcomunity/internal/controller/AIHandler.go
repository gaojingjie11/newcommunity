package controller

import (
	"strconv"
	"strings"

	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	Service service.AIService
}

func (h *AIHandler) Send(c *gin.Context) {
	var req struct {
		Content         string `json:"content"`
		Message         string `json:"message"`
		PaymentPassword string `json:"payment_password"`
		PayType         string `json:"pay_type"`
		FaceImageURL    string `json:"face_image_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		content = strings.TrimSpace(req.Message)
	}
	if content == "" {
		response.Fail(c, "content cannot be empty")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		response.FailWithCode(c, 401, "please login first")
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok || userID <= 0 {
		response.FailWithCode(c, 401, "invalid login status")
		return
	}

	reply, err := h.Service.ChatWithMemory(
		userID,
		content,
		strings.TrimSpace(req.PaymentPassword),
		strings.TrimSpace(req.PayType),
		strings.TrimSpace(req.FaceImageURL),
	)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, gin.H{"reply": reply})
}

func (h *AIHandler) History(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.FailWithCode(c, 401, "please login first")
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok || userID <= 0 {
		response.FailWithCode(c, 401, "invalid login status")
		return
	}

	limit := 50
	if raw := strings.TrimSpace(c.DefaultQuery("limit", "50")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			limit = parsed
		}
	}

	list, err := h.Service.GetChatHistory(userID, limit)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, gin.H{"list": list})
}
