package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallReceiveOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallReceiveOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallReceiveOrderLogic {
	return &MallReceiveOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallReceiveOrderLogic) MallReceiveOrder(req *types.ReceiveOrderReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.ReceiveOrder(l.ctx, &mall.OrderIDReq{
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
