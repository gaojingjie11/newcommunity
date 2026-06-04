package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPromotionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPromotionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPromotionsLogic {
	return &ListPromotionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPromotionsLogic) ListPromotions(in *mall.ListPromotionsReq) (*mall.PromotionListResp, error) {
	promos, total, err := l.svcCtx.PromotionSvc.List(int(in.Page), int(in.Size))
	if err != nil {
		return nil, err
	}
	var list []*mall.PromotionInfo
	for _, p := range promos {
		list = append(list, toProtoPromotion(&p))
	}
	return &mall.PromotionListResp{List: list, Total: total}, nil
}
