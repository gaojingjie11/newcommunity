package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallOrderDetailLogic {
	return &MallOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallOrderDetailLogic) MallOrderDetail(req *types.OrderIDReq) (resp *types.OrderInfo, err error) {
	rpcResp, err := l.svcCtx.MallRpc.GetOrderDetail(l.ctx, &mall.OrderIDReq{
		Id:     req.Id,
		UserId: getUserIDFromCtx(l.ctx),
	})
	if err != nil {
		return nil, err
	}
	info := toAPIOrderInfo(rpcResp)
	return &info, nil
}
