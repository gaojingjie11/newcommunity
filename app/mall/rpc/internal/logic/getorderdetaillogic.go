package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GetOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderDetailLogic {
	return &GetOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderDetailLogic) GetOrderDetail(in *mall.OrderIDReq) (*mall.OrderInfo, error) {
	order, err := l.svcCtx.OrderRepo.FindByID(in.Id)
	if err != nil {
		return nil, err
	}

	isAdmin := false
	md, ok := metadata.FromIncomingContext(l.ctx)
	if ok {
		vals := md.Get("x-is-admin")
		if len(vals) > 0 && vals[0] == "true" {
			isAdmin = true
		}
	}

	if isAdmin {
		if err := checkStoreAccess(l.ctx, order.StoreID); err != nil {
			return nil, err
		}
	} else {
		if order.UserID != in.UserId {
			return nil, status.Error(codes.PermissionDenied, "无权查看此订单")
		}
	}

	return toProtoOrder(order), nil
}
