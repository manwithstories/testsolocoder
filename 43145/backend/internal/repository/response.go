package repository

import (
	"errors"
	"survey-platform/internal/model"
	"time"

	"gorm.io/gorm"
)

type ResponseRepository struct {
	db *gorm.DB
}

func NewResponseRepository(db *gorm.DB) *ResponseRepository {
	return &ResponseRepository{db: db}
}

func (r *ResponseRepository) Create(response *model.Response) error {
	return r.db.Create(response).Error
}

func (r *ResponseRepository) FindByID(id uint) (*model.Response, error) {
	var response model.Response
	err := r.db.Preload("Answers").First(&response, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &response, nil
}

func (r *ResponseRepository) FindBySessionID(surveyID uint, sessionID string) (*model.Response, error) {
	var response model.Response
	err := r.db.Preload("Answers").
		Where("survey_id = ? AND session_id = ?", surveyID, sessionID).
		First(&response).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &response, nil
}

func (r *ResponseRepository) Update(response *model.Response) error {
	return r.db.Save(response).Error
}

func (r *ResponseRepository) SaveAnswer(answer *model.Answer) error {
	return r.db.Save(answer).Error
}

func (r *ResponseRepository) SaveAnswers(responseID uint, answers []*model.Answer) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("response_id = ?", responseID).Delete(&model.Answer{}).Error; err != nil {
			return err
		}
		for i := range answers {
			answers[i].ResponseID = responseID
			if err := tx.Create(answers[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ResponseRepository) CompleteResponse(responseID uint) error {
	now := time.Now()
	var response model.Response
	if err := r.db.First(&response, responseID).Error; err != nil {
		return err
	}

	duration := 0
	if response.StartTime != nil {
		duration = int(now.Sub(*response.StartTime).Seconds())
	}

	return r.db.Model(&model.Response{}).Where("id = ?", responseID).
		Updates(map[string]interface{}{
			"status":        2,
			"complete_time": now,
			"duration":      duration,
		}).Error
}

func (r *ResponseRepository) List(surveyID uint, page, pageSize int, status int, startDate, endDate *time.Time, channel, sortBy, sortOrder string) ([]*model.Response, int64, error) {
	var responses []*model.Response
	var total int64

	query := r.db.Model(&model.Response{}).Where("survey_id = ?", surveyID)
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := sortBy + " " + sortOrder
	if err := query.Preload("Answers").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order(order).Find(&responses).Error; err != nil {
		return nil, 0, err
	}

	return responses, total, nil
}

func (r *ResponseRepository) CountBySurveyID(surveyID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Response{}).Where("survey_id = ?", surveyID).Count(&count).Error
	return count, err
}

func (r *ResponseRepository) CountByStatus(surveyID uint, status int) (int64, error) {
	var count int64
	err := r.db.Model(&model.Response{}).Where("survey_id = ? AND status = ?", surveyID, status).Count(&count).Error
	return count, err
}

func (r *ResponseRepository) Delete(id uint) error {
	return r.db.Delete(&model.Response{}, id).Error
}

func (r *ResponseRepository) FindByUserAndSurvey(userID uint, surveyID uint) ([]*model.Response, error) {
	var responses []*model.Response
	err := r.db.Where("user_id = ? AND survey_id = ?", userID, surveyID).Find(&responses).Error
	return responses, err
}

func (r *ResponseRepository) CountByIPAndSurvey(ipAddress string, surveyID uint, since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&model.Response{}).
		Where("ip_address = ? AND survey_id = ? AND created_at >= ?", ipAddress, surveyID, since).
		Count(&count).Error
	return count, err
}
