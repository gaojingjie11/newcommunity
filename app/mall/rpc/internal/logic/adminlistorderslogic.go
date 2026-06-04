package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListOrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListOrdersLogic {
	return &AdminListOrdersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminListOrdersLogic) AdminListOrders(in *mall.AdminListOrdersReq) (*mall.OrderListResp, error) {
	var status *int
	if in.Status >= 0 {
		val := int(in.Status)
		status = &val
	}
	orders, total, err := l.svcCtx.OrderSvc.AdminListOrders(int(in.Page), int(in.Size), status, in.OrderNo)
	if err != nil {
		return nil, err
	}
	var list []*mall.OrderInfo
	for _, o := range orders {
		list = append(list, toProtoOrder(&o))
	}
	return &mall.OrderListResp{List: list, Total: total}, nil
}
