package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"consultation-platform/config"
	"consultation-platform/models"
	"consultation-platform/repositories"
	"consultation-platform/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppointmentService struct {
	appointmentRepo *repositories.AppointmentRepository
	paymentRepo     *repositories.PaymentRepository
	scheduleRepo    *repositories.ScheduleRepository
	serviceRepo     *repositories.ServiceRepository
	userRepo        *repositories.UserRepository
	cfg             *config.Config
}

func NewAppointmentService(cfg *config.Config) *AppointmentService {
	return &AppointmentService{
		appointmentRepo: repositories.NewAppointmentRepository(),
		paymentRepo:     repositories.NewPaymentRepository(),
		scheduleRepo:    repositories.NewScheduleRepository(),
		serviceRepo:     repositories.NewServiceRepository(),
		userRepo:        repositories.NewUserRepository(),
		cfg:             cfg,
	}
}

type CreateAppointmentRequest struct {
	ServiceID    uuid.UUID `json:"service_id" binding:"required"`
	ScheduleID   uuid.UUID `json:"schedule_id" binding:"required"`
	Notes        string    `json:"notes"`
}

type ConfirmAppointmentRequest struct {
	AppointmentID uuid.UUID `json:"appointment_id" binding:"required"`
}

type CancelAppointmentRequest struct {
	AppointmentID uuid.UUID `json:"appointment_id" binding:"required"`
	Reason        string    `json:"reason" binding:"required"`
}

type AppointmentServiceInterface interface {
	CreateAppointment(clientID uuid.UUID, req *CreateAppointmentRequest) (*models.Appointment, *models.Payment, error)
	GetAppointmentByID(id uuid.UUID) (*models.Appointment, error)
	GetClientAppointments(clientID uuid.UUID, page, pageSize int, status string) ([]models.Appointment, int64, error)
	GetProfessionalAppointments(professionalID uuid.UUID, page, pageSize int, status string) ([]models.Appointment, int64, error)
	ConfirmAppointment(professionalID uuid.UUID, req *ConfirmAppointmentRequest) (*models.Appointment, error)
	CancelAppointment(userID uuid.UUID, userRole string, req *CancelAppointmentRequest) (*models.Appointment, error)
	CompleteAppointment(appointmentID uuid.UUID) (*models.Appointment, error)
	ProcessPayment(appointmentID uuid.UUID, transactionID string) error
	RefundPayment(appointmentID uuid.UUID, reason string) error
	CleanupExpiredAppointments() error
}

