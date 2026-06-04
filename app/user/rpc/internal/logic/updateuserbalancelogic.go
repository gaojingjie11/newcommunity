package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserBalanceLogic {
	return &UpdateUserBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserBalanceLogic) UpdateUserBalance(in *user.UpdateUserBalanceReq) (*user.BaseResp, error) {
	err := l.svcCtx.AdminService.UpdateUserBalance(in.UserId, in.Amount, 0)
	if err != nil {
		return nil, err
	}

	return &user.BaseResp{
		Code:    0,
		Message: "余额更新成功",
	}, nil
}
