package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"smartcommunity/internal/config"
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	qwenVisionModel = "qwen-vl-plus"

	chatContextWindowSize = 10
	maxChatHistoryLimit   = 200
	maxAgentToolRounds    = 5
	facePayMinConfidence  = 85.0

	toolGetRecentNotices   = "get_recent_notices"
	toolCreateRepairTicket = "create_repair_ticket"
	toolSearchProduct      = "search_product"
	toolCreateOrder        = "create_order"
	toolPayOrder           = "pay_order"
)

const communityAssistantSystemPrompt = "You are a smart community assistant. Use tools for real data operations. " +
	"For order payment, auth can be password or face, and payment must be explicitly confirmed by user. " +
	"Never claim order creation/payment success unless backend tool execution succeeds. " +
	"Never claim repair/complaint ticket creation success unless tool result has success=true and repair_id. " +
	"Complaint submission is supported via create_repair_ticket with type=complaint."

type AIService struct{}

type chatExecutionContext struct {
	paymentPassword string
	payType         string
	faceImageURL    string
}

type dashScopeMessage struct {
	Role       string              `json:"role"`
	Content    interface{}         `json:"content"`
	Name       string              `json:"name,omitempty"`
	ToolCalls  []dashScopeToolCall `json:"tool_calls,omitempty"`
	ToolCallID string              `json:"tool_call_id,omitempty"`
}

type dashScopeTool struct {
	Type     string                    `json:"type"`
	Function dashScopeToolFunctionSpec `json:"function"`
}

type dashScopeToolFunctionSpec struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type dashScopeToolCall struct {
	ID       string                    `json:"id"`
	Type     string                    `json:"type"`
	Function dashScopeToolCallFunction `json:"function"`
}

type dashScopeToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type dashScopeRequest struct {
	Model          string             `json:"model"`
	Messages       []dashScopeMessage `json:"messages"`
	ResponseFormat map[string]string  `json:"response_format,omitempty"`
	Tools          []dashScopeTool    `json:"tools,omitempty"`
	ToolChoice     interface{}        `json:"tool_choice,omitempty"`
}

