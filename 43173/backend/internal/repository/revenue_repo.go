package repository

import (
	"music-platform/internal/model"
	"music-platform/pkg/database"
	"music-platform/pkg/utils"
	"time"
)

type RevenueRepository struct{}

func NewRevenueRepository() *RevenueRepository {
	return &RevenueRepository{}
}

func (r *RevenueRepository) CreateRevenueRecord(record *model.RevenueRecord) error {
	return database.DB.Create(record).Error
}

func (r *RevenueRepository) BatchCreateRevenueRecords(records []model.RevenueRecord) error {
	return database.DB.Create(&records).Error
}

func (r *RevenueRepository) GetRevenueRecords(userID uint, page, pageSize int, startDate, endDate *time.Time) ([]model.RevenueRecord, int64, error) {
	var records []model.RevenueRecord
	var total int64

	query := database.DB.Model(&model.RevenueRecord{}).Where("user_id = ?", userID)

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&records).Error
	return records, total, err
}

func (r *RevenueRepository) GetRevenueRecordsByArtist(artistID uint, page, pageSize int, startDate, endDate *time.Time) ([]model.RevenueRecord, int64, error) {
	var records []model.RevenueRecord
	var total int64

	query := database.DB.Model(&model.RevenueRecord{}).Where("artist_id = ?", artistID)

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&records).Error
	return records, total, err
}

func (r *RevenueRepository) GetTotalRevenue(userID uint, startDate, endDate *time.Time) (float64, error) {
	var total float64
	query := database.DB.Model(&model.RevenueRecord{}).Where("user_id = ?", userID)

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}

func (r *RevenueRepository) UpdateRevenueStatus(id uint, status int) error {
	return database.DB.Model(&model.RevenueRecord{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     status,
		"settled_at": time.Now(),
	}).Error
}

func (r *RevenueRepository) CreateWithdrawRequest(request *model.WithdrawRequest) error {
	return database.DB.Create(request).Error
}

func (r *RevenueRepository) GetWithdrawRequestByID(id uint) (*model.WithdrawRequest, error) {
	var request model.WithdrawRequest
	err := database.DB.First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *RevenueRepository) GetWithdrawRequests(userID uint, page, pageSize int, status int) ([]model.WithdrawRequest, int64, error) {
	var requests []model.WithdrawRequest
	var total int64

	query := database.DB.Model(&model.WithdrawRequest{}).Where("user_id = ?", userID)

	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&requests).Error
	return requests, total, err
}

func (r *RevenueRepository) GetAllWithdrawRequests(page, pageSize int, status int) ([]model.WithdrawRequest, int64, error) {
	var requests []model.WithdrawRequest
	var total int64

	query := database.DB.Model(&model.WithdrawRequest{})

	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&requests).Error
	return requests, total, err
}

func (r *RevenueRepository) UpdateWithdrawStatus(id uint, status model.WithdrawStatus, updates map[string]interface{}) error {
	return database.DB.Model(&model.WithdrawRequest{}).Where("id = ?", id).Updates(updates).Error
}

func (r *RevenueRepository) GetTotalWithdraw(userID uint) (float64, error) {
	var total float64
	err := database.DB.Model(&model.WithdrawRequest{}).
		Where("user_id = ? AND status IN ?", userID, []int{1, 3}).
		Select("COALESCE(SUM(actual_amount), 0)").Scan(&total).Error
	return total, err
}

func (r *RevenueRepository) GetPendingWithdraw(userID uint) (float64, error) {
	var total float64
	err := database.DB.Model(&model.WithdrawRequest{}).
		Where("user_id = ? AND status = ?", userID, 0).
		Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}

func (r *RevenueRepository) CreateSubscription(subscription *model.Subscription) error {
	return database.DB.Create(subscription).Error
}

func (r *RevenueRepository) GetSubscriptions(userID uint, page, pageSize int) ([]model.Subscription, int64, error) {
	var subscriptions []model.Subscription
	var total int64

	query := database.DB.Model(&model.Subscription{}).Where("user_id = ?", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&subscriptions).Error
	return subscriptions, total, err
}

func (r *RevenueRepository) GetArtistSubscribers(artistID uint, page, pageSize int) ([]model.Subscription, int64, error) {
	var subscriptions []model.Subscription
	var total int64

	query := database.DB.Model(&model.Subscription{}).Where("artist_id = ?", artistID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&subscriptions).Error
	return subscriptions, total, err
}

func (r *RevenueRepository) CreateDailyStats(stats *model.DailyStats) error {
	return database.DB.Create(stats).Error
}

func (r *RevenueRepository) UpdateDailyStats(id uint, updates map[string]interface{}) error {
	return database.DB.Model(&model.DailyStats{}).Where("id = ?", id).Updates(updates).Error
}

func (r *RevenueRepository) GetDailyStats(userID uint, startDate, endDate time.Time) ([]model.DailyStats, error) {
	var stats []model.DailyStats
	err := database.DB.Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
		Order("date ASC").Find(&stats).Error
	return stats, err
}

func (r *RevenueRepository) GetDailyStatsByArtist(artistID uint, startDate, endDate time.Time) ([]model.DailyStats, error) {
	var stats []model.DailyStats
	err := database.DB.Where("artist_id = ? AND date >= ? AND date <= ?", artistID, startDate, endDate).
		Order("date ASC").Find(&stats).Error
	return stats, err
}

func (r *RevenueRepository) CreateOperationLog(log *model.OperationLog) error {
	return database.DB.Create(log).Error
}

func (r *RevenueRepository) GetOperationLogs(page, pageSize int, userID uint, module string, keyword string) ([]model.OperationLog, int64, error) {
	var logs []model.OperationLog
	var total int64

	query := database.DB.Model(&model.OperationLog{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if keyword != "" {
		query = query.Where("operation LIKE ? OR path LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := utils.GetOffset(page, pageSize)
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}

func (r *RevenueRepository) GetUnsettledRevenueRecords(userID uint) ([]model.RevenueRecord, error) {
	var records []model.RevenueRecord
	err := database.DB.Where("user_id = ? AND status = ?", userID, 0).Find(&records).Error
	return records, err
}

func (r *RevenueRepository) BatchUpdateRevenueStatus(ids []uint, status int) error {
	return database.DB.Model(&model.RevenueRecord{}).Where("id IN ?", ids).Updates(map[string]interface{}{
		"status":     status,
		"settled_at": time.Now(),
	}).Error
}

func (r *RevenueRepository) GetRevenueByPeriod(userID uint, periodStart, periodEnd time.Time) ([]model.RevenueRecord, error) {
	var records []model.RevenueRecord
	err := database.DB.Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, periodStart, periodEnd).
		Find(&records).Error
	return records, err
}

func (r *RevenueRepository) SumRevenueByType(userID uint, revenueType model.RevenueType, startDate, endDate *time.Time) (float64, error) {
	var total float64
	query := database.DB.Model(&model.RevenueRecord{}).
		Where("user_id = ? AND type = ?", userID, revenueType)

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}
