package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sports-league/models"
	"sports-league/pkg/auth"
)

type TeamHandler struct {
	db *gorm.DB
}

func NewTeamHandler(db *gorm.DB) *TeamHandler {
	return &TeamHandler{db: db}
}

type CreateTeamRequest struct {
	Name         string `json:"name" binding:"required"`
	Logo         string `json:"logo"`
	Description  string `json:"description"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
}

func (h *TeamHandler) List(c *gin.Context) {
	var teams []models.Team
	h.db.Preload("Captain").Preload("Players").Find(&teams)
	c.JSON(http.StatusOK, teams)
}

func (h *TeamHandler) Create(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*auth.Claims)
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	team := models.Team{
		Name:         req.Name,
		Logo:         req.Logo,
		CaptainID:    userClaims.UserID,
		Description:  req.Description,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
	}
	h.db.Create(&team)
	h.db.Model(&models.User{}).Where("id = ?", userClaims.UserID).Update("team_id", team.ID)
	c.JSON(http.StatusCreated, team)
}

func (h *TeamHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var team models.Team
	if err := h.db.Preload("Captain").Preload("Players").First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var team models.Team
	if err := h.db.First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.db.Save(&team)
	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.db.Delete(&models.Team{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

type AddPlayerRequest struct {
	Name      string `json:"name" binding:"required"`
	Number    int    `json:"number" binding:"required"`
	Position  string `json:"position"`
	BirthDate string `json:"birth_date"`
	UserID    *uint  `json:"user_id"`
}

func (h *TeamHandler) AddPlayer(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("id"))
	var req AddPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	player := models.Player{
		TeamID:    uint(teamID),
		Name:      req.Name,
		Number:    req.Number,
		Position:  req.Position,
		BirthDate: parseDate(req.BirthDate),
		UserID:    req.UserID,
		IsActive:  true,
	}
	h.db.Create(&player)
	c.JSON(http.StatusCreated, player)
}

func (h *TeamHandler) UpdatePlayer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("player_id"))
	var player models.Player
	if err := h.db.First(&player, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.db.Save(&player)
	c.JSON(http.StatusOK, player)
}

func (h *TeamHandler) DeletePlayer(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("player_id"))
	h.db.Delete(&models.Player{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

type RegisterTeamRequest struct {
	SeasonID  uint   `json:"season_id" binding:"required"`
	GroupName string `json:"group_name"`
}

func (h *TeamHandler) RegisterSeason(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("id"))
	var req RegisterTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var existing models.Registration
	if h.db.Where("team_id = ? AND season_id = ?", teamID, req.SeasonID).First(&existing).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "already registered"})
		return
	}
	reg := models.Registration{
		SeasonID:  req.SeasonID,
		TeamID:    uint(teamID),
		GroupName: req.GroupName,
		Status:    "pending",
	}
	h.db.Create(&reg)
	c.JSON(http.StatusCreated, reg)
}

func (h *TeamHandler) ListRegistrations(c *gin.Context) {
	seasonID := c.Query("season_id")
	var regs []models.Registration
	q := h.db.Preload("Team").Preload("Season")
	if seasonID != "" {
		q = q.Where("season_id = ?", seasonID)
	}
	q.Find(&regs)
	c.JSON(http.StatusOK, regs)
}

func (h *TeamHandler) ApproveRegistration(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("reg_id"))
	var reg models.Registration
	if err := h.db.First(&reg, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	reg.Status = "approved"
	h.db.Save(&reg)
	c.JSON(http.StatusOK, reg)
}
