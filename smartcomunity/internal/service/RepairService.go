package service

import (
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"strings"
)

type RepairService struct{}

// Create 提交报修/投诉
func (s *RepairService) Create(repair *model.Repair) error {
	repair.Status = 0
	repair.Category = normalizeRepairCategoryForDisplay(repair.Category)
	return global.DB.Create(repair).Error
}

// GetUserList 获取我的报修记录（分页）
func (s *RepairService) GetUserList(userID int64, page, size int) ([]model.Repair, int64, error) {
	var list []model.Repair
	var total int64
	db := global.DB.Model(&model.Repair{}).Where("user_id = ?", userID)
	db.Count(&total)

	offset := (page - 1) * size
	err := db.Order("created_at desc").Offset(offset).Limit(size).Find(&list).Error
	applyRepairCategoryLabels(list)
	return list, total, err
}

// UpdateStatus 更新报修状态（派单/完成）
func (s *RepairService) UpdateStatus(id int64, status int, feedback string) error {
	return global.DB.Model(&model.Repair{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": status,
		"result": feedback,
	}).Error
}

// GetAllList 管理员查看所有报修
func (s *RepairService) GetAllList(limit int) ([]model.Repair, error) {
	var list []model.Repair
	err := global.DB.Order("id desc").Limit(limit).Find(&list).Error
	applyRepairCategoryLabels(list)
	return list, err
}

// GetPageList 分页获取所有报修
func (s *RepairService) GetPageList(page, size int) ([]model.Repair, int64, error) {
	var list []model.Repair
	var total int64
	offset := (page - 1) * size

	db := global.DB.Model(&model.Repair{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Preload("User").Order("id desc").Offset(offset).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	applyRepairCategoryLabels(list)
	return list, total, nil
}

func applyRepairCategoryLabels(list []model.Repair) {
	for i := range list {
		list[i].Category = normalizeRepairCategoryForDisplay(list[i].Category)
	}
}

func normalizeRepairCategoryForDisplay(raw string) string {
	v := strings.TrimSpace(raw)
	if v == "" {
		return ""
	}
	switch strings.ToLower(v) {
	case "plumbing", "water", "漏水", "下水", "水暖":
		return "水暖"
	case "door_window", "door", "window", "lock", "门窗":
		return "门窗"
	case "electrical", "electric", "power", "电路", "用电":
		return "电路"
	case "air_conditioner", "aircon", "ac", "空调":
		return "空调"
	case "heating", "radiator", "暖气", "供暖":
		return "供暖"
	case "noise", "扰民", "噪音":
		return "噪音"
	case "sanitation", "clean", "卫生", "垃圾":
		return "卫生"
	case "other", "其他":
		return "其他"
	default:
		return v
	}
}
