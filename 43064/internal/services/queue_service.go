package services

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/notification-center/internal/config"
	"github.com/notification-center/internal/database"
	"github.com/notification-center/internal/errors"
	"github.com/notification-center/internal/logger"
	"github.com/notification-center/internal/models"
	"github.com/notification-center/internal/ratelimit"
	"github.com/notification-center/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type QueueService struct {
	cfg          *config.QueueConfig
	retryCfg     *config.RetryConfig
	rateLimiter   *ratelimit.RateLimiter
	channelSvc    *ChannelService
	senderSvc     *SenderService
	webhookSvc    *WebhookService
	workerWg       sync.WaitGroup
	stopCh        chan struct{}
	isRunning     bool
	mu            sync.Mutex
}

func NewQueueService(cfg *config.QueueConfig, retryCfg *config.RetryConfig, rateLimiter *ratelimit.RateLimiter,
	channelSvc *ChannelService, senderSvc *SenderService, webhookSvc *WebhookService) *QueueService {
	return &QueueService{
		cfg:        cfg,
		retryCfg:    retryCfg,
		rateLimiter:  rateLimiter,
		channelSvc:   channelSvc,
		senderSvc:    senderSvc,
		webhookSvc:   webhookSvc,
		stopCh:       make(chan struct{}),
	}
}

func (qs *QueueService) Start() {
	qs.mu.Lock()
	if qs.isRunning {
		qs.mu.Unlock()
		return
	}
	qs.isRunning = true
	qs.mu.Unlock()

	logger.Info("starting queue workers", zap.Int("worker_count", qs.cfg.WorkerCount))

	for i := 0; i < qs.cfg.WorkerCount; i++ {
		qs.workerWg.Add(1)
		go qs.worker(i)
	}

	qs.workerWg.Add(1)
	go qs.retryWorker()
}

func (qs *QueueService) Stop() {
	qs.mu.Lock()
	if !qs.isRunning {
		qs.mu.Unlock()
		return
	}
	qs.isRunning = false
	close(qs.stopCh)
	qs.mu.Unlock()

	qs.workerWg.Wait()
	logger.Info("queue workers stopped")
}

func (qs *QueueService) Enqueue(message *models.Message) error {
	if message.MessageID == "" {
		message.MessageID = models.GenerateMessageID()
	}

	message.Status = models.MessageStatusQueued

	if err := database.DB.Create(message).Error; err != nil {
		logger.Error("create message failed", zap.Error(err))
		return errors.DatabaseError(err)
	}

	queueItem := &models.MessageQueue{
		MessageID:   message.MessageID,
		Priority:  message.Priority,
		ScheduledAt: utils.ValueOrDefault(message.ScheduledAt, time.Now()),
	}

	if err := database.DB.Create(queueItem).Error; err != nil {
		logger.Error("enqueue message failed", zap.Error(err))
		return errors.QueueError("enqueue failed", err)
	}

	logger.Info("message enqueued", zap.String("message_id", message.MessageID), zap.Uint("channel_id", message.ChannelID))

	return nil
}

func (qs *QueueService) worker(id int) {
	defer qs.workerWg.Done()

	logger.Debug("queue worker started", zap.Int("worker_id", id))

	ticker := time.NewTicker(time.Duration(qs.cfg.PollInterval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-qs.stopCh:
			logger.Debug("queue worker stopping", zap.Int("worker_id", id))
			return
		case <-ticker.C:
			qs.processNext(id)
		}
	}
}

func (qs *QueueService) processNext(workerID int) {
	queueItem, err := qs.pickNextMessage(workerID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Error("pick next message failed", zap.Error(err))
		}
		return
	}

	if queueItem == nil {
		return
	}

	defer qs.releaseQueueItem(queueItem)

	var message models.Message
	if err := database.DB.Where("message_id = ?", queueItem.MessageID).First(&message).Error; err != nil {
		logger.Error("find message failed", zap.Error(err), zap.String("message_id", queueItem.MessageID))
		return
	}

	channel, err := qs.channelSvc.GetByID(message.ChannelID)
	if err != nil {
		logger.Error("get channel failed", zap.Error(err), zap.Uint("channel_id", message.ChannelID))
		qs.markMessageFailed(&message, err)
		return
	}

	if !channel.Enabled {
		logger.Warn("channel is disabled", zap.Uint("channel_id", message.ChannelID))
		qs.markMessageFailed(&message, errors.ChannelError("channel is disabled", nil))
		return
	}

	rps := 0.0
	burst := 0
	if channel.RateLimit != nil && channel.RateLimit.Enabled {
		rps = channel.RateLimit.RequestsPerSecond
		burst = channel.RateLimit.Burst
	}

	if !qs.rateLimiter.Allow(message.ChannelID, rps, burst) {
		logger.Warn("rate limit exceeded, requeue", zap.String("message_id", message.MessageID))
		nextTry := time.Now().Add(time.Second)
		queueItem.ScheduledAt = nextTry
		database.DB.Model(queueItem).Update("scheduled_at", nextTry)
		return
	}

	startTime := time.Now()
	message.Status = models.MessageStatusSending
	database.DB.Model(&message).Update("status", models.MessageStatusSending)

	sendErr := qs.senderSvc.Send(&message, channel)

	duration := time.Since(startTime).Milliseconds()
	message.DurationMs = duration

	if sendErr != nil {
		logger.Error("send message failed", zap.Error(sendErr), zap.String("message_id", message.MessageID))
		qs.handleSendFailure(&message, sendErr)
	} else {
		message.Status = models.MessageStatusSent
		now := time.Now()
		message.SentAt = &now
		database.DB.Model(&message).Updates(map[string]interface{}{
			"status":   models.MessageStatusSent,
			"sent_at":  now,
			"duration_ms": duration,
		})
		logger.Info("message sent successfully", zap.String("message_id", message.MessageID))
	}

	if message.WebhookURL != "" {
		go qs.webhookSvc.Notify(message.WebhookURL, &message)
	}

	go qs.webhookSvc.NotifySubscribers(&message)
}

