package controllers

import (
	"health-platform/services"
	"health-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BillingController struct {
	billingService *services.BillingService
}

func NewBillingController() *BillingController {
	return &BillingController{
		billingService: services.NewBillingService(),
	}
}

func (ctrl *BillingController) GenerateMonthlyBilling(c *gin.Context) {
	var req services.GenerateBillingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	billing, err := ctrl.billingService.GenerateMonthlyBilling(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, billing)
}

func (ctrl *BillingController) GetBilling(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	billing, err := ctrl.billingService.GetBilling(uint(id))
	if err != nil {
		utils.Error(c, 404, "账单不存在")
		return
	}

	utils.Success(c, billing)
}

func (ctrl *BillingController) GetCompanyBillings(c *gin.Context) {
	companyID := c.GetUint("company_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	billings, total, err := ctrl.billingService.GetCompanyBillings(companyID, page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, billings)
}

func (ctrl *BillingController) GetAgencyBillings(c *gin.Context) {
	agencyID := c.GetUint("agency_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	billings, total, err := ctrl.billingService.GetAgencyBillings(agencyID, page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, billings)
}

func (ctrl *BillingController) PayBilling(c *gin.Context) {
	var req services.PayBillingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.billingService.PayBilling(&req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *BillingController) Recharge(c *gin.Context) {
	var req services.RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.billingService.Recharge(&req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *BillingController) GetTransactions(c *gin.Context) {
	companyID := c.GetUint("company_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	transactions, total, err := ctrl.billingService.GetTransactions(companyID, page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, transactions)
}

func (ctrl *BillingController) GetCompanyBalance(c *gin.Context) {
	companyID := c.GetUint("company_id")
	
	balance, err := ctrl.billingService.GetCompanyBalance(companyID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, balance)
}
