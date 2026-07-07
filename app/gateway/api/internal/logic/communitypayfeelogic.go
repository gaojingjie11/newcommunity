package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/user/rpc/userrpc"
	"smartcommunity-microservices/common/faceauth"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type CommunityPayFeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommunityPayFeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommunityPayFeeLogic {
	return &CommunityPayFeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommunityPayFeeLogic) CommunityPayFee(req *types.PayPropertyFeeReq) (resp *types.BaseResp, err error) {
	userID := getUserIDFromCtx(l.ctx)

	// Defer cleanup of the temporary pay face image if one was uploaded
	if req.PayType == "face" && req.FaceImageUrl != "" && l.svcCtx.MinioClient != nil {
		defer func() {
			bucketName := l.svcCtx.Config.MinIO.Bucket
			url := req.FaceImageUrl
			searchStr := "/" + bucketName + "/"
			idx := strings.Index(url, searchStr)
			var objectKey string
			if idx != -1 {
				objectKey = url[idx+len(searchStr):]
			} else {
				parts := strings.Split(url, "/")
				if len(parts) >= 5 {
					objectKey = strings.Join(parts[4:], "/")
				}
			}

			if objectKey != "" {
				l.Logger.Infof("[Property Fee Pay] Deleting temp face image: bucket=%s, key=%s", bucketName, objectKey)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = l.svcCtx.MinioClient.RemoveObject(ctx, bucketName, objectKey, minio.RemoveObjectOptions{})
			}
		}()
	}

	// 1. Payment Verification
	if req.PayType == "password" {
		if req.Password == "" {
			return nil, fmt.Errorf("请输入支付密码")
		}
		// Verify password using UserRpc.VerifyPassword (no session side-effects)
		verifyResp, err := l.svcCtx.UserRpc.VerifyPassword(l.ctx, &userrpc.VerifyPasswordReq{
			UserId:   userID,
			Password: req.Password,
		})
		if err != nil {
			return nil, err
		}
		if !verifyResp.Success {
			return nil, fmt.Errorf("支付密码错误")
		}
	} else if req.PayType == "face" {
		profile, err := l.svcCtx.UserRpc.GetProfile(l.ctx, &userrpc.UserIDReq{UserId: userID})
		if err != nil {
			return nil, err
		}
		if !profile.FaceRegistered {
			return nil, fmt.Errorf("当前账号未录入人脸")
		}
		if req.FaceImageUrl == "" {
			return nil, fmt.Errorf("请先完成刷脸验证")
		}
		if _, err := faceauth.VerifyMatch(l.ctx, profile.FaceImageUrl, req.FaceImageUrl); err != nil {
			return nil, err
		}
	}

	// 2. Execute payment via CommunityRpc
	rpcResp, err := l.svcCtx.CommunityRpc.PayPropertyFee(l.ctx, &communityrpc.PayPropertyFeeReq{
		Id:     req.Id,
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}
