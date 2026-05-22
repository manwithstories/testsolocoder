package redis

import (
	"context"
	"fmt"
	"strconv"
	"ticket-system/config"
	"ticket-system/internal/logger"
	"time"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client
var Ctx = context.Background()

func Init() {
	db, _ := strconv.Atoi(strconv.Itoa(config.App.RedisDB))
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.App.RedisHost, config.App.RedisPort),
		Password: config.App.RedisPassword,
		DB:       db,
	})

	if err := Client.Ping(Ctx).Err(); err != nil {
		logger.Log.Fatalf("Failed to connect to Redis: %v", err)
	}

	logger.Log.Info("Redis connected successfully")
}

func DecrementStock(ticketTypeID uint, quantity int) (int64, error) {
	key := fmt.Sprintf("ticket:stock:%d", ticketTypeID)
	return Client.DecrBy(Ctx, key, int64(quantity)).Result()
}

func IncrementStock(ticketTypeID uint, quantity int) error {
	key := fmt.Sprintf("ticket:stock:%d", ticketTypeID)
	return Client.IncrBy(Ctx, key, int64(quantity)).Err()
}

func GetStock(ticketTypeID uint) (int64, error) {
	key := fmt.Sprintf("ticket:stock:%d", ticketTypeID)
	return Client.Get(Ctx, key).Int64()
}

func SetStock(ticketTypeID uint, stock int) error {
	key := fmt.Sprintf("ticket:stock:%d", ticketTypeID)
	return Client.Set(Ctx, key, stock, 0).Err()
}

func LockOrder(orderNo string, expiration time.Duration) (bool, error) {
	key := fmt.Sprintf("order:lock:%s", orderNo)
	return Client.SetNX(Ctx, key, 1, expiration).Result()
}

func UnlockOrder(orderNo string) error {
	key := fmt.Sprintf("order:lock:%s", orderNo)
	return Client.Del(Ctx, key).Err()
}
