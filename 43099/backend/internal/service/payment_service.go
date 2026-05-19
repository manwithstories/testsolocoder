package service

import (
	"errors"
	"time"
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/internal/repository"
)

type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	orderRepo   *repository.OrderRepository
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		paymentRepo: repository.NewPaymentRepository(),
		orderRepo:   repository.NewOrderRepository(),
	}
}

func (s *PaymentService) ConfirmPayment(orderID uint, req *dto.ConfirmPaymentRequest) (*model.Payment, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.Status != model.OrderStatusConfirmed {
		return nil, errors.New("only confirmed orders can be paid")
	}

	existing, _ := s.paymentRepo.GetByOrderID(orderID)
	if existing != nil && existing.Status == model.PaymentStatusSuccess {
		return nil, errors.New("payment already exists for this order")
	}

	now := time.Now()
	payment := &model.Payment{
		OrderID:       orderID,
		TransactionNo: req.TransactionNo,
		Amount:        req.Amount,
		PaymentMethod: model.PaymentMethod(req.PaymentMethod),
		Status:        model.PaymentStatusSuccess,
		PaidAt:        &now,
	}

	err = s.paymentRepo.Create(payment)
	if err != nil {
		return nil, err
	}

	order.Status = model.OrderStatusPaid
	s.orderRepo.Update(order)

	return payment, nil
}

func (s *PaymentService) GetByID(id uint) (*model.Payment, error) {
	return s.paymentRepo.GetByID(id)
}

func (s *PaymentService) List(req *dto.PaymentListRequest) ([]model.Payment, int64, error) {
	return s.paymentRepo.List(req)
}

func (s *PaymentService) ExportForDateRange(startDate, endDate string) ([]model.Payment, error) {
	return s.paymentRepo.GetByDateRange(startDate, endDate)
}
