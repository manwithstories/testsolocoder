package service

import (
	"gym-management/internal/models"
	"gym-management/internal/pkg/database"
	"time"

	"gorm.io/gorm"
)

type StatsService interface {
	GetMemberStats(startDate, endDate time.Time) (map[string]interface{}, error)
	GetCourseStats(startDate, endDate time.Time) (map[string]interface{}, error)
	GetCoachStats(startDate, endDate time.Time) (map[string]interface{}, error)
	GetDashboardStats() (map[string]interface{}, error)
}

type statsService struct {
	db *gorm.DB
}

func NewStatsService() StatsService {
	return &statsService{
		db: database.GetDB(),
	}
}

func (s *statsService) GetMemberStats(startDate, endDate time.Time) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var totalMembers int64
	s.db.Model(&models.Member{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Count(&totalMembers)
	result["new_members"] = totalMembers

	var activeMembers int64
	s.db.Model(&models.CheckIn{}).
		Where("check_in_time BETWEEN ? AND ?", startDate, endDate).
		Distinct("member_id").
		Count(&activeMembers)
	result["active_members"] = activeMembers

	type DailyNewMember struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var dailyNew []DailyNewMember
	s.db.Model(&models.Member{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("DATE(created_at)").
		Scan(&dailyNew)
	result["daily_new"] = dailyNew

	return result, nil
}

func (s *statsService) GetCourseStats(startDate, endDate time.Time) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var totalBookings int64
	s.db.Model(&models.Booking{}).
		Where("status IN (1, 3) AND created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&totalBookings)
	result["total_bookings"] = totalBookings

	var totalCheckIns int64
	s.db.Model(&models.CheckIn{}).
		Where("check_in_time BETWEEN ? AND ?", startDate, endDate).
		Count(&totalCheckIns)
	result["total_check_ins"] = totalCheckIns

	type CourseAttendance struct {
		CourseID   uint   `json:"course_id"`
		CourseName string `json:"course_name"`
		Booked     int64  `json:"booked"`
		CheckedIn  int64  `json:"checked_in"`
		Capacity   int    `json:"capacity"`
	}
	var courseStats []CourseAttendance
	s.db.Table("bookings").
		Select("courses.id as course_id, courses.name as course_name, "+
			"COUNT(DISTINCT bookings.id) as booked, "+
			"COUNT(DISTINCT check_ins.id) as checked_in, "+
			"course_schedules.capacity as capacity").
		Joins("JOIN course_schedules ON course_schedules.id = bookings.schedule_id").
		Joins("JOIN courses ON courses.id = course_schedules.course_id").
		Joins("LEFT JOIN check_ins ON check_ins.schedule_id = course_schedules.id AND check_ins.member_id = bookings.member_id").
		Where("bookings.created_at BETWEEN ? AND ?", startDate, endDate).
		Group("courses.id, courses.name, course_schedules.capacity").
		Scan(&courseStats)
	result["course_attendance"] = courseStats

	type DailyBooking struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var dailyBookings []DailyBooking
	s.db.Model(&models.Booking{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Group("DATE(created_at)").
		Scan(&dailyBookings)
	result["daily_bookings"] = dailyBookings

	return result, nil
}

func (s *statsService) GetCoachStats(startDate, endDate time.Time) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	type CoachStat struct {
		CoachID     uint   `json:"coach_id"`
		CoachName   string `json:"coach_name"`
		TotalCourses int64 `json:"total_courses"`
		TotalStudents int64 `json:"total_students"`
		TotalHours  float64 `json:"total_hours"`
	}
	var coachStats []CoachStat

	s.db.Table("courses").
		Select("coaches.id as coach_id, coaches.name as coach_name, "+
			"COUNT(DISTINCT course_schedules.id) as total_courses, "+
			"COUNT(DISTINCT bookings.member_id) as total_students, "+
			"SUM(courses.duration) / 60.0 as total_hours").
		Joins("JOIN coaches ON coaches.id = courses.coach_id").
		Joins("JOIN course_schedules ON course_schedules.course_id = courses.id").
		Joins("LEFT JOIN bookings ON bookings.schedule_id = course_schedules.id AND bookings.status IN (1, 3)").
		Where("course_schedules.start_time BETWEEN ? AND ?", startDate, endDate).
		Group("coaches.id, coaches.name").
		Scan(&coachStats)

	result["coach_stats"] = coachStats
	return result, nil
}

func (s *statsService) GetDashboardStats() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var totalMembers int64
	s.db.Model(&models.Member{}).Where("status = 1").Count(&totalMembers)
	result["total_members"] = totalMembers

	var activeMemberships int64
	s.db.Model(&models.Membership{}).Where("status = 1").Count(&activeMemberships)
	result["active_memberships"] = activeMemberships

	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var todayCheckIns int64
	s.db.Model(&models.CheckIn{}).Where("check_in_time BETWEEN ? AND ?", startOfDay, endOfDay).Count(&todayCheckIns)
	result["today_check_ins"] = todayCheckIns

	var todayBookings int64
	s.db.Model(&models.Booking{}).
		Joins("JOIN course_schedules ON course_schedules.id = bookings.schedule_id").
		Where("bookings.status IN (1, 3) AND course_schedules.start_time BETWEEN ? AND ?", startOfDay, endOfDay).
		Count(&todayBookings)
	result["today_bookings"] = todayBookings

	var upcomingCourses int64
	s.db.Model(&models.CourseSchedule{}).
		Where("status = 1 AND start_time > ? AND start_time <= ?", today, today.Add(24*time.Hour)).
		Count(&upcomingCourses)
	result["upcoming_courses"] = upcomingCourses

	return result, nil
}
