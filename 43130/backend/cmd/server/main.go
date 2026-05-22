package main

import (
	"fmt"
	"log"
	"wedding-planner/config"
	"wedding-planner/internal/middleware"
	"wedding-planner/internal/services"
	"wedding-planner/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	gin.SetMode(cfg.Server.Mode)

	database.InitDB()

	scheduler := services.NewSchedulerService()
	scheduler.StartScheduler()

	r := gin.Default()

	r.Use(middleware.OperationLogger())

	SetupRouter(r)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("API Documentation available at http://localhost%s/api", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
