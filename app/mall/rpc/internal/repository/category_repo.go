package repository

import (
	"smartcommunity-microservices/app/mall/rpc/internal/model"

	"gorm.io/gorm"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) List() ([]model.ProductCategory, error) {
	var categories []model.ProductCategory
	err := r.db.Order("sort asc, id asc").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepo) FindByID(id int64) (*model.ProductCategory, error) {
	var cat model.ProductCategory
	err := r.db.First(&cat, id).Error
	return &cat, err
}

func (r *CategoryRepo) Create(cat *model.ProductCategory) error {
	return r.db.Create(cat).Error
}

func (r *CategoryRepo) Update(cat *model.ProductCategory) error {
	return r.db.Save(cat).Error
}

func (r *CategoryRepo) Delete(id int64) error {
	return r.db.Delete(&model.ProductCategory{}, id).Error
}
