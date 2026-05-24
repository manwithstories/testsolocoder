package service

import (
	"errors"

	"github.com/google/uuid"
	"skillshare/internal/models"
	"skillshare/internal/repository"
	"skillshare/pkg/validator"
)

type SkillService struct {
	skillRepo *repository.SkillRepository
}

func NewSkillService(skillRepo *repository.SkillRepository) *SkillService {
	return &SkillService{skillRepo: skillRepo}
}

func (s *SkillService) CreateCategory(name, icon string, sortOrder int) (*models.SkillCategory, error) {
	category := &models.SkillCategory{
		Name:      name,
		Icon:      icon,
		SortOrder: sortOrder,
	}

	if err := s.skillRepo.CreateCategory(category); err != nil {
		return nil, errors.New("创建分类失败")
	}
	return category, nil
}

func (s *SkillService) GetCategories() ([]*models.SkillCategory, error) {
	return s.skillRepo.GetCategories()
}

func (s *SkillService) UpdateCategory(id uuid.UUID, name, icon string, sortOrder int) (*models.SkillCategory, error) {
	category := &models.SkillCategory{
		ID:        id,
		Name:      name,
		Icon:      icon,
		SortOrder: sortOrder,
	}

	if err := s.skillRepo.UpdateCategory(category); err != nil {
		return nil, errors.New("更新分类失败")
	}
	return category, nil
}

func (s *SkillService) DeleteCategory(id uuid.UUID) error {
	return s.skillRepo.DeleteCategory(id)
}

func (s *SkillService) CreateTag(name string, categoryID *uuid.UUID) (*models.SkillTag, error) {
	tag := &models.SkillTag{
		Name:       name,
		CategoryID: *categoryID,
	}

	if err := s.skillRepo.CreateTag(tag); err != nil {
		return nil, errors.New("创建标签失败")
	}
	return tag, nil
}

func (s *SkillService) GetTags(categoryID *uuid.UUID) ([]*models.SkillTag, error) {
	return s.skillRepo.GetTags(categoryID)
}

func (s *SkillService) UpdateTag(id uuid.UUID, name string, categoryID *uuid.UUID) (*models.SkillTag, error) {
	tag := &models.SkillTag{
		ID:         id,
		Name:       name,
		CategoryID: *categoryID,
	}

	if err := s.skillRepo.UpdateTag(tag); err != nil {
		return nil, errors.New("更新标签失败")
	}
	return tag, nil
}

func (s *SkillService) DeleteTag(id uuid.UUID) error {
	return s.skillRepo.DeleteTag(id)
}

type CreateSkillInput struct {
	Title         string              `json:"title"`
	Description   string              `json:"description"`
	CategoryID    uuid.UUID           `json:"category_id"`
	TagIDs        []uuid.UUID         `json:"tag_ids"`
	Difficulty    models.DifficultyLevel `json:"difficulty"`
	CoverImage    string              `json:"cover_image"`
	VideoURL      string              `json:"video_url"`
	Prerequisites string              `json:"prerequisites"`
	Outcomes      string              `json:"outcomes"`
}

func (s *SkillService) CreateSkill(input *CreateSkillInput) (*models.Skill, error) {
	skill := &models.Skill{
		Title:         input.Title,
		Description:   input.Description,
		CategoryID:    input.CategoryID,
		Difficulty:    input.Difficulty,
		CoverImage:    input.CoverImage,
		VideoURL:      input.VideoURL,
		Prerequisites: input.Prerequisites,
		Outcomes:      input.Outcomes,
	}

	for _, tagID := range input.TagIDs {
		tag, err := s.skillRepo.FindTagByID(tagID)
		if err == nil {
			skill.Tags = append(skill.Tags, *tag)
		}
	}

	if err := s.skillRepo.CreateSkill(skill); err != nil {
		return nil, errors.New("创建技能失败")
	}
	return skill, nil
}

func (s *SkillService) GetSkills(page, pageSize int, categoryID *uuid.UUID, keyword string) ([]*models.Skill, int64, error) {
	return s.skillRepo.GetSkills(page, pageSize, categoryID, keyword)
}

func (s *SkillService) GetSkill(id uuid.UUID) (*models.Skill, error) {
	skill, err := s.skillRepo.FindSkillByID(id)
	if err != nil {
		return nil, errors.New("技能不存在")
	}
	return skill, nil
}

