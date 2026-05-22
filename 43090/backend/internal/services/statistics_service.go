package services

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"auction-system/internal/dto"
	"auction-system/internal/models"
)

type StatisticsService struct{}

func NewStatisticsService() *StatisticsService {
	return &StatisticsService{}
}

func (s *StatisticsService) GetOverallStatistics(query *dto.StatisticsQuery) (*dto.StatisticsResponse, error) {
	startTime, endTime := s.parseDateRange(query.StartDate, query.EndDate)

	var totalAuctions int64
	auctionQuery := models.DB.Model(&models.AuctionItem{})
	if startTime != nil {
		auctionQuery = auctionQuery.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		auctionQuery = auctionQuery.Where("created_at <= ?", *endTime)
	}
	auctionQuery.Count(&totalAuctions)

	var totalBids int64
	bidQuery := models.DB.Model(&models.Bid{})
	if startTime != nil {
		bidQuery = bidQuery.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		bidQuery = bidQuery.Where("created_at <= ?", *endTime)
	}
	bidQuery.Count(&totalBids)

	var totalOrders int64
	var totalAmount float64
	orderQuery := models.DB.Model(&models.Order{})
	if startTime != nil {
		orderQuery = orderQuery.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		orderQuery = orderQuery.Where("created_at <= ?", *endTime)
	}
	orderQuery.Count(&totalOrders)
	orderQuery.Select("IFNULL(SUM(price), 0)").Row().Scan(&totalAmount)

	var successOrders int64
	successQuery := models.DB.Model(&models.Order{}).Where("status >= ?", OrderStatusPaid)
	if startTime != nil {
		successQuery = successQuery.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		successQuery = successQuery.Where("created_at <= ?", *endTime)
	}
	successQuery.Count(&successOrders)

	successRate := 0.0
	if totalOrders > 0 {
		successRate = float64(successOrders) / float64(totalOrders) * 100
	}

	var activeUsers int64
	userQuery := models.DB.Model(&models.User{}).Where("status = ?", 1)
	userQuery.Count(&activeUsers)

	var newUsers int64
	newUserQuery := models.DB.Model(&models.User{})
	if startTime != nil {
		newUserQuery = newUserQuery.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		newUserQuery = newUserQuery.Where("created_at <= ?", *endTime)
	}
	newUserQuery.Count(&newUsers)

	var averageBidAmount float64
	if totalBids > 0 {
		avgQuery := models.DB.Model(&models.Bid{})
		if startTime != nil {
			avgQuery = avgQuery.Where("created_at >= ?", *startTime)
		}
		if endTime != nil {
			avgQuery = avgQuery.Where("created_at <= ?", *endTime)
		}
		avgQuery.Select("IFNULL(AVG(amount), 0)").Row().Scan(&averageBidAmount)
	}

	return &dto.StatisticsResponse{
		TotalAuctions:    totalAuctions,
		TotalBids:        totalBids,
		TotalOrders:      totalOrders,
		TotalAmount:      totalAmount,
		SuccessRate:      successRate,
		ActiveUsers:      activeUsers,
		NewUsers:         newUsers,
		AverageBidAmount: averageBidAmount,
	}, nil
}

