package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"ship-rental-platform/internal/database"
	"ship-rental-platform/internal/model"
	"ship-rental-platform/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VoyageHandler struct{}

func NewVoyageHandler() *VoyageHandler {
	return &VoyageHandler{}
}

func (h *VoyageHandler) CreateVoyageLog(c *gin.Context) {
	userID, _ := c.Get("user_id")
	captainID, _ := uuid.Parse(userID.(string))

	var req model.CreateVoyageLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	rentalID, _ := uuid.Parse(req.RentalID)

	var rental model.Rental
	if err := database.DB.First(&rental, "id = ?", rentalID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Rental not found")
		return
	}

	if rental.Status != model.RentalStatusActive && rental.Status != model.RentalStatusCompleted {
		utils.Error(c, http.StatusBadRequest, "Can only create voyage logs for active or completed rentals")
		return
	}

	log := model.VoyageLog{
		RentalID:         rentalID,
		ShipID:           rental.ShipID,
		CaptainID:        captainID,
		LogDate:          req.LogDate,
		DeparturePort:    req.DeparturePort,
		ArrivalPort:      req.ArrivalPort,
		DepartureTime:    req.DepartureTime,
		ArrivalTime:      req.ArrivalTime,
		Distance:         req.Distance,
		EngineHours:      req.EngineHours,
		FuelConsumed:     req.FuelConsumed,
		FuelUnit:         req.FuelUnit,
		AvgSpeed:         req.AvgSpeed,
		MaxSpeed:         req.MaxSpeed,
		Weather:          req.Weather,
		WindSpeed:        req.WindSpeed,
		WindDirection:    req.WindDirection,
		WaveHeight:       req.WaveHeight,
		WaterTemperature: req.WaterTemperature,
		AirTemperature:   req.AirTemperature,
		PassengerCount:   req.PassengerCount,
		CrewCount:        req.CrewCount,
		Notes:            req.Notes,
		Incidents:        req.Incidents,
	}

	if err := database.DB.Create(&log).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create voyage log")
		return
	}

	utils.Created(c, log)
}

func (h *VoyageHandler) GetVoyageLogs(c *gin.Context) {
	rentalID := c.Query("rental_id")
	shipID := c.Query("ship_id")

	var logs []model.VoyageLog
	query := database.DB.Preload("Ship").Preload("Captain")

	if rentalID != "" {
		query = query.Where("rental_id = ?", rentalID)
	}
	if shipID != "" {
		query = query.Where("ship_id = ?", shipID)
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate != "" && endDate != "" {
		query = query.Where("log_date BETWEEN ? AND ?", startDate, endDate)
	}

	var total int64
	query.Model(&model.VoyageLog{}).Count(&total)

	page := utils.ParseInt(c.DefaultQuery("page", "1"), 1)
	pageSize := utils.ParseInt(c.DefaultQuery("page_size", "20"), 20)
	offset := (page - 1) * pageSize

	if err := query.Order("log_date DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch voyage logs")
		return
	}

	utils.Paginated(c, logs, total, page, pageSize)
}

func (h *VoyageHandler) GetVoyageLog(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid voyage log ID")
		return
	}

	var log model.VoyageLog
	if err := database.DB.Preload("Ship").Preload("Captain").Preload("Rental").First(&log, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Voyage log not found")
		return
	}

	utils.Success(c, log)
}

func (h *VoyageHandler) UpdateVoyageLog(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid voyage log ID")
		return
	}

	var log model.VoyageLog
	if err := database.DB.First(&log, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Voyage log not found")
		return
	}

	var req model.CreateVoyageLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.DeparturePort != "" {
		log.DeparturePort = req.DeparturePort
	}
	if req.ArrivalPort != "" {
		log.ArrivalPort = req.ArrivalPort
	}
	if req.DepartureTime != nil {
		log.DepartureTime = req.DepartureTime
	}
	if req.ArrivalTime != nil {
		log.ArrivalTime = req.ArrivalTime
	}
	if req.Distance != 0 {
		log.Distance = req.Distance
	}
	if req.EngineHours != 0 {
		log.EngineHours = req.EngineHours
	}
	if req.FuelConsumed != 0 {
		log.FuelConsumed = req.FuelConsumed
	}
	if req.Weather != "" {
		log.Weather = req.Weather
	}
	if req.WindSpeed != 0 {
		log.WindSpeed = req.WindSpeed
	}
	if req.Notes != "" {
		log.Notes = req.Notes
	}
	if req.Incidents != "" {
		log.Incidents = req.Incidents
	}

	if err := database.DB.Save(&log).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update voyage log")
		return
	}

	utils.Success(c, log)
}

func (h *VoyageHandler) DeleteVoyageLog(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid voyage log ID")
		return
	}

	if err := database.DB.Delete(&model.VoyageLog{}, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete voyage log")
		return
	}

	utils.Success(c, gin.H{"message": "Voyage log deleted successfully"})
}

func (h *VoyageHandler) ExportVoyageLogs(c *gin.Context) {
	var req model.ExportVoyageLogRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	rentalID, _ := uuid.Parse(req.RentalID)

	var logs []model.VoyageLog
	query := database.DB.Where("rental_id = ?", rentalID)

	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("log_date BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	if err := query.Order("log_date ASC").Find(&logs).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch voyage logs")
		return
	}

	if req.Format == "csv" {
		csvData := generateVoyageCSV(logs)
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=voyage_logs_%s.csv", time.Now().Format("20060102")))
		c.String(http.StatusOK, csvData)
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=voyage_logs_%s.pdf", time.Now().Format("20060102")))
	c.String(http.StatusOK, "PDF export placeholder - install unipdf for full PDF support")
}

func generateVoyageCSV(logs []model.VoyageLog) string {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{
		"Date", "Departure Port", "Arrival Port", "Departure Time", "Arrival Time",
		"Distance (nm)", "Engine Hours", "Fuel Consumed", "Weather",
		"Wind Speed", "Wave Height", "Air Temp", "Water Temp", "Notes",
	}
	writer.Write(headers)

	for _, log := range logs {
		row := []string{
			log.LogDate.Format("2006-01-02"),
			log.DeparturePort,
			log.ArrivalPort,
			formatTime(log.DepartureTime),
			formatTime(log.ArrivalTime),
			fmt.Sprintf("%.2f", log.Distance),
			fmt.Sprintf("%.1f", log.EngineHours),
			fmt.Sprintf("%.2f %s", log.FuelConsumed, log.FuelUnit),
			string(log.Weather),
			fmt.Sprintf("%.1f", log.WindSpeed),
			fmt.Sprintf("%.1f", log.WaveHeight),
			fmt.Sprintf("%.1f", log.AirTemperature),
			fmt.Sprintf("%.1f", log.WaterTemperature),
			log.Notes,
		}
		writer.Write(row)
	}

	writer.Flush()
	return buf.String()
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("15:04")
}
