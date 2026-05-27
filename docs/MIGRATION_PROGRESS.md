# Migration Progress

Date: 2026-05-21

## Functional Baseline Rule

- From this stage forward, `需求陈述书-东软智慧社区项目.docx` is the functional baseline and highest source of truth.
- The legacy `smartcomunity` code is migration reference only. If legacy behavior conflicts with the Word requirements, implement the Word requirement.
- All future business migration must be traced through `docs/REQUIREMENTS_TRACEABILITY_MATRIX.md`.
- If a feature has no traceability matrix ID, it must not be directly developed. Add or update the requirement record first.

## Completed In Requirements Baseline Stage

- Added `docs/REQUIREMENTS_BASELINE.md`.
- Added `docs/REQUIREMENTS_TRACEABILITY_MATRIX.md`.
- Added `docs/API_CONTRACT_BY_REQUIREMENTS.md`.
- Added `docs/MIGRATION_GAP_ANALYSIS.md`.
- Added `docs/DATA_MODEL_PLAN.md`.
- Scanned Word requirements and compared them with legacy controller/service/model/router coverage.
- Marked implementation states as `TODO`, `PARTIAL`, `DONE`, `GAP`, or `EXTENSION`.

## Completed In This Stage

- Scanned the existing project and preserved `frontend` plus `smartcomunity` as legacy reference code.
- Added root microservice module and shared Go packages under `pkg/`.
- Added service skeletons:
  - `gateway-service`
  - `user-service`
  - `mall-service`
  - `community-service`
  - `workorder-service`
  - `statistics-service`
  - `agent-service`
- Added health checks and placeholder APIs for each service.
- Added Nacos best-effort service registration.
- Added RabbitMQ event publish placeholder for `repair.created` and `complaint.created`.
- Added FastAPI `agent-service` with four placeholder Agent endpoints.
- Added Dockerfiles for all new services.
- Added Docker Compose infrastructure and service orchestration under `deploy/docker-compose`.
- Added Kubernetes skeleton under `deploy/k8s`.
- Added a lightweight frontend `/agent` placeholder page and service entry.
- Added migration, architecture, infrastructure, gateway, agent and deployment docs.

## Verification Performed

- `go mod tidy` completed and generated `go.sum`.
- `go build ./...` passed for all new Go packages and services.
- `go test ./...` passed for all new Go packages and services.
- `python3 -m py_compile services/agent-service/app/*.py` passed with a writable pycache path.
- `docker compose -f deploy/docker-compose/docker-compose.yml config --quiet` passed.
- K8s YAML files passed local YAML parsing. `kubectl apply --dry-run=client` could not complete because this machine has no reachable Kubernetes API server and `kubectl` attempted discovery against `localhost:8080`.
- `gateway-service` was started locally and verified:
  - `GET /health`
  - `GET /api/gateway/services`

## Completed In User-Service Migration Stage

- Implemented `pkg/auth/password.go` with `HashPassword` and `CheckPasswordHash` (bcrypt cost=14).
- Implemented `pkg/middleware/jwt_auth.go` — JWTAuth middleware with Redis single-session validation.
- Implemented `pkg/middleware/require_role.go` — RequireRole middleware for role-based access control.
- Implemented user-service data models: `SysUser`, `SysRole`, `SysMenu`, `SysUserRole`, `SysRoleMenu`, `UserLoginLog`, `AdminLoginLog`, `PasswordResetCode`.
- Implemented repository layer: `UserRepo`, `RoleRepo`, `LoginLogRepo`, `PasswordResetRepo`.
- Implemented service layer: `AuthService`, `UserService`, `AdminService`, `LoginLogService`.
- Implemented handler layer: `AuthHandler`, `UserHandler`, `AdminHandler`.
- Implemented router with public, authenticated, and admin route groups.
- Rewrote `cmd/server/main.go` with full DI wiring: config → MySQL → AutoMigrate → Redis → Nacos → repos → services → handlers → router.
- Added JWT config section to `configs/config.yaml`.
- Added SQL init scripts: `002_user_service_tables.sql`, `003_user_service_seed.sql`.
- Requirements covered: AUTH-001~007, LOG-001~002, ADMIN-MALL-001~004 (13 IDs total).
- All `go build ./...` passed.

