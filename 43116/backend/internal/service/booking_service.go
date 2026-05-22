package service

import (
	"car-rental/internal/model"
	"car-rental/internal/repository"
	"car-rental/internal/utils"
	"errors"
	"time"
)

type BookingService struct {
	bookingRepo   *repository.BookingRepository
	carRepo       *repository.CarRepository
	orderRepo     *repository.OrderRepository
	promoRepo     *repository.PromoCodeRepository
	pricingRepo   *repository.PricingRuleRepository
	userRepo      *repository.UserRepository
	messageService *MessageService
}

func NewBookingService(messageService *MessageService) *BookingService {
	return &BookingService{
		bookingRepo:    repository.NewBookingRepository(),
		carRepo:        repository.NewCarRepository(),
		orderRepo:      repository.NewOrderRepository(),
		promoRepo:      repository.NewPromoCodeRepository(),
		pricingRepo:    repository.NewPricingRuleRepository(),
		userRepo:       repository.NewUserRepository(),
		messageService: messageService,
	}
}

type CreateBookingRequest struct {
	CarID         uint      `json:"car_id" binding:"required"`
	PickupStoreID uint      `json:"pickup_store_id" binding:"required"`
	ReturnStoreID uint      `json:"return_store_id" binding:"required"`
	PickupTime    time.Time `json:"pickup_time" binding:"required"`
	ReturnTime    time.Time `json:"return_time" binding:"required"`
	PromoCode     string    `json:"promo_code"`
}

type PriceCalculation struct {
	BasePrice    float64 `json:"base_price"`
	TotalPrice   float64 `json:"total_price"`
	Discount     float64 `json:"discount"`
	FinalPrice   float64 `json:"final_price"`
	TotalDays    int     `json:"total_days"`
	PriceDetails []DayPrice `json:"price_details"`
}

type DayPrice struct {
	Date       time.Time `json:"date"`
	BasePrice  float64   `json:"base_price"`
	Multiplier float64   `json:"multiplier"`
	DayPrice   float64   `json:"day_price"`
}

func (s *BookingService) CalculatePrice(carID uint, pickupTime, returnTime time.Time, promoCode string) (*PriceCalculation, error) {
	car, err := s.carRepo.FindByID(carID)
	if err != nil {
		return nil, errors.New("车辆不存在")
	}

	if !returnTime.After(pickupTime) {
		return nil, errors.New("还车时间必须晚于取车时间")
	}

	totalDays := utils.CalculateDays(pickupTime, returnTime)
	pricingRules, err := s.pricingRepo.GetActiveRules()
	if err != nil {
		return nil, err
	}

	var priceDetails []DayPrice
	var totalPrice float64

	for i := 0; i < totalDays; i++ {
		currentDate := pickupTime.AddDate(0, 0, i)
		dayBasePrice := car.DailyRent
		multiplier := 1.0

		for _, rule := range pricingRules {
			if rule.RuleType == "weekend" && utils.IsWeekend(currentDate) {
				multiplier = rule.Multiplier
				break
			} else if rule.RuleType == "holiday" {
				if rule.StartDate != nil && rule.EndDate != nil {
					if !currentDate.Before(*rule.StartDate) && !currentDate.After(*rule.EndDate) {
						multiplier = rule.Multiplier
						break
					}
				}
			}
		}

		dayPrice := dayBasePrice * multiplier
		priceDetails = append(priceDetails, DayPrice{
			Date:       currentDate,
			BasePrice:  dayBasePrice,
			Multiplier: multiplier,
			DayPrice:   dayPrice,
		})
		totalPrice += dayPrice
	}

	calc := &PriceCalculation{
		BasePrice:    car.DailyRent * float64(totalDays),
		TotalPrice:   totalPrice,
		TotalDays:    totalDays,
		PriceDetails: priceDetails,
	}

	if promoCode != "" {
		promo, err := s.promoRepo.FindByCode(promoCode)
		if err == nil && promo.IsActive {
			if promo.UsedCount < promo.UsageLimit {
				if totalPrice >= promo.MinAmount {
					discount := totalPrice * promo.Value / 100
					if promo.MaxDiscount > 0 && discount > promo.MaxDiscount {
						discount = promo.MaxDiscount
					}
					calc.Discount = discount
					calc.FinalPrice = totalPrice - discount
				}
			}
		}
	}

	if calc.FinalPrice == 0 {
		calc.FinalPrice = totalPrice
	}

	return calc, nil
}

