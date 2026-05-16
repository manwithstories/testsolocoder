package controllers

import (
	"fmt"

	"finance-api/middleware"
	"finance-api/models"
	"finance-api/utils"

	"github.com/gin-gonic/gin"
)

type AccountController struct{}

func NewAccountController() *AccountController {
	return &AccountController{}
}

type CreateAccountRequest struct {
	Name           string  `json:"name" binding:"required,max=100"`
	Currency       string  `json:"currency" binding:"required,max=10"`
	InitialBalance float64 `json:"initial_balance"`
	Remark         string  `json:"remark" binding:"max=255"`
}

type UpdateAccountRequest struct {
	Name           string   `json:"name" binding:"omitempty,max=100"`
	Currency       string   `json:"currency" binding:"omitempty,max=10"`
	InitialBalance *float64 `json:"initial_balance"`
	Remark         string   `json:"remark" binding:"max=255"`
}

func (ctrl *AccountController) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters: "+err.Error())
		return
	}

	account := models.Account{
		UserID:         userID,
		Name:           req.Name,
		Currency:       req.Currency,
		InitialBalance: req.InitialBalance,
		Balance:        req.InitialBalance,
		Remark:         req.Remark,
	}

	if result := utils.DB.Create(&account); result.Error != nil {
		utils.InternalError(c, "Failed to create account: "+result.Error.Error())
		return
	}

	utils.Success(c, account)
}

func (ctrl *AccountController) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var accounts []models.Account
	if result := utils.DB.Where("user_id = ?", userID).Find(&accounts); result.Error != nil {
		utils.InternalError(c, "Failed to fetch accounts: "+result.Error.Error())
		return
	}

	utils.Success(c, accounts)
}

func (ctrl *AccountController) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var account models.Account
	if result := utils.DB.Where("id = ? AND user_id = ?", id, userID).First(&account); result.Error != nil {
		utils.NotFound(c, "Account not found")
		return
	}

	utils.Success(c, account)
}

func (ctrl *AccountController) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var account models.Account
	if result := utils.DB.Where("id = ? AND user_id = ?", id, userID).First(&account); result.Error != nil {
		utils.NotFound(c, "Account not found")
		return
	}

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Currency != "" {
		updates["currency"] = req.Currency
	}
	if req.InitialBalance != nil {
		balanceDiff := *req.InitialBalance - account.InitialBalance
		updates["initial_balance"] = *req.InitialBalance
		updates["balance"] = account.Balance + balanceDiff
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	if result := utils.DB.Model(&account).Updates(updates); result.Error != nil {
		utils.InternalError(c, "Failed to update account: "+result.Error.Error())
		return
	}

	utils.DB.First(&account, id)
	utils.Success(c, account)
}

func (ctrl *AccountController) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var account models.Account
	if result := utils.DB.Where("id = ? AND user_id = ?", id, userID).First(&account); result.Error != nil {
		utils.NotFound(c, "Account not found")
		return
	}

	var transactionCount int64
	utils.DB.Model(&models.Transaction{}).Where("account_id = ?", id).Count(&transactionCount)
	if transactionCount > 0 {
		utils.BadRequest(c, fmt.Sprintf("Cannot delete account: there are %d associated transaction records", transactionCount))
		return
	}

	if result := utils.DB.Delete(&account); result.Error != nil {
		utils.InternalError(c, "Failed to delete account: "+result.Error.Error())
		return
	}

	utils.Success(c, nil)
}
