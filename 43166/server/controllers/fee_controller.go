package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"business-registration-platform/models"
	"business-registration-platform/services"
	"business-registration-platform/utils"
)

type FeeController struct {
	feeService *services.FeeService
}

func NewFeeController() *FeeController {
	return &FeeController{
		feeService: services.NewFeeService(),
	}
}

func (ctrl *FeeController) CalculateFee(c *gin.Context) {
	var req services.CalculateFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	fee, err := ctrl.feeService.CalculateFee(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, fee)
}

func (ctrl *FeeController) CreateApplicationFee(c *gin.Context) {
	var req struct {
		ApplicationID uint              `json:"applicationId" binding:"required"`
		CompanyType   models.CompanyType `json:"companyType" binding:"required"`
		Capital       float64           `json:"capital" binding:"required"`
		DiscountCode  string            `json:"discountCode"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	fee, err := ctrl.feeService.CreateApplicationFee(req.ApplicationID, req.CompanyType, req.Capital, req.DiscountCode)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, fee)
}

func (ctrl *FeeController) PayFee(c *gin.Context) {
	var req services.PayFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	fee, err := ctrl.feeService.PayFee(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, fee)
}

func (ctrl *FeeController) GetApplicationFee(c *gin.Context) {
	applicationID, err := strconv.ParseUint(c.Param("applicationId"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid application ID")
		return
	}

	fee, err := ctrl.feeService.GetApplicationFee(uint(applicationID))
	if err != nil {
		utils.NotFound(c, "Fee not found")
		return
	}

	utils.Success(c, fee)
}

func (ctrl *FeeController) GetFeeList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")

	fees, total, err := ctrl.feeService.GetFeeList(page, pageSize, status)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":     fees,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (ctrl *FeeController) GetFeeStandards(c *gin.Context) {
	standards, err := ctrl.feeService.GetFeeStandards()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, standards)
}

func (ctrl *FeeController) UpdateFeeStandard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid standard ID")
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.feeService.UpdateFeeStandard(uint(id), data); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *FeeController) CreateFeeStandard(c *gin.Context) {
	var standard models.FeeStandard
	if err := c.ShouldBindJSON(&standard); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.feeService.CreateFeeStandard(&standard); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, standard)
}

func (ctrl *FeeController) GetDiscountPolicies(c *gin.Context) {
	policies, err := ctrl.feeService.GetDiscountPolicies()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, policies)
}

func (ctrl *FeeController) CreateDiscountPolicy(c *gin.Context) {
	var policy models.DiscountPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.feeService.CreateDiscountPolicy(&policy); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, policy)
}

func (ctrl *FeeController) UpdateDiscountPolicy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid policy ID")
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.feeService.UpdateDiscountPolicy(uint(id), data); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *FeeController) DeleteDiscountPolicy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid policy ID")
		return
	}

	if err := ctrl.feeService.DeleteDiscountPolicy(uint(id)); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}
