package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"smartcommunity-microservices/app/agent/rpc/internal/model"
	"smartcommunity-microservices/app/agent/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/google/uuid"
)

// CtxKey type for context values
type CtxKey string

const (
	CtxKeyUserID          CtxKey = "user_id"
	CtxKeyConversationID  CtxKey = "conversation_id"
	CtxKeyPayType         CtxKey = "pay_type"
	CtxKeyPaymentPassword CtxKey = "payment_password"
	CtxKeyFaceImageURL    CtxKey = "face_image_url"
	CtxKeyStreamCallback  CtxKey = "stream_callback"
)

func truncateText(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 || len([]rune(s)) <= n {
		return s
	}
	r := []rune(s)
	return string(r[:n]) + "..."
}

type StreamCallback func(eventType string, payload map[string]interface{})

func triggerToolStart(ctx context.Context, toolName string, args interface{}) {
	if cb, ok := ctx.Value(CtxKeyStreamCallback).(StreamCallback); ok {
		cb("tool_call_start", map[string]interface{}{
			"tool":         toolName,
			"args_preview": args,
		})
	}
}

func triggerToolEnd(ctx context.Context, toolName string, result interface{}) {
	if cb, ok := ctx.Value(CtxKeyStreamCallback).(StreamCallback); ok {
		cb("tool_call_end", map[string]interface{}{
			"tool":   toolName,
			"result": result,
		})
	}
}

func proposeAction(ctx context.Context, svcCtx *svc.ServiceContext, actionType string, input interface{}) (string, error) {
	userID, _ := ctx.Value(CtxKeyUserID).(int64)
	convID, _ := ctx.Value(CtxKeyConversationID).(string)

	if userID <= 0 || convID == "" {
		return "错误：未检测到有效登录用户或会话ID，请在登录后重试。", nil
	}

	payloadBytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Sprintf("序列化参数失败: %v", err), nil
	}

	actionID := "act-" + uuid.NewString()
	approval := &model.AgentActionApproval{
		ID:             actionID,
		ConversationID: convID,
		UserID:         userID,
		ActionType:     actionType,
		RiskLevel:      "high",
		ActionPayload:  string(payloadBytes),
		Status:         "pending",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if errDb := svcCtx.DB.Create(approval).Error; errDb != nil {
		return fmt.Sprintf("创建审批记录失败: %v", errDb), nil
	}

	// Trigger callbacks
	if cb, ok := ctx.Value(CtxKeyStreamCallback).(StreamCallback); ok {
		cb("approval_required", map[string]interface{}{
			"action_id":   actionID,
			"action_type": actionType,
			"payload":     input,
		})
	}

	return fmt.Sprintf("[APPROVAL_REQUIRED: %s]", actionID), nil
}

// 1. Get Current Time Tool
type GetCurrentTimeInput struct{}

func NewGetCurrentTimeTool() tool.InvokableTool {
	t, err := utils.InferTool(
		"get_current_time",
		"获取当前系统的日期和时间。在需要处理时间维度计算（如判断今天、昨天或当前时间）时调用本工具。",
		func(ctx context.Context, input *GetCurrentTimeInput) (output string, err error) {
			triggerToolStart(ctx, "get_current_time", nil)
			out := time.Now().Format("2006-01-02 15:04:05")
			triggerToolEnd(ctx, "get_current_time", out)
			return out, nil
		},
	)
	if err != nil {
		log.Fatalf("failed to create get_current_time tool: %v", err)
	}
	return t
}

// 2. Query Notices Tool
type QueryNoticesInput struct {
	Keyword string `json:"keyword" jsonschema:"description=可选，公告标题或正文中的简单关键字。若只是查看最新公告列表可留空；若需要按主题、语义、历史内容检索，请改用 search_knowledge"`
}

