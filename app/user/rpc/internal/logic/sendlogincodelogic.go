package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendLoginCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendLoginCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendLoginCodeLogic {
	return &SendLoginCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendLoginCodeLogic) SendLoginCode(in *user.SendCodeReq) (*user.BaseResp, error) {
	_, err := l.svcCtx.AuthService.SendLoginCode(in.Mobile)
	if err != nil {
		return nil, err
	}

	return &user.BaseResp{
		Code:    0,
		Message: "验证码发送成功",
	}, nil
}
