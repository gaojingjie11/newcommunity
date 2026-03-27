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

	toolGetRecentNotices   = "get_recent_notices"
	toolCreateRepairTicket = "create_repair_ticket"
	toolSearchProduct      = "search_product"
	toolCreateOrder        = "create_order"
	toolPayOrder           = "pay_order"
)

const communityAssistantSystemPrompt = "You are a smart community assistant. Use tools for real data operations. " +
	"For order payment, user password is required and payment must be explicitly confirmed by user. " +
	"Never claim order creation/payment success unless backend tool execution succeeds."

type AIService struct{}

type chatExecutionContext struct {
	paymentPassword string
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
}

func (s *AIService) Chat(userContent string) (string, error) {
	messages := []dashScopeMessage{
		{Role: "system", Content: communityAssistantSystemPrompt},
		{Role: "user", Content: userContent},
	}
	return s.callTextModel(config.Conf.AI.Model, messages, nil)
}

func (s *AIService) ChatWithMemory(userID int64, userContent string, paymentPassword string) (string, error) {
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

	reply, err := s.runAgentWithTools(userID, chatExecutionContext{paymentPassword: strings.TrimSpace(paymentPassword)}, messages)
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
		if isCriticalOrderTool(toolCall.Function.Name) {
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
		if strings.TrimSpace(args.Description) == "" {
			return nil, errors.New("description cannot be empty")
		}
		tp := strings.TrimSpace(args.Type)
		if tp == "" {
			tp = "other"
		}
		repair := &model.Repair{UserID: userID, Type: mapRepairType(tp), Category: tp, Content: strings.TrimSpace(args.Description)}
		if err := (&RepairService{}).Create(repair); err != nil {
			return nil, err
		}
		return map[string]interface{}{"success": true, "repair_id": repair.ID, "status": "pending"}, nil

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
		password := strings.TrimSpace(args.Password)
		if password == "" {
			password = strings.TrimSpace(execCtx.paymentPassword)
		}
		if password == "" {
			return nil, errors.New("payment password is required")
		}
		payResult, err := (&FinanceService{}).UnifiedPay(userID, orderID, PayTypeOrder, password)
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
				Description: "Create repair ticket",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type":        map[string]interface{}{"type": "string"},
						"description": map[string]interface{}{"type": "string"},
					},
					"required": []string{"type", "description"},
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
				Description: "Pay order with login password",
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
						"password": map[string]interface{}{"type": "string"},
					},
					"required": []string{"password"},
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

func mapRepairType(typeLabel string) int {
	normalized := strings.ToLower(strings.TrimSpace(typeLabel))
	switch {
	case strings.Contains(normalized, "水"), strings.Contains(normalized, "plumb"):
		return 1
	case strings.Contains(normalized, "电"), strings.Contains(normalized, "power"):
		return 2
	case strings.Contains(normalized, "门禁"), strings.Contains(normalized, "access"):
		return 3
	case strings.Contains(normalized, "电梯"), strings.Contains(normalized, "lift"), strings.Contains(normalized, "elevator"):
		return 4
	case strings.Contains(normalized, "网络"), strings.Contains(normalized, "网"), strings.Contains(normalized, "network"):
		return 5
	default:
		return 0
	}
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

func isCriticalOrderTool(name string) bool {
	return name == toolCreateOrder || name == toolPayOrder
}

func formatCriticalToolError(toolName string, err error) error {
	switch toolName {
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

	password := strings.TrimSpace(execCtx.paymentPassword)
	if password == "" && looksLikePasswordOnly(lastUser) {
		password = strings.TrimSpace(lastUser)
	}
	if password == "" {
		return "支付失败：请先点击“确认支付”并输入登录密码。", true, nil
	}

	var order model.Order
	if err := global.DB.Where("user_id = ? AND status = ?", userID, 0).Order("id DESC").First(&order).Error; err != nil {
		return "支付失败：未找到待支付订单。", true, nil
	}

	payResult, err := (&FinanceService{}).UnifiedPay(userID, order.ID, PayTypeOrder, password)
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
	prompt := `请识别图片中的垃圾是否已正确分类，并根据垃圾数量和分类准确度返回 1 到 50 的整数奖励积分。
请严格只返回 JSON：{"points": 20, "reason": "分类准确，包含可回收物"}`
	messages := []dashScopeMessage{
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

	content, err := s.callTextModel(qwenVisionModel, messages, map[string]string{"type": "json_object"})
	if err != nil {
		return nil, err
	}

	result, err := parseGarbageRecognitionResult(content)
	if err != nil {
		return nil, err
	}
	return result, nil
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
	clean := strings.TrimSpace(raw)
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimSuffix(clean, "```")
	clean = strings.TrimSpace(clean)

	if !strings.HasPrefix(clean, "{") {
		re := regexp.MustCompile(`\{[\s\S]*\}`)
		match := re.FindString(clean)
		if match != "" {
			clean = match
		}
	}

	var result GarbageRecognitionResult
	if err := json.Unmarshal([]byte(clean), &result); err != nil {
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
