package repository

import (
	"smartcommunity-microservices/app/mall/rpc/internal/model"

	"gorm.io/gorm"
)

type PromotionRepo struct {
	db *gorm.DB
}

func NewPromotionRepo(db *gorm.DB) *PromotionRepo {
	return &PromotionRepo{db: db}
}

func (r *PromotionRepo) List(page, size int) ([]model.Promotion, int64, error) {
	var promotions []model.Promotion
	var total int64

	query := r.db.Model(&model.Promotion{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id desc").Preload("Products").Find(&promotions).Error; err != nil {
		return nil, 0, err
	}
	return promotions, total, nil
}

func (r *PromotionRepo) FindByID(id int64) (*model.Promotion, error) {
	var promo model.Promotion
	err := r.db.Preload("Products").First(&promo, id).Error
	return &promo, err
}

func (r *PromotionRepo) Create(promo *model.Promotion) error {
	return r.db.Create(promo).Error
}

func (r *PromotionRepo) Update(promo *model.Promotion) error {
	return r.db.Save(promo).Error
}

func (r *PromotionRepo) Delete(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("promotion_id = ?", id).Delete(&model.PromotionProduct{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Promotion{}, id).Error
	})
}

func (r *PromotionRepo) BindProducts(promotionID int64, productIDs []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("promotion_id = ?", promotionID).Delete(&model.PromotionProduct{}).Error; err != nil {
			return err
		}
		for _, pid := range productIDs {
			pp := model.PromotionProduct{PromotionID: promotionID, ProductID: pid}
			if err := tx.Create(&pp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
