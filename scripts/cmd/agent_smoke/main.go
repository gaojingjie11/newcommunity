package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	pb "smartcommunity-microservices/app/agent/rpc/agent"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type testCase struct {
	Name    string
	UserID  int64
	Mode    string
	Message string
}

func main() {
	addr := flag.String("addr", "127.0.0.1:9006", "agent rpc address")
	only := flag.String("only", "", "comma-separated case names")
	timeout := flag.Duration("timeout", 25*time.Second, "per-case timeout")
	flag.Parse()

	cases := []testCase{
		{
			Name:    "latest_notices_direct",
			UserID:  5,
			Mode:    "fast",
			Message: "最近有什么公告",
		},
		{
			Name:    "semantic_notice_rag",
			UserID:  5,
			Mode:    "smart",
			Message: "最近有没有停水通知",
		},
		{
			Name:    "admin_ai_report_strict",
			UserID:  2,
			Mode:    "smart",
			Message: "请只从最近一期AI报告中总结社区运营风险和建议，不要查公告。",
		},
		{
			Name:    "user_ai_report_strict",
			UserID:  5,
			Mode:    "smart",
			Message: "请只从最近一期AI报告中总结社区运营风险和建议，不要查公告。",
		},
		{
			Name:    "user_ai_report_broad",
			UserID:  5,
			Mode:    "smart",
			Message: "请在AI报告里检索并总结最近一份运营报表的核心风险和建议",
		},
		{
			Name:    "list_products",
			UserID:  5,
			Mode:    "smart",
			Message: "帮我推荐一些便利店商品",
		},
		{
			Name:    "query_orders",
			UserID:  5,
			Mode:    "fast",
			Message: "查一下我最近的订单状态",
		},
		{
			Name:    "submit_repair",
			UserID:  5,
			Mode:    "smart",
			Message: "我要报修，二楼路灯坏了",
		},
		{
			Name:    "submit_complaint",
			UserID:  5,
			Mode:    "smart",
			Message: "帮我投诉下小区的垃圾桶，最近经常爆满没人清理。",
		},
		{
			Name:    "buy_apples",
			UserID:  5,
			Mode:    "fast",
			Message: "帮我买苹果1kg，直接下单",
		},
		{
			Name:    "want_oranges",
			UserID:  5,
			Mode:    "smart",
			Message: "我想买橘子",
		},
		{
			Name:    "buy_cup_phrase",
			UserID:  5,
			Mode:    "fast",
			Message: "帮我买个水杯吧",
		},
		{
			Name:    "shopping_capability_1",
			UserID:  5,
			Mode:    "fast",
			Message: "那你能帮我买东西不",
		},
		{
			Name:    "shopping_capability_2",
			UserID:  5,
			Mode:    "fast",
			Message: "你可以帮我买东西吗",
		},
		{
			Name:    "banana_browse_1",
			UserID:  5,
			Mode:    "smart",
			Message: "商场有没有香蕉呢",
		},
		{
			Name:    "banana_browse_2",
			UserID:  5,
			Mode:    "smart",
			Message: "商场有没有卖香蕉呢",
		},
	}

	selected := buildSelection(*only)

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewAgentRpcClient(conn)

	for _, tc := range cases {
		if len(selected) > 0 && !selected[tc.Name] {
			continue
		}
		runCase(client, tc, *timeout)
	}
}

func buildSelection(raw string) map[string]bool {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	selected := make(map[string]bool, len(parts))
	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name != "" {
			selected[name] = true
		}
	}
	return selected
}

func runCase(client pb.AgentRpcClient, tc testCase, timeout time.Duration) {
	fmt.Printf("=== CASE %s ===\n", tc.Name)
	fmt.Printf("USER=%d MODE=%s MESSAGE=%s\n", tc.UserID, tc.Mode, tc.Message)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if tc.Mode != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-agent-mode", tc.Mode)
	}

	start := time.Now()
	stream, err := client.ChatStream(ctx, &pb.ChatReq{
		UserId:  tc.UserID,
		Message: tc.Message,
	})
	if err != nil {
		fmt.Printf("STATUS=CALL_ERROR DURATION_MS=%d ERROR=%v\n\n", time.Since(start).Milliseconds(), err)
		return
	}

	var (
		reply      strings.Builder
		eventTypes []string
	)

	for {
		chunk, recvErr := stream.Recv()
		if recvErr == io.EOF {
			fmt.Printf("STATUS=OK DURATION_MS=%d EVENTS=%s\n", time.Since(start).Milliseconds(), strings.Join(eventTypes, ","))
			fmt.Printf("REPLY=%s\n\n", compactReply(reply.String()))
			return
		}
		if recvErr != nil {
			fmt.Printf("STATUS=STREAM_ERROR DURATION_MS=%d EVENTS=%s ERROR=%v\n", time.Since(start).Milliseconds(), strings.Join(eventTypes, ","), recvErr)
			fmt.Printf("REPLY_PARTIAL=%s\n\n", compactReply(reply.String()))
			return
		}

		if chunk.EventType != "" {
			eventTypes = append(eventTypes, chunk.EventType)
		}
		reply.WriteString(chunk.Chunk)
	}
}

func compactReply(reply string) string {
	reply = strings.TrimSpace(reply)
	reply = strings.ReplaceAll(reply, "\n", " ")
	reply = strings.Join(strings.Fields(reply), " ")
	if len([]rune(reply)) > 220 {
		runes := []rune(reply)
		return string(runes[:220]) + "..."
	}
	return reply
}
