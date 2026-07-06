package logic

import "strings"

type AgentIntent string

const (
	IntentFallbackAgent      AgentIntent = "fallback_agent"
	IntentLatestAIReport     AgentIntent = "latest_ai_report"
	IntentLatestNotices      AgentIntent = "latest_notices"
	IntentNoticeSemantic     AgentIntent = "notice_semantic"
	IntentSubmitRepair       AgentIntent = "submit_repair"
	IntentShoppingCapability AgentIntent = "shopping_capability"
	IntentProductBrowse      AgentIntent = "product_browse"
	IntentProductBuyConfirm  AgentIntent = "product_buy_confirm"
)

type IntentRouteResult struct {
	Intent      AgentIntent
	Confidence  float64
	UseFastPath bool
	Keyword     string
	Reason      string
}

func RouteAgentIntent(message string) IntentRouteResult {
	text := strings.TrimSpace(message)
	if text == "" {
		return IntentRouteResult{Intent: IntentFallbackAgent}
	}

	if shouldDirectLatestAIReport(text) {
		return IntentRouteResult{
			Intent:      IntentLatestAIReport,
			Confidence:  0.99,
			UseFastPath: true,
			Reason:      "strict_ai_report",
		}
	}

	if shouldDirectNoticeKnowledge(text) {
		return IntentRouteResult{
			Intent:      IntentNoticeSemantic,
			Confidence:  0.95,
			UseFastPath: true,
			Reason:      "semantic_notice_search",
		}
	}

	if shouldDirectLatestNotices(text) {
		return IntentRouteResult{
			Intent:      IntentLatestNotices,
			Confidence:  0.97,
			UseFastPath: true,
			Keyword:     extractNoticeKeyword(text),
			Reason:      "latest_notice_list",
		}
	}

	if shouldDirectShoppingCapability(text) {
		return IntentRouteResult{
			Intent:      IntentShoppingCapability,
			Confidence:  0.98,
			UseFastPath: true,
			Reason:      "shopping_capability_question",
		}
	}

	if route := routeServiceRequestIntent(text); route.UseFastPath {
		return route
	}

	if route := routeShoppingBuyIntent(text); route.UseFastPath {
		return route
	}

	if route := routeShoppingBrowseIntent(text); route.UseFastPath {
		return route
	}

	return IntentRouteResult{
		Intent:      IntentFallbackAgent,
		Confidence:  0.0,
		UseFastPath: false,
		Reason:      "fallback_to_agent",
	}
}

func routeServiceRequestIntent(message string) IntentRouteResult {
	text := strings.ToLower(strings.TrimSpace(message))
	if text == "" {
		return IntentRouteResult{Intent: IntentFallbackAgent}
	}
	if strings.Contains(text, "流程") || strings.Contains(text, "怎么投诉") || strings.Contains(text, "如何投诉") ||
		strings.Contains(text, "怎么报修") || strings.Contains(text, "如何报修") || strings.Contains(text, "能投诉吗") {
		return IntentRouteResult{Intent: IntentFallbackAgent}
	}

	complaintPhrases := []string{"帮我投诉", "投诉下", "投诉一下", "我要投诉", "我想投诉", "提交投诉", "帮我反馈", "我要反馈", "我想反馈"}
	for _, phrase := range complaintPhrases {
		if strings.Contains(text, phrase) {
			return IntentRouteResult{
				Intent:      IntentSubmitRepair,
				Confidence:  0.97,
				UseFastPath: true,
				Reason:      "explicit_complaint_submit",
			}
		}
	}

	repairPhrases := []string{"帮我报修", "报修下", "报修一下", "我要报修", "我想报修", "提交报修", "帮我维修", "维修一下"}
	for _, phrase := range repairPhrases {
		if strings.Contains(text, phrase) {
			return IntentRouteResult{
				Intent:      IntentSubmitRepair,
				Confidence:  0.97,
				UseFastPath: true,
				Reason:      "explicit_repair_submit",
			}
		}
	}

	if strings.Contains(text, "帮我") || strings.Contains(text, "麻烦") || strings.Contains(text, "请") {
		issueKeywords := []string{"漏水", "跳闸", "坏了", "故障", "不亮", "堵", "垃圾桶", "爆满", "没人清理", "太吵", "噪音", "电梯"}
		for _, keyword := range issueKeywords {
			if strings.Contains(text, keyword) {
				return IntentRouteResult{
					Intent:      IntentSubmitRepair,
					Confidence:  0.88,
					UseFastPath: true,
					Reason:      "issue_submit_with_clear_problem",
				}
			}
		}
	}

	return IntentRouteResult{Intent: IntentFallbackAgent}
}

