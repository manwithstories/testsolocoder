package handlers

import (
	"errand-service/internal/models"
	"errand-service/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type SubmitVerificationRequest struct {
	IDCardNo    string `json:"id_card_no" binding:"required"`
	IDCardName  string `json:"id_card_name" binding:"required"`
	IDCardFront string `json:"id_card_front" binding:"required"`
	IDCardBack  string `json:"id_card_back" binding:"required"`
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
		return
	}

	var courierProfile *models.CourierProfile
	if user.Role == models.RoleCourier {
		courierProfile = &models.CourierProfile{}
		h.db.Where("user_id = ?", userID).First(courierProfile)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":              user.ID,
			"phone":           user.Phone,
			"nickname":        user.Nickname,
			"avatar":          user.Avatar,
			"role":            user.Role,
			"status":          user.Status,
			"balance":         user.Balance,
			"rating":          user.Rating,
			"order_count":     user.OrderCount,
			"created_at":      user.CreatedAt,
			"courier_profile": courierProfile,
		},
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request parameters"})
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "No changes made"})
		return
	}

	result := h.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	if result.Error != nil {
		logger.Errorf("Failed to update profile: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Profile updated successfully"})
}

func (h *UserHandler) SubmitVerification(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req SubmitVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request parameters", "error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
		return
	}

	var existing models.Verification
	h.db.Where("user_id = ? AND status = ?", userID, "pending").First(&existing)
	if existing.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "You already have a pending verification request"})
		return
	}

	verification := models.Verification{
		UserID:      userID,
		IDCardNo:    req.IDCardNo,
		IDCardName:  req.IDCardName,
		IDCardFront: req.IDCardFront,
		IDCardBack:  req.IDCardBack,
		Status:      "pending",
	}

	if err := h.db.Create(&verification).Error; err != nil {
		logger.Errorf("Failed to submit verification: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to submit verification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Verification request submitted successfully",
		"data":    verification,
	})
}
