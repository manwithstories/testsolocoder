package api

import (
	"gym-management/internal/middleware"
	"gym-management/internal/models"
	"gym-management/internal/pkg/utils"
	"gym-management/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courseService   service.CourseService
	scheduleService service.ScheduleService
}

func NewCourseHandler() *CourseHandler {
	return &CourseHandler{
		courseService:   service.NewCourseService(),
		scheduleService: service.NewScheduleService(),
	}
}

func (h *CourseHandler) RegisterRoutes(r *gin.RouterGroup) {
	course := r.Group("/courses")
	course.Use(middleware.JWTAuth())
	{
		course.POST("/", h.Create)
		course.GET("/", h.List)
		course.GET("/:id", h.GetByID)
		course.PUT("/:id", h.Update)
		course.DELETE("/:id", h.Delete)
		course.PATCH("/:id/status", h.UpdateStatus)
		course.POST("/:id/generate-schedules", h.GenerateSchedules)

		schedule := course.Group("/schedules")
		{
			schedule.GET("/", h.ListSchedules)
			schedule.GET("/available", h.GetAvailableSchedules)
			schedule.GET("/:id", h.GetScheduleByID)
			schedule.PATCH("/:id/status", h.UpdateScheduleStatus)
		}
	}
}

type CreateCourseRequest struct {
	Name        string               `json:"name" binding:"required"`
	Description string               `json:"description"`
	CoachID     uint                 `json:"coach_id" binding:"required"`
	Capacity    int                  `json:"capacity" binding:"required,min=1"`
	Duration    int                  `json:"duration" binding:"required,min=1"`
	Type        models.CourseType    `json:"type" binding:"oneof=single weekly monthly"`
	Weekdays    string               `json:"weekdays"`
	StartDate   string               `json:"start_date" binding:"required"`
	EndDate     string               `json:"end_date"`
	StartTime   string               `json:"start_time" binding:"required"`
	Location    string               `json:"location"`
}

func (h *CourseHandler) Create(c *gin.Context) {
	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	startDate, err := time.ParseInLocation("2006-01-02", req.StartDate, time.Local)
	if err != nil {
		utils.BadRequest(c, "开始日期格式错误，请使用 YYYY-MM-DD", nil)
		return
	}

	course := &models.Course{
		Name:        req.Name,
		Description: req.Description,
		CoachID:     req.CoachID,
		Capacity:    req.Capacity,
		Duration:    req.Duration,
		Type:        req.Type,
		Weekdays:    req.Weekdays,
		StartDate:   startDate,
		StartTime:   req.StartTime,
		Location:    req.Location,
	}

	if req.EndDate != "" {
		endDate, err := time.ParseInLocation("2006-01-02", req.EndDate, time.Local)
		if err != nil {
			utils.BadRequest(c, "结束日期格式错误，请使用 YYYY-MM-DD", nil)
			return
		}
		course.EndDate = &endDate
	}

	if err := h.courseService.Create(course); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, course)
}

func (h *CourseHandler) List(c *gin.Context) {
	page, pageSize := utils.GetPageInfo(c)
	keyword := c.Query("keyword")
	coachID, _ := strconv.ParseUint(c.Query("coach_id"), 10, 32)
	courseType := c.Query("type")

	courses, total, err := h.courseService.List(page, pageSize, keyword, uint(coachID), courseType)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPagination(c, courses, page, pageSize, total)
}

func (h *CourseHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	course, err := h.courseService.GetByID(uint(id))
	if err != nil {
		utils.NotFound(c, "课程不存在")
		return
	}

	utils.Success(c, course)
}

func (h *CourseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "updated_at")

	if err := h.courseService.Update(uint(id), updates); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}

func (h *CourseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	if err := h.courseService.Delete(uint(id)); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}

func (h *CourseHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=1 2 3"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := h.courseService.UpdateStatus(uint(id), req.Status); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}

func (h *CourseHandler) GenerateSchedules(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	if err := h.courseService.GenerateSchedules(uint(id)); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}

func (h *CourseHandler) ListSchedules(c *gin.Context) {
	page, pageSize := utils.GetPageInfo(c)
	courseID, _ := strconv.ParseUint(c.Query("course_id"), 10, 32)

	var startDate, endDate *time.Time
	if startStr := c.Query("start_date"); startStr != "" {
		if t, err := time.ParseInLocation("2006-01-02", startStr, time.Local); err == nil {
			startDate = &t
		}
	}
	if endStr := c.Query("end_date"); endStr != "" {
		if t, err := time.ParseInLocation("2006-01-02", endStr, time.Local); err == nil {
			endDate = &t
		}
	}

	schedules, total, err := h.scheduleService.List(page, pageSize, uint(courseID), startDate, endDate)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPagination(c, schedules, page, pageSize, total)
}

func (h *CourseHandler) GetAvailableSchedules(c *gin.Context) {
	schedules, err := h.scheduleService.FindAvailable()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, schedules)
}

func (h *CourseHandler) GetScheduleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	schedule, err := h.scheduleService.GetByID(uint(id))
	if err != nil {
		utils.NotFound(c, "课程排期不存在")
		return
	}

	utils.Success(c, schedule)
}

func (h *CourseHandler) UpdateScheduleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的ID", nil)
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=1 2 3 4"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := h.scheduleService.UpdateStatus(uint(id), req.Status); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.Success(c, nil)
}
