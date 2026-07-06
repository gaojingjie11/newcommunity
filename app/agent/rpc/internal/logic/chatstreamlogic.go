package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"smartcommunity-microservices/app/agent/rpc/agent"
	"smartcommunity-microservices/app/agent/rpc/internal/model"
	"smartcommunity-microservices/app/agent/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ChatStreamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatStreamLogic {
	return &ChatStreamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChatStreamLogic) ChatStream(in *agent.ChatReq, stream agent.AgentRpc_ChatStreamServer) error {
	const maxPromptHistoryMessages = 6

	if strings.TrimSpace(in.ConversationId) == "" {
		in.ConversationId = uuid.NewString()
	}

	// 1. Check if conversation exists, or create it dynamically
	var conv model.SysUserConversation
	err := l.svcCtx.DB.Where("id = ? AND user_id = ?", in.ConversationId, in.UserId).First(&conv).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Errorf("failed to load conversation: %v", err)
			return err
		}
		conv = model.SysUserConversation{
			ID:        in.ConversationId,
			UserID:    in.UserId,
			Title:     "新对话",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if errDb := l.svcCtx.DB.Create(&conv).Error; errDb != nil {
			l.Errorf("failed to auto-create conversation in DB: %v", errDb)
			return errDb
		}
	}

	// 2. Bind credentials/IDs to context for use inside Eino tools
	agentCtx := context.WithValue(l.ctx, CtxKeyUserID, in.UserId)
	agentCtx = context.WithValue(agentCtx, CtxKeyConversationID, in.ConversationId)
	if in.PayType != "" {
		agentCtx = context.WithValue(agentCtx, CtxKeyPayType, in.PayType)
		agentCtx = context.WithValue(agentCtx, CtxKeyPaymentPassword, in.PaymentPassword)
		agentCtx = context.WithValue(agentCtx, CtxKeyFaceImageURL, in.FaceImageUrl)
	}

	var lastApprovalPayload string
	streamCallback := func(eventType string, payload map[string]interface{}) {
		payloadBytes, errMarshal := json.Marshal(payload)
		if errMarshal != nil {
			return
		}
		if eventType == "approval_required" {
			lastApprovalPayload = string(payloadBytes)
		}
		_ = stream.Send(&agent.ChatResp{
			EventType:    eventType,
			EventPayload: string(payloadBytes),
		})
	}
	agentCtx = context.WithValue(agentCtx, CtxKeyStreamCallback, StreamCallback(streamCallback))

	if handled, reply, hasApprovalRequired, err := l.tryDirectResponse(agentCtx, in, stream); handled {
		if err != nil {
			return err
		}
		var errSave error
		if hasApprovalRequired || lastApprovalPayload != "" {
			errSave = l.saveChatMessagesTx(in.ConversationId, in.UserId, in.Message, reply, "approval_required", lastApprovalPayload)
		} else {
			errSave = l.saveChatMessagesTx(in.ConversationId, in.UserId, in.Message, reply, "", "")
		}
		if errSave != nil {
			l.Errorf("failed to save direct chat messages in transaction: %v", errSave)
			return errSave
		}
		return nil
	}

	// 3. Fetch only the prompt-relevant history window instead of loading all messages every turn.
	history, err := l.loadPromptHistory(in.ConversationId, in.UserId, conv.SummaryUntil, maxPromptHistoryMessages)
	if err != nil {
		l.Errorf("failed to load prompt history: %v", err)
		return err
	}

	l.Infof("Agent Config status: globalKeyConfigured=%t, globalUrl=%q, globalModel=%q", l.svcCtx.Config.Agent.LlmApiKey != "", l.svcCtx.Config.Agent.LlmBaseUrl, l.svcCtx.Config.Agent.LlmModel)
	requestedMode := requestedChatModeFromContext(l.ctx)
	resolvedProfile := resolveChatProfile(requestedMode, in.Message)
	l.Infof("agent chat mode requested=%q resolved=%q", requestedMode, resolvedProfile)

	// 5. Assemble Eino Prompt Messages
	var einoMessages []*schema.Message
	einoMessages = append(einoMessages, schema.SystemMessage(SystemPrompt))

	// Inject existing summary if present
	if conv.Summary != "" {
		einoMessages = append(einoMessages, schema.SystemMessage("历史对话摘要："+conv.Summary))
	}

	// Append active history turns since last summary
	for _, m := range history {
		if m.Role == "user" {
			einoMessages = append(einoMessages, schema.UserMessage(m.Content))
		} else if m.Role == "assistant" {
			einoMessages = append(einoMessages, schema.AssistantMessage(m.Content, nil))
		}
	}

	// Append current user query
	einoMessages = append(einoMessages, schema.UserMessage(in.Message))

	// 6. Invoke Eino Streaming Chat
	fullReply, hasApprovalRequired, err := l.runAgentStreamWithFallback(agentCtx, resolvedProfile, einoMessages, stream)
	if err != nil {
		return err
	}

	// 7. Save messages to database via Transaction
	var errSave error
	if lastApprovalPayload != "" {
		errSave = l.saveChatMessagesTx(in.ConversationId, in.UserId, in.Message, fullReply, "approval_required", lastApprovalPayload)
	} else {
		errSave = l.saveChatMessagesTx(in.ConversationId, in.UserId, in.Message, fullReply, "", "")
	}
	if errSave != nil {
		l.Errorf("failed to save chat messages in transaction: %v", errSave)
		return errSave
	}

	// 8. Auto-summarize old history turns (only if not approval required)
	if !hasApprovalRequired {
		go func(userID int64, convID string) {
			bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			l.autoSummarizeSession(bgCtx, convID, userID)
		}(in.UserId, in.ConversationId)
	}

	return nil
}

