package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	token, u, err := l.svcCtx.AuthService.Login(in.Mobile, in.Password, in.ClientIp, in.UserAgent)
	if err != nil {
		return nil, err
	}

	profileCompleted := !(u.RealName == "" || u.RealName == "未完善资料" || u.Age <= 1)

	roles, _ := l.svcCtx.AdminService.ListRoles()
	roleMap := make(map[string]string)
	for _, r := range roles {
		roleMap[r.Code] = r.Name
	}

	return &user.LoginResp{
		Token:            token,
		UserInfo:         mapUserInfo(u, roleMap),
		IsNewUser:        false,
		ProfileCompleted: profileCompleted,
	}, nil
}
