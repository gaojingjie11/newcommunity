package repository

import (
	"smartcommunity-microservices/app/mall/rpc/internal/model"

	"gorm.io/gorm"
)

type CommentRepo struct {
	db *gorm.DB
}

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(comment *model.ProductComment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepo) ListByProduct(productID int64, page, size int) ([]model.ProductComment, int64, error) {
	var comments []model.ProductComment
	var total int64

	query := r.db.Model(&model.ProductComment{}).Where("product_id = ?", productID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := query.
		Preload("User").
		Order("created_at desc, id desc").
		Offset(offset).
		Limit(size).
		Find(&comments).Error
	return comments, total, err
}
