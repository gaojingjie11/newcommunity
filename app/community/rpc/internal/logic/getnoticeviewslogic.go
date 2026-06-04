package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoticeViewsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNoticeViewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoticeViewsLogic {
	return &GetNoticeViewsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNoticeViewsLogic) GetNoticeViews(in *community.NoticeIDReq) (*community.NoticeViewsResp, error) {
	notice, err := l.svcCtx.NoticeRepo.Get(in.Id)
	if err != nil {
		return nil, err
	}

	return &community.NoticeViewsResp{
		Views: notice.ViewCount,
	}, nil
}
