package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommunityPingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommunityPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommunityPingLogic {
	return &CommunityPingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommunityPingLogic) CommunityPing() (resp *types.BaseResp, err error) {
	return &types.BaseResp{
		Code:    0,
		Message: "pong",
	}, nil
}
