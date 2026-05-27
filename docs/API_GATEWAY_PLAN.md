# API Gateway Plan

`gateway-service` is the unified entry for new microservice APIs.

## Completed Features

- `GET /health` — health check
- `GET /api/gateway/services` — service list with source (nacos/config)
- Unified JWT validation at gateway level
- RBAC permission checking for `/api/admin/**` and `/api/statistics/**`
- Nacos service discovery with local config fallback
- X-Internal-Token injection for service-to-service calls
- X-User-ID and X-User-Role injection from JWT claims
- Three-tier access control: public, authenticated, permission_required

## Architecture

```
客户端
  │
  ▼
gateway-service (:8000)
  ├─ CORS / RequestID / Logger / Recovery
  ├─ Public routes → 直接代理
  ├─ Authenticated routes → JWTAuth → 代理
  └─ Permission routes → JWTAuth → RequirePermission → 代理
  │
  ▼
下游服务 (user/mall/community/workorder/statistics/agent)
  └─ 保留二次校验 (JWTAuth + RequirePermission)
```

## Route Mapping

| 路径前缀 | 下游服务 |
|----------|----------|
| `/api/users/**` | user-service |
| `/api/admin/roles/**`, `/api/admin/users/**`, `/api/admin/members/**`, `/api/admin/permissions/**`, `/api/admin/menus/**`, `/api/admin/*-logs` | user-service |
| `/api/mall/**`, `/api/admin/mall/**` | mall-service |
| `/api/community/**`, `/api/admin/community/**`, `/api/workorders/**`, `/api/admin/workorders/**` | community-service |
| `/api/statistics/**` | community-service |
| `/agent/**` | agent-service |

## Permission Mapping

Gateway 维护一份 path → permission code 映射表（`internal/perm/mapping.go`），用于 `/api/admin/**` 和 `/api/statistics/**` 路由的权限校验。

权限码定义见 `docs/RBAC_DESIGN.md`。

## Service Discovery

1. **优先**：从 Nacos 获取服务实例地址（每 30 秒刷新）
2. **兜底**：如果 Nacos 不可用，从 `configs/config.yaml` 的 `gateway.services` 读取地址
3. Nacos 不可用不导致 gateway 启动失败，仅输出 warning

## Future Enhancements

- Add request aggregation for dashboard and profile views
- Add route-level timeouts and retry policies
- Add rate limiting
- Add circuit breaker for downstream services
- Keep frontend migration incremental by moving one API domain at a time
