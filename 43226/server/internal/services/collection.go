package services

import (
	"fmt"

	"museum-server/internal/dto"
	"museum-server/internal/models"
	"museum-server/internal/repository"
)

type CollectionService struct {
	collectionRepo *repository.CollectionRepository
}

func NewCollectionService(collectionRepo *repository.CollectionRepository) *CollectionService {
	return &CollectionService{collectionRepo: collectionRepo}
}

func (s *CollectionService) Create(museumID uint, req *dto.CollectionRequest) (*models.Collection, error) {
	existing, _ := s.collectionRepo.FindByCode(req.Code)
	if existing != nil {
		return nil, fmt.Errorf("collection code already exists")
	}

	collection := &models.Collection{
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		Code:        req.Code,
		Era:         req.Era,
		Material:    req.Material,
		Size:        req.Size,
		Source:      req.Source,
		Condition:   req.Condition,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		Status:      models.CollectionStatusActive,
		Tags:        req.Tags,
		MuseumID:    museumID,
	}

	if req.Status != "" {
		collection.Status = req.Status
	}

	if err := s.collectionRepo.Create(collection); err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	return collection, nil
}

func (s *CollectionService) GetByID(id uint) (*models.Collection, error) {
	collection, err := s.collectionRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("collection not found")
	}

	s.collectionRepo.IncrementViewCount(id)
	return collection, nil
}

func (s *CollectionService) Update(id uint, req *dto.CollectionRequest) error {
	collection, err := s.collectionRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("collection not found")
	}

	collection.Name = req.Name
	collection.CategoryID = req.CategoryID
	collection.Code = req.Code
	collection.Era = req.Era
	collection.Material = req.Material
	collection.Size = req.Size
	collection.Source = req.Source
	collection.Condition = req.Condition
	collection.Description = req.Description
	if req.ImageUrl != "" {
		collection.ImageUrl = req.ImageUrl
	}
	if req.Status != "" {
		collection.Status = req.Status
	}
	if req.Tags != "" {
		collection.Tags = req.Tags
	}

	return s.collectionRepo.Update(collection)
}

func (s *CollectionService) Delete(id uint) error {
	return s.collectionRepo.Delete(id)
}

func (s *CollectionService) List(query *dto.CollectionListQuery) ([]models.Collection, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}
	return s.collectionRepo.List(query)
}

func (s *CollectionService) CreateCategory(req *dto.CollectionCategoryRequest) (*models.CollectionCategory, error) {
	category := &models.CollectionCategory{
		Name:      req.Name,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
	}

	if err := s.collectionRepo.CreateCategory(category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (s *CollectionService) UpdateCategory(id uint, req *dto.CollectionCategoryRequest) error {
	category, err := s.collectionRepo.FindCategoryByID(id)
	if err != nil {
		return fmt.Errorf("category not found")
	}

	category.Name = req.Name
	category.ParentID = req.ParentID
	category.SortOrder = req.SortOrder

	return s.collectionRepo.UpdateCategory(category)
}

func (s *CollectionService) DeleteCategory(id uint) error {
	return s.collectionRepo.DeleteCategory(id)
}

func (s *CollectionService) ListCategories() ([]models.CollectionCategory, error) {
	return s.collectionRepo.ListCategories()
}

func (s *CollectionService) ListTags() ([]models.CollectionTag, error) {
	return s.collectionRepo.ListTags()
}

func (s *CollectionService) CreateTag(name string) (*models.CollectionTag, error) {
	tag := &models.CollectionTag{Name: name}
	if err := s.collectionRepo.CreateTag(tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}
	return tag, nil
}

func (s *CollectionService) BatchImport(museumID uint, collections []models.Collection) error {
	for i := range collections {
		collections[i].MuseumID = museumID
		if collections[i].Status == "" {
			collections[i].Status = models.CollectionStatusActive
		}
	}
	return s.collectionRepo.BatchCreate(collections)
}
