package service

import (
	"time"

	"music-platform/internal/model"
	"music-platform/internal/repository"
	apperrors "music-platform/pkg/errors"
	"music-platform/pkg/database"
	"music-platform/pkg/utils"
	"gorm.io/gorm"
)

type RevenueService struct {
	revenueRepo *repository.RevenueRepository
	userRepo    *repository.UserRepository
}

func NewRevenueService() *RevenueService {
	return &RevenueService{
		revenueRepo: repository.NewRevenueRepository(),
		userRepo:    repository.NewUserRepository(),
	}
}

type WithdrawRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Method      string  `json:"method" binding:"required"`
	Account     string  `json:"account" binding:"required"`
	AccountName string  `json:"account_name"`
	BankName    string  `json:"bank_name"`
}

type SettleRevenueRequest struct {
	UserIDs []uint `json:"user_ids"`
	Period  string `json:"period" binding:"required"`
}

func (s *RevenueService) GetRevenueRecords(userID uint, page, pageSize int, startDate, endDate *time.Time) ([]model.RevenueRecord, int64, error) {
	return s.revenueRepo.GetRevenueRecords(userID, page, pageSize, startDate, endDate)
}

func (s *RevenueService) GetArtistRevenueRecords(artistID uint, page, pageSize int, startDate, endDate *time.Time) ([]model.RevenueRecord, int64, error) {
	return s.revenueRepo.GetRevenueRecordsByArtist(artistID, page, pageSize, startDate, endDate)
}

func (s *RevenueService) GetTotalRevenue(userID uint, startDate, endDate *time.Time) (float64, error) {
	return s.revenueRepo.GetTotalRevenue(userID, startDate, endDate)
}

func (s *RevenueService) GetRevenueSummary(userID uint) (map[string]interface{}, error) {
	now := time.Now()
	monthStart := utils.GetMonthStart()
	lastMonthStart := monthStart.AddDate(0, -1, 0)

	totalRevenue, _ := s.revenueRepo.GetTotalRevenue(userID, nil, nil)
	monthRevenue, _ := s.revenueRepo.GetTotalRevenue(userID, &monthStart, &now)
	lastMonthRevenue, _ := s.revenueRepo.GetTotalRevenue(userID, &lastMonthStart, &monthStart)

	playRevenue, _ := s.revenueRepo.SumRevenueByType(userID, model.RevenueTypePlay, nil, nil)
	subscriptionRevenue, _ := s.revenueRepo.SumRevenueByType(userID, model.RevenueTypeSubscription, nil, nil)
	ticketRevenue, _ := s.revenueRepo.SumRevenueByType(userID, model.RevenueTypeTicket, nil, nil)

	totalWithdraw, _ := s.revenueRepo.GetTotalWithdraw(userID)
	pendingWithdraw, _ := s.revenueRepo.GetPendingWithdraw(userID)

	artistInfo, _ := s.userRepo.FindArtistInfoByUserID(userID)

	var availableBalance float64
	if artistInfo != nil {
		availableBalance = artistInfo.Balance
	}

	return map[string]interface{}{
		"total_revenue":        totalRevenue,
		"month_revenue":        monthRevenue,
		"last_month_revenue":   lastMonthRevenue,
		"play_revenue":         playRevenue,
		"subscription_revenue": subscriptionRevenue,
		"ticket_revenue":       ticketRevenue,
		"available_balance":    availableBalance,
		"frozen_balance":       artistInfo.FrozenBalance,
		"total_withdraw":       totalWithdraw,
		"pending_withdraw":     pendingWithdraw,
	}, nil
}

