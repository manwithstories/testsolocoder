package service

import (
	"errors"
	"fmt"
	"hotel-system/internal/dto"
	"hotel-system/internal/model"
	"hotel-system/internal/pkg/logger"
	"hotel-system/internal/repository"

	"gorm.io/gorm"
)

type RoomService interface {
	CreateRoomType(req *dto.RoomTypeCreateRequest) (*model.RoomType, error)
	GetRoomType(id uint) (*model.RoomType, error)
	UpdateRoomType(id uint, req *dto.RoomTypeUpdateRequest) (*model.RoomType, error)
	DeleteRoomType(id uint) error
	ListRoomTypes() ([]model.RoomType, error)

	CreateRoom(req *dto.RoomCreateRequest) (*model.Room, error)
	GetRoom(id uint) (*model.Room, error)
	UpdateRoom(id uint, req *dto.RoomUpdateRequest) (*model.Room, error)
	DeleteRoom(id uint) error
	ListRooms(req *dto.RoomListRequest) ([]model.Room, int64, error)

	BatchImportRooms(req *dto.RoomBatchImportRequest) (*dto.RoomBatchImportResponse, error)
	GetAvailableRooms(roomTypeID uint) ([]model.Room, error)
	UpdateRoomStatus(id uint, status model.RoomStatus) error
	GetRoomDashboard() (*dto.RoomDashboardResponse, error)
}

type roomService struct {
	roomTypeRepo repository.RoomTypeRepository
	roomRepo     repository.RoomRepository
	db           *gorm.DB
}

func NewRoomService(roomTypeRepo repository.RoomTypeRepository, roomRepo repository.RoomRepository, db *gorm.DB) RoomService {
	return &roomService{
		roomTypeRepo: roomTypeRepo,
		roomRepo:     roomRepo,
		db:           db,
	}
}

func (s *roomService) CreateRoomType(req *dto.RoomTypeCreateRequest) (*model.RoomType, error) {
	existing, _ := s.roomTypeRepo.GetByName(req.Name)
	if existing != nil {
		logger.Warnf("房型名称已存在: %s", req.Name)
		return nil, errors.New("房型名称已存在")
	}

	roomType := &model.RoomType{
		Name:        req.Name,
		Description: req.Description,
		BasePrice:   req.BasePrice,
		BedCount:    req.BedCount,
		MaxGuests:   req.MaxGuests,
		Facilities:  model.StringArray(req.Facilities),
	}

	err := s.roomTypeRepo.Create(roomType)
	if err != nil {
		logger.Errorf("创建房型失败: %v", err)
		return nil, errors.New("创建房型失败")
	}

	logger.Infof("房型创建成功: %s", roomType.Name)
	return roomType, nil
}

func (s *roomService) GetRoomType(id uint) (*model.RoomType, error) {
	roomType, err := s.roomTypeRepo.GetByID(id)
	if err != nil {
		logger.Errorf("获取房型失败: id=%d, err=%v", id, err)
		return nil, errors.New("房型不存在")
	}
	return roomType, nil
}

func (s *roomService) UpdateRoomType(id uint, req *dto.RoomTypeUpdateRequest) (*model.RoomType, error) {
	roomType, err := s.roomTypeRepo.GetByID(id)
	if err != nil {
		logger.Errorf("更新房型失败，房型不存在: id=%d, err=%v", id, err)
		return nil, errors.New("房型不存在")
	}

	if req.Name != "" && req.Name != roomType.Name {
		existing, _ := s.roomTypeRepo.GetByName(req.Name)
		if existing != nil {
			logger.Warnf("房型名称已存在: %s", req.Name)
			return nil, errors.New("房型名称已存在")
		}
		roomType.Name = req.Name
	}

	if req.Description != "" {
		roomType.Description = req.Description
	}
	if req.BasePrice > 0 {
		roomType.BasePrice = req.BasePrice
	}
	if req.BedCount > 0 {
		roomType.BedCount = req.BedCount
	}
	if req.MaxGuests > 0 {
		roomType.MaxGuests = req.MaxGuests
	}
	if req.Facilities != nil {
		roomType.Facilities = model.StringArray(req.Facilities)
	}

	err = s.roomTypeRepo.Update(roomType)
	if err != nil {
		logger.Errorf("更新房型失败: %v", err)
		return nil, errors.New("更新房型失败")
	}

	logger.Infof("房型更新成功: id=%d, name=%s", roomType.ID, roomType.Name)
	return roomType, nil
}

