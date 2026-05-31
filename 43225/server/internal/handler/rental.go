package handler

import (
	"net/http"
	"time"

	"ship-rental-platform/internal/config"
	"ship-rental-platform/internal/database"
	"ship-rental-platform/internal/model"
	"ship-rental-platform/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RentalHandler struct{}

func NewRentalHandler() *RentalHandler {
	return &RentalHandler{}
}

func (h *RentalHandler) CreateRental(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tenantID, _ := uuid.Parse(userID.(string))

	var req model.CreateRentalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.EndDate.Before(req.StartDate) {
		utils.Error(c, http.StatusBadRequest, "End date must be after start date")
		return
	}

	shipID, _ := uuid.Parse(req.ShipID)

	var ship model.Ship
	if err := database.DB.First(&ship, "id = ?", shipID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Ship not found")
		return
	}

	if ship.Status != model.ShipStatusAvailable {
		utils.Error(c, http.StatusBadRequest, "Ship is not available for rental")
		return
	}

	var conflictingRentals []model.Rental
	database.DB.Where(
		"ship_id = ? AND status IN ? AND (start_date < ? AND end_date > ?)",
		shipID,
		[]model.RentalStatus{model.RentalStatusPending, model.RentalStatusConfirmed, model.RentalStatusActive},
		req.EndDate,
		req.StartDate,
	).Find(&conflictingRentals)

	if len(conflictingRentals) > 0 {
		utils.Error(c, http.StatusConflict, "Ship is already booked for this time period")
		return
	}

	cfg := config.GetConfig()
	currency := req.Currency
	if currency == "" {
		currency = cfg.Finance.DefaultCurrency
	}

	duration := req.EndDate.Sub(req.StartDate)
	hours := duration.Hours()
	var baseAmount float64

	switch req.RentalType {
	case model.RentalTypeHourly:
		baseAmount = ship.HourlyRate * hours
	case model.RentalTypeDaily, model.RentalTypeVoyage:
		days := int(hours / 24)
		if hours-float64(days*24) > 0 {
			days++
		}
		baseAmount = ship.DailyRate * float64(days)
	}

	insuranceAmount := 0.0
	if ship.InsuranceRequired || req.InsuranceType != model.InsuranceTypeNone {
		switch req.InsuranceType {
		case model.InsuranceTypeBasic:
			insuranceAmount = baseAmount * 0.05
		case model.InsuranceTypePremium:
			insuranceAmount = baseAmount * 0.10
		}
	}

	platformFee := (baseAmount + insuranceAmount) * cfg.Finance.PlatformFeeRate
	totalAmount := baseAmount + insuranceAmount + platformFee

	rental := model.Rental{
		TenantID:          tenantID,
		ShipID:            shipID,
		RentalType:        req.RentalType,
		StartDate:         req.StartDate,
		EndDate:           req.EndDate,
		StartLocation:     req.StartLocation,
		EndLocation:       req.EndLocation,
		BaseAmount:        utils.RoundTo(baseAmount, 2),
		InsuranceType:     req.InsuranceType,
		InsuranceAmount:   utils.RoundTo(insuranceAmount, 2),
		PlatformFee:       utils.RoundTo(platformFee, 2),
		TotalAmount:       utils.RoundTo(totalAmount, 2),
		DepositAmount:     ship.DepositAmount,
		DepositStatus:     "unpaid",
		Currency:          currency,
		Status:            model.RentalStatusPending,
		EmergencyContact:  req.EmergencyContact,
		EmergencyPhone:    req.EmergencyPhone,
		Notes:             req.Notes,
		CrewCount:         req.CrewCount,
		PassengerCount:    req.PassengerCount,
	}

	if err := database.DB.Create(&rental).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create rental")
		return
	}

	utils.Created(c, rental)
}

func (h *RentalHandler) GetRentals(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var rentals []model.Rental
	query := database.DB.Preload("Ship").Preload("Tenant")

	if role.(string) == string(model.RoleTenant) {
		query = query.Where("tenant_id = ?", userID)
	} else if role.(string) == string(model.RoleOwner) {
		query = query.Joins("JOIN ships ON ships.id = rentals.ship_id").Where("ships.owner_id = ?", userID)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("rentals.status = ?", status)
	}
	if shipID := c.Query("ship_id"); shipID != "" {
		query = query.Where("rentals.ship_id = ?", shipID)
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate != "" && endDate != "" {
		query = query.Where("start_date BETWEEN ? AND ?", startDate, endDate)
	}

	var total int64
	query.Model(&model.Rental{}).Count(&total)

	page := utils.ParseInt(c.DefaultQuery("page", "1"), 1)
	pageSize := utils.ParseInt(c.DefaultQuery("page_size", "10"), 10)
	offset := (page - 1) * pageSize

	if err := query.Order("rentals.created_at DESC").Offset(offset).Limit(pageSize).Find(&rentals).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch rentals")
		return
	}

	utils.Paginated(c, rentals, total, page, pageSize)
}

