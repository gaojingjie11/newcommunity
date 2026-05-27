# Data Model Plan

数据模型以 Word 需求为基线，结合当前 Go 微服务拆分规划。当前阶段不强制拆多数据库，允许开发期共用 `smart_community`，但表归属必须按服务边界维护。

## user-service

### `sys_user`

- 用户主体信息。
- 字段：`id`、`username`、`password`、`real_name`、`mobile`（uniqueIndex）、`age`、`gender`、`email`、`avatar`、`green_points`、`balance`、`face_registered`、`face_image_url`、`role`（legacy 兼容，不作为主权限依据）、`status`、`created_at`、`updated_at`。
- 约束：`mobile` 唯一且注册必填。

### `sys_role`

- 角色信息。
- 字段：`id`、`name`、`code`（uniqueIndex）、`remark`、`created_at`。
- 预置角色：`admin`、`property`、`store`、`user`。

### `sys_menu`

- 前端菜单，控制管理端页面可见性。
- 字段：`id`、`parent_id`、`name`、`path`、`component`、`sort`、`type`、`created_at`。

### `sys_permission`

- 后端 API 权限点，每个权限点对应一个或多个 API 接口。
- 字段：`id`、`code`（uniqueIndex，格式 `domain:resource:action`）、`name`、`resource`、`method`（HTTP 方法）、`path`（API 路径模式）、`type`（1=菜单权限, 2=操作权限）、`status`（1=启用）、`created_at`、`updated_at`。

### `sys_user_role`

- 用户与角色的多对多关系。一个用户可拥有多个角色，权限取并集。
- 字段：`id`、`user_id`、`role_id`。
- 约束：`unique(user_id, role_id)`。

### `sys_role_menu`

- 角色与菜单的多对多关系。
- 字段：`id`、`role_id`、`menu_id`。

### `sys_role_permission`

- 角色与权限点的多对多关系。决定角色可以访问哪些后端 API。
- 字段：`id`、`role_id`、`permission_id`。
- 约束：`unique(role_id, permission_id)`。

### `user_login_logs`

- 用户端登录日志。
- 字段：`id`、`user_id`、`mobile`、`login_time`、`ip`、`user_agent`、`success`、`failure_reason`。

### `admin_login_logs`

- 管理端登录日志。
- 字段：`id`、`admin_user_id`、`mobile`、`login_time`、`ip`、`user_agent`、`success`、`failure_reason`。

### `password_reset_codes`

- 密码重置验证码记录。
- 字段：`id`、`mobile`、`code_hash`（bcrypt 哈希，不存明文）、`expires_at`、`used_at`、`created_at`。
- 验证码校验以 Redis 为主，DB 用于审计。

## mall-service

### `products`

- 商品基础信息。
- 字段建议：`id`、`category_id`、`name`、`description`、`price`、`original_price`、`stock`、`image_url`、`status`、`sales_count`、`created_at`、`updated_at`。

### `product_categories`

- 商品分类。
- 字段建议：`id`、`name`、`icon`、`sort`、`created_at`、`updated_at`。

### `promotions`

- 促销活动或促销类型。
- 字段建议：`id`、`title`、`type`、`start_at`、`end_at`、`status`、`created_at`、`updated_at`。

### `pms_promotion_product`

- 促销与商品绑定关系。
- 字段建议：`id`、`promotion_id`、`product_id`。

### `stores`

- 取货门店。
- 字段建议：`id`、`service_area_id`、`name`、`address`、`phone`、`business_hours`、`longitude`、`latitude`、`status`。

### `service_areas`

- 服务区域。
- 字段建议：`id`、`name`、`city`、`district`、`address_scope`、`status`。

### `store_products`

- 门店商品绑定和库存。
- 字段建议：`id`、`store_id`、`product_id`、`stock`、`status`、`created_at`、`updated_at`。
- `status` 表示上架/下架。

### `carts`

- 购物车。
- 字段建议：`id`、`user_id`、`product_id`、`store_id`、`quantity`、`created_at`、`updated_at`。

### `favorites`

- 商品收藏。
- 字段建议：`id`、`user_id`、`product_id`、`created_at`。

### `pms_product_comment` (已实现)

