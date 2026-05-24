package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"skillshare/internal/models"
	"skillshare/internal/repository"
	"skillshare/pkg/validator"
)

type BookingService struct {
	bookingRepo *repository.BookingRepository
	reviewRepo  *repository.ReviewRepository
	skillRepo   *repository.SkillRepository
	userRepo    *repository.UserRepository
}

func NewBookingService(bookingRepo *repository.BookingRepository, reviewRepo *repository.ReviewRepository, skillRepo *repository.SkillRepository, userRepo *repository.UserRepository) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		reviewRepo:  reviewRepo,
		skillRepo:   skillRepo,
		userRepo:    userRepo,
	}
}

type CreateBookingInput struct {
	PostingID      uuid.UUID `json:"posting_id"`
	ScheduledStart time.Time `json:"scheduled_start"`
	ScheduledEnd   time.Time `json:"scheduled_end"`
	Note           string    `json:"note"`
}

func (s *BookingService) CreateBooking(studentID uuid.UUID, input *CreateBookingInput) (*models.Booking, error) {
	if err := validator.ValidateTimeSlot(input.ScheduledStart, input.ScheduledEnd); err != nil {
		return nil, err
	}

	posting, err := s.skillRepo.FindPostingByID(input.PostingID)
	if err != nil {
		return nil, errors.New("课程不存在")
	}

	if !posting.IsActive {
		return nil, errors.New("课程已下架")
	}

	conflictingBookings, err := s.bookingRepo.GetConflictingBookings(posting.TeacherID, input.ScheduledStart, input.ScheduledEnd)
	if err != nil {
		return nil, err
	}
	if len(conflictingBookings) > 0 {
		return nil, errors.New("该时间段已有预约")
	}

	duration := input.ScheduledEnd.Sub(input.ScheduledStart).Hours()
	price := posting.PricePerHour * duration
	platformFee := price * 0.05
	teacherEarnings := price - platformFee

	booking := &models.Booking{
		PostingID:       input.PostingID,
		StudentID:       studentID,
		TeacherID:       posting.TeacherID,
		ScheduledStart:  input.ScheduledStart,
		ScheduledEnd:    input.ScheduledEnd,
		Status:          models.BookingStatusPending,
		Price:           price,
		PlatformFee:     platformFee,
		TeacherEarnings: teacherEarnings,
		Note:            input.Note,
	}

	if err := s.bookingRepo.Create(booking); err != nil {
		return nil, errors.New("创建预约失败")
	}

	return booking, nil
}

func (s *BookingService) GetBooking(id uuid.UUID) (*models.Booking, error) {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("预约不存在")
	}
	return booking, nil
}

func (s *BookingService) ListBookings(userID uuid.UUID, role string, page, pageSize int, status *string) ([]*models.Booking, int64, error) {
	return s.bookingRepo.ListByUser(userID, role, page, pageSize, status)
}

func (s *BookingService) ConfirmBooking(id uuid.UUID, userID uuid.UUID) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return errors.New("预约不存在")
	}

	if booking.TeacherID != userID {
		return errors.New("无权限操作")
	}

	if booking.Status != models.BookingStatusPending {
		return errors.New("预约状态不允许确认")
	}

	booking.Status = models.BookingStatusConfirmed
	return s.bookingRepo.Update(booking)
}

func (s *BookingService) RejectBooking(id uuid.UUID, userID uuid.UUID, reason string) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return errors.New("预约不存在")
	}

	if booking.TeacherID != userID {
		return errors.New("无权限操作")
	}

	if booking.Status != models.BookingStatusPending {
		return errors.New("预约状态不允许拒绝")
	}

	booking.Status = models.BookingStatusRejected
	booking.RejectReason = reason
	return s.bookingRepo.Update(booking)
}

func (s *BookingService) CancelBooking(id uuid.UUID, userID uuid.UUID, reason string) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return errors.New("预约不存在")
	}

	if booking.StudentID != userID && booking.TeacherID != userID {
		return errors.New("无权限操作")
	}

	if booking.Status != models.BookingStatusPending && booking.Status != models.BookingStatusConfirmed {
		return errors.New("预约状态不允许取消")
	}

	booking.Status = models.BookingStatusCancelled
	booking.CancelReason = reason
	booking.CancelledBy = &userID
	return s.bookingRepo.Update(booking)
}

