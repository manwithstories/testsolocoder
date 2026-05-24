package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"temp-staff-platform/database"
	"temp-staff-platform/models"
	"temp-staff-platform/utils"
)

type SalaryHandler struct{}

func NewSalaryHandler() *SalaryHandler {
	return &SalaryHandler{}
}

type CalculateSalaryRequest struct {
	TemporaryID   uuid.UUID `json:"temporary_id" binding:"required"`
	EmployerID    uuid.UUID `json:"employer_id" binding:"required"`
	JobID         uuid.UUID `json:"job_id" binding:"required"`
	PeriodStart   string    `json:"period_start" binding:"required"`
	PeriodEnd     string    `json:"period_end" binding:"required"`
	OvertimeRate  float64   `json:"overtime_rate"`
	Deductions    float64   `json:"deductions"`
	DeductionNote string    `json:"deduction_note"`
}

type BatchPaymentRequest struct {
	SalaryIDs     []uuid.UUID `json:"salary_ids" binding:"required"`
	PaymentMethod string      `json:"payment_method" binding:"required"`
}

func (h *SalaryHandler) CalculateSalary(c *gin.Context) {
	var req CalculateSalaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	periodStart, err := time.Parse("2006-01-02", req.PeriodStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid period_start date format",
		})
		return
	}

	periodEnd, err := time.Parse("2006-01-02", req.PeriodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid period_end date format",
		})
		return
	}

	var job models.JobPosting
	if err := database.DB.First(&job, req.JobID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Job posting not found",
		})
		return
	}

	var checkIns []models.CheckIn
	database.DB.Joins("JOIN schedules ON check_ins.schedule_id = schedules.id").
		Where("check_ins.temporary_id = ? AND check_ins.status = ? AND schedules.job_id = ? AND check_ins.check_in_time >= ? AND check_ins.check_in_time <= ?",
			req.TemporaryID, "checked_out", req.JobID, periodStart, periodEnd.Add(24*time.Hour)).
		Preload("Schedule").
		Find(&checkIns)

	if len(checkIns) == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "No completed check-in records found for this period",
		})
		return
	}

	var totalHours float64
	var baseSalary float64
	var overtimeHours float64
	var overtimePay float64
	var details []models.SalaryDetail

	regularHours := 8.0
	overtimeRate := req.OvertimeRate
	if overtimeRate <= 0 {
		overtimeRate = 1.5
	}

	for _, ci := range checkIns {
		hours := ci.WorkHours
		if hours <= 0 {
			if ci.CheckOutTime != nil {
				hours = ci.CheckOutTime.Sub(ci.CheckInTime).Hours()
			}
		}

		detail := models.SalaryDetail{
			CheckInID:   ci.ID,
			Date:        ci.CheckInTime,
			WorkHours:   hours,
			HourlyRate:  job.SalaryPerHour,
			Type:        "regular",
			Description: fmt.Sprintf("Work on %s", ci.CheckInTime.Format("2006-01-02")),
		}

		if hours > regularHours {
			otHours := hours - regularHours
			overtimeHours += otHours
			otAmount := otHours * job.SalaryPerHour * overtimeRate
			overtimePay += otAmount

			detail.Amount = regularHours * job.SalaryPerHour
			detail.WorkHours = regularHours

			otDetail := models.SalaryDetail{
				CheckInID:   ci.ID,
				Date:        ci.CheckInTime,
				WorkHours:   otHours,
				HourlyRate:  job.SalaryPerHour * overtimeRate,
				Amount:      otAmount,
				Type:        "overtime",
				Description: fmt.Sprintf("Overtime on %s", ci.CheckInTime.Format("2006-01-02")),
			}
			details = append(details, detail, otDetail)
		} else {
			detail.Amount = hours * job.SalaryPerHour
			details = append(details, detail)
		}

		totalHours += hours
		baseSalary += hours * job.SalaryPerHour
	}

	totalSalary := baseSalary + overtimePay - req.Deductions

	salaryRecord := models.SalaryRecord{
		TemporaryID:   req.TemporaryID,
		EmployerID:    req.EmployerID,
		JobID:         req.JobID,
		PeriodStart:   periodStart,
		PeriodEnd:     periodEnd,
		TotalHours:    totalHours,
		BaseSalary:    baseSalary,
		OvertimeHours: overtimeHours,
		OvertimePay:   overtimePay,
		Deductions:    req.Deductions,
		TotalSalary:   totalSalary,
		Status:        "pending",
		Remark:        req.DeductionNote,
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&salaryRecord).Error; err != nil {
			return err
		}

		for i := range details {
			details[i].SalaryID = salaryRecord.ID
			if err := tx.Create(&details[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to calculate salary: " + err.Error(),
		})
		return
	}

	database.DB.Preload("Temporary").Preload("Employer").Preload("JobPosting").First(&salaryRecord, salaryRecord.ID)

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "Salary calculated successfully",
		Data: gin.H{
			"salary":  salaryRecord,
			"details": details,
		},
	})
}

