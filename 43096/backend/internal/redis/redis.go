package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"ticket-system/internal/config"
)

var Client *redis.Client
var ctx = context.Background()

func InitRedis() {
	cfg := config.AppConfig.Redis
	Client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	fmt.Println("Redis connected successfully")
}

func GetCtx() context.Context {
	return ctx
}

const (
	KeyPrefixSeatLock    = "seat:lock:"
	KeyPrefixSeatStatus  = "seat:status:"
	KeyPrefixSessionSeats = "session:seats:"
	KeyPrefixUserInfo    = "user:info:"
	KeyPrefixOrderLock   = "order:lock:"
)

func LockSeat(sessionID, seatID uint64, userID uint64, lockMinutes int) (bool, error) {
	key := fmt.Sprintf("%s%d:%d", KeyPrefixSeatLock, sessionID, seatID)
	locked, err := Client.SetNX(ctx, key, userID, time.Duration(lockMinutes)*time.Minute).Result()
	if err != nil {
		return false, err
	}
	return locked, nil
}

func UnlockSeat(sessionID, seatID uint64) error {
	key := fmt.Sprintf("%s%d:%d", KeyPrefixSeatLock, sessionID, seatID)
	return Client.Del(ctx, key).Err()
}

func GetSeatLockUser(sessionID, seatID uint64) (uint64, error) {
	key := fmt.Sprintf("%s%d:%d", KeyPrefixSeatLock, sessionID, seatID)
	result, err := Client.Get(ctx, key).Uint64()
	if err == redis.Nil {
		return 0, nil
	}
	return result, err
}

func SetSeatStatus(sessionID, seatID uint64, status int) error {
	key := fmt.Sprintf("%s%d", KeyPrefixSeatStatus, sessionID)
	return Client.HSet(ctx, key, fmt.Sprintf("%d", seatID), status).Err()
}

func GetSeatStatus(sessionID, seatID uint64) (int, error) {
	key := fmt.Sprintf("%s%d", KeyPrefixSeatStatus, sessionID)
	result, err := Client.HGet(ctx, key, fmt.Sprintf("%d", seatID)).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return result, err
}

func GetAllSeatStatuses(sessionID uint64) (map[string]string, error) {
	key := fmt.Sprintf("%s%d", KeyPrefixSeatStatus, sessionID)
	return Client.HGetAll(ctx, key).Result()
}

func CacheUserInfo(userID uint64, userJSON string) error {
	key := fmt.Sprintf("%s%d", KeyPrefixUserInfo, userID)
	return Client.Set(ctx, key, userJSON, time.Hour).Err()
}

func GetCachedUserInfo(userID uint64) (string, error) {
	key := fmt.Sprintf("%s%d", KeyPrefixUserInfo, userID)
	return Client.Get(ctx, key).Result()
}

func DeleteCachedUserInfo(userID uint64) error {
	key := fmt.Sprintf("%s%d", KeyPrefixUserInfo, userID)
	return Client.Del(ctx, key).Err()
}

func ExtendSeatLock(sessionID, seatID uint64, lockMinutes int) error {
	key := fmt.Sprintf("%s%d:%d", KeyPrefixSeatLock, sessionID, seatID)
	return Client.Expire(ctx, key, time.Duration(lockMinutes)*time.Minute).Err()
}

func GetSeatLockTTL(sessionID, seatID uint64) (time.Duration, error) {
	key := fmt.Sprintf("%s%d:%d", KeyPrefixSeatLock, sessionID, seatID)
	return Client.TTL(ctx, key).Result()
}
