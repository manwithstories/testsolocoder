package services

import (
	"errors"
	"fmt"

	"pet-adoption-platform/database"
	"pet-adoption-platform/models"
	"pet-adoption-platform/utils"
)

func GenerateArchiveNumber() string {
	var count int64
	database.DB.Model(&models.Pet{}).Count(&count)
	return fmt.Sprintf("PET-%06d", count+1)
}

func CreatePet(req *models.CreatePetRequest, rescueID uint) (*models.Pet, error) {
	pet := &models.Pet{
		ArchiveNumber:  GenerateArchiveNumber(),
		Name:           req.Name,
		Species:        req.Species,
		Breed:          req.Breed,
		Age:            req.Age,
		Gender:         models.PetGender(req.Gender),
		Weight:         req.Weight,
		Color:          req.Color,
		Description:    req.Description,
		Status:         models.PetStatusAdoptable,
		HealthStatus:   req.HealthStatus,
		Vaccinated:     req.Vaccinated,
		Neutered:       req.Neutered,
		RescueID:       rescueID,
		FoundLocation:  req.FoundLocation,
		FoundDate:      req.FoundDate,
		MedicalHistory: req.MedicalHistory,
		Personality:    req.Personality,
		SpecialNeeds:   req.SpecialNeeds,
	}

	if err := database.DB.Create(pet).Error; err != nil {
		return nil, fmt.Errorf("failed to create pet: %w", err)
	}

	return pet, nil
}

func GetPetByID(id uint) (*models.Pet, error) {
	var pet models.Pet
	if err := database.DB.Preload("Rescue").Preload("Adopter").First(&pet, id).Error; err != nil {
		return nil, err
	}
	return &pet, nil
}

func UpdatePet(id uint, req *models.UpdatePetRequest) (*models.Pet, error) {
	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Species != "" {
		updates["species"] = req.Species
	}
	if req.Breed != "" {
		updates["breed"] = req.Breed
	}
	if req.Age != "" {
		updates["age"] = req.Age
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}
	if req.Color != "" {
		updates["color"] = req.Color
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.HealthStatus != "" {
		updates["health_status"] = req.HealthStatus
	}
	if req.Vaccinated != nil {
		updates["vaccinated"] = *req.Vaccinated
	}
	if req.Neutered != nil {
		updates["neutered"] = *req.Neutered
	}
	if req.FoundLocation != "" {
		updates["found_location"] = req.FoundLocation
	}
	if req.FoundDate != nil {
		updates["found_date"] = req.FoundDate
	}
	if req.MedicalHistory != "" {
		updates["medical_history"] = req.MedicalHistory
	}
	if req.Personality != "" {
		updates["personality"] = req.Personality
	}
	if req.SpecialNeeds != "" {
		updates["special_needs"] = req.SpecialNeeds
	}

	result := database.DB.Model(&models.Pet{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	return GetPetByID(id)
}

func DeletePet(id uint) error {
	return database.DB.Delete(&models.Pet{}, id).Error
}

func ListPets(query *models.PetListQuery) ([]models.Pet, int64, error) {
	var pets []models.Pet
	var total int64

	db := database.DB.Model(&models.Pet{})

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Species != "" {
		db = db.Where("species = ?", query.Species)
	}
	if query.Gender != "" {
		db = db.Where("gender = ?", query.Gender)
	}
	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		db = db.Where("name ILIKE ? OR breed ILIKE ? OR description ILIKE ? OR archive_number ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}
	if query.RescueID > 0 {
		db = db.Where("rescue_id = ?", query.RescueID)
	}
	if query.Breed != "" {
		db = db.Where("breed ILIKE ?", "%"+query.Breed+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Rescue").Preload("Adopter").
		Order("created_at DESC").
		Offset(offset).Limit(query.PageSize).
		Find(&pets).Error; err != nil {
		return nil, 0, err
	}

	return pets, total, nil
}

func UpdatePetPhotos(id uint, photos string) error {
	return database.DB.Model(&models.Pet{}).Where("id = ?", id).Update("photos", photos).Error
}

func UpdatePetVideos(id uint, videos string) error {
	return database.DB.Model(&models.Pet{}).Where("id = ?", id).Update("videos", videos).Error
}

func GetPetAdoptionHistory(petID uint) ([]models.AdoptionApplication, error) {
	var applications []models.AdoptionApplication
	err := database.DB.Preload("Adopter").
		Where("pet_id = ?", petID).
		Order("created_at DESC").
		Find(&applications).Error
	return applications, err
}

func ValidatePetOwnership(petID, rescueID uint) (*models.Pet, error) {
	pet, err := GetPetByID(petID)
	if err != nil {
		return nil, errors.New("pet not found")
	}
	if pet.RescueID != rescueID {
		return nil, errors.New("you don't have permission to manage this pet")
	}
	return pet, nil
}

func UpdatePetStatus(petID uint, status models.PetStatus) error {
	return database.DB.Model(&models.Pet{}).Where("id = ?", petID).Update("status", status).Error
}

func InitPetService() {
	utils.Logger.Info("Pet service initialized")
}
