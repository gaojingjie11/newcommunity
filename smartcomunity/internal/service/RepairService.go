package service

import (
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"strings"
)

type RepairService struct{}

// Create submits a repair/complaint ticket.
func (s *RepairService) Create(repair *model.Repair) error {
	repair.Status = 0
	repair.Category = normalizeRepairCategoryForDisplay(repair.Category)
	return global.DB.Create(repair).Error
}

// GetUserList returns current user's tickets.
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

// UpdateStatus updates ticket status.
func (s *RepairService) UpdateStatus(id int64, status int, feedback string) error {
	return global.DB.Model(&model.Repair{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": status,
		"result": feedback,
	}).Error
}

// GetAllList returns latest tickets for admin.
func (s *RepairService) GetAllList(limit int) ([]model.Repair, error) {
	var list []model.Repair
	err := global.DB.Order("id desc").Limit(limit).Find(&list).Error
	applyRepairCategoryLabels(list)
	return list, err
}

// GetPageList returns paged tickets for admin.
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
		cat := normalizeRepairCategoryForDisplay(list[i].Category)
		if cat == "" || cat == "\u5176\u4ed6" {
			if inferred := inferRepairCategoryFromContent(list[i].Content); inferred != "" {
				cat = inferred
			}
		}
		if cat == "" {
			cat = "\u5176\u4ed6"
		}
		list[i].Category = cat
	}
}

func normalizeRepairCategoryForDisplay(raw string) string {
	v := strings.TrimSpace(raw)
	if v == "" {
		return ""
	}
	switch strings.ToLower(v) {
	case "plumbing", "water", "pipe", "\u6f0f\u6c34", "\u4e0b\u6c34", "\u6c34\u7ba1", "\u7ba1\u9053", "\u6c34\u6696":
		return "\u6c34\u6696"
	case "door_window", "door", "window", "lock", "\u95e8\u7a97":
		return "\u95e8\u7a97"
	case "electrical", "electric", "power", "\u7535\u8def", "\u7528\u7535":
		return "\u7535\u8def"
	case "air_conditioner", "aircon", "ac", "\u7a7a\u8c03":
		return "\u7a7a\u8c03"
	case "heating", "radiator", "\u6696\u6c14", "\u4f9b\u6696":
		return "\u4f9b\u6696"
	case "noise", "\u6270\u6c11", "\u566a\u97f3":
		return "\u566a\u97f3"
	case "sanitation", "clean", "\u536b\u751f", "\u5783\u573e":
		return "\u536b\u751f"
	case "other", "\u5176\u4ed6":
		return "\u5176\u4ed6"
	default:
		return v
	}
}

func inferRepairCategoryFromContent(content string) string {
	text := strings.ToLower(strings.TrimSpace(content))
	if text == "" {
		return ""
	}
	switch {
	case containsRepairKeyword(text, "\u6c34\u7ba1", "\u7ba1\u9053", "\u6c34\u9f99\u5934", "\u6f0f\u6c34", "\u4e0b\u6c34", "pipe", "plumb", "faucet"):
		return "\u6c34\u6696"
	case containsRepairKeyword(text, "\u95e8", "\u7a97", "\u95e8\u7a97", "door", "window", "lock"):
		return "\u95e8\u7a97"
	case containsRepairKeyword(text, "\u7535\u8def", "\u8df3\u95f8", "\u63d2\u5ea7", "\u706f", "power", "electric"):
		return "\u7535\u8def"
	case containsRepairKeyword(text, "\u6696\u6c14", "\u4f9b\u6696", "\u4e0d\u70ed", "heating", "radiator"):
		return "\u4f9b\u6696"
	case containsRepairKeyword(text, "\u7a7a\u8c03", "aircon", "ac"):
		return "\u7a7a\u8c03"
	case containsRepairKeyword(text, "\u566a\u97f3", "\u6270\u6c11", "noise"):
		return "\u566a\u97f3"
	case containsRepairKeyword(text, "\u536b\u751f", "\u5783\u573e", "clean"):
		return "\u536b\u751f"
	default:
		return ""
	}
}

func containsRepairKeyword(text string, keywords ...string) bool {
	for _, kw := range keywords {
		if kw != "" && strings.Contains(text, kw) {
			return true
		}
	}
	return false
}
