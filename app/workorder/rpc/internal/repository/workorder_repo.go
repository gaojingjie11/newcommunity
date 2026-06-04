package repository

import (
	"time"

	"smartcommunity-microservices/app/workorder/rpc/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WorkorderRepo struct {
	db *gorm.DB
}

func NewWorkorderRepo(db *gorm.DB) *WorkorderRepo {
	return &WorkorderRepo{db: db}
}

func (r *WorkorderRepo) Create(item *model.WorkOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		return tx.Create(&model.WorkorderLog{
			TargetType: item.Type,
			TargetID:   item.ID,
			FromStatus: -1,
			ToStatus:   item.Status,
			OperatorID: item.UserID,
			Action:     "created",
			Remark:     item.Description,
		}).Error
	})
}

func (r *WorkorderRepo) ListByUser(userID int64, page, size int) ([]model.WorkOrder, int64, error) {
	var items []model.WorkOrder
	var total int64
	q := r.db.Model(&model.WorkOrder{}).Where("user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
	return items, total, err
}

func (r *WorkorderRepo) ListAll(workorderType string, status *int, page, size int) ([]model.WorkOrderAdminView, int64, error) {
	var items []model.WorkOrderAdminView
	var total int64
	q := r.db.Model(&model.WorkOrder{})
	if workorderType != "" {
		q = q.Where("workorders.type = ?", workorderType)
	}
	if status != nil {
		q = q.Where("workorders.status = ?", *status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Select(`
			workorders.*,
			COALESCE(NULLIF(sys_user.real_name, ''), NULLIF(sys_user.username, ''), CONCAT('用户#', workorders.user_id)) AS user_name,
			COALESCE(sys_user.mobile, '') AS user_mobile
		`).
		Joins("LEFT JOIN sys_user ON sys_user.id = workorders.user_id").
		Order("workorders.created_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Scan(&items).Error
	return items, total, err
}

func (r *WorkorderRepo) Process(id, operatorID int64, status int, result string) (*model.WorkOrder, error) {
	var item model.WorkOrder
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, id).Error; err != nil {
			return err
		}
		from := item.Status
		item.Status = status
		item.Result = result
		item.ProcessorID = operatorID
		now := nowPtr()
		item.ProcessedAt = now
		if err := tx.Save(&item).Error; err != nil {
			return err
		}
		return tx.Create(&model.WorkorderLog{
			TargetType: item.Type,
			TargetID:   item.ID,
			FromStatus: from,
			ToStatus:   status,
			OperatorID: operatorID,
			Action:     "processed",
			Remark:     result,
		}).Error
	})
	return &item, err
}

func (r *WorkorderRepo) ListLogs(targetID int64) ([]model.WorkorderLog, error) {
	var workorder model.WorkOrder
	if err := r.db.Select("id", "type", "created_at").First(&workorder, targetID).Error; err != nil {
		return nil, err
	}
	var items []model.WorkorderLog
	err := r.db.Where("target_type = ? AND target_id = ? AND created_at >= ?", workorder.Type, targetID, workorder.CreatedAt).
		Order("id ASC").Find(&items).Error
	return items, err
}

func nowPtr() *time.Time {
	now := time.Now()
	return &now
}
