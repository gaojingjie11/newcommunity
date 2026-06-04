package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStoreDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStoreDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStoreDetailLogic {
	return &GetStoreDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetStoreDetailLogic) GetStoreDetail(in *mall.StoreIDReq) (*mall.StoreDetailResp, error) {
	store, err := l.svcCtx.StoreSvc.GetDetail(in.Id)
	if err != nil {
		return nil, err
	}
	return &mall.StoreDetailResp{Store: toProtoStore(store)}, nil
}
