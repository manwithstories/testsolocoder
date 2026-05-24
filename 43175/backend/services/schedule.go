package services

import (
	"fmt"
	"log"
	"smart-energy-platform/models"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	cronInstance *cron.Cron
	scheduleMap  = make(map[uint]cron.EntryID)
	scheduleLock sync.Mutex
)

func InitScheduleExecutor() {
	log.Println("Schedule executor started")

	cronInstance = cron.New()
	cronInstance.Start()

	loadSchedules()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		loadSchedules()
	}
}

func loadSchedules() {
	var schedules []models.Schedule
	if err := models.DB.Where("is_enabled = ?", true).Find(&schedules).Error; err != nil {
		log.Printf("Error loading schedules: %v", err)
		return
	}

	scheduleLock.Lock()
	defer scheduleLock.Unlock()

	existingIDs := make(map[uint]bool)

	for _, sched := range schedules {
		existingIDs[sched.ID] = true

		if _, exists := scheduleMap[sched.ID]; exists {
			continue
		}

		schedCopy := sched
		entryID, err := cronInstance.AddFunc(sched.CronExpr, func() {
			executeSchedule(schedCopy)
		})
		if err != nil {
			log.Printf("Error adding schedule %d: %v", sched.ID, err)
			continue
		}

		scheduleMap[sched.ID] = entryID
		log.Printf("Schedule loaded: %d - %s", sched.ID, sched.Name)
	}

	for id, entryID := range scheduleMap {
		if !existingIDs[id] {
			cronInstance.Remove(entryID)
			delete(scheduleMap, id)
		}
	}
}

func executeSchedule(schedule models.Schedule) {
	log.Printf("Executing schedule: %s (ID: %d)", schedule.Name, schedule.ID)

	var device models.Device
	if err := models.DB.First(&device, schedule.DeviceID).Error; err != nil {
		logScheduleError(schedule, "Device not found")
		return
	}

	success := true
	message := fmt.Sprintf("Schedule executed successfully, action: %s", schedule.Action)

	oldStatus := device.Status
	device.Status = schedule.Action

	if err := models.DB.Save(&device).Error; err != nil {
		success = false
		message = fmt.Sprintf("Failed to update device: %v", err)
	} else {
		UpdateDeviceStatus(device.ID, schedule.Action)
	}

	energyDelta := 0.0
	if schedule.Action == "on" && oldStatus != "on" {
		energyDelta = device.Power * 0.5
	}

	saveScheduleLog(schedule, device, success, message, energyDelta)

	now := time.Now()
	models.DB.Model(&models.Schedule{}).Where("id = ?", schedule.ID).Update("last_run", &now)

	log.Printf("Schedule %s executed: %s", schedule.Name, message)
}

func logScheduleError(schedule models.Schedule, errMsg string) {
	saveScheduleLog(schedule, models.Device{}, false, errMsg, 0)
}

func saveScheduleLog(schedule models.Schedule, device models.Device, success bool, message string, energyDelta float64) {
	logEntry := models.ScheduleLog{
		ScheduleID:  schedule.ID,
		DeviceID:    schedule.DeviceID,
		Action:      schedule.Action,
		Value:       schedule.Value,
		Success:     success,
		Message:     message,
		EnergyDelta: energyDelta,
		ExecutedAt:  time.Now(),
	}
	models.DB.Create(&logEntry)
}

func ReloadSchedule(scheduleID uint) {
	scheduleLock.Lock()
	defer scheduleLock.Unlock()

	if entryID, exists := scheduleMap[scheduleID]; exists {
		cronInstance.Remove(entryID)
		delete(scheduleMap, scheduleID)
	}

	var schedule models.Schedule
	if err := models.DB.First(&schedule, scheduleID).Error; err != nil {
		return
	}

	if schedule.IsEnabled {
		schedCopy := schedule
		entryID, err := cronInstance.AddFunc(schedCopy.CronExpr, func() {
			executeSchedule(schedCopy)
		})
		if err != nil {
			log.Printf("Error reloading schedule %d: %v", schedule.ID, err)
			return
		}
		scheduleMap[schedule.ID] = entryID
	}
}

func RemoveSchedule(scheduleID uint) {
	scheduleLock.Lock()
	defer scheduleLock.Unlock()

	if entryID, exists := scheduleMap[scheduleID]; exists {
		cronInstance.Remove(entryID)
		delete(scheduleMap, scheduleID)
	}
}

func CheckTimeConflict(cronExpr string, excludeID *uint) bool {
	var schedules []models.Schedule
	query := models.DB.Where("cron_expr = ? AND is_enabled = ?", cronExpr, true)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	query.Find(&schedules)
	return len(schedules) > 0
}

func ParseCronField(expr string) error {
	_, err := cron.ParseStandard(expr)
	return err
}