func (s *StatisticsService) ExportOrdersToCSV(query *dto.StatisticsQuery) (string, error) {
	startTime, endTime := s.parseDateRange(query.StartDate, query.EndDate)

	var orders []models.Order
	orderQuery := models.DB.Model(&models.Order{}).Preload("AuctionItem").Preload("Buyer").Preload("Seller")
	if startTime != nil {
		orderQuery = orderQuery.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		orderQuery = orderQuery.Where("created_at <= ?", *endTime)
	}
	orderQuery.Order("created_at DESC").Find(&orders)

	exportDir := "./exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("%s/orders_%s.csv", exportDir, time.Now().Format("20060102150405"))

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"订单号", "拍卖品", "买家", "卖家", "价格", "状态", "创建时间", "支付时间"}
	writer.Write(headers)

	statusMap := map[int]string{
		0: "待支付",
		1: "已支付",
		2: "已发货",
		3: "已送达",
		4: "已完成",
		5: "已取消",
	}

	for _, order := range orders {
		status := statusMap[order.Status]
		paymentTime := ""
		if order.PaymentTime != nil {
			paymentTime = order.PaymentTime.Format("2006-01-02 15:04:05")
		}

		row := []string{
			order.OrderNo,
			order.AuctionItem.Title,
			order.Buyer.Username,
			order.Seller.Username,
			fmt.Sprintf("%.2f", order.Price),
			status,
			order.CreatedAt.Format("2006-01-02 15:04:05"),
			paymentTime,
		}
		writer.Write(row)
	}

	return filename, nil
}

func (s *StatisticsService) ExportBidsToCSV(query *dto.StatisticsQuery) (string, error) {
	startTime, endTime := s.parseDateRange(query.StartDate, query.EndDate)

	var bids []models.Bid
	bidQuery := models.DB.Model(&models.Bid{}).Preload("AuctionItem").Preload("User")
	if startTime != nil {
		bidQuery = bidQuery.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		bidQuery = bidQuery.Where("created_at <= ?", *endTime)
	}
	bidQuery.Order("created_at DESC").Find(&bids)

	exportDir := "./exports"
	os.MkdirAll(exportDir, 0755)
	filename := fmt.Sprintf("%s/bids_%s.csv", exportDir, time.Now().Format("20060102150405"))

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"出价时间", "拍卖品", "用户", "出价金额", "是否自动出价", "是否领先"}
	writer.Write(headers)

	for _, bid := range bids {
		row := []string{
			bid.CreatedAt.Format("2006-01-02 15:04:05"),
			bid.AuctionItem.Title,
			bid.User.Username,
			fmt.Sprintf("%.2f", bid.Amount),
			map[bool]string{true: "是", false: "否"}[bid.IsAutoBid == 1],
			map[bool]string{true: "是", false: "否"}[bid.IsWinning == 1],
		}
		writer.Write(row)
	}

	return filename, nil
}

func (s *StatisticsService) GetUserBidStatistics(userID uint) (map[string]interface{}, error) {
	var totalBids int64
	models.DB.Model(&models.Bid{}).Where("user_id = ?", userID).Count(&totalBids)

	var winningBids int64
	models.DB.Model(&models.Bid{}).Where("user_id = ? AND is_winning = 1", userID).Count(&winningBids)

	var totalOrders int64
	models.DB.Model(&models.Order{}).Where("buyer_id = ?", userID).Count(&totalOrders)

	var totalSpent float64
	models.DB.Model(&models.Order{}).Where("buyer_id = ? AND status >= ?", userID, OrderStatusPaid).Select("IFNULL(SUM(price), 0)").Row().Scan(&totalSpent)

	var totalSold float64
	models.DB.Model(&models.Order{}).Where("seller_id = ? AND status = ?", userID, OrderStatusCompleted).Select("IFNULL(SUM(price), 0)").Row().Scan(&totalSold)

	return map[string]interface{}{
		"total_bids":     totalBids,
		"winning_bids":   winningBids,
		"total_orders":   totalOrders,
		"total_spent":    totalSpent,
		"total_sold":     totalSold,
	}, nil
}

func (s *StatisticsService) parseDateRange(startDate, endDate string) (*time.Time, *time.Time) {
	var startTime, endTime *time.Time

	if startDate != "" {
		t, err := time.ParseInLocation("2006-01-02", startDate, time.Local)
		if err == nil {
			startTime = &t
		}
	}

	if endDate != "" {
		t, err := time.ParseInLocation("2006-01-02", endDate, time.Local)
		if err == nil {
			t = t.Add(24 * time.Hour)
			endTime = &t
		}
	}

	return startTime, endTime
}
