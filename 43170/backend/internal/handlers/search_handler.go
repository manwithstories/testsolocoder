package handlers

import (
	"net/http"
	"photo-rental/internal/models"
	"photo-rental/internal/utils"
	"photo-rental/pkg/database"
	redispkg "photo-rental/pkg/redis"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct{}

type SearchRequest struct {
	Category   string  `json:"category"`
	Brand      string  `json:"brand"`
	MinPrice   float64 `json:"minPrice"`
	MaxPrice   float64 `json:"maxPrice"`
	Status     string  `json:"status"`
	StartDate  string  `json:"startDate"`
	EndDate    string  `json:"endDate"`
	Page       int     `json:"page"`
	PageSize   int     `json:"pageSize"`
	SortBy     string  `json:"sortBy"`
	SortOrder  string  `json:"sortOrder"`
}

func NewSearchHandler() *SearchHandler {
	return &SearchHandler{}
}

func (h *SearchHandler) SearchEquipments(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	db := database.DB.Model(&models.Equipment{})

	if req.Category != "" {
		db = db.Where("category = ?", req.Category)
	}

	if req.Brand != "" {
		db = db.Where("brand = ?", req.Brand)
	}

	if req.MinPrice > 0 {
		db = db.Where("daily_rent >= ?", req.MinPrice)
	}

	if req.MaxPrice > 0 {
		db = db.Where("daily_rent <= ?", req.MaxPrice)
	}

	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	} else {
		db = db.Where("status = ?", "available")
	}

	var total int64
	db.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	db = db.Offset(offset).Limit(req.PageSize)

	orderClause := req.SortBy + " " + req.SortOrder
	db = db.Order(orderClause)

	var equipments []models.Equipment
	if err := db.Preload("Images").Find(&equipments).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search equipments")
		return
	}

	if req.StartDate != "" && req.EndDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid start date format")
			return
		}

		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid end date format")
			return
		}

		var availableEquipments []models.Equipment
		for _, eq := range equipments {
			conflicts, err := redispkg.CheckEquipmentConflict(eq.ID, startDate, endDate)
			if err != nil {
				continue
			}
			if len(conflicts) == 0 {
				availableEquipments = append(availableEquipments, eq)
			}
		}
		equipments = availableEquipments
	}

	utils.PaginatedSuccessResponse(c, equipments, total, req.Page, req.PageSize)
}

func (h *SearchHandler) SearchOrders(c *gin.Context) {
	userID := c.GetUint("userId")
	userRole := c.GetString("role")

	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	db := database.DB.Model(&models.Order{})

	switch userRole {
	case "owner":
		db = db.Where("owner_id = ?", userID)
	case "renter":
		db = db.Where("renter_id = ?", userID)
	}

	if req.Category != "" {
		db = db.Joins("JOIN equipments ON equipments.id = orders.equipment_id").
			Where("equipments.category = ?", req.Category)
	}

	if req.Brand != "" {
		db = db.Joins("JOIN equipments ON equipments.id = orders.equipment_id").
			Where("equipments.brand = ?", req.Brand)
	}

	if req.MinPrice > 0 {
		db = db.Where("total_rent >= ?", req.MinPrice)
	}

	if req.MaxPrice > 0 {
		db = db.Where("total_rent <= ?", req.MaxPrice)
	}

	if req.Status != "" {
		db = db.Where("orders.status = ?", req.Status)
	}

	if req.StartDate != "" {
		db = db.Where("start_date >= ?", req.StartDate)
	}

	if req.EndDate != "" {
		db = db.Where("end_date <= ?", req.EndDate)
	}

	var total int64
	db.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	db = db.Offset(offset).Limit(req.PageSize)

	orderClause := "orders." + req.SortBy + " " + req.SortOrder
	db = db.Order(orderClause)

	var orders []models.Order
	if err := db.Preload("Equipment.Images").Preload("Renter").Preload("Owner").Find(&orders).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search orders")
		return
	}

	for i := range orders {
		orders[i].Renter.Password = ""
		orders[i].Renter.IDCard = ""
		orders[i].Owner.Password = ""
		orders[i].Owner.IDCard = ""
	}

	utils.PaginatedSuccessResponse(c, orders, total, req.Page, req.PageSize)
}

func (h *SearchHandler) GetCategories(c *gin.Context) {
	var categories []string
	database.DB.Model(&models.Equipment{}).Distinct("category").Pluck("category", &categories)

	utils.SuccessResponse(c, categories)
}

func (h *SearchHandler) GetBrands(c *gin.Context) {
	var brands []string
	database.DB.Model(&models.Equipment{}).Distinct("brand").Pluck("brand", &brands)

	utils.SuccessResponse(c, brands)
}

func (h *SearchHandler) GetEquipmentReservedDates(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid equipment ID")
		return
	}

	dates, err := redispkg.GetEquipmentReservedDates(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get reserved dates")
		return
	}

	utils.SuccessResponse(c, dates)
}
