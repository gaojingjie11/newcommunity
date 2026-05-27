# API 访问控制策略

Date: 2026-05-21

## 访问级别定义

| 级别 | 说明 | 中间件 | 校验逻辑 |
|------|------|--------|----------|
| `public` | 公开接口，无需认证 | 无 | 直接访问 |
| `authenticated` | 需登录认证 | `JWTAuth` | 解析 JWT + Redis 验证 token 有效性 |
| `permission_required` | 需特定权限 | `JWTAuth` + `RequirePermission` | 认证 + RBAC 权限码校验 |

## user-service

### Public 接口

| Method | Path | 说明 |
|--------|------|------|
| POST | `/api/users/register` | 用户注册 |
| POST | `/api/users/login` | 用户登录 |
| POST | `/api/users/password-reset/code` | 发送重置验证码 |
| POST | `/api/users/password-reset` | 重置密码 |

### Authenticated 接口

| Method | Path | 说明 |
|--------|------|------|
| POST | `/api/users/logout` | 退出登录 |
| GET | `/api/users/me` | 查看个人资料 |
| PUT | `/api/users/me` | 修改个人资料 |
| PUT | `/api/users/me/password` | 修改密码 |

### Permission Required 接口

| Method | Path | Permission Code | 说明 |
|--------|------|-----------------|------|
| POST | `/api/admin/roles` | `rbac:role:create` | 创建角色 |
| PUT | `/api/admin/roles/:id` | `rbac:role:update` | 更新角色 |
| DELETE | `/api/admin/roles/:id` | `rbac:role:delete` | 删除角色 |
| GET | `/api/admin/roles` | `rbac:role:list` | 查询角色列表 |
| POST | `/api/admin/roles/:id/menus` | `rbac:role:bind_menu` | 角色绑定菜单 |
| POST | `/api/admin/roles/:id/permissions` | `rbac:role:bind_permission` | 角色绑定权限 |
| GET | `/api/admin/roles/:id/permissions` | `rbac:role:get_permissions` | 查询角色权限 |
| GET | `/api/admin/users` | `rbac:user:list` | 查询管理员列表 |
| POST | `/api/admin/users/freeze` | `rbac:user:freeze` | 冻结/解冻用户 |
| POST | `/api/admin/users/assign-role` | `rbac:user:assign_role` | 分配用户角色(旧) |
| POST | `/api/admin/users/:id/roles` | `rbac:user:assign_roles` | 分配用户多角色 |
| GET | `/api/admin/users/:id/roles` | `rbac:user:get_roles` | 查询用户角色 |
| GET | `/api/admin/members` | `rbac:member:list` | 查询会员列表 |
| GET | `/api/admin/permissions` | `rbac:permission:list` | 查询权限列表 |
| GET | `/api/admin/menus` | `rbac:menu:list` | 查询菜单列表 |
| GET | `/api/admin/user-login-logs` | `log:user_login:list` | 查询用户登录日志 |
| GET | `/api/admin/admin-login-logs` | `log:admin_login:list` | 查询管理员登录日志 |

---

## mall-service

### Public 接口

| Method | Path | 说明 |
|--------|------|------|
| GET | `/api/mall/products` | 商品列表（MALL-001/004） |
| GET | `/api/mall/products/search` | 商品搜索（MALL-002） |
| GET | `/api/mall/products/promotions` | 促销商品（MALL-003） |
| GET | `/api/mall/stores` | 门店列表（MALL-006） |
| GET | `/api/mall/stores/:id` | 门店详情 |
| GET | `/api/mall/promotions` | 促销列表 |
| GET | `/api/mall/promotions/:id` | 促销详情 |
| GET | `/api/mall/categories` | 分类列表 |
| GET | `/api/mall/categories/:id` | 分类详情 |
| GET | `/api/mall/service-areas` | 服务区域列表 |

### Authenticated 接口

