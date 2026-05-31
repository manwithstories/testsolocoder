package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"pet-board/internal/dto"
	"pet-board/internal/models"
	"pet-board/internal/repository"
	"pet-board/internal/utils"
)

type ReservationService struct {
	resRepo   *repository.ReservationRepository
	pkgRepo   *repository.PackageRepository
	petRepo   *repository.PetRepository
	orderRepo *repository.OrderRepository
	logger    *logrus.Logger
}

func NewReservationService(
	resRepo *repository.ReservationRepository,
	pkgRepo *repository.PackageRepository,
	petRepo *repository.PetRepository,
	orderRepo *repository.OrderRepository,
	logger *logrus.Logger,
) *ReservationService {
	return &ReservationService{
		resRepo:   resRepo,
		pkgRepo:   pkgRepo,
		petRepo:   petRepo,
		orderRepo: orderRepo,
		logger:    logger,
	}
}

func (s *ReservationService) Create(ownerID uuid.UUID, req *dto.ReservationRequest) (*models.Reservation, error) {
	petID, err := uuid.Parse(req.PetID)
	if err != nil {
		return nil, errors.New("invalid pet ID")
	}
	storeID, err := uuid.Parse(req.StoreID)
	if err != nil {
		return nil, errors.New("invalid store ID")
	}
	pkgID, err := uuid.Parse(req.PackageID)
	if err != nil {
		return nil, errors.New("invalid package ID")
	}

	pet, err := s.petRepo.GetByID(petID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pet: %w", err)
	}
	if pet == nil {
		return nil, errors.New("pet not found")
	}
	if pet.OwnerID != ownerID {
		return nil, errors.New("not authorized to create reservation for this pet")
	}

	if !s.petRepo.HasValidVaccine(petID) {
		return nil, errors.New("pet does not have valid vaccine records")
	}

	pkg, err := s.pkgRepo.GetByID(pkgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get package: %w", err)
	}
	if pkg == nil {
		return nil, errors.New("package not found")
	}
	if pkg.StoreID != storeID {
		return nil, errors.New("package does not belong to this store")
	}
	if !pkg.IsAvailable {
		return nil, errors.New("package is not available")
	}

	if req.CheckInDate.After(req.CheckOutDate) {
		return nil, errors.New("check-in date must be before check-out date")
	}
	if req.CheckInDate.Before(time.Now()) {
		return nil, errors.New("check-in date cannot be in the past")
	}

	conflictCount, err := s.resRepo.CheckConflict(storeID, pkgID, req.CheckInDate, req.CheckOutDate, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check conflict: %w", err)
	}
	if conflictCount >= int64(pkg.Capacity) {
		return nil, errors.New("no available slots for the selected date range")
	}

	totalDays := utils.DaysBetween(req.CheckInDate, req.CheckOutDate)
	totalAmount := float64(totalDays) * pkg.PricePerDay

	orderNo := utils.GenerateOrderNo("RES")

	reservation := &models.Reservation{
		ID:            uuid.New(),
		OrderNo:       orderNo,
		OwnerID:       ownerID,
		PetID:         petID,
		StoreID:       storeID,
		PackageID:     pkgID,
		PackageType:   pkg.Type,
		CheckInDate:   req.CheckInDate,
		CheckOutDate:  req.CheckOutDate,
		TotalDays:     totalDays,
		TotalAmount:   totalAmount,
		Status:        models.ReservationStatusPending,
		Remark:        req.Remark,
	}

	if req.KeeperID != "" {
		keeperID, err := uuid.Parse(req.KeeperID)
		if err != nil {
			return nil, errors.New("invalid keeper ID")
		}
		reservation.KeeperID = &keeperID
	}

	tx := s.orderRepo.BeginTransaction()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	if err := s.resRepo.Create(reservation); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	depositAmount := totalAmount * 0.3
	if depositAmount > 0 {
		amountHash := utils.HashAmount(depositAmount, reservation.ID.String())
		order := &models.Order{
			ID:            uuid.New(),
			OrderNo:       utils.GenerateOrderNo("ORD"),
			ReservationID: reservation.ID,
			OwnerID:       ownerID,
			StoreID:       storeID,
			Type:          models.OrderTypePrepay,
			Amount:        depositAmount,
			PayStatus:     models.PayStatusUnpaid,
			AmountHash:    amountHash,
		}
		if err := s.orderRepo.Create(order, tx); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create order: %w", err)
		}

		reservation.DepositAmount = depositAmount
		if err := s.resRepo.Update(reservation); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update reservation: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Infof("Reservation created: %s (owner: %s, pet: %s)", reservation.OrderNo, ownerID, petID)
	return reservation, nil
}

