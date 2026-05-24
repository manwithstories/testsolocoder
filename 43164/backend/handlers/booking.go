package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"tutoring-platform/database"
	"tutoring-platform/models"
)

type BookingRequest struct {
	TeacherID uuid.UUID `json:"teacherId" binding:"required"`
	SubjectID uuid.UUID `json:"subjectId" binding:"required"`
	StartTime time.Time `json:"startTime" binding:"required"`
	EndTime   time.Time `json:"endTime" binding:"required"`
	Notes     string    `json:"notes"`
}

type RescheduleRequest struct {
	BookingID uuid.UUID `json:"bookingId" binding:"required"`
	NewStartTime time.Time `json:"newStartTime" binding:"required"`
	NewEndTime   time.Time `json:"newEndTime" binding:"required"`
	Reason      string    `json:"reason"`
}

type CancelBookingRequest struct {
	BookingID  uuid.UUID `json:"bookingId" binding:"required"`
	Reason     string    `json:"reason"`
}

func CreateBooking(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.StartTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot book in the past"})
		return
	}

	if req.EndTime.Before(req.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
		return
	}

	duration := int(req.EndTime.Sub(req.StartTime).Minutes())
	if duration < 30 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Minimum booking duration is 30 minutes"})
		return
	}

	var teacherProfile models.TeacherProfile
	if err := database.DB.Preload("User").Where("user_id = ? AND approval_status = ?", req.TeacherID, "approved").First(&teacherProfile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found or not approved"})
		return
	}

	hasConflict, err := checkBookingConflict(req.TeacherID, req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check booking conflict"})
		return
	}
	if hasConflict {
		c.JSON(http.StatusConflict, gin.H{"error": "Teacher already has a booking at this time"})
		return
	}

	studentConflict, err := checkStudentBookingConflict(userID.(uuid.UUID), req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check student booking conflict"})
		return
	}
	if studentConflict {
		c.JSON(http.StatusConflict, gin.H{"error": "You already have a booking at this time"})
		return
	}

	isAvailable, err := checkTeacherAvailability(req.TeacherID, req.StartTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check availability"})
		return
	}
	if !isAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Teacher is not available at this time"})
		return
	}

	var subject models.Subject
	if err := database.DB.Where("id = ?", req.SubjectID).First(&subject).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}

	hourlyRate := teacherProfile.HourlyRate
	var teacherSubject models.TeacherSubject
	if database.DB.Where("teacher_id = ? AND subject_id = ?", teacherProfile.ID, req.SubjectID).First(&teacherSubject).Error == nil {
		if teacherSubject.CustomRate > 0 {
			hourlyRate = teacherSubject.CustomRate
		}
	}

	totalAmount := hourlyRate * float64(duration) / 60

	booking := models.Booking{
		StudentID:   userID.(uuid.UUID),
		TeacherID:   req.TeacherID,
		SubjectID:   req.SubjectID,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Duration:    duration,
		HourlyRate:  hourlyRate,
		TotalAmount: totalAmount,
		Status:      models.BookingStatusPending,
		Notes:       req.Notes,
	}

	if err := database.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	createNotification(c, req.TeacherID, models.NotificationTypeBookingCreated, "New Booking Request", "You have a new booking request from a student")

	c.JSON(http.StatusCreated, gin.H{"message": "Booking created successfully", "booking": booking})
}

func checkBookingConflict(teacherID uuid.UUID, startTime, endTime time.Time) (bool, error) {
	var count int64
	err := database.DB.Model(&models.Booking{}).
		Where("teacher_id = ? AND status IN ? AND start_time < ? AND end_time > ?",
			teacherID,
			[]models.BookingStatus{models.BookingStatusPending, models.BookingStatusConfirmed},
			endTime,
			startTime,
		).Count(&count).Error
	return count > 0, err
}

func checkStudentBookingConflict(studentID uuid.UUID, startTime, endTime time.Time) (bool, error) {
	var count int64
	err := database.DB.Model(&models.Booking{}).
		Where("student_id = ? AND status IN ? AND start_time < ? AND end_time > ?",
			studentID,
			[]models.BookingStatus{models.BookingStatusPending, models.BookingStatusConfirmed},
			endTime,
			startTime,
		).Count(&count).Error
	return count > 0, err
}

func checkTeacherAvailability(teacherID uuid.UUID, startTime time.Time) (bool, error) {
	dayOfWeek := int(startTime.Weekday())
	timeStr := startTime.Format("15:04")

	var profile models.TeacherProfile
	if err := database.DB.Where("user_id = ?", teacherID).First(&profile).Error; err != nil {
		return false, err
	}

	var count int64
	err := database.DB.Model(&models.AvailabilitySlot{}).
		Where("teacher_id = ? AND day_of_week = ? AND start_time <= ? AND end_time > ? AND is_recurring = ?",
			profile.ID, dayOfWeek, timeStr, timeStr, true,
		).Count(&count).Error

	return count > 0, err
}

