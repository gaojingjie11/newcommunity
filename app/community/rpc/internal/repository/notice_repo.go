package repository

import (
	"context"
	"fmt"

	"smartcommunity-microservices/app/community/rpc/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type NoticeRepo struct {
	db *gorm.DB
}

func NewNoticeRepo(db *gorm.DB) *NoticeRepo {
	return &NoticeRepo{db: db}
}

func (r *NoticeRepo) List(page, size int, includeHidden bool) ([]model.Notice, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
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

func (r *NoticeRepo) View(ctx context.Context, rdb *redis.Client, id int64) (*model.Notice, error) {
	var notice model.Notice
	if err := r.db.First(&notice, id).Error; err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf("notice:view_count:%d", id)

	v, err := rdb.Incr(ctx, redisKey).Result()
	if err == nil {
		if v == 1 {
			// Initialize Redis with MySQL's view_count + 1
			rdb.Set(ctx, redisKey, notice.ViewCount+1, 0)
			v = notice.ViewCount + 1
		}
		notice.ViewCount = v

		// Quantitative writeback: every 10 clicks, update DB
		if v%10 == 0 {
			r.db.Model(&model.Notice{}).Where("id = ?", id).UpdateColumn("view_count", v)
		}
	} else {
		// Fallback to direct DB increment if Redis is down
		r.db.Model(&model.Notice{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
		notice.ViewCount++
	}

	return &notice, nil
}
