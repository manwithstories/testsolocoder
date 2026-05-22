package services

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"meeting-room/internal/config"
	"meeting-room/internal/models"
	"meeting-room/internal/repositories"
	"meeting-room/internal/utils"

	"github.com/google/uuid"
)

type RoomService struct {
	roomRepo *repositories.RoomRepository
}

func NewRoomService() *RoomService {
	return &RoomService{
		roomRepo: repositories.NewRoomRepository(),
	}
}

type CreateRoomRequest struct {
	Name         string   `json:"name" binding:"required"`
	Floor        string   `json:"floor" binding:"required"`
	Capacity     int      `json:"capacity" binding:"required,min=1"`
	Location     string   `json:"location"`
	PricePerHour float64  `json:"price_per_hour" binding:"required,min=0"`
	Equipment    []string `json:"equipment"`
	Description  string   `json:"description"`
	AvailableStart string `json:"available_start"`
	AvailableEnd   string `json:"available_end"`
}

type UpdateRoomRequest struct {
	Name           string   `json:"name"`
	Floor          string   `json:"floor"`
	Capacity       int      `json:"capacity"`
	Location       string   `json:"location"`
	PricePerHour   float64  `json:"price_per_hour"`
	Equipment      []string `json:"equipment"`
	Description    string   `json:"description"`
	Status         int      `json:"status"`
	AvailableStart string   `json:"available_start"`
	AvailableEnd   string   `json:"available_end"`
}

func (s *RoomService) CreateRoom(req *CreateRoomRequest, userID uint) (*models.Room, error) {
	room := &models.Room{
		Name:           req.Name,
		Floor:          req.Floor,
		Capacity:       req.Capacity,
		Location:       req.Location,
		PricePerHour:   req.PricePerHour,
		Description:    req.Description,
		AvailableStart: req.AvailableStart,
		AvailableEnd:   req.AvailableEnd,
		Status:         models.RoomStatusActive,
		CreatedBy:      userID,
	}
	room.SetEquipmentList(req.Equipment)

	if room.AvailableStart == "" {
		room.AvailableStart = "08:00"
	}
	if room.AvailableEnd == "" {
		room.AvailableEnd = "22:00"
	}

	err := s.roomRepo.Create(room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *RoomService) GetRoom(id uint) (*models.Room, error) {
	room, err := s.roomRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("会议室不存在")
	}
	return room, nil
}

func (s *RoomService) UpdateRoom(id uint, req *UpdateRoomRequest) (*models.Room, error) {
	room, err := s.roomRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("会议室不存在")
	}

	if req.Name != "" {
		room.Name = req.Name
	}
	if req.Floor != "" {
		room.Floor = req.Floor
	}
	if req.Capacity > 0 {
		room.Capacity = req.Capacity
	}
	if req.Location != "" {
		room.Location = req.Location
	}
	if req.PricePerHour > 0 {
		room.PricePerHour = req.PricePerHour
	}
	if req.Description != "" {
		room.Description = req.Description
	}
	if req.Status >= 0 {
		room.Status = models.RoomStatus(req.Status)
	}
	if req.AvailableStart != "" {
		room.AvailableStart = req.AvailableStart
	}
	if req.AvailableEnd != "" {
		room.AvailableEnd = req.AvailableEnd
	}
	if len(req.Equipment) > 0 {
		room.SetEquipmentList(req.Equipment)
	}

	err = s.roomRepo.Update(room)
	return room, err
}

func (s *RoomService) DeleteRoom(id uint) error {
	return s.roomRepo.Delete(id)
}

func (s *RoomService) ListRooms(page, pageSize int, floor string, equipment string) ([]models.Room, int64, error) {
	return s.roomRepo.List(page, pageSize, floor, equipment)
}

func (s *RoomService) ListAllRooms() ([]models.Room, error) {
	return s.roomRepo.ListAll()
}

func (s *RoomService) GetAllFloors() ([]string, error) {
	return s.roomRepo.GetAllFloors()
}

func (s *RoomService) UploadPhoto(roomID uint, file io.Reader, filename string) (*models.RoomPhoto, error) {
	_, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return nil, errors.New("会议室不存在")
	}

	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	uploadDir := config.Cfg.Upload.Dir + "/rooms/"
	os.MkdirAll(uploadDir, 0755)

	filepath := uploadDir + newFilename
	dst, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, err
	}

	photo := &models.RoomPhoto{
		RoomID: roomID,
		URL:    "/uploads/rooms/" + newFilename,
	}

	err = s.roomRepo.AddPhoto(photo)
	if err != nil {
		os.Remove(filepath)
		return nil, err
	}

	return photo, nil
}

func (s *RoomService) DeletePhoto(id uint) error {
	utils.Logger.Infof("删除会议室照片: %d", id)
	return s.roomRepo.DeletePhoto(id)
}

func (s *RoomService) GetRoomCacheKey(roomID uint) string {
	return fmt.Sprintf("room:%d", roomID)
}

func (s *RoomService) CacheRoom(room *models.Room) {
	key := s.GetRoomCacheKey(room.ID)
	utils.RedisSet(key, "1", 1*time.Hour)
}
