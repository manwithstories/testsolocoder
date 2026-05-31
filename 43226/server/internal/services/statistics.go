package services

import (
	"bytes"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"

	"museum-server/internal/dto"
	"museum-server/internal/models"
	"museum-server/internal/repository"
)

type StatisticsService struct {
	reservationRepo *repository.ReservationRepository
	exhibitionRepo  *repository.ExhibitionRepository
}

func NewStatisticsService(
	reservationRepo *repository.ReservationRepository,
	exhibitionRepo *repository.ExhibitionRepository,
) *StatisticsService {
	return &StatisticsService{
		reservationRepo: reservationRepo,
		exhibitionRepo:  exhibitionRepo,
	}
}

func (s *StatisticsService) GetStatistics(query *dto.StatisticsQuery) (map[string]interface{}, error) {
	stats, err := s.reservationRepo.GetStatistics(query)
	if err != nil {
		return nil, err
	}

	totalVisitors := int64(0)
	totalReservations := int64(0)
	totalRevenue := 0.0
	dailyData := make([]map[string]interface{}, 0)

	for _, stat := range stats {
		totalVisitors += int64(stat.VisitorCount)
		totalReservations += int64(stat.ReservationCount)
		totalRevenue += stat.Revenue

		dailyData = append(dailyData, map[string]interface{}{
			"date":             stat.StatDate.Format("2006-01-02"),
			"visitor_count":    stat.VisitorCount,
			"reservation_count": stat.ReservationCount,
			"revenue":          stat.Revenue,
		})
	}

	exhibitionListQuery := &dto.ExhibitionListQuery{
		Page:     1,
		PageSize: 100,
	}
	exhibitions, _, _ := s.exhibitionRepo.List(exhibitionListQuery)

	exhibitionStats := make([]map[string]interface{}, 0)
	for _, e := range exhibitions {
		visitorCount, reservationCount, revenue, _ := s.reservationRepo.AggregateDailyStats(
			e.ID, time.Now(),
		)
		_ = visitorCount
		exhibitionStats = append(exhibitionStats, map[string]interface{}{
			"id":                e.ID,
			"title":             e.Title,
			"view_count":        e.ViewCount,
			"reservation_count": reservationCount,
			"revenue":           revenue,
		})
	}

	return map[string]interface{}{
		"summary": map[string]interface{}{
			"total_visitors":     totalVisitors,
			"total_reservations": totalReservations,
			"total_revenue":      totalRevenue,
		},
		"daily_data":    dailyData,
		"exhibitions":   exhibitionStats,
	}, nil
}

func (s *StatisticsService) ExportExcel(query *dto.StatisticsQuery) ([]byte, error) {
	stats, err := s.reservationRepo.GetStatistics(query)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Statistics"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Date", "Exhibition ID", "Visitor Count", "Reservation Count", "Revenue"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, stat := range stats {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), stat.StatDate.Format("2006-01-02"))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), stat.ExhibitionID)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), stat.VisitorCount)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), stat.ReservationCount)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), stat.Revenue)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write excel: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *StatisticsService) ExportPDF(query *dto.StatisticsQuery) (map[string]interface{}, error) {
	stats, err := s.reservationRepo.GetStatistics(query)
	if err != nil {
		return nil, err
	}

	totalVisitors := 0
	totalReservations := 0
	totalRevenue := 0.0
	rows := make([]map[string]interface{}, 0)

	for _, stat := range stats {
		totalVisitors += stat.VisitorCount
		totalReservations += stat.ReservationCount
		totalRevenue += stat.Revenue

		rows = append(rows, map[string]interface{}{
			"date":              stat.StatDate.Format("2006-01-02"),
			"exhibition_id":     stat.ExhibitionID,
			"visitor_count":     stat.VisitorCount,
			"reservation_count": stat.ReservationCount,
			"revenue":           stat.Revenue,
		})
	}

	return map[string]interface{}{
		"title": "Museum Statistics Report",
		"generated_at": time.Now().Format("2006-01-02 15:04:05"),
		"summary": map[string]interface{}{
			"total_visitors":     totalVisitors,
			"total_reservations": totalReservations,
			"total_revenue":      totalRevenue,
		},
		"rows": rows,
	}, nil
}

type MuseumService struct {
	museumRepo *repository.MuseumRepository
}

func NewMuseumService(museumRepo *repository.MuseumRepository) *MuseumService {
	return &MuseumService{museumRepo: museumRepo}
}

func (s *MuseumService) Create(req *dto.MuseumRequest) (*models.Museum, error) {
	museum := &models.Museum{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Contact:     req.Contact,
		Phone:       req.Phone,
		Email:       req.Email,
		Logo:        req.Logo,
		OpenTime:    req.OpenTime,
		CloseTime:   req.CloseTime,
	}

	if err := s.museumRepo.Create(museum); err != nil {
		return nil, fmt.Errorf("failed to create museum: %w", err)
	}

	return museum, nil
}

func (s *MuseumService) Update(id uint, req *dto.MuseumRequest) error {
	museum, err := s.museumRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("museum not found")
	}

	museum.Name = req.Name
	museum.Description = req.Description
	museum.Address = req.Address
	museum.Contact = req.Contact
	museum.Phone = req.Phone
	museum.Email = req.Email
	if req.Logo != "" {
		museum.Logo = req.Logo
	}
	museum.OpenTime = req.OpenTime
	museum.CloseTime = req.CloseTime

	return s.museumRepo.Update(museum)
}

func (s *MuseumService) Delete(id uint) error {
	return s.museumRepo.Delete(id)
}

func (s *MuseumService) GetByID(id uint) (*models.Museum, error) {
	return s.museumRepo.FindByID(id)
}

func (s *MuseumService) List() ([]models.Museum, error) {
	return s.museumRepo.List()
}
