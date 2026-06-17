package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/user/rpc/userrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommunityMyFeesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommunityMyFeesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommunityMyFeesLogic {
	return &CommunityMyFeesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommunityMyFeesLogic) CommunityMyFees(req *types.ListPropertyFeesReq) (resp *types.PropertyFeeListResp, err error) {
	userID := getUserIDFromCtx(l.ctx)
	rpcResp, err := l.svcCtx.CommunityRpc.ListMyFees(l.ctx, &communityrpc.ListPropertyFeesReq{
		UserId: userID,
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		return nil, err
	}

	var name, mobile string
	profile, err := l.svcCtx.UserRpc.GetProfile(l.ctx, &userrpc.UserIDReq{UserId: userID})
	if err == nil {
		name = profile.RealName
		mobile = profile.Mobile
	} else {
		name = "未知用户"
	}

	list := make([]types.PropertyFeeInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		apiInfo := toAPIPropertyFeeInfo(item)
		apiInfo.UserName = name
		apiInfo.UserMobile = mobile
		list = append(list, apiInfo)
	}
	return &types.PropertyFeeListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
