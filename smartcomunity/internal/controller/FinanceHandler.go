package controller

import (
	"fmt"
	"strconv"
	"strings"

	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"

	"github.com/gin-gonic/gin"
)

type FinanceHandler struct {
	Service service.FinanceService
}

type financePayRequest struct {
	BusinessID   int64
	BusinessType int
	PayType      string
	Password     string
	FaceImageURL string
}

func (h *FinanceHandler) Pay(c *gin.Context) {
	userID, _ := c.Get("userID")
	uid := userID.(int64)

	req, err := parseFinancePayRequest(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	var confidence *float32
	if req.PayType == service.AuthTypeFace {
		var user model.SysUser
		if err := global.DB.Select("id", "face_registered", "face_image_url").First(&user, uid).Error; err != nil {
			response.Fail(c, "user not found")
			return
		}
		if !user.FaceRegistered || strings.TrimSpace(user.FaceImageURL) == "" {
			response.Fail(c, "face is not registered, please register first")
			return
		}
		if strings.TrimSpace(req.FaceImageURL) == "" {
			response.Fail(c, "face image is required")
			return
		}

		faceService, err := service.NewFaceService()
		if err != nil {
			response.Fail(c, "face service init failed: "+err.Error())
			return
		}

		score, err := faceService.CompareFace(user.FaceImageURL, req.FaceImageURL)
		if err != nil {
			response.Fail(c, "人脸验证失败，请稍后重试")
			return
		}
		if score < 85.0 {
			response.Fail(c, "人脸不匹配，请重试")
			return
		}
		confidence = &score
	}

	result, err := h.Service.UnifiedPayWithAuth(uid, req.BusinessID, req.BusinessType, req.Password, req.PayType)
	if err != nil {
		response.Fail(c, "payment failed: "+err.Error())
		return
	}

	if confidence != nil {
		response.Success(c, gin.H{
			"payment_result":  result,
			"face_confidence": *confidence,
		})
		return
	}
	response.Success(c, result)
}

func (h *FinanceHandler) ListPropertyFee(c *gin.Context) {
	userID, _ := c.Get("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.Service.GetPropertyFeeList(userID.(int64), page, size)
	if err != nil {
		response.Fail(c, "failed to fetch property fee list")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

func (h *FinanceHandler) Recharge(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}

	if err := h.Service.Recharge(userID.(int64), req.Amount); err != nil {
		response.Fail(c, "recharge failed: "+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *FinanceHandler) Transfer(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ToMobile string  `json:"to_mobile"`
		Amount   float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}

	if err := h.Service.Transfer(userID.(int64), req.ToMobile, req.Amount); err != nil {
		response.Fail(c, "transfer failed: "+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *FinanceHandler) ListTransactions(c *gin.Context) {
	userID, _ := c.Get("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	list, total, err := h.Service.GetTransactionList(userID.(int64), page, size)
	if err != nil {
		response.Fail(c, "failed to fetch transaction list")
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

func (h *FinanceHandler) CreatePropertyFee(c *gin.Context) {
	var req model.PropertyFee
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}

	if err := h.Service.CreatePropertyFee(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *FinanceHandler) ListAllPropertyFees(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	list, total, err := h.Service.ListAllPropertyFees(page, size)
	if err != nil {
		response.Fail(c, "failed to fetch property fee list")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

func parseFinancePayRequest(c *gin.Context) (*financePayRequest, error) {
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		return nil, fmt.Errorf("invalid request parameters")
	}

	req := &financePayRequest{
		BusinessID:   asInt64(raw["business_id"]),
		BusinessType: asInt(raw["business_type"]),
		PayType:      strings.ToLower(strings.TrimSpace(asString(raw["pay_type"]))),
		Password:     strings.TrimSpace(asString(raw["password"])),
		FaceImageURL: strings.TrimSpace(asString(raw["face_image_url"])),
	}

	// 兼容旧版：pay_type 传数字时表示业务类型(1=订单,2=物业)
	if legacyType := asInt(raw["pay_type"]); legacyType > 0 {
		if req.BusinessType == 0 {
			req.BusinessType = legacyType
		}
		if req.PayType == strconv.Itoa(legacyType) {
			req.PayType = service.AuthTypePassword
		}
	}
	if req.BusinessType == 0 {
		return nil, fmt.Errorf("business_type is required")
	}
	if req.PayType == "" {
		req.PayType = service.AuthTypePassword
	}

	if req.BusinessID <= 0 {
		return nil, fmt.Errorf("business_id is required")
	}
	if req.PayType != service.AuthTypePassword && req.PayType != service.AuthTypeFace {
		return nil, fmt.Errorf("pay_type must be password or face")
	}

	return req, nil
}

func asInt64(v interface{}) int64 {
	switch t := v.(type) {
	case float64:
		return int64(t)
	case float32:
		return int64(t)
	case int:
		return int64(t)
	case int64:
		return t
	case int32:
		return int64(t)
	case string:
		if t == "" {
			return 0
		}
		n, err := strconv.ParseInt(strings.TrimSpace(t), 10, 64)
		if err == nil {
			return n
		}
	}
	return 0
}

func asInt(v interface{}) int {
	return int(asInt64(v))
}

func asString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", t)
	}
}
