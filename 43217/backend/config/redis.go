package config

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     GlobalConfig.Redis.GetAddr(),
		Password: GlobalConfig.Redis.Password,
		DB:       GlobalConfig.Redis.DB,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}

	fmt.Println("Redis connection successful")
}