func (s *ReservationService) GetByID(id uuid.UUID) (*models.Reservation, error) {
	res, err := s.resRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}
	if res == nil {
		return nil, errors.New("reservation not found")
	}
	return res, nil
}

func (s *ReservationService) List(query dto.StatisticsQuery, page, pageSize int, ownerID, storeID *uuid.UUID, status string) ([]models.Reservation, int64, error) {
	return s.resRepo.List(query, page, pageSize, ownerID, storeID, status)
}

func (s *ReservationService) Confirm(id, storeID uuid.UUID, req *dto.ReservationConfirmRequest) error {
	res, err := s.resRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get reservation: %w", err)
	}
	if res == nil {
		return errors.New("reservation not found")
	}
	if res.StoreID != storeID {
		return errors.New("not authorized to confirm this reservation")
	}
	if res.Status != models.ReservationStatusPending {
		return errors.New("reservation is not in pending status")
	}

	if req.Status == models.ReservationStatusConfirmed {
		res.Status = models.ReservationStatusConfirmed
		s.logger.Infof("Reservation confirmed: %s", res.OrderNo)
	} else if req.Status == models.ReservationStatusCancelled {
		res.Status = models.ReservationStatusCancelled
		res.CancelReason = req.Reason
		now := time.Now()
		res.CancelledAt = &now
		s.logger.Infof("Reservation cancelled: %s (reason: %s)", res.OrderNo, req.Reason)
	} else {
		return errors.New("invalid status")
	}

	return s.resRepo.Update(res)
}

func (s *ReservationService) CheckIn(id, storeID uuid.UUID, req *dto.CheckInOutRequest) error {
	res, err := s.resRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get reservation: %w", err)
	}
	if res == nil {
		return errors.New("reservation not found")
	}
	if res.StoreID != storeID {
		return errors.New("not authorized")
	}
	if res.Status != models.ReservationStatusConfirmed {
		return errors.New("reservation is not in confirmed status")
	}

	res.Status = models.ReservationStatusCheckedIn
	if req.Remark != "" {
		res.Remark = res.Remark + "\n[Check-in]: " + req.Remark
	}

	s.logger.Infof("Pet checked in: %s (pet: %s)", res.OrderNo, res.PetID)
	return s.resRepo.Update(res)
}

func (s *ReservationService) CheckOut(id, storeID uuid.UUID, req *dto.CheckInOutRequest) error {
	res, err := s.resRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get reservation: %w", err)
	}
	if res == nil {
		return errors.New("reservation not found")
	}
	if res.StoreID != storeID {
		return errors.New("not authorized")
	}
	if res.Status != models.ReservationStatusCheckedIn {
		return errors.New("reservation is not in checked-in status")
	}

	res.Status = models.ReservationStatusCompleted
	now := time.Now()
	res.CompletedAt = &now
	if req.Remark != "" {
		res.Remark = res.Remark + "\n[Check-out]: " + req.Remark
	}

	settlementAmount := res.TotalAmount - res.DepositAmount
	if settlementAmount > 0 {
		amountHash := utils.HashAmount(settlementAmount, res.ID.String())
		order := &models.Order{
			ID:            uuid.New(),
			OrderNo:       utils.GenerateOrderNo("ORD"),
			ReservationID: res.ID,
			OwnerID:       res.OwnerID,
			StoreID:       res.StoreID,
			Type:          models.OrderTypeSettlement,
			Amount:        settlementAmount,
			PayStatus:     models.PayStatusUnpaid,
			AmountHash:    amountHash,
		}
		if err := s.orderRepo.Create(order, nil); err != nil {
			s.logger.Errorf("Failed to create settlement order: %v", err)
		}
	}

	s.logger.Infof("Pet checked out: %s (pet: %s)", res.OrderNo, res.PetID)
	return s.resRepo.Update(res)
}

