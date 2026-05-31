package handlers

import (
	"net/http"
	"strconv"

	"print3d-platform/internal/middleware"
	"print3d-platform/internal/models"
	"print3d-platform/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientIP := c.ClientIP()
	resp, err := h.userService.Login(c.Request.Context(), &req, clientIP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.userService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	user, err := h.userService.GetUser(c.Request.Context(), authUser.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	var req struct {
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
		RealName string `json:"real_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), authUser.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Phone = req.Phone
	user.Avatar = req.Avatar
	user.RealName = req.RealName

	err = h.userService.UpdateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "user": user})
}

func (h *UserHandler) UpdateDesignerProfile(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	if authUser.Role != models.RoleDesigner {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only designers can update designer profile"})
		return
	}

	var req struct {
		Nickname        string   `json:"nickname"`
		Bio             string   `json:"bio"`
		PortfolioURL    string   `json:"portfolio_url"`
		Specialties     []string `json:"specialties"`
		ExperienceYears int      `json:"experience_years"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), authUser.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.DesignerProfile == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Designer profile not found"})
		return
	}

	user.DesignerProfile.Nickname = req.Nickname
	user.DesignerProfile.Bio = req.Bio
	user.DesignerProfile.PortfolioURL = req.PortfolioURL
	user.DesignerProfile.Specialties = req.Specialties
	user.DesignerProfile.ExperienceYears = req.ExperienceYears

	err = h.userService.UpdateDesignerProfile(c.Request.Context(), user.DesignerProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update designer profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Designer profile updated successfully"})
}

func (h *UserHandler) UpdatePrinterProfile(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	if authUser.Role != models.RolePrinter {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only printers can update printer profile"})
		return
	}

	var req struct {
		CompanyName        string   `json:"company_name"`
		Address            string   `json:"address"`
		MaxPrintSize       string   `json:"max_print_size"`
		SupportedMaterials []string `json:"supported_materials"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), authUser.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.PrinterProfile == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Printer profile not found"})
		return
	}

	user.PrinterProfile.CompanyName = req.CompanyName
	user.PrinterProfile.Address = req.Address
	user.PrinterProfile.MaxPrintSize = req.MaxPrintSize
	user.PrinterProfile.SupportedMaterials = req.SupportedMaterials

	err = h.userService.UpdatePrinterProfile(c.Request.Context(), user.PrinterProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update printer profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Printer profile updated successfully"})
}

func (h *UserHandler) ListDesigners(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	designers, total, err := h.userService.ListDesigners(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list designers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  designers,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *UserHandler) ListPrinters(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	printers, total, err := h.userService.ListPrinters(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list printers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  printers,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *UserHandler) GetUserStats(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	stats, err := h.userService.GetUserStats(c.Request.Context(), authUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *UserHandler) GetNotifications(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	notifications, total, err := h.userService.GetNotifications(c.Request.Context(), authUser.UserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  notifications,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *UserHandler) MarkNotificationRead(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	err = h.userService.MarkNotificationRead(c.Request.Context(), id, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (h *UserHandler) GetTransactions(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	transactions, total, err := h.userService.GetTransactions(c.Request.Context(), authUser.UserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  transactions,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}
