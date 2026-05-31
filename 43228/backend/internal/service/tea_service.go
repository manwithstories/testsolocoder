package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tea-platform/config"
	"tea-platform/internal/models"
	"tea-platform/internal/repository"
)

type CreateTeaRequest struct {
	Name             string  `json:"name" binding:"required,max=128"`
	Type             string  `json:"type" binding:"required,oneof=green_tea black_tea oolong puer white_tea yellow_tea dark_tea flower_tea"`
	Origin           string  `json:"origin" binding:"omitempty,max=128"`
	Year             int     `json:"year" binding:"omitempty,min=1900,max=2100"`
	Grade            string  `json:"grade" binding:"omitempty,oneof=special grade1 grade2 grade3"`
	ProcessType      string  `json:"process_type" binding:"omitempty,oneof=manual semi_manual mechanism"`
	StorageCondition string  `json:"storage_condition" binding:"omitempty,max=255"`
	Description      string  `json:"description"`
	Price            float64 `json:"price" binding:"required,gte=0"`
	Stock            int     `json:"stock" binding:"omitempty,gte=0"`
}

type UpdateTeaRequest struct {
	Name             *string  `json:"name" binding:"omitempty,max=128"`
	Type             *string  `json:"type" binding:"omitempty,oneof=green_tea black_tea oolong puer white_tea yellow_tea dark_tea flower_tea"`
	Origin           *string  `json:"origin" binding:"omitempty,max=128"`
	Year             *int     `json:"year" binding:"omitempty,min=1900,max=2100"`
	Grade            *string  `json:"grade" binding:"omitempty,oneof=special grade1 grade2 grade3"`
	ProcessType      *string  `json:"process_type" binding:"omitempty,oneof=manual semi_manual mechanism"`
	StorageCondition *string  `json:"storage_condition" binding:"omitempty,max=255"`
	Description      *string  `json:"description"`
	Price            *float64 `json:"price" binding:"omitempty,gte=0"`
	Stock            *int     `json:"stock" binding:"omitempty,gte=0"`
}

type AddTraceabilityRequest struct {
	Stage       string    `json:"stage" binding:"required,max=64"`
	Description string    `json:"description"`
	Location    string    `json:"location" binding:"omitempty,max=255"`
	Operator    string    `json:"operator" binding:"omitempty,max=128"`
	RecordTime  time.Time `json:"record_time" binding:"required"`
}

