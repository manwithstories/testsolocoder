package main

import (
	"flag"
	"log"
	"time"

	"housekeeping/config"
	"housekeeping/database"
	"housekeeping/handlers"
	"housekeeping/router"
	"housekeeping/utils"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "config file path")
	flag.Parse()
	if _, err := config.Load(*cfgPath); err != nil {
		log.Fatalf("load config failed: %v", err)
	}
	utils.InitLogger(config.C.App.Mode)
	if _, err := database.Init(config.C.Database); err != nil {
		log.Fatalf("init db failed: %v", err)
	}
	go schedule()
	r := router.New()
	addr := ":" + config.C.App.Port
	utils.Logger.Infow("starting server", "addr", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func schedule() {
	t := time.NewTicker(10 * time.Minute)
	for range t.C {
		handlers.AutoConfirmOrders()
		handlers.EscalateTickets()
	}
}
