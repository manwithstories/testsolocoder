package repository

import (
	"errors"

	"gorm.io/gorm"

	"furniture-platform/internal/model"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户数据访问层
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据 ID 查询用户
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名查询用户
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 根据 ID 删除用户（软删除由调用方决定，此处物理删除）
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// ListParams 列表查询参数
type ListParams struct {
	Page     int
	PageSize int
	Keyword  string
	Role     string
	Status   *int
}

// List 分页搜索用户列表
func (r *UserRepository) List(params *ListParams) ([]*model.User, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	query := r.db.Model(&model.User{})

	if params.Keyword != "" {
		like := "%" + params.Keyword + "%"
		query = query.Where("username LIKE ? OR nickname LIKE ? OR phone LIKE ? OR email LIKE ?",
			like, like, like, like)
	}
	if params.Role != "" {
		query = query.Where("role = ?", params.Role)
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []*model.User
	err := query.Order("id DESC").
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
