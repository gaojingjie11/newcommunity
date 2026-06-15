package logic

import (
	"context"
	"fmt"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListStoresLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListStoresLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListStoresLogic {
	return &ListStoresLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListStoresLogic) ListStores(in *mall.ListStoresReq) (*mall.StoreListResp, error) {
	fmt.Printf("[DEBUG] ListStores: in.StoreIds=%v, isNil=%t, len=%d\n", in.StoreIds, in.StoreIds == nil, len(in.StoreIds))
	storeIDs := in.StoreIds
	if ids, ok := getStoreIDsFromCtx(l.ctx); ok {
		fmt.Printf("[DEBUG] ListStores: getStoreIDsFromCtx ok=true, ids=%v\n", ids)
		if len(in.StoreIds) > 0 {
			var intersect []int64
			allowed := make(map[int64]bool)
			for _, id := range ids {
				allowed[id] = true
			}
			for _, id := range in.StoreIds {
				if allowed[id] {
					intersect = append(intersect, id)
				}
			}
			if intersect == nil {
				intersect = []int64{}
			}
			storeIDs = intersect
		} else {
			storeIDs = ids
		}
	}
	stores, total, err := l.svcCtx.StoreSvc.List(int(in.Page), int(in.Size), 0, storeIDs)
	if err != nil {
		return nil, err
	}
	var list []*mall.StoreInfo
	for _, s := range stores {
		list = append(list, toProtoStore(&s))
	}
	return &mall.StoreListResp{List: list, Total: total}, nil
}
