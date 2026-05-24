package handlers

import (
	"os"

	"recruitment-platform/config"
	"recruitment-platform/middleware"
	"recruitment-platform/models"
	"recruitment-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResumeHandler struct {
	DB *gorm.DB
}

func NewResumeHandler(db *gorm.DB) *ResumeHandler {
	return &ResumeHandler{DB: db}
}

type ResumeRequest struct {
	Title     string               `json:"title" binding:"required"`
	FullName  string               `json:"full_name" binding:"required"`
	Email     string               `json:"email"`
	Phone     string               `json:"phone"`
	Location  string               `json:"location"`
	Summary   string               `json:"summary"`
	IsDefault bool                 `json:"is_default"`
	Education []EducationRequest   `json:"education_list"`
	Work      []WorkExperienceReq  `json:"work_experiences"`
	Skills    []SkillRequest       `json:"skills"`
}

type EducationRequest struct {
	School      string `json:"school" binding:"required"`
	Degree      string `json:"degree"`
	Major       string `json:"major"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
}

type WorkExperienceReq struct {
	Company     string `json:"company" binding:"required"`
	Position    string `json:"position" binding:"required"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
}

type SkillRequest struct {
	Name  string `json:"name" binding:"required"`
	Level string `json:"level"`
}

func (h *ResumeHandler) CreateResume(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	var req ResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	tx := h.DB.Begin()

	resume := models.Resume{
		JobSeekerID: jobSeeker.ID,
		Title:       req.Title,
		FullName:    req.FullName,
		Email:       req.Email,
		Phone:       req.Phone,
		Location:    req.Location,
		Summary:     req.Summary,
		IsDefault:   req.IsDefault,
	}

	if req.IsDefault {
		tx.Model(&models.Resume{}).Where("jobseeker_id = ?", jobSeeker.ID).Update("is_default", false)
	}

	if err := tx.Create(&resume).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to create resume")
		return
	}

	for _, edu := range req.Education {
		e := models.Education{
			ResumeID:    resume.ID,
			School:      edu.School,
			Degree:      edu.Degree,
			Major:       edu.Major,
			StartDate:   edu.StartDate,
			EndDate:     edu.EndDate,
			Description: edu.Description,
		}
		tx.Create(&e)
	}

	for _, work := range req.Work {
		w := models.WorkExperience{
			ResumeID:    resume.ID,
			Company:     work.Company,
			Position:    work.Position,
			StartDate:   work.StartDate,
			EndDate:     work.EndDate,
			Description: work.Description,
		}
		tx.Create(&w)
	}

	for _, skill := range req.Skills {
		s := models.Skill{
			ResumeID: resume.ID,
			Name:     skill.Name,
			Level:    skill.Level,
		}
		tx.Create(&s)
	}

	tx.Commit()

	h.DB.Preload("EducationList").Preload("WorkExperiences").Preload("Skills").First(&resume, resume.ID)
	utils.Success(c, resume)
}

