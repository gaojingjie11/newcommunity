// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/stats/rpc/statsrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGreenPointsLeaderboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGreenPointsLeaderboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGreenPointsLeaderboardLogic {
	return &GetGreenPointsLeaderboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGreenPointsLeaderboardLogic) GetGreenPointsLeaderboard(req *types.RankLimitReq) (resp *types.EcoStatsResp, err error) {
	rpcResp, err := l.svcCtx.StatsRpc.GetEcoLeaderboard(l.ctx, &statsrpc.BaseResp{})
	if err != nil {
		return nil, err
	}

	var list []types.EcoLeaderboardInfo
	for _, item := range rpcResp.List {
		list = append(list, types.EcoLeaderboardInfo{
			UserId:   item.UserId,
			Username: item.Username,
			RealName: item.RealName,
			Nickname: item.Nickname,
			Points:   item.Points,
		})
	}

	return &types.EcoStatsResp{
		List:              list,
		TotalPointsIssued: rpcResp.TotalPointsIssued,
	}, nil
}
