package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"smartcommunity-microservices/common/db"

	"gorm.io/gorm"
)

type userRow struct {
	ID int64
}

type idMap struct {
	OldID  int64
	NewID  int64
	TempID int64
}

type refSpec struct {
	Table string
	Col   string
	Where string
}

func main() {
	cfg := db.PostgresConfig{
		Host:     "42.193.104.173",
		Port:     5432,
		Database: "smart_community",
		Username: "admin",
		Password: "sdl@admin",
		SSLMode:  "disable",
		TimeZone: "Asia/Shanghai",
	}

	conn, err := db.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("connect postgres failed: %v", err)
	}

	var users []userRow
	if err := conn.Table("sys_user").Select("id").Order("id asc").Find(&users).Error; err != nil {
		log.Fatalf("load users failed: %v", err)
	}
	if len(users) == 0 {
		log.Fatal("sys_user is empty")
	}

	maps := buildMapping(users)
	fmt.Println("planned mapping:")
	for _, m := range maps {
		fmt.Printf("%d -> %d\n", m.OldID, m.NewID)
	}

	if isAlreadySequential(maps) {
		fmt.Println("user ids are already sequential, nothing to do")
		return
	}

	specs := []refSpec{
		{Table: "sys_user", Col: "id"},
		{Table: "sys_user_role", Col: "user_id"},
		{Table: "user_login_logs", Col: "user_id"},
		{Table: "admin_login_logs", Col: "admin_user_id"},
		{Table: "oms_cart", Col: "user_id"},
		{Table: "oms_order", Col: "user_id"},
		{Table: "pms_favorite", Col: "user_id"},
		{Table: "pms_product_comment", Col: "user_id"},
		{Table: "wallets", Col: "user_id"},
		{Table: "wallet_transactions", Col: "user_id"},
		{Table: "wallet_transactions", Col: "related_id", Where: `"type" = 2 OR "biz_type" IN ('admin_recharge', 'admin_deduct')`},
		{Table: "payment_records", Col: "user_id"},
		{Table: "product_view_logs", Col: "user_id", Where: `"user_id" > 0`},
		{Table: "notice_view_logs", Col: "user_id"},
		{Table: "visitors", Col: "user_id"},
		{Table: "user_parking_bindings", Col: "user_id"},
		{Table: "property_fees", Col: "user_id"},
		{Table: "property_fee_payments", Col: "user_id"},
		{Table: "workorders", Col: "user_id"},
		{Table: "workorders", Col: "processor_id", Where: `"processor_id" > 0`},
		{Table: "workorder_logs", Col: "operator_id", Where: `"operator_id" > 0`},
		{Table: "community_messages", Col: "user_id"},
		{Table: "cms_ai_report", Col: "generated_by", Where: `"generated_by" > 0`},
		{Table: "pms_user_store", Col: "user_id"},
		{Table: "sys_user_conversation", Col: "user_id"},
		{Table: "sys_user_chat_message", Col: "user_id"},
		{Table: "agent_action_approval", Col: "user_id"},
	}

	if err := conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
CREATE TABLE IF NOT EXISTS sys_user_id_resequence_history (
  id BIGSERIAL PRIMARY KEY,
  executed_at TIMESTAMP NOT NULL DEFAULT NOW(),
  old_id BIGINT NOT NULL,
  new_id BIGINT NOT NULL
)`).Error; err != nil {
			return err
		}

		if err := tx.Exec(`CREATE TEMP TABLE tmp_user_id_map (old_id BIGINT PRIMARY KEY, new_id BIGINT NOT NULL UNIQUE, temp_id BIGINT NOT NULL UNIQUE) ON COMMIT DROP`).Error; err != nil {
			return err
		}

		for _, m := range maps {
			if err := tx.Exec(`INSERT INTO tmp_user_id_map (old_id, new_id, temp_id) VALUES (?, ?, ?)`, m.OldID, m.NewID, m.TempID).Error; err != nil {
				return err
			}
			if err := tx.Exec(`INSERT INTO sys_user_id_resequence_history (executed_at, old_id, new_id) VALUES (?, ?, ?)`, time.Now(), m.OldID, m.NewID).Error; err != nil {
				return err
			}
		}

		for _, spec := range specs {
			if err := remapIfExists(tx, spec, "old_id", "temp_id"); err != nil {
				return err
			}
		}
		for _, spec := range specs {
			if err := remapIfExists(tx, spec, "temp_id", "new_id"); err != nil {
				return err
			}
		}

		var seqName *string
		if err := tx.Raw(`SELECT pg_get_serial_sequence('sys_user', 'id')`).Scan(&seqName).Error; err != nil {
			return err
		}
		if seqName != nil && *seqName != "" {
			var maxID int64
			if err := tx.Raw(`SELECT COALESCE(MAX(id), 0) FROM sys_user`).Scan(&maxID).Error; err != nil {
				return err
			}
			sql := fmt.Sprintf(`SELECT setval('%s', %d, true)`, strings.ReplaceAll(*seqName, `'`, `''`), maxID)
			if err := tx.Exec(sql).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		log.Fatalf("resequence failed: %v", err)
	}

	fmt.Println("resequence committed successfully")
}

func buildMapping(users []userRow) []idMap {
	result := make([]idMap, 0, len(users))
	nextID := int64(0)
	hasOne := false
	for _, u := range users {
		if u.ID == 1 {
			hasOne = true
			break
		}
	}
	if hasOne {
		nextID = 1
	}

	for _, u := range users {
		if u.ID == 1 {
			result = append(result, idMap{OldID: u.ID, NewID: 1, TempID: 1000001})
			continue
		}
		nextID++
		result = append(result, idMap{OldID: u.ID, NewID: nextID, TempID: 1000000 + nextID})
	}
	return result
}

func isAlreadySequential(maps []idMap) bool {
	for _, m := range maps {
		if m.OldID != m.NewID {
			return false
		}
	}
	return true
}

func remapIfExists(tx *gorm.DB, spec refSpec, fromCol, toCol string) error {
	exists, err := columnExists(tx, spec.Table, spec.Col)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	sql := fmt.Sprintf(`UPDATE "%s" t SET "%s" = m.%s FROM tmp_user_id_map m WHERE t."%s" = m.%s`,
		spec.Table, spec.Col, toCol, spec.Col, fromCol)
	if spec.Where != "" {
		sql += " AND (" + spec.Where + ")"
	}
	return tx.Exec(sql).Error
}

func columnExists(tx *gorm.DB, table, column string) (bool, error) {
	var count int64
	err := tx.Raw(`
SELECT COUNT(*)
FROM information_schema.columns
WHERE table_schema = 'public'
  AND table_name = ?
  AND column_name = ?
`, table, column).Scan(&count).Error
	return count > 0, err
}