func (s *SkillService) UpdateSkill(id uuid.UUID, input *CreateSkillInput) (*models.Skill, error) {
	skill, err := s.skillRepo.FindSkillByID(id)
	if err != nil {
		return nil, errors.New("技能不存在")
	}

	skill.Title = input.Title
	skill.Description = input.Description
	skill.CategoryID = input.CategoryID
	skill.Difficulty = input.Difficulty
	skill.CoverImage = input.CoverImage
	skill.VideoURL = input.VideoURL
	skill.Prerequisites = input.Prerequisites
	skill.Outcomes = input.Outcomes

	skill.Tags = skill.Tags[:0]
	for _, tagID := range input.TagIDs {
		tag, err := s.skillRepo.FindTagByID(tagID)
		if err == nil {
			skill.Tags = append(skill.Tags, *tag)
		}
	}

	if err := s.skillRepo.UpdateSkill(skill); err != nil {
		return nil, errors.New("更新技能失败")
	}
	return skill, nil
}

func (s *SkillService) DeleteSkill(id uuid.UUID) error {
	return s.skillRepo.DeleteSkill(id)
}

func (s *SkillService) GetPopularSkills(limit int) ([]*models.Skill, error) {
	return s.skillRepo.GetPopularSkills(limit)
}

type CreatePostingInput struct {
	SkillID         uuid.UUID              `json:"skill_id"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	TeachingMethod  models.TeachingMethod  `json:"teaching_method"`
	TeachingMode    models.TeachingMode    `json:"teaching_mode"`
	MaxStudents     int                    `json:"max_students"`
	PricePerHour    float64                `json:"price_per_hour"`
	SessionDuration int                    `json:"session_duration"`
	Location        string                 `json:"location"`
	Latitude        float64                `json:"latitude"`
	Longitude       float64                `json:"longitude"`
	Availability    string                 `json:"availability"`
}

func (s *SkillService) CreatePosting(teacherID uuid.UUID, input *CreatePostingInput) (*models.SkillPosting, error) {
	if err := validator.ValidatePrice(input.PricePerHour); err != nil {
		return nil, err
	}

	posting := &models.SkillPosting{
		TeacherID:       teacherID,
		SkillID:         input.SkillID,
		Title:           input.Title,
		Description:     input.Description,
		TeachingMethod:  input.TeachingMethod,
		TeachingMode:    input.TeachingMode,
		MaxStudents:     input.MaxStudents,
		PricePerHour:    input.PricePerHour,
		SessionDuration: input.SessionDuration,
		Location:        input.Location,
		Latitude:        input.Latitude,
		Longitude:       input.Longitude,
		Availability:    input.Availability,
	}

	if err := s.skillRepo.CreatePosting(posting); err != nil {
		return nil, errors.New("发布技能失败")
	}
	return posting, nil
}

func (s *SkillService) GetPostings(page, pageSize int, skillID, teacherID *uuid.UUID) ([]*models.SkillPosting, int64, error) {
	return s.skillRepo.GetPostings(page, pageSize, skillID, teacherID)
}

func (s *SkillService) GetPosting(id uuid.UUID) (*models.SkillPosting, error) {
	posting, err := s.skillRepo.FindPostingByID(id)
	if err != nil {
		return nil, errors.New("课程不存在")
	}
	return posting, nil
}

func (s *SkillService) UpdatePosting(id uuid.UUID, input *CreatePostingInput) (*models.SkillPosting, error) {
	posting, err := s.skillRepo.FindPostingByID(id)
	if err != nil {
		return nil, errors.New("课程不存在")
	}

	if err := validator.ValidatePrice(input.PricePerHour); err != nil {
		return nil, err
	}

	posting.Title = input.Title
	posting.Description = input.Description
	posting.TeachingMethod = input.TeachingMethod
	posting.TeachingMode = input.TeachingMode
	posting.MaxStudents = input.MaxStudents
	posting.PricePerHour = input.PricePerHour
	posting.SessionDuration = input.SessionDuration
	posting.Location = input.Location
	posting.Latitude = input.Latitude
	posting.Longitude = input.Longitude
	posting.Availability = input.Availability

	if err := s.skillRepo.UpdatePosting(posting); err != nil {
		return nil, errors.New("更新课程失败")
	}
	return posting, nil
}

func (s *SkillService) DeletePosting(id uuid.UUID) error {
	return s.skillRepo.DeletePosting(id)
}

type MatchFilter struct {
	SkillID   *uuid.UUID
	Location  string
	MinRating float64
	Method    *string
	MaxPrice  *float64
}

func (s *SkillService) MatchSkills(userID uuid.UUID, filter *MatchFilter, page, pageSize int) ([]*models.SkillPosting, int64, error) {
	return s.skillRepo.GetPostings(page, pageSize, filter.SkillID, nil)
}
