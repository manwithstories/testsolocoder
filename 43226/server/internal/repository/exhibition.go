package repository

import (
	"time"

	"gorm.io/gorm"

	"museum-server/internal/dto"
	"museum-server/internal/models"
)

type ExhibitionRepository struct {
	db *gorm.DB
}

func NewExhibitionRepository(db *gorm.DB) *ExhibitionRepository {
	return &ExhibitionRepository{db: db}
}

func (r *ExhibitionRepository) Create(exhibition *models.Exhibition) error {
	return r.db.Create(exhibition).Error
}

func (r *ExhibitionRepository) FindByID(id uint) (*models.Exhibition, error) {
	var exhibition models.Exhibition
	err := r.db.Preload("Collections").Preload("Museum").First(&exhibition, id).Error
	if err != nil {
		return nil, err
	}
	return &exhibition, nil
}

func (r *ExhibitionRepository) Update(exhibition *models.Exhibition) error {
	return r.db.Save(exhibition).Error
}

func (r *ExhibitionRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("exhibition_id = ?", id).Delete(&models.ExhibitionCollection{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Exhibition{}, id).Error
	})
}

func (r *ExhibitionRepository) List(query *dto.ExhibitionListQuery) ([]models.Exhibition, int64, error) {
	var exhibitions []models.Exhibition
	var total int64

	db := r.db.Model(&models.Exhibition{})

	if query.Keyword != "" {
		db = db.Where("title LIKE ? OR description LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if !query.StartDate.IsZero() {
		db = db.Where("start_date >= ?", query.StartDate)
	}
	if !query.EndDate.IsZero() {
		db = db.Where("end_date <= ?", query.EndDate)
	}

	db.Count(&total)
	db.Preload("Collections").Preload("Museum").
		Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize).
		Order("created_at DESC").
		Find(&exhibitions)

	return exhibitions, total, nil
}

func (r *ExhibitionRepository) AddCollections(exhibitionID uint, collectionIDs []uint, sortOrders map[uint]int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, cid := range collectionIDs {
			sortOrder := 0
			if so, ok := sortOrders[cid]; ok {
				sortOrder = so
			}
			ec := models.ExhibitionCollection{
				ExhibitionID: exhibitionID,
				CollectionID: cid,
				SortOrder:    sortOrder,
				CreatedAt:    time.Now(),
			}
			if err := tx.Create(&ec).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ExhibitionRepository) RemoveCollections(exhibitionID uint, collectionIDs []uint) error {
	return r.db.Where("exhibition_id = ? AND collection_id IN ?", exhibitionID, collectionIDs).
		Delete(&models.ExhibitionCollection{}).Error
}

func (r *ExhibitionRepository) GetCollections(exhibitionID uint) ([]models.Collection, error) {
	var collections []models.Collection
	err := r.db.Joins("JOIN exhibition_collections ON collections.id = exhibition_collections.collection_id").
		Where("exhibition_collections.exhibition_id = ?", exhibitionID).
		Order("exhibition_collections.sort_order ASC").
		Find(&collections).Error
	return collections, err
}

func (r *ExhibitionRepository) CreateTimeSlot(slot *models.TimeSlot) error {
	return r.db.Create(slot).Error
}

func (r *ExhibitionRepository) FindTimeSlotByID(id uint) (*models.TimeSlot, error) {
	var slot models.TimeSlot
	err := r.db.First(&slot, id).Error
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *ExhibitionRepository) ListTimeSlots(exhibitionID uint, date time.Time) ([]models.TimeSlot, error) {
	var slots []models.TimeSlot
	db := r.db.Where("exhibition_id = ?", exhibitionID)
	if !date.IsZero() {
		db = db.Where("date = ?", date)
	}
	err := db.Order("start_time ASC").Find(&slots).Error
	return slots, err
}

func (r *ExhibitionRepository) UpdateTimeSlot(slot *models.TimeSlot) error {
	return r.db.Save(slot).Error
}

func (r *ExhibitionRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Exhibition{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
