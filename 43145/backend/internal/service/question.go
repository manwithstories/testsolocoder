package service

import (
	"errors"
	"fmt"
	"survey-platform/internal/dto"
	"survey-platform/internal/model"
	"survey-platform/internal/repository"
	"survey-platform/internal/utils"
)

type QuestionService struct {
	questionRepo *repository.QuestionRepository
	surveyRepo   *repository.SurveyRepository
	responseRepo *repository.ResponseRepository
}

func NewQuestionService(questionRepo *repository.QuestionRepository, surveyRepo *repository.SurveyRepository, responseRepo *repository.ResponseRepository) *QuestionService {
	return &QuestionService{
		questionRepo: questionRepo,
		surveyRepo:   surveyRepo,
		responseRepo: responseRepo,
	}
}

func (s *QuestionService) Create(surveyID, userID uint, req *dto.CreateQuestionRequest) (*model.Question, error) {
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return nil, err
	}
	if survey == nil {
		return nil, errors.New("survey not found")
	}
	if survey.UserID != userID {
		return nil, errors.New("permission denied")
	}
	if survey.Status == 3 {
		return nil, errors.New("cannot modify closed survey")
	}

	validTypes := map[string]bool{
		"single_choice": true, "multi_choice": true, "fill_in": true,
		"rating": true, "ranking": true, "matrix": true,
	}
	if !validTypes[req.Type] {
		return nil, errors.New("invalid question type")
	}

	question := &model.Question{
		SurveyID:       surveyID,
		Title:          req.Title,
		Type:           req.Type,
		IsRequired:     req.IsRequired,
		OrderIndex:     req.OrderIndex,
		Description:    req.Description,
		Placeholder:    req.Placeholder,
		MinValue:       req.MinValue,
		MaxValue:       req.MaxValue,
		DefaultValue:   req.DefaultValue,
		ValidationRule: req.ValidationRule,
		DisplayLogic:   req.DisplayLogic,
		Status:         1,
	}

	for _, opt := range req.Options {
		question.Options = append(question.Options, model.Option{
			Text:       opt.Text,
			OrderIndex: opt.OrderIndex,
			IsOther:    opt.IsOther,
			JumpTarget: opt.JumpTarget,
			Score:      opt.Score,
		})
	}

	for _, lj := range req.LogicJumps {
		question.LogicJumps = append(question.LogicJumps, model.LogicJump{
			Condition: lj.Condition,
			Value:     lj.Value,
			JumpTo:    lj.JumpTo,
		})
	}

	if err := s.validateQuestion(question); err != nil {
		return nil, err
	}

	if err := s.questionRepo.Create(question); err != nil {
		return nil, err
	}

	utils.CacheDelete(fmt.Sprintf("survey:%d", surveyID))
	return question, nil
}

func (s *QuestionService) GetByID(id uint) (*model.Question, error) {
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if question == nil {
		return nil, errors.New("question not found")
	}
	return question, nil
}

func (s *QuestionService) GetBySurveyID(surveyID uint) ([]*model.Question, error) {
	return s.questionRepo.FindBySurveyID(surveyID)
}

func (s *QuestionService) Update(id, userID uint, req *dto.UpdateQuestionRequest) error {
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		return err
	}
	if question == nil {
		return errors.New("question not found")
	}

	survey, err := s.surveyRepo.FindByID(question.SurveyID)
	if err != nil {
		return err
	}
	if survey == nil {
		return errors.New("survey not found")
	}
	if survey.UserID != userID {
		return errors.New("permission denied")
	}
	if survey.Status == 3 {
		return errors.New("cannot modify closed survey")
	}

	hasResponses, _ := s.responseRepo.CountBySurveyID(question.SurveyID)
	if hasResponses > 0 && req.Type != "" && req.Type != question.Type {
		return errors.New("cannot change question type after responses collected")
	}

	if req.Title != "" {
		question.Title = req.Title
	}
	if req.Type != "" {
		question.Type = req.Type
	}
	if req.IsRequired != nil {
		question.IsRequired = *req.IsRequired
	}
	if req.OrderIndex != nil {
		question.OrderIndex = *req.OrderIndex
	}
	if req.Description != "" {
		question.Description = req.Description
	}
	if req.Placeholder != "" {
		question.Placeholder = req.Placeholder
	}
	if req.MinValue != nil {
		question.MinValue = *req.MinValue
	}
	if req.MaxValue != nil {
		question.MaxValue = *req.MaxValue
	}
	if req.DefaultValue != "" {
		question.DefaultValue = req.DefaultValue
	}
	if req.ValidationRule != "" {
		question.ValidationRule = req.ValidationRule
	}
	if req.DisplayLogic != "" {
		question.DisplayLogic = req.DisplayLogic
	}

	question.Options = nil
	for _, opt := range req.Options {
		question.Options = append(question.Options, model.Option{
			Text:       opt.Text,
			OrderIndex: opt.OrderIndex,
			IsOther:    opt.IsOther,
			JumpTarget: opt.JumpTarget,
			Score:      opt.Score,
		})
	}

	question.LogicJumps = nil
	for _, lj := range req.LogicJumps {
		question.LogicJumps = append(question.LogicJumps, model.LogicJump{
			Condition: lj.Condition,
			Value:     lj.Value,
			JumpTo:    lj.JumpTo,
		})
	}

	if err := s.validateQuestion(question); err != nil {
		return err
	}

	utils.CacheDelete(fmt.Sprintf("survey:%d", question.SurveyID))
	return s.questionRepo.Update(question)
}

