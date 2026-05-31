package repository

import (
	"time"

	"gorm.io/gorm"

	"museum-server/internal/dto"
	"museum-server/internal/models"
)

type ReservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) Create(reservation *models.Reservation) error {
	return r.db.Create(reservation).Error
}

func (r *ReservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.Preload("User").Preload("Exhibition").Preload("TimeSlot").First(&reservation, id).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *ReservationRepository) FindByQRCode(qrCode string) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.Preload("User").Preload("Exhibition").Preload("TimeSlot").
		Where("qr_code = ?", qrCode).First(&reservation).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *ReservationRepository) Update(reservation *models.Reservation) error {
	return r.db.Save(reservation).Error
}

func (r *ReservationRepository) ListByUser(userID uint, page, pageSize int) ([]models.Reservation, int64, error) {
	var reservations []models.Reservation
	var total int64

	r.db.Model(&models.Reservation{}).Where("user_id = ?", userID).Count(&total)
	r.db.Preload("User").Preload("Exhibition").Preload("TimeSlot").
		Where("user_id = ?", userID).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&reservations)

	return reservations, total, nil
}

func (r *ReservationRepository) ListByExhibition(exhibitionID uint, page, pageSize int, status string) ([]models.Reservation, int64, error) {
	var reservations []models.Reservation
	var total int64

	db := r.db.Model(&models.Reservation{}).Where("exhibition_id = ?", exhibitionID)
	if status != "" {
		db = db.Where("status = ?", status)
	}
	db.Count(&total)
	db.Preload("User").Preload("Exhibition").Preload("TimeSlot").
		Where("exhibition_id = ?", exhibitionID).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&reservations)

	return reservations, total, nil
}

func (r *ReservationRepository) CountByTimeSlot(timeSlotID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Reservation{}).
		Where("time_slot_id = ? AND status IN ?", timeSlotID, []string{"pending", "confirmed"}).
		Count(&count).Error
	return count, err
}

func (r *ReservationRepository) CreateVisitRecord(record *models.VisitRecord) error {
	return r.db.Create(record).Error
}

func (r *ReservationRepository) UpdateVisitRecord(record *models.VisitRecord) error {
	return r.db.Save(record).Error
}

func (r *ReservationRepository) FindVisitRecordByReservation(reservationID uint) (*models.VisitRecord, error) {
	var record models.VisitRecord
	err := r.db.Where("reservation_id = ?", reservationID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *ReservationRepository) ListVisitRecords(userID uint, page, pageSize int) ([]models.VisitRecord, int64, error) {
	var records []models.VisitRecord
	var total int64

	r.db.Model(&models.VisitRecord{}).Where("user_id = ?", userID).Count(&total)
	r.db.Where("user_id = ?", userID).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&records)

	return records, total, nil
}

func (r *ReservationRepository) GetUserVisitStats(userID uint) (map[string]interface{}, error) {
	var totalVisits int64
	r.db.Model(&models.VisitRecord{}).Where("user_id = ?", userID).Count(&totalVisits)

	var records []models.VisitRecord
	r.db.Where("user_id = ? AND rating > 0", userID).Find(&records)

	type ExhibitionCount struct {
		ExhibitionID uint
		Count        int64
	}
	var exhibitionCounts []ExhibitionCount
	r.db.Model(&models.VisitRecord{}).
		Select("exhibition_id, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("exhibition_id").
		Scan(&exhibitionCounts)

	stats := map[string]interface{}{
		"total_visits":      totalVisits,
		"exhibition_counts": exhibitionCounts,
	}
	return stats, nil
}

func (r *ReservationRepository) GetStatistics(query *dto.StatisticsQuery) ([]models.Statistic, error) {
	var stats []models.Statistic
	db := r.db.Model(&models.Statistic{})

	if query.ExhibitionID > 0 {
		db = db.Where("exhibition_id = ?", query.ExhibitionID)
	}
	if !query.StartDate.IsZero() {
		db = db.Where("stat_date >= ?", query.StartDate)
	}
	if !query.EndDate.IsZero() {
		db = db.Where("stat_date <= ?", query.EndDate)
	}

	err := db.Order("stat_date ASC").Find(&stats).Error
	return stats, err
}

func (r *ReservationRepository) CreateStatistic(stat *models.Statistic) error {
	return r.db.Create(stat).Error
}

func (r *ReservationRepository) AggregateDailyStats(exhibitionID uint, date time.Time) (visitorCount, reservationCount int64, revenue float64, err error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	r.db.Model(&models.Reservation{}).
		Where("exhibition_id = ? AND created_at >= ? AND created_at < ? AND status = ?",
			exhibitionID, startOfDay, endOfDay, models.ReservationStatusCompleted).
		Count(&reservationCount)

	r.db.Model(&models.VisitRecord{}).
		Where("check_in_time >= ? AND check_in_time < ?", startOfDay, endOfDay).
		Count(&visitorCount)

	r.db.Table("reservations").
		Where("exhibition_id = ? AND created_at >= ? AND created_at < ? AND status = ?",
			exhibitionID, startOfDay, endOfDay, models.ReservationStatusCompleted).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&revenue)

	return
}
