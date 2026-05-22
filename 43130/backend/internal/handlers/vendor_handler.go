package handlers

import (
	"wedding-planner/internal/models"
	"wedding-planner/pkg/database"
	"wedding-planner/pkg/response"

	"github.com/gin-gonic/gin"
)

type VendorHandler struct{}

func NewVendorHandler() *VendorHandler {
	return &VendorHandler{}
}

type VendorRequest struct {
	Name          string  `json:"name" binding:"required,max=200"`
	Category      string  `json:"category" binding:"required,max=100"`
	ContactPerson string  `json:"contact_person"`
	Phone         string  `json:"phone"`
	Email         string  `json:"email"`
	Address       string  `json:"address"`
	Website       string  `json:"website"`
	ServiceArea   string  `json:"service_area"`
	PriceRange    string  `json:"price_range"`
	Rating        float64 `json:"rating"`
	Notes         string  `json:"notes"`
	WeddingID     *uint   `json:"wedding_id"`
}

type VendorReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Content string `json:"content"`
}

func (h *VendorHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req VendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	vendor := models.Vendor{
		UserID:        userID,
		WeddingID:     req.WeddingID,
		Name:          req.Name,
		Category:      req.Category,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
		Email:         req.Email,
		Address:       req.Address,
		Website:       req.Website,
		ServiceArea:   req.ServiceArea,
		PriceRange:    req.PriceRange,
		Rating:        req.Rating,
		Notes:         req.Notes,
		Status:        "active",
	}

	if err := db.Create(&vendor).Error; err != nil {
		response.InternalError(c, "Failed to create vendor")
		return
	}

	response.Created(c, vendor)
}

func (h *VendorHandler) GetList(c *gin.Context) {
	userID := c.GetUint("user_id")

	db := database.GetDB()

	var vendors []models.Vendor
	var total int64

	page := c.GetInt("page")
	pageSize := c.GetInt("page_size")
	search := c.Query("search")
	category := c.Query("category")
	weddingID := c.Query("wedding_id")

	query := db.Model(&models.Vendor{}).Where("user_id = ?", userID)

	if search != "" {
		query = query.Where("name LIKE ? OR contact_person LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if weddingID != "" {
		query = query.Where("wedding_id = ?", weddingID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&vendors)

	response.Paginated(c, vendors, total, page, pageSize)
}

func (h *VendorHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.GetUint("id")

	db := database.GetDB()

	var vendor models.Vendor
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&vendor).Error; err != nil {
		response.NotFound(c, "Vendor not found")
		return
	}

	var reviews []models.VendorReview
	db.Where("vendor_id = ?", id).Order("created_at DESC").Limit(10).Find(&reviews)

	response.Success(c, gin.H{
		"vendor":  vendor,
		"reviews": reviews,
	})
}

func (h *VendorHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.GetUint("id")

	var req VendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	var vendor models.Vendor
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&vendor).Error; err != nil {
		response.NotFound(c, "Vendor not found")
		return
	}

	updates := map[string]interface{}{
		"name":           req.Name,
		"category":       req.Category,
		"contact_person": req.ContactPerson,
		"phone":          req.Phone,
		"email":          req.Email,
		"address":        req.Address,
		"website":        req.Website,
		"service_area":   req.ServiceArea,
		"price_range":    req.PriceRange,
		"notes":          req.Notes,
		"wedding_id":     req.WeddingID,
	}

	db.Model(&vendor).Updates(updates)

	response.Success(c, vendor)
}

func (h *VendorHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.GetUint("id")

	db := database.GetDB()

	var vendor models.Vendor
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&vendor).Error; err != nil {
		response.NotFound(c, "Vendor not found")
		return
	}

	db.Delete(&vendor)

	response.Success(c, gin.H{"message": "Vendor deleted successfully"})
}

func (h *VendorHandler) AddReview(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.GetUint("id")

	var req VendorReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request parameters")
		return
	}

	db := database.GetDB()

	var vendor models.Vendor
	if err := db.Where("id = ?", id).First(&vendor).Error; err != nil {
		response.NotFound(c, "Vendor not found")
		return
	}

	review := models.VendorReview{
		VendorID: id,
		UserID:   userID,
		Rating:   req.Rating,
		Content:  req.Content,
	}

	if err := db.Create(&review).Error; err != nil {
		response.InternalError(c, "Failed to create review")
		return
	}

	var avgRating float64
	var reviewCount int64
	db.Model(&models.VendorReview{}).Where("vendor_id = ?", id).Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)
	db.Model(&models.VendorReview{}).Where("vendor_id = ?", id).Count(&reviewCount)

	db.Model(&vendor).Updates(map[string]interface{}{
		"rating":       avgRating,
		"review_count": reviewCount,
	})

	response.Created(c, review)
}

func (h *VendorHandler) GetCategories(c *gin.Context) {
	categories := []string{
		"摄影摄像", "花艺布置", "司仪主持", "化妆造型", "场地布置",
		"婚车租赁", "婚礼蛋糕", "婚纱礼服", "婚礼策划", "婚礼跟拍",
		"婚礼音乐", "婚礼顾问", "蜜月旅行", "其他",
	}

	response.Success(c, categories)
}