func (l *ChatStreamLogic) runAgentStreamWithFallback(
	agentCtx context.Context,
	resolvedProfile string,
	einoMessages []*schema.Message,
	stream agent.AgentRpc_ChatStreamServer,
) (string, bool, error) {
	reply, hasApprovalRequired, emittedChunks, err := l.runAgentStream(agentCtx, resolvedProfile, einoMessages, stream)
	if err == nil {
		return reply, hasApprovalRequired, nil
	}

	if !isMaxStepsError(err) || resolvedProfile == chatModeFast {
		return "", false, err
	}
	if emittedChunks > 0 {
		return "", false, fmt.Errorf("复杂模式推理未完整结束，请重试或切换简洁模式: %w", err)
	}

	l.Errorf("agent run exceeded max steps under profile=%q, falling back to fast profile: %v", resolvedProfile, err)
	reply, hasApprovalRequired, _, err = l.runAgentStream(agentCtx, chatModeFast, einoMessages, stream)
	return reply, hasApprovalRequired, err
}

func (l *ChatStreamLogic) runAgentStream(
	agentCtx context.Context,
	profile string,
	einoMessages []*schema.Message,
	stream agent.AgentRpc_ChatStreamServer,
) (string, bool, int, error) {
	runner, err := GetOrBuildAgent(agentCtx, l.svcCtx, profile)

	// Mock Fallback if LLM config is missing (for local testing/robustness)
	if err != nil {
		l.Errorf("LLM not configured (error: %v), running in mock fallback mode", err)
		reply, err := l.runMockFallback(&agent.ChatReq{
			ConversationId: anyValue[string](agentCtx, CtxKeyConversationID),
			UserId:         anyValue[int64](agentCtx, CtxKeyUserID),
			Message:        lastUserMessage(einoMessages),
		}, stream)
		return reply, false, len([]rune(reply)), err
	}

	sr, err := runner.Stream(agentCtx, einoMessages)
	if err != nil {
		l.Errorf("Eino Stream call failed under profile=%q: %v", profile, err)
		return "", false, 0, err
	}
	defer sr.Close()

	var fullReply strings.Builder
	hasApprovalRequired := false
	emittedChunks := 0

	for {
		chunk, errRecv := sr.Recv()
		if errors.Is(errRecv, io.EOF) {
			break
		}
		if errRecv != nil {
			l.Errorf("error reading Eino chunk under profile=%q: %v", profile, errRecv)
			return fullReply.String(), false, emittedChunks, errRecv
		}

		if strings.Contains(chunk.Content, "[APPROVAL_REQUIRED:") {
			hasApprovalRequired = true
			break
		}

		fullReply.WriteString(chunk.Content)

		errSend := stream.Send(&agent.ChatResp{
			Chunk:        chunk.Content,
			EventType:    "message_delta",
			EventPayload: fmt.Sprintf(`{"chunk":%q}`, chunk.Content),
		})
		if errSend != nil {
			l.Errorf("error writing gRPC stream chunk to gateway: %v", errSend)
			return fullReply.String(), false, emittedChunks, errSend
		}
		emittedChunks++
	}

	return fullReply.String(), hasApprovalRequired, emittedChunks, nil
}

