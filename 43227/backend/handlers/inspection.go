package handlers

import (
	"net/http"
	"strconv"
	"time"

	"beehive-platform/database"
	"beehive-platform/models"
	"beehive-platform/utils"

	"github.com/gin-gonic/gin"
)

type InspectionHandler struct{}

func NewInspectionHandler() *InspectionHandler {
	return &InspectionHandler{}
}

func (h *InspectionHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var inventory models.Inventory
	if err := database.DB.Where("id = ? AND user_id = ?", req.InventoryID, userID).First(&inventory).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "inventory not found")
		return
	}

	appointmentDate, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid date format")
		return
	}

	inspection := models.Inspection{
		UserID:          userID.(uint),
		InventoryID:     req.InventoryID,
		BatchCode:       req.BatchCode,
		AppointmentDate: appointmentDate,
		Status:          "pending",
		Notes:           req.Notes,
	}

	if err := database.DB.Create(&inspection).Error; err != nil {
		utils.FailWithError(c, http.StatusInternalServerError, "failed to create inspection", err)
		return
	}

	utils.Success(c, inspection)
}

func (h *InspectionHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")

	var pageParams utils.PageParams
	if err := c.ShouldBindQuery(&pageParams); err != nil {
		pageParams = utils.PageParams{Page: 1, PageSize: 10}
	}
	if pageParams.Page < 1 {
		pageParams.Page = 1
	}
	if pageParams.PageSize < 1 {
		pageParams.PageSize = 10
	}

	query := database.DB.Model(&models.Inspection{})

	if role == "beekeeper" {
		query = query.Where("user_id = ?", userID)
	} else if role == "inspector" {
		inspectorID := c.Query("inspector_id")
		if inspectorID != "" {
			query = query.Where("inspector_id = ?", inspectorID)
		}
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	sortBy := c.DefaultQuery("sort_by", "appointment_date")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	query = query.Order(sortBy + " " + sortOrder)

	var total int64
	query.Count(&total)

	var inspections []models.Inspection
	query.Offset(pageParams.GetOffset()).Limit(pageParams.PageSize).
		Preload("User").Preload("Inventory").Find(&inspections)

	utils.SuccessWithTotal(c, inspections, total)
}

func (h *InspectionHandler) Get(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var inspection models.Inspection
	query := database.DB.Preload("User").Preload("Inventory.Harvest.Beehive").Where("id = ?", id)

	if role == "beekeeper" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&inspection).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "inspection not found")
		return
	}

	utils.Success(c, inspection)
}

func (h *InspectionHandler) AssignInspector(c *gin.Context) {
	role, _ := c.Get("user_role")
	if role != "inspector" {
		utils.Fail(c, http.StatusForbidden, "only inspectors can assign themselves")
		return
	}

	inspectorID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var inspection models.Inspection
	if err := database.DB.First(&inspection, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "inspection not found")
		return
	}

	if inspection.Status != "pending" {
		utils.Fail(c, http.StatusBadRequest, "inspection has already been assigned")
		return
	}

	uid := inspectorID.(uint)
	database.DB.Model(&inspection).Updates(map[string]interface{}{
		"inspector_id": &uid,
		"status":       "confirmed",
	})

	userNotification := models.Notification{
		UserID:    inspection.UserID,
		Type:      "inspection_result",
		Title:     "检测已预约成功",
		Content:   "您的检测预约已被接受，请在预约时间准备好样品",
		RelatedID: &inspection.ID,
	}
	database.DB.Create(&userNotification)

	utils.Success(c, nil)
}

func (h *InspectionHandler) SubmitResult(c *gin.Context) {
	inspectorID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req models.UpdateInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var inspection models.Inspection
	if err := database.DB.Where("id = ? AND inspector_id = ?", id, inspectorID).First(&inspection).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "inspection not found or unauthorized")
		return
	}

	tx := database.DB.Begin()

	updates := map[string]interface{}{
		"status": "completed",
	}
	now := time.Now()
	updates["inspection_date"] = &now

	if req.ReportURL != nil {
		updates["report_url"] = *req.ReportURL
	}
	if req.Result != nil {
		updates["result"] = *req.Result
	}
	if req.Grade != nil {
		updates["grade"] = *req.Grade
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}

	tx.Model(&inspection).Updates(updates)

	if req.Grade != nil {
		tx.Model(&models.Inventory{}).Where("id = ?", inspection.InventoryID).
			Update("grade", *req.Grade)
	}

	userNotification := models.Notification{
		UserID:    inspection.UserID,
		Type:      "inspection_result",
		Title:     "检测结果已出",
		Content:   "您的蜂蜜检测已完成，检测等级：" + inspection.Grade,
		RelatedID: &inspection.ID,
	}
	tx.Create(&userNotification)

	tx.Commit()

	utils.Success(c, nil)
}

func (h *InspectionHandler) Cancel(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var inspection models.Inspection
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&inspection).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "inspection not found")
		return
	}

	if inspection.Status != "pending" && inspection.Status != "confirmed" {
		utils.Fail(c, http.StatusBadRequest, "inspection cannot be cancelled")
		return
	}

	database.DB.Model(&inspection).Update("status", "cancelled")

	utils.Success(c, nil)
}
