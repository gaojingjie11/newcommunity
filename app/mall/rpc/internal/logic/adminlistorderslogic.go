package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListOrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListOrdersLogic {
	return &AdminListOrdersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminListOrdersLogic) AdminListOrders(in *mall.AdminListOrdersReq) (*mall.OrderListResp, error) {
	storeIDs := in.StoreIds
	if ids, ok := getStoreIDsFromCtx(l.ctx); ok {
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
	var status *int
	if in.Status >= 0 {
		val := int(in.Status)
		status = &val
	}
	orders, total, err := l.svcCtx.OrderSvc.AdminListOrders(int(in.Page), int(in.Size), status, in.OrderNo, storeIDs)
	if err != nil {
		return nil, err
	}
	var list []*mall.OrderInfo
	for _, o := range orders {
		list = append(list, toProtoOrder(&o))
	}
	return &mall.OrderListResp{List: list, Total: total}, nil
}
