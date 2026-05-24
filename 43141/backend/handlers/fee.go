package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sports-league/models"
)

type FeeHandler struct {
	db *gorm.DB
}

func NewFeeHandler(db *gorm.DB) *FeeHandler {
	return &FeeHandler{db: db}
}

type CreateFeeRequest struct {
	SeasonID uint    `json:"season_id" binding:"required"`
	TeamID   uint    `json:"team_id" binding:"required"`
	Type     string  `json:"type" binding:"required"`
	Amount   float64 `json:"amount" binding:"required"`
	Note     string  `json:"note"`
}

func (h *FeeHandler) Create(c *gin.Context) {
	var req CreateFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	invoiceNo := fmt.Sprintf("INV-%d-%d", req.SeasonID, time.Now().Unix())
	fee := models.Fee{
		SeasonID:  req.SeasonID,
		TeamID:    req.TeamID,
		Type:      req.Type,
		Amount:    req.Amount,
		Status:    "unpaid",
		InvoiceNo: invoiceNo,
		Note:      req.Note,
	}
	h.db.Create(&fee)
	c.JSON(http.StatusCreated, fee)
}

func (h *FeeHandler) List(c *gin.Context) {
	seasonID := c.Query("season_id")
	teamID := c.Query("team_id")
	var fees []models.Fee
	q := h.db.Preload("Team").Preload("Season")
	if seasonID != "" {
		q = q.Where("season_id = ?", seasonID)
	}
	if teamID != "" {
		q = q.Where("team_id = ?", teamID)
	}
	q.Find(&fees)
	c.JSON(http.StatusOK, fees)
}

func (h *FeeHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var fee models.Fee
	if err := h.db.Preload("Team").Preload("Season").First(&fee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, fee)
}

func (h *FeeHandler) MarkPaid(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var fee models.Fee
	if err := h.db.First(&fee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	now := time.Now()
	fee.Status = "paid"
	fee.PaidAt = &now
	h.db.Save(&fee)

	h.db.Model(&models.Registration{}).Where("team_id = ? AND season_id = ?", fee.TeamID, fee.SeasonID).Update("paid", true)

	c.JSON(http.StatusOK, fee)
}

func (h *FeeHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.db.Delete(&models.Fee{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

type FeeReport struct {
	SeasonID    uint    `json:"season_id"`
	TotalAmount float64 `json:"total_amount"`
	PaidAmount  float64 `json:"paid_amount"`
	UnpaidAmount float64 `json:"unpaid_amount"`
	InvoiceCount int64  `json:"invoice_count"`
}

func (h *FeeHandler) Report(c *gin.Context) {
	seasonID, _ := strconv.Atoi(c.Param("season_id"))
	var fees []models.Fee
	h.db.Where("season_id = ?", seasonID).Find(&fees)

	var report FeeReport
	report.SeasonID = uint(seasonID)
	report.InvoiceCount = int64(len(fees))
	for _, f := range fees {
		report.TotalAmount += f.Amount
		if f.Status == "paid" {
			report.PaidAmount += f.Amount
		} else {
			report.UnpaidAmount += f.Amount
		}
	}
	c.JSON(http.StatusOK, report)
}
