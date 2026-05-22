package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"auction-system/config"
	"auction-system/pkg/logger"
)

var Client *redis.Client
var ctx = context.Background()

func InitRedis() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.AppConfig.Redis.Host, config.AppConfig.Redis.Port),
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})

	if err := Client.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis: %v", err)
		return err
	}

	logger.Info("Redis connected successfully")
	return nil
}

func Get(key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

func Set(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

func Del(key string) error {
	return Client.Del(ctx, key).Err()
}

func Exists(key string) (bool, error) {
	result, err := Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func IncrBy(key string, value int64) (int64, error) {
	return Client.IncrBy(ctx, key, value).Result()
}

func ZAdd(key string, score float64, member string) error {
	return Client.ZAdd(ctx, key, &redis.Z{Score: score, Member: member}).Err()
}

func ZRevRange(key string, start, stop int64) ([]string, error) {
	return Client.ZRevRange(ctx, key, start, stop).Result()
}

func ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	return Client.ZRangeWithScores(ctx, key, start, stop).Result()
}

func Lock(key string, expiration time.Duration) (bool, error) {
	return Client.SetNX(ctx, key, "locked", expiration).Result()
}

func Unlock(key string) error {
	return Del(key)
}

func Publish(channel string, message interface{}) error {
	return Client.Publish(ctx, channel, message).Err()
}

func Subscribe(channels ...string) *redis.PubSub {
	return Client.Subscribe(ctx, channels...)
}
