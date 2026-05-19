package service

import (
	"errors"
	"gym-management/internal/models"
	"gym-management/internal/repository"
	"time"

	"gorm.io/gorm"
)

type BookingService interface {
	Book(memberID, scheduleID uint) (*models.Booking, error)
	Cancel(bookingID uint, memberID uint) error
	GetByID(id uint) (*models.Booking, error)
	ListByMember(memberID uint, page, pageSize int) ([]models.Booking, int64, error)
	ListBySchedule(scheduleID uint, page, pageSize int) ([]models.Booking, int64, error)
	AddToWaitlist(memberID, scheduleID uint) (*models.Waitlist, error)
	RemoveFromWaitlist(waitlistID uint, memberID uint) error
	ProcessWaitlist(scheduleID uint) error
}

type bookingService struct {
	bookingRepo    repository.BookingRepository
	waitlistRepo   repository.WaitlistRepository
	scheduleRepo   repository.ScheduleRepository
	membershipRepo repository.MembershipRepository
	db             *gorm.DB
}

func NewBookingService() BookingService {
	return &bookingService{
		bookingRepo:    repository.NewBookingRepository(),
		waitlistRepo:   repository.NewWaitlistRepository(),
		scheduleRepo:   repository.NewScheduleRepository(),
		membershipRepo: repository.NewMembershipRepository(),
		db:             repository.GetDB(),
	}
}

func (s *bookingService) Book(memberID, scheduleID uint) (*models.Booking, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	membership, err := s.membershipRepo.GetByMemberID(memberID)
	if err != nil || !membership.IsValid() {
		tx.Rollback()
		return nil, errors.New("会员卡无效或已过期")
	}

	schedule, err := s.scheduleRepo.GetByID(scheduleID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("课程排期不存在")
	}

	if schedule.Status != 1 {
		tx.Rollback()
		return nil, errors.New("该课程不可预约")
	}

	if time.Now().After(schedule.StartTime) {
		tx.Rollback()
		return nil, errors.New("课程已开始，无法预约")
	}

	existing, _ := s.bookingRepo.GetByMemberAndSchedule(memberID, scheduleID)
	if existing != nil && existing.ID > 0 && existing.Status != 2 {
		tx.Rollback()
		return nil, errors.New("您已预约该课程")
	}

	hasConflict, err := s.bookingRepo.CheckMemberTimeConflict(memberID, schedule.StartTime, schedule.EndTime)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if hasConflict {
		tx.Rollback()
		return nil, errors.New("该时间段您已有其他预约")
	}

	var bookedCount int64
	err = tx.Model(&models.Booking{}).
		Where("schedule_id = ? AND status IN (1, 3)", scheduleID).
		Count(&bookedCount).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if int(bookedCount) >= schedule.Capacity {
		tx.Rollback()
		return nil, errors.New("课程已满，可加入等待列表")
	}

	booking := &models.Booking{
		MemberID:    memberID,
		ScheduleID:  scheduleID,
		Status:      models.BookingStatusBooked,
		BookingTime: time.Now(),
	}

	if err := tx.Create(booking).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(&models.CourseSchedule{}).
		Where("id = ?", scheduleID).
		Update("booked_count", bookedCount+1).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return booking, nil
}

func (s *bookingService) Cancel(bookingID uint, memberID uint) error {
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return errors.New("预约不存在")
	}

	if booking.MemberID != memberID {
		return errors.New("无权取消他人预约")
	}

	if booking.Status == models.BookingStatusCancelled {
		return errors.New("该预约已取消")
	}

	schedule, err := s.scheduleRepo.GetByID(booking.ScheduleID)
	if err != nil {
		return errors.New("课程排期不存在")
	}

	hoursBefore := time.Until(schedule.StartTime).Hours()
	if hoursBefore < 2 {
		return errors.New("开课前2小时内不能取消预约")
	}

	now := time.Now()
	booking.Status = models.BookingStatusCancelled
	booking.CancelTime = &now

	err = s.bookingRepo.Update(booking)
	if err != nil {
		return err
	}

	if schedule.BookedCount > 0 {
		schedule.BookedCount--
		_ = s.scheduleRepo.Update(schedule)
	}

	go s.ProcessWaitlist(booking.ScheduleID)

	return nil
}

