package controllers

import (
	"fmt"
	"health-platform/config"
	"health-platform/services"
	"health-platform/utils"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	reportService *services.ReportService
}

func NewReportController() *ReportController {
	return &ReportController{
		reportService: services.NewReportService(),
	}
}

func (ctrl *ReportController) CreateReport(c *gin.Context) {
	var req services.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	report, err := ctrl.reportService.CreateReport(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	ctrl.reportService.UpdateHealthRecordFromReport(report.ID)

	utils.Success(c, report)
}

func (ctrl *ReportController) GetReport(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	report, err := ctrl.reportService.GetReport(uint(id))
	if err != nil {
		utils.Error(c, 404, "报告不存在")
		return
	}

	utils.Success(c, report)
}

func (ctrl *ReportController) GetReportByAppointment(c *gin.Context) {
	appointmentID, _ := strconv.ParseUint(c.Param("appointment_id"), 10, 64)
	
	report, err := ctrl.reportService.GetReportByAppointment(uint(appointmentID))
	if err != nil {
		utils.Error(c, 404, "报告不存在")
		return
	}

	utils.Success(c, report)
}

func (ctrl *ReportController) GetEmployeeReports(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	reports, total, err := ctrl.reportService.GetEmployeeReports(uint(employeeID), page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, reports)
}

func (ctrl *ReportController) GetCompanyReports(c *gin.Context) {
	companyID := c.GetUint("company_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	reports, total, err := ctrl.reportService.GetCompanyReports(companyID, page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, reports)
}

func (ctrl *ReportController) UploadReportFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请上传文件")
		return
	}

	ext := filepath.Ext(file.Filename)
	allowedExts := []string{".pdf", ".jpg", ".jpeg", ".png"}
	
	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	
	if !allowed {
		utils.BadRequest(c, "不支持的文件格式")
		return
	}

	uploadPath := config.GlobalConfig.Upload.Path
	os.MkdirAll(uploadPath, 0755)

	reportNo := c.DefaultPostForm("report_no", "")
	if reportNo == "" {
		reportNo = utils.GenerateOrderNo("RPT")
	}

	filename := fmt.Sprintf("%s%s", reportNo, ext)
	filePath := filepath.Join(uploadPath, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.Error(c, 500, "文件上传失败")
		return
	}

	utils.Success(c, gin.H{
		"file_url":  filePath,
		"file_name": filename,
	})
}

func (ctrl *ReportController) DownloadReportFile(c *gin.Context) {
	reportID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	report, err := ctrl.reportService.GetReport(uint(reportID))
	if err != nil {
		utils.Error(c, 404, "报告不存在")
		return
	}

	if report.PdfFile == "" {
		utils.Error(c, 400, "报告文件不存在")
		return
	}

	ctrl.reportService.MarkReportViewed(uint(reportID))

	filePath := report.PdfFile
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.Error(c, 404, "文件不存在")
		return
	}

	fileName := fmt.Sprintf("体检报告_%s.pdf", report.ReportNo)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/pdf")
	
	file, err := os.Open(filePath)
	if err != nil {
		utils.Error(c, 500, "文件读取失败")
		return
	}
	defer file.Close()

	io.Copy(c.Writer, file)
}

func (ctrl *ReportController) GetAbnormalReports(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)
	
	reports, err := ctrl.reportService.GetAbnormalReports(uint(employeeID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, reports)
}

func (ctrl *ReportController) GetAbnormalItems(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)
	
	items, err := ctrl.reportService.GetAbnormalItems(uint(employeeID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, items)
}
