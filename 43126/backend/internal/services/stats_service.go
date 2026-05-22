package services

import (
	"fmt"
	"time"

	"meeting-room/internal/models"
	"meeting-room/internal/repositories"

	"github.com/xuri/excelize/v2"
)

type StatsService struct {
	bookingRepo *repositories.BookingRepository
	roomRepo    *repositories.RoomRepository
}

func NewStatsService() *StatsService {
	return &StatsService{
		bookingRepo: repositories.NewBookingRepository(),
		roomRepo:    repositories.NewRoomRepository(),
	}
}

type UtilizationStats struct {
	RoomID     uint    `json:"room_id"`
	RoomName   string  `json:"room_name"`
	Floor      string  `json:"floor"`
	TotalHours float64 `json:"total_hours"`
	Bookings   int     `json:"bookings"`
	UtilRate   float64 `json:"util_rate"`
	Revenue    float64 `json:"revenue"`
}

type HourlyStats struct {
	Hour    string `json:"hour"`
	Count   int    `json:"count"`
	Percent float64 `json:"percent"`
}

type RevenueStats struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Bookings int    `json:"bookings"`
}

type StatsResponse struct {
	Utilization    []UtilizationStats `json:"utilization"`
	HourlyPopular  []HourlyStats      `json:"hourly_popular"`
	RevenueTrend   []RevenueStats     `json:"revenue_trend"`
	Summary        StatsSummary       `json:"summary"`
}

type StatsSummary struct {
	TotalBookings int     `json:"total_bookings"`
	TotalHours    float64 `json:"total_hours"`
	TotalRevenue  float64 `json:"total_revenue"`
	UtilRate      float64 `json:"util_rate"`
}

func (s *StatsService) GetStats(startTime, endTime time.Time, department string) (*StatsResponse, error) {
	bookings, err := s.bookingRepo.GetStats(startTime, endTime, department)
	if err != nil {
		return nil, err
	}

	rooms, _ := s.roomRepo.ListAll()

	response := &StatsResponse{
		Utilization:   s.calculateUtilization(bookings, rooms, startTime, endTime),
		HourlyPopular: s.calculateHourlyPopular(bookings),
		RevenueTrend:  s.calculateRevenueTrend(bookings, startTime, endTime),
		Summary:       s.calculateSummary(bookings, startTime, endTime),
	}

	return response, nil
}

func (s *StatsService) calculateUtilization(bookings []models.Booking, rooms []models.Room, startTime, endTime time.Time) []UtilizationStats {
	roomStats := make(map[uint]*UtilizationStats)
	roomMap := make(map[uint]models.Room)

	for _, room := range rooms {
		roomMap[room.ID] = room
		roomStats[room.ID] = &UtilizationStats{
			RoomID:   room.ID,
			RoomName: room.Name,
			Floor:    room.Floor,
		}
	}

	for _, booking := range bookings {
		stats, exists := roomStats[booking.RoomID]
		if !exists {
			roomName := "未知"
			floor := ""
			if room, ok := roomMap[booking.RoomID]; ok {
				roomName = room.Name
				floor = room.Floor
			}
			stats = &UtilizationStats{
				RoomID:   booking.RoomID,
				RoomName: roomName,
				Floor:    floor,
			}
			roomStats[booking.RoomID] = stats
		}

		hours := booking.EndTime.Sub(booking.StartTime).Hours()
		stats.TotalHours += hours
		stats.Bookings++
		stats.Revenue += booking.TotalPrice
	}

	days := endTime.Sub(startTime).Hours() / 24
	if days < 1 {
		days = 1
	}
	totalAvailableHours := days * 14

	result := make([]UtilizationStats, 0, len(roomStats))
	for _, stats := range roomStats {
		if totalAvailableHours > 0 {
			stats.UtilRate = float64(int(stats.TotalHours/totalAvailableHours*10000)) / 100
		}
		result = append(result, *stats)
	}

	return result
}

func (s *StatsService) calculateHourlyPopular(bookings []models.Booking) []HourlyStats {
	hourlyCount := make(map[string]int)
	total := 0

	for _, booking := range bookings {
		startHour := booking.StartTime.Hour()
		endHour := booking.EndTime.Hour()
		if booking.EndTime.Minute() > 0 {
			endHour++
		}
		for hour := startHour; hour < endHour; hour++ {
			key := fmt.Sprintf("%02d:00", hour)
			hourlyCount[key]++
			total++
		}
	}

	result := make([]HourlyStats, 0, len(hourlyCount))
	for hour, count := range hourlyCount {
		percent := 0.0
		if total > 0 {
			percent = float64(int(float64(count)/float64(total)*10000)) / 100
		}
		result = append(result, HourlyStats{
			Hour:    hour,
			Count:   count,
			Percent: percent,
		})
	}

	return result
}

