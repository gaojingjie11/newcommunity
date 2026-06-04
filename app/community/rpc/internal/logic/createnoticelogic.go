package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/model"
	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateNoticeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNoticeLogic {
	return &CreateNoticeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateNoticeLogic) CreateNotice(in *community.CreateNoticeReq) (*community.BaseResp, error) {
	publisher := in.Publisher
	if publisher == "" {
		publisher = "管理员"
	}

	notice := &model.Notice{
		Title:     in.Title,
		Content:   in.Content,
		Publisher: publisher,
		Status:    1,
	}

	if err := l.svcCtx.NoticeRepo.Create(notice); err != nil {
		return nil, err
	}

	return &community.BaseResp{Code: 0, Message: "success"}, nil
}
