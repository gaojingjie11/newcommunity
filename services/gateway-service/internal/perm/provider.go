package perm

import (
	"smartcommunity-microservices/services/gateway-service/internal/model"

	"gorm.io/gorm"
)

// PermissionProvider implements middleware.PermissionProvider by querying
// sys_user_role + sys_role_permission + sys_permission directly from MySQL.
type PermissionProvider struct {
	db *gorm.DB
}

func NewPermissionProvider(db *gorm.DB) *PermissionProvider {
	return &PermissionProvider{db: db}
}

// GetPermissionCodesByUserID returns all permission codes for a user.
func (p *PermissionProvider) GetPermissionCodesByUserID(userID int64) ([]string, error) {
	var codes []string
	err := p.db.Model(&model.SysUserRole{}).
		Select("DISTINCT sys_permission.code").
		Joins("JOIN sys_role_permission ON sys_role_permission.role_id = sys_user_role.role_id").
		Joins("JOIN sys_permission ON sys_permission.id = sys_role_permission.permission_id").
		Where("sys_user_role.user_id = ? AND sys_permission.status = 1", userID).
		Pluck("sys_permission.code", &codes).Error
	return codes, err
}
