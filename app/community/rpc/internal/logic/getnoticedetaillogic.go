package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoticeDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNoticeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoticeDetailLogic {
	return &GetNoticeDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNoticeDetailLogic) GetNoticeDetail(in *community.NoticeIDReq) (*community.NoticeDetailResp, error) {
	notice, err := l.svcCtx.NoticeRepo.View(l.ctx, l.svcCtx.Redis, in.Id)
	if err != nil {
		return nil, err
	}

	return &community.NoticeDetailResp{
		Notice: &community.NoticeInfo{
			Id:        notice.ID,
			Title:     notice.Title,
			Content:   notice.Content,
			Publisher: notice.Publisher,
			ViewCount: notice.ViewCount,
			Status:    int32(notice.Status),
			CreatedAt: notice.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}
