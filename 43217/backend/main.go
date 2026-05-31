package main

import (
	"fmt"
	"log"

	"health-platform/config"
	"health-platform/cron"
	"health-platform/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	config.InitDatabase()

	config.InitRedis()

	gin.SetMode(config.GlobalConfig.Server.Mode)

	r := gin.Default()

	routes.SetupRouter(r)

	scheduler := cron.NewScheduler()
	scheduler.Start()
	defer scheduler.Stop()

	addr := fmt.Sprintf(":%s", config.GlobalConfig.Server.Port)
	log.Printf("Server starting on port %s", config.GlobalConfig.Server.Port)
	
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
