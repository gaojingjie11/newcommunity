package logic

import (
	"context"
	"errors"

	"smartcommunity-microservices/app/user/rpc/internal/model"
	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"
	"smartcommunity-microservices/common/auth"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type VerifyPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyPasswordLogic {
	return &VerifyPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Password Verification (no side effects)
func (l *VerifyPasswordLogic) VerifyPassword(in *user.VerifyPasswordReq) (*user.VerifyPasswordResp, error) {
	var u model.SysUser
	err := l.svcCtx.DB.Where("id = ?", in.UserId).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.VerifyPasswordResp{Success: false}, nil
		}
		return nil, err
	}

	success := auth.CheckPasswordHash(in.Password, u.Password)
	return &user.VerifyPasswordResp{Success: success}, nil
}
