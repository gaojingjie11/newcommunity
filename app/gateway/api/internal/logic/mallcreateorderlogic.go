package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallCreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallCreateOrderLogic {
	return &MallCreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallCreateOrderLogic) MallCreateOrder(req *types.CreateOrderReq) (resp *types.OrderInfo, err error) {
	rpcResp, err := l.svcCtx.MallRpc.CreateOrder(l.ctx, &mall.CreateOrderReq{
		UserId:  getUserIDFromCtx(l.ctx),
		CartIds: req.CartIds,
		StoreId: req.StoreId,
	})
	if err != nil {
		return nil, err
	}
	info := toAPIOrderInfo(rpcResp)
	return &info, nil
}