func GetBookings(c *gin.Context) {
	userID, _ := c.Get("userId")
	userRole, _ := c.Get("userRole")
	status := c.Query("status")

	var bookings []models.Booking
	query := database.DB.Preload("Student").Preload("Teacher").Preload("Subject")

	if userRole == models.RoleStudent {
		query = query.Where("student_id = ?", userID)
	} else if userRole == models.RoleTeacher {
		query = query.Where("teacher_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("start_time DESC").Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func GetBookingByID(c *gin.Context) {
	id := c.Param("id")
	bookingID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var booking models.Booking
	if err := database.DB.Preload("Student").Preload("Teacher").Preload("Subject").Preload("VideoSession").Where("id = ?", bookingID).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func ConfirmBooking(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userId")

	var booking models.Booking
	if err := database.DB.Where("id = ? AND teacher_id = ?", id, userID).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if booking.Status != models.BookingStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking is not pending"})
		return
	}

	booking.Status = models.BookingStatusConfirmed
	if err := database.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm booking"})
		return
	}

	createNotification(c, booking.StudentID, models.NotificationTypeBookingCreated, "Booking Confirmed", "Your booking has been confirmed by the teacher")

	c.JSON(http.StatusOK, gin.H{"message": "Booking confirmed successfully", "booking": booking})
}

func RescheduleBooking(c *gin.Context) {
	userID, _ := c.Get("userId")
	userRole, _ := c.Get("userRole")

	var req RescheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.NewStartTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot reschedule to past time"})
		return
	}

	var booking models.Booking
	if err := database.DB.Where("id = ?", req.BookingID).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if userRole == models.RoleStudent && booking.StudentID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}
	if userRole == models.RoleTeacher && booking.TeacherID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if booking.StartTime.Before(time.Now().Add(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot reschedule within 24 hours of the booking"})
		return
	}

	hasConflict, _ := checkBookingConflict(booking.TeacherID, req.NewStartTime, req.NewEndTime)
	if hasConflict {
		c.JSON(http.StatusConflict, gin.H{"error": "Teacher already has a booking at this time"})
		return
	}

	studentConflict, _ := checkStudentBookingConflict(booking.StudentID, req.NewStartTime, req.NewEndTime)
	if studentConflict {
		c.JSON(http.StatusConflict, gin.H{"error": "Student already has a booking at this time"})
		return
	}

	history := models.BookingRescheduleHistory{
		BookingID:     booking.ID,
		OldStartTime:  booking.StartTime,
		OldEndTime:    booking.EndTime,
		NewStartTime:  req.NewStartTime,
		NewEndTime:    req.NewEndTime,
		RescheduledBy: userID.(uuid.UUID),
		Reason:        req.Reason,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record reschedule"})
		return
	}

	duration := int(req.NewEndTime.Sub(req.NewStartTime).Minutes())
	booking.StartTime = req.NewStartTime
	booking.EndTime = req.NewEndTime
	booking.Duration = duration
	booking.Status = models.BookingStatusRescheduled

	if err := tx.Save(&booking).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reschedule booking"})
		return
	}

	tx.Commit()

	var notifyUserID uuid.UUID
	if userRole == models.RoleStudent {
		notifyUserID = booking.TeacherID
	} else {
		notifyUserID = booking.StudentID
	}
	createNotification(c, notifyUserID, models.NotificationTypeBookingUpdated, "Booking Rescheduled", "A booking has been rescheduled")

	c.JSON(http.StatusOK, gin.H{"message": "Booking rescheduled successfully", "booking": booking})
}

func CancelBooking(c *gin.Context) {
	userID, _ := c.Get("userId")
	userRole, _ := c.Get("userRole")

	var req CancelBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var booking models.Booking
	if err := database.DB.Where("id = ?", req.BookingID).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if userRole == models.RoleStudent && booking.StudentID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}
	if userRole == models.RoleTeacher && booking.TeacherID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if booking.StartTime.Before(time.Now().Add(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot cancel within 24 hours of the booking"})
		return
	}

	now := time.Now()
	uid := userID.(uuid.UUID)
	booking.Status = models.BookingStatusCancelled
	booking.CancelledBy = &uid
	booking.CancelledAt = &now
	booking.CancelReason = req.Reason

	if err := database.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking"})
		return
	}

	var notifyUserID uuid.UUID
	if userRole == models.RoleStudent {
		notifyUserID = booking.TeacherID
	} else {
		notifyUserID = booking.StudentID
	}
	createNotification(c, notifyUserID, models.NotificationTypeBookingCancelled, "Booking Cancelled", "A booking has been cancelled")

	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})
}

