package controller

import (
	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Service service.ProductService
}

// Request struct for Create/Update
type ProductReq struct {
	model.Product
}

// List 商品列表接口
func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	name := c.Query("name")
	sort := c.Query("sort")
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	// Added category_id filter
	categoryID, _ := strconv.ParseInt(c.Query("category_id"), 10, 64)
	// Support true/1 for is_promotion query parameter.
	isPromotion := c.Query("is_promotion") == "true" || c.Query("is_promotion") == "1"
	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			status = &s
		}
	}

	list, total, err := h.Service.GetList(page, size, name, minPrice, maxPrice, sort, categoryID, isPromotion, status)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// Detail 商品详情接口
func (h *ProductHandler) Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	product, err := h.Service.GetDetail(id)
	if err != nil {
		response.Fail(c, "商品不存在")
		return
	}
	response.Success(c, product)
}

// Create 发布商品 (Admin)
func (h *ProductHandler) Create(c *gin.Context) {
	var req ProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	req.Status = 1
	req.CreatedAt = time.Now()

	// Handler now just passes request payload which includes OriginalPrice
	// Service handles logic
	if err := h.Service.Create(&req.Product); err != nil {
		response.Fail(c, "发布失败")
		return
	}
	response.Success(c, req.Product)
}

// Update 修改商品 (Admin)
func (h *ProductHandler) Update(c *gin.Context) {
	var req ProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := h.Service.Update(&req.Product); err != nil {
		response.Fail(c, "更新失败")
		return
	}
	response.Success(c, nil)
}

// Delete 删除商品 (Admin)
func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.Service.Delete(id); err != nil {
		response.Fail(c, "删除失败")
		return
	}
	response.Success(c, nil)
}

// GetRank 获取销量排行
func (h *ProductHandler) GetRank(c *gin.Context) {
	list, err := h.Service.GetSalesRank()
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, list)
}

// GetCategories 获取分类列表 (New)
func (h *ProductHandler) GetCategories(c *gin.Context) {
	list, err := h.Service.GetCategories()
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, list)
}
