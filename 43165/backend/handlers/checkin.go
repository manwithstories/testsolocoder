package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"temp-staff-platform/database"
	"temp-staff-platform/models"
	"temp-staff-platform/utils"
)

type CheckInHandler struct{}

func NewCheckInHandler() *CheckInHandler {
	return &CheckInHandler{}
}

type CheckInRequest struct {
	ScheduleID uuid.UUID `json:"schedule_id" binding:"required"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Location   string    `json:"location"`
	CheckInType string   `json:"check_in_type" binding:"required,oneof=qr location face"`
	QRCode     string    `json:"qr_code"`
}

type CheckOutRequest struct {
	CheckInID uuid.UUID `json:"check_in_id" binding:"required"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

type FaceVerifyRequest struct {
	CheckInID  uuid.UUID `json:"check_in_id" binding:"required"`
	FaceData   string    `json:"face_data" binding:"required"`
}

func (h *CheckInHandler) GenerateQRCode(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	scheduleID := id.(uuid.UUID)

	var schedule models.Schedule
	if err := database.DB.First(&schedule, scheduleID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Schedule not found",
		})
		return
	}

	hash := sha256.New()
	hash.Write([]byte(schedule.ID.String() + schedule.TemporaryID.String() + time.Now().String()))
	qrCode := hex.EncodeToString(hash.Sum(nil))[:32]

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "QR code generated successfully",
		Data: gin.H{
			"qr_code": qrCode,
			"schedule_id": scheduleID,
		},
	})
}

func (h *CheckInHandler) CheckIn(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	var schedule models.Schedule
	if err := database.DB.First(&schedule, req.ScheduleID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Schedule not found",
		})
		return
	}

	if schedule.TemporaryID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    403,
			Message: "This schedule does not belong to you",
		})
		return
	}

	if schedule.Status == "completed" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "This shift has already been completed",
		})
		return
	}

	var existingCheckIn models.CheckIn
	if err := database.DB.Where("schedule_id = ? AND status = ?", req.ScheduleID, "checked_in").First(&existingCheckIn).Error; err == nil {
		c.JSON(http.StatusConflict, models.Response{
			Code:    409,
			Message: "You have already checked in for this shift",
		})
		return
	}

	var faceVerified bool
	if req.CheckInType == "face" {
		faceVerified = true
	}

	checkIn := models.CheckIn{
		ScheduleID:   req.ScheduleID,
		TemporaryID:  userID.(uuid.UUID),
		CheckInType:  req.CheckInType,
		CheckInTime:  time.Now(),
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Location:     req.Location,
		FaceVerified: faceVerified,
		QRCode:       req.QRCode,
		Status:       "checked_in",
	}

	if err := database.DB.Create(&checkIn).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to check in: " + err.Error(),
		})
		return
	}

	database.DB.Model(&schedule).Update("status", "in_progress")

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "Check-in successful",
		Data:    checkIn,
	})
}

func (h *CheckInHandler) CheckOut(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CheckOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	var checkIn models.CheckIn
	if err := database.DB.First(&checkIn, req.CheckInID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Check-in record not found",
		})
		return
	}

	if checkIn.TemporaryID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    403,
			Message: "This check-in record does not belong to you",
		})
		return
	}

	if checkIn.Status != "checked_in" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid check-in status for checkout",
		})
		return
	}

	now := time.Now()
	workHours := now.Sub(checkIn.CheckInTime).Hours()

	updates := map[string]interface{}{
		"check_out_time": now,
		"latitude":       req.Latitude,
		"longitude":      req.Longitude,
		"status":         "checked_out",
		"work_hours":     workHours,
	}

	database.DB.Model(&checkIn).Updates(updates)

	var schedule models.Schedule
	database.DB.First(&schedule, checkIn.ScheduleID)
	database.DB.Model(&schedule).Update("status", "completed")

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Check-out successful",
		Data: gin.H{
			"check_in": checkIn,
			"work_hours": workHours,
		},
	})
}

func (h *CheckInHandler) VerifyFace(c *gin.Context) {
	var req FaceVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	var checkIn models.CheckIn
	if err := database.DB.First(&checkIn, req.CheckInID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Check-in record not found",
		})
		return
	}

	var tempUser models.User
	database.DB.First(&tempUser, checkIn.TemporaryID)

	if tempUser.FaceData == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Face data not registered. Please register face data first.",
		})
		return
	}

	verified := tempUser.FaceData == req.FaceData
	if !verified {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    401,
			Message: "Face verification failed",
		})
		return
	}

	database.DB.Model(&checkIn).Update("face_verified", true)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Face verification successful",
	})
}

func (h *CheckInHandler) GetCheckInRecords(c *gin.Context) {
	pagination := utils.GetPagination(c)

	var checkIns []models.CheckIn
	var total int64

	query := database.DB.Model(&models.CheckIn{})

	userRole, _ := c.Get("user_role")
	userID, _ := c.Get("user_id")

	if userRole == models.RoleTemporary {
		query = query.Where("temporary_id = ?", userID.(uuid.UUID))
	}

	if scheduleID := c.Query("schedule_id"); scheduleID != "" {
		uid, err := uuid.Parse(scheduleID)
		if err == nil {
			query = query.Where("schedule_id = ?", uid)
		}
	}
	if tempID := c.Query("temporary_id"); tempID != "" {
		uid, err := uuid.Parse(tempID)
		if err == nil {
			query = query.Where("temporary_id = ?", uid)
		}
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("check_in_time >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("check_in_time <= ?", endDate)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Preload("Temporary").
		Preload("Schedule").
		Preload("Schedule.JobPosting").
		Order("check_in_time DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&checkIns)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: models.PaginatedResponse{
			Data:       checkIns,
			Total:      total,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: utils.GetTotalPages(total, pagination.PageSize),
		},
	})
}

func (h *CheckInHandler) GetCheckInStats(c *gin.Context) {
	userRole, _ := c.Get("user_role")
	userID, _ := c.Get("user_id")

	query := database.DB.Model(&models.CheckIn{})

	if userRole == models.RoleTemporary {
		query = query.Where("temporary_id = ?", userID.(uuid.UUID))
	}

	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("check_in_time >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("check_in_time <= ?", endDate)
	}

	var totalCount int64
	query.Count(&totalCount)

	var checkedInCount int64
	database.DB.Model(&models.CheckIn{}).Where("status = ?", "checked_in").Count(&checkedInCount)

	var checkedOutCount int64
	database.DB.Model(&models.CheckIn{}).Where("status = ?", "checked_out").Count(&checkedOutCount)

	var totalHours float64
	database.DB.Model(&models.CheckIn{}).Where("status = ?", "checked_out").Select("COALESCE(SUM(work_hours), 0)").Scan(&totalHours)

	var absentCount int64
	database.DB.Raw(`
		SELECT COUNT(*) FROM schedules s
		LEFT JOIN check_ins ci ON ci.schedule_id = s.id
		WHERE ci.id IS NULL AND s.shift_date < CURRENT_DATE
	`).Scan(&absentCount)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: gin.H{
			"total_check_ins":    totalCount,
			"currently_checked_in": checkedInCount,
			"completed_check_outs": checkedOutCount,
			"total_work_hours":    fmt.Sprintf("%.2f", totalHours),
			"absent_count":        absentCount,
		},
	})
}

func (h *CheckInHandler) RegisterFaceData(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		FaceData string `json:"face_data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters",
		})
		return
	}

	database.DB.Model(&models.User{}).Where("id = ?", userID.(uuid.UUID)).Update("face_data", req.FaceData)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Face data registered successfully",
	})
}
