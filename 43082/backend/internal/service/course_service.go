package service

import (
	"errors"
	"gym-management/internal/models"
	"gym-management/internal/repository"
	"strconv"
	"strings"
	"time"
)

type CourseService interface {
	Create(course *models.Course) error
	GetByID(id uint) (*models.Course, error)
	List(page, pageSize int, keyword string, coachID uint, courseType string) ([]models.Course, int64, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	UpdateStatus(id uint, status int) error
	GenerateSchedules(courseID uint) error
}

type courseService struct {
	courseRepo   repository.CourseRepository
	scheduleRepo repository.ScheduleRepository
}

func NewCourseService() CourseService {
	return &courseService{
		courseRepo:   repository.NewCourseRepository(),
		scheduleRepo: repository.NewScheduleRepository(),
	}
}

func (s *courseService) Create(course *models.Course) error {
	if course.Capacity <= 0 {
		return errors.New("课程容量必须大于0")
	}
	if course.Duration <= 0 {
		return errors.New("课程时长必须大于0")
	}
	if course.StartTime == "" {
		return errors.New("请设置课程开始时间")
	}

	course.Status = 1
	err := s.courseRepo.Create(course)
	if err != nil {
		return err
	}

	if course.Type != models.CourseTypeSingle {
		return s.GenerateSchedules(course.ID)
	}

	return nil
}

func (s *courseService) GetByID(id uint) (*models.Course, error) {
	return s.courseRepo.GetByID(id)
}

func (s *courseService) List(page, pageSize int, keyword string, coachID uint, courseType string) ([]models.Course, int64, error) {
	return s.courseRepo.List(page, pageSize, keyword, coachID, courseType)
}

func (s *courseService) Update(id uint, updates map[string]interface{}) error {
	_, err := s.courseRepo.GetByID(id)
	if err != nil {
		return errors.New("课程不存在")
	}
	return s.courseRepo.DB().Model(&models.Course{}).Where("id = ?", id).Updates(updates).Error
}

func (s *courseService) Delete(id uint) error {
	_, err := s.courseRepo.GetByID(id)
	if err != nil {
		return errors.New("课程不存在")
	}
	return s.courseRepo.Delete(id)
}

func (s *courseService) UpdateStatus(id uint, status int) error {
	_, err := s.courseRepo.GetByID(id)
	if err != nil {
		return errors.New("课程不存在")
	}
	return s.courseRepo.UpdateStatus(id, status)
}

func (s *courseService) GenerateSchedules(courseID uint) error {
	course, err := s.courseRepo.GetByID(courseID)
	if err != nil {
		return err
	}

	if course.Type == models.CourseTypeSingle {
		startDateTime, err := parseDateTime(course.StartDate, course.StartTime)
		if err != nil {
			return err
		}
		endDateTime := startDateTime.Add(time.Duration(course.Duration) * time.Minute)

		schedule := &models.CourseSchedule{
			CourseID:   course.ID,
			StartTime:  startDateTime,
			EndTime:    endDateTime,
			Capacity:   course.Capacity,
			BookedCount: 0,
			Status:     1,
		}
		return s.scheduleRepo.Create(schedule)
	}

	weekdays := parseWeekdays(course.Weekdays)
	if len(weekdays) == 0 {
		return errors.New("请设置周几上课")
	}

	endDate := course.EndDate
	if endDate == nil {
		oneYearLater := course.StartDate.AddDate(1, 0, 0)
		endDate = &oneYearLater
	}

	currentDate := course.StartDate
	for currentDate.Before(*endDate) || currentDate.Equal(*endDate) {
		weekday := int(currentDate.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		if containsInt(weekdays, weekday) {
			startDateTime, err := parseDateTime(currentDate, course.StartTime)
			if err == nil {
				endDateTime := startDateTime.Add(time.Duration(course.Duration) * time.Minute)

				schedule := &models.CourseSchedule{
					CourseID:    course.ID,
					StartTime:   startDateTime,
					EndTime:     endDateTime,
					Capacity:    course.Capacity,
					BookedCount: 0,
					Status:      1,
				}
				_ = s.scheduleRepo.Create(schedule)
			}
		}

		if course.Type == models.CourseTypeWeekly {
			currentDate = currentDate.AddDate(0, 0, 1)
		} else if course.Type == models.CourseTypeMonthly {
			currentDate = currentDate.AddDate(0, 1, 0)
		}
	}

	return nil
}

func parseDateTime(date time.Time, timeStr string) (time.Time, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return time.Time{}, errors.New("时间格式错误")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil || hour < 0 || hour > 23 {
		return time.Time{}, errors.New("小时无效")
	}

	minute, err := strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute > 59 {
		return time.Time{}, errors.New("分钟无效")
	}

	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, time.Local), nil
}

func parseWeekdays(weekdaysStr string) []int {
	var weekdays []int
	parts := strings.Split(weekdaysStr, ",")
	for _, part := range parts {
		day, err := strconv.Atoi(strings.TrimSpace(part))
		if err == nil && day >= 1 && day <= 7 {
			weekdays = append(weekdays, day)
		}
	}
	return weekdays
}

func containsInt(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type ScheduleService interface {
	GetByID(id uint) (*models.CourseSchedule, error)
	List(page, pageSize int, courseID uint, startDate, endDate *time.Time) ([]models.CourseSchedule, int64, error)
	FindAvailable() ([]models.CourseSchedule, error)
	UpdateStatus(id uint, status int) error
}

type scheduleService struct {
	scheduleRepo repository.ScheduleRepository
}

func NewScheduleService() ScheduleService {
	return &scheduleService{
		scheduleRepo: repository.NewScheduleRepository(),
	}
}

func (s *scheduleService) GetByID(id uint) (*models.CourseSchedule, error) {
	return s.scheduleRepo.GetByID(id)
}

func (s *scheduleService) List(page, pageSize int, courseID uint, startDate, endDate *time.Time) ([]models.CourseSchedule, int64, error) {
	return s.scheduleRepo.List(page, pageSize, courseID, startDate, endDate)
}

func (s *scheduleService) FindAvailable() ([]models.CourseSchedule, error) {
	return s.scheduleRepo.FindAvailable()
}

func (s *scheduleService) UpdateStatus(id uint, status int) error {
	schedule, err := s.scheduleRepo.GetByID(id)
	if err != nil {
		return errors.New("课程排期不存在")
	}
	schedule.Status = status
	return s.scheduleRepo.Update(schedule)
}
