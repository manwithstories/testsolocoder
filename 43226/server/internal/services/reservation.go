package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"museum-server/internal/dto"
	"museum-server/internal/models"
	"museum-server/internal/repository"
	redisPkg "museum-server/pkg/redis"
	"museum-server/pkg/utils"
)

type ReservationService struct {
	reservationRepo *repository.ReservationRepository
	exhibitionRepo  *repository.ExhibitionRepository
	db              *gorm.DB
}

func NewReservationService(
	reservationRepo *repository.ReservationRepository,
	exhibitionRepo *repository.ExhibitionRepository,
	db *gorm.DB,
) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		exhibitionRepo:  exhibitionRepo,
		db:              db,
	}
}

func (s *ReservationService) Create(userID uint, req *dto.ReservationRequest) (*models.Reservation, error) {
	slot, err := s.exhibitionRepo.FindTimeSlotByID(req.TimeSlotID)
	if err != nil {
		return nil, fmt.Errorf("time slot not found")
	}

	exhibition, err := s.exhibitionRepo.FindByID(req.ExhibitionID)
	if err != nil {
		return nil, fmt.Errorf("exhibition not found")
	}

	if slot.ExhibitionID != req.ExhibitionID {
		return nil, fmt.Errorf("time slot does not belong to this exhibition")
	}

	ctx := context.Background()
	client := redisPkg.GetClient()
	lockKey := fmt.Sprintf("reservation_lock:slot:%d", req.TimeSlotID)
	locked, err := client.SetNX(ctx, lockKey, "1", 5*time.Second).Result()
	if err != nil || !locked {
		return nil, fmt.Errorf("too many concurrent requests, please try again")
	}
	defer client.Del(ctx, lockKey)

	bookedCount, err := s.reservationRepo.CountByTimeSlot(req.TimeSlotID)
	if err != nil {
		return nil, fmt.Errorf("failed to check slot availability")
	}

	if bookedCount+int64(req.VisitorCount) > int64(slot.MaxCapacity) {
		return nil, fmt.Errorf("this time slot is full")
	}

	var reservation *models.Reservation
	err = s.db.Transaction(func(tx *gorm.DB) error {
		qrCode := utils.GenerateQRCode(fmt.Sprintf("%d-%d-%d", userID, req.ExhibitionID, time.Now().UnixNano()))

		totalPrice := exhibition.TicketPrice * float64(req.VisitorCount)

		reservation = &models.Reservation{
			UserID:       userID,
			ExhibitionID: req.ExhibitionID,
			TimeSlotID:   req.TimeSlotID,
			VisitorCount: req.VisitorCount,
			GuideType:    req.GuideType,
			TotalPrice:   totalPrice,
			Status:       models.ReservationStatusPending,
			QRCode:       qrCode,
		}

		if err := tx.Create(reservation).Error; err != nil {
			return fmt.Errorf("failed to create reservation: %w", err)
		}

		slot.BookedCount = int(bookedCount) + req.VisitorCount
		if err := tx.Save(slot).Error; err != nil {
			return fmt.Errorf("failed to update time slot: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.cacheReservationStatus(req.ExhibitionID)

	return reservation, nil
}

func (s *ReservationService) Confirm(id uint) error {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("reservation not found")
	}

	if reservation.Status != models.ReservationStatusPending {
		return fmt.Errorf("reservation cannot be confirmed")
	}

	reservation.Status = models.ReservationStatusConfirmed
	return s.reservationRepo.Update(reservation)
}

func (s *ReservationService) Cancel(id uint, reason string) error {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("reservation not found")
	}

	if reservation.Status != models.ReservationStatusPending &&
		reservation.Status != models.ReservationStatusConfirmed {
		return fmt.Errorf("reservation cannot be cancelled")
	}

	now := time.Now()
	reservation.Status = models.ReservationStatusCancelled
	reservation.CancelledAt = &now
	reservation.Remark = reason

	slot, _ := s.exhibitionRepo.FindTimeSlotByID(reservation.TimeSlotID)
	if slot != nil && slot.BookedCount >= reservation.VisitorCount {
		slot.BookedCount -= reservation.VisitorCount
		s.exhibitionRepo.UpdateTimeSlot(slot)
	}

	return s.reservationRepo.Update(reservation)
}

func (s *ReservationService) Reschedule(id uint, newTimeSlotID uint) error {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("reservation not found")
	}

	if reservation.Status == models.ReservationStatusCancelled ||
		reservation.Status == models.ReservationStatusCompleted {
		return fmt.Errorf("reservation cannot be rescheduled")
	}

	newSlot, err := s.exhibitionRepo.FindTimeSlotByID(newTimeSlotID)
	if err != nil {
		return fmt.Errorf("new time slot not found")
	}

	if newSlot.ExhibitionID != reservation.ExhibitionID {
		return fmt.Errorf("new time slot does not belong to this exhibition")
	}

	bookedCount, _ := s.reservationRepo.CountByTimeSlot(newTimeSlotID)
	if bookedCount+int64(reservation.VisitorCount) > int64(newSlot.MaxCapacity) {
		return fmt.Errorf("new time slot is full")
	}

	oldSlot, _ := s.exhibitionRepo.FindTimeSlotByID(reservation.TimeSlotID)
	if oldSlot != nil && oldSlot.BookedCount >= reservation.VisitorCount {
		oldSlot.BookedCount -= reservation.VisitorCount
		s.exhibitionRepo.UpdateTimeSlot(oldSlot)
	}

	newSlot.BookedCount += reservation.VisitorCount
	s.exhibitionRepo.UpdateTimeSlot(newSlot)

	reservation.TimeSlotID = newTimeSlotID
	reservation.QRCode = utils.GenerateQRCode(fmt.Sprintf("%d-%d-%d", reservation.UserID, reservation.ExhibitionID, time.Now().UnixNano()))

	return s.reservationRepo.Update(reservation)
}

func (s *ReservationService) GetByID(id uint) (*models.Reservation, error) {
	return s.reservationRepo.FindByID(id)
}

func (s *ReservationService) GetByQRCode(qrCode string) (*models.Reservation, error) {
	return s.reservationRepo.FindByQRCode(qrCode)
}

func (s *ReservationService) ListByUser(userID uint, page, pageSize int) ([]models.Reservation, int64, error) {
	return s.reservationRepo.ListByUser(userID, page, pageSize)
}

func (s *ReservationService) ListByExhibition(exhibitionID uint, page, pageSize int, status string) ([]models.Reservation, int64, error) {
	return s.reservationRepo.ListByExhibition(exhibitionID, page, pageSize, status)
}

func (s *ReservationService) CheckIn(qrCode string) error {
	reservation, err := s.reservationRepo.FindByQRCode(qrCode)
	if err != nil {
		return fmt.Errorf("reservation not found")
	}

	if reservation.Status != models.ReservationStatusConfirmed {
		return fmt.Errorf("reservation is not confirmed")
	}

	now := time.Now()
	reservation.Status = models.ReservationStatusCompleted
	s.reservationRepo.Update(reservation)

	record := &models.VisitRecord{
		UserID:        reservation.UserID,
		ReservationID: reservation.ID,
		ExhibitionID:  reservation.ExhibitionID,
		CheckInTime:   &now,
	}

	return s.reservationRepo.CreateVisitRecord(record)
}

func (s *ReservationService) CheckOut(reservationID uint) error {
	record, err := s.reservationRepo.FindVisitRecordByReservation(reservationID)
	if err != nil {
		return fmt.Errorf("visit record not found")
	}

	now := time.Now()
	record.CheckOutTime = &now

	return s.reservationRepo.UpdateVisitRecord(record)
}

func (s *ReservationService) RateVisit(reservationID uint, rating int, comment string, favorite bool) error {
	record, err := s.reservationRepo.FindVisitRecordByReservation(reservationID)
	if err != nil {
		return fmt.Errorf("visit record not found")
	}

	record.Rating = rating
	record.Comment = comment
	record.Favorite = favorite

	return s.reservationRepo.UpdateVisitRecord(record)
}

func (s *ReservationService) ListVisitRecords(userID uint, page, pageSize int) ([]models.VisitRecord, int64, error) {
	return s.reservationRepo.ListVisitRecords(userID, page, pageSize)
}

func (s *ReservationService) GetUserVisitStats(userID uint) (map[string]interface{}, error) {
	return s.reservationRepo.GetUserVisitStats(userID)
}

func (s *ReservationService) cacheReservationStatus(exhibitionID uint) {
	ctx := context.Background()
	client := redisPkg.GetClient()

	slotCount, _ := s.reservationRepo.CountByTimeSlot(exhibitionID)
	key := fmt.Sprintf("reservation:exhibition:%d", exhibitionID)
	data, _ := json.Marshal(map[string]interface{}{
		"exhibition_id": exhibitionID,
		"reserved_count": slotCount,
		"updated_at":    time.Now(),
	})
	client.Set(ctx, key, data, 5*time.Minute)
}

func (s *ReservationService) GetCachedReservationStatus(exhibitionID uint) (map[string]interface{}, error) {
	ctx := context.Background()
	client := redisPkg.GetClient()

	key := fmt.Sprintf("reservation:exhibition:%d", exhibitionID)
	val, err := client.Get(ctx, key).Result()
	if err != nil || val == "" {
		slotCount, _ := s.reservationRepo.CountByTimeSlot(exhibitionID)
		return map[string]interface{}{
			"exhibition_id": exhibitionID,
			"reserved_count": slotCount,
		}, nil
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(val), &result)
	return result, nil
}
