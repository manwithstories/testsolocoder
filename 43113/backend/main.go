package main

import (
	"fmt"
	"os"
	"qa-platform/config"
	"qa-platform/models"
	"qa-platform/repository"
	"qa-platform/router"
	"qa-platform/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	cfg := config.LoadConfig(env)

	utils.InitJWT(cfg.JWT.Secret)
	utils.InitLogger(cfg.Log.Level, cfg.Log.Output)

	gin.SetMode(cfg.Server.Mode)

	if err := repository.InitDB(cfg.Database); err != nil {
		fmt.Printf("数据库初始化失败: %v\n", err)
		os.Exit(1)
	}

	if err := repository.SeedData(); err != nil {
		utils.LogError("数据初始化失败: %v", err)
	}

	var sensitiveWords []models.SensitiveWord
	repository.DB.Find(&sensitiveWords)
	utils.InitSensitiveFilter(sensitiveWords)
	utils.LogInfo("敏感词过滤器初始化完成，共加载 %d 个敏感词", len(sensitiveWords))

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			questionService := &struct{}{}
			_ = questionService
			var questions []models.Question
			repository.DB.Where("status = ?", "published").Find(&questions)
			now := time.Now()
			for _, q := range questions {
				hours := now.Sub(q.CreatedAt).Hours()
				decay := 1.0
				if hours > 24 {
					decay = 1.0 / (hours / 24)
				}
				hotScore := float64(q.Views)*0.1 + float64(q.LikeCount)*2 + float64(q.AnswerCount)*5 + float64(q.CollectCount)*3
				hotScore *= decay
				repository.DB.Model(&q).Update("hot_score", hotScore)
			}
		}
	}()

	r := router.SetupRouter()

	utils.LogInfo("服务器启动在端口 %s", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
		os.Exit(1)
	}
}
