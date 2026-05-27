# RBAC 权限体系设计

Date: 2026-05-21

## 1. 概述

当前 user-service 已有 `users.role` 字符串字段和 `RequireRole("admin")` 中间件，但这不是完整 RBAC。本文档定义基于角色的访问控制（RBAC）完整设计，支持：

- 用户多角色
- 角色绑定菜单（前端可见性）
- 角色绑定权限点（后端 API 访问控制）
- 权限中间件（RequirePermission）
- Redis 缓存加速权限校验

## 2. 实体关系

```
SysUser (用户主体)
  ├── SysUserRole (多对多) ── SysRole (角色)
  │                              ├── SysRoleMenu (多对多) ── SysMenu (前端菜单)
  │                              └── SysRolePermission (多对多) ── SysPermission (后端权限点)
  └── users.role (legacy 兼容字段，不作为主权限依据)
```

### 2.1 SysUser

- 仅表示用户主体（身份信息、登录凭证）。
- `role` 字段保留为 legacy 兼容，不作为主权限依据。
- 用户实际拥有的角色通过 `sys_user_role` 表关联。

### 2.2 SysRole

- 表示角色，如 `admin`（系统管理员）、`property`（物业管理员）、`store`（门店管理员）、`user`（普通用户）。
- `code` 字段为角色唯一标识，用于程序判断。

### 2.3 SysUserRole

- 用户与角色的多对多关系。
- 一个用户可以拥有多个角色，权限取并集。

### 2.4 SysMenu

- 前端菜单，控制用户在管理端能看到哪些页面。
- 通过 `sys_role_menu` 与角色绑定。

### 2.5 SysPermission

- 后端 API 权限点，每个权限点对应一个或多个 API 接口。
- `code` 为权限唯一标识，格式如 `rbac:role:create`。
- `resource` 表示资源类型，`method` 表示 HTTP 方法，`path` 表示 API 路径。

### 2.6 SysRolePermission

- 角色与权限点的多对多关系。
- 决定角色可以访问哪些后端 API。

## 3. 权限校验流程

### 3.1 JWT 与会话

- JWT Claims 仅包含 `UserID` 和 `Role`（legacy 兼容）。
- JWT 不包含完整权限列表（避免 token 过大）。
- 登录成功后，token 存储在 Redis `login:token:{userID}`，强制单会话。

### 3.2 权限中间件 (RequirePermission)

```
请求 → JWTAuth (解析 token, 验证 Redis) → RequirePermission(permissionCode) → 业务 handler
```

校验流程：
1. 从 gin context 获取 `userID`。
2. 优先从 Redis 缓存 `rbac:permissions:{userID}` 读取权限列表。
3. 缓存未命中时，通过 `PermissionProvider` 回调从 DB 查询：
   - 查 `sys_user_role` 获取用户所有 roleID。
   - 查 `sys_role_permission` 获取所有 permission code。
4. 检查请求的 `permissionCode` 是否在用户权限列表中。
5. 无权限返回 403。

### 3.3 权限缓存

- Redis key: `rbac:permissions:{userID}`
- 类型: SET（存储 permission code 字符串）
- TTL: 10 分钟
- 失效时机:
  - 用户角色变更时，删除缓存。
  - 角色权限变更时，删除该角色所有用户的缓存（延迟策略：设置较短 TTL 自然过期）。

## 4. 权限点定义

### 4.1 权限编码规范

格式: `{domain}:{resource}:{action}`

### 4.2 当前权限点

| Permission Code | 描述 | HTTP Method | Path |
|---|---|---|---|
| rbac:role:create | 创建角色 | POST | /api/admin/roles |
| rbac:role:update | 更新角色 | PUT | /api/admin/roles |
| rbac:role:delete | 删除角色 | DELETE | /api/admin/roles |
| rbac:role:list | 查询角色列表 | GET | /api/admin/roles |
| rbac:role:bind_menu | 角色绑定菜单 | POST | /api/admin/roles/:id/menus |
| rbac:role:bind_permission | 角色绑定权限 | POST | /api/admin/roles/:id/permissions |
| rbac:user:list | 查询管理员列表 | GET | /api/admin/users |
| rbac:user:freeze | 冻结/解冻用户 | POST | /api/admin/users/freeze |
| rbac:user:assign_role | 分配用户角色 | POST | /api/admin/users/assign-role |
| rbac:user:assign_roles | 分配用户多角色 | POST | /api/admin/users/:id/roles |
| rbac:user:get_roles | 查询用户角色 | GET | /api/admin/users/:id/roles |
| rbac:member:list | 查询会员列表 | GET | /api/admin/members |
| rbac:permission:list | 查询权限列表 | GET | /api/admin/permissions |
| rbac:menu:list | 查询菜单列表 | GET | /api/admin/menus |
| rbac:role:get_permissions | 查询角色权限 | GET | /api/admin/roles/:id/permissions |
| log:user_login:list | 查询用户登录日志 | GET | /api/admin/user-login-logs |
| log:admin_login:list | 查询管理员登录日志 | GET | /api/admin/admin-login-logs |