func isMaxStepsError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "exceeds max steps")
}

func lastUserMessage(messages []*schema.Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i] == nil {
			continue
		}
		if messages[i].Role == schema.User {
			return messages[i].Content
		}
	}
	return ""
}

func anyValue[T any](ctx context.Context, key interface{}) T {
	var zero T
	value := ctx.Value(key)
	if value == nil {
		return zero
	}
	typed, ok := value.(T)
	if !ok {
		return zero
	}
	return typed
}

func (l *ChatStreamLogic) tryDirectResponse(
	agentCtx context.Context,
	in *agent.ChatReq,
	stream agent.AgentRpc_ChatStreamServer,
) (bool, string, bool, error) {
	message := strings.TrimSpace(in.Message)
	if message == "" {
		return false, "", false, nil
	}

	if shouldDirectLatestAIReport(message) {
		reply, err := l.directLatestAIReport(agentCtx, in.UserId)
		if err != nil {
			return true, "", false, err
		}
		if err := sendDirectReply(stream, reply); err != nil {
			return true, "", false, err
		}
		return true, reply, false, nil
	}

	if shouldDirectNoticeKnowledge(message) {
		reply, err := l.directNoticeKnowledge(agentCtx, in.UserId, message)
		if err != nil {
			return true, "", false, err
		}
		if err := sendDirectReply(stream, reply); err != nil {
			return true, "", false, err
		}
		return true, reply, false, nil
	}

	if shouldDirectLatestNotices(message) {
		reply, err := l.directLatestNotices(agentCtx, message)
		if err != nil {
			return true, "", false, err
		}
		if reply != "" {
			if err := sendDirectReply(stream, reply); err != nil {
				return true, "", false, err
			}
			return true, reply, false, nil
		}
	}

	if shouldDirectCreateOrder(message) {
		reply, hasApprovalRequired, err := l.directCreateOrder(agentCtx, message)
		if err != nil {
			return true, "", false, err
		}
		if err := sendDirectReply(stream, reply); err != nil {
			return true, "", false, err
		}
		return true, reply, hasApprovalRequired, nil
	}

	if shouldDirectProductIntent(message) || shouldDirectProductBrowse(message) {
		reply, err := l.directProductBrowse(agentCtx, message)
		if err != nil {
			return true, "", false, err
		}
		if err := sendDirectReply(stream, reply); err != nil {
			return true, "", false, err
		}
		return true, reply, false, nil
	}

	return false, "", false, nil
}

func (l *ChatStreamLogic) directLatestAIReport(ctx context.Context, userID int64) (string, error) {
	if l.svcCtx.KnowledgeSvc == nil {
		return "当前知识库检索尚未启用，暂时无法分析 AI 报告。", nil
	}

	allowed, err := l.svcCtx.KnowledgeSvc.CanAccessAdminKnowledge(ctx, userID)
	if err != nil {
		return "", err
	}
	if !allowed {
		return "您当前无权查看 AI 报告内容。", nil
	}

	report, err := l.svcCtx.KnowledgeSvc.GetLatestAIReport(ctx, userID)
	if err != nil {
		return "", err
	}
	if report == nil {
		return "当前暂无可用的 AI 报告。", nil
	}

	summary := strings.TrimSpace(report.Summary)
	if summary == "" {
		summary = strings.TrimSpace(report.Content)
	}
	summary = truncateText(summary, 320)

	return fmt.Sprintf(
		"我已为您检索最近一期 AI 报告。\n\n标题：%s\n时间：%s\n\n核心摘要：%s",
		report.Title,
		report.UpdatedAt.Format("2006-01-02 15:04:05"),
		summary,
	), nil
}

