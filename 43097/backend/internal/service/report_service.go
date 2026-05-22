package service

import (
	"errors"
	"hotel-system/internal/dto"
	"hotel-system/internal/model"
	"hotel-system/internal/pkg/logger"
	"hotel-system/pkg/excel"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ReportService interface {
	GetOccupancyRateReport(startDate, endDate string) ([]dto.OccupancyRateResponse, error)
	GetRevenueReport(startDate, endDate string) ([]dto.RevenueResponse, error)
	ExportReportToExcel(startDate, endDate string) ([]byte, error)
}

type reportService struct {
	db *gorm.DB
}

func NewReportService(db *gorm.DB) ReportService {
	return &reportService{db: db}
}

func (s *reportService) GetOccupancyRateReport(startDate, endDate string) ([]dto.OccupancyRateResponse, error) {
	start, err := time.ParseInLocation("2006-01-02", startDate, time.Local)
	if err != nil {
		logger.Warnf("日期格式错误: %v", err)
		return nil, errors.New("开始日期格式错误，请使用YYYY-MM-DD格式")
	}

	end, err := time.ParseInLocation("2006-01-02", endDate, time.Local)
	if err != nil {
		logger.Warnf("日期格式错误: %v", err)
		return nil, errors.New("结束日期格式错误，请使用YYYY-MM-DD格式")
	}

	if start.After(end) {
		return nil, errors.New("开始日期不能晚于结束日期")
	}

	var totalRooms int64
	err = s.db.Model(&model.Room{}).Where("status != ?", model.RoomStatusMaintenance).Count(&totalRooms).Error
	if err != nil {
		logger.Errorf("获取总房间数失败: %v", err)
		return nil, errors.New("获取房间统计失败")
	}

	if totalRooms == 0 {
		return nil, errors.New("没有可用的房间数据")
	}

	var results []dto.OccupancyRateResponse

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		nextDay := d.AddDate(0, 0, 1)

		var occupiedCount int64

		err = s.db.Model(&model.CheckIn{}).
			Where("status = ? OR (status = ? AND actual_check_out >= ? AND actual_check_out < ?)",
				model.CheckInStatusActive,
				model.CheckInStatusCheckedOut,
				d,
				nextDay).
			Distinct("room_id").
			Count(&occupiedCount).Error
		if err != nil {
			logger.Errorf("获取入住房间数失败: date=%s, err=%v", dateStr, err)
			return nil, errors.New("获取入住统计失败")
		}

		var bookingOccupiedCount int64
		err = s.db.Model(&model.Booking{}).
			Where("status = ? AND check_in_date <= ? AND check_out_date > ?",
				model.BookingStatusConfirmed,
				d,
				d).
			Distinct("room_id").
			Count(&bookingOccupiedCount).Error
		if err != nil {
			logger.Errorf("获取预订入住房间数失败: date=%s, err=%v", dateStr, err)
			return nil, errors.New("获取预订统计失败")
		}

		totalOccupied := occupiedCount
		if totalOccupied < bookingOccupiedCount {
			totalOccupied = bookingOccupiedCount
		}

		if totalOccupied > totalRooms {
			totalOccupied = totalRooms
		}

		occupancyRate := float64(totalOccupied) / float64(totalRooms) * 100

		results = append(results, dto.OccupancyRateResponse{
			Date:          dateStr,
			TotalRooms:    totalRooms,
			OccupiedRooms: totalOccupied,
			OccupancyRate: occupancyRate,
		})
	}

	logger.Infof("入住率统计完成: %s 至 %s, 共%d天", startDate, endDate, len(results))
	return results, nil
}

func (s *reportService) GetRevenueReport(startDate, endDate string) ([]dto.RevenueResponse, error) {
	start, err := time.ParseInLocation("2006-01-02", startDate, time.Local)
	if err != nil {
		logger.Warnf("日期格式错误: %v", err)
		return nil, errors.New("开始日期格式错误，请使用YYYY-MM-DD格式")
	}

	end, err := time.ParseInLocation("2006-01-02", endDate, time.Local)
	if err != nil {
		logger.Warnf("日期格式错误: %v", err)
		return nil, errors.New("结束日期格式错误，请使用YYYY-MM-DD格式")
	}

	if start.After(end) {
		return nil, errors.New("开始日期不能晚于结束日期")
	}

	end = end.AddDate(0, 0, 1)

	type RevenueResult struct {
		Date         string
		RoomRevenue  float64
		OtherRevenue float64
		CheckOutCount int64
	}

	var rawResults []RevenueResult

	err = s.db.Model(&model.CheckIn{}).
		Select("DATE(actual_check_out) as date, " +
			"SUM(total_amount - extra_charges) as room_revenue, " +
			"SUM(extra_charges) as other_revenue, " +
			"COUNT(*) as check_out_count").
		Where("status = ? AND actual_check_out >= ? AND actual_check_out < ?",
			model.CheckInStatusCheckedOut,
			start,
			end).
		Group("DATE(actual_check_out)").
		Order("date ASC").
		Scan(&rawResults).Error
	if err != nil {
		logger.Errorf("获取营收统计失败: %v", err)
		return nil, errors.New("获取营收统计失败")
	}

	revenueMap := make(map[string]RevenueResult)
	for _, r := range rawResults {
		revenueMap[r.Date] = r
	}

	var results []dto.RevenueResponse
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		item, exists := revenueMap[dateStr]
		if !exists {
			results = append(results, dto.RevenueResponse{
				Date:          dateStr,
				RoomRevenue:   0,
				OtherRevenue:  0,
				TotalRevenue:  0,
				CheckOutCount: 0,
			})
		} else {
			results = append(results, dto.RevenueResponse{
				Date:          dateStr,
				RoomRevenue:   item.RoomRevenue,
				OtherRevenue:  item.OtherRevenue,
				TotalRevenue:  item.RoomRevenue + item.OtherRevenue,
				CheckOutCount: item.CheckOutCount,
			})
		}
	}

	logger.Infof("营收统计完成: %s 至 %s, 共%d天", startDate, endDate, len(results))
	return results, nil
}

func (s *reportService) ExportReportToExcel(startDate, endDate string) ([]byte, error) {
	occupancyData, err := s.GetOccupancyRateReport(startDate, endDate)
	if err != nil {
		return nil, err
	}

	revenueData, err := s.GetRevenueReport(startDate, endDate)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	defer f.Close()

	f.DeleteSheet("Sheet1")

	err = excel.ExportOccupancyRateSheet(f, occupancyData)
	if err != nil {
		logger.Errorf("生成入住率报表失败: %v", err)
		return nil, errors.New("生成Excel报表失败")
	}

	err = excel.ExportRevenueSheet(f, revenueData)
	if err != nil {
		logger.Errorf("生成营收报表失败: %v", err)
		return nil, errors.New("生成Excel报表失败")
	}

	data, err := excel.SaveExcel(f)
	if err != nil {
		logger.Errorf("保存Excel文件失败: %v", err)
		return nil, errors.New("生成Excel报表失败")
	}

	logger.Infof("Excel报表导出完成: %s 至 %s", startDate, endDate)
	return data, nil
}
