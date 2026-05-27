package repository

import "gorm.io/gorm"

type PermissionRepo struct {
	db *gorm.DB
}

func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{db: db}
}

func (r *PermissionRepo) GetPermissionCodesByUserID(userID int64) ([]string, error) {
	var roleIDs []int64
	if err := r.db.Table("sys_user_role").Where("user_id = ?", userID).Pluck("role_id", &roleIDs).Error; err != nil {
		return nil, err
	}
	if len(roleIDs) == 0 {
		return nil, nil
	}

	var codes []string
	err := r.db.Table("sys_role_permission rp").
		Select("DISTINCT p.code").
		Joins("JOIN sys_permission p ON rp.permission_id = p.id").
		Where("rp.role_id IN ? AND p.status = 1", roleIDs).
		Pluck("p.code", &codes).Error
	return codes, err
}
