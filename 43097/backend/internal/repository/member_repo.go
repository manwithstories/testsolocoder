package repository

import (
	"hotel-system/internal/model"

	"gorm.io/gorm"
)

type MemberRepository interface {
	Create(member *model.Member) error
	GetByID(id uint) (*model.Member, error)
	GetByMemberNo(memberNo string) (*model.Member, error)
	GetByPhone(phone string) (*model.Member, error)
	Update(member *model.Member) error
	Delete(id uint) error
	List(page, pageSize int, name, phone string, levelID uint, status model.MemberStatus) ([]model.Member, int64, error)
	UpdatePoints(memberID uint, pointsChange int, description string, logType model.PointsLogType, orderNo string) error
	GetTodayMemberCount(dateStr string) (int64, error)
}

type MemberLevelRepository interface {
	Create(level *model.MemberLevel) error
	GetByID(id uint) (*model.MemberLevel, error)
	Update(level *model.MemberLevel) error
	Delete(id uint) error
	List() ([]model.MemberLevel, int64, error)
	GetLevelByPoints(points int) (*model.MemberLevel, error)
}

type memberRepository struct {
	db *gorm.DB
}

type memberLevelRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{db: db}
}

func NewMemberLevelRepository(db *gorm.DB) MemberLevelRepository {
	return &memberLevelRepository{db: db}
}

func (r *memberRepository) Create(member *model.Member) error {
	return r.db.Create(member).Error
}

func (r *memberRepository) GetByID(id uint) (*model.Member, error) {
	var member model.Member
	err := r.db.Preload("Level").First(&member, id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) GetByMemberNo(memberNo string) (*model.Member, error) {
	var member model.Member
	err := r.db.Preload("Level").Where("member_no = ?", memberNo).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) GetByPhone(phone string) (*model.Member, error) {
	var member model.Member
	err := r.db.Preload("Level").Where("phone = ?", phone).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) Update(member *model.Member) error {
	return r.db.Save(member).Error
}

func (r *memberRepository) Delete(id uint) error {
	return r.db.Delete(&model.Member{}, id).Error
}

func (r *memberRepository) List(page, pageSize int, name, phone string, levelID uint, status model.MemberStatus) ([]model.Member, int64, error) {
	var members []model.Member
	var total int64

	query := r.db.Model(&model.Member{}).Preload("Level")

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if levelID > 0 {
		query = query.Where("level_id = ?", levelID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&members).Error
	if err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

func (r *memberRepository) UpdatePoints(memberID uint, pointsChange int, description string, logType model.PointsLogType, orderNo string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var member model.Member
		if err := tx.First(&member, memberID).Error; err != nil {
			return err
		}

		newPoints := member.Points + pointsChange
		if newPoints < 0 {
			newPoints = 0
		}

		if err := tx.Model(&member).Update("points", newPoints).Error; err != nil {
			return err
		}

		log := &model.MemberPointsLog{
			MemberID:    memberID,
			Points:      pointsChange,
			Balance:     newPoints,
			Type:        logType,
			Description: description,
			OrderNo:     orderNo,
		}
		if err := tx.Create(log).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *memberRepository) GetTodayMemberCount(dateStr string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Member{}).Where("DATE(created_at) = ?", dateStr).Count(&count).Error
	return count, err
}

func (r *memberLevelRepository) Create(level *model.MemberLevel) error {
	return r.db.Create(level).Error
}

func (r *memberLevelRepository) GetByID(id uint) (*model.MemberLevel, error) {
	var level model.MemberLevel
	err := r.db.First(&level, id).Error
	if err != nil {
		return nil, err
	}
	return &level, nil
}

func (r *memberLevelRepository) Update(level *model.MemberLevel) error {
	return r.db.Save(level).Error
}

func (r *memberLevelRepository) Delete(id uint) error {
	return r.db.Delete(&model.MemberLevel{}, id).Error
}

func (r *memberLevelRepository) List() ([]model.MemberLevel, int64, error) {
	var levels []model.MemberLevel
	var total int64

	err := r.db.Model(&model.MemberLevel{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Order("min_points ASC").Find(&levels).Error
	if err != nil {
		return nil, 0, err
	}

	return levels, total, nil
}

func (r *memberLevelRepository) GetLevelByPoints(points int) (*model.MemberLevel, error) {
	var level model.MemberLevel
	err := r.db.Where("min_points <= ? AND (max_points >= ? OR max_points = 0)", points, points).
		Order("min_points DESC").
		First(&level).Error
	if err != nil {
		return nil, err
	}
	return &level, nil
}
