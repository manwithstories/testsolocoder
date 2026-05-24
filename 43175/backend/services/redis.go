package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"smart-energy-platform/config"
	"smart-energy-platform/models"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb      *redis.Client
	Ctx      = context.Background()
	Cfg      *config.Config
)

func InitRedis(cfg *config.Config) {
	Cfg = cfg
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v, will use fallback", err)
	} else {
		log.Println("Redis connected successfully")
	}
}

func DeviceStatusKey(deviceID uint) string {
	return fmt.Sprintf("device:status:%d", deviceID)
}

func DeviceEnergyKey(deviceID uint) string {
	return fmt.Sprintf("device:energy:%d", deviceID)
}

func FamilyDevicesKey(familyID uint) string {
	return fmt.Sprintf("family:devices:%d", familyID)
}

func UpdateDeviceStatus(deviceID uint, status string) {
	if Rdb == nil {
		return
	}
	data := map[string]interface{}{
		"deviceId": deviceID,
		"status":   status,
		"time":     time.Now().Format(time.RFC3339),
	}
	jsonData, _ := json.Marshal(data)
	Rdb.Set(Ctx, DeviceStatusKey(deviceID), jsonData, 24*time.Hour)
	Rdb.Publish(Ctx, fmt.Sprintf("device_status_%d", deviceID), jsonData)
}

func UpdateDeviceEnergy(deviceID uint, energy float64) {
	if Rdb == nil {
		return
	}
	data := map[string]interface{}{
		"deviceId": deviceID,
		"energy":   energy,
		"time":     time.Now().Format(time.RFC3339),
	}
	jsonData, _ := json.Marshal(data)
	Rdb.Set(Ctx, DeviceEnergyKey(deviceID), jsonData, 5*time.Minute)
}

func GetDeviceStatus(deviceID uint) (string, error) {
	if Rdb == nil {
		return "", nil
	}
	val, err := Rdb.Get(Ctx, DeviceStatusKey(deviceID)).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func GetDeviceEnergy(deviceID uint) (string, error) {
	if Rdb == nil {
		return "", nil
	}
	val, err := Rdb.Get(Ctx, DeviceEnergyKey(deviceID)).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func CacheFamilyDevices(familyID uint, devices []models.Device) {
	if Rdb == nil {
		return
	}
	jsonData, _ := json.Marshal(devices)
	Rdb.Set(Ctx, FamilyDevicesKey(familyID), jsonData, 5*time.Minute)
}

func GetCachedFamilyDevices(familyID uint) ([]models.Device, error) {
	if Rdb == nil {
		return nil, nil
	}
	val, err := Rdb.Get(Ctx, FamilyDevicesKey(familyID)).Result()
	if err != nil {
		return nil, err
	}
	var devices []models.Device
	json.Unmarshal([]byte(val), &devices)
	return devices, nil
}

func InvalidateFamilyDevicesCache(familyID uint) {
	if Rdb == nil {
		return
	}
	Rdb.Del(Ctx, FamilyDevicesKey(familyID))
}
