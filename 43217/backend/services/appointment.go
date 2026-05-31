package services

import (
	"errors"
	"fmt"
	"time"

	"health-platform/config"
	"health-platform/models"
	"health-platform/repository"
	"health-platform/utils"
)

type AppointmentService struct {
	appointmentRepo *repository.AppointmentRepository
	employeeRepo    *repository.EmployeeRepository
	packageRepo     *repository.PackageRepository
	timeSlotRepo    *repository.TimeSlotRepository
	companyRepo     *repository.CompanyRepository
	deptApptRepo    *repository.DepartmentAppointmentRepository
}

func NewAppointmentService() *AppointmentService {
	return &AppointmentService{
		appointmentRepo: repository.NewAppointmentRepository(),
		employeeRepo:    repository.NewEmployeeRepository(),
		packageRepo:     repository.NewPackageRepository(),
		timeSlotRepo:    repository.NewTimeSlotRepository(),
		companyRepo:     repository.NewCompanyRepository(),
		deptApptRepo:    repository.NewDepartmentAppointmentRepository(),
	}
}

type CreateAppointmentRequest struct {
	EmployeeID uint   `json:"employee_id" binding:"required"`
	AgencyID   uint   `json:"agency_id" binding:"required"`
	PackageID  uint   `json:"package_id" binding:"required"`
	TimeSlotID uint   `json:"time_slot_id" binding:"required"`
}

type RescheduleAppointmentRequest struct {
	TimeSlotID      uint   `json:"time_slot_id" binding:"required"`
	AppointmentDate string `json:"appointment_date" binding:"required"`
	StartTime       string `json:"start_time" binding:"required"`
	EndTime         string `json:"end_time" binding:"required"`
}

type CancelAppointmentRequest struct {
	Reason string `json:"reason" binding:"required"`
}

func (s *AppointmentService) CreateAppointment(req *CreateAppointmentRequest) (*models.Appointment, error) {
	var employee models.Employee
	if err := s.employeeRepo.FindByID(&employee, req.EmployeeID); err != nil {
		return nil, errors.New("员工不存在")
	}

	var pkg models.Package
	if err := s.packageRepo.FindByID(&pkg, req.PackageID); err != nil {
		return nil, errors.New("套餐不存在")
	}

	if pkg.Status != models.PackageStatusOnline {
		return nil, errors.New("该套餐已下架")
	}

	age := utils.AgeFromBirthday(employee.Birthday)
	if pkg.MinAge > 0 && age < pkg.MinAge {
		return nil, fmt.Errorf("该套餐适合年龄%d岁以上", pkg.MinAge)
	}
	if pkg.MaxAge > 0 && age > pkg.MaxAge {
		return nil, fmt.Errorf("该套餐适合年龄%d岁以下", pkg.MaxAge)
	}

	if pkg.GenderLimit > 0 && pkg.GenderLimit != employee.Gender {
		return nil, errors.New("该套餐不适合当前性别")
	}

	yearCount, _ := s.appointmentRepo.GetEmployeeAppointmentCount(req.EmployeeID, time.Now().Year())
	if employee.UsedQuota >= employee.Quota {
		return nil, errors.New("您的体检额度已用完")
	}

	deptAppointments, _ := s.deptApptRepo.FindByDepartmentID(employee.DepartmentID, time.Now().Year())
	hasValidDeptAppt := false
	for _, da := range deptAppointments {
		if da.AgencyID == req.AgencyID && da.PackageID == req.PackageID {
			if da.UsedQuota < da.TotalQuota {
				hasValidDeptAppt = true
			} else {
				return nil, errors.New("该部门此套餐的预约名额已满")
			}
		}
	}
	if len(deptAppointments) > 0 && !hasValidDeptAppt {
		return nil, errors.New("您的部门未分配此机构/套餐的预约名额")
	}

	var timeSlot models.TimeSlot
	if err := s.timeSlotRepo.FindByID(&timeSlot, req.TimeSlotID); err != nil {
		return nil, errors.New("预约时段不存在")
	}

	if timeSlot.Booked >= timeSlot.Total {
		return nil, errors.New("该时段已约满")
	}

	appointmentNo := utils.GenerateOrderNo("APT")

	appointment := &models.Appointment{
		EmployeeID:      req.EmployeeID,
		CompanyID:       employee.CompanyID,
		AgencyID:        req.AgencyID,
		PackageID:       req.PackageID,
		TimeSlotID:      req.TimeSlotID,
		AppointmentNo:   appointmentNo,
		AppointmentDate: timeSlot.Date,
		StartTime:       timeSlot.StartTime,
		EndTime:         timeSlot.EndTime,
		Status:          models.AppointmentStatusConfirmed,
	}

	if err := s.appointmentRepo.Create(appointment); err != nil {
		return nil, fmt.Errorf("创建预约失败: %w", err)
	}

	s.timeSlotRepo.IncrementBooked(req.TimeSlotID)
	s.employeeRepo.UpdateQuota(req.EmployeeID, 1)
	s.packageRepo.IncrementSaleCount(req.PackageID)
	s.companyRepo.UpdateBudget(employee.CompanyID, pkg.Price)

	for _, da := range deptAppointments {
		if da.AgencyID == req.AgencyID && da.PackageID == req.PackageID {
			s.deptApptRepo.UpdateUsedQuota(da.ID, 1)
			break
		}
	}

	cacheKey := fmt.Sprintf("hot_packages_%d", req.PackageID)
	config.RedisClient.Del(config.Ctx, cacheKey)

	return appointment, nil
}

