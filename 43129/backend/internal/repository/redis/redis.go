package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client
var Ctx = context.Background()

func Init(host string, port int, password string, database int, poolSize int) error {
	addr := fmt.Sprintf("%s:%d", host, port)

	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
		PoolSize: poolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
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

func Exists(key string) (bool, error) {
	result, err := Client.Exists(Ctx, key).Result()
	return result > 0, err
}

func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return Client.SetNX(Ctx, key, value, expiration).Result()
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

func HDel(key, field string) error {
	return Client.HDel(Ctx, key, field).Err()
}

func SAdd(key string, members ...interface{}) error {
	return Client.SAdd(Ctx, key, members...).Err()
}

func SMembers(key string) ([]string, error) {
	return Client.SMembers(Ctx, key).Result()
}

func SRem(key string, members ...interface{}) error {
	return Client.SRem(Ctx, key, members...).Err()
}