func (s *ReservationService) Cancel(id, ownerID uuid.UUID, reason string) error {
	res, err := s.resRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get reservation: %w", err)
	}
	if res == nil {
		return errors.New("reservation not found")
	}
	if res.OwnerID != ownerID {
		return errors.New("not authorized to cancel this reservation")
	}
	if res.Status != models.ReservationStatusPending && res.Status != models.ReservationStatusConfirmed {
		return errors.New("reservation cannot be cancelled in current status")
	}

	res.Status = models.ReservationStatusCancelled
	res.CancelReason = reason
	now := time.Now()
	res.CancelledAt = &now

	s.logger.Infof("Reservation cancelled by owner: %s (reason: %s)", res.OrderNo, reason)
	return s.resRepo.Update(res)
}

type DailyRecordService struct {
	recordRepo *repository.DailyRecordRepository
	resRepo    *repository.ReservationRepository
	logger     *logrus.Logger
}

func NewDailyRecordService(
	recordRepo *repository.DailyRecordRepository,
	resRepo *repository.ReservationRepository,
	logger *logrus.Logger,
) *DailyRecordService {
	return &DailyRecordService{
		recordRepo: recordRepo,
		resRepo:    resRepo,
		logger:     logger,
	}
}

func (s *DailyRecordService) Create(keeperID uuid.UUID, req *dto.DailyRecordRequest) (*models.DailyRecord, error) {
	resID, err := uuid.Parse(req.ReservationID)
	if err != nil {
		return nil, errors.New("invalid reservation ID")
	}
	petID, err := uuid.Parse(req.PetID)
	if err != nil {
		return nil, errors.New("invalid pet ID")
	}

	res, err := s.resRepo.GetByID(resID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}
	if res == nil {
		return nil, errors.New("reservation not found")
	}
	if res.Status != models.ReservationStatusCheckedIn {
		return nil, errors.New("can only create records for checked-in reservations")
	}

	record := &models.DailyRecord{
		ID:            uuid.New(),
		ReservationID: resID,
		PetID:         petID,
		KeeperID:      keeperID,
		RecordDate:    req.RecordDate,
		FeedStatus:    req.FeedStatus,
		Activity:      req.Activity,
		HealthStatus:  req.HealthStatus,
		Mood:          req.Mood,
		Photos:        req.Photos,
		Remark:        req.Remark,
	}

	if err := s.recordRepo.Create(record); err != nil {
		return nil, fmt.Errorf("failed to create daily record: %w", err)
	}

	s.logger.Infof("Daily record created: reservation %s, pet %s", resID, petID)
	return record, nil
}

func (s *DailyRecordService) GetByID(id uuid.UUID) (*models.DailyRecord, error) {
	record, err := s.recordRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get record: %w", err)
	}
	if record == nil {
		return nil, errors.New("record not found")
	}
	return record, nil
}

func (s *DailyRecordService) ListByReservation(reservationID uuid.UUID, page, pageSize int) ([]models.DailyRecord, int64, error) {
	return s.recordRepo.ListByReservation(reservationID, page, pageSize)
}

func (s *DailyRecordService) ListByPet(petID uuid.UUID, page, pageSize int) ([]models.DailyRecord, int64, error) {
	return s.recordRepo.ListByPet(petID, page, pageSize)
}

func (s *DailyRecordService) Update(id, keeperID uuid.UUID, req *dto.DailyRecordRequest) error {
	record, err := s.recordRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get record: %w", err)
	}
	if record == nil {
		return errors.New("record not found")
	}
	if record.KeeperID != keeperID {
		return errors.New("not authorized to update this record")
	}

	record.FeedStatus = req.FeedStatus
	record.Activity = req.Activity
	record.HealthStatus = req.HealthStatus
	record.Mood = req.Mood
	record.Photos = req.Photos
	record.Remark = req.Remark

	return s.recordRepo.Update(record)
}

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	resRepo    *repository.ReservationRepository
	userRepo   *repository.UserRepository
	logger     *logrus.Logger
}

func NewReviewService(
	reviewRepo *repository.ReviewRepository,
	resRepo *repository.ReservationRepository,
	userRepo *repository.UserRepository,
	logger *logrus.Logger,
) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		resRepo:    resRepo,
		userRepo:   userRepo,
		logger:     logger,
	}
}

