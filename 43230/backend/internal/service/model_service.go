package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"print3d-platform/internal/config"
	"print3d-platform/internal/models"
	"print3d-platform/internal/repository"
	"print3d-platform/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelService struct {
	modelRepo  *repository.ModelRepository
	userRepo   *repository.UserRepository
	orderRepo  *repository.OrderRepository
	cfg        *config.Config
}

func NewModelService(modelRepo *repository.ModelRepository, userRepo *repository.UserRepository, orderRepo *repository.OrderRepository, cfg *config.Config) *ModelService {
	return &ModelService{
		modelRepo: modelRepo,
		userRepo:  userRepo,
		orderRepo: orderRepo,
		cfg:       cfg,
	}
}

type CreateModelRequest struct {
	Title             string              `json:"title" binding:"required"`
	Description       string              `json:"description"`
	Category          string              `json:"category"`
	Tags              []string            `json:"tags"`
	Price             float64             `json:"price" binding:"required,min=0"`
	LicenseType       models.LicenseType  `json:"license_type" binding:"required"`
	SubscriptionPrice float64             `json:"subscription_price"`
	Volume            float64             `json:"volume"`
	BoundingBox       string              `json:"bounding_box"`
	PrintTimeHours    float64             `json:"print_time_hours"`
	RecommendedMaterials []string         `json:"recommended_materials"`
	PolygonCount      int                 `json:"polygon_count"`
}

type UpdateModelRequest struct {
	Title             string              `json:"title"`
	Description       string              `json:"description"`
	Category          string              `json:"category"`
	Tags              []string            `json:"tags"`
	Price             float64             `json:"price"`
	LicenseType       models.LicenseType  `json:"license_type"`
	SubscriptionPrice float64             `json:"subscription_price"`
}

type PurchaseModelRequest struct {
	PurchaseType models.LicenseType `json:"purchase_type" binding:"required"`
}

