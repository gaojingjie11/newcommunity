package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/service"
	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Authentication & Session
func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	u, err := l.svcCtx.AuthService.Register(service.RegisterRequest{
		Mobile:   in.Mobile,
		Password: in.Password,
	})
	if err != nil {
		return nil, err
	}

	return &user.RegisterResp{
		UserId: u.ID,
	}, nil
}
