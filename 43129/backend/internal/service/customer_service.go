package service

import (
	"fmt"

	"beauty-salon-system/internal/model"
	"beauty-salon-system/internal/repository"
)

type CustomerService struct {
	customerRepo *repository.CustomerRepository
	userRepo     *repository.UserRepository
}

func NewCustomerService(customerRepo *repository.CustomerRepository, userRepo *repository.UserRepository) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
		userRepo:     userRepo,
	}
}

type CreateCustomerRequest struct {
	UserID         uint   `json:"user_id" binding:"required"`
	Name           string `json:"name"`
	Gender         string `json:"gender"`
	Age            int    `json:"age"`
	SkinType       string `json:"skin_type"`
	HairPreference string `json:"hair_preference"`
	AllergyHistory string `json:"allergy_history"`
	Notes          string `json:"notes"`
}

type UpdateCustomerRequest struct {
	Name           string `json:"name"`
	Gender         string `json:"gender"`
	Age            int    `json:"age"`
	SkinType       string `json:"skin_type"`
	HairPreference string `json:"hair_preference"`
	AllergyHistory string `json:"allergy_history"`
	Notes          string `json:"notes"`
}

func (s *CustomerService) Create(req *CreateCustomerRequest) (*model.Customer, error) {
	user, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	existing, _ := s.customerRepo.GetByUserID(req.UserID)
	if existing != nil {
		return nil, fmt.Errorf("customer profile already exists")
	}

	customer := &model.Customer{
		UserID:         req.UserID,
		Name:           req.Name,
		Gender:         req.Gender,
		Age:            req.Age,
		SkinType:       req.SkinType,
		HairPreference: req.HairPreference,
		AllergyHistory: req.AllergyHistory,
		Notes:          req.Notes,
		MemberLevel:    1,
		Points:         0,
		User:           user,
	}

	if err := s.customerRepo.Create(customer); err != nil {
		return nil, fmt.Errorf("create customer: %w", err)
	}

	return customer, nil
}

func (s *CustomerService) GetByID(id uint) (*model.Customer, error) {
	return s.customerRepo.GetByID(id)
}

func (s *CustomerService) GetByUserID(userID uint) (*model.Customer, error) {
	return s.customerRepo.GetByUserID(userID)
}

func (s *CustomerService) Update(id uint, req *UpdateCustomerRequest) (*model.Customer, error) {
	customer, err := s.customerRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.Gender != "" {
		customer.Gender = req.Gender
	}
	if req.Age > 0 {
		customer.Age = req.Age
	}
	if req.SkinType != "" {
		customer.SkinType = req.SkinType
	}
	if req.HairPreference != "" {
		customer.HairPreference = req.HairPreference
	}
	if req.AllergyHistory != "" {
		customer.AllergyHistory = req.AllergyHistory
	}
	if req.Notes != "" {
		customer.Notes = req.Notes
	}

	if err := s.customerRepo.Update(customer); err != nil {
		return nil, fmt.Errorf("update customer: %w", err)
	}

	return customer, nil
}

func (s *CustomerService) List(page, pageSize int, keyword string) ([]model.Customer, int64, error) {
	return s.customerRepo.List(page, pageSize, keyword)
}

func (s *CustomerService) AddPoints(id uint, points int) error {
	return s.customerRepo.AddPoints(id, points)
}

func (s *CustomerService) DeductPoints(id uint, points int) error {
	return s.customerRepo.DeductPoints(id, points)
}

func (s *CustomerService) CalculateLevel(totalSpent float64) int {
	switch {
	case totalSpent >= 50000:
		return 5
	case totalSpent >= 20000:
		return 4
	case totalSpent >= 10000:
		return 3
	case totalSpent >= 5000:
		return 2
	default:
		return 1
	}
}

func (s *CustomerService) UpdateLevel(id uint, totalSpent float64) error {
	level := s.CalculateLevel(totalSpent)
	return s.customerRepo.UpdateLevel(id, level)
}

func (s *CustomerService) AddVisit(id uint) error {
	return s.customerRepo.AddVisit(id)
}

func (s *CustomerService) AddTotalSpent(id uint, amount float64) error {
	customer, err := s.customerRepo.GetByID(id)
	if err != nil {
		return err
	}
	newTotal := customer.TotalSpent + amount
	if err := s.customerRepo.AddTotalSpent(id, amount); err != nil {
		return err
	}
	return s.UpdateLevel(id, newTotal)
}
