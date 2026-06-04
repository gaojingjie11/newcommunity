package logic

import (
	"context"
	"errors"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkNoticeReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkNoticeReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkNoticeReadLogic {
	return &MarkNoticeReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MarkNoticeReadLogic) MarkNoticeRead(in *community.NoticeIDReq) (*community.BaseResp, error) {
	if in.UserId <= 0 {
		return nil, errors.New("invalid user id")
	}

	if _, err := l.svcCtx.NoticeRepo.Get(in.Id); err != nil {
		return nil, err
	}

	if err := l.svcCtx.NoticeRepo.MarkRead(in.Id, in.UserId); err != nil {
		return nil, err
	}

	return &community.BaseResp{Code: 0, Message: "success"}, nil
}
