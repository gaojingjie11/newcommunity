package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/stats/rpc/statsrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsProductViewRankLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatsProductViewRankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsProductViewRankLogic {
	return &StatsProductViewRankLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatsProductViewRankLogic) StatsProductViewRank(req *types.RankLimitReq) (resp *types.ViewRankResp, err error) {
	rpcResp, err := l.svcCtx.StatsRpc.GetProductViewRank(l.ctx, &statsrpc.RankLimitReq{
		Limit: req.Limit,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.ViewRankInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, toAPIViewRankInfo(item))
	}
	return &types.ViewRankResp{
		List: list,
	}, nil
}