func (l *ChatStreamLogic) directLatestNotices(ctx context.Context, message string) (string, error) {
	keyword := extractNoticeKeyword(message)
	resp, err := l.svcCtx.CommunityRpc.ListNotices(ctx, &communityrpc.ListNoticesReq{
		Page: 1,
		Size: 10,
	})
	if err != nil {
		return "", err
	}

	lines := make([]string, 0, len(resp.List))
	keywordLower := strings.ToLower(keyword)
	for _, item := range resp.List {
		if keywordLower != "" {
			title := strings.ToLower(item.Title)
			content := strings.ToLower(item.Content)
			if !strings.Contains(title, keywordLower) && !strings.Contains(content, keywordLower) {
				continue
			}
		}
		lines = append(lines, fmt.Sprintf(
			"- 标题：%s\n  发布时间：%s\n  内容摘要：%s",
			item.Title,
			item.CreatedAt,
			truncateText(item.Content, 72),
		))
	}

	if len(lines) == 0 {
		if keyword != "" {
			return "最新公告中暂未找到匹配内容。如需按主题或历史内容深入检索，请继续说明具体问题。", nil
		}
		return "当前社区暂无公告。", nil
	}

	return "这是最近的社区公告：\n\n" + strings.Join(lines, "\n\n"), nil
}

func (l *ChatStreamLogic) directNoticeKnowledge(ctx context.Context, userID int64, message string) (string, error) {
	if l.svcCtx.KnowledgeSvc == nil {
		return "当前知识库检索尚未启用，暂时无法按主题检索公告。", nil
	}

	hits, err := l.svcCtx.KnowledgeSvc.Search(ctx, userID, message, "notice", 4)
	if err != nil {
		return "", err
	}
	if len(hits) == 0 {
		return "知识库中暂未找到相关公告内容。", nil
	}

	lines := make([]string, 0, len(hits))
	for _, hit := range hits {
		snippet := hit.Content
		if strings.TrimSpace(hit.Summary) != "" {
			snippet = hit.Summary
		}
		lines = append(lines, fmt.Sprintf(
			"- 标题：%s\n  时间：%s\n  内容：%s",
			hit.Title,
			hit.UpdatedAt.Format("2006-01-02 15:04:05"),
			truncateText(snippet, 120),
		))
	}
	return "我为您检索到以下相关公告：\n\n" + strings.Join(lines, "\n\n"), nil
}

func (l *ChatStreamLogic) directCreateOrder(ctx context.Context, message string) (string, bool, error) {
	keyword := extractOrderKeyword(message)
	if isGenericKeyword(keyword) {
		return "我可以直接帮您创建订单，请告诉我具体商品名称，例如“帮我买苹果1kg，直接下单”。", false, nil
	}

	var (
		resp *mall.ProductListResp
		err  error
	)
	for _, candidate := range buildProductSearchCandidates(keyword) {
		resp, err = l.svcCtx.MallRpc.ListProducts(ctx, &mall.ListProductsReq{
			Name: candidate,
			Page: 1,
			Size: 5,
		})
		if err != nil {
			return "", false, err
		}
		if len(resp.List) > 0 {
			keyword = candidate
			break
		}
	}
	if len(resp.List) == 0 {
		return fmt.Sprintf("未找到与“%s”相关的商品，您可以换个关键词再试。", keyword), false, nil
	}

	product := pickBestMatchedProduct(resp.List, keyword)
	if product == nil {
		product = resp.List[0]
	}
	if product == nil {
		return "当前未找到可下单商品。", false, nil
	}

	quantity := extractPurchaseQuantity(message)
	approvalText, err := proposeAction(ctx, l.svcCtx, "create_order", &CreateOrderInput{
		ProductID: product.Id,
		Quantity:  quantity,
	})
	if err != nil {
		return "", false, err
	}
	if strings.HasPrefix(approvalText, "错误：") {
		return approvalText, false, nil
	}

	reply := fmt.Sprintf(
		"已为您准备下单确认：\n\n- 商品：%s\n- 数量：%d\n- 单价：￥%.2f\n- 预计金额：￥%.2f\n\n请在前端确认创建订单。",
		product.Name,
		quantity,
		float64(product.Price)/100.0,
		float64(product.Price*int64(quantity))/100.0,
	)
	return reply, true, nil
}

