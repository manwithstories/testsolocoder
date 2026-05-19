package repository

import (
	"gym-management/internal/models"
	"gym-management/internal/pkg/database"
	"time"

	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *models.Course) error
	GetByID(id uint) (*models.Course, error)
	List(page, pageSize int, keyword string, coachID uint, courseType string) ([]models.Course, int64, error)
	Update(course *models.Course) error
	Delete(id uint) error
	UpdateStatus(id uint, status int) error
	DB() *gorm.DB
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository() CourseRepository {
	return &courseRepository{db: database.GetDB()}
}

func (r *courseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) GetByID(id uint) (*models.Course, error) {
	var course models.Course
	err := r.db.Preload("Coach").Preload("CourseSchedules").First(&course, id).Error
	return &course, err
}

func (r *courseRepository) List(page, pageSize int, keyword string, coachID uint, courseType string) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	query := r.db.Model(&models.Course{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if coachID > 0 {
		query = query.Where("coach_id = ?", coachID)
	}
	if courseType != "" {
		query = query.Where("type = ?", courseType)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Preload("Coach").Find(&courses).Error
	return courses, total, err
}

func (r *courseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

func (r *courseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Course{}, id).Error
}

func (r *courseRepository) UpdateStatus(id uint, status int) error {
	return r.db.Model(&models.Course{}).Where("id = ?", id).Update("status", status).Error
}

func (r *courseRepository) DB() *gorm.DB {
	return r.db
}

type ScheduleRepository interface {
	Create(schedule *models.CourseSchedule) error
	GetByID(id uint) (*models.CourseSchedule, error)
	GetByCourseID(courseID uint) ([]models.CourseSchedule, error)
	List(page, pageSize int, courseID uint, startDate, endDate *time.Time) ([]models.CourseSchedule, int64, error)
	Update(schedule *models.CourseSchedule) error
	Delete(id uint) error
	UpdateBookedCount(id uint, count int) error
	FindAvailable() ([]models.CourseSchedule, error)
	GetUpcoming(hours int) ([]models.CourseSchedule, error)
}

type scheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository() ScheduleRepository {
	return &scheduleRepository{db: database.GetDB()}
}

func (r *scheduleRepository) Create(schedule *models.CourseSchedule) error {
	return r.db.Create(schedule).Error
}

func (r *scheduleRepository) GetByID(id uint) (*models.CourseSchedule, error) {
	var schedule models.CourseSchedule
	err := r.db.Preload("Course.Coach").First(&schedule, id).Error
	return &schedule, err
}

func (r *scheduleRepository) GetByCourseID(courseID uint) ([]models.CourseSchedule, error) {
	var schedules []models.CourseSchedule
	err := r.db.Where("course_id = ?", courseID).Find(&schedules).Error
	return schedules, err
}

func (r *scheduleRepository) List(page, pageSize int, courseID uint, startDate, endDate *time.Time) ([]models.CourseSchedule, int64, error) {
	var schedules []models.CourseSchedule
	var total int64

	query := r.db.Model(&models.CourseSchedule{})
	if courseID > 0 {
		query = query.Where("course_id = ?", courseID)
	}
	if startDate != nil {
		query = query.Where("start_time >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("start_time <= ?", *endDate)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Preload("Course.Coach").Order("start_time ASC").Find(&schedules).Error
	return schedules, total, err
}

func (r *scheduleRepository) Update(schedule *models.CourseSchedule) error {
	return r.db.Save(schedule).Error
}

func (r *scheduleRepository) Delete(id uint) error {
	return r.db.Delete(&models.CourseSchedule{}, id).Error
}

func (r *scheduleRepository) UpdateBookedCount(id uint, count int) error {
	return r.db.Model(&models.CourseSchedule{}).Where("id = ?", id).Update("booked_count", count).Error
}

func (r *scheduleRepository) FindAvailable() ([]models.CourseSchedule, error) {
	var schedules []models.CourseSchedule
	now := time.Now()
	err := r.db.Where("status = 1 AND booked_count < capacity AND start_time > ?", now).
		Preload("Course.Coach").
		Order("start_time ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *scheduleRepository) GetUpcoming(hours int) ([]models.CourseSchedule, error) {
	var schedules []models.CourseSchedule
	now := time.Now()
	endTime := now.Add(time.Duration(hours) * time.Hour)
	err := r.db.Where("status = 1 AND start_time > ? AND start_time <= ?", now, endTime).
		Preload("Course.Coach").
		Order("start_time ASC").
		Find(&schedules).Error
	return schedules, err
}
