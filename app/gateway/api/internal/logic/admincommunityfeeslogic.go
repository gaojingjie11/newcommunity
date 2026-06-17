package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/user/rpc/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCommunityFeesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCommunityFeesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCommunityFeesLogic {
	return &AdminCommunityFeesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCommunityFeesLogic) AdminCommunityFees(req *types.ListPropertyFeesReq) (resp *types.PropertyFeeListResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.AdminListFees(l.ctx, &communityrpc.ListPropertyFeesReq{
		UserId: 0,
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		return nil, err
	}

	type userInfo struct {
		name   string
		mobile string
	}
	userCache := make(map[int64]userInfo)

	list := make([]types.PropertyFeeInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		apiInfo := toAPIPropertyFeeInfo(item)

		if uInfo, ok := userCache[item.UserId]; ok {
			apiInfo.UserName = uInfo.name
			apiInfo.UserMobile = uInfo.mobile
		} else {
			profile, err := l.svcCtx.UserRpc.GetProfile(l.ctx, &userrpc.UserIDReq{UserId: item.UserId})
			if err == nil {
				userCache[item.UserId] = userInfo{name: profile.RealName, mobile: profile.Mobile}
				apiInfo.UserName = profile.RealName
				apiInfo.UserMobile = profile.Mobile
			} else {
				userCache[item.UserId] = userInfo{name: "未知用户", mobile: ""}
				apiInfo.UserName = "未知用户"
				apiInfo.UserMobile = ""
			}
		}

		list = append(list, apiInfo)
	}

	return &types.PropertyFeeListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
