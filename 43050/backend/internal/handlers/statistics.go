package handlers

import (
	"fmt"
	"net/http"
	"splitwise-clone/internal/database"
	"splitwise-clone/internal/models"
	"splitwise-clone/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StatisticsSummary struct {
	TotalPaid     float64 `json:"totalPaid"`
	TotalOwed     float64 `json:"totalOwed"`
	NetBalance    float64 `json:"netBalance"`
	ExpenseCount  int64   `json:"expenseCount"`
}

type MonthlyStats struct {
	Month     string  `json:"month"`
	TotalPaid float64 `json:"totalPaid"`
	TotalOwed float64 `json:"totalOwed"`
}

func GetUserStatistics(c *gin.Context) {
	userID := c.GetUint("userID")

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	groupID := c.Query("groupId")
	memberID := c.Query("memberId")

	type filter struct {
		startDate *time.Time
		endDate   *time.Time
		groupID   *uint
		memberID  *uint
	}

	var f filter

	if startDate != "" {
		if parsed, err := time.Parse("2006-01-02", startDate); err == nil {
			f.startDate = &parsed
		}
	}
	if endDate != "" {
		if parsed, err := time.Parse("2006-01-02", endDate); err == nil {
			f.endDate = &parsed
		}
	}
	if groupID != "" {
		var gid uint
		if _, err := fmt.Sscanf(groupID, "%d", &gid); err == nil {
			f.groupID = &gid
		}
	}
	if memberID != "" {
		var mid uint
		if _, err := fmt.Sscanf(memberID, "%d", &mid); err == nil {
			f.memberID = &mid
		}
	}

	buildExpenseConditions := func(query *gorm.DB) *gorm.DB {
		if f.startDate != nil {
			query = query.Where("expenses.expense_date >= ?", *f.startDate)
		}
		if f.endDate != nil {
			query = query.Where("expenses.expense_date <= ?", *f.endDate)
		}
		if f.groupID != nil {
			query = query.Where("expenses.group_id = ?", *f.groupID)
		}
		if f.memberID != nil {
			query = query.Where(`
				expenses.paid_by = ? OR 
				EXISTS (
					SELECT 1 FROM expense_participants ep 
					WHERE ep.expense_id = expenses.id AND ep.user_id = ?
				)
			`, *f.memberID, *f.memberID)
		}
		return query
	}

	var totalPaid float64
	paidQuery := database.DB.Model(&models.Expense{}).
		Where("paid_by = ?", userID)
	paidQuery = buildExpenseConditions(paidQuery)
	paidQuery.Select("COALESCE(SUM(amount), 0)").Scan(&totalPaid)

	var totalOwed float64
	owedQuery := database.DB.Model(&models.ExpenseParticipant{}).
		Joins("JOIN expenses ON expense_participants.expense_id = expenses.id").
		Where("expense_participants.user_id = ? AND expense_participants.user_id != expenses.paid_by", userID)
	owedQuery = buildExpenseConditions(owedQuery)
	owedQuery.Select("COALESCE(SUM(expense_participants.amount), 0)").Scan(&totalOwed)

	var expenseCount int64
	countQuery := database.DB.Model(&models.Expense{}).
		Joins("JOIN expense_participants ON expenses.id = expense_participants.expense_id").
		Where("expenses.paid_by = ? OR expense_participants.user_id = ?", userID, userID)
	countQuery = buildExpenseConditions(countQuery)
	countQuery.Distinct("expenses.id").Count(&expenseCount)

	summary := StatisticsSummary{
		TotalPaid:    utils.RoundFloat(totalPaid, 2),
		TotalOwed:    utils.RoundFloat(totalOwed, 2),
		NetBalance:   utils.RoundFloat(totalPaid-totalOwed, 2),
		ExpenseCount: expenseCount,
	}

	c.JSON(http.StatusOK, summary)
}

