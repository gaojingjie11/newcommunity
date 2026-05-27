package repository

import (
	"time"

	"smartcommunity-microservices/services/user-service/internal/model"

	"gorm.io/gorm"
)

type PasswordResetRepo struct {
	db *gorm.DB
}

func NewPasswordResetRepo(db *gorm.DB) *PasswordResetRepo {
	return &PasswordResetRepo{db: db}
}

func (r *PasswordResetRepo) Create(code *model.PasswordResetCode) error {
	return r.db.Create(code).Error
}

func (r *PasswordResetRepo) MarkUsed(id int64) error {
	now := time.Now()
	return r.db.Model(&model.PasswordResetCode{}).Where("id = ?", id).Update("used_at", &now).Error
}

func (r *PasswordResetRepo) MarkUsedByMobile(mobile string) error {
	now := time.Now()
	return r.db.Model(&model.PasswordResetCode{}).
		Where("mobile = ? AND used_at IS NULL", mobile).
		Order("id desc").
		Limit(1).
		Update("used_at", &now).Error
}
