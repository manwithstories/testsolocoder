package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"freelancer-management/internal/config"
	"freelancer-management/internal/database"
	"freelancer-management/internal/logger"
	"freelancer-management/internal/middleware"
	"freelancer-management/internal/models"
	"freelancer-management/internal/utils"
	pdfpkg "freelancer-management/pkg/pdf"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InvoiceHandler struct {
	db *gorm.DB
}

func NewInvoiceHandler() *InvoiceHandler {
	return &InvoiceHandler{db: database.GetDB()}
}

type CreateInvoiceRequest struct {
	ClientID    uint            `json:"client_id" binding:"required"`
	ProjectID   *uint           `json:"project_id"`
	TimeEntryIDs []uint         `json:"time_entry_ids"`
	IssueDate   string          `json:"issue_date" binding:"required"`
	DueDate     string          `json:"due_date" binding:"required"`
	TaxRate     float64         `json:"tax_rate"`
	Notes       string          `json:"notes"`
	CustomItems []CustomItemReq `json:"custom_items"`
}

type CustomItemReq struct {
	Description string  `json:"description" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
}

type UpdateInvoiceStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *InvoiceHandler) generateInvoiceNumber(tx *gorm.DB, year int) (string, error) {
	cfg := config.GetConfig()
	var counter models.InvoiceCounter

	result := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("year = ?", year).
		First(&counter)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		counter = models.InvoiceCounter{
			Year:  year,
			Count: 1,
		}
		if err := tx.Create(&counter).Error; err != nil {
			return "", err
		}
	} else if result.Error != nil {
		return "", result.Error
	} else {
		counter.Count++
		if err := tx.Save(&counter).Error; err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s-%d-%06d", cfg.App.InvoicePrefix, year, counter.Count), nil
}

func (h *InvoiceHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	var client models.Client
	if err := h.db.Where("id = ? AND user_id = ?", req.ClientID, userID).First(&client).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Client not found")
		return
	}

	if req.ProjectID != nil {
		var project models.Project
		if err := h.db.Where("id = ? AND user_id = ?", *req.ProjectID, userID).First(&project).Error; err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "Project not found")
			return
		}
	}

	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid issue date format, use YYYY-MM-DD")
		return
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid due date format, use YYYY-MM-DD")
		return
	}

	tx := h.db.Begin()
	if tx.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	invoiceNumber, err := h.generateInvoiceNumber(tx, issueDate.Year())
	if err != nil {
		tx.Rollback()
		logger.LogError("Failed to generate invoice number: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate invoice number")
		return
	}

	hourlyRate := client.DefaultRate
	if req.ProjectID != nil {
		var project models.Project
		tx.First(&project, *req.ProjectID)
		if project.HourlyRate > 0 {
			hourlyRate = project.HourlyRate
		}
	}

	invoice := models.Invoice{
		UserID:        userID,
		ClientID:      req.ClientID,
		ProjectID:     req.ProjectID,
		InvoiceNumber: invoiceNumber,
		Status:        models.InvoiceStatusDraft,
		IssueDate:     issueDate,
		DueDate:       dueDate,
		TaxRate:       req.TaxRate,
		Notes:         req.Notes,
	}

	var subtotal float64
	var invoiceItems []models.InvoiceItem

	for _, timeEntryID := range req.TimeEntryIDs {
		var timeEntry models.TimeEntry
		if err := tx.Where("id = ? AND user_id = ? AND billable = ?", timeEntryID, userID, true).First(&timeEntry).Error; err != nil {
			tx.Rollback()
			utils.ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Time entry %d not found or not billable", timeEntryID))
			return
		}

		rate := hourlyRate
		amount := timeEntry.Hours * rate
		subtotal += amount

		invoiceItems = append(invoiceItems, models.InvoiceItem{
			Description: fmt.Sprintf("Work on %s - %s", timeEntry.Date.Format("2006-01-02"), timeEntry.Description),
			Quantity:    timeEntry.Hours,
			UnitPrice:   rate,
			Amount:      amount,
		})
	}

	for _, item := range req.CustomItems {
		amount := item.Quantity * item.UnitPrice
		subtotal += amount
		invoiceItems = append(invoiceItems, models.InvoiceItem{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Amount:      amount,
		})
	}

	invoice.Subtotal = subtotal
	invoice.TaxAmount = subtotal * (req.TaxRate / 100)
	invoice.Total = subtotal + invoice.TaxAmount

	if err := tx.Create(&invoice).Error; err != nil {
		tx.Rollback()
		logger.LogError("Failed to create invoice: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create invoice")
		return
	}

	for i := range invoiceItems {
		invoiceItems[i].InvoiceID = invoice.ID
	}
	if len(invoiceItems) > 0 {
		if err := tx.Create(&invoiceItems).Error; err != nil {
			tx.Rollback()
			logger.LogError("Failed to create invoice items: %v", err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create invoice items")
			return
		}
	}

	tx.Commit()

	invoice.Items = invoiceItems
	invoice.Client = &client

	logger.LogOperation(userID, "create_invoice", "Invoice created: "+invoiceNumber)
	utils.SuccessResponse(c, invoice)
}

func (h *InvoiceHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	status := c.Query("status")
	clientID := c.Query("client_id")
	year := c.Query("year")
	offset := (page - 1) * perPage

	var invoices []models.Invoice
	var total int64

	query := h.db.Model(&models.Invoice{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if clientID != "" {
		query = query.Where("client_id = ?", clientID)
	}
	if year != "" {
		query = query.Where("strftime('%Y', issue_date) = ?", year)
	}

	query.Count(&total)
	query.Preload("Client").Preload("Items").Offset(offset).Limit(perPage).Order("issue_date DESC, created_at DESC").Find(&invoices)

	utils.PaginatedSuccessResponse(c, invoices, utils.PaginationMeta{
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: int((total + int64(perPage) - 1) / int64(perPage)),
	})
}

func (h *InvoiceHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var invoice models.Invoice
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Client").
		Preload("Project").
		Preload("Items").
		First(&invoice).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Invoice not found")
		return
	}

	utils.SuccessResponse(c, invoice)
}

func (h *InvoiceHandler) UpdateStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var invoice models.Invoice
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&invoice).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Invoice not found")
		return
	}

	var req UpdateInvoiceStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	validStatuses := map[string]bool{
		"draft":     true,
		"sent":      true,
		"paid":      true,
		"overdue":   true,
		"cancelled": true,
	}

	if !validStatuses[req.Status] {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid status")
		return
	}

	invoice.Status = models.InvoiceStatus(req.Status)

	if err := h.db.Save(&invoice).Error; err != nil {
		logger.LogError("Failed to update invoice status: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update invoice status")
		return
	}

	logger.LogOperation(userID, "update_invoice_status", fmt.Sprintf("Invoice %s status updated to %s", invoice.InvoiceNumber, req.Status))
	utils.SuccessResponse(c, invoice)
}

func (h *InvoiceHandler) DownloadPDF(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var invoice models.Invoice
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Client").
		Preload("Items").
		First(&invoice).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Invoice not found")
		return
	}

	filename, err := pdfpkg.GenerateInvoicePDF(&invoice)
	if err != nil {
		logger.LogError("Failed to generate PDF: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate PDF")
		return
	}

	defer os.Remove(filename)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+invoice.InvoiceNumber+".pdf")
	c.Header("Content-Type", "application/pdf")
	c.File(filename)

	logger.LogOperation(userID, "download_invoice_pdf", "Invoice PDF downloaded: "+invoice.InvoiceNumber)
}

func (h *InvoiceHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var invoice models.Invoice
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&invoice).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Invoice not found")
		return
	}

	tx := h.db.Begin()

	if err := tx.Where("invoice_id = ?", id).Delete(&models.InvoiceItem{}).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete invoice items")
		return
	}

	if err := tx.Delete(&invoice).Error; err != nil {
		tx.Rollback()
		logger.LogError("Failed to delete invoice: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete invoice")
		return
	}

	tx.Commit()

	logger.LogOperation(userID, "delete_invoice", "Invoice deleted: "+invoice.InvoiceNumber)
	utils.SuccessResponseWithMessage(c, "Invoice deleted successfully", nil)
}