## Completed In Mall-Service Migration Stage

- Implemented mall-service data models (10 files): `Product`, `ProductCategory`, `Promotion`, `PromotionProduct`, `Cart`, `Order`, `OrderItem`, `Store`, `StoreProduct`, `Favorite`, `Wallet`, `WalletTransaction`, `ServiceArea`, `SysUser` (read-only ref).
- Implemented repository layer (10 files): `ProductRepo`, `CategoryRepo`, `PromotionRepo`, `CartRepo`, `OrderRepo`, `StoreRepo`, `StoreProductRepo`, `FavoriteRepo`, `WalletRepo`, `ServiceAreaRepo`.
- Implemented service layer (9 files): `ProductService`, `CartService`, `OrderService`, `StoreService`, `FavoriteService`, `WalletService`, `PromotionService`, `CategoryService`, `ServiceAreaService`.
- Implemented handler layer (10 files): `ProductHandler`, `CartHandler`, `OrderHandler`, `StoreHandler`, `FavoriteHandler`, `WalletHandler`, `PromotionHandler`, `CategoryHandler`, `ServiceAreaHandler`, `AdminOrderHandler`.
- Implemented router with public, authenticated user, and admin route groups.
- Rewrote `cmd/server/main.go` with full DI wiring: config → MySQL → AutoMigrate(13 tables) → Redis → Nacos → repos → services → handlers → router.
- Added JWT config section to `configs/config.yaml`.
- Added SQL init scripts: `004_mall_service_tables.sql`, `005_mall_service_seed.sql`.
- Key business logic: atomic stock deduction (WHERE stock>=qty), wallet payment with SELECT FOR UPDATE, transfer with ordered locking to prevent deadlock, promotion M2M binding.
- Requirements covered: MALL-001~021, ADMIN-MALL-005~011 (28 IDs total).
- All `go build ./...` passed.

## Completed In RBAC Enhancement Stage

- Created `docs/RBAC_DESIGN.md` — complete RBAC architecture: users, roles, menus, permissions, user_roles, role_menus, role_permissions.
- Updated `docs/DATA_MODEL_PLAN.md` — added `sys_permission` and `sys_role_permission` table definitions.
- Added `SysPermission` and `SysRolePermission` models in `services/user-service/internal/model/permission.go`.
- Updated `SysUser.Role` field comment: legacy compat, not primary RBAC authority.
- Extended `RoleRepo` with: `BindPermissions`, `FindPermissionsByRoleID`, `FindPermissionsByRoleIDs`, `ListAllPermissions`, `FindRolesByUserID`, `BindUserRoles`, `FindRoleCodesByUserID`. Role deletion now cascades to role_menus, role_permissions, user_roles.
- Extended `AdminService` with: `BindRolePermissions`, `GetRolePermissions`, `GetPermissionsByRoleIDs`, `ListAllPermissions`, `AssignUserRoles`, `GetUserRoles`, `ListAllMenus`. `AssignRole` now writes both `users.role` and `sys_user_role`. Added Redis cache invalidation on role/permission changes. `AdminService` now requires `*goredis.Client`.
- Implemented `pkg/middleware/require_permission.go` — `RequirePermission` middleware with `PermissionProvider` interface, Redis SET cache (`rbac:permissions:{userID}`, 10-min TTL), admin role shortcut.
- Updated router: all admin endpoints now use `RequirePermission` instead of `RequireRole("admin")`. Added new endpoints: `POST/GET /api/admin/roles/:id/permissions`, `POST/GET /api/admin/users/:id/roles`, `GET /api/admin/permissions`, `GET /api/admin/menus`.
- 17 permission codes defined covering RBAC management and login log queries.
- Updated SQL seed: added `store` role (4 roles total), 25 menus (hierarchical), 17 permissions, role-permission bindings, role-menu bindings. Uses `ON DUPLICATE KEY UPDATE` for idempotency.
- Bug fixes:
  - `password_reset_codes.code_hash` now stores bcrypt hash instead of plaintext.
  - `ResetPassword` now calls `MarkUsedByMobile` to write `used_at`.
  - Login failure now writes to `admin_login_logs` when user role is `admin`.
  - Registration validates `age > 0` and `gender in {0, 1, 2}` (explicit check, not binding required).
  - Added `PasswordResetRepo.MarkUsedByMobile` method.