func (s *BookingService) CreateBooking(userID uint, req *CreateBookingRequest) (*model.Booking, error) {
	car, err := s.carRepo.FindByID(req.CarID)
	if err != nil {
		return nil, errors.New("车辆不存在")
	}

	if car.Status != model.CarStatusAvailable {
		return nil, errors.New("车辆当前不可用")
	}

	if !s.bookingRepo.IsCarAvailable(req.CarID, req.PickupTime, req.ReturnTime) {
		return nil, errors.New("所选时间段车辆已被预订")
	}

	if !req.ReturnTime.After(req.PickupTime) {
		return nil, errors.New("还车时间必须晚于取车时间")
	}

	priceCalc, err := s.CalculatePrice(req.CarID, req.PickupTime, req.ReturnTime, req.PromoCode)
	if err != nil {
		return nil, err
	}

	var promoCodeID *uint
	if req.PromoCode != "" {
		promo, err := s.promoRepo.FindByCode(req.PromoCode)
		if err == nil && promo.IsActive {
			id := promo.ID
			promoCodeID = &id
			err = s.promoRepo.IncrementUsedCount(promo.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	bookingNo := utils.GenerateBookingNo()
	booking := &model.Booking{
		BookingNo:     bookingNo,
		UserID:        userID,
		CarID:         req.CarID,
		PickupStoreID: req.PickupStoreID,
		ReturnStoreID: req.ReturnStoreID,
		PickupTime:    req.PickupTime,
		ReturnTime:    req.ReturnTime,
		TotalDays:     priceCalc.TotalDays,
		BasePrice:     priceCalc.BasePrice,
		TotalPrice:    priceCalc.TotalPrice,
		Discount:     priceCalc.Discount,
		FinalPrice:   priceCalc.FinalPrice,
		Deposit:      car.Deposit,
		PromoCodeID:  promoCodeID,
		Status:       model.BookingStatusPending,
	}

	err = s.bookingRepo.Create(booking)
	if err != nil {
		return nil, err
	}

	orderNo := utils.GenerateOrderNo()
	order := &model.Order{
		OrderNo:       orderNo,
		UserID:        userID,
		CarID:         req.CarID,
		BookingID:     booking.ID,
		TotalAmount:   priceCalc.TotalPrice,
		Discount:      priceCalc.Discount,
		FinalAmount:   priceCalc.FinalPrice,
		PaymentStatus: model.PaymentStatusPending,
		Status:        model.OrderStatusPending,
	}

	err = s.orderRepo.Create(order)
	if err != nil {
		return nil, err
	}

	booking.OrderID = &order.ID
	err = s.bookingRepo.Update(booking)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userID)
	if err == nil && s.messageService != nil {
		s.messageService.SendBookingConfirmation(userID, user.Email, booking.BookingNo, booking.ID)
	}

	return booking, nil
}

func (s *BookingService) GetBookingByID(id uint) (*model.Booking, error) {
	return s.bookingRepo.FindByID(id)
}

func (s *BookingService) GetBookingByNo(bookingNo string) (*model.Booking, error) {
	return s.bookingRepo.FindByBookingNo(bookingNo)
}

func (s *BookingService) GetAllBookings(page, pageSize int, userID uint, status string, carID uint, startDate, endDate *time.Time) ([]model.Booking, int64, error) {
	return s.bookingRepo.FindAll(page, pageSize, userID, status, carID, startDate, endDate)
}

func (s *BookingService) ConfirmBooking(id uint) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return errors.New("预订不存在")
	}

	if booking.Status != model.BookingStatusPending {
		return errors.New("预订状态不允许确认")
	}

	err = s.bookingRepo.UpdateStatus(id, model.BookingStatusConfirmed)
	if err != nil {
		return err
	}

	return s.carRepo.UpdateStatus(booking.CarID, model.CarStatusRented)
}

func (s *BookingService) CancelBooking(id uint, reason string) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return errors.New("预订不存在")
	}

	if booking.Status == model.BookingStatusCancelled || booking.Status == model.BookingStatusCompleted {
		return errors.New("预订已取消或完成")
	}

	err = s.bookingRepo.Cancel(id, reason)
	if err != nil {
		return err
	}

	return s.carRepo.UpdateStatus(booking.CarID, model.CarStatusAvailable)
}

func (s *BookingService) CompleteBooking(id uint) error {
	booking, err := s.bookingRepo.FindByID(id)
	if err != nil {
		return errors.New("预订不存在")
	}

	if booking.Status != model.BookingStatusConfirmed {
		return errors.New("预订状态不允许完成")
	}

	now := time.Now()
	booking.ActualReturnTime = &now
	booking.Status = model.BookingStatusCompleted
	err = s.bookingRepo.Update(booking)
	if err != nil {
		return err
	}

	return s.carRepo.UpdateStatus(booking.CarID, model.CarStatusAvailable)
}

func (s *BookingService) CheckCarAvailability(carID uint, pickupTime, returnTime time.Time) bool {
	return s.bookingRepo.IsCarAvailable(carID, pickupTime, returnTime)
}

func (s *BookingService) GetUserBookings(userID uint, page, pageSize int) ([]model.Booking, int64, error) {
	return s.bookingRepo.FindAll(page, pageSize, userID, "", 0, nil, nil)
}