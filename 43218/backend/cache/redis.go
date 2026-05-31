package cache

import (
	"context"
	"fmt"
	"time"

	"secondhand-platform/config"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client

func InitRedis(cfg *config.RedisConfig) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.GetAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	logrus.Info("Redis connection established successfully")
	return nil
}

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

func Get(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

func Delete(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

func Exists(ctx context.Context, key string) (bool, error) {
	result, err := RedisClient.Exists(ctx, key).Result()
	return result > 0, err
}

func SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

func GetJSON(ctx context.Context, key string, dest interface{}) error {
	return RedisClient.Get(ctx, key).Scan(dest)
}

func HSet(ctx context.Context, key, field string, value interface{}) error {
	return RedisClient.HSet(ctx, key, field, value).Err()
}

func HGet(ctx context.Context, key, field string) (string, error) {
	return RedisClient.HGet(ctx, key, field).Result()
}

func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return RedisClient.HGetAll(ctx, key).Result()
}

func HDel(ctx context.Context, key string, fields ...string) error {
	return RedisClient.HDel(ctx, key, fields...).Err()
}

func LPush(ctx context.Context, key string, values ...interface{}) error {
	return RedisClient.LPush(ctx, key, values...).Err()
}

func LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return RedisClient.LRange(ctx, key, start, stop).Result()
}

func Incr(ctx context.Context, key string) (int64, error) {
	return RedisClient.Incr(ctx, key).Result()
}

func Decr(ctx context.Context, key string) (int64, error) {
	return RedisClient.Decr(ctx, key).Result()
}

func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return RedisClient.Expire(ctx, key, expiration).Err()
}

func TTL(ctx context.Context, key string) (time.Duration, error) {
	return RedisClient.TTL(ctx, key).Result()
}
