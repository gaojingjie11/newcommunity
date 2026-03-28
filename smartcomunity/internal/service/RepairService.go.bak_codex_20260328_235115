package service

import (
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
)

type RepairService struct{}

// Create 提交报修/投诉
func (s *RepairService) Create(repair *model.Repair) error {
	// 默认状态 0 (待处理)
	repair.Status = 0
	return global.DB.Create(repair).Error
}

// GetUserList 获取我的报修记录 (分页)
func (s *RepairService) GetUserList(userID int64, page, size int) ([]model.Repair, int64, error) {
	var list []model.Repair
	var total int64
	db := global.DB.Model(&model.Repair{}).Where("user_id = ?", userID)
	db.Count(&total)

	offset := (page - 1) * size
	err := db.Order("created_at desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

// UpdateStatus 更新报修状态 (派单/完成)
func (s *RepairService) UpdateStatus(id int64, status int, feedback string) error {
	// 状态: 0待处理 1处理中 2已完成
	return global.DB.Model(&model.Repair{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": status,
		"result": feedback, // 管理员的处理反馈
	}).Error
}

// GetAllList 管理员查看所有报修
func (s *RepairService) GetAllList(limit int) ([]model.Repair, error) {
	var list []model.Repair
	err := global.DB.Order("id desc").Limit(limit).Find(&list).Error
	return list, err
}

// GetPageList 分页获取所有报修
func (s *RepairService) GetPageList(page, size int) ([]model.Repair, int64, error) {
	var list []model.Repair
	var total int64
	offset := (page - 1) * size

	db := global.DB.Model(&model.Repair{})

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Preload("User").Order("id desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}
