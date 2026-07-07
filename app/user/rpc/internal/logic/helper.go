package logic

import (
	"smartcommunity-microservices/app/user/rpc/internal/model"
	"smartcommunity-microservices/app/user/rpc/user"
)

func mapUserInfo(u *model.SysUser, roleNames map[string]string) *user.UserInfo {
	if u == nil {
		return nil
	}
	roleName := ""
	if roleNames != nil {
		roleName = roleNames[u.Role]
	}
	if roleName == "" {
		mapFallback := map[string]string{
			"admin":    "系统管理员",
			"store":    "商户",
			"property": "物业人员",
			"user":     "居民",
		}
		if name, ok := mapFallback[u.Role]; ok {
			roleName = name
		} else {
			roleName = "居民"
		}
	}
	return &user.UserInfo{
		Id:             u.ID,
		Username:       u.Username,
		RealName:       u.RealName,
		Mobile:         u.Mobile,
		Avatar:         u.Avatar,
		GreenPoints:    int32(u.GreenPoints),
		Role:           u.Role,
		RoleName:       roleName,
		Status:         int32(u.Status),
		FaceRegistered: u.FaceRegistered,
		FaceImageUrl:   u.FaceImageURL,
		Balance:        u.Balance,
		Gender:         int32(u.Gender),
		Email:          u.Email,
		Age:            int32(u.Age),
	}
}
