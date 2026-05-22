package redis

import (
	"context"
	"fmt"
	"log"
	"property-management/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var Ctx = context.Background()

func Init(cfg *config.Config) {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect Redis:", err)
	}

	log.Println("Redis connected successfully")
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

func CachePropertyStatus(propertyID uint, status int) {
	key := fmt.Sprintf("property:status:%d", propertyID)
	Client.Set(Ctx, key, status, 24*time.Hour)
}

func GetPropertyStatus(propertyID uint) (int, error) {
	key := fmt.Sprintf("property:status:%d", propertyID)
	return Client.Get(Ctx, key).Int()
}

func AddContractReminder(contractID uint, endDate time.Time) {
	member := fmt.Sprintf("contract:%d", contractID)
	score := float64(endDate.Unix())
	Client.ZAdd(Ctx, "contract:expiry", redis.Z{Score: score, Member: member})
}

func GetExpiringContracts(from, to time.Time) []string {
	min := fmt.Sprintf("%d", from.Unix())
	max := fmt.Sprintf("%d", to.Unix())
	members, _ := Client.ZRangeByScore(Ctx, "contract:expiry", &redis.ZRangeBy{Min: min, Max: max}).Result()
	return members
}

func RemoveContractReminder(contractID uint) {
	member := fmt.Sprintf("contract:%d", contractID)
	Client.ZRem(Ctx, "contract:expiry", member)
}
