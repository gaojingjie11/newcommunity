package handler

import (
	"errors"
	"net/http"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/community-service/internal/repository"
	"smartcommunity-microservices/services/community-service/internal/service"

	"github.com/gin-gonic/gin"
)

type ParkingHandler struct {
	svc *service.ParkingService
}

func NewParkingHandler(svc *service.ParkingService) *ParkingHandler {
	return &ParkingHandler{svc: svc}
}

func (h *ParkingHandler) MyBindings(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	items, err := h.svc.MyBindings(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询车位失败", nil)
		return
	}
	response.Success(c, gin.H{"list": items})
}

func (h *ParkingHandler) BindPlate(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req service.BindPlateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}
	item, err := h.svc.BindPlate(id, userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "绑定车牌失败: "+err.Error(), nil)
		return
	}
	response.Success(c, item)
}

func (h *ParkingHandler) AdminList(c *gin.Context) {
	page, size := response.ParsePage(c)
	items, total, err := h.svc.List(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询车位失败", nil)
		return
	}
	response.SuccessPaged(c, items, total, page, size)
}

func (h *ParkingHandler) Create(c *gin.Context) {
	var req service.CreateParkingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}
	item, err := h.svc.Create(req)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateParkingNo) {
			response.Error(c, http.StatusConflict, 409, "车位编号已存在，请换一个", nil)
			return
		}
		if err.Error() == "parking_no required" {
			response.Error(c, http.StatusBadRequest, 400, "请填写车位编号", nil)
			return
		}
		response.Error(c, http.StatusBadRequest, 400, "创建车位失败: "+err.Error(), nil)
		return
	}
	response.Success(c, item)
}

func (h *ParkingHandler) Assign(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req service.AssignParkingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}
	item, err := h.svc.Assign(id, req)
	if err != nil {
		if errors.Is(err, repository.ErrParkingSpaceUnavailable) {
			response.Error(c, http.StatusConflict, 409, "车位已被占用", nil)
			return
		}
		if errors.Is(err, repository.ErrParkingUserNotFound) {
			response.Error(c, http.StatusBadRequest, 400, "用户手机号不存在，请检查后重试", nil)
			return
		}
		if err.Error() == "mobile required" {
			response.Error(c, http.StatusBadRequest, 400, "请填写用户手机号", nil)
			return
		}
		response.Error(c, http.StatusBadRequest, 400, "分配车位失败: "+err.Error(), nil)
		return
	}
	response.Success(c, item)
}

func (h *ParkingHandler) Stats(c *gin.Context) {
	stats, err := h.svc.Stats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询车位统计失败", nil)
		return
	}
	response.Success(c, stats)
}
