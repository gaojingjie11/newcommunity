package logic

import (
	"context"

	"smartcommunity-microservices/app/user/rpc/internal/model"
	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateUserPointsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserPointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserPointsLogic {
	return &UpdateUserPointsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Points Management
func (l *UpdateUserPointsLogic) UpdateUserPoints(in *user.UpdateUserPointsReq) (*user.BaseResp, error) {
	updates := map[string]interface{}{
		"green_points": gorm.Expr("green_points + ?", in.Points),
	}
	if in.Points > 0 {
		updates["green_points_total_earned"] = gorm.Expr("green_points_total_earned + ?", in.Points)
	}

	err := l.svcCtx.DB.Model(&model.SysUser{}).
		Where("id = ?", in.UserId).
		Updates(updates).Error
	if err != nil {
		l.Errorf("failed to update user green_points: %v", err)
		return &user.BaseResp{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return &user.BaseResp{
		Code:    0,
		Message: "success",
	}, nil
}
