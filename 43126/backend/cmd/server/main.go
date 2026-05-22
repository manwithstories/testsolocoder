package main

import (
	"fmt"
	"log"

	"meeting-room/internal/config"
	"meeting-room/internal/router"
	"meeting-room/internal/services"
	"meeting-room/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	config.Init()
	utils.LogDir()
	utils.InitLogger()
	utils.InitDB()
	utils.InitRedis()

	gin.SetMode(config.Cfg.Server.Mode)

	r := router.SetupRouter()

	startBackgroundJobs()

	addr := fmt.Sprintf(":%s", config.Cfg.Server.Port)
	log.Printf("服务器启动中，监听端口: %s", config.Cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func startBackgroundJobs() {
	notificationService := services.NewNotificationService()
	materialService := services.NewMaterialService()

	go notificationService.ProcessNotificationQueue()

	c := cron.New()

	c.AddFunc("0 * * * *", func() {
		notificationService.CheckAndSendReminders()
	})

	c.AddFunc("0 2 * * *", func() {
		materialService.CleanupExpiredMaterials()
	})

	c.AddFunc("0 */5 * * * *", func() {
		materialService.CleanupExpiredFiles()
	})

	c.Start()

	log.Println("后台任务已启动")
}
