package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"time"

	"print3d-platform/internal/config"
	"print3d-platform/internal/models"
	"print3d-platform/internal/repository"
	"print3d-platform/internal/utils"

	"github.com/google/uuid"
)

type FileService struct {
	fileRepo *repository.FileRepository
	cfg      *config.Config
}

func NewFileService(fileRepo *repository.FileRepository, cfg *config.Config) *FileService {
	return &FileService{
		fileRepo: fileRepo,
		cfg:      cfg,
	}
}

type InitiateUploadRequest struct {
	FileName string `json:"file_name" binding:"required"`
	FileSize int64  `json:"file_size" binding:"required,min=1"`
	FileType string `json:"file_type"`
}

type UploadChunkRequest struct {
	UploadID    uuid.UUID `json:"upload_id" binding:"required"`
	ChunkNumber int       `json:"chunk_number" binding:"required"`
	TotalChunks int       `json:"total_chunks" binding:"required"`
	ChunkSize   int64     `json:"chunk_size" binding:"required"`
}

type InitiateUploadResponse struct {
	UploadID    uuid.UUID `json:"upload_id"`
	ChunkSize   int64     `json:"chunk_size"`
	TotalChunks int       `json:"total_chunks"`
	ExpiresAt   time.Time `json:"expires_at"`
}

const DefaultChunkSize = 5 * 1024 * 1024

func (s *FileService) InitiateUpload(ctx context.Context, userID uuid.UUID, req *InitiateUploadRequest) (*InitiateUploadResponse, error) {
	ext := utils.GetFileExtension(req.FileName)
	if !utils.IsAllowedExtension(ext, s.cfg.Storage.AllowedExtensions) {
		return nil, fmt.Errorf("file type %s not allowed", ext)
	}

	if req.FileSize > s.cfg.Storage.MaxFileSize {
		return nil, errors.New("file size exceeds maximum allowed")
	}

	totalChunks := int(math.Ceil(float64(req.FileSize) / float64(DefaultChunkSize)))
	expiresAt := time.Now().Add(24 * time.Hour)

	upload := &models.FileUpload{
		ID:             uuid.New(),
		UserID:         userID,
		FileName:       req.FileName,
		OriginalName:   req.FileName,
		FileType:       req.FileType,
		FileSize:       req.FileSize,
		UploadStatus:   "uploading",
		TotalChunks:    totalChunks,
		UploadedChunks: 0,
		ChunkSize:      DefaultChunkSize,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      &expiresAt,
	}

	err := s.fileRepo.CreateUpload(ctx, upload)
	if err != nil {
		return nil, err
	}

	uploadDir := filepath.Join(s.cfg.Storage.UploadPath, "chunks", upload.ID.String())
	err = utils.EnsureDir(uploadDir)
	if err != nil {
		return nil, err
	}

	return &InitiateUploadResponse{
		UploadID:    upload.ID,
		ChunkSize:   DefaultChunkSize,
		TotalChunks: totalChunks,
		ExpiresAt:   expiresAt,
	}, nil
}

func (s *FileService) UploadChunk(ctx context.Context, userID uuid.UUID, uploadID uuid.UUID, chunkNumber int, chunkData io.Reader) error {
	upload, err := s.fileRepo.GetUpload(ctx, uploadID)
	if err != nil {
		return errors.New("upload not found")
	}

	if upload.UserID != userID {
		return errors.New("not authorized")
	}

	if upload.UploadStatus != "uploading" {
		return errors.New("upload already completed or expired")
	}

	if upload.ExpiresAt != nil && upload.ExpiresAt.Before(time.Now()) {
		return errors.New("upload has expired")
	}

	chunkDir := filepath.Join(s.cfg.Storage.UploadPath, "chunks", upload.ID.String())
	chunkPath := filepath.Join(chunkDir, fmt.Sprintf("chunk_%d", chunkNumber))

	tmpFile, err := os.Create(chunkPath)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, chunkData)
	if err != nil {
		return err
	}

	newChunkCount := upload.UploadedChunks + 1
	status := "uploading"
	if newChunkCount >= upload.TotalChunks {
		status = "chunks_completed"
	}

	return s.fileRepo.UpdateUploadProgress(ctx, uploadID, newChunkCount, status)
}

