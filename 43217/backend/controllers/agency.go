package controllers

import (
	"health-platform/services"
	"health-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AgencyController struct {
	agencyService *services.AgencyService
}

func NewAgencyController() *AgencyController {
	return &AgencyController{
		agencyService: services.NewAgencyService(),
	}
}

func (ctrl *AgencyController) RegisterAgency(c *gin.Context) {
	var req services.RegisterAgencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	agency, adminUser, err := ctrl.agencyService.RegisterAgency(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"agency":     agency,
		"admin_user": adminUser,
	})
}

func (ctrl *AgencyController) GetAgency(c *gin.Context) {
	agencyID := c.GetUint("agency_id")
	if agencyID == 0 {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		agencyID = uint(id)
	}

	agency, err := ctrl.agencyService.GetAgency(agencyID)
	if err != nil {
		utils.Error(c, 404, "机构不存在")
		return
	}

	utils.Success(c, agency)
}

func (ctrl *AgencyController) UpdateAgency(c *gin.Context) {
	agencyID := c.GetUint("agency_id")
	
	var req services.UpdateAgencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.agencyService.UpdateAgency(agencyID, &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *AgencyController) ListAgencies(c *gin.Context) {
	agencies, err := ctrl.agencyService.ListAgencies()
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, agencies)
}

func (ctrl *AgencyController) CreatePackage(c *gin.Context) {
	var req services.CreatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	pkg, err := ctrl.agencyService.CreatePackage(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, pkg)
}

func (ctrl *AgencyController) GetPackage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	pkg, err := ctrl.agencyService.GetPackage(uint(id))
	if err != nil {
		utils.Error(c, 404, "套餐不存在")
		return
	}

	utils.Success(c, pkg)
}

func (ctrl *AgencyController) UpdatePackage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	var req services.UpdatePackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.agencyService.UpdatePackage(uint(id), &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *AgencyController) UpdatePackagePrice(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	var req struct {
		Price float64 `json:"price" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.agencyService.UpdatePackagePrice(uint(id), req.Price); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *AgencyController) UpdatePackageStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.agencyService.UpdatePackageStatus(uint(id), 1); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *AgencyController) GetAgencyPackages(c *gin.Context) {
	agencyID := c.GetUint("agency_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	packages, total, err := ctrl.agencyService.GetAgencyPackages(agencyID, page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, packages)
}

func (ctrl *AgencyController) ListOnlinePackages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.DefaultQuery("keyword", "")

	packages, total, err := ctrl.agencyService.ListOnlinePackages(page, pageSize, keyword)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, packages)
}

func (ctrl *AgencyController) GetHotPackages(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	packages, err := ctrl.agencyService.GetHotPackages(limit)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, packages)
}

func (ctrl *AgencyController) CreateTimeSlot(c *gin.Context) {
	var req services.CreateTimeSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.agencyService.CreateTimeSlot(&req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *AgencyController) GetPackageTimeSlots(c *gin.Context) {
	packageID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	timeSlots, err := ctrl.agencyService.GetAvailableTimeSlots(uint(packageID))
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, timeSlots)
}
