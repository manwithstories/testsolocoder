package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"beehive-platform/config"
	"beehive-platform/database"
	"beehive-platform/models"
	"beehive-platform/routes"
	"beehive-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	utils.InitJWT(cfg.JWT.Secret)

	logger, err := utils.InitLogger(utils.LoggerConfig{
		Path:       cfg.Log.Path,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	log.SetOutput(logger)

	if err := os.MkdirAll(cfg.Upload.Path, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	if err := database.Init(cfg.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db := database.GetDB()
	runMigrations(db)

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = logger

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithWriter(logger))
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Disposition")
		c.Header("Access-Control-Max-Age", "86400")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	routes.SetupRoutes(r, cfg)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runMigrations(db *gorm.DB) {
	tables := []interface{}{
		&models.User{},
		&models.Beehive{},
		&models.HealthRecord{},
		&models.Harvest{},
		&models.Inventory{},
		&models.Inspection{},
		&models.Product{},
		&models.Order{},
		&models.Post{},
		&models.Comment{},
		&models.Notification{},
		&models.OperationLog{},
		&models.Upload{},
	}

	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			log.Printf("Migration warning for %T: %v", table, err)
		}
	}

	log.Println("Database migrations completed")
}