func (s *AppointmentService) GetAppointment(appointmentID uint) (*models.Appointment, error) {
	return s.appointmentRepo.GetWithReport(appointmentID)
}

func (s *AppointmentService) GetEmployeeAppointments(employeeID uint, page, pageSize int) ([]models.Appointment, int64, error) {
	return s.appointmentRepo.FindByEmployeeID(employeeID, page, pageSize)
}

func (s *AppointmentService) GetCompanyAppointments(companyID uint, page, pageSize int) ([]models.Appointment, int64, error) {
	return s.appointmentRepo.FindByCompanyID(companyID, page, pageSize)
}

func (s *AppointmentService) GetAgencyAppointments(agencyID uint, page, pageSize int) ([]models.Appointment, int64, error) {
	return s.appointmentRepo.FindByAgencyID(agencyID, page, pageSize)
}

func (s *AppointmentService) RescheduleAppointment(appointmentID uint, req *RescheduleAppointmentRequest) error {
	var appointment models.Appointment
	if err := s.appointmentRepo.FindByID(&appointment, appointmentID); err != nil {
		return errors.New("预约不存在")
	}

	if appointment.Status != models.AppointmentStatusConfirmed {
		return errors.New("当前预约状态不支持改约")
	}

	var newTimeSlot models.TimeSlot
	if err := s.timeSlotRepo.FindByID(&newTimeSlot, req.TimeSlotID); err != nil {
		return errors.New("新预约时段不存在")
	}

	if newTimeSlot.Booked >= newTimeSlot.Total {
		return errors.New("新时段已约满")
	}

	appointmentDate, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		return errors.New("日期格式错误")
	}

	s.timeSlotRepo.DecrementBooked(appointment.TimeSlotID)
	s.timeSlotRepo.IncrementBooked(req.TimeSlotID)

	return s.appointmentRepo.RescheduleAppointment(appointmentID, req.TimeSlotID, appointmentDate, req.StartTime, req.EndTime)
}

func (s *AppointmentService) CancelAppointment(appointmentID uint, req *CancelAppointmentRequest) error {
	var appointment models.Appointment
	if err := s.appointmentRepo.FindByID(&appointment, appointmentID); err != nil {
		return errors.New("预约不存在")
	}

	if appointment.Status != models.AppointmentStatusConfirmed {
		return errors.New("当前预约状态不支持取消")
	}

	if appointment.AppointmentDate.Before(time.Now()) {
		return errors.New("已过预约时间，无法取消")
	}

	s.timeSlotRepo.DecrementBooked(appointment.TimeSlotID)

	return s.appointmentRepo.CancelAppointment(appointmentID, req.Reason)
}

func (s *AppointmentService) CompleteAppointment(appointmentID uint) error {
	return s.appointmentRepo.UpdateStatus(appointmentID, models.AppointmentStatusCompleted)
}

func (s *AppointmentService) GetEmployeeAppointmentStatus(employeeID uint) (map[string]interface{}, error) {
	var employee models.Employee
	if err := s.employeeRepo.FindByID(&employee, employeeID); err != nil {
		return nil, errors.New("员工不存在")
	}

	year := time.Now().Year()
	totalCount, _ := s.appointmentRepo.GetEmployeeAppointmentCount(employeeID, year)

	result := map[string]interface{}{
		"employee_id":   employeeID,
		"quota":         employee.Quota,
		"used_quota":    employee.UsedQuota,
		"remaining":     employee.Quota - employee.UsedQuota,
		"total_appointments": totalCount,
		"year":          year,
	}

	return result, nil
}

func (s *AppointmentService) CheckQuota(employeeID uint, packageID uint) (bool, error) {
	var employee models.Employee
	if err := s.employeeRepo.FindByID(&employee, employeeID); err != nil {
		return false, errors.New("员工不存在")
	}

	if employee.UsedQuota >= employee.Quota {
		return false, nil
	}

	var pkg models.Package
	if err := s.packageRepo.FindByID(&pkg, packageID); err != nil {
		return false, errors.New("套餐不存在")
	}

	if pkg.Status != models.PackageStatusOnline {
		return false, nil
	}

	return true, nil
}
