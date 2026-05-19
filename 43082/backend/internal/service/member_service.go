package service

import (
	"errors"
	"gym-management/internal/models"
	"gym-management/internal/repository"
)

type MemberService interface {
	Register(member *models.Member) error
	Login(phone, password string) (*models.Member, error)
	GetByID(id uint) (*models.Member, error)
	List(page, pageSize int, keyword string) ([]models.Member, int64, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	UpdateStatus(id uint, status int) error
}

type memberService struct {
	memberRepo     repository.MemberRepository
	membershipRepo repository.MembershipRepository
}

func NewMemberService() MemberService {
	return &memberService{
		memberRepo:     repository.NewMemberRepository(),
		membershipRepo: repository.NewMembershipRepository(),
	}
}

func (s *memberService) Register(member *models.Member) error {
	existing, _ := s.memberRepo.GetByPhone(member.Phone)
	if existing != nil && existing.ID > 0 {
		return errors.New("手机号已存在")
	}

	if member.Email != "" {
		existing, _ = s.memberRepo.GetByEmail(member.Email)
		if existing != nil && existing.ID > 0 {
			return errors.New("邮箱已存在")
		}
	}

	if member.Password == "" {
		member.Password = "123456"
	}
	err := member.HashPassword(member.Password)
	if err != nil {
		return err
	}

	member.Status = 1
	return s.memberRepo.Create(member)
}

func (s *memberService) Login(phone, password string) (*models.Member, error) {
	member, err := s.memberRepo.GetByPhone(phone)
	if err != nil || member.ID == 0 {
		return nil, errors.New("用户不存在")
	}

	if member.Status != 1 {
		return nil, errors.New("账号已被冻结")
	}

	if !member.CheckPassword(password) {
		return nil, errors.New("密码错误")
	}

	return member, nil
}

func (s *memberService) GetByID(id uint) (*models.Member, error) {
	return s.memberRepo.GetByID(id)
}

func (s *memberService) List(page, pageSize int, keyword string) ([]models.Member, int64, error) {
	return s.memberRepo.List(page, pageSize, keyword)
}

func (s *memberService) Update(id uint, updates map[string]interface{}) error {
	_, err := s.memberRepo.GetByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	if phone, ok := updates["phone"].(string); ok {
		existing, _ := s.memberRepo.GetByPhone(phone)
		if existing != nil && existing.ID != id {
			return errors.New("手机号已存在")
		}
	}

	if email, ok := updates["email"].(string); ok {
		existing, _ := s.memberRepo.GetByEmail(email)
		if existing != nil && existing.ID != id {
			return errors.New("邮箱已存在")
		}
	}

	return s.memberRepo.DB().Model(&models.Member{}).Where("id = ?", id).Updates(updates).Error
}

func (s *memberService) Delete(id uint) error {
	_, err := s.memberRepo.GetByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}
	return s.memberRepo.Delete(id)
}

func (s *memberService) UpdateStatus(id uint, status int) error {
	_, err := s.memberRepo.GetByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}
	return s.memberRepo.UpdateStatus(id, status)
}
