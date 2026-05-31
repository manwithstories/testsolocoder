package handlers

import (
	"net/http"
	"strconv"
	"time"

	"print3d-platform/internal/middleware"
	"print3d-platform/internal/models"
	"print3d-platform/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PrinterHandler struct {
	printerService *service.PrinterService
}

func NewPrinterHandler(printerService *service.PrinterService) *PrinterHandler {
	return &PrinterHandler{printerService: printerService}
}

func (h *PrinterHandler) CreateDevice(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	if authUser.Role != models.RolePrinter {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only printers can manage devices"})
		return
	}

	var req service.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := h.printerService.CreateDevice(c.Request.Context(), authUser.UserID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, device)
}

func (h *PrinterHandler) GetDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	device, err := h.printerService.GetDevice(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	c.JSON(http.StatusOK, device)
}

func (h *PrinterHandler) GetPrinterDevices(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	devices, err := h.printerService.GetPrinterDevices(c.Request.Context(), authUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get devices"})
		return
	}

	c.JSON(http.StatusOK, devices)
}

func (h *PrinterHandler) UpdateDevice(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	var req service.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := h.printerService.UpdateDevice(c.Request.Context(), id, authUser.UserID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device)
}

func (h *PrinterHandler) DeleteDevice(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	err = h.printerService.DeleteDevice(c.Request.Context(), id, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}

func (h *PrinterHandler) CreateInventory(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	if authUser.Role != models.RolePrinter {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only printers can manage inventory"})
		return
	}

	var req service.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inventory, err := h.printerService.CreateInventory(c.Request.Context(), authUser.UserID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, inventory)
}

func (h *PrinterHandler) GetPrinterInventory(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	inventory, err := h.printerService.GetPrinterInventory(c.Request.Context(), authUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inventory"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (h *PrinterHandler) UpdateInventoryQuantity(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory ID"})
		return
	}

	var req struct {
		QuantityGrams float64 `json:"quantity_grams" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.printerService.UpdateInventoryQuantity(c.Request.Context(), id, authUser.UserID, req.QuantityGrams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory updated successfully"})
}

func (h *PrinterHandler) DeleteInventory(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory ID"})
		return
	}

	err = h.printerService.DeleteInventory(c.Request.Context(), id, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory item deleted successfully"})
}

func (h *PrinterHandler) CreateSchedule(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	if authUser.Role != models.RolePrinter {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only printers can create schedules"})
		return
	}

	var req struct {
		DeviceID       uuid.UUID `json:"device_id" binding:"required"`
		OrderID        uuid.UUID `json:"order_id" binding:"required"`
		ScheduledStart string    `json:"scheduled_start"`
		ScheduledEnd   string    `json:"scheduled_end"`
		Priority       int       `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var scheduledStart, scheduledEnd *time.Time
	if req.ScheduledStart != "" {
		t, err := time.Parse(time.RFC3339, req.ScheduledStart)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled start time"})
			return
		}
		scheduledStart = &t
	}
	if req.ScheduledEnd != "" {
		t, err := time.Parse(time.RFC3339, req.ScheduledEnd)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled end time"})
			return
		}
		scheduledEnd = &t
	}

	schedule, err := h.printerService.CreateSchedule(
		c.Request.Context(),
		authUser.UserID,
		req.DeviceID,
		req.OrderID,
		scheduledStart,
		scheduledEnd,
		req.Priority,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, schedule)
}

func (h *PrinterHandler) GetPrinterSchedules(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	schedules, err := h.printerService.GetPrinterSchedules(c.Request.Context(), authUser.UserID, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get schedules"})
		return
	}

	c.JSON(http.StatusOK, schedules)
}

func (h *PrinterHandler) CreateReview(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	var req service.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.printerService.CreateReview(c.Request.Context(), authUser.UserID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func (h *PrinterHandler) GetReview(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := h.printerService.GetReview(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, review)
}

func (h *PrinterHandler) GetModelReviews(c *gin.Context) {
	modelID, err := uuid.Parse(c.Param("model_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	reviews, total, err := h.printerService.GetModelReviews(c.Request.Context(), modelID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  reviews,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *PrinterHandler) GetPrinterReviews(c *gin.Context) {
	printerID, err := uuid.Parse(c.Param("printer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid printer ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	reviews, total, err := h.printerService.GetPrinterReviews(c.Request.Context(), printerID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  reviews,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *PrinterHandler) GetMaterials(c *gin.Context) {
	materials, err := h.printerService.GetMaterials(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get materials"})
		return
	}

	c.JSON(http.StatusOK, materials)
}

func (h *PrinterHandler) GetMaterial(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
		return
	}

	material, err := h.printerService.GetMaterial(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
		return
	}

	c.JSON(http.StatusOK, material)
}

func (h *PrinterHandler) GetIdleDevices(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	devices, err := h.printerService.GetIdleDevices(c.Request.Context(), authUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get idle devices"})
		return
	}

	c.JSON(http.StatusOK, devices)
}
