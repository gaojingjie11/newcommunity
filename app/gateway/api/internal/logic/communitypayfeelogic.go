package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommunityPayFeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommunityPayFeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommunityPayFeeLogic {
	return &CommunityPayFeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommunityPayFeeLogic) CommunityPayFee(req *types.PayPropertyFeeReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.PayPropertyFee(l.ctx, &communityrpc.PayPropertyFeeReq{
		Id:     req.Id,
		UserId: getUserIDFromCtx(l.ctx),
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
