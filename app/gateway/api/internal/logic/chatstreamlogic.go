package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"smartcommunity-microservices/app/agent/rpc/agentrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type ChatStreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatStreamLogic {
	return &ChatStreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatStreamLogic) ChatStream(w http.ResponseWriter, req *types.ChatStreamReq) error {
	userID := getUserIDFromCtx(l.ctx)
	if userID == 0 {
		return errors.New("请先登录")
	}

	// 1. Set SSE Headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return errors.New("streaming unsupported by client connection")
	}

	// 2. Call Agent RPC stream
	rpcCtx := l.ctx
	if mode := strings.TrimSpace(req.Mode); mode != "" {
		rpcCtx = metadata.AppendToOutgoingContext(rpcCtx, "x-agent-mode", mode)
	}

	streamClient, err := l.svcCtx.AgentRpc.ChatStream(rpcCtx, &agentrpc.ChatReq{
		UserId:          userID,
		ConversationId:  req.ConversationId,
		Message:         req.Message,
		PayType:         req.PayType,
		PaymentPassword: req.PaymentPassword,
		FaceImageUrl:    req.FaceImageUrl,
	})
	if err != nil {
		l.Errorf("failed to call AgentRpc ChatStream: %v", err)
		return err
	}

	// 3. Loop and pipe responses to SSE
	for {
		chunk, errRecv := streamClient.Recv()
		if errors.Is(errRecv, io.EOF) {
			break
		}
		if errRecv != nil {
			l.Errorf("error reading from AgentRpc stream: %v", errRecv)
			_, _ = fmt.Fprintf(w, "data: [ERROR] %v\n\n", errRecv)
			flusher.Flush()
			return errRecv
		}

		var payload []byte
		var errMarshal error

		if chunk.EventType != "" {
			var eventData interface{}
			if chunk.EventPayload != "" {
				_ = json.Unmarshal([]byte(chunk.EventPayload), &eventData)
			}
			payload, errMarshal = json.Marshal(map[string]interface{}{
				"type": chunk.EventType,
				"data": eventData,
			})
		} else {
			payload, errMarshal = json.Marshal(map[string]interface{}{
				"type": "message_delta",
				"data": map[string]string{
					"chunk": chunk.Chunk,
				},
			})
		}

		if errMarshal != nil {
			continue
		}

		_, errWrite := fmt.Fprintf(w, "data: %s\n\n", string(payload))
		if errWrite != nil {
			l.Errorf("failed to write to HTTP client: %v", errWrite)
			return errWrite
		}
		flusher.Flush()
	}

	// Write EOF marker
	_, _ = fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()

	return nil
}
