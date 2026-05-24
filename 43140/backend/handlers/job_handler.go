package handlers

import (
	"strconv"

	"recruitment-platform/middleware"
	"recruitment-platform/models"
	"recruitment-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JobHandler struct {
	DB *gorm.DB
}

func NewJobHandler(db *gorm.DB) *JobHandler {
	return &JobHandler{DB: db}
}

type JobRequest struct {
	Title          string         `json:"title" binding:"required"`
	Location       string         `json:"location" binding:"required"`
	SalaryMin      float64        `json:"salary_min"`
	SalaryMax      float64        `json:"salary_max"`
	SalaryType     string         `json:"salary_type"`
	Description    string         `json:"description" binding:"required"`
	Requirements   string         `json:"requirements"`
	Skills         string         `json:"skills"`
	JobType        models.JobType `json:"job_type" binding:"required,oneof=full-time part-time contract internship remote"`
	DepartmentID   *uint          `json:"department_id"`
	Status         models.JobStatus `json:"status"`
}

func (h *JobHandler) CreateJob(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var req JobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	job := models.Job{
		CompanyID:    company.ID,
		Title:        req.Title,
		Location:     req.Location,
		SalaryMin:    req.SalaryMin,
		SalaryMax:    req.SalaryMax,
		SalaryType:   req.SalaryType,
		Description:  req.Description,
		Requirements: req.Requirements,
		Skills:       req.Skills,
		JobType:      req.JobType,
		DepartmentID: req.DepartmentID,
		Status:       models.JobStatusOpen,
	}

	if req.Status != "" {
		job.Status = req.Status
	}

	if err := h.DB.Create(&job).Error; err != nil {
		utils.InternalError(c, "Failed to create job")
		return
	}

	utils.Success(c, job)
}

func (h *JobHandler) UpdateJob(c *gin.Context) {
	userID := middleware.GetUserID(c)
	jobID := c.Param("id")

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var job models.Job
	if err := h.DB.Where("id = ? AND company_id = ?", jobID, company.ID).First(&job).Error; err != nil {
		utils.NotFound(c, "Job not found")
		return
	}

	var req JobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{
		"title":        req.Title,
		"location":     req.Location,
		"salary_min":   req.SalaryMin,
		"salary_max":   req.SalaryMax,
		"salary_type":  req.SalaryType,
		"description":  req.Description,
		"requirements": req.Requirements,
		"skills":       req.Skills,
		"job_type":     req.JobType,
	}
	if req.DepartmentID != nil {
		updates["department_id"] = *req.DepartmentID
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := h.DB.Model(&job).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Failed to update job")
		return
	}

	h.DB.Preload("Company").Preload("Department").First(&job, job.ID)
	utils.Success(c, job)
}

func (h *JobHandler) DeleteJob(c *gin.Context) {
	userID := middleware.GetUserID(c)
	jobID := c.Param("id")

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	result := h.DB.Where("id = ? AND company_id = ?", jobID, company.ID).Delete(&models.Job{})
	if result.RowsAffected == 0 {
		utils.NotFound(c, "Job not found")
		return
	}

	utils.Success(c, nil)
}

func (h *JobHandler) GetJob(c *gin.Context) {
	jobID := c.Param("id")

	var job models.Job
	if err := h.DB.Preload("Company").Preload("Department").
		First(&job, jobID).Error; err != nil {
		utils.NotFound(c, "Job not found")
		return
	}

	h.DB.Model(&job).UpdateColumn("views", gorm.Expr("views + ?", 1))

	h.DB.Create(&models.JobView{
		JobID:    job.ID,
		UserID:   nil,
		IP:       c.ClientIP(),
		ViewedAt: job.CreatedAt,
	})

	utils.Success(c, job)
}

func (h *JobHandler) ListJobs(c *gin.Context) {
	var jobs []models.Job
	query := h.DB.Preload("Company").Preload("Department")

	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR skills ILIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if location := c.Query("location"); location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	if jobType := c.Query("job_type"); jobType != "" {
		query = query.Where("job_type = ?", jobType)
	}

	if salaryMin := c.Query("salary_min"); salaryMin != "" {
		if min, err := strconv.ParseFloat(salaryMin, 64); err == nil {
			query = query.Where("salary_min >= ? OR salary_max >= ?", min, min)
		}
	}

	if salaryMax := c.Query("salary_max"); salaryMax != "" {
		if max, err := strconv.ParseFloat(salaryMax, 64); err == nil {
			query = query.Where("salary_max <= ?", max)
		}
	}

	if departmentID := c.Query("department_id"); departmentID != "" {
		query = query.Where("department_id = ?", departmentID)
	}

	query = query.Where("status = ?", models.JobStatusOpen)

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	query = query.Order(sortBy + " " + sortOrder)

	var total int64
	query.Model(&models.Job{}).Count(&total)

	page := 1
	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if psNum, err := strconv.Atoi(ps); err == nil && psNum > 0 {
			pageSize = psNum
		}
	}

	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&jobs)

	utils.SuccessWithPagination(c, jobs, page, pageSize, total)
}

func (h *JobHandler) ListCompanyJobs(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var jobs []models.Job
	query := h.DB.Preload("Company").Preload("Department").Where("company_id = ?", company.ID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if departmentID := c.Query("department_id"); departmentID != "" {
		query = query.Where("department_id = ?", departmentID)
	}

	var total int64
	query.Model(&models.Job{}).Count(&total)

	page := 1
	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if psNum, err := strconv.Atoi(ps); err == nil && psNum > 0 {
			pageSize = psNum
		}
	}

	query.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&jobs)

	utils.SuccessWithPagination(c, jobs, page, pageSize, total)
}

func (h *JobHandler) GetJobStats(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var stats struct {
		TotalJobs      int64 `json:"total_jobs"`
		OpenJobs       int64 `json:"open_jobs"`
		TotalViews     int64 `json:"total_views"`
		TotalApplications int64 `json:"total_applications"`
	}

	h.DB.Model(&models.Job{}).Where("company_id = ?", company.ID).Count(&stats.TotalJobs)
	h.DB.Model(&models.Job{}).Where("company_id = ? AND status = ?", company.ID, models.JobStatusOpen).Count(&stats.OpenJobs)
	h.DB.Model(&models.Job{}).Where("company_id = ?", company.ID).Select("COALESCE(SUM(views), 0)").Scan(&stats.TotalViews)
	h.DB.Table("applications").Joins("JOIN jobs ON applications.job_id = jobs.id").Where("jobs.company_id = ?", company.ID).Count(&stats.TotalApplications)

	utils.Success(c, stats)
}

func (h *JobHandler) UpdateJobStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	jobID := c.Param("id")

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var req struct {
		Status models.JobStatus `json:"status" binding:"required,oneof=open paused closed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result := h.DB.Model(&models.Job{}).Where("id = ? AND company_id = ?", jobID, company.ID).Update("status", req.Status)
	if result.RowsAffected == 0 {
		utils.NotFound(c, "Job not found")
		return
	}

	utils.Success(c, nil)
}
