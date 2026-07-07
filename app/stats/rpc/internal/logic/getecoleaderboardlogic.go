package logic

import (
	"context"

	"smartcommunity-microservices/app/stats/rpc/internal/svc"
	"smartcommunity-microservices/app/stats/rpc/types/stats"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEcoLeaderboardLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEcoLeaderboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEcoLeaderboardLogic {
	return &GetEcoLeaderboardLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEcoLeaderboardLogic) GetEcoLeaderboard(in *stats.BaseResp) (*stats.EcoStatsResp, error) {
	list, total, err := l.svcCtx.StatsSvc.GreenPointsLeaderboard(15)
	if err != nil {
		return nil, err
	}

	var pbList []*stats.EcoLeaderboardInfo
	for _, item := range list {
		nickname := item.RealName
		if nickname == "" {
			nickname = item.Username
		}
		pbList = append(pbList, &stats.EcoLeaderboardInfo{
			UserId:      item.UserID,
			Username:    item.Username,
			RealName:    item.RealName,
			Nickname:    nickname,
			Points:      item.GreenPoints,
			Avatar:      item.Avatar,
		})
	}

	return &stats.EcoStatsResp{
		List:              pbList,
		TotalPointsIssued: total,
	}, nil
}
