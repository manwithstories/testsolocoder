package handlers

import (
	"smart-energy-platform/models"
	"smart-energy-platform/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EnergyStatisticsResult struct {
	TotalEnergy     float64                `json:"totalEnergy"`
	ByDevice        []DeviceEnergyStat     `json:"byDevice"`
	ByHour          []HourlyEnergyStat     `json:"byHour"`
	PeakHour        int                    `json:"peakHour"`
	AvgPower        float64                `json:"avgPower"`
	ComparedYesterday string              `json:"comparedYesterday"`
}

type DeviceEnergyStat struct {
	DeviceID   uint    `json:"deviceId"`
	DeviceName string  `json:"deviceName"`
	TotalEnergy float64 `json:"totalEnergy"`
	Percentage float64 `json:"percentage"`
}

type HourlyEnergyStat struct {
	Hour   int     `json:"hour"`
	Energy float64 `json:"energy"`
}

type TrendPoint struct {
	Date   string  `json:"date"`
	Energy float64 `json:"energy"`
}

func GetRealtimeEnergy(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseQueryUint(c, "familyId")

	if familyID > 0 && !hasFamilyAccess(userID, familyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	var familyIDs []uint
	if familyID > 0 {
		familyIDs = []uint{familyID}
	} else {
		models.DB.Model(&models.FamilyMember{}).
			Where("user_id = ? AND status = ?", userID, 1).
			Pluck("family_id", &familyIDs)
	}

	var devices []models.Device
	if len(familyIDs) > 0 {
		models.DB.Where("family_id IN ?", familyIDs).Find(&devices)
	}

	var totalPower float64
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)

	deviceData := make([]gin.H, 0, len(devices))
	for _, device := range devices {
		var latestEnergy models.EnergyData
		models.DB.Where("device_id = ?", device.ID).
			Order("timestamp DESC").
			First(&latestEnergy)

		var hourlyEnergy float64
		models.DB.Model(&models.EnergyData{}).
			Where("device_id = ? AND timestamp >= ?", device.ID, oneHourAgo).
			Select("COALESCE(SUM(energy_used), 0)").
			Scan(&hourlyEnergy)

		totalPower += latestEnergy.Power

		deviceData = append(deviceData, gin.H{
			"deviceId":     device.ID,
			"deviceName":   device.Name,
			"deviceType":   device.DeviceType,
			"status":       device.Status,
			"power":        latestEnergy.Power,
			"hourlyEnergy": hourlyEnergy,
		})
	}

	utils.Success(c, gin.H{
		"timestamp":   now,
		"totalPower":  totalPower,
		"devices":     deviceData,
		"deviceCount": len(devices),
	})
}

