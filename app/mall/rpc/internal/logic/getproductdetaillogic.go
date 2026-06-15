package logic

import (
	"context"
	"fmt"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetProductDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductDetailLogic {
	return &GetProductDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductDetailLogic) GetProductDetail(in *mall.ProductIDReq) (*mall.ProductDetailResp, error) {
	product, err := l.svcCtx.ProductSvc.GetDetail(in.Id)
	if err != nil {
		return nil, err
	}

	redisKey := fmt.Sprintf("product:view_count:%d", in.Id)
	v, err := l.svcCtx.Redis.Incr(l.ctx, redisKey).Result()
	if err == nil {
		if v == 1 {
			// Initialize Redis with MySQL's view_count + 1
			l.svcCtx.Redis.Set(l.ctx, redisKey, product.ViewCount+1, 0)
			v = product.ViewCount + 1
		}
		product.ViewCount = v

		// Quantitative writeback: every 10 clicks, update DB
		if v%10 == 0 {
			l.svcCtx.DB.Model(&model.Product{}).Where("id = ?", in.Id).UpdateColumn("view_count", v)
		}
	} else {
		// Fallback to direct DB increment if Redis is down
		l.svcCtx.DB.Model(&model.Product{}).Where("id = ?", in.Id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
		product.ViewCount++
	}

	return &mall.ProductDetailResp{Product: toProtoProduct(product)}, nil
}
