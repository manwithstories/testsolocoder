package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"meeting-room/internal/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Cfg.Redis.Host + ":" + config.Cfg.Redis.Port,
		Password: config.Cfg.Redis.Password,
		DB:       config.Cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(Ctx, 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("连接Redis失败: %v", err)
	}

	log.Println("Redis连接成功")
}

func RedisSet(key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(Ctx, key, value, expiration).Err()
}

func RedisGet(key string) (string, error) {
	return RedisClient.Get(Ctx, key).Result()
}

func RedisDel(key string) error {
	return RedisClient.Del(Ctx, key).Err()
}

func RedisHSet(key, field string, value interface{}) error {
	return RedisClient.HSet(Ctx, key, field, value).Err()
}

func RedisHGetAll(key string) (map[string]string, error) {
	return RedisClient.HGetAll(Ctx, key).Result()
}

func RedisLPush(key string, value interface{}) error {
	return RedisClient.LPush(Ctx, key, value).Err()
}

func RedisBRPop(key string, timeout time.Duration) ([]string, error) {
	return RedisClient.BRPop(Ctx, timeout, key).Result()
}

func RedisZAdd(key string, score float64, member interface{}) error {
	return RedisClient.ZAdd(Ctx, key, redis.Z{Score: score, Member: fmt.Sprint(member)}).Err()
}

func RedisZRange(key string, start, stop int64) ([]string, error) {
	return RedisClient.ZRange(Ctx, key, start, stop).Result()
}

func RedisExpire(key string, expiration time.Duration) error {
	return RedisClient.Expire(Ctx, key, expiration).Err()
}