func NewQueryNoticesTool(svcCtx *svc.ServiceContext) tool.InvokableTool {
	t, err := utils.InferTool(
		"query_notices",
		"查询社区最新公告列表。适合“看看最近公告、最新通知、公告列表”这类直接拉取最新公告的场景。若用户是按主题、问题、历史内容、AI报告、制度说明做语义检索或总结，请不要使用本工具，应改用 search_knowledge。",
		func(ctx context.Context, input *QueryNoticesInput) (output string, err error) {
			triggerToolStart(ctx, "query_notices", input)
			resp, err := svcCtx.CommunityRpc.ListNotices(ctx, &communityrpc.ListNoticesReq{
				Page: 1,
				Size: 10,
			})
			if err != nil {
				out := fmt.Sprintf("查询公告失败: %v", err)
				triggerToolEnd(ctx, "query_notices", out)
				return out, nil
			}

			keyword := strings.TrimSpace(strings.ToLower(input.Keyword))
			var list []string
			for _, item := range resp.List {
				if keyword != "" {
					title := strings.ToLower(item.Title)
					content := strings.ToLower(item.Content)
					if !strings.Contains(title, keyword) && !strings.Contains(content, keyword) {
						continue
					}
				}
				list = append(list, fmt.Sprintf("- 标题: %s\n  内容摘要: %s\n  发布时间: %s", item.Title, truncateText(item.Content, 48), item.CreatedAt))
			}
			var out string
			if len(list) == 0 {
				if keyword != "" {
					out = "最新公告中暂未找到匹配该关键字的内容。如需按主题或历史内容进一步检索，请使用知识库检索。"
				} else {
					out = "当前社区暂无公告。"
				}
			} else {
				out = strings.Join(list, "\n\n")
			}
			triggerToolEnd(ctx, "query_notices", out)
			return out, nil
		},
	)
	if err != nil {
		log.Fatalf("failed to create query_notices tool: %v", err)
	}
	return t
}

type SearchKnowledgeInput struct {
	Query string `json:"query" jsonschema:"description=要检索的知识问题或主题，例如：最近社区停水通知、最近AI报告提到的主要问题、近期公告重点总结"`
	Scope string `json:"scope,omitempty" jsonschema:"description=检索范围，可选 notice、report、all。查AI报告可传 report；查公告语义内容可传 notice；留空默认全局检索"`
}

func NewSearchKnowledgeTool(svcCtx *svc.ServiceContext) tool.InvokableTool {
	t, err := utils.InferTool(
		"search_knowledge",
		"检索社区知识库。适用于按主题、问题、历史内容去语义化查询公告、制度说明、管理员AI报告等文档型内容。需要解释、总结、找依据、按主题查找时优先使用此工具。若用户只是想直接查看最新公告列表，请优先使用 query_notices。不要用本工具查询商品、订单、支付、库存、报修进度等实时业务数据。",
		func(ctx context.Context, input *SearchKnowledgeInput) (output string, err error) {
			triggerToolStart(ctx, "search_knowledge", input)

			if svcCtx.KnowledgeSvc == nil {
				out := "当前知识库检索尚未启用，请先配置 RAG 向量化能力。"
				triggerToolEnd(ctx, "search_knowledge", out)
				return out, nil
			}

			userID, _ := ctx.Value(CtxKeyUserID).(int64)
			hits, err := svcCtx.KnowledgeSvc.Search(ctx, userID, input.Query, input.Scope, 4)
			if err != nil {
				out := fmt.Sprintf("知识库检索失败: %v", err)
				triggerToolEnd(ctx, "search_knowledge", out)
				return out, nil
			}
			if len(hits) == 0 {
				out := "知识库中暂未找到相关内容。"
				triggerToolEnd(ctx, "search_knowledge", out)
				return out, nil
			}

			lines := make([]string, 0, len(hits))
			for _, hit := range hits {
				sourceLabel := "知识条目"
				switch hit.SourceType {
				case "notice":
					sourceLabel = "社区公告"
				case "ai_report":
					sourceLabel = "AI 报告"
				}

				snippet := hit.Content
				if strings.TrimSpace(hit.Summary) != "" {
					snippet = hit.Summary
				}

				lines = append(lines, fmt.Sprintf(
					"- 来源: %s\n  标题: %s\n  时间: %s\n  内容: %s",
					sourceLabel,
					hit.Title,
					hit.UpdatedAt.Format("2006-01-02 15:04:05"),
					truncateText(snippet, 180),
				))
			}

			out := strings.Join(lines, "\n\n")
			triggerToolEnd(ctx, "search_knowledge", out)
			return out, nil
		},
	)
	if err != nil {
		log.Fatalf("failed to create search_knowledge tool: %v", err)
	}
	return t
}

func isGenericKeyword(kw string) bool {
	kw = strings.TrimSpace(strings.ToLower(kw))
	if kw == "" {
		return true
	}
	// Common generic patterns
	generics := []string{
		"商品", "所有", "全部", "东西", "列表", "all", "product", "products", "stuff", "goods", "items", "shop", "store", "商铺", "商店", "超市", "便利店",
	}
	// If the keyword is purely composed of/contains only generic words, it's generic.
	temp := kw
	for _, g := range generics {
		temp = strings.ReplaceAll(temp, g, "")
	}
	if strings.TrimSpace(temp) == "" {
		return true
	}
	return false
}

