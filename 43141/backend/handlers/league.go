package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sports-league/models"
	"sports-league/pkg/auth"
)

type LeagueHandler struct {
	db *gorm.DB
}

func NewLeagueHandler(db *gorm.DB) *LeagueHandler {
	return &LeagueHandler{db: db}
}

type CreateLeagueRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Sport       string `json:"sport"`
}

func (h *LeagueHandler) List(c *gin.Context) {
	var leagues []models.League
	h.db.Preload("Seasons").Find(&leagues)
	c.JSON(http.StatusOK, leagues)
}

func (h *LeagueHandler) Create(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*auth.Claims)
	var req CreateLeagueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	league := models.League{
		Name:        req.Name,
		Description: req.Description,
		Sport:       req.Sport,
		Status:      "planning",
		CreatedBy:   userClaims.UserID,
	}
	h.db.Create(&league)
	c.JSON(http.StatusCreated, league)
}

func (h *LeagueHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var league models.League
	if err := h.db.Preload("Seasons").First(&league, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, league)
}

func (h *LeagueHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var league models.League
	if err := h.db.First(&league, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := c.ShouldBindJSON(&league); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.db.Save(&league)
	c.JSON(http.StatusOK, league)
}

func (h *LeagueHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.db.Delete(&models.League{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

type CreateSeasonRequest struct {
	Name            string  `json:"name" binding:"required"`
	StartDate       string  `json:"start_date" binding:"required"`
	EndDate         string  `json:"end_date" binding:"required"`
	Format          string  `json:"format"`
	GroupCount      int     `json:"group_count"`
	PointsForWin    int     `json:"points_for_win"`
	PointsForDraw   int     `json:"points_for_draw"`
	PointsForLoss   int     `json:"points_for_loss"`
	MaxTeams        int     `json:"max_teams"`
	RegistrationFee float64 `json:"registration_fee"`
	VenueFee        float64 `json:"venue_fee"`
	CustomRules     string  `json:"custom_rules"`
	Awards          string  `json:"awards"`
}

func (h *LeagueHandler) CreateSeason(c *gin.Context) {
	leagueID, _ := strconv.Atoi(c.Param("id"))
	var req CreateSeasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	startDate := parseDate(req.StartDate)
	endDate := parseDate(req.EndDate)
	season := models.Season{
		LeagueID:        uint(leagueID),
		Name:            req.Name,
		StartDate:       startDate,
		EndDate:         endDate,
		Format:          req.Format,
		GroupCount:      req.GroupCount,
		PointsForWin:    req.PointsForWin,
		PointsForDraw:   req.PointsForDraw,
		PointsForLoss:   req.PointsForLoss,
		MaxTeams:        req.MaxTeams,
		RegistrationFee: req.RegistrationFee,
		VenueFee:        req.VenueFee,
		CustomRules:     req.CustomRules,
		Awards:          req.Awards,
		Status:          "planning",
	}
	if season.Format == "" {
		season.Format = "round_robin"
	}
	if season.GroupCount == 0 {
		season.GroupCount = 1
	}
	if season.PointsForWin == 0 {
		season.PointsForWin = 3
	}
	if season.PointsForDraw == 0 {
		season.PointsForDraw = 1
	}
	if season.MaxTeams == 0 {
		season.MaxTeams = 16
	}
	h.db.Create(&season)
	c.JSON(http.StatusCreated, season)
}

func (h *LeagueHandler) UpdateSeason(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("season_id"))
	var season models.Season
	if err := h.db.First(&season, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := c.ShouldBindJSON(&season); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.db.Save(&season)
	c.JSON(http.StatusOK, season)
}

func (h *LeagueHandler) GetSeason(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("season_id"))
	var season models.Season
	if err := h.db.Preload("Matches.HomeTeam").Preload("Matches.AwayTeam").Preload("Matches.Venue").First(&season, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, season)
}
