package services

import (
	"time"

	"github.com/notification-center/internal/database"
	"github.com/notification-center/internal/errors"
	"github.com/notification-center/internal/logger"
	"github.com/notification-center/internal/models"
	"go.uber.org/zap"
)

type StatisticsService struct {
	channelService *ChannelService
}

func NewStatisticsService(channelService *ChannelService) *StatisticsService {
	return &StatisticsService{
		channelService: channelService,
	}
}

type MessageStats struct {
	Total        int64   `json:"total"`
	Success      int64   `json:"success"`
	Failed       int64   `json:"failed"`
	Pending      int64   `json:"pending"`
	SuccessRate  float64 `json:"success_rate"`
	TotalDuration int64  `json:"total_duration_ms"`
	AvgDuration  int64   `json:"avg_duration_ms"`
}

type ChannelStats struct {
	ChannelID   uint   `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	Total       int64  `json:"total"`
	Success     int64  `json:"success"`
	Failed      int64  `json:"failed"`
}

func (ss *StatisticsService) GetMessageStats(channelID uint, start, end time.Time) (*MessageStats, error) {
	query := database.DB.Model(&models.Message{})

	if channelID > 0 {
		query = query.Where("channel_id = ?", channelID)
	}
	if !start.IsZero() {
		query = query.Where("created_at >= ?", start)
	}
	if !end.IsZero() {
		query = query.Where("created_at <= ?", end)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error("count total messages failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	var success int64
	if err := query.Where("status = ?", models.MessageStatusSent).Count(&success).Error; err != nil {
		logger.Error("count success messages failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	var failed int64
	if err := query.Where("status = ?", models.MessageStatusFailed).Count(&failed).Error; err != nil {
		logger.Error("count failed messages failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	var pending int64
	if err := query.Where("status IN (?)", []models.MessageStatus{
		models.MessageStatusQueued,
		models.MessageStatusSending,
		models.MessageStatusRetrying,
	}).Count(&pending).Error; err != nil {
		logger.Error("count pending messages failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	type DurationResult struct {
		Total   int64
		Count   int64
	}
	var durationResult DurationResult
	if err := query.Select("COALESCE(SUM(duration_ms), 0) as total, COUNT(*) as count").
		Where("status = ? AND duration_ms > 0", models.MessageStatusSent).
		Scan(&durationResult).Error; err != nil {
		logger.Error("get duration stats failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	successRate := 0.0
	if total > 0 {
		successRate = float64(success) / float64(total) * 100
	}

	avgDuration := int64(0)
	if durationResult.Count > 0 {
		avgDuration = durationResult.Total / durationResult.Count
	}

	return &MessageStats{
		Total:         total,
		Success:       success,
		Failed:        failed,
		Pending:       pending,
		SuccessRate:   successRate,
		TotalDuration: durationResult.Total,
		AvgDuration:   avgDuration,
	}, nil
}

func (ss *StatisticsService) GetChannelStats(start, end time.Time) ([]ChannelStats, error) {
	type Result struct {
		ChannelID uint
		Total     int64
		Success   int64
		Failed    int64
	}

	var results []Result
	query := database.DB.Model(&models.Message{}).
		Select("channel_id, COUNT(*) as total, " +
			"SUM(CASE WHEN status = 'sent' THEN 1 ELSE 0 END) as success, " +
			"SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed")

	if !start.IsZero() {
		query = query.Where("created_at >= ?", start)
	}
	if !end.IsZero() {
		query = query.Where("created_at <= ?", end)
	}

	if err := query.Group("channel_id").Scan(&results).Error; err != nil {
		logger.Error("get channel stats failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	stats := make([]ChannelStats, 0, len(results))
	for _, r := range results {
		name := ""
		if ss.channelService != nil {
			channel, err := ss.channelService.GetByID(r.ChannelID)
			if err == nil {
				name = channel.Name
			}
		}
		stats = append(stats, ChannelStats{
			ChannelID:   r.ChannelID,
			ChannelName: name,
			Total:       r.Total,
			Success:     r.Success,
			Failed:      r.Failed,
		})
	}

	return stats, nil
}

func (ss *StatisticsService) GetMessageList(channelID uint, status models.MessageStatus, start, end time.Time, page, pageSize int) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64

	query := database.DB.Model(&models.Message{}).Preload("Channel")

	if channelID > 0 {
		query = query.Where("channel_id = ?", channelID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if !start.IsZero() {
		query = query.Where("created_at >= ?", start)
	}
	if !end.IsZero() {
		query = query.Where("created_at <= ?", end)
	}

	if err := query.Count(&total).Error; err != nil {
		logger.Error("count messages failed", zap.Error(err))
		return nil, 0, errors.DatabaseError(err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		logger.Error("list messages failed", zap.Error(err))
		return nil, 0, errors.DatabaseError(err)
	}

	return messages, total, nil
}

func (ss *StatisticsService) GetMessageByID(messageID string) (*models.Message, error) {
	var message models.Message
	if err := database.DB.Preload("Channel").Preload("Template").Where("message_id = ?", messageID).First(&message).Error; err != nil {
		logger.Error("get message failed", zap.Error(err), zap.String("message_id", messageID))
		return nil, errors.NotFound("message", err)
	}
	return &message, nil
}

type DailyStats struct {
	Date    string `json:"date"`
	Total   int64  `json:"total"`
	Success int64  `json:"success"`
	Failed  int64  `json:"failed"`
}

func (ss *StatisticsService) GetDailyStats(channelID uint, days int) ([]DailyStats, error) {
	type Result struct {
		Date    string
		Total   int64
		Success int64
		Failed  int64
	}

	var results []Result
	startDate := time.Now().AddDate(0, 0, -days)

	query := database.DB.Model(&models.Message{}).
		Select("DATE(created_at) as date, COUNT(*) as total, " +
			"SUM(CASE WHEN status = 'sent' THEN 1 ELSE 0 END) as success, " +
			"SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed")

	if channelID > 0 {
		query = query.Where("channel_id = ?", channelID)
	}

	if err := query.Where("created_at >= ?", startDate).
		Group("DATE(created_at)").
		Order("date DESC").
		Scan(&results).Error; err != nil {
		logger.Error("get daily stats failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	stats := make([]DailyStats, len(results))
	for i, r := range results {
		stats[i] = DailyStats{
			Date:    r.Date,
			Total:   r.Total,
			Success: r.Success,
			Failed:  r.Failed,
		}
	}

	return stats, nil
}
