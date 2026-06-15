// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallListStoresLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallListStoresLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallListStoresLogic {
	return &AdminMallListStoresLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallListStoresLogic) AdminMallListStores(req *types.AdminListStoresReq) (resp *types.StoreListResp, err error) {
	page := req.Page
	size := req.Size
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}

	rpcResp, err := l.svcCtx.MallRpc.ListStores(l.ctx, &mall.ListStoresReq{
		Page: page,
		Size: size,
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
