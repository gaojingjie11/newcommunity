package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrdersLogic {
	return &ListOrdersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListOrdersLogic) ListOrders(in *mall.ListOrdersReq) (*mall.OrderListResp, error) {
	var status *int
	if in.Status >= 0 {
		val := int(in.Status)
		status = &val
	}
	orders, total, err := l.svcCtx.OrderSvc.ListOrders(in.UserId, int(in.Page), int(in.Size), status)
	if err != nil {
		return nil, err
	}
	var list []*mall.OrderInfo
	for _, o := range orders {
		list = append(list, toProtoOrder(&o))
	}
	return &mall.OrderListResp{List: list, Total: total}, nil
}