#### Mall-Service 权限点

| Permission Code | 描述 | HTTP Method | Path |
|---|---|---|---|
| mall:category:create | 创建商品分类 | POST | /api/admin/mall/categories |
| mall:category:update | 更新商品分类 | PUT | /api/admin/mall/categories/:id |
| mall:category:delete | 删除商品分类 | DELETE | /api/admin/mall/categories/:id |
| mall:product:create | 创建商品 | POST | /api/admin/mall/products |
| mall:product:update | 更新商品 | PUT | /api/admin/mall/products/:id |
| mall:product:delete | 删除商品 | DELETE | /api/admin/mall/products/:id |
| mall:promotion:create | 创建促销 | POST | /api/admin/mall/promotions |
| mall:promotion:update | 更新促销 | PUT | /api/admin/mall/promotions/:id |
| mall:promotion:delete | 删除促销 | DELETE | /api/admin/mall/promotions/:id |
| mall:promotion:bind_product | 促销绑定商品 | POST | /api/admin/mall/promotions/:id/products |
| mall:service_area:create | 创建服务区域 | POST | /api/admin/mall/service-areas |
| mall:service_area:update | 更新服务区域 | PUT | /api/admin/mall/service-areas/:id |
| mall:service_area:delete | 删除服务区域 | DELETE | /api/admin/mall/service-areas/:id |
| mall:store:create | 创建门店 | POST | /api/admin/mall/stores |
| mall:store:update | 更新门店 | PUT | /api/admin/mall/stores/:id |
| mall:store:delete | 删除门店 | DELETE | /api/admin/mall/stores/:id |
| mall:store_product:bind | 绑定门店商品 | POST | /api/admin/mall/store-products |
| mall:store_product:unbind | 解绑门店商品 | DELETE | /api/admin/mall/store-products |
| mall:store_product:status | 上下架门店商品 | PUT | /api/admin/mall/store-products/status |
| mall:store_product:stock | 门店商品库存 | PUT | /api/admin/mall/store-products/stock |
| mall:store_product:list | 查询门店商品 | GET | /api/admin/mall/store-products/:store_id |
| mall:order:list | 查询订单 | GET | /api/admin/mall/orders |
| mall:order:ship | 订单发货 | POST | /api/admin/mall/orders/:id/ship |
| mall:order:cancel | 订单作废 | POST | /api/admin/mall/orders/:id/cancel |

community-service 已定义并接入以下权限点：

| Permission Code | 说明 | Method | Path |
|---|---|---|---|
| community:notice:list | 查询公告管理列表 | GET | /api/admin/community/notices |
| community:notice:create | 发布公告 | POST | /api/admin/community/notices |
| community:notice:delete | 删除公告 | DELETE | /api/admin/community/notices/:id |
| community:notice:views | 查询公告浏览状态 | GET | /api/admin/community/notices/:id/views |
| community:visitor:list | 查询访客记录 | GET | /api/admin/community/visitors |
| community:visitor:audit | 访客审核放行 | POST | /api/admin/community/visitors/:id/audit |
| community:parking:list | 查询车位列表 | GET | /api/admin/community/parking-spaces |
| community:parking:create | 创建车位 | POST | /api/admin/community/parking-spaces |
| community:parking:assign | 分配车位 | POST | /api/admin/community/parking-spaces/:id/assign |
| community:parking:statistics | 查询车位统计 | GET | /api/admin/community/parking-spaces/statistics |
| community:fee:list | 查询物业费 | GET | /api/admin/community/property-fees |
| community:fee:create | 创建物业费 | POST | /api/admin/community/property-fees |
| community:fee:payment_list | 查询缴费记录 | GET | /api/admin/community/property-fees/payments |

workorder-service 已定义并接入以下权限点：

| Permission Code | 说明 | Method | Path |
|---|---|---|---|
| workorder:repair:list | 查询报修列表 | GET | /api/admin/workorders/repairs |
| workorder:repair:process | 处理报事维修 | POST | /api/admin/workorders/repairs/:id/process |
| workorder:complaint:list | 查询投诉列表 | GET | /api/admin/workorders/complaints |
| workorder:complaint:process | 处理事项投诉 | POST | /api/admin/workorders/complaints/:id/process |

statistics-service 已定义并接入以下权限点：

