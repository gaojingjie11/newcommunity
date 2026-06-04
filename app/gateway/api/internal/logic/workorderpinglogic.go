package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkorderPingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkorderPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkorderPingLogic {
	return &WorkorderPingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkorderPingLogic) WorkorderPing() (resp *types.BaseResp, err error) {
	return &types.BaseResp{
		Code:    0,
		Message: "pong",
	}, nil
}
