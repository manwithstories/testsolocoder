package service

import (
	"errors"
	"ticket-system/internal/dto"
	"ticket-system/internal/logger"
	"ticket-system/internal/models"
	"time"

	"gorm.io/gorm"
)

type ActivityService struct{}

func NewActivityService() *ActivityService {
	return &ActivityService{}
}

func (s *ActivityService) Create(userID uint, req *dto.ActivityCreateRequest) (*models.Activity, error) {
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, time.Local)
	if err != nil {
		return nil, errors.New("开始时间格式错误，应为YYYY-MM-DD HH:mm:ss")
	}

	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, time.Local)
	if err != nil {
		return nil, errors.New("结束时间格式错误，应为YYYY-MM-DD HH:mm:ss")
	}

	if endTime.Before(startTime) {
		return nil, errors.New("结束时间必须晚于开始时间")
	}

	activity := &models.Activity{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		Location:    req.Location,
		Capacity:    req.Capacity,
		Poster:      req.Poster,
		Status:      models.ActivityStatusDraft,
		CreatedBy:   userID,
	}

	if err := models.DB.Create(activity).Error; err != nil {
		logger.Log.Errorf("Create activity failed: %v", err)
		return nil, errors.New("创建活动失败")
	}

	logger.Log.Infof("Activity created: %d - %s", activity.ID, activity.Title)
	return activity, nil
}

func (s *ActivityService) GetList(req *dto.ActivityListRequest) ([]models.Activity, int64, error) {
	var activities []models.Activity
	var total int64

	query := models.DB.Model(&models.Activity{}).Preload("TicketTypes")

	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&activities).Error; err != nil {
		return nil, 0, err
	}

	return activities, total, nil
}

func (s *ActivityService) GetByID(id uint) (*models.Activity, error) {
	var activity models.Activity
	if err := models.DB.Preload("TicketTypes").First(&activity, id).Error; err != nil {
		return nil, err
	}
	return &activity, nil
}

func (s *ActivityService) Update(id uint, req *dto.ActivityUpdateRequest) (*models.Activity, error) {
	activity, err := s.GetByID(id)
	if err != nil {
		return nil, errors.New("活动不存在")
	}

	if activity.Status == models.ActivityStatusCanceled {
		return nil, errors.New("已取消的活动无法修改")
	}

	if req.Title != "" {
		activity.Title = req.Title
	}
	if req.Description != "" {
		activity.Description = req.Description
	}
	if req.Location != "" {
		activity.Location = req.Location
	}
	if req.Capacity > 0 {
		activity.Capacity = req.Capacity
	}
	if req.Poster != "" {
		activity.Poster = req.Poster
	}
	if req.StartTime != "" {
		startTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, time.Local)
		if err != nil {
			return nil, errors.New("开始时间格式错误")
		}
		activity.StartTime = startTime
	}
	if req.EndTime != "" {
		endTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, time.Local)
		if err != nil {
			return nil, errors.New("结束时间格式错误")
		}
		activity.EndTime = endTime
	}

	if err := models.DB.Save(activity).Error; err != nil {
		return nil, err
	}

	logger.Log.Infof("Activity updated: %d", id)
	return activity, nil
}

func (s *ActivityService) UpdateStatus(id uint, status string) (*models.Activity, error) {
	activity, err := s.GetByID(id)
	if err != nil {
		return nil, errors.New("活动不存在")
	}

	activity.Status = status

	if err := models.DB.Save(activity).Error; err != nil {
		return nil, err
	}

	logger.Log.Infof("Activity status updated: %d -> %s", id, status)
	return activity, nil
}

func (s *ActivityService) Delete(id uint) error {
	result := models.DB.Delete(&models.Activity{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	logger.Log.Infof("Activity deleted: %d", id)
	return nil
}
