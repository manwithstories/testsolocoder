package handlers

import (
	"fmt"
	"time"

	"recruitment-platform/middleware"
	"recruitment-platform/models"
	"recruitment-platform/utils"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ExportHandler struct {
	DB *gorm.DB
}

func NewExportHandler(db *gorm.DB) *ExportHandler {
	return &ExportHandler{DB: db}
}

func (h *ExportHandler) ExportApplications(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	jobID := c.Query("job_id")
	status := c.Query("status")

	var applications []models.Application
	query := h.DB.Preload("Job").
		Preload("JobSeeker").Preload("JobSeeker.User").
		Preload("Resume").
		Joins("JOIN jobs ON applications.job_id = jobs.id").
		Where("jobs.company_id = ?", company.ID)

	if jobID != "" {
		query = query.Where("applications.job_id = ?", jobID)
	}
	if status != "" {
		query = query.Where("applications.status = ?", status)
	}

	query.Order("applications.created_at desc").Find(&applications)

	f := excelize.NewFile()
	sheetName := "Applications"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"ID", "Job Title", "Applicant Name", "Email", "Phone", "Status", "Applied Date"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, app := range applications {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), app.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), app.Job.Title)
		if app.JobSeeker.User != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), app.JobSeeker.User.Name)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), app.JobSeeker.User.Email)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), app.JobSeeker.User.Phone)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), string(app.Status))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), app.CreatedAt.Format("2006-01-02 15:04"))
	}

	filename := fmt.Sprintf("applications_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("File-Name", filename)

	f.Write(c.Writer)
}

func (h *ExportHandler) ExportInterviews(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	status := c.Query("status")

	var interviews []models.Interview
	query := h.DB.Preload("Application").Preload("Application.Job").
		Preload("Application.JobSeeker").Preload("Application.JobSeeker.User").
		Joins("JOIN applications ON interviews.application_id = applications.id").
		Joins("JOIN jobs ON applications.job_id = jobs.id").
		Where("jobs.company_id = ?", company.ID)

	if status != "" {
		query = query.Where("interviews.status = ?", status)
	}

	query.Order("interviews.scheduled_at desc").Find(&interviews)

	f := excelize.NewFile()
	sheetName := "Interviews"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"ID", "Job Title", "Candidate", "Scheduled At", "Duration", "Location", "Interviewer", "Status"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, interview := range interviews {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), interview.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), interview.Application.Job.Title)
		if interview.Application.JobSeeker.User != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), interview.Application.JobSeeker.User.Name)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), interview.ScheduledAt.Format("2006-01-02 15:04"))
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), interview.Duration)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), interview.Location)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), interview.Interviewer)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), string(interview.Status))
	}

	filename := fmt.Sprintf("interviews_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("File-Name", filename)

	f.Write(c.Writer)
}

func (h *ExportHandler) ExportJobs(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var company models.Company
	if err := h.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		utils.NotFound(c, "Company not found")
		return
	}

	var jobs []models.Job
	h.DB.Preload("Department").
		Where("company_id = ?", company.ID).
		Order("created_at desc").Find(&jobs)

	f := excelize.NewFile()
	sheetName := "Jobs"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"ID", "Title", "Department", "Location", "Salary Min", "Salary Max", "Job Type", "Status", "Views", "Applications", "Created At"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, job := range jobs {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), job.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), job.Title)
		if job.Department != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), job.Department.Name)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), job.Location)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), job.SalaryMin)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), job.SalaryMax)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), string(job.JobType))
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), string(job.Status))
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), job.Views)

		var appCount int64
		h.DB.Model(&models.Application{}).Where("job_id = ?", job.ID).Count(&appCount)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), appCount)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), job.CreatedAt.Format("2006-01-02 15:04"))
	}

	filename := fmt.Sprintf("jobs_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("File-Name", filename)

	f.Write(c.Writer)
}
