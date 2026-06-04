package logic

import (
	"context"

	"smartcommunity-microservices/app/stats/rpc/internal/svc"
	"smartcommunity-microservices/app/stats/rpc/types/stats"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkorderStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWorkorderStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkorderStatsLogic {
	return &GetWorkorderStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWorkorderStatsLogic) GetWorkorderStats(in *stats.BaseResp) (*stats.WorkorderStatsResp, error) {
	summary, err := l.svcCtx.StatsSvc.WorkorderSummary()
	if err != nil {
		return nil, err
	}

	var pbSummary []*stats.WorkorderSummaryInfo
	for _, s := range summary {
		pbSummary = append(pbSummary, &stats.WorkorderSummaryInfo{
			Type:   s.Type,
			Status: int32(s.Status),
			Count:  s.Count,
		})
	}

	return &stats.WorkorderStatsResp{
		Summaries: pbSummary,
	}, nil
}
