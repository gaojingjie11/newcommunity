package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNoticeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNoticeLogic {
	return &DeleteNoticeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteNoticeLogic) DeleteNotice(in *community.NoticeIDReq) (*community.BaseResp, error) {
	if err := l.svcCtx.NoticeRepo.Delete(in.Id); err != nil {
		return nil, err
	}

	return &community.BaseResp{Code: 0, Message: "success"}, nil
}
