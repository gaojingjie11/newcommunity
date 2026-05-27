package handler

import (
	"net/http"
	"strconv"
	"strings"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/mall-service/internal/model"
	"smartcommunity-microservices/services/mall-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminOrderHandler struct {
	orderSvc   *service.OrderService
	productSvc *service.ProductService
	storeSvc   *service.StoreService
	promoSvc   *service.PromotionService
}

func NewAdminOrderHandler(
	orderSvc *service.OrderService,
	productSvc *service.ProductService,
	storeSvc *service.StoreService,
	promoSvc *service.PromotionService,
) *AdminOrderHandler {
	return &AdminOrderHandler{
		orderSvc:   orderSvc,
		productSvc: productSvc,
		storeSvc:   storeSvc,
		promoSvc:   promoSvc,
	}
}

// GET /api/admin/mall/orders (ADMIN-MALL-011)
func (h *AdminOrderHandler) ListOrders(c *gin.Context) {
	page, size := response.ParsePage(c)
	keyword := c.Query("keyword")

	var status *int
	if s := c.Query("status"); s != "" {
		v, err := strconv.Atoi(s)
		if err == nil {
			status = &v
		}
	}

	orders, total, err := h.orderSvc.AdminListOrders(page, size, status, keyword)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, toOrderResponses(orders), total, page, size)
}

// POST /api/admin/mall/orders/:id/ship (ADMIN-MALL-011)
func (h *AdminOrderHandler) ShipOrder(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	if err := h.orderSvc.ShipOrder(orderID); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// POST /api/admin/mall/orders/:id/cancel (ADMIN-MALL-011)
func (h *AdminOrderHandler) CancelOrder(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.Reason == "" {
		req.Reason = "管理员取消"
	}

	if err := h.orderSvc.CancelOrder(orderID, req.Reason); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// GET /api/admin/mall/products (ADMIN-MALL-006)
func (h *AdminOrderHandler) ListProducts(c *gin.Context) {
	page, size := response.ParsePage(c)
	name := strings.TrimSpace(c.Query("name"))
	categoryID, _ := strconv.ParseInt(c.Query("category_id"), 10, 64)

	var isPromotion *bool
	if raw := c.Query("is_promotion"); raw != "" {
		value, err := strconv.ParseBool(raw)
		if err != nil {
			response.Error(c, http.StatusBadRequest, 400, "is_promotion参数无效", nil)
			return
		}
		isPromotion = &value
	}

	products, total, err := h.productSvc.AdminList(page, size, name, categoryID, isPromotion)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, toProductResponses(products), total, page, size)
}

// POST /api/admin/mall/products (ADMIN-MALL-006)
func (h *AdminOrderHandler) CreateProduct(c *gin.Context) {
	var req productPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	product := req.toModel(0)
	if err := h.productSvc.Create(&product); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{"id": product.ID})
}

// PUT /api/admin/mall/products/:id (ADMIN-MALL-006)
func (h *AdminOrderHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	var req productPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	product := req.toModel(id)
	if err := h.productSvc.Update(&product); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// DELETE /api/admin/mall/products/:id (ADMIN-MALL-006)
func (h *AdminOrderHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	if err := h.productSvc.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// POST /api/admin/mall/promotions (ADMIN-MALL-007)
func (h *AdminOrderHandler) CreatePromotion(c *gin.Context) {
	var promo model.Promotion
	if err := c.ShouldBindJSON(&promo); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.promoSvc.Create(&promo); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{"id": promo.ID})
}

// PUT /api/admin/mall/promotions/:id (ADMIN-MALL-007)
func (h *AdminOrderHandler) UpdatePromotion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	var promo model.Promotion
	if err := c.ShouldBindJSON(&promo); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	promo.ID = id

	if err := h.promoSvc.Update(&promo); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// DELETE /api/admin/mall/promotions/:id (ADMIN-MALL-007)
func (h *AdminOrderHandler) DeletePromotion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	if err := h.promoSvc.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// POST /api/admin/mall/promotions/:id/products (ADMIN-MALL-007)
func (h *AdminOrderHandler) BindPromotionProducts(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	var req struct {
		ProductIDs []int64 `json:"product_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.promoSvc.BindProducts(id, req.ProductIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// POST /api/admin/mall/stores (ADMIN-MALL-009)
func (h *AdminOrderHandler) CreateStore(c *gin.Context) {
	var store model.Store
	if err := c.ShouldBindJSON(&store); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.storeSvc.Create(&store); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{"id": store.ID})
}

// PUT /api/admin/mall/stores/:id (ADMIN-MALL-009)
func (h *AdminOrderHandler) UpdateStore(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	var store model.Store
	if err := c.ShouldBindJSON(&store); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	store.ID = id

	if err := h.storeSvc.Update(&store); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// DELETE /api/admin/mall/stores/:id (ADMIN-MALL-009)
func (h *AdminOrderHandler) DeleteStore(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	if err := h.storeSvc.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// POST /api/admin/mall/store-products (ADMIN-MALL-010)
func (h *AdminOrderHandler) BindStoreProduct(c *gin.Context) {
	var req struct {
		StoreID   int64 `json:"store_id" binding:"required"`
		ProductID int64 `json:"product_id" binding:"required"`
		Stock     int   `json:"stock"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.storeSvc.BindProduct(req.StoreID, req.ProductID, req.Stock); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// DELETE /api/admin/mall/store-products (ADMIN-MALL-010)
func (h *AdminOrderHandler) UnbindStoreProduct(c *gin.Context) {
	var req struct {
		StoreID   int64 `json:"store_id" binding:"required"`
		ProductID int64 `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.storeSvc.UnbindProduct(req.StoreID, req.ProductID); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// PUT /api/admin/mall/store-products/status (ADMIN-MALL-010)
func (h *AdminOrderHandler) UpdateStoreProductStatus(c *gin.Context) {
	var req struct {
		StoreID   int64 `json:"store_id" binding:"required"`
		ProductID int64 `json:"product_id" binding:"required"`
		Status    int   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.storeSvc.UpdateProductStatus(req.StoreID, req.ProductID, req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// PUT /api/admin/mall/store-products/stock (ADMIN-MALL-010)
func (h *AdminOrderHandler) UpdateStoreProductStock(c *gin.Context) {
	var req struct {
		StoreID   int64 `json:"store_id" binding:"required"`
		ProductID int64 `json:"product_id" binding:"required"`
		Stock     int   `json:"stock"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}

	if err := h.storeSvc.UpdateProductStock(req.StoreID, req.ProductID, req.Stock); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// GET /api/admin/mall/store-products/:store_id (ADMIN-MALL-010)
func (h *AdminOrderHandler) ListStoreProducts(c *gin.Context) {
	storeID, err := strconv.ParseInt(c.Param("store_id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	items, err := h.storeSvc.ListProducts(storeID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, toStoreProductResponses(items), int64(len(items)), 1, len(items))
}
