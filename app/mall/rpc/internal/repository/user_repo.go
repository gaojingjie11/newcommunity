package repository

import (
	"smartcommunity-microservices/app/mall/rpc/internal/model"

	"gorm.io/gorm"
)

// UserRepo is a read-only repository for user-service owned tables.
type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByID(id int64) (*model.SysUser, error) {
	var user model.SysUser
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}
