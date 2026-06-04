package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FreezeUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFreezeUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FreezeUserLogic {
	return &FreezeUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FreezeUserLogic) FreezeUser(in *user.FreezeUserReq) (*user.BaseResp, error) {
	err := l.svcCtx.AdminService.FreezeUser(in.UserId, int(in.Status))
	if err != nil {
		return nil, err
	}

	return &user.BaseResp{
		Code:    0,
		Message: "用户状态更新成功",
	}, nil
}
