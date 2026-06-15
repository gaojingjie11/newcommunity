package logic

import (
	"context"
	"errors"

	"smartcommunity-microservices/app/agent/rpc/agentrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatHistoryLogic {
	return &GetChatHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatHistoryLogic) GetChatHistory(req *types.ConversationIDReq) (resp *types.ChatHistoryResp, err error) {
	userID := getUserIDFromCtx(l.ctx)
	if userID == 0 {
		return nil, errors.New("请先登录")
	}

	rpcResp, err := l.svcCtx.AgentRpc.GetChatHistory(l.ctx, &agentrpc.GetChatHistoryReq{
		UserId:         userID,
		ConversationId: req.Id,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.ChatMessageInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, types.ChatMessageInfo{
			Id:             item.Id,
			Role:           item.Role,
			Content:        item.Content,
			EventType:      item.EventType,
			EventPayload:   item.EventPayload,
			ActionResolved: item.ActionResolved,
			ResultPayload:  item.ResultPayload,
			CreatedAt:      item.CreatedAt,
		})
	}

	return &types.ChatHistoryResp{
		List: list,
	}, nil
}
