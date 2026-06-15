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

func (r *StoreRepo) List(page, size int, areaID int64, storeIDs []int64) ([]model.Store, int64, error) {
	var stores []model.Store
	var total int64

	query := r.db.Model(&model.Store{})
	if areaID > 0 {
		query = query.Where("area_id = ?", areaID)
	}
	if len(storeIDs) > 0 {
		query = query.Where("id IN ?", storeIDs)
	} else if storeIDs != nil {
		// If storeIDs is non-nil but empty, it means the store admin is not bound to any store.
		// Return empty list immediately.
		return []model.Store{}, 0, nil
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id desc").Find(&stores).Error; err != nil {
		return nil, 0, err
	}
	return stores, total, nil
}

func (r *StoreRepo) FindByID(id int64) (*model.Store, error) {
	var store model.Store
	err := r.db.First(&store, id).Error
	return &store, err
}

func (r *StoreRepo) ListByIDs(ids []int64) ([]model.Store, error) {
	var stores []model.Store
	if len(ids) == 0 {
		return stores, nil
	}
	err := r.db.Where("id IN ?", ids).Order("id desc").Find(&stores).Error
	return stores, err
}

func (r *StoreRepo) Create(store *model.Store) error {
	return r.db.Create(store).Error
}

func (r *StoreRepo) Update(store *model.Store) error {
	return r.db.Save(store).Error
}

func (r *StoreRepo) Delete(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Cascade delete user store relationship bindings
		if err := tx.Where("store_id = ?", id).Delete(&model.UserStore{}).Error; err != nil {
			return err
		}
		// Delete the store itself
		return tx.Delete(&model.Store{}, id).Error
	})
}

func (r *StoreRepo) BindUserStores(userID int64, storeIDs []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Delete existing bindings for this user
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserStore{}).Error; err != nil {
			return err
		}
		// 2. Insert new bindings
		for _, storeID := range storeIDs {
			us := model.UserStore{
				UserID:  userID,
				StoreID: storeID,
			}
			if err := tx.Create(&us).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *StoreRepo) GetUserStores(userID int64) ([]int64, error) {
	var storeIDs []int64
	err := r.db.Model(&model.UserStore{}).Where("user_id = ?", userID).Pluck("store_id", &storeIDs).Error
	return storeIDs, err
}