func (s *BookingService) CompleteBooking(id uuid.UUID, userID uuid.UUID) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return errors.New("预约不存在")
	}

	if booking.TeacherID != userID {
		return errors.New("无权限操作")
	}

	if booking.Status != models.BookingStatusConfirmed {
		return errors.New("预约状态不允许完成")
	}

	now := time.Now()
	actualStart := booking.ScheduledStart
	actualEnd := now

	if actualEnd.Before(booking.ScheduledEnd) {
		actualEnd = booking.ScheduledEnd
	}

	duration := actualEnd.Sub(actualStart).Hours()

	if err := s.bookingRepo.CompleteBooking(id, actualStart, actualEnd); err != nil {
		return err
	}

	s.skillRepo.IncrementPostingStats(booking.PostingID, duration)
	s.userRepo.UpdateRating(booking.TeacherID)

	return nil
}

type CreateReviewInput struct {
	BookingID uuid.UUID `json:"booking_id"`
	Rating    int       `json:"rating"`
	Content   string    `json:"content"`
	IsPublic  bool      `json:"is_public"`
}

func (s *BookingService) CreateReview(reviewerID uuid.UUID, input *CreateReviewInput) (*models.Review, error) {
	if err := validator.ValidateRating(input.Rating); err != nil {
		return nil, err
	}

	if err := validator.ValidateReviewContent(input.Content); err != nil {
		return nil, err
	}

	booking, err := s.bookingRepo.FindByID(input.BookingID)
	if err != nil {
		return nil, errors.New("预约不存在")
	}

	if booking.Status != models.BookingStatusCompleted {
		return nil, errors.New("只有已完成的预约才能评价")
	}

	var revieweeID uuid.UUID
	if booking.StudentID == reviewerID {
		if booking.ReviewedByStudent {
			return nil, errors.New("已评价过")
		}
		revieweeID = booking.TeacherID
	} else if booking.TeacherID == reviewerID {
		if booking.ReviewedByTeacher {
			return nil, errors.New("已评价过")
		}
		revieweeID = booking.StudentID
	} else {
		return nil, errors.New("无权限评价")
	}

	review := &models.Review{
		BookingID:  input.BookingID,
		ReviewerID: reviewerID,
		RevieweeID: revieweeID,
		PostingID:  booking.PostingID,
		Rating:     input.Rating,
		Content:    input.Content,
		IsPublic:   input.IsPublic,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, errors.New("创建评价失败")
	}

	if booking.StudentID == reviewerID {
		booking.ReviewedByStudent = true
	} else {
		booking.ReviewedByTeacher = true
	}
	s.bookingRepo.Update(booking)

	s.userRepo.UpdateRating(revieweeID)
	s.skillRepo.UpdatePostingRating(booking.PostingID)

	return review, nil
}

func (s *BookingService) GetReviewsByPosting(postingID uuid.UUID, page, pageSize int) ([]*models.Review, int64, error) {
	return s.bookingRepo.GetReviewsByPosting(postingID, page, pageSize)
}

func (s *BookingService) GetReviewsByUser(userID uuid.UUID, page, pageSize int) ([]*models.Review, int64, error) {
	return s.bookingRepo.GetReviewsByUser(userID, page, pageSize)
}

type CreateComplaintInput struct {
	Type        models.ComplaintType `json:"type"`
	TargetID    uuid.UUID            `json:"target_id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Evidence    string               `json:"evidence"`
}

func (s *BookingService) CreateComplaint(reporterID uuid.UUID, input *CreateComplaintInput) (*models.Complaint, error) {
	complaint := &models.Complaint{
		ReporterID:  reporterID,
		Type:        input.Type,
		TargetID:    input.TargetID,
		Title:       input.Title,
		Description: input.Description,
		Evidence:    input.Evidence,
		Status:      models.ComplaintStatusPending,
	}

	if err := s.bookingRepo.CreateComplaint(complaint); err != nil {
		return nil, errors.New("创建投诉失败")
	}
	return complaint, nil
}

func (s *BookingService) GetComplaints(page, pageSize int, status *string) ([]*models.Complaint, int64, error) {
	return s.bookingRepo.GetComplaints(page, pageSize, status)
}

func (s *BookingService) HandleComplaint(id uuid.UUID, handlerID uuid.UUID, result string) error {
	return s.bookingRepo.HandleComplaint(id, handlerID, result)
}
