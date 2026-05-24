package handlers

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sports-league/models"
)

type MatchHandler struct {
	db *gorm.DB
}

func NewMatchHandler(db *gorm.DB) *MatchHandler {
	return &MatchHandler{db: db}
}

type GenerateScheduleRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	Interval  string `json:"interval"`
	VenueIDs  []uint `json:"venue_ids"`
}

type ReportScoreRequest struct {
	HomeScore    int  `json:"home_score" binding:"required"`
	AwayScore    int  `json:"away_score" binding:"required"`
	HasOT        bool `json:"has_ot"`
	OTHomeScore  int  `json:"ot_home_score"`
	OTAwayScore  int  `json:"ot_away_score"`
	HasPenalty   bool `json:"has_penalty"`
	PenHomeScore int  `json:"pen_home_score"`
	PenAwayScore int  `json:"pen_away_score"`
}

type GenerateKnockoutRequest struct {
	Format string `json:"format" binding:"required"`
	TeamIDs []uint `json:"team_ids" binding:"required"`
}

func (h *MatchHandler) GenerateSchedule(c *gin.Context) {
	seasonID, _ := strconv.Atoi(c.Param("season_id"))
	var req GenerateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var season models.Season
	if err := h.db.First(&season, seasonID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "season not found"})
		return
	}

	var registrations []models.Registration
	h.db.Where("season_id = ? AND status = ?", seasonID, "approved").Preload("Team").Find(&registrations)

	if len(registrations) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need at least 2 approved teams"})
		return
	}

	h.db.Where("season_id = ?", seasonID).Delete(&models.Match{})

	interval := 7
	if req.Interval == "3days" {
		interval = 3
	} else if req.Interval == "daily" {
		interval = 1
	}

	startDate := parseDate(req.StartDate)

	venues := req.VenueIDs
	if len(venues) == 0 {
		var allVenues []models.Venue
		h.db.Where("is_active = ?", true).Find(&allVenues)
		for _, v := range allVenues {
			venues = append(venues, v.ID)
		}
	}

	groups := make(map[string][]models.Registration)
	for _, reg := range registrations {
		gn := reg.GroupName
		if gn == "" {
			gn = "A"
		}
		groups[gn] = append(groups[gn], reg)
	}

	matchDate := startDate
	var matches []models.Match
	venueIdx := 0

	for _, teamRegs := range groups {
		teams := teamRegs
		n := len(teams)
		isOdd := n%2 != 0
		if isOdd {
			teams = append(teams, models.Registration{})
			n++
		}
		totalRounds := n - 1
		half := n / 2
		for round := 1; round <= totalRounds; round++ {
			for i := 0; i < half; i++ {
				home := teams[i]
				away := teams[n-1-i]
				if home.TeamID == 0 || away.TeamID == 0 {
					continue
				}
				var venueID *uint
				if venueIdx < len(venues) {
					venueID = &venues[venueIdx]
					venueIdx++
				}
				mt := matchDate
				m := models.Match{
					SeasonID:   uint(seasonID),
					Round:      round,
					GroupName:  home.GroupName,
					HomeTeamID: home.TeamID,
					AwayTeamID: away.TeamID,
					VenueID:    venueID,
					MatchTime:  &mt,
					Status:     "scheduled",
				}
				matches = append(matches, m)
			}
			teams = rotateTeams(teams)
			if round%2 == 0 {
				matchDate = matchDate.AddDate(0, 0, interval)
			}
		}
	}

	for _, m := range matches {
		h.db.Create(&m)
	}

	c.JSON(http.StatusCreated, gin.H{"matches": matches, "count": len(matches)})
}

func rotateTeams(teams []models.Registration) []models.Registration {
	if len(teams) <= 2 {
		return teams
	}
	last := teams[len(teams)-1]
	for i := len(teams) - 1; i > 1; i-- {
		teams[i] = teams[i-1]
	}
	teams[1] = last
	return teams
}

func (h *MatchHandler) ListMatches(c *gin.Context) {
	seasonID := c.Query("season_id")
	teamID := c.Query("team_id")
	var matches []models.Match
	q := h.db.Preload("HomeTeam").Preload("AwayTeam").Preload("Venue").Preload("Referee")
	if seasonID != "" {
		q = q.Where("season_id = ?", seasonID)
	}
	if teamID != "" {
		q = q.Where("home_team_id = ? OR away_team_id = ?", teamID, teamID)
	}
	q.Order("match_time ASC").Find(&matches)
	c.JSON(http.StatusOK, matches)
}

