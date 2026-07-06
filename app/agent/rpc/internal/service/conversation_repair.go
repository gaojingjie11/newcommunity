package service

import (
	"fmt"
	"strings"

	"smartcommunity-microservices/app/agent/rpc/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RepairLegacyConversationIDs migrates malformed historical conversation IDs
// created by earlier buggy builds, then rewrites dependent chat/approval rows.
func RepairLegacyConversationIDs(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	var badConversations []model.SysUserConversation
	if err := db.Where("COALESCE(BTRIM(id), '') = ''").Order("updated_at DESC").Find(&badConversations).Error; err != nil {
		return fmt.Errorf("query invalid conversations failed: %w", err)
	}

	for _, conv := range badConversations {
		oldID := conv.ID
		newID := uuid.NewString()

		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(
				`UPDATE sys_user_chat_message
				 SET conversation_id = ?
				 WHERE conversation_id = ? AND user_id = ?`,
				newID, oldID, conv.UserID,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`UPDATE agent_action_approval
				 SET conversation_id = ?
				 WHERE conversation_id = ? AND user_id = ?`,
				newID, oldID, conv.UserID,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`UPDATE sys_user_conversation
				 SET id = ?
				 WHERE id = ? AND user_id = ?`,
				newID, oldID, conv.UserID,
			).Error; err != nil {
				return err
			}

			return nil
		}); err != nil {
			return fmt.Errorf("repair invalid conversation id for user %d failed: %w", conv.UserID, err)
		}
	}

	// Extra guard for old rows that may have stray whitespace IDs.
	var whitespaceConversations []model.SysUserConversation
	if err := db.Where("id <> '' AND id <> BTRIM(id)").Find(&whitespaceConversations).Error; err != nil {
		return fmt.Errorf("query whitespace conversation ids failed: %w", err)
	}

	for _, conv := range whitespaceConversations {
		trimmedID := strings.TrimSpace(conv.ID)
		if trimmedID == "" || trimmedID == conv.ID {
			continue
		}

		if err := db.Transaction(func(tx *gorm.DB) error {
			var existing int64
			if err := tx.Model(&model.SysUserConversation{}).
				Where("id = ? AND user_id = ?", trimmedID, conv.UserID).
				Count(&existing).Error; err != nil {
				return err
			}
			if existing > 0 {
				trimmedID = uuid.NewString()
			}

			if err := tx.Exec(
				`UPDATE sys_user_chat_message
				 SET conversation_id = ?
				 WHERE conversation_id = ? AND user_id = ?`,
				trimmedID, conv.ID, conv.UserID,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`UPDATE agent_action_approval
				 SET conversation_id = ?
				 WHERE conversation_id = ? AND user_id = ?`,
				trimmedID, conv.ID, conv.UserID,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				`UPDATE sys_user_conversation
				 SET id = ?
				 WHERE id = ? AND user_id = ?`,
				trimmedID, conv.ID, conv.UserID,
			).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			return fmt.Errorf("repair whitespace conversation id for user %d failed: %w", conv.UserID, err)
		}
	}

	return nil
}
