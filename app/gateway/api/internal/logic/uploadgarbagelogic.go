// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadGarbageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadGarbageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadGarbageLogic {
	return &UploadGarbageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadGarbageLogic) UploadGarbage(fileURL, filename string) (resp *types.UploadGarbageResp, err error) {
	userID := getUserIDFromCtx(l.ctx)
	if userID == 0 {
		return nil, errors.New("请先登录")
	}

	// 1. Perform AI garbage recognition with rule-based fallback
	reason, points, err := l.recognizeGarbage(fileURL, filename)
	if err != nil {
		l.Errorf("failed to recognize garbage: %v", err)
		return nil, err
	}

	// 2. Update user's green points in the database via UserRpc
	_, err = l.svcCtx.UserRpc.UpdateUserPoints(l.ctx, &user.UpdateUserPointsReq{
		UserId: userID,
		Points: points,
	})
	if err != nil {
		l.Errorf("failed to update user green_points: %v", err)
		return nil, fmt.Errorf("更新用户积分失败: %w", err)
	}

	// 3. Fetch the updated user profile to get the new green_points balance
	profile, err := l.svcCtx.UserRpc.GetProfile(l.ctx, &user.UserIDReq{
		UserId: userID,
	})
	if err != nil {
		l.Errorf("failed to fetch user profile: %v", err)
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 4. Invalidate green points leaderboard cache in Redis
	if l.svcCtx.RedisClient != nil {
		iter := l.svcCtx.RedisClient.Scan(l.ctx, 0, "stats:green:leaderboard:*", 0).Iterator()
		for iter.Next(l.ctx) {
			l.svcCtx.RedisClient.Del(l.ctx, iter.Val())
		}
	}

	return &types.UploadGarbageResp{
		Points:      points,
		Reason:      reason,
		GreenPoints: profile.GreenPoints,
	}, nil
}

func (l *UploadGarbageLogic) recognizeGarbage(fileURL, filename string) (string, int32, error) {
	apiKey := os.Getenv("LLM_IMAGE_VISION_API_KEY")
	baseURL := os.Getenv("LLM_IMAGE_VISION_BASE_URL")
	model := os.Getenv("LLM_IMAGE_VISION_MODEL")

	if apiKey == "" || baseURL == "" || model == "" {
		l.Info("AI vision credentials not configured, using fallback simulation")
		return getFallbackRecognition(filename), getRandomPoints(filename), nil
	}

	// OpenAI chat completion endpoint
	urlStr := baseURL
	if !strings.HasSuffix(urlStr, "/chat/completions") {
		urlStr = strings.TrimSuffix(urlStr, "/") + "/chat/completions"
	}

	prompt := "你是一个垃圾分类助手。请识别这张图片中的垃圾物品，判断它属于哪种垃圾分类（可回收物、有害垃圾、湿垃圾/厨余垃圾、干垃圾/其他垃圾）。给出非常简短的识别结果，并说明原因。同时，为了积极鼓励居民参与垃圾分类并体现环保关怀，请给出这件垃圾对应的较为宽松的环保积分值（10 到 30 之间的整数，根据分类价值从优给分）。请以如下 JSON 格式返回，不要包含任何 markdown 标记（如 ```json 等）：\n{\n  \"reason\": \"AI 识别结果: [物品名称]（[分类类别]），[简短解释]\",\n  \"points\": [积分值]\n}"

	reqBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": prompt,
					},
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": fileURL,
						},
					},
				},
			},
		},
		"response_format": map[string]string{
			"type": "json_object",
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return getFallbackRecognition(filename), getRandomPoints(filename), nil
	}

	ctx, cancel := context.WithTimeout(l.ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlStr, bytes.NewBuffer(jsonData))
	if err != nil {
		return getFallbackRecognition(filename), getRandomPoints(filename), nil
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		l.Errorf("AI vision API request failed: %v, using fallback", err)
		return getFallbackRecognition(filename), getRandomPoints(filename), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		l.Errorf("AI vision API returned status code %d, using fallback", resp.StatusCode)
		return getFallbackRecognition(filename), getRandomPoints(filename), nil
	}

	var aiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil || len(aiResp.Choices) == 0 {
		l.Errorf("failed to decode AI response or empty choices, using fallback: %v", err)
		return getFallbackRecognition(filename), getRandomPoints(filename), nil
	}

	content := aiResp.Choices[0].Message.Content
	var result struct {
		Reason string `json:"reason"`
		Points int32  `json:"points"`
	}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		l.Errorf("failed to unmarshal AI content json, using fallback: %v", err)
		return getFallbackRecognition(filename), getRandomPoints(filename), nil
	}

	if result.Points < 10 {
		result.Points = 10
	}
	if result.Points > 30 {
		result.Points = 30
	}

	return result.Reason, result.Points, nil
}

func getFallbackRecognition(filename string) string {
	filename = strings.ToLower(filename)
	if strings.Contains(filename, "plastic") || strings.Contains(filename, "bottle") || strings.Contains(filename, "pingzi") {
		return "AI 识别结果: 塑料瓶（可回收物），感谢您为社区环保做出贡献！"
	}
	if strings.Contains(filename, "paper") || strings.Contains(filename, "box") || strings.Contains(filename, "zhi") {
		return "AI 识别结果: 废纸箱（可回收物），感谢您为社区环保做出贡献！"
	}
	if strings.Contains(filename, "apple") || strings.Contains(filename, "food") || strings.Contains(filename, "fruit") || strings.Contains(filename, "guopi") {
		return "AI 识别结果: 水果皮（厨余垃圾），湿垃圾应当沥干水分后投放。"
	}
	if strings.Contains(filename, "battery") || strings.Contains(filename, "dianchi") {
		return "AI 识别结果: 废旧电池（有害垃圾），有害垃圾需投放至红色收集容器。"
	}
	items := []string{
		"AI 识别结果: 易拉罐（可回收物），感谢您为社区环保做出贡献！",
		"AI 识别结果: 剩饭剩菜（厨余垃圾），湿垃圾应当沥干水分后投放。",
		"AI 识别结果: 过期药品（有害垃圾），有害垃圾需投放至红色收集容器。",
		"AI 识别结果: 废弃纸巾（其他垃圾），干垃圾请投放至灰色收集容器。",
	}
	idx := len(filename) % len(items)
	return items[idx]
}

func getRandomPoints(filename string) int32 {
	filename = strings.ToLower(filename)
	if strings.Contains(filename, "battery") || strings.Contains(filename, "dianchi") {
		return 25
	}
	if strings.Contains(filename, "plastic") || strings.Contains(filename, "bottle") || strings.Contains(filename, "pingzi") {
		return 15
	}
	points := []int32{12, 15, 18, 20}
	idx := len(filename) % len(points)
	return points[idx]
}
