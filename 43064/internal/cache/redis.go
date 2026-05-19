package cache

import (
	"context"
	"time"

	"github.com/notification-center/internal/config"
	"github.com/notification-center/internal/logger"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var Client *redis.Client
var Ctx = context.Background()

func Init(cfg *config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(Ctx, 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		logger.Error("redis connection failed", zap.Error(err))
		return err
	}

	logger.Info("redis connected successfully")
	return nil
}

func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

func Get(key string) (string, error) {
	return Client.Get(Ctx, key).Result()
}

func Set(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(Ctx, key, value, expiration).Err()
}

func Del(key string) error {
	return Client.Del(Ctx, key).Err()
}

func Exists(key string) (bool, error) {
	result, err := Client.Exists(Ctx, key).Result()
	return result > 0, err
}

func Incr(key string) (int64, error) {
	return Client.Incr(Ctx, key).Result()
}

func Expire(key string, expiration time.Duration) (bool, error) {
	return Client.Expire(Ctx, key, expiration).Result()
}

func ZAdd(key string, score float64, member string) (int64, error) {
	return Client.ZAdd(Ctx, key, &redis.Z{Score: score, Member: member}).Result()
}

func ZRangeByScore(key string, min, max string, offset, count int64) ([]string, error) {
	return Client.ZRangeByScore(Ctx, key, &redis.ZRangeBy{
		Min: min, Max: max, Offset: offset, Count: count,
	}).Result()
}

func ZRem(key string, members ...interface{}) (int64, error) {
	return Client.ZRem(Ctx, key, members...).Result()
}

func LPush(key string, values ...interface{}) (int64, error) {
	return Client.LPush(Ctx, key, values...).Result()
}

func RPop(key string) (string, error) {
	return Client.RPop(Ctx, key).Result()
}

func LLen(key string) (int64, error) {
	return Client.LLen(Ctx, key).Result()
}
