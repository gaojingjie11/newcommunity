package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordCodeLogic {
	return &ResetPasswordCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordCodeLogic) ResetPasswordCode(req *types.SendCodeReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.SendResetCode(l.ctx, &user.SendCodeReq{
		Mobile: req.Mobile,
	})
	if err != nil {
		return nil, err
	}

	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