func (s *bookingService) GetByID(id uint) (*models.Booking, error) {
	return s.bookingRepo.GetByID(id)
}

func (s *bookingService) ListByMember(memberID uint, page, pageSize int) ([]models.Booking, int64, error) {
	return s.bookingRepo.ListByMember(memberID, page, pageSize)
}

func (s *bookingService) ListBySchedule(scheduleID uint, page, pageSize int) ([]models.Booking, int64, error) {
	return s.bookingRepo.ListBySchedule(scheduleID, page, pageSize)
}

func (s *bookingService) AddToWaitlist(memberID, scheduleID uint) (*models.Waitlist, error) {
	membership, err := s.membershipRepo.GetByMemberID(memberID)
	if err != nil || !membership.IsValid() {
		return nil, errors.New("会员卡无效或已过期")
	}

	schedule, err := s.scheduleRepo.GetByID(scheduleID)
	if err != nil {
		return nil, errors.New("课程排期不存在")
	}

	if time.Now().After(schedule.StartTime) {
		return nil, errors.New("课程已开始，无法加入等待列表")
	}

	existing, _ := s.bookingRepo.GetByMemberAndSchedule(memberID, scheduleID)
	if existing != nil && existing.ID > 0 && existing.Status != 2 {
		return nil, errors.New("您已预约该课程")
	}

	existingWaitlist, _ := s.waitlistRepo.GetByMemberAndSchedule(memberID, scheduleID)
	if existingWaitlist != nil && existingWaitlist.ID > 0 {
		return nil, errors.New("您已在等待列表中")
	}

	maxPos, _ := s.waitlistRepo.GetMaxPosition(scheduleID)

	waitlist := &models.Waitlist{
		MemberID:   memberID,
		ScheduleID: scheduleID,
		Position:   maxPos + 1,
		Status:     1,
		Notified:   false,
	}

	err = s.waitlistRepo.Create(waitlist)
	if err != nil {
		return nil, err
	}

	return waitlist, nil
}

func (s *bookingService) RemoveFromWaitlist(waitlistID uint, memberID uint) error {
	waitlist, err := s.waitlistRepo.GetByID(waitlistID)
	if err != nil {
		return errors.New("等待列表记录不存在")
	}

	if waitlist.MemberID != memberID {
		return errors.New("无权操作")
	}

	waitlist.Status = 3
	return s.waitlistRepo.Update(waitlist)
}

func (s *bookingService) ProcessWaitlist(scheduleID uint) error {
	schedule, err := s.scheduleRepo.GetByID(scheduleID)
	if err != nil {
		return err
	}

	bookedCount, _ := s.bookingRepo.CountBySchedule(scheduleID)
	availableSlots := schedule.Capacity - int(bookedCount)

	if availableSlots <= 0 {
		return nil
	}

	waitlists, err := s.waitlistRepo.ListBySchedule(scheduleID)
	if err != nil {
		return err
	}

	for i := 0; i < availableSlots && i < len(waitlists); i++ {
		waitlist := waitlists[i]

		booking := &models.Booking{
			MemberID:    waitlist.MemberID,
			ScheduleID:  scheduleID,
			Status:      models.BookingStatusBooked,
			BookingTime: time.Now(),
		}

		err := s.bookingRepo.Create(booking)
		if err == nil {
			waitlist.Status = 2
			waitlist.Notified = true
			_ = s.waitlistRepo.Update(&waitlist)

			schedule.BookedCount++
			_ = s.scheduleRepo.Update(schedule)
		}
	}

	return nil
}