type dashScopeResponse struct {
	Choices []struct {
		Message dashScopeMessage `json:"message"`
	} `json:"choices"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Error   *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type GarbageRecognitionResult struct {
	Points int    `json:"points"`
	Reason string `json:"reason"`
}

type getRecentNoticesArgs struct {
	Limit int `json:"limit"`
}

type createRepairTicketArgs struct {
	Type        string `json:"type"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type searchProductArgs struct {
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit"`
}

type createOrderItemArg struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type createOrderArgs struct {
	Items   []createOrderItemArg `json:"items"`
	StoreID int64                `json:"store_id"`
}

type payOrderArgs struct {
	OrderIDRaw json.RawMessage `json:"order_id"`
	OrderNo    string          `json:"order_no"`
	Password   string          `json:"password"`
	PayType    string          `json:"pay_type"`
	FaceImage  string          `json:"face_image_url"`
}

func (s *AIService) Chat(userContent string) (string, error) {
	messages := []dashScopeMessage{
		{Role: "system", Content: communityAssistantSystemPrompt},
		{Role: "user", Content: userContent},
	}
	return s.callTextModel(config.Conf.AI.Model, messages, nil)
}

func (s *AIService) ChatWithMemory(userID int64, userContent string, paymentPassword string, payType string, faceImageURL string) (string, error) {
	userContent = strings.TrimSpace(userContent)
	if userID <= 0 {
		return "", errors.New("invalid user id")
	}
	if userContent == "" {
		return "", errors.New("content cannot be empty")
	}

	userMsg := model.ChatMessage{UserID: userID, Role: "user", Content: userContent}
	if err := global.DB.Create(&userMsg).Error; err != nil {
		return "", errors.New("failed to save user message")
	}

	var recent []model.ChatMessage
	if err := global.DB.Where("user_id = ?", userID).
		Order("id DESC").Limit(chatContextWindowSize).Find(&recent).Error; err != nil {
		return "", errors.New("failed to load chat context")
	}
	reverseChatMessages(recent)

	messages := make([]dashScopeMessage, 0, len(recent)+1)
	messages = append(messages, dashScopeMessage{Role: "system", Content: communityAssistantSystemPrompt})
	for _, m := range recent {
		if m.Role == "user" || m.Role == "assistant" {
			messages = append(messages, dashScopeMessage{Role: m.Role, Content: m.Content})
		}
	}

	reply, err := s.runAgentWithTools(userID, chatExecutionContext{
		paymentPassword: strings.TrimSpace(paymentPassword),
		payType:         strings.ToLower(strings.TrimSpace(payType)),
		faceImageURL:    strings.TrimSpace(faceImageURL),
	}, messages)
	if err != nil {
		return "", err
	}

	assistantMsg := model.ChatMessage{UserID: userID, Role: "assistant", Content: reply}
	if err := global.DB.Create(&assistantMsg).Error; err != nil {
		return "", errors.New("failed to save assistant message")
	}
	return reply, nil
}

func (s *AIService) runAgentWithTools(userID int64, execCtx chatExecutionContext, messages []dashScopeMessage) (string, error) {
	tools := buildCommunityAgentTools()
	working := append([]dashScopeMessage(nil), messages...)
	lastUser := latestUserMessage(messages)

	// Deterministic path: create repair/complaint ticket for explicit requests.
	if reply, handled, err := s.tryHandleRepairIntent(userID, lastUser); handled {
		return reply, err
	}

	// Deterministic path: create order directly for simple "下单X瓶" instructions.
	if reply, handled, err := s.tryHandleCreateOrderIntent(userID, messages, lastUser); handled {
		return reply, err
	}
	// Deterministic path: pay latest pending order when user confirms payment / enters password.
	if reply, handled, err := s.tryHandlePayIntent(userID, execCtx, lastUser); handled {
		return reply, err
	}

	if forced := buildForcedSearchToolCall(lastUser); forced != nil {
		working = append(working, dashScopeMessage{
			Role:      "assistant",
			Content:   "",
			ToolCalls: []dashScopeToolCall{*forced},
		})
		if err := s.executeAndAppendToolResult(&working, userID, execCtx, *forced); err != nil {
			return "", err
		}
	}

	for round := 0; round < maxAgentToolRounds; round++ {
		assistantMsg, err := s.callModel(config.Conf.AI.Model, working, nil, tools)
		if err != nil {
			return "", err
		}

		if len(assistantMsg.ToolCalls) == 0 {
			finalReply := strings.TrimSpace(extractMessageContent(assistantMsg.Content))
			if finalReply == "" {
				return "", errors.New("empty AI response")
			}
			return finalReply, nil
		}

		content := assistantMsg.Content
		if content == nil {
			content = ""
		}
		working = append(working, dashScopeMessage{Role: "assistant", Content: content, ToolCalls: assistantMsg.ToolCalls})

		for _, call := range assistantMsg.ToolCalls {
			if err := s.executeAndAppendToolResult(&working, userID, execCtx, call); err != nil {
				return "", err
			}
		}
	}
	return "", errors.New("tool calling exceeded max rounds")
}

func (s *AIService) executeAndAppendToolResult(working *[]dashScopeMessage, userID int64, execCtx chatExecutionContext, toolCall dashScopeToolCall) error {
	result, err := s.dispatchToolCall(userID, execCtx, toolCall)
	if err != nil {
		log.Printf("tool execution failed, tool=%s userID=%d err=%v", toolCall.Function.Name, userID, err)
		if isCriticalTool(toolCall.Function.Name) {
			return formatCriticalToolError(toolCall.Function.Name, err)
		}
		result = map[string]interface{}{"success": false, "error": err.Error()}
	}

	payload, _ := json.Marshal(result)
	*working = append(*working, dashScopeMessage{
		Role:       "tool",
		Name:       toolCall.Function.Name,
		ToolCallID: toolCall.ID,
		Content:    string(payload),
	})
	return nil
}

func (s *AIService) dispatchToolCall(userID int64, execCtx chatExecutionContext, toolCall dashScopeToolCall) (interface{}, error) {
	switch toolCall.Function.Name {
	case toolGetRecentNotices:
		var args getRecentNoticesArgs
		if err := parseToolArguments(toolCall.Function.Arguments, &args); err != nil {
			return nil, err
		}
		if args.Limit <= 0 {
			args.Limit = 5
		}
		if args.Limit > 20 {
			args.Limit = 20
		}
		list, err := (&NoticeService{}).GetList(args.Limit)
		if err != nil {
			return nil, err
		}
		items := make([]map[string]interface{}, 0, len(list))
		for _, n := range list {
			items = append(items, map[string]interface{}{
				"id": n.ID, "title": n.Title, "content": n.Content, "publisher": n.Publisher, "created_at": n.CreatedAt,
			})
		}
		return map[string]interface{}{"success": true, "count": len(items), "items": items}, nil

	case toolCreateRepairTicket:
		var args createRepairTicketArgs
		if err := parseToolArguments(toolCall.Function.Arguments, &args); err != nil {
			return nil, err
		}
		description := strings.TrimSpace(args.Description)
		if description == "" {
			return nil, errors.New("description cannot be empty")
		}
		repairType, ticketTypeLabel, _ := classifyRepairIntent(args.Type, description)
		category := normalizeRepairCategory(args.Category, args.Type, description)
		repair, err := s.createAndVerifyRepairTicket(userID, repairType, category, description)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"success":     true,
			"repair_id":   repair.ID,
			"status":      "pending",
			"ticket_type": ticketTypeLabel,
			"category":    category,
		}, nil

	case toolSearchProduct:
		var args searchProductArgs
		if err := parseToolArguments(toolCall.Function.Arguments, &args); err != nil {
			return nil, err
		}
		args.Keyword = strings.TrimSpace(args.Keyword)
		if args.Keyword == "" {
			return nil, errors.New("keyword cannot be empty")
		}
		if args.Limit <= 0 {
			args.Limit = 5
		}
		if args.Limit > 10 {
			args.Limit = 10
		}

		onShelf := 1
		list, total, err := (&ProductService{}).GetList(1, args.Limit, args.Keyword, 0, 0, "sales_desc", 0, false, &onShelf)
		if err != nil {
			return nil, err
		}
		items := make([]map[string]interface{}, 0, len(list))
		for _, p := range list {
			items = append(items, map[string]interface{}{
				"id": p.ID, "name": p.Name, "price": p.Price, "original": p.OriginalPrice, "stock": p.Stock,
				"sales": p.Sales, "status": p.Status, "category": p.CategoryName, "image_url": p.ImageURL,
			})
		}
		return map[string]interface{}{"success": true, "keyword": args.Keyword, "total": total, "count": len(items), "items": items}, nil

	case toolCreateOrder:
		var args createOrderArgs
		if err := parseToolArguments(toolCall.Function.Arguments, &args); err != nil {
			return nil, err
		}
		if len(args.Items) == 0 {
			return nil, errors.New("items cannot be empty")
		}
		order, err := s.createOrderFromProducts(userID, args.StoreID, args.Items)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"success": true, "order_id": order.ID, "order_no": order.OrderNo, "store_id": order.StoreID,
			"total_amount": order.TotalAmount, "status": order.Status,
		}, nil

	case toolPayOrder:
		var args payOrderArgs
		if err := parseToolArguments(toolCall.Function.Arguments, &args); err != nil {
			return nil, err
		}
		orderID, err := s.resolveOrderIDForPayment(userID, args)
		if err != nil {
			return nil, err
		}
		authType := normalizeAIAuthType(args.PayType, execCtx.payType)
		password := strings.TrimSpace(args.Password)
		if password == "" {
			password = strings.TrimSpace(execCtx.paymentPassword)
		}
		faceImageURL := strings.TrimSpace(args.FaceImage)
		if faceImageURL == "" {
			faceImageURL = strings.TrimSpace(execCtx.faceImageURL)
		}
		payResult, err := s.payOrderWithAIAuth(userID, orderID, authType, password, faceImageURL)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{"success": true, "order_id": orderID, "payment_result": payResult}, nil
	}
	return nil, fmt.Errorf("unsupported tool call: %s", toolCall.Function.Name)
}

