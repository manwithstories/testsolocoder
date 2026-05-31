package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"consultation-platform/config"
)

var RedisClient *redis.Client

func InitRedis(cfg config.RedisConfig) (*redis.Client, error) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	return RedisClient, nil
}

type DistributedLock struct {
	client *redis.Client
	key    string
	value  string
}

func NewDistributedLock(key string) *DistributedLock {
	return &DistributedLock{
		client: RedisClient,
		key:    key,
		value:  time.Now().String(),
	}
}

func (dl *DistributedLock) Lock(ctx context.Context, expiration time.Duration) (bool, error) {
	ok, err := dl.client.SetNX(ctx, dl.key, dl.value, expiration).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (dl *DistributedLock) Unlock(ctx context.Context) error {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
	_, err := dl.client.Eval(ctx, script, []string{dl.key}, dl.value).Result()
	return err
}
