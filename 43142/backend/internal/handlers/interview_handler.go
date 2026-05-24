package handlers

import (
	"github.com/gin-gonic/gin"

	"recruitment-platform/internal/models"
	"recruitment-platform/internal/services"
	"recruitment-platform/internal/utils"
	"recruitment-platform/pkg/logger"
)

type InterviewHandler struct {
	interviewService *services.InterviewService
}

func NewInterviewHandler(interviewService *services.InterviewService) *InterviewHandler {
	return &InterviewHandler{interviewService: interviewService}
}

func (h *InterviewHandler) ScheduleInterview(c *gin.Context) {
	companyID := c.GetUint("user_id")

	var req services.ScheduleInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	interview, err := h.interviewService.ScheduleInterview(companyID, &req)
	if err != nil {
		logger.Error("安排面试失败: %v", err)
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, interview)
}

func (h *InterviewHandler) GetInterview(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的面试ID")
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("user_role")

	interview, err := h.interviewService.GetInterview(id, userID, models.UserRole(role))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, interview)
}

func (h *InterviewHandler) ListCompanyInterviews(c *gin.Context) {
	companyID := c.GetUint("user_id")
	page, pageSize := utils.GetPagination(c)
	status := c.Query("status")

	interviews, total, err := h.interviewService.ListByCompany(companyID, page, pageSize, status)
	if err != nil {
		logger.Error("获取面试列表失败: %v", err)
		utils.InternalError(c, "获取面试列表失败")
		return
	}

	utils.Paginated(c, interviews, page, pageSize, total)
}

func (h *InterviewHandler) ListMyInterviews(c *gin.Context) {
	applicantID := c.GetUint("user_id")
	page, pageSize := utils.GetPagination(c)
	status := c.Query("status")

	interviews, total, err := h.interviewService.ListByApplicant(applicantID, page, pageSize, status)
	if err != nil {
		logger.Error("获取面试列表失败: %v", err)
		utils.InternalError(c, "获取面试列表失败")
		return
	}

	utils.Paginated(c, interviews, page, pageSize, total)
}

func (h *InterviewHandler) UpdateInterview(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的面试ID")
		return
	}

	companyID := c.GetUint("user_id")

	var req services.UpdateInterviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	interview, err := h.interviewService.UpdateInterview(id, companyID, &req)
	if err != nil {
		logger.Error("更新面试失败: %v", err)
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, interview)
}

func (h *InterviewHandler) AcceptInterview(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的面试ID")
		return
	}

	applicantID := c.GetUint("user_id")

	if err := h.interviewService.AcceptInterview(id, applicantID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "已接受面试", nil)
}

func (h *InterviewHandler) RejectInterview(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的面试ID")
		return
	}

	applicantID := c.GetUint("user_id")

	if err := h.interviewService.RejectInterview(id, applicantID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "已拒绝面试", nil)
}

func (h *InterviewHandler) CompleteInterview(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的面试ID")
		return
	}

	companyID := c.GetUint("user_id")

	var req struct {
		Feedback string `json:"feedback"`
		Rating   int    `json:"rating"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.interviewService.CompleteInterview(id, companyID, req.Feedback, req.Rating); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "面试已完成", nil)
}

func (h *InterviewHandler) CancelInterview(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的面试ID")
		return
	}

	companyID := c.GetUint("user_id")

	if err := h.interviewService.CancelInterview(id, companyID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "面试已取消", nil)
}