func (s *roomService) DeleteRoomType(id uint) error {
	_, err := s.roomTypeRepo.GetByID(id)
	if err != nil {
		logger.Errorf("删除房型失败，房型不存在: id=%d, err=%v", id, err)
		return errors.New("房型不存在")
	}

	hasRooms, err := s.roomTypeRepo.HasRooms(id)
	if err != nil {
		logger.Errorf("检查房型关联房间失败: %v", err)
		return errors.New("删除房型失败")
	}
	if hasRooms {
		logger.Warnf("房型下存在关联房间，无法删除: id=%d", id)
		return errors.New("该房型下存在关联房间，无法删除")
	}

	err = s.roomTypeRepo.Delete(id)
	if err != nil {
		logger.Errorf("删除房型失败: %v", err)
		return errors.New("删除房型失败")
	}

	logger.Infof("房型删除成功: id=%d", id)
	return nil
}

func (s *roomService) ListRoomTypes() ([]model.RoomType, error) {
	roomTypes, err := s.roomTypeRepo.List()
	if err != nil {
		logger.Errorf("获取房型列表失败: %v", err)
		return nil, errors.New("获取房型列表失败")
	}
	return roomTypes, nil
}

func (s *roomService) CreateRoom(req *dto.RoomCreateRequest) (*model.Room, error) {
	_, err := s.roomTypeRepo.GetByID(req.RoomTypeID)
	if err != nil {
		logger.Errorf("创建房间失败，房型不存在: room_type_id=%d, err=%v", req.RoomTypeID, err)
		return nil, errors.New("房型不存在")
	}

	existing, _ := s.roomRepo.GetByRoomNo(req.RoomNo)
	if existing != nil {
		logger.Warnf("房间号已存在: %s", req.RoomNo)
		return nil, errors.New("房间号已存在")
	}

	status := model.RoomStatusAvailable
	if req.Status != "" {
		status = req.Status
	}

	room := &model.Room{
		RoomNo:     req.RoomNo,
		Floor:      req.Floor,
		RoomTypeID: req.RoomTypeID,
		Status:     status,
		Price:      req.Price,
		Facilities: model.StringArray(req.Facilities),
	}

	err = s.roomRepo.Create(room)
	if err != nil {
		logger.Errorf("创建房间失败: %v", err)
		return nil, errors.New("创建房间失败")
	}

	logger.Infof("房间创建成功: %s", room.RoomNo)
	return room, nil
}

func (s *roomService) GetRoom(id uint) (*model.Room, error) {
	room, err := s.roomRepo.GetByID(id)
	if err != nil {
		logger.Errorf("获取房间失败: id=%d, err=%v", id, err)
		return nil, errors.New("房间不存在")
	}
	return room, nil
}

func (s *roomService) UpdateRoom(id uint, req *dto.RoomUpdateRequest) (*model.Room, error) {
	room, err := s.roomRepo.GetByID(id)
	if err != nil {
		logger.Errorf("更新房间失败，房间不存在: id=%d, err=%v", id, err)
		return nil, errors.New("房间不存在")
	}

	if req.RoomNo != "" && req.RoomNo != room.RoomNo {
		existing, _ := s.roomRepo.GetByRoomNo(req.RoomNo)
		if existing != nil {
			logger.Warnf("房间号已存在: %s", req.RoomNo)
			return nil, errors.New("房间号已存在")
		}
		room.RoomNo = req.RoomNo
	}

	if req.RoomTypeID > 0 && req.RoomTypeID != room.RoomTypeID {
		_, err := s.roomTypeRepo.GetByID(req.RoomTypeID)
		if err != nil {
			logger.Errorf("更新房间失败，房型不存在: room_type_id=%d, err=%v", req.RoomTypeID, err)
			return nil, errors.New("房型不存在")
		}
		room.RoomTypeID = req.RoomTypeID
	}

	if req.Floor > 0 {
		room.Floor = req.Floor
	}
	if req.Price > 0 {
		room.Price = req.Price
	}
	if req.Status != "" {
		room.Status = req.Status
	}
	if req.Facilities != nil {
		room.Facilities = model.StringArray(req.Facilities)
	}

	err = s.roomRepo.Update(room)
	if err != nil {
		logger.Errorf("更新房间失败: %v", err)
		return nil, errors.New("更新房间失败")
	}

	logger.Infof("房间更新成功: id=%d, room_no=%s", room.ID, room.RoomNo)
	return room, nil
}

