package repository

import (
	"encoding/json"
	"venue-booking/internal/dto"
	"venue-booking/internal/model"
	"venue-booking/pkg/database"
)

type VenueRepository struct{}

func NewVenueRepository() *VenueRepository {
	return &VenueRepository{}
}

func (r *VenueRepository) Create(venue *model.Venue) error {
	return database.DB.Create(venue).Error
}

func (r *VenueRepository) GetByID(id uint) (*model.Venue, error) {
	var venue model.Venue
	err := database.DB.First(&venue, id).Error
	return &venue, err
}

func (r *VenueRepository) List(req *dto.VenueListRequest) ([]model.Venue, int64, error) {
	var venues []model.Venue
	var total int64

	query := database.DB.Model(&model.Venue{})

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	err := query.Offset(offset).Limit(req.PageSize).Order("id DESC").Find(&venues).Error
	return venues, total, err
}

func (r *VenueRepository) Update(venue *model.Venue) error {
	return database.DB.Save(venue).Error
}

func (r *VenueRepository) Delete(id uint) error {
	return database.DB.Delete(&model.Venue{}, id).Error
}

func (r *VenueRepository) UpdateStatus(id uint, status model.VenueStatus) error {
	return database.DB.Model(&model.Venue{}).Where("id = ?", id).Update("status", status).Error
}

func (r *VenueRepository) SetPrice(venueID uint, dayOfWeek int, timeSlots []dto.TimeSlotPrice) error {
	var price model.VenuePrice
	err := database.DB.Where("venue_id = ? AND day_of_week = ?", venueID, dayOfWeek).First(&price).Error

	timeSlotsJSON, _ := json.Marshal(timeSlots)

	if err != nil {
		price = model.VenuePrice{
			VenueID:   venueID,
			DayOfWeek: dayOfWeek,
			TimeSlots: string(timeSlotsJSON),
		}
		return database.DB.Create(&price).Error
	}

	price.TimeSlots = string(timeSlotsJSON)
	return database.DB.Save(&price).Error
}

func (r *VenueRepository) GetPrice(venueID uint, dayOfWeek int) ([]dto.TimeSlotPrice, error) {
	var price model.VenuePrice
	err := database.DB.Where("venue_id = ? AND day_of_week = ?", venueID, dayOfWeek).First(&price).Error
	if err != nil {
		return nil, err
	}

	var timeSlots []dto.TimeSlotPrice
	err = json.Unmarshal([]byte(price.TimeSlots), &timeSlots)
	return timeSlots, err
}

func (r *VenueRepository) GetAllPrices(venueID uint) (map[int][]dto.TimeSlotPrice, error) {
	var prices []model.VenuePrice
	err := database.DB.Where("venue_id = ?", venueID).Find(&prices).Error
	if err != nil {
		return nil, err
	}

	result := make(map[int][]dto.TimeSlotPrice)
	for _, p := range prices {
		var timeSlots []dto.TimeSlotPrice
		json.Unmarshal([]byte(p.TimeSlots), &timeSlots)
		result[p.DayOfWeek] = timeSlots
	}

	return result, nil
}
