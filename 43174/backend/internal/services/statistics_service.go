package services

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"campus-trade-platform/internal/models"
)

type StatisticsService struct {
	db *gorm.DB
}

func NewStatisticsService(db *gorm.DB) *StatisticsService {
	return &StatisticsService{db: db}
}

type TextbookStats struct {
	TotalCount     int64   `json:"total_count"`
	AvailableCount int64   `json:"available_count"`
	SoldCount      int64   `json:"sold_count"`
	TotalValue     float64 `json:"total_value"`
}

type UserStats struct {
	TotalCount     int64 `json:"total_count"`
	StudentCount   int64 `json:"student_count"`
	MerchantCount  int64 `json:"merchant_count"`
	PendingCount   int64 `json:"pending_count"`
}

type OrderStats struct {
	TotalCount    int64   `json:"total_count"`
	TotalRevenue  float64 `json:"total_revenue"`
	PendingCount  int64   `json:"pending_count"`
	CompletedCount int64  `json:"completed_count"`
}

type MonthlyStats struct {
	Month         string  `json:"month"`
	OrderCount    int64   `json:"order_count"`
	Revenue       float64 `json:"revenue"`
	NewUsers      int64   `json:"new_users"`
	NewTextbooks  int64   `json:"new_textbooks"`
}

func (s *StatisticsService) GetTextbookStats() (*TextbookStats, error) {
	var stats TextbookStats

	s.db.Model(&models.Textbook{}).Count(&stats.TotalCount)
	s.db.Model(&models.Textbook{}).Where("status = ?", models.TextbookStatusAvailable).Count(&stats.AvailableCount)
	s.db.Model(&models.Textbook{}).Where("status = ?", models.TextbookStatusSold).Count(&stats.SoldCount)
	s.db.Model(&models.Textbook{}).Where("status = ?", models.TextbookStatusSold).
		Select("COALESCE(SUM(price), 0)").Scan(&stats.TotalValue)

	return &stats, nil
}

func (s *StatisticsService) GetUserStats() (*UserStats, error) {
	var stats UserStats

	s.db.Model(&models.User{}).Count(&stats.TotalCount)
	s.db.Model(&models.User{}).Where("role = ?", models.RoleStudent).Count(&stats.StudentCount)
	s.db.Model(&models.User{}).Where("role = ?", models.RoleMerchant).Count(&stats.MerchantCount)
	s.db.Model(&models.User{}).Where("status = ?", models.UserStatusPending).Count(&stats.PendingCount)

	return &stats, nil
}

func (s *StatisticsService) GetOrderStats() (*OrderStats, error) {
	var stats OrderStats

	s.db.Model(&models.Order{}).Count(&stats.TotalCount)
	s.db.Model(&models.Order{}).Where("status = ?", models.OrderStatusPending).Count(&stats.PendingCount)
	s.db.Model(&models.Order{}).Where("status = ?", models.OrderStatusCompleted).Count(&stats.CompletedCount)
	s.db.Model(&models.Order{}).Where("status = ?", models.OrderStatusCompleted).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.TotalRevenue)

	return &stats, nil
}

func (s *StatisticsService) GetPopularTextbooks(limit int) ([]models.Textbook, error) {
	var textbooks []models.Textbook
	err := s.db.Preload("Seller").Preload("Category").
		Where("status = ?", models.TextbookStatusAvailable).
		Order("view_count DESC").
		Limit(limit).
		Find(&textbooks).Error
	return textbooks, err
}

func (s *StatisticsService) GetTopUsers(limit int) ([]models.User, error) {
	var users []models.User
	err := s.db.Where("role IN ?", []string{"student", "merchant"}).
		Order("rating DESC").
		Limit(limit).
		Find(&users).Error
	return users, err
}

