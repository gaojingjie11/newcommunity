package perm

import (
	"net/http"
	"sort"
	"strings"
)

// RoutePermission maps HTTP method + path pattern to a required permission code.
// Patterns support wildcard (*) segments that match exactly one path segment.
// Matching is segment-exact: "/api/admin/roles" matches "/api/admin/roles" but NOT "/api/admin/roles/1/permissions".
type RoutePermission struct {
	Method  string
	Pattern string
	Code    string
}

// routeTable is the merged and sorted permission table, built at init time.
var routeTable []RoutePermission

func init() {
	table := []RoutePermission{
		// ── user-service admin ──
		{http.MethodPost, "/api/admin/roles/*/menus", "rbac:role:bind_menu"},
		{http.MethodPost, "/api/admin/roles/*/permissions", "rbac:role:bind_permission"},
		{http.MethodGet, "/api/admin/roles/*/permissions", "rbac:role:get_permissions"},
		{http.MethodPost, "/api/admin/roles", "rbac:role:create"},
		{http.MethodPut, "/api/admin/roles", "rbac:role:update"},
		{http.MethodDelete, "/api/admin/roles", "rbac:role:delete"},
		{http.MethodGet, "/api/admin/roles", "rbac:role:list"},
		{http.MethodPost, "/api/admin/users/*/roles", "rbac:user:assign_roles"},
		{http.MethodGet, "/api/admin/users/*/roles", "rbac:user:get_roles"},
		{http.MethodGet, "/api/admin/users", "rbac:user:list"},
		{http.MethodPost, "/api/admin/users/freeze", "rbac:user:freeze"},
		{http.MethodPost, "/api/admin/users/assign-role", "rbac:user:assign_role"},
		{http.MethodPost, "/api/admin/users/update-balance", "rbac:user:update_balance"},
		{http.MethodGet, "/api/admin/members", "rbac:member:list"},
		{http.MethodGet, "/api/admin/permissions", "rbac:permission:list"},
		{http.MethodGet, "/api/admin/menus", "rbac:menu:list"},
		{http.MethodGet, "/api/admin/user-login-logs", "log:user_login:list"},
		{http.MethodGet, "/api/admin/admin-login-logs", "log:admin_login:list"},

		// ── mall-service admin ──
		{http.MethodPost, "/api/admin/mall/promotions/*/products", "mall:promotion:bind_product"},
		{http.MethodPost, "/api/admin/mall/orders/*/ship", "mall:order:ship"},
		{http.MethodPost, "/api/admin/mall/orders/*/cancel", "mall:order:cancel"},
		{http.MethodPost, "/api/admin/mall/categories", "mall:category:create"},
		{http.MethodPut, "/api/admin/mall/categories/*", "mall:category:update"},
		{http.MethodDelete, "/api/admin/mall/categories/*", "mall:category:delete"},
		{http.MethodGet, "/api/admin/mall/products", "mall:product:list"},
		{http.MethodPost, "/api/admin/mall/products", "mall:product:create"},
		{http.MethodPut, "/api/admin/mall/products/*", "mall:product:update"},
		{http.MethodDelete, "/api/admin/mall/products/*", "mall:product:delete"},
		{http.MethodPost, "/api/admin/mall/promotions", "mall:promotion:create"},
		{http.MethodPut, "/api/admin/mall/promotions/*", "mall:promotion:update"},
		{http.MethodDelete, "/api/admin/mall/promotions/*", "mall:promotion:delete"},
		{http.MethodPost, "/api/admin/mall/service-areas", "mall:service_area:create"},
		{http.MethodPut, "/api/admin/mall/service-areas/*", "mall:service_area:update"},
		{http.MethodDelete, "/api/admin/mall/service-areas/*", "mall:service_area:delete"},
		{http.MethodPost, "/api/admin/mall/stores", "mall:store:create"},
		{http.MethodPut, "/api/admin/mall/stores/*", "mall:store:update"},
		{http.MethodDelete, "/api/admin/mall/stores/*", "mall:store:delete"},
		{http.MethodPut, "/api/admin/mall/store-products/status", "mall:store_product:status"},
		{http.MethodPut, "/api/admin/mall/store-products/stock", "mall:store_product:stock"},
		{http.MethodPost, "/api/admin/mall/store-products", "mall:store_product:bind"},
		{http.MethodDelete, "/api/admin/mall/store-products", "mall:store_product:unbind"},
		{http.MethodGet, "/api/admin/mall/store-products/*", "mall:store_product:list"},
		{http.MethodGet, "/api/admin/mall/orders", "mall:order:list"},

		// ── community-service admin ──
		{http.MethodGet, "/api/admin/community/notices/*/views", "community:notice:views"},
		{http.MethodPost, "/api/admin/community/visitors/*/audit", "community:visitor:audit"},
		{http.MethodPost, "/api/admin/community/parking-spaces/*/assign", "community:parking:assign"},
		{http.MethodGet, "/api/admin/community/parking-spaces/statistics", "community:parking:statistics"},
		{http.MethodGet, "/api/admin/community/property-fees/payments", "community:fee:payment_list"},
		{http.MethodGet, "/api/admin/community/notices", "community:notice:list"},
		{http.MethodPost, "/api/admin/community/notices", "community:notice:create"},
		{http.MethodDelete, "/api/admin/community/notices/*", "community:notice:delete"},
		{http.MethodGet, "/api/admin/community/visitors", "community:visitor:list"},
		{http.MethodGet, "/api/admin/community/parking-spaces", "community:parking:list"},
		{http.MethodPost, "/api/admin/community/parking-spaces", "community:parking:create"},
		{http.MethodGet, "/api/admin/community/property-fees", "community:fee:list"},
		{http.MethodPost, "/api/admin/community/property-fees", "community:fee:create"},

		// ── workorder-service admin ──
		{http.MethodGet, "/api/admin/workorders", "workorder:repair:list"},
		{http.MethodPost, "/api/admin/workorders/*/process", "workorder:repair:process"},

		// ── statistics-service (explicit per-endpoint) ──
		{http.MethodPost, "/api/statistics/ai-report/generate", "statistics:ai_report:generate"},
		{http.MethodGet, "/api/statistics/ai-report/latest", "statistics:ai_report:read"},
		{http.MethodGet, "/api/statistics/ai-report/list", "statistics:ai_report:read"},
		{http.MethodGet, "/api/statistics/ai-report/*", "statistics:ai_report:read"},
		{http.MethodGet, "/api/statistics/products/sales-rank", "statistics:product:sales_rank"},
		{http.MethodGet, "/api/statistics/products/view-rank", "statistics:product:view_rank"},
		{http.MethodGet, "/api/statistics/community/overview", "statistics:community:overview"},
		{http.MethodGet, "/api/statistics/orders", "statistics:order:summary"},
		{http.MethodGet, "/api/statistics/workorders", "statistics:workorder:summary"},
	}

	// Sort by pattern segment count descending (longer/more specific first).
	sort.Slice(table, func(i, j int) bool {
		pi := len(strings.Split(strings.Trim(table[i].Pattern, "/"), "/"))
		pj := len(strings.Split(strings.Trim(table[j].Pattern, "/"), "/"))
		if pi != pj {
			return pi > pj
		}
		// Same length: prefer patterns with fewer wildcards (more specific)
		wi := strings.Count(table[i].Pattern, "*")
		wj := strings.Count(table[j].Pattern, "*")
		return wi < wj
	})
	routeTable = table
}