func (h *RentalHandler) GetRental(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid rental ID")
		return
	}

	var rental model.Rental
	if err := database.DB.Preload("Ship").Preload("Tenant").Preload("Ship.Owner").First(&rental, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Rental not found")
		return
	}

	utils.Success(c, rental)
}

func (h *RentalHandler) UpdateRentalStatus(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid rental ID")
		return
	}

	var req model.UpdateRentalStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var rental model.Rental
	if err := database.DB.First(&rental, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Rental not found")
		return
	}

	validTransitions := map[model.RentalStatus][]model.RentalStatus{
		model.RentalStatusPending:   {model.RentalStatusConfirmed, model.RentalStatusCancelled},
		model.RentalStatusConfirmed: {model.RentalStatusActive, model.RentalStatusCancelled},
		model.RentalStatusActive:    {model.RentalStatusCompleted, model.RentalStatusCancelled},
		model.RentalStatusCompleted: {model.RentalStatusRefunded},
	}

	validNextStatuses, exists := validTransitions[rental.Status]
	if !exists {
		utils.Error(c, http.StatusBadRequest, "Invalid status transition")
		return
	}

	isValid := false
	for _, status := range validNextStatuses {
		if status == req.Status {
			isValid = true
			break
		}
	}

	if !isValid {
		utils.Error(c, http.StatusBadRequest, "Invalid status transition from "+string(rental.Status)+" to "+string(req.Status))
		return
	}

	rental.Status = req.Status

	switch req.Status {
	case model.RentalStatusCancelled:
		now := time.Now()
		rental.CancelledAt = &now
		rental.CancellationReason = req.CancellationReason
	case model.RentalStatusCompleted:
		now := time.Now()
		rental.CompletedAt = &now
	case model.RentalStatusConfirmed:
		rental.DepositStatus = "paid"
	}

	if err := database.DB.Save(&rental).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update rental status")
		return
	}

	if req.Status == model.RentalStatusConfirmed {
		database.DB.Model(&model.Ship{}).Where("id = ?", rental.ShipID).Update("status", model.ShipStatusRented)
	} else if req.Status == model.RentalStatusCompleted || req.Status == model.RentalStatusCancelled {
		var activeRentals int64
		database.DB.Model(&model.Rental{}).Where(
			"ship_id = ? AND status IN ?",
			rental.ShipID,
			[]model.RentalStatus{model.RentalStatusConfirmed, model.RentalStatusActive},
		).Count(&activeRentals)

		if activeRentals == 0 {
			database.DB.Model(&model.Ship{}).Where("id = ?", rental.ShipID).Update("status", model.ShipStatusAvailable)
		}
	}

	utils.Success(c, rental)
}

func (h *RentalHandler) CancelRental(c *gin.Context) {
	id := c.Param("id")
	if !utils.IsValidUUID(id) {
		utils.Error(c, http.StatusBadRequest, "Invalid rental ID")
		return
	}

	userID, _ := c.Get("user_id")
	tenantID, _ := uuid.Parse(userID.(string))

	var rental model.Rental
	if err := database.DB.First(&rental, "id = ?", id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Rental not found")
		return
	}

	if rental.TenantID != tenantID {
		role, _ := c.Get("role")
		if role.(string) != string(model.RoleAdmin) {
			utils.Error(c, http.StatusForbidden, "You can only cancel your own rentals")
			return
		}
	}

	if rental.Status == model.RentalStatusCompleted || rental.Status == model.RentalStatusCancelled {
		utils.Error(c, http.StatusBadRequest, "Rental cannot be cancelled in current status")
		return
	}

	now := time.Now()
	rental.Status = model.RentalStatusCancelled
	rental.CancelledAt = &now
	rental.CancellationReason = c.PostForm("reason")

	if err := database.DB.Save(&rental).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to cancel rental")
		return
	}

	utils.Success(c, gin.H{"message": "Rental cancelled successfully"})
}

func (h *RentalHandler) GetMyRentals(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tenantID, _ := uuid.Parse(userID.(string))

	var rentals []model.Rental
	if err := database.DB.Preload("Ship").Where("tenant_id = ?", tenantID).Order("created_at DESC").Find(&rentals).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch rentals")
		return
	}

	utils.Success(c, rentals)
}
