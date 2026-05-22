package service

import (
	"bytes"
	"errors"
	"fmt"
	"ticket-system/internal/dto"
	"ticket-system/internal/models"
	"time"

	"github.com/xuri/excelize/v2"
)

type StatisticsService struct{}

func NewStatisticsService() *StatisticsService {
	return &StatisticsService{}
}

func (s *StatisticsService) GetActivityStatistics(req *dto.StatisticsRequest) ([]dto.ActivityStatistics, error) {
	var results []struct {
		ActivityID    uint
		ActivityTitle string
		TotalOrders   int64
		TotalAmount   float64
		TotalTickets  int64
	}

	query := models.DB.Table("orders").
		Select("orders.activity_id, activities.title as activity_title, " +
			"COUNT(DISTINCT orders.id) as total_orders, " +
			"SUM(orders.pay_amount) as total_amount, " +
			"SUM(order_items.quantity) as total_tickets").
		Joins("JOIN activities ON activities.id = orders.activity_id").
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Where("orders.status = ?", models.OrderStatusPaid)

	if req.StartDate != "" {
		startTime, err := time.ParseInLocation("2006-01-02", req.StartDate, time.Local)
		if err == nil {
			query = query.Where("orders.created_at >= ?", startTime)
		}
	}

	if req.EndDate != "" {
		endTime, err := time.ParseInLocation("2006-01-02", req.EndDate, time.Local)
		if err == nil {
			query = query.Where("orders.created_at <= ?", endTime.Add(24*time.Hour))
		}
	}

	if req.ActivityID > 0 {
		query = query.Where("orders.activity_id = ?", req.ActivityID)
	}

	if err := query.Group("orders.activity_id, activities.title").Find(&results).Error; err != nil {
		return nil, err
	}

	var stats []dto.ActivityStatistics
	for _, r := range results {
		stats = append(stats, dto.ActivityStatistics{
			ActivityID:    r.ActivityID,
			ActivityTitle: r.ActivityTitle,
			TotalOrders:   r.TotalOrders,
			TotalAmount:   r.TotalAmount,
			TotalTickets:  r.TotalTickets,
		})
	}

	return stats, nil
}

func (s *StatisticsService) GetTicketTypeStatistics(req *dto.StatisticsRequest) ([]dto.TicketTypeStatistics, error) {
	var results []struct {
		TicketTypeID   uint
		TicketTypeName string
		TotalSold      int64
		TotalAmount    float64
	}

	query := models.DB.Table("order_items").
		Select("order_items.ticket_type_id, ticket_types.name as ticket_type_name, " +
			"SUM(order_items.quantity) as total_sold, " +
			"SUM(order_items.subtotal) as total_amount").
		Joins("JOIN ticket_types ON ticket_types.id = order_items.ticket_type_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.status = ?", models.OrderStatusPaid)

	if req.StartDate != "" {
		startTime, err := time.ParseInLocation("2006-01-02", req.StartDate, time.Local)
		if err == nil {
			query = query.Where("orders.created_at >= ?", startTime)
		}
	}

	if req.EndDate != "" {
		endTime, err := time.ParseInLocation("2006-01-02", req.EndDate, time.Local)
		if err == nil {
			query = query.Where("orders.created_at <= ?", endTime.Add(24*time.Hour))
		}
	}

	if req.ActivityID > 0 {
		query = query.Where("ticket_types.activity_id = ?", req.ActivityID)
	}

	if err := query.Group("order_items.ticket_type_id, ticket_types.name").Find(&results).Error; err != nil {
		return nil, err
	}

	var stats []dto.TicketTypeStatistics
	for _, r := range results {
		stats = append(stats, dto.TicketTypeStatistics{
			TicketTypeID:   r.TicketTypeID,
			TicketTypeName: r.TicketTypeName,
			TotalSold:      r.TotalSold,
			TotalAmount:    r.TotalAmount,
		})
	}

	return stats, nil
}