func (s *QuestionService) Delete(id, userID uint) error {
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		return err
	}
	if question == nil {
		return errors.New("question not found")
	}

	survey, err := s.surveyRepo.FindByID(question.SurveyID)
	if err != nil {
		return err
	}
	if survey == nil {
		return errors.New("survey not found")
	}
	if survey.UserID != userID {
		return errors.New("permission denied")
	}

	utils.CacheDelete(fmt.Sprintf("survey:%d", question.SurveyID))
	return s.questionRepo.Delete(id)
}

func (s *QuestionService) Reorder(id, userID uint, orderIndex int) error {
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		return err
	}
	if question == nil {
		return errors.New("question not found")
	}

	survey, err := s.surveyRepo.FindByID(question.SurveyID)
	if err != nil {
		return err
	}
	if survey == nil {
		return errors.New("survey not found")
	}
	if survey.UserID != userID {
		return errors.New("permission denied")
	}

	return s.questionRepo.UpdateOrderIndex(id, orderIndex)
}

func (s *QuestionService) BatchCreate(surveyID, userID uint, req *dto.BatchUpdateQuestionsRequest) error {
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return err
	}
	if survey == nil {
		return errors.New("survey not found")
	}
	if survey.UserID != userID {
		return errors.New("permission denied")
	}
	if survey.Status == 3 {
		return errors.New("cannot modify closed survey")
	}

	var questions []*model.Question
	for _, q := range req.Questions {
		question := &model.Question{
			SurveyID:       surveyID,
			Title:          q.Title,
			Type:           q.Type,
			IsRequired:     q.IsRequired,
			OrderIndex:     q.OrderIndex,
			Description:    q.Description,
			Placeholder:    q.Placeholder,
			MinValue:       q.MinValue,
			MaxValue:       q.MaxValue,
			DefaultValue:   q.DefaultValue,
			ValidationRule: q.ValidationRule,
			DisplayLogic:   q.DisplayLogic,
			Status:         1,
		}

		for _, opt := range q.Options {
			question.Options = append(question.Options, model.Option{
				Text:       opt.Text,
				OrderIndex: opt.OrderIndex,
				IsOther:    opt.IsOther,
				JumpTarget: opt.JumpTarget,
				Score:      opt.Score,
			})
		}

		for _, lj := range q.LogicJumps {
			question.LogicJumps = append(question.LogicJumps, model.LogicJump{
				Condition: lj.Condition,
				Value:     lj.Value,
				JumpTo:    lj.JumpTo,
			})
		}

		questions = append(questions, question)
	}

	utils.CacheDelete(fmt.Sprintf("survey:%d", surveyID))
	return s.questionRepo.BatchCreate(surveyID, questions)
}

func (s *QuestionService) validateQuestion(q *model.Question) error {
	if q.Title == "" {
		return errors.New("question title is required")
	}

	validTypes := map[string]bool{
		"single_choice": true, "multi_choice": true, "fill_in": true,
		"rating": true, "ranking": true, "matrix": true,
	}
	if !validTypes[q.Type] {
		return errors.New("invalid question type")
	}

	if (q.Type == "single_choice" || q.Type == "multi_choice" || q.Type == "ranking" || q.Type == "matrix") && len(q.Options) < 2 {
		return errors.New("at least 2 options required for this question type")
	}

	if q.Type == "rating" {
		if q.MinValue >= q.MaxValue {
			return errors.New("min value must be less than max value for rating")
		}
	}

	return nil
}
