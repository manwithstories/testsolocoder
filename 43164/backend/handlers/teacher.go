package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tutoring-platform/database"
	"tutoring-platform/models"
)

type TeacherProfileRequest struct {
	Bio            string  `json:"bio"`
	Education      string  `json:"education"`
	Experience     string  `json:"experience"`
	Certifications string  `json:"certifications"`
	HourlyRate     float64 `json:"hourlyRate" binding:"required,min=0"`
	Currency       string  `json:"currency"`
	ResumeURL      string  `json:"resumeUrl"`
	IDCardURL      string  `json:"idCardUrl"`
}

type TeacherSubjectRequest struct {
	SubjectID  uuid.UUID `json:"subjectId" binding:"required"`
	Level      string    `json:"level" binding:"required"`
	CustomRate float64   `json:"customRate"`
	IsPrimary  bool      `json:"isPrimary"`
}

type AvailabilitySlotRequest struct {
	DayOfWeek   int        `json:"dayOfWeek" binding:"required,min=0,max=6"`
	StartTime   string     `json:"startTime" binding:"required"`
	EndTime     string     `json:"endTime" binding:"required"`
	IsRecurring bool       `json:"isRecurring"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
}

func CreateTeacherProfile(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req TeacherProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.TeacherProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	if profile.ApprovalStatus == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile already approved, cannot modify"})
		return
	}

	updates := map[string]interface{}{
		"bio":             req.Bio,
		"education":       req.Education,
		"experience":      req.Experience,
		"certifications":  req.Certifications,
		"hourly_rate":     req.HourlyRate,
		"currency":        req.Currency,
		"resume_url":      req.ResumeURL,
		"id_card_url":     req.IDCardURL,
		"approval_status": "pending",
	}

	if err := database.DB.Model(&profile).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully, pending approval", "profile": profile})
}

func GetTeacherProfile(c *gin.Context) {
	userID, _ := c.Get("userId")

	var profile models.TeacherProfile
	if err := database.DB.Preload("Subjects.Subject").Preload("Availabilities").Preload("User").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func GetTeacherByID(c *gin.Context) {
	id := c.Param("id")
	teacherID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid teacher ID"})
		return
	}

	var profile models.TeacherProfile
	if err := database.DB.Preload("Subjects.Subject").Preload("Availabilities").Preload("User").Where("user_id = ? AND approval_status = ?", teacherID, "approved").First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found or not approved"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func AddTeacherSubject(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req TeacherSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.TeacherProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	var subject models.Subject
	if err := database.DB.Where("id = ?", req.SubjectID).First(&subject).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}

	var existing models.TeacherSubject
	if database.DB.Where("teacher_id = ? AND subject_id = ?", profile.ID, req.SubjectID).First(&existing).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Subject already added"})
		return
	}

	teacherSubject := models.TeacherSubject{
		TeacherID:  profile.ID,
		SubjectID:  req.SubjectID,
		Level:      req.Level,
		CustomRate: req.CustomRate,
		IsPrimary:  req.IsPrimary,
	}

	if err := database.DB.Create(&teacherSubject).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add subject"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Subject added successfully", "subject": teacherSubject})
}

func RemoveTeacherSubject(c *gin.Context) {
	userID, _ := c.Get("userId")
	subjectID := c.Param("subjectId")

	var profile models.TeacherProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	if err := database.DB.Where("teacher_id = ? AND subject_id = ?", profile.ID, subjectID).Delete(&models.TeacherSubject{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove subject"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subject removed successfully"})
}

func AddAvailabilitySlot(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req AvailabilitySlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.TeacherProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	slot := models.AvailabilitySlot{
		TeacherID:   profile.ID,
		DayOfWeek:   req.DayOfWeek,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsRecurring: req.IsRecurring,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := database.DB.Create(&slot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add availability slot"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Availability slot added successfully", "slot": slot})
}

func RemoveAvailabilitySlot(c *gin.Context) {
	slotID := c.Param("slotId")
	userID, _ := c.Get("userId")

	var profile models.TeacherProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	if err := database.DB.Where("id = ? AND teacher_id = ?", slotID, profile.ID).Delete(&models.AvailabilitySlot{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove availability slot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Availability slot removed successfully"})
}

func ApproveTeacherProfile(c *gin.Context) {
	teacherID := c.Param("id")
	adminID, _ := c.Get("userId")

	var profile models.TeacherProfile
	if err := database.DB.Where("user_id = ?", teacherID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"approval_status": "approved",
		"is_verified":     true,
		"approved_at":     &now,
		"approved_by":     adminID,
	}

	if err := database.DB.Model(&profile).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Teacher profile approved successfully"})
}

func RejectTeacherProfile(c *gin.Context) {
	teacherID := c.Param("id")
	adminID, _ := c.Get("userId")

	var req struct {
		Notes string `json:"notes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.TeacherProfile
	if err := database.DB.Where("user_id = ?", teacherID).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher profile not found"})
		return
	}

	updates := map[string]interface{}{
		"approval_status": "rejected",
		"approval_notes":  req.Notes,
		"approved_by":     adminID,
	}

	if err := database.DB.Model(&profile).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Teacher profile rejected"})
}

func GetPendingTeacherApprovals(c *gin.Context) {
	var profiles []models.TeacherProfile
	if err := database.DB.Preload("User").Where("approval_status = ?", "pending").Find(&profiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending approvals"})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

func ListTeachers(c *gin.Context) {
	subject := c.Query("subject")
	minRate := c.Query("minRate")
	maxRate := c.Query("maxRate")
	minRating := c.Query("minRating")

	var teachers []models.TeacherProfile
	query := database.DB.Preload("Subjects.Subject").Preload("User").Where("approval_status = ?", "approved")

	if subject != "" {
		query = query.Joins("JOIN teacher_subjects ON teacher_subjects.teacher_id = teacher_profiles.id").
			Joins("JOIN subjects ON subjects.id = teacher_subjects.subject_id").
			Where("subjects.name ILIKE ? OR subjects.id = ?", "%"+subject+"%", subject)
	}

	if minRate != "" {
		query = query.Where("hourly_rate >= ?", minRate)
	}
	if maxRate != "" {
		query = query.Where("hourly_rate <= ?", maxRate)
	}
	if minRating != "" {
		query = query.Where("rating >= ?", minRating)
	}

	if err := query.Find(&teachers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teachers"})
		return
	}

	c.JSON(http.StatusOK, teachers)
}
