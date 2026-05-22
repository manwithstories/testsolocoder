package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"ticket-system/internal/config"
	"ticket-system/internal/dto"
	"ticket-system/internal/logging"
	"ticket-system/internal/models"
	"ticket-system/internal/redis"
	"ticket-system/internal/repository"
)

type ShowService struct {
	showRepo *repository.ShowRepository
}

func NewShowService() *ShowService {
	return &ShowService{
		showRepo: repository.NewShowRepository(),
	}
}

func (s *ShowService) CreateShow(req *dto.ShowCreateRequest) (*models.Show, error) {
	show := &models.Show{
		Name:        req.Name,
		Description: req.Description,
		Poster:      req.Poster,
		Artist:      req.Artist,
		Duration:    req.Duration,
		Status:      req.Status,
		OrganizerID:  req.Organizer,
		Venue:       req.Venue,
		Address:     req.Address,
	}

	err := s.showRepo.Create(show)
	if err != nil {
		logging.Errorf("Failed to create show: %v", err)
		return nil, err
	}

	logging.Infof("Show created: %d, %s", show.ID, show.Name)
	return show, nil
}

func (s *ShowService) GetShow(id uint64) (*models.Show, error) {
	return s.showRepo.GetByID(id)
}

func (s *ShowService) ListShows(req *dto.ShowListRequest) (*dto.PaginatedResponse, error) {
	shows, total, err := s.showRepo.List(req.Page, req.PageSize, req.Status, req.Keyword)
	if err != nil {
		return nil, err
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	return &dto.PaginatedResponse{
		List: shows,
		Pagination: dto.Pagination{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    total,
			Pages:    pages,
		},
	}, nil
}

func (s *ShowService) UpdateShow(id uint64, req *dto.ShowUpdateRequest) (*models.Show, error) {
	show, err := s.showRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("演出不存在")
	}

	if req.Name != "" {
		show.Name = req.Name
	}
	if req.Description != "" {
		show.Description = req.Description
	}
	if req.Poster != "" {
		show.Poster = req.Poster
	}
	if req.Artist != "" {
		show.Artist = req.Artist
	}
	if req.Duration > 0 {
		show.Duration = req.Duration
	}
	if req.Status >= 0 {
		show.Status = req.Status
	}
	if req.Organizer != "" {
		show.OrganizerID = req.Organizer
	}
	if req.Venue != "" {
		show.Venue = req.Venue
	}
	if req.Address != "" {
		show.Address = req.Address
	}

	err = s.showRepo.Update(show)
	if err != nil {
		return nil, err
	}

	logging.Infof("Show updated: %d", id)
	return show, nil
}

func (s *ShowService) DeleteShow(id uint64) error {
	return s.showRepo.Delete(id)
}

func (s *ShowService) CreateSession(req *dto.SessionCreateRequest) (*models.Session, error) {
	_, err := s.showRepo.GetByID(req.ShowID)
	if err != nil {
		return nil, errors.New("演出不存在")
	}

	session := &models.Session{
		ShowID:    req.ShowID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Status:    req.Status,
	}

	err = s.showRepo.CreateSession(session)
	if err != nil {
		return nil, err
	}

	logging.Infof("Session created: %d for show %d", session.ID, req.ShowID)
	return session, nil
}

func (s *ShowService) GetSessions(showID uint64) ([]models.Session, error) {
	return s.showRepo.GetSessionsByShowID(showID)
}

func (s *ShowService) CreateSeatArea(req *dto.SeatAreaCreateRequest) (*models.SeatArea, error) {
	_, err := s.showRepo.GetSessionByID(req.SessionID)
	if err != nil {
		return nil, errors.New("场次不存在")
	}

	area := &models.SeatArea{
		SessionID: req.SessionID,
		Name:      req.Name,
		Color:     req.Color,
		Price:     req.Price,
		SortOrder: req.SortOrder,
	}

	err = s.showRepo.CreateSeatArea(area)
	if err != nil {
		return nil, err
	}

	logging.Infof("Seat area created: %d for session %d", area.ID, req.SessionID)
	return area, nil
}

