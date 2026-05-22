package service

import (
	"encoding/json"
	"fmt"
	"time"

	"beauty-salon-system/internal/model"
	"beauty-salon-system/internal/repository"
	"beauty-salon-system/internal/repository/redis"
	"beauty-salon-system/internal/utils"
)

type ServiceProduct struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type AppointmentService struct {
	appointmentRepo *repository.AppointmentRepository
	technicianRepo  *repository.TechnicianRepository
	customerRepo    *repository.CustomerRepository
	serviceRepo     *repository.ServiceRepository
	packageRepo     *repository.CustomerPackageRepository
	productRepo     *repository.ProductRepository
	productRecordRepo *repository.ProductRecordRepository
	cancelFreeHours int
	deductPoints    int
}

func NewAppointmentService(
	appointmentRepo *repository.AppointmentRepository,
	technicianRepo *repository.TechnicianRepository,
	customerRepo *repository.CustomerRepository,
	serviceRepo *repository.ServiceRepository,
	packageRepo *repository.CustomerPackageRepository,
	productRepo *repository.ProductRepository,
	productRecordRepo *repository.ProductRecordRepository,
	cancelFreeHours int,
	deductPoints int,
) *AppointmentService {
	return &AppointmentService{
		appointmentRepo:   appointmentRepo,
		technicianRepo:    technicianRepo,
		customerRepo:      customerRepo,
		serviceRepo:       serviceRepo,
		packageRepo:       packageRepo,
		productRepo:       productRepo,
		productRecordRepo: productRecordRepo,
		cancelFreeHours:   cancelFreeHours,
		deductPoints:      deductPoints,
	}
}

type CreateAppointmentRequest struct {
	CustomerID      uint   `json:"customer_id" binding:"required"`
	TechnicianID    uint   `json:"technician_id" binding:"required"`
	ServiceID       uint   `json:"service_id" binding:"required"`
	PackageID       *uint  `json:"package_id"`
	AppointmentDate string `json:"appointment_date" binding:"required"`
	StartTime       string `json:"start_time" binding:"required"`
	Remark          string `json:"remark"`
}

type CancelAppointmentRequest struct {
	ID           uint   `json:"id" binding:"required"`
	CancelReason string `json:"cancel_reason"`
}

type RescheduleRequest struct {
	ID              uint   `json:"id" binding:"required"`
	AppointmentDate string `json:"appointment_date" binding:"required"`
	StartTime       string `json:"start_time" binding:"required"`
}

type AvailableSlot struct {
	Date  string   `json:"date"`
	Slots []string `json:"slots"`
}

func (s *AppointmentService) Create(req *CreateAppointmentRequest) (*model.Appointment, error) {
	technician, err := s.technicianRepo.GetByID(req.TechnicianID)
	if err != nil {
		return nil, fmt.Errorf("technician not found: %w", err)
	}
	if technician.Status != 1 {
		return nil, fmt.Errorf("technician is not available")
	}

	serviceItem, err := s.serviceRepo.GetByID(req.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}

	if err := s.CheckProductStock(req.ServiceID); err != nil {
		return nil, err
	}

	customer, err := s.customerRepo.GetByID(req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	appointmentDate, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}

	isLeave, err := s.technicianRepo.IsOnLeave(req.TechnicianID, appointmentDate)
	if err != nil {
		return nil, fmt.Errorf("check leave: %w", err)
	}
	if isLeave {
		return nil, fmt.Errorf("technician is on leave on this date")
	}

	endTime := utils.AddMinutesToTime(req.StartTime, serviceItem.Duration)

	lockKey := fmt.Sprintf("appointment:lock:%d:%s", req.TechnicianID, req.AppointmentDate)
	locked, err := redis.SetNX(lockKey, "1", 30*time.Second)
	if err != nil || !locked {
		return nil, fmt.Errorf("system busy, please try again")
	}
	defer redis.Delete(lockKey)

	conflict, err := s.appointmentRepo.CheckTimeConflict(req.TechnicianID, appointmentDate, req.StartTime, endTime, nil)
	if err != nil {
		return nil, fmt.Errorf("check conflict: %w", err)
	}
	if conflict {
		availableSlots, _ := s.GetAvailableSlots(req.TechnicianID, req.AppointmentDate, serviceItem.Duration)
		return nil, fmt.Errorf("time conflict, available slots: %v", availableSlots)
	}

	appointment := &model.Appointment{
		CustomerID:      req.CustomerID,
		TechnicianID:    req.TechnicianID,
		ServiceID:       req.ServiceID,
		PackageID:       req.PackageID,
		AppointmentDate: appointmentDate,
		StartTime:       req.StartTime,
		EndTime:         endTime,
		Status:          "confirmed",
		Remark:          req.Remark,
		Customer:        customer,
		Technician:      technician,
		Service:         serviceItem,
	}

	if req.PackageID != nil {
		customerPkg, err := s.packageRepo.GetByID(*req.PackageID)
		if err != nil {
			return nil, fmt.Errorf("package not found: %w", err)
		}
		if customerPkg.UsedCount >= customerPkg.TotalCount {
			return nil, fmt.Errorf("package has been fully used")
		}
		if customerPkg.ExpireDate.Before(time.Now()) {
			return nil, fmt.Errorf("package has expired")
		}
	}

	if err := s.appointmentRepo.Create(appointment); err != nil {
		return nil, fmt.Errorf("create appointment: %w", err)
	}

	redisKey := fmt.Sprintf("appointment:%d:%s", req.TechnicianID, req.AppointmentDate)
	redis.SAdd(redisKey, fmt.Sprintf("%s-%s", req.StartTime, endTime))

	return appointment, nil
}

