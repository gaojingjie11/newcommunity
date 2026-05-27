# Agent Architecture V2

本文档定义智慧社区项目下一阶段的 Agent 完整设计方案。目标不是只把 4 个占位接口换成大模型调用，而是把 Agent 做成可落地、可审计、可控权限、可持续演进的业务智能层。

## 1. 目标

Agent 的定位是“智能交互和智能分析层”，不是替代业务微服务。

必须遵守以下边界：

- Word 需求中的注册、登录、下单、支付、访客审核、物业费缴纳、工单处理等主流程，仍由现有 Go 业务服务负责。
- Agent 可以理解意图、调用工具、组织回答、生成报告、辅助创建工单或订单，但不能绕过业务服务直接改库。
- 所有写操作必须走既有业务接口或新增的内部受控接口。
- 所有敏感分析能力必须受角色和权限控制。

## 2. 推荐技术选型

推荐将当前 `agent-service` 从 FastAPI 占位升级为：

`GoFrame + Eino`

选择原因：

- 你参考的 `SuperBizAgent` 项目已经验证了 `GoFrame + Eino` 的组合，包含 Agent 编排、SSE 流式输出、Tool Calling、记忆、RAG 预留等结构，迁移思路最顺。
- `Eino` 适合做 Agent 编排、Tool 调用、RAG、ReAct、多阶段图流程；它不是 Gin 的替代，而是 Agent 运行时框架。
- `GoFrame` 比 Gin 更适合承载你参考项目的 controller / middleware / config / cron 风格，也方便复用既有设计-规划-执行的组织方式。
- 现有业务微服务保持 `Gin + Nacos` 不动，Agent 单独演进，不会牵连主业务。

结论：

- 业务服务：继续 `Gin + Nacos`
- Agent 服务：升级为 `GoFrame + Eino`
- LLM 编排：`Eino`
- 记忆缓存：`Redis`
- 会话和审计：`MySQL`
- 文档 / 知识文件：`MinIO` + 后续 RAG 索引

## 3. 部署建议

### 3.1 推荐部署形态

推荐把 Agent 拆成两个运行职责：

- `agent-api`
  - 对外提供 `/agent/chat`、`/agent/chat/stream`、`/agent/recommend` 等接口
  - 处理用户会话、意图识别、工具编排
- `agent-worker`
  - 执行定时任务
  - 生成日报 / 周报
  - 异步摘要、记忆整理、RAG 索引更新

二者可以是同一个服务进程中的不同模块，也可以后续拆成独立 worker。

### 3.2 关于“Agent 在本机运行”的风险

你之前的设想是业务服务在服务器，Agent 可能跑在自己电脑上。这个方式适合开发调试，但不适合承载以下长期任务：

- 每天 9:00 自动生成报告
- 用户长期记忆整理
- 夜间摘要压缩
- RAG 索引增量更新

因为本机可能关机、断网、内网穿透不稳定。

推荐策略：

- 交互式 Agent 可以本机调试
- 定时报告和记忆整理必须部署在服务器常驻环境

如果必须保留本机 Agent：

- 将“每日 9:00 报告任务”放到服务器侧的 `community-service` 调度器或独立 `agent-worker`
- 本机 Agent 仅作为聊天入口，不承担定时任务

## 4. Agent 能力总览

Agent 最终形态不是 4 个互相孤立的接口，而是“1 个主 Chat Agent + 3 类专用子能力 + 1 个报告能力 + 1 套记忆系统”。

### 4.1 主 Chat Agent

主入口：`POST /agent/chat`

职责：

- 用户自然语言对话
- 意图识别
- 多轮上下文管理
- 调用下游工具
- 统一输出回复
- 根据角色决定可用工具范围

主 Chat Agent 支持的主要意图：

- 公告解读 / 通知问答
- 物业服务咨询
- 报修创建
- 投诉创建
- 商品推荐
- 商品下单辅助
- 订单状态问答
- 物业费说明
- 报表解读
- 数据分析

### 4.2 报修 / 投诉能力

不再把 `/agent/repair-classify` 和 `/agent/complaint-risk` 视为独立产品入口，而是视为主 Chat Agent 的工具能力。

用途：

- 用户聊天中表达“我要报修”或“我要投诉”
- Agent 先做分类、紧急度、风险判断
- 再引导用户补全关键信息
- 最终调用业务服务创建工单

