package logic

import (
	"context"
	"errors"

	"smartcommunity-microservices/app/agent/rpc/agentrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RejectActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRejectActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectActionLogic {
	return &RejectActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RejectActionLogic) RejectAction(req *types.RejectActionReq) (resp *types.BaseResp, err error) {
	userID := getUserIDFromCtx(l.ctx)
	if userID == 0 {
		return nil, errors.New("请先登录")
	}

	rpcResp, err := l.svcCtx.AgentRpc.RejectAction(l.ctx, &agentrpc.RejectActionReq{
		UserId:         userID,
		ConversationId: req.ConversationId,
		ActionId:       req.ActionId,
	})
	if err != nil {
		return nil, err
	}

	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