func CompleteBooking(c *gin.Context) {
	id := c.Param("id")

	var booking models.Booking
	if err := database.DB.Where("id = ?", id).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if booking.Status != models.BookingStatusConfirmed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking is not confirmed"})
		return
	}

	booking.Status = models.BookingStatusCompleted
	if err := database.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete booking"})
		return
	}

	var teacherProfile models.TeacherProfile
	hourlyRate := 0.0
	if err := database.DB.Where("user_id = ?", booking.TeacherID).First(&teacherProfile).Error; err == nil {
		hourlyRate = teacherProfile.HourlyRate
	}

	duration := booking.EndTime.Sub(booking.StartTime).Hours()
	totalAmount := duration * hourlyRate

	platformCommission := totalAmount * 0.10
	teacherEarnings := totalAmount - platformCommission

	var studentWallet models.Wallet
	if err := database.DB.Where("user_id = ?", booking.StudentID).First(&studentWallet).Error; err != nil {
		studentWallet = models.Wallet{
			UserID:   booking.StudentID,
			Balance:  0,
			Currency: "USD",
		}
		database.DB.Create(&studentWallet)
	}

	var teacherWallet models.Wallet
	if err := database.DB.Where("user_id = ?", booking.TeacherID).First(&teacherWallet).Error; err != nil {
		teacherWallet = models.Wallet{
			UserID:   booking.TeacherID,
			Balance:  0,
			Currency: "USD",
		}
		database.DB.Create(&teacherWallet)
	}

	var platformWallet models.Wallet
	platformUserID := uuid.Nil
	if err := database.DB.Where("user_id = ?", platformUserID).First(&platformWallet).Error; err != nil {
		platformWallet = models.Wallet{
			UserID:   platformUserID,
			Balance:  0,
			Currency: "USD",
		}
		database.DB.Create(&platformWallet)
	}

	now := time.Now()
	bookingID := booking.ID

	studentTransaction := models.Transaction{
		WalletID:      studentWallet.ID,
		UserID:        booking.StudentID,
		Type:          models.TransactionTypePayment,
		Amount:        -totalAmount,
		Currency:      "USD",
		BalanceAfter:  studentWallet.Balance - totalAmount,
		Status:        models.TransactionStatusCompleted,
		Description:   "Payment for completed lesson",
		ReferenceID:   booking.ID.String(),
		ReferenceType: "booking",
		BookingID:     &bookingID,
		CompletedAt:   &now,
	}
	database.DB.Create(&studentTransaction)

	studentWallet.Balance -= totalAmount
	studentWallet.TotalSpent += totalAmount
	database.DB.Save(&studentWallet)

	teacherTransaction := models.Transaction{
		WalletID:      teacherWallet.ID,
		UserID:        booking.TeacherID,
		Type:          models.TransactionTypePayment,
		Amount:        teacherEarnings,
		Currency:      "USD",
		BalanceAfter:  teacherWallet.Balance + teacherEarnings,
		Status:        models.TransactionStatusCompleted,
		Description:   "Earnings from completed lesson",
		ReferenceID:   booking.ID.String(),
		ReferenceType: "booking",
		BookingID:     &bookingID,
		CompletedAt:   &now,
	}
	database.DB.Create(&teacherTransaction)

	teacherWallet.Balance += teacherEarnings
	teacherWallet.TotalIncome += teacherEarnings
	database.DB.Save(&teacherWallet)

	platformTransaction := models.Transaction{
		WalletID:      platformWallet.ID,
		UserID:        platformUserID,
		Type:          models.TransactionTypeCommission,
		Amount:        platformCommission,
		Currency:      "USD",
		BalanceAfter:  platformWallet.Balance + platformCommission,
		Status:        models.TransactionStatusCompleted,
		Description:   "Platform commission",
		ReferenceID:   booking.ID.String(),
		ReferenceType: "booking",
		BookingID:     &bookingID,
		CompletedAt:   &now,
	}
	database.DB.Create(&platformTransaction)

	platformWallet.Balance += platformCommission
	database.DB.Save(&platformWallet)

	createNotification(c, booking.StudentID, models.NotificationTypeSystem, "Lesson Completed", "Your lesson has been completed and payment processed")
	createNotification(c, booking.TeacherID, models.NotificationTypeSystem, "Lesson Completed", "Your lesson earnings have been added to your wallet")

	c.JSON(http.StatusOK, gin.H{"message": "Booking completed successfully"})
}

func createNotification(c *gin.Context, userID uuid.UUID, notificationType models.NotificationType, title, content string) {
	notification := models.Notification{
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Content: content,
	}
	database.DB.Create(&notification)
}

var _ = errors.New("")