- Updated `cmd/server/main.go`: AutoMigrate includes `SysPermission` and `SysRolePermission`.
- All `go build ./...` passed.

## Known RBAC Limitations (Future Work)

- `RequirePermission` admin shortcut (JWT role=="admin") is acceptable short-term but should be removed when gateway handles auth.
- Role permission cache invalidation is TTL-based (10 min); no fan-out to invalidate all affected users immediately.
- Other services (mall, community, workorder) admin endpoints still use `RequireRole("admin")`; should migrate to `RequirePermission` with service-specific permission codes.
- `users.role` field still exists for legacy compat; can be removed once all consumers use `sys_user_role`.
- Gateway-service not yet implementing unified JWT/permission validation.

## RBAC Verification Completed (2026-05-21)

- Updated `docs/RBAC_DESIGN.md` with verification section (§8-§9): authority data sources, permission chain, cache invalidation, known limitations.
- Created `docs/API_ACCESS_CONTROL_POLICY.md` — three-tier access control table for all services.
- Verified: `go build ./...` passes, `go test ./...` passes, `gofmt` clean.
- Verified: AutoMigrate includes SysPermission + SysRolePermission.
- Verified: Seed SQL has 4 roles, 25 menus, 17 permissions, role-permission/role-menu bindings, admin user-role binding.
- Verified: Permission chain `AdminHandler → GetUserRoles → GetPermissionsByRoleIDs → FindPermissionsByRoleIDs` is correct.
- Verified: Redis cache pattern (SMEMBERS/SADD/EXPIRE) and user-level DEL on role change.
- RBAC for user-service: DONE.
- RBAC for mall-service: DONE (migrated RequireRole("admin") → RequirePermission, 24 permission codes).
- RBAC for community-service: DONE (migrated admin community APIs to RequirePermission, 13 permission codes).
- RBAC for workorder-service: DONE (migrated admin workorder APIs to RequirePermission, 4 permission codes).
- RBAC for gateway-service: TODO (unified JWT/permission validation).

## Completed In Payment Consistency Overhaul (2026-05-21)

