package redis

import (
	"context"
	"errand-service/internal/config"
	"errand-service/pkg/logger"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func Init(cfg config.RedisConfig) (*redis.Client, error) {
	client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	logger.Info("Redis connected successfully")
	return client, nil
}

func GetClient() *redis.Client {
	return client
}

func Close() {
	if client != nil {
		client.Close()
		logger.Info("Redis connection closed")
	}
}

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}

func Get(ctx context.Context, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func Delete(ctx context.Context, key string) error {
	return client.Del(ctx, key).Err()
}

func GeoAdd(ctx context.Context, key string, longitude, latitude float64, member string) error {
	return client.GeoAdd(ctx, key, &redis.GeoLocation{
		Longitude: longitude,
		Latitude:  latitude,
		Name:      member,
	}).Err()
}

func GeoRadius(ctx context.Context, key string, longitude, latitude float64, radius float64, unit string) ([]redis.GeoLocation, error) {
	return client.GeoRadius(ctx, key, longitude, latitude, &redis.GeoRadiusQuery{
		Radius: radius,
		Unit:   unit,
		Sort:   "ASC",
	}).Result()
}

func GeoRemove(ctx context.Context, key string, member string) error {
	return client.ZRem(ctx, key, member).Err()
}

func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return client.SetNX(ctx, key, value, expiration).Result()
}
