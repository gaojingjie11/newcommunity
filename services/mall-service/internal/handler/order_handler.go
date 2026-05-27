package handler

import (
	"net/http"
	"strconv"
	"strings"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/mall-service/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderSvc *service.OrderService
}

func NewOrderHandler(orderSvc *service.OrderService) *OrderHandler {
	return &OrderHandler{orderSvc: orderSvc}
}

// GET /api/mall/orders/available-stores?cart_ids=1,2
func (h *OrderHandler) ListAvailableStores(c *gin.Context) {
	userID := c.GetInt64("userID")
	cartIDs, err := parseIDList(c.Query("cart_ids"))
	if err != nil || len(cartIDs) == 0 {
		response.Error(c, http.StatusBadRequest, 400, "请选择要结算的购物车商品", nil)
		return
	}

	stores, err := h.orderSvc.ListAvailableStores(userID, cartIDs)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, stores, int64(len(stores)), 1, len(stores))
}

func parseIDList(raw string) ([]int64, error) {
	parts := strings.Split(raw, ",")
	ids := make([]int64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil || id <= 0 {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// POST /api/mall/orders (MALL-012)
func (h *OrderHandler) Create(c *gin.Context) {
	userID := c.GetInt64("userID")
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	order, err := h.orderSvc.CreateOrder(userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{
		"order_id":           order.ID,
		"order_no":           order.OrderNo,
		"total_amount":       yuanFromCents(order.TotalAmount),
		"expire_at":          order.ExpireAt,
		"expires_in_seconds": secondsUntil(order.ExpireAt),
	})
}

// POST /api/mall/orders/:id/cancel (MALL-014)
func (h *OrderHandler) Cancel(c *gin.Context) {
	userID := c.GetInt64("userID")
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	// Verify ownership
	order, err := h.orderSvc.GetOrder(orderID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, "订单不存在", nil)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.Reason == "" {
		req.Reason = "用户取消"
	}

	if err := h.orderSvc.CancelOrder(order.ID, req.Reason); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// POST /api/mall/orders/:id/receive (MALL-017)
func (h *OrderHandler) Receive(c *gin.Context) {
	userID := c.GetInt64("userID")
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	order, err := h.orderSvc.GetOrder(orderID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, "订单不存在", nil)
		return
	}

	if err := h.orderSvc.ReceiveOrder(order.ID); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// GET /api/mall/orders/:id (MALL-014)
func (h *OrderHandler) GetDetail(c *gin.Context) {
	userID := c.GetInt64("userID")
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	order, err := h.orderSvc.GetOrder(orderID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, err.Error(), nil)
		return
	}
	response.Success(c, toOrderResponse(*order))
}

// GET /api/mall/orders (MALL-015~017)
func (h *OrderHandler) List(c *gin.Context) {
	userID := c.GetInt64("userID")
	page, size := response.ParsePage(c)

	var status *int
	if s := c.Query("status"); s != "" {
		v, err := strconv.Atoi(s)
		if err == nil {
			status = &v
		}
	}

	orders, total, err := h.orderSvc.ListOrders(userID, page, size, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, toOrderResponses(orders), total, page, size)
}
