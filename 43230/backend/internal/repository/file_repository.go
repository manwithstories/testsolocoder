package repository

import (
	"context"
	"errors"
	"time"

	"print3d-platform/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) CreateUpload(ctx context.Context, upload *models.FileUpload) error {
	return r.db.WithContext(ctx).Create(upload).Error
}

func (r *FileRepository) GetUpload(ctx context.Context, id uuid.UUID) (*models.FileUpload, error) {
	var upload models.FileUpload
	err := r.db.WithContext(ctx).First(&upload, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &upload, nil
}

func (r *FileRepository) UpdateUpload(ctx context.Context, upload *models.FileUpload) error {
	return r.db.WithContext(ctx).Save(upload).Error
}

func (r *FileRepository) UpdateUploadProgress(ctx context.Context, id uuid.UUID, uploadedChunks int, status string) error {
	updates := map[string]interface{}{
		"uploaded_chunks": uploadedChunks,
		"upload_status":   status,
	}
	if status == "completed" {
		updates["expires_at"] = nil
	}
	return r.db.WithContext(ctx).Model(&models.FileUpload{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *FileRepository) GetUserUploads(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.FileUpload, int64, error) {
	var uploads []models.FileUpload
	var total int64

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&uploads).Error
	return uploads, total, err
}

func (r *FileRepository) DeleteExpiredUploads(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("upload_status = ? AND expires_at < ?", "uploading", time.Now()).
		Delete(&models.FileUpload{}).Error
}

func (r *FileRepository) CreateAccessLog(ctx context.Context, log *models.FileAccessLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *FileRepository) GetAccessLogs(ctx context.Context, fileID uuid.UUID, page, pageSize int) ([]models.FileAccessLog, int64, error) {
	var logs []models.FileAccessLog
	var total int64

	query := r.db.WithContext(ctx).Where("file_id = ?", fileID)
	query.Count(&total)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, total, err
}

func (r *FileRepository) CheckFileHashExists(ctx context.Context, hash string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.FileUpload{}).
		Where("file_hash = ?", hash).
		Count(&count).Error
	return count > 0, err
}

func (r *FileRepository) DeleteUpload(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.FileUpload{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("upload not found")
	}
	return nil
}

func (r *FileRepository) GetFileByHash(ctx context.Context, hash string) (*models.FileUpload, error) {
	var upload models.FileUpload
	err := r.db.WithContext(ctx).
		Where("file_hash = ? AND upload_status = ?", hash, "completed").
		First(&upload).Error
	if err != nil {
		return nil, err
	}
	return &upload, nil
}
