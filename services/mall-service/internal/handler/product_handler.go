package handler

import (
	"net/http"
	"strconv"
	"time"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/mall-service/internal/model"
	"smartcommunity-microservices/services/mall-service/internal/repository"
	"smartcommunity-microservices/services/mall-service/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productSvc  *service.ProductService
	viewLogRepo *repository.ViewLogRepo
}

func NewProductHandler(productSvc *service.ProductService, viewLogRepo *repository.ViewLogRepo) *ProductHandler {
	return &ProductHandler{productSvc: productSvc, viewLogRepo: viewLogRepo}
}

// GET /api/mall/products (MALL-001/004)
func (h *ProductHandler) List(c *gin.Context) {
	page, size := response.ParsePage(c)
	categoryID, _ := strconv.ParseInt(c.Query("category_id"), 10, 64)
	sort := c.DefaultQuery("sort", "")
	keyword := c.Query("name")
	if keyword == "" {
		keyword = c.Query("keyword")
	}

	products, total, err := h.productSvc.List(page, size, categoryID, sort, keyword)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, toProductResponses(products), total, page, size)
}

// GET /api/mall/products/search (MALL-002)
func (h *ProductHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	page, size := response.ParsePage(c)

	products, total, err := h.productSvc.Search(keyword, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, toProductResponses(products), total, page, size)
}

// GET /api/mall/products/promotions (MALL-003)
func (h *ProductHandler) GetPromotions(c *gin.Context) {
	page, size := response.ParsePage(c)

	products, total, err := h.productSvc.GetPromotions(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.SuccessPaged(c, toProductResponses(products), total, page, size)
}

// GET /api/mall/products/:id (MALL-005)
func (h *ProductHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误", nil)
		return
	}

	product, err := h.productSvc.GetDetail(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, "商品不存在", nil)
		return
	}

	userID, _ := c.Get("userID")
	uid, _ := userID.(int64)
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()
	viewedAt := time.Now()

	// Record view log asynchronously (best-effort, non-blocking). Copy request
	// data before starting the goroutine; gin.Context is not goroutine-safe.
	go h.recordViewLog(id, uid, ip, userAgent, viewedAt)

	response.Success(c, toProductResponse(*product))
}

func (h *ProductHandler) recordViewLog(productID, userID int64, ip, userAgent string, viewedAt time.Time) {
	if h.viewLogRepo == nil {
		return
	}
	_ = h.viewLogRepo.Create(&model.ProductViewLog{
		ProductID: productID,
		UserID:    userID,
		IP:        ip,
		UserAgent: userAgent,
		ViewedAt:  viewedAt,
	})
}
