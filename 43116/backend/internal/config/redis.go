package config

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis(cfg *RedisConfig) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接失败: %w", err)
	}

	return nil
}

func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}

func Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return RedisClient.Set(Ctx, key, data, expiration).Err()
}

func Get(key string, dest interface{}) error {
	data, err := RedisClient.Get(Ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func Delete(key string) error {
	return RedisClient.Del(Ctx, key).Err()
}

func Exists(key string) bool {
	result, err := RedisClient.Exists(Ctx, key).Result()
	if err != nil {
		return false
	}
	return result > 0
}

func SetNX(key string, value interface{}, expiration time.Duration) bool {
	data, err := json.Marshal(value)
	if err != nil {
		return false
	}
	result, err := RedisClient.SetNX(Ctx, key, data, expiration).Result()
	if err != nil {
		return false
	}
	return result
}

func LPush(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return RedisClient.LPush(Ctx, key, data).Err()
}

func BLPop(key string, timeout time.Duration) ([]byte, error) {
	result, err := RedisClient.BLPop(Ctx, timeout, key).Result()
	if err != nil {
		return nil, err
	}
	if len(result) < 2 {
		return nil, fmt.Errorf("BLPop returned unexpected result")
	}
	return []byte(result[1]), nil
}