package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByCodeLogic {
	return &LoginByCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginByCodeLogic) LoginByCode(in *user.LoginByCodeReq) (*user.LoginResp, error) {
	res, err := l.svcCtx.AuthService.LoginByCode(in.Mobile, in.Code, in.ClientIp, in.UserAgent)
	if err != nil {
		return nil, err
	}

	return &user.LoginResp{
		Token:            res.Token,
		UserInfo:         mapUserInfo(res.User),
		IsNewUser:        res.IsNewUser,
		ProfileCompleted: res.ProfileCompleted,
	}, nil
}
