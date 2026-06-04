package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNoticesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListNoticesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNoticesLogic {
	return &ListNoticesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListNoticesLogic) ListNotices(in *community.ListNoticesReq) (*community.NoticeListResp, error) {
	notices, total, err := l.svcCtx.NoticeRepo.List(int(in.Page), int(in.Size), false)
	if err != nil {
		return nil, err
	}

	var list []*community.NoticeInfo
	for _, n := range notices {
		list = append(list, &community.NoticeInfo{
			Id:        n.ID,
			Title:     n.Title,
			Content:   n.Content,
			Publisher: n.Publisher,
			ViewCount: n.ViewCount,
			Status:    int32(n.Status),
			CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &community.NoticeListResp{
		List:  list,
		Total: total,
	}, nil
}