// 3. List Products Tool
type ListProductsInput struct {
	Keyword string `json:"keyword" jsonschema:"description=要查询的商品名称或关键字，如：可乐、泡面"`
}

func NewListProductsTool(svcCtx *svc.ServiceContext) tool.InvokableTool {
	t, err := utils.InferTool(
		"list_products",
		"查询社区便利店的商品列表。大模型在用户询问商品、想买东西或打折促销时，应调用此工具查询商品名称、价格、库存和图片。",
		func(ctx context.Context, input *ListProductsInput) (output string, err error) {
			triggerToolStart(ctx, "list_products", input)

			searchKeyword := input.Keyword
			if isGenericKeyword(searchKeyword) {
				searchKeyword = ""
			}

			resp, err := svcCtx.MallRpc.ListProducts(ctx, &mall.ListProductsReq{
				Name: searchKeyword,
				Page: 1,
				Size: 8,
			})
			if err != nil {
				out := fmt.Sprintf("查询商品失败: %v", err)
				triggerToolEnd(ctx, "list_products", out)
				return out, nil
			}

			var list []string
			for _, p := range resp.List {
				priceYuan := float64(p.Price) / 100.0
				list = append(list, fmt.Sprintf("- 商品ID: %d\n  名称: %s\n  价格: ￥%.2f\n  库存: %d\n  描述: %s",
					p.Id, p.Name, priceYuan, p.Stock, truncateText(p.Description, 36)))
			}
			var out string
			if len(list) == 0 {
				if input.Keyword == "" || isGenericKeyword(input.Keyword) {
					out = "当前社区便利店暂无商品。"
				} else {
					out = fmt.Sprintf("未找到关键字为 '%s' 的商品。", input.Keyword)
				}
			} else {
				out = strings.Join(list, "\n\n")
			}
			triggerToolEnd(ctx, "list_products", out)
			return out, nil
		},
	)
	if err != nil {
		log.Fatalf("failed to create list_products tool: %v", err)
	}
	return t
}

// 4. Create Order Tool
type CreateOrderInput struct {
	ProductID int64 `json:"product_id" jsonschema:"description=要购买的商品的ID"`
	Quantity  int32 `json:"quantity" jsonschema:"description=购买的数量，必须大于等于1"`
}

func NewCreateOrderTool(svcCtx *svc.ServiceContext) tool.InvokableTool {
	t, err := utils.InferTool(
		"create_order",
		"为用户创建商城购买订单。调用此工具代表用户确认下单。该动作为高风险，必须先触发用户确认。确认后会生成订单ID用于后续支付。",
		func(ctx context.Context, input *CreateOrderInput) (output string, err error) {
			triggerToolStart(ctx, "create_order", input)
			out, err := proposeAction(ctx, svcCtx, "create_order", input)
			triggerToolEnd(ctx, "create_order", out)
			return out, err
		},
	)
	if err != nil {
		log.Fatalf("failed to create create_order tool: %v", err)
	}
	return t
}

// 5. Pay Order Tool
type PayOrderInput struct {
	OrderID int64 `json:"order_id" jsonschema:"description=需要付款支付的订单ID"`
}

func NewPayOrderTool(svcCtx *svc.ServiceContext) tool.InvokableTool {
	t, err := utils.InferTool(
		"pay_order",
		"对商城订单进行余额扣款支付。该动作为高风险，必须先触发用户付款授权确认卡片。不可随意静默扣款。",
		func(ctx context.Context, input *PayOrderInput) (output string, err error) {
			triggerToolStart(ctx, "pay_order", input)
			out, err := proposeAction(ctx, svcCtx, "pay_order", input)
			triggerToolEnd(ctx, "pay_order", out)
			return out, err
		},
	)
	if err != nil {
		log.Fatalf("failed to create pay_order tool: %v", err)
	}
	return t
}

// 6. Submit Repair Tool
type SubmitRepairInput struct {
	Type        string `json:"type" jsonschema:"description=服务单类别，只允许 'repair'(报修) 或 'complaint'(投诉)"`
	Category    string `json:"category" jsonschema:"description=具体的分类或投诉/报修原因小分类（由你根据用户的话来提炼或者用户直接给出的词，例如：水暖、电工、电梯、噪音、卫生、邻里纠纷等）"`
	Description string `json:"description" jsonschema:"description=业主对故障或投诉细节的具体文字描述"`
}

