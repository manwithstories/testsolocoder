package services

import (
	"medical-platform/internal/config"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Scheduler struct {
	cron *cron.Cron
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

func (s *Scheduler) Start() {
	s.cron.AddFunc("0 0 9 * * *", func() {
		config.Logger.Info("开始发送就诊提醒")
		notificationService := NewNotificationService()
		if err := notificationService.SendAppointmentReminders(); err != nil {
			config.Logger.Error("发送就诊提醒失败", zap.Error(err))
		} else {
			config.Logger.Info("就诊提醒发送完成")
		}
	})

	s.cron.AddFunc("0 30 20 * * *", func() {
		config.Logger.Info("开始发送明天就诊提醒")
		notificationService := NewNotificationService()
		if err := notificationService.SendAppointmentReminders(); err != nil {
			config.Logger.Error("发送明天就诊提醒失败", zap.Error(err))
		} else {
			config.Logger.Info("明天就诊提醒发送完成")
		}
	})

	s.cron.Start()
	config.Logger.Info("定时任务调度器已启动")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	config.Logger.Info("定时任务调度器已停止")
}
