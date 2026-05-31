package handlers

import (
	"garden-planner/database"
	"garden-planner/middleware"
	"garden-planner/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateSeedExchangeRequest struct {
	Title         string `json:"title" binding:"required,max=200"`
	SeedName      string `json:"seed_name" binding:"required"`
	Description   string `json:"description"`
	ImageURLs     string `json:"image_urls"`
	Quantity      int    `json:"quantity"`
	ExchangeType  string `json:"exchange_type"`
	WantSeeds     string `json:"want_seeds"`
	Location      string `json:"location"`
}

type UpdateSeedExchangeRequest struct {
	Title        string `json:"title"`
	SeedName     string `json:"seed_name"`
	Description  string `json:"description"`
	ImageURLs    string `json:"image_urls"`
	Quantity     int    `json:"quantity"`
	ExchangeType string `json:"exchange_type"`
	WantSeeds    string `json:"want_seeds"`
	Location     string `json:"location"`
	Status       string `json:"status"`
}

type CreateExchangeOfferRequest struct {
	OfferSeeds string `json:"offer_seeds" binding:"required"`
	Message    string `json:"message"`
}

type UpdateExchangeOfferRequest struct {
	Status string `json:"status"`
}

func CreateSeedExchange(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateSeedExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exchange := models.SeedExchange{
		ID:            uuid.New(),
		OwnerID:       userID,
		Title:         req.Title,
		SeedName:      req.SeedName,
		Description:   req.Description,
		ImageURLs:     req.ImageURLs,
		Quantity:      req.Quantity,
		ExchangeType:  req.ExchangeType,
		WantSeeds:     req.WantSeeds,
		Location:      req.Location,
		Status:        "available",
	}

	if err := database.DB.Create(&exchange).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seed exchange"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Seed exchange created successfully",
		"exchange": exchange,
	})
}

func GetSeedExchanges(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	seedName := c.Query("seed_name")
	exchangeType := c.Query("exchange_type")
	status := c.DefaultQuery("status", "available")
	location := c.Query("location")
	search := c.Query("search")

	var exchanges []models.SeedExchange
	var total int64

	query := database.DB.Model(&models.SeedExchange{})

	if seedName != "" {
		query = query.Where("seed_name ILIKE ?", "%"+seedName+"%")
	}
	if exchangeType != "" {
		query = query.Where("exchange_type = ?", exchangeType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR seed_name ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	query.Preload("Owner").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&exchanges)

	c.JSON(http.StatusOK, gin.H{
		"exchanges": exchanges,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetSeedExchange(c *gin.Context) {
	id := c.Param("id")

	var exchange models.SeedExchange
	if err := database.DB.Preload("Owner").
		Preload("ExchangeOffers.Offerer").
		First(&exchange, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seed exchange not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exchange": exchange})
}

func UpdateSeedExchange(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var exchange models.SeedExchange
	if err := database.DB.Where("id = ? AND owner_id = ?", id, userID).First(&exchange).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seed exchange not found"})
		return
	}

	var req UpdateSeedExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != "" {
		exchange.Title = req.Title
	}
	if req.SeedName != "" {
		exchange.SeedName = req.SeedName
	}
	if req.Description != "" {
		exchange.Description = req.Description
	}
	if req.ImageURLs != "" {
		exchange.ImageURLs = req.ImageURLs
	}
	if req.Quantity > 0 {
		exchange.Quantity = req.Quantity
	}
	if req.ExchangeType != "" {
		exchange.ExchangeType = req.ExchangeType
	}
	if req.WantSeeds != "" {
		exchange.WantSeeds = req.WantSeeds
	}
	if req.Location != "" {
		exchange.Location = req.Location
	}
	if req.Status != "" {
		exchange.Status = req.Status
	}

	if err := database.DB.Save(&exchange).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seed exchange"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Seed exchange updated successfully",
		"exchange": exchange,
	})
}

func DeleteSeedExchange(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	result := database.DB.Where("id = ? AND owner_id = ?", id, userID).Delete(&models.SeedExchange{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete seed exchange"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seed exchange not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seed exchange deleted successfully"})
}

func CreateExchangeOffer(c *gin.Context) {
	userID := middleware.GetUserID(c)
	exchangeID := c.Param("id")

	exchangeUUID, err := uuid.Parse(exchangeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exchange ID"})
		return
	}

	var req CreateExchangeOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exchange models.SeedExchange
	if err := database.DB.First(&exchange, exchangeUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seed exchange not found"})
		return
	}

	if exchange.OwnerID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot make offer on your own exchange"})
		return
	}

	offer := models.ExchangeOffer{
		ID:             uuid.New(),
		SeedExchangeID: exchangeUUID,
		OffererID:      userID,
		OfferSeeds:     req.OfferSeeds,
		Message:        req.Message,
		Status:         "pending",
	}

	if err := database.DB.Create(&offer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exchange offer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Exchange offer created successfully",
		"offer":   offer,
	})
}

func GetExchangeOffers(c *gin.Context) {
	exchangeID := c.Param("id")

	var offers []models.ExchangeOffer
	database.DB.Where("seed_exchange_id = ?", exchangeID).
		Preload("Offerer").
		Order("created_at DESC").
		Find(&offers)

	c.JSON(http.StatusOK, gin.H{"offers": offers})
}

func UpdateExchangeOffer(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var offer models.ExchangeOffer
	if err := database.DB.First(&offer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exchange offer not found"})
		return
	}

	var exchange models.SeedExchange
	if err := database.DB.First(&exchange, offer.SeedExchangeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seed exchange not found"})
		return
	}

	if exchange.OwnerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only exchange owner can update offer status"})
		return
	}

	var req UpdateExchangeOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status != "" {
		offer.Status = req.Status

		if req.Status == "accepted" {
			exchange.Status = "exchanged"
			database.DB.Save(&exchange)
		}
	}

	if err := database.DB.Save(&offer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exchange offer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exchange offer updated successfully",
		"offer":   offer,
	})
}

func DeleteExchangeOffer(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	result := database.DB.Where("id = ? AND offerer_id = ?", id, userID).Delete(&models.ExchangeOffer{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exchange offer"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exchange offer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exchange offer deleted successfully"})
}
