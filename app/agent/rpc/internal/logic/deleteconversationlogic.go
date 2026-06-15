package logic

import (
	"context"

	"smartcommunity-microservices/app/agent/rpc/agent"
	"smartcommunity-microservices/app/agent/rpc/internal/model"
	"smartcommunity-microservices/app/agent/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConversationLogic {
	return &DeleteConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteConversationLogic) DeleteConversation(in *agent.DeleteConversationReq) (*agent.BaseResp, error) {
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Delete messages
		if err := tx.Where("conversation_id = ? AND user_id = ?", in.Id, in.UserId).
			Delete(&model.SysUserChatMessage{}).Error; err != nil {
			return err
		}
		// 2. Delete conversation
		if err := tx.Where("id = ? AND user_id = ?", in.Id, in.UserId).
			Delete(&model.SysUserConversation{}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		l.Errorf("failed to delete conversation %s: %v", in.Id, err)
		return &agent.BaseResp{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &agent.BaseResp{
		Code:    0,
		Message: "success",
	}, nil
}

