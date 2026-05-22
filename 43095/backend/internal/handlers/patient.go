package handlers

import (
	"medical-platform/internal/middleware"
	"medical-platform/internal/services"
	"medical-platform/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	patientService *services.PatientService
}

func NewPatientHandler() *PatientHandler {
	return &PatientHandler{
		patientService: services.NewPatientService(),
	}
}

func (h *PatientHandler) GetPatientList(c *gin.Context) {
	var query services.PatientListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	patients, total, err := h.patientService.GetPatientList(query)
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.SuccessWithPagination(c, patients, total, query.Page, query.PageSize)
}

func (h *PatientHandler) GetPatientDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的患者ID")
		return
	}

	currentUser := utils.GetCurrentUser(c)
	if currentUser == nil {
		utils.Unauthorized(c, "请先登录")
		return
	}

	if string(currentUser.Role) == "patient" {
		patient, err := h.patientService.GetPatientByUserID(currentUser.UserID)
		if err != nil {
			utils.Forbidden(c, "权限不足")
			return
		}
		if patient.ID != uint(id) {
			utils.Forbidden(c, "只能查看自己的信息")
			return
		}
	}

	patient, err := h.patientService.GetPatientByID(uint(id))
	if err != nil {
		if err.Error() == "患者不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, patient)
}

func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的患者ID")
		return
	}

	currentUser := utils.GetCurrentUser(c)
	if currentUser == nil {
		utils.Unauthorized(c, "请先登录")
		return
	}

	patient, err := h.patientService.GetPatientByUserID(currentUser.UserID)
	if err != nil {
		utils.Forbidden(c, "权限不足")
		return
	}
	if patient.ID != uint(id) {
		utils.Forbidden(c, "只能修改自己的信息")
		return
	}

	var req services.UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	updatedPatient, err := h.patientService.UpdatePatient(uint(id), req)
	if err != nil {
		if err.Error() == "患者不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, updatedPatient)
}

func (h *PatientHandler) DeletePatient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的患者ID")
		return
	}

	if err := h.patientService.DeletePatient(uint(id)); err != nil {
		if err.Error() == "患者不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, nil)
}

func RegisterPatientRoutes(api *gin.RouterGroup) {
	handler := NewPatientHandler()

	patients := api.Group("/patients")
	{
		patients.Use(middleware.Auth())
		{
			patients.GET("/:id", handler.GetPatientDetail)
			patients.PUT("/:id", handler.UpdatePatient)
		}

		admin := patients.Group("")
		admin.Use(middleware.AdminRequired())
		{
			admin.GET("", handler.GetPatientList)
			admin.DELETE("/:id", handler.DeletePatient)
		}
	}
}