- **Constants**: Created `internal/consts/order_status.go` — order states (0/1/2/3/40), wallet tx types, biz types, payment statuses, 30-min expiry constant.
- **Models**: Converted all monetary fields from float64 to int64 (cents): `Order.TotalAmount`, `Order.UsedBalance`, `Product.Price`, `Product.OriginalPrice`, `Wallet.Balance`, `WalletTransaction.Amount`, `OrderItem.Price`.
- **Models**: Added `Order.ExpireAt`, `Order.CancelReason`, `Order.CancelledAt`, `Order.PaidAt`, `Order.Version`, `Order.UpdatedAt`, `OrderItem.ProductSnapshot`, `Product.Version`, `Wallet.Version`, `StoreProduct.LockedStock`, `StoreProduct.SoldCount`, `StoreProduct.Version`, `WalletTransaction.BalanceBefore/After/BizType/BizID/IdempotencyKey`.
- **New Model**: `PaymentRecord` — idempotent payment audit trail with `idempotency_key` UNIQUE constraint.
- **Repository**: Added `WithTx(tx)` pattern to all repos for transaction-scoped instances.
- **Repository**: `OrderRepo` — conditional `UpdateStatus(from→to)`, `MarkAsPaid`, `MarkAsCancelled`, `FindByIDForUpdate`, `FindExpiredPendingOrders`.
- **Repository**: `ProductRepo` — `DeductStock(tx, id, qty)` with `WHERE stock>=qty`, `RestoreStock(tx, id, qty)`.
- **Repository**: `WalletRepo` — `Debit/Credit/CreateTransaction` now take tx parameter; all amounts int64.
- **New Repository**: `PaymentRepo` — `Create`, `FindByIdempotencyKey`, `FindByOrderID`, `UpdateStatus`.
- **Service**: `OrderService` — fixed transaction leaks (all mutations use WithTx(tx)); `CancelOrder` restores stock and refunds wallet for paid orders; removed `PayOrder` (moved to `PaymentService`).
- **New Service**: `PaymentService` — idempotent `PayOrder` with SELECT FOR UPDATE on order+wallet, atomic debit, conditional MarkAsPaid, best-effort event publish.
- **New Service**: `OrderTimeoutService` — polling every 60s for expired orders, conditional cancel with stock restore.
- **New Service**: `EventBus` — wraps `rabbitmq.Client` for best-effort event publishing (order.created, order.paid, order.cancelled).
- **Service**: `WalletService` — rewritten for int64 amounts, `BalanceBefore/After` tracking.
- **Handler**: `OrderHandler` — removed `Pay` (moved to `PaymentHandler`), added `Cancel` with ownership check.
- **New Handler**: `PaymentHandler` — `Pay` (requires idempotency_key), `GetPaymentStatus`.
- **Handler**: `WalletHandler` — amount type float64→int64.
- **Handler**: `AdminOrderHandler.CancelOrder` — now accepts reason parameter.
- **Router**: Added `PaymentHandler`, new routes: `POST /orders/:id/cancel`, `GET /orders/:id/payment-status`.
- **Main**: Wired `EventBus`, `PaymentRepo`, `PaymentService`, `PaymentHandler`, `OrderTimeoutService`; AutoMigrate includes `PaymentRecord`.
- **Migration**: `migrations/001_order_payment_overhaul.sql` — ALTER statements for float64→bigint, new columns, new `payment_records` table.
- **Tests**: Basic unit tests for constants, request validation, status transitions.
- **Docs**: `docs/PAYMENT_CONSISTENCY_DESIGN.md` — complete design document.
- All `go build ./...` passed, `go test ./...` passed, `gofmt` clean.

## Completed In Community-Service Migration Stage (2026-05-21)

- Implemented community-service data models: `Notice`, `NoticeViewLog`, `Visitor`, `ParkingSpace`, `UserParkingBinding`, `PropertyFee`, `PropertyFeePayment`.
- Implemented repository layer: notices with browse/read logs, visitors with audit flow, parking spaces with binding/statistics, property fees with idempotent payment records, shared RBAC permission lookup.
- Implemented service layer: announcement list/detail/create/delete/read, visitor registration/audit, parking list/create/assign/bind plate/statistics, property fee create/list/pay/payment records.
- Implemented handler and router layers with public, authenticated user, and `RequirePermission` admin route groups.
- Rewrote `community-service/cmd/server/main.go` with full DI wiring: config → MySQL → AutoMigrate → Redis → Nacos → repos → services → handlers → router.
- Added JWT config section to `services/community-service/configs/config.yaml` and `.example`.
- Added Docker Compose init script `006_community_service_tables.sql` with community tables and seed data.
- Updated `003_user_service_seed.sql` with 13 community permission codes and role bindings for `admin` and `property`.
- Requirements now covered in community-service: `COMM-001~004`, `ADMIN-COMM-001~004`, `ADMIN-COMM-007`. `COMM-007` is PARTIAL because wallet deduction is not yet wired.
- Updated docs: `REQUIREMENTS_TRACEABILITY_MATRIX.md`, `API_ACCESS_CONTROL_POLICY.md`, `RBAC_DESIGN.md`, `DATA_MODEL_PLAN.md`, and community-service README.
- Verification: `go test ./...` passed.

## Completed In Workorder-Service Migration Stage (2026-05-21)