保留独立接口的原因：

- 便于前端单独调试
- 便于管理端后续单独接入
- 便于规则回归测试

### 4.3 购买辅助 Agent

用途：

- 基于用户需求推荐商品
- 帮助筛选门店
- 解释促销和库存情况
- 在用户确认后调用商城服务创建订单

注意：

- Agent 不直接支付
- 支付仍由现有商城支付流程完成
- Agent 只负责“辅助选购”和“预下单”

### 4.4 报告 / 数据分析能力

用途：

- 管理员、物业管理员、门店管理员查看日报 / 周报 / 最新统计摘要
- 默认优先读取“最近一份已生成报告”
- 只有当用户明确要求“实时查数据库”时，才触发在线聚合分析

这部分要区分：

- `report_read`: 读取已生成报告，低成本、快响应
- `report_generate`: 即时生成新报告，成本较高
- `db_analysis`: 实时查询统计接口或受控 SQL 分析，权限更高

## 5. 角色与权限边界

当前系统角色：

- `user`：普通用户
- `property`：物业管理员
- `store`：门店管理员
- `admin`：系统管理员

Agent 权限建议如下：

### 5.1 普通用户

允许：

- 公告 / 通知解读
- 社区服务问答
- 报修创建辅助
- 投诉创建辅助
- 商品推荐
- 下单辅助
- 订单、物业费、公告等“与本人相关”的状态问答

不允许：

- 查看全站运营统计
- 查看他人订单 / 他人工单 / 他人缴费
- 直接分析数据库
- 查看管理报告全文

### 5.2 物业管理员

允许：

- 普通用户全部非敏感问答能力
- 查看物业侧日报 / 周报
- 查询社区公告、访客、车位、物业费、工单统计
- 分析社区相关数据
- 读取和总结社区/工单数据

不允许：

- 商城经营数据深度分析
- 全局用户与权限数据分析

### 5.3 门店管理员

允许：

- 普通用户全部非敏感问答能力
- 查看商城经营报告
- 查询商品、门店、订单、销售排行、访客排行
- 分析商城相关数据

不允许：

- 全局社区治理数据分析
- 全局用户权限数据分析

### 5.4 系统管理员

允许：

- 全部 Agent 能力
- 默认读取最新综合报告
- 明确指定时可触发实时数据库分析
- 可跨商城、社区、工单、日志做综合分析

## 6. 工具设计

Agent 不直接查库，默认通过工具访问现有微服务。

工具分为四层：

### 6.1 查询型工具

- `get_my_profile`
- `get_my_orders`
- `get_my_property_fees`
- `get_latest_notices`
- `get_my_workorders`
- `get_store_products`
- `get_sales_rank`
- `get_view_rank`
- `get_community_overview`
- `get_latest_report`
- `get_report_detail`

### 6.2 执行型工具

- `create_workorder`
- `create_complaint_workorder`
- `create_cart_order`
- `recommend_products`
- `mark_notice_read`

### 6.3 分析型工具

- `analyze_notice_content`
- `classify_repair_request`
- `assess_complaint_risk`
- `analyze_property_operations`
- `analyze_store_operations`
- `analyze_admin_dashboard`

### 6.4 报告型工具

- `load_latest_cached_report`
- `generate_fresh_property_report`
- `generate_fresh_store_report`
- `generate_fresh_admin_report`

## 7. 主 Chat Agent 的执行流程

主 Chat Agent 建议采用“设计-规划-执行”三段式，而不是用户一句话就直接调用大模型回答。

### 7.1 设计阶段

输入：

- 用户消息
- 用户身份
- 最近若干轮会话
- 当前页面上下文

输出：

- 意图
- 是否需要工具
- 是否需要追问
- 是否涉及敏感数据
- 是否命中“默认看最新报告”

### 7.2 规划阶段

如果需要工具，则产出：

- 计划调用哪些工具
- 调用顺序
- 是否先读缓存报告
- 是否需要转写成业务结构化参数

### 7.3 执行阶段

- 调用工具
- 收集结果
- 输出自然语言答复
- 写入动作日志
- 更新用户记忆

## 8. 关键智能场景

### 8.1 聊天中创建报修 / 投诉

目标体验：

