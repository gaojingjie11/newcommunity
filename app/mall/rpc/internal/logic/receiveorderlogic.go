package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReceiveOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReceiveOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveOrderLogic {
	return &ReceiveOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReceiveOrderLogic) ReceiveOrder(in *mall.OrderIDReq) (*mall.BaseResp, error) {
	err := l.svcCtx.OrderSvc.ReceiveOrder(in.Id)
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
