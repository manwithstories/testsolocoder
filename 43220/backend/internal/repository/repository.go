package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"pet-board/internal/database"
	"pet-board/internal/models"
	"pet-board/internal/dto"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("StoreInfo").Preload("KeeperInfo").First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("StoreInfo").Preload("KeeperInfo").First(&user, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("StoreInfo").Preload("KeeperInfo").First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return database.DB.Save(user).Error
}

func (r *UserRepository) List(page, pageSize int, role string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{})
	if role != "" {
		query = query.Where("role = ?", role)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("StoreInfo").Preload("KeeperInfo").
		Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) CreateStoreInfo(info *models.StoreInfo) error {
	return database.DB.Create(info).Error
}

func (r *UserRepository) UpdateStoreInfo(info *models.StoreInfo) error {
	return database.DB.Save(info).Error
}

func (r *UserRepository) CreateKeeperInfo(info *models.KeeperInfo) error {
	return database.DB.Create(info).Error
}

func (r *UserRepository) UpdateKeeperInfo(info *models.KeeperInfo) error {
	return database.DB.Save(info).Error
}

func (r *UserRepository) GetStoreInfoByUserID(userID uuid.UUID) (*models.StoreInfo, error) {
	var info models.StoreInfo
	err := database.DB.First(&info, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func (r *UserRepository) GetKeeperInfoByUserID(userID uuid.UUID) (*models.KeeperInfo, error) {
	var info models.KeeperInfo
	err := database.DB.Preload("StoreInfo").First(&info, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

type PetRepository struct{}

func NewPetRepository() *PetRepository {
	return &PetRepository{}
}

func (r *PetRepository) Create(pet *models.Pet) error {
	return database.DB.Create(pet).Error
}

func (r *PetRepository) GetByID(id uuid.UUID) (*models.Pet, error) {
	var pet models.Pet
	err := database.DB.Preload("VaccineRecords").First(&pet, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pet, nil
}

func (r *PetRepository) ListByOwner(ownerID uuid.UUID, page, pageSize int) ([]models.Pet, int64, error) {
	var pets []models.Pet
	var total int64

	query := database.DB.Model(&models.Pet{}).Where("owner_id = ?", ownerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("VaccineRecords").
		Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&pets).Error; err != nil {
		return nil, 0, err
	}

	return pets, total, nil
}

func (r *PetRepository) Update(pet *models.Pet) error {
	return database.DB.Save(pet).Error
}

func (r *PetRepository) Delete(id uuid.UUID) error {
	return database.DB.Delete(&models.Pet{}, "id = ?", id).Error
}

func (r *PetRepository) CreateVaccineRecord(record *models.VaccineRecord) error {
	return database.DB.Create(record).Error
}

func (r *PetRepository) GetVaccineRecords(petID uuid.UUID) ([]models.VaccineRecord, error) {
	var records []models.VaccineRecord
	err := database.DB.Where("pet_id = ?", petID).Order("vaccinated_at DESC").Find(&records).Error
	return records, err
}

func (r *PetRepository) HasValidVaccine(petID uuid.UUID) bool {
	var count int64
	now := time.Now()
	database.DB.Model(&models.VaccineRecord{}).
		Where("pet_id = ? AND expire_at >= ? AND is_valid = ?", petID, now, true).
		Count(&count)
	return count > 0
}

func (r *PetRepository) CreateDewormRecord(record *models.DewormRecord) error {
	return database.DB.Create(record).Error
}

func (r *PetRepository) GetDewormRecords(petID uuid.UUID) ([]models.DewormRecord, error) {
	var records []models.DewormRecord
	err := database.DB.Where("pet_id = ?", petID).Order("dewormed_at DESC").Find(&records).Error
	return records, err
}

type PackageRepository struct{}

func NewPackageRepository() *PackageRepository {
	return &PackageRepository{}
}

func (r *PackageRepository) Create(pkg *models.BoardingPackage) error {
	return database.DB.Create(pkg).Error
}

func (r *PackageRepository) GetByID(id uuid.UUID) (*models.BoardingPackage, error) {
	var pkg models.BoardingPackage
	err := database.DB.First(&pkg, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pkg, nil
}

func (r *PackageRepository) ListByStore(storeID uuid.UUID, pkgType string, page, pageSize int) ([]models.BoardingPackage, int64, error) {
	var packages []models.BoardingPackage
	var total int64

	query := database.DB.Model(&models.BoardingPackage{}).Where("store_id = ?", storeID)
	if pkgType != "" {
		query = query.Where("type = ?", pkgType)
	}
	query = query.Where("is_available = ?", true)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("sort_order ASC, created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&packages).Error; err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

func (r *PackageRepository) Update(pkg *models.BoardingPackage) error {
	return database.DB.Save(pkg).Error
}

func (r *PackageRepository) Delete(id uuid.UUID) error {
	return database.DB.Delete(&models.BoardingPackage{}, "id = ?", id).Error
}

type ReservationRepository struct{}

func NewReservationRepository() *ReservationRepository {
	return &ReservationRepository{}
}

func (r *ReservationRepository) Create(reservation *models.Reservation) error {
	return database.DB.Create(reservation).Error
}

func (r *ReservationRepository) GetByID(id uuid.UUID) (*models.Reservation, error) {
	var res models.Reservation
	err := database.DB.First(&res, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *ReservationRepository) GetByOrderNo(orderNo string) (*models.Reservation, error) {
	var res models.Reservation
	err := database.DB.First(&res, "order_no = ?", orderNo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *ReservationRepository) List(query dto.StatisticsQuery, page, pageSize int, ownerID, storeID *uuid.UUID, status string) ([]models.Reservation, int64, error) {
	var reservations []models.Reservation
	var total int64

	q := database.DB.Model(&models.Reservation{})
	if ownerID != nil {
		q = q.Where("owner_id = ?", *ownerID)
	}
	if storeID != nil {
		q = q.Where("store_id = ?", *storeID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if query.StartDate != nil {
		q = q.Where("check_in_date >= ?", *query.StartDate)
	}
	if query.EndDate != nil {
		q = q.Where("check_out_date <= ?", *query.EndDate)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := q.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&reservations).Error; err != nil {
		return nil, 0, err
	}

	return reservations, total, nil
}

func (r *ReservationRepository) Update(reservation *models.Reservation) error {
	return database.DB.Save(reservation).Error
}

func (r *ReservationRepository) CheckConflict(storeID, packageID uuid.UUID, checkIn, checkOut time.Time, excludeID *uuid.UUID) (int64, error) {
	var count int64
	q := database.DB.Model(&models.Reservation{}).
		Where("store_id = ? AND package_id = ?", storeID, packageID).
		Where("status IN ?", []string{models.ReservationStatusPending, models.ReservationStatusConfirmed, models.ReservationStatusCheckedIn}).
		Where("check_in_date < ? AND check_out_date > ?", checkOut, checkIn)

	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}

	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

type DailyRecordRepository struct{}

func NewDailyRecordRepository() *DailyRecordRepository {
	return &DailyRecordRepository{}
}

func (r *DailyRecordRepository) Create(record *models.DailyRecord) error {
	return database.DB.Create(record).Error
}

func (r *DailyRecordRepository) GetByID(id uuid.UUID) (*models.DailyRecord, error) {
	var record models.DailyRecord
	err := database.DB.First(&record, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *DailyRecordRepository) ListByReservation(reservationID uuid.UUID, page, pageSize int) ([]models.DailyRecord, int64, error) {
	var records []models.DailyRecord
	var total int64

	query := database.DB.Model(&models.DailyRecord{}).Where("reservation_id = ?", reservationID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("record_date DESC").
		Limit(pageSize).Offset(offset).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (r *DailyRecordRepository) ListByPet(petID uuid.UUID, page, pageSize int) ([]models.DailyRecord, int64, error) {
	var records []models.DailyRecord
	var total int64

	query := database.DB.Model(&models.DailyRecord{}).Where("pet_id = ?", petID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("record_date DESC").
		Limit(pageSize).Offset(offset).
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (r *DailyRecordRepository) Update(record *models.DailyRecord) error {
	return database.DB.Save(record).Error
}

type ReviewRepository struct{}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	return database.DB.Create(review).Error
}

func (r *ReviewRepository) GetByID(id uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := database.DB.First(&review, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) GetByReservationID(reservationID uuid.UUID) (*models.Review, error) {
	var review models.Review
	err := database.DB.First(&review, "reservation_id = ?", reservationID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepository) ListByStore(storeID uuid.UUID, page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := database.DB.Model(&models.Review{}).Where("store_id = ?", storeID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *ReviewRepository) ListByKeeper(keeperID uuid.UUID, page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := database.DB.Model(&models.Review{}).Where("keeper_id = ?", keeperID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *ReviewRepository) Update(review *models.Review) error {
	return database.DB.Save(review).Error
}

func (r *ReviewRepository) GetStoreRating(storeID uuid.UUID) (float64, int, error) {
	type result struct {
		AvgRating float64
		Count     int
	}
	var res result
	err := database.DB.Model(&models.Review{}).
		Where("store_id = ?", storeID).
		Select("COALESCE(AVG(store_rating), 0) as avg_rating, COUNT(*) as count").
		Scan(&res).Error
	if err != nil {
		return 0, 0, err
	}
	return res.AvgRating, res.Count, nil
}

func (r *ReviewRepository) GetKeeperRating(keeperID uuid.UUID) (float64, int, error) {
	type result struct {
		AvgRating float64
		Count     int
	}
	var res result
	err := database.DB.Model(&models.Review{}).
		Where("keeper_id = ?", keeperID).
		Select("COALESCE(AVG(keeper_rating), 0) as avg_rating, COUNT(*) as count").
		Scan(&res).Error
	if err != nil {
		return 0, 0, err
	}
	return res.AvgRating, res.Count, nil
}

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) Create(order *models.Order, tx *gorm.DB) error {
	if tx != nil {
		return tx.Create(order).Error
	}
	return database.DB.Create(order).Error
}

func (r *OrderRepository) GetByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := database.DB.First(&order, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByOrderNo(orderNo string) (*models.Order, error) {
	var order models.Order
	err := database.DB.First(&order, "order_no = ?", orderNo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByReservationID(reservationID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	err := database.DB.Where("reservation_id = ?", reservationID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) List(ownerID, storeID *uuid.UUID, payStatus, orderType string, query dto.StatisticsQuery, page, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	q := database.DB.Model(&models.Order{})
	if ownerID != nil {
		q = q.Where("owner_id = ?", *ownerID)
	}
	if storeID != nil {
		q = q.Where("store_id = ?", *storeID)
	}
	if payStatus != "" {
		q = q.Where("pay_status = ?", payStatus)
	}
	if orderType != "" {
		q = q.Where("type = ?", orderType)
	}
	if query.StartDate != nil {
		q = q.Where("created_at >= ?", *query.StartDate)
	}
	if query.EndDate != nil {
		q = q.Where("created_at <= ?", *query.EndDate)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := q.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) Update(order *models.Order, tx *gorm.DB) error {
	if tx != nil {
		return tx.Save(order).Error
	}
	return database.DB.Save(order).Error
}

func (r *OrderRepository) BeginTransaction() *gorm.DB {
	return database.DB.Begin()
}

type HealthAlertRepository struct{}

func NewHealthAlertRepository() *HealthAlertRepository {
	return &HealthAlertRepository{}
}

func (r *HealthAlertRepository) Create(alert *models.HealthAlert) error {
	return database.DB.Create(alert).Error
}

func (r *HealthAlertRepository) ListByUser(userID uuid.UUID, isRead *bool, page, pageSize int) ([]models.HealthAlert, int64, error) {
	var alerts []models.HealthAlert
	var total int64

	query := database.DB.Model(&models.HealthAlert{}).Where("user_id = ?", userID)
	if isRead != nil {
		query = query.Where("is_read = ?", *isRead)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&alerts).Error; err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

func (r *HealthAlertRepository) MarkAsRead(id uuid.UUID) error {
	return database.DB.Model(&models.HealthAlert{}).
		Where("id = ?", id).
		Update("is_read", true).Error
}

func (r *HealthAlertRepository) MarkAllAsRead(userID uuid.UUID) error {
	return database.DB.Model(&models.HealthAlert{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

func (r *HealthAlertRepository) GetExpiringVaccines(days int) ([]models.VaccineRecord, error) {
	var records []models.VaccineRecord
	threshold := time.Now().AddDate(0, 0, days)
	err := database.DB.Where("expire_at <= ? AND expire_at > ? AND is_valid = ?", threshold, time.Now(), true).
		Find(&records).Error
	return records, err
}

func (r *HealthAlertRepository) GetExpiringDeworms(days int) ([]models.DewormRecord, error) {
	var records []models.DewormRecord
	threshold := time.Now().AddDate(0, 0, days)
	err := database.DB.Where("expire_at <= ? AND expire_at > ? AND is_valid = ?", threshold, time.Now(), true).
		Find(&records).Error
	return records, err
}

type OperationLogRepository struct{}

func NewOperationLogRepository() *OperationLogRepository {
	return &OperationLogRepository{}
}

func (r *OperationLogRepository) Create(log *models.OperationLog) error {
	return database.DB.Create(log).Error
}

func (r *OperationLogRepository) List(page, pageSize int, userID *uuid.UUID) ([]models.OperationLog, int64, error) {
	var logs []models.OperationLog
	var total int64

	query := database.DB.Model(&models.OperationLog{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

type StatisticsRepository struct{}

func NewStatisticsRepository() *StatisticsRepository {
	return &StatisticsRepository{}
}

func (r *StatisticsRepository) GetRevenueTrend(storeID uuid.UUID, start, end time.Time) ([]map[string]interface{}, error) {
	type row struct {
		Date  string  `gorm:"column:date"`
		Total float64 `gorm:"column:total"`
	}
	var rows []row

	err := database.DB.Model(&models.Order{}).
		Select("DATE(created_at) as date, COALESCE(SUM(amount), 0) as total").
		Where("store_id = ? AND pay_status = ? AND created_at >= ? AND created_at <= ?",
			storeID, models.PayStatusPaid, start, end).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(rows))
	for i, r := range rows {
		result[i] = map[string]interface{}{
			"date":  r.Date,
			"total": r.Total,
		}
	}
	return result, nil
}

func (r *StatisticsRepository) GetOccupancyRate(storeID uuid.UUID, start, end time.Time) (float64, error) {
	type row struct {
		TotalDays     float64 `gorm:"column:total_days"`
		OccupiedDays  float64 `gorm:"column:occupied_days"`
	}

	var totalReservations float64
	err := database.DB.Model(&models.Reservation{}).
		Where("store_id = ? AND status IN ? AND check_in_date >= ? AND check_in_date <= ?",
			storeID,
			[]string{models.ReservationStatusConfirmed, models.ReservationStatusCheckedIn, models.ReservationStatusCompleted},
			start, end).
		Select("COALESCE(SUM(total_days), 0)").
		Scan(&totalReservations).Error
	if err != nil {
		return 0, err
	}

	var totalCapacity int
	err = database.DB.Model(&models.BoardingPackage{}).
		Where("store_id = ?", storeID).
		Select("COALESCE(SUM(capacity), 0)").
		Scan(&totalCapacity).Error
	if err != nil {
		return 0, err
	}

	days := end.Sub(start).Hours()/24 + 1
	if days <= 0 || totalCapacity == 0 {
		return 0, nil
	}

	return totalReservations / (float64(totalCapacity) * days), nil
}

func (r *StatisticsRepository) GetPetTypeDistribution(storeID uuid.UUID, start, end time.Time) ([]map[string]interface{}, error) {
	type row struct {
		Species string `gorm:"column:species"`
		Count   int64  `gorm:"column:count"`
	}

	var rows []row
	err := database.DB.Table("reservations r").
		Select("p.species, COUNT(DISTINCT r.pet_id) as count").
		Joins("JOIN pets p ON p.id = r.pet_id").
		Where("r.store_id = ? AND r.status != ? AND r.created_at >= ? AND r.created_at <= ?",
			storeID, models.ReservationStatusCancelled, start, end).
		Group("p.species").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(rows))
	for i, r := range rows {
		result[i] = map[string]interface{}{
			"species": r.Species,
			"count":   r.Count,
		}
	}
	return result, nil
}

func (r *StatisticsRepository) GetOrderStatistics(storeID uuid.UUID, start, end time.Time) (map[string]interface{}, error) {
	type row struct {
		TotalAmount    float64 `gorm:"column:total_amount"`
		OrderCount     int64   `gorm:"column:order_count"`
		RefundAmount   float64 `gorm:"column:refund_amount"`
		RefundCount    int64   `gorm:"column:refund_count"`
	}

	var stat row
	err := database.DB.Model(&models.Order{}).
		Select(`COALESCE(SUM(CASE WHEN pay_status = 'paid' THEN amount ELSE 0 END), 0) as total_amount,
			COUNT(CASE WHEN pay_status = 'paid' THEN 1 END) as order_count,
			COALESCE(SUM(CASE WHEN pay_status = 'refunded' THEN refund_amount ELSE 0 END), 0) as refund_amount,
			COUNT(CASE WHEN pay_status = 'refunded' THEN 1 END) as refund_count`).
		Where("store_id = ? AND created_at >= ? AND created_at <= ?", storeID, start, end).
		Scan(&stat).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_amount":   stat.TotalAmount,
		"order_count":    stat.OrderCount,
		"refund_amount":  stat.RefundAmount,
		"refund_count":   stat.RefundCount,
	}, nil
}
