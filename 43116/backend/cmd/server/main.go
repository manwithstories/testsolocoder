package main

import (
	"car-rental/internal/config"
	cachedb "car-rental/internal/config"
	"car-rental/internal/model"
	"car-rental/internal/router"
	"car-rental/internal/service"
	"car-rental/internal/utils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)

	err = cachedb.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer cachedb.CloseDB()

	err = cachedb.InitRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}
	defer cachedb.CloseRedis()

	utils.InitJWT(&cfg.JWT)

	if err := os.MkdirAll(cfg.Upload.Path+"/cars", 0755); err != nil {
		log.Printf("创建上传目录失败: %v", err)
	}
	if err := os.MkdirAll(cfg.Upload.Path+"/licenses", 0755); err != nil {
		log.Printf("创建上传目录失败: %v", err)
	}
	if err := os.MkdirAll(cfg.Upload.Path+"/avatars", 0755); err != nil {
		log.Printf("创建上传目录失败: %v", err)
	}
	if err := os.MkdirAll(cfg.Upload.Path+"/exports", 0755); err != nil {
		log.Printf("创建导出目录失败: %v", err)
	}

	initDefaultData()

	messageService := service.NewMessageService(&cfg.Email)
	go messageService.ProcessNotificationQueue()
	log.Println("消息队列消费者已启动")

	schedulerService := service.NewSchedulerService(messageService)
	schedulerService.Start()
	defer schedulerService.Stop()
	log.Println("定时任务调度器已启动")

	r := router.SetupRouter(cfg)

	log.Printf("服务器启动于端口 %s", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func initDefaultData() {
	db := cachedb.DB

	var adminRole model.Role
	if db.Where("name = ?", "admin").First(&adminRole).Error != nil {
		adminRole = model.Role{
			Name:        "admin",
			Description: "系统管理员",
		}
		db.Create(&adminRole)
	}

	var userRole model.Role
	if db.Where("name = ?", "user").First(&userRole).Error != nil {
		userRole = model.Role{
			Name:        "user",
			Description: "普通用户",
		}
		db.Create(&userRole)
	}

	var adminUser model.User
	if db.Where("username = ?", "admin").First(&adminUser).Error != nil {
		hashedPassword, _ := utils.HashPassword("admin123")
		adminUser = model.User{
			Username:   "admin",
			Password:   hashedPassword,
			Email:      "admin@example.com",
			RealName:   "系统管理员",
			RoleID:     adminRole.ID,
			AuthStatus: model.UserStatusActive,
			Status:     model.UserStatusActive,
		}
		db.Create(&adminUser)
	}

	log.Println("默认数据初始化完成")
}
