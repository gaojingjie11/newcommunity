package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"smartcommunity-microservices/app/gateway/api/internal/perm"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	mall "smartcommunity-microservices/app/mall/rpc/types/mall"
	user "smartcommunity-microservices/app/user/rpc/user"
	"smartcommunity-microservices/app/user/rpc/userrpc"
	"smartcommunity-microservices/common/auth"

	goredis "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
)

const (
	CtxKeyUserId   = "x-user-id"
	CtxKeyUserRole = "x-user-role"
	CtxKeyStoreIds = "x-store-ids"
	CtxKeyIsAdmin  = "x-is-admin"
)

func isLocalRoute(path string) bool {
	if path == "/health" || path == "/api/upload" {
		return true
	}
	if strings.HasPrefix(path, "/api/users") || strings.HasPrefix(path, "/api/agent") {
		return true
	}
	return false
}

func requiresAuth(method, path string) bool {
	// Admin and statistics always require auth
	if strings.HasPrefix(path, "/api/admin/") || strings.HasPrefix(path, "/api/statistics/") {
		return true
	}
	// Python agent routes (except health) require auth
	if strings.HasPrefix(path, "/agent/") {
		if path == "/agent/health" {
			return false
		}
		return true
	}
	// Local routes:
	// 1. Users endpoints
	if strings.HasPrefix(path, "/api/users/") {
		// Only /api/users/me is authenticated
		if strings.HasPrefix(path, "/api/users/me") {
			return true
		}
		return false
	}
	// 2. Mall endpoints
	if strings.HasPrefix(path, "/api/mall/") {
		// Public GET endpoints: products, categories, stores, comments
		if method == http.MethodGet {
			if strings.HasPrefix(path, "/api/mall/products") ||
				path == "/api/mall/categories" ||
				path == "/api/mall/stores" ||
				strings.HasPrefix(path, "/api/mall/comments") {
				return false
			}
		}
		// Alipay notify callback is public
		if path == "/api/mall/payments/alipay/notify" {
			return false
		}
		return true
	}
	// 3. Community endpoints
	if strings.HasPrefix(path, "/api/community/") {
		// Public GET endpoints: ping, notices
		if method == http.MethodGet {
			if path == "/api/community/ping" ||
				strings.HasPrefix(path, "/api/community/notices") {
				return false
			}
		}
		return true
	}
	// Default to false for health/upload/etc.
	if path == "/health" || path == "/api/upload" {
		return false
	}

	return false
}

func resolveServiceName(path string) string {
	return ""
}

func proxyMiddleware(svcCtx *svc.ServiceContext) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isLocalRoute(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			serveProxy(w, r, svcCtx)
		})
	}
}

func permCheckMiddleware(svcCtx *svc.ServiceContext) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if requiresAuth(r.Method, path) {
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					writeJSONError(w, 401, "请先登录")
					return
				}

				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) != 2 || parts[0] != "Bearer" {
					writeJSONError(w, 401, "Token格式错误")
					return
				}

				tokenStr := parts[1]
				claims, err := auth.ParseToken(svcCtx.Config.Auth.AccessSecret, tokenStr)
				if err != nil {
					writeJSONError(w, 401, "Token无效或已过期")
					return
				}

				redisKey := fmt.Sprintf("login:token:%d", claims.UserID)
				cachedToken, err := svcCtx.RedisClient.Get(r.Context(), redisKey).Result()
				if err != nil || cachedToken != tokenStr {
					writeJSONError(w, 401, "登录已失效，请重新登录")
					return
				}

				userID := claims.UserID
				userRole := claims.Role

				// Set default ctx values
				var storeIdsStr string = ""
				var isAdmin bool = false

				// Check if it's admin/statistics area
				if strings.HasPrefix(path, "/api/admin/") || strings.HasPrefix(path, "/api/statistics/") {
					if userRole == "admin" {
						isAdmin = true
					} else {
						requiredPerm := perm.LookupPermission(r.Method, path)
						if requiredPerm == "" {
							writeJSONError(w, 403, "当前管理接口未配置权限点，已拒绝访问")
							return
						}

						codes, err := loadPermissions(r.Context(), svcCtx.RedisClient, svcCtx.UserRpc, userID)
						if err != nil {
							writeJSONError(w, 500, "权限校验失败")
							return
						}

						hasFull := false
						hasRestricted := false
						requiredPermAll := requiredPerm + "_all"

						for _, code := range codes {
							if code == requiredPermAll {
								hasFull = true
								break
							}
							if code == requiredPerm {
								hasRestricted = true
							}
						}

						if hasFull {
							isAdmin = true
						} else if hasRestricted {
							isAdmin = true
							// Restricted view: fetch store ids
							storesResp, err := svcCtx.MallRpc.GetUserStores(r.Context(), &mall.UserIDReq{UserId: userID})
							if err != nil {
								writeJSONError(w, 500, fmt.Sprintf("获取绑定门店失败: %v", err))
								return
							}
							if len(storesResp.StoreIds) == 0 {
								storeIdsStr = "none"
							} else {
								var idStrs []string
								for _, sid := range storesResp.StoreIds {
									idStrs = append(idStrs, fmt.Sprintf("%d", sid))
								}
								storeIdsStr = strings.Join(idStrs, ",")
							}
						} else {
							writeJSONError(w, 403, "无权限访问此资源")
							return
						}
					}
				}

				// Inject to context
				ctx := r.Context()
				ctx = context.WithValue(ctx, CtxKeyUserId, userID)
				ctx = context.WithValue(ctx, CtxKeyUserRole, userRole)
				ctx = context.WithValue(ctx, CtxKeyIsAdmin, isAdmin)
				if storeIdsStr != "" {
					ctx = context.WithValue(ctx, CtxKeyStoreIds, storeIdsStr)
				}
				r = r.WithContext(ctx)
				fmt.Printf("[DEBUG] permCheckMiddleware: path=%s, userID=%d, role=%s, isAdmin=%t, storeIds=%s\n", path, userID, userRole, isAdmin, storeIdsStr)
			}

			next(w, r)
		}
	}
}

