package handler

import (
	"net/http"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/community-service/internal/service"

	"github.com/gin-gonic/gin"
)

type WorkorderHandler struct {
	svc *service.WorkorderService
}

func NewWorkorderHandler(svc *service.WorkorderService) *WorkorderHandler {
	return &WorkorderHandler{svc: svc}
}

func (h *WorkorderHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	var req service.CreateWorkorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}
	result, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "提交工单失败: "+err.Error(), nil)
		return
	}
	response.Success(c, result)
}

func (h *WorkorderHandler) MyList(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	page, size := response.ParsePage(c)
	items, total, err := h.svc.MyList(userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询工单失败", nil)
		return
	}
	response.SuccessPaged(c, items, total, page, size)
}

func (h *WorkorderHandler) AdminList(c *gin.Context) {
	status, ok := parseStatusQuery(c)
	if !ok {
		return
	}
	page, size := response.ParsePage(c)
	items, total, err := h.svc.AdminList(c.Query("type"), status, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询工单失败", nil)
		return
	}
	response.SuccessPaged(c, items, total, page, size)
}

func (h *WorkorderHandler) Process(c *gin.Context) {
	operatorID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req service.ProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}
	item, err := h.svc.Process(id, operatorID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "处理工单失败: "+err.Error(), nil)
		return
	}
	response.Success(c, item)
}

func (h *WorkorderHandler) Logs(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	items, err := h.svc.Logs(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "查询状态日志失败: "+err.Error(), nil)
		return
	}
	response.Success(c, gin.H{"list": items})
}
