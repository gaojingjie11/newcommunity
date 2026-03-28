package main

import (
	"fmt"
	"log"
	"os"
	"smartcommunity/internal/config"
	"smartcommunity/internal/global"
	"smartcommunity/internal/middleware"
	"smartcommunity/internal/model"
	"smartcommunity/internal/router"
	"smartcommunity/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	config.Init(env)

	global.InitDB(config.Conf.DB.DSN)
	global.InitRedis(config.Conf.Redis.Addr, "")
	global.InitMinio(config.Conf.MinIO)

	if err := global.DB.AutoMigrate(
		&model.SysUser{},
		&model.SysRole{},
		&model.SysMenu{},
		&model.SysUserRole{},
		&model.SysRoleMenu{},
		&model.Product{},
		&model.ProductCategory{},
		&model.HotProduct{},
		&model.Promotion{},
		&model.Store{},
		&model.StoreProduct{},
		&model.Cart{},
		&model.Order{},
		&model.OrderItem{},
		&model.Notice{},
		&model.NoticeRead{},
		&model.Repair{},
		&model.Visitor{},
		&model.Parking{},
		&model.PropertyFee{},
		&model.Favorite{},
		&model.ProductComment{},
		&model.SysTransaction{},
		&model.GreenPointRecord{},
		&model.AIReport{},
		&model.ChatMessage{},
		&model.CommunityMessage{},
	); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	service.StartAIReportDailyScheduler()

	r := gin.Default()
	r.Use(middleware.CORS())

	router.InitRouter(r)

	port := config.Conf.Server.Port
	if port == "" {
		port = "8080"
	}
	fmt.Printf("service started on :%s\n", port)
	r.Run(":" + port)
}
