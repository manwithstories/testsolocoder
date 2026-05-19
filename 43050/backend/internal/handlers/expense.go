package handlers

import (
	"errors"
	"net/http"
	"splitwise-clone/internal/database"
	"splitwise-clone/internal/models"
	"splitwise-clone/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ParticipantInput struct {
	UserID uint    `json:"userId" binding:"required"`
	Amount float64 `json:"amount"`
	Ratio  float64 `json:"ratio"`
}

type CreateExpenseRequest struct {
	Title        string              `json:"title" binding:"required"`
	Amount       float64             `json:"amount" binding:"required,gt=0"`
	PaidBy       uint                `json:"paidBy" binding:"required"`
	SplitType    models.SplitType    `json:"splitType" binding:"required,oneof=equal ratio custom"`
	ExpenseDate  string              `json:"expenseDate"`
	Participants []ParticipantInput  `json:"participants" binding:"required,min=1"`
	Version      int                 `json:"version"`
}

func CreateExpense(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("groupId")

	var req CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var member models.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ? AND is_active = ?", groupID, userID, true).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a member of this group"})
		return
	}

	expenseDate := time.Now()
	if req.ExpenseDate != "" {
		if parsed, err := time.Parse(time.RFC3339, req.ExpenseDate); err == nil {
			expenseDate = parsed
		}
	}

	var participants []models.ExpenseParticipant
	totalSplitAmount := 0.0

	switch req.SplitType {
	case models.SplitEqual:
		equalAmount := utils.RoundFloat(req.Amount/float64(len(req.Participants)), 2)
		for i, p := range req.Participants {
			amount := equalAmount
			if i == len(req.Participants)-1 {
				amount = utils.RoundFloat(req.Amount-equalAmount*float64(len(req.Participants)-1), 2)
			}
			participants = append(participants, models.ExpenseParticipant{
				UserID: p.UserID,
				Amount: amount,
			})
			totalSplitAmount += amount
		}

	case models.SplitRatio:
		totalRatio := 0.0
		for _, p := range req.Participants {
			totalRatio += p.Ratio
		}
		if totalRatio <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Total ratio must be greater than 0"})
			return
		}
		cumulativeAmount := 0.0
		for i, p := range req.Participants {
			amount := utils.RoundFloat(req.Amount*p.Ratio/totalRatio, 2)
			if i == len(req.Participants)-1 {
				amount = utils.RoundFloat(req.Amount-cumulativeAmount, 2)
			}
			participants = append(participants, models.ExpenseParticipant{
				UserID: p.UserID,
				Amount: amount,
				Ratio:  p.Ratio,
			})
			cumulativeAmount += amount
			totalSplitAmount += amount
		}

	case models.SplitCustom:
		for _, p := range req.Participants {
			if p.Amount <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Custom amount must be greater than 0"})
				return
			}
			participants = append(participants, models.ExpenseParticipant{
				UserID: p.UserID,
				Amount: p.Amount,
			})
			totalSplitAmount += p.Amount
		}
		if utils.RoundFloat(totalSplitAmount, 2) != utils.RoundFloat(req.Amount, 2) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Sum of custom amounts must equal total amount"})
			return
		}
	}

	expense := models.Expense{
		GroupID:     parseUint(groupID),
		Title:       req.Title,
		Amount:      req.Amount,
		PaidBy:      req.PaidBy,
		SplitType:   req.SplitType,
		CreatedBy:   userID,
		ExpenseDate: expenseDate,
		Version:     1,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&expense).Error; err != nil {
			return err
		}
		for i := range participants {
			participants[i].ExpenseID = expense.ID
		}
		if len(participants) > 0 {
			return tx.Create(&participants).Error
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	database.DB.Preload("Payer").Preload("Participants.User").First(&expense, expense.ID)
	c.JSON(http.StatusCreated, expense)
}

func GetGroupExpenses(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("groupId")

	var member models.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ? AND is_active = ?", groupID, userID, true).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a member of this group"})
		return
	}

	var expenses []models.Expense
	database.DB.Where("group_id = ?", groupID).
		Preload("Payer").
		Preload("Participants.User").
		Order("expense_date DESC, created_at DESC").
		Find(&expenses)

	c.JSON(http.StatusOK, expenses)
}

