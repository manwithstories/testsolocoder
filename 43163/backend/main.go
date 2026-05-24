package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"printshop/internal/config"
	"printshop/internal/database"
	"printshop/internal/handlers"
	"printshop/internal/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	seedData(db)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	h := handlers.NewHandler(db, cfg)
	h.RegisterRoutes(r)

	r.StaticFS("/uploads", gin.Dir("./uploads", true))

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func seedData(db *gorm.DB) {
	var roleCount int64
	db.Model(&models.Role{}).Count(&roleCount)
	if roleCount == 0 {
		roles := []models.Role{
			{Name: "admin", Description: "System Administrator"},
			{Name: "operator", Description: "Production Operator"},
			{Name: "sales", Description: "Sales Representative"},
			{Name: "customer", Description: "Customer User"},
		}
		db.Create(&roles)
	}

	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		var adminRole models.Role
		db.Where("name = ?", "admin").First(&adminRole)
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		db.Create(&models.User{
			Username:     "admin",
			PasswordHash: string(hash),
			RealName:     "Administrator",
			Email:        "admin@printshop.com",
			RoleID:       adminRole.ID,
			Active:       true,
		})
	}

	os.MkdirAll("./uploads", 0755)
}
