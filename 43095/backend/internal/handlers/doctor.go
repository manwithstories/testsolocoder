package handlers

import (
	"medical-platform/internal/middleware"
	"medical-platform/internal/services"
	"medical-platform/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	doctorService *services.DoctorService
}

func NewDoctorHandler() *DoctorHandler {
	return &DoctorHandler{
		doctorService: services.NewDoctorService(),
	}
}

func (h *DoctorHandler) GetDoctorList(c *gin.Context) {
	var query services.DoctorListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	doctors, total, err := h.doctorService.GetDoctorList(query)
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.SuccessWithPagination(c, doctors, total, query.Page, query.PageSize)
}

func (h *DoctorHandler) GetDoctorDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的医生ID")
		return
	}

	doctor, err := h.doctorService.GetDoctorByID(uint(id))
	if err != nil {
		if err.Error() == "医生不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, doctor)
}

func (h *DoctorHandler) CreateDoctor(c *gin.Context) {
	var req services.CreateDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	doctor, err := h.doctorService.CreateDoctor(req)
	if err != nil {
		if err.Error() == "该用户已是医生" || err.Error() == "用户不存在" || err.Error() == "科室不存在" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, doctor)
}

func (h *DoctorHandler) UpdateDoctor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的医生ID")
		return
	}

	currentUser := utils.GetCurrentUser(c)
	if currentUser == nil {
		utils.Unauthorized(c, "请先登录")
		return
	}

	if string(currentUser.Role) == "doctor" {
		doctor, err := h.doctorService.GetDoctorByUserID(currentUser.UserID)
		if err != nil {
			utils.Forbidden(c, "权限不足")
			return
		}
		if doctor.ID != uint(id) {
			utils.Forbidden(c, "只能修改自己的信息")
			return
		}
	}

	var req services.UpdateDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	doctor, err := h.doctorService.UpdateDoctor(uint(id), req)
	if err != nil {
		if err.Error() == "医生不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, doctor)
}

func (h *DoctorHandler) DeleteDoctor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的医生ID")
		return
	}

	if err := h.doctorService.DeleteDoctor(uint(id)); err != nil {
		if err.Error() == "医生不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, nil)
}

func (h *DoctorHandler) GetSchedules(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的医生ID")
		return
	}

	schedules, err := h.doctorService.GetSchedules(uint(id))
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, schedules)
}

func (h *DoctorHandler) CreateSchedule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的医生ID")
		return
	}

	currentUser := utils.GetCurrentUser(c)
	if currentUser == nil {
		utils.Unauthorized(c, "请先登录")
		return
	}

	if string(currentUser.Role) == "doctor" {
		doctor, err := h.doctorService.GetDoctorByUserID(currentUser.UserID)
		if err != nil {
			utils.Forbidden(c, "权限不足")
			return
		}
		if doctor.ID != uint(id) {
			utils.Forbidden(c, "只能管理自己的排班")
			return
		}
	}

	var req services.ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	schedule, err := h.doctorService.CreateSchedule(uint(id), req)
	if err != nil {
		if err.Error() == "医生不存在" || err.Error() == "该时间段已有排班" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, schedule)
}

func (h *DoctorHandler) UpdateSchedule(c *gin.Context) {
	scheduleIDStr := c.Param("scheduleId")
	scheduleID, err := strconv.ParseUint(scheduleIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的排班ID")
		return
	}

	var req services.ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	schedule, err := h.doctorService.UpdateSchedule(uint(scheduleID), req)
	if err != nil {
		if err.Error() == "排班不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, schedule)
}

func (h *DoctorHandler) DeleteSchedule(c *gin.Context) {
	scheduleIDStr := c.Param("scheduleId")
	scheduleID, err := strconv.ParseUint(scheduleIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的排班ID")
		return
	}

	if err := h.doctorService.DeleteSchedule(uint(scheduleID)); err != nil {
		if err.Error() == "排班不存在" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, nil)
}

func (h *DoctorHandler) GetAvailableTimeSlots(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的医生ID")
		return
	}

	date := c.Query("date")
	if date == "" {
		utils.BadRequest(c, "请提供日期参数 date (格式: YYYY-MM-DD)")
		return
	}

	availableSlots, err := h.doctorService.GetAvailableTimeSlots(uint(id), date)
	if err != nil {
		if err.Error() == "日期格式错误，应为 YYYY-MM-DD" {
			utils.BadRequest(c, err.Error())
			return
		}
		utils.InternalError(c, err)
		return
	}

	utils.Success(c, availableSlots)
}

func RegisterDoctorRoutes(api *gin.RouterGroup) {
	handler := NewDoctorHandler()

	doctors := api.Group("/doctors")
	{
		doctors.GET("", handler.GetDoctorList)
		doctors.GET("/:id", handler.GetDoctorDetail)
		doctors.GET("/:id/schedules", handler.GetSchedules)
		doctors.GET("/:id/available-slots", handler.GetAvailableTimeSlots)

		admin := doctors.Group("")
		admin.Use(middleware.AdminRequired())
		{
			admin.POST("", handler.CreateDoctor)
			admin.DELETE("/:id", handler.DeleteDoctor)
		}

		adminOrDoctor := doctors.Group("")
		adminOrDoctor.Use(middleware.DoctorRequired())
		{
			adminOrDoctor.PUT("/:id", handler.UpdateDoctor)
			adminOrDoctor.POST("/:id/schedules", handler.CreateSchedule)
			adminOrDoctor.PUT("/schedules/:scheduleId", handler.UpdateSchedule)
			adminOrDoctor.DELETE("/schedules/:scheduleId", handler.DeleteSchedule)
		}
	}
}
