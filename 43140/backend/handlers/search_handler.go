package handlers

import (
	"strconv"
	"strings"

	"recruitment-platform/models"
	"recruitment-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SearchHandler struct {
	DB *gorm.DB
}

func NewSearchHandler(db *gorm.DB) *SearchHandler {
	return &SearchHandler{DB: db}
}

func (h *SearchHandler) SearchJobs(c *gin.Context) {
	var jobs []models.Job
	query := h.DB.Preload("Company").Preload("Department")

	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR requirements ILIKE ? OR skills ILIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if location := c.Query("location"); location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	if jobType := c.Query("job_type"); jobType != "" {
		jobTypes := strings.Split(jobType, ",")
		query = query.Where("job_type IN ?", jobTypes)
	}

	if salaryMin := c.Query("salary_min"); salaryMin != "" {
		if min, err := strconv.ParseFloat(salaryMin, 64); err == nil {
			query = query.Where("salary_min >= ? OR (salary_min = 0 AND salary_max >= ?)", min, min)
		}
	}

	if salaryMax := c.Query("salary_max"); salaryMax != "" {
		if max, err := strconv.ParseFloat(salaryMax, 64); err == nil {
			query = query.Where("salary_max <= ? OR salary_max = 0", max)
		}
	}

	if companyID := c.Query("company_id"); companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	if departmentID := c.Query("department_id"); departmentID != "" {
		query = query.Where("department_id = ?", departmentID)
	}

	if skills := c.Query("skills"); skills != "" {
		skillList := strings.Split(skills, ",")
		for _, skill := range skillList {
			skill = strings.TrimSpace(skill)
			if skill != "" {
				query = query.Where("skills ILIKE ?", "%"+skill+"%")
			}
		}
	}

	query = query.Where("status = ?", models.JobStatusOpen)

	sortBy := c.DefaultQuery("sort_by", "created_at")
	allowedSortFields := map[string]bool{
		"created_at": true,
		"salary_min": true,
		"salary_max": true,
		"views":      true,
		"title":      true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}

	sortOrder := c.DefaultQuery("sort_order", "desc")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
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
		if psNum, err := strconv.Atoi(ps); err == nil && psNum > 0 && psNum <= 100 {
			pageSize = psNum
		}
	}

	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&jobs)

	utils.SuccessWithPagination(c, jobs, page, pageSize, total)
}

func (h *SearchHandler) GetJobTypes(c *gin.Context) {
	types := []map[string]string{
		{"value": "full-time", "label": "Full Time"},
		{"value": "part-time", "label": "Part Time"},
		{"value": "contract", "label": "Contract"},
		{"value": "internship", "label": "Internship"},
		{"value": "remote", "label": "Remote"},
	}
	utils.Success(c, types)
}

func (h *SearchHandler) GetLocations(c *gin.Context) {
	var locations []string
	h.DB.Model(&models.Job{}).
		Where("status = ? AND location != ?", models.JobStatusOpen, "").
		Distinct("location").
		Pluck("location", &locations)

	utils.Success(c, locations)
}

func (h *SearchHandler) GetSalaryRanges(c *gin.Context) {
	ranges := []map[string]interface{}{
		{"value": "0-50000", "label": "< $50k", "min": 0, "max": 50000},
		{"value": "50000-80000", "label": "$50k - $80k", "min": 50000, "max": 80000},
		{"value": "80000-120000", "label": "$80k - $120k", "min": 80000, "max": 120000},
		{"value": "120000-180000", "label": "$120k - $180k", "min": 120000, "max": 180000},
		{"value": "180000-300000", "label": "$180k - $300k", "min": 180000, "max": 300000},
		{"value": "300000+", "label": "$300k+", "min": 300000, "max": 0},
	}
	utils.Success(c, ranges)
}

func (h *SearchHandler) GetAllSkills(c *gin.Context) {
	var skills []string
	h.DB.Model(&models.Skill{}).Distinct("name").Pluck("name", &skills)
	utils.Success(c, skills)
}