func (s *ShowService) GetSeatAreas(sessionID uint64) ([]models.SeatArea, error) {
	return s.showRepo.GetSeatAreasBySessionID(sessionID)
}

func (s *ShowService) CreateSeat(req *dto.SeatCreateRequest) (*models.Seat, error) {
	_, err := s.showRepo.GetSessionByID(req.SessionID)
	if err != nil {
		return nil, errors.New("场次不存在")
	}
	_, err = s.showRepo.GetSeatAreaByID(req.AreaID)
	if err != nil {
		return nil, errors.New("区域不存在")
	}

	seat := &models.Seat{
		SessionID: req.SessionID,
		AreaID:    req.AreaID,
		Row:       strings.ToUpper(req.Row),
		Col:       req.Col,
		SeatNo:    strings.ToUpper(req.Row) + "-" + uintToString(uint64(req.Col)),
		Status:    0,
		X:         req.X,
		Y:         req.Y,
		Width:     req.Width,
		Height:    req.Height,
	}

	err = s.showRepo.CreateSeat(seat)
	if err != nil {
		return nil, err
	}

	return seat, nil
}

func (s *ShowService) BatchCreateSeats(req *dto.SeatBatchCreateRequest) (int, error) {
	_, err := s.showRepo.GetSessionByID(req.SessionID)
	if err != nil {
		return 0, errors.New("场次不存在")
	}
	_, err = s.showRepo.GetSeatAreaByID(req.AreaID)
	if err != nil {
		return 0, errors.New("区域不存在")
	}

	startRow := strings.ToUpper(req.StartRow)
	endRow := strings.ToUpper(req.EndRow)

	startRowNum := int(startRow[0] - 'A')
	endRowNum := int(endRow[0] - 'A')

	if endRowNum < startRowNum {
		return 0, errors.New("结束排不能小于开始排")
	}
	if req.EndCol < req.StartCol {
		return 0, errors.New("结束列不能小于开始列")
	}

	var seats []models.Seat
	for r := startRowNum; r <= endRowNum; r++ {
		row := string(rune('A' + r))
		for c := req.StartCol; c <= req.EndCol; c++ {
			seat := models.Seat{
				SessionID: req.SessionID,
				AreaID:    req.AreaID,
				Row:       row,
				Col:       c,
				SeatNo:    row + "-" + uintToString(uint64(c)),
				Status:    0,
				X:         req.XOffset + (c-req.StartCol)*(req.SeatWidth+req.ColGap),
				Y:         req.YOffset + (r-startRowNum)*(req.SeatHeight+req.RowGap),
				Width:     req.SeatWidth,
				Height:    req.SeatHeight,
			}
			seats = append(seats, seat)
		}
	}

	err = s.showRepo.BatchCreateSeats(seats)
	if err != nil {
		return 0, err
	}

	return len(seats), nil
}

func (s *ShowService) GetSeats(sessionID uint64) ([]models.Seat, error) {
	seats, err := s.showRepo.GetSeatsBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	for i := range seats {
		status, _ := redis.GetSeatStatus(sessionID, seats[i].ID)
		if status == 1 {
			seats[i].Status = 1
		}

		lockUser, _ := redis.GetSeatLockUser(sessionID, seats[i].ID)
		if lockUser > 0 {
			seats[i].Status = 1
		}
	}

	return seats, nil
}

