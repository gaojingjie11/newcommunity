package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/user-service/internal/model"
	"smartcommunity-microservices/services/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService    *service.AdminService
	loginLogService *service.LoginLogService
}

func NewAdminHandler(adminService *service.AdminService, loginLogService *service.LoginLogService) *AdminHandler {
	return &AdminHandler{adminService: adminService, loginLogService: loginLogService}
}

// --- PermissionProvider implementation for RequirePermission middleware ---

func (h *AdminHandler) GetPermissionCodesByUserID(userID int64) ([]string, error) {
	roles, err := h.adminService.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}
	roleIDs := make([]int64, len(roles))
	for i, r := range roles {
		roleIDs[i] = r.ID
	}
	perms, err := h.adminService.GetPermissionsByRoleIDs(roleIDs)
	if err != nil {
		return nil, err
	}
	codes := make([]string, len(perms))
	for i, p := range perms {
		codes[i] = p.Code
	}
	return codes, nil
}

// POST /api/admin/roles (ADMIN-MALL-001)
func (h *AdminHandler) CreateRole(c *gin.Context) {
	var role model.SysRole
	if err := c.ShouldBindJSON(&role); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	if err := h.adminService.CreateRole(&role); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "创建角色失败", nil)
		return
	}
	response.Success(c, gin.H{"id": role.ID})
}

// PUT /api/admin/roles (ADMIN-MALL-001)
func (h *AdminHandler) UpdateRole(c *gin.Context) {
	var role model.SysRole
	if err := c.ShouldBindJSON(&role); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	if err := h.adminService.UpdateRole(&role); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "更新角色失败", nil)
		return
	}
	response.Success(c, nil)
}

// DELETE /api/admin/roles (ADMIN-MALL-001)
func (h *AdminHandler) DeleteRole(c *gin.Context) {
	var req struct {
		ID int64 `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	if err := h.adminService.DeleteRole(req.ID); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "删除角色失败", nil)
		return
	}
	response.Success(c, nil)
}

// GET /api/admin/roles (ADMIN-MALL-001)
func (h *AdminHandler) ListRoles(c *gin.Context) {
	roles, err := h.adminService.ListRoles()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询角色失败", nil)
		return
	}
	response.Success(c, gin.H{"list": roles, "total": len(roles)})
}

// POST /api/admin/roles/:id/menus (ADMIN-MALL-002)
func (h *AdminHandler) BindRoleMenu(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "角色ID格式错误", nil)
		return
	}
	var req struct {
		MenuIDs []int64 `json:"menu_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	if err := h.adminService.BindRoleMenu(roleID, req.MenuIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "绑定菜单失败", nil)
		return
	}
	response.Success(c, nil)
}

// GET /api/admin/users (ADMIN-MALL-003)
func (h *AdminHandler) ListAdminUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")

	users, total, err := h.adminService.ListAdminUsers(page, size, keyword)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询失败", nil)
		return
	}
	response.Success(c, gin.H{"list": users, "total": total})
}

// POST /api/admin/users/freeze (ADMIN-MALL-003)
func (h *AdminHandler) FreezeUser(c *gin.Context) {
	var req struct {
		ID     int64 `json:"id" binding:"required"`
		Status int   `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	if err := h.adminService.FreezeUser(req.ID, req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "操作失败", nil)
		return
	}
	response.Success(c, nil)
}

// POST /api/admin/users/assign-role (ADMIN-MALL-003)
func (h *AdminHandler) AssignRole(c *gin.Context) {
	var req struct {
		UserID   int64  `json:"user_id" binding:"required"`
		RoleCode string `json:"role_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	if err := h.adminService.AssignRole(req.UserID, req.RoleCode); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "分配角色失败", nil)
		return
	}
	response.Success(c, nil)
}

// GET /api/admin/members (ADMIN-MALL-004)
func (h *AdminHandler) ListMembers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")

	users, total, err := h.adminService.ListMembers(page, size, keyword)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询失败", nil)
		return
	}
	response.Success(c, gin.H{"list": users, "total": total})
}

// --- RBAC endpoints ---

// POST /api/admin/roles/:id/permissions
func (h *AdminHandler) BindRolePermissions(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "角色ID格式错误", nil)
		return
	}
	var req struct {
		PermissionIDs []int64 `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	if err := h.adminService.BindRolePermissions(roleID, req.PermissionIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "绑定权限失败", nil)
		return
	}
	response.Success(c, nil)
}

// GET /api/admin/roles/:id/permissions
func (h *AdminHandler) GetRolePermissions(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "角色ID格式错误", nil)
		return
	}
	perms, err := h.adminService.GetRolePermissions(roleID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询失败", nil)
		return
	}
	response.Success(c, gin.H{"list": perms, "total": len(perms)})
}

// POST /api/admin/users/:id/roles
func (h *AdminHandler) AssignUserRoles(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "用户ID格式错误", nil)
		return
	}
	var req struct {
		RoleIDs []int64 `json:"role_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	if err := h.adminService.AssignUserRoles(userID, req.RoleIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "分配角色失败", nil)
		return
	}
	response.Success(c, nil)
}

// GET /api/admin/users/:id/roles
func (h *AdminHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "用户ID格式错误", nil)
		return
	}
	roles, err := h.adminService.GetUserRoles(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询失败", nil)
		return
	}
	response.Success(c, gin.H{"list": roles, "total": len(roles)})
}

// GET /api/admin/permissions
func (h *AdminHandler) ListPermissions(c *gin.Context) {
	perms, err := h.adminService.ListAllPermissions()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询失败", nil)
		return
	}
	response.Success(c, gin.H{"list": perms, "total": len(perms)})
}

// GET /api/admin/menus
func (h *AdminHandler) ListMenus(c *gin.Context) {
	menus, err := h.adminService.ListAllMenus()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询失败", nil)
		return
	}
	response.Success(c, gin.H{"list": menus, "total": len(menus)})
}

// GET /api/admin/user-login-logs (LOG-001)
func (h *AdminHandler) QueryUserLoginLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	mobile := c.Query("mobile")
	var success *bool
	if s := c.Query("success"); s != "" {
		v := s == "true"
		success = &v
	}

	logs, total, err := h.loginLogService.QueryUserLogs(page, size, mobile, success)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询失败", nil)
		return
	}
	response.Success(c, gin.H{"list": logs, "total": total})
}

// POST /api/admin/users/update-balance
func (h *AdminHandler) UpdateUserBalance(c *gin.Context) {
	operatorID := c.GetInt64("userID")
	var req struct {
		UserID int64   `json:"user_id" binding:"required"`
		Amount float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数错误: "+err.Error(), nil)
		return
	}
	if err := h.adminService.UpdateUserBalance(req.UserID, req.Amount, operatorID); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "余额调整失败: "+err.Error(), nil)
		return
	}
	response.Success(c, nil)
}

// GET /api/admin/admin-login-logs (LOG-002)
func (h *AdminHandler) QueryAdminLoginLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	mobile := c.Query("mobile")
	var success *bool
	if s := c.Query("success"); s != "" {
		v := s == "true"
		success = &v
	}

	logs, total, err := h.loginLogService.QueryAdminLogs(page, size, mobile, success)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询失败", nil)
		return
	}
	response.Success(c, gin.H{"list": logs, "total": total})
}