func (h *MatchHandler) GetMatch(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var m models.Match
	if err := h.db.Preload("HomeTeam").Preload("AwayTeam").Preload("Venue").Preload("Referee").First(&m, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *MatchHandler) UpdateMatch(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var m models.Match
	if err := h.db.First(&m, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.db.Save(&m)
	c.JSON(http.StatusOK, m)
}

func (h *MatchHandler) DeleteMatch(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.db.Delete(&models.Match{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *MatchHandler) ReportScore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req ReportScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var m models.Match
	if err := h.db.First(&m, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if m.Status == "finished" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "match already finished"})
		return
	}
	m.HomeScore = &req.HomeScore
	m.AwayScore = &req.AwayScore
	m.HasOT = req.HasOT
	m.HasPenalty = req.HasPenalty
	if req.HasOT {
		m.OTHomeScore = &req.OTHomeScore
		m.OTAwayScore = &req.OTAwayScore
	}
	if req.HasPenalty {
		m.PenHomeScore = &req.PenHomeScore
		m.PenAwayScore = &req.PenAwayScore
	}
	var winnerID uint
	if req.HomeScore > req.AwayScore {
		winnerID = m.HomeTeamID
	} else if req.AwayScore > req.HomeScore {
		winnerID = m.AwayTeamID
	} else {
		if req.HasPenalty {
			if req.PenHomeScore > req.PenAwayScore {
				winnerID = m.HomeTeamID
			} else {
				winnerID = m.AwayTeamID
			}
		}
	}
	if winnerID != 0 {
		m.WinnerTeamID = &winnerID
	}
	m.Status = "finished"
	h.db.Save(&m)

	h.sendMatchResultNotification(m)
	c.JSON(http.StatusOK, m)
}

func (h *MatchHandler) sendMatchResultNotification(m models.Match) {
	homeTeamName := ""
	awayTeamName := ""
	var home models.Team
	var away models.Team
	h.db.First(&home, m.HomeTeamID)
	h.db.First(&away, m.AwayTeamID)
	homeTeamName = home.Name
	awayTeamName = away.Name
	content := fmt.Sprintf("%s %d - %d %s", homeTeamName, *m.HomeScore, *m.AwayScore, awayTeamName)
	h.db.Create(&models.Notification{
		UserID:  home.CaptainID,
		Title:   "比赛结果",
		Content: content,
		Type:    "match_result",
	})
	h.db.Create(&models.Notification{
		UserID:  away.CaptainID,
		Title:   "比赛结果",
		Content: content,
		Type:    "match_result",
	})
}

func (h *MatchHandler) GetStandings(c *gin.Context) {
	seasonID, _ := strconv.Atoi(c.Param("season_id"))
	var season models.Season
	if err := h.db.First(&season, seasonID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "season not found"})
		return
	}

	var matches []models.Match
	h.db.Where("season_id = ? AND status = ?", seasonID, "finished").Find(&matches)

	type teamStat struct {
		TeamID       uint
		TeamName     string
		GroupName    string
		Played       int
		Wins         int
		Draws        int
		Losses       int
		GoalsFor     int
		GoalsAgainst int
		Points       int
	}
	stats := make(map[uint]*teamStat)

	var regs []models.Registration
	h.db.Where("season_id = ?", seasonID).Find(&regs)
	regMap := make(map[uint]models.Registration)
	for _, r := range regs {
		regMap[r.TeamID] = r
	}

	teamNames := make(map[uint]string)
	var teams []models.Team
	h.db.Find(&teams)
	for _, t := range teams {
		teamNames[t.ID] = t.Name
	}

	for _, m := range matches {
		if m.HomeScore == nil || m.AwayScore == nil {
			continue
		}
		if _, ok := stats[m.HomeTeamID]; !ok {
			stats[m.HomeTeamID] = &teamStat{TeamID: m.HomeTeamID, TeamName: teamNames[m.HomeTeamID], GroupName: regMap[m.HomeTeamID].GroupName}
		}
		if _, ok := stats[m.AwayTeamID]; !ok {
			stats[m.AwayTeamID] = &teamStat{TeamID: m.AwayTeamID, TeamName: teamNames[m.AwayTeamID], GroupName: regMap[m.AwayTeamID].GroupName}
		}
		hs, as := *m.HomeScore, *m.AwayScore
		stats[m.HomeTeamID].Played++
		stats[m.AwayTeamID].Played++
		stats[m.HomeTeamID].GoalsFor += hs
		stats[m.HomeTeamID].GoalsAgainst += as
		stats[m.AwayTeamID].GoalsFor += as
		stats[m.AwayTeamID].GoalsAgainst += hs
		if hs > as {
			stats[m.HomeTeamID].Wins++
			stats[m.AwayTeamID].Losses++
			stats[m.HomeTeamID].Points += season.PointsForWin
			stats[m.AwayTeamID].Points += season.PointsForLoss
		} else if hs < as {
			stats[m.AwayTeamID].Wins++
			stats[m.HomeTeamID].Losses++
			stats[m.AwayTeamID].Points += season.PointsForWin
			stats[m.HomeTeamID].Points += season.PointsForLoss
		} else {
			stats[m.HomeTeamID].Draws++
			stats[m.AwayTeamID].Draws++
			stats[m.HomeTeamID].Points += season.PointsForDraw
			stats[m.AwayTeamID].Points += season.PointsForDraw
		}
	}

	var result []models.Standings
	for _, s := range stats {
		result = append(result, models.Standings{
			SeasonID:     uint(seasonID),
			TeamID:       s.TeamID,
			TeamName:     s.TeamName,
			GroupName:    s.GroupName,
			Played:       s.Played,
			Wins:         s.Wins,
			Draws:        s.Draws,
			Losses:       s.Losses,
			GoalsFor:     s.GoalsFor,
			GoalsAgainst: s.GoalsAgainst,
			GoalDiff:     s.GoalsFor - s.GoalsAgainst,
			Points:       s.Points,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Points != result[j].Points {
			return result[i].Points > result[j].Points
		}
		return result[i].GoalDiff > result[j].GoalDiff
	})

	c.JSON(http.StatusOK, result)
}

