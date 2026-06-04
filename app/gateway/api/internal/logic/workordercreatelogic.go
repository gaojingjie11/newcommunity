package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/workorder/rpc/workorderrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkorderCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkorderCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkorderCreateLogic {
	return &WorkorderCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkorderCreateLogic) WorkorderCreate(req *types.CreateWorkorderReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.WorkorderRpc.CreateWorkorder(l.ctx, &workorderrpc.CreateWorkorderReq{
		UserId:      getUserIDFromCtx(l.ctx),
		Type:        req.Type,
		Category:    req.Category,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