func serveProxy(w http.ResponseWriter, r *http.Request, svcCtx *svc.ServiceContext) {
	if handleUserAdminLocal(w, r, svcCtx) {
		return
	}
	serviceName := resolveServiceName(r.URL.Path)
	if serviceName == "" {
		writeJSONError(w, 404, fmt.Sprintf("no upstream service for path: %s", r.URL.Path))
		return
	}

	baseTarget := svcCtx.Config.Gateway.Services[serviceName]
	if baseTarget == "" {
		writeJSONError(w, 502, fmt.Sprintf("service unavailable: %s", serviceName))
		return
	}

	target, err := url.Parse(baseTarget)
	if err != nil {
		writeJSONError(w, 500, fmt.Sprintf("invalid upstream url for service %s: %v", serviceName, err))
		return
	}

	userID, _ := r.Context().Value(CtxKeyUserId).(int64)
	userRole, _ := r.Context().Value(CtxKeyUserRole).(string)
	storeIds, _ := r.Context().Value(CtxKeyStoreIds).(string)
	isAdmin, _ := r.Context().Value(CtxKeyIsAdmin).(bool)

	// Prepare proxy call
	proxy := httputil.NewSingleHostReverseProxy(target)
	r.Host = target.Host

	// Inject tracing and authentication headers
	r.Header.Set("X-Gateway-Proxy", "gateway-api")
	r.Header.Set("X-Gateway-Time", time.Now().Format(time.RFC3339))

	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		r.Header.Set("X-Request-ID", requestID)
	}

	if svcCtx.Config.Gateway.InternalToken != "" {
		r.Header.Set("X-Internal-Token", svcCtx.Config.Gateway.InternalToken)
	}

	if userID > 0 {
		r.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
	}
	if userRole != "" {
		r.Header.Set("X-User-Role", userRole)
	}
	if storeIds != "" {
		r.Header.Set("X-Store-Ids", storeIds)
	}
	if isAdmin {
		r.Header.Set("X-Is-Admin", "true")
	}

	proxy.ServeHTTP(w, r)
}

func loadPermissions(ctx context.Context, rdb *goredis.Client, userRpc userrpc.UserRpc, userID int64) ([]string, error) {
	cacheKey := fmt.Sprintf("rbac:permissions:%d", userID)

	// Try cache first
	cached, err := rdb.SMembers(ctx, cacheKey).Result()
	if err == nil && len(cached) > 0 {
		return cached, nil
	}

	// Cache miss — call user-rpc!
	resp, err := userRpc.GetUserPermissions(ctx, &user.UserIDReq{UserId: userID})
	if err != nil {
		return nil, err
	}
	codes := resp.Permissions

	// Write to cache with 10-min TTL
	if len(codes) > 0 {
		members := make([]interface{}, len(codes))
		for i, c := range codes {
			members[i] = c
		}
		_ = rdb.SAdd(ctx, cacheKey, members...).Err()
		_ = rdb.Expire(ctx, cacheKey, 600*time.Second).Err()
	}

	return codes, nil
}
func writeJSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    nil,
	})
}

