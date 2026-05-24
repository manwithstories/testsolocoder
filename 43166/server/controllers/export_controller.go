package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"business-registration-platform/services"
	"business-registration-platform/utils"
)

type ExportController struct {
	exportService *services.ExportService
}

func NewExportController() *ExportController {
	return &ExportController{
		exportService: services.NewExportService(),
	}
}

func (ctrl *ExportController) CreateExportTask(c *gin.Context) {
	userID, _ := c.Get("userID")

	var params services.ExportParams
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	task, err := ctrl.exportService.CreateExportTask(userID.(uint), &params)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, task)
}

func (ctrl *ExportController) GetExportTasks(c *gin.Context) {
	userID, _ := c.Get("userID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")

	tasks, total, err := ctrl.exportService.GetExportTasks(userID.(uint), page, pageSize, status)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":     tasks,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (ctrl *ExportController) GetExportTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid task ID")
		return
	}

	task, err := ctrl.exportService.GetExportTask(uint(taskID))
	if err != nil {
		utils.NotFound(c, "Task not found")
		return
	}

	utils.Success(c, task)
}

func (ctrl *ExportController) DownloadExport(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid task ID")
		return
	}

	userID, _ := c.Get("userID")

	filePath, err := ctrl.exportService.DownloadExport(uint(taskID), userID.(uint))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.NotFound(c, "File not found")
		return
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		utils.InternalServerError(c, "Failed to read file")
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.Name())
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	c.File(filePath)
	c.Status(http.StatusOK)
}