func (s *FileService) CompleteUpload(ctx context.Context, userID uuid.UUID, uploadID uuid.UUID) (*models.FileUpload, error) {
	upload, err := s.fileRepo.GetUpload(ctx, uploadID)
	if err != nil {
		return nil, errors.New("upload not found")
	}

	if upload.UserID != userID {
		return nil, errors.New("not authorized")
	}

	if upload.UploadedChunks < upload.TotalChunks {
		return nil, errors.New("not all chunks uploaded")
	}

	chunkDir := filepath.Join(s.cfg.Storage.UploadPath, "chunks", upload.ID.String())
	mergedDir := filepath.Join(s.cfg.Storage.UploadPath, "completed")
	err = utils.EnsureDir(mergedDir)
	if err != nil {
		return nil, err
	}

	finalFilename := utils.GenerateFileName(upload.FileName)
	finalPath := filepath.Join(mergedDir, finalFilename)

	mergedFile, err := os.Create(finalPath)
	if err != nil {
		return nil, err
	}
	defer mergedFile.Close()

	for i := 0; i < upload.TotalChunks; i++ {
		chunkPath := filepath.Join(chunkDir, fmt.Sprintf("chunk_%d", i))
		chunk, err := os.Open(chunkPath)
		if err != nil {
			return nil, fmt.Errorf("missing chunk %d", i)
		}

		_, err = io.Copy(mergedFile, chunk)
		chunk.Close()
		if err != nil {
			return nil, err
		}
	}

	fileHash, err := utils.HashFilePath(finalPath)
	if err != nil {
		return nil, err
	}

	exists, err := s.fileRepo.CheckFileHashExists(ctx, fileHash)
	if err != nil {
		return nil, err
	}
	if exists {
		os.Remove(finalPath)
		os.RemoveAll(chunkDir)
		return nil, errors.New("file already exists")
	}

	storagePath := filepath.Join("completed", finalFilename)

	upload.FileHash = fileHash
	upload.StoragePath = storagePath
	upload.FileName = finalFilename
	upload.UploadStatus = "completed"
	upload.ExpiresAt = nil
	upload.UpdatedAt = time.Now()

	err = s.fileRepo.UpdateUpload(ctx, upload)
	if err != nil {
		return nil, err
	}

	os.RemoveAll(chunkDir)

	accessLog := &models.FileAccessLog{
		ID:         uuid.New(),
		FileID:     upload.ID,
		UserID:     userID,
		AccessType: "upload",
		CreatedAt:  time.Now(),
	}
	_ = s.fileRepo.CreateAccessLog(ctx, accessLog)

	utils.LogInfo("File upload completed: %s by user %s", upload.OriginalName, userID)
	return upload, nil
}

func (s *FileService) GetUpload(ctx context.Context, id uuid.UUID) (*models.FileUpload, error) {
	return s.fileRepo.GetUpload(ctx, id)
}

func (s *FileService) GetUserUploads(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.FileUpload, int64, error) {
	return s.fileRepo.GetUserUploads(ctx, userID, page, pageSize)
}

func (s *FileService) DeleteUpload(ctx context.Context, id, userID uuid.UUID) error {
	upload, err := s.fileRepo.GetUpload(ctx, id)
	if err != nil {
		return err
	}

	if upload.UserID != userID {
		return errors.New("not authorized")
	}

	if upload.StoragePath != "" {
		filePath := filepath.Join(s.cfg.Storage.UploadPath, upload.StoragePath)
		os.Remove(filePath)
	}

	chunkDir := filepath.Join(s.cfg.Storage.UploadPath, "chunks", id.String())
	os.RemoveAll(chunkDir)

	return s.fileRepo.DeleteUpload(ctx, id)
}

func (s *FileService) LogFileAccess(ctx context.Context, fileID, userID uuid.UUID, accessType, ipAddress, userAgent string) error {
	accessLog := &models.FileAccessLog{
		ID:         uuid.New(),
		FileID:     fileID,
		UserID:     userID,
		AccessType: accessType,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		CreatedAt:  time.Now(),
	}
	return s.fileRepo.CreateAccessLog(ctx, accessLog)
}

func (s *FileService) GetAccessLogs(ctx context.Context, fileID uuid.UUID, page, pageSize int) ([]models.FileAccessLog, int64, error) {
	return s.fileRepo.GetAccessLogs(ctx, fileID, page, pageSize)
}

func (s *FileService) CleanupExpiredUploads(ctx context.Context) error {
	return s.fileRepo.DeleteExpiredUploads(ctx)
}
