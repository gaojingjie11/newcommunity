package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/community-service/internal/service"

	"github.com/gin-gonic/gin"
)

type VisitorHandler struct {
	svc *service.VisitorService
}

func NewVisitorHandler(svc *service.VisitorService) *VisitorHandler {
	return &VisitorHandler{svc: svc}
}

func (h *VisitorHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	var req service.CreateVisitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}
	item, err := h.svc.Create(userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "访客登记失败: "+err.Error(), nil)
		return
	}
	response.Success(c, item)
}

func (h *VisitorHandler) MyList(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	page, size := response.ParsePage(c)
	items, total, err := h.svc.ListByUser(userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询访客记录失败", nil)
		return
	}
	response.SuccessPaged(c, items, total, page, size)
}

func (h *VisitorHandler) AdminList(c *gin.Context) {
	page, size := response.ParsePage(c)
	var status *int
	if raw := c.Query("status"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			response.Error(c, http.StatusBadRequest, 400, "status参数无效", nil)
			return
		}
		status = &value
	}
	items, total, err := h.svc.ListAll(status, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询访客记录失败", nil)
		return
	}
	response.SuccessPaged(c, items, total, page, size)
}

func (h *VisitorHandler) Audit(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req service.AuditVisitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}
	item, err := h.svc.Audit(id, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "访客审核失败: "+err.Error(), nil)
		return
	}
	response.Success(c, item)
}
