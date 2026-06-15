package logic

import (
	"context"
	"errors"

	"smartcommunity-microservices/app/agent/rpc/agentrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateConversationLogic {
	return &CreateConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateConversationLogic) CreateConversation(req *types.CreateConversationReq) (resp *types.CreateConversationResp, err error) {
	userID := getUserIDFromCtx(l.ctx)
	if userID == 0 {
		return nil, errors.New("请先登录")
	}

	rpcResp, err := l.svcCtx.AgentRpc.CreateConversation(l.ctx, &agentrpc.CreateConversationReq{
		UserId: userID,
		Title:  req.Title,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateConversationResp{
		Id: rpcResp.Id,
	}, nil
}