func NewSubmitRepairTool(svcCtx *svc.ServiceContext) tool.InvokableTool {
	t, err := utils.InferTool(
		"submit_repair",
		"创建物业报修工单或投诉工单。该动作为高风险状态变更，必须触发用户审批卡片确认内容后才提交物业工单系统。",
		func(ctx context.Context, input *SubmitRepairInput) (output string, err error) {
			triggerToolStart(ctx, "submit_repair", input)
			out, err := proposeAction(ctx, svcCtx, "submit_repair", input)
			triggerToolEnd(ctx, "submit_repair", out)
			return out, err
		},
	)
	if err != nil {
		log.Fatalf("failed to create submit_repair tool: %v", err)
	}
	return t
}

// 7. Query Order Status Tool
type QueryOrderStatusInput struct {
	OrderID int64 `json:"order_id" jsonschema:"description=要查询的订单ID，如果不指定或为0则默认获取最近创建的订单列表以检查其状态"`
}

func NewQueryOrderStatusTool(svcCtx *svc.ServiceContext) tool.InvokableTool {
	t, err := utils.InferTool(
		"query_order_status",
		"查询用户订单的最新支付、发货或退款状态。如果提供了特定的订单ID，则查询该订单的详细状态；若未指定ID，则查询自己最近的订单列表。",
		func(ctx context.Context, input *QueryOrderStatusInput) (output string, err error) {
			userID, _ := ctx.Value(CtxKeyUserID).(int64)
			if userID <= 0 {
				return "错误：未检测到有效登录用户，请在登录后重试。", nil
			}

			triggerToolStart(ctx, "query_order_status", input)

			if input.OrderID > 0 {
				resp, err := svcCtx.MallRpc.GetOrderDetail(ctx, &mall.OrderIDReq{
					Id:     input.OrderID,
					UserId: userID,
				})
				if err != nil {
					out := fmt.Sprintf("查询订单 #%d 详情失败: %v", input.OrderID, err)
					triggerToolEnd(ctx, "query_order_status", out)
					return out, nil
				}

				statusStr := map[int32]string{
					0:  "待支付 (Pending Payment)",
					1:  "待发货 / 已完成支付 (Paid / Pending Shipment)",
					2:  "待收货 / 已发货 (Shipped / Pending Receipt)",
					3:  "已完成 (Completed)",
					40: "已取消 (Cancelled)",
				}[resp.Status]
				if statusStr == "" {
					statusStr = fmt.Sprintf("未知状态 (%d)", resp.Status)
				}

				out := fmt.Sprintf("订单ID: %d\n订单号: %s\n商品数量: %d 件\n总金额: ￥%.2f\n当前状态: %s\n创建时间: %s",
					resp.Id, resp.OrderNo, len(resp.Items), float64(resp.TotalAmount)/100.0, statusStr, resp.CreatedAt)
				triggerToolEnd(ctx, "query_order_status", out)
				return out, nil
			}

			// Query order list for the user
			resp, err := svcCtx.MallRpc.ListOrders(ctx, &mall.ListOrdersReq{
				UserId: userID,
				Page:   1,
				Size:   5,
			})
			if err != nil {
				out := fmt.Sprintf("查询订单列表失败: %v", err)
				triggerToolEnd(ctx, "query_order_status", out)
				return out, nil
			}

			var list []string
			for _, o := range resp.List {
				statusStr := map[int32]string{
					0:  "待支付",
					1:  "待发货 / 已付款",
					2:  "待收货 / 已发货",
					3:  "已完成",
					40: "已取消",
				}[o.Status]
				if statusStr == "" {
					statusStr = fmt.Sprintf("未知 (%d)", o.Status)
				}
				list = append(list, fmt.Sprintf("- 订单ID: %d, 订单号: %s, 金额: ￥%.2f, 状态: %s, 创建时间: %s",
					o.Id, o.OrderNo, float64(o.TotalAmount)/100.0, statusStr, o.CreatedAt))
			}

			var out string
			if len(list) == 0 {
				out = "您当前没有任何商城订单。"
			} else {
				out = "您最近的订单列表如下：\n" + strings.Join(list, "\n")
			}
			triggerToolEnd(ctx, "query_order_status", out)
			return out, nil
		},
	)
	if err != nil {
		log.Fatalf("failed to create query_order_status tool: %v", err)
	}
	return t
}
