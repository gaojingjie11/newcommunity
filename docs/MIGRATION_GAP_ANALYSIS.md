# Migration Gap Analysis

基线来源：`需求陈述书-东软智慧社区项目.docx`。

结论原则：

- Word 需求为准。
- legacy `smartcomunity` 仅用于识别可迁移代码，不作为需求依据。
- 当前微服务大多仍是占位接口，业务迁移必须先匹配追踪矩阵编号。

## A. legacy 已有且可迁移

### 注册登录与用户

- 手机号唯一识别：`model.SysUser.Mobile` 有 `uniqueIndex`。
- 手机号密码登录：`UserService.Login(mobile, password, ip, userAgent)`。
- 修改密码：`UserService.ChangePassword`。
- 个人资料查看与修改：`UserHandler.Info`、`UserHandler.Update`。
- 退出登录：Redis token 删除逻辑可迁移。
- 角色、菜单、用户角色、角色菜单模型：`SysRole`、`SysMenu`、`SysUserRole`、`SysRoleMenu`。
- 角色绑定菜单：`AdminService.BindRoleMenu`。

### 商城

- 商品列表、价格排序、销量排序、价格范围筛选：`ProductService.GetList`。
- 促销商品筛选：legacy 使用 `original_price > price` 和 `is_promotion`。
- 商品增删改：`ProductService.Create/Update/Delete`。
- 商品销量排行：`ProductService.GetSalesRank`。
- 购物车添加、删除、数量修改、列表：`CartService`。
- 收藏、取消收藏、我的收藏：`FavoriteService`。
- 订单创建、支付、发货、收货、取消、订单状态列表：`OrderService`。
- 充值、转账、交易记录：`FinanceService.Recharge/Transfer/GetTransactionList`。
- 门店列表、门店 CRUD、门店绑定商品和库存：`StoreService`。

### 社区

- 公告列表、公告详情、公告发布、删除、标记已读：`NoticeService`。
- 访客登记和审核：`SecurityService.CreateVisitor/AuditVisitor`。
- 车位查询、车牌绑定、车位统计、管理员分配车位：`SecurityService`。
- 物业费创建、列表、缴费：`FinanceService`。
- 报事维修提交和处理：`RepairService`。

### 统计与扩展

- Dashboard 基础统计：`AdminService.GetDashboardStats`。
- AI 报表与社区运行摘要：`AdminService.GenerateAIReport`、`AIReportScheduler`。
- 社区群聊、绿色积分、垃圾识别、人脸支付属于 legacy 扩展，可保留但不是 Word 主流程。

## B. legacy 已有但逻辑需要修正

### 注册字段校验

Word 要求注册必填：手机号、密码、真实姓名、年龄、性别。legacy `UserHandler.Register` 只强校验手机号和密码，真实姓名、年龄、性别未强制校验。

### 忘记密码验证码

Word 要求手机验证码重置。legacy 有 `SendCode`，但 `ResetPassword` 中验证码校验写死为 `123456`，未读取 Redis 验证码。

### 登录日志

`UserService.Login` 接收 IP 和 UserAgent，但没有写入用户登录日志或管理员登录日志。

### 商品搜索

Word 要求商品名称、商品简介均可搜索。legacy `ProductService.GetList` 只按 `name LIKE` 搜索，缺 `description` 搜索。

### 商品详情和门店库存

Word 要求商品详情展示价格、库存、取货门店。legacy 商品详情只返回商品本身；门店库存存在 `StoreProduct`，但详情接口未聚合门店可取货库存。

### 订单库存

Word 要求创建订单操作库存。legacy 扣减 `Product.Stock` 和 `Product.Sales`，但未严格扣减所选门店的 `StoreProduct.Stock`。

### 钱包模型

Word 要求系统钱包。legacy 钱包余额放在 `sys_user.balance`，交易流水为 `sys_transaction`。微服务阶段应拆为 `wallets`、`wallet_transactions`。

### 促销模型

Word 要求自定义促销类型并为促销绑定商品。legacy `Promotion` 直接含 `ProductID`，不支持一个促销绑定多个商品，也缺完整修改接口。

### 分类管理

legacy 有分类模型和分类列表，但后台分类新增、修改、删除接口不完整。

### 访客登记字段

Word 要求来访目的、放行时间及有效日期。legacy `Visitor` 有 `Reason`、`VisitTime`，但没有独立 `valid_date` 或有效期字段。

### 报事维修和投诉

Word 将报事维修、事项投诉列为不同流程。legacy 通过 `Repair.Type` 区分报修/投诉，未建立独立 `complaints` 表，投诉处理、风险、状态日志不清晰。

### 公告浏览状态

legacy 有 `view_count` 和 `NoticeRead`，但缺完整后台浏览状态查询接口，例如已读用户、未读用户、阅读时间列表。

