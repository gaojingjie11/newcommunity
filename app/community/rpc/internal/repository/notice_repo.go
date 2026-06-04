package repository

import (
	"time"

	"smartcommunity-microservices/app/community/rpc/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NoticeRepo struct {
	db *gorm.DB
}

func NewNoticeRepo(db *gorm.DB) *NoticeRepo {
	return &NoticeRepo{db: db}
}

func (r *NoticeRepo) List(page, size int, includeHidden bool) ([]model.Notice, int64, error) {
	var items []model.Notice
	var total int64
	q := r.db.Model(&model.Notice{})
	if !includeHidden {
		q = q.Where("status = ?", 1)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
	return items, total, err
}

func (r *NoticeRepo) Get(id int64) (*model.Notice, error) {
	var item model.Notice
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *NoticeRepo) Create(item *model.Notice) error {
	return r.db.Create(item).Error
}

func (r *NoticeRepo) Delete(id int64) error {
	return r.db.Model(&model.Notice{}).Where("id = ?", id).Update("status", 0).Error
}

func (r *NoticeRepo) View(id, userID int64) (*model.Notice, error) {
	var notice model.Notice
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&notice, id).Error; err != nil {
			return err
		}
		if notice.Status != 1 {
			return gorm.ErrRecordNotFound
		}
		if err := tx.Model(&model.Notice{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
			return err
		}
		notice.ViewCount++
		if userID > 0 {
			now := time.Now()
			log := model.NoticeViewLog{NoticeID: id, UserID: userID, ViewedAt: now}
			return tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "notice_id"}, {Name: "user_id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{"viewed_at": now, "updated_at": now}),
			}).Create(&log).Error
		}
		return nil
	})
	return &notice, err
}

func (r *NoticeRepo) MarkRead(id, userID int64) error {
	now := time.Now()
	log := model.NoticeViewLog{NoticeID: id, UserID: userID, ViewedAt: now, ReadAt: &now}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "notice_id"}, {Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"read_at": now, "updated_at": now}),
	}).Create(&log).Error
}

func (r *NoticeRepo) ListViews(noticeID int64, page, size int) ([]model.NoticeViewLog, int64, error) {
	var items []model.NoticeViewLog
	var total int64
	q := r.db.Model(&model.NoticeViewLog{}).Where("notice_id = ?", noticeID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("updated_at DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error
	return items, total, err
}
