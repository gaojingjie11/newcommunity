package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindUserStoresLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindUserStoresLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUserStoresLogic {
	return &BindUserStoresLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindUserStoresLogic) BindUserStores(in *mall.BindUserStoresReq) (*mall.BaseResp, error) {
	err := l.svcCtx.StoreRepo.BindUserStores(in.UserId, in.StoreIds)
	if err != nil {
		return &mall.BaseResp{
			Code:    500,
			Message: err.Error(),
		}, nil
	}
	return &mall.BaseResp{
		Code:    0,
		Message: "success",
	}, nil
}
