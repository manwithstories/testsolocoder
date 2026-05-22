package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"

	"travel-planner/internal/database"
	"travel-planner/internal/logger"
	"travel-planner/internal/models"
	"travel-planner/internal/utils"
)

func ExportJSON(c *gin.Context) {
	planID := c.Param("id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	var plan models.TravelPlan
	if err := database.DB.Preload("Owner").
		Preload("Participants.User").
		Preload("Activities").
		Preload("Files").
		Preload("Checklists.Items").
		First(&plan, "id = ?", planUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Plan not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	exportData := struct {
		Plan         models.TravelPlan `json:"plan"`
		ExportedAt   time.Time        `json:"exported_at"`
		Version       string           `json:"version"`
	}{
		Plan:       plan,
		ExportedAt: time.Now(),
		Version:    "1.0",
	}

	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		logger.Errorf("Failed to marshal JSON: %v", err)
		utils.InternalServerError(c, "Failed to export")
		return
	}

	filename := fmt.Sprintf("%s_travel_plan.json", plan.Title)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/json")
	c.String(200, string(jsonData))

	logger.Infof("Plan exported as JSON: %s", planID)
}

func ExportPDF(c *gin.Context) {
	planID := c.Param("id")
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		utils.BadRequest(c, "Invalid plan ID")
		return
	}

	var plan models.TravelPlan
	if err := database.DB.Preload("Owner").
		Preload("Participants.User").
		Preload("Activities").
		Preload("Checklists.Items").
		First(&plan, "id = ?", planUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Plan not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle(plan.Title, true)
	pdf.SetAuthor("Travel Planner", true)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, plan.Title)
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 8, fmt.Sprintf("Destination: %s", plan.Destination))
	pdf.Ln(8)
	pdf.Cell(40, 8, fmt.Sprintf("Dates: %s to %s", 
		plan.StartDate.Format("2006-01-02"), plan.EndDate.Format("2006-01-02"))
	pdf.Ln(8)
	pdf.Cell(40, 8, fmt.Sprintf("Budget: %.2f %s", plan.Budget, plan.Currency))
	pdf.Ln(8)
	pdf.Cell(40, 8, fmt.Sprintf("Owner: %s %s", plan.Owner.FirstName, plan.Owner.LastName))
	pdf.Ln(12)

	if plan.Description != "" {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 8, "Description:")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 6, plan.Description, "", "", false)
		pdf.Ln(8)
	}

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Activities")
	pdf.Ln(10)

	currentDate := ""
	for _, activity := range plan.Activities {
		dateStr := activity.Date.Format("2006-01-02")
		if dateStr != currentDate {
			pdf.SetFont("Arial", "B", 12)
			pdf.Cell(40, 8, dateStr)
			pdf.Ln(8)
			currentDate = dateStr
		}

		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(40, 6, fmt.Sprintf("- %s", activity.Title))
		pdf.Ln(6)
		pdf.SetFont("Arial", "", 8)
		if activity.StartTime != "" {
			pdf.Cell(40, 5, fmt.Sprintf("  Time: %s - %s", activity.StartTime, activity.EndTime))
			pdf.Ln(5)
		}
		if activity.Location != "" {
			pdf.Cell(40, 5, fmt.Sprintf("  Location: %s", activity.Location))
			pdf.Ln(5)
		}
		if activity.Cost > 0 {
			pdf.Cell(40, 5, fmt.Sprintf("  Cost: %.2f %s", activity.Cost, activity.Currency))
			pdf.Ln(5)
		}
		if activity.Notes != "" {
			pdf.Cell(40, 5, fmt.Sprintf("  Notes: %s", activity.Notes))
			pdf.Ln(5)
		}
		pdf.Ln(2)
	}

	if len(plan.Checklists) > 0 {
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "Checklists")
		pdf.Ln(10)

		for _, checklist := range plan.Checklists {
			pdf.SetFont("Arial", "B", 12)
			pdf.Cell(40, 8, checklist.Title)
			pdf.Ln(8)

			for _, item := range checklist.Items {
				pdf.SetFont("Arial", "", 10)
				status := "[ ]"
				if item.IsCompleted {
					status = "[X]"
				}
				pdf.Cell(40, 6, fmt.Sprintf("%s %s", status, item.Title))
				pdf.Ln(6)
			}
			pdf.Ln(4)
		}
	}

	filename := fmt.Sprintf("%s_travel_plan.pdf", plan.Title)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/pdf")

	err = pdf.Output(c.Writer)
	if err != nil {
		logger.Errorf("Failed to generate PDF: %v", err)
		utils.InternalServerError(c, "Failed to generate PDF")
		return
	}

	logger.Infof("Plan exported as PDF: %s", planID)
}
