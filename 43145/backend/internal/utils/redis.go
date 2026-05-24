package utils

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisCtx context.Context

func InitRedis() {
	cfg := GetConfig()
	RedisCtx = context.Background()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	_, err := RedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Printf("Redis connection failed: %v, continuing without cache", err)
	} else {
		log.Println("Redis connection established")
	}
}

func CacheSet(key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.Set(RedisCtx, key, value, expiration).Err()
}

func CacheGet(key string) (string, error) {
	if RedisClient == nil {
		return "", nil
	}
	return RedisClient.Get(RedisCtx, key).Result()
}

func CacheDelete(key string) error {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.Del(RedisCtx, key).Err()
}

func CacheDeleteByPattern(pattern string) error {
	if RedisClient == nil {
		return nil
	}
	var cursor uint64
	for {
		keys, nextCursor, err := RedisClient.Scan(RedisCtx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := RedisClient.Del(RedisCtx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

func CacheGetJSON(key string, dest interface{}) error {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.Get(RedisCtx, key).Scan(dest)
}

func CacheSetJSON(key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.Set(RedisCtx, key, value, expiration).Err()
}
