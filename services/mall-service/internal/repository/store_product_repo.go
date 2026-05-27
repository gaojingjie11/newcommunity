package repository

import (
	"smartcommunity-microservices/services/mall-service/internal/model"

	"gorm.io/gorm"
)

type StoreProductRepo struct {
	db *gorm.DB
}

func NewStoreProductRepo(db *gorm.DB) *StoreProductRepo {
	return &StoreProductRepo{db: db}
}

func (r *StoreProductRepo) WithTx(tx *gorm.DB) *StoreProductRepo {
	return &StoreProductRepo{db: tx}
}

func (r *StoreProductRepo) Bind(storeID, productID int64, stock int) error {
	sp := model.StoreProduct{StoreID: storeID, ProductID: productID, Stock: stock, Status: 1}
	return r.db.Create(&sp).Error
}

func (r *StoreProductRepo) Unbind(storeID, productID int64) error {
	return r.db.Where("store_id = ? AND product_id = ?", storeID, productID).Delete(&model.StoreProduct{}).Error
}

func (r *StoreProductRepo) UpdateStatus(storeID, productID int64, status int) error {
	return r.db.Model(&model.StoreProduct{}).Where("store_id = ? AND product_id = ?", storeID, productID).Update("status", status).Error
}

func (r *StoreProductRepo) UpdateStock(storeID, productID int64, stock int) error {
	return r.db.Model(&model.StoreProduct{}).Where("store_id = ? AND product_id = ?", storeID, productID).Update("stock", stock).Error
}

func (r *StoreProductRepo) ListByStore(storeID int64) ([]model.StoreProduct, error) {
	var items []model.StoreProduct
	err := r.db.Preload("Product").Where("store_id = ?", storeID).Find(&items).Error
	return items, err
}

func (r *StoreProductRepo) Find(storeID, productID int64) (*model.StoreProduct, error) {
	var sp model.StoreProduct
	err := r.db.Where("store_id = ? AND product_id = ?", storeID, productID).First(&sp).Error
	return &sp, err
}

func (r *StoreProductRepo) ListAvailableStoreIDs(productID int64, qty int) ([]int64, error) {
	var storeIDs []int64
	err := r.db.Model(&model.StoreProduct{}).
		Where("product_id = ? AND status = ? AND stock >= ?", productID, 1, qty).
		Pluck("store_id", &storeIDs).Error
	return storeIDs, err
}

func (r *StoreProductRepo) ListAvailableDetails(storeIDs, productIDs []int64) ([]model.StoreProduct, error) {
	var items []model.StoreProduct
	if len(storeIDs) == 0 || len(productIDs) == 0 {
		return items, nil
	}
	err := r.db.
		Where("store_id IN ? AND product_id IN ? AND status = ?", storeIDs, productIDs, 1).
		Preload("Product").
		Find(&items).Error
	return items, err
}

// DeductStock atomically decrements store product stock.
// Returns RowsAffected; 0 means insufficient stock.
func (r *StoreProductRepo) DeductStock(tx *gorm.DB, storeID, productID int64, qty int) (int64, error) {
	result := tx.Model(&model.StoreProduct{}).
		Where("store_id = ? AND product_id = ? AND status = ? AND stock >= ?", storeID, productID, 1, qty).
		UpdateColumns(map[string]interface{}{
			"stock":      gorm.Expr("stock - ?", qty),
			"sold_count": gorm.Expr("sold_count + ?", qty),
			"version":    gorm.Expr("version + 1"),
		})
	return result.RowsAffected, result.Error
}

// RestoreStock atomically increments store product stock.
func (r *StoreProductRepo) RestoreStock(tx *gorm.DB, storeID, productID int64, qty int) error {
	return tx.Model(&model.StoreProduct{}).
		Where("store_id = ? AND product_id = ?", storeID, productID).
		UpdateColumns(map[string]interface{}{
			"stock":      gorm.Expr("stock + ?", qty),
			"sold_count": gorm.Expr("sold_count - ?", qty),
			"version":    gorm.Expr("version + 1"),
		}).Error
}
