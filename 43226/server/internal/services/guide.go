package services

import (
	"fmt"
	"time"

	"museum-server/internal/dto"
	"museum-server/internal/models"
	"museum-server/internal/repository"
)

type GuideService struct {
	guideRepo *repository.GuideRepository
}

func NewGuideService(guideRepo *repository.GuideRepository) *GuideService {
	return &GuideService{guideRepo: guideRepo}
}

func (s *GuideService) CreateSchedule(guideID uint, req *dto.GuideScheduleRequest) (*models.GuideSchedule, error) {
	schedule := &models.GuideSchedule{
		GuideID:     guideID,
		Date:        req.Date,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsAvailable: req.IsAvailable,
	}

	if err := s.guideRepo.CreateSchedule(schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return schedule, nil
}

func (s *GuideService) UpdateSchedule(id uint, req *dto.GuideScheduleRequest) error {
	schedule, err := s.guideRepo.FindScheduleByID(id)
	if err != nil {
		return fmt.Errorf("schedule not found")
	}

	schedule.Date = req.Date
	schedule.StartTime = req.StartTime
	schedule.EndTime = req.EndTime
	schedule.IsAvailable = req.IsAvailable

	return s.guideRepo.UpdateSchedule(schedule)
}

func (s *GuideService) DeleteSchedule(id uint) error {
	return s.guideRepo.DeleteSchedule(id)
}

func (s *GuideService) ListSchedules(guideID uint, startDate, endDate time.Time) ([]models.GuideSchedule, error) {
	return s.guideRepo.ListSchedules(guideID, startDate, endDate)
}

func (s *GuideService) CreateContent(req *dto.GuideContentRequest) (*models.GuideContent, error) {
	content := &models.GuideContent{
		CollectionID: req.CollectionID,
		ExhibitionID: req.ExhibitionID,
		Language:     req.Language,
		Content:      req.Content,
		AudioUrl:     req.AudioUrl,
		SortOrder:    req.SortOrder,
	}

	if err := s.guideRepo.CreateContent(content); err != nil {
		return nil, fmt.Errorf("failed to create content: %w", err)
	}

	return content, nil
}

func (s *GuideService) UpdateContent(id uint, req *dto.GuideContentRequest) error {
	content, err := s.guideRepo.FindContentByID(id)
	if err != nil {
		return fmt.Errorf("content not found")
	}

	content.CollectionID = req.CollectionID
	content.ExhibitionID = req.ExhibitionID
	content.Language = req.Language
	content.Content = req.Content
	content.AudioUrl = req.AudioUrl
	content.SortOrder = req.SortOrder

	return s.guideRepo.UpdateContent(content)
}

func (s *GuideService) DeleteContent(id uint) error {
	return s.guideRepo.DeleteContent(id)
}

func (s *GuideService) ListContents(collectionID, exhibitionID uint, language string) ([]models.GuideContent, error) {
	return s.guideRepo.ListContents(collectionID, exhibitionID, language)
}

type ResearchService struct {
	researchRepo *repository.ResearchRepository
}

func NewResearchService(researchRepo *repository.ResearchRepository) *ResearchService {
	return &ResearchService{researchRepo: researchRepo}
}

func (s *ResearchService) CreateApplication(userID uint, req *dto.ResearchApplicationRequest) (*models.ResearchApplication, error) {
	app := &models.ResearchApplication{
		UserID:       userID,
		CollectionID: req.CollectionID,
		Purpose:      req.Purpose,
		Institution:  req.Institution,
		Status:       models.ApplicationStatusPending,
	}

	if err := s.researchRepo.CreateApplication(app); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return app, nil
}

func (s *ResearchService) ReviewApplication(id, reviewerID uint, req *dto.ApplicationReviewRequest) error {
	app, err := s.researchRepo.FindApplicationByID(id)
	if err != nil {
		return fmt.Errorf("application not found")
	}

	if app.Status != models.ApplicationStatusPending {
		return fmt.Errorf("application already reviewed")
	}

	now := time.Now()
	app.Status = req.Status
	app.ReviewerID = &reviewerID
	app.ReviewComment = req.ReviewComment
	app.ReviewedAt = &now

	if req.Status == models.ApplicationStatusApproved {
		app.ApprovedAt = &now
	}

	return s.researchRepo.UpdateApplication(app)
}

func (s *ResearchService) ListApplications(status string, page, pageSize int) ([]models.ResearchApplication, int64, error) {
	return s.researchRepo.ListApplications(status, page, pageSize)
}

func (s *ResearchService) ListApplicationsByUser(userID uint, page, pageSize int) ([]models.ResearchApplication, int64, error) {
	return s.researchRepo.ListApplicationsByUser(userID, page, pageSize)
}

func (s *ResearchService) GetApplication(id uint) (*models.ResearchApplication, error) {
	return s.researchRepo.FindApplicationByID(id)
}
