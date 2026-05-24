package main

import (
	"log"
	"survey-platform/internal/router"
	"survey-platform/internal/utils"
)

func main() {
	utils.LoadConfig()

	db := utils.InitDB()
	utils.InitRedis()

	r := router.SetupRouter(db)

	addr := utils.GetConfig().ServerPort
	log.Printf("Server starting on port %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
