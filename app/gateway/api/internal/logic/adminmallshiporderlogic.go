package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallShipOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallShipOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallShipOrderLogic {
	return &AdminMallShipOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallShipOrderLogic) AdminMallShipOrder(req *types.AdminShipOrderReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.ShipOrder(l.ctx, &mall.ShipOrderReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