// LookupPermission returns the permission code for the given method+path, or "" if not mapped.
func LookupPermission(method, path string) string {
	for _, rp := range routeTable {
		if rp.Method == method && matchPath(path, rp.Pattern) {
			return rp.Code
		}
	}
	return ""
}

// IsAdminOrStatsPath returns true if the path should be protected by permission_required.
// All /api/admin/** and /api/statistics/** paths are protected.
func IsAdminOrStatsPath(path string) bool {
	return strings.HasPrefix(path, "/api/admin/") || strings.HasPrefix(path, "/api/statistics/")
}

// matchPath checks if requestPath matches the route pattern with segment-exact matching.
// Each segment must match exactly, except:
//   - "*" matches exactly one arbitrary segment
//   - "**" at the end matches one or more remaining segments (greedy)
//
// Examples:
//
//	matchPath("/api/admin/roles/1/permissions", "/api/admin/roles/*/permissions") = true
//	matchPath("/api/admin/roles", "/api/admin/roles")                            = true
//	matchPath("/api/admin/roles/1/permissions", "/api/admin/roles")              = false (segment count differs)
//	matchPath("/api/admin/mall/orders/1/ship", "/api/admin/mall/orders/*/ship")  = true
func matchPath(requestPath, pattern string) bool {
	reqParts := strings.Split(strings.Trim(requestPath, "/"), "/")
	patParts := strings.Split(strings.Trim(pattern, "/"), "/")

	for i, pat := range patParts {
		if pat == "**" {
			// "**" at the end matches everything remaining
			return true
		}
		if i >= len(reqParts) {
			return false
		}
		if pat == "*" {
			continue // matches exactly one segment
		}
		if reqParts[i] != pat {
			return false
		}
	}
	// Pattern consumed — request must have the same number of segments
	return len(reqParts) == len(patParts)
}
