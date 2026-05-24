package services

import (
	"log"
	"smart-energy-platform/models"
	"strconv"
	"strings"
	"time"
)

func InitSceneTriggerWatcher() {
	log.Println("Scene trigger watcher started")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		checkSceneTriggers()
	}
}

func checkSceneTriggers() {
	var scenes []models.Scene
	if err := models.DB.Where("is_active = ?", true).
		Preload("Conditions").
		Preload("Actions").
		Find(&scenes).Error; err != nil {
		log.Printf("Error loading scenes: %v", err)
		return
	}

	now := time.Now()

	for _, scene := range scenes {
		if len(scene.Conditions) == 0 {
			continue
		}

		allConditionsMet := true

		for _, condition := range scene.Conditions {
			if !evaluateCondition(condition, now) {
				allConditionsMet = false
				break
			}
		}

		if allConditionsMet {
			executeSceneActions(scene)
		}
	}
}

func evaluateCondition(condition models.SceneCondition, now time.Time) bool {
	switch condition.Type {
	case "time":
		return evaluateTimeCondition(condition, now)
	case "device":
		return evaluateDeviceCondition(condition)
	case "sensor":
		return evaluateSensorCondition(condition)
	default:
		return false
	}
}

func evaluateTimeCondition(condition models.SceneCondition, now time.Time) bool {
	if condition.TimeExpr == "" {
		return false
	}

	parts := strings.Split(condition.TimeExpr, ":")
	if len(parts) != 2 {
		return false
	}

	targetHour, _ := strconv.Atoi(parts[0])
	targetMinute, _ := strconv.Atoi(parts[1])

	currentMinuteOfDay := now.Hour()*60 + now.Minute()
	targetMinuteOfDay := targetHour*60 + targetMinute

	diff := currentMinuteOfDay - targetMinuteOfDay
	return diff >= 0 && diff < 1
}

func evaluateDeviceCondition(condition models.SceneCondition) bool {
	if condition.DeviceID == nil {
		return false
	}

	var device models.Device
	if err := models.DB.First(&device, *condition.DeviceID).Error; err != nil {
		return false
	}

	deviceStatus := device.Status

	switch condition.Operator {
	case "eq":
		return deviceStatus == condition.Value
	case "neq":
		return deviceStatus != condition.Value
	default:
		return deviceStatus == condition.Value
	}
}

func evaluateSensorCondition(condition models.SceneCondition) bool {
	if condition.DeviceID == nil {
		return false
	}

	var device models.Device
	if err := models.DB.First(&device, *condition.DeviceID).Error; err != nil {
		return false
	}

	return device.Status == condition.Value
}

func executeSceneActions(scene models.Scene) {
	log.Printf("Executing scene: %s (ID: %d)", scene.Name, scene.ID)

	for _, action := range scene.Actions {
		var device models.Device
		if err := models.DB.First(&device, action.DeviceID).Error; err != nil {
			log.Printf("Error finding device %d for scene: %v", action.DeviceID, err)
			continue
		}

		oldStatus := device.Status
		device.Status = action.Action

		if err := models.DB.Save(&device).Error; err != nil {
			log.Printf("Error updating device %d: %v", action.DeviceID, err)
			continue
		}

		UpdateDeviceStatus(device.ID, action.Action)

		log.Printf("Scene action executed: device %d -> %s (was: %s)", device.ID, action.Action, oldStatus)
	}
}

func ExecuteSceneNow(sceneID uint) error {
	var scene models.Scene
	if err := models.DB.Preload("Actions").First(&scene, sceneID).Error; err != nil {
		return err
	}

	executeSceneActions(scene)
	return nil
}
