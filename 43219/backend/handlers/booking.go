package handlers

import (
	"time"

	"housekeeping/config"
	"housekeeping/database"
	"housekeeping/models"
	"housekeeping/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingInput struct {
	ServiceID uint      `json:"service_id" binding:"required"`
	StartAt   time.Time `json:"start_at" binding:"required"`
	EndAt     time.Time `json:"end_at" binding:"required"`
	Address   string    `json:"address" binding:"required"`
	Remark    string    `json:"remark"`
	StaffID   *uint     `json:"staff_id,omitempty"`
	Price     float64   `json:"price"`
}

func hasConflict(tx *gorm.DB, staffID uint, start, end time.Time, excludeBookingID *uint) (bool, error) {
	q := tx.Model(&models.Booking{}).Where(
		"staff_id = ? AND status IN ? AND start_at < ? AND end_at > ?",
		staffID,
		[]models.BookingStatus{models.BookingPending, models.BookingConfirmed},
		end,
		start,
	)
	if excludeBookingID != nil {
		q = q.Where("id <> ?", *excludeBookingID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func findAvailableStaff(tx *gorm.DB, svc models.Service, start, end time.Time) (*models.User, error) {
	var staffs []models.User
	if err := tx.Where("role = ? AND company_id = ? AND suspended = ?", models.RoleStaff, svc.CompanyID, false).
		Find(&staffs).Error; err != nil {
		return nil, err
	}
	for _, s := range staffs {
		ok, err := hasConflict(tx, s.ID, start, end, nil)
		if err != nil {
			return nil, err
		}
		if !ok {
			return &s, nil
		}
	}
	return nil, nil
}

func CreateBooking(c *gin.Context) {
	uid, _ := c.Get("uid")
	var in BookingInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if in.EndAt.Before(in.StartAt) || in.EndAt.Equal(in.StartAt) {
		utils.BadRequest(c, "invalid time range")
		return
	}
	var svc models.Service
	if err := database.DB.First(&svc, in.ServiceID).Error; err != nil {
		utils.NotFound(c, "service not found")
		return
	}
	price := in.Price
	if price <= 0 {
		price = (svc.MinPrice + svc.MaxPrice) / 2
	}

	var booking models.Booking
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var staffID *uint
		if in.StaffID != nil {
			ok, err := hasConflict(tx, *in.StaffID, in.StartAt, in.EndAt, nil)
			if err != nil {
				return err
			}
			if ok {
				return errStaffBusy
			}
			staffID = in.StaffID
		} else {
			s, err := findAvailableStaff(tx, svc, in.StartAt, in.EndAt)
			if err != nil {
				return err
			}
			if s == nil {
				return errNoAvailable
			}
			staffID = &s.ID
		}
		booking = models.Booking{
			CustomerID: uid.(uint),
			ServiceID:  svc.ID,
			StaffID:    staffID,
			StartAt:    in.StartAt,
			EndAt:      in.EndAt,
			Address:    in.Address,
			Remark:     in.Remark,
			Price:      price,
			Status:     models.BookingPending,
		}
		return tx.Create(&booking).Error
	})
	if err != nil {
		if err == errStaffBusy {
			utils.BadRequest(c, "selected staff busy")
			return
		}
		if err == errNoAvailable {
			utils.BadRequest(c, "no available staff at selected time")
			return
		}
		utils.ServerError(c, "create booking failed")
		return
	}
	utils.Logger.Infow("booking created", "id", booking.ID, "customer", uid)
	utils.OK(c, booking)
}

func ConfirmBooking(c *gin.Context) {
	id := c.Param("id")
	var bk models.Booking
	if err := database.DB.First(&bk, id).Error; err != nil {
		utils.NotFound(c, "booking not found")
		return
	}
	uid, _ := c.Get("uid")
	role, _ := c.Get("role")
	if role.(string) == string(models.RoleStaff) && bk.StaffID != nil && *bk.StaffID != uid.(uint) {
		utils.Forbidden(c, "not your booking")
		return
	}
	if role.(string) == string(models.RoleCustomer) && bk.CustomerID != uid.(uint) {
		utils.Forbidden(c, "not your booking")
		return
	}
	bk.Status = models.BookingConfirmed
	database.DB.Save(&bk)
	utils.OK(c, bk)
}

func RescheduleBooking(c *gin.Context) {
	id := c.Param("id")
	uid, _ := c.Get("uid")
	var bk models.Booking
	if err := database.DB.First(&bk, id).Error; err != nil {
		utils.NotFound(c, "booking not found")
		return
	}
	if bk.CustomerID != uid.(uint) {
		utils.Forbidden(c, "not your booking")
		return
	}
	var body struct {
		StartAt time.Time `json:"start_at" binding:"required"`
		EndAt   time.Time `json:"end_at" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if bk.RescheduleCount >= config.C.Order.MaxReschedule {
		bk.NeedReview = true
		database.DB.Save(&bk)
		utils.BadRequest(c, "reschedule exceeds limit, pending manual review")
		return
	}
	var svc models.Service
	database.DB.First(&svc, bk.ServiceID)
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if bk.StaffID != nil {
			ok, err := hasConflict(tx, *bk.StaffID, body.StartAt, body.EndAt, &bk.ID)
			if err != nil {
				return err
			}
			if !ok {
				bk.StartAt = body.StartAt
				bk.EndAt = body.EndAt
				bk.RescheduleCount++
				return tx.Save(&bk).Error
			}
		}
		s, err := findAvailableStaff(tx, svc, body.StartAt, body.EndAt)
		if err != nil {
			return err
		}
		if s == nil {
			return errNoAvailable
		}
		bk.StaffID = &s.ID
		bk.StartAt = body.StartAt
		bk.EndAt = body.EndAt
		bk.RescheduleCount++
		return tx.Save(&bk).Error
	})
	if err != nil {
		if err == errNoAvailable {
			utils.BadRequest(c, "no available staff at new time")
			return
		}
		utils.ServerError(c, "reschedule failed")
		return
	}
	utils.OK(c, bk)
}

func CancelBooking(c *gin.Context) {
	id := c.Param("id")
	uid, _ := c.Get("uid")
	var bk models.Booking
	if err := database.DB.First(&bk, id).Error; err != nil {
		utils.NotFound(c, "booking not found")
		return
	}
	if bk.CustomerID != uid.(uint) {
		utils.Forbidden(c, "not your booking")
		return
	}
	if bk.Status != models.BookingPending && bk.Status != models.BookingConfirmed {
		utils.BadRequest(c, "cannot cancel")
		return
	}
	bk.Status = models.BookingCanceled
	database.DB.Save(&bk)
	utils.OK(c, "canceled")
}

func ListBookings(c *gin.Context) {
	uid, _ := c.Get("uid")
	role, _ := c.Get("role")
	var list []models.Booking
	q := database.DB.Preload("Service").Preload("Staff")
	switch role.(string) {
	case string(models.RoleCustomer):
		q = q.Where("customer_id = ?", uid)
	case string(models.RoleStaff):
		q = q.Where("staff_id = ?", uid)
	case string(models.RoleCompany):
		var staffIDs []uint
		database.DB.Model(&models.User{}).Where("company_id = ? AND role = ?", uid, models.RoleStaff).Pluck("id", &staffIDs)
		q = q.Where("staff_id IN ? OR EXISTS (SELECT 1 FROM services WHERE services.id = bookings.service_id AND services.company_id = ?)", staffIDs, uid)
	}
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("id desc").Find(&list).Error; err != nil {
		utils.ServerError(c, "query failed")
		return
	}
	utils.OK(c, list)
}

func GetBooking(c *gin.Context) {
	id := c.Param("id")
	var bk models.Booking
	if err := database.DB.Preload("Service").Preload("Staff").Preload("Customer").First(&bk, id).Error; err != nil {
		utils.NotFound(c, "booking not found")
		return
	}
	utils.OK(c, bk)
}

func ReviewNeedReschedule(c *gin.Context) {
	id := c.Param("id")
	var bk models.Booking
	if err := database.DB.First(&bk, id).Error; err != nil {
		utils.NotFound(c, "booking not found")
		return
	}
	var body struct {
		Approve bool      `json:"approve"`
		StartAt time.Time `json:"start_at"`
		EndAt   time.Time `json:"end_at"`
	}
	c.ShouldBindJSON(&body)
	if !body.Approve {
		bk.NeedReview = false
		bk.Status = models.BookingCanceled
		database.DB.Save(&bk)
		utils.OK(c, "rejected")
		return
	}
	bk.StartAt = body.StartAt
	bk.EndAt = body.EndAt
	bk.NeedReview = false
	bk.RescheduleCount = 0
	database.DB.Save(&bk)
	utils.OK(c, bk)
}

var (
	errStaffBusy   = &errMsg{"selected staff is busy"}
	errNoAvailable = &errMsg{"no available staff"}
)

type errMsg struct{ s string }

func (e *errMsg) Error() string { return e.s }