func (s *AIService) createOrderFromProducts(userID, storeID int64, rawItems []createOrderItemArg) (*model.Order, error) {
	merged := make(map[int64]int)
	for _, it := range rawItems {
		if it.ProductID <= 0 {
			return nil, errors.New("invalid product_id")
		}
		if it.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than 0")
		}
		merged[it.ProductID] += it.Quantity
	}

	productIDs := make([]int64, 0, len(merged))
	for pid := range merged {
		productIDs = append(productIDs, pid)
	}
	sort.Slice(productIDs, func(i, j int) bool { return productIDs[i] < productIDs[j] })

	var products []model.Product
	if err := global.DB.Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		return nil, err
	}
	pmap := make(map[int64]model.Product, len(products))
	for _, p := range products {
		pmap[p.ID] = p
	}
	for _, pid := range productIDs {
		p, ok := pmap[pid]
		if !ok {
			return nil, fmt.Errorf("product_id %d does not exist", pid)
		}
		if p.Status == 0 {
			return nil, fmt.Errorf("product_id %d is unavailable", pid)
		}
	}

	resolvedStoreID := storeID
	if resolvedStoreID <= 0 {
		var store model.Store
		if err := global.DB.Order("id ASC").First(&store).Error; err != nil {
			return nil, errors.New("no available store found")
		}
		resolvedStoreID = store.ID
	}

	cartIDs := make([]int64, 0, len(productIDs))
	orderItems := make([]model.CartItemParam, 0, len(productIDs))
	for _, pid := range productIDs {
		cart := model.Cart{UserID: userID, ProductID: pid, Quantity: merged[pid]}
		if err := global.DB.Create(&cart).Error; err != nil {
			s.cleanupTempCarts(userID, cartIDs)
			return nil, err
		}
		cartIDs = append(cartIDs, cart.ID)
		orderItems = append(orderItems, model.CartItemParam{CartID: cart.ID, Quantity: merged[pid]})
	}

	order, err := (&OrderService{}).CreateOrder(userID, resolvedStoreID, orderItems)
	if err != nil {
		s.cleanupTempCarts(userID, cartIDs)
		return nil, err
	}
	return order, nil
}

func (s *AIService) cleanupTempCarts(userID int64, cartIDs []int64) {
	if len(cartIDs) == 0 {
		return
	}
	_ = global.DB.Where("user_id = ? AND id IN ?", userID, cartIDs).Delete(&model.Cart{}).Error
}

func (s *AIService) resolveOrderIDForPayment(userID int64, args payOrderArgs) (int64, error) {
	orderNo := strings.TrimSpace(args.OrderNo)
	orderIDText := strings.TrimSpace(string(args.OrderIDRaw))
	orderIDText = strings.Trim(orderIDText, "\"")
	if orderIDText == "null" {
		orderIDText = ""
	}
	if orderIDText != "" {
		if parsed, err := strconv.ParseInt(orderIDText, 10, 64); err == nil && parsed > 0 {
			var order model.Order
			if err := global.DB.Where("id = ? AND user_id = ?", parsed, userID).First(&order).Error; err == nil {
				return order.ID, nil
			}
		}
		if orderNo == "" {
			orderNo = orderIDText
		}
	}

	if orderNo != "" {
		var order model.Order
		if err := global.DB.Where("user_id = ? AND order_no = ?", userID, orderNo).First(&order).Error; err == nil {
			return order.ID, nil
		}
	}

	// Fallback for LLM confusion: if a pending order exists, pay the newest one.
	var pending model.Order
	if err := global.DB.Where("user_id = ? AND status = ?", userID, 0).Order("id DESC").First(&pending).Error; err == nil {
		return pending.ID, nil
	}

	if orderNo == "" && orderIDText == "" {
		return 0, errors.New("order_id cannot be empty")
	}
	return 0, errors.New("order not found")
}

func normalizeAIAuthType(primary string, fallback string) string {
	value := strings.ToLower(strings.TrimSpace(primary))
	if value == "" {
		value = strings.ToLower(strings.TrimSpace(fallback))
	}
	if value == AuthTypeFace {
		return AuthTypeFace
	}
	return AuthTypePassword
}

func (s *AIService) payOrderWithAIAuth(userID int64, orderID int64, authType, password, faceImageURL string) (*MixedPaymentResult, error) {
	authType = normalizeAIAuthType(authType, AuthTypePassword)
	password = strings.TrimSpace(password)

	if authType == AuthTypeFace {
		if err := s.verifyAIFacePayment(userID, strings.TrimSpace(faceImageURL)); err != nil {
			return nil, err
		}
		password = ""
	} else if password == "" {
		return nil, errors.New("请输入支付密码")
	}

	result, err := (&FinanceService{}).UnifiedPayWithAuth(userID, orderID, PayTypeOrder, password, authType)
	if err != nil {
		return nil, localizeAIPaymentError(err)
	}
	return result, nil
}

func (s *AIService) verifyAIFacePayment(userID int64, faceImageURL string) error {
	if faceImageURL == "" {
		return errors.New("请先完成刷脸抓拍")
	}

	var user model.SysUser
	if err := global.DB.Select("id", "face_registered", "face_image_url").First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}
	if !user.FaceRegistered || strings.TrimSpace(user.FaceImageURL) == "" {
		return errors.New("当前账号未录入人脸，请先到个人中心录入")
	}

	faceService, err := NewFaceService()
	if err != nil {
		log.Printf("face service init failed in ai pay, userID=%d err=%v", userID, err)
		return errors.New("人脸服务暂不可用，请稍后重试")
	}
	score, err := faceService.CompareFace(user.FaceImageURL, faceImageURL)
	if err != nil {
		log.Printf("face verification failed in ai pay, userID=%d err=%v", userID, err)
		return errors.New("人脸验证失败，请稍后重试")
	}
	if score < facePayMinConfidence {
		return errors.New("人脸不匹配，请重试")
	}
	return nil
}

