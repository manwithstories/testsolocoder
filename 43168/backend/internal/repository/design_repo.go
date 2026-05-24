package repository

import (
	"errors"

	"gorm.io/gorm"

	"furniture-platform/internal/model"
)

type DesignRepository struct {
	db *gorm.DB
}

func NewDesignRepository(db *gorm.DB) *DesignRepository {
	return &DesignRepository{db: db}
}

func (r *DesignRepository) CreateProject(project *model.DesignProject) error {
	return r.db.Create(project).Error
}

func (r *DesignRepository) GetProjectByID(id uint) (*model.DesignProject, error) {
	var p model.DesignProject
	err := r.db.First(&p, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *DesignRepository) UpdateProject(project *model.DesignProject) error {
	return r.db.Save(project).Error
}

func (r *DesignRepository) DeleteProject(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ?", id).Delete(&model.DesignComment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("project_id = ?", id).Delete(&model.DesignImage{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.DesignProject{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

type ProjectListParams struct {
	Page     int
	PageSize int
	Keyword  string
	Status   string
	Role     string
	UserID   uint
}

func (r *DesignRepository) ListProjects(params *ProjectListParams) ([]*model.DesignProject, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	query := r.db.Model(&model.DesignProject{})

	if params.Keyword != "" {
		like := "%" + params.Keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", like, like)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Role == model.RoleDesigner {
		query = query.Where("designer_id = ?", params.UserID)
	} else if params.Role == model.RoleOwner {
		query = query.Where("owner_id = ?", params.UserID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*model.DesignProject
	err := query.Order("id DESC").
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *DesignRepository) CreateImage(img *model.DesignImage) error {
	return r.db.Create(img).Error
}

func (r *DesignRepository) ListImages(projectID uint) ([]*model.DesignImage, error) {
	var list []*model.DesignImage
	err := r.db.Where("project_id = ?", projectID).Order("sort ASC, id ASC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *DesignRepository) GetImageByID(id uint) (*model.DesignImage, error) {
	var img model.DesignImage
	err := r.db.First(&img, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &img, nil
}

func (r *DesignRepository) DeleteImage(id uint) error {
	return r.db.Delete(&model.DesignImage{}, id).Error
}

func (r *DesignRepository) CreateComment(comment *model.DesignComment) error {
	return r.db.Create(comment).Error
}

func (r *DesignRepository) ListComments(projectID uint) ([]*model.DesignComment, error) {
	var list []*model.DesignComment
	err := r.db.Where("project_id = ?", projectID).Order("id ASC").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *DesignRepository) GetCommentByID(id uint) (*model.DesignComment, error) {
	var c model.DesignComment
	err := r.db.First(&c, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}
