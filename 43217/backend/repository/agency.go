package repository

import (
	"time"

	"health-platform/models"
)

type AgencyRepository struct {
	*BaseRepository
}

func NewAgencyRepository() *AgencyRepository {
	return &AgencyRepository{
		BaseRepository: NewBaseRepository(),
	}
}

func (r *AgencyRepository) FindByName(name string) (*models.Agency, error) {
	var agency models.Agency
	err := r.DB.Where("name = ?", name).First(&agency).Error
	if err != nil {
		return nil, err
	}
	return &agency, nil
}

func (r *AgencyRepository) FindByUnifiedCode(code string) (*models.Agency, error) {
	var agency models.Agency
	err := r.DB.Where("unified_code = ?", code).First(&agency).Error
	if err != nil {
		return nil, err
	}
	return &agency, nil
}

func (r *AgencyRepository) GetWithPackages(agencyID uint) (*models.Agency, error) {
	var agency models.Agency
	err := r.DB.Preload("Packages").First(&agency, agencyID).Error
	if err != nil {
		return nil, err
	}
	return &agency, nil
}

func (r *AgencyRepository) ListActiveAgencies() ([]models.Agency, error) {
	var agencies []models.Agency
	err := r.DB.Where("status = ?", models.AgencyStatusActive).Find(&agencies).Error
	return agencies, err
}
