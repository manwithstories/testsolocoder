package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"skillshare/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Preload("SkillTags").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) UpdateLastLogin(id uuid.UUID, ip string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at": gorm.Expr("NOW()"),
			"last_login_ip": ip,
		}).Error
}

func (r *UserRepository) List(page, pageSize int, keyword string) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	query := r.db.Model(&models.User{})

	if keyword != "" {
		query = query.Where("nickname LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)

	return users, total, nil
}

func (r *UserRepository) UpdateRating(id uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).
		Update("rating", r.db.Model(&models.Review{}).
			Where("reviewee_id = ?", id).
			Select("COALESCE(AVG(rating), 0)"),
		).Error
}

func (r *UserRepository) AddSkillTags(userID uuid.UUID, tagIDs []uuid.UUID) error {
	var tags []models.SkillTag
	for _, id := range tagIDs {
		tags = append(tags, models.SkillTag{ID: id})
	}
	return r.db.Model(&models.User{ID: userID}).Association("SkillTags").Append(tags)
}

func (r *UserRepository) RemoveSkillTag(userID, tagID uuid.UUID) error {
	return r.db.Model(&models.User{ID: userID}).Association("SkillTags").
		Delete(&models.SkillTag{ID: tagID})
}
