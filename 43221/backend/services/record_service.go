package services

import (
	"errors"
	"strings"
	"time"

	"consultation-platform/models"
	"consultation-platform/repositories"

	"github.com/google/uuid"
)

type RecordService struct {
	consultRecordRepo *repositories.ConsultRecordRepository
	reviewRepo        *repositories.ReviewRepository
	appointmentRepo   *repositories.AppointmentRepository
	serviceRepo       *repositories.ServiceRepository
	sensitiveWords    []string
}

func NewRecordService(sensitiveWords []string) *RecordService {
	return &RecordService{
		consultRecordRepo: repositories.NewConsultRecordRepository(),
		reviewRepo:        repositories.NewReviewRepository(),
		appointmentRepo:   repositories.NewAppointmentRepository(),
		serviceRepo:       repositories.NewServiceRepository(),
		sensitiveWords:    sensitiveWords,
	}
}

type CreateConsultRecordRequest struct {
	AppointmentID  uuid.UUID `json:"appointment_id" binding:"required"`
	Summary        string    `json:"summary"`
	Advice         string    `json:"advice"`
	FollowUpDate   string    `json:"follow_up_date"`
	IsConfidential bool      `json:"is_confidential"`
}

type CreateReviewRequest struct {
	AppointmentID uuid.UUID `json:"appointment_id" binding:"required"`
	Rating        int       `json:"rating" binding:"required,gte=1,lte=5"`
	Content       string    `json:"content"`
}

type UpdateReviewStatusRequest struct {
	ReviewID     uuid.UUID `json:"review_id" binding:"required"`
	Status       string    `json:"status" binding:"required,oneof=approved rejected"`
	RejectReason string    `json:"reject_reason"`
}

type RecordServiceInterface interface {
	CreateConsultRecord(professionalID uuid.UUID, req *CreateConsultRecordRequest) (*models.ConsultRecord, error)
	GetConsultRecordByID(id uuid.UUID) (*models.ConsultRecord, error)
	GetClientConsultRecords(clientID uuid.UUID, page, pageSize int) ([]models.ConsultRecord, int64, error)
	GetProfessionalConsultRecords(professionalID uuid.UUID, page, pageSize int) ([]models.ConsultRecord, int64, error)
	CreateReview(clientID uuid.UUID, req *CreateReviewRequest) (*models.Review, error)
	GetProfessionalReviews(professionalID uuid.UUID, page, pageSize int, status string) ([]models.Review, int64, error)
	GetServiceReviews(serviceID uuid.UUID, page, pageSize int) ([]models.Review, int64, error)
	GetPendingReviews(page, pageSize int) ([]models.Review, int64, error)
	UpdateReviewStatus(req *UpdateReviewStatusRequest) error
	GetProfessionalReviewStats(professionalID uuid.UUID) (float64, int, error)
	FilterSensitiveWords(content string) string
}

func (s *RecordService) CreateConsultRecord(professionalID uuid.UUID, req *CreateConsultRecordRequest) (*models.ConsultRecord, error) {
	appointment, err := s.appointmentRepo.FindByID(req.AppointmentID)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	if appointment.ProfessionalID != professionalID {
		return nil, errors.New("you do not have permission to create consult record for this appointment")
	}

	if appointment.Status != models.AppointmentCompleted {
		return nil, errors.New("can only create consult record for completed appointments")
	}

	existingRecord, _ := s.consultRecordRepo.FindByAppointmentID(req.AppointmentID)
	if existingRecord != nil {
		return nil, errors.New("consult record already exists for this appointment")
	}

	record := &models.ConsultRecord{
		AppointmentID:  req.AppointmentID,
		ClientID:       appointment.ClientID,
		ProfessionalID: professionalID,
		Summary:        req.Summary,
		Advice:         req.Advice,
		IsConfidential: req.IsConfidential,
	}

	if req.FollowUpDate != "" {
		record.FollowUpDate = parseTimeOrNil(req.FollowUpDate)
	}

	if err := s.consultRecordRepo.Create(record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *RecordService) GetConsultRecordByID(id uuid.UUID) (*models.ConsultRecord, error) {
	return s.consultRecordRepo.FindByID(id)
}

func (s *RecordService) GetClientConsultRecords(clientID uuid.UUID, page, pageSize int) ([]models.ConsultRecord, int64, error) {
	return s.consultRecordRepo.FindByClientID(clientID, page, pageSize)
}

func (s *RecordService) GetProfessionalConsultRecords(professionalID uuid.UUID, page, pageSize int) ([]models.ConsultRecord, int64, error) {
	return s.consultRecordRepo.FindByProfessionalID(professionalID, page, pageSize)
}

func (s *RecordService) CreateReview(clientID uuid.UUID, req *CreateReviewRequest) (*models.Review, error) {
	appointment, err := s.appointmentRepo.FindByID(req.AppointmentID)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	if appointment.ClientID != clientID {
		return nil, errors.New("you do not have permission to review this appointment")
	}

	if appointment.Status != models.AppointmentCompleted {
		return nil, errors.New("can only review completed appointments")
	}

	existingReview, _ := s.reviewRepo.FindByAppointmentID(req.AppointmentID)
	if existingReview != nil {
		return nil, errors.New("you have already reviewed this appointment")
	}

	filteredContent := s.FilterSensitiveWords(req.Content)

	review := &models.Review{
		AppointmentID:  req.AppointmentID,
		ClientID:       clientID,
		ProfessionalID: appointment.ProfessionalID,
		ServiceID:      appointment.ServiceID,
		Rating:         req.Rating,
		Content:        filteredContent,
		Status:         models.ReviewPending,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *RecordService) GetProfessionalReviews(professionalID uuid.UUID, page, pageSize int, status string) ([]models.Review, int64, error) {
	return s.reviewRepo.FindByProfessionalID(professionalID, page, pageSize, status)
}

func (s *RecordService) GetServiceReviews(serviceID uuid.UUID, page, pageSize int) ([]models.Review, int64, error) {
	return s.reviewRepo.FindByServiceID(serviceID, page, pageSize)
}

func (s *RecordService) GetPendingReviews(page, pageSize int) ([]models.Review, int64, error) {
	return s.reviewRepo.FindPendingReviews(page, pageSize)
}

func (s *RecordService) UpdateReviewStatus(req *UpdateReviewStatusRequest) error {
	status := models.ReviewStatus(req.Status)
	if err := s.reviewRepo.UpdateStatus(req.ReviewID, status, req.RejectReason); err != nil {
		return err
	}

	if status == models.ReviewApproved {
		review, err := s.reviewRepo.FindByID(req.ReviewID)
		if err != nil {
			return err
		}

		avgRating, count, err := s.reviewRepo.GetAverageRatingByServiceID(review.ServiceID)
		if err != nil {
			return err
		}

		_ = s.serviceRepo.UpdateRating(review.ServiceID, avgRating, count)
	}

	return nil
}

func (s *RecordService) GetProfessionalReviewStats(professionalID uuid.UUID) (float64, int, error) {
	return s.reviewRepo.GetAverageRatingByProfessionalID(professionalID)
}

func (s *RecordService) FilterSensitiveWords(content string) string {
	filtered := content
	for _, word := range s.sensitiveWords {
		filtered = strings.ReplaceAll(filtered, word, strings.Repeat("*", len([]rune(word))))
	}
	return filtered
}

func parseTimeOrNil(value string) *time.Time {
	return nil
}
