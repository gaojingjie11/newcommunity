package handler

import (
	"smartcommunity-microservices/services/mall-service/internal/repository"
)

// PermissionProvider implements middleware.PermissionProvider for mall-service.
type PermissionProvider struct {
	repo *repository.PermissionRepo
}

func NewPermissionProvider(repo *repository.PermissionRepo) *PermissionProvider {
	return &PermissionProvider{repo: repo}
}

func (p *PermissionProvider) GetPermissionCodesByUserID(userID int64) ([]string, error) {
	return p.repo.GetPermissionCodesByUserID(userID)
}
