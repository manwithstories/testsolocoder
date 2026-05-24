package repository

import (
	"matchmaking-platform/internal/model"
	"matchmaking-platform/internal/utils"
)

type ProfileRepo struct{}

func NewProfileRepo() *ProfileRepo {
	return &ProfileRepo{}
}

func (r *ProfileRepo) Create(profile *model.Profile) error {
	return utils.DB.Create(profile).Error
}

func (r *ProfileRepo) FindByUserID(userID uint) (*model.Profile, error) {
	var profile model.Profile
	err := utils.DB.Where("user_id = ?", userID).First(&profile).Error
	return &profile, err
}

func (r *ProfileRepo) Update(userID uint, updates map[string]interface{}) error {
	return utils.DB.Model(&model.Profile{}).Where("user_id = ?", userID).Updates(updates).Error
}

func (r *ProfileRepo) Upsert(profile *model.Profile) error {
	var existing model.Profile
	err := utils.DB.Where("user_id = ?", profile.UserID).First(&existing).Error
	if err != nil {
		return utils.DB.Create(profile).Error
	}
	return utils.DB.Model(&existing).Updates(profile).Error
}

func (r *ProfileRepo) ListByFilter(filter map[string]interface{}, excludeUserIDs []uint, page, pageSize int) ([]model.Profile, int64, error) {
	var profiles []model.Profile
	var total int64
	db := utils.DB.Model(&model.Profile{})
	for key, val := range filter {
		db = db.Where(key+" = ?", val)
	}
	if len(excludeUserIDs) > 0 {
		db = db.Where("user_id NOT IN ?", excludeUserIDs)
	}
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&profiles).Error
	return profiles, total, err
}

func (r *ProfileRepo) FindByIDs(ids []uint) ([]model.Profile, error) {
	var profiles []model.Profile
	err := utils.DB.Where("user_id IN ?", ids).Find(&profiles).Error
	return profiles, err
}
