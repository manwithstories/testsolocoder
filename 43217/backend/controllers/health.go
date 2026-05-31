package controllers

import (
	"health-platform/services"
	"health-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	healthService *services.HealthService
}

func NewHealthController() *HealthController {
	return &HealthController{
		healthService: services.NewHealthService(),
	}
}

func (ctrl *HealthController) GetHealthRecords(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	records, total, err := ctrl.healthService.GetHealthRecords(uint(employeeID), page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, records)
}

func (ctrl *HealthController) GetHealthRecordByYear(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)
	year, _ := strconv.Atoi(c.Param("year"))

	record, err := ctrl.healthService.GetHealthRecordByYear(uint(employeeID), year)
	if err != nil {
		utils.Error(c, 404, "健康档案不存在")
		return
	}

	utils.Success(c, record)
}

func (ctrl *HealthController) GetAllHealthRecords(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)

	records, err := ctrl.healthService.GetAllHealthRecords(uint(employeeID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, records)
}

func (ctrl *HealthController) GetTrendData(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)

	data, err := ctrl.healthService.GetTrendData(uint(employeeID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, data)
}

func (ctrl *HealthController) CreateAbnormalItem(c *gin.Context) {
	var req services.CreateAbnormalItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	item, err := ctrl.healthService.CreateAbnormalItem(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, item)
}

func (ctrl *HealthController) GetAbnormalItems(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	items, total, err := ctrl.healthService.GetAbnormalItems(uint(employeeID), page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, items)
}

func (ctrl *HealthController) GetAllAbnormalItems(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)

	items, err := ctrl.healthService.GetAllAbnormalItems(uint(employeeID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, items)
}

func (ctrl *HealthController) SetRecheckDate(c *gin.Context) {
	var req services.SetRecheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.healthService.SetRecheckDate(&req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *HealthController) UpdateRecheckStatus(c *gin.Context) {
	var req services.UpdateRecheckStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.healthService.UpdateRecheckStatus(&req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *HealthController) GetNeedRecheckItems(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)

	items, err := ctrl.healthService.GetNeedRecheckItems(uint(employeeID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, items)
}

func (ctrl *HealthController) GetReminders(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	reminders, total, err := ctrl.healthService.GetReminders(uint(employeeID), page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, reminders)
}

func (ctrl *HealthController) GetUnreadReminders(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)

	reminders, err := ctrl.healthService.GetUnreadReminders(uint(employeeID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, reminders)
}

func (ctrl *HealthController) MarkReminderAsRead(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := ctrl.healthService.MarkReminderAsRead(uint(id)); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *HealthController) GetHealthSummary(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Param("employee_id"), 10, 64)

	summary, err := ctrl.healthService.GenerateHealthSummary(uint(employeeID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, summary)
}
