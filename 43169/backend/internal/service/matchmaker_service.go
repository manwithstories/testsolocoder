package service

import (
	"errors"
	"matchmaking-platform/internal/dto"
	"matchmaking-platform/internal/model"
	"matchmaking-platform/internal/repository"
	"matchmaking-platform/internal/utils"
	"time"

	"gorm.io/gorm"
)

type MatchmakerService struct {
	matchmakerRepo *repository.MatchmakerRepo
	userRepo       *repository.UserRepo
	dateRepo       *repository.DateRepo
	logRepo        *repository.SystemLogRepo
}

func NewMatchmakerService() *MatchmakerService {
	return &MatchmakerService{
		matchmakerRepo: repository.NewMatchmakerRepo(),
		userRepo:       repository.NewUserRepo(),
		dateRepo:       repository.NewDateRepo(),
		logRepo:        repository.NewSystemLogRepo(),
	}
}

func (s *MatchmakerService) AddMember(matchmakerID uint, req *dto.MemberManageRequest) error {
	isMember, _ := s.matchmakerRepo.IsMember(matchmakerID, req.MemberID)
	if isMember {
		return errors.New("该用户已是您的会员")
	}

	user, err := s.userRepo.FindByID(req.MemberID)
	if err != nil {
		return errors.New("用户不存在")
	}
	if user.Role != model.RoleUser {
		return errors.New("只能添加普通用户为会员")
	}

	member := &model.MatchmakerMember{
		MatchmakerID: matchmakerID,
		MemberID:     req.MemberID,
		Status:       "active",
		JoinedAt:     time.Now(),
	}

	return utils.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(member).Error; err != nil {
			return err
		}
		stats, _ := s.matchmakerRepo.GetOrCreateStats(matchmakerID)
		return tx.Model(&model.MatchmakerStats{}).
			Where("matchmaker_id = ?", matchmakerID).
			Update("total_members", stats.TotalMembers+1).Error
	})
}

func (s *MatchmakerService) RemoveMember(matchmakerID uint, memberID uint) error {
	isMember, _ := s.matchmakerRepo.IsMember(matchmakerID, memberID)
	if !isMember {
		return errors.New("该用户不是您的会员")
	}

	return utils.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("matchmaker_id = ? AND member_id = ?", matchmakerID, memberID).
			Delete(&model.MatchmakerMember{}).Error; err != nil {
			return err
		}
		stats, _ := s.matchmakerRepo.GetOrCreateStats(matchmakerID)
		return tx.Model(&model.MatchmakerStats{}).
			Where("matchmaker_id = ?", matchmakerID).
			Update("total_members", utils.Max(stats.TotalMembers-1, 0)).Error
	})
}

func (s *MatchmakerService) ListMembers(matchmakerID uint, page, pageSize int) ([]model.MatchmakerMember, int64, error) {
	return s.matchmakerRepo.ListMembers(matchmakerID, page, pageSize)
}

func (s *MatchmakerService) CreateService(matchmakerID uint, req *dto.MatchmakerServiceRequest) error {
	isMemberA, _ := s.matchmakerRepo.IsMember(matchmakerID, req.MemberAID)
	isMemberB, _ := s.matchmakerRepo.IsMember(matchmakerID, req.MemberBID)
	if !isMemberA || !isMemberB {
		return errors.New("双方必须都是您的会员")
	}

	service := &model.MatchmakerService{
		MatchmakerID: matchmakerID,
		MemberAID:   req.MemberAID,
		MemberBID:   req.MemberBID,
		ServiceType: req.ServiceType,
		Note:        req.Note,
		Status:      "progress",
		DateID:      req.DateID,
	}

	return utils.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(service).Error; err != nil {
			return err
		}
		stats, _ := s.matchmakerRepo.GetOrCreateStats(matchmakerID)
		return tx.Model(&model.MatchmakerStats{}).
			Where("matchmaker_id = ?", matchmakerID).
			Update("total_services", stats.TotalServices+1).Error
	})
}

func (s *MatchmakerService) UpdateServiceProgress(matchmakerID, serviceID uint, progress int) error {
	_, _, err := s.matchmakerRepo.ListServices(matchmakerID, 1, 100)
	if err != nil {
		return err
	}

	return s.matchmakerRepo.UpdateService(serviceID, map[string]interface{}{
		"progress": progress,
		"status":   "completed",
	})
}

func (s *MatchmakerService) ListServices(matchmakerID uint, page, pageSize int) ([]model.MatchmakerService, int64, error) {
	return s.matchmakerRepo.ListServices(matchmakerID, page, pageSize)
}

func (s *MatchmakerService) GetStats(matchmakerID uint) (*model.MatchmakerStats, error) {
	stats, err := s.matchmakerRepo.GetOrCreateStats(matchmakerID)
	if err != nil {
		return nil, err
	}

	dates, _, _ := s.dateRepo.ListByUser(matchmakerID, 1, 10000)
	totalDates := len(dates)
	successDates := 0
	for _, d := range dates {
		if d.Status == model.DateStatusCompleted {
			successDates++
		}
	}

	s.matchmakerRepo.UpdateStats(matchmakerID, map[string]interface{}{
		"total_dates":   totalDates,
		"success_dates": successDates,
	})

	stats.TotalDates = totalDates
	stats.SuccessDates = successDates

	return stats, nil
}

func (s *MatchmakerService) ListAllMatchmakers(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	db := utils.DB.Model(&model.User{}).Where("role = ?", model.RoleMatchmaker)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&users).Error
	return users, total, err
}