- Implemented workorder-service data models: `Repair`, `Complaint`, `WorkorderLog`.
- Split legacy `cms_repair type=1/2` into Word-aligned `repairs` and `complaints` tables.
- Implemented repository/service/handler/router layers for repair submit/list/process, complaint submit/list/process, and status logs.
- Replaced placeholder submit endpoints with authenticated user APIs and real persistence.
- Preserved RabbitMQ event publishing for `repair.created` and `complaint.created`; RabbitMQ remains best-effort and service still accepts requests when unavailable.
- Added `RequirePermission` RBAC on management APIs with `workorder:repair:list`, `workorder:repair:process`, `workorder:complaint:list`, `workorder:complaint:process`.
- Rewrote `workorder-service/cmd/server/main.go` with config → MySQL → AutoMigrate → Redis → Nacos → RabbitMQ → DI → router.
- Added JWT config section to `services/workorder-service/configs/config.yaml` and `.example`.
- Added Docker Compose init script `007_workorder_service_tables.sql`.
- Updated `003_user_service_seed.sql` with workorder permission codes and bindings for `admin` and `property`.
- Requirements now covered in workorder-service: `COMM-005`, `COMM-006`, `ADMIN-COMM-005`, `ADMIN-COMM-006`.

## Completed In Gateway Enhancement Stage (2026-05-21)

- **Service Discovery**: Created `internal/discovery/discovery.go` — Nacos service discovery with local config fallback. Refreshes every 30s. Nacos unavailability does not block startup.
- **Permission Provider**: Created `internal/perm/provider.go` — queries `sys_user_role` + `sys_role_permission` + `sys_permission` directly from MySQL.
- **Permission Mapping**: Created `internal/perm/mapping.go` — maps HTTP method + path to permission codes for all admin and statistics routes (60+ mappings covering user/mall/community/workorder/statistics services).
- **Models**: Created `internal/model/permission.go` — minimal `SysUserRole`, `SysRolePermission`, `SysPermission` models for gateway permission queries.
- **Main**: Rewrote `cmd/server/main.go` with:
  - MySQL + Redis connections
  - JWT validation via `pkg/middleware.JWTAuth`
  - RBAC permission checking via `gatewayRequirePermission` middleware
  - Three-tier route groups: public, authenticated, permission_required
  - Service discovery-based proxy with header injection (X-User-ID, X-User-Role, X-Request-ID, X-Internal-Token)
- **Config**: Updated `configs/config.yaml` and `.example` with `jwt.secret` and `jwt.ttl` sections.
- **Route Mapping**:
  - Public: `/api/users/register`, `/api/users/login`, mall list/search/promotion APIs, `/api/community/ping`, `/api/community/notices`, `/api/workorders/ping`, `/agent/health`
  - Authenticated: `/api/users/me`, `/api/mall/products/:id`, `/api/mall/cart/**`, `/api/mall/orders/**`, `/api/community/visitors`, `/api/workorders/repairs`, etc.
  - Permission Required: `/api/admin/**`, `/api/statistics/**`
- Gateway now serves as the unified API entry point. Downstream services retain their own JWTAuth and RequirePermission for二次校验.
- Updated docs: `API_ACCESS_CONTROL_POLICY.md`, `API_GATEWAY_PLAN.md`.
- All `go build ./...` passed.

## Completed In Statistics-Service Migration Stage (2026-05-21)

- **Models**: Created `internal/model/models.go` — read-only models for cross-table aggregation: `Product`, `Order`, `OrderItem`, `PropertyFee`, `Repair`, `Complaint`, `SysUser`, plus aggregation result structs (`ProductSalesRank`, `OrderSummary`, `OrderTrend`, `WorkorderSummary`, `CommunityOverview`).
- **Repository**: Created `internal/repository/stats_repo.go` — SQL aggregation queries:
  - `ProductSalesRank`: joins `oms_order_item` + `oms_order` + `pms_product`, groups by product, orders by total amount.
  - `CommunityOverview`: counts from `sys_user`, `oms_order`, `repairs`, `complaints`, `property_fees`.
  - `OrderSummary`: groups orders by status with amount totals.
  - `OrderTrend`: daily order counts for last N days.
  - `WorkorderSummary`: groups repairs and complaints by type and status.
- **Service**: Created `internal/service/stats_service.go` — business logic with input validation (limit/days bounds).
- **Handler**: Created `internal/handler/stats_handler.go` — 5 endpoints:
  - `GET /api/statistics/products/sales-rank` (STAT-001)
  - `GET /api/statistics/products/view-rank` (STAT-002)
  - `GET /api/statistics/community/overview` (STAT-003)
  - `GET /api/statistics/orders` (STAT-004)
  - `GET /api/statistics/workorders` (STAT-005)
