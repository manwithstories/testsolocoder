package services

import (
	"errors"
	"medical-platform/internal/models"
	"medical-platform/pkg/database"

	"gorm.io/gorm"
)

type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"omitempty"`
	Location    string `json:"location" binding:"omitempty,max=200"`
}

type UpdateDepartmentRequest struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Description string `json:"description" binding:"omitempty"`
	Location    string `json:"location" binding:"omitempty,max=200"`
}

type DepartmentService struct {
	db *gorm.DB
}

func NewDepartmentService() *DepartmentService {
	return &DepartmentService{
		db: database.GetDB(),
	}
}

func (s *DepartmentService) GetList(page, pageSize int, keyword string) ([]models.Department, int64, error) {
	var departments []models.Department
	var total int64

	query := s.db.Model(&models.Department{})

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("获取科室总数失败")
	}

	if err := query.Scopes(database.Paginate(page, pageSize)).Find(&departments).Error; err != nil {
		return nil, 0, errors.New("获取科室列表失败")
	}

	return departments, total, nil
}

func (s *DepartmentService) GetByID(id uint) (*models.Department, error) {
	var department models.Department
	if err := s.db.Preload("Doctors").First(&department, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("科室不存在")
		}
		return nil, errors.New("获取科室详情失败")
	}
	return &department, nil
}

func (s *DepartmentService) Create(req *CreateDepartmentRequest) (*models.Department, error) {
	var count int64
	s.db.Model(&models.Department{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		return nil, errors.New("科室名称已存在")
	}

	department := &models.Department{
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
	}

	if err := s.db.Create(department).Error; err != nil {
		return nil, errors.New("创建科室失败")
	}

	return department, nil
}

func (s *DepartmentService) Update(id uint, req *UpdateDepartmentRequest) (*models.Department, error) {
	department, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" && req.Name != department.Name {
		var count int64
		s.db.Model(&models.Department{}).Where("name = ? AND id != ?", req.Name, id).Count(&count)
		if count > 0 {
			return nil, errors.New("科室名称已存在")
		}
		department.Name = req.Name
	}

	if req.Description != "" {
		department.Description = req.Description
	}

	if req.Location != "" {
		department.Location = req.Location
	}

	if err := s.db.Save(department).Error; err != nil {
		return nil, errors.New("更新科室失败")
	}

	return department, nil
}

func (s *DepartmentService) Delete(id uint) error {
	department, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.db.Delete(department).Error; err != nil {
		return errors.New("删除科室失败")
	}

	return nil
}
