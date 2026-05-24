package service

import (
	"errors"
	"fmt"
	"survey-platform/internal/dto"
	"survey-platform/internal/model"
	"survey-platform/internal/repository"
	"survey-platform/internal/utils"
	"time"

	"github.com/google/uuid"
)

type ResponseService struct {
	responseRepo *repository.ResponseRepository
	surveyRepo   *repository.SurveyRepository
	distRepo     *repository.DistributionRepository
}

func NewResponseService(
	responseRepo *repository.ResponseRepository,
	surveyRepo *repository.SurveyRepository,
	distRepo *repository.DistributionRepository,
) *ResponseService {
	return &ResponseService{
		responseRepo: responseRepo,
		surveyRepo:   surveyRepo,
		distRepo:     distRepo,
	}
}

func (s *ResponseService) StartResponse(surveyID uint, userID *uint, sessionID, ipAddress, userAgent string, distributionID *uint) (*model.Response, error) {
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return nil, err
	}
	if survey == nil {
		return nil, errors.New("survey not found")
	}

	if survey.IsClosed() {
		return nil, errors.New("survey is closed")
	}

	if !survey.IsStarted() {
		return nil, errors.New("survey has not started yet")
	}

	if survey.MaxResponses > 0 && survey.ResponseCount >= survey.MaxResponses {
		return nil, errors.New("survey has reached maximum responses")
	}

	if !survey.Anonymous && userID != nil && survey.MaxPerUser > 0 {
		userResponses, _ := s.responseRepo.FindByUserAndSurvey(*userID, surveyID)
		completedCount := 0
		for _, r := range userResponses {
			if r.Status == 2 {
				completedCount++
			}
		}
		if completedCount >= survey.MaxPerUser {
			return nil, errors.New("you have reached the maximum number of responses for this survey")
		}
	}

	rateLimitKey := "rate_limit:" + ipAddress
	countStr, _ := utils.CacheGet(rateLimitKey)
	if countStr == "1" {
		return nil, errors.New("too many requests, please try again later")
	}
	utils.CacheSet(rateLimitKey, "1", 5*time.Second)

	if sessionID != "" {
		existing, _ := s.responseRepo.FindBySessionID(surveyID, sessionID)
		if existing != nil && existing.Status == 1 && survey.AllowResume {
			return existing, nil
		}
	}

	now := time.Now()
	response := &model.Response{
		SurveyID:       surveyID,
		UserID:         userID,
		SessionID:      sessionID,
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		DistributionID: distributionID,
		Status:         1,
		StartTime:      &now,
	}

	if err := s.responseRepo.Create(response); err != nil {
		return nil, err
	}

	if distributionID != nil {
		s.distRepo.IncrementUseCount(*distributionID)
	}

	return response, nil
}

func (s *ResponseService) SaveProgress(surveyID uint, sessionID string, answers []dto.SaveAnswerRequest) error {
	response, err := s.responseRepo.FindBySessionID(surveyID, sessionID)
	if err != nil {
		return err
	}
	if response == nil {
		return errors.New("response not found")
	}
	if response.Status != 1 {
		return errors.New("response already completed")
	}

	var answerModels []*model.Answer
	for _, a := range answers {
		answerModels = append(answerModels, &model.Answer{
			QuestionID:   a.QuestionID,
			OptionID:     a.OptionID,
			TextValue:    a.TextValue,
			NumericValue: a.NumericValue,
			MatrixValues: a.MatrixValues,
			RankingOrder: a.RankingOrder,
		})
	}

	return s.responseRepo.SaveAnswers(response.ID, answerModels)
}

func (s *ResponseService) SubmitResponse(surveyID uint, req *dto.SubmitResponseRequest) error {
	response, err := s.responseRepo.FindBySessionID(surveyID, req.SessionID)
	if err != nil {
		return err
	}
	if response == nil {
		return errors.New("response not found")
	}
	if response.Status == 2 {
		return errors.New("response already submitted")
	}

	var answerModels []*model.Answer
	for _, a := range req.Answers {
		answerModels = append(answerModels, &model.Answer{
			QuestionID:   a.QuestionID,
			OptionID:     a.OptionID,
			TextValue:    a.TextValue,
			NumericValue: a.NumericValue,
			MatrixValues: a.MatrixValues,
			RankingOrder: a.RankingOrder,
		})
	}

	if err := s.responseRepo.SaveAnswers(response.ID, answerModels); err != nil {
		return err
	}

	if err := s.responseRepo.CompleteResponse(response.ID); err != nil {
		return err
	}

	if err := s.surveyRepo.IncrementResponseCount(surveyID); err != nil {
		return err
	}

	utils.CacheDeleteByPattern(fmt.Sprintf("stats:%d:*", surveyID))

	return nil
}

func (s *ResponseService) GetByID(id uint) (*model.Response, error) {
	response, err := s.responseRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("response not found")
	}
	return response, nil
}

func (s *ResponseService) GetDetail(id uint) (*dto.ResponseDetailResponse, error) {
	response, err := s.responseRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("response not found")
	}

	var answers []dto.AnswerResponse
	for _, a := range response.Answers {
		answers = append(answers, dto.AnswerResponse{
			QuestionID:   a.QuestionID,
			OptionID:     a.OptionID,
			TextValue:    a.TextValue,
			NumericValue: a.NumericValue,
			MatrixValues: a.MatrixValues,
			RankingOrder: a.RankingOrder,
		})
	}

	return &dto.ResponseDetailResponse{
		ID:           response.ID,
		SurveyID:     response.SurveyID,
		UserID:       response.UserID,
		SessionID:    response.SessionID,
		Status:       response.Status,
		StartTime:    response.StartTime,
		CompleteTime: response.CompleteTime,
		Duration:     response.Duration,
		Answers:      answers,
		CreatedAt:    response.CreatedAt,
	}, nil
}

func (s *ResponseService) List(surveyID uint, query *dto.ResponseListQuery) ([]*model.Response, int64, error) {
	return s.responseRepo.List(surveyID, query.Page, query.PageSize, query.Status,
		query.StartDate, query.EndDate, query.Channel, query.SortBy, query.SortOrder)
}

func (s *ResponseService) Delete(id uint) error {
	return s.responseRepo.Delete(id)
}

func (s *ResponseService) ValidateSurveyAccess(surveyID uint, password string, userID *uint) error {
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return err
	}
	if survey == nil {
		return errors.New("survey not found")
	}

	if survey.IsClosed() {
		return errors.New("survey is closed")
	}
	if !survey.IsStarted() {
		return errors.New("survey has not started yet")
	}
	if survey.RequiresLogin && userID == nil {
		return errors.New("login required to access this survey")
	}
	if survey.Password != "" && survey.Password != password {
		return errors.New("incorrect password")
	}

	return nil
}

func (s *ResponseService) GenerateSessionID() string {
	return uuid.New().String()
}