func (l *ChatStreamLogic) directProductBrowse(ctx context.Context, message string) (string, error) {
	keyword := extractProductKeyword(message)
	if isGenericKeyword(keyword) && shouldDirectProductIntent(message) {
		keyword = extractOrderKeyword(message)
	}
	if isGenericKeyword(keyword) {
		keyword = ""
	}

	resp, err := l.svcCtx.MallRpc.ListProducts(ctx, &mall.ListProductsReq{
		Name: keyword,
		Page: 1,
		Size: 8,
	})
	if err != nil {
		return "", err
	}

	lines := make([]string, 0, len(resp.List))
	for _, p := range resp.List {
		lines = append(lines, fmt.Sprintf(
			"- 商品ID：%d\n  名称：%s\n  价格：￥%.2f\n  库存：%d\n  描述：%s",
			p.Id,
			p.Name,
			float64(p.Price)/100.0,
			p.Stock,
			truncateText(p.Description, 48),
		))
	}

	if len(lines) == 0 {
		if keyword == "" {
			return "当前社区便利店暂无商品。", nil
		}
		return fmt.Sprintf("未找到与“%s”相关的商品。", keyword), nil
	}

	prefix := "这是当前社区便利店可选的商品："
	if keyword != "" {
		prefix = fmt.Sprintf("我为您找到这些与“%s”相关的商品：", keyword)
	}
	reply := prefix + "\n\n" + strings.Join(lines, "\n\n")
	if keyword != "" {
		reply += fmt.Sprintf("\n\n如果您已经决定购买，直接回复“帮我买%s，直接下单”就行。", keyword)
	}
	return reply, nil
}

func sendDirectReply(stream agent.AgentRpc_ChatStreamServer, reply string) error {
	return stream.Send(&agent.ChatResp{
		Chunk:        reply,
		EventType:    "message_delta",
		EventPayload: fmt.Sprintf(`{"chunk":%q}`, reply),
	})
}

func shouldDirectLatestAIReport(message string) bool {
	text := strings.ToLower(strings.TrimSpace(message))
	if strings.Contains(text, "生成") && strings.Contains(text, "报告") {
		return false
	}
	if !(strings.Contains(text, "ai报告") || strings.Contains(text, "运营报表") || strings.Contains(text, "运营周报") || strings.Contains(text, "运营日报") || strings.Contains(text, "月报") || strings.Contains(text, "周报")) {
		return false
	}
	return strings.Contains(text, "最近") ||
		strings.Contains(text, "最新") ||
		strings.Contains(text, "最近一期") ||
		strings.Contains(text, "最近一份") ||
		strings.Contains(text, "最新一期") ||
		strings.Contains(text, "最新一份")
}

func shouldDirectLatestNotices(message string) bool {
	text := strings.ToLower(strings.TrimSpace(message))
	if strings.Contains(text, "ai报告") || strings.Contains(text, "运营报表") || strings.Contains(text, "报告") {
		return false
	}
	if strings.Contains(text, "公告列表") || strings.Contains(text, "最新公告") || strings.Contains(text, "最近公告") || strings.Contains(text, "最近通知") {
		return true
	}
	return strings.Contains(text, "最近有什么公告") || strings.Contains(text, "最近发了什么公告")
}

func shouldDirectNoticeKnowledge(message string) bool {
	text := strings.ToLower(strings.TrimSpace(message))
	if text == "" {
		return false
	}
	if strings.Contains(text, "ai报告") || strings.Contains(text, "运营报表") || strings.Contains(text, "周报") || strings.Contains(text, "月报") {
		return false
	}
	return (strings.Contains(text, "通知") || strings.Contains(text, "公告")) &&
		(strings.Contains(text, "有没有") || strings.Contains(text, "是否有") || strings.Contains(text, "查一下") || strings.Contains(text, "检索") || strings.Contains(text, "搜索"))
}

