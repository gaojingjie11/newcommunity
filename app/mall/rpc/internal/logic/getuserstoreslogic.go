package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserStoresLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserStoresLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserStoresLogic {
	return &GetUserStoresLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserStoresLogic) GetUserStores(in *mall.UserIDReq) (*mall.StoreIDListResp, error) {
	storeIDs, err := l.svcCtx.StoreRepo.GetUserStores(in.UserId)
	if err != nil {
		return nil, err
	}
	return &mall.StoreIDListResp{
		StoreIds: storeIDs,
	}, nil
}