### 商品销售排行

Word 要求按销售额排行。legacy `GetSalesRank` 按 `sales` 数量排行，不是销售额排行。

### 服务区域管理

legacy `Store.Region` 是字符串字段，没有独立 `service_areas` 表和区域 CRUD。

## C. legacy 没有或明显缺失

- `user_login_logs`：用户登录日志表和管理查询接口。
- `admin_login_logs`：管理员登录日志表和管理查询接口。
- `password_reset_codes`：密码重置验证码独立持久化表。当前主要依赖 Redis/mock。
- `wallets`：独立钱包表。
- `wallet_transactions`：统一钱包流水表。
- `product_view_logs`：商品浏览日志，支撑商品访客排行。
- `service_areas`：服务区域表。
- `pms_promotion_product`：促销与商品多对多绑定表。
- `notice_view_logs`：公告浏览状态明细表。
- `user_parking_bindings`：用户与车位绑定关系表。legacy 直接写在 parking 上。
- `property_fee_payments`：物业费缴费记录明细表。
- `workorders`：统一承载报修/投诉，使用 `type=repair|complaint` 区分，不再保留独立 `repairs`、`complaints` 表。
- `workorder_logs`：报修/投诉状态流转日志。
- `agent_sessions`、`agent_messages`、`agent_action_logs`：Agent 运行记录表。

## D. legacy 有但 Word 未明确要求，可作为扩展

- 绿色积分和垃圾分类识别：`GreenPointService`。
- 人脸录入和人脸支付：`FaceService`、`RegisterFace`、AI 支付分支。
- AI 聊天与工具调用：`AIService.ChatWithMemory`。
- 社区群聊：`CommunityMessageService`。
- AI 社区运营日报：`AIReportScheduler` 和 `AIReport`。
- MinIO 通用上传：`StorageService`。
- 商品评论：`CommentService`。
- 收藏检查接口：`FavoriteService.Check`。

这些能力可以保留为扩展，但不能改变 Word 主流程验收。

## E. 当前微服务骨架已完成但仍是占位

- `gateway-service`
  - 已有 `/health`、`/api/gateway/services`、配置式代理占位。
  - 未实现真实统一鉴权、Nacos 服务发现、管理端/用户端路由策略。

- `user-service`
  - 已有 `/health`、`/api/users/ping`、登录占位。
  - 未迁移注册、忘记密码、修改密码、资料、角色菜单、登录日志。

- `mall-service`
  - 已有 `/health`、`/api/mall/ping`、商品列表占位。
  - 未迁移商品、购物车、订单、钱包、门店、营销。

- `community-service`
  - 已有 `/health`、`/api/community/ping`、公告列表占位。
  - 未迁移公告、访客、车位、物业费、公告浏览状态。

- `workorder-service`
  - 已有 `/health`、`/api/workorders/ping`、报修/投诉提交占位。
  - 已预留 `repair.created`、`complaint.created` RabbitMQ 发布能力。
  - 未实现真实报事、投诉、状态流转、处理记录。

- `statistics-service`
  - 已有 `/health`、`/api/statistics/overview` 占位。
  - 未实现商品销售排行、商品访客排行、订单统计、报修投诉统计。

- `agent-service`
  - 已有 4 个 Agent 占位接口。
  - 未接 LLM，未接业务服务，未记录 Agent 会话和动作日志。

## 重点检查结论

| 检查点 | 结论 |
| --- | --- |
| 注册是否以手机号作为唯一身份识别 | 已有，`mobile uniqueIndex` |
| 注册字段是否包含真实姓名、年龄、性别 | 字段已有，必填校验不足 |
| 登录是否使用手机号和密码 | 已有 |
| 是否有忘记密码 | 部分有，验证码校验需修正 |
| 是否有钱包充值、转账、账单 | 已有，但需拆钱包模型 |
| 是否有商品销量排序、价格排序、价格范围筛选 | 已有 |
| 是否有促销商品 | 部分有，促销模型简化 |
| 是否有商品访客排行 | 缺失 |
| 是否有门店、服务区域、门店商品库存绑定 | 门店和门店商品有，服务区域缺独立模型，上下架不足 |
| 是否有待付款、已付款、待取货、已完成订单 | 已有状态能力 |
| 是否有访客登记审核流程 | 已有，但有效日期缺失 |
| 是否有车位绑定车牌号 | 已有 |
| 是否有物业费缴纳和缴费记录 | 部分有，缺独立 payment 明细 |
| 是否有管理员登录日志和用户登录日志 | 缺失 |
| 是否有角色绑定菜单 | 已有 |
| 是否有公告浏览状态 | 部分有，后台状态明细不足 |
| 是否有投诉处理和报事处理状态流转 | 部分有，投诉复用 repair，缺状态日志 |