func (s *StatisticsService) GetDailyStatistics(req *dto.StatisticsRequest) ([]dto.DailyStatistics, error) {
	var results []struct {
		Date         string
		TotalOrders  int64
		TotalAmount  float64
		TotalTickets int64
	}

	query := models.DB.Table("orders").
		Select("DATE(orders.created_at) as date, " +
			"COUNT(DISTINCT orders.id) as total_orders, " +
			"SUM(orders.pay_amount) as total_amount, " +
			"SUM(order_items.quantity) as total_tickets").
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Where("orders.status = ?", models.OrderStatusPaid)

	if req.StartDate != "" {
		startTime, err := time.ParseInLocation("2006-01-02", req.StartDate, time.Local)
		if err == nil {
			query = query.Where("orders.created_at >= ?", startTime)
		}
	}

	if req.EndDate != "" {
		endTime, err := time.ParseInLocation("2006-01-02", req.EndDate, time.Local)
		if err == nil {
			query = query.Where("orders.created_at <= ?", endTime.Add(24*time.Hour))
		}
	}

	if req.ActivityID > 0 {
		query = query.Where("orders.activity_id = ?", req.ActivityID)
	}

	if err := query.Group("DATE(orders.created_at)").Order("date DESC").Find(&results).Error; err != nil {
		return nil, err
	}

	var stats []dto.DailyStatistics
	for _, r := range results {
		stats = append(stats, dto.DailyStatistics{
			Date:         r.Date,
			TotalOrders:  r.TotalOrders,
			TotalAmount:  r.TotalAmount,
			TotalTickets: r.TotalTickets,
		})
	}

	return stats, nil
}

func (s *StatisticsService) ExportExcel(req *dto.StatisticsRequest) ([]byte, error) {
	activityStats, err := s.GetActivityStatistics(req)
	if err != nil {
		return nil, errors.New("获取活动统计失败")
	}

	ticketStats, err := s.GetTicketTypeStatistics(req)
	if err != nil {
		return nil, errors.New("获取票型统计失败")
	}

	dailyStats, err := s.GetDailyStatistics(req)
	if err != nil {
		return nil, errors.New("获取每日统计失败")
	}

	f := excelize.NewFile()

	sheet1 := "活动统计"
	f.SetSheetName("Sheet1", sheet1)
	headers1 := []string{"活动ID", "活动名称", "订单数", "收入(元)", "售票数"}
	for i, h := range headers1 {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet1, cell, h)
	}
	for i, stat := range activityStats {
		row := i + 2
		f.SetCellValue(sheet1, fmt.Sprintf("A%d", row), stat.ActivityID)
		f.SetCellValue(sheet1, fmt.Sprintf("B%d", row), stat.ActivityTitle)
		f.SetCellValue(sheet1, fmt.Sprintf("C%d", row), stat.TotalOrders)
		f.SetCellValue(sheet1, fmt.Sprintf("D%d", row), stat.TotalAmount)
		f.SetCellValue(sheet1, fmt.Sprintf("E%d", row), stat.TotalTickets)
	}

	sheet2 := "票型统计"
	f.NewSheet(sheet2)
	headers2 := []string{"票型ID", "票型名称", "售出数量", "收入(元)"}
	for i, h := range headers2 {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet2, cell, h)
	}
	for i, stat := range ticketStats {
		row := i + 2
		f.SetCellValue(sheet2, fmt.Sprintf("A%d", row), stat.TicketTypeID)
		f.SetCellValue(sheet2, fmt.Sprintf("B%d", row), stat.TicketTypeName)
		f.SetCellValue(sheet2, fmt.Sprintf("C%d", row), stat.TotalSold)
		f.SetCellValue(sheet2, fmt.Sprintf("D%d", row), stat.TotalAmount)
	}

	sheet3 := "每日统计"
	f.NewSheet(sheet3)
	headers3 := []string{"日期", "订单数", "收入(元)", "售票数"}
	for i, h := range headers3 {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet3, cell, h)
	}
	for i, stat := range dailyStats {
		row := i + 2
		f.SetCellValue(sheet3, fmt.Sprintf("A%d", row), stat.Date)
		f.SetCellValue(sheet3, fmt.Sprintf("B%d", row), stat.TotalOrders)
		f.SetCellValue(sheet3, fmt.Sprintf("C%d", row), stat.TotalAmount)
		f.SetCellValue(sheet3, fmt.Sprintf("D%d", row), stat.TotalTickets)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
