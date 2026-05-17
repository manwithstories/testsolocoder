package handlers

import (
	"net/http"
	"splitwise-clone/internal/database"
	"splitwise-clone/internal/models"
	"splitwise-clone/internal/settlement"
	"splitwise-clone/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func GetGroupBalances(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("groupId")

	var member models.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ? AND is_active = ?", groupID, userID, true).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a member of this group"})
		return
	}

	var expenses []models.Expense
	database.DB.Where("group_id = ?", groupID).
		Preload("Participants").
		Find(&expenses)

	var settlements []models.Settlement
	database.DB.Where("group_id = ? AND is_paid = ?", groupID, true).
		Find(&settlements)

	balanceMap := make(map[uint]float64)

	for _, exp := range expenses {
		balanceMap[exp.PaidBy] += exp.Amount
		for _, p := range exp.Participants {
			balanceMap[p.UserID] -= p.Amount
		}
	}

	for _, s := range settlements {
		balanceMap[s.FromUserID] += s.Amount
		balanceMap[s.ToUserID] -= s.Amount
	}

	var balances []settlement.Balance
	for userID, bal := range balanceMap {
		balances = append(balances, settlement.Balance{
			UserID:  userID,
			Balance: utils.RoundFloat(bal, 2),
		})
	}

	c.JSON(http.StatusOK, balances)
}

func GetOptimalTransfers(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("groupId")

	var member models.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ? AND is_active = ?", groupID, userID, true).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a member of this group"})
		return
	}

	var expenses []models.Expense
	database.DB.Where("group_id = ?", groupID).
		Preload("Participants").
		Find(&expenses)

	var settlements []models.Settlement
	database.DB.Where("group_id = ? AND is_paid = ?", groupID, true).
		Find(&settlements)

	balanceMap := make(map[uint]float64)
	for _, exp := range expenses {
		balanceMap[exp.PaidBy] += exp.Amount
		for _, p := range exp.Participants {
			balanceMap[p.UserID] -= p.Amount
		}
	}

	for _, s := range settlements {
		balanceMap[s.FromUserID] += s.Amount
		balanceMap[s.ToUserID] -= s.Amount
	}

	var balances []settlement.Balance
	for uid, bal := range balanceMap {
		balances = append(balances, settlement.Balance{
			UserID:  uid,
			Balance: utils.RoundFloat(bal, 2),
		})
	}

	transfers := settlement.CalculateOptimalTransfers(balances)

	type TransferWithNames struct {
		FromUserID   uint    `json:"fromUserId"`
		FromUsername string  `json:"fromUsername"`
		ToUserID     uint    `json:"toUserId"`
		ToUsername   string  `json:"toUsername"`
		Amount       float64 `json:"amount"`
	}

	var result []TransferWithNames
	userCache := make(map[uint]string)

	for _, t := range transfers {
		if _, ok := userCache[t.From]; !ok {
			var user models.User
			database.DB.Select("username").First(&user, t.From)
			userCache[t.From] = user.Username
		}
		if _, ok := userCache[t.To]; !ok {
			var user models.User
			database.DB.Select("username").First(&user, t.To)
			userCache[t.To] = user.Username
		}

		result = append(result, TransferWithNames{
			FromUserID:   t.From,
			FromUsername: userCache[t.From],
			ToUserID:     t.To,
			ToUsername:   userCache[t.To],
			Amount:       t.Amount,
		})
	}

	c.JSON(http.StatusOK, result)
}

func CreateSettlement(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("groupId")

	var req struct {
		FromUserID uint    `json:"fromUserId" binding:"required"`
		ToUserID   uint    `json:"toUserId" binding:"required"`
		Amount     float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var member models.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ? AND is_active = ?", groupID, userID, true).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a member of this group"})
		return
	}

	settlementRecord := models.Settlement{
		GroupID:    parseUint(groupID),
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		Amount:     req.Amount,
		IsPaid:     false,
	}

	if err := database.DB.Create(&settlementRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create settlement"})
		return
	}

	database.DB.Preload("FromUser").Preload("ToUser").First(&settlementRecord, settlementRecord.ID)
	c.JSON(http.StatusCreated, settlementRecord)
}

func MarkSettlementPaid(c *gin.Context) {
	userID := c.GetUint("userID")
	settlementID := c.Param("id")

	var settlementRecord models.Settlement
	if err := database.DB.First(&settlementRecord, settlementID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Settlement not found"})
		return
	}

	if settlementRecord.FromUserID != userID && settlementRecord.ToUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not involved in this settlement"})
		return
	}

	now := time.Now()
	settlementRecord.IsPaid = true
	settlementRecord.PaidAt = &now

	if err := database.DB.Save(&settlementRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as paid"})
		return
	}

	c.JSON(http.StatusOK, settlementRecord)
}

func GetGroupSettlements(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("groupId")

	var member models.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ? AND is_active = ?", groupID, userID, true).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a member of this group"})
		return
	}

	var settlements []models.Settlement
	database.DB.Where("group_id = ?", groupID).
		Preload("FromUser").
		Preload("ToUser").
		Order("created_at DESC").
		Find(&settlements)

	c.JSON(http.StatusOK, settlements)
}
