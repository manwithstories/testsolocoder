package service

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"ticket-system/internal/database"
	"ticket-system/internal/dto"
	"ticket-system/internal/models"
)

type StatisticsService struct{}

func NewStatisticsService() *StatisticsService {
	return &StatisticsService{}
}

type SalesStatistics struct {
	TotalRevenue    float64 `json:"total_revenue"`
	TotalOrders     int64   `json:"total_orders"`
	TotalTickets    int64   `json:"total_tickets"`
	AvgOrderAmount  float64 `json:"avg_order_amount"`
}

type AreaSales struct {
	AreaID    uint64  `json:"area_id"`
	AreaName  string  `json:"area_name"`
	SoldCount int64   `json:"sold_count"`
	TotalSeats int64  `json:"total_seats"`
	Revenue   float64 `json:"revenue"`
}

type SeatHeatmapData struct {
	SeatID   uint64 `json:"seat_id"`
	SeatNo   string `json:"seat_no"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Sold     bool   `json:"sold"`
}

type AudienceProfile struct {
	MaleCount     int64            `json:"male_count"`
	FemaleCount   int64            `json:"female_count"`
	AgeDistribution map[string]int64 `json:"age_distribution"`
	LevelDistribution []map[string]interface{} `json:"level_distribution"`
}

func (s *StatisticsService) GetSalesStatistics(req *dto.StatisticsRequest) (*SalesStatistics, error) {
	var stats SalesStatistics

	query := database.DB.Model(&models.Order{}).Where("status = ?", models.OrderStatusPaid)

	if req.StartDate != "" {
		query = query.Where("pay_time >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("pay_time <= ?", req.EndDate+" 23:59:59")
	}
	if req.ShowID > 0 {
		query = query.Where("show_id = ?", req.ShowID)
	}
	if req.SessionID > 0 {
		query = query.Where("session_id = ?", req.SessionID)
	}

	var result struct {
		TotalRevenue float64
		TotalOrders  int64
	}
	query.Select("COALESCE(SUM(pay_amount), 0) as total_revenue, COUNT(*) as total_orders").
		Scan(&result)

	stats.TotalRevenue = result.TotalRevenue
	stats.TotalOrders = result.TotalOrders

	if result.TotalOrders > 0 {
		stats.AvgOrderAmount = result.TotalRevenue / float64(result.TotalOrders)
	}

	ticketQuery := database.DB.Model(&models.Ticket{}).
		Joins("JOIN orders ON tickets.order_id = orders.id").
		Where("orders.status = ?", models.OrderStatusPaid)

	if req.ShowID > 0 {
		ticketQuery = ticketQuery.Where("tickets.show_id = ?", req.ShowID)
	}
	if req.SessionID > 0 {
		ticketQuery = ticketQuery.Where("tickets.session_id = ?", req.SessionID)
	}
	if req.AreaID > 0 {
		ticketQuery = ticketQuery.Where("tickets.area_id = ?", req.AreaID)
	}
	ticketQuery.Count(&stats.TotalTickets)

	return &stats, nil
}

func (s *StatisticsService) GetAreaSales(sessionID uint64) ([]AreaSales, error) {
	var areas []AreaSales

	database.DB.Table("seat_areas").
		Select("seat_areas.id as area_id, seat_areas.name as area_name, "+
			"COUNT(tickets.id) as sold_count, seat_areas.total_seats, "+
			"COALESCE(SUM(tickets.price), 0) as revenue").
		Joins("LEFT JOIN tickets ON tickets.area_id = seat_areas.id AND tickets.status = ?", models.TicketStatusValid).
		Where("seat_areas.session_id = ?", sessionID).
		Group("seat_areas.id").
		Scan(&areas)

	return areas, nil
}

func (s *StatisticsService) GetSeatHeatmap(sessionID uint64) ([]SeatHeatmapData, error) {
	var data []SeatHeatmapData

	database.DB.Table("seats").
		Select("seats.id as seat_id, seats.seat_no, seats.x, seats.y, "+
			"CASE WHEN seats.status = ? THEN true ELSE false END as sold", models.SeatStatusSold).
		Where("seats.session_id = ?", sessionID).
		Scan(&data)

	return data, nil
}

func (s *StatisticsService) GetAudienceProfile(req *dto.StatisticsRequest) (*AudienceProfile, error) {
	var profile AudienceProfile
	profile.AgeDistribution = make(map[string]int64)

	query := database.DB.Model(&models.User{}).
		Joins("JOIN orders ON users.id = orders.user_id").
		Where("orders.status = ?", models.OrderStatusPaid)

	if req.StartDate != "" {
		query = query.Where("orders.pay_time >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("orders.pay_time <= ?", req.EndDate+" 23:59:59")
	}
	if req.ShowID > 0 {
		query = query.Where("orders.show_id = ?", req.ShowID)
	}

	var users []models.User
	query.Distinct("users.id").Select("users.id_card, users.member_level").Find(&users)

	for _, user := range users {
		if user.IDCard != "" && len(user.IDCard) == 18 {
			genderDigit := int(user.IDCard[16] - '0')
			if genderDigit%2 == 1 {
				profile.MaleCount++
			} else {
				profile.FemaleCount++
			}

			age := GetAgeFromIDCard(user.IDCard)
			ageGroup := getAgeGroup(age)
			profile.AgeDistribution[ageGroup]++
		}
	}

	levelCounts := make(map[int]int64)
	for _, user := range users {
		levelCounts[user.MemberLevel]++
	}

	levelNames := map[int]string{
		1: "普通会员",
		2: "银卡会员",
		3: "金卡会员",
		4: "铂金会员",
	}
	for level, count := range levelCounts {
		profile.LevelDistribution = append(profile.LevelDistribution, map[string]interface{}{
			"level": level,
			"name":  levelNames[level],
			"count": count,
		})
	}

	return &profile, nil
}

func getAgeGroup(age int) string {
	switch {
	case age < 18:
		return "18岁以下"
	case age >= 18 && age < 25:
		return "18-24岁"
	case age >= 25 && age < 35:
		return "25-34岁"
	case age >= 35 && age < 45:
		return "35-44岁"
	case age >= 45 && age < 55:
		return "45-54岁"
	default:
		return "55岁以上"
	}
}

func GetAgeFromIDCard(idCard string) int {
	if len(idCard) != 18 {
		return 0
	}
	birthYear := int(idCard[6]-'0')*1000 + int(idCard[7]-'0')*100 + int(idCard[8]-'0')*10 + int(idCard[9]-'0')
	currentYear := 2026
	return currentYear - birthYear
}

func (s *StatisticsService) GetDailySales(startDate, endDate string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	rows, err := database.DB.Model(&models.Order{}).
		Select("DATE(pay_time) as date, COALESCE(SUM(pay_amount), 0) as revenue, COUNT(*) as orders").
		Where("status = ? AND pay_time >= ? AND pay_time <= ?",
			models.OrderStatusPaid, startDate, endDate+" 23:59:59").
		Group("DATE(pay_time)").
		Order("date asc").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		var revenue float64
		var orders int64
		rows.Scan(&date, &revenue, &orders)
		results = append(results, map[string]interface{}{
			"date":    date,
			"revenue": revenue,
			"orders":  orders,
		})
	}

	return results, nil
}

func (s *StatisticsService) GenerateStatisticsPDF(req *dto.StatisticsRequest) ([]byte, error) {
	salesStats, err := s.GetSalesStatistics(req)
	if err != nil {
		return nil, err
	}

	dailySales, err := s.GetDailySales(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	audienceProfile, err := s.GetAudienceProfile(req)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(40, 10, "演唱会票务系统 - 销售统计报告")
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 10)
	dateRange := "全部日期"
	if req.StartDate != "" && req.EndDate != "" {
		dateRange = fmt.Sprintf("%s 至 %s", req.StartDate, req.EndDate)
	}
	pdf.Cell(40, 8, "统计日期: "+dateRange)
	pdf.Ln(10)
	pdf.Cell(40, 8, "生成时间: "+time.Now().Format("2006-01-02 15:04:05"))
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "一、销售概览")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	overviewData := [][]string{
		{"指标", "数值"},
		{"总销售额", fmt.Sprintf("¥%.2f", salesStats.TotalRevenue)},
		{"总订单数", fmt.Sprintf("%d", salesStats.TotalOrders)},
		{"售票总数", fmt.Sprintf("%d", salesStats.TotalTickets)},
		{"平均订单金额", fmt.Sprintf("¥%.2f", salesStats.AvgOrderAmount)},
	}

	for _, row := range overviewData {
		pdf.CellFormat(60, 8, row[0], "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 8, row[1], "1", 1, "", false, 0, "")
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "二、日销售趋势")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	if len(dailySales) > 0 {
		pdf.CellFormat(40, 8, "日期", "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, "销售额", "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, "订单数", "1", 1, "C", false, 0, "")

		for _, item := range dailySales {
			dateStr, _ := item["date"].(string)
			revenue, _ := item["revenue"].(float64)
			orders, _ := item["orders"].(int64)

			pdf.CellFormat(40, 8, dateStr, "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 8, fmt.Sprintf("¥%.2f", revenue), "1", 0, "", false, 0, "")
			pdf.CellFormat(40, 8, fmt.Sprintf("%d", orders), "1", 1, "", false, 0, "")
		}
	} else {
		pdf.Cell(40, 8, "暂无销售数据")
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "三、观众画像")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(40, 8, "指标", "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 8, "数值", "1", 1, "", false, 0, "")
	pdf.CellFormat(40, 8, "男性观众", "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 8, fmt.Sprintf("%d人", audienceProfile.MaleCount), "1", 1, "", false, 0, "")
	pdf.CellFormat(40, 8, "女性观众", "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 8, fmt.Sprintf("%d人", audienceProfile.FemaleCount), "1", 1, "", false, 0, "")
	pdf.Ln(5)

	pdf.Cell(40, 8, "年龄分布:")
	pdf.Ln(8)
	for ageGroup, count := range audienceProfile.AgeDistribution {
		pdf.CellFormat(60, 8, ageGroup, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("%d人", count), "1", 1, "", false, 0, "")
	}
	pdf.Ln(5)

	pdf.Cell(40, 8, "会员等级分布:")
	pdf.Ln(8)
	for _, level := range audienceProfile.LevelDistribution {
		name, _ := level["name"].(string)
		count, _ := level["count"].(int64)
		pdf.CellFormat(60, 8, name, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("%d人", count), "1", 1, "", false, 0, "")
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
