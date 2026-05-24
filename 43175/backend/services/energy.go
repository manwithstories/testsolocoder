package services

import (
	"fmt"
	"log"
	"math/rand"
	"smart-energy-platform/models"
	"time"
)

func InitEnergyCollector() {
	log.Println("Energy data collector started")

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		collectEnergyData()
		checkEnergyAlerts()
	}
}

func collectEnergyData() {
	var devices []models.Device
	if err := models.DB.Where("status = ?", "online").Find(&devices).Error; err != nil {
		log.Printf("Error collecting energy data: %v", err)
		return
	}

	now := time.Now()
	dateStr := now.Format("2006-01-02")

	for _, device := range devices {
		if device.Power <= 0 {
			continue
		}

		power := device.Power * (0.8 + rand.Float64()*0.4)
		voltage := 220.0 + rand.Float64()*10 - 5
		current := power / voltage
		energyUsed := power * (5.0 / 60.0)

		energyData := models.EnergyData{
			DeviceID:   device.ID,
			FamilyID:   device.FamilyID,
			Power:      power,
			Voltage:    voltage,
			Current:    current,
			EnergyUsed: energyUsed,
			Timestamp:  now,
			Date:       dateStr,
			Hour:       now.Hour(),
		}

		if err := models.DB.Create(&energyData).Error; err != nil {
			log.Printf("Error saving energy data for device %d: %v", device.ID, err)
		}

		UpdateDeviceEnergy(device.ID, power)
	}
}

func checkEnergyAlerts() {
	var devices []models.Device
	if err := models.DB.Find(&devices).Error; err != nil {
		return
	}

	now := time.Now()
	dateStr := now.Format("2006-01-02")

	for _, device := range devices {
		var totalEnergy float64
		models.DB.Model(&models.EnergyData{}).
			Where("device_id = ? AND date = ?", device.ID, dateStr).
			Select("COALESCE(SUM(energy_used), 0)").
			Scan(&totalEnergy)

		dailyThreshold := device.Power * 24 * 0.8

		if totalEnergy > dailyThreshold {
			var existingAlert models.EnergyAlert
			result := models.DB.Where("device_id = ? AND alert_type = ? AND resolved = ?",
				device.ID, "high_consumption", false).First(&existingAlert)

			if result.Error != nil {
				alert := models.EnergyAlert{
					FamilyID:  device.FamilyID,
					DeviceID:  device.ID,
					AlertType: "high_consumption",
					Level:     "warning",
					Message:   "Device daily energy consumption exceeds threshold",
					Value:     totalEnergy,
					Threshold: dailyThreshold,
				}
				models.DB.Create(&alert)

				var members []models.FamilyMember
				models.DB.Where("family_id = ? AND status = ?", device.FamilyID, 1).Find(&members)

				for _, member := range members {
					notification := models.Notification{
						UserID:  member.UserID,
						Type:    "alert",
						Title:   "High Energy Consumption Alert",
						Content: fmt.Sprintf("Device '%s' has consumed %.2f kWh today, exceeding the daily threshold.", device.Name, totalEnergy),
					}
					models.DB.Create(&notification)
				}
			}
		}
	}
}