type TeaDetailResponse struct {
	ID               uint                     `json:"id"`
	Name             string                   `json:"name"`
	Type             string                   `json:"type"`
	Origin           string                   `json:"origin"`
	Year             int                      `json:"year"`
	Grade            string                   `json:"grade"`
	ProcessType      string                   `json:"process_type"`
	StorageCondition string                   `json:"storage_condition"`
	Description      string                   `json:"description"`
	Price            float64                  `json:"price"`
	Stock            int                      `json:"stock"`
	SellerID         uint                     `json:"seller_id"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
	Images           []models.TeaImage        `json:"images"`
	Traceability     []models.Traceability    `json:"traceability"`
}

type TeaService struct {
	repo *repository.TeaRepository
}

func NewTeaService() *TeaService {
	return &TeaService{
		repo: repository.NewTeaRepository(),
	}
}

func (s *TeaService) CreateTea(userID uint, req *CreateTeaRequest) (*models.Tea, error) {
	tea := &models.Tea{
		Name:             req.Name,
		Type:             models.TeaType(req.Type),
		Origin:           req.Origin,
		Year:             req.Year,
		Grade:            models.TeaGrade(req.Grade),
		ProcessType:      models.ProcessType(req.ProcessType),
		StorageCondition: req.StorageCondition,
		Description:      req.Description,
		Price:            req.Price,
		Stock:            req.Stock,
		SellerID:         userID,
	}

	if err := s.repo.CreateTea(tea); err != nil {
		return nil, errors.New("创建茶叶失败")
	}

	return tea, nil
}

func (s *TeaService) UpdateTea(userID uint, teaID uint, req *UpdateTeaRequest) error {
	tea, err := s.repo.GetTeaByID(teaID)
	if err != nil {
		return errors.New("茶叶不存在")
	}

	if tea.SellerID != userID {
		return errors.New("无权修改该茶叶")
	}

	if req.Name != nil {
		tea.Name = *req.Name
	}
	if req.Type != nil {
		tea.Type = models.TeaType(*req.Type)
	}
	if req.Origin != nil {
		tea.Origin = *req.Origin
	}
	if req.Year != nil {
		tea.Year = *req.Year
	}
	if req.Grade != nil {
		tea.Grade = models.TeaGrade(*req.Grade)
	}
	if req.ProcessType != nil {
		tea.ProcessType = models.ProcessType(*req.ProcessType)
	}
	if req.StorageCondition != nil {
		tea.StorageCondition = *req.StorageCondition
	}
	if req.Description != nil {
		tea.Description = *req.Description
	}
	if req.Price != nil {
		tea.Price = *req.Price
	}
	if req.Stock != nil {
		tea.Stock = *req.Stock
	}

	if err := s.repo.UpdateTea(tea); err != nil {
		return errors.New("更新茶叶失败")
	}

	return nil
}

func (s *TeaService) DeleteTea(userID uint, teaID uint) error {
	tea, err := s.repo.GetTeaByID(teaID)
	if err != nil {
		return errors.New("茶叶不存在")
	}

	if tea.SellerID != userID {
		return errors.New("无权删除该茶叶")
	}

	if err := s.repo.DeleteTea(teaID); err != nil {
		return errors.New("删除茶叶失败")
	}

	return nil
}

func (s *TeaService) GetTeaByID(id uint) (*TeaDetailResponse, error) {
	tea, err := s.repo.GetTeaByID(id)
	if err != nil {
		return nil, errors.New("茶叶不存在")
	}

	images, _ := s.repo.GetTeaImages(id)
	traces, _ := s.repo.GetTraceabilityByTeaID(id)

	resp := &TeaDetailResponse{
		ID:               tea.ID,
		Name:             tea.Name,
		Type:             string(tea.Type),
		Origin:           tea.Origin,
		Year:             tea.Year,
		Grade:            string(tea.Grade),
		ProcessType:      string(tea.ProcessType),
		StorageCondition: tea.StorageCondition,
		Description:      tea.Description,
		Price:            tea.Price,
		Stock:            tea.Stock,
		SellerID:         tea.SellerID,
		CreatedAt:        tea.CreatedAt,
		UpdatedAt:        tea.UpdatedAt,
		Images:           images,
		Traceability:     traces,
	}

	return resp, nil
}

func (s *TeaService) GetTeaList(page, pageSize int, filters repository.TeaFilters) ([]models.Tea, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.GetTeaList(page, pageSize, filters)
}

func (s *TeaService) UploadTeaImage(teaID uint, imageType string, file *multipart.FileHeader) error {
	tea, err := s.repo.GetTeaByID(teaID)
	if err != nil {
		return errors.New("茶叶不存在")
	}

	imageTypeStr := string(models.TeaImageTypeDetail)
	switch imageType {
	case "main":
		imageTypeStr = string(models.TeaImageTypeMain)
	case "detail":
		imageTypeStr = string(models.TeaImageTypeDetail)
	case "packaging":
		imageTypeStr = string(models.TeaImageTypePackaging)
	}

	cfg := config.Get()
	uploadPath := cfg.Upload.TeaImagesPath

	maxSize := int64(cfg.Upload.MaxSize) * 1024 * 1024
	if file.Size > maxSize {
		return fmt.Errorf("文件大小超过限制 (最大 %dMB)", cfg.Upload.MaxSize)
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))
	allowed := false
	for _, t := range cfg.Upload.AllowedTypes {
		if ext == t {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("不支持的文件类型: %s", ext)
	}

	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return fmt.Errorf("创建上传目录失败: %w", err)
	}

	now := time.Now()
	fileName := fmt.Sprintf("%d_%d.%s", now.UnixNano(), now.Nanosecond(), ext)
	fullPath := filepath.Join(uploadPath, fileName)

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("保存文件失败: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("保存文件失败: %w", err)
	}

	image := &models.TeaImage{
		TeaID:     tea.ID,
		ImageURL:  fullPath,
		ImageType: models.TeaImageType(imageTypeStr),
	}

	if err := s.repo.CreateTeaImage(image); err != nil {
		return errors.New("保存图片失败")
	}

	return nil
}

func (s *TeaService) DeleteTeaImage(imageID uint) error {
	if err := s.repo.DeleteTeaImage(imageID); err != nil {
		return errors.New("删除图片失败")
	}
	return nil
}

func (s *TeaService) AddTraceability(teaID uint, req *AddTraceabilityRequest) error {
	_, err := s.repo.GetTeaByID(teaID)
	if err != nil {
		return errors.New("茶叶不存在")
	}

	trace := &models.Traceability{
		TeaID:       teaID,
		Stage:       req.Stage,
		Description: req.Description,
		Location:    req.Location,
		Operator:    req.Operator,
		RecordTime:  req.RecordTime,
	}

	if err := s.repo.CreateTraceability(trace); err != nil {
		return errors.New("添加溯源记录失败")
	}

	return nil
}

func (s *TeaService) GetTraceability(teaID uint) ([]models.Traceability, error) {
	traces, err := s.repo.GetTraceabilityByTeaID(teaID)
	if err != nil {
		return nil, errors.New("查询溯源记录失败")
	}
	return traces, nil
}

func (s *TeaService) SearchTea(keyword string, page, pageSize int) ([]models.Tea, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	filters := repository.TeaFilters{
		Keyword: keyword,
	}
	return s.repo.GetTeaList(page, pageSize, filters)
}
