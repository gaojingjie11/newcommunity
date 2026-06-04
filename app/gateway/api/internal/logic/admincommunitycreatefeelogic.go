package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCommunityCreateFeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCommunityCreateFeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCommunityCreateFeeLogic {
	return &AdminCommunityCreateFeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCommunityCreateFeeLogic) AdminCommunityCreateFee(req *types.CreatePropertyFeeReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.AdminCreateFee(l.ctx, &communityrpc.CreatePropertyFeeReq{
		UserId:  req.UserId,
		Month:   req.Month,
		Amount:  req.Amount,
		DueDate: req.DueDate,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