| Method | Path | 说明 |
|--------|------|------|
| POST | `/api/mall/cart/items` | 添加购物车（MALL-007） |
| DELETE | `/api/mall/cart/items/:id` | 移除购物车（MALL-010） |
| PUT | `/api/mall/cart/items/:id` | 更新数量（MALL-011） |
| GET | `/api/mall/cart/items` | 购物车列表 |
| POST | `/api/mall/favorites` | 收藏商品（MALL-008） |
| DELETE | `/api/mall/favorites/:product_id` | 取消收藏（MALL-009） |
| GET | `/api/mall/favorites` | 我的收藏（MALL-018） |
| GET | `/api/mall/favorites/check/:product_id` | 检查收藏状态 |
| GET | `/api/mall/products/:id` | 商品详情（MALL-005，进入详情需登录） |
| POST | `/api/mall/orders` | 创建订单（MALL-012） |
| POST | `/api/mall/orders/:id/pay` | 订单支付（MALL-013） |
| GET | `/api/mall/orders/:id` | 订单详情 |
| GET | `/api/mall/orders` | 订单列表（MALL-014~017） |
| POST | `/api/mall/wallet/recharge` | 充值（MALL-019） |
| POST | `/api/mall/wallet/transfer` | 转账（MALL-020） |
| GET | `/api/mall/wallet/balance` | 查询余额（MALL-021） |
| GET | `/api/mall/wallet/transactions` | 交易记录（MALL-021） |

### Permission Required 接口（已迁移至 RequirePermission）

| Method | Path | Permission Code | 说明 |
|--------|------|-----------------|------|
| POST | `/api/admin/mall/categories` | `mall:category:create` | 创建分类（ADMIN-MALL-005） |
| PUT | `/api/admin/mall/categories/:id` | `mall:category:update` | 更新分类 |
| DELETE | `/api/admin/mall/categories/:id` | `mall:category:delete` | 删除分类 |
| POST | `/api/admin/mall/products` | `mall:product:create` | 创建商品（ADMIN-MALL-006） |
| PUT | `/api/admin/mall/products/:id` | `mall:product:update` | 更新商品 |
| DELETE | `/api/admin/mall/products/:id` | `mall:product:delete` | 删除商品 |
| POST | `/api/admin/mall/promotions` | `mall:promotion:create` | 创建促销（ADMIN-MALL-007） |
| PUT | `/api/admin/mall/promotions/:id` | `mall:promotion:update` | 更新促销 |
| DELETE | `/api/admin/mall/promotions/:id` | `mall:promotion:delete` | 删除促销 |
| POST | `/api/admin/mall/promotions/:id/products` | `mall:promotion:bind_product` | 绑定促销商品 |
| POST | `/api/admin/mall/service-areas` | `mall:service_area:create` | 创建服务区域（ADMIN-MALL-008） |
| PUT | `/api/admin/mall/service-areas/:id` | `mall:service_area:update` | 更新服务区域 |
| DELETE | `/api/admin/mall/service-areas/:id` | `mall:service_area:delete` | 删除服务区域 |
| POST | `/api/admin/mall/stores` | `mall:store:create` | 创建门店（ADMIN-MALL-009） |
| PUT | `/api/admin/mall/stores/:id` | `mall:store:update` | 更新门店 |
| DELETE | `/api/admin/mall/stores/:id` | `mall:store:delete` | 删除门店 |
| POST | `/api/admin/mall/store-products` | `mall:store_product:bind` | 绑定门店商品（ADMIN-MALL-010） |
| DELETE | `/api/admin/mall/store-products` | `mall:store_product:unbind` | 解绑门店商品 |
| PUT | `/api/admin/mall/store-products/status` | `mall:store_product:status` | 上下架门店商品 |
| PUT | `/api/admin/mall/store-products/stock` | `mall:store_product:stock` | 门店商品库存 |
| GET | `/api/admin/mall/store-products/:store_id` | `mall:store_product:list` | 查询门店商品 |
| GET | `/api/admin/mall/orders` | `mall:order:list` | 订单列表（ADMIN-MALL-011） |
| POST | `/api/admin/mall/orders/:id/ship` | `mall:order:ship` | 订单发货 |
| POST | `/api/admin/mall/orders/:id/cancel` | `mall:order:cancel` | 订单作废 |