func (s *AppointmentService) GetByID(id uint) (*model.Appointment, error) {
	return s.appointmentRepo.GetByID(id)
}

func (s *AppointmentService) List(page, pageSize int, filters map[string]interface{}) ([]model.Appointment, int64, error) {
	return s.appointmentRepo.List(page, pageSize, filters)
}

func (s *AppointmentService) GetByTechnicianAndDate(technicianID uint, dateStr string) ([]model.Appointment, error) {
	date, _ := time.Parse("2006-01-02", dateStr)
	return s.appointmentRepo.GetByTechnicianAndDate(technicianID, date)
}

func (s *AppointmentService) GetByCustomer(customerID uint, page, pageSize int) ([]model.Appointment, int64, error) {
	return s.appointmentRepo.GetByCustomer(customerID, page, pageSize)
}

func (s *AppointmentService) Cancel(req *CancelAppointmentRequest) error {
	appointment, err := s.appointmentRepo.GetByID(req.ID)
	if err != nil {
		return fmt.Errorf("appointment not found: %w", err)
	}

	if appointment.Status == "cancelled" || appointment.Status == "completed" {
		return fmt.Errorf("appointment cannot be cancelled")
	}

	now := time.Now()
	appointmentTime := time.Date(
		appointment.AppointmentDate.Year(),
		appointment.AppointmentDate.Month(),
		appointment.AppointmentDate.Day(),
		0, 0, 0, 0,
		appointment.AppointmentDate.Location(),
	)

	hoursDiff := appointmentTime.Sub(now).Hours()

	pointsDeducted := 0
	if hoursDiff < float64(s.cancelFreeHours) {
		pointsDeducted = s.deductPoints
		if err := s.customerRepo.DeductPoints(appointment.CustomerID, pointsDeducted); err != nil {
			return fmt.Errorf("deduct points: %w", err)
		}
	}

	appointment.Status = "cancelled"
	appointment.CancelReason = req.CancelReason
	appointment.PointsDeducted = pointsDeducted

	if err := s.appointmentRepo.Update(appointment); err != nil {
		return fmt.Errorf("update appointment: %w", err)
	}

	redisKey := fmt.Sprintf("appointment:%d:%s", appointment.TechnicianID, appointment.AppointmentDate.Format("2006-01-02"))
	redis.SRem(redisKey, fmt.Sprintf("%s-%s", appointment.StartTime, appointment.EndTime))

	return nil
}

