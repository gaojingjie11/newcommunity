package logic

import (
	"context"

	"smartcommunity-microservices/app/stats/rpc/internal/svc"
	"smartcommunity-microservices/app/stats/rpc/types/stats"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductSalesRankLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductSalesRankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductSalesRankLogic {
	return &GetProductSalesRankLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductSalesRankLogic) GetProductSalesRank(in *stats.RankLimitReq) (*stats.SalesRankResp, error) {
	rank, err := l.svcCtx.StatsSvc.ProductSalesRank(int(in.Limit))
	if err != nil {
		return nil, err
	}

	var list []*stats.SalesRankInfo
	for _, r := range rank {
		list = append(list, &stats.SalesRankInfo{
			ProductId:   r.ProductID,
			ProductName: r.ProductName,
			TotalSales:  r.TotalSales,
			TotalAmount: float64(r.TotalAmount) / 100.0,
		})
	}

	return &stats.SalesRankResp{
		List: list,
	}, nil
}
