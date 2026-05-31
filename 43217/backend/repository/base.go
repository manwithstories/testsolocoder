package repository

import (
	"health-platform/config"

	"gorm.io/gorm"
)

type BaseRepository struct {
	DB *gorm.DB
}

func NewBaseRepository() *BaseRepository {
	return &BaseRepository{
		DB: config.DB,
	}
}

func (r *BaseRepository) Create(model interface{}) error {
	return r.DB.Create(model).Error
}

func (r *BaseRepository) Save(model interface{}) error {
	return r.DB.Save(model).Error
}

func (r *BaseRepository) Update(model interface{}) error {
	return r.DB.Save(model).Error
}

func (r *BaseRepository) Delete(model interface{}) error {
	return r.DB.Delete(model).Error
}

func (r *BaseRepository) FindByID(model interface{}, id uint) error {
	return r.DB.First(model, id).Error
}

func (r *BaseRepository) FindAll(model interface{}) error {
	return r.DB.Find(model).Error
}

func (r *BaseRepository) FindWithConditions(model interface{}, conditions map[string]interface{}) error {
	return r.DB.Where(conditions).Find(model).Error
}

func (r *BaseRepository) Count(model interface{}) (int64, error) {
	var count int64
	err := r.DB.Model(model).Count(&count).Error
	return count, err
}

func (r *BaseRepository) CountWithConditions(model interface{}, conditions map[string]interface{}) (int64, error) {
	var count int64
	err := r.DB.Model(model).Where(conditions).Count(&count).Error
	return count, err
}

func (r *BaseRepository) Paginate(model interface{}, page, pageSize int) (int64, error) {
	var total int64
	err := r.DB.Model(model).Count(&total).Error
	if err != nil {
		return 0, err
	}
	
	offset := (page - 1) * pageSize
	err = r.DB.Offset(offset).Limit(pageSize).Find(model).Error
	return total, err
}

func (r *BaseRepository) PaginateWithConditions(model interface{}, page, pageSize int, conditions map[string]interface{}) (int64, error) {
	var total int64
	query := r.DB.Model(model).Where(conditions)
	err := query.Count(&total).Error
	if err != nil {
		return 0, err
	}
	
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(model).Error
	return total, err
}
