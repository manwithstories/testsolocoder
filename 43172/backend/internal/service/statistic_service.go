package service

import (
	"context"
	"fmt"
	"time"

	"luxury-trading-platform/internal/cache"
	"luxury-trading-platform/internal/model"
	"luxury-trading-platform/internal/repository"

	"gorm.io/gorm"
)

type StatisticService struct {
	statRepo    *repository.StatisticRepository
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
	userRepo    *repository.UserRepository
	redisClient *cache.RedisClient
	db          *gorm.DB
}

func NewStatisticService(statRepo *repository.StatisticRepository, orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository, userRepo *repository.UserRepository, redisClient *cache.RedisClient, db *gorm.DB) *StatisticService {
	return &StatisticService{
		statRepo:    statRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
		redisClient: redisClient,
		db:          db,
	}
}

type DashboardStats struct {
	TotalOrders    int64              `json:"total_orders"`
	TotalAmount    float64            `json:"total_amount"`
	TotalUsers     int64              `json:"total_users"`
	TotalProducts  int64              `json:"total_products"`
	TransactionTrend []TransactionTrend `json:"transaction_trend"`
	BrandRankings  []BrandRanking     `json:"brand_rankings"`
	AuthStats      AuthStats          `json:"auth_stats"`
	RecentOrders   []model.Order      `json:"recent_orders"`
}

type TransactionTrend struct {
	Date        string  `json:"date"`
	OrderCount  int     `json:"order_count"`
	TotalAmount float64 `json:"total_amount"`
}

type BrandRanking struct {
	BrandName    string `json:"brand_name"`
	ProductCount int64  `json:"product_count"`
	OrderCount   int64  `json:"order_count"`
}

type AuthStats struct {
	Total     int64   `json:"total"`
	Completed int64   `json:"completed"`
	Passed    int64   `json:"passed"`
	PassRate  float64 `json:"pass_rate"`
}

func (s *StatisticService) GetDashboardStats(ctx context.Context, days int) (*DashboardStats, error) {
	cacheKey := fmt.Sprintf("dashboard:stats:%d", days)

	if s.redisClient != nil {
		cached, err := s.redisClient.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			return s.getStatsFromDB(ctx, days)
		}
	}

	stats, err := s.getStatsFromDB(ctx, days)
	if err != nil {
		return nil, err
	}

	if s.redisClient != nil {
		_ = s.redisClient.Set(ctx, cacheKey, stats, 5*time.Minute)
	}

	return stats, nil
}

func (s *StatisticService) getStatsFromDB(ctx context.Context, days int) (*DashboardStats, error) {
	totalOrders, err := s.statRepo.GetTotalOrders()
	if err != nil {
		return nil, fmt.Errorf("failed to get total orders: %w", err)
	}

	totalAmount, err := s.statRepo.GetTotalAmount()
	if err != nil {
		return nil, fmt.Errorf("failed to get total amount: %w", err)
	}

	totalUsers, err := s.statRepo.GetTotalUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}

	totalProducts, err := s.statRepo.GetTotalProducts()
	if err != nil {
		return nil, fmt.Errorf("failed to get total products: %w", err)
	}

	transactionTrend, err := s.getTransactionTrend(days)
	if err != nil {
		return nil, err
	}

	brandRankings, err := s.getBrandRankings()
	if err != nil {
		return nil, err
	}

	authStats, err := s.getAuthStats()
	if err != nil {
		return nil, err
	}

	recentOrders, _, err := s.orderRepo.List(1, 10, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get recent orders: %w", err)
	}

	return &DashboardStats{
		TotalOrders:       totalOrders,
		TotalAmount:       totalAmount,
		TotalUsers:        totalUsers,
		TotalProducts:     totalProducts,
		TransactionTrend:  transactionTrend,
		BrandRankings:     brandRankings,
		AuthStats:         authStats,
		RecentOrders:      recentOrders,
	}, nil
}

func (s *StatisticService) getTransactionTrend(days int) ([]TransactionTrend, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	stats, err := s.statRepo.GetTransactionTrend(startDate, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction trend: %w", err)
	}

	var trends []TransactionTrend
	for _, stat := range stats {
		trends = append(trends, TransactionTrend{
			Date:        stat.Date.Format("2006-01-02"),
			OrderCount:  stat.TotalOrders,
			TotalAmount: stat.TotalAmount,
		})
	}

	return trends, nil
}

func (s *StatisticService) getBrandRankings() ([]BrandRanking, error) {
	results, err := s.statRepo.GetBrandRankings()
	if err != nil {
		return nil, fmt.Errorf("failed to get brand rankings: %w", err)
	}

	var rankings []BrandRanking
	for _, r := range results {
		rankings = append(rankings, BrandRanking{
			BrandName:    r.BrandName,
			ProductCount: r.ProductCount,
			OrderCount:   r.OrderCount,
		})
	}

	return rankings, nil
}

func (s *StatisticService) getAuthStats() (AuthStats, error) {
	result, err := s.statRepo.GetAuthStatistics()
	if err != nil {
		return AuthStats{}, fmt.Errorf("failed to get auth stats: %w", err)
	}

	return AuthStats{
		Total:     result.Total,
		Completed: result.Completed,
		Passed:    result.Passed,
		PassRate:  result.PassRate,
	}, nil
}

func (s *StatisticService) InvalidateCache(ctx context.Context) {
	if s.redisClient != nil {
		_ = s.redisClient.Del(ctx, "dashboard:stats:7")
		_ = s.redisClient.Del(ctx, "dashboard:stats:30")
	}
}
