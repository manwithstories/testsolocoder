package redis

import (
	"context"
	"fmt"
	"log"
	"photo-rental/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var Ctx = context.Background()

func InitRedis(cfg config.RedisConfig) {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(Ctx, 5*time.Second)
	defer cancel()

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connection established successfully")
}

func SetEquipmentAvailability(equipmentID uint, dates []string, orderID uint) error {
	ctx := Ctx
	pipe := Client.Pipeline()

	for _, date := range dates {
		key := fmt.Sprintf("equipment:%d:availability", equipmentID)
		pipe.Set(ctx, fmt.Sprintf("%s:%s", key, date), orderID, 24*time.Hour*365)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func CheckEquipmentConflict(equipmentID uint, startDate, endDate time.Time) ([]string, error) {
	ctx := Ctx
	var conflicts []string

	current := startDate
	for !current.After(endDate) {
		dateStr := current.Format("2006-01-02")
		key := fmt.Sprintf("equipment:%d:availability:%s", equipmentID, dateStr)

		val, err := Client.Exists(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		if val > 0 {
			conflicts = append(conflicts, dateStr)
		}

		current = current.AddDate(0, 0, 1)
	}

	return conflicts, nil
}

func GetEquipmentReservedDates(equipmentID uint) ([]string, error) {
	ctx := Ctx
	key := fmt.Sprintf("equipment:%d:availability:*", equipmentID)

	var dates []string
	iter := Client.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		dates = append(dates, iter.Val())
	}

	return dates, iter.Err()
}

func RemoveEquipmentAvailability(equipmentID uint, startDate, endDate time.Time) error {
	ctx := Ctx
	pipe := Client.Pipeline()

	current := startDate
	for !current.After(endDate) {
		dateStr := current.Format("2006-01-02")
		key := fmt.Sprintf("equipment:%d:availability:%s", equipmentID, dateStr)
		pipe.Del(ctx, key)
		current = current.AddDate(0, 0, 1)
	}

	_, err := pipe.Exec(ctx)
	return err
}