func shouldDirectProductBrowse(message string) bool {
	text := strings.ToLower(strings.TrimSpace(message))
	if text == "" {
		return false
	}
	if strings.Contains(text, "订单") || strings.Contains(text, "支付") || strings.Contains(text, "报修") || strings.Contains(text, "投诉") {
		return false
	}
	return strings.Contains(text, "便利店商品") ||
		strings.Contains(text, "推荐商品") ||
		strings.Contains(text, "推荐一些商品") ||
		strings.Contains(text, "商城有什么") ||
		strings.Contains(text, "有什么商品") ||
		strings.Contains(text, "看看商品")
}

func shouldDirectProductIntent(message string) bool {
	text := strings.ToLower(strings.TrimSpace(message))
	if text == "" {
		return false
	}
	if shouldDirectCreateOrder(text) {
		return false
	}
	if strings.Contains(text, "订单") || strings.Contains(text, "支付") || strings.Contains(text, "报修") || strings.Contains(text, "投诉") {
		return false
	}
	intentPhrases := []string{
		"我想买", "想买", "我要买", "买点", "来点", "有没有", "想看看", "想要",
	}
	for _, phrase := range intentPhrases {
		if strings.Contains(text, phrase) {
			keyword := extractOrderKeyword(message)
			return !isGenericKeyword(keyword)
		}
	}
	return false
}

func shouldDirectCreateOrder(message string) bool {
	text := strings.ToLower(strings.TrimSpace(message))
	if text == "" {
		return false
	}
	if strings.Contains(text, "支付") || strings.Contains(text, "订单状态") || strings.Contains(text, "取消订单") {
		return false
	}
	return strings.Contains(text, "直接下单") ||
		(strings.Contains(text, "帮我买") && !strings.Contains(text, "推荐")) ||
		strings.Contains(text, "帮我下单") ||
		(strings.Contains(text, "我要买") && strings.Contains(text, "下单"))
}

func extractNoticeKeyword(message string) string {
	text := strings.TrimSpace(message)
	replacers := []string{
		"请帮我", "帮我", "请", "看看", "查询", "检索", "最新", "最近", "社区", "公告", "通知", "一下", "有哪些", "有没有", "发了什么", "有什么",
	}
	for _, item := range replacers {
		text = strings.ReplaceAll(text, item, "")
	}
	return strings.TrimSpace(text)
}

func extractOrderKeyword(message string) string {
	text := strings.TrimSpace(message)
	replacers := []string{
		"帮我", "请", "直接", "下单", "购买", "买", "我要", "来一份", "来点", "给我", "安排", "一下",
	}
	for _, item := range replacers {
		text = strings.ReplaceAll(text, item, "")
	}
	return sanitizeProductKeyword(text)
}

func extractProductKeyword(message string) string {
	text := strings.TrimSpace(message)
	replacers := []string{
		"帮我", "请", "推荐", "一些", "看下", "看看", "查询", "搜索", "商品", "便利店", "商城", "有什么", "有啥", "推荐下", "一下",
		"我想买", "想买", "我要买", "买点", "来点", "想要", "想看看", "有没有",
	}
	for _, item := range replacers {
		text = strings.ReplaceAll(text, item, "")
	}
	return sanitizeProductKeyword(text)
}

func extractPurchaseQuantity(message string) int32 {
	text := strings.ToLower(strings.TrimSpace(message))
	patterns := []string{
		`(?:买|下单|来|要)\s*(\d+)\s*(?:份|件|个|袋|盒|瓶|箱)`,
		`x\s*(\d+)`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(text)
		if len(match) < 2 {
			continue
		}
		value, err := strconv.Atoi(match[1])
		if err != nil || value <= 0 {
			continue
		}
		return int32(value)
	}
	return 1
}

func pickBestMatchedProduct(list []*mall.ProductInfo, keyword string) *mall.ProductInfo {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		if len(list) == 0 {
			return nil
		}
		return list[0]
	}

	var best *mall.ProductInfo
	bestScore := -1
	for _, item := range list {
		if item == nil {
			continue
		}
		name := strings.ToLower(item.Name)
		score := 0
		if name == keyword {
			score += 4
		}
		if strings.Contains(name, keyword) {
			score += 3
		}
		for _, token := range strings.Fields(keyword) {
			if token != "" && strings.Contains(name, token) {
				score++
			}
		}
		if score > bestScore {
			bestScore = score
			best = item
		}
	}
	return best
}

