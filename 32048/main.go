package main

import (
	"log"
	"secondhand-trading/config"
	"secondhand-trading/middleware"
	"secondhand-trading/models"
	"secondhand-trading/routes"
	"secondhand-trading/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	models.InitDB()
	utils.InitJWT()

	go utils.StartAutoCloseTransactionJob()

	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())

	routes.SetupRoutes(r)

	log.Printf("Server starting on port %s", config.AppConfig.ServerPort)
	if err := r.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
