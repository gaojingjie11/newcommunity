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
	"smartcommunity-microservices/app/user/rpc/user"
	"smartcommunity-microservices/app/user/rpc/userrpc"
	"smartcommunity-microservices/common/auth"

	goredis "github.com/redis/go-redis/v9"
)

type RouteType int

const (
	RouteTypePublic RouteType = iota
	RouteTypeAuthenticated
	RouteTypeAdmin
)

func isLocalRoute(path string) bool {
	if path == "/health" || path == "/api/upload" {
		return true
	}
	if strings.HasPrefix(path, "/api/users") {
		return true
	}
	return false
}

func getRouteType(path string) RouteType {
	// Agent health
	if path == "/agent/health" {
		return RouteTypePublic
	}

	return RouteTypeAuthenticated
}

func resolveServiceName(path string) string {
	switch {
	case strings.HasPrefix(path, "/agent"):
		return "agent-service"
	default:
		return ""
	}
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

	// Authorization and identification
	routeType := getRouteType(r.URL.Path)
	var userID int64
	var userRole string

	if routeType == RouteTypeAuthenticated || routeType == RouteTypeAdmin {
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

		userID = claims.UserID
		userRole = claims.Role

		// Enforce admin permission required
		if routeType == RouteTypeAdmin {
			if userRole != "admin" {
				requiredPerm := perm.LookupPermission(r.Method, r.URL.Path)
				if requiredPerm == "" {
					writeJSONError(w, 403, "该管理接口未配置权限映射")
					return
				}

				codes, err := loadPermissions(r.Context(), svcCtx.RedisClient, svcCtx.UserRpc, userID)
				if err != nil {
					writeJSONError(w, 500, "权限校验失败")
					return
				}

				hasPerm := false
				for _, code := range codes {
					if code == requiredPerm {
						hasPerm = true
						break
					}
				}

				if !hasPerm {
					writeJSONError(w, 403, "无权限访问此资源")
					return
				}
			}
		}
	}

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
		_ = rdb.Expire(ctx, cacheKey, 600).Err()
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
	if !strings.HasPrefix(r.URL.Path, "/api/admin/users") {
		return false
	}

	// 1. Authenticate admin user
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeJSONError(w, 401, "请先登录")
		return true
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		writeJSONError(w, 401, "Token格式错误")
		return true
	}

	tokenStr := parts[1]
	claims, err := auth.ParseToken(svcCtx.Config.Auth.AccessSecret, tokenStr)
	if err != nil {
		writeJSONError(w, 401, "Token无效或已过期")
		return true
	}

	redisKey := fmt.Sprintf("login:token:%d", claims.UserID)
	cachedToken, err := svcCtx.RedisClient.Get(r.Context(), redisKey).Result()
	if err != nil || cachedToken != tokenStr {
		writeJSONError(w, 401, "登录已失效，请重新登录")
		return true
	}

	if claims.Role != "admin" {
		writeJSONError(w, 403, "无权限访问此资源")
		return true
	}

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
			UserId   int64  `json:"user_id"`
			RoleCode string `json:"role_code"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, 400, "invalid body parameter format")
			return true
		}

		// Query all roles to match role code to ID
		rolesResp, err := svcCtx.UserRpc.ListRoles(r.Context(), &user.EmptyReq{})
		if err != nil {
			writeJSONError(w, 500, fmt.Sprintf("gRPC error listing roles: %v", err))
			return true
		}

		var targetRoleId int64 = 0
		for _, rl := range rolesResp.Roles {
			if rl.Code == req.RoleCode {
				targetRoleId = rl.Id
				break
			}
		}

		if targetRoleId == 0 {
			writeJSONError(w, 400, fmt.Sprintf("role code '%s' not found", req.RoleCode))
			return true
		}

		_, err = svcCtx.UserRpc.AssignRole(r.Context(), &user.AssignRoleReq{
			UserId: req.UserId,
			RoleId: targetRoleId,
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

	writeJSONError(w, 404, "admin subroute not found")
	return true
}