func routeShoppingBuyIntent(message string) IntentRouteResult {
	text := strings.ToLower(strings.TrimSpace(message))
	if text == "" {
		return IntentRouteResult{Intent: IntentFallbackAgent}
	}
	if shouldDirectShoppingCapability(text) {
		return IntentRouteResult{Intent: IntentFallbackAgent}
	}
	if strings.Contains(text, "支付") || strings.Contains(text, "订单状态") || strings.Contains(text, "取消订单") {
		return IntentRouteResult{Intent: IntentFallbackAgent}
	}

	keyword := extractOrderKeyword(message)
	if isGenericKeyword(keyword) {
		return IntentRouteResult{Intent: IntentFallbackAgent}
	}

	if strings.Contains(text, "直接下单") || strings.Contains(text, "帮我下单") {
		return IntentRouteResult{
			Intent:      IntentProductBuyConfirm,
			Confidence:  0.99,
			UseFastPath: true,
			Keyword:     keyword,
			Reason:      "explicit_create_order",
		}
	}

	buyPhrases := []string{"帮我买", "我要买", "给我买", "买个", "买一", "来个"}
	for _, phrase := range buyPhrases {
		if strings.Contains(text, phrase) {
			return IntentRouteResult{
				Intent:      IntentProductBuyConfirm,
				Confidence:  0.92,
				UseFastPath: true,
				Keyword:     keyword,
				Reason:      "specific_product_buy",
			}
		}
	}

	return IntentRouteResult{Intent: IntentFallbackAgent}
}

func routeShoppingBrowseIntent(message string) IntentRouteResult {
	text := strings.ToLower(strings.TrimSpace(message))
	if text == "" {
		return IntentRouteResult{Intent: IntentFallbackAgent}
	}
	if strings.Contains(text, "订单") || strings.Contains(text, "支付") || strings.Contains(text, "报修") || strings.Contains(text, "投诉") {
		return IntentRouteResult{Intent: IntentFallbackAgent}
	}

	if shouldDirectProductBrowse(text) {
		return IntentRouteResult{
			Intent:      IntentProductBrowse,
			Confidence:  0.96,
			UseFastPath: true,
			Keyword:     extractProductKeyword(message),
			Reason:      "explicit_product_browse",
		}
	}

	if route := routeSpecificProductDiscovery(text, message); route.UseFastPath {
		return route
	}

	return IntentRouteResult{Intent: IntentFallbackAgent}
}

func routeSpecificProductDiscovery(lowerText, rawMessage string) IntentRouteResult {
	intentPhrases := []string{"我想买", "想买", "我要买", "买点", "来点", "想要", "想看看", "有没有"}
	matchedPhrase := false
	for _, phrase := range intentPhrases {
		if strings.Contains(lowerText, phrase) {
			matchedPhrase = true
			break
		}
	}
	if !matchedPhrase {
		return IntentRouteResult{Intent: IntentFallbackAgent}
	}

	keyword := extractProductKeyword(rawMessage)
	if isGenericKeyword(keyword) {
		keyword = extractOrderKeyword(rawMessage)
	}
	if isGenericKeyword(keyword) {
		return IntentRouteResult{Intent: IntentFallbackAgent}
	}

	return IntentRouteResult{
		Intent:      IntentProductBrowse,
		Confidence:  0.90,
		UseFastPath: true,
		Keyword:     keyword,
		Reason:      "specific_product_discovery",
	}
}
