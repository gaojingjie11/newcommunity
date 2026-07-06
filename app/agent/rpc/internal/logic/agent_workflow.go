package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"smartcommunity-microservices/app/agent/rpc/internal/svc"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"google.golang.org/grpc/metadata"
)

const SystemPrompt = `你是一个非常专业、友好、高效的智慧社区管家（AI助手）。
你可以通过调用底层工具来直接帮用户处理以下事务：
1. 查询社区公告通知 (使用 query_notices 工具)
2. 检索知识库中的公告、AI报告和说明文档 (使用 search_knowledge 工具)
3. 浏览/搜索社区商城便民商店的商品 (使用 list_products 工具)
4. 帮用户在便民商店里直接下单 (使用 create_order 工具)
5. 帮用户对已创建好的订单进行钱包余额支付 (使用 pay_order 工具)
6. 提交物业报修单或投诉建议 (使用 submit_repair 工具)
7. 查询用户已下单的商品订单的最新状态和详情 (使用 query_order_status 工具)

【工具选择规则（非常重要）】：
- 当用户只是想看“最新公告、最近通知、公告列表、最近发了什么”，优先使用 query_notices。它适合直接拉取最新公告列表，成本更低、速度更快。
- 当用户想按主题、语义、问题去检索文档内容时，使用 search_knowledge。例如：“最近有没有停水通知”“AI报告提到了哪些风险”“帮我总结近期公告重点”。它适合做语义检索、归纳总结、找依据。
- 当问题涉及 AI 报告、制度说明、历史文档归纳、跨多篇公告找答案时，优先使用 search_knowledge，不要只用 query_notices。
- 当用户明确要求“只查 AI 报告”“不要查公告”“基于 AI 报告回答”时，如果没有可访问的 AI 报告内容，就直接明确说明无权限或未找到，不要回退去分析公告。
- 当问题涉及商品、订单、支付、报修、投诉、用户自己的实时状态时，不要使用 search_knowledge，应直接使用对应业务工具。
- 不要用 search_knowledge 去查询实时订单状态、商品库存、支付结果、钱包余额等强实时业务数据。
- 如果用户既要“看看最近公告”，又要“总结重点/查某个主题”，可以先用 query_notices 看最新列表，再按需要使用 search_knowledge 做补充。

【报修与投诉服务单提交要求】：
- 对于提交物业报修单或投诉建议（submit_repair）：一旦确定了类型（报修还是投诉）、提炼出具体的故障或投诉分类（由你根据用户的话自动提炼简短的分类词，例如：水暖、电工、电梯、噪音、卫生、服务态度、邻里纠纷等），并获得基本的文本描述（如：“楼上晚上太吵了”、“二楼路灯坏了”），你就应该【立即】调用 submit_repair 工具提交审批。
- **不要**向用户索要多余的楼栋号、时间段、具体噪音类型等细节信息！不需要过于死板，应以最快的速度帮用户触发工单审批卡片，让用户在弹出的工单卡片中确认内容即可。

【关键交易支付指南（非常重要）】：
- 用户的付款验证凭证（如支付密码、人脸图片）会自动绑定在后台上下文中。
- 当用户表达想要购买某件商品时，你应该：
  1. 先使用 list_products 查询该商品，并告知用户商品信息及价格。
  2. 确认用户购买后，使用 create_order 生成订单，并获得订单ID（OrderID）。
  3. 大模型在下单成功后，必须返回包含订单ID的响应，并提示用户输入支付密码或进行人脸扫描验证以完成支付扣款。
  4. 只有当用户输入了密码/完成了人脸验证（上下文的 pay_type 不为空），且你拥有订单ID时，你才能调用 pay_order 工具进行支付。不要在大门紧锁（没有支付参数）时空跑支付接口。

【商品ID展示与匹配要求（极度重要）】：
- 当你调用 list_products 检索商品并在聊天回复中列给用户看时，你必须显式、完整地保留并展示商品真实的“商品ID: <id>”（例如 “商品ID: 5”），并明确告知用户可以基于这个 ID 或名称进行购买。
- 绝对不要省略商品 ID，也绝对不能混淆为列表序号（如 1.、2. 等）！在调用 create_order 下单时，Items 数组中的 product_id 必须使用商品的真实数据库 ID，不能填错。

【通用要求】：
- 当前系统时间，请使用 get_current_time 工具获取。
- 无论进行任何查询（公告、商品等），若工具返回了多条记录，请以清晰的 Markdown 列表整理，价格显示为两位小数的元（例如 ￥3.00）。
- 回答要礼貌简洁，输出格式美观，必要时换行。不需要输出 Markdown 标记，直接输出文字。`

const (
	chatModeAuto = "auto"
	chatModeFast = "fast"
	chatModeDeep = "deep"

	agentModeMetadataKey = "x-agent-mode"
)

var (
	agentInstances = map[string]*react.Agent{}
	agentMu        sync.Mutex
)