---

## community-service

### Public 接口

| Method | Path | 说明 |
|--------|------|------|
| GET | `/api/community/ping` | 服务连通性检查 |
| GET | `/api/community/notices` | 公告列表（COMM-001） |
| GET | `/api/community/notices/:id` | 公告详情与浏览计数（COMM-001） |

### Authenticated 接口

| Method | Path | 说明 |
|--------|------|------|
| POST | `/api/community/notices/:id/read` | 标记公告已读 |
| POST | `/api/community/visitors` | 访客登记（COMM-002） |
| GET | `/api/community/visitors` | 我的访客登记记录 |
| GET | `/api/community/parking-spaces/my` | 我的车位（COMM-003） |
| PUT | `/api/community/parking-spaces/:id/plate` | 绑定车牌（COMM-004） |
| GET | `/api/community/property-fees` | 我的物业费账单 |
| POST | `/api/community/property-fees/:id/pay` | 缴纳物业费（COMM-007） |
| GET | `/api/community/property-fees/payments` | 我的缴费记录 |

### Permission Required 接口（已迁移至 RequirePermission）

| Method | Path | Permission Code | 说明 |
|--------|------|-------------|------|
| GET | `/api/admin/community/notices` | `community:notice:list` | 管理端公告列表 |
| POST | `/api/admin/community/notices` | `community:notice:create` | 公告发布（ADMIN-COMM-001） |
| DELETE | `/api/admin/community/notices/:id` | `community:notice:delete` | 删除公告 |
| GET | `/api/admin/community/notices/:id/views` | `community:notice:views` | 公告浏览状态（ADMIN-COMM-002） |
| GET | `/api/admin/community/visitors` | `community:visitor:list` | 访客列表 |
| POST | `/api/admin/community/visitors/:id/audit` | `community:visitor:audit` | 访客审核（ADMIN-COMM-003） |
| GET | `/api/admin/community/parking-spaces` | `community:parking:list` | 车位列表 |
| POST | `/api/admin/community/parking-spaces` | `community:parking:create` | 创建车位 |
| POST | `/api/admin/community/parking-spaces/:id/assign` | `community:parking:assign` | 分配车位 |
| GET | `/api/admin/community/parking-spaces/statistics` | `community:parking:statistics` | 车位统计（ADMIN-COMM-004） |
| GET | `/api/admin/community/property-fees` | `community:fee:list` | 物业费列表 |
| POST | `/api/admin/community/property-fees` | `community:fee:create` | 创建物业费 |
| GET | `/api/admin/community/property-fees/payments` | `community:fee:payment_list` | 缴费记录（ADMIN-COMM-007） |

---

## workorder-service

### Public 接口

| Method | Path | 说明 |
|--------|------|------|
| GET | `/api/workorders/ping` | 服务连通性检查 |

### Authenticated 接口

| Method | Path | 说明 |
|--------|------|------|
| POST | `/api/workorders/repairs` | 提交报修（COMM-005） |
| GET | `/api/workorders/repairs` | 我的报修记录 |
| POST | `/api/workorders/complaints` | 提交投诉（COMM-006） |
| GET | `/api/workorders/complaints` | 我的投诉记录 |
| GET | `/api/workorders/:type/:id/logs` | 报修/投诉状态流转日志 |

### Permission Required 接口（已迁移至 RequirePermission）

