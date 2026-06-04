package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/stats/rpc/statsrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsProductSalesRankLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatsProductSalesRankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsProductSalesRankLogic {
	return &StatsProductSalesRankLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatsProductSalesRankLogic) StatsProductSalesRank(req *types.RankLimitReq) (resp *types.SalesRankResp, err error) {
	rpcResp, err := l.svcCtx.StatsRpc.GetProductSalesRank(l.ctx, &statsrpc.RankLimitReq{
		Limit: req.Limit,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.SalesRankInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, toAPISalesRankInfo(item))
	}
	return &types.SalesRankResp{
		List: list,
	}, nil
}
