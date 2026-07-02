package perm

import (
	"net/http"
	"testing"
)

func TestLookupPermissionMappedRoute(t *testing.T) {
	code := LookupPermission(http.MethodPut, "/api/admin/mall/products/123")
	if code != "mall:product:update" {
		t.Fatalf("LookupPermission returned %q, want %q", code, "mall:product:update")
	}
}

func TestLookupPermissionUnmappedAdminRoute(t *testing.T) {
	code := LookupPermission(http.MethodPatch, "/api/admin/mall/products/123")
	if code != "" {
		t.Fatalf("LookupPermission returned %q for unmapped route, want empty string", code)
	}
}
