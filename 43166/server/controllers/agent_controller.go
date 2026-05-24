package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"business-registration-platform/services"
	"business-registration-platform/utils"
)

type AgentController struct {
	agentService *services.AgentService
}

func NewAgentController() *AgentController {
	return &AgentController{
		agentService: services.NewAgentService(),
	}
}

func (ctrl *AgentController) GetAgentList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	agents, total, err := ctrl.agentService.GetAgentList(page, pageSize, status, keyword)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":     agents,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (ctrl *AgentController) GetAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid agent ID")
		return
	}

	agent, err := ctrl.agentService.GetAgentByID(uint(id))
	if err != nil {
		utils.NotFound(c, "Agent not found")
		return
	}

	utils.Success(c, agent)
}

func (ctrl *AgentController) CreateAgent(c *gin.Context) {
	var req services.CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	agent, err := ctrl.agentService.CreateAgent(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, agent)
}

func (ctrl *AgentController) UpdateAgentProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid agent ID")
		return
	}

	var req services.UpdateAgentProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.agentService.UpdateAgentProfile(uint(id), &req); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *AgentController) DeleteAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid agent ID")
		return
	}

	if err := ctrl.agentService.DeleteAgent(uint(id)); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *AgentController) AutoAssignAgent(c *gin.Context) {
	applicationID, err := strconv.ParseUint(c.Param("applicationId"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid application ID")
		return
	}

	agent, err := ctrl.agentService.AutoAssignAgent(uint(applicationID))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, agent)
}

func (ctrl *AgentController) GetAgentApplications(c *gin.Context) {
	userID, _ := c.Get("userID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")

	applications, total, err := ctrl.agentService.GetAgentApplications(userID.(uint), page, pageSize, status)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":     applications,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (ctrl *AgentController) GetAgentStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid agent ID")
		return
	}

	stats, err := ctrl.agentService.GetAgentStats(uint(id))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

func (ctrl *AgentController) GetAgentPerformanceReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid agent ID")
		return
	}

	var req struct {
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
	}
	c.ShouldBindJSON(&req)

	report, err := ctrl.agentService.GetAgentPerformanceReport(uint(id), parseTime(req.StartDate), parseTime(req.EndDate))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, report)
}

func (ctrl *AgentController) GetAvailableAgents(c *gin.Context) {
	agents, err := ctrl.agentService.GetAvailableAgents()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, agents)
}

func (ctrl *AgentController) UpdateAgentWorkSchedule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid agent ID")
		return
	}

	var req struct {
		StartTime string `json:"startTime" binding:"required"`
		EndTime   string `json:"endTime" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.agentService.UpdateAgentWorkSchedule(uint(id), req.StartTime, req.EndTime); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *AgentController) UpdateAgentMaxApps(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid agent ID")
		return
	}

	var req struct {
		MaxApps int `json:"maxApps" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.agentService.UpdateAgentMaxApps(uint(id), req.MaxApps); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func parseTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}
