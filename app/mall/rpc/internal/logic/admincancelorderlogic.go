package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCancelOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCancelOrderLogic {
	return &AdminCancelOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminCancelOrderLogic) AdminCancelOrder(in *mall.AdminCancelOrderReq) (*mall.BaseResp, error) {
	order, err := l.svcCtx.OrderRepo.FindByID(in.Id)
	if err != nil {
		return nil, err
	}
	if err := checkStoreAccess(l.ctx, order.StoreID); err != nil {
		return nil, err
	}

	err = l.svcCtx.OrderSvc.CancelOrder(in.Id, "管理员取消")
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}
