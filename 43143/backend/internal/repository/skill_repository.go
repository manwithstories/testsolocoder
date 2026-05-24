package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"skillshare/internal/models"
)

type SkillRepository struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) *SkillRepository {
	return &SkillRepository{db: db}
}

func (r *SkillRepository) CreateCategory(category *models.SkillCategory) error {
	return r.db.Create(category).Error
}

func (r *SkillRepository) GetCategories() ([]*models.SkillCategory, error) {
	var categories []*models.SkillCategory
	err := r.db.Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

func (r *SkillRepository) UpdateCategory(category *models.SkillCategory) error {
	return r.db.Save(category).Error
}

func (r *SkillRepository) DeleteCategory(id uuid.UUID) error {
	return r.db.Delete(&models.SkillCategory{}, id).Error
}

func (r *SkillRepository) CreateTag(tag *models.SkillTag) error {
	return r.db.Create(tag).Error
}

func (r *SkillRepository) GetTags(categoryID *uuid.UUID) ([]*models.SkillTag, error) {
	var tags []*models.SkillTag
	query := r.db
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}
	err := query.Find(&tags).Error
	return tags, err
}

func (r *SkillRepository) FindTagByID(id uuid.UUID) (*models.SkillTag, error) {
	var tag models.SkillTag
	err := r.db.First(&tag, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *SkillRepository) UpdateTag(tag *models.SkillTag) error {
	return r.db.Save(tag).Error
}

func (r *SkillRepository) DeleteTag(id uuid.UUID) error {
	return r.db.Delete(&models.SkillTag{}, id).Error
}

func (r *SkillRepository) CreateSkill(skill *models.Skill) error {
	return r.db.Create(skill).Error
}

func (r *SkillRepository) GetSkills(page, pageSize int, categoryID *uuid.UUID, keyword string) ([]*models.Skill, int64, error) {
	var skills []*models.Skill
	var total int64

	query := r.db.Model(&models.Skill{}).Where("is_active = ?", true)

	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	query.Preload("Category").Preload("Tags").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&skills)

	return skills, total, nil
}

func (r *SkillRepository) FindSkillByID(id uuid.UUID) (*models.Skill, error) {
	var skill models.Skill
	err := r.db.Preload("Category").Preload("Tags").First(&skill, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *SkillRepository) UpdateSkill(skill *models.Skill) error {
	return r.db.Save(skill).Error
}

func (r *SkillRepository) DeleteSkill(id uuid.UUID) error {
	return r.db.Model(&models.Skill{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *SkillRepository) GetPopularSkills(limit int) ([]*models.Skill, error) {
	var skills []*models.Skill
	err := r.db.Where("is_active = ?", true).
		Order("review_count DESC, rating DESC").
		Limit(limit).
		Preload("Category").
		Find(&skills).Error
	return skills, err
}

func (r *SkillRepository) CreatePosting(posting *models.SkillPosting) error {
	return r.db.Create(posting).Error
}

func (r *SkillRepository) GetPostings(page, pageSize int, skillID *uuid.UUID, teacherID *uuid.UUID) ([]*models.SkillPosting, int64, error) {
	var postings []*models.SkillPosting
	var total int64

	query := r.db.Model(&models.SkillPosting{}).Where("is_active = ?", true)

	if skillID != nil {
		query = query.Where("skill_id = ?", *skillID)
	}

	if teacherID != nil {
		query = query.Where("teacher_id = ?", *teacherID)
	}

	query.Count(&total)
	query.Preload("Skill").Preload("Teacher").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&postings)

	return postings, total, nil
}

func (r *SkillRepository) FindPostingByID(id uuid.UUID) (*models.SkillPosting, error) {
	var posting models.SkillPosting
	err := r.db.Preload("Skill").Preload("Teacher").First(&posting, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &posting, nil
}

func (r *SkillRepository) UpdatePosting(posting *models.SkillPosting) error {
	return r.db.Save(posting).Error
}

func (r *SkillRepository) DeletePosting(id uuid.UUID) error {
	return r.db.Model(&models.SkillPosting{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *SkillRepository) IncrementPostingStats(id uuid.UUID, hours float64) error {
	return r.db.Model(&models.SkillPosting{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"booking_count": gorm.Expr("booking_count + 1"),
			"total_hours":   gorm.Expr("total_hours + ?", hours),
		}).Error
}

func (r *SkillRepository) UpdatePostingRating(id uuid.UUID) error {
	return r.db.Model(&models.SkillPosting{}).Where("id = ?", id).
		Update("rating", r.db.Model(&models.Review{}).
			Where("posting_id = ?", id).
			Select("COALESCE(AVG(rating), 0)")).Error
}
