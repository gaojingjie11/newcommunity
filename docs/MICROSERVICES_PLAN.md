# Microservices Plan

This stage only creates the skeleton. Business migration should be incremental and service by service.

## Requirement Governance

- `需求陈述书-东软智慧社区项目.docx` is the functional baseline and highest source of truth.
- Legacy `smartcomunity` code is migration reference only, not the requirement source.
- Every migrated feature must map to a row in `docs/REQUIREMENTS_TRACEABILITY_MATRIX.md`.
- If a feature has no traceability matrix ID, do not implement it directly. Add the requirement record first, then implement.
- Java technologies mentioned by the Word document are treated as functionally equivalent to the chosen Go microservice stack; functional behavior remains mandatory.

## gateway-service

Future scope: authentication, request routing, aggregation, service discovery, gateway-level rate limiting.

Requirement IDs: `GATEWAY-001`, `GATEWAY-002`.

Migration candidates:

- `internal/middleware/jwt.go`
- `internal/middleware/role.go`
- selected router aggregation from `internal/router/router.go`

## user-service

Future scope: users, login, registration, roles, JWT, SMS captcha, face registration metadata.

Requirement IDs: `AUTH-*`, `LOG-*`, `ADMIN-MALL-001` to `ADMIN-MALL-004`.

Migration candidates:

- `internal/controller/UserHandler.go`
- `internal/service/UserService.go`
- `internal/model/user.go`
- `pkg/utils/jwt.go`
- `pkg/utils/password.go`

## mall-service

Future scope: products, categories, carts, orders, payments, wallet, comments, favorites and transactions.

Requirement IDs: `MALL-*`, `ADMIN-MALL-005` to `ADMIN-MALL-011`, `STAT-001`, `STAT-002`.

Migration candidates:

- `ProductHandler.go`, `CartHandler.go`, `OrderHandler.go`, `CommentHandler.go`, `FavoriteHandler.go`, `FinanceHandler.go`
- `ProductService.go`, `CartService.go`, `OrderService.go`, `CommentService.go`, `FavoriteService.go`, `FinanceService.go`
- product, cart, order, transaction, comment and favorite models

## community-service

Future scope: notices, visitors, parking, property fees, community chat/messages, repairs, complaints, workorder dispatch, statistics, AI reports and basic community services.

Requirement IDs: `COMM-001` to `COMM-007`, `ADMIN-COMM-001` to `ADMIN-COMM-007`, `STAT-*`.

Migration candidates:

- `NoticeHandler.go`, `SecurityHandler.go`, `CommunityMessageHandler.go`
- `RepairHandler.go`, `RepairService.go`, `repair.go`
- `AdminHandler.go`, `AdminService.go`, `AIReportScheduler.go`
- notice, visitor, parking, property fee, community message, repair, complaint and workorder log services/models

Required event publishing:

- `repair.created`
- `complaint.created`

## agent-service

Future scope: LLM-powered customer service, repair classification, complaint risk analysis and recommendations.

Requirement IDs: `AGENT-*`. These are extensions to the Word baseline and must not replace required Word flows.

Migration candidates:

- AI chat and report prompts from `AIHandler.go` and `AIService.go`
- service client calls to all Go services through gateway or direct service URLs