func localizeAIPaymentError(err error) error {
	if err == nil {
		return nil
	}
	msg := strings.TrimSpace(err.Error())
	lower := strings.ToLower(msg)
	switch {
	case strings.Contains(lower, "insufficient balance"):
		return errors.New("余额不足")
	case strings.Contains(lower, "payment password") && strings.Contains(lower, "incorrect"):
		return errors.New("支付密码错误")
	case strings.Contains(lower, "payment password"):
		return errors.New("请输入支付密码")
	case strings.Contains(lower, "order not found"):
		return errors.New("未找到订单")
	case strings.Contains(lower, "status") && strings.Contains(lower, "pay"):
		return errors.New("订单状态不支持支付")
	case strings.Contains(lower, "user not found"), strings.Contains(lower, "payer not found"):
		return errors.New("用户不存在")
	case strings.Contains(lower, "unsupported"):
		return errors.New("不支持的支付类型")
	default:
		return err
	}
}

func buildCommunityAgentTools() []dashScopeTool {
	return []dashScopeTool{
		{
			Type: "function",
			Function: dashScopeToolFunctionSpec{
				Name:        toolGetRecentNotices,
				Description: "Get recent community notices",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"limit": map[string]interface{}{"type": "integer", "minimum": 1, "maximum": 20},
					},
				},
			},
		},
		{
			Type: "function",
			Function: dashScopeToolFunctionSpec{
				Name:        toolCreateRepairTicket,
				Description: "Create repair or complaint ticket. type can be repair/complaint.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type":        map[string]interface{}{"type": "string", "description": "repair or complaint"},
						"category":    map[string]interface{}{"type": "string", "description": "ticket category, e.g. faucet, door, power"},
						"description": map[string]interface{}{"type": "string"},
					},
					"required": []string{"description"},
				},
			},
		},
		{
			Type: "function",
			Function: dashScopeToolFunctionSpec{
				Name:        toolSearchProduct,
				Description: "Search products by keyword",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"keyword": map[string]interface{}{"type": "string"},
						"limit":   map[string]interface{}{"type": "integer", "minimum": 1, "maximum": 10},
					},
					"required": []string{"keyword"},
				},
			},
		},
		{
			Type: "function",
			Function: dashScopeToolFunctionSpec{
				Name:        toolCreateOrder,
				Description: "Create order by product list",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"store_id": map[string]interface{}{"type": "integer"},
						"items": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"product_id": map[string]interface{}{"type": "integer"},
									"quantity":   map[string]interface{}{"type": "integer", "minimum": 1},
								},
								"required": []string{"product_id", "quantity"},
							},
						},
					},
					"required": []string{"items"},
				},
			},
		},
		{
			Type: "function",
			Function: dashScopeToolFunctionSpec{
				Name:        toolPayOrder,
				Description: "Pay order with password or face auth",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"order_id": map[string]interface{}{
							"type":        "string",
							"description": "Order primary key ID, if known",
						},
						"order_no": map[string]interface{}{
							"type":        "string",
							"description": "Order number, can be used instead of order_id",
						},
						"password": map[string]interface{}{
							"type":        "string",
							"description": "required when pay_type=password",
						},
						"pay_type": map[string]interface{}{
							"type":        "string",
							"description": "password or face",
						},
						"face_image_url": map[string]interface{}{
							"type":        "string",
							"description": "required when pay_type=face",
						},
					},
				},
			},
		},
	}
}

func parseToolArguments(raw string, target interface{}) error {
	clean := strings.TrimSpace(raw)
	if clean == "" || clean == "{}" || clean == "null" {
		return nil
	}
	if err := json.Unmarshal([]byte(clean), target); err == nil {
		return nil
	}
	var wrapped string
	if err := json.Unmarshal([]byte(clean), &wrapped); err == nil {
		wrapped = strings.TrimSpace(wrapped)
		if wrapped == "" {
			return nil
		}
		return json.Unmarshal([]byte(wrapped), target)
	}
	return fmt.Errorf("invalid JSON arguments: %s", clean)
}

