package repository

import (
	"gorm.io/gorm"
	"ticket-system/internal/database"
	"ticket-system/internal/models"
)

type ShowRepository struct{}

func NewShowRepository() *ShowRepository {
	return &ShowRepository{}
}

func (r *ShowRepository) Create(show *models.Show) error {
	return database.DB.Create(show).Error
}

func (r *ShowRepository) GetByID(id uint64) (*models.Show, error) {
	var show models.Show
	err := database.DB.Preload("Sessions").First(&show, id).Error
	return &show, err
}

func (r *ShowRepository) List(page, pageSize int, status int, keyword string) ([]models.Show, int64, error) {
	var shows []models.Show
	var total int64

	query := database.DB.Model(&models.Show{})
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR artist LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Preload("Sessions").Offset(offset).Limit(pageSize).Order("created_at desc").Find(&shows).Error
	return shows, total, err
}

func (r *ShowRepository) Update(show *models.Show) error {
	return database.DB.Save(show).Error
}

func (r *ShowRepository) Delete(id uint64) error {
	return database.DB.Delete(&models.Show{}, id).Error
}

func (r *ShowRepository) CreateSession(session *models.Session) error {
	return database.DB.Create(session).Error
}

func (r *ShowRepository) GetSessionByID(id uint64) (*models.Session, error) {
	var session models.Session
	err := database.DB.Preload("SeatAreas").First(&session, id).Error
	return &session, err
}

func (r *ShowRepository) GetSessionsByShowID(showID uint64) ([]models.Session, error) {
	var sessions []models.Session
	err := database.DB.Where("show_id = ?", showID).Preload("SeatAreas").Order("start_time asc").Find(&sessions).Error
	return sessions, err
}

func (r *ShowRepository) UpdateSession(session *models.Session) error {
	return database.DB.Save(session).Error
}

func (r *ShowRepository) CreateSeatArea(area *models.SeatArea) error {
	return database.DB.Create(area).Error
}

func (r *ShowRepository) GetSeatAreaByID(id uint64) (*models.SeatArea, error) {
	var area models.SeatArea
	err := database.DB.First(&area, id).Error
	return &area, err
}

func (r *ShowRepository) GetSeatAreasBySessionID(sessionID uint64) ([]models.SeatArea, error) {
	var areas []models.SeatArea
	err := database.DB.Where("session_id = ?", sessionID).Order("sort_order asc").Find(&areas).Error
	return areas, err
}

func (r *ShowRepository) CreateSeat(seat *models.Seat) error {
	return database.DB.Create(seat).Error
}

func (r *ShowRepository) GetSeatByID(id uint64) (*models.Seat, error) {
	var seat models.Seat
	err := database.DB.First(&seat, id).Error
	return &seat, err
}

func (r *ShowRepository) GetSeatsBySessionID(sessionID uint64) ([]models.Seat, error) {
	var seats []models.Seat
	err := database.DB.Where("session_id = ?", sessionID).Find(&seats).Error
	return seats, err
}

func (r *ShowRepository) GetSeatsByIDs(sessionID uint64, seatIDs []uint64) ([]models.Seat, error) {
	var seats []models.Seat
	err := database.DB.Where("session_id = ? AND id IN ?", sessionID, seatIDs).Find(&seats).Error
	return seats, err
}

func (r *ShowRepository) UpdateSeat(seat *models.Seat) error {
	return database.DB.Save(seat).Error
}

func (r *ShowRepository) UpdateSeatsStatus(sessionID uint64, seatIDs []uint64, status int) error {
	return database.DB.Model(&models.Seat{}).
		Where("session_id = ? AND id IN ?", sessionID, seatIDs).
		Update("status", status).Error
}

func (r *ShowRepository) CreateSeatChart(chart *models.SeatChart) error {
	return database.DB.Create(chart).Error
}

func (r *ShowRepository) GetSeatChartBySessionID(sessionID uint64) (*models.SeatChart, error) {
	var chart models.SeatChart
	err := database.DB.Where("session_id = ?", sessionID).First(&chart).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &chart, err
}

func (r *ShowRepository) UpdateSeatChart(chart *models.SeatChart) error {
	return database.DB.Save(chart).Error
}

func (r *ShowRepository) BatchCreateSeats(seats []models.Seat) error {
	return database.DB.Create(&seats).Error
}

func (r *ShowRepository) UpdateSessionSoldSeats(sessionID uint64, count int) error {
	return database.DB.Model(&models.Session{}).Where("id = ?", sessionID).
		UpdateColumn("sold_seats", gorm.Expr("sold_seats + ?", count)).Error
}

func (r *ShowRepository) UpdateAreaSoldSeats(areaID uint64, count int) error {
	return database.DB.Model(&models.SeatArea{}).Where("id = ?", areaID).
		UpdateColumn("sold_seats", gorm.Expr("sold_seats + ?", count)).Error
}
