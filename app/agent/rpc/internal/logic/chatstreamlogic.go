package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"smartcommunity-microservices/app/agent/rpc/agent"
	"smartcommunity-microservices/app/agent/rpc/internal/model"
	"smartcommunity-microservices/app/agent/rpc/internal/svc"

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
	const maxPromptHistoryMessages = 8

	// 1. Check if conversation exists, or create it dynamically
	var conv model.SysUserConversation
	err := l.svcCtx.DB.Where("id = ? AND user_id = ?", in.ConversationId, in.UserId).First(&conv).Error
	if err != nil {
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

	// 2. Fetch only the prompt-relevant history window instead of loading all messages every turn.
	history, err := l.loadPromptHistory(in.ConversationId, in.UserId, conv.SummaryUntil, maxPromptHistoryMessages)
	if err != nil {
		l.Errorf("failed to load prompt history: %v", err)
		return err
	}

	// 3. Bind credentials/IDs to context for use inside Eino tools
	agentCtx := context.WithValue(l.ctx, CtxKeyUserID, in.UserId)
	agentCtx = context.WithValue(agentCtx, CtxKeyConversationID, in.ConversationId)
	if in.PayType != "" {
		agentCtx = context.WithValue(agentCtx, CtxKeyPayType, in.PayType)
		agentCtx = context.WithValue(agentCtx, CtxKeyPaymentPassword, in.PaymentPassword)
		agentCtx = context.WithValue(agentCtx, CtxKeyFaceImageURL, in.FaceImageUrl)
	}

	l.Infof("Agent Config status: globalKeyConfigured=%t, globalUrl=%q, globalModel=%q", l.svcCtx.Config.Agent.LlmApiKey != "", l.svcCtx.Config.Agent.LlmBaseUrl, l.svcCtx.Config.Agent.LlmModel)
	requestedMode := requestedChatModeFromContext(l.ctx)
	resolvedProfile := resolveChatProfile(requestedMode, in.Message)
	l.Infof("agent chat mode requested=%q resolved=%q", requestedMode, resolvedProfile)

	// Inject StreamCallback to context
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
