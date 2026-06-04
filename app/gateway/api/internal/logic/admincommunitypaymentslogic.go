package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCommunityPaymentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCommunityPaymentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCommunityPaymentsLogic {
	return &AdminCommunityPaymentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCommunityPaymentsLogic) AdminCommunityPayments(req *types.ListPaymentsReq) (resp *types.PaymentListResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.AdminListPayments(l.ctx, &communityrpc.ListPaymentsReq{
		UserId: 0,
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.PropertyFeePaymentInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, toAPIPropertyFeePaymentInfo(item))
	}
	return &types.PaymentListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