func (s *ReviewService) Create(ownerID uuid.UUID, req *dto.ReviewRequest) (*models.Review, error) {
	resID, err := uuid.Parse(req.ReservationID)
	if err != nil {
		return nil, errors.New("invalid reservation ID")
	}

	res, err := s.resRepo.GetByID(resID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}
	if res == nil {
		return nil, errors.New("reservation not found")
	}
	if res.OwnerID != ownerID {
		return nil, errors.New("not authorized to review this reservation")
	}
	if res.Status != models.ReservationStatusCompleted {
		return nil, errors.New("can only review completed reservations")
	}

	existingReview, err := s.reviewRepo.GetByReservationID(resID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing review: %w", err)
	}
	if existingReview != nil {
		return nil, errors.New("review already exists for this reservation")
	}

	review := &models.Review{
		ID:            uuid.New(),
		ReservationID: resID,
		OwnerID:       ownerID,
		StoreID:       res.StoreID,
		StoreRating:   req.StoreRating,
		KeeperRating:  req.KeeperRating,
		Content:       req.Content,
		Images:        req.Images,
	}

	if res.KeeperID != nil {
		review.KeeperID = res.KeeperID
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	storeRating, storeCount, _ := s.reviewRepo.GetStoreRating(res.StoreID)
	storeInfo, _ := s.userRepo.GetStoreInfoByUserID(res.StoreID)
	if storeInfo != nil {
		storeInfo.Rating = storeRating
		storeInfo.ReviewCount = storeCount
		s.userRepo.UpdateStoreInfo(storeInfo)
	}

	if res.KeeperID != nil {
		keeperRating, keeperCount, _ := s.reviewRepo.GetKeeperRating(*res.KeeperID)
		keeperInfo, _ := s.userRepo.GetKeeperInfoByUserID(*res.KeeperID)
		if keeperInfo != nil {
			keeperInfo.Rating = keeperRating
			keeperInfo.ReviewCount = keeperCount
			s.userRepo.UpdateKeeperInfo(keeperInfo)
		}
	}

	s.logger.Infof("Review created: reservation %s, store rating %d", resID, req.StoreRating)
	return review, nil
}

func (s *ReviewService) GetByID(id uuid.UUID) (*models.Review, error) {
	review, err := s.reviewRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get review: %w", err)
	}
	if review == nil {
		return nil, errors.New("review not found")
	}
	return review, nil
}

func (s *ReviewService) ListByStore(storeID uuid.UUID, page, pageSize int) ([]models.Review, int64, error) {
	return s.reviewRepo.ListByStore(storeID, page, pageSize)
}

func (s *ReviewService) ListByKeeper(keeperID uuid.UUID, page, pageSize int) ([]models.Review, int64, error) {
	return s.reviewRepo.ListByKeeper(keeperID, page, pageSize)
}

func (s *ReviewService) Reply(id, storeID uuid.UUID, req *dto.ReviewReplyRequest) error {
	review, err := s.reviewRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get review: %w", err)
	}
	if review == nil {
		return errors.New("review not found")
	}
	if review.StoreID != storeID {
		return errors.New("not authorized to reply to this review")
	}

	review.Reply = req.Reply
	now := time.Now()
	review.ReplyAt = &now

	return s.reviewRepo.Update(review)
}

type OrderService struct {
	orderRepo *repository.OrderRepository
	resRepo   *repository.ReservationRepository
	logger    *logrus.Logger
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	resRepo *repository.ReservationRepository,
	logger *logrus.Logger,
) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		resRepo:   resRepo,
		logger:    logger,
	}
}