- 用户直接说“我家厨房水管漏水了”
- Agent 自动判断为报修
- 追问必要信息，例如楼栋、时间、图片、联系方式
- 调用分类工具判断紧急度
- 生成标准化工单参数
- 调用社区侧工单创建接口
- 返回创建结果

### 8.2 聊天中辅助购买

目标体验：

- 用户说“帮我看看今天有什么便宜点的水果”
- Agent 查询促销商品、库存、可配送门店
- 结合用户历史偏好推荐
- 用户确认后帮助生成订单

注意：

- 支付必须跳回商城支付流程
- Agent 不能绕过钱包密码 / 人脸 / 幂等校验

### 8.3 默认读最新报告

管理员 / 物业 / 门店管理员问：

- “今天情况怎么样”
- “帮我总结下最近运营情况”
- “看看最近社区问题”

默认动作：

- 先读最新报告
- 对报告做二次解释

只有明确出现以下语义时才做实时分析：

- “查最新数据库”
- “不要报告，直接看实时数据”
- “重新生成一份”
- “按今天最新订单统计”

### 8.4 实时数据库分析

实时分析不建议让 LLM 直接拼 SQL 访问主库。

推荐方式：

- Agent 调用受控统计工具
- 统计工具优先走已有 `/api/statistics/**`
- 确实没有现成统计接口时，再新增“白名单 SQL 分析工具”

白名单 SQL 分析工具约束：

- 仅 `admin` 可用，`property` / `store` 可在各自命名空间内用
- 只允许 `SELECT`
- 只允许访问白名单表或白名单视图
- 强制分页和超时
- 全量记录审计日志

## 9. 报告系统设计

### 9.1 报告类型

建议拆成三类：

- `admin_daily_report`
- `property_daily_report`
- `store_daily_report`

必要时再扩展：

- `weekly_summary_report`
- `exception_alert_report`

### 9.2 每天 9:00 自动任务

建议由 `agent-worker` 执行：

- 每天 09:00 生成全局管理员日报
- 每天 09:00 生成物业日报
- 每天 09:00 生成门店日报

如果门店很多，门店日报可按门店聚合分发。

### 9.3 报告生成流程

1. 从统计服务和业务服务取结构化数据
2. 生成统一 `report_context`
3. 构造不同角色的 Prompt
4. 调用大模型生成 Markdown
5. 保存报告摘要、全文、标签、统计快照
6. 写入 `cms_ai_report` 或新的 Agent 报表表

### 9.4 现有实现如何演进

当前 `community-service` 已有：

- `cms_ai_report`
- `GenerateReport`
- `GetLatestReport`
- `ListReports`

建议不要直接废弃，而是演进为：

- 社区服务保留“报告存储与展示”
- Agent 侧接管“报告生成策略、Prompt、角色区分、定时任务、解释能力”

也可以第二阶段再将报告能力完全迁入 `agent-service`。

## 10. 记忆系统设计

每个用户都应有自己的大模型记忆，但记忆必须分层。

### 10.1 短期记忆

存储位置：

- Redis

内容：

- 最近 10 到 20 轮对话
- 当前任务上下文
- 待确认槽位信息

用途：

- 多轮聊天连续性
- 追问补全
- 当前会话内上下文保持

### 10.2 中期记忆

存储位置：

- MySQL

内容：

- 会话摘要
- 用户偏好
- 常问问题
- 常见业务场景

例如：

- 喜欢买哪类商品
- 偏好哪个门店
- 经常关注物业费或访客问题

### 10.3 长期记忆

存储位置：

- MySQL
- 后续可加向量索引

内容：

- 稳定偏好
- 长期画像
- 重要事实

例如：

- 常住业主
- 常用车牌
- 常购物类
- 常关注的社区事项

### 10.4 记忆写入原则

- 用户显式确认的重要事实才进入长期记忆
- 普通闲聊不直接写长期记忆
- 每轮工具调用结果都不应原样塞入记忆
- 记忆写入必须先摘要再入库

## 11. RAG 预留设计

后续 RAG 不建议一开始就上复杂向量数据库。

第一阶段建议：

- 知识源放 MinIO / 本地挂载
- 文档元数据放 MySQL
- 先支持公告、制度文档、物业规则、商城帮助文档

知识类型建议：

