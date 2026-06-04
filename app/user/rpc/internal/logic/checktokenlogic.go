package logic

import (
	"context"
	"fmt"
	"strings"

	"smartcommunity-microservices/app/user/rpc/internal/svc"
	"smartcommunity-microservices/app/user/rpc/user"
	"smartcommunity-microservices/common/auth"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CheckTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckTokenLogic {
	return &CheckTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckTokenLogic) CheckToken(in *user.CheckTokenReq) (*user.CheckTokenResp, error) {
	tokenStr := in.Token
	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = tokenStr[7:]
	}

	claims, err := auth.ParseToken(l.svcCtx.Config.Jwt.Secret, tokenStr)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Token无效或已过期")
	}

	redisKey := fmt.Sprintf("login:token:%d", claims.UserID)
	cachedToken, err := l.svcCtx.RedisClient.Get(l.ctx, redisKey).Result()
	if err != nil || cachedToken != tokenStr {
		return nil, status.Error(codes.Unauthenticated, "登录已失效，请重新登录")
	}

	return &user.CheckTokenResp{
		UserId: claims.UserID,
		Role:   claims.Role,
	}, nil
}
