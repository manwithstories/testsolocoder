package handlers

import (
	"wedding-planner/internal/models"
	"wedding-planner/pkg/database"
	"wedding-planner/pkg/response"

	"github.com/gin-gonic/gin"
)

type WeddingHandler struct{}

func NewWeddingHandler() *WeddingHandler {
	return &WeddingHandler{}
}

type WeddingRequest struct {
	Title         string  `json:"title" binding:"required,max=200"`
	GroomName     string  `json:"groom_name" binding:"required,max=100"`
	BrideName     string  `json:"bride_name" binding:"required,max=100"`
	WeddingDate   string  `json:"wedding_date" binding:"required"`
	Budget        float64 `json:"budget"`
	Style         string  `json:"style"`
	ThemeColor    string  `json:"theme_color"`
	GuestCount    int     `json:"guest_count"`
	Venue         string  `json:"venue"`
	VenueAddress  string  `json:"venue_address"`
	VenueCapacity int     `json:"venue_capacity"`
	Description   string  `json:"description"`
}

type WeddingListQuery struct {
	Search string `form:"search"`
	Status string `form:"status"`
	Sort   string `form:"sort"`
}

func (h *WeddingHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req WeddingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	weddingDate, err := parseDate(req.WeddingDate)
	if err != nil {
		response.BadRequest(c, "Invalid wedding date format")
		return
	}

	wedding := models.Wedding{
		UserID:        userID,
		Title:         req.Title,
		GroomName:     req.GroomName,
		BrideName:     req.BrideName,
		WeddingDate:   weddingDate,
		Budget:        req.Budget,
		Style:         req.Style,
		ThemeColor:    req.ThemeColor,
		GuestCount:    req.GuestCount,
		Venue:         req.Venue,
		VenueAddress:  req.VenueAddress,
		VenueCapacity: req.VenueCapacity,
		Description:   req.Description,
		Status:        "planning",
	}

	if err := db.Create(&wedding).Error; err != nil {
		response.InternalError(c, "Failed to create wedding")
		return
	}

	response.Created(c, wedding)
}

func (h *WeddingHandler) GetList(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole, _ := c.Get("user_role")

	db := database.GetDB()

	var weddings []models.Wedding
	var total int64

	page := c.GetInt("page")
	pageSize := c.GetInt("page_size")
	search := c.Query("search")
	status := c.Query("status")

	query := db.Model(&models.Wedding{})

	if userRole != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if search != "" {
		query = query.Where("title LIKE ? OR groom_name LIKE ? OR bride_name LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("wedding_date DESC").Find(&weddings)

	response.Paginated(c, weddings, total, page, pageSize)
}

func (h *WeddingHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole, _ := c.Get("user_role")
	id := c.GetUint("id")

	db := database.GetDB()

	var wedding models.Wedding
	query := db.Where("id = ?", id)
	if userRole != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&wedding).Error; err != nil {
		response.NotFound(c, "Wedding not found")
		return
	}

	response.Success(c, wedding)
}

func (h *WeddingHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole, _ := c.Get("user_role")
	id := c.GetUint("id")

	var req WeddingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	var wedding models.Wedding
	query := db.Where("id = ?", id)
	if userRole != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&wedding).Error; err != nil {
		response.NotFound(c, "Wedding not found")
		return
	}

	weddingDate, err := parseDate(req.WeddingDate)
	if err != nil {
		response.BadRequest(c, "Invalid wedding date format")
		return
	}

	updates := map[string]interface{}{
		"title":          req.Title,
		"groom_name":     req.GroomName,
		"bride_name":     req.BrideName,
		"wedding_date":   weddingDate,
		"budget":         req.Budget,
		"style":          req.Style,
		"theme_color":    req.ThemeColor,
		"guest_count":    req.GuestCount,
		"venue":          req.Venue,
		"venue_address":  req.VenueAddress,
		"venue_capacity": req.VenueCapacity,
		"description":    req.Description,
	}

	if err := db.Model(&wedding).Updates(updates).Error; err != nil {
		response.InternalError(c, "Failed to update wedding")
		return
	}

	response.Success(c, wedding)
}

func (h *WeddingHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole, _ := c.Get("user_role")
	id := c.GetUint("id")

	db := database.GetDB()

	var wedding models.Wedding
	query := db.Where("id = ?", id)
	if userRole != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&wedding).Error; err != nil {
		response.NotFound(c, "Wedding not found")
		return
	}

	if err := db.Delete(&wedding).Error; err != nil {
		response.InternalError(c, "Failed to delete wedding")
		return
	}

	response.Success(c, gin.H{"message": "Wedding deleted successfully"})
}

func (h *WeddingHandler) UpdateStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole, _ := c.Get("user_role")
	id := c.GetUint("id")

	type StatusRequest struct {
		Status string `json:"status" binding:"required,oneof=planning confirmed cancelled completed"`
	}

	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid status")
		return
	}

	db := database.GetDB()

	var wedding models.Wedding
	query := db.Where("id = ?", id)
	if userRole != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&wedding).Error; err != nil {
		response.NotFound(c, "Wedding not found")
		return
	}

	wedding.Status = req.Status
	db.Save(&wedding)

	response.Success(c, wedding)
}