func buildProductSearchCandidates(keyword string) []string {
	seen := make(map[string]bool)
	add := func(list []string, value string) []string {
		value = sanitizeProductKeyword(value)
		if value == "" || seen[value] {
			return list
		}
		seen[value] = true
		return append(list, value)
	}

	var candidates []string
	candidates = add(candidates, keyword)
	candidates = add(candidates, strings.ReplaceAll(keyword, "kg", " kg"))
	candidates = add(candidates, regexp.MustCompile(`\d+\s*(kg|公斤|g|克|ml|l|升|斤)$`).ReplaceAllString(keyword, ""))
	return candidates
}

func sanitizeProductKeyword(value string) string {
	replacer := strings.NewReplacer(
		"，", " ",
		",", " ",
		"。", " ",
		".", " ",
		"！", " ",
		"!", " ",
		"？", " ",
		"?", " ",
		"；", " ",
		";", " ",
		"：", " ",
		":", " ",
	)
	value = replacer.Replace(value)
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func (l *ChatStreamLogic) loadPromptHistory(convID string, userID int64, summaryUntil, maxPromptHistoryMessages int) ([]model.SysUserChatMessage, error) {
	if maxPromptHistoryMessages <= 0 {
		return nil, nil
	}

	var total int64
	if err := l.svcCtx.DB.Model(&model.SysUserChatMessage{}).
		Where("conversation_id = ? AND user_id = ?", convID, userID).
		Count(&total).Error; err != nil {
		return nil, err
	}

	startIdx := summaryUntil
	if startIdx < 0 {
		startIdx = 0
	}

	if totalInt := int(total); totalInt > startIdx+maxPromptHistoryMessages {
		startIdx = totalInt - maxPromptHistoryMessages
	}

	var history []model.SysUserChatMessage
	err := l.svcCtx.DB.
		Select("role", "content", "created_at").
		Where("conversation_id = ? AND user_id = ?", convID, userID).
		Order("created_at ASC").
		Offset(startIdx).
		Limit(maxPromptHistoryMessages).
		Find(&history).Error
	if err != nil {
		return nil, err
	}

	return history, nil
}

func (l *ChatStreamLogic) runMockFallback(in *agent.ChatReq, stream agent.AgentRpc_ChatStreamServer) (string, error) {
	reply := "您好！由于尚未配置大模型 API 密钥，当前正处于本地模拟演示模式。\n我可以为您演示社区公告总结、商城商品选购、订单支付以及报修工单创建："

	lowerMsg := strings.ToLower(in.Message)
	if strings.Contains(lowerMsg, "公告") || strings.Contains(lowerMsg, "通知") {
		reply = "【智能通告助手 - 模拟结果】\n最近社区有以下通知公告：\n1. 粽叶飘香：社区端午节包粽子与手工香囊制作活动将于本周五下午在居委会活动中心举办，欢迎报名。\n2. 关于近期雷雨及台风天气的安全出行提示，请勿在窗台外侧堆放杂物。\n3. 小区电梯机房定期例行安全维保通知。"
	} else if strings.Contains(lowerMsg, "可乐") || strings.Contains(lowerMsg, "泡面") || strings.Contains(lowerMsg, "商品") || strings.Contains(lowerMsg, "买") {
		reply = "【商城助手 - 模拟查询结果】\n为您找到以下商品：\n- 商品ID: 1\n  名称: 可口可乐罐装 330ml\n  价格: ￥3.00\n  库存: 99\n- 商品ID: 2\n  名称: 康师傅红烧牛肉面方便面\n  价格: ￥4.50\n  库存: 45\n\n您可以发送“模拟下单 商品ID”来下单体验。"
	} else if strings.Contains(lowerMsg, "下单") || strings.Contains(lowerMsg, "模拟下单") {
		reply = "已为您模拟创建订单！订单ID: 20261001，状态: 待付款，金额: ￥3.00。\n由于涉及资金交易，请在输入框左侧进行人脸识别或输入支付密码后，输入“确认支付 20261001”完成支付。"
	} else if strings.Contains(lowerMsg, "支付") || strings.Contains(lowerMsg, "付款") {
		points := 10
		reply = fmt.Sprintf("模拟支付成功！流水号: PAY991827361\n实付金额: ￥3.00\n订单状态已变更为: 待发货。\n绿色积分已为您自动抵扣并赠送了 %d 积分作为环保奖励！", points)
	} else if strings.Contains(lowerMsg, "报修") || strings.Contains(lowerMsg, "漏水") || strings.Contains(lowerMsg, "修") {
		reply = "已为您模拟提交物业报修单！\n工单类别：水暖修缮\n内容描述：业主反馈卫生间水龙头漏水\n工单号：WP20260610001\n系统已指派水工组张师傅预计于下午两点上门检修。"
	}

	runes := []rune(reply)
	for i := 0; i < len(runes); i++ {
		_ = stream.Send(&agent.ChatResp{
			Chunk:        string(runes[i]),
			EventType:    "message_delta",
			EventPayload: fmt.Sprintf(`{"chunk":%q}`, string(runes[i])),
		})
		time.Sleep(15 * time.Millisecond) // Simulated typing speed
	}

	return reply, nil
}

func (l *ChatStreamLogic) saveChatMessagesTx(convID string, userID int64, userMsg, botMsg, eventType, eventPayload string) error {
	return l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		userChat := &model.SysUserChatMessage{
			ID:             uuid.NewString(),
			UserID:         userID,
			ConversationID: convID,
			Role:           "user",
			Content:        userMsg,
			CreatedAt:      now,
		}
		if err := tx.Create(userChat).Error; err != nil {
			return err
		}

		botChat := &model.SysUserChatMessage{
			ID:             uuid.NewString(),
			UserID:         userID,
			ConversationID: convID,
			Role:           "assistant",
			Content:        botMsg,
			EventType:      eventType,
			EventPayload:   eventPayload,
			CreatedAt:      now.Add(1 * time.Millisecond),
		}
		if err := tx.Create(botChat).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.SysUserConversation{}).
			Where("id = ? AND user_id = ?", convID, userID).
			Update("updated_at", now).Error; err != nil {
			return err
		}

		return nil
	})
}

