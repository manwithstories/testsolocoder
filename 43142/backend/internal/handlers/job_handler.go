package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"recruitment-platform/internal/services"
	"recruitment-platform/internal/utils"
	"recruitment-platform/pkg/logger"
)

type JobHandler struct {
	jobService *services.JobService
}

func NewJobHandler(jobService *services.JobService) *JobHandler {
	return &JobHandler{jobService: jobService}
}

func parseUintParam(s string) uint {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(n)
}

func (h *JobHandler) CreateJob(c *gin.Context) {
	companyID := c.GetUint("user_id")

	var req services.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	job, err := h.jobService.CreateJob(companyID, &req)
	if err != nil {
		logger.Error("创建职位失败: %v", err)
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, job)
}

func (h *JobHandler) GetJob(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的职位ID")
		return
	}

	userID, _ := c.Get("user_id")
	var uidPtr *uint
	if uid, ok := userID.(uint); ok && uid > 0 {
		uidPtr = &uid
	}

	job, err := h.jobService.GetJobWithView(id, c.ClientIP(), uidPtr)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, job)
}

func (h *JobHandler) UpdateJob(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的职位ID")
		return
	}

	companyID := c.GetUint("user_id")

	var req services.UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	job, err := h.jobService.UpdateJob(id, companyID, &req)
	if err != nil {
		logger.Error("更新职位失败: %v", err)
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, job)
}

func (h *JobHandler) DeleteJob(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的职位ID")
		return
	}

	companyID := c.GetUint("user_id")

	if err := h.jobService.DeleteJob(id, companyID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *JobHandler) ListJobs(c *gin.Context) {
	page, pageSize := utils.GetPagination(c)
	status := c.Query("status")
	keyword := c.Query("keyword")

	var companyIDPtr *uint
	if cid := c.Query("company_id"); cid != "" {
		cidUint := parseUintParam(cid)
		if cidUint > 0 {
			companyIDPtr = &cidUint
		}
	}

	jobs, total, err := h.jobService.ListJobs(page, pageSize, companyIDPtr, status, keyword)
	if err != nil {
		logger.Error("获取职位列表失败: %v", err)
		utils.InternalError(c, "获取职位列表失败")
		return
	}

	utils.Paginated(c, jobs, page, pageSize, total)
}

func (h *JobHandler) ListMyJobs(c *gin.Context) {
	page, pageSize := utils.GetPagination(c)
	companyID := c.GetUint("user_id")
	status := c.Query("status")

	jobs, total, err := h.jobService.ListJobs(page, pageSize, &companyID, status, "")
	if err != nil {
		logger.Error("获取职位列表失败: %v", err)
		utils.InternalError(c, "获取职位列表失败")
		return
	}

	utils.Paginated(c, jobs, page, pageSize, total)
}

func (h *JobHandler) PublishJob(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的职位ID")
		return
	}

	companyID := c.GetUint("user_id")

	if err := h.jobService.PublishJob(id, companyID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "发布成功", nil)
}

func (h *JobHandler) CloseJob(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的职位ID")
		return
	}

	companyID := c.GetUint("user_id")

	if err := h.jobService.CloseJob(id, companyID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "关闭成功", nil)
}

func (h *JobHandler) BulkImport(c *gin.Context) {
	companyID := c.GetUint("user_id")

	var req struct {
		Jobs []services.BulkJobItem `json:"jobs" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	successCount, errors := h.jobService.BulkImport(companyID, req.Jobs)

	utils.Success(c, gin.H{
		"success_count": successCount,
		"errors":        errors,
	})
}

func (h *JobHandler) BulkDelete(c *gin.Context) {
	companyID := c.GetUint("user_id")

	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.jobService.BulkDelete(req.IDs, companyID); err != nil {
		utils.InternalError(c, "批量删除失败")
		return
	}

	utils.SuccessWithMessage(c, "批量删除成功", nil)
}

func (h *JobHandler) ExportJobs(c *gin.Context) {
	companyID := c.GetUint("user_id")

	data, err := h.jobService.ExportJobs(companyID)
	if err != nil {
		utils.InternalError(c, "导出失败")
		return
	}

	utils.Success(c, data)
}

func (h *JobHandler) GetViewStats(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的职位ID")
		return
	}

	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 30 {
			days = parsed
		}
	}

	count, err := h.jobService.GetViewStats(id, days)
	if err != nil {
		utils.InternalError(c, "获取统计数据失败")
		return
	}

	utils.Success(c, gin.H{
		"days":  days,
		"views": count,
	})
}
