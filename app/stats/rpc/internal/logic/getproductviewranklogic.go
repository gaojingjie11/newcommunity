package logic

import (
	"context"

	"smartcommunity-microservices/app/stats/rpc/internal/svc"
	"smartcommunity-microservices/app/stats/rpc/types/stats"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductViewRankLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductViewRankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductViewRankLogic {
	return &GetProductViewRankLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductViewRankLogic) GetProductViewRank(in *stats.RankLimitReq) (*stats.ViewRankResp, error) {
	rank, err := l.svcCtx.StatsSvc.ProductViewRank(int(in.Limit))
	if err != nil {
		return nil, err
	}

	var list []*stats.ViewRankInfo
	for _, r := range rank {
		list = append(list, &stats.ViewRankInfo{
			ProductId:   r.ProductID,
			ProductName: r.ProductName,
			ViewCount:   r.ViewCount,
			UniqueUsers: r.UniqueUsers,
		})
	}

	return &stats.ViewRankResp{
		List: list,
	}, nil
}