func (s *AppointmentService) Reschedule(req *RescheduleRequest) (*model.Appointment, error) {
	appointment, err := s.appointmentRepo.GetByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("appointment not found: %w", err)
	}

	if appointment.Status == "cancelled" || appointment.Status == "completed" {
		return nil, fmt.Errorf("appointment cannot be rescheduled")
	}

	appointmentDate, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}

	isLeave, err := s.technicianRepo.IsOnLeave(appointment.TechnicianID, appointmentDate)
	if err != nil {
		return nil, fmt.Errorf("check leave: %w", err)
	}
	if isLeave {
		return nil, fmt.Errorf("technician is on leave on this date")
	}

	serviceItem, _ := s.serviceRepo.GetByID(appointment.ServiceID)
	endTime := utils.AddMinutesToTime(req.StartTime, serviceItem.Duration)

	conflict, err := s.appointmentRepo.CheckTimeConflict(appointment.TechnicianID, appointmentDate, req.StartTime, endTime, &req.ID)
	if err != nil {
		return nil, fmt.Errorf("check conflict: %w", err)
	}
	if conflict {
		availableSlots, _ := s.GetAvailableSlots(appointment.TechnicianID, req.AppointmentDate, serviceItem.Duration)
		return nil, fmt.Errorf("time conflict, available slots: %v", availableSlots)
	}

	oldRedisKey := fmt.Sprintf("appointment:%d:%s", appointment.TechnicianID, appointment.AppointmentDate.Format("2006-01-02"))
	redis.SRem(oldRedisKey, fmt.Sprintf("%s-%s", appointment.StartTime, appointment.EndTime))

	appointment.AppointmentDate = appointmentDate
	appointment.StartTime = req.StartTime
	appointment.EndTime = endTime

	if err := s.appointmentRepo.Update(appointment); err != nil {
		return nil, fmt.Errorf("update appointment: %w", err)
	}

	newRedisKey := fmt.Sprintf("appointment:%d:%s", appointment.TechnicianID, req.AppointmentDate)
	redis.SAdd(newRedisKey, fmt.Sprintf("%s-%s", req.StartTime, endTime))

	return appointment, nil
}

func (s *AppointmentService) GetAvailableSlots(technicianID uint, dateStr string, duration int) (*AvailableSlot, error) {
	technician, err := s.technicianRepo.GetByID(technicianID)
	if err != nil {
		return nil, fmt.Errorf("technician not found: %w", err)
	}

	date, _ := time.Parse("2006-01-02", dateStr)

	isLeave, _ := s.technicianRepo.IsOnLeave(technicianID, date)
	if isLeave {
		return &AvailableSlot{Date: dateStr, Slots: []string{}}, nil
	}

	appointments, err := s.appointmentRepo.GetByTechnicianAndDate(technicianID, date)
	if err != nil {
		return nil, err
	}

	workStart := utils.TimeToMinutes(technician.WorkStartTime)
	workEnd := utils.TimeToMinutes(technician.WorkEndTime)

	var slots []string
	for start := workStart; start+duration <= workEnd; start += 30 {
		slotStart := utils.MinutesToTime(start)
		slotEnd := utils.MinutesToTime(start + duration)

		conflict := false
		for _, appt := range appointments {
			if appt.Status != "cancelled" && utils.TimeOverlap(slotStart, slotEnd, appt.StartTime, appt.EndTime) {
				conflict = true
				break
			}
		}

		if !conflict {
			slots = append(slots, fmt.Sprintf("%s-%s", slotStart, slotEnd))
		}
	}

	return &AvailableSlot{
		Date:  dateStr,
		Slots: slots,
	}, nil
}

func (s *AppointmentService) Complete(id uint) error {
	appointment, err := s.appointmentRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("appointment not found: %w", err)
	}

	if appointment.Status != "confirmed" {
		return fmt.Errorf("appointment cannot be completed")
	}

	appointment.Status = "completed"
	if err := s.appointmentRepo.Update(appointment); err != nil {
		return fmt.Errorf("update appointment: %w", err)
	}

	if err := s.customerRepo.AddVisit(appointment.CustomerID); err != nil {
		return fmt.Errorf("add visit: %w", err)
	}

	return nil
}

func (s *AppointmentService) CheckProductStock(serviceID uint) error {
	serviceItem, err := s.serviceRepo.GetByID(serviceID)
	if err != nil {
		return fmt.Errorf("service not found: %w", err)
	}

	if serviceItem.Products == "" {
		return nil
	}

	var serviceProducts []ServiceProduct
	if err := json.Unmarshal([]byte(serviceItem.Products), &serviceProducts); err != nil {
		return nil
	}

	var insufficientProducts []string
	for _, sp := range serviceProducts {
		product, err := s.productRepo.GetByID(sp.ProductID)
		if err != nil {
			continue
		}
		if product.Stock < sp.Quantity {
			insufficientProducts = append(insufficientProducts,
				fmt.Sprintf("%s(需要%d, 现有%d)", product.Name, sp.Quantity, product.Stock))
		}
	}

	if len(insufficientProducts) > 0 {
		return fmt.Errorf("库存不足：%s", joinStrings(insufficientProducts, ", "))
	}

	return nil
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