- **Router**: Created `internal/router/router.go` — all endpoints require JWT + `RequirePermission`.
- **Main**: Rewrote `cmd/server/main.go` with DI wiring: config → MySQL → Redis → Nacos → permProvider → StatsRepo → StatsService → StatsHandler → router.
- **Permission Provider**: `permProvider` struct queries `sys_user_role` + `sys_role_permission` + `sys_permission` directly.
- **Seed SQL**: Updated `003_user_service_seed.sql` with 5 statistics permission codes (IDs 80-84) and admin role bindings.
- **Gateway**: Statistics permission mapping was already correct in `internal/perm/mapping.go`.
- Updated docs: `API_ACCESS_CONTROL_POLICY.md` (statistics section marked as implemented).
- Requirements covered: `STAT-001`, `STAT-003`, `STAT-004`, `STAT-005`. `STAT-002` is GAP (product_view_logs table not implemented).
- All `go build ./...` passed.

## Completed In STAT-002 Product View Rank Stage (2026-05-21)

- **mall-service model**: Created `internal/model/view_log.go` — `ProductViewLog` model mapping to `product_view_logs` table (product_id, user_id, ip, user_agent, viewed_at).
- **mall-service repository**: Created `internal/repository/view_log_repo.go` — `ViewLogRepo` with `Create` method.
- **mall-service handler**: Modified `internal/handler/product_handler.go` — `GetDetail` now accepts `ViewLogRepo`, records view log asynchronously via goroutine after returning response. Supports anonymous users (user_id=0 from c.ClientIP).
- **mall-service main.go**: Added `ProductViewLog` to AutoMigrate, wired `ViewLogRepo` into `ProductHandler` DI.
- **SQL DDL**: Added `product_view_logs` table to `004_mall_service_tables.sql` with indexes on product_id, user_id, viewed_at.
- **statistics-service model**: Added `ProductViewLog` read-only model and `ProductViewRank` result struct (product_id, product_name, view_count, unique_users) to `internal/model/models.go`.
- **statistics-service repository**: Added `ProductViewRank` method to `internal/repository/stats_repo.go` — joins `product_view_logs` + `pms_product`, groups by product, counts total views and distinct logged-in users.
- **statistics-service service**: Added `ProductViewRank` method with limit validation (default 20, max 100).
- **statistics-service handler**: Replaced stub `ProductViewRank` with real implementation calling service layer.
- Updated docs: `REQUIREMENTS_TRACEABILITY_MATRIX.md` (STAT-002: GAP→DONE), `DATA_MODEL_PLAN.md` (product_view_logs marked implemented), `API_ACCESS_CONTROL_POLICY.md` (removed GAP note), `MIGRATION_PROGRESS.md` (this entry).

## Completed In Property Fee Wallet Integration Stage (2026-05-21)

