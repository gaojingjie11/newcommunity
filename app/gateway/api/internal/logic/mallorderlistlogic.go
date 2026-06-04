package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallOrderListLogic {
	return &MallOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallOrderListLogic) MallOrderList(req *types.ListOrdersReq) (resp *types.OrderListResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.ListOrders(l.ctx, &mall.ListOrdersReq{
		UserId: getUserIDFromCtx(l.ctx),
		Status: req.Status,
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.OrderInfo, 0, len(rpcResp.List))
	for _, o := range rpcResp.List {
		list = append(list, toAPIOrderInfo(o))
	}
	return &types.OrderListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
