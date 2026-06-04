package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommunityMyPaymentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommunityMyPaymentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommunityMyPaymentsLogic {
	return &CommunityMyPaymentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommunityMyPaymentsLogic) CommunityMyPayments(req *types.ListPaymentsReq) (resp *types.PaymentListResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.ListMyPayments(l.ctx, &communityrpc.ListPaymentsReq{
		UserId: getUserIDFromCtx(l.ctx),
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