| Method | Path | Permission Code | 说明 |
|--------|------|-------------|------|
| GET | `/api/admin/workorders/repairs` | `workorder:repair:list` | 报修列表 |
| POST | `/api/admin/workorders/repairs/:id/process` | `workorder:repair:process` | 处理报修（ADMIN-COMM-005） |
| GET | `/api/admin/workorders/complaints` | `workorder:complaint:list` | 投诉列表（ADMIN-COMM-006） |
| POST | `/api/admin/workorders/complaints/:id/process` | `workorder:complaint:process` | 处理投诉 |

---

## statistics-service

### Permission Required 接口（已接入 RequirePermission）

| Method | Path | Permission Code | 说明 |
|--------|------|-----------------|------|
| GET | `/api/statistics/products/sales-rank` | `statistics:product:sales_rank` | 商品销售排行（STAT-001） |
| GET | `/api/statistics/products/view-rank` | `statistics:product:view_rank` | 商品访客排行（STAT-002） |
| GET | `/api/statistics/community/overview` | `statistics:community:overview` | 社区运营概览（STAT-003） |
| GET | `/api/statistics/orders` | `statistics:order:summary` | 订单统计（STAT-004） |
| GET | `/api/statistics/workorders` | `statistics:workorder:summary` | 报修投诉统计（STAT-005） |

---

## gateway-service

### Public 接口

| Method | Path | 说明 |
|--------|------|------|
| GET | `/health` | 健康检查 |
| GET | `/api/gateway/services` | 服务列表（GATEWAY-002） |

### Gateway 统一鉴权说明

gateway-service 作为统一 API 入口，实现三级访问控制：

1. **Public**：无需 JWT，直接代理到下游服务
2. **Authenticated**：需要 JWT，但不需要特定权限码
3. **Permission Required**：需要 JWT + RBAC 权限码校验

#### 访问级别路由

| 级别 | 路径模式 | 说明 |
|------|----------|------|
| `public` | `/api/users/register`, `/api/users/login`, `/api/mall/products`, `/api/mall/products/search`, `/api/mall/products/promotions`, `/api/community/ping`, `/api/community/notices`, `/api/workorders/ping`, `/agent/health` | 无需认证 |
| `authenticated` | `/api/users/me`, `/api/mall/products/:id`, `/api/mall/cart/**`, `/api/mall/orders/**`, `/api/community/visitors`, `/api/workorders/repairs` | 需要 JWT |
| `permission_required` | `/api/admin/**`, `/api/statistics/**` | 需要 JWT + 权限码 |

#### Gateway 中间件链

```
请求 → CORS → RequestID → Logger → Recovery
  ├─ public 路径 → 直接代理
  ├─ authenticated 路径 → JWTAuth → 代理
  └─ permission_required 路径 → JWTAuth → RequirePermission → 代理
```

#### 注入的请求头

Gateway 代理时自动注入以下请求头：

| Header | 说明 |
|--------|------|
| `X-Request-ID` | 请求追踪 ID |
| `X-Gateway-Proxy` | 标识来源为 gateway |
| `X-Gateway-Time` | 代理时间 |
| `X-Internal-Token` | 服务间调用令牌 |
| `X-User-ID` | 用户 ID（来自 JWT） |
| `X-User-Role` | 用户角色（来自 JWT） |

#### 服务发现

Gateway 支持两种服务发现方式：

1. **优先**：从 Nacos 获取服务实例地址
2. **兜底**：如果 Nacos 不可用，从 `configs/config.yaml` 的 `gateway.services` 读取地址

`/api/gateway/services` 接口返回当前可用服务地址及来源（`nacos` 或 `config`）。

---

## 后续迁移计划

后续服务迁移到 `RequirePermission` 时，需：

1. 在 `sys_permission` 表新增业务权限码
2. 在 seed SQL 中绑定 admin/业务管理员角色权限
3. 修改 `router.go`：`RequireRole("admin")` → `RequirePermission(rdb, provider, code)`
4. 服务 handler 实现 `PermissionProvider` 接口或复用同构 repository 查询
5. 更新本文档：将 "当前守卫" 列改为 `RequirePermission`