func (qs *QueueService) pickNextMessage(workerID int) (*models.MessageQueue, error) {
	now := time.Now()
	workerIDStr := fmt.Sprintf("worker-%d", workerID)

	var queueItem models.MessageQueue

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		subQuery := tx.Model(&models.MessageQueue{}).
			Where("scheduled_at <= ? AND (picked_at IS NULL OR picked_at < ?)", now, now.Add(-5*time.Minute)).
			Order("priority DESC, scheduled_at ASC, id ASC").
			Limit(1)

		result := tx.Model(&models.MessageQueue{}).
			Where("id IN (?)", subQuery.Select("id")).
			Updates(map[string]interface{}{
				"picked_at": now,
				"worker_id": workerIDStr,
			})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return tx.Where("worker_id = ? AND picked_at = ?", workerIDStr, now).
			First(&queueItem).Error
	})

	if err != nil {
		return nil, err
	}

	return &queueItem, nil
}

func (qs *QueueService) releaseQueueItem(item *models.MessageQueue) {
	database.DB.Delete(item)
}

func (qs *QueueService) handleSendFailure(message *models.Message, sendErr error) {
	message.RetryCount++
	message.LastError = sendErr.Error()

	if message.RetryCount < message.MaxRetries && qs.isRetryableError(sendErr) {
		backoff := utils.CalculateBackoff(
			message.RetryCount,
			qs.retryCfg.InitialBackoff,
			qs.retryCfg.MaxBackoff,
			qs.retryCfg.BackoffMultiplier,
		)

		nextRetry := time.Now().Add(backoff)
		message.NextRetryAt = &nextRetry
		message.Status = models.MessageStatusRetrying

		database.DB.Model(message).Updates(map[string]interface{}{
			"status":        models.MessageStatusRetrying,
			"retry_count":    message.RetryCount,
			"last_error":    message.LastError,
			"next_retry_at": nextRetry,
		})

		queueItem := &models.MessageQueue{
			MessageID:   message.MessageID,
			Priority:    message.Priority,
			ScheduledAt: nextRetry,
			RetryCount:  message.RetryCount,
		}
		database.DB.Create(queueItem)

		logger.Info("message scheduled for retry",
			zap.String("message_id", message.MessageID),
			zap.Int("retry_count", message.RetryCount),
			zap.Time("next_retry_at", nextRetry),
		)

		go qs.webhookSvc.NotifySubscribers(message)
	} else {
		qs.markMessageFailed(message, sendErr)
	}
}

func (qs *QueueService) markMessageFailed(message *models.Message, err error) {
	message.Status = models.MessageStatusFailed
	message.LastError = err.Error()
	message.ErrorStack = utils.GetStackTrace(err)

	database.DB.Model(message).Updates(map[string]interface{}{
		"status":      models.MessageStatusFailed,
		"retry_count": message.RetryCount,
		"last_error":   message.LastError,
		"error_stack": message.ErrorStack,
	})

	logger.Error("message failed permanently",
		zap.String("message_id", message.MessageID),
		zap.Error(err),
	)

	go qs.webhookSvc.NotifySubscribers(message)
}

func (qs *QueueService) isRetryableError(err error) bool {
	errStr := err.Error()
	for _, retryable := range qs.retryCfg.RetryableErrors {
		if strings.Contains(strings.ToLower(errStr), strings.ToLower(retryable)) {
			return true
		}
	}
	return errors.IsRetryable(err)
}

func (qs *QueueService) retryWorker() {
	defer qs.workerWg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-qs.stopCh:
			return
		case <-ticker.C:
			qs.cleanupStuckItems()
		}
	}
}

func (qs *QueueService) cleanupStuckItems() {
	cutoff := time.Now().Add(-10 * time.Minute)

	var stuckItems []models.MessageQueue
	database.DB.Where("picked_at IS NOT NULL AND picked_at < ?", cutoff).Find(&stuckItems)

	for _, item := range stuckItems {
		database.DB.Model(&item).Updates(map[string]interface{}{
			"picked_at": nil,
			"worker_id":  "",
		})
		logger.Warn("released stuck queue item", zap.String("message_id", item.MessageID))
	}
}

func (qs *QueueService) GetQueueStats() (pending int64, processing int64, err error) {
	database.DB.Model(&models.MessageQueue{}).Where("picked_at IS NULL").Count(&pending)
	database.DB.Model(&models.MessageQueue{}).Where("picked_at IS NOT NULL").Count(&processing)
	return
}