| Permission Code | 说明 | Method | Path |
|---|---|---|---|
| statistics:product:sales_rank | 商品销售排行 | GET | /api/statistics/products/sales-rank |
| statistics:product:view_rank | 商品访客排行 | GET | /api/statistics/products/view-rank |
| statistics:community:overview | 社区运营概览 | GET | /api/statistics/community/overview |
| statistics:order:summary | 订单统计 | GET | /api/statistics/orders |
| statistics:workorder:summary | 报修投诉统计 | GET | /api/statistics/workorders |

## 5. 默认角色与权限

| 角色 | Code | 说明 | 默认权限 |
|---|---|---|---|
| 系统管理员 | admin | 全部权限 | 所有 `*` |
| 物业管理员 | property | 物业与社区管理 | 社区、工单相关 |
| 门店管理员 | store | 门店与商品管理 | 商城管理相关 |
| 普通用户 | user | 居民用户 | 无管理端权限 |

admin 角色默认绑定所有权限点。其他角色按需绑定。

## 6. 与 gateway-service 的关系

- 当前 user-service 负责权限校验（JWT + RequirePermission）。
- 后续 gateway-service 可统一做 JWT 验证和权限校验，业务服务保留二次校验。
- 权限数据由 user-service 管理，gateway 通过 RPC 或缓存同步。

## 7. Legacy 兼容

- `users.role` 字段保留，`AssignRole` 旧接口同步写 `sys_user_role` 表和 `users.role` 字段。
- `RequireRole` 中间件保留，但管理端新接口逐步迁移到 `RequirePermission`。
- JWT Claims 中的 `Role` 字段保留，用于快速判断（如 admin 角色可跳过权限检查）。

## 8. 权威数据源说明

| 维度 | 权威表 | Legacy 字段 | 说明 |
|------|--------|-------------|------|
| 用户角色 | `sys_user_role` | `users.role` | `users.role` 仅做 legacy 兼容，`AssignRole`/`AssignUserRoles` 双写 |
| 角色权限 | `sys_role_permission` | 无 | 全新表，通过 `RequirePermission` 中间件校验 |
| 角色菜单 | `sys_role_menu` | 无 | 控制前端菜单可见性，不影响后端 API 访问 |
| 权限点定义 | `sys_permission` | 无 | 全新表，code 格式 `{domain}:{resource}:{action}` |

## 9. 验证记录 (2026-05-21)

### 9.1 编译与格式验证

- `go build ./...` 通过
- `gofmt -l services/ pkg/` 无输出
- `go test ./...` 通过（均为 `[no test files]`）

### 9.2 AutoMigrate 验证

`user-service/cmd/server/main.go` AutoMigrate 包含：
- `&model.SysPermission{}`
- `&model.SysRolePermission{}`

### 9.3 Seed 数据验证

`deploy/docker-compose/mysql/init/003_user_service_seed.sql` 包含：
- 4 个角色：admin / property / store / user
- 25 个菜单（层级结构，parent_id 关联）
- 17 个权限点（ON DUPLICATE KEY UPDATE 幂等）
- 角色-权限绑定：admin 17 个全量、property 4 个、store 2 个
- 角色-菜单绑定：admin 25 个全量、property 12 个、store 9 个
- admin 用户（id=1）+ sys_user_role 绑定（user_id=1, role_id=1）

### 9.4 权限校验链路验证

```
AdminHandler.GetPermissionCodesByUserID(userID)
  → AdminService.GetUserRoles(userID)
    → RoleRepo.FindRolesByUserID(userID)  [JOIN sys_user_role + sys_role]
  → AdminService.GetPermissionsByRoleIDs(roleIDs)
    → RoleRepo.FindPermissionsByRoleIDs(roleIDs)  [JOIN sys_role_permission + sys_permission]
  → 返回 []string permission codes
```

RequirePermission 中间件：
1. 从 gin.Context 获取 userID
2. admin 快捷路径：JWT role=="admin" 直接放行（临时方案，gateway 接管后移除）
3. Redis SMEMBERS `rbac:permissions:{userID}` 命中则直接校验
4. 未命中则调用 PermissionProvider 回调查询 DB → SADD 写入 Redis SET + EXPIRE 600s
5. 权限码不在列表中则返回 403

### 9.5 缓存失效验证

- 用户角色变更：`AssignUserRoles` → `invalidateUserPermissionCache` → DEL `rbac:permissions:{userID}`
- 角色权限变更：`BindRolePermissions` → `invalidateRolePermissionCache` → 短 TTL 自然过期策略

### 9.6 已知限制

1. admin 快捷路径（JWT role=="admin" 跳过 DB 查询）为临时方案，gateway 接管后移除
2. 角色权限变更不立即失效所有相关用户的缓存（10 分钟 TTL 自然过期）
3. gateway-service 尚未统一做 JWT/权限校验，当前仍由各服务在管理端路由内二次校验
