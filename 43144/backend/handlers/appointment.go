package handlers

import (
	"strconv"

	"pet-adoption-platform/models"
	"pet-adoption-platform/services"
	"pet-adoption-platform/utils"

	"github.com/gin-gonic/gin"
)

func CreateAppointment(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	appointment, err := services.CreateAppointment(&req, userID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, appointment)
}

func GetAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid appointment id")
		return
	}

	appointment, err := services.GetAppointmentByID(uint(id))
	if err != nil {
		utils.NotFound(c, "appointment not found")
		return
	}

	utils.Success(c, appointment)
}

func ListAppointments(c *gin.Context) {
	var query models.AppointmentListQuery
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

	appointments, total, err := services.ListAppointments(&query, userID, role, rescueID)
	if err != nil {
		utils.InternalError(c, "failed to list appointments")
		return
	}

	utils.PaginatedSuccess(c, appointments, total, query.Page, query.PageSize)
}

func UpdateAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid appointment id")
		return
	}

	var req models.UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	appointment, err := services.UpdateAppointment(uint(id), &req)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, appointment)
}

func CancelAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid appointment id")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	appointment, err := services.CancelAppointment(uint(id), req.Reason)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, appointment)
}

func RescheduleAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid appointment id")
		return
	}

	var req models.UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	appointment, err := services.RescheduleAppointment(uint(id), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, appointment)
}

func ConfirmAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid appointment id")
		return
	}

	appointment, err := services.ConfirmAppointment(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, appointment)
}

func CompleteAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid appointment id")
		return
	}

	appointment, err := services.CompleteAppointment(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, appointment)
}
