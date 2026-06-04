package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Profile & Face
func (l *GetProfileLogic) GetProfile(in *user.UserIDReq) (*user.UserInfo, error) {
	u, err := l.svcCtx.UserService.GetProfile(in.UserId)
	if err != nil {
		return nil, err
	}

	return mapUserInfo(u), nil
}
