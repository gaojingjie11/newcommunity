//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/zeromicro/go-zero/core/conf"
	"smartcommunity-microservices/app/agent/rpc/internal/config"
	"smartcommunity-microservices/app/agent/rpc/internal/logic"
	"smartcommunity-microservices/app/agent/rpc/internal/svc"

	"github.com/cloudwego/eino/schema"
)

func main() {
	configFile := "etc/agent.yaml"
	var c config.Config

	err := conf.Load(configFile, &c, conf.UseEnv())
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	svcCtx := svc.NewServiceContext(c)

	ctx := context.Background()

	streamCallback := func(eventType string, payload map[string]interface{}) {
		fmt.Printf("[StreamCallback] Event: %s, Payload: %+v\n", eventType, payload)
	}

	ctx = context.WithValue(ctx, logic.CtxKeyUserID, int64(102))
	ctx = context.WithValue(ctx, logic.CtxKeyConversationID, "3756e197-dfb0-4229-9201-79b88f005552")
	ctx = context.WithValue(ctx, logic.CtxKeyStreamCallback, logic.StreamCallback(streamCallback))

	agent, err := logic.BuildEinoAgent(ctx, svcCtx, "fast")
	if err != nil {
		log.Fatalf("failed to build agent: %v", err)
	}

	var einoMessages []*schema.Message
	einoMessages = append(einoMessages, schema.SystemMessage(logic.SystemPrompt))

	einoMessages = append(einoMessages, schema.AssistantMessage("已为您取消了**物业工单提交**操作。", nil))
	einoMessages = append(einoMessages, schema.UserMessage("你看看商场又撒谎"))
	einoMessages = append(einoMessages, schema.AssistantMessage("您的投诉工单已提交成功！\n\n工单类型：投诉\n分类：虚假宣传\n描述：商场又撒谎\n\n请在弹出的工单卡片中确认内容，我们会尽快处理您的投诉。", nil))
	einoMessages = append(einoMessages, schema.UserMessage("看看商城有什么商品"))

	fmt.Println("Starting local ReAct Agent execution...")

	sr, err := agent.Stream(ctx, einoMessages)
	if err != nil {
		log.Fatalf("agent execution failed: %v", err)
	}
	defer sr.Close()

	for {
		chunk, err := sr.Recv()
		if err != nil {
			fmt.Printf("\nStream finished/err: %v\n", err)
			break
		}
		fmt.Printf("Chunk: %q\n", chunk.Content)
	}
}
