package handler

import (
	"drone-rental/internal/dto"
	"drone-rental/internal/middleware"
	"drone-rental/internal/model"
	"drone-rental/internal/pkg/response"
	"drone-rental/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FlightHandler struct {
	flightService *service.FlightService
}

func NewFlightHandler() *FlightHandler {
	return &FlightHandler{
		flightService: service.NewFlightService(),
	}
}

func (h *FlightHandler) Create(c *gin.Context) {
	pilotID := middleware.GetUserID(c)
	var req dto.CreateFlightRecordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	flight, err := h.flightService.Create(pilotID, &req)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, flight)
}

func (h *FlightHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	flight, err := h.flightService.GetByID(uint(id))
	if err != nil {
		response.ErrNotFound(c, "飞行记录不存在")
		return
	}
	response.Success(c, flight)
}

func (h *FlightHandler) Update(c *gin.Context) {
	pilotID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.CreateFlightRecordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.flightService.Update(uint(id), pilotID, &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *FlightHandler) Delete(c *gin.Context) {
	pilotID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.flightService.Delete(uint(id), pilotID); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *FlightHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	droneID, _ := strconv.ParseUint(c.Query("drone_id"), 10, 64)
	pilotID, _ := strconv.ParseUint(c.Query("pilot_id"), 10, 64)
	orderID, _ := strconv.ParseUint(c.Query("order_id"), 10, 64)
	serviceID, _ := strconv.ParseUint(c.Query("service_id"), 10, 64)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	flights, total, err := h.flightService.List(page, pageSize, uint(droneID), uint(pilotID), uint(orderID), uint(serviceID), startDate, endDate)
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Page(c, flights, total, page, pageSize)
}

type InsuranceHandler struct {
	insuranceService *service.InsuranceService
}

func NewInsuranceHandler() *InsuranceHandler {
	return &InsuranceHandler{
		insuranceService: service.NewInsuranceService(),
	}
}

func (h *InsuranceHandler) CreateClaim(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.CreateClaimReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	claim, err := h.insuranceService.CreateClaim(userID, &req)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, claim)
}

func (h *InsuranceHandler) GetClaimByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	claim, err := h.insuranceService.GetClaimByID(uint(id))
	if err != nil {
		response.ErrNotFound(c, "理赔申请不存在")
		return
	}
	response.Success(c, claim)
}

func (h *InsuranceHandler) ListClaims(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	orderID, _ := strconv.ParseUint(c.Query("order_id"), 10, 64)
	status := model.ClaimStatus(c.Query("status"))
	claims, total, err := h.insuranceService.ListClaims(page, pageSize, uint(orderID), status)
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Page(c, claims, total, page, pageSize)
}

func (h *InsuranceHandler) ReviewClaim(c *gin.Context) {
	reviewerID := middleware.GetUserID(c)
	var req dto.ReviewClaimReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.insuranceService.ReviewClaim(reviewerID, &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{
		reviewService: service.NewReviewService(),
	}
}

func (h *ReviewHandler) Create(c *gin.Context) {
	reviewerID := middleware.GetUserID(c)
	var req dto.CreateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	review, err := h.reviewService.Create(reviewerID, &req)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, review)
}

func (h *ReviewHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	review, err := h.reviewService.GetByID(uint(id))
	if err != nil {
		response.ErrNotFound(c, "评价不存在")
		return
	}
	response.Success(c, review)
}

func (h *ReviewHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	orderID, _ := strconv.ParseUint(c.Query("order_id"), 10, 64)
	serviceID, _ := strconv.ParseUint(c.Query("service_id"), 10, 64)
	revieweeID, _ := strconv.ParseUint(c.Query("reviewee_id"), 10, 64)
	droneID, _ := strconv.ParseUint(c.Query("drone_id"), 10, 64)
	reviewType := model.ReviewType(c.Query("type"))
	reviews, total, err := h.reviewService.List(page, pageSize, uint(orderID), uint(serviceID), uint(revieweeID), uint(droneID), reviewType)
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Page(c, reviews, total, page, pageSize)
}

func (h *ReviewHandler) Reply(c *gin.Context) {
	reviewerID := middleware.GetUserID(c)
	var req dto.ReplyReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.reviewService.Reply(reviewerID, &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{
		statsService: service.NewStatsService(),
	}
}

func (h *StatsHandler) Revenue(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")
	stats, err := h.statsService.GetRevenueStats(startDate, endDate)
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Success(c, stats)
}

func (h *StatsHandler) Region(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")
	stats, err := h.statsService.GetRegionStats(startDate, endDate)
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Success(c, stats)
}

func (h *StatsHandler) Drone(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")
	stats, err := h.statsService.GetDroneStats(startDate, endDate)
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Success(c, stats)
}
