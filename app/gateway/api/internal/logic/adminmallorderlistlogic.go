package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallOrderListLogic {
	return &AdminMallOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallOrderListLogic) AdminMallOrderList(req *types.AdminListOrdersReq) (resp *types.OrderListResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.AdminListOrders(l.ctx, &mall.AdminListOrdersReq{
		Page:    req.Page,
		Size:    req.Size,
		OrderNo: req.OrderNo,
		Status:  req.Status,
	})
	if err != nil {
		return nil, err
	}
	
	userMobileCache := make(map[int64]string)
	list := make([]types.OrderInfo, 0, len(rpcResp.List))
	for _, o := range rpcResp.List {
		info := toAPIOrderInfo(o)
		if o.UserId > 0 {
			if mobile, cached := userMobileCache[o.UserId]; cached {
				info.UserMobile = mobile
			} else {
				if profile, err := l.svcCtx.UserRpc.GetProfile(l.ctx, &user.UserIDReq{UserId: o.UserId}); err == nil && profile != nil {
					userMobileCache[o.UserId] = profile.Mobile
					info.UserMobile = profile.Mobile
				} else {
					userMobileCache[o.UserId] = ""
				}
			}
		}
		list = append(list, info)
	}
	
	return &types.OrderListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
