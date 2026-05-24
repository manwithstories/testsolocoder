package repository

import (
	"survey-platform/internal/model"
	"time"

	"gorm.io/gorm"
)

type DistributionRepository struct {
	db *gorm.DB
}

func NewDistributionRepository(db *gorm.DB) *DistributionRepository {
	return &DistributionRepository{db: db}
}

func (r *DistributionRepository) CreateLink(link *model.DistributionLink) error {
	return r.db.Create(link).Error
}

func (r *DistributionRepository) FindByToken(token string) (*model.DistributionLink, error) {
	var link model.DistributionLink
	err := r.db.Where("link_token = ?", token).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *DistributionRepository) ListBySurveyID(surveyID uint) ([]*model.DistributionLink, error) {
	var links []*model.DistributionLink
	err := r.db.Where("survey_id = ?", surveyID).Order("created_at DESC").Find(&links).Error
	return links, err
}

func (r *DistributionRepository) IncrementUseCount(id uint) error {
	return r.db.Model(&model.DistributionLink{}).Where("id = ?", id).
		UpdateColumn("use_count", gorm.Expr("use_count + 1")).Error
}

func (r *DistributionRepository) CreateInvitation(invitation *model.Invitation) error {
	return r.db.Create(invitation).Error
}

func (r *DistributionRepository) BatchCreateInvitations(invitations []*model.Invitation) error {
	return r.db.Create(&invitations).Error
}

func (r *DistributionRepository) FindInvitationByToken(token string) (*model.Invitation, error) {
	var invitation model.Invitation
	err := r.db.Where("link_token = ?", token).First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *DistributionRepository) UpdateInvitationStatus(id uint, status int, errMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == 2 {
		now := time.Now()
		updates["sent_at"] = now
	}
	if errMsg != "" {
		updates["error_message"] = errMsg
		updates["retry_count"] = gorm.Expr("retry_count + 1")
	}
	return r.db.Model(&model.Invitation{}).Where("id = ?", id).Updates(updates).Error
}

func (r *DistributionRepository) ListInvitations(surveyID uint, page, pageSize int) ([]*model.Invitation, int64, error) {
	var invitations []*model.Invitation
	var total int64

	query := r.db.Model(&model.Invitation{}).Where("survey_id = ?", surveyID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&invitations).Error; err != nil {
		return nil, 0, err
	}

	return invitations, total, nil
}

func (r *DistributionRepository) DeleteLink(id uint) error {
	return r.db.Delete(&model.DistributionLink{}, id).Error
}
