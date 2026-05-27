package handler

import (
	"errors"
	"net/http"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/community-service/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NoticeHandler struct {
	svc *service.NoticeService
}

func NewNoticeHandler(svc *service.NoticeService) *NoticeHandler {
	return &NoticeHandler{svc: svc}
}

func (h *NoticeHandler) List(c *gin.Context) {
	page, size := response.ParsePage(c)
	items, total, err := h.svc.List(page, size, false)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询公告失败", nil)
		return
	}
	response.SuccessPaged(c, items, total, page, size)
}

func (h *NoticeHandler) AdminList(c *gin.Context) {
	page, size := response.ParsePage(c)
	items, total, err := h.svc.List(page, size, true)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询公告失败", nil)
		return
	}
	response.SuccessPaged(c, items, total, page, size)
}

func (h *NoticeHandler) Detail(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var viewerID int64
	if v, exists := c.Get("userID"); exists {
		viewerID, _ = v.(int64)
	}
	item, err := h.svc.Get(id, viewerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, 404, "公告不存在", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, 500, "查询公告失败", nil)
		return
	}
	response.Success(c, item)
}

func (h *NoticeHandler) Create(c *gin.Context) {
	var req service.CreateNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}
	item, err := h.svc.Create(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "创建公告失败", nil)
		return
	}
	response.Success(c, item)
}

func (h *NoticeHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "删除公告失败", nil)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *NoticeHandler) MarkRead(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.svc.MarkRead(id, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "标记公告已读失败", nil)
		return
	}
	response.Success(c, gin.H{"read": true})
}

func (h *NoticeHandler) Views(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	page, size := response.ParsePage(c)
	items, total, err := h.svc.ListViews(id, page, size)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, 404, "公告不存在", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, 500, "查询浏览状态失败", nil)
		return
	}
	response.SuccessPaged(c, items, total, page, size)
}
