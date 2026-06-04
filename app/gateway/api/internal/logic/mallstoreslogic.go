package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallStoresLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallStoresLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallStoresLogic {
	return &MallStoresLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallStoresLogic) MallStores() (resp *types.StoreListResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.ListStores(l.ctx, &mall.ListStoresReq{
		Page: 1,
		Size: 1000,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.StoreInfo, 0, len(rpcResp.List))
	for _, s := range rpcResp.List {
		list = append(list, toAPIStoreInfo(s))
	}
	return &types.StoreListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
