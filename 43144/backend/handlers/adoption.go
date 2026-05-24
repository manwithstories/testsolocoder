package handlers

import (
	"strconv"

	"pet-adoption-platform/models"
	"pet-adoption-platform/services"
	"pet-adoption-platform/utils"

	"github.com/gin-gonic/gin"
)

func CreateAdoptionApplication(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.CreateAdoptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	app, err := services.CreateAdoptionApplication(&req, userID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, app)
}

func GetAdoptionApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid application id")
		return
	}

	app, err := services.GetAdoptionApplicationByID(uint(id))
	if err != nil {
		utils.NotFound(c, "application not found")
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")
	rescueID := c.GetUint("rescue_id")

	if role == "adopter" && app.AdopterID != userID {
		utils.Forbidden(c, "you can only view your own applications")
		return
	}
	if role == "rescue" && app.RescueID != rescueID {
		utils.Forbidden(c, "you can only view your own rescue station's applications")
		return
	}

	utils.Success(c, app)
}

func ListAdoptionApplications(c *gin.Context) {
	var query models.AdoptionListQuery
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

	userID := c.GetUint("user_id")
	role := c.GetString("role")

	var rescueID *uint
	if role == "rescue" {
		rid := c.GetUint("rescue_id")
		rescueID = &rid
	}

	apps, total, err := services.ListAdoptionApplications(&query, userID, role, rescueID)
	if err != nil {
		utils.InternalError(c, "failed to list applications")
		return
	}

	utils.PaginatedSuccess(c, apps, total, query.Page, query.PageSize)
}

func ReviewAdoptionApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid application id")
		return
	}

	userID := c.GetUint("user_id")
	rescueID := c.GetUint("rescue_id")

	app, err := services.GetAdoptionApplicationByID(uint(id))
	if err != nil {
		utils.NotFound(c, "application not found")
		return
	}

	if app.RescueID != rescueID {
		utils.Forbidden(c, "you can only review your own rescue station's applications")
		return
	}

	var req models.ReviewAdoptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	app, err = services.ReviewAdoptionApplication(uint(id), userID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, app)
}

func SignAdoptionAgreement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid application id")
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")

	app, err := services.GetAdoptionApplicationByID(uint(id))
	if err != nil {
		utils.NotFound(c, "application not found")
		return
	}

	if role == "adopter" && app.AdopterID != userID {
		utils.Forbidden(c, "you can only sign your own agreement")
		return
	}

	agreement, err := services.SignAdoptionAgreement(uint(id), role, userID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, agreement)
}

func GetAdoptionAgreement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid application id")
		return
	}

	agreement, err := services.GetAdoptionAgreement(uint(id))
	if err != nil {
		utils.NotFound(c, "agreement not found")
		return
	}

	utils.Success(c, agreement)
}

func CompleteAdoption(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid application id")
		return
	}

	app, err := services.CompleteAdoption(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, app)
}

func CreateFollowUpRecord(c *gin.Context) {
	userID := c.GetUint("user_id")
	rescueID := c.GetUint("rescue_id")

	var req models.CreateFollowUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	record, err := services.CreateFollowUpRecord(&req, userID, rescueID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, record)
}

func ListFollowUpRecords(c *gin.Context) {
	petID, err := strconv.ParseUint(c.Param("pet_id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid pet id")
		return
	}

	records, err := services.ListFollowUpRecords(uint(petID))
	if err != nil {
		utils.InternalError(c, "failed to list follow-up records")
		return
	}

	utils.Success(c, records)
}
