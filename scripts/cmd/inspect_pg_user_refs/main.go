package main

import (
	"fmt"
	"log"
	"sort"

	"smartcommunity-microservices/common/db"

	"gorm.io/gorm"
)

type tableInfo struct {
	TableName  string
	ColumnName string
}

type refSpec struct {
	Table string
	Col   string
	Where string
}

type userRow struct {
	ID       int64
	Username string
	Mobile   string
	Role     string
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
	if err := conn.Table("sys_user").Select("id, username, mobile, role").Order("id asc").Find(&users).Error; err != nil {
		log.Fatalf("load users failed: %v", err)
	}

	fmt.Println("== users ==")
	for _, u := range users {
		fmt.Printf("id=%d username=%s mobile=%s role=%s\n", u.ID, u.Username, u.Mobile, u.Role)
	}

	var refs []tableInfo
	if err := conn.Raw(`
SELECT table_name, column_name
FROM information_schema.columns
WHERE table_schema = 'public'
  AND column_name IN ('user_id', 'admin_user_id', 'processor_id', 'operator_id', 'generated_by')
ORDER BY table_name, column_name
`).Scan(&refs).Error; err != nil {
		log.Fatalf("scan information_schema failed: %v", err)
	}

	fmt.Println()
	fmt.Println("== direct user-id reference columns ==")
	for _, ref := range refs {
		var count int64
		query := fmt.Sprintf(`SELECT COUNT(*) FROM "%s" WHERE "%s" IS NOT NULL`, ref.TableName, ref.ColumnName)
		if err := conn.Raw(query).Scan(&count).Error; err != nil {
			log.Fatalf("count %s.%s failed: %v", ref.TableName, ref.ColumnName, err)
		}
		fmt.Printf("%s.%s count=%d\n", ref.TableName, ref.ColumnName, count)
	}

	extra := map[string]string{
		"wallet_transactions.related_id": "SELECT COUNT(*) FROM wallet_transactions WHERE type = 2 OR biz_type IN ('admin_recharge', 'admin_deduct')",
	}

	keys := make([]string, 0, len(extra))
	for k := range extra {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println()
	fmt.Println("== indirect / conditional references ==")
	for _, key := range keys {
		var count int64
		if err := conn.Raw(extra[key]).Scan(&count).Error; err != nil {
			log.Fatalf("count %s failed: %v", key, err)
		}
		fmt.Printf("%s count=%d\n", key, count)
	}

	fmt.Println()
	fmt.Println("== unique indexes on sys_user ==")
	type idxRow struct {
		IndexName string
		IndexDef  string
	}
	var idxRows []idxRow
	if err := conn.Raw(`
SELECT indexname, indexdef
FROM pg_indexes
WHERE schemaname = 'public' AND tablename = 'sys_user'
ORDER BY indexname
`).Scan(&idxRows).Error; err != nil {
		log.Fatalf("scan indexes failed: %v", err)
	}
	for _, row := range idxRows {
		fmt.Printf("%s => %s\n", row.IndexName, row.IndexDef)
	}

	fmt.Println()
	fmt.Println("== foreign keys referencing sys_user ==")
	type fkRow struct {
		ConstraintName string
		TableName      string
		ColumnName     string
	}
	var fkRows []fkRow
	if err := conn.Raw(`
SELECT
  tc.constraint_name,
  tc.table_name,
  kcu.column_name
FROM information_schema.table_constraints tc
JOIN information_schema.key_column_usage kcu
  ON tc.constraint_name = kcu.constraint_name
 AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage ccu
  ON ccu.constraint_name = tc.constraint_name
 AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY'
  AND tc.table_schema = 'public'
  AND ccu.table_name = 'sys_user'
  AND ccu.column_name = 'id'
ORDER BY tc.table_name, kcu.column_name
`).Scan(&fkRows).Error; err != nil {
		log.Fatalf("scan foreign keys failed: %v", err)
	}
	if len(fkRows) == 0 {
		fmt.Println("(none)")
	} else {
		for _, row := range fkRows {
			fmt.Printf("%s => %s.%s\n", row.ConstraintName, row.TableName, row.ColumnName)
		}
	}

	specs := []refSpec{
		{Table: "sys_user_role", Col: "user_id"},
		{Table: "user_login_logs", Col: "user_id", Where: `"user_id" > 0`},
		{Table: "admin_login_logs", Col: "admin_user_id", Where: `"admin_user_id" > 0`},
		{Table: "oms_cart", Col: "user_id"},
		{Table: "oms_order", Col: "user_id"},
		{Table: "pms_favorite", Col: "user_id"},
		{Table: "pms_product_comment", Col: "user_id"},
		{Table: "wallets", Col: "user_id"},
		{Table: "wallet_transactions", Col: "user_id"},
		{Table: "wallet_transactions", Col: "related_id", Where: `("type" = 2 OR "biz_type" IN ('admin_recharge', 'admin_deduct')) AND "related_id" > 0`},
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

	fmt.Println()
	fmt.Println("== orphan user references ==")
	for _, spec := range specs {
		exists, err := columnExists(conn, spec.Table, spec.Col)
		if err != nil {
			log.Fatalf("column check %s.%s failed: %v", spec.Table, spec.Col, err)
		}
		if !exists {
			fmt.Printf("%s.%s orphan_check_skipped=missing\n", spec.Table, spec.Col)
			continue
		}
		query := fmt.Sprintf(`SELECT COUNT(*) FROM "%s" t LEFT JOIN sys_user u ON t."%s" = u.id WHERE u.id IS NULL`, spec.Table, spec.Col)
		if spec.Where != "" {
			query += " AND (" + spec.Where + ")"
		}
		var count int64
		if err := conn.Raw(query).Scan(&count).Error; err != nil {
			log.Fatalf("orphan check %s.%s failed: %v", spec.Table, spec.Col, err)
		}
		fmt.Printf("%s.%s orphan_count=%d\n", spec.Table, spec.Col, count)
	}

	_ = conn
	_ = (*gorm.DB)(nil)
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
