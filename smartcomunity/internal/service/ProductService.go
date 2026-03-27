package service

import (
	"errors"
	"fmt"
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"time"

	"gorm.io/gorm"
)

type ProductService struct{}

func (s *ProductService) GetList(page, size int, name string, minPrice, maxPrice float64, sort string, categoryID int64, isPromotion bool, status *int) ([]model.Product, int64, error) {
	var list []model.Product
	var total int64

	tx := global.DB.Model(&model.Product{})
	if name != "" {
		tx = tx.Where("name LIKE ?", "%"+name+"%")
	}
	if minPrice > 0 {
		tx = tx.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		tx = tx.Where("price <= ?", maxPrice)
	}
	if categoryID > 0 {
		tx = tx.Where("category_id = ?", categoryID)
	}
	if status != nil {
		tx = tx.Where("status = ?", *status)
	}
	if isPromotion {
		// Promotion is defined by price relation, not stale stored flag.
		tx = tx.Where("original_price > price")
	}

	switch sort {
	case "price_asc":
		tx = tx.Order("price asc")
	case "price_desc":
		tx = tx.Order("price desc")
	case "sales_desc":
		tx = tx.Order("sales desc")
	default:
		tx = tx.Order("id desc")
	}

	tx.Count(&total)

	offset := (page - 1) * size
	if err := tx.Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		list[i].IsPromotion = calcIsPromotion(list[i].OriginalPrice, list[i].Price)
	}
	return list, total, nil
}

func (s *ProductService) GetDetail(id int64) (*model.Product, error) {
	var product model.Product
	if err := global.DB.First(&product, id).Error; err != nil {
		return nil, err
	}
	product.IsPromotion = calcIsPromotion(product.OriginalPrice, product.Price)
	return &product, nil
}

func (s *ProductService) GetSalesRank() ([]model.Product, error) {
	var list []model.Product
	err := global.DB.Select("id, name, sales, price").
		Order("sales desc").
		Limit(10).
		Find(&list).Error
	return list, err
}

func (s *ProductService) Create(product *model.Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}

	if product.CategoryID > 0 {
		var cat model.ProductCategory
		if err := global.DB.First(&cat, product.CategoryID).Error; err == nil {
			product.CategoryName = cat.Name
		}
	}

	product.IsPromotion = calcIsPromotion(product.OriginalPrice, product.Price)

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return fmt.Errorf("create product failed: %w", err)
		}
		return s.syncPromotionRecord(tx, product.ID, product.IsPromotion)
	})
}

func (s *ProductService) Update(product *model.Product) error {
	if product.ID == 0 {
		return errors.New("product id is required")
	}
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}

	if product.CategoryID > 0 {
		var cat model.ProductCategory
		if err := global.DB.First(&cat, product.CategoryID).Error; err == nil {
			product.CategoryName = cat.Name
		}
	}

	product.IsPromotion = calcIsPromotion(product.OriginalPrice, product.Price)

	return global.DB.Transaction(func(tx *gorm.DB) error {
		// Use map to ensure zero values (e.g. status=0) are persisted.
		updates := map[string]interface{}{
			"name":           product.Name,
			"price":          product.Price,
			"original_price": product.OriginalPrice,
			"stock":          product.Stock,
			"category_id":    product.CategoryID,
			"category_name":  product.CategoryName,
			"description":    product.Description,
			"image_url":      product.ImageURL,
			"status":         product.Status,
			"is_promotion":   product.IsPromotion,
		}
		if err := tx.Model(&model.Product{}).Where("id = ?", product.ID).Updates(updates).Error; err != nil {
			return err
		}
		return s.syncPromotionRecord(tx, product.ID, product.IsPromotion)
	})
}

func (s *ProductService) Delete(id int64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		tx.Delete(&model.Promotion{}, "product_id = ?", id)
		tx.Delete(&model.Product{}, id)
		return nil
	})
}

func (s *ProductService) GetCategories() ([]model.ProductCategory, error) {
	var list []model.ProductCategory
	err := global.DB.Find(&list).Error
	return list, err
}

func calcIsPromotion(originalPrice, price float64) int {
	if originalPrice > price {
		return 1
	}
	return 0
}

func (s *ProductService) syncPromotionRecord(tx *gorm.DB, productID int64, isPromotion int) error {
	if isPromotion == 1 {
		var count int64
		if err := tx.Model(&model.Promotion{}).Where("product_id = ?", productID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			promotion := model.Promotion{
				Title:     "special promotion",
				Type:      2,
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(1, 0, 0),
				Status:    1,
				ProductID: productID,
			}
			if err := tx.Create(&promotion).Error; err != nil {
				return fmt.Errorf("create promotion failed: %w", err)
			}
		}
		return nil
	}
	return tx.Where("product_id = ?", productID).Delete(&model.Promotion{}).Error
}
