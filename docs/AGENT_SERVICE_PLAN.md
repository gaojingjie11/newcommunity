# Agent Service Plan

当前仓库中的 `agent-service` 仍是 FastAPI 占位实现，只提供 4 个占位接口。

下一阶段的目标架构不再建议继续以“FastAPI 占位 + 零散功能追加”的方式演进，而是升级为：

`GoFrame + Eino`

完整方案见：

- [AGENT_ARCHITECTURE_V2.md](/Users/gao/Downloads/shequfuwu/docs/AGENT_ARCHITECTURE_V2.md)

## Requirement Governance

- `需求陈述书-东软智慧社区项目.docx` is the functional baseline.
- Agent capabilities are currently planned as extensions and must not replace mandatory Word flows such as registration, order payment, visitor approval, repair handling, complaint handling, or property fee payment.
- Every Agent-related implementation must reference `AGENT-*` rows in `docs/REQUIREMENTS_TRACEABILITY_MATRIX.md`.
- If an Agent proposes or automates a business action, the underlying business service must still satisfy its own `AUTH-*`, `MALL-*`, `COMM-*`, `ADMIN-*`, `STAT-*`, or `LOG-*` requirement ID.
- If a feature has no traceability matrix ID, do not develop it directly.

## Community Customer Service Agent

Endpoint: `POST /agent/chat`

Boundary: identify user intent, answer common community service questions, and call Go services when needed.

Target evolution:

- upgrade this endpoint into the main Chat Agent entry
- support role-aware tool calling
- support repair/complaint creation from chat
- support default latest-report reading for management roles
- support per-user memory

## Repair Dispatch Agent

Endpoint: `POST /agent/repair-classify`

Boundary: classify repair text/images, infer urgency, and suggest maintenance department.

Target evolution:

- keep as an independently testable tool endpoint
- use internally by the main Chat Agent before creating a repair workorder

## Complaint Risk Agent

Endpoint: `POST /agent/complaint-risk`

Boundary: detect complaint category and escalation risk, recommend follow-up action.

Target evolution:

- keep as an independently testable tool endpoint
- use internally by the main Chat Agent before creating a complaint workorder

## Recommendation Agent

Endpoint: `POST /agent/recommend`

Boundary: recommend products or community services from user context and service data.

Target evolution:

- support mall recommendation and service recommendation
- support order-assist scenarios
- remain unable to bypass the normal payment workflow

## Reporting And Analysis

Target additions:

- daily scheduled reports at `09:00`
- role-specific reports for `admin`, `property`, and `store`
- default behavior for management questions should read the latest generated report first
- real-time database analysis should only run when the user explicitly asks for fresh live data
- report generation and database analysis must be permission-gated

## Memory Strategy

Target memory layers:

- short-term session memory in Redis
- medium-term summarized memory in MySQL
- long-term user preference/profile memory with future RAG integration

Memory must be user-scoped and role-aware, and should never directly persist raw tool payloads as long-term memory.

## Next Integration Points

- Replace the FastAPI placeholder with a GoFrame service.
- Introduce Eino for Chat Agent orchestration, tool calling, and future RAG.
- Add LLM provider client using `LLM_API_KEY`, `LLM_BASE_URL`, and `LLM_MODEL`.
- Add typed clients for gateway/user/mall/community/workorder/statistics services.
- Add session/message/action-log persistence.
- Add scheduled report generation worker.
- Add traceable prompt templates, guardrails, and deterministic fallback rules.