func (s *ShowService) LockSeats(sessionID uint64, seatIDs []uint64, userID uint64) ([]uint64, error) {
	seats, err := s.showRepo.GetSeatsByIDs(sessionID, seatIDs)
	if err != nil {
		return nil, err
	}

	if len(seats) != len(seatIDs) {
		return nil, errors.New("部分座位不存在")
	}

	session, err := s.showRepo.GetSessionByID(sessionID)
	if err != nil {
		return nil, errors.New("场次不存在")
	}

	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	memberLevel, err := userRepo.GetMemberLevelByLevel(user.MemberLevel)
	if err != nil {
		return nil, errors.New("获取会员等级失败")
	}

	now := time.Now()
	if !session.PresaleStartTime.IsZero() && now.Before(session.PresaleStartTime) {
		presaleDays := int(session.PresaleStartTime.Sub(now).Hours() / 24)
		if memberLevel.Priority < presaleDays {
			return nil, fmt.Errorf("预售阶段，%s可提前%d天购票，您的等级可提前%d天，请在开放后再试", 
				memberLevel.Name, memberLevel.Priority, presaleDays)
		}
	}

	lockMinutes := config.AppConfig.Seat.LockMinutes
	var failedSeats []uint64

	for _, seat := range seats {
		if seat.Status == models.SeatStatusSold {
			failedSeats = append(failedSeats, seat.ID)
			continue
		}

		locked, err := redis.LockSeat(sessionID, seat.ID, userID, lockMinutes)
		if err != nil || !locked {
			failedSeats = append(failedSeats, seat.ID)
			continue
		}

		_ = redis.SetSeatStatus(sessionID, seat.ID, models.SeatStatusLocked)
	}

	if len(failedSeats) > 0 {
		for _, seatID := range seatIDs {
			locked := false
			for _, failed := range failedSeats {
				if seatID == failed {
					locked = true
					break
				}
			}
			if !locked {
				_ = redis.UnlockSeat(sessionID, seatID)
				_ = redis.SetSeatStatus(sessionID, seatID, models.SeatStatusAvailable)
			}
		}
		return failedSeats, errors.New("部分座位锁定失败")
	}

	logging.Infof("Seats locked: session=%d, seats=%v, user=%d", sessionID, seatIDs, userID)
	return nil, nil
}

func (s *ShowService) UnlockSeats(sessionID uint64, seatIDs []uint64) error {
	for _, seatID := range seatIDs {
		_ = redis.UnlockSeat(sessionID, seatID)
		_ = redis.SetSeatStatus(sessionID, seatID, models.SeatStatusAvailable)
	}
	return nil
}

func (s *ShowService) ExtendSeatLock(sessionID uint64, seatID uint64, userID uint64) error {
	lockUser, err := redis.GetSeatLockUser(sessionID, seatID)
	if err != nil || lockUser != userID {
		return errors.New("无权限延长锁定")
	}

	lockMinutes := config.AppConfig.Seat.LockMinutes
	return redis.ExtendSeatLock(sessionID, seatID, lockMinutes)
}

func (s *ShowService) UpdateSeatChart(req *dto.SeatChartUpdateRequest) (*models.SeatChart, error) {
	chart, err := s.showRepo.GetSeatChartBySessionID(req.SessionID)
	if err != nil {
		return nil, err
	}

	if chart == nil {
		chart = &models.SeatChart{
			SessionID:  req.SessionID,
			Config:     req.Config,
			Background: req.Background,
			Width:      req.Width,
			Height:     req.Height,
		}
		err = s.showRepo.CreateSeatChart(chart)
	} else {
		chart.Config = req.Config
		chart.Background = req.Background
		if req.Width > 0 {
			chart.Width = req.Width
		}
		if req.Height > 0 {
			chart.Height = req.Height
		}
		err = s.showRepo.UpdateSeatChart(chart)
	}

	return chart, err
}

func (s *ShowService) GetSeatChart(sessionID uint64) (*models.SeatChart, error) {
	return s.showRepo.GetSeatChartBySessionID(sessionID)
}

func uintToString(n uint64) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	for n > 0 {
		result = append(result, byte(n%10)+'0')
		n /= 10
	}
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}
