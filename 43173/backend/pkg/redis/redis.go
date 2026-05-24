package redis

import (
	"context"
	"fmt"
	"time"

	"music-platform/pkg/config"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis(cfg *config.RedisConfig) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	return nil
}

func CloseRedis() error {
	return RDB.Close()
}

func Set(key string, value interface{}, expiration time.Duration) error {
	return RDB.Set(Ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return RDB.Get(Ctx, key).Result()
}

func Delete(key string) error {
	return RDB.Del(Ctx, key).Err()
}

func Incr(key string) (int64, error) {
	return RDB.Incr(Ctx, key).Result()
}

func ZAdd(key string, score float64, member interface{}) error {
	memberStr := fmt.Sprintf("%v", member)
	return RDB.ZAdd(Ctx, key, redis.Z{
		Score:  score,
		Member: memberStr,
	}).Err()
}

func ZIncrBy(key string, increment float64, member interface{}) (float64, error) {
	memberStr := fmt.Sprintf("%v", member)
	return RDB.ZIncrBy(Ctx, key, increment, memberStr).Result()
}

func ZRevRange(key string, start, stop int64) ([]string, error) {
	return RDB.ZRevRange(Ctx, key, start, stop).Result()
}

func ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	return RDB.ZRevRangeWithScores(Ctx, key, start, stop).Result()
}

func ZRank(key string, member interface{}) (int64, error) {
	memberStr := fmt.Sprintf("%v", member)
	return RDB.ZRevRank(Ctx, key, memberStr).Result()
}

func ZScore(key string, member interface{}) (float64, error) {
	memberStr := fmt.Sprintf("%v", member)
	return RDB.ZScore(Ctx, key, memberStr).Result()
}

func AcquireLock(key string, value string, expiration time.Duration) (bool, error) {
	return RDB.SetNX(Ctx, key, value, expiration).Result()
}

func ReleaseLock(key string, value string) error {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
	return RDB.Eval(Ctx, script, []string{key}, value).Err()
}

func HSet(key string, field string, value interface{}) error {
	return RDB.HSet(Ctx, key, field, value).Err()
}

func HGet(key string, field string) (string, error) {
	return RDB.HGet(Ctx, key, field).Result()
}

func HGetAll(key string) (map[string]string, error) {
	return RDB.HGetAll(Ctx, key).Result()
}

func Expire(key string, expiration time.Duration) error {
	return RDB.Expire(Ctx, key, expiration).Err()
}

func Keys(pattern string) ([]string, error) {
	return RDB.Keys(Ctx, pattern).Result()
}
