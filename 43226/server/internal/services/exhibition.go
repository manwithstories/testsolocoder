package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"museum-server/internal/dto"
	"museum-server/internal/models"
	"museum-server/internal/repository"
	redisPkg "museum-server/pkg/redis"
)

type ExhibitionService struct {
	exhibitionRepo *repository.ExhibitionRepository
}

func NewExhibitionService(exhibitionRepo *repository.ExhibitionRepository) *ExhibitionService {
	return &ExhibitionService{exhibitionRepo: exhibitionRepo}
}

func (s *ExhibitionService) Create(museumID uint, req *dto.ExhibitionRequest) (*models.Exhibition, error) {
	exhibition := &models.Exhibition{
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Location:    req.Location,
		HallNumber:  req.HallNumber,
		TicketPrice: req.TicketPrice,
		MaxVisitors: req.MaxVisitors,
		ImageUrl:    req.ImageUrl,
		Status:      models.ExhibitionStatusDraft,
		IsVirtual:   req.IsVirtual,
		VirtualUrl:  req.VirtualUrl,
		MuseumID:    museumID,
	}

	if req.Status != "" {
		exhibition.Status = req.Status
	}

	if err := s.exhibitionRepo.Create(exhibition); err != nil {
		return nil, fmt.Errorf("failed to create exhibition: %w", err)
	}

	if len(req.CollectionIDs) > 0 {
		sortOrders := make(map[uint]int)
		for i, cid := range req.CollectionIDs {
			sortOrders[cid] = i
		}
		s.exhibitionRepo.AddCollections(exhibition.ID, req.CollectionIDs, sortOrders)
	}

	return exhibition, nil
}

func (s *ExhibitionService) GetByID(id uint) (*models.Exhibition, error) {
	exhibition, err := s.exhibitionRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("exhibition not found")
	}

	s.exhibitionRepo.IncrementViewCount(id)

	s.cacheHotExhibition(exhibition)

	return exhibition, nil
}

func (s *ExhibitionService) Update(id uint, req *dto.ExhibitionRequest) error {
	exhibition, err := s.exhibitionRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("exhibition not found")
	}

	exhibition.Title = req.Title
	exhibition.Description = req.Description
	exhibition.StartDate = req.StartDate
	exhibition.EndDate = req.EndDate
	exhibition.Location = req.Location
	exhibition.HallNumber = req.HallNumber
	exhibition.TicketPrice = req.TicketPrice
	exhibition.MaxVisitors = req.MaxVisitors
	if req.ImageUrl != "" {
		exhibition.ImageUrl = req.ImageUrl
	}
	if req.Status != "" {
		exhibition.Status = req.Status
	}
	exhibition.IsVirtual = req.IsVirtual
	exhibition.VirtualUrl = req.VirtualUrl

	return s.exhibitionRepo.Update(exhibition)
}

func (s *ExhibitionService) Delete(id uint) error {
	return s.exhibitionRepo.Delete(id)
}

func (s *ExhibitionService) List(query *dto.ExhibitionListQuery) ([]models.Exhibition, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}
	return s.exhibitionRepo.List(query)
}

func (s *ExhibitionService) AddCollections(id uint, collectionIDs []uint) error {
	sortOrders := make(map[uint]int)
	for i, cid := range collectionIDs {
		sortOrders[cid] = i
	}
	return s.exhibitionRepo.AddCollections(id, collectionIDs, sortOrders)
}

func (s *ExhibitionService) RemoveCollections(id uint, collectionIDs []uint) error {
	return s.exhibitionRepo.RemoveCollections(id, collectionIDs)
}

func (s *ExhibitionService) GetCollections(id uint) ([]models.Collection, error) {
	return s.exhibitionRepo.GetCollections(id)
}

func (s *ExhibitionService) CreateTimeSlot(req *dto.TimeSlotRequest) (*models.TimeSlot, error) {
	slot := &models.TimeSlot{
		ExhibitionID: 0,
		Date:         req.Date,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		MaxCapacity:  req.MaxCapacity,
	}

	if err := s.exhibitionRepo.CreateTimeSlot(slot); err != nil {
		return nil, fmt.Errorf("failed to create time slot: %w", err)
	}

	return slot, nil
}

func (s *ExhibitionService) BatchCreateTimeSlots(req *dto.BatchTimeSlotRequest) error {
	_, err := s.exhibitionRepo.FindByID(req.ExhibitionID)
	if err != nil {
		return fmt.Errorf("exhibition not found")
	}

	exhibition, _ := s.exhibitionRepo.FindByID(req.ExhibitionID)
	for date := req.StartDate; !date.After(req.EndDate); date = date.AddDate(0, 0, 1) {
		startHour, startMin := parseTime(req.StartTime)
		endHour, endMin := parseTime(req.EndTime)

		for {
			slotStart := fmt.Sprintf("%02d:%02d", startHour, startMin)
			slotEndHour, slotEndMin := startHour, startMin+req.Interval
			if slotEndMin >= 60 {
				slotEndHour += slotEndMin / 60
				slotEndMin %= 60
			}
			slotEnd := fmt.Sprintf("%02d:%02d", slotEndHour, slotEndMin)

			if slotEndHour > endHour || (slotEndHour == endHour && slotEndMin > endMin) {
				break
			}

			slot := &models.TimeSlot{
				ExhibitionID: req.ExhibitionID,
				Date:         date,
				StartTime:    slotStart,
				EndTime:      slotEnd,
				MaxCapacity:  req.MaxCapacity,
			}
			s.exhibitionRepo.CreateTimeSlot(slot)

			startHour, startMin = slotEndHour, slotEndMin
			if startMin >= 60 {
				startHour += startMin / 60
				startMin %= 60
			}
		}

		_ = exhibition
	}

	return nil
}

func parseTime(t string) (int, int) {
	var h, m int
	fmt.Sscanf(t, "%d:%d", &h, &m)
	return h, m
}

func (s *ExhibitionService) ListTimeSlots(exhibitionID uint, date time.Time) ([]models.TimeSlot, error) {
	return s.exhibitionRepo.ListTimeSlots(exhibitionID, date)
}

func (s *ExhibitionService) GetHotExhibitions() ([]models.Exhibition, error) {
	ctx := context.Background()
	client := redisPkg.GetClient()

	val, err := client.Get(ctx, "hot_exhibitions").Result()
	if err == nil && val != "" {
		var exhibitions []models.Exhibition
		if json.Unmarshal([]byte(val), &exhibitions) == nil {
			return exhibitions, nil
		}
	}

	query := &dto.ExhibitionListQuery{
		Page:     1,
		PageSize: 5,
		Status:   models.ExhibitionStatusPublished,
	}
	exhibitions, _, err := s.exhibitionRepo.List(query)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(exhibitions)
	client.Set(ctx, "hot_exhibitions", data, 5*time.Minute)

	return exhibitions, nil
}

func (s *ExhibitionService) cacheHotExhibition(exhibition *models.Exhibition) {
	ctx := context.Background()
	client := redisPkg.GetClient()
	key := fmt.Sprintf("exhibition:%d", exhibition.ID)
	data, _ := json.Marshal(exhibition)
	client.Set(ctx, key, data, 10*time.Minute)
}