func containsAny(text string, keywords ...string) bool {
	for _, kw := range keywords {
		if kw != "" && strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

func countMatches(text string, keywords ...string) int {
	count := 0
	for _, kw := range keywords {
		if kw != "" && strings.Contains(text, kw) {
			count++
		}
	}
	return count
}

func hasExplicitComplaintIntent(text string) bool {
	return containsAny(text,
		"\u6295\u8bc9", "\u6211\u8981\u6295\u8bc9", "\u5e2e\u6211\u6295\u8bc9", "\u8bf7\u5e2e\u6211\u6295\u8bc9",
		"complaint", "report complaint",
	)
}

func hasExplicitRepairIntent(text string) bool {
	return containsAny(text,
		"\u62a5\u4fee", "\u7ef4\u4fee", "\u6211\u8981\u62a5\u4fee", "\u5e2e\u6211\u62a5\u4fee", "\u8bf7\u5e2e\u6211\u62a5\u4fee",
		"repair", "fix", "maintenance",
	)
}

func classifyRepairIntent(typeHint, description string) (int, string, bool) {
	typeText := strings.ToLower(strings.TrimSpace(typeHint))
	descText := strings.ToLower(strings.TrimSpace(description))
	text := strings.TrimSpace(typeText + " " + descText)

	repairAction := hasExplicitRepairIntent(typeText) || hasExplicitRepairIntent(descText)
	complaintAction := hasExplicitComplaintIntent(typeText) || hasExplicitComplaintIntent(descText)
	if complaintAction && !repairAction {
		return 2, "complaint", false
	}
	if repairAction && !complaintAction {
		return 1, "repair", false
	}

	repairSymptoms := []string{
		"\u574f\u4e86", "\u6545\u969c", "\u6f0f\u6c34", "\u4e0d\u70ed", "\u4e0d\u4eae", "\u8df3\u95f8", "\u65e0\u6c34", "\u65e0\u7535",
		"\u6696\u6c14", "\u4f9b\u6696", "\u6c34\u9f99\u5934", "\u95e8", "\u7a97", "\u7535\u8def", "\u63d2\u5ea7", "\u9a6c\u6876", "\u7a7a\u8c03",
		"repair", "fix", "broken", "leak", "heating", "radiator", "power",
	}
	complaintReasons := []string{
		"\u6001\u5ea6\u5dee", "\u670d\u52a1\u5dee", "\u4e71\u6536\u8d39", "\u6536\u8d39\u4e0d\u5408\u7406", "\u6270\u6c11", "\u566a\u97f3",
		"\u4e0d\u5904\u7406", "\u4e0d\u4f5c\u4e3a", "\u957f\u671f\u4e0d\u89e3\u51b3", "\u5f88\u4e0d\u6ee1\u610f", "\u6295\u8bc9",
		"complaint", "service bad", "overcharge", "noise",
	}

	repairScore := countMatches(text, repairSymptoms...)
	complaintScore := countMatches(text, complaintReasons...)

	if complaintScore > repairScore {
		return 2, "complaint", false
	}
	if repairScore > complaintScore {
		return 1, "repair", false
	}

	return 1, "repair", true
}

func normalizeRepairCategory(categoryHint, typeHint, description string) string {
	explicit := normalizeRepairCategoryLabel(categoryHint)
	if explicit != "" {
		return explicit
	}

	text := strings.ToLower(strings.TrimSpace(typeHint + " " + description))
	switch {
	case containsAny(text, "\u6c34\u9f99\u5934", "\u6f0f\u6c34", "\u4e0b\u6c34", "\u6c34\u7ba1", "\u7ba1\u9053", "\u9a6c\u6876", "faucet", "plumb", "pipe"):
		return "\u6c34\u6696"
	case containsAny(text, "\u95e8", "\u7a97", "\u95e8\u7a97", "door", "window", "lock"):
		return "\u95e8\u7a97"
	case containsAny(text, "\u7535", "\u7535\u8def", "\u8df3\u95f8", "\u63d2\u5ea7", "\u706f", "power", "electric"):
		return "\u7535\u8def"
	case containsAny(text, "\u7a7a\u8c03", "aircon", "ac"):
		return "\u7a7a\u8c03"
	case containsAny(text, "\u6696\u6c14", "\u4f9b\u6696", "\u4e0d\u70ed", "heating", "radiator"):
		return "\u4f9b\u6696"
	case containsAny(text, "\u566a\u97f3", "noise"):
		return "\u566a\u97f3"
	case containsAny(text, "\u536b\u751f", "\u5783\u573e", "clean"):
		return "\u536b\u751f"
	default:
		hint := normalizeRepairCategoryLabel(typeHint)
		if hint != "" {
			return hint
		}
		return "\u5176\u4ed6"
	}
}

func normalizeRepairCategoryLabel(raw string) string {
	v := strings.TrimSpace(raw)
	if v == "" {
		return ""
	}
	switch strings.ToLower(v) {
	case "plumbing", "water", "pipe", "\u6c34\u6696", "\u6f0f\u6c34", "\u4e0b\u6c34", "\u6c34\u7ba1", "\u7ba1\u9053":
		return "\u6c34\u6696"
	case "door_window", "door", "window", "lock", "\u95e8\u7a97":
		return "\u95e8\u7a97"
	case "electrical", "electric", "power", "\u7535\u8def", "\u7528\u7535":
		return "\u7535\u8def"
	case "air_conditioner", "aircon", "ac", "\u7a7a\u8c03":
		return "\u7a7a\u8c03"
	case "heating", "radiator", "\u6696\u6c14", "\u4f9b\u6696":
		return "\u4f9b\u6696"
	case "noise", "\u566a\u97f3", "\u6270\u6c11":
		return "\u566a\u97f3"
	case "sanitation", "clean", "\u536b\u751f", "\u5783\u573e":
		return "\u536b\u751f"
	case "other", "\u5176\u4ed6":
		return "\u5176\u4ed6"
	default:
		return v
	}
}

func extractRepairDescription(text string) string {
	content := strings.TrimSpace(text)
	if content == "" {
		return ""
	}

	prefixes := []string{
		"\u5e2e\u6211\u62a5\u4fee\u4e00\u4e0b", "\u5e2e\u6211\u62a5\u4fee", "\u8bf7\u5e2e\u6211\u62a5\u4fee", "\u62a5\u4fee\u4e00\u4e0b", "\u6211\u8981\u62a5\u4fee",
		"\u5e2e\u6211\u6295\u8bc9\u4e00\u4e0b", "\u5e2e\u6211\u6295\u8bc9", "\u8bf7\u5e2e\u6211\u6295\u8bc9", "\u6295\u8bc9\u4e00\u4e0b", "\u6211\u8981\u6295\u8bc9",
	}
	for _, p := range prefixes {
		if strings.HasPrefix(content, p) {
			content = strings.TrimSpace(strings.TrimPrefix(content, p))
			break
		}
	}
	if content == "" {
		content = strings.TrimSpace(text)
	}
	return content
}

func (s *AIService) createAndVerifyRepairTicket(userID int64, repairType int, category, content string) (*model.Repair, error) {
	repair := &model.Repair{
		UserID:   userID,
		Type:     repairType,
		Category: normalizeRepairCategoryLabel(category),
		Content:  strings.TrimSpace(content),
	}
	if err := (&RepairService{}).Create(repair); err != nil {
		return nil, err
	}
	if repair.ID <= 0 {
		return nil, errors.New("repair ticket create failed: empty id")
	}

	var saved model.Repair
	if err := global.DB.Where("id = ? AND user_id = ?", repair.ID, userID).First(&saved).Error; err != nil {
		return nil, fmt.Errorf("repair ticket verification failed: %w", err)
	}
	return &saved, nil
}

func (s *AIService) tryHandleRepairIntent(userID int64, lastUser string) (string, bool, error) {
	text := strings.TrimSpace(lastUser)
	if text == "" || isPayIntent(text) || isCreateOrderIntent(text) {
		return "", false, nil
	}
	if !containsAny(strings.ToLower(text),
		"\u62a5\u4fee", "\u7ef4\u4fee", "repair", "fix", "\u6295\u8bc9", "complaint",
		"\u574f\u4e86", "\u6545\u969c", "\u6f0f\u6c34", "\u566a\u97f3", "\u6270\u6c11", "\u6696\u6c14", "\u4f9b\u6696", "\u4e0d\u70ed",
	) {
		return "", false, nil
	}

	description := extractRepairDescription(text)
	if strings.TrimSpace(description) == "" {
		return "\u8bf7\u8865\u5145\u5177\u4f53\u95ee\u9898\u63cf\u8ff0\uff08\u4f8b\u5982\u4f4d\u7f6e\u3001\u6545\u969c\u73b0\u8c61\uff09\uff0c\u6211\u518d\u4e3a\u60a8\u521b\u5efa\u5de5\u5355\u3002", true, nil
	}

	repairType, typeLabel, ambiguous := classifyRepairIntent("", text)
	if ambiguous {
		return "\u8bf7\u786e\u8ba4\u60a8\u662f\u8981\u63d0\u4ea4\u300c\u62a5\u4fee\u300d\u8fd8\u662f\u300c\u6295\u8bc9\u300d\uff1f\u56de\u590d\u201c\u62a5\u4fee\u201d\u6216\u201c\u6295\u8bc9\u201d\u5373\u53ef\u3002", true, nil
	}
	category := normalizeRepairCategory("", "", description)
	repair, err := s.createAndVerifyRepairTicket(userID, repairType, category, description)
	if err != nil {
		return fmt.Sprintf("\u521b\u5efa\u5de5\u5355\u5931\u8d25\uff1a%v", err), true, nil
	}

	if typeLabel == "complaint" {
		return fmt.Sprintf("\u2705 \u5df2\u4e3a\u60a8\u521b\u5efa\u6295\u8bc9\u5de5\u5355\uff01\n- \u5de5\u5355\u53f7\uff1a#%d\n- \u5de5\u5355\u7c7b\u578b\uff1a\u6295\u8bc9\n- \u5f53\u524d\u72b6\u6001\uff1apending", repair.ID), true, nil
	}
	return fmt.Sprintf("\u2705 \u5df2\u4e3a\u60a8\u521b\u5efa\u62a5\u4fee\u5de5\u5355\uff01\n- \u5de5\u5355\u53f7\uff1a#%d\n- \u5de5\u5355\u7c7b\u578b\uff1a\u62a5\u4fee\n- \u5f53\u524d\u72b6\u6001\uff1apending", repair.ID), true, nil
}
func requiresForcedSearchIntent(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" || isPayIntent(text) || looksLikePasswordOnly(text) {
		return false
	}
	searchKeywords := []string{
		"买", "购买", "搜索", "查", "查询", "找", "推荐", "商品", "商城", "洗衣液", "日用品", "下单",
		"product", "search", "shop",
	}
	for _, kw := range searchKeywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

func extractProductKeyword(text string) string {
	clean := strings.TrimSpace(text)
	if clean == "" {
		return ""
	}
	replacer := strings.NewReplacer(
		"帮我", "", "请", "", "我要", "", "我想", "", "给我", "",
		"买", "", "购买", "", "搜索", "", "查", "", "查询", "", "找", "", "推荐", "",
		"商品", "", "商城", "", "下单", "", "自动", "", "支付", "", "付款", "", "结算", "", "确认", "", "并", "", "然后", "",
		"马上", "", "立即", "", "那款", "", "这款", "", "这个", "", "那个", "",
		"*", "", "#", "", "，", "", "。", "", "！", "", "？", "", ",", "", ".", "", "!", "", "?", "",
		" ", "",
	)
	clean = strings.TrimSpace(strings.ToLower(replacer.Replace(clean)))
	if len([]rune(clean)) < 2 || regexp.MustCompile(`^[0-9一二两三四五六七八九十百千瓶件个份盒袋支包]+$`).MatchString(clean) {
		return ""
	}
	if len([]rune(clean)) > 24 {
		runes := []rune(clean)
		clean = string(runes[:24])
	}
	return clean
}

func buildForcedSearchToolCall(lastUserText string) *dashScopeToolCall {
	if !requiresForcedSearchIntent(lastUserText) {
		return nil
	}
	keyword := extractProductKeyword(lastUserText)
	if keyword == "" {
		return nil
	}
	argsJSON, _ := json.Marshal(searchProductArgs{Keyword: keyword, Limit: 5})
	return &dashScopeToolCall{
		ID:   fmt.Sprintf("forced_search_%d", time.Now().UnixNano()),
		Type: "function",
		Function: dashScopeToolCallFunction{
			Name:      toolSearchProduct,
			Arguments: string(argsJSON),
		},
	}
}

func latestUserMessage(messages []dashScopeMessage) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role != "user" {
			continue
		}
		if s, ok := messages[i].Content.(string); ok {
			return strings.TrimSpace(s)
		}
		raw, _ := json.Marshal(messages[i].Content)
		return strings.TrimSpace(string(raw))
	}
	return ""
}

