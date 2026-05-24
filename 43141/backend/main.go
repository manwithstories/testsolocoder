package main

import (
	"log"
	"os"

	"sports-league/config"
	"sports-league/models"
	"sports-league/pkg/database"
	"sports-league/pkg/logger"
	"sports-league/routes"
	"sports-league/seeds"
)

func main() {
	config.Load()
	logger.Init()

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("connect database failed: %v", err)
	}
	if err := models.Migrate(db); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}
	seeds.SeedAdmin(db)

	r := routes.Setup(db)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("server listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}