func (s *OrderService) Pay(ownerID uuid.UUID, req *dto.OrderRequest) (*models.Order, error) {
	resID, err := uuid.Parse(req.ReservationID)
	if err != nil {
		return nil, errors.New("invalid reservation ID")
	}

	res, err := s.resRepo.GetByID(resID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}
	if res == nil {
		return nil, errors.New("reservation not found")
	}
	if res.OwnerID != ownerID {
		return nil, errors.New("not authorized to pay for this reservation")
	}

	expectedHash := utils.HashAmount(req.Amount, res.ID.String())
	if req.AmountHash != expectedHash {
		s.logger.Errorf("Amount tampering detected: expected %s, got %s", expectedHash, req.AmountHash)
		return nil, errors.New("amount verification failed, possible tampering")
	}

	orders, err := s.orderRepo.GetByReservationID(resID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	var targetOrder *models.Order
	for _, o := range orders {
		if o.PayStatus == models.PayStatusUnpaid && o.Amount == req.Amount {
			targetOrder = &o
			break
		}
	}

	if targetOrder == nil {
		targetOrder = &models.Order{
			ID:            uuid.New(),
			OrderNo:       utils.GenerateOrderNo("ORD"),
			ReservationID: resID,
			OwnerID:       ownerID,
			StoreID:       res.StoreID,
			Type:          models.OrderTypePrepay,
			Amount:        req.Amount,
			AmountHash:    req.AmountHash,
		}
	}

	targetOrder.PayStatus = models.PayStatusPaid
	targetOrder.PayMethod = req.PayMethod
	now := time.Now()
	targetOrder.PaidAt = &now
	targetOrder.TransactionID = fmt.Sprintf("TX_%s", utils.GenerateOrderNo(""))
	targetOrder.Remark = req.Remark

	if err := s.orderRepo.Update(targetOrder, nil); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	s.logger.Infof("Order paid: %s, amount: %.2f, method: %s", targetOrder.OrderNo, req.Amount, req.PayMethod)
	return targetOrder, nil
}

func (s *OrderService) Settle(storeID uuid.UUID, req *dto.OrderSettlementRequest) (*models.Order, error) {
	resID, err := uuid.Parse(req.ReservationID)
	if err != nil {
		return nil, errors.New("invalid reservation ID")
	}

	res, err := s.resRepo.GetByID(resID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}
	if res == nil {
		return nil, errors.New("reservation not found")
	}
	if res.StoreID != storeID {
		return nil, errors.New("not authorized")
	}

	expectedHash := utils.HashAmount(req.Amount, res.ID.String())
	if req.AmountHash != expectedHash {
		s.logger.Errorf("Settlement amount tampering detected: expected %s, got %s", expectedHash, req.AmountHash)
		return nil, errors.New("amount verification failed, possible tampering")
	}

	order := &models.Order{
		ID:            uuid.New(),
		OrderNo:       utils.GenerateOrderNo("ORD"),
		ReservationID: resID,
		OwnerID:       res.OwnerID,
		StoreID:       storeID,
		Type:          models.OrderTypeSettlement,
		Amount:        req.Amount,
		PayStatus:     models.PayStatusPaid,
		PayMethod:     req.PayMethod,
		PaidAt:        &time.Time{},
		TransactionID: fmt.Sprintf("TX_%s", utils.GenerateOrderNo("")),
		AmountHash:    req.AmountHash,
		Remark:        req.Remark,
	}
	now := time.Now()
	order.PaidAt = &now

	if err := s.orderRepo.Create(order, nil); err != nil {
		return nil, fmt.Errorf("failed to create settlement order: %w", err)
	}

	s.logger.Infof("Settlement order created: %s, amount: %.2f", order.OrderNo, req.Amount)
	return order, nil
}

func (s *OrderService) Refund(orderID, storeID uuid.UUID, amount float64, reason string) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return errors.New("order not found")
	}
	if order.StoreID != storeID {
		return errors.New("not authorized")
	}
	if order.PayStatus != models.PayStatusPaid {
		return errors.New("can only refund paid orders")
	}
	if amount > order.Amount {
		return errors.New("refund amount exceeds order amount")
	}

	order.RefundAmount = amount
	order.PayStatus = models.PayStatusRefund
	now := time.Now()
	order.RefundAt = &now
	order.Remark = order.Remark + " Refund reason: " + reason

	return s.orderRepo.Update(order, nil)
}

func (s *OrderService) GetByID(id uuid.UUID) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *OrderService) List(ownerID, storeID *uuid.UUID, payStatus, orderType string, query dto.StatisticsQuery, page, pageSize int) ([]models.Order, int64, error) {
	return s.orderRepo.List(ownerID, storeID, payStatus, orderType, query, page, pageSize)
}

func (s *OrderService) GetByReservationID(reservationID uuid.UUID) ([]models.Order, error) {
	return s.orderRepo.GetByReservationID(reservationID)
}

type HealthAlertService struct {
	alertRepo *repository.HealthAlertRepository
	petRepo   *repository.PetRepository
	logger    *logrus.Logger
}

func NewHealthAlertService(
	alertRepo *repository.HealthAlertRepository,
	petRepo *repository.PetRepository,
	logger *logrus.Logger,
) *HealthAlertService {
	return &HealthAlertService{
		alertRepo: alertRepo,
		petRepo:   petRepo,
		logger:    logger,
	}
}