func (s *AppointmentService) CreateAppointment(clientID uuid.UUID, req *CreateAppointmentRequest) (*models.Appointment, *models.Payment, error) {
	lockKey := fmt.Sprintf("appointment:lock:%s", req.ScheduleID.String())
	lock := utils.NewDistributedLock(lockKey)

	ctx := context.Background()
	locked, err := lock.Lock(ctx, 30*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !locked {
		return nil, nil, errors.New("this time slot is being booked by another user, please try again")
	}
	defer lock.Unlock(ctx)

	service, err := s.serviceRepo.FindByID(req.ServiceID)
	if err != nil {
		return nil, nil, errors.New("service not found")
	}

	if service.Status != models.ServiceStatusActive {
		return nil, nil, errors.New("service is not available")
	}

	schedule, err := s.scheduleRepo.FindByID(req.ScheduleID)
	if err != nil {
		return nil, nil, errors.New("schedule not found")
	}

	if schedule.ServiceID != req.ServiceID {
		return nil, nil, errors.New("schedule does not belong to this service")
	}

	if !schedule.IsAvailable || schedule.IsBooked {
		return nil, nil, errors.New("this time slot is no longer available")
	}

	existingAppointment, _ := s.appointmentRepo.FindByScheduleID(req.ScheduleID)
	if existingAppointment != nil && existingAppointment.Status != models.AppointmentCancelled {
		return nil, nil, errors.New("this time slot is already booked")
	}

	var appointment *models.Appointment
	var payment *models.Payment

	err = utils.DB.Transaction(func(tx *gorm.DB) error {
		appointmentRepo := s.appointmentRepo.WithTx(tx)
		paymentRepo := s.paymentRepo.WithTx(tx)
		scheduleRepo := s.scheduleRepo.WithTx(tx)

		appointment = &models.Appointment{
			ClientID:       clientID,
			ProfessionalID: service.ProfessionalID,
			ServiceID:      req.ServiceID,
			ScheduleID:     req.ScheduleID,
			Status:         models.AppointmentPending,
			Notes:          req.Notes,
		}

		if err := appointmentRepo.Create(appointment); err != nil {
			return err
		}

		payment = &models.Payment{
			AppointmentID:  appointment.ID,
			ClientID:       clientID,
			ProfessionalID: service.ProfessionalID,
			Amount:         service.Price,
			PaymentMethod:  models.PaymentOnline,
			Status:         models.PaymentPending,
			ExpiresAt:      time.Now().Add(time.Duration(s.cfg.Payment.TimeoutMinutes) * time.Minute),
		}

		if err := paymentRepo.Create(payment); err != nil {
			return err
		}

		if err := scheduleRepo.MarkAsBooked(req.ScheduleID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return appointment, payment, nil
}

func (s *AppointmentService) GetAppointmentByID(id uuid.UUID) (*models.Appointment, error) {
	return s.appointmentRepo.FindByID(id)
}

func (s *AppointmentService) GetClientAppointments(clientID uuid.UUID, page, pageSize int, status string) ([]models.Appointment, int64, error) {
	return s.appointmentRepo.FindByClientID(clientID, page, pageSize, status)
}

func (s *AppointmentService) GetProfessionalAppointments(professionalID uuid.UUID, page, pageSize int, status string) ([]models.Appointment, int64, error) {
	return s.appointmentRepo.FindByProfessionalID(professionalID, page, pageSize, status)
}

func (s *AppointmentService) ConfirmAppointment(professionalID uuid.UUID, req *ConfirmAppointmentRequest) (*models.Appointment, error) {
	appointment, err := s.appointmentRepo.FindByID(req.AppointmentID)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	if appointment.ProfessionalID != professionalID {
		return nil, errors.New("you do not have permission to confirm this appointment")
	}

	if appointment.Status != models.AppointmentPending {
		return nil, errors.New("appointment cannot be confirmed")
	}

	err = s.appointmentRepo.UpdateStatus(appointment.ID, models.AppointmentConfirmed)
	if err != nil {
		return nil, err
	}

	appointment.Status = models.AppointmentConfirmed
	return appointment, nil
}

func (s *AppointmentService) CancelAppointment(userID uuid.UUID, userRole string, req *CancelAppointmentRequest) (*models.Appointment, error) {
	appointment, err := s.appointmentRepo.FindByID(req.AppointmentID)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	if userRole == string(models.RoleClient) && appointment.ClientID != userID {
		return nil, errors.New("you do not have permission to cancel this appointment")
	}

	if userRole == string(models.RoleProfessional) && appointment.ProfessionalID != userID {
		return nil, errors.New("you do not have permission to cancel this appointment")
	}

	if appointment.Status != models.AppointmentPending && appointment.Status != models.AppointmentConfirmed {
		return nil, errors.New("appointment cannot be cancelled")
	}

	err = utils.DB.Transaction(func(tx *gorm.DB) error {
		appointmentRepo := s.appointmentRepo.WithTx(tx)
		scheduleRepo := s.scheduleRepo.WithTx(tx)
		paymentRepo := s.paymentRepo.WithTx(tx)

		appointment.Status = models.AppointmentCancelled
		appointment.CancelReason = req.Reason
		now := time.Now()
		appointment.CancelledAt = &now

		if err := appointmentRepo.Update(appointment); err != nil {
			return err
		}

		if err := scheduleRepo.MarkAsAvailable(appointment.ScheduleID); err != nil {
			return err
		}

		payment, err := paymentRepo.FindByAppointmentID(appointment.ID)
		if err == nil && payment.Status == models.PaymentPaid {
			payment.Status = models.PaymentRefunded
			payment.RefundReason = req.Reason
			now := time.Now()
			payment.RefundedAt = &now
			if err := paymentRepo.Update(payment); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to cancel appointment: %w", err)
	}

	return appointment, nil
}

func (s *AppointmentService) CompleteAppointment(appointmentID uuid.UUID) (*models.Appointment, error) {
	appointment, err := s.appointmentRepo.FindByID(appointmentID)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	if appointment.Status != models.AppointmentConfirmed {
		return nil, errors.New("appointment cannot be completed")
	}

	err = s.appointmentRepo.UpdateStatus(appointmentID, models.AppointmentCompleted)
	if err != nil {
		return nil, err
	}

	appointment.Status = models.AppointmentCompleted
	return appointment, nil
}

func (s *AppointmentService) ProcessPayment(appointmentID uuid.UUID, transactionID string) error {
	appointment, err := s.appointmentRepo.FindByID(appointmentID)
	if err != nil {
		return errors.New("appointment not found")
	}

	payment, err := s.paymentRepo.FindByAppointmentID(appointmentID)
	if err != nil {
		return errors.New("payment not found")
	}

	if payment.Status != models.PaymentPending {
		return errors.New("payment has already been processed")
	}

	if time.Now().After(payment.ExpiresAt) {
		return errors.New("payment has expired")
	}

	err = utils.DB.Transaction(func(tx *gorm.DB) error {
		paymentRepo := s.paymentRepo.WithTx(tx)
		appointmentRepo := s.appointmentRepo.WithTx(tx)

		payment.Status = models.PaymentPaid
		payment.TransactionID = transactionID
		now := time.Now()
		payment.PaidAt = &now
		if err := paymentRepo.Update(payment); err != nil {
			return err
		}

		appointment.Status = models.AppointmentConfirmed
		if err := appointmentRepo.Update(appointment); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *AppointmentService) RefundPayment(appointmentID uuid.UUID, reason string) error {
	appointment, err := s.appointmentRepo.FindByID(appointmentID)
	if err != nil {
		return errors.New("appointment not found")
	}

	if appointment.Status != models.AppointmentConfirmed && appointment.Status != models.AppointmentCompleted {
		return errors.New("appointment is not eligible for refund")
	}

	payment, err := s.paymentRepo.FindByAppointmentID(appointmentID)
	if err != nil {
		return errors.New("payment not found")
	}

	if payment.Status != models.PaymentPaid {
		return errors.New("payment has not been made")
	}

	return utils.DB.Transaction(func(tx *gorm.DB) error {
		paymentRepo := s.paymentRepo.WithTx(tx)
		appointmentRepo := s.appointmentRepo.WithTx(tx)

		payment.Status = models.PaymentRefunded
		payment.RefundReason = reason
		now := time.Now()
		payment.RefundedAt = &now
		if err := paymentRepo.Update(payment); err != nil {
			return err
		}

		appointment.Status = models.AppointmentRefunded
		if err := appointmentRepo.Update(appointment); err != nil {
			return err
		}

		return nil
	})
}

func (s *AppointmentService) CleanupExpiredAppointments() error {
	expireDuration := time.Duration(s.cfg.Payment.TimeoutMinutes) * time.Minute
	expiredAppointments, err := s.appointmentRepo.FindExpiredPendingAppointments(expireDuration)
	if err != nil {
		return err
	}

	for _, appointment := range expiredAppointments {
		_ = utils.DB.Transaction(func(tx *gorm.DB) error {
			appointmentRepo := s.appointmentRepo.WithTx(tx)
			scheduleRepo := s.scheduleRepo.WithTx(tx)

			appointment.Status = models.AppointmentCancelled
			appointment.CancelReason = "Payment timeout"
			now := time.Now()
			appointment.CancelledAt = &now

			if err := appointmentRepo.Update(&appointment); err != nil {
				return err
			}

			if err := scheduleRepo.MarkAsAvailable(appointment.ScheduleID); err != nil {
				return err
			}

			return nil
		})
	}

	return nil
}
