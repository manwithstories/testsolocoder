package service

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
)

// StatisticsService 数据统计业务逻辑层
type StatisticsService struct {
	db *gorm.DB
}

// NewStatisticsService 创建数据统计服务
func NewStatisticsService(db *gorm.DB) *StatisticsService {
	return &StatisticsService{db: db}
}

// parseDateRange 解析时间范围
func parseDateRange(start, end string) (time.Time, time.Time, error) {
	loc := time.Local
	startTime := time.Date(2000, 1, 1, 0, 0, 0, 0, loc)
	endTime := time.Now()

	if start != "" {
		t, err := time.ParseInLocation("2006-01-02", start, loc)
		if err != nil {
			return startTime, endTime, fmt.Errorf("开始日期格式错误: %w", err)
		}
		startTime = t
	}
	if end != "" {
		t, err := time.ParseInLocation("2006-01-02", end, loc)
		if err != nil {
			return startTime, endTime, fmt.Errorf("结束日期格式错误: %w", err)
		}
		endTime = t.Add(24*time.Hour - time.Nanosecond)
	}
	return startTime, endTime, nil
}

// GetSalesTrend 按日/周/月统计销售额
func (s *StatisticsService) GetSalesTrend(req *dto.StatisticsRequest) (*dto.SalesTrendResponse, error) {
	startTime, endTime, err := parseDateRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	granularity := req.Granularity
	if granularity == "" {
		granularity = "month"
	}

	var dateFormat string
	var groupExpr string
	switch granularity {
	case "day":
		dateFormat = "%Y-%m-%d"
	case "week":
		dateFormat = "%Y-%u"
	default:
		dateFormat = "%Y-%m"
	}
	groupExpr = fmt.Sprintf("DATE_FORMAT(created_at, '%s')", dateFormat)

	type row struct {
		Period     string  `gorm:"column:period"`
		Amount     float64 `gorm:"column:amount"`
		OrderCount int64   `gorm:"column:order_count"`
	}

	var rows []row
	if err := s.db.Model(&model.Order{}).
		Select(fmt.Sprintf("%s as period, COALESCE(SUM(final_amount),0) as amount, COUNT(*) as order_count", groupExpr)).
		Where("created_at >= ? AND created_at <= ?", startTime, endTime).
		Where("status != ?", model.OrderStatusCancelled).
		Group(groupExpr).
		Order("period ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	list := make([]dto.SalesTrendItem, 0, len(rows))
	var totalAmount float64
	var totalOrders int64
	for _, r := range rows {
		list = append(list, dto.SalesTrendItem{
			Period:     r.Period,
			Amount:     r.Amount,
			OrderCount: r.OrderCount,
		})
		totalAmount += r.Amount
		totalOrders += r.OrderCount
	}

	return &dto.SalesTrendResponse{
		Granularity: granularity,
		StartDate:   startTime.Format("2006-01-02"),
		EndDate:     endTime.Format("2006-01-02"),
		List:        list,
		TotalAmount: totalAmount,
		TotalOrders: totalOrders,
	}, nil
}

// GetCustomerProfile 客户画像统计（角色分布、地区分布）
func (s *StatisticsService) GetCustomerProfile() (*dto.CustomerProfileResponse, error) {
	var total int64
	if err := s.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, err
	}

	type roleRow struct {
		Role  string `gorm:"column:role"`
		Count int64  `gorm:"column:count"`
	}
	var roleRows []roleRow
	if err := s.db.Model(&model.User{}).
		Select("role as role, COUNT(*) as count").
		Group("role").
		Scan(&roleRows).Error; err != nil {
		return nil, err
	}

	roleDistribution := make([]dto.RoleDistribution, 0, len(roleRows))
	for _, r := range roleRows {
		roleDistribution = append(roleDistribution, dto.RoleDistribution{
			Role:  r.Role,
			Count: r.Count,
		})
	}

	// 地区分布：根据订单中的地址统计（按首字前缀聚合）
	type regionRow struct {
		Region string `gorm:"column:region"`
		Count  int64  `gorm:"column:count"`
	}
	var regionRows []regionRow
	if err := s.db.Model(&model.Order{}).
		Select("CASE WHEN address = '' OR address IS NULL THEN '未知' ELSE SUBSTRING_INDEX(address, ' ', 1) END as region, COUNT(DISTINCT owner_id) as count").
		Group("region").
		Order("count DESC").
		Limit(20).
		Scan(&regionRows).Error; err != nil {
		return nil, err
	}

	regionDistribution := make([]dto.RegionDistribution, 0, len(regionRows))
	for _, r := range regionRows {
		regionDistribution = append(regionDistribution, dto.RegionDistribution{
			Region: r.Region,
			Count:  r.Count,
		})
	}

	return &dto.CustomerProfileResponse{
		TotalUsers:         total,
		RoleDistribution:   roleDistribution,
		RegionDistribution: regionDistribution,
	}, nil
}

