package logic

import (
	"context"
	"fmt"
	"time"

	"smartcommunity-microservices/app/agent/rpc/agent"
	"smartcommunity-microservices/app/agent/rpc/internal/model"
	"smartcommunity-microservices/app/agent/rpc/internal/svc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type RejectActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRejectActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectActionLogic {
	return &RejectActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RejectActionLogic) RejectAction(in *agent.RejectActionReq) (*agent.BaseResp, error) {
	// 1. Fetch approval record
	var approval model.AgentActionApproval
	err := l.svcCtx.DB.Where("id = ? AND conversation_id = ? AND user_id = ?", in.ActionId, in.ConversationId, in.UserId).First(&approval).Error
	if err != nil {
		l.Errorf("approval action not found: %v", err)
		return &agent.BaseResp{Code: 404, Message: "审批动作未找到"}, nil
	}

	if approval.Status != "pending" {
		return &agent.BaseResp{Code: 400, Message: fmt.Sprintf("无法拒绝该审批动作（当前状态: %s）", approval.Status)}, nil
	}

	// 2. Mark status as rejected
	approval.Status = "rejected"
	approval.UpdatedAt = time.Now()
	l.svcCtx.DB.Save(&approval)

	// 3. Save assistant cancelled message
	actionName := "操作"
	switch approval.ActionType {
	case "create_order":
		actionName = "商品下单"
	case "pay_order":
		actionName = "订单支付"
	case "submit_repair":
		actionName = "物业工单提交"
	}
	resultMsg := fmt.Sprintf("已为您取消了**%s**操作。", actionName)
	l.saveChatMessage(in.ConversationId, in.UserId, "assistant", resultMsg)

	return &agent.BaseResp{Code: 0, Message: "审批已成功拒绝"}, nil
}

func (l *RejectActionLogic) saveChatMessage(convID string, userID int64, role, content string) {
	msg := &model.SysUserChatMessage{
		ID:             uuid.NewString(),
		UserID:         userID,
		ConversationID: convID,
		Role:           role,
		Content:        content,
		CreatedAt:      time.Now(),
	}
	l.svcCtx.DB.Create(msg)

	l.svcCtx.DB.Model(&model.SysUserConversation{}).
		Where("id = ? AND user_id = ?", convID, userID).
		Update("updated_at", time.Now())
}