func GetOrBuildAgent(ctx context.Context, svcCtx *svc.ServiceContext, profile string) (*react.Agent, error) {
	agentMu.Lock()
	defer agentMu.Unlock()

	profile = normalizeChatMode(profile)
	if profile == chatModeAuto {
		profile = chatModeFast
	}
	if agent, ok := agentInstances[profile]; ok && agent != nil {
		return agent, nil
	}

	agent, err := BuildEinoAgent(ctx, svcCtx, profile)
	if err != nil {
		return nil, err
	}
	agentInstances[profile] = agent
	return agent, nil
}

func BuildEinoAgent(ctx context.Context, svcCtx *svc.ServiceContext, profile string) (*react.Agent, error) {
	cfg := svcCtx.Config.Agent
	profile = normalizeChatMode(profile)

	apiKey, baseUrl, modelName := cfg.GetModelConfig(cfg.Models.ChatDefault)
	if profile == chatModeDeep {
		apiKey, baseUrl, modelName = cfg.GetModelConfig(cfg.Models.AgentReasoning)
	} else if apiKey == "" || baseUrl == "" || strings.TrimSpace(modelName) == "" {
		apiKey, baseUrl, modelName = cfg.GetModelConfig(cfg.Models.AgentReasoning)
	}
	if apiKey == "" || baseUrl == "" {
		return nil, errors.New("LLM API key or Base URL is not configured for chat agent")
	}

	// 1. Initialize OpenAI-compatible Chat Model
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   modelName,
		APIKey:  apiKey,
		BaseURL: baseUrl,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Eino OpenAI chat model: %w", err)
	}

	// 2. Load all tools
	tools := []tool.BaseTool{
		NewGetCurrentTimeTool(),
		NewQueryNoticesTool(svcCtx),
		NewSearchKnowledgeTool(svcCtx),
		NewListProductsTool(svcCtx),
		NewCreateOrderTool(svcCtx),
		NewPayOrderTool(svcCtx),
		NewSubmitRepairTool(svcCtx),
		NewQueryOrderStatusTool(svcCtx),
	}

	// 3. Instantiate ReAct Agent
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		MaxStep:          maxAgentSteps(profile),
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to compile Eino ReAct agent: %w", err)
	}

	return agent, nil
}

func maxAgentSteps(profile string) int {
	switch normalizeChatMode(profile) {
	case chatModeDeep:
		return 8
	case chatModeFast:
		return 4
	default:
		return 6
	}
}

func requestedChatModeFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return chatModeAuto
	}
	values := md.Get(agentModeMetadataKey)
	if len(values) == 0 {
		return chatModeAuto
	}
	return normalizeChatMode(values[0])
}

func resolveChatProfile(requestedMode, message string) string {
	// Backend retains final decision authority and overrides frontend preference (e.g. fast)
	// for specific heavy scenarios requiring advanced reasoning capacity.
	if isScenarioRequiringDeep(message) {
		return chatModeDeep
	}

	switch normalizeChatMode(requestedMode) {
	case chatModeFast:
		return chatModeFast
	case chatModeDeep:
		return chatModeDeep
	default:
		if isComplexAgentQuestion(message) {
			return chatModeDeep
		}
		return chatModeFast
	}
}

func isScenarioRequiringDeep(message string) bool {
	text := strings.TrimSpace(strings.ToLower(message))

	// Enforce deep model for specific complex/high-risk scenarios:
	deepKeywords := []string{
		// 1. 生成日报
		"日报", "周报", "月报", "工作周报", "生成报告", "数据报告", "daily report", "weekly report", "report generation",
		// 2. 实时数据库分析
		"数据库", "数据分析", "实时分析", "运营分析", "统计数据", "运营数据", "database", "db analysis", "query db", "sql query",
		// 3. 多步审批推理
		"审批", "多步审批", "审批流", "授权确认", "审批推理", "approval", "approve action", "multi-step",
		// 4. 高风险动作解释
		"高风险", "风险解释", "为什么购买", "为什么付款", "为什么扣款", "扣款原因", "安全验证", "risk", "high-risk", "security check",
		// 5. 跨服务复杂总结
		"跨服务", "服务总结", "综合总结", "汇总数据", "多系统", "cross-service", "summary", "summarize",
	}

	for _, kw := range deepKeywords {
		if strings.Contains(text, kw) {
			return true
		}
	}

	// Override if message length is very long (indicates complex input requiring strong contextual summary)
	if len([]rune(text)) >= 220 {
		return true
	}

	return false
}

func normalizeChatMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case chatModeFast:
		return chatModeFast
	case chatModeDeep:
		return chatModeDeep
	default:
		return chatModeAuto
	}
}

func isComplexAgentQuestion(message string) bool {
	text := strings.TrimSpace(strings.ToLower(message))
	if len([]rune(text)) >= 120 {
		return true
	}

	complexKeywords := []string{
		"分析", "总结", "报告", "报表", "统计", "趋势", "原因", "建议",
		"对比", "方案", "优化", "数据库", "近7天", "近30天", "本周", "本月",
		"同比", "环比", "风险", "异常", "复盘", "analyze", "analysis", "report",
	}
	for _, keyword := range complexKeywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}