func (h *ResumeHandler) UpdateResume(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	var resume models.Resume
	if err := h.DB.Where("id = ? AND jobseeker_id = ?", id, jobSeeker.ID).First(&resume).Error; err != nil {
		utils.NotFound(c, "Resume not found")
		return
	}

	var req ResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	tx := h.DB.Begin()

	if req.IsDefault {
		tx.Model(&models.Resume{}).Where("jobseeker_id = ? AND id != ?", jobSeeker.ID, resume.ID).Update("is_default", false)
	}

	updates := map[string]interface{}{
		"title":     req.Title,
		"full_name": req.FullName,
		"email":     req.Email,
		"phone":     req.Phone,
		"location":  req.Location,
		"summary":   req.Summary,
		"is_default": req.IsDefault,
	}

	if err := tx.Model(&resume).Updates(updates).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to update resume")
		return
	}

	tx.Where("resume_id = ?", resume.ID).Delete(&models.Education{})
	tx.Where("resume_id = ?", resume.ID).Delete(&models.WorkExperience{})
	tx.Where("resume_id = ?", resume.ID).Delete(&models.Skill{})

	for _, edu := range req.Education {
		e := models.Education{
			ResumeID:    resume.ID,
			School:      edu.School,
			Degree:      edu.Degree,
			Major:       edu.Major,
			StartDate:   edu.StartDate,
			EndDate:     edu.EndDate,
			Description: edu.Description,
		}
		tx.Create(&e)
	}

	for _, work := range req.Work {
		w := models.WorkExperience{
			ResumeID:    resume.ID,
			Company:     work.Company,
			Position:    work.Position,
			StartDate:   work.StartDate,
			EndDate:     work.EndDate,
			Description: work.Description,
		}
		tx.Create(&w)
	}

	for _, skill := range req.Skills {
		s := models.Skill{
			ResumeID: resume.ID,
			Name:     skill.Name,
			Level:    skill.Level,
		}
		tx.Create(&s)
	}

	tx.Commit()

	h.DB.Preload("EducationList").Preload("WorkExperiences").Preload("Skills").First(&resume, resume.ID)
	utils.Success(c, resume)
}

func (h *ResumeHandler) DeleteResume(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	result := h.DB.Where("id = ? AND jobseeker_id = ?", id, jobSeeker.ID).Delete(&models.Resume{})
	if result.RowsAffected == 0 {
		utils.NotFound(c, "Resume not found")
		return
	}

	utils.Success(c, nil)
}

func (h *ResumeHandler) GetResume(c *gin.Context) {
	id := c.Param("id")

	var resume models.Resume
	if err := h.DB.Preload("EducationList").Preload("WorkExperiences").Preload("Skills").
		First(&resume, id).Error; err != nil {
		utils.NotFound(c, "Resume not found")
		return
	}

	utils.Success(c, resume)
}

func (h *ResumeHandler) ListResumes(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	var resumes []models.Resume
	h.DB.Preload("EducationList").Preload("WorkExperiences").Preload("Skills").
		Where("jobseeker_id = ?", jobSeeker.ID).Order("is_default desc, created_at desc").Find(&resumes)

	utils.Success(c, resumes)
}

func (h *ResumeHandler) UploadResumeFile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	var resume models.Resume
	if err := h.DB.Where("id = ? AND jobseeker_id = ?", id, jobSeeker.ID).First(&resume).Error; err != nil {
		utils.NotFound(c, "Resume not found")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "No file uploaded")
		return
	}

	if file.Size > 10*1024*1024 {
		utils.BadRequest(c, "File size exceeds 10MB limit")
		return
	}

	ext := ""
	if len(file.Filename) > 4 {
		ext = file.Filename[len(file.Filename)-4:]
	}
	if ext != ".pdf" && ext != ".PDF" {
		utils.BadRequest(c, "Only PDF files are allowed")
		return
	}

	uploadDir := config.AppConfig.UploadDir + "/resumes"
	os.MkdirAll(uploadDir, 0755)

	filePath, fileName, err := utils.SaveUploadedFile(c, file, uploadDir)
	if err != nil {
		utils.InternalError(c, "Failed to save file")
		return
	}

	h.DB.Model(&resume).Updates(map[string]interface{}{
		"file_path": filePath,
		"file_name": fileName,
	})

	utils.Success(c, gin.H{"file_path": filePath, "file_name": fileName})
}

func (h *ResumeHandler) SetDefaultResume(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	tx := h.DB.Begin()
	tx.Model(&models.Resume{}).Where("jobseeker_id = ?", jobSeeker.ID).Update("is_default", false)
	tx.Model(&models.Resume{}).Where("id = ? AND jobseeker_id = ?", id, jobSeeker.ID).Update("is_default", true)
	tx.Commit()

	utils.Success(c, nil)
}