func GetExpenseByID(c *gin.Context) {
	userID := c.GetUint("userID")
	expenseID := c.Param("id")

	var expense models.Expense
	if err := database.DB.Preload("Payer").Preload("Participants.User").First(&expense, expenseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	var member models.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ? AND is_active = ?", expense.GroupID, userID, true).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a member of this group"})
		return
	}

	c.JSON(http.StatusOK, expense)
}

func UpdateExpense(c *gin.Context) {
	userID := c.GetUint("userID")
	expenseID := c.Param("id")

	var req CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expense models.Expense
	if err := database.DB.First(&expense, expenseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	if expense.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only creator can edit this expense"})
		return
	}

	var member models.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ? AND is_active = ?", expense.GroupID, userID, true).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a member of this group"})
		return
	}

	expID := parseUint(expenseID)

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if req.Version > 0 {
			if err := validateExpenseVersion(tx, expID, req.Version); err != nil {
				return err
			}
		}

		var existingExpense models.Expense
		if err := tx.First(&existingExpense, expenseID).Error; err != nil {
			return err
		}

		expenseDate := time.Now()
		if req.ExpenseDate != "" {
			if parsed, err := time.Parse(time.RFC3339, req.ExpenseDate); err == nil {
				expenseDate = parsed
			}
		}

		existingExpense.Title = req.Title
		existingExpense.Amount = req.Amount
		existingExpense.PaidBy = req.PaidBy
		existingExpense.SplitType = req.SplitType
		existingExpense.ExpenseDate = expenseDate
		existingExpense.Version++

		if err := tx.Save(&existingExpense).Error; err != nil {
			return err
		}

		if err := tx.Where("expense_id = ?", expenseID).Delete(&models.ExpenseParticipant{}).Error; err != nil {
			return err
		}

		var participants []models.ExpenseParticipant
		switch req.SplitType {
		case models.SplitEqual:
			equalAmount := utils.RoundFloat(req.Amount/float64(len(req.Participants)), 2)
			for i, p := range req.Participants {
				amount := equalAmount
				if i == len(req.Participants)-1 {
					amount = utils.RoundFloat(req.Amount-equalAmount*float64(len(req.Participants)-1), 2)
				}
				participants = append(participants, models.ExpenseParticipant{
					ExpenseID: existingExpense.ID,
					UserID:    p.UserID,
					Amount:    amount,
				})
			}
		case models.SplitRatio:
			totalRatio := 0.0
			for _, p := range req.Participants {
				totalRatio += p.Ratio
			}
			cumulativeAmount := 0.0
			for i, p := range req.Participants {
				amount := utils.RoundFloat(req.Amount*p.Ratio/totalRatio, 2)
				if i == len(req.Participants)-1 {
					amount = utils.RoundFloat(req.Amount-cumulativeAmount, 2)
				}
				participants = append(participants, models.ExpenseParticipant{
					ExpenseID: existingExpense.ID,
					UserID:    p.UserID,
					Amount:    amount,
					Ratio:     p.Ratio,
				})
				cumulativeAmount += amount
			}
		case models.SplitCustom:
			for _, p := range req.Participants {
				participants = append(participants, models.ExpenseParticipant{
					ExpenseID: existingExpense.ID,
					UserID:    p.UserID,
					Amount:    p.Amount,
				})
			}
		}

		if len(participants) > 0 {
			return tx.Create(&participants).Error
		}
		return nil
	})

	if err != nil {
		if err.Error() == "expense has been modified by another user" {
			c.JSON(http.StatusConflict, gin.H{"error": "账单已被其他用户修改，请刷新后重试"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}

	database.DB.Preload("Payer").Preload("Participants.User").First(&expense, expenseID)
	c.JSON(http.StatusOK, expense)
}

func DeleteExpense(c *gin.Context) {
	userID := c.GetUint("userID")
	expenseID := c.Param("id")

	var expense models.Expense
	if err := database.DB.First(&expense, expenseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	if expense.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only creator can delete this expense"})
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("expense_id = ?", expenseID).Delete(&models.ExpenseParticipant{}).Error; err != nil {
			return err
		}
		return tx.Delete(&expense).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}

func parseUint(s string) uint {
	var result uint
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + uint(c-'0')
		}
	}
	return result
}

func validateExpenseVersion(tx *gorm.DB, expenseID uint, expectedVersion int) error {
	var expense models.Expense
	if err := tx.First(&expense, expenseID).Error; err != nil {
		return err
	}
	if expense.Version != expectedVersion {
		return errors.New("expense has been modified by another user")
	}
	return nil
}
