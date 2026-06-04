package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMessagesLogic {
	return &ListMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMessagesLogic) ListMessages(in *community.ListMessagesReq) (*community.MessageListResp, error) {
	messages, total, err := l.svcCtx.MessageRepo.List(int(in.Page), int(in.Size))
	if err != nil {
		return nil, err
	}

	var list []*community.MessageInfo
	for _, m := range messages {
		username := ""
		avatar := ""
		if m.User != nil {
			username = m.User.Username
			avatar = m.User.Avatar
		}

		list = append(list, &community.MessageInfo{
			Id:        m.ID,
			UserId:    m.UserID,
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
			Username:  username,
			Avatar:    avatar,
		})
	}

	return &community.MessageListResp{
		List:  list,
		Total: total,
	}, nil
}