func (s *RevenueService) RequestWithdraw(userID uint, req *WithdrawRequest) (*model.WithdrawRequest, error) {
	artistInfo, err := s.userRepo.FindArtistInfoByUserID(userID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	if artistInfo.Balance < req.Amount {
		return nil, apperrors.ErrInsufficientBalance
	}

	fee := req.Amount * 0.01
	actualAmount := req.Amount - fee

	err = s.userRepo.UpdateArtistBalance(artistInfo.ID, -req.Amount)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.UpdateFrozenBalance(artistInfo.ID, req.Amount)
	if err != nil {
		return nil, err
	}

	withdrawRequest := &model.WithdrawRequest{
		UserID:       userID,
		ArtistID:     artistInfo.ID,
		Amount:       req.Amount,
		Fee:          fee,
		ActualAmount: actualAmount,
		Method:       req.Method,
		Account:      req.Account,
		AccountName:  req.AccountName,
		BankName:     req.BankName,
		Status:       model.WithdrawStatusPending,
	}

	err = s.revenueRepo.CreateWithdrawRequest(withdrawRequest)
	if err != nil {
		_ = s.userRepo.UpdateArtistBalance(artistInfo.ID, req.Amount)
		_ = s.userRepo.UpdateFrozenBalance(artistInfo.ID, -req.Amount)
		return nil, err
	}

	return withdrawRequest, nil
}

func (s *RevenueService) GetWithdrawRequests(userID uint, page, pageSize int, status int) ([]model.WithdrawRequest, int64, error) {
	return s.revenueRepo.GetWithdrawRequests(userID, page, pageSize, status)
}

func (s *RevenueService) GetAllWithdrawRequests(page, pageSize int, status int) ([]model.WithdrawRequest, int64, error) {
	return s.revenueRepo.GetAllWithdrawRequests(page, pageSize, status)
}

func (s *RevenueService) ApproveWithdraw(withdrawID uint, approvedBy uint) error {
	withdrawRequest, err := s.revenueRepo.GetWithdrawRequestByID(withdrawID)
	if err != nil {
		return apperrors.NewAppError(5001, "提现申请不存在")
	}

	if withdrawRequest.Status != model.WithdrawStatusPending {
		return apperrors.NewAppError(5002, "申请状态异常")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":      model.WithdrawStatusApproved,
		"approved_at": &now,
		"approved_by": &approvedBy,
	}

	return s.revenueRepo.UpdateWithdrawStatus(withdrawID, model.WithdrawStatusApproved, updates)
}

func (s *RevenueService) RejectWithdraw(withdrawID uint, reason string) error {
	withdrawRequest, err := s.revenueRepo.GetWithdrawRequestByID(withdrawID)
	if err != nil {
		return apperrors.NewAppError(5001, "提现申请不存在")
	}

	if withdrawRequest.Status != model.WithdrawStatusPending {
		return apperrors.NewAppError(5002, "申请状态异常")
	}

	artistInfo, _ := s.userRepo.FindArtistInfoByUserID(withdrawRequest.UserID)
	if artistInfo != nil {
		_ = s.userRepo.UpdateFrozenBalance(artistInfo.ID, -withdrawRequest.Amount)
		_ = s.userRepo.UpdateArtistBalance(artistInfo.ID, withdrawRequest.Amount)
	}

	updates := map[string]interface{}{
		"status":        model.WithdrawStatusRejected,
		"reject_reason": reason,
	}

	return s.revenueRepo.UpdateWithdrawStatus(withdrawID, model.WithdrawStatusRejected, updates)
}

func (s *RevenueService) MarkWithdrawPaid(withdrawID uint, transactionNo string) error {
	withdrawRequest, err := s.revenueRepo.GetWithdrawRequestByID(withdrawID)
	if err != nil {
		return apperrors.NewAppError(5001, "提现申请不存在")
	}

	if withdrawRequest.Status != model.WithdrawStatusApproved {
		return apperrors.NewAppError(5002, "申请状态异常")
	}

	artistInfo, _ := s.userRepo.FindArtistInfoByUserID(withdrawRequest.UserID)
	if artistInfo != nil {
		_ = s.userRepo.UpdateFrozenBalance(artistInfo.ID, -withdrawRequest.Amount)
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":         model.WithdrawStatusPaid,
		"paid_at":        &now,
		"transaction_no": transactionNo,
	}

	return s.revenueRepo.UpdateWithdrawStatus(withdrawID, model.WithdrawStatusPaid, updates)
}

func (s *RevenueService) SettleRevenue(period string) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		periodStart, _ := time.Parse("2006-01", period)
		periodEnd := periodStart.AddDate(0, 1, 0)

		records, err := s.revenueRepo.GetRevenueByPeriod(0, periodStart, periodEnd)
		if err != nil {
			return err
		}

		userRevenueMap := make(map[uint]float64)
		var recordIDs []uint

		for _, record := range records {
			if record.Status == 0 {
				userRevenueMap[record.UserID] += record.Amount
				recordIDs = append(recordIDs, record.ID)
			}
		}

		for userID, amount := range userRevenueMap {
			artistInfo, err := s.userRepo.FindArtistInfoByUserID(userID)
			if err != nil {
				continue
			}
			err = s.userRepo.UpdateArtistBalance(artistInfo.ID, amount)
			if err != nil {
				return err
			}
		}

		if len(recordIDs) > 0 {
			err = s.revenueRepo.BatchUpdateRevenueStatus(recordIDs, 1)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *RevenueService) CreateRevenueRecord(record *model.RevenueRecord) error {
	return s.revenueRepo.CreateRevenueRecord(record)
}

func (s *RevenueService) GetSubscriptions(userID uint, page, pageSize int) ([]model.Subscription, int64, error) {
	return s.revenueRepo.GetSubscriptions(userID, page, pageSize)
}

func (s *RevenueService) GetArtistSubscribers(artistID uint, page, pageSize int) ([]model.Subscription, int64, error) {
	return s.revenueRepo.GetArtistSubscribers(artistID, page, pageSize)
}

func (s *RevenueService) GetDailyStats(userID uint, startDate, endDate time.Time) ([]model.DailyStats, error) {
	return s.revenueRepo.GetDailyStats(userID, startDate, endDate)
}

func (s *RevenueService) GetArtistDailyStats(artistID uint, startDate, endDate time.Time) ([]model.DailyStats, error) {
	return s.revenueRepo.GetDailyStatsByArtist(artistID, startDate, endDate)
}

func (s *RevenueService) GetArtistStats(artistID uint) (map[string]interface{}, error) {
	now := time.Now()
	weekStart := utils.GetWeekStart()
	monthStart := utils.GetMonthStart()

	dailyStats, _ := s.revenueRepo.GetDailyStatsByArtist(artistID, monthStart, now)

	var totalPlays, totalLikes, totalShares, totalComments int64
	var totalRevenue float64
	var dailyData []map[string]interface{}

	for _, stat := range dailyStats {
		totalPlays += stat.NewPlays
		totalLikes += stat.NewLikes
		totalShares += stat.NewShares
		totalComments += stat.NewComments
		totalRevenue += stat.Revenue

		dailyData = append(dailyData, map[string]interface{}{
			"date":         stat.Date.Format("2006-01-02"),
			"new_followers": stat.NewFollowers,
			"new_plays":    stat.NewPlays,
			"new_likes":    stat.NewLikes,
			"new_shares":   stat.NewShares,
			"new_comments": stat.NewComments,
			"revenue":      stat.Revenue,
		})
	}

	weeklyStats, _ := s.revenueRepo.GetDailyStatsByArtist(artistID, weekStart, now)
	var weeklyPlays, weeklyLikes, weeklyShares, weeklyComments int64
	var weeklyRevenue float64

	for _, stat := range weeklyStats {
		weeklyPlays += stat.NewPlays
		weeklyLikes += stat.NewLikes
		weeklyShares += stat.NewShares
		weeklyComments += stat.NewComments
		weeklyRevenue += stat.Revenue
	}

	return map[string]interface{}{
		"monthly_summary": map[string]interface{}{
			"total_plays":    totalPlays,
			"total_likes":    totalLikes,
			"total_shares":   totalShares,
			"total_comments": totalComments,
			"total_revenue":  totalRevenue,
		},
		"weekly_summary": map[string]interface{}{
			"total_plays":    weeklyPlays,
			"total_likes":    weeklyLikes,
			"total_shares":   weeklyShares,
			"total_comments": weeklyComments,
			"total_revenue":  weeklyRevenue,
		},
		"daily_data": dailyData,
	}, nil
}

func (s *RevenueService) GetOperationLogs(page, pageSize int, userID uint, module string, keyword string) ([]model.OperationLog, int64, error) {
	return s.revenueRepo.GetOperationLogs(page, pageSize, userID, module, keyword)
}

func (s *RevenueService) ExportRevenueExcel(userID uint, startDate, endDate *time.Time) ([]model.RevenueRecord, error) {
	records, _, err := s.revenueRepo.GetRevenueRecords(userID, 1, 10000, startDate, endDate)
	return records, err
}

func (s *RevenueService) ExportWithdrawExcel(userID uint, status int) ([]model.WithdrawRequest, error) {
	records, _, err := s.revenueRepo.GetWithdrawRequests(userID, 1, 10000, status)
	return records, err
}