func (s *HealthAlertService) CheckAndCreateAlerts() error {
	vaccines, err := s.alertRepo.GetExpiringVaccines(30)
	if err != nil {
		return fmt.Errorf("failed to check expiring vaccines: %w", err)
	}

	for _, v := range vaccines {
		pet, _ := s.petRepo.GetByID(v.PetID)
		if pet == nil {
			continue
		}

		alert := &models.HealthAlert{
			ID:        uuid.New(),
			UserID:    pet.OwnerID,
			PetID:     v.PetID,
			AlertType: models.AlertTypeVaccineExpire,
			Title:     fmt.Sprintf("疫苗到期提醒 - %s", v.VaccineName),
			Content:   fmt.Sprintf("宠物%s的%s疫苗将于%s到期，请及时接种。", pet.Name, v.VaccineName, v.ExpireAt.Format("2006-01-02")),
			RecordID:  v.ID,
			ExpireAt:  v.ExpireAt,
		}
		if err := s.alertRepo.Create(alert); err != nil {
			s.logger.Errorf("Failed to create vaccine alert: %v", err)
		}
	}

	deworms, err := s.alertRepo.GetExpiringDeworms(30)
	if err != nil {
		return fmt.Errorf("failed to check expiring deworms: %w", err)
	}

	for _, d := range deworms {
		pet, _ := s.petRepo.GetByID(d.PetID)
		if pet == nil {
			continue
		}

		alert := &models.HealthAlert{
			ID:        uuid.New(),
			UserID:    pet.OwnerID,
			PetID:     d.PetID,
			AlertType: models.AlertTypeDewormExpire,
			Title:     fmt.Sprintf("驱虫到期提醒 - %s", d.DewormType),
			Content:   fmt.Sprintf("宠物%s的%s驱虫将于%s到期，请及时驱虫。", pet.Name, d.DewormType, d.ExpireAt.Format("2006-01-02")),
			RecordID:  d.ID,
			ExpireAt:  d.ExpireAt,
		}
		if err := s.alertRepo.Create(alert); err != nil {
			s.logger.Errorf("Failed to create deworm alert: %v", err)
		}
	}

	s.logger.Infof("Health alerts checked: %d vaccine alerts, %d deworm alerts", len(vaccines), len(deworms))
	return nil
}

func (s *HealthAlertService) ListByUser(userID uuid.UUID, isRead *bool, page, pageSize int) ([]models.HealthAlert, int64, error) {
	return s.alertRepo.ListByUser(userID, isRead, page, pageSize)
}

func (s *HealthAlertService) MarkAsRead(id, userID uuid.UUID) error {
	alert, err := s.alertRepo.ListByUser(userID, nil, 1, 100)
	if err != nil {
		return err
	}
	for _, a := range alert {
		if a.ID == id {
			return s.alertRepo.MarkAsRead(id)
		}
	}
	return errors.New("alert not found or not authorized")
}

func (s *HealthAlertService) MarkAllAsRead(userID uuid.UUID) error {
	return s.alertRepo.MarkAllAsRead(userID)
}

type StatisticsService struct {
	statRepo *repository.StatisticsRepository
	logger   *logrus.Logger
}

func NewStatisticsService(statRepo *repository.StatisticsRepository, logger *logrus.Logger) *StatisticsService {
	return &StatisticsService{
		statRepo: statRepo,
		logger:   logger,
	}
}

func (s *StatisticsService) GetRevenueTrend(storeID uuid.UUID, start, end time.Time) ([]map[string]interface{}, error) {
	return s.statRepo.GetRevenueTrend(storeID, start, end)
}

func (s *StatisticsService) GetOccupancyRate(storeID uuid.UUID, start, end time.Time) (float64, error) {
	return s.statRepo.GetOccupancyRate(storeID, start, end)
}

func (s *StatisticsService) GetPetTypeDistribution(storeID uuid.UUID, start, end time.Time) ([]map[string]interface{}, error) {
	return s.statRepo.GetPetTypeDistribution(storeID, start, end)
}

func (s *StatisticsService) GetOrderStatistics(storeID uuid.UUID, start, end time.Time) (map[string]interface{}, error) {
	return s.statRepo.GetOrderStatistics(storeID, start, end)
}

type OperationLogService struct {
	logRepo *repository.OperationLogRepository
	logger  *logrus.Logger
}

func NewOperationLogService(logRepo *repository.OperationLogRepository, logger *logrus.Logger) *OperationLogService {
	return &OperationLogService{
		logRepo: logRepo,
		logger:  logger,
	}
}

func (s *OperationLogService) Create(log *models.OperationLog) error {
	return s.logRepo.Create(log)
}

func (s *OperationLogService) List(page, pageSize int, userID *uuid.UUID) ([]models.OperationLog, int64, error) {
	return s.logRepo.List(page, pageSize, userID)
}
