package main

import (
	"fmt"
	"log"
	"translation-platform/internal/config"
	"translation-platform/internal/database"
	"translation-platform/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)

	if err := database.Init(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	routes.SetupRoutes(r)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("服务器启动于 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
