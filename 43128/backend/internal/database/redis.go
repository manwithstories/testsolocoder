package database

import (
	"context"
	"time"

	"event-platform/internal/config"
	applog "event-platform/internal/logger"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis(cfg config.RedisConfig) (*redis.Client, error) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := RDB.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	applog.Infof("redis connected: %s", cfg.Addr)
	return RDB, nil
}
