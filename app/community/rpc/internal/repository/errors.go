package repository

import "errors"

var (
	ErrParkingSpaceUnavailable = errors.New("parking space unavailable")
	ErrDuplicateParkingNo      = errors.New("duplicate parking number")
	ErrParkingUserNotFound     = errors.New("parking user not found")
	ErrInvalidCarPlate         = errors.New("车牌号格式不正确，请输入标准蓝牌格式，如：辽A12345")
	ErrCarPlateAlreadyExists   = errors.New("该车牌号已被绑定，请检查后重试")
	ErrPropertyFeePaid         = errors.New("property fee already paid")
)
