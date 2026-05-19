package service

import (
	"errors"
	"gym-management/internal/models"
	"gym-management/internal/repository"
	"time"
)

type MembershipService interface {
	Create(membership *models.Membership) error
	GetByID(id uint) (*models.Membership, error)
	GetByMemberID(memberID uint) (*models.Membership, error)
	List(page, pageSize int, memberID uint) ([]models.Membership, int64, error)
	Renew(id uint, newType models.MembershipType) (*models.Membership, error)
	Upgrade(memberID uint, newType models.MembershipType, price float64) (*models.Membership, error)
	UpdateStatus(id uint, status int) error
	CheckValidity(memberID uint) (bool, error)
}

type membershipService struct {
	membershipRepo repository.MembershipRepository
	memberRepo     repository.MemberRepository
}

func NewMembershipService() MembershipService {
	return &membershipService{
		membershipRepo: repository.NewMembershipRepository(),
		memberRepo:     repository.NewMemberRepository(),
	}
}

func (s *membershipService) Create(membership *models.Membership) error {
	days, valid := models.ValidMembershipTypes[membership.Type]
	if !valid {
		return errors.New("无效的会员卡类型")
	}

	if membership.StartDate.IsZero() {
		membership.StartDate = time.Now()
	}
	membership.EndDate = membership.StartDate.AddDate(0, 0, days)
	membership.Status = 1

	_, err := s.memberRepo.GetByID(membership.MemberID)
	if err != nil {
		return errors.New("会员不存在")
	}

	return s.membershipRepo.Create(membership)
}

func (s *membershipService) GetByID(id uint) (*models.Membership, error) {
	return s.membershipRepo.GetByID(id)
}

func (s *membershipService) GetByMemberID(memberID uint) (*models.Membership, error) {
	return s.membershipRepo.GetByMemberID(memberID)
}

func (s *membershipService) List(page, pageSize int, memberID uint) ([]models.Membership, int64, error) {
	return s.membershipRepo.List(page, pageSize, memberID)
}

func (s *membershipService) Renew(id uint, newType models.MembershipType) (*models.Membership, error) {
	oldMembership, err := s.membershipRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("会员卡不存在")
	}

	days, valid := models.ValidMembershipTypes[newType]
	if !valid {
		return nil, errors.New("无效的会员卡类型")
	}

	newMembership := &models.Membership{
		MemberID: oldMembership.MemberID,
		Type:     newType,
	}

	if oldMembership.IsValid() {
		newMembership.StartDate = oldMembership.EndDate
	} else {
		newMembership.StartDate = time.Now()
	}
	newMembership.EndDate = newMembership.StartDate.AddDate(0, 0, days)
	newMembership.Status = 1
	newMembership.AutoRenew = oldMembership.AutoRenew

	price := calculateMembershipPrice(newType)
	newMembership.Price = price

	if oldMembership.IsValid() {
		oldMembership.Status = 1
		_ = s.membershipRepo.Update(oldMembership)
	}

	err = s.membershipRepo.Create(newMembership)
	if err != nil {
		return nil, err
	}

	oldMembership.Status = 2
	_ = s.membershipRepo.Update(oldMembership)

	member, _ := s.memberRepo.GetByID(oldMembership.MemberID)
	if member != nil {
		member.MembershipID = &newMembership.ID
		_ = s.memberRepo.Update(member)
	}

	return newMembership, nil
}

func (s *membershipService) Upgrade(memberID uint, newType models.MembershipType, price float64) (*models.Membership, error) {
	days, valid := models.ValidMembershipTypes[newType]
	if !valid {
		return nil, errors.New("无效的会员卡类型")
	}

	oldMembership, _ := s.membershipRepo.GetByMemberID(memberID)

	newMembership := &models.Membership{
		MemberID: memberID,
		Type:     newType,
		Price:    price,
	}

	if oldMembership != nil && oldMembership.IsValid() {
		newMembership.StartDate = oldMembership.EndDate
	} else {
		newMembership.StartDate = time.Now()
	}
	newMembership.EndDate = newMembership.StartDate.AddDate(0, 0, days)
	newMembership.Status = 1

	err := s.membershipRepo.Create(newMembership)
	if err != nil {
		return nil, err
	}

	if oldMembership != nil {
		oldMembership.Status = 2
		_ = s.membershipRepo.Update(oldMembership)
	}

	member, _ := s.memberRepo.GetByID(memberID)
	if member != nil {
		member.MembershipID = &newMembership.ID
		_ = s.memberRepo.Update(member)
	}

	return newMembership, nil
}

func (s *membershipService) UpdateStatus(id uint, status int) error {
	membership, err := s.membershipRepo.GetByID(id)
	if err != nil {
		return errors.New("会员卡不存在")
	}
	membership.Status = status
	return s.membershipRepo.Update(membership)
}

func (s *membershipService) CheckValidity(memberID uint) (bool, error) {
	membership, err := s.membershipRepo.GetByMemberID(memberID)
	if err != nil {
		return false, errors.New("会员没有有效会员卡")
	}
	return membership.IsValid(), nil
}

func calculateMembershipPrice(membershipType models.MembershipType) float64 {
	switch membershipType {
	case models.MembershipTypeMonthly:
		return 299.0
	case models.MembershipTypeQuarter:
		return 799.0
	case models.MembershipTypeYearly:
		return 2888.0
	default:
		return 0
	}
}

func (s *membershipService) GetExpiringSoon(days int) ([]models.Membership, error) {
	return s.membershipRepo.GetExpiringSoon(days)
}