- **mall-service consts**: Added `WalletTxTypeFee = 5` and `BizTypePropertyFee = "property_fee"` to `internal/consts/order_status.go`.
- **mall-service repository**: Added `FindTransactionByIdempotencyKey(key)` to `internal/repository/wallet_repo.go`.
- **mall-service service**: Added `DebitForExternal(userID, amount, bizType, bizID, idempotencyKey, remark)` to `internal/service/wallet_service.go` — idempotency check → transaction with SELECT FOR UPDATE → atomic debit (`WHERE balance>=amount`) → create wallet transaction. Returns `(walletTxID, balanceBefore, balanceAfter, error)`.
- **mall-service handler**: Added `DebitWallet` handler to `internal/handler/internal_handler.go` — validates request, calls `DebitForExternal`, returns wallet_transaction_id. Added `*WalletService` to `NewInternalHandler` constructor.
- **mall-service router**: Added `POST /wallet/debit` to internal routes in `internal/router/router.go`.
- **mall-service main.go**: Updated `NewInternalHandler(orderSvc, paymentSvc)` → `NewInternalHandler(orderSvc, paymentSvc, walletSvc)`.
- **community-service client**: Created `internal/client/mall_client.go` — HTTP client with `DebitWallet(userID, amount, idempotencyKey, remark)` method. Sets `X-Internal-Token` header, parses response, returns `wallet_transaction_id`.
- **community-service repository**: Modified `PropertyFeeRepo.Pay` — now accepts `walletTxID int64`, writes it to `PropertyFeePayment.WalletTransactionID`. Added `FindByID(id)`. Changed from `tx.Save` to conditional `UPDATE ... WHERE status=0` with `RowsAffected==0` check.
- **community-service service**: Modified `PropertyFeeService.Pay` — validates fee → calls mall-client debit (idempotency key = `community-fee:{feeID}`) → passes `walletTxID` to repo. Added `mallClient *client.MallClient` field.
- **community-service main.go**: Added `mall_internal` config reading from `cfg.Raw`, creates `MallClient`, passes to `NewPropertyFeeService`.
- **community-service config**: Added `mall_internal` section to `config.yaml` and `config.yaml.example` with `base_url` and `internal_token`.
- **Tests**: Created `services/mall-service/internal/service/wallet_service_test.go` — 8 tests covering constants, idempotency key pattern, atomic balance deduction, insufficient balance, idempotent unique constraint, wallet tx type uniqueness, cross-service flow documentation, concurrent property fee pay.
- **Docs**: Updated `docs/PAYMENT_CONSISTENCY_DESIGN.md` with cross-service deduction architecture, transaction strategy, idempotency design, concurrency guarantees.
- Requirements now covered: `COMM-007` upgraded from PARTIAL to DONE.
- All `go build ./...` passed, `go test ./...` passed, `gofmt` clean.

## Completed In Property Fee Payment Hardening Stage (2026-05-21)

- **mall-service config**: Added `internal_token` to `services/mall-service/configs/config.yaml` and `config.yaml.example`, so `/api/internal/mall/wallet/debit` is registered and protected by `X-Internal-Token`.
- **community-service client**: Fixed mall-service response parsing to follow the shared response contract: `code=0`, `message`, `data`. The client now also rejects non-2xx HTTP responses before parsing business data.
- **community-service service**: Property fee wallet idempotency is now service-generated and bill-scoped: `community-fee:{feeID}`. Client-provided idempotency keys are no longer used for wallet deduction, preventing duplicate wallet debits for the same property fee when callers send different keys concurrently.
- **community-service service**: Paying a property fee now fails fast if the mall-service wallet client is not configured, and refuses to mark the fee as paid when no wallet transaction ID is returned.
- **mall-service service**: Hardened `DebitForExternal` idempotency. Existing wallet transactions are only reused when `user_id`, `amount`, `biz_type`, and `biz_id` match; conflicting reuse of an idempotency key is rejected. Concurrent duplicate inserts now re-check the existing transaction and return it when it matches.
- **Tests**: Added tests for mall-client shared response parsing, property-fee wallet key format, and external wallet debit idempotency conflict validation.
- Verification: `go test ./...` passed; `docker compose -f deploy/docker-compose/docker-compose.yml config --quiet` passed.

## Completed In Redis Cache Stage (2026-05-22)

- **StatsService**: Injected `*goredis.Client` and `*slog.Logger` into `StatsService` struct. Updated `NewStatsService` constructor to accept Redis client and logger.
- **Cache helpers**: Added `getJSONCache(ctx, key, dest)` and `setJSONCache(ctx, key, value, ttl)` internal methods with nil-safe Redis checks and JSON serialization. Cache failures log a warning and degrade to MySQL transparently.
- **Cached endpoints** (5 total):
  - `ProductSalesRank` — key `stats:product:sales-rank:{limit}`, TTL 60s
  - `ProductViewRank` — key `stats:product:view-rank:{limit}`, TTL 30s
  - `CommunityOverview` — key `stats:community:overview`, TTL 30s
  - `OrderStatsCombined` (new method) — key `stats:orders:{days}`, TTL 60s, caches summary+trend together
  - `WorkorderSummary` — key `stats:workorders:summary`, TTL 60s
