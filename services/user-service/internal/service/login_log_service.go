package service

import (
	"smartcommunity-microservices/services/user-service/internal/model"
	"smartcommunity-microservices/services/user-service/internal/repository"
)

type LoginLogService struct {
	logRepo *repository.LoginLogRepo
}

func NewLoginLogService(logRepo *repository.LoginLogRepo) *LoginLogService {
	return &LoginLogService{logRepo: logRepo}
}

// LOG-001: Query user login logs
func (s *LoginLogService) QueryUserLogs(page, size int, mobile string, success *bool) ([]model.UserLoginLog, int64, error) {
	return s.logRepo.QueryUserLogs(page, size, mobile, success)
}

// LOG-002: Query admin login logs
func (s *LoginLogService) QueryAdminLogs(page, size int, mobile string, success *bool) ([]model.AdminLoginLog, int64, error) {
	return s.logRepo.QueryAdminLogs(page, size, mobile, success)
}
