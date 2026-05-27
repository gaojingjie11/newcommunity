# API Contract By Requirements

本文档是基于 Word 需求陈述书的接口契约草案。字段和路径用于后续迁移约束，最终实现前必须回填到 `REQUIREMENTS_TRACEABILITY_MATRIX.md`。

统一响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## user-service

| 功能 | 方法 | 路径 | 请求要点 | 响应要点 |
| --- | --- | --- | --- | --- |
| 注册 | POST | `/api/users/register` | `mobile`、`password`、`real_name`、`age`、`gender` 必填；`mobile` 唯一 | 用户 ID、基础资料 |
| 登录 | POST | `/api/users/login` | `mobile`、`password` | JWT、用户信息、角色 |
| 忘记密码发送验证码 | POST | `/api/users/password-reset/code` | `mobile` | 验证码发送状态 |
| 忘记密码重置 | POST | `/api/users/password-reset` | `mobile`、`code`、`new_password` | 重置结果 |
| 修改密码 | PUT | `/api/users/me/password` | `old_password`、`new_password` | 修改结果 |
| 个人资料查看 | GET | `/api/users/me` | JWT | 用户资料 |
| 个人资料修改 | PUT | `/api/users/me` | `avatar`、`mobile`、`username`、`gender`、`email` 等 | 修改结果 |
| 用户登录日志记录 | POST | internal `/api/users/login-logs` | 登录成功/失败、IP、UA、时间 | 记录结果 |
| 用户登录日志查询 | GET | `/api/admin/user-login-logs` | 管理端分页筛选 | 日志列表 |
| 管理员登录日志查询 | GET | `/api/admin/admin-login-logs` | 管理端分页筛选 | 日志列表 |

## mall-service

| 功能 | 方法 | 路径 | 请求要点 | 响应要点 |
| --- | --- | --- | --- | --- |
| 全部商品 | GET | `/api/mall/products` | `page`、`size`、`category_id` | 商品列表 |
| 商品搜索 | GET | `/api/mall/products/search` | `keyword` 匹配商品名称、商品简介 | 搜索结果 |
| 促销商品 | GET | `/api/mall/products/promotions` | 促销类型、分页 | 促销商品 |
| 商品列表排序筛选 | GET | `/api/mall/products` | `sort=sales_desc/price_asc/price_desc`、`min_price`、`max_price` | 商品列表 |
| 商品详情 | GET | `/api/mall/products/{id}` | 商品 ID | 价格、库存、取货门店 |
| 门店列表 | GET | `/api/mall/stores` | `product_id` 可选 | 门店列表、库存 |
| 添加购物车 | POST | `/api/mall/cart/items` | `product_id`、`quantity`、`store_id` 可选 | 购物车项 |
| 移除购物车 | DELETE | `/api/mall/cart/items/{id}` | 购物车项 ID | 删除结果 |
| 修改购物车数量 | PUT | `/api/mall/cart/items/{id}` | `quantity` | 购物车金额重算 |
| 商品收藏 | POST | `/api/mall/favorites` | `product_id` | 收藏结果 |
| 取消收藏 | DELETE | `/api/mall/favorites/{product_id}` | 商品 ID | 取消结果 |
| 我的收藏 | GET | `/api/mall/favorites` | 分页 | 收藏商品 |
| 创建订单 | POST | `/api/mall/orders` | 购物车项、数量、取货门店 | 订单、库存扣减结果 |
| 订单支付 | POST | `/api/mall/orders/{id}/pay` | 钱包支付认证信息 | 支付结果 |
| 待付款订单 | GET | `/api/mall/orders?status=pending_payment` | 分页 | 订单列表 |
| 已付款订单 | GET | `/api/mall/orders?status=paid` | 分页 | 订单列表 |
| 待取货订单 | GET | `/api/mall/orders?status=ready_for_pickup` | 分页 | 订单列表 |
| 已完成订单 | GET | `/api/mall/orders?status=completed` | 分页 | 订单列表 |
| 充值 | POST | `/api/mall/wallet/recharge` | `amount` | 钱包余额、交易流水 |
| 转账 | POST | `/api/mall/wallet/transfer` | `target_mobile`、`amount` | 双方流水 |
| 账单 | GET | `/api/mall/wallet/transactions` | 类型、时间、分页 | 支付/转账/充值记录 |
| 钱包余额 | GET | `/api/mall/wallet/balance` | JWT | 当前余额 |
| 分类管理 | POST/PUT/DELETE | `/api/admin/mall/categories` | 类别信息 | 管理结果 |
| 商品信息管理 | POST/PUT/DELETE | `/api/admin/mall/products` | 商品基础信息、图片、库存 | 管理结果 |
| 营销管理 | POST/PUT/DELETE | `/api/admin/mall/promotions` | 促销类型、绑定商品 | 管理结果 |
| 服务区域管理 | POST/PUT/DELETE | `/api/admin/mall/service-areas` | 区域信息 | 管理结果 |
| 门店管理 | POST/PUT/DELETE | `/api/admin/mall/stores` | 区域、营业时间、位置 | 管理结果 |
| 门店商品管理 | POST/PUT/DELETE | `/api/admin/mall/store-products` | 门店、商品、库存、上下架状态 | 管理结果 |
| 后台订单管理 | GET/POST | `/api/admin/mall/orders` | 查询、发货、作废 | 管理结果 |