func isCriticalTool(name string) bool {
	return name == toolCreateOrder || name == toolPayOrder || name == toolCreateRepairTicket
}

func formatCriticalToolError(toolName string, err error) error {
	switch toolName {
	case toolCreateRepairTicket:
		return fmt.Errorf("工单创建失败：%s", err.Error())
	case toolCreateOrder:
		return fmt.Errorf("订单创建失败：%s", err.Error())
	case toolPayOrder:
		return fmt.Errorf("订单支付失败：%s", err.Error())
	default:
		return err
	}
}

func (s *AIService) tryHandleCreateOrderIntent(userID int64, messages []dashScopeMessage, lastUser string) (string, bool, error) {
	if !isCreateOrderIntent(lastUser) {
		return "", false, nil
	}

	keyword := extractProductKeyword(lastUser)
	if keyword == "" {
		keyword = inferProductKeywordFromHistory(messages)
	}
	if keyword == "" {
		return "请先告诉我要购买的具体商品名称，例如“下单1瓶洗衣液”。", true, nil
	}

	quantity := extractOrderQuantity(lastUser)
	onShelf := 1
	products, _, err := (&ProductService{}).GetList(1, 1, keyword, 0, 0, "sales_desc", 0, false, &onShelf)
	if err != nil {
		return fmt.Sprintf("下单失败：搜索商品失败：%v", err), true, nil
	}
	if len(products) == 0 {
		return fmt.Sprintf("下单失败：未找到商品“%s”。", keyword), true, nil
	}

	order, err := s.createOrderFromProducts(userID, 0, []createOrderItemArg{{
		ProductID: products[0].ID,
		Quantity:  quantity,
	}})
	if err != nil {
		return fmt.Sprintf("下单失败：%v", err), true, nil
	}

	return fmt.Sprintf(
		"✅ 订单已创建成功\n订单号：%s\n应付金额：￥%.2f\n商品：%s x%d\n\n请发送“确认支付”并输入登录密码完成支付。",
		order.OrderNo,
		order.TotalAmount,
		products[0].Name,
		quantity,
	), true, nil
}