func GetUserMonthlyStats(c *gin.Context) {
	userID := c.GetUint("userID")

	type monthlyResult struct {
		Month     string
		TotalPaid float64
		TotalOwed float64
	}

	var results []monthlyResult

	database.DB.Raw(`
		SELECT 
			month,
			COALESCE(SUM(total_paid), 0) as total_paid,
			COALESCE(SUM(total_owed), 0) as total_owed
		FROM (
			SELECT
				strftime('%Y-%m', e.expense_date) as month,
				e.amount as total_paid,
				0 as total_owed
			FROM expenses e
			WHERE e.paid_by = ?
			
			UNION ALL
			
			SELECT
				strftime('%Y-%m', e.expense_date) as month,
				0 as total_paid,
				ep.amount as total_owed
			FROM expense_participants ep
			JOIN expenses e ON ep.expense_id = e.id
			WHERE ep.user_id = ? AND ep.user_id != e.paid_by
		) combined
		GROUP BY month
		ORDER BY month DESC
		LIMIT 12
	`, userID, userID).Scan(&results)

	var monthlyStats []MonthlyStats
	for _, r := range results {
		monthlyStats = append(monthlyStats, MonthlyStats{
			Month:     r.Month,
			TotalPaid: utils.RoundFloat(r.TotalPaid, 2),
			TotalOwed: utils.RoundFloat(r.TotalOwed, 2),
		})
	}

	c.JSON(http.StatusOK, monthlyStats)
}

func GetUserExpenseHistory(c *gin.Context) {
	userID := c.GetUint("userID")

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	groupID := c.Query("groupId")
	memberID := c.Query("memberId")

	db := database.DB.Model(&models.Expense{}).
		Joins("JOIN expense_participants ON expenses.id = expense_participants.expense_id").
		Where("expenses.paid_by = ? OR expense_participants.user_id = ?", userID, userID)

	if startDate != "" {
		if parsed, err := time.Parse("2006-01-02", startDate); err == nil {
			db = db.Where("expenses.expense_date >= ?", parsed)
		}
	}
	if endDate != "" {
		if parsed, err := time.Parse("2006-01-02", endDate); err == nil {
			db = db.Where("expenses.expense_date <= ?", parsed)
		}
	}
	if groupID != "" {
		db = db.Where("expenses.group_id = ?", groupID)
	}
	if memberID != "" {
		var mid uint
		if _, err := fmt.Sscanf(memberID, "%d", &mid); err == nil {
			db = db.Where(`
				expenses.paid_by = ? OR 
				EXISTS (
					SELECT 1 FROM expense_participants ep 
					WHERE ep.expense_id = expenses.id AND ep.user_id = ?
				)
			`, mid, mid)
		}
	}

	var expenses []models.Expense
	db.Distinct("expenses.id").
		Preload("Payer").
		Preload("Participants.User").
		Order("expenses.expense_date DESC").
		Find(&expenses)

	c.JSON(http.StatusOK, expenses)
}

func GetGroupExpenseStats(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("groupId")

	var member models.GroupMember
	if err := database.DB.Where("group_id = ? AND user_id = ? AND is_active = ?", groupID, userID, true).First(&member).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not a member of this group"})
		return
	}

	type MemberStats struct {
		UserID    uint    `json:"userId"`
		Username  string  `json:"username"`
		TotalPaid float64 `json:"totalPaid"`
		TotalOwed float64 `json:"totalOwed"`
		Balance   float64 `json:"balance"`
	}

	var results []struct {
		UserID    uint
		Username  string
		TotalPaid float64
		TotalOwed float64
	}

	database.DB.Raw(`
		SELECT
			u.id as user_id,
			u.username,
			COALESCE((
				SELECT SUM(e.amount)
				FROM expenses e
				WHERE e.group_id = gm.group_id AND e.paid_by = u.id
			), 0) as total_paid,
			COALESCE((
				SELECT SUM(ep.amount)
				FROM expense_participants ep
				JOIN expenses e ON ep.expense_id = e.id
				WHERE e.group_id = gm.group_id AND ep.user_id = u.id
			), 0) as total_owed
		FROM users u
		JOIN group_members gm ON u.id = gm.user_id
		WHERE gm.group_id = ? AND gm.is_active = 1
		GROUP BY u.id, u.username
	`, groupID).Scan(&results)

	var stats []MemberStats
	for _, r := range results {
		stats = append(stats, MemberStats{
			UserID:    r.UserID,
			Username:  r.Username,
			TotalPaid: utils.RoundFloat(r.TotalPaid, 2),
			TotalOwed: utils.RoundFloat(r.TotalOwed, 2),
			Balance:   utils.RoundFloat(r.TotalPaid-r.TotalOwed, 2),
		})
	}

	c.JSON(http.StatusOK, stats)
}
