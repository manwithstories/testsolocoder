package handlers

import (
	"strconv"

	"pet-adoption-platform/models"
	"pet-adoption-platform/services"
	"pet-adoption-platform/utils"

	"github.com/gin-gonic/gin"
)

func ListRescueStations(c *gin.Context) {
	var query models.RescueListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}

	rescues, total, err := services.ListRescueStations(&query)
	if err != nil {
		utils.InternalError(c, "failed to list rescue stations")
		return
	}

	utils.PaginatedSuccess(c, rescues, total, query.Page, query.PageSize)
}

func GetRescueStation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid rescue station id")
		return
	}

	rescue, err := services.GetRescueStationByID(uint(id))
	if err != nil {
		utils.NotFound(c, "rescue station not found")
		return
	}

	utils.Success(c, rescue)
}

func ReviewRescueStation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid rescue station id")
		return
	}

	userID := c.GetUint("user_id")

	var req models.ReviewRescueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	rescue, err := services.ReviewRescueStation(uint(id), userID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, rescue)
}

func GetRescueStats(c *gin.Context) {
	rescueID := c.GetUint("rescue_id")

	stats, err := services.GetRescueStats(rescueID)
	if err != nil {
		utils.InternalError(c, "failed to get stats")
		return
	}

	utils.Success(c, stats)
}

func GetRescueStatsByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid rescue station id")
		return
	}

	stats, err := services.GetRescueStats(uint(id))
	if err != nil {
		utils.InternalError(c, "failed to get stats")
		return
	}

	utils.Success(c, stats)
}

func GetAllRescuesStats(c *gin.Context) {
	stats, err := services.GetAllRescuesStats()
	if err != nil {
		utils.InternalError(c, "failed to get stats")
		return
	}

	utils.Success(c, stats)
}