func GetEnergyStatistics(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseQueryUint(c, "familyId")
	period := c.DefaultQuery("period", "day")

	if familyID > 0 && !hasFamilyAccess(userID, familyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	var familyIDs []uint
	if familyID > 0 {
		familyIDs = []uint{familyID}
	} else {
		models.DB.Model(&models.FamilyMember{}).
			Where("user_id = ? AND status = ?", userID, 1).
			Pluck("family_id", &familyIDs)
	}

	var startTime time.Time
	now := time.Now()
	switch period {
	case "week":
		startTime = now.Add(-7 * 24 * time.Hour)
	case "month":
		startTime = now.Add(-30 * 24 * time.Hour)
	default:
		startTime = now.Add(-24 * time.Hour)
	}

	var totalEnergy float64
	if len(familyIDs) > 0 {
		models.DB.Model(&models.EnergyData{}).
			Where("family_id IN ? AND timestamp >= ?", familyIDs, startTime).
			Select("COALESCE(SUM(energy_used), 0)").
			Scan(&totalEnergy)
	}

	var byDevice []DeviceEnergyStat
	if len(familyIDs) > 0 {
		type result struct {
			DeviceID    uint
			TotalEnergy float64
		}
		var results []result
		models.DB.Model(&models.EnergyData{}).
			Select("device_id, COALESCE(SUM(energy_used), 0) as total_energy").
			Where("family_id IN ? AND timestamp >= ?", familyIDs, startTime).
			Group("device_id").
			Scan(&results)

		for _, r := range results {
			var device models.Device
			models.DB.First(&device, r.DeviceID)
			pct := 0.0
			if totalEnergy > 0 {
				pct = r.TotalEnergy / totalEnergy * 100
			}
			byDevice = append(byDevice, DeviceEnergyStat{
				DeviceID:    r.DeviceID,
				DeviceName:  device.Name,
				TotalEnergy: r.TotalEnergy,
				Percentage:  pct,
			})
		}
	}

	var byHour []HourlyEnergyStat
	if len(familyIDs) > 0 && period == "day" {
		type hourResult struct {
			Hour   int
			Energy float64
		}
		var hourResults []hourResult
		models.DB.Model(&models.EnergyData{}).
			Select("hour, COALESCE(SUM(energy_used), 0) as energy").
			Where("family_id IN ? AND timestamp >= ?", familyIDs, startTime).
			Group("hour").
			Order("hour ASC").
			Scan(&hourResults)

		for _, r := range hourResults {
			byHour = append(byHour, HourlyEnergyStat{
				Hour:   r.Hour,
				Energy: r.Energy,
			})
		}
	}

	peakHour := 0
	maxEnergy := 0.0
	for _, h := range byHour {
		if h.Energy > maxEnergy {
			maxEnergy = h.Energy
			peakHour = h.Hour
		}
	}

	utils.Success(c, gin.H{
		"totalEnergy": totalEnergy,
		"byDevice":    byDevice,
		"byHour":      byHour,
		"peakHour":    peakHour,
		"period":      period,
	})
}

func GetEnergyTrend(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseQueryUint(c, "familyId")
	days := parseQueryInt(c, "days", 7)

	if familyID > 0 && !hasFamilyAccess(userID, familyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	var familyIDs []uint
	if familyID > 0 {
		familyIDs = []uint{familyID}
	} else {
		models.DB.Model(&models.FamilyMember{}).
			Where("user_id = ? AND status = ?", userID, 1).
			Pluck("family_id", &familyIDs)
	}

	startDate := time.Now().AddDate(0, 0, -days)

	type dailyResult struct {
		Date   string
		Energy float64
	}

	var results []dailyResult
	if len(familyIDs) > 0 {
		models.DB.Model(&models.EnergyData{}).
			Select("date, COALESCE(SUM(energy_used), 0) as energy").
			Where("family_id IN ? AND date >= ?", familyIDs, startDate.Format("2006-01-02")).
			Group("date").
			Order("date ASC").
			Scan(&results)
	}

	trendData := make([]TrendPoint, len(results))
	for i, r := range results {
		trendData[i] = TrendPoint{
			Date:   r.Date,
			Energy: r.Energy,
		}
	}

	var totalEnergy float64
	var avgEnergy float64
	for _, t := range trendData {
		totalEnergy += t.Energy
	}
	if len(trendData) > 0 {
		avgEnergy = totalEnergy / float64(len(trendData))
	}

	utils.Success(c, gin.H{
		"trend":      trendData,
		"total":      totalEnergy,
		"average":    avgEnergy,
		"days":       days,
	})
}

func ListEnergyAlerts(c *gin.Context) {
	userID := c.GetUint("userId")
	familyID := parseQueryUint(c, "familyId")
	resolved := c.Query("resolved")

	if familyID > 0 && !hasFamilyAccess(userID, familyID) {
		utils.Forbidden(c, "No access to this family")
		return
	}

	var familyIDs []uint
	if familyID > 0 {
		familyIDs = []uint{familyID}
	} else {
		models.DB.Model(&models.FamilyMember{}).
			Where("user_id = ? AND status = ?", userID, 1).
			Pluck("family_id", &familyIDs)
	}

	var alerts []models.EnergyAlert
	query := models.DB

	if len(familyIDs) > 0 {
		query = query.Where("family_id IN ?", familyIDs)
	}
	if resolved != "" {
		query = query.Where("resolved = ?", resolved == "true")
	}

	query.Order("created_at DESC").Find(&alerts)

	utils.Success(c, alerts)
}

func parseQueryInt(c *gin.Context, name string, defaultVal int) int {
	val := c.Query(name)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil || n == 0 {
		return defaultVal
	}
	return n
}
