package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"smartcommunity-microservices/common/db"
)

type conversationRow struct {
	ID        string    `gorm:"column:id"`
	UserID    int64     `gorm:"column:user_id"`
	Title     string    `gorm:"column:title"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func main() {
	conn, err := db.InitPostgres(db.PostgresConfig{
		Host:     "42.193.104.173",
		Port:     5432,
		Database: "smart_community",
		Username: "admin",
		Password: "sdl@admin",
		SSLMode:  "disable",
		TimeZone: "Asia/Shanghai",
	})
	if err != nil {
		log.Fatalf("connect postgres failed: %v", err)
	}

	var conversations []conversationRow
	if err := conn.Table("sys_user_conversation").Order("updated_at desc").Find(&conversations).Error; err != nil {
		log.Fatalf("load conversations failed: %v", err)
	}

	var invalid []conversationRow
	for _, conv := range conversations {
		if strings.TrimSpace(conv.ID) == "" {
			invalid = append(invalid, conv)
		}
	}

	fmt.Printf("total_conversations=%d invalid_conversations=%d\n", len(conversations), len(invalid))
	for _, conv := range invalid {
		fmt.Printf("invalid id=%q user_id=%d title=%q updated_at=%s\n",
			conv.ID, conv.UserID, conv.Title, conv.UpdatedAt.Format("2006-01-02 15:04:05"))
	}
}
