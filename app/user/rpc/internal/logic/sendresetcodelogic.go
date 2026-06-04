package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendResetCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendResetCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendResetCodeLogic {
	return &SendResetCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendResetCodeLogic) SendResetCode(in *user.SendCodeReq) (*user.BaseResp, error) {
	err := l.svcCtx.AuthService.SendResetCode(in.Mobile)
	if err != nil {
		return nil, err
	}

	return &user.BaseResp{
		Code:    0,
		Message: "重置密码验证码已发送",
	}, nil
}
