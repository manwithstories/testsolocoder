package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"event-platform/internal/config"
	"event-platform/internal/database"
	applog "event-platform/internal/logger"
	"event-platform/internal/middleware"
	"event-platform/internal/model"
	"event-platform/internal/queue"
	"event-platform/internal/repository"
	"event-platform/internal/router"
	"event-platform/pkg/jwt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfgPath := flag.String("config", "configs/config.yaml", "config file path")
	flag.Parse()

	cfg := config.Load(*cfgPath)
	applog.Init(cfg.Server.LogLevel, "./logs")

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(applog.GinLog, middleware.Recovery(), cors.Default())

	if _, err := database.Init(cfg.MySQL); err != nil {
		applog.Errorf("mysql init failed: %v", err)
	} else if err := database.DB.AutoMigrate(
		&model.User{},
		&model.Event{},
		&model.EventItem{},
		&model.Registration{},
		&model.Score{},
		&model.Certificate{},
		&model.Message{},
		&model.OperationLog{},
	); err != nil {
		applog.Errorf("auto migrate error: %v", err)
	}

	if _, err := database.InitRedis(cfg.Redis); err != nil {
		applog.Warnf("redis init failed: %v", err)
	}

	jm := jwt.New(cfg.JWT.Secret, cfg.JWT.ExpireMinutes)
	msgRepo := repository.NewMessageRepo()
	certRepo := repository.NewCertificateRepo()
	q := queue.New(msgRepo, certRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go q.ConsumeMessages(ctx)

	router.Register(r, cfg, jm, q)

	srv := &http.Server{Addr: ":" + cfg.Server.Port, Handler: r}
	go func() {
		applog.Infof("server listening on :%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			applog.Errorf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	applog.Infof("shutdown...")
	cancel()
	shutdownCtx, release := context.WithTimeout(context.Background(), 5*time.Second)
	defer release()
	_ = srv.Shutdown(shutdownCtx)
}
