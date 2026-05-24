package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sports-league/models"
	"sports-league/pkg/auth"
)

type RefereeHandler struct {
	db *gorm.DB
}

func NewRefereeHandler(db *gorm.DB) *RefereeHandler {
	return &RefereeHandler{db: db}
}

type AssignRefereeRequest struct {
	MatchID   uint `json:"match_id" binding:"required"`
	RefereeID uint `json:"referee_id" binding:"required"`
}

func (h *RefereeHandler) Assign(c *gin.Context) {
	var req AssignRefereeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var match models.Match
	if err := h.db.First(&match, req.MatchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}
	match.RefereeID = &req.RefereeID
	h.db.Save(&match)

	assignment := models.RefereeAssignment{
		MatchID:   req.MatchID,
		RefereeID: req.RefereeID,
		Status:    "assigned",
	}
	h.db.Create(&assignment)

	var referee models.User
	h.db.First(&referee, req.RefereeID)
	h.db.Create(&models.Notification{
		UserID:  req.RefereeID,
		Title:   "新的执法任务",
		Content: "您已被指派执法一场比赛，请确认。",
		Type:    "referee_assignment",
	})

	c.JSON(http.StatusCreated, assignment)
}

func (h *RefereeHandler) Respond(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	claims, _ := c.Get("claims")
	userClaims := claims.(*auth.Claims)

	var assignment models.RefereeAssignment
	if err := h.db.First(&assignment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if assignment.RefereeID != userClaims.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your assignment"})
		return
	}

	action := c.Query("action")
	if action == "accept" {
		assignment.Status = "accepted"
	} else if action == "reject" {
		assignment.Status = "rejected"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "action must be accept or reject"})
		return
	}
	now := time.Now()
	assignment.RespondedAt = &now
	h.db.Save(&assignment)
	c.JSON(http.StatusOK, assignment)
}

func (h *RefereeHandler) ListAssignments(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*auth.Claims)
	var assignments []models.RefereeAssignment
	q := h.db.Preload("Match.HomeTeam").Preload("Match.AwayTeam").Preload("Match.Venue")
	if userClaims.Role == models.RoleReferee {
		q = q.Where("referee_id = ?", userClaims.UserID)
	}
	q.Find(&assignments)
	c.JSON(http.StatusOK, assignments)
}

func (h *RefereeHandler) ListReferees(c *gin.Context) {
	var referees []models.User
	h.db.Where("role = ? AND is_active = ?", models.RoleReferee, true).Find(&referees)
	c.JSON(http.StatusOK, referees)
}
