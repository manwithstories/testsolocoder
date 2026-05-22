package service

import (
	"beauty-salon-system/internal/model"
	"beauty-salon-system/internal/repository"
)

type AuditService struct {
	auditRepo *repository.AuditLogRepository
}

func NewAuditService(auditRepo *repository.AuditLogRepository) *AuditService {
	return &AuditService{
		auditRepo: auditRepo,
	}
}

func (s *AuditService) Log(userID uint, action, module, detail, ip string) error {
	log := &model.AuditLog{
		UserID:  userID,
		Action:  action,
		Module:  module,
		Detail:  detail,
		IP:      ip,
	}
	return s.auditRepo.Create(log)
}

func (s *AuditService) List(page, pageSize int, filters map[string]interface{}) ([]model.AuditLog, int64, error) {
	return s.auditRepo.List(page, pageSize, filters)
}

type ReviewService struct {
	reviewRepo    *repository.ReviewRepository
	technicianRepo *repository.TechnicianRepository
}

func NewReviewService(reviewRepo *repository.ReviewRepository, technicianRepo *repository.TechnicianRepository) *ReviewService {
	return &ReviewService{
		reviewRepo:    reviewRepo,
		technicianRepo: technicianRepo,
	}
}

type CreateReviewRequest struct {
	AppointmentID uint   `json:"appointment_id" binding:"required"`
	CustomerID    uint   `json:"customer_id" binding:"required"`
	TechnicianID  uint   `json:"technician_id" binding:"required"`
	ServiceID     uint   `json:"service_id" binding:"required"`
	Rating        int    `json:"rating" binding:"required,min=1,max=5"`
	Content       string `json:"content"`
}

func (s *ReviewService) Create(req *CreateReviewRequest) (*model.Review, error) {
	existing, _ := s.reviewRepo.GetByAppointmentID(req.AppointmentID)
	if existing != nil {
		return nil, nil
	}

	review := &model.Review{
		AppointmentID: req.AppointmentID,
		CustomerID:    req.CustomerID,
		TechnicianID:  req.TechnicianID,
		ServiceID:     req.ServiceID,
		Rating:        req.Rating,
		Content:       req.Content,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	avgRating, totalCount, err := s.reviewRepo.GetAverageRating(req.TechnicianID)
	if err != nil {
		return review, nil
	}

	s.technicianRepo.UpdateRating(req.TechnicianID, avgRating, totalCount)

	return review, nil
}

func (s *ReviewService) GetByID(id uint) (*model.Review, error) {
	return s.reviewRepo.GetByID(id)
}

func (s *ReviewService) GetByTechnicianID(technicianID uint, page, pageSize int) ([]model.Review, int64, error) {
	return s.reviewRepo.GetByTechnicianID(technicianID, page, pageSize)
}