func (s *AIService) tryHandlePayIntent(userID int64, execCtx chatExecutionContext, lastUser string) (string, bool, error) {
	if !isPayIntent(lastUser) && !looksLikePasswordOnly(lastUser) {
		return "", false, nil
	}

	authType := normalizeAIAuthType(execCtx.payType, AuthTypePassword)
	password := strings.TrimSpace(execCtx.paymentPassword)
	if authType == AuthTypePassword && password == "" && looksLikePasswordOnly(lastUser) {
		password = strings.TrimSpace(lastUser)
	}
	if authType == AuthTypePassword && password == "" {
		return "支付失败：请先点击“确认支付”并输入登录密码。", true, nil
	}
	if authType == AuthTypeFace && strings.TrimSpace(execCtx.faceImageURL) == "" {
		return "支付失败：请先完成刷脸抓拍。", true, nil
	}

	var order model.Order
	if err := global.DB.Where("user_id = ? AND status = ?", userID, 0).Order("id DESC").First(&order).Error; err != nil {
		return "支付失败：未找到待支付订单。", true, nil
	}

	payResult, err := s.payOrderWithAIAuth(userID, order.ID, authType, password, execCtx.faceImageURL)
	if err != nil {
		return fmt.Sprintf("支付失败：%v", err), true, nil
	}

	return fmt.Sprintf(
		"✅ 支付成功\n订单号：%s\n支付金额：￥%.2f\n积分抵扣：%d\n余额支付：￥%.2f",
		order.OrderNo,
		payResult.TotalAmount,
		payResult.UsedPoints,
		payResult.UsedBalance,
	), true, nil
}

func isCreateOrderIntent(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return false
	}
	if looksLikePasswordOnly(text) {
		return false
	}

	keywords := []string{"下单", "创建订单", "生成订单", "提交订单", "立即购买", "去结算", "拍下"}
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}

	re := regexp.MustCompile(`^(买|购买).*(\d+|一|二|两|三|四|五|六|七|八|九|十).*(瓶|件|个|份|盒|袋|支|包)`)
	return re.MatchString(text)
}

func isPayIntent(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return false
	}
	keywords := []string{"确认支付", "立即支付", "去支付", "支付", "付款", "结算", "付钱"}
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

func looksLikePasswordOnly(text string) bool {
	text = strings.TrimSpace(text)
	if text == "" {
		return false
	}
	if len(text) < 4 || len(text) > 32 {
		return false
	}
	for _, ch := range text {
		if (ch < '0' || ch > '9') && (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') {
			return false
		}
	}
	return true
}

func extractOrderQuantity(text string) int {
	text = strings.TrimSpace(text)
	if text == "" {
		return 1
	}
	re := regexp.MustCompile(`(\d+)`)
	if m := re.FindStringSubmatch(text); len(m) == 2 {
		if parsed, err := strconv.Atoi(m[1]); err == nil && parsed > 0 {
			return parsed
		}
	}

	chineseQuantityMap := []struct {
		token string
		value int
	}{
		{"十", 10},
		{"九", 9},
		{"八", 8},
		{"七", 7},
		{"六", 6},
		{"五", 5},
		{"四", 4},
		{"三", 3},
		{"两", 2},
		{"俩", 2},
		{"二", 2},
		{"一", 1},
	}
	for _, item := range chineseQuantityMap {
		if strings.Contains(text, item.token) {
			return item.value
		}
	}

	switch {
	case strings.Contains(strings.ToLower(text), "two"):
		return 2
	case strings.Contains(strings.ToLower(text), "three"):
		return 3
	case strings.Contains(strings.ToLower(text), "four"):
		return 4
	case strings.Contains(strings.ToLower(text), "five"):
		return 5
	default:
		return 1
	}
}

func inferProductKeywordFromHistory(messages []dashScopeMessage) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role != "user" {
			continue
		}
		text, ok := messages[i].Content.(string)
		if !ok {
			continue
		}
		keyword := extractProductKeyword(text)
		if keyword != "" {
			return keyword
		}
	}
	return ""
}

func (s *AIService) GetChatHistory(userID int64, limit int) ([]model.ChatMessage, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > maxChatHistoryLimit {
		limit = maxChatHistoryLimit
	}

	var list []model.ChatMessage
	if err := global.DB.Where("user_id = ?", userID).Order("id DESC").Limit(limit).Find(&list).Error; err != nil {
		return nil, errors.New("failed to fetch chat history")
	}
	reverseChatMessages(list)
	return list, nil
}

