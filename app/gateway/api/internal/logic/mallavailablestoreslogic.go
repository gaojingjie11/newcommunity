package logic

import (
	"context"
	"fmt"
	"strings"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallAvailableStoresLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallAvailableStoresLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallAvailableStoresLogic {
	return &MallAvailableStoresLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallAvailableStoresLogic) MallAvailableStores(req *types.AvailableStoresReq) (resp *types.StoreListResp, err error) {
	var cartIDs []int64
	for _, idStr := range strings.Split(req.CartIds, ",") {
		if idStr == "" {
			continue
		}
		var id int64
		_, _ = fmt.Sscan(idStr, &id)
		if id > 0 {
			cartIDs = append(cartIDs, id)
		}
	}
	rpcResp, err := l.svcCtx.MallRpc.ListAvailableStores(l.ctx, &mall.ListAvailableStoresReq{
		UserId:  getUserIDFromCtx(l.ctx),
		CartIds: cartIDs,
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
