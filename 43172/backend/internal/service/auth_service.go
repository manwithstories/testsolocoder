package service

import (
	"context"
	"errors"
	"fmt"

	"luxury-trading-platform/internal/cache"
	"luxury-trading-platform/internal/model"
	"luxury-trading-platform/internal/repository"
	"luxury-trading-platform/internal/utils"

	"gorm.io/gorm"
)

type AuthenticationService struct {
	authRepo    *repository.AuthenticationRepository
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
	userRepo    *repository.UserRepository
	redisClient *cache.RedisClient
	db          *gorm.DB
}

func NewAuthenticationService(authRepo *repository.AuthenticationRepository, orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository, userRepo *repository.UserRepository, redisClient *cache.RedisClient, db *gorm.DB) *AuthenticationService {
	return &AuthenticationService{
		authRepo:    authRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
		redisClient: redisClient,
		db:          db,
	}
}

type CompleteAuthenticationRequest struct {
	Result         model.AuthenticationResult `json:"result" binding:"required,oneof=genuine counterfeit inconclusive"`
	ReportFile     string                     `json:"report_file"`
	ReportContent  string                     `json:"report_content"`
	AuthenticatorNotes string                `json:"authenticator_notes"`
}

func (s *AuthenticationService) CreateAuthentication(ctx context.Context, buyerID uint, orderID uint) (*model.Authentication, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order: %w", err)
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	if order.BuyerID != buyerID {
		return nil, errors.New("permission denied: this order does not belong to you")
	}

	if order.Status != model.OrderStatusPaid && order.Status != model.OrderStatusShipped {
		return nil, errors.New("authentication can only be requested for paid/shipped orders")
	}

	existing, err := s.authRepo.FindByOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing authentication: %w", err)
	}
	if existing != nil {
		return nil, errors.New("authentication already exists for this order")
	}

	auth := &model.Authentication{
		OrderID:   orderID,
		ProductID: order.ProductID,
		BuyerID:   buyerID,
		Status:    model.AuthenticationStatusPending,
	}

	err = s.authRepo.Create(auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create authentication: %w", err)
	}

	return auth, nil
}

func (s *AuthenticationService) GetAuthentication(ctx context.Context, id uint) (*model.Authentication, error) {
	auth, err := s.authRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find authentication: %w", err)
	}
	if auth == nil {
		return nil, errors.New("authentication not found")
	}
	return auth, nil
}

func (s *AuthenticationService) GetAuthenticationByOrder(ctx context.Context, orderID uint) (*model.Authentication, error) {
	auth, err := s.authRepo.FindByOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find authentication: %w", err)
	}
	if auth == nil {
		return nil, errors.New("authentication not found")
	}
	return auth, nil
}

func (s *AuthenticationService) ListAuthentications(page, pageSize int, status model.AuthenticationStatus, authenticatorID *uint, buyerID *uint) ([]model.Authentication, int64, error) {
	page, pageSize = utils.ValidatePage(page, pageSize)
	return s.authRepo.List(page, pageSize, status, authenticatorID, buyerID)
}

func (s *AuthenticationService) AcceptAuthentication(ctx context.Context, authID uint, authenticatorID uint) (*model.Authentication, error) {
	auth, err := s.authRepo.FindByID(authID)
	if err != nil {
		return nil, fmt.Errorf("failed to find authentication: %w", err)
	}
	if auth == nil {
		return nil, errors.New("authentication not found")
	}

	if auth.Status != model.AuthenticationStatusPending {
		return nil, errors.New("authentication is not in pending status")
	}

	profile, err := s.userRepo.FindAuthenticatorProfileByUserID(authenticatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find authenticator profile: %w", err)
	}
	if profile == nil || profile.Status != model.AuthenticatorStatusApproved {
		return nil, errors.New("authenticator is not approved")
	}

	err = s.authRepo.Accept(authID, authenticatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to accept authentication: %w", err)
	}

	return s.authRepo.FindByID(authID)
}

func (s *AuthenticationService) CompleteAuthentication(ctx context.Context, authID uint, authenticatorID uint, req *CompleteAuthenticationRequest) (*model.Authentication, error) {
	auth, err := s.authRepo.FindByID(authID)
	if err != nil {
		return nil, fmt.Errorf("failed to find authentication: %w", err)
	}
	if auth == nil {
		return nil, errors.New("authentication not found")
	}

	if auth.AuthenticatorID == nil || *auth.AuthenticatorID != authenticatorID {
		return nil, errors.New("permission denied: you are not assigned to this authentication")
	}

	if auth.Status != model.AuthenticationStatusAccepted {
		return nil, errors.New("authentication is not in accepted status")
	}

	err = s.authRepo.Complete(authID, req.Result, req.ReportFile, req.ReportContent, req.AuthenticatorNotes)
	if err != nil {
		return nil, fmt.Errorf("failed to complete authentication: %w", err)
	}

	if req.Result == model.AuthenticationResultGenuine {
		_ = s.productRepo.SetAuthenticated(auth.ProductID, true)
	}

	profile, _ := s.userRepo.FindAuthenticatorProfileByUserID(authenticatorID)
	if profile != nil {
		profile.CompletedCount++
		_ = s.userRepo.UpdateAuthenticatorProfile(profile)
	}

	return s.authRepo.FindByID(authID)
}

func (s *AuthenticationService) RejectAuthentication(ctx context.Context, authID uint, authenticatorID uint, reason string) (*model.Authentication, error) {
	auth, err := s.authRepo.FindByID(authID)
	if err != nil {
		return nil, fmt.Errorf("failed to find authentication: %w", err)
	}
	if auth == nil {
		return nil, errors.New("authentication not found")
	}

	if auth.AuthenticatorID == nil || *auth.AuthenticatorID != authenticatorID {
		return nil, errors.New("permission denied")
	}

	if auth.Status != model.AuthenticationStatusAccepted {
		return nil, errors.New("authentication is not in accepted status")
	}

	err = s.authRepo.Reject(authID, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to reject authentication: %w", err)
	}

	return s.authRepo.FindByID(authID)
}

func (s *AuthenticationService) CancelAuthentication(ctx context.Context, authID uint, buyerID uint) error {
	auth, err := s.authRepo.FindByID(authID)
	if err != nil {
		return fmt.Errorf("failed to find authentication: %w", err)
	}
	if auth == nil {
		return errors.New("authentication not found")
	}

	if auth.BuyerID != buyerID {
		return errors.New("permission denied")
	}

	if auth.Status != model.AuthenticationStatusPending {
		return errors.New("authentication cannot be cancelled")
	}

	return s.authRepo.Cancel(authID)
}
