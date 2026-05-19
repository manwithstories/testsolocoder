package handler

import (
	"net/http"
	"strconv"
	"venue-booking/internal/dto"
	"venue-booking/internal/service"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	deviceService *service.DeviceService
	logService    *service.OperationLogService
}

func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{
		deviceService: service.NewDeviceService(),
		logService:    service.NewOperationLogService(),
	}
}

func (h *DeviceHandler) CreateCategory(c *gin.Context) {
	var req dto.DeviceCategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	category, err := h.deviceService.CreateCategory(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to create category"))
		return
	}

	userID, _ := c.Get("userID")
	h.logService.Log(c, userID.(uint), "create_device_category", "device", map[string]interface{}{
		"category_id":   category.ID,
		"category_name": category.Name,
	})

	c.JSON(http.StatusOK, dto.Success(category))
}

func (h *DeviceHandler) ListCategories(c *gin.Context) {
	categories, err := h.deviceService.ListCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to get categories"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(categories))
}

func (h *DeviceHandler) Create(c *gin.Context) {
	var req dto.DeviceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	device, err := h.deviceService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	userID, _ := c.Get("userID")
	h.logService.Log(c, userID.(uint), "create_device", "device", map[string]interface{}{
		"device_id":   device.ID,
		"device_name": device.Name,
	})

	c.JSON(http.StatusOK, dto.Success(device))
}

func (h *DeviceHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid device ID"))
		return
	}

	device, err := h.deviceService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Error(404, "Device not found"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(device))
}

func (h *DeviceHandler) List(c *gin.Context) {
	var req dto.DeviceListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	devices, total, err := h.deviceService.List(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to get devices"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.PaginationResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     devices,
	}))
}

func (h *DeviceHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid device ID"))
		return
	}

	var req dto.DeviceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	device, err := h.deviceService.Update(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	userID, _ := c.Get("userID")
	h.logService.Log(c, userID.(uint), "update_device", "device", map[string]interface{}{
		"device_id": id,
	})

	c.JSON(http.StatusOK, dto.Success(device))
}

func (h *DeviceHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid device ID"))
		return
	}

	var req dto.DeviceStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	err = h.deviceService.UpdateStatus(uint(id), req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to update device status"))
		return
	}

	userID, _ := c.Get("userID")
	h.logService.Log(c, userID.(uint), "update_device_status", "device", map[string]interface{}{
		"device_id": id,
		"status":    req.Status,
	})

	c.JSON(http.StatusOK, dto.SuccessNoData())
}

func (h *DeviceHandler) GetAvailability(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid device ID"))
		return
	}

	var req dto.DeviceAvailabilityRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	availability, err := h.deviceService.GetAvailability(uint(id), req.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(availability))
}

func (h *DeviceHandler) BatchImport(c *gin.Context) {
	var items []dto.DeviceBatchImportItem
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	result, err := h.deviceService.BatchImport(items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to batch import devices"))
		return
	}

	userID, _ := c.Get("userID")
	h.logService.Log(c, userID.(uint), "batch_import_devices", "device", map[string]interface{}{
		"success_count": result.SuccessCount,
		"fail_count":    result.FailCount,
	})

	c.JSON(http.StatusOK, dto.Success(result))
}
