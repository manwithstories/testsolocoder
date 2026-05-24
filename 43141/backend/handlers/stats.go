package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sports-league/models"
)

type StatsHandler struct {
	db *gorm.DB
}

func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

type AddPlayerStatRequest struct {
	PlayerID   uint `json:"player_id" binding:"required"`
	MatchID    uint `json:"match_id" binding:"required"`
	SeasonID   uint `json:"season_id" binding:"required"`
	TeamID     uint `json:"team_id" binding:"required"`
	Goals      int  `json:"goals"`
	Assists    int  `json:"assists"`
	Fouls      int  `json:"fouls"`
	YellowCard int  `json:"yellow_card"`
	RedCard    int  `json:"red_card"`
	Minutes    int  `json:"minutes"`
}

func (h *StatsHandler) AddStat(c *gin.Context) {
	var req AddPlayerStatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stat := models.PlayerStat{
		PlayerID:   req.PlayerID,
		MatchID:    req.MatchID,
		SeasonID:   req.SeasonID,
		TeamID:     req.TeamID,
		Goals:      req.Goals,
		Assists:    req.Assists,
		Fouls:      req.Fouls,
		YellowCard: req.YellowCard,
		RedCard:    req.RedCard,
		Minutes:    req.Minutes,
	}
	h.db.Create(&stat)
	c.JSON(http.StatusCreated, stat)
}

type PlayerRanking struct {
	PlayerID   uint   `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	Matches    int    `json:"matches"`
	Goals      int    `json:"goals"`
	Assists    int    `json:"assists"`
	Fouls      int    `json:"fouls"`
	YellowCard int    `json:"yellow_card"`
	RedCard    int    `json:"red_card"`
	Minutes    int    `json:"minutes"`
}

func (h *StatsHandler) GetRankings(c *gin.Context) {
	seasonID := c.Query("season_id")
	teamID := c.Query("team_id")

	var stats []models.PlayerStat
	q := h.db.Preload("Player").Preload("Match")
	if seasonID != "" {
		q = q.Where("season_id = ?", seasonID)
	}
	if teamID != "" {
		q = q.Where("team_id = ?", teamID)
	}
	q.Find(&stats)

	rankings := make(map[uint]*PlayerRanking)
	for _, s := range stats {
		if _, ok := rankings[s.PlayerID]; !ok {
			teamName := ""
			var team models.Team
			h.db.First(&team, s.TeamID)
			teamName = team.Name
			playerName := ""
			if s.Player != nil {
				playerName = s.Player.Name
			}
			rankings[s.PlayerID] = &PlayerRanking{
				PlayerID:   s.PlayerID,
				PlayerName: playerName,
				TeamName:   teamName,
			}
		}
		rankings[s.PlayerID].Matches++
		rankings[s.PlayerID].Goals += s.Goals
		rankings[s.PlayerID].Assists += s.Assists
		rankings[s.PlayerID].Fouls += s.Fouls
		rankings[s.PlayerID].YellowCard += s.YellowCard
		rankings[s.PlayerID].RedCard += s.RedCard
		rankings[s.PlayerID].Minutes += s.Minutes
	}

	var result []PlayerRanking
	for _, r := range rankings {
		result = append(result, *r)
	}
	c.JSON(http.StatusOK, result)
}

func (h *StatsHandler) GetPlayerStats(c *gin.Context) {
	playerID, _ := strconv.Atoi(c.Param("id"))
	var stats []models.PlayerStat
	h.db.Preload("Match").Where("player_id = ?", playerID).Find(&stats)
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) DeleteStat(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.db.Delete(&models.PlayerStat{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
