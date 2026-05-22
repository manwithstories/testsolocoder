package service

import (
	"car-rental/internal/model"
	"car-rental/internal/repository"
	"errors"
	"time"

	"github.com/xuri/excelize/v2"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	bookingRepo *repository.BookingRepository
	carRepo     *repository.CarRepository
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:   repository.NewOrderRepository(),
		bookingRepo: repository.NewBookingRepository(),
		carRepo:     repository.NewCarRepository(),
	}
}

func (s *OrderService) GetOrderByID(id uint) (*model.Order, error) {
	return s.orderRepo.FindByID(id)
}

func (s *OrderService) GetOrderByNo(orderNo string) (*model.Order, error) {
	return s.orderRepo.FindByOrderNo(orderNo)
}

func (s *OrderService) GetAllOrders(page, pageSize int, userID uint, status string, startDate, endDate *time.Time) ([]model.Order, int64, error) {
	return s.orderRepo.FindAll(page, pageSize, userID, status, startDate, endDate)
}

func (s *OrderService) GetUserOrders(userID uint, page, pageSize int) ([]model.Order, int64, error) {
	return s.orderRepo.FindAll(page, pageSize, userID, "", nil, nil)
}

func (s *OrderService) UpdateOrderStatus(id uint, status model.OrderStatus) error {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return errors.New("订单不存在")
	}

	if status == model.OrderStatusPaid && order.PaymentStatus != model.PaymentStatusPaid {
		err = s.orderRepo.UpdatePaymentStatus(id, model.PaymentStatusPaid, order.FinalAmount)
		if err != nil {
			return err
		}

		if order.Booking != nil {
			err = s.bookingRepo.UpdateStatus(order.BookingID, model.BookingStatusConfirmed)
			if err != nil {
				return err
			}
			err = s.carRepo.UpdateStatus(order.CarID, model.CarStatusRented)
			if err != nil {
				return err
			}
		}
	}

	return s.orderRepo.UpdateStatus(id, status)
}

func (s *OrderService) RefundOrder(id uint, reason string) error {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return errors.New("订单不存在")
	}

	if order.PaymentStatus != model.PaymentStatusPaid {
		return errors.New("订单未支付，无法退款")
	}

	now := time.Now()
	order.RefundedAt = &now
	order.RefundReason = reason
	err = s.orderRepo.Update(order)
	if err != nil {
		return err
	}

	return s.orderRepo.UpdateStatus(id, model.OrderStatusRefunded)
}

func (s *OrderService) ExportOrders(status string, startDate, endDate *time.Time) (string, error) {
	orders, err := s.orderRepo.FindAllForExport(status, startDate, endDate)
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()
	sheet := "Orders"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"订单号", "用户名", "邮箱", "车型", "车牌号", "取车时间", "还车时间", "金额", "状态", "支付状态", "创建时间"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	for i, order := range orders {
		row := i + 2
		username := ""
		email := ""
		if order.User != nil {
			username = order.User.Username
			email = order.User.Email
		}

		carModel := ""
		licensePlate := ""
		if order.Car != nil {
			carModel = order.Car.Brand + " " + order.Car.Model
			licensePlate = order.Car.LicensePlate
		}

		pickupTime := ""
		returnTime := ""
		if order.Booking != nil {
			pickupTime = order.Booking.PickupTime.Format("2006-01-02 15:04")
			returnTime = order.Booking.ReturnTime.Format("2006-01-02 15:04")
		}

		f.SetCellValue(sheet, "A"+itoa(row), order.OrderNo)
		f.SetCellValue(sheet, "B"+itoa(row), username)
		f.SetCellValue(sheet, "C"+itoa(row), email)
		f.SetCellValue(sheet, "D"+itoa(row), carModel)
		f.SetCellValue(sheet, "E"+itoa(row), licensePlate)
		f.SetCellValue(sheet, "F"+itoa(row), pickupTime)
		f.SetCellValue(sheet, "G"+itoa(row), returnTime)
		f.SetCellValue(sheet, "H"+itoa(row), order.FinalAmount)
		f.SetCellValue(sheet, "I"+itoa(row), order.Status)
		f.SetCellValue(sheet, "J"+itoa(row), order.PaymentStatus)
		f.SetCellValue(sheet, "K"+itoa(row), order.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	fileName := "orders_" + time.Now().Format("20060102150405") + ".xlsx"
	filePath := "./uploads/exports/" + fileName

	err = f.SaveAs(filePath)
	if err != nil {
		return "", err
	}

	return "/uploads/exports/" + fileName, nil
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}