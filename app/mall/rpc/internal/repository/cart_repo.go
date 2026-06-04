package repository

import (
	"smartcommunity-microservices/app/mall/rpc/internal/model"

	"gorm.io/gorm"
)

type CartRepo struct {
	db *gorm.DB
}

func NewCartRepo(db *gorm.DB) *CartRepo {
	return &CartRepo{db: db}
}

func (r *CartRepo) WithTx(tx *gorm.DB) *CartRepo {
	return &CartRepo{db: tx}
}

func (r *CartRepo) FindByUserProduct(userID, productID int64) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error
	return &cart, err
}

func (r *CartRepo) Create(cart *model.Cart) error {
	return r.db.Create(cart).Error
}

func (r *CartRepo) UpdateQuantity(id int64, qty int) error {
	return r.db.Model(&model.Cart{}).Where("id = ?", id).Update("quantity", qty).Error
}

func (r *CartRepo) ListByUser(userID int64) ([]model.Cart, error) {
	var items []model.Cart
	err := r.db.Where("user_id = ?", userID).Preload("Product").Order("id desc").Find(&items).Error
	return items, err
}

func (r *CartRepo) FindByIDs(ids []int64, userID int64) ([]model.Cart, error) {
	var items []model.Cart
	err := r.db.Where("id IN ? AND user_id = ?", ids, userID).Preload("Product").Find(&items).Error
	return items, err
}

func (r *CartRepo) Delete(id int64, userID int64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Cart{}).Error
}

func (r *CartRepo) DeleteByIDs(ids []int64, userID int64) error {
	return r.db.Where("id IN ? AND user_id = ?", ids, userID).Delete(&model.Cart{}).Error
}
