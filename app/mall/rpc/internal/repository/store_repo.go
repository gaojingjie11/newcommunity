package repository

import (
	"smartcommunity-microservices/app/mall/rpc/internal/model"

	"gorm.io/gorm"
)

type StoreRepo struct {
	db *gorm.DB
}

func NewStoreRepo(db *gorm.DB) *StoreRepo {
	return &StoreRepo{db: db}
}

func (r *StoreRepo) List(page, size int, areaID int64) ([]model.Store, int64, error) {
	var stores []model.Store
	var total int64

	query := r.db.Model(&model.Store{})
	if areaID > 0 {
		query = query.Where("area_id = ?", areaID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id desc").Preload("ServiceArea").Find(&stores).Error; err != nil {
		return nil, 0, err
	}
	return stores, total, nil
}

func (r *StoreRepo) FindByID(id int64) (*model.Store, error) {
	var store model.Store
	err := r.db.Preload("ServiceArea").First(&store, id).Error
	return &store, err
}

func (r *StoreRepo) ListByIDs(ids []int64) ([]model.Store, error) {
	var stores []model.Store
	if len(ids) == 0 {
		return stores, nil
	}
	err := r.db.Where("id IN ?", ids).Order("id desc").Preload("ServiceArea").Find(&stores).Error
	return stores, err
}

func (r *StoreRepo) Create(store *model.Store) error {
	return r.db.Create(store).Error
}

func (r *StoreRepo) Update(store *model.Store) error {
	return r.db.Save(store).Error
}

func (r *StoreRepo) Delete(id int64) error {
	return r.db.Delete(&model.Store{}, id).Error
}
