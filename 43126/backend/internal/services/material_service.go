package services

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"meeting-room/internal/config"
	"meeting-room/internal/models"
	"meeting-room/internal/repositories"
	"meeting-room/internal/utils"

	"github.com/google/uuid"
)

type MaterialService struct {
	materialRepo *repositories.MaterialRepository
	bookingRepo  *repositories.BookingRepository
}

func NewMaterialService() *MaterialService {
	return &MaterialService{
		materialRepo: repositories.NewMaterialRepository(),
		bookingRepo:  repositories.NewBookingRepository(),
	}
}

type UploadMaterialRequest struct {
	BookingID uint
	UserID    uint
	FileName  string
	FileSize  int64
	FileType  string
	FileData  io.Reader
}

func (s *MaterialService) UploadMaterial(req *UploadMaterialRequest) (*models.MeetingMaterial, error) {
	booking, err := s.bookingRepo.FindByID(req.BookingID)
	if err != nil {
		return nil, errors.New("预订不存在")
	}

	if booking.UserID != req.UserID {
		return nil, errors.New("无权上传材料到此预订")
	}

	ext := filepath.Ext(req.FileName)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	uploadDir := config.Cfg.Upload.Dir + "/materials/"
	os.MkdirAll(uploadDir, 0755)

	filePath := uploadDir + newFilename
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, req.FileData)
	if err != nil {
		return nil, err
	}

	expireAt := booking.EndTime.AddDate(0, 0, 30)

	material := &models.MeetingMaterial{
		BookingID:   req.BookingID,
		UserID:      req.UserID,
		FileName:    req.FileName,
		FilePath:    "/uploads/materials/" + newFilename,
		FileSize:    req.FileSize,
		FileType:    req.FileType,
		MeetingDate: booking.StartTime,
		ExpireAt:    &expireAt,
	}

	err = s.materialRepo.Create(material)
	if err != nil {
		os.Remove(filePath)
		return nil, err
	}

	return material, nil
}

func (s *MaterialService) GetMaterialsByBooking(bookingID uint) ([]models.MeetingMaterial, error) {
	return s.materialRepo.GetByBooking(bookingID)
}

func (s *MaterialService) DownloadMaterial(id uint) (*models.MeetingMaterial, string, error) {
	material, err := s.materialRepo.FindByID(id)
	if err != nil {
		return nil, "", errors.New("文件不存在")
	}

	fullPath := config.Cfg.Upload.Dir + material.FilePath
	return material, fullPath, nil
}

func (s *MaterialService) DeleteMaterial(id uint, userID uint) error {
	material, err := s.materialRepo.FindByID(id)
	if err != nil {
		return errors.New("文件不存在")
	}

	if material.UserID != userID {
		return errors.New("无权删除此文件")
	}

	fullPath := config.Cfg.Upload.Dir + material.FilePath
	os.Remove(fullPath)

	return s.materialRepo.Delete(id)
}

func (s *MaterialService) CleanupExpiredMaterials() {
	materials, err := s.materialRepo.GetExpired()
	if err != nil {
		utils.Logger.Errorf("获取过期材料失败: %v", err)
		return
	}

	for _, material := range materials {
		fullPath := config.Cfg.Upload.Dir + material.FilePath
		os.Remove(fullPath)
		s.materialRepo.Delete(material.ID)
		utils.Logger.Infof("清理过期材料: %s", material.FileName)
	}
}

func (s *MaterialService) CleanupExpiredFiles() {
	s.CleanupExpiredMaterials()
}
