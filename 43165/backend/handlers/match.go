package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"temp-staff-platform/database"
	"temp-staff-platform/models"
	"temp-staff-platform/utils"
)

type MatchHandler struct{}

func NewMatchHandler() *MatchHandler {
	return &MatchHandler{}
}

type MatchTemporariesRequest struct {
	JobID      uuid.UUID `json:"job_id" binding:"required"`
	SkillTags  []string  `json:"skill_tags"`
	MinCredit  int       `json:"min_credit"`
	MaxDistance float64  `json:"max_distance"`
	Limit      int       `json:"limit"`
}

type MatchResult struct {
	Temporary    models.User `json:"temporary"`
	MatchScore   float64     `json:"match_score"`
	MatchReasons []string    `json:"match_reasons"`
}

func (h *MatchHandler) MatchTemporaries(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	if userRole != models.RoleAgent && userRole != models.RoleEmployer {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    403,
			Message: "Only agents or employers can perform matching",
		})
		return
	}

	var req MatchTemporariesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
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

	if job.EmployerID != userID.(uuid.UUID) && userRole == models.RoleEmployer {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    403,
			Message: "You can only match for your own job postings",
		})
		return
	}

	var alreadyApplied []uuid.UUID
	database.DB.Model(&models.JobApplication{}).
		Where("job_id = ?", req.JobID).
		Pluck("temporary_id", &alreadyApplied)

	var alreadyScheduled []uuid.UUID
	database.DB.Model(&models.Schedule{}).
		Where("job_id = ?", req.JobID).
		Pluck("temporary_id", &alreadyScheduled)

	excludeIDs := append(alreadyApplied, alreadyScheduled...)

	var temporaries []models.User
	query := database.DB.Model(&models.User{}).
		Where("role = ? AND status = ?", models.RoleTemporary, models.UserStatusActive)

	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN ?", excludeIDs)
	}

	if req.MinCredit > 0 {
		query = query.Where("credit_score >= ?", req.MinCredit)
	}

	query = query.Order("credit_score DESC")

	limit := req.Limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	query.Limit(limit).Find(&temporaries)

	var matchResults []MatchResult
	for _, temp := range temporaries {
		score := 0.0
		var reasons []string

		if temp.CreditScore >= 80 {
			score += 30
			reasons = append(reasons, "高信用分")
		} else if temp.CreditScore >= 60 {
			score += 20
			reasons = append(reasons, "良好信用分")
		} else {
			score += 10
			reasons = append(reasons, "信用分一般")
		}

		if job.ActivityType != "" {
			score += 20
			reasons = append(reasons, "活动类型匹配")
		}

		if job.Location != "" {
			score += 20
			reasons = append(reasons, "地点匹配")
		}

		if job.SalaryPerHour >= 50 {
			score += 15
			reasons = append(reasons, "高薪吸引")
		} else if job.SalaryPerHour >= 30 {
			score += 10
			reasons = append(reasons, "薪资合理")
		}

		if job.IsUrgent {
			score += 15
			reasons = append(reasons, "紧急岗位优先")
		}

		matchResults = append(matchResults, MatchResult{
			Temporary:    temp,
			MatchScore:   score,
			MatchReasons: reasons,
		})
	}

	for i := range matchResults {
		for j := i + 1; j < len(matchResults); j++ {
			if matchResults[i].MatchScore < matchResults[j].MatchScore {
				matchResults[i], matchResults[j] = matchResults[j], matchResults[i]
			}
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Matching completed successfully",
		Data: gin.H{
			"total_matches": len(matchResults),
			"matches":       matchResults,
		},
	})
}

func (h *MatchHandler) QuickAssign(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	if userRole != models.RoleAgent && userRole != models.RoleEmployer {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    403,
			Message: "Only agents or employers can perform assignment",
		})
		return
	}

	var req struct {
		JobID        uuid.UUID   `json:"job_id" binding:"required"`
		TemporaryIDs []uuid.UUID `json:"temporary_ids" binding:"required"`
		ShiftDate    string      `json:"shift_date" binding:"required"`
		StartTime    string      `json:"start_time" binding:"required"`
		EndTime      string      `json:"end_time" binding:"required"`
		Location     string      `json:"location"`
		Notes        string      `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Invalid request parameters: " + err.Error(),
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

	shiftDate, _ := time.Parse("2006-01-02", req.ShiftDate)

	var createdSchedules []models.Schedule
	var failedAssignments []gin.H

	for _, tempID := range req.TemporaryIDs {
		var existing models.Schedule
		if err := database.DB.Where("temporary_id = ? AND shift_date = ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))",
			tempID, shiftDate, req.EndTime, req.StartTime, req.EndTime, req.StartTime).First(&existing).Error; err == nil {
			failedAssignments = append(failedAssignments, gin.H{
				"temporary_id": tempID,
				"reason":       "Schedule conflict",
			})
			continue
		}

		var temp models.User
		database.DB.First(&temp, tempID)

		var existingApp models.JobApplication
		appExists := false
		if err := database.DB.Where("job_id = ? AND temporary_id = ?", req.JobID, tempID).First(&existingApp).Error; err == nil {
			appExists = true
		}

		if !appExists {
			agentID := userID.(uuid.UUID)
			application := models.JobApplication{
				JobID:       req.JobID,
				TemporaryID: tempID,
				AgentID:     &agentID,
				Status:      "approved",
				AppliedAt:   time.Now(),
				ReviewedAt:  &[]time.Time{time.Now()}[0],
			}
			database.DB.Create(&application)
		} else if existingApp.Status == "pending" {
			now := time.Now()
			database.DB.Model(&existingApp).Updates(map[string]interface{}{
				"status":      "approved",
				"reviewed_at": now,
				"agent_id":    userID.(uuid.UUID),
			})
		}

		schedule := models.Schedule{
			JobID:       req.JobID,
			TemporaryID: tempID,
			ShiftDate:   shiftDate,
			StartTime:   req.StartTime,
			EndTime:     req.EndTime,
			Location:    req.Location,
			Notes:       req.Notes,
			Status:      "scheduled",
			CreatedBy:   userID.(uuid.UUID),
		}

		if err := database.DB.Create(&schedule).Error; err == nil {
			database.DB.Preload("Temporary").First(&schedule, schedule.ID)
			createdSchedules = append(createdSchedules, schedule)
		} else {
			failedAssignments = append(failedAssignments, gin.H{
				"temporary_id": tempID,
				"reason":       err.Error(),
			})
		}
	}

	database.DB.Model(&job).Update("hired_count", job.HiredCount+len(createdSchedules))

	c.JSON(http.StatusCreated, models.Response{
		Code:    201,
		Message: "Assignment completed",
		Data: gin.H{
			"created_count":   len(createdSchedules),
			"schedules":       createdSchedules,
			"failed_count":    len(failedAssignments),
			"failed_details":  failedAssignments,
		},
	})
}

func (h *MatchHandler) GetMatchHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	pagination := utils.GetPagination(c)

	var applications []models.JobApplication
	var total int64

	query := database.DB.Model(&models.JobApplication{}).Where("agent_id = ?", userID.(uuid.UUID))

	if jobID := c.Query("job_id"); jobID != "" {
		uid, err := uuid.Parse(jobID)
		if err == nil {
			query = query.Where("job_id = ?", uid)
		}
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Preload("Temporary").
		Preload("JobPosting").
		Order("applied_at DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&applications)

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "Success",
		Data: models.PaginatedResponse{
			Data:       applications,
			Total:      total,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			TotalPages: utils.GetTotalPages(total, pagination.PageSize),
		},
	})
}
