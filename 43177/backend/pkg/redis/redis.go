package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"repair-platform/pkg/config"
)

var Client *redis.Client
var Ctx = context.Background()

func InitRedis(cfg *config.Config) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(Ctx, 5*time.Second)
	defer cancel()

	_, err := Client.Ping(ctx).Result()
	return err
}

func Set(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(Ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return Client.Get(Ctx, key).Result()
}

func Delete(key string) error {
	return Client.Del(Ctx, key).Err()
}

func HSet(key, field string, value interface{}) error {
	return Client.HSet(Ctx, key, field, value).Err()
}

func HGet(key, field string) (string, error) {
	return Client.HGet(Ctx, key, field).Result()
}

func HGetAll(key string) (map[string]string, error) {
	return Client.HGetAll(Ctx, key).Result()
}

func Incr(key string) (int64, error) {
	return Client.Incr(Ctx, key).Result()
}

func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return Client.SetNX(Ctx, key, value, expiration).Result()
}

func Keys(pattern string) ([]string, error) {
	return Client.Keys(Ctx, pattern).Result()
}
