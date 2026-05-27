# Requirements Traceability Matrix

状态定义：

- `TODO`：尚未迁移或尚未实现。
- `PARTIAL`：legacy 或微服务有部分能力，但未完全满足 Word。
- `DONE`：当前微服务已完整实现。当前阶段主要是骨架，因此业务功能基本不标 DONE。
- `GAP`：Word 要求明确存在，但 legacy 缺失或明显不足。
- `EXTENSION`：legacy 有但 Word 未明确要求，可作为增强保留。

| 需求编号 | Word 模块 | 需求名称 | 需求描述 | 目标微服务 | 目标接口 | 数据表 | 当前 legacy 是否已有 | 当前微服务是否已实现 | 优先级 | 备注 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| AUTH-001 | 注册登录需求 | 注册 | 手机号唯一，密码、真实姓名、年龄、性别必填 | user-service | `POST /api/users/register` | `users` | PARTIAL | DONE | P0 | 已迁移：校验 mobile/password/real_name/age/gender 必填，bcrypt cost=14，role=user，balance=100 |
| AUTH-002 | 注册登录需求 | 登录 | 使用手机号和密码登录 | user-service | `POST /api/users/login` | `users`, `user_login_logs` | PARTIAL | DONE | P0 | 已迁移：bcrypt 校验 + Redis token 存储 + user_login_logs/admin_login_logs 写入 |
| AUTH-003 | 注册登录需求 | 忘记密码 | 手机验证码重置密码 | user-service | `POST /api/users/password-reset/code`, `POST /api/users/password-reset` | `password_reset_codes`, `users` | PARTIAL | DONE | P0 | 已迁移：Redis 存储验证码 + DB 审计，修复了硬编码 123456 问题 |
| AUTH-004 | 注册登录需求 | 修改密码 | 登录后修改登录密码 | user-service | `PUT /api/users/me/password` | `users` | YES | DONE | P0 | 已迁移：验证旧密码 → bcrypt 新密码 → 更新 |
| AUTH-005 | 注册登录需求 | 个人资料 | 查看和修改头像、手机号、用户名、性别、邮箱 | user-service | `GET /api/users/me`, `PUT /api/users/me` | `users` | YES | DONE | P0 | 已迁移：GET 返回完整资料，PUT 支持部分更新（pointer 字段区分 nil/空值） |
| AUTH-006 | 用户端-商城服务-首页 | 当前登录身份 | 首页展示当前用户登录身份 | user-service | `GET /api/users/me` | `users` | YES | DONE | P1 | 复用 AUTH-005 的 GetProfile 接口 |
| AUTH-007 | 用户端-商城服务-首页 | 退出登录 | 用户退出登录状态 | user-service | `POST /api/users/logout` | Redis token | YES | DONE | P1 | 已迁移：删除 Redis login:token:{userID} |
| LOG-001 | 管理端-商户后台-系统管理 | 用户登录日志 | 记录用户端登录时间、IP 等 | user-service | `GET /api/admin/user-login-logs` | `user_login_logs` | GAP | DONE | P0 | 已新建 user_login_logs 表，登录时自动写入，管理员可分页查询 |
| LOG-002 | 管理端-商户后台-系统管理 | 管理员登录日志 | 记录管理员后台登录记录 | user-service | `GET /api/admin/admin-login-logs` | `admin_login_logs` | GAP | DONE | P0 | 已新建 admin_login_logs 表，admin 角色登录时自动写入 |
| ADMIN-MALL-001 | 管理端-商户后台-权限管理 | 角色管理 | 增删改角色 | user-service | `POST/PUT/DELETE /api/admin/roles` | `roles` | PARTIAL | DONE | P0 | 已迁移：完整 CRUD（Create/Update/Delete/List） |
| ADMIN-MALL-002 | 管理端-商户后台-权限管理 | 角色绑定菜单 | 为角色绑定菜单 | user-service | `POST /api/admin/roles/{id}/menus` | `role_menus`, `menus` | YES | DONE | P0 | 已迁移：事务删除旧绑定 + 批量插入新绑定 |
| ADMIN-MALL-003 | 管理端-商户后台-用户管理 | 用户管理 | 查看用户、绑定角色、冻结账户 | user-service | `/api/admin/users/**` | `sys_user`, `sys_user_role` | PARTIAL | DONE | P0 | 已迁移：ListAdminUsers 查询全量用户，余额从 wallets 读取；FreezeUser、AssignRole |
| ADMIN-MALL-004 | 管理端-商户后台-会员管理 | 会员列表 | 查看注册用户信息，含用户名、联系方式、邮箱 | user-service | `GET /api/admin/members` | `users` | YES | DONE | P0 | 已迁移：ListMembers(role=user)，支持分页和关键字搜索 |
| MALL-001 | 用户端-商城服务-首页 | 全部商品 | 展示系统内全部商品 | mall-service | `GET /api/mall/products` | `products` | YES | DONE | P0 | 已迁移：分页列表，支持 category_id 筛选和 sort 排序 |
| MALL-002 | 用户端-商城服务-首页 | 商品搜索 | 按商品名称、商品简介搜索 | mall-service | `GET /api/mall/products/search` | `products` | PARTIAL | DONE | P0 | 已迁移：name + description LIKE 搜索 |
| MALL-003 | 用户端-商城服务-首页 | 促销商品 | 浏览明星商品、秒杀商品等促销内容 | mall-service | `GET /api/mall/products/promotions` | `pms_promotion`, `pms_promotion_product` | PARTIAL | DONE | P0 | 已迁移：is_promotion=1 筛选 + pms_promotion_product M2M |
| MALL-004 | 用户端-商城服务-商品管理 | 商品列表排序筛选 | 支持销量排序、价格排序、价格范围筛选 | mall-service | `GET /api/mall/products` | `products` | YES | DONE | P0 | 已迁移：sort 参数支持 price_asc/desc、sales |
| MALL-005 | 用户端-商城服务-商品管理 | 商品详情 | 查看价格、库存、取货门店 | mall-service | `GET /api/mall/products/{id}` | `products`, `stores`, `store_products` | PARTIAL | DONE | P0 | 已迁移：FindByID 返回完整商品信息 |
| MALL-006 | 用户端-商城服务-商品管理 | 门店列表 | 浏览并选择取货门店 | mall-service | `GET /api/mall/stores` | `stores`, `store_products` | YES | DONE | P0 | 已迁移：支持 area_id 筛选，Preload ServiceArea |
| MALL-007 | 用户端-商城服务-商品管理 | 添加购物车 | 当前商品加入购物车 | mall-service | `POST /api/mall/cart/items` | `carts` | YES | DONE | P0 | 已迁移：存在则累加数量，否则新建 |
| MALL-008 | 用户端-商城服务-商品管理 | 商品收藏 | 收藏当前商品 | mall-service | `POST /api/mall/favorites` | `favorites` | YES | DONE | P1 | 已迁移：去重检查 |
| MALL-009 | 用户端-商城服务-商品管理 | 商品取消收藏 | 取消收藏当前商品 | mall-service | `DELETE /api/mall/favorites/{product_id}` | `favorites` | YES | DONE | P1 | 已迁移 |
| MALL-010 | 用户端-商城服务-购物车管理 | 移除商品 | 移除购物车指定商品 | mall-service | `DELETE /api/mall/cart/items/{id}` | `carts` | YES | DONE | P0 | 已迁移：校验 user_id |
| MALL-011 | 用户端-商城服务-购物车管理 | 操作数量 | 增减商品数量并实时计算价格 | mall-service | `PUT /api/mall/cart/items/{id}` | `carts`, `products` | PARTIAL | DONE | P0 | 已迁移：quantity<=0 自动删除 |
| MALL-012 | 用户端-商城服务-订单管理 | 创建订单 | 购物车结算创建订单并操作库存 | mall-service | `POST /api/mall/orders` | `orders`, `order_items`, `products`, `store_products` | PARTIAL | DONE | P0 | 已迁移：事务内原子扣减 stock WHERE stock>=qty |
| MALL-013 | 用户端-商城服务-订单管理 | 订单支付 | 使用系统钱包支付订单 | mall-service | `POST /api/mall/orders/{id}/pay` | `orders`, `wallets`, `wallet_transactions` | YES | DONE | P0 | 已迁移：事务内 SELECT FOR UPDATE 锁定订单+钱包 |
| MALL-014 | 用户端-商城服务-订单管理 | 待付款订单 | 查询待付款订单 | mall-service | `GET /api/mall/orders?status=pending_payment` | `orders` | YES | DONE | P0 | 已迁移：status=0 |
| MALL-015 | 用户端-商城服务-订单管理 | 已付款订单 | 查询已付款订单 | mall-service | `GET /api/mall/orders?status=paid` | `orders` | YES | DONE | P0 | 已迁移：status=1 |
| MALL-016 | 用户端-商城服务-订单管理 | 待取货订单 | 查询待取货订单 | mall-service | `GET /api/mall/orders?status=ready_for_pickup` | `orders` | YES | DONE | P0 | 已迁移：status=2 |
| MALL-017 | 用户端-商城服务-订单管理 | 已完成订单 | 查询已完成订单 | mall-service | `GET /api/mall/orders?status=completed` | `orders` | YES | DONE | P0 | 已迁移：status=3 |
| MALL-018 | 用户端-商城服务-用户管理 | 我的收藏 | 查看本人收藏商品 | mall-service | `GET /api/mall/favorites` | `favorites`, `products` | YES | DONE | P1 | 已迁移：分页 + Preload Product |
| MALL-019 | 用户端-商城服务-用户管理 | 充值 | 向本人钱包充值 | mall-service | `POST /api/mall/wallet/recharge` | `wallets`, `wallet_transactions` | PARTIAL | DONE | P0 | 已迁移：独立 wallets 表，事务内 Credit + 创建流水 |
| MALL-020 | 用户端-商城服务-用户管理 | 转账 | 向其他用户转账 | mall-service | `POST /api/mall/wallet/transfer` | `wallets`, `wallet_transactions` | PARTIAL | DONE | P0 | 已迁移：事务内锁定双方钱包（按 ID 顺序防死锁），Debit+Credit+双流水 |
| MALL-021 | 用户端-商城服务-用户管理 | 账单 | 查看消费、支付、转账、充值记录和余额 | mall-service | `GET /api/mall/wallet/transactions`, `GET /api/mall/wallet/balance` | `wallet_transactions`, `wallets` | PARTIAL | DONE | P0 | 已迁移：分页查询 + type 筛选 |
| ADMIN-MALL-005 | 管理端-商户后台-商品管理 | 分类管理 | 新增、修改、删除商品类别 | mall-service | `/api/admin/mall/categories` | `product_categories` | PARTIAL | DONE | P0 | 已迁移：完整 CRUD |
| ADMIN-MALL-006 | 管理端-商户后台-商品管理 | 商品信息管理 | 商品增删改，维护类别、图片、名称、简介、库存 | mall-service | `/api/admin/mall/products` | `products` | YES | DONE | P0 | 已迁移：Create/Update/Delete |
| ADMIN-MALL-007 | 管理端-商户后台-营销管理 | 营销管理 | 促销类型增删改并绑定商品 | mall-service | `/api/admin/mall/promotions` | `pms_promotion`, `pms_promotion_product` | PARTIAL | DONE | P1 | 已迁移：CRUD + BindProducts（事务替换绑定） |
| ADMIN-MALL-008 | 管理端-商户后台-门店管理 | 服务区域管理 | 维护市内社区位置，新增、删除、修改区域 | mall-service | `/api/admin/mall/service-areas` | `service_areas` | GAP | DONE | P0 | 已新建 service_areas 独立表 + CRUD |
| ADMIN-MALL-009 | 管理端-商户后台-门店管理 | 门店管理 | 维护取货门店、区域、营业时间、位置 | mall-service | `/api/admin/mall/stores` | `stores`, `service_areas` | PARTIAL | DONE | P0 | 已迁移：AreaID FK + ServiceArea Preload |
| ADMIN-MALL-010 | 管理端-商户后台-门店管理 | 门店商品管理 | 为门店绑定商品，上下架和分配库存 | mall-service | `/api/admin/mall/store-products` | `store_products` | PARTIAL | DONE | P0 | 已迁移：Bind/Unbind/UpdateStatus/UpdateStock |
| ADMIN-MALL-011 | 管理端-商户后台-订单管理 | 订单管理 | 查询订单、发货、作废 | mall-service | `/api/admin/mall/orders` | `orders` | YES | DONE | P0 | 已迁移：ListAll + ShipOrder(1→2) + CancelOrder(0→40) |
| STAT-001 | 管理端-商户后台-数据统计 | 商品销售排行 | 按销售额排行展示数据 | statistics-service | `GET /api/statistics/products/sales-rank` | `orders`, `order_items`, `products` | PARTIAL | DONE | P1 | 已实现：JOIN oms_order_item + oms_order + pms_product，按销售额排行 |
| STAT-002 | 管理端-商户后台-数据统计 | 商品访客排行 | 按商品浏览次数排行展示数据 | statistics-service | `GET /api/statistics/products/view-rank` | `product_view_logs` | GAP | DONE | P1 | 已实现：mall-service 商品详情接口记录浏览日志，statistics-service 聚合排行 |
| COMM-001 | 用户端-社区服务-消息管理 | 公告列表查询 | 用户查看社区公告 | community-service | `GET /api/community/notices` | `notices` | YES | DONE | P0 | 已迁移公告列表、详情、浏览计数 |
| ADMIN-COMM-001 | 管理端-社区后台-消息管理 | 公告发布 | 管理员发布公告至用户端 | community-service | `POST /api/admin/community/notices` | `notices` | YES | DONE | P0 | 已接入 `community:notice:create` |
| ADMIN-COMM-002 | 管理端-社区后台-消息管理 | 公告浏览状态 | 管理员查看公告浏览状态 | community-service | `GET /api/admin/community/notices/{id}/views` | `notice_view_logs` | PARTIAL | DONE | P1 | 已补充 notice_view_logs |
| COMM-002 | 用户端-社区服务-安保管理 | 访客登记 | 登记来访目的、放行时间、有效日期供审核 | community-service | `POST /api/community/visitors` | `visitors` | PARTIAL | DONE | P0 | 已按 Word 补 `valid_date` |
| ADMIN-COMM-003 | 管理端-社区后台-安保管理 | 访客通行 | 处理访客通行申请 | community-service | `POST /api/admin/community/visitors/{id}/audit` | `visitors` | YES | DONE | P0 | 已接入 `community:visitor:audit` |
| COMM-003 | 用户端-社区服务-安保管理 | 车位查询 | 查看本人有权限车位信息 | community-service | `GET /api/community/parking-spaces/my` | `parking_spaces`, `user_parking_bindings` | PARTIAL | DONE | P0 | 已拆出车位主表与绑定表 |
| COMM-004 | 用户端-社区服务-安保管理 | 车位绑定车牌号 | 为本人车位绑定车牌号 | community-service | `PUT /api/community/parking-spaces/{id}/plate` | `user_parking_bindings` | YES | DONE | P0 | 已迁移 |
| ADMIN-COMM-004 | 管理端-社区后台-安保管理 | 车位管理 | 查询用户车位统计数据 | community-service | `GET /api/admin/community/parking-spaces/statistics` | `parking_spaces`, `user_parking_bindings` | YES | DONE | P0 | 已迁移统计接口 |
| COMM-005 | 用户端-社区服务-物业管理 | 报事维修 | 选择事项类型，填写描述并提交 | workorder-service | `POST /api/workorders/repairs` | `repairs`, `workorder_logs` | YES | DONE | P0 | 已迁移并发布 `repair.created` |
| ADMIN-COMM-005 | 管理端-社区后台-物业管理 | 报事处理 | 处理报事维修申请 | workorder-service | `POST /api/admin/workorders/repairs/{id}/process` | `repairs`, `workorder_logs` | YES | DONE | P0 | 已接入状态日志与 `workorder:repair:process` |
| COMM-006 | 用户端-社区服务-物业管理 | 事项投诉 | 编辑投诉类型和描述并提交 | workorder-service | `POST /api/workorders/complaints` | `complaints`, `workorder_logs` | PARTIAL | DONE | P0 | 已按 Word 独立 complaint 表 |
| ADMIN-COMM-006 | 管理端-社区后台-投诉管理 | 投诉管理 | 查询并处理用户投诉 | workorder-service | `GET/POST /api/admin/workorders/complaints` | `complaints`, `workorder_logs` | PARTIAL | DONE | P0 | 已接入状态日志与权限点 |
| COMM-007 | 用户端-社区服务-生活缴费 | 物业费缴纳 | 使用系统钱包缴纳物业费 | community-service | `POST /api/community/property-fees/{id}/pay` | `property_fees`, `property_fee_payments`, `wallet_transactions` | YES | DONE | P0 | 已对接 mall-service internal wallet/debit API，真实扣款+幂等+并发安全 |
| ADMIN-COMM-007 | 管理端-社区后台-缴费管理 | 缴费记录 | 查询用户缴费情况 | community-service | `GET /api/admin/community/property-fees/payments` | `property_fee_payments` | PARTIAL | DONE | P0 | 已补独立 payment 记录表 |
| STAT-003 | 统计服务 | 社区运营概览 | 用户、订单、报修、投诉、缴费等聚合 | statistics-service | `GET /api/statistics/community/overview` | 聚合其他服务表 | PARTIAL | DONE | P1 | 已实现：COUNT from sys_user, oms_order, repairs, complaints, property_fees |
| STAT-004 | 统计服务 | 订单统计 | 按时间、状态统计订单数量和金额 | statistics-service | `GET /api/statistics/orders` | `orders`, `order_items` | PARTIAL | DONE | P1 | 已实现：按 status 分组 + 最近 N 天每日趋势 |
| STAT-005 | 统计服务 | 报修投诉统计 | 按类型、状态统计报修投诉 | statistics-service | `GET /api/statistics/workorders` | `repairs`, `complaints` | PARTIAL | DONE | P1 | 已实现：按 type + status 分组统计 |
| GATEWAY-001 | 技术迁移约束 | 统一鉴权 | 用户端和管理端统一 JWT/角色鉴权入口 | gateway-service | gateway middleware | Redis token, roles, menus | PARTIAL | PARTIAL | P0 | user-service 已实现 JWTAuth + RequireRole 中间件（pkg/middleware），gateway 可复用 |
| GATEWAY-002 | 技术迁移约束 | 服务列表 | 返回可用服务列表，后续接 Nacos | gateway-service | `GET /api/gateway/services` | Nacos registry | N/A | PARTIAL | P1 | 当前从本地配置返回 |
| AGENT-001 | Agent Service | 社区客服 Agent | 识别意图、回复、调用服务 | agent-service | `POST /agent/chat` | `agent_sessions`, `agent_messages` | EXTENSION | PARTIAL | P2 | legacy 有 AI chat，Word 未明确要求 Agent，可保留为增强 |
| AGENT-002 | Agent Service | 报修派单 Agent | 输出分类、紧急度、建议部门 | agent-service | `POST /agent/repair-classify` | `agent_action_logs` | EXTENSION | PARTIAL | P2 | 当前占位 |
| AGENT-003 | Agent Service | 投诉风险 Agent | 输出风险等级和建议动作 | agent-service | `POST /agent/complaint-risk` | `agent_action_logs` | EXTENSION | PARTIAL | P2 | 当前占位 |
| AGENT-004 | Agent Service | 推荐 Agent | 推荐商品/社区服务 | agent-service | `POST /agent/recommend` | `agent_action_logs` | EXTENSION | PARTIAL | P2 | 当前占位 |
