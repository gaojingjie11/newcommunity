package repository

import (
	"smartcommunity-microservices/app/user/rpc/internal/model"

	"gorm.io/gorm"
)

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) Create(role *model.SysRole) error {
	return r.db.Create(role).Error
}

func (r *RoleRepo) Update(role *model.SysRole) error {
	return r.db.Model(&model.SysRole{}).Where("id = ?", role.ID).Updates(role).Error
}

func (r *RoleRepo) Delete(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&model.SysRoleMenu{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&model.SysRolePermission{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&model.SysUserRole{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.SysRole{}, id).Error
	})
}

func (r *RoleRepo) FindAll() ([]model.SysRole, error) {
	var roles []model.SysRole
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) FindByID(id int64) (*model.SysRole, error) {
	var role model.SysRole
	err := r.db.Where("id = ?", id).First(&role).Error
	return &role, err
}

func (r *RoleRepo) FindByCode(code string) (*model.SysRole, error) {
	var role model.SysRole
	err := r.db.Where("code = ?", code).First(&role).Error
	return &role, err
}

// --- Menu binding ---

func (r *RoleRepo) BindMenus(roleID int64, menuIDs []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.SysRoleMenu{}).Error; err != nil {
			return err
		}
		for _, menuID := range menuIDs {
			rm := model.SysRoleMenu{RoleID: roleID, MenuID: menuID}
			if err := tx.Create(&rm).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RoleRepo) FindMenusByRoleID(roleID int64) ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := r.db.Joins("JOIN sys_role_menu ON sys_role_menu.menu_id = sys_menu.id").
		Where("sys_role_menu.role_id = ?", roleID).
		Order("sys_menu.sort asc").
		Find(&menus).Error
	return menus, err
}

func (r *RoleRepo) ListAllMenus() ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := r.db.Order("sort asc").Find(&menus).Error
	return menus, err
}

// --- Permission binding ---

func (r *RoleRepo) BindPermissions(roleID int64, permissionIDs []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.SysRolePermission{}).Error; err != nil {
			return err
		}
		for _, pid := range permissionIDs {
			rp := model.SysRolePermission{RoleID: roleID, PermissionID: pid}
			if err := tx.Create(&rp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RoleRepo) FindPermissionsByRoleID(roleID int64) ([]model.SysPermission, error) {
	var perms []model.SysPermission
	err := r.db.Joins("JOIN sys_role_permission ON sys_role_permission.permission_id = sys_permission.id").
		Where("sys_role_permission.role_id = ? AND sys_permission.status = 1", roleID).
		Find(&perms).Error
	return perms, err
}

func (r *RoleRepo) FindPermissionsByRoleIDs(roleIDs []int64) ([]model.SysPermission, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	var perms []model.SysPermission
	err := r.db.Distinct().
		Joins("JOIN sys_role_permission ON sys_role_permission.permission_id = sys_permission.id").
		Where("sys_role_permission.role_id IN ? AND sys_permission.status = 1", roleIDs).
		Find(&perms).Error
	return perms, err
}

func (r *RoleRepo) ListAllPermissions() ([]model.SysPermission, error) {
	var perms []model.SysPermission
	err := r.db.Where("status = 1").Order("id asc").Find(&perms).Error
	return perms, err
}

// --- User-role binding ---

func (r *RoleRepo) FindRolesByUserID(userID int64) ([]model.SysRole, error) {
	var roles []model.SysRole
	err := r.db.Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.id").
		Where("sys_user_role.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) BindUserRoles(userID int64, roleIDs []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.SysUserRole{}).Error; err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			ur := model.SysUserRole{UserID: userID, RoleID: roleID}
			if err := tx.Create(&ur).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RoleRepo) FindRoleCodesByUserID(userID int64) ([]string, error) {
	var codes []string
	err := r.db.Model(&model.SysRole{}).
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.id").
		Where("sys_user_role.user_id = ?", userID).
		Pluck("sys_role.code", &codes).Error
	return codes, err
}