func (h *MatchHandler) GenerateKnockout(c *gin.Context) {
	seasonID, _ := strconv.Atoi(c.Param("season_id"))
	var req GenerateKnockoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	n := len(req.TeamIDs)
	if n < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "need at least 2 teams"})
		return
	}

	pow2 := int(math.Pow(2, math.Ceil(math.Log2(float64(n)))))
	teams := make([]uint, pow2)
	copy(teams, req.TeamIDs)

	h.db.Where("season_id = ? AND knockout_stage <> ?", seasonID, "").Delete(&models.Match{})

	var matches []models.Match
	round := 1
	stage := fmt.Sprintf("Round of %d", pow2)
	if pow2 == 4 {
		stage = "Semi-final"
	} else if pow2 == 2 {
		stage = "Final"
	}

	for i := 0; i < len(teams); i += 2 {
		var venueID *uint
		var venue models.Venue
		h.db.Where("is_active = ?", true).First(&venue)
		if venue.ID != 0 {
			venueID = &venue.ID
		}
		homeID := teams[i]
		awayID := teams[i+1]
		if homeID == 0 || awayID == 0 {
			continue
		}
		m := models.Match{
			SeasonID:      uint(seasonID),
			Round:         round,
			HomeTeamID:    homeID,
			AwayTeamID:    awayID,
			VenueID:       venueID,
			Status:        "scheduled",
			KnockoutStage: stage,
		}
		matches = append(matches, m)
		h.db.Create(&m)
	}

	c.JSON(http.StatusCreated, gin.H{"matches": matches, "format": req.Format})
}

func (h *MatchHandler) CreateVenue(c *gin.Context) {
	var v models.Venue
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v.IsActive = true
	h.db.Create(&v)
	c.JSON(http.StatusCreated, v)
}

func (h *MatchHandler) ListVenues(c *gin.Context) {
	var venues []models.Venue
	h.db.Find(&venues)
	c.JSON(http.StatusOK, venues)
}

func (h *MatchHandler) UpdateVenue(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var v models.Venue
	if err := h.db.First(&v, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.db.Save(&v)
	c.JSON(http.StatusOK, v)
}

func (h *MatchHandler) CheckVenueConflict(c *gin.Context) {
	venueID, _ := strconv.Atoi(c.Query("venue_id"))
	matchTime := c.Query("match_time")
	if venueID == 0 || matchTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "venue_id and match_time required"})
		return
	}
	t, err := time.Parse(time.RFC3339, matchTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match_time format"})
		return
	}
	var count int64
	h.db.Model(&models.Match{}).Where(
		"venue_id = ? AND match_time = ? AND status <> ?",
		venueID, t, "cancelled",
	).Count(&count)
	c.JSON(http.StatusOK, gin.H{"conflict": count > 0, "count": count})
}

func parseDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		t, _ = time.Parse(time.RFC3339, s)
	}
	return t
}
