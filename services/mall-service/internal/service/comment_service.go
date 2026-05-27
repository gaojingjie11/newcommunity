package service

import (
	"errors"
	"strings"
	"time"

	"smartcommunity-microservices/services/mall-service/internal/model"
	"smartcommunity-microservices/services/mall-service/internal/repository"
)

type CommentService struct {
	commentRepo *repository.CommentRepo
	productRepo *repository.ProductRepo
}

func NewCommentService(commentRepo *repository.CommentRepo, productRepo *repository.ProductRepo) *CommentService {
	return &CommentService{commentRepo: commentRepo, productRepo: productRepo}
}

func (s *CommentService) Create(userID, productID int64, content string, rating int) error {
	if userID <= 0 {
		return errors.New("请先登录")
	}
	if productID <= 0 {
		return errors.New("商品参数错误")
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("评价内容不能为空")
	}
	if rating < 1 || rating > 5 {
		return errors.New("评分必须在1到5之间")
	}
	if _, err := s.productRepo.FindByID(productID); err != nil {
		return errors.New("商品不存在")
	}

	return s.commentRepo.Create(&model.ProductComment{
		UserID:    userID,
		ProductID: productID,
		Content:   content,
		Rating:    rating,
		CreatedAt: time.Now(),
	})
}

func (s *CommentService) List(productID int64, page, size int) ([]model.ProductComment, int64, error) {
	if productID <= 0 {
		return nil, 0, errors.New("商品参数错误")
	}
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	return s.commentRepo.ListByProduct(productID, page, size)
}