- 社区公告
- 物业制度
- 报修 / 投诉处理规范
- 购物与支付说明
- 后台操作说明

后续可演进：

- 文档切片
- Embedding
- 向量检索
- Agent 与 RAG 混合回答

## 12. 建议的数据模型

在现有 `agent_sessions`、`agent_messages`、`agent_action_logs` 基础上，新增：

### `agent_sessions`

- `id`
- `user_id`
- `role_code`
- `agent_type`
- `source_page`
- `status`
- `last_message_at`
- `created_at`
- `updated_at`

### `agent_messages`

- `id`
- `session_id`
- `user_id`
- `role`
- `message_type`
- `content`
- `tool_calls_json`
- `token_usage_json`
- `created_at`

### `agent_action_logs`

- `id`
- `session_id`
- `user_id`
- `agent_type`
- `action`
- `tool_name`
- `request_payload`
- `response_payload`
- `success`
- `latency_ms`
- `created_at`

### `agent_memory_profiles`

- `id`
- `user_id`
- `profile_json`
- `summary_text`
- `updated_at`

### `agent_memory_items`

- `id`
- `user_id`
- `memory_type`
- `content`
- `importance`
- `source_session_id`
- `confirmed`
- `created_at`

### `agent_report_jobs`

- `id`
- `report_type`
- `scope_json`
- `status`
- `trigger_type`
- `started_at`
- `finished_at`
- `error_message`

### `agent_reports`

- `id`
- `report_type`
- `role_scope`
- `report_date`
- `title`
- `summary`
- `report_markdown`
- `metrics_snapshot_json`
- `generated_by`
- `created_at`

## 13. 接口设计建议

保留现有 4 个接口，但增加主入口和管理能力。

### 面向前端

- `GET /health`
- `POST /agent/chat`
- `POST /agent/chat/stream`
- `GET /agent/sessions`
- `GET /agent/sessions/{id}/messages`
- `POST /agent/repair-classify`
- `POST /agent/complaint-risk`
- `POST /agent/recommend`

### 面向管理端

- `GET /agent/reports/latest?type=admin_daily_report`
- `GET /agent/reports`
- `GET /agent/reports/{id}`
- `POST /agent/reports/generate`

### 面向内部调度

- `POST /internal/agent/reports/generate-daily`
- `POST /internal/agent/memory/compact`
- `POST /internal/agent/rag/reindex`

## 14. 权限控制建议

Agent 自身也要做权限校验，不能只靠前端隐藏按钮。

建议在 Agent 层增加三类权限：

- `agent:chat`
- `agent:report:read`
- `agent:report:generate`
- `agent:db_analysis`

并通过角色映射：

- `user` -> `agent:chat`
- `property` -> `agent:chat`, `agent:report:read`
- `store` -> `agent:chat`, `agent:report:read`
- `admin` -> 全部

## 15. 推荐实现阶段

### Phase 1

- 将 `agent-service` 从 FastAPI 占位升级为 `GoFrame + Eino`
- 实现 `/agent/chat`
- 实现会话、消息、动作日志落库
- 实现基础短期记忆
- 接入公告问答、报修创建、投诉创建、商品推荐

### Phase 2

- 接入 `/agent/chat/stream`
- 接入角色化报告读取
- 接入每天 9:00 自动报告
- 接入商城购买辅助

### Phase 3

- 接入中长期记忆
- 接入 RAG 文档检索
- 接入实时数据分析工具
- 加强 Prompt、评测和审计

## 16. 最终推荐结论

你的想法是对的，但建议做成下面这个形态：

- 用 **一个主 Chat Agent** 统一对话入口
- 报修 / 投诉 / 推荐 / 报告都作为 **工具能力或子 Agent**
- 普通用户只能做个人服务问答和通知解读
- `property`、`store`、`admin` 根据角色开放不同分析能力
- 管理角色默认读取最新报告，明确要求时才实时查库
- 每天 9:00 自动生成角色化报告
- 记忆分短期 / 中期 / 长期三层
- RAG 第二阶段接入，不和第一阶段耦死

从技术上，最适合你的落地组合是：

`GoFrame + Eino + Redis + MySQL + MinIO`

当前已有 `community-service` 报表与统计能力、群聊数据、RBAC 和现成微服务接口，完全足够支撑这一版 Agent 方案逐步落地。
