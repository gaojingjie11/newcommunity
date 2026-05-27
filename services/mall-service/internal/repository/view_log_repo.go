package repository

import (
	"smartcommunity-microservices/services/mall-service/internal/model"

	"gorm.io/gorm"
)

type ViewLogRepo struct {
	db *gorm.DB
}

func NewViewLogRepo(db *gorm.DB) *ViewLogRepo {
	return &ViewLogRepo{db: db}
}

func (r *ViewLogRepo) Create(log *model.ProductViewLog) error {
	return r.db.Create(log).Error
}