## community-service

| 功能 | 方法 | 路径 | 请求要点 | 响应要点 |
| --- | --- | --- | --- | --- |
| 公告列表查询 | GET | `/api/community/notices` | 分页 | 公告列表 |
| 公告发布 | POST | `/api/admin/community/notices` | 标题、内容、发布人 | 发布结果 |
| 公告浏览状态 | GET | `/api/admin/community/notices/{id}/views` | 公告 ID | 浏览次数、已读用户 |
| 公告标记已读 | POST | `/api/community/notices/{id}/read` | 公告 ID | 标记结果 |
| 访客登记 | POST | `/api/community/visitors` | `visit_purpose`、`release_time`、`valid_date`、访客信息 | 待审核申请 |
| 访客审核/通行 | POST | `/api/admin/community/visitors/{id}/audit` | `status`、`remark` | 审核结果 |
| 车位查询 | GET | `/api/community/parking-spaces/my` | JWT | 本人车位 |
| 车位绑定车牌号 | PUT | `/api/community/parking-spaces/{id}/plate` | `car_plate` | 绑定结果 |
| 车位统计 | GET | `/api/admin/community/parking-spaces/statistics` | 管理端 | 总数、已用、空闲 |
| 物业费缴纳 | POST | `/api/community/property-fees/{id}/pay` | `idempotency_key` 或 `Idempotency-Key`；后续接钱包扣款 | 缴费结果 |
| 缴费记录 | GET | `/api/community/property-fees/payments` | 用户或管理端筛选 | 缴费记录 |
| 后台缴费记录 | GET | `/api/admin/community/property-fees/payments` | 分页筛选 | 全部缴费记录 |

## workorder-service

| 功能 | 方法 | 路径 | 请求要点 | 响应要点 |
| --- | --- | --- | --- | --- |
| 报事维修提交 | POST | `/api/workorders/repairs` | `type`、`description`、附件可选 | 工单、`repair.created` 事件 |
| 报事处理 | POST | `/api/admin/workorders/repairs/{id}/process` | `status`、`result`、处理人 | 状态流转 |
| 事项投诉提交 | POST | `/api/workorders/complaints` | `complaint_type`、`description`、附件可选 | 投诉单、`complaint.created` 事件 |
| 投诉处理 | POST | `/api/admin/workorders/complaints/{id}/process` | `status`、`result`、处理人 | 状态流转 |
| 报修/投诉状态流转 | GET | `/api/workorders/{type}/{id}/logs` | 类型、ID | 状态日志 |
| RabbitMQ repair.created | event | `repair.created` | 工单 ID、用户 ID、类型、描述、时间 | 消费者异步处理 |
| RabbitMQ complaint.created | event | `complaint.created` | 投诉 ID、用户 ID、类型、描述、时间 | 消费者异步处理 |

## statistics-service

| 功能 | 方法 | 路径 | 请求要点 | 响应要点 |
| --- | --- | --- | --- | --- |
| 商品销售排行 | GET | `/api/statistics/products/sales-rank` | 时间范围、limit | 按销售额排行 |
| 商品访客排行 | GET | `/api/statistics/products/view-rank` | 时间范围、limit | 按浏览次数排行 |
| 社区运营概览 | GET | `/api/statistics/community/overview` | 时间范围 | 用户、订单、报修、投诉、缴费等摘要 |
| 订单统计 | GET | `/api/statistics/orders` | 时间范围、状态 | 订单数量、金额 |
| 报修投诉统计 | GET | `/api/statistics/workorders` | 时间范围、状态、类型 | 报修/投诉统计 |

## gateway-service

| 功能 | 方法 | 路径 | 请求要点 | 响应要点 |
| --- | --- | --- | --- | --- |
| 统一鉴权 | middleware | all protected APIs | JWT、角色、菜单权限 | 放行或拒绝 |
| 用户端 API 路由 | ANY | `/api/app/**` | 路由到用户端微服务 | 代理响应 |
| 管理端 API 路由 | ANY | `/api/admin/**` | 管理员 JWT、角色权限 | 代理响应 |
| 服务列表 | GET | `/api/gateway/services` | 无 | Nacos 或本地配置服务列表 |
| 健康检查 | GET | `/health` | 无 | gateway 状态 |

## agent-service

| 功能 | 方法 | 路径 | 请求要点 | 响应要点 |
| --- | --- | --- | --- | --- |
| 社区客服 Agent | POST | `/agent/chat` | 用户问题、上下文 | `intent`、`reply`、`called_services` |
| 报修派单 Agent | POST | `/agent/repair-classify` | 报修描述、图片 | `category`、`urgency`、建议部门 |
| 投诉风险 Agent | POST | `/agent/complaint-risk` | 投诉描述、用户上下文 | 分类、风险等级、建议动作 |
| 商品/社区服务推荐 Agent | POST | `/agent/recommend` | 用户、场景、上下文 | 推荐项、原因 |
