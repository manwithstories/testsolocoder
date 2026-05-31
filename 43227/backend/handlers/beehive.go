package handlers

import (
	"net/http"
	"strconv"

	"beehive-platform/database"
	"beehive-platform/models"
	"beehive-platform/utils"

	"github.com/gin-gonic/gin"
)

type BeehiveHandler struct{}

func NewBeehiveHandler() *BeehiveHandler {
	return &BeehiveHandler{}
}

func (h *BeehiveHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateBeehiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var existing models.Beehive
	if err := database.DB.Where("code = ?", req.Code).First(&existing).Error; err == nil {
		utils.Fail(c, http.StatusConflict, "beehive code already exists")
		return
	}

	beehive := models.Beehive{
		UserID:       userID.(uint),
		Name:         req.Name,
		Code:         req.Code,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Region:       req.Region,
		BeeSpecies:   req.BeeSpecies,
		GroupName:    req.GroupName,
		QueenStatus:  req.QueenStatus,
		WorkerCount:  req.WorkerCount,
		Notes:        req.Notes,
		Status:       "active",
		HealthStatus: "healthy",
	}

	if err := database.DB.Create(&beehive).Error; err != nil {
		utils.FailWithError(c, http.StatusInternalServerError, "failed to create beehive", err)
		return
	}

	utils.Success(c, beehive)
}

func (h *BeehiveHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")

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

	query := database.DB.Model(&models.Beehive{}).Where("user_id = ?", userID)

	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name ILIKE ? OR code ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if region := c.Query("region"); region != "" {
		query = query.Where("region = ?", region)
	}
	if beeSpecies := c.Query("bee_species"); beeSpecies != "" {
		query = query.Where("bee_species = ?", beeSpecies)
	}
	if groupName := c.Query("group_name"); groupName != "" {
		query = query.Where("group_name = ?", groupName)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if healthStatus := c.Query("health_status"); healthStatus != "" {
		query = query.Where("health_status = ?", healthStatus)
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	query = query.Order(sortBy + " " + sortOrder)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to count beehives")
		return
	}

	var beehives []models.Beehive
	if err := query.Offset(pageParams.GetOffset()).Limit(pageParams.PageSize).Find(&beehives).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to fetch beehives")
		return
	}

	utils.SuccessWithTotal(c, beehives, total)
}

func (h *BeehiveHandler) Get(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var beehive models.Beehive
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&beehive).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "beehive not found")
		return
	}

	utils.Success(c, beehive)
}

func (h *BeehiveHandler) Update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req models.UpdateBeehiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var beehive models.Beehive
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&beehive).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "beehive not found")
		return
	}

	if req.Name != nil {
		beehive.Name = *req.Name
	}
	if req.Latitude != nil {
		beehive.Latitude = *req.Latitude
	}
	if req.Longitude != nil {
		beehive.Longitude = *req.Longitude
	}
	if req.Region != nil {
		beehive.Region = *req.Region
	}
	if req.BeeSpecies != nil {
		beehive.BeeSpecies = *req.BeeSpecies
	}
	if req.GroupName != nil {
		beehive.GroupName = *req.GroupName
	}
	if req.Status != nil {
		beehive.Status = *req.Status
	}
	if req.HealthStatus != nil {
		beehive.HealthStatus = *req.HealthStatus
	}
	if req.QueenStatus != nil {
		beehive.QueenStatus = *req.QueenStatus
	}
	if req.WorkerCount != nil {
		beehive.WorkerCount = *req.WorkerCount
	}
	if req.Notes != nil {
		beehive.Notes = *req.Notes
	}

	database.DB.Save(&beehive)

	utils.Success(c, beehive)
}

func (h *BeehiveHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Beehive{})
	if result.Error != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to delete beehive")
		return
	}
	if result.RowsAffected == 0 {
		utils.Fail(c, http.StatusNotFound, "beehive not found")
		return
	}

	utils.Success(c, nil)
}

func (h *BeehiveHandler) Groups(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var groups []struct {
		GroupName string `json:"group_name"`
		Count     int64  `json:"count"`
	}

	database.DB.Model(&models.Beehive{}).
		Where("user_id = ? AND group_name != ''", userID).
		Select("group_name, count(*) as count").
		Group("group_name").
		Find(&groups)

	utils.Success(c, groups)
}
