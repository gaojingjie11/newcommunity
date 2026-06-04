package repository

import (
	"smartcommunity-microservices/app/mall/rpc/internal/model"

	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) WithTx(tx *gorm.DB) *ProductRepo {
	return &ProductRepo{db: tx}
}

func (r *ProductRepo) List(page, size int, categoryID int64, sort string, keyword string) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{}).Where("status = ?", 1)
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "id desc"
	switch sort {
	case "price_asc":
		order = "price asc"
	case "price_desc":
		order = "price desc"
	case "sales":
		order = "sales desc"
	}

	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order(order).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductRepo) AdminList(page, size int, name string, categoryID int64, isPromotion *bool) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{})
	if name != "" {
		like := "%" + name + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", like, like)
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if isPromotion != nil {
		if *isPromotion {
			query = query.Where("original_price > price AND price > 0")
		} else {
			query = query.Where("original_price <= price OR price <= 0")
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id asc").Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductRepo) Search(keyword string, page, size int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	like := "%" + keyword + "%"
	query := r.db.Model(&model.Product{}).Where("status = ? AND (name LIKE ? OR description LIKE ?)", 1, like, like)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id desc").Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductRepo) FindByID(id int64) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *ProductRepo) ListPromotionProducts(page, size int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{}).Where("status = ? AND original_price > price AND price > 0", 1)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id desc").Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductRepo) Create(product *model.Product) error {
	r.normalizeProduct(product)
	return r.db.Create(product).Error
}

func (r *ProductRepo) Update(product *model.Product) error {
	var existing model.Product
	if err := r.db.First(&existing, product.ID).Error; err != nil {
		return err
	}
	existing.Name = product.Name
	existing.Description = product.Description
	existing.Price = product.Price
	existing.OriginalPrice = product.OriginalPrice
	if existing.OriginalPrice <= 0 {
		existing.OriginalPrice = existing.Price
	}
	if existing.OriginalPrice > existing.Price && existing.Price > 0 {
		existing.IsPromotion = 1
	} else {
		existing.IsPromotion = 0
	}
	existing.Stock = product.Stock
	existing.ImageURL = product.ImageURL
	existing.Status = product.Status
	existing.CategoryID = product.CategoryID

	if product.CategoryID > 0 {
		var cat model.ProductCategory
		if err := r.db.First(&cat, product.CategoryID).Error; err == nil {
			existing.CategoryName = cat.Name
		} else {
			existing.CategoryName = ""
		}
	} else {
		existing.CategoryName = ""
	}

	return r.db.Save(&existing).Error
}

func (r *ProductRepo) normalizeProduct(product *model.Product) {
	if product.OriginalPrice <= 0 {
		product.OriginalPrice = product.Price
	}
	if product.OriginalPrice > product.Price && product.Price > 0 {
		product.IsPromotion = 1
	} else {
		product.IsPromotion = 0
	}
	if product.CategoryID > 0 {
		var cat model.ProductCategory
		if err := r.db.First(&cat, product.CategoryID).Error; err == nil {
			product.CategoryName = cat.Name
		}
	}
}

func (r *ProductRepo) UpdateFields(id int64, fields map[string]interface{}) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ProductRepo) Delete(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Delete bindings in pms_store_product
		if err := tx.Where("product_id = ?", id).Delete(&model.StoreProduct{}).Error; err != nil {
			return err
		}
		// 2. Delete cart items in oms_cart
		if err := tx.Where("product_id = ?", id).Delete(&model.Cart{}).Error; err != nil {
			return err
		}
		// 3. Delete favorites in pms_favorite
		if err := tx.Where("product_id = ?", id).Delete(&model.Favorite{}).Error; err != nil {
			return err
		}
		// 4. Delete the product itself
		return tx.Delete(&model.Product{}, id).Error
	})
}

// DeductStock atomically decrements stock and increments sales.
// Uses the caller's tx to ensure the operation is within the correct transaction.
// Returns RowsAffected; 0 means insufficient stock.
func (r *ProductRepo) DeductStock(tx *gorm.DB, id int64, qty int) (int64, error) {
	result := tx.Model(&model.Product{}).
		Where("id = ? AND stock >= ?", id, qty).
		UpdateColumns(map[string]interface{}{
			"stock":   gorm.Expr("stock - ?", qty),
			"sales":   gorm.Expr("sales + ?", qty),
			"version": gorm.Expr("version + 1"),
		})
	return result.RowsAffected, result.Error
}

// RestoreStock atomically increments stock and decrements sales.
// Used when cancelling an order.
func (r *ProductRepo) RestoreStock(tx *gorm.DB, id int64, qty int) (int64, error) {
	result := tx.Model(&model.Product{}).
		Where("id = ?", id).
		UpdateColumns(map[string]interface{}{
			"stock":   gorm.Expr("stock + ?", qty),
			"sales":   gorm.Expr("sales - ?", qty),
			"version": gorm.Expr("version + 1"),
		})
	return result.RowsAffected, result.Error
}
