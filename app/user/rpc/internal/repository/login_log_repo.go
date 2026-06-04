package repository

import (
	"smartcommunity-microservices/app/user/rpc/internal/model"

	"gorm.io/gorm"
)

type LoginLogRepo struct {
	db *gorm.DB
}

func NewLoginLogRepo(db *gorm.DB) *LoginLogRepo {
	return &LoginLogRepo{db: db}
}

func (r *LoginLogRepo) CreateUserLog(log *model.UserLoginLog) error {
	return r.db.Create(log).Error
}

func (r *LoginLogRepo) CreateAdminLog(log *model.AdminLoginLog) error {
	return r.db.Create(log).Error
}

func (r *LoginLogRepo) QueryUserLogs(page, size int, mobile string, success *bool) ([]model.UserLoginLog, int64, error) {
	var logs []model.UserLoginLog
	var total int64

	query := r.db.Model(&model.UserLoginLog{})
	if mobile != "" {
		query = query.Where("mobile = ?", mobile)
	}
	if success != nil {
		query = query.Where("success = ?", *success)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id desc").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (r *LoginLogRepo) QueryAdminLogs(page, size int, mobile string, success *bool) ([]model.AdminLoginLog, int64, error) {
	var logs []model.AdminLoginLog
	var total int64

	query := r.db.Model(&model.AdminLoginLog{})
	if mobile != "" {
		query = query.Where("mobile = ?", mobile)
	}
	if success != nil {
		query = query.Where("success = ?", *success)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id desc").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
