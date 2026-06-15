package service

import (
	"context"
	"fmt"
	"time"

	"smartcommunity-microservices/app/user/rpc/internal/model"
	"smartcommunity-microservices/app/user/rpc/internal/repository"

	goredis "github.com/redis/go-redis/v9"
)

type AdminService struct {
	userRepo *repository.UserRepo
	roleRepo *repository.RoleRepo
	rdb      *goredis.Client
}

func NewAdminService(userRepo *repository.UserRepo, roleRepo *repository.RoleRepo, rdb *goredis.Client) *AdminService {
	return &AdminService{userRepo: userRepo, roleRepo: roleRepo, rdb: rdb}
}

// ADMIN-MALL-001: Role CRUD

func (s *AdminService) CreateRole(role *model.SysRole) error {
	return s.roleRepo.Create(role)
}

func (s *AdminService) UpdateRole(role *model.SysRole) error {
	return s.roleRepo.Update(role)
}

func (s *AdminService) DeleteRole(id int64) error {
	return s.roleRepo.Delete(id)
}

func (s *AdminService) ListRoles() ([]model.SysRole, error) {
	return s.roleRepo.FindAll()
}

// ADMIN-MALL-002: Bind menus to role

func (s *AdminService) BindRoleMenu(roleID int64, menuIDs []int64) error {
	return s.roleRepo.BindMenus(roleID, menuIDs)
}

// RBAC: Bind permissions to role

func (s *AdminService) BindRolePermissions(roleID int64, permissionIDs []int64) error {
	if err := s.roleRepo.BindPermissions(roleID, permissionIDs); err != nil {
		return err
	}
	s.invalidateRolePermissionCache(roleID)
	return nil
}

func (s *AdminService) GetRolePermissions(roleID int64) ([]model.SysPermission, error) {
	return s.roleRepo.FindPermissionsByRoleID(roleID)
}

func (s *AdminService) ListAllPermissions() ([]model.SysPermission, error) {
	return s.roleRepo.ListAllPermissions()
}

func (s *AdminService) GetPermissionsByRoleIDs(roleIDs []int64) ([]model.SysPermission, error) {
	return s.roleRepo.FindPermissionsByRoleIDs(roleIDs)
}

// RBAC: User-role management

func (s *AdminService) AssignUserRoles(userID int64, roleIDs []int64) error {
	if err := s.roleRepo.BindUserRoles(userID, roleIDs); err != nil {
		return err
	}
	// Sync legacy users.role field
	if len(roleIDs) > 0 {
		var hasAdmin bool
		var firstRoleCode string
		for _, rid := range roleIDs {
			role, err := s.roleRepo.FindByID(rid)
			if err == nil {
				if firstRoleCode == "" {
					firstRoleCode = role.Code
				}
				if role.Code == "admin" {
					hasAdmin = true
				}
			}
		}
		if hasAdmin {
			_ = s.userRepo.UpdateRole(userID, "admin")
		} else if firstRoleCode != "" {
			_ = s.userRepo.UpdateRole(userID, firstRoleCode)
		}
	} else {
		// Default to user if no roles assigned
		_ = s.userRepo.UpdateRole(userID, "user")
	}
	s.invalidateUserPermissionCache(userID)
	return nil
}

func (s *AdminService) GetUserRoles(userID int64) ([]model.SysRole, error) {
	return s.roleRepo.FindRolesByUserID(userID)
}

// ADMIN-MALL-003: Admin user management

func (s *AdminService) ListAdminUsers(page, size int, keyword string) ([]model.SysUser, int64, error) {
	return s.userRepo.ListUsers(page, size, keyword, "")
}

func (s *AdminService) FreezeUser(id int64, status int) error {
	return s.userRepo.UpdateStatus(id, status)
}

// AssignRole: legacy compat — writes users.role AND sys_user_role
func (s *AdminService) AssignRole(userID int64, roleCode string) error {
	if err := s.userRepo.UpdateRole(userID, roleCode); err != nil {
		return err
	}
	// Find role by code and bind in sys_user_role
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return err
	}
	for _, r := range roles {
		if r.Code == roleCode {
			_ = s.roleRepo.BindUserRoles(userID, []int64{r.ID})
			break
		}
	}
	s.invalidateUserPermissionCache(userID)
	return nil
}

// ADMIN-MALL-004: Member list

func (s *AdminService) ListMembers(page, size int, keyword string) ([]model.SysUser, int64, error) {
	return s.userRepo.ListMembers(page, size, keyword)
}

// RBAC: Menu list

func (s *AdminService) ListAllMenus() ([]model.SysMenu, error) {
	return s.roleRepo.ListAllMenus()
}

// --- Cache invalidation ---

func (s *AdminService) invalidateUserPermissionCache(userID int64) {
	ctx := context.Background()
	key := fmt.Sprintf("rbac:permissions:%d", userID)
	_ = s.rdb.Del(ctx, key).Err()
}

func (s *AdminService) invalidateRolePermissionCache(roleID int64) {
	// Short TTL strategy: cache expires naturally in 10 min.
	// For immediate effect, we'd need to enumerate all users with this role,
	// which is expensive. The 10-min TTL is acceptable for admin operations.
	_ = roleID
	_ = time.Now()
}

// UpdateUserBalance adjusts a user's balance (positive=recharge, negative=deduct).
func (s *AdminService) UpdateUserBalance(userID int64, amount float64, operatorID int64) error {
	if amount == 0 {
		return nil
	}
	remark := "系统余额充值"
	if amount < 0 {
		remark = "系统余额扣减"
	}
	return s.userRepo.AdjustBalance(userID, amount, operatorID, remark)
}
