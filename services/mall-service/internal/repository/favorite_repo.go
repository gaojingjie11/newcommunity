package repository

import (
	"smartcommunity-microservices/services/mall-service/internal/model"

	"gorm.io/gorm"
)

type FavoriteRepo struct {
	db *gorm.DB
}

func NewFavoriteRepo(db *gorm.DB) *FavoriteRepo {
	return &FavoriteRepo{db: db}
}

func (r *FavoriteRepo) Add(userID, productID int64) error {
	fav := model.Favorite{UserID: userID, ProductID: productID}
	return r.db.Create(&fav).Error
}

func (r *FavoriteRepo) Remove(userID, productID int64) error {
	return r.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&model.Favorite{}).Error
}

func (r *FavoriteRepo) Exists(userID, productID int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Favorite{}).Where("user_id = ? AND product_id = ?", userID, productID).Count(&count).Error
	return count > 0, err
}

func (r *FavoriteRepo) ListByUser(userID int64, page, size int) ([]model.Favorite, int64, error) {
	var favorites []model.Favorite
	var total int64

	query := r.db.Model(&model.Favorite{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id desc").Preload("Product").Find(&favorites).Error; err != nil {
		return nil, 0, err
	}
	return favorites, total, nil
}
