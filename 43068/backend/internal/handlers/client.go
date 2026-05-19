package handlers

import (
	"net/http"
	"strconv"

	"freelancer-management/internal/database"
	"freelancer-management/internal/logger"
	"freelancer-management/internal/middleware"
	"freelancer-management/internal/models"
	"freelancer-management/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ClientHandler struct {
	db *gorm.DB
}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{db: database.GetDB()}
}

type CreateClientRequest struct {
	Name        string  `json:"name" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	Phone       string  `json:"phone"`
	Address     string  `json:"address"`
	Company     string  `json:"company"`
	ContractURL string `json:"contract_url"`
	DefaultRate float64 `json:"default_rate"`
}

type UpdateClientRequest struct {
	Name        string  `json:"name"`
	Email       string  `json:"email" binding:"omitempty,email"`
	Phone       string  `json:"phone"`
	Address     string  `json:"address"`
	Company     string  `json:"company"`
	ContractURL string `json:"contract_url"`
	DefaultRate float64 `json:"default_rate"`
}

func (h *ClientHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if !utils.IsValidEmail(req.Email) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid email format")
		return
	}

	client := models.Client{
		UserID:      userID,
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Address:     req.Address,
		Company:     req.Company,
		ContractURL: req.ContractURL,
		DefaultRate: req.DefaultRate,
	}

	if err := h.db.Create(&client).Error; err != nil {
		logger.LogError("Failed to create client: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create client")
		return
	}

	logger.LogOperation(userID, "create_client", "Client created: "+client.Name)
	utils.SuccessResponse(c, client)
}

func (h *ClientHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	var clients []models.Client
	var total int64

	query := h.db.Model(&models.Client{}).Where("user_id = ?", userID)
	query.Count(&total)
	query.Offset(offset).Limit(perPage).Find(&clients)

	utils.PaginatedSuccessResponse(c, clients, utils.PaginationMeta{
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	})
}

func (h *ClientHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var client models.Client
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Preload("Projects").First(&client).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Client not found")
		return
	}

	utils.SuccessResponse(c, client)
}

func (h *ClientHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var client models.Client
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&client).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Client not found")
		return
	}

	var req UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email != "" && !utils.IsValidEmail(req.Email) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid email format")
		return
	}

	if req.Name != "" {
		client.Name = req.Name
	}
	if req.Email != "" {
		client.Email = req.Email
	}
	client.Phone = req.Phone
	client.Address = req.Address
	client.Company = req.Company
	client.ContractURL = req.ContractURL
	if req.DefaultRate != 0 {
		client.DefaultRate = req.DefaultRate
	}

	if err := h.db.Save(&client).Error; err != nil {
		logger.LogError("Failed to update client: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update client")
		return
	}

	logger.LogOperation(userID, "update_client", "Client updated: "+client.Name)
	utils.SuccessResponse(c, client)
}

func (h *ClientHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var client models.Client
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&client).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Client not found")
		return
	}

	var projectCount int64
	h.db.Model(&models.Project{}).Where("client_id = ?", id).Count(&projectCount)
	if projectCount > 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Cannot delete client with associated projects")
		return
	}

	if err := h.db.Delete(&client).Error; err != nil {
		logger.LogError("Failed to delete client: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete client")
		return
	}

	logger.LogOperation(userID, "delete_client", "Client deleted: "+client.Name)
	utils.SuccessResponseWithMessage(c, "Client deleted successfully", nil)
}
