package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SmsCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSmsCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SmsCodeLogic {
	return &SmsCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SmsCodeLogic) SmsCode(req *types.SendCodeReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.UserRpc.SendLoginCode(l.ctx, &user.SendCodeReq{
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
