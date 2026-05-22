package handler

import (
	"car-rental/internal/middleware"
	"car-rental/internal/service"
	"car-rental/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MaintenanceHandler struct {
	maintenanceService *service.MaintenanceService
}

func NewMaintenanceHandler() *MaintenanceHandler {
	return &MaintenanceHandler{
		maintenanceService: service.NewMaintenanceService(),
	}
}

func (h *MaintenanceHandler) CreateMaintenance(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	var req service.CreateMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	plan, err := h.maintenanceService.CreateMaintenance(user.UserID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, plan)
}

func (h *MaintenanceHandler) GetMaintenanceByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	plan, err := h.maintenanceService.GetMaintenanceByID(uint(id))
	if err != nil {
		utils.NotFound(c, "维护计划不存在")
		return
	}

	utils.Success(c, plan)
}

func (h *MaintenanceHandler) GetAllMaintenance(c *gin.Context) {
	page, pageSize, _, _ := utils.ParsePageParams(c)
	carID, _ := strconv.ParseUint(c.Query("car_id"), 10, 64)
	status := c.Query("status")

	plans, total, err := h.maintenanceService.GetAllMaintenance(page, pageSize, uint(carID), status)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, plans, total, page, pageSize)
}

func (h *MaintenanceHandler) GetCarMaintenance(c *gin.Context) {
	carID, _ := strconv.ParseUint(c.Param("carId"), 10, 64)
	page, pageSize, _, _ := utils.ParsePageParams(c)

	plans, total, err := h.maintenanceService.GetCarMaintenance(uint(carID), page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, plans, total, page, pageSize)
}

func (h *MaintenanceHandler) StartMaintenance(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.maintenanceService.StartMaintenance(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *MaintenanceHandler) CompleteMaintenance(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.maintenanceService.CompleteMaintenance(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *MaintenanceHandler) CancelMaintenance(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.maintenanceService.CancelMaintenance(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *MaintenanceHandler) UpdateMaintenance(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.maintenanceService.UpdateMaintenance(uint(id), updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *MaintenanceHandler) DeleteMaintenance(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.maintenanceService.DeleteMaintenance(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *MaintenanceHandler) GetUpcomingMaintenance(c *gin.Context) {
	plans, err := h.maintenanceService.GetUpcomingMaintenance()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, plans)
}