- 商品评价，mall-service 拥有。
- 字段：`id`、`user_id`、`product_id`、`content`、`rating`、`created_at`。
- 关联展示：通过 `user_id` 读取 `sys_user` 的 `username`、`real_name`、`avatar` 作为评价区用户信息。

### `orders`

- 订单主表。
- 字段建议：`id`、`order_no`、`user_id`、`store_id`、`total_amount`、`status`、`paid_at`、`picked_up_at`、`completed_at`、`cancelled_at`、`created_at`。
- 状态建议：`pending_payment`、`paid`、`ready_for_pickup`、`completed`、`cancelled`。

### `order_items`

- 订单明细。
- 字段建议：`id`、`order_id`、`product_id`、`product_name_snapshot`、`price`、`quantity`、`amount`。

### `wallets`

- 用户钱包。
- 字段建议：`id`、`user_id`、`balance`、`status`、`created_at`、`updated_at`。

### `wallet_transactions`

- 钱包流水。
- 字段建议：`id`、`wallet_id`、`user_id`、`type`、`amount`、`direction`、`related_type`、`related_id`、`remark`、`created_at`。
- 类型：充值、转账入、转账出、订单支付、物业费支付。

### `product_view_logs` (已实现)

- 商品浏览日志，mall-service 拥有，statistics-service 只读聚合。
- 字段：`id`、`product_id`、`user_id`（0=匿名）、`ip`、`user_agent`、`viewed_at`。
- 索引：`idx_product_id`、`idx_user_id`、`idx_viewed_at`。
- 写入时机：`GET /api/mall/products/:id` 商品详情被访问时异步记录。
- 读取方：`statistics-service` 的 `ProductViewRank` 聚合查询，按浏览次数排行。

## community-service

### `notices`

- 公告主表。
- 字段建议：`id`、`title`、`content`、`publisher`、`view_count`、`status`、`created_at`、`updated_at`。

### `notice_view_logs`

- 公告浏览状态明细。
- 字段建议：`id`、`notice_id`、`user_id`、`viewed_at`、`read_at`。

### `visitors`

- 访客登记。
- 字段建议：`id`、`user_id`、`visitor_name`、`visitor_phone`、`visit_purpose`、`release_time`、`valid_date`、`status`、`audit_remark`、`audit_at`、`created_at`、`updated_at`。

### `parking_spaces`

- 车位主表。
- 字段建议：`id`、`parking_no`、`status`、`created_at`、`updated_at`。

### `user_parking_bindings`

- 用户与车位绑定关系。
- 字段建议：`id`、`user_id`、`parking_space_id`、`car_plate`、`status`、`bound_at`、`unbound_at`。

### `property_fees`

- 物业费账单。
- 字段建议：`id`、`user_id`、`month`、`amount`（分）、`status`、`due_date`、`paid_at`、`created_at`、`updated_at`。

### `property_fee_payments`

- 物业费缴费记录。
- 字段建议：`id`、`property_fee_id`、`user_id`、`amount`（分）、`wallet_transaction_id`、`idempotency_key`、`paid_at`、`status`。

## workorder-service

### `repairs`

- 报事维修。
- 字段建议：`id`、`user_id`、`matter_type`、`description`、`status`、`result`、`processor_id`、`created_at`、`updated_at`、`processed_at`。

### `complaints`

- 事项投诉。
- 字段建议：`id`、`user_id`、`complaint_type`、`description`、`status`、`result`、`processor_id`、`created_at`、`updated_at`、`processed_at`。

### `workorder_logs`

- 报修/投诉状态流转日志。
- 字段建议：`id`、`target_type`、`target_id`、`from_status`、`to_status`、`operator_id`、`action`、`remark`、`created_at`。

## statistics-service

统计服务优先基于其他服务表聚合，不强制独立建表。

可选缓存表：

- `statistics_snapshots`：按日保存运营概览。
- `product_rank_snapshots`：按日保存商品销售/访客排行。

## agent-service

### `agent_sessions`

- Agent 会话。
- 字段建议：`id`、`user_id`、`agent_type`、`status`、`created_at`、`updated_at`。

### `agent_messages`

- Agent 消息。
- 字段建议：`id`、`session_id`、`role`、`content`、`created_at`。

### `agent_action_logs`

- Agent 调用业务服务、工具或模型的记录。
- 字段建议：`id`、`session_id`、`agent_type`、`action`、`request_payload`、`response_payload`、`success`、`created_at`。
