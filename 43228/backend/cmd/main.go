package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"tea-platform/config"
	"tea-platform/internal/router"
)

var (
	DB  *gorm.DB
	JWTSecret []byte
)

func initConfig() {
	path := "config/config.yaml"
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		path = p
	}
	cfg, err := config.Load(path)
	if err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}
	JWTSecret = []byte(cfg.JWT.Secret)
	log.Printf("[main] 配置加载成功, 监听端口: %s", cfg.Server.Port)
}

func initLogger() {
	cfg := config.Get()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("[main] 日志级别: %s, 格式: %s, 输出: %s", cfg.Log.Level, cfg.Log.Format, cfg.Log.Output)
}

func initDB() {
	cfg := config.Get()
	gormLevel := gormlogger.Warn
	switch cfg.Log.Level {
	case "debug":
		gormLevel = gormlogger.Info
	case "info":
		gormLevel = gormlogger.Warn
	case "error":
		gormLevel = gormlogger.Error
	}

	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var err error
	DB, err = gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取底层 sql.DB 失败: %v", err)
	}
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)
	log.Printf("[main] 数据库连接成功: %s@%s:%s/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
}

func initJWT() {
	_ = jwt.SigningMethodHS256
	log.Printf("[main] JWT 已初始化, 过期时间: %d 小时", config.Get().JWT.ExpireHours)
}

func setupRouter() *gin.Engine {
	cfg := config.Get()
	gin.SetMode(cfg.Server.Mode)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "tea-platform",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	return r
}

func main() {
	initConfig()
	initLogger()
	initDB()
	initJWT()

	r := setupRouter()
	router.SetupRoutes(r)
	cfg := config.Get()

	addr := cfg.Server.Port
	if addr == "" {
		addr = ":8080"
	}

	fmt.Printf(`
╔══════════════════════════════════════╗
║       Tea Platform Server            ║
║       Listening on %s           ║
╚══════════════════════════════════════╝
`, addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
