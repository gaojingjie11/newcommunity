package repository

import (
	"time"

	"smartcommunity-microservices/services/community-service/internal/model"

	"gorm.io/gorm"
)

type VisitorRepo struct {
	db *gorm.DB
}

func NewVisitorRepo(db *gorm.DB) *VisitorRepo {
	return &VisitorRepo{db: db}
}

func (r *VisitorRepo) Create(item *model.Visitor) error {
	return r.db.Create(item).Error
}

func (r *VisitorRepo) ListByUser(userID int64, page, size int) ([]model.Visitor, int64, error) {
	var items []model.Visitor
	var total int64
	q := r.db.Model(&model.Visitor{}).Where("user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
	return items, total, err
}

func (r *VisitorRepo) ListAll(status *int, page, size int) ([]model.VisitorAdminView, int64, error) {
	var items []model.VisitorAdminView
	var total int64
	q := r.db.Model(&model.Visitor{})
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Select(`
			visitors.*,
			COALESCE(NULLIF(sys_user.real_name, ''), NULLIF(sys_user.username, ''), CONCAT('用户#', visitors.user_id)) AS user_name,
			COALESCE(sys_user.mobile, '') AS user_mobile
		`).
		Joins("LEFT JOIN sys_user ON sys_user.id = visitors.user_id").
		Order("visitors.created_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Scan(&items).Error
	return items, total, err
}

func (r *VisitorRepo) Audit(id int64, status int, remark string) (*model.Visitor, error) {
	now := time.Now()
	var item model.Visitor
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&item, id).Error; err != nil {
			return err
		}
		item.Status = status
		item.AuditRemark = remark
		item.AuditAt = &now
		return tx.Save(&item).Error
	})
	return &item, err
}