func (s *StatisticsService) GetMonthlyStats(months int) ([]MonthlyStats, error) {
	var result []MonthlyStats

	currentDate := time.Now()

	for i := months - 1; i >= 0; i-- {
		monthStart := time.Date(currentDate.Year(), currentDate.Month()-time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := monthStart.AddDate(0, 1, 0)

		var stat MonthlyStats
		stat.Month = monthStart.Format("2006-01")

		s.db.Model(&models.Order{}).
			Where("created_at >= ? AND created_at < ?", monthStart, monthEnd).
			Count(&stat.OrderCount)

		s.db.Model(&models.Order{}).
			Where("created_at >= ? AND created_at < ? AND status = ?", monthStart, monthEnd, models.OrderStatusCompleted).
			Select("COALESCE(SUM(total_amount), 0)").Scan(&stat.Revenue)

		s.db.Model(&models.User{}).
			Where("created_at >= ? AND created_at < ?", monthStart, monthEnd).
			Count(&stat.NewUsers)

		s.db.Model(&models.Textbook{}).
			Where("created_at >= ? AND created_at < ?", monthStart, monthEnd).
			Count(&stat.NewTextbooks)

		result = append(result, stat)
	}

	return result, nil
}

func (s *StatisticsService) ExportMonthlyReport(month string) (string, error) {
	f := excelize.NewFile()

	ordersSheet := "Orders"
	f.SetSheetName("Sheet1", ordersSheet)

	headers := []string{"Order No", "Buyer", "Seller", "Amount", "Status", "Created At"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(ordersSheet, cell, header)
	}

	var orders []models.Order
	s.db.Preload("Buyer").Preload("Seller").
		Where("to_char(created_at, 'YYYY-MM') = ?", month).
		Find(&orders)

	for i, order := range orders {
		row := i + 2
		f.SetCellValue(ordersSheet, fmt.Sprintf("A%d", row), order.OrderNo)
		if order.Buyer != nil {
			f.SetCellValue(ordersSheet, fmt.Sprintf("B%d", row), order.Buyer.Username)
		}
		if order.Seller != nil {
			f.SetCellValue(ordersSheet, fmt.Sprintf("C%d", row), order.Seller.Username)
		}
		f.SetCellValue(ordersSheet, fmt.Sprintf("D%d", row), order.TotalAmount)
		f.SetCellValue(ordersSheet, fmt.Sprintf("E%d", row), order.Status)
		f.SetCellValue(ordersSheet, fmt.Sprintf("F%d", row), order.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	usersSheet := "New Users"
	f.NewSheet(usersSheet)

	userHeaders := []string{"Username", "Email", "Role", "Status", "Registered At"}
	for i, header := range userHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(usersSheet, cell, header)
	}

	var users []models.User
	s.db.Where("to_char(created_at, 'YYYY-MM') = ?", month).
		Find(&users)

	for i, user := range users {
		row := i + 2
		f.SetCellValue(usersSheet, fmt.Sprintf("A%d", row), user.Username)
		f.SetCellValue(usersSheet, fmt.Sprintf("B%d", row), user.Email)
		f.SetCellValue(usersSheet, fmt.Sprintf("C%d", row), user.Role)
		f.SetCellValue(usersSheet, fmt.Sprintf("D%d", row), user.Status)
		f.SetCellValue(usersSheet, fmt.Sprintf("E%d", row), user.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	textbooksSheet := "Textbooks"
	f.NewSheet(textbooksSheet)

	bookHeaders := []string{"Title", "ISBN", "Price", "Status", "Seller", "Listed At"}
	for i, header := range bookHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(textbooksSheet, cell, header)
	}

	var textbooks []models.Textbook
	s.db.Preload("Seller").
		Where("to_char(created_at, 'YYYY-MM') = ?", month).
		Find(&textbooks)

	for i, textbook := range textbooks {
		row := i + 2
		f.SetCellValue(textbooksSheet, fmt.Sprintf("A%d", row), textbook.Title)
		f.SetCellValue(textbooksSheet, fmt.Sprintf("B%d", row), textbook.ISBN)
		f.SetCellValue(textbooksSheet, fmt.Sprintf("C%d", row), textbook.Price)
		f.SetCellValue(textbooksSheet, fmt.Sprintf("D%d", row), textbook.Status)
		if textbook.Seller != nil {
			f.SetCellValue(textbooksSheet, fmt.Sprintf("E%d", row), textbook.Seller.Username)
		}
		f.SetCellValue(textbooksSheet, fmt.Sprintf("F%d", row), textbook.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	fileName := fmt.Sprintf("monthly_report_%s.xlsx", month)
	filePath := fmt.Sprintf("./exports/%s", fileName)

	if err := f.SaveAs(filePath); err != nil {
		return "", fmt.Errorf("failed to save report: %w", err)
	}

	return filePath, nil
}
