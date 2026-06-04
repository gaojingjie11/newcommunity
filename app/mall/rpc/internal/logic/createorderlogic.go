package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/service"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *mall.CreateOrderReq) (*mall.OrderInfo, error) {
	order, err := l.svcCtx.OrderSvc.CreateOrder(in.UserId, service.CreateOrderRequest{
		CartIDs: in.CartIds,
		StoreID: in.StoreId,
	})
	if err != nil {
		return nil, err
	}
	return toProtoOrder(order), nil
}
