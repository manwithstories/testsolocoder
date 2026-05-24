package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"sports-league/models"
)

type ExportHandler struct {
	db *gorm.DB
}

func NewExportHandler(db *gorm.DB) *ExportHandler {
	return &ExportHandler{db: db}
}

func (h *ExportHandler) ExportSchedule(c *gin.Context) {
	seasonID, _ := strconv.Atoi(c.Param("season_id"))
	var season models.Season
	h.db.First(&season, seasonID)

	var matches []models.Match
	h.db.Preload("HomeTeam").Preload("AwayTeam").Preload("Venue").
		Where("season_id = ?", seasonID).Order("match_time ASC").Find(&matches)

	f := excelize.NewFile()
	sheet := "Schedule"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Round", "Group", "Home Team", "Away Team", "Venue", "Match Time", "Status"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, m := range matches {
		row := i + 2
		homeName := ""
		awayName := ""
		venueName := ""
		matchTime := ""
		if m.HomeTeam != nil {
			homeName = m.HomeTeam.Name
		}
		if m.AwayTeam != nil {
			awayName = m.AwayTeam.Name
		}
		if m.Venue != nil {
			venueName = m.Venue.Name
		}
		if m.MatchTime != nil {
			matchTime = m.MatchTime.Format("2006-01-02 15:04")
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), m.Round)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), m.GroupName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), homeName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), awayName)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), venueName)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), matchTime)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), m.Status)
	}

	filename := fmt.Sprintf("schedule_season_%d_%s.xlsx", seasonID, time.Now().Format("20060102"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	f.Write(c.Writer)
}

func (h *ExportHandler) ExportStandings(c *gin.Context) {
	seasonID, _ := strconv.Atoi(c.Param("season_id"))

	standings := h.getStandingsData(seasonID)

	f := excelize.NewFile()
	sheet := "Standings"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Team", "Group", "Pld", "W", "D", "L", "GF", "GA", "GD", "Pts"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, s := range standings {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), s.TeamName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), s.GroupName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), s.Played)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), s.Wins)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), s.Draws)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), s.Losses)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), s.GoalsFor)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), s.GoalsAgainst)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), s.GoalDiff)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), s.Points)
	}

	filename := fmt.Sprintf("standings_season_%d_%s.xlsx", seasonID, time.Now().Format("20060102"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	f.Write(c.Writer)
}

func (h *ExportHandler) ExportStats(c *gin.Context) {
	seasonID := c.Query("season_id")
	var stats []models.PlayerStat
	q := h.db.Preload("Player")
	if seasonID != "" {
		q = q.Where("season_id = ?", seasonID)
	}
	q.Find(&stats)

	f := excelize.NewFile()
	sheet := "PlayerStats"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Player", "Goals", "Assists", "Fouls", "Yellow", "Red", "Minutes"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, s := range stats {
		row := i + 2
		playerName := ""
		if s.Player != nil {
			playerName = s.Player.Name
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), playerName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), s.Goals)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), s.Assists)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), s.Fouls)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), s.YellowCard)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), s.RedCard)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), s.Minutes)
	}

	filename := fmt.Sprintf("player_stats_%s.xlsx", time.Now().Format("20060102"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	f.Write(c.Writer)
}

func (h *ExportHandler) ExportPDF(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "PDF export",
		"note":    "PDF export is available via Excel conversion or dedicated PDF library integration",
	})
}

func (h *ExportHandler) getStandingsData(seasonID int) []models.Standings {
	var season models.Season
	h.db.First(&season, seasonID)

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
	statsMap := make(map[uint]*teamStat)

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
		if _, ok := statsMap[m.HomeTeamID]; !ok {
			statsMap[m.HomeTeamID] = &teamStat{TeamID: m.HomeTeamID, TeamName: teamNames[m.HomeTeamID], GroupName: regMap[m.HomeTeamID].GroupName}
		}
		if _, ok := statsMap[m.AwayTeamID]; !ok {
			statsMap[m.AwayTeamID] = &teamStat{TeamID: m.AwayTeamID, TeamName: teamNames[m.AwayTeamID], GroupName: regMap[m.AwayTeamID].GroupName}
		}
		hs, as := *m.HomeScore, *m.AwayScore
		statsMap[m.HomeTeamID].Played++
		statsMap[m.AwayTeamID].Played++
		statsMap[m.HomeTeamID].GoalsFor += hs
		statsMap[m.HomeTeamID].GoalsAgainst += as
		statsMap[m.AwayTeamID].GoalsFor += as
		statsMap[m.AwayTeamID].GoalsAgainst += hs
		if hs > as {
			statsMap[m.HomeTeamID].Wins++
			statsMap[m.AwayTeamID].Losses++
			statsMap[m.HomeTeamID].Points += season.PointsForWin
			statsMap[m.AwayTeamID].Points += season.PointsForLoss
		} else if hs < as {
			statsMap[m.AwayTeamID].Wins++
			statsMap[m.HomeTeamID].Losses++
			statsMap[m.AwayTeamID].Points += season.PointsForWin
			statsMap[m.HomeTeamID].Points += season.PointsForLoss
		} else {
			statsMap[m.HomeTeamID].Draws++
			statsMap[m.AwayTeamID].Draws++
			statsMap[m.HomeTeamID].Points += season.PointsForDraw
			statsMap[m.AwayTeamID].Points += season.PointsForDraw
		}
	}

	var result []models.Standings
	for _, s := range statsMap {
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

	return result
}
