package service

import (
	"errors"
	"fmt"
	"survey-platform/internal/dto"
	"survey-platform/internal/model"
	"survey-platform/internal/repository"
	"survey-platform/internal/utils"
	"time"
)

type SurveyService struct {
	surveyRepo *repository.SurveyRepository
}

func NewSurveyService(surveyRepo *repository.SurveyRepository) *SurveyService {
	return &SurveyService{surveyRepo: surveyRepo}
}

func (s *SurveyService) Create(userID uint, req *dto.CreateSurveyRequest) (*model.Survey, error) {
	survey := &model.Survey{
		Title:         req.Title,
		Description:   req.Description,
		CoverImage:    req.CoverImage,
		UserID:        userID,
		Status:        1,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Anonymous:     req.Anonymous,
		Password:      req.Password,
		MaxResponses:  req.MaxResponses,
		MaxPerUser:    req.MaxPerUser,
		RequiresLogin: req.RequiresLogin,
		AllowResume:   req.AllowResume,
		Category:      req.Category,
		Tags:          req.Tags,
	}

	if err := s.surveyRepo.Create(survey); err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("survey:%d", survey.ID)
	utils.CacheSetJSON(cacheKey, survey, 30*time.Minute)

	return survey, nil
}

func (s *SurveyService) GetByID(id uint) (*model.Survey, error) {
	cacheKey := fmt.Sprintf("survey:%d", id)
	var survey model.Survey
	if err := utils.CacheGetJSON(cacheKey, &survey); err == nil && survey.ID > 0 {
		return &survey, nil
	}

	surveyResult, err := s.surveyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if surveyResult == nil {
		return nil, errors.New("survey not found")
	}

	utils.CacheSetJSON(cacheKey, surveyResult, 30*time.Minute)
	return surveyResult, nil
}

func (s *SurveyService) Update(id, userID uint, req *dto.UpdateSurveyRequest) error {
	survey, err := s.surveyRepo.FindByID(id)
	if err != nil {
		return err
	}
	if survey == nil {
		return errors.New("survey not found")
	}

	if survey.UserID != userID {
		return errors.New("permission denied")
	}

	if survey.Status == 2 || survey.Status == 3 {
		if req.Title != "" || req.Description != "" {
			return errors.New("cannot modify title/description after publishing")
		}
	}

	if req.Title != "" {
		survey.Title = req.Title
	}
	if req.Description != "" {
		survey.Description = req.Description
	}
	if req.CoverImage != "" {
		survey.CoverImage = req.CoverImage
	}
	if req.StartTime != nil {
		survey.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		survey.EndTime = req.EndTime
	}
	if req.Anonymous != nil {
		survey.Anonymous = *req.Anonymous
	}
	if req.Password != "" {
		survey.Password = req.Password
	}
	survey.MaxResponses = req.MaxResponses
	survey.MaxPerUser = req.MaxPerUser
	if req.RequiresLogin != nil {
		survey.RequiresLogin = *req.RequiresLogin
	}
	if req.AllowResume != nil {
		survey.AllowResume = *req.AllowResume
	}
	if req.Category != "" {
		survey.Category = req.Category
	}
	if req.Tags != "" {
		survey.Tags = req.Tags
	}

	if err := s.surveyRepo.Update(survey); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("survey:%d", id)
	utils.CacheSetJSON(cacheKey, survey, 30*time.Minute)

	return nil
}

func (s *SurveyService) Delete(id, userID uint) error {
	survey, err := s.surveyRepo.FindByID(id)
	if err != nil {
		return err
	}
	if survey == nil {
		return errors.New("survey not found")
	}
	if survey.UserID != userID {
		return errors.New("permission denied")
	}

	utils.CacheDelete(fmt.Sprintf("survey:%d", id))
	return s.surveyRepo.Delete(id)
}

func (s *SurveyService) List(userID uint, query *dto.SurveyListQuery) ([]*model.Survey, int64, error) {
	return s.surveyRepo.List(userID, query.Page, query.PageSize, query.Status, query.Category, query.Keyword, query.SortBy, query.SortOrder)
}

func (s *SurveyService) ListAll(query *dto.SurveyListQuery) ([]*model.Survey, int64, error) {
	return s.surveyRepo.ListAll(query.Page, query.PageSize, query.Status, query.Category, query.Keyword, query.SortBy, query.SortOrder)
}

func (s *SurveyService) Publish(id, userID uint, req *dto.PublishSurveyRequest) error {
	survey, err := s.surveyRepo.FindByID(id)
	if err != nil {
		return err
	}
	if survey == nil {
		return errors.New("survey not found")
	}
	if survey.UserID != userID {
		return errors.New("permission denied")
	}
	if survey.Status != 1 {
		return errors.New("only draft surveys can be published")
	}

	if len(survey.Questions) == 0 {
		return errors.New("survey must have at least one question")
	}

	if req.StartTime != nil {
		survey.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		survey.EndTime = req.EndTime
	}
	survey.Status = 2

	if err := s.surveyRepo.Update(survey); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("survey:%d", id)
	utils.CacheSetJSON(cacheKey, survey, 30*time.Minute)

	return nil
}

func (s *SurveyService) Close(id, userID uint) error {
	survey, err := s.surveyRepo.FindByID(id)
	if err != nil {
		return err
	}
	if survey == nil {
		return errors.New("survey not found")
	}
	if survey.UserID != userID {
		return errors.New("permission denied")
	}

	survey.Status = 3

	if err := s.surveyRepo.Update(survey); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("survey:%d", id)
	utils.CacheSetJSON(cacheKey, survey, 30*time.Minute)

	return nil
}

func (s *SurveyService) Copy(id, userID uint, req *dto.CopySurveyRequest) (*model.Survey, error) {
	survey, err := s.surveyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if survey == nil {
		return nil, errors.New("survey not found")
	}

	newTitle := req.Title
	if newTitle == "" {
		newTitle = survey.Title + " (Copy)"
	}

	return s.surveyRepo.Copy(id, newTitle, userID)
}

func (s *SurveyService) CheckAccess(surveyID uint, password string, userID *uint) (bool, error) {
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return false, err
	}
	if survey == nil {
		return false, errors.New("survey not found")
	}

	if survey.IsClosed() {
		return false, errors.New("survey is closed")
	}

	if !survey.IsStarted() {
		return false, errors.New("survey has not started yet")
	}

	if survey.RequiresLogin && userID == nil {
		return false, errors.New("login required")
	}

	if survey.Password != "" && survey.Password != password {
		return false, errors.New("incorrect password")
	}

	return true, nil
}