func (s *AIService) GenerateCommunityReport(prompt string) (string, error) {
	messages := []dashScopeMessage{
		{
			Role: "system",
			Content: "You are a senior community operations manager. " +
				"Generate concise, professional, management-oriented analysis in Chinese.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}
	return s.callTextModel(config.Conf.AI.Model, messages, nil)
}

func (s *AIService) RecognizeGarbage(imageURL string) (*GarbageRecognitionResult, error) {
	basePrompt := `请识别图片中的垃圾是否已正确分类，并返回 1 到 50 的整数奖励积分。
只返回 JSON，例如：{"points": 20, "reason": "分类准确"}。`

	buildMessages := func(prompt string) []dashScopeMessage {
		return []dashScopeMessage{
			{
				Role: "user",
				Content: []map[string]interface{}{
					{
						"type": "text",
						"text": prompt,
					},
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": imageURL,
						},
					},
				},
			},
		}
	}

	prompts := []string{
		basePrompt,
		basePrompt + "\n只允许返回一个合法 JSON 对象，禁止返回任何额外字符。",
	}
	var lastErr error
	for idx, prompt := range prompts {
		content, err := s.callTextModel(qwenVisionModel, buildMessages(prompt), map[string]string{"type": "json_object"})
		if err != nil {
			return nil, err
		}
		result, parseErr := parseGarbageRecognitionResult(content)
		if parseErr == nil {
			return result, nil
		}
		lastErr = parseErr
		log.Printf("garbage recognition parse failed, attempt=%d content=%q err=%v", idx+1, content, parseErr)
	}

	if lastErr == nil {
		lastErr = errors.New("empty response")
	}
	return nil, fmt.Errorf("garbage recognition response format invalid: %w", lastErr)
}

func (s *AIService) callTextModel(model string, messages []dashScopeMessage, responseFormat map[string]string) (string, error) {
	message, err := s.callModel(model, messages, responseFormat, nil)
	if err != nil {
		return "", err
	}
	content := strings.TrimSpace(extractMessageContent(message.Content))
	if content == "" {
		return "", errors.New("no response from AI")
	}
	return content, nil
}

func (s *AIService) callModel(model string, messages []dashScopeMessage, responseFormat map[string]string, tools []dashScopeTool) (dashScopeMessage, error) {
	apiKey := strings.TrimSpace(config.Conf.AI.APIKey)
	baseURL := strings.TrimSpace(config.Conf.AI.BaseURL)
	if apiKey == "" || baseURL == "" {
		return dashScopeMessage{}, errors.New("AI service is not configured")
	}

	reqBody := dashScopeRequest{
		Model:          model,
		Messages:       messages,
		ResponseFormat: responseFormat,
		Tools:          tools,
	}
	if len(tools) > 0 {
		reqBody.ToolChoice = "auto"
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return dashScopeMessage{}, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return dashScopeMessage{}, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return dashScopeMessage{}, err
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return dashScopeMessage{}, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		log.Printf("AI request failed with status=%d body=%s", resp.StatusCode, string(bodyText))
		return dashScopeMessage{}, fmt.Errorf("AI request failed with status %d", resp.StatusCode)
	}

	var aiResp dashScopeResponse
	if err := json.Unmarshal(bodyText, &aiResp); err != nil {
		return dashScopeMessage{}, err
	}
	if aiResp.Error != nil && aiResp.Error.Message != "" {
		return dashScopeMessage{}, errors.New(aiResp.Error.Message)
	}
	if aiResp.Code != "" {
		return dashScopeMessage{}, errors.New(aiResp.Message)
	}
	if len(aiResp.Choices) == 0 {
		return dashScopeMessage{}, errors.New("no response from AI")
	}
	return aiResp.Choices[0].Message, nil
}

func extractMessageContent(content interface{}) string {
	switch value := content.(type) {
	case string:
		return value
	case []interface{}:
		var parts []string
		for _, item := range value {
			partMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if text, ok := partMap["text"].(string); ok {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, "\n")
	default:
		raw, _ := json.Marshal(value)
		return string(raw)
	}
}

func parseGarbageRecognitionResult(raw string) (*GarbageRecognitionResult, error) {
	candidates := buildGarbageJSONCandidates(raw)
	var lastErr error
	for _, clean := range candidates {
		result, err := decodeGarbageRecognitionResult(clean)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = errors.New("empty response")
	}
	return nil, lastErr
}

func buildGarbageJSONCandidates(raw string) []string {
	clean := strings.TrimSpace(raw)
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimSuffix(clean, "```")
	clean = strings.TrimSpace(clean)
	if clean == "" {
		return nil
	}

	candidates := make([]string, 0, 4)
	seen := make(map[string]struct{}, 4)
	appendCandidate := func(v string) {
		v = strings.TrimSpace(v)
		if v == "" {
			return
		}
		if _, ok := seen[v]; ok {
			return
		}
		seen[v] = struct{}{}
		candidates = append(candidates, v)
	}

	appendCandidate(clean)

	start := strings.Index(clean, "{")
	end := strings.LastIndex(clean, "}")
	if start >= 0 && end > start {
		appendCandidate(clean[start : end+1])
	}
	if start >= 0 && end < start {
		appendCandidate(clean[start:] + "}")
	}

	re := regexp.MustCompile(`\{[\s\S]*\}`)
	if match := re.FindString(clean); match != "" {
		appendCandidate(match)
	}

	return candidates
}

func decodeGarbageRecognitionResult(clean string) (*GarbageRecognitionResult, error) {
	normalized := strings.NewReplacer("“", "\"", "”", "\"", "，", ",", "：", ":").Replace(clean)
	normalized = regexp.MustCompile(`,\s*}`).ReplaceAllString(normalized, "}")
	normalized = strings.TrimSpace(normalized)

	var result GarbageRecognitionResult
	if err := json.Unmarshal([]byte(normalized), &result); err != nil {
		return nil, err
	}
	if result.Points < 1 {
		result.Points = 1
	}
	if result.Points > 50 {
		result.Points = 50
	}
	result.Reason = strings.TrimSpace(result.Reason)
	if result.Reason == "" {
		result.Reason = "AI completed garbage classification analysis"
	}
	return &result, nil
}

func reverseChatMessages(messages []model.ChatMessage) {
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
}
