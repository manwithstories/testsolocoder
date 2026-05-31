package handlers

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"housekeeping/database"
	"housekeeping/models"
	"housekeeping/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WithdrawalInput struct {
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	Account string  `json:"account" binding:"required"`
}

func RequestWithdrawal(c *gin.Context) {
	uid, _ := c.Get("uid")
	var in WithdrawalInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var w models.Wallet
		if err := tx.Where("user_id = ?", uid).First(&w).Error; err != nil {
			return err
		}
		if w.Balance < in.Amount {
			return fmt.Errorf("insufficient balance")
		}
		if err := tx.Model(&w).UpdateColumn("balance", gorm.Expr("balance - ?", in.Amount)).Error; err != nil {
			return err
		}
		wd := models.Withdrawal{
			StaffID: uid.(uint),
			Amount:  in.Amount,
			Account: in.Account,
			Status:  "pending",
		}
		return tx.Create(&wd).Error
	})
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, "ok")
}

func MyEarnings(c *gin.Context) {
	uid, _ := c.Get("uid")
	var settlements []models.Settlement
	database.DB.Where("staff_id = ?", uid).Order("id desc").Find(&settlements)
	var total float64
	for _, s := range settlements {
		total += s.StaffShare
	}
	var withdrawals []models.Withdrawal
	database.DB.Where("staff_id = ?", uid).Order("id desc").Find(&withdrawals)
	var wdTotal float64
	for _, w := range withdrawals {
		if w.Status == "paid" {
			wdTotal += w.Amount
		}
	}
	var wallet models.Wallet
	database.DB.Where("user_id = ?", uid).First(&wallet)
	utils.OK(c, gin.H{
		"total_earned": total,
		"withdrawn":    wdTotal,
		"balance":      wallet.Balance,
		"settlements":  settlements,
		"withdrawals":  withdrawals,
	})
}

func CompanyMonthlyStats(c *gin.Context) {
	uid, _ := c.Get("uid")
	month := c.DefaultQuery("month", time.Now().Format("2006-01"))
	start, _ := time.Parse("2006-01", month)
	end := start.AddDate(0, 1, 0)
	var settles []models.Settlement
	database.DB.Where("company_id = ? AND created_at >= ? AND created_at < ?", uid, start, end).Find(&settles)
	var income, payout float64
	for _, s := range settles {
		income += s.TotalAmount
		payout += s.StaffShare
	}
	var orderCount int64
	database.DB.Model(&models.Order{}).
		Where("company_id = ? AND status = ? AND created_at >= ? AND created_at < ?", uid, models.OrderPaid, start, end).
		Count(&orderCount)
	utils.OK(c, gin.H{
		"month":       month,
		"order_count": orderCount,
		"income":      income,
		"payout":      payout,
		"net":         income - payout,
		"settlements": settles,
	})
}

func ExportFinanceCSV(c *gin.Context) {
	uid, _ := c.Get("uid")
	role, _ := c.Get("role")
	var settles []models.Settlement
	q := database.DB
	switch role.(string) {
	case string(models.RoleCompany):
		q = q.Where("company_id = ?", uid)
	case string(models.RoleStaff):
		q = q.Where("staff_id = ?", uid)
	}
	q.Find(&settles)
	c.Writer.Header().Set("Content-Type", "text/csv")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=finance.csv")
	w := csv.NewWriter(c.Writer)
	w.Write([]string{"ID", "OrderID", "Total", "CompanyShare", "StaffShare", "Status", "CreatedAt"})
	for _, s := range settles {
		w.Write([]string{
			strconv.Itoa(int(s.ID)),
			strconv.Itoa(int(s.OrderID)),
			fmt.Sprintf("%.2f", s.TotalAmount),
			fmt.Sprintf("%.2f", s.CompanyShare),
			fmt.Sprintf("%.2f", s.StaffShare),
			s.Status,
			s.CreatedAt.Format(time.RFC3339),
		})
	}
	w.Flush()
}

func WalletInfo(c *gin.Context) {
	uid, _ := c.Get("uid")
	var w models.Wallet
	if err := database.DB.Where("user_id = ?", uid).First(&w).Error; err != nil {
		utils.NotFound(c, "wallet not found")
		return
	}
	utils.OK(c, w)
}