- **Handler**: `OrderStats` handler now calls `OrderStatsCombined` instead of separate `OrderSummary`+`OrderTrend` calls, reducing MySQL queries per request from 2 to 1 (on cache miss).
- **Main DI**: Updated `community-service/cmd/server/main.go` to pass `rdb` and `logr` to `NewStatsService`.
- **Cache invalidation strategy**: Current stage uses short TTL auto-expiry. Cache key constants include comments documenting future event-driven invalidation points (order events, view logs, workorder events, fee payments).
- **Existing Redis features confirmed**:
  - SMS reset code: `sms:reset:{mobile}`, TTL 5min — stored in user-service auth_service.go
  - Login token: `login:token:{userID}`, TTL jwtTTL — stored on login, validated in JWT middleware, deleted on logout
  - RBAC permission cache: `rbac:permissions:{userID}`, TTL 10min — Redis SET, populated on miss, invalidated on role/permission changes
- **Tests**: Added `stats_service_test.go` with 10 tests covering cache miss, cache hit, invalid JSON handling, nil Redis safety, TTL expiry, and JSON round-trip serialization for all cached types. Uses `miniredis/v2` for in-memory Redis testing.
- **Current cache design**: Result-level caching (short TTL query result snapshots), not real-time排行榜. Deliberately avoids Redis ZSET complexity.
- **Future upgrade path**: Redis ZSET for real-time product view ranking, event-driven cache invalidation on order/workorder/fee mutations, cache warming for dashboard.
- `go test ./...` passed, `gofmt` clean.

## Completed In User/Schema Cleanup Stage (2026-05-24)

- **user-service admin user list**: `GET /api/admin/users` now lists all users instead of excluding `role=user`, matching the current unified user management screen.
- **wallet source of truth**: User list, member list, login user payload, and `/api/users/me` now read displayed balance from `wallets.balance` instead of legacy `sys_user.balance`.
- **admin balance adjustment**: `POST /api/admin/users/update-balance` now writes to `wallets` and `wallet_transactions`; `user_balance_logs` is no longer used.
- **schema cleanup**: Removed `UserBalanceLog` AutoMigrate/model usage and removed `user_balance_logs` creation from full reset SQL.
- **comment schema completion**: Added `pms_product_comment` to both incremental mall SQL and full reset SQL.
- **redundant remote tables removed**: Dropped obsolete `repairs`, `complaints`, `promotion_products`, and `user_balance_logs` from the current remote MySQL database. Unified `workorders` remains the active报修/投诉表.
- **docs cleanup**: Updated table naming from legacy `promotion_products` to active `pms_promotion_product`.
- Verification: `go test ./services/user-service/... ./services/mall-service/... ./services/community-service/... ./services/gateway-service/...` passed; rebuilt `user-service` and `/health` is healthy.

## Not Completed Yet

- Full business migration from `smartcomunity`.
- Full alignment of legacy behavior with Word requirements.
- Frontend API switch from legacy backend to `gateway-service`.
- Real Agent LLM integration.
- Production-grade K8s persistence, probes, security and scaling.
- End-to-end Docker Compose verification against all services in this environment.
- ~~Property fee payment currently marks the bill paid and records an idempotent payment row; it still needs cross-service wallet deduction through mall-service/internal APIs or a dedicated payment service.~~ DONE

## Handoff Tasks For Next Model

- Before implementing any business feature, locate its requirement ID in `REQUIREMENTS_TRACEABILITY_MATRIX.md`.
- ~~Migrate `UserHandler`, `UserService`, user models and JWT middleware into `user-service`.~~ DONE
- ~~Migrate product/category/cart/order/payment/comment/favorite modules into `mall-service`.~~ DONE
- ~~Migrate notice/visitor/parking/property fee modules into `community-service`.~~ DONE
- ~~Migrate repair and complaint domain logic into `workorder-service`.~~ DONE
- ~~Implement gateway JWT validation and Nacos service discovery.~~ DONE
- ~~Implement statistics-service aggregation endpoints.~~ DONE
- Add `order.created` event publishing from `mall-service`.
- Wire frontend API modules to `gateway-service` gradually.
- Add real LLM provider client and service clients to `agent-service`.
- ~~Implement property fee real wallet deduction via mall-service internal API.~~ DONE
- ~~Implement product_view_logs table for STAT-002 (商品访客排行).~~ DONE
