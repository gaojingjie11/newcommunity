package handler

import (
	"net/http"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/community-service/internal/service"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	svc *service.MessageService
}

func NewMessageHandler(svc *service.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

func (h *MessageHandler) List(c *gin.Context) {
	page, size := response.ParsePage(c)
	items, total, err := h.svc.ListMessages(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询消息失败", nil)
		return
	}
	response.SuccessPaged(c, items, total, page, size)
}

func (h *MessageHandler) Send(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "消息内容不能为空", nil)
		return
	}
	msg, err := h.svc.SendMessage(userID, req.Content)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "发送消息失败", nil)
		return
	}
	response.Success(c, msg)
}
