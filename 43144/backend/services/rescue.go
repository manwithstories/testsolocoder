package services

import (
	"errors"
	"fmt"
	"time"

	"pet-adoption-platform/database"
	"pet-adoption-platform/models"
)

func GetRescueStationByID(id uint) (*models.RescueStation, error) {
	var rescue models.RescueStation
	if err := database.DB.First(&rescue, id).Error; err != nil {
		return nil, err
	}
	return &rescue, nil
}

func ListRescueStations(query *models.RescueListQuery) ([]models.RescueStation, int64, error) {
	var rescues []models.RescueStation
	var total int64

	db := database.DB.Model(&models.RescueStation{})

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		db = db.Where("name ILIKE ? OR contact_person ILIKE ? OR contact_email ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Order("created_at DESC").
		Offset(offset).Limit(query.PageSize).
		Find(&rescues).Error; err != nil {
		return nil, 0, err
	}

	return rescues, total, nil
}

func ReviewRescueStation(id uint, reviewerID uint, req *models.ReviewRescueRequest) (*models.RescueStation, error) {
	rescue, err := GetRescueStationByID(id)
	if err != nil {
		return nil, errors.New("rescue station not found")
	}

	if rescue.Status != models.RescueStatusPending {
		return nil, errors.New("rescue station is not in pending status")
	}

	now := time.Now()
	if req.Action == "approve" {
		rescue.Status = models.RescueStatusApproved
		rescue.VerifiedBy = &reviewerID
		rescue.VerifiedAt = &now
	} else if req.Action == "reject" {
		rescue.Status = models.RescueStatusRejected
		rescue.VerifiedBy = &reviewerID
		rescue.VerifiedAt = &now
		rescue.RejectReason = req.RejectReason
	}

	if err := database.DB.Save(rescue).Error; err != nil {
		return nil, err
	}

	return rescue, nil
}

func GetRescueStats(rescueID uint) (*models.RescueStats, error) {
	stats := &models.RescueStats{}

	database.DB.Model(&models.Pet{}).Where("rescue_id = ?", rescueID).Count(&stats.TotalPets)
	database.DB.Model(&models.Pet{}).Where("rescue_id = ? AND status = ?", rescueID, models.PetStatusAdoptable).Count(&stats.AdoptablePets)
	database.DB.Model(&models.Pet{}).Where("rescue_id = ? AND status = ?", rescueID, models.PetStatusAdopted).Count(&stats.AdoptedPets)
	database.DB.Model(&models.Pet{}).Where("rescue_id = ? AND status = ?", rescueID, models.PetStatusTreatment).Count(&stats.TreatmentPets)
	database.DB.Model(&models.Pet{}).Where("rescue_id = ? AND status = ?", rescueID, models.PetStatusDeceased).Count(&stats.DeceasedPets)

	database.DB.Model(&models.AdoptionApplication{}).Where("rescue_id = ?", rescueID).Count(&stats.TotalAdoptions)
	database.DB.Model(&models.AdoptionApplication{}).Where("rescue_id = ? AND status = ?", rescueID, models.AdoptionStatusPending).Count(&stats.PendingApplications)

	database.DB.Model(&models.FollowUpRecord{}).Where("rescue_id = ?", rescueID).Count(&stats.TotalFollowUps)
	database.DB.Model(&models.FollowUpRecord{}).Where("rescue_id = ? AND health_status != ?", rescueID, "").Count(&stats.CompletedFollowUps)

	if stats.TotalPets > 0 {
		stats.AdoptionRate = float64(stats.AdoptedPets) / float64(stats.TotalPets) * 100
	}
	if stats.TotalFollowUps > 0 {
		stats.FollowUpRate = float64(stats.CompletedFollowUps) / float64(stats.TotalFollowUps) * 100
	}

	return stats, nil
}

func GetAllRescuesStats() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	rows, err := database.DB.Table("rescue_stations r").
		Select("r.id, r.name, r.status, " +
			"(SELECT COUNT(*) FROM pets p WHERE p.rescue_id = r.id) as total_pets, " +
			"(SELECT COUNT(*) FROM pets p WHERE p.rescue_id = r.id AND p.status = 'adopted') as adopted_pets, " +
			"(SELECT COUNT(*) FROM adoption_applications a WHERE a.rescue_id = r.id) as total_adoptions").
		Where("r.deleted_at IS NULL").
		Order("r.created_at DESC").
		Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id uint
		var name, status string
		var totalPets, adoptedPets, totalAdoptions int64
		rows.Scan(&id, &name, &status, &totalPets, &adoptedPets, &totalAdoptions)

		result := map[string]interface{}{
			"id":              id,
			"name":            name,
			"status":          status,
			"total_pets":      totalPets,
			"adopted_pets":    adoptedPets,
			"total_adoptions": totalAdoptions,
		}
		results = append(results, result)
	}

	return results, nil
}
