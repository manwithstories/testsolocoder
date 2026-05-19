package service

import (
	"errors"
	"gym-management/internal/models"
	"gym-management/internal/repository"
	"time"
)

type CheckInService interface {
	CheckIn(memberID uint, scheduleID *uint) (*models.CheckIn, error)
	GetByID(id uint) (*models.CheckIn, error)
	ListByMember(memberID uint, page, pageSize int) ([]models.CheckIn, int64, error)
	ListByDate(date time.Time) ([]models.CheckIn, error)
	ValidateAndCheckIn(memberID uint, scheduleID *uint) (*models.CheckIn, error)
}

type checkInService struct {
	checkInRepo  repository.CheckInRepository
	bookingRepo  repository.BookingRepository
	scheduleRepo repository.ScheduleRepository
}

func NewCheckInService() CheckInService {
	return &checkInService{
		checkInRepo:  repository.NewCheckInRepository(),
		bookingRepo:  repository.NewBookingRepository(),
		scheduleRepo: repository.NewScheduleRepository(),
	}
}

func (s *checkInService) CheckIn(memberID uint, scheduleID *uint) (*models.CheckIn, error) {
	checkIn := &models.CheckIn{
		MemberID:    memberID,
		ScheduleID:  scheduleID,
		CheckInTime: time.Now(),
		CheckType:   1,
	}

	err := s.checkInRepo.Create(checkIn)
	if err != nil {
		return nil, err
	}

	if scheduleID != nil {
		booking, _ := s.bookingRepo.GetByMemberAndSchedule(memberID, *scheduleID)
		if booking != nil && booking.Status == models.BookingStatusBooked {
			booking.Status = models.BookingStatusCheckedIn
			booking.CheckInID = &checkIn.ID
			_ = s.bookingRepo.Update(booking)
		}
	}

	return checkIn, nil
}

func (s *checkInService) ValidateAndCheckIn(memberID uint, scheduleID *uint) (*models.CheckIn, error) {
	now := time.Now()

	if scheduleID != nil {
		booking, err := s.bookingRepo.GetByMemberAndSchedule(memberID, *scheduleID)
		if err != nil || booking.ID == 0 {
			return nil, errors.New("您没有预约该课程")
		}

		if booking.Status == models.BookingStatusCancelled {
			return nil, errors.New("该预约已取消")
		}

		schedule, err := s.scheduleRepo.GetByID(*scheduleID)
		if err != nil {
			return nil, errors.New("课程不存在")
		}

		scheduleDate := time.Date(schedule.StartTime.Year(), schedule.StartTime.Month(), schedule.StartTime.Day(), 0, 0, 0, 0, time.Local)
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

		if !scheduleDate.Equal(today) {
			return nil, errors.New("只能签到当天的课程")
		}

		existingCheckIns, _ := s.checkInRepo.GetByMemberAndDate(memberID, now)
		for _, ci := range existingCheckIns {
			if ci.ScheduleID != nil && *ci.ScheduleID == *scheduleID {
				return nil, errors.New("您已完成该课程签到")
			}
		}
	} else {
		existingCheckIns, _ := s.checkInRepo.GetByMemberAndDate(memberID, now)
		if len(existingCheckIns) > 0 {
			return nil, errors.New("您今日已签到")
		}
	}

	return s.CheckIn(memberID, scheduleID)
}

func (s *checkInService) GetByID(id uint) (*models.CheckIn, error) {
	return s.checkInRepo.GetByID(id)
}

func (s *checkInService) ListByMember(memberID uint, page, pageSize int) ([]models.CheckIn, int64, error) {
	return s.checkInRepo.ListByMember(memberID, page, pageSize)
}

func (s *checkInService) ListByDate(date time.Time) ([]models.CheckIn, error) {
	return s.checkInRepo.ListByDate(date)
}
