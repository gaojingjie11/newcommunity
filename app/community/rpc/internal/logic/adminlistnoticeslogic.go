package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListNoticesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminListNoticesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListNoticesLogic {
	return &AdminListNoticesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminListNoticesLogic) AdminListNotices(in *community.ListNoticesReq) (*community.NoticeListResp, error) {
	notices, total, err := l.svcCtx.NoticeRepo.List(int(in.Page), int(in.Size), true)
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
