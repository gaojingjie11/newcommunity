package controller

import (
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommunityMessageHandler struct {
	Service service.CommunityMessageService
}

func (h *CommunityMessageHandler) Send(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		response.FailWithCode(c, 401, "please login first")
		return
	}
	uid, ok := userID.(int64)
	if !ok {
		response.FailWithCode(c, 401, "invalid login status")
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}

	msg, err := h.Service.Send(uid, req.Content)
	if err != nil {
		response.Fail(c, "send message failed: "+err.Error())
		return
	}
	response.Success(c, msg)
}

func (h *CommunityMessageHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	list, total, err := h.Service.List(page, size)
	if err != nil {
		response.Fail(c, "failed to fetch community messages")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}