func (l *ChatStreamLogic) autoSummarizeSession(ctx context.Context, convID string, userID int64) {
	var conv model.SysUserConversation
	if err := l.svcCtx.DB.Where("id = ? AND user_id = ?", convID, userID).First(&conv).Error; err != nil {
		return
	}

	var messages []model.SysUserChatMessage
	l.svcCtx.DB.Where("conversation_id = ? AND user_id = ?", convID, userID).Order("created_at ASC").Find(&messages)

	// Summary compression trigger
	const maxMessages = 20
	const keepRecent = 6
	if len(messages) <= maxMessages {
		return
	}

	summaryEnd := len(messages) - keepRecent
	if summaryEnd <= conv.SummaryUntil {
		return
	}

	var toSummarize []string
	if conv.Summary != "" {
		toSummarize = append(toSummarize, "原有摘要: "+conv.Summary)
	}
	for _, msg := range messages[conv.SummaryUntil:summaryEnd] {
		toSummarize = append(toSummarize, fmt.Sprintf("%s: %s", msg.Role, msg.Content))
	}

	summaryPrompt := fmt.Sprintf("请简短总结以下对话内容，着重归纳用户的要求、创建的订单号、报修的类型及结论。总结应在150字以内，保持客观精炼：\n\n%s", strings.Join(toSummarize, "\n"))

	cfg := l.svcCtx.Config.Agent
	apiKey, baseUrl, modelName := cfg.GetModelConfig(cfg.Models.ChatDefault)
	if apiKey == "" || baseUrl == "" {
		return
	}

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   modelName,
		APIKey:  apiKey,
		BaseURL: baseUrl,
	})
	if err != nil {
		return
	}

	resp, err := chatModel.Generate(ctx, []*schema.Message{
		schema.UserMessage(summaryPrompt),
	})
	if err != nil {
		return
	}

	l.svcCtx.DB.Model(&conv).Updates(map[string]interface{}{
		"summary":       resp.Content,
		"summary_until": summaryEnd,
		"updated_at":    time.Now(),
	})
}