func handleUserAdminLocal(w http.ResponseWriter, r *http.Request, svcCtx *svc.ServiceContext) bool {
	if !strings.HasPrefix(r.URL.Path, "/api/admin/users") &&
		!strings.HasPrefix(r.URL.Path, "/api/admin/roles") &&
		!strings.HasPrefix(r.URL.Path, "/api/admin/permissions") {
		return false
	}

	isAdmin, _ := r.Context().Value(CtxKeyIsAdmin).(bool)
	if !isAdmin {
		writeJSONError(w, 403, "无权限访问此资源")
		return true
	}

	var err error

	// 2. Dispatch routes
	if r.Method == http.MethodGet && r.URL.Path == "/api/admin/users" {
		// Get query parameters
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		size, _ := strconv.Atoi(q.Get("size"))
		keyword := q.Get("keyword")

		if page < 1 {
			page = 1
		}
		if size < 1 {
			size = 20
		}

		// Call gRPC UserRpc.ListAdminUsers (large page size to enable memory filtering)
		rpcResp, err := svcCtx.UserRpc.ListAdminUsers(r.Context(), &user.ListAdminUsersReq{
			Page: 1,
			Size: 10000,
		})
		if err != nil {
			writeJSONError(w, 500, fmt.Sprintf("gRPC error: %v", err))
			return true
		}

		// Filter users in memory
		var filtered []*user.UserInfo
		keyword = strings.TrimSpace(strings.ToLower(keyword))
		for _, u := range rpcResp.List {
			if keyword == "" {
				filtered = append(filtered, u)
			} else {
				match := strings.Contains(strings.ToLower(u.Username), keyword) ||
					strings.Contains(strings.ToLower(u.RealName), keyword) ||
					strings.Contains(strings.ToLower(u.Mobile), keyword)
				if match {
					filtered = append(filtered, u)
				}
			}
		}

		total := int64(len(filtered))

		// Paginate
		start := (page - 1) * size
		if start > len(filtered) {
			start = len(filtered)
		}
		end := start + size
		if end > len(filtered) {
			end = len(filtered)
		}

		var paginated []*user.UserInfo
		if start < end {
			paginated = filtered[start:end]
		} else {
			paginated = make([]*user.UserInfo, 0)
		}

		var paginatedMaps []map[string]interface{}
		for _, u := range paginated {
			m := map[string]interface{}{
				"id":              u.Id,
				"username":        u.Username,
				"real_name":       u.RealName,
				"mobile":          u.Mobile,
				"avatar":          u.Avatar,
				"green_points":    u.GreenPoints,
				"role":            u.Role,
				"status":          u.Status,
				"face_registered": u.FaceRegistered,
				"face_image_url":  u.FaceImageUrl,
				"balance":         u.Balance,
			}
			paginatedMaps = append(paginatedMaps, m)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"list":  paginatedMaps,
				"total": total,
			},
		})
		return true
	}

	if r.Method == http.MethodPost && r.URL.Path == "/api/admin/users/freeze" {
		var req struct {
			Id     int64 `json:"id"`
			Status int   `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, 400, "invalid body parameter format")
			return true
		}

		_, err := svcCtx.UserRpc.FreezeUser(r.Context(), &user.FreezeUserReq{
			UserId: req.Id,
			Status: int32(req.Status),
		})
		if err != nil {
			writeJSONError(w, 500, fmt.Sprintf("gRPC error: %v", err))
			return true
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    nil,
		})
		return true
	}

	if r.Method == http.MethodPost && r.URL.Path == "/api/admin/users/assign-role" {
		var req struct {
			UserId   int64   `json:"user_id"`
			RoleIds  []int64 `json:"role_ids"`
			StoreIds []int64 `json:"store_ids"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, 400, "invalid body parameter format")
			return true
		}

		// 1. Fetch all system roles to match role codes
		rolesResp, err := svcCtx.UserRpc.ListRoles(r.Context(), &user.EmptyReq{})
		if err != nil {
			writeJSONError(w, 500, fmt.Sprintf("gRPC error listing roles: %v", err))
			return true
		}

		roleCodeMap := make(map[int64]string)
		for _, rl := range rolesResp.Roles {
			roleCodeMap[rl.Id] = rl.Code
		}

		// Check if "store" is one of the assigned roles
		isStore := false
		for _, rid := range req.RoleIds {
			if roleCodeMap[rid] == "store" {
				isStore = true
				break
			}
		}

		// 2. Assign multiple roles to the user
		_, err = svcCtx.UserRpc.AssignUserRoles(r.Context(), &user.AssignUserRolesReq{
			UserId:  req.UserId,
			RoleIds: req.RoleIds,
		})
		if err != nil {
			writeJSONError(w, 500, fmt.Sprintf("gRPC error: %v", err))
			return true
		}

		// 3. Save store bindings if role is 'store'
		if isStore {
			_, err = svcCtx.MallRpc.BindUserStores(r.Context(), &mall.BindUserStoresReq{
				UserId:   req.UserId,
				StoreIds: req.StoreIds,
			})
			if err != nil {
				writeJSONError(w, 500, fmt.Sprintf("gRPC error binding user stores: %v", err))
				return true
			}
		} else {
			// Purge store bindings if role is changed to something else
			_, err = svcCtx.MallRpc.BindUserStores(r.Context(), &mall.BindUserStoresReq{
				UserId:   req.UserId,
				StoreIds: []int64{},
			})
			if err != nil {
				writeJSONError(w, 500, fmt.Sprintf("gRPC error cleaning user stores: %v", err))
				return true
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    nil,
		})
		return true
	}

	if r.Method == http.MethodGet && r.URL.Path == "/api/admin/users/roles" {
		q := r.URL.Query()
		targetUserID, _ := strconv.ParseInt(q.Get("user_id"), 10, 64)
		if targetUserID == 0 {
			writeJSONError(w, 400, "user_id parameter is required")
			return true
		}

		rpcResp, err := svcCtx.UserRpc.GetUserRoles(r.Context(), &user.UserIDReq{UserId: targetUserID})
		if err != nil {
			writeJSONError(w, 500, fmt.Sprintf("gRPC error: %v", err))
			return true
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    rpcResp.RoleIds,
		})
		return true
	}

	if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/admin/users/") && strings.HasSuffix(r.URL.Path, "/stores") {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) >= 6 {
			userIdStr := parts[4]
			userId, _ := strconv.ParseInt(userIdStr, 10, 64)
			if userId > 0 {
				resp, err := svcCtx.MallRpc.GetUserStores(r.Context(), &mall.UserIDReq{UserId: userId})
				if err != nil {
					writeJSONError(w, 500, fmt.Sprintf("gRPC error: %v", err))
					return true
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    0,
					"message": "success",
					"data":    resp.StoreIds,
				})
				return true
			}
		}
		writeJSONError(w, 400, "invalid user id")
		return true
	}

	if r.Method == http.MethodPost && r.URL.Path == "/api/admin/users/update-balance" {
		var req struct {
			UserId int64   `json:"user_id"`
			Amount float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, 400, "invalid body parameter format")
			return true
		}

		_, err := svcCtx.UserRpc.UpdateUserBalance(r.Context(), &user.UpdateUserBalanceReq{
			UserId: req.UserId,
			Amount: req.Amount,
		})
		if err != nil {
			writeJSONError(w, 500, fmt.Sprintf("gRPC error: %v", err))
			return true
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    nil,
		})
		return true
	}

	// GET /api/admin/permissions
	if r.Method == http.MethodGet && r.URL.Path == "/api/admin/permissions" {
		resp, err := svcCtx.UserRpc.ListPermissions(r.Context(), &user.EmptyReq{})
		if err != nil {
			writeJSONError(w, 500, fmt.Sprintf("gRPC error: %v", err))
			return true
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    resp.Permissions,
		})
		return true
	}

	// GET /api/admin/roles
	if r.Method == http.MethodGet && r.URL.Path == "/api/admin/roles" {
		resp, err := svcCtx.UserRpc.ListRoles(r.Context(), &user.EmptyReq{})
		if err != nil {
			writeJSONError(w, 500, fmt.Sprintf("gRPC error: %v", err))
			return true
		}
		var list []map[string]interface{}
		for _, rl := range resp.Roles {
			list = append(list, map[string]interface{}{
				"id":     rl.Id,
				"name":   rl.Name,
				"code":   rl.Code,
				"remark": rl.Remark,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    list,
		})
		return true
	}

	// POST /api/admin/roles
	if r.Method == http.MethodPost && r.URL.Path == "/api/admin/roles" {
		var req struct {
			Name   string `json:"name"`
			Code   string `json:"code"`
			Remark string `json:"remark"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, 400, "invalid body parameter format")
			return true
		}
		res, err := svcCtx.UserRpc.CreateRole(r.Context(), &user.CreateRoleReq{
			Name:   req.Name,
			Code:   req.Code,
			Remark: req.Remark,
		})
		if err != nil {
			writeJSONError(w, 500, fmt.Sprintf("gRPC error: %v", err))
			return true
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"id":     res.Id,
				"name":   res.Name,
				"code":   res.Code,
				"remark": res.Remark,
			},
		})
		return true
	}

	// PUT /api/admin/roles
	if r.Method == http.MethodPut && r.URL.Path == "/api/admin/roles" {
		var req struct {
			Id     int64  `json:"id"`
			Name   string `json:"name"`
			Code   string `json:"code"`
			Remark string `json:"remark"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, 400, "invalid body parameter format")
			return true
		}
		_, err := svcCtx.UserRpc.UpdateRole(r.Context(), &user.UpdateRoleReq{
			Id:     req.Id,
			Name:   req.Name,
			Code:   req.Code,
			Remark: req.Remark,
		})
		if err != nil {
			writeJSONError(w, 500, fmt.Sprintf("gRPC error: %v", err))
			return true
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    nil,
		})
		return true
	}

	// DELETE /api/admin/roles
	if r.Method == http.MethodDelete && (r.URL.Path == "/api/admin/roles" || strings.HasPrefix(r.URL.Path, "/api/admin/roles/")) {
		var roleId int64
		if r.URL.Path == "/api/admin/roles" {
			idStr := r.URL.Query().Get("id")
			roleId, _ = strconv.ParseInt(idStr, 10, 64)
		} else {
			idStr := strings.TrimPrefix(r.URL.Path, "/api/admin/roles/")
			roleId, _ = strconv.ParseInt(idStr, 10, 64)
		}
		if roleId == 0 {
			var req struct {
				Id int64 `json:"id"`
			}
			_ = json.NewDecoder(r.Body).Decode(&req)
			roleId = req.Id
		}
		if roleId == 0 {
			writeJSONError(w, 400, "invalid role id")
			return true
		}
		_, err := svcCtx.UserRpc.DeleteRole(r.Context(), &user.DeleteRoleReq{Id: roleId})
		if err != nil {
			writeJSONError(w, 500, fmt.Sprintf("gRPC error: %v", err))
			return true
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    nil,
		})
		return true
	}

	// GET /api/admin/roles/:id/permissions
	if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/admin/roles/") && strings.HasSuffix(r.URL.Path, "/permissions") {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) >= 6 {
			roleIdStr := parts[4]
			roleId, _ := strconv.ParseInt(roleIdStr, 10, 64)
			if roleId > 0 {
				resp, err := svcCtx.UserRpc.GetRolePermissions(r.Context(), &user.GetRolePermissionsReq{RoleId: roleId})
				if err != nil {
					writeJSONError(w, 500, fmt.Sprintf("gRPC error: %v", err))
					return true
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    0,
					"message": "success",
					"data":    resp.Permissions,
				})
				return true
			}
		}
		writeJSONError(w, 400, "invalid role id")
		return true
	}

	// POST /api/admin/roles/:id/permissions
	if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/admin/roles/") && strings.HasSuffix(r.URL.Path, "/permissions") {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) >= 6 {
			roleIdStr := parts[4]
			roleId, _ := strconv.ParseInt(roleIdStr, 10, 64)
			if roleId > 0 {
				var req struct {
					Permissions []string `json:"permissions"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					writeJSONError(w, 400, "invalid body parameter format")
					return true
				}
				_, err = svcCtx.UserRpc.BindRolePermissions(r.Context(), &user.BindRolePermissionsReq{
					RoleId:          roleId,
					PermissionCodes: req.Permissions,
				})
				if err != nil {
					writeJSONError(w, 500, fmt.Sprintf("gRPC error: %v", err))
					return true
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    0,
					"message": "success",
					"data":    nil,
				})
				return true
			}
		}
		writeJSONError(w, 400, "invalid role id")
		return true
	}

	writeJSONError(w, 404, "admin subroute not found")
	return true
}