func (s *StatsService) calculateRevenueTrend(bookings []models.Booking, startTime, endTime time.Time) []RevenueStats {
	dailyRevenue := make(map[string]*RevenueStats)

	for d := startTime; d.Before(endTime); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		dailyRevenue[dateStr] = &RevenueStats{
			Date: dateStr,
		}
	}

	for _, booking := range bookings {
		dateStr := booking.StartTime.Format("2006-01-02")
		stats, exists := dailyRevenue[dateStr]
		if !exists {
			stats = &RevenueStats{Date: dateStr}
			dailyRevenue[dateStr] = stats
		}
		stats.Revenue += booking.TotalPrice
		stats.Bookings++
	}

	result := make([]RevenueStats, 0, len(dailyRevenue))
	for _, stats := range dailyRevenue {
		result = append(result, *stats)
	}

	return result
}

func (s *StatsService) calculateSummary(bookings []models.Booking, startTime, endTime time.Time) StatsSummary {
	var summary StatsSummary
	summary.TotalBookings = len(bookings)

	for _, booking := range bookings {
		hours := booking.EndTime.Sub(booking.StartTime).Hours()
		summary.TotalHours += hours
		summary.TotalRevenue += booking.TotalPrice
	}

	days := endTime.Sub(startTime).Hours() / 24
	if days < 1 {
		days = 1
	}
	rooms, _ := s.roomRepo.ListAll()
	totalAvailable := days * 14 * float64(len(rooms))
	if totalAvailable > 0 {
		summary.UtilRate = float64(int(summary.TotalHours/totalAvailable*10000)) / 100
	}

	return summary
}

func (s *StatsService) ExportToExcel(startTime, endTime time.Time, department string) (string, error) {
	stats, err := s.GetStats(startTime, endTime, department)
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()

	sheet1 := "利用率统计"
	f.SetSheetName("Sheet1", sheet1)

	f.SetCellValue(sheet1, "A1", "会议室名称")
	f.SetCellValue(sheet1, "B1", "楼层")
	f.SetCellValue(sheet1, "C1", "使用时长(小时)")
	f.SetCellValue(sheet1, "D1", "预订次数")
	f.SetCellValue(sheet1, "E1", "利用率(%)")
	f.SetCellValue(sheet1, "F1", "收入(元)")

	for i, stat := range stats.Utilization {
		row := i + 2
		f.SetCellValue(sheet1, fmt.Sprintf("A%d", row), stat.RoomName)
		f.SetCellValue(sheet1, fmt.Sprintf("B%d", row), stat.Floor)
		f.SetCellValue(sheet1, fmt.Sprintf("C%d", row), stat.TotalHours)
		f.SetCellValue(sheet1, fmt.Sprintf("D%d", row), stat.Bookings)
		f.SetCellValue(sheet1, fmt.Sprintf("E%d", row), stat.UtilRate)
		f.SetCellValue(sheet1, fmt.Sprintf("F%d", row), stat.Revenue)
	}

	sheet2 := "热门时段"
	f.NewSheet(sheet2)
	f.SetCellValue(sheet2, "A1", "时段")
	f.SetCellValue(sheet2, "B1", "使用次数")
	f.SetCellValue(sheet2, "C1", "占比(%)")

	for i, stat := range stats.HourlyPopular {
		row := i + 2
		f.SetCellValue(sheet2, fmt.Sprintf("A%d", row), stat.Hour)
		f.SetCellValue(sheet2, fmt.Sprintf("B%d", row), stat.Count)
		f.SetCellValue(sheet2, fmt.Sprintf("C%d", row), stat.Percent)
	}

	sheet3 := "收入趋势"
	f.NewSheet(sheet3)
	f.SetCellValue(sheet3, "A1", "日期")
	f.SetCellValue(sheet3, "B1", "收入(元)")
	f.SetCellValue(sheet3, "C1", "预订次数")

	for i, stat := range stats.RevenueTrend {
		row := i + 2
		f.SetCellValue(sheet3, fmt.Sprintf("A%d", row), stat.Date)
		f.SetCellValue(sheet3, fmt.Sprintf("B%d", row), stat.Revenue)
		f.SetCellValue(sheet3, fmt.Sprintf("C%d", row), stat.Bookings)
	}

	sheet4 := "汇总"
	f.NewSheet(sheet4)
	f.SetCellValue(sheet4, "A1", "统计项")
	f.SetCellValue(sheet4, "B1", "数值")
	f.SetCellValue(sheet4, "A2", "总预订次数")
	f.SetCellValue(sheet4, "B2", stats.Summary.TotalBookings)
	f.SetCellValue(sheet4, "A3", "总使用时长(小时)")
	f.SetCellValue(sheet4, "B3", stats.Summary.TotalHours)
	f.SetCellValue(sheet4, "A4", "总收入(元)")
	f.SetCellValue(sheet4, "B4", stats.Summary.TotalRevenue)
	f.SetCellValue(sheet4, "A5", "整体利用率(%)")
	f.SetCellValue(sheet4, "B5", stats.Summary.UtilRate)

	filename := fmt.Sprintf("stats_%s_%s.xlsx", startTime.Format("20060102"), endTime.Format("20060102"))
	filepath := "./uploads/" + filename

	if err := f.SaveAs(filepath); err != nil {
		return "", err
	}

	return filename, nil
}