func (s *ModelService) CreateModel(ctx context.Context, designerID uuid.UUID, req *CreateModelRequest) (*models.Model3D, error) {
	now := time.Now()
	model := &models.Model3D{
		ID:                   uuid.New(),
		DesignerID:           designerID,
		Title:                req.Title,
		Description:          req.Description,
		Category:             req.Category,
		Tags:                 req.Tags,
		Price:                req.Price,
		LicenseType:          req.LicenseType,
		SubscriptionPrice:    req.SubscriptionPrice,
		Status:               models.ModelStatusDraft,
		Version:              "1.0.0",
		Volume:               req.Volume,
		BoundingBox:          req.BoundingBox,
		PrintTimeHours:       req.PrintTimeHours,
		RecommendedMaterials: req.RecommendedMaterials,
		PolygonCount:         req.PolygonCount,
		Rating:               5.0,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	err := s.modelRepo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	utils.LogInfo("New model created: %s by designer %s", model.Title, designerID)
	return model, nil
}

func (s *ModelService) GetModel(ctx context.Context, id uuid.UUID) (*models.Model3D, error) {
	err := s.modelRepo.IncrementViewCount(ctx, id)
	if err != nil {
		utils.LogWarn("Failed to increment view count for model %s: %v", id, err)
	}
	return s.modelRepo.FindByID(ctx, id)
}

func (s *ModelService) UpdateModel(ctx context.Context, id, designerID uuid.UUID, req *UpdateModelRequest) (*models.Model3D, error) {
	model, err := s.modelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if model.DesignerID != designerID {
		return nil, errors.New("not authorized to update this model")
	}

	if req.Title != "" {
		model.Title = req.Title
	}
	if req.Description != "" {
		model.Description = req.Description
	}
	if req.Category != "" {
		model.Category = req.Category
	}
	if req.Tags != nil {
		model.Tags = req.Tags
	}
	if req.Price >= 0 {
		model.Price = req.Price
	}
	if req.LicenseType != "" {
		model.LicenseType = req.LicenseType
	}
	if req.SubscriptionPrice >= 0 {
		model.SubscriptionPrice = req.SubscriptionPrice
	}

	model.UpdatedAt = time.Now()
	err = s.modelRepo.Update(ctx, model)
	return model, err
}

func (s *ModelService) DeleteModel(ctx context.Context, id, designerID uuid.UUID) error {
	model, err := s.modelRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if model.DesignerID != designerID {
		return errors.New("not authorized to delete this model")
	}

	return s.modelRepo.Delete(ctx, id)
}

func (s *ModelService) ListModels(ctx context.Context, filter repository.ModelFilter, page, pageSize int) ([]models.Model3D, int64, error) {
	return s.modelRepo.List(ctx, filter, page, pageSize)
}

func (s *ModelService) GetDesignerModels(ctx context.Context, designerID uuid.UUID, page, pageSize int) ([]models.Model3D, int64, error) {
	return s.modelRepo.GetDesignerModels(ctx, designerID, page, pageSize)
}

func (s *ModelService) UploadModelFile(ctx context.Context, modelID, designerID uuid.UUID, file io.Reader, filename string, fileSize int64) error {
	model, err := s.modelRepo.FindByID(ctx, modelID)
	if err != nil {
		return err
	}

	if model.DesignerID != designerID {
		return errors.New("not authorized to upload file for this model")
	}

	if fileSize > s.cfg.Storage.MaxFileSize {
		return errors.New("file size exceeds maximum allowed")
	}

	ext := utils.GetFileExtension(filename)
	if !utils.IsAllowedExtension(ext, s.cfg.Storage.AllowedExtensions) {
		return fmt.Errorf("file type %s not allowed", ext)
	}

	uploadDir := filepath.Join(s.cfg.Storage.UploadPath, "models", modelID.String())
	err = utils.EnsureDir(uploadDir)
	if err != nil {
		return err
	}

	newFilename := utils.GenerateFileName(filename)
	filePath := filepath.Join(uploadDir, newFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	fileHash, err := utils.HashFilePath(filePath)
	if err != nil {
		return err
	}

	model.FileURL = fmt.Sprintf("/uploads/models/%s/%s", modelID.String(), newFilename)
	model.FileSize = fileSize
	model.FileType = strings.TrimPrefix(ext, ".")
	model.FileHash = fileHash
	model.UpdatedAt = time.Now()

	return s.modelRepo.Update(ctx, model)
}

func (s *ModelService) UploadThumbnail(ctx context.Context, modelID, designerID uuid.UUID, file io.Reader, filename string) error {
	model, err := s.modelRepo.FindByID(ctx, modelID)
	if err != nil {
		return err
	}

	if model.DesignerID != designerID {
		return errors.New("not authorized to upload thumbnail for this model")
	}

	ext := utils.GetFileExtension(filename)
	allowedImageExts := []string{".jpg", ".jpeg", ".png"}
	if !utils.IsAllowedExtension(ext, allowedImageExts) {
		return errors.New("only JPG and PNG images are allowed")
	}

	thumbnailDir := filepath.Join(s.cfg.Storage.UploadPath, "thumbnails", modelID.String())
	err = utils.EnsureDir(thumbnailDir)
	if err != nil {
		return err
	}

	newFilename := "thumbnail" + ext
	filePath := filepath.Join(thumbnailDir, newFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}

	model.ThumbnailURL = fmt.Sprintf("/uploads/thumbnails/%s/%s", modelID.String(), newFilename)
	model.UpdatedAt = time.Now()

	return s.modelRepo.Update(ctx, model)
}

func (s *ModelService) PurchaseModel(ctx context.Context, modelID, userID uuid.UUID, req *PurchaseModelRequest) (*models.ModelPurchase, error) {
	model, err := s.modelRepo.FindByID(ctx, modelID)
	if err != nil {
		return nil, err
	}

	if model.Status != models.ModelStatusPublished {
		return nil, errors.New("model is not available for purchase")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	var amount float64
	var expiresAt *time.Time

	switch req.PurchaseType {
	case models.LicensePerPurchase:
		amount = model.Price
	case models.LicenseSubscription:
		if model.SubscriptionPrice <= 0 {
			return nil, errors.New("subscription not available for this model")
		}
		amount = model.SubscriptionPrice
		exp := time.Now().AddDate(0, 1, 0)
		expiresAt = &exp
	default:
		return nil, errors.New("invalid purchase type")
	}

	if user.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	err = s.userRepo.ExecTx(ctx, func(tx *gorm.DB) error {
		err = tx.WithContext(ctx).Model(&models.User{}).
			Where("id = ?", userID).
			Update("balance", gorm.Expr("balance - ?", amount)).Error
		if err != nil {
			return err
		}

		err = tx.WithContext(ctx).Model(&models.User{}).
			Where("id = ?", model.DesignerID).
			Update("balance", gorm.Expr("balance + ?", amount*s.cfg.Pricing.DesignerFeeRate)).Error
		if err != nil {
			return err
		}

		err = tx.WithContext(ctx).Model(&models.DesignerProfile{}).
			Where("user_id = ?", model.DesignerID).
			Updates(map[string]interface{}{
				"total_models": gorm.Expr("total_models + 1"),
				"total_sales":  gorm.Expr("total_sales + ?", amount),
			}).Error
		if err != nil {
			return err
		}

		purchase := &models.ModelPurchase{
			ID:            uuid.New(),
			ModelID:       modelID,
			UserID:        userID,
			PurchaseType:  req.PurchaseType,
			Amount:        amount,
			TransactionID: utils.GenerateTransactionNo(),
			ExpiresAt:     expiresAt,
			CreatedAt:     time.Now(),
		}

		err = tx.WithContext(ctx).Create(purchase).Error
		if err != nil {
			return err
		}

		err = tx.WithContext(ctx).Model(&models.Model3D{}).
			Where("id = ?", modelID).
			Update("purchase_count", gorm.Expr("purchase_count + 1")).Error

		return err
	})

	if err != nil {
		return nil, err
	}

	purchase, err := s.modelRepo.GetPurchase(ctx, modelID, userID)
	utils.LogInfo("Model %s purchased by user %s, amount: %.2f", modelID, userID, amount)
	return purchase, nil
}

func (s *ModelService) DownloadModel(ctx context.Context, modelID, userID uuid.UUID) (string, error) {
	model, err := s.modelRepo.FindByID(ctx, modelID)
	if err != nil {
		return "", err
	}

	purchased, err := s.modelRepo.IsPurchased(ctx, modelID, userID)
	if err != nil {
		return "", err
	}

	if !purchased && model.DesignerID != userID {
		return "", errors.New("you need to purchase this model first")
	}

	err = s.modelRepo.IncrementDownloadCount(ctx, modelID)
	if err != nil {
		utils.LogWarn("Failed to increment download count: %v", err)
	}

	downloadRecord := &models.DownloadRecord{
		ID:        uuid.New(),
		ModelID:   modelID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	err = s.modelRepo.CreateDownloadRecord(ctx, downloadRecord)
	if err != nil {
		utils.LogWarn("Failed to create download record: %v", err)
	}

	return model.FileURL, nil
}

func (s *ModelService) AddFavorite(ctx context.Context, modelID, userID uuid.UUID) error {
	isFav, err := s.modelRepo.IsFavorite(ctx, modelID, userID)
	if err != nil {
		return err
	}
	if isFav {
		return errors.New("already favorited")
	}

	favorite := &models.ModelFavorite{
		ID:        uuid.New(),
		ModelID:   modelID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	err = s.modelRepo.AddFavorite(ctx, favorite)
	if err != nil {
		return err
	}

	return s.modelRepo.Update(ctx, &models.Model3D{
		ID:            modelID,
		FavoriteCount: 1,
	})
}

func (s *ModelService) RemoveFavorite(ctx context.Context, modelID, userID uuid.UUID) error {
	return s.modelRepo.RemoveFavorite(ctx, modelID, userID)
}

func (s *ModelService) GetUserFavorites(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Model3D, int64, error) {
	return s.modelRepo.GetUserFavorites(ctx, userID, page, pageSize)
}

func (s *ModelService) GetHotModels(ctx context.Context, limit int) ([]models.Model3D, error) {
	return s.modelRepo.GetHotModels(ctx, limit)
}

func (s *ModelService) GetUserPurchases(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.ModelPurchase, int64, error) {
	return s.modelRepo.GetUserPurchases(ctx, userID, page, pageSize)
}

func (s *ModelService) CreateNewVersion(ctx context.Context, modelID, designerID uuid.UUID, changeLog string, file io.Reader, filename string) error {
	model, err := s.modelRepo.FindByID(ctx, modelID)
	if err != nil {
		return err
	}

	if model.DesignerID != designerID {
		return errors.New("not authorized")
	}

	oldVersion := &models.ModelVersion{
		ID:            uuid.New(),
		ModelID:       modelID,
		VersionNumber: model.Version,
		FileURL:       model.FileURL,
		FileHash:      model.FileHash,
		ChangeLog:     changeLog,
		CreatedAt:     time.Now(),
	}

	err = s.modelRepo.CreateVersion(ctx, oldVersion)
	if err != nil {
		return err
	}

	return s.UploadModelFile(ctx, modelID, designerID, file, filename, model.FileSize)
}

func (s *ModelService) GetVersions(ctx context.Context, modelID uuid.UUID) ([]models.ModelVersion, error) {
	return s.modelRepo.GetVersions(ctx, modelID)
}

func (s *ModelService) GetStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, error) {
	return s.modelRepo.GetStats(ctx, startDate, endDate)
}

func (s *ModelService) roundFloat(val float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(val*pow) / pow
}