// ExportExcel 导出统计报表 Excel
func (s *StatisticsService) ExportExcel(req *dto.ExportRequest) ([]byte, string, error) {
	startTime, endTime, err := parseDateRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, "", err
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			_ = err
		}
	}()

	filename := fmt.Sprintf("%s_%s_%s.xlsx", req.Type, startTime.Format("20060102"), endTime.Format("20060102"))

	switch req.Type {
	case "sales":
		if err := exportSalesSheet(s.db, f, startTime, endTime); err != nil {
			return nil, "", err
		}
	case "orders":
		if err := exportOrdersSheet(s.db, f, startTime, endTime); err != nil {
			return nil, "", err
		}
	case "users":
		if err := exportUsersSheet(s.db, f); err != nil {
			return nil, "", err
		}
	default:
		return nil, "", errors.New("不支持的导出类型")
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), filename, nil
}

// exportSalesSheet 导出销售统计
func exportSalesSheet(db *gorm.DB, f *excelize.File, start, end time.Time) error {
	sheet := "销售统计"
	if err := f.SetSheetName("Sheet1", sheet); err != nil {
		return err
	}

	headers := []string{"月份", "订单数", "销售额"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			return err
		}
	}

	type row struct {
		Period     string  `gorm:"column:period"`
		Amount     float64 `gorm:"column:amount"`
		OrderCount int64   `gorm:"column:order_count"`
	}
	var rows []row
	if err := db.Model(&model.Order{}).
		Select("DATE_FORMAT(created_at, '%Y-%m') as period, COALESCE(SUM(final_amount),0) as amount, COUNT(*) as order_count").
		Where("created_at >= ? AND created_at <= ?", start, end).
		Where("status != ?", model.OrderStatusCancelled).
		Group("period").
		Order("period ASC").
		Scan(&rows).Error; err != nil {
		return err
	}

	for i, r := range rows {
		rowNum := i + 2
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), r.Period)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), r.OrderCount)
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), r.Amount)
	}
	return nil
}

// exportOrdersSheet 导出订单明细
func exportOrdersSheet(db *gorm.DB, f *excelize.File, start, end time.Time) error {
	sheet := "订单明细"
	if err := f.SetSheetName("Sheet1", sheet); err != nil {
		return err
	}

	headers := []string{"订单号", "业主ID", "厂商ID", "状态", "订单金额", "实付金额", "创建时间"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			return err
		}
	}

	var orders []model.Order
	if err := db.Where("created_at >= ? AND created_at <= ?", start, end).
		Order("id ASC").
		Find(&orders).Error; err != nil {
		return err
	}

	for i, o := range orders {
		rowNum := i + 2
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), o.OrderNo)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), o.OwnerID)
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), o.ManufacturerID)
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), o.Status)
		_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), o.TotalAmount)
		_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), o.FinalAmount)
		_ = f.SetCellValue(sheet, fmt.Sprintf("G%d", rowNum), o.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	return nil
}

// exportUsersSheet 导出用户列表
func exportUsersSheet(db *gorm.DB, f *excelize.File) error {
	sheet := "用户列表"
	if err := f.SetSheetName("Sheet1", sheet); err != nil {
		return err
	}

	headers := []string{"用户ID", "用户名", "角色", "昵称", "手机", "邮箱", "状态", "注册时间"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheet, cell, h); err != nil {
			return err
		}
	}

	var users []model.User
	if err := db.Order("id ASC").Find(&users).Error; err != nil {
		return err
	}

	for i, u := range users {
		rowNum := i + 2
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), u.ID)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), u.Username)
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), u.Role)
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), u.Nickname)
		_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), u.Phone)
		_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), u.Email)
		_ = f.SetCellValue(sheet, fmt.Sprintf("G%d", rowNum), u.Status)
		_ = f.SetCellValue(sheet, fmt.Sprintf("H%d", rowNum), u.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	return nil
}
