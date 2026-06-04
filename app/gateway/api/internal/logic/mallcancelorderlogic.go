package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallCancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallCancelOrderLogic {
	return &MallCancelOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallCancelOrderLogic) MallCancelOrder(req *types.CancelOrderReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.CancelOrder(l.ctx, &mall.CancelOrderReq{
		Id:     req.Id,
		UserId: getUserIDFromCtx(l.ctx),
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