func (h *SalaryHandler) GetSalaryRecords(c *gin.Context) {
	pagination := utils.GetPagination(c)

	var salaries []models.SalaryRecord
	var total int64

	query := database.DB.Model(&models.SalaryRecord{})

	userRole, _ := c.Get("user_role")
	userID, _ := c.Get("user_id")

	if userRole == models.RoleTemporary {
		query = query.Where("temporary_id = ?", userID.(uuid.UUID))
	} else if userRole == models.RoleEmployer {
		query = query.Where("employer_id = ?", userID.(uuid.UUID))
	}

	if tempID := c.Query("temporary_id"); tempID != "" {
		uid, err := uuid.Parse(tempID)
		if err == nil {
			query = query.Where("temporary_id = ?", uid)
		}
	}
	if jobID := c.Query("job_id"); jobID != "" {
		uid, err := uuid.Parse(jobID)
		if err == nil {
			query = query.Where("job_id = ?", uid)
		}
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("period_start >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("period_end <= ?", endDate)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Preload("Temporary").
		Preload("Employer").
		Preload("JobPosting").
		Order("created_at DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&salaries)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: models.PaginatedResponse{
			Data:       salaries,
			Total:      total,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: utils.GetTotalPages(total, pagination.PageSize),
		},
	})
}

func (h *SalaryHandler) GetSalaryRecord(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	salaryID := id.(uuid.UUID)

	var salary models.SalaryRecord
	if err := database.DB.Preload("Temporary").Preload("Employer").Preload("JobPosting").First(&salary, salaryID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Salary record not found",
		})
		return
	}

	var details []models.SalaryDetail
	database.DB.Where("salary_id = ?", salaryID).Preload("CheckIn").Find(&details)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: gin.H{
			"salary":  salary,
			"details": details,
		},
	})
}

func (h *SalaryHandler) PaySalary(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	salaryID := id.(uuid.UUID)

	var salary models.SalaryRecord
	if err := database.DB.First(&salary, salaryID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Salary record not found",
		})
		return
	}

	if salary.Status != "pending" {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Salary has already been paid",
		})
		return
	}

	now := time.Now()
	transactionID := fmt.Sprintf("PAY%s%d", salaryID.String()[:8], now.Unix())

	updates := map[string]interface{}{
		"status":         "paid",
		"payment_method": "bank_transfer",
		"payment_at":     now,
		"transaction_id": transactionID,
	}

	database.DB.Model(&salary).Updates(updates)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Salary paid successfully",
		Data: gin.H{
			"transaction_id": transactionID,
			"payment_at":     now,
		},
	})
}

func (h *SalaryHandler) BatchPaySalary(c *gin.Context) {
	var req BatchPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	if len(req.SalaryIDs) == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "At least one salary ID is required",
		})
		return
	}

	now := time.Now()
	var results []gin.H

	for _, sid := range req.SalaryIDs {
		var salary models.SalaryRecord
		if err := database.DB.First(&salary, sid).Error; err != nil {
			results = append(results, gin.H{"id": sid, "status": "failed", "message": "Not found"})
			continue
		}

		if salary.Status != "pending" {
			results = append(results, gin.H{"id": sid, "status": "skipped", "message": "Already paid"})
			continue
		}

		transactionID := fmt.Sprintf("PAY%s%d", sid.String()[:8], now.Unix())

		database.DB.Model(&salary).Updates(map[string]interface{}{
			"status":         "paid",
			"payment_method": req.PaymentMethod,
			"payment_at":     now,
			"transaction_id": transactionID,
		})

		results = append(results, gin.H{"id": sid, "status": "success", "transaction_id": transactionID})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: fmt.Sprintf("Processed %d salary records", len(req.SalaryIDs)),
		Data:    results,
	})
}

func (h *SalaryHandler) ExportSalary(c *gin.Context) {
	id, _ := c.Get("id_uuid")
	salaryID := id.(uuid.UUID)

	var salary models.SalaryRecord
	if err := database.DB.Preload("Temporary").Preload("Employer").Preload("JobPosting").First(&salary, salaryID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    404,
			Message: "Salary record not found",
		})
		return
	}

	var details []models.SalaryDetail
	database.DB.Where("salary_id = ?", salaryID).Find(&details)

	filename := fmt.Sprintf("salary_%s.pdf", salaryID.String()[:8])
	filepath, err := utils.ExportSalaryToPDF(&salary, details, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Failed to export salary: " + err.Error(),
		})
		return
	}

	utils.ServeFile(c, filepath, filename)
}