func (s *roomService) DeleteRoom(id uint) error {
	_, err := s.roomRepo.GetByID(id)
	if err != nil {
		logger.Errorf("删除房间失败，房间不存在: id=%d, err=%v", id, err)
		return errors.New("房间不存在")
	}

	err = s.roomRepo.Delete(id)
	if err != nil {
		logger.Errorf("删除房间失败: %v", err)
		return errors.New("删除房间失败")
	}

	logger.Infof("房间删除成功: id=%d", id)
	return nil
}

func (s *roomService) ListRooms(req *dto.RoomListRequest) ([]model.Room, int64, error) {
	page := req.GetPage()
	pageSize := req.GetPageSize()

	rooms, total, err := s.roomRepo.List(page, pageSize, req.RoomNo, req.Floor, req.RoomTypeID, req.Status)
	if err != nil {
		logger.Errorf("获取房间列表失败: %v", err)
		return nil, 0, errors.New("获取房间列表失败")
	}

	return rooms, total, nil
}

func (s *roomService) BatchImportRooms(req *dto.RoomBatchImportRequest) (*dto.RoomBatchImportResponse, error) {
	response := &dto.RoomBatchImportResponse{}
	var rooms []model.Room
	var failReasons []string

	for i, roomReq := range req.Rooms {
		_, err := s.roomTypeRepo.GetByID(roomReq.RoomTypeID)
		if err != nil {
			failReasons = append(failReasons, fmt.Sprintf("第%d条数据：房型不存在, 房型ID=%d", i+1, roomReq.RoomTypeID))
			response.FailCount++
			continue
		}

		existing, _ := s.roomRepo.GetByRoomNo(roomReq.RoomNo)
		if existing != nil {
			failReasons = append(failReasons, fmt.Sprintf("第%d条数据：房间号已存在, 房间号=%s", i+1, roomReq.RoomNo))
			response.FailCount++
			continue
		}

		status := model.RoomStatusAvailable
		if roomReq.Status != "" {
			status = roomReq.Status
		}

		room := model.Room{
			RoomNo:     roomReq.RoomNo,
			Floor:      roomReq.Floor,
			RoomTypeID: roomReq.RoomTypeID,
			Status:     status,
			Price:      roomReq.Price,
			Facilities: model.StringArray(roomReq.Facilities),
		}
		rooms = append(rooms, room)
		response.SuccessCount++
	}

	if len(rooms) > 0 {
		err := s.db.Transaction(func(tx *gorm.DB) error {
			for _, room := range rooms {
				if err := tx.Create(&room).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			logger.Errorf("批量导入房间失败: %v", err)
			return nil, errors.New("批量导入房间失败")
		}
	}

	response.FailReasons = failReasons
	logger.Infof("批量导入房间完成: 成功%d条, 失败%d条", response.SuccessCount, response.FailCount)
	return response, nil
}

func (s *roomService) GetAvailableRooms(roomTypeID uint) ([]model.Room, error) {
	rooms, err := s.roomRepo.GetAvailableRooms(roomTypeID)
	if err != nil {
		logger.Errorf("获取可预订房间失败: %v", err)
		return nil, errors.New("获取可预订房间失败")
	}
	return rooms, nil
}

func (s *roomService) UpdateRoomStatus(id uint, status model.RoomStatus) error {
	_, err := s.roomRepo.GetByID(id)
	if err != nil {
		logger.Errorf("更新房间状态失败，房间不存在: id=%d, err=%v", id, err)
		return errors.New("房间不存在")
	}

	err = s.roomRepo.UpdateStatus(id, status)
	if err != nil {
		logger.Errorf("更新房间状态失败: %v", err)
		return errors.New("更新房间状态失败")
	}

	logger.Infof("房间状态更新成功: id=%d, status=%s", id, status)
	return nil
}

func (s *roomService) GetRoomDashboard() (*dto.RoomDashboardResponse, error) {
	total, err := s.roomRepo.CountAll()
	if err != nil {
		logger.Errorf("获取房间总数失败: %v", err)
		return nil, errors.New("获取房间统计失败")
	}

	counts, err := s.roomRepo.CountByStatus()
	if err != nil {
		logger.Errorf("获取房间状态统计失败: %v", err)
		return nil, errors.New("获取房间统计失败")
	}

	dashboard := &dto.RoomDashboardResponse{
		Total:       total,
		Available:   counts[string(model.RoomStatusAvailable)],
		Occupied:    counts[string(model.RoomStatusOccupied)],
		Reserved:    counts[string(model.RoomStatusReserved)],
		Maintenance: counts[string(model.RoomStatusMaintenance)],
	}

	return dashboard, nil
}
