package repository

import (
	"smartcommunity-microservices/services/mall-service/internal/model"

	"gorm.io/gorm"
)

type ServiceAreaRepo struct {
	db *gorm.DB
}

func NewServiceAreaRepo(db *gorm.DB) *ServiceAreaRepo {
	return &ServiceAreaRepo{db: db}
}

func (r *ServiceAreaRepo) List() ([]model.ServiceArea, error) {
	var areas []model.ServiceArea
	err := r.db.Where("status = ?", 1).Order("sort asc, id asc").Find(&areas).Error
	return areas, err
}

func (r *ServiceAreaRepo) ListAll() ([]model.ServiceArea, error) {
	var areas []model.ServiceArea
	err := r.db.Order("sort asc, id asc").Find(&areas).Error
	return areas, err
}

func (r *ServiceAreaRepo) FindByID(id int64) (*model.ServiceArea, error) {
	var area model.ServiceArea
	err := r.db.First(&area, id).Error
	return &area, err
}

func (r *ServiceAreaRepo) Create(area *model.ServiceArea) error {
	return r.db.Create(area).Error
}

func (r *ServiceAreaRepo) Update(area *model.ServiceArea) error {
	return r.db.Save(area).Error
}

func (r *ServiceAreaRepo) Delete(id int64) error {
	return r.db.Delete(&model.ServiceArea{}, id).Error
}
