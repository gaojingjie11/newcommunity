package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelOrderLogic) CancelOrder(in *mall.CancelOrderReq) (*mall.BaseResp, error) {
	order, err := l.svcCtx.OrderRepo.FindByID(in.Id)
	if err != nil {
		return &mall.BaseResp{Code: 404, Message: "订单不存在"}, nil
	}
	if order.UserID != in.UserId {
		return &mall.BaseResp{Code: 403, Message: "无权取消此订单"}, nil
	}
	err = l.svcCtx.OrderSvc.CancelOrder(in.Id, "用户主动取消")
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
