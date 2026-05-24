package handlers

import (
	"github.com/gin-gonic/gin"

	"recruitment-platform/internal/models"
	"recruitment-platform/internal/services"
	"recruitment-platform/internal/utils"
	"recruitment-platform/pkg/logger"
)

type ApplicationHandler struct {
	applicationService *services.ApplicationService
}

func NewApplicationHandler(applicationService *services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{applicationService: applicationService}
}

func (h *ApplicationHandler) Apply(c *gin.Context) {
	applicantID := c.GetUint("user_id")

	var req services.ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	application, err := h.applicationService.Apply(applicantID, &req)
	if err != nil {
		logger.Error("投递失败: %v", err)
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, application)
}

func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的投递ID")
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("user_role")

	application, err := h.applicationService.GetApplication(id, userID, models.UserRole(role))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, application)
}

func (h *ApplicationHandler) ListMyApplications(c *gin.Context) {
	applicantID := c.GetUint("user_id")
	page, pageSize := utils.GetPagination(c)
	status := c.Query("status")

	applications, total, err := h.applicationService.ListByApplicant(applicantID, page, pageSize, status)
	if err != nil {
		logger.Error("获取投递列表失败: %v", err)
		utils.InternalError(c, "获取投递列表失败")
		return
	}

	utils.Paginated(c, applications, page, pageSize, total)
}

func (h *ApplicationHandler) ListJobApplications(c *gin.Context) {
	jobID := parseUintParam(c.Param("jobId"))
	if jobID == 0 {
		utils.BadRequest(c, "无效的职位ID")
		return
	}

	companyID := c.GetUint("user_id")
	page, pageSize := utils.GetPagination(c)
	status := c.Query("status")

	applications, total, err := h.applicationService.ListByJob(jobID, companyID, page, pageSize, status)
	if err != nil {
		logger.Error("获取投递列表失败: %v", err)
		utils.InternalError(c, "获取投递列表失败")
		return
	}

	utils.Paginated(c, applications, page, pageSize, total)
}

func (h *ApplicationHandler) ListCompanyApplications(c *gin.Context) {
	companyID := c.GetUint("user_id")
	page, pageSize := utils.GetPagination(c)
	status := c.Query("status")
	keyword := c.Query("keyword")

	applications, total, err := h.applicationService.ListByCompany(companyID, page, pageSize, status, keyword)
	if err != nil {
		logger.Error("获取投递列表失败: %v", err)
		utils.InternalError(c, "获取投递列表失败")
		return
	}

	utils.Paginated(c, applications, page, pageSize, total)
}

func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的投递ID")
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("user_role")

	var req services.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	application, err := h.applicationService.UpdateStatus(id, userID, models.UserRole(role), &req)
	if err != nil {
		logger.Error("更新状态失败: %v", err)
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, application)
}

func (h *ApplicationHandler) BulkUpdateStatus(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req services.BulkUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	count, err := h.applicationService.BulkUpdateStatus(req.IDs, req.Status, req.Reason, userID)
	if err != nil {
		logger.Error("批量更新状态失败: %v", err)
		utils.InternalError(c, "批量更新失败")
		return
	}

	utils.Success(c, gin.H{"updated_count": count})
}

func (h *ApplicationHandler) Withdraw(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的投递ID")
		return
	}

	applicantID := c.GetUint("user_id")

	if err := h.applicationService.Withdraw(id, applicantID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "撤回成功", nil)
}

func (h *ApplicationHandler) GetHistory(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的投递ID")
		return
	}

	history, err := h.applicationService.GetHistory(id)
	if err != nil {
		logger.Error("获取历史记录失败: %v", err)
		utils.InternalError(c, "获取历史记录失败")
		return
	}

	utils.Success(c, history)
}

func (h *ApplicationHandler) GetStatusCount(c *gin.Context) {
	jobID := parseUintParam(c.Param("jobId"))
	if jobID == 0 {
		utils.BadRequest(c, "无效的职位ID")
		return
	}

	counts, err := h.applicationService.GetStatusCountByJob(jobID)
	if err != nil {
		logger.Error("获取状态统计失败: %v", err)
		utils.InternalError(c, "获取状态统计失败")
		return
	}

	utils.Success(c, counts)
}
