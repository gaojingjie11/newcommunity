package logic

import (
	"smartcommunity-microservices/app/user/rpc/internal/model"
	"smartcommunity-microservices/app/user/rpc/user"
)

func mapUserInfo(u *model.SysUser) *user.UserInfo {
	if u == nil {
		return nil
	}
	return &user.UserInfo{
		Id:             u.ID,
		Username:       u.Username,
		RealName:       u.RealName,
		Mobile:         u.Mobile,
		Avatar:         u.Avatar,
		GreenPoints:    int32(u.GreenPoints),
		Role:           u.Role,
		Status:         int32(u.Status),
		FaceRegistered: u.FaceRegistered,
		FaceImageUrl:   u.FaceImageURL,
		Balance:        u.Balance,
	}
}
