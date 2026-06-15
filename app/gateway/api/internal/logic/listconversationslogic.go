package logic

import (
	"context"
	"errors"

	"smartcommunity-microservices/app/agent/rpc/agentrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListConversationsLogic {
	return &ListConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListConversationsLogic) ListConversations() (resp *types.ListConversationsResp, err error) {
	userID := getUserIDFromCtx(l.ctx)
	if userID == 0 {
		return nil, errors.New("请先登录")
	}

	rpcResp, err := l.svcCtx.AgentRpc.ListConversations(l.ctx, &agentrpc.ListConversationsReq{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.ConversationInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, types.ConversationInfo{
			Id:        item.Id,
			Title:     item.Title,
			Summary:   item.Summary,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return &types.ListConversationsResp{
		List: list,
	}, nil
}
