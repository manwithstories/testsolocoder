package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"property-management/internal/database"
	"property-management/internal/model"
	redisclient "property-management/internal/redis"
	"property-management/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PropertyHandler struct {
	UploadDir string
}

func NewPropertyHandler(uploadDir string) *PropertyHandler {
	return &PropertyHandler{UploadDir: uploadDir}
}

type PropertyRequest struct {
	Title       string  `json:"title" binding:"required"`
	Community   string  `json:"community"`
	Address     string  `json:"address"`
	Area        float64 `json:"area"`
	Layout      string  `json:"layout"`
	Floor       string  `json:"floor"`
	Rent        float64 `json:"rent" binding:"required"`
	Deposit     float64 `json:"deposit"`
	PaymentType string  `json:"paymentType"`
	Description string  `json:"description"`
	Region      string  `json:"region"`
	Building    string  `json:"building"`
	RoomNo      string  `json:"roomNo"`
	ImageUrls   []string `json:"imageUrls"`
	FacilityIDs []uint   `json:"facilityIds"`
}

func (h *PropertyHandler) Create(c *gin.Context) {
	var req PropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	userID := utils.GetUserIDFromContext(c)
	property := model.Property{
		Title:       req.Title,
		Community:   req.Community,
		Address:     req.Address,
		Area:        req.Area,
		Layout:      req.Layout,
		Floor:       req.Floor,
		Rent:        req.Rent,
		Deposit:     req.Deposit,
		PaymentType: req.PaymentType,
		Description: req.Description,
		Region:      req.Region,
		Building:    req.Building,
		RoomNo:      req.RoomNo,
		Status:      1,
		OwnerID:     userID,
	}

	if err := database.DB.Create(&property).Error; err != nil {
		utils.Error(c, 500, "Failed to create property")
		return
	}

	for i, url := range req.ImageUrls {
		database.DB.Create(&model.PropertyImage{
			PropertyID: property.ID,
			URL:        url,
			Sort:       i,
		})
	}

	if len(req.FacilityIDs) > 0 {
		var facilities []model.Facility
		database.DB.Where("id IN ?", req.FacilityIDs).Find(&facilities)
		database.DB.Model(&property).Association("Facilities").Append(facilities)
	}

	redisclient.CachePropertyStatus(property.ID, property.Status)

	utils.Success(c, property)
}

func (h *PropertyHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var property model.Property
	if err := database.DB.First(&property, id).Error; err != nil {
		utils.Error(c, 404, "Property not found")
		return
	}

	var req PropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	updates := map[string]interface{}{
		"title":       req.Title,
		"community":   req.Community,
		"address":     req.Address,
		"area":        req.Area,
		"layout":      req.Layout,
		"floor":       req.Floor,
		"rent":        req.Rent,
		"deposit":     req.Deposit,
		"paymentType": req.PaymentType,
		"description": req.Description,
		"region":      req.Region,
		"building":    req.Building,
		"room_no":     req.RoomNo,
	}

	database.DB.Model(&property).Updates(updates)

	database.DB.Where("property_id = ?", id).Delete(&model.PropertyImage{})
	for i, url := range req.ImageUrls {
		database.DB.Create(&model.PropertyImage{
			PropertyID: property.ID,
			URL:        url,
			Sort:       i,
		})
	}

	if len(req.FacilityIDs) > 0 {
		var facilities []model.Facility
		database.DB.Where("id IN ?", req.FacilityIDs).Find(&facilities)
		database.DB.Model(&property).Association("Facilities").Replace(facilities)
	}

	redisclient.CachePropertyStatus(property.ID, property.Status)

	utils.Success(c, nil)
}

func (h *PropertyHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&model.Property{}, id)
	database.DB.Where("property_id = ?", id).Delete(&model.PropertyImage{})
	utils.Success(c, nil)
}

func (h *PropertyHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	database.DB.Model(&model.Property{}).Where("id = ?", id).Update("status", req.Status)
	redisclient.CachePropertyStatus(uint(id), req.Status)
	utils.Success(c, nil)
}

func (h *PropertyHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query := database.DB.Model(&model.Property{})

	if region := c.Query("region"); region != "" {
		query = query.Where("region = ?", region)
	}
	if layout := c.Query("layout"); layout != "" {
		query = query.Where("layout = ?", layout)
	}
	if minRent := c.Query("minRent"); minRent != "" {
		rent, _ := strconv.ParseFloat(minRent, 64)
		query = query.Where("rent >= ?", rent)
	}
	if maxRent := c.Query("maxRent"); maxRent != "" {
		rent, _ := strconv.ParseFloat(maxRent, 64)
		query = query.Where("rent <= ?", rent)
	}
	if status := c.Query("status"); status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title LIKE ? OR community LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")
	query = query.Order(sortBy + " " + sortOrder)

	var total int64
	query.Count(&total)

	var properties []model.Property
	offset := (page - 1) * pageSize
	query.Preload("Images").Preload("Facilities").Preload("Owner").
		Offset(offset).Limit(pageSize).Find(&properties)

	utils.Success(c, gin.H{
		"list":     properties,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *PropertyHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var property model.Property
	if err := database.DB.Preload("Images").Preload("Facilities").Preload("Owner").
		First(&property, id).Error; err != nil {
		utils.Error(c, 404, "Property not found")
		return
	}
	utils.Success(c, property)
}

func (h *PropertyHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, 400, "No file received")
		return
	}

	if h.UploadDir == "" {
		h.UploadDir = "./uploads"
	}
	os.MkdirAll(h.UploadDir, 0755)

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixMilli(), ext)
	filepath := filepath.Join(h.UploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		utils.Error(c, 500, "Failed to save file")
		return
	}

	url := "/uploads/" + filename
	utils.Success(c, gin.H{"url": url})
}

func (h *PropertyHandler) ListFacilities(c *gin.Context) {
	var facilities []model.Facility
	database.DB.Find(&facilities)
	utils.Success(c, facilities)
}

func (h *PropertyHandler) MyProperties(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query := database.DB.Model(&model.Property{}).Where("owner_id = ?", userID)

	var total int64
	query.Count(&total)

	var properties []model.Property
	offset := (page - 1) * pageSize
	query.Preload("Images").Preload("Facilities").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&properties)

	utils.Success(c, gin.H{
		"list":     properties,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}
