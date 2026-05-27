package repository

import (
	"errors"
	"math"

	"smartcommunity-microservices/services/user-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByMobile(mobile string) (*model.SysUser, error) {
	var user model.SysUser
	err := r.db.Model(&model.SysUser{}).
		Joins("LEFT JOIN wallets ON wallets.user_id = sys_user.id").
		Select("sys_user.*, COALESCE(wallets.balance, 0) / 100 AS balance").
		Where("sys_user.mobile = ?", mobile).
		First(&user).Error
	return &user, err
}

func (r *UserRepo) FindByID(id int64) (*model.SysUser, error) {
	var user model.SysUser
	err := r.db.Model(&model.SysUser{}).
		Joins("LEFT JOIN wallets ON wallets.user_id = sys_user.id").
		Select("sys_user.*, COALESCE(wallets.balance, 0) / 100 AS balance").
		Where("sys_user.id = ?", id).
		First(&user).Error
	return &user, err
}

func (r *UserRepo) Create(user *model.SysUser) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) UpdatePassword(id int64, hash string) error {
	return r.db.Model(&model.SysUser{}).Where("id = ?", id).Update("password", hash).Error
}

func (r *UserRepo) UpdateFields(id int64, fields map[string]interface{}) error {
	return r.db.Model(&model.SysUser{}).Where("id = ?", id).Updates(fields).Error
}

func (r *UserRepo) ListUsers(page, size int, keyword string, excludeRole string) ([]model.SysUser, int64, error) {
	var users []model.SysUser
	var total int64

	query := r.db.Model(&model.SysUser{}).Joins("LEFT JOIN wallets ON wallets.user_id = sys_user.id")
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR mobile LIKE ? OR real_name LIKE ?", like, like, like)
	}
	if excludeRole != "" {
		query = query.Where("role != ?", excludeRole)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Select("sys_user.*, COALESCE(wallets.balance, 0) / 100 AS balance").
		Offset(offset).Limit(size).Order("sys_user.id asc").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepo) ListMembers(page, size int, keyword string) ([]model.SysUser, int64, error) {
	var users []model.SysUser
	var total int64

	query := r.db.Model(&model.SysUser{}).Joins("LEFT JOIN wallets ON wallets.user_id = sys_user.id").Where("role = ?", "user")
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR mobile LIKE ? OR real_name LIKE ?", like, like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Select("sys_user.*, COALESCE(wallets.balance, 0) / 100 AS balance").
		Offset(offset).Limit(size).Order("sys_user.id asc").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepo) UpdateStatus(id int64, status int) error {
	return r.db.Model(&model.SysUser{}).Where("id = ?", id).Update("status", status).Error
}

func (r *UserRepo) UpdateRole(id int64, role string) error {
	return r.db.Model(&model.SysUser{}).Where("id = ?", id).Update("role", role).Error
}

func (r *UserRepo) CountByMobile(mobile string) (int64, error) {
	var count int64
	err := r.db.Model(&model.SysUser{}).Where("mobile = ?", mobile).Count(&count).Error
	return count, err
}

// AdjustBalance atomically adjusts the mall wallet balance and writes wallet_transactions.
// sys_user.balance is legacy-only and is no longer the source of truth.
func (r *UserRepo) AdjustBalance(userID int64, amount float64, operatorID int64, remark string) error {
	cents := int64(math.Round(amount * 100))
	if cents == 0 {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		var userIDCheck int64
		if err := tx.Model(&model.SysUser{}).Select("id").Where("id = ?", userID).First(&userIDCheck).Error; err != nil {
			return err
		}
		if err := tx.Exec(
			"INSERT INTO wallets (user_id, balance, version, created_at, updated_at) VALUES (?, 0, 0, NOW(), NOW()) ON DUPLICATE KEY UPDATE user_id = user_id",
			userID,
		).Error; err != nil {
			return err
		}

		var wallet struct {
			ID      int64
			Balance int64
		}
		if err := tx.Table("wallets").Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id, balance").Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return err
		}

		after := wallet.Balance + cents
		if after < 0 {
			return errors.New("余额不足")
		}
		if err := tx.Table("wallets").Where("id = ?", wallet.ID).Updates(map[string]interface{}{
			"balance":    after,
			"version":    gorm.Expr("version + 1"),
			"updated_at": gorm.Expr("NOW()"),
		}).Error; err != nil {
			return err
		}

		txType := 3
		bizType := "admin_recharge"
		if cents < 0 {
			txType = 4
			bizType = "admin_deduct"
		}
		return tx.Table("wallet_transactions").Create(map[string]interface{}{
			"user_id":        userID,
			"type":           txType,
			"amount":         cents,
			"balance_before": wallet.Balance,
			"balance_after":  after,
			"related_id":     operatorID,
			"remark":         remark,
			"biz_type":       bizType,
			"biz_id":         "",
			"created_at":     gorm.Expr("NOW()"),
		}).Error
	})
}
