package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"pet-board/internal/dto"
	"pet-board/internal/models"
	"pet-board/internal/service"
	"pet-board/internal/utils"
)

type ReservationHandler struct {
	resService *service.ReservationService
}

func NewReservationHandler(resService *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{resService: resService}
}

func (h *ReservationHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var req dto.ReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	reservation, err := h.resService.Create(uid, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, reservation)
}

func (h *ReservationHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	resID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid reservation ID")
		return
	}

	reservation, err := h.resService.GetByID(resID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, reservation)
}

func (h *ReservationHandler) List(c *gin.Context) {
	var query dto.StatisticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	page := utils.AtoiOrZero(c.DefaultQuery("page", "1"))
	pageSize := utils.AtoiOrZero(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var ownerID, storeID *uuid.UUID
	role := c.GetString("role")

	if ownerIDStr := c.Query("owner_id"); ownerIDStr != "" {
		id, err := uuid.Parse(ownerIDStr)
		if err == nil {
			ownerID = &id
		}
	}

	if storeIDStr := c.Query("store_id"); storeIDStr != "" {
		id, err := uuid.Parse(storeIDStr)
		if err == nil {
			storeID = &id
		}
	}

	if role == string(models.RoleOwner) && ownerID == nil {
		userID := c.GetString("user_id")
		uid, _ := uuid.Parse(userID)
		ownerID = &uid
	}

	reservations, total, err := h.resService.List(query, page, pageSize, ownerID, storeID, status)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, dto.NewPagedResult(total, page, pageSize, reservations))
}

func (h *ReservationHandler) Confirm(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	id := c.Param("id")
	resID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid reservation ID")
		return
	}

	var req dto.ReservationConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.resService.Confirm(resID, uid, &req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *ReservationHandler) CheckIn(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	id := c.Param("id")
	resID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid reservation ID")
		return
	}

	var req dto.CheckInOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.resService.CheckIn(resID, uid, &req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *ReservationHandler) CheckOut(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	id := c.Param("id")
	resID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid reservation ID")
		return
	}

	var req dto.CheckInOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.resService.CheckOut(resID, uid, &req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *ReservationHandler) Cancel(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	id := c.Param("id")
	resID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid reservation ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.resService.Cancel(resID, uid, req.Reason); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

type DailyRecordHandler struct {
	recordService *service.DailyRecordService
}

func NewDailyRecordHandler(recordService *service.DailyRecordService) *DailyRecordHandler {
	return &DailyRecordHandler{recordService: recordService}
}

func (h *DailyRecordHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var req dto.DailyRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	record, err := h.recordService.Create(uid, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, record)
}

func (h *DailyRecordHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	recID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid record ID")
		return
	}

	record, err := h.recordService.GetByID(recID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, record)
}

func (h *DailyRecordHandler) ListByReservation(c *gin.Context) {
	resID := c.Query("reservation_id")
	uid, err := uuid.Parse(resID)
	if err != nil {
		utils.BadRequest(c, "invalid reservation ID")
		return
	}

	page := utils.AtoiOrZero(c.DefaultQuery("page", "1"))
	pageSize := utils.AtoiOrZero(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	records, total, err := h.recordService.ListByReservation(uid, page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, dto.NewPagedResult(total, page, pageSize, records))
}

func (h *DailyRecordHandler) ListByPet(c *gin.Context) {
	petID := c.Param("petId")
	uid, err := uuid.Parse(petID)
	if err != nil {
		utils.BadRequest(c, "invalid pet ID")
		return
	}

	page := utils.AtoiOrZero(c.DefaultQuery("page", "1"))
	pageSize := utils.AtoiOrZero(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	records, total, err := h.recordService.ListByPet(uid, page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, dto.NewPagedResult(total, page, pageSize, records))
}

func (h *DailyRecordHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	id := c.Param("id")
	recID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid record ID")
		return
	}

	var req dto.DailyRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.recordService.Update(recID, uid, &req); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

func (h *ReviewHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var req dto.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	review, err := h.reviewService.Create(uid, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, review)
}

func (h *ReviewHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	reviewID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid review ID")
		return
	}

	review, err := h.reviewService.GetByID(reviewID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, review)
}

func (h *ReviewHandler) ListByStore(c *gin.Context) {
	storeID := c.Query("store_id")
	uid, err := uuid.Parse(storeID)
	if err != nil {
		utils.BadRequest(c, "invalid store ID")
		return
	}

	page := utils.AtoiOrZero(c.DefaultQuery("page", "1"))
	pageSize := utils.AtoiOrZero(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	reviews, total, err := h.reviewService.ListByStore(uid, page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, dto.NewPagedResult(total, page, pageSize, reviews))
}

func (h *ReviewHandler) ListByKeeper(c *gin.Context) {
	keeperID := c.Query("keeper_id")
	uid, err := uuid.Parse(keeperID)
	if err != nil {
		utils.BadRequest(c, "invalid keeper ID")
		return
	}

	page := utils.AtoiOrZero(c.DefaultQuery("page", "1"))
	pageSize := utils.AtoiOrZero(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	reviews, total, err := h.reviewService.ListByKeeper(uid, page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, dto.NewPagedResult(total, page, pageSize, reviews))
}

func (h *ReviewHandler) Reply(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	id := c.Param("id")
	reviewID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid review ID")
		return
	}

	var req dto.ReviewReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.reviewService.Reply(reviewID, uid, &req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Pay(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var req dto.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	order, err := h.orderService.Pay(uid, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, order)
}

func (h *OrderHandler) Settle(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var req dto.OrderSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	order, err := h.orderService.Settle(uid, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, order)
}

func (h *OrderHandler) Refund(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	id := c.Param("id")
	orderID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid order ID")
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
		Reason string  `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.orderService.Refund(orderID, uid, req.Amount, req.Reason); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	orderID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid order ID")
		return
	}

	order, err := h.orderService.GetByID(orderID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, order)
}

func (h *OrderHandler) List(c *gin.Context) {
	var query dto.StatisticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	page := utils.AtoiOrZero(c.DefaultQuery("page", "1"))
	pageSize := utils.AtoiOrZero(c.DefaultQuery("page_size", "10"))
	payStatus := c.Query("pay_status")
	orderType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var ownerID, storeID *uuid.UUID
	role := c.GetString("role")

	if ownerIDStr := c.Query("owner_id"); ownerIDStr != "" {
		id, err := uuid.Parse(ownerIDStr)
		if err == nil {
			ownerID = &id
		}
	}

	if storeIDStr := c.Query("store_id"); storeIDStr != "" {
		id, err := uuid.Parse(storeIDStr)
		if err == nil {
			storeID = &id
		}
	}

	if role == string(models.RoleOwner) && ownerID == nil {
		userID := c.GetString("user_id")
		uid, _ := uuid.Parse(userID)
		ownerID = &uid
	}

	orders, total, err := h.orderService.List(ownerID, storeID, payStatus, orderType, query, page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, dto.NewPagedResult(total, page, pageSize, orders))
}

func (h *OrderHandler) GetByReservationID(c *gin.Context) {
	resID := c.Param("reservationId")
	uid, err := uuid.Parse(resID)
	if err != nil {
		utils.BadRequest(c, "invalid reservation ID")
		return
	}

	orders, err := h.orderService.GetByReservationID(uid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, orders)
}

type HealthAlertHandler struct {
	alertService *service.HealthAlertService
}

func NewHealthAlertHandler(alertService *service.HealthAlertService) *HealthAlertHandler {
	return &HealthAlertHandler{alertService: alertService}
}

func (h *HealthAlertHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var isRead *bool
	if isReadStr := c.Query("is_read"); isReadStr != "" {
		b := isReadStr == "true"
		isRead = &b
	}

	page := utils.AtoiOrZero(c.DefaultQuery("page", "1"))
	pageSize := utils.AtoiOrZero(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	alerts, total, err := h.alertService.ListByUser(uid, isRead, page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, dto.NewPagedResult(total, page, pageSize, alerts))
}

func (h *HealthAlertHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	id := c.Param("id")
	alertID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid alert ID")
		return
	}

	if err := h.alertService.MarkAsRead(alertID, uid); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *HealthAlertHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	if err := h.alertService.MarkAllAsRead(uid); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

type StatisticsHandler struct {
	statService *service.StatisticsService
}

func NewStatisticsHandler(statService *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{statService: statService}
}

func (h *StatisticsHandler) GetRevenueTrend(c *gin.Context) {
	storeID := c.Query("store_id")
	uid, err := uuid.Parse(storeID)
	if err != nil {
		utils.BadRequest(c, "invalid store ID")
		return
	}

	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	trend, err := h.statService.GetRevenueTrend(uid, start, end)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, trend)
}

func (h *StatisticsHandler) GetOccupancyRate(c *gin.Context) {
	storeID := c.Query("store_id")
	uid, err := uuid.Parse(storeID)
	if err != nil {
		utils.BadRequest(c, "invalid store ID")
		return
	}

	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	rate, err := h.statService.GetOccupancyRate(uid, start, end)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{"occupancy_rate": rate})
}

func (h *StatisticsHandler) GetPetTypeDistribution(c *gin.Context) {
	storeID := c.Query("store_id")
	uid, err := uuid.Parse(storeID)
	if err != nil {
		utils.BadRequest(c, "invalid store ID")
		return
	}

	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	distribution, err := h.statService.GetPetTypeDistribution(uid, start, end)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, distribution)
}

func (h *StatisticsHandler) GetOrderStatistics(c *gin.Context) {
	storeID := c.Query("store_id")
	uid, err := uuid.Parse(storeID)
	if err != nil {
		utils.BadRequest(c, "invalid store ID")
		return
	}

	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	stats, err := h.statService.GetOrderStatistics(uid, start, end)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, stats)
}

type ExportHandler struct {
	statService *service.StatisticsService
}

func NewExportHandler(statService *service.StatisticsService) *ExportHandler {
	return &ExportHandler{statService: statService}
}

func (h *ExportHandler) ExportExcel(c *gin.Context) {
	storeID := c.Query("store_id")
	uid, err := uuid.Parse(storeID)
	if err != nil {
		utils.BadRequest(c, "invalid store ID")
		return
	}

	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	revenue, err := h.statService.GetRevenueTrend(uid, start, end)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	occupancy, _ := h.statService.GetOccupancyRate(uid, start, end)
	petDist, _ := h.statService.GetPetTypeDistribution(uid, start, end)
	orderStats, _ := h.statService.GetOrderStatistics(uid, start, end)

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=statistics.xlsx")

	excelContent := generateExcelReport(revenue, occupancy, petDist, orderStats, startDate, endDate)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelContent)
}

func (h *ExportHandler) ExportPDF(c *gin.Context) {
	storeID := c.Query("store_id")
	uid, err := uuid.Parse(storeID)
	if err != nil {
		utils.BadRequest(c, "invalid store ID")
		return
	}

	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	revenue, err := h.statService.GetRevenueTrend(uid, start, end)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	occupancy, _ := h.statService.GetOccupancyRate(uid, start, end)
	petDist, _ := h.statService.GetPetTypeDistribution(uid, start, end)
	orderStats, _ := h.statService.GetOrderStatistics(uid, start, end)

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=statistics.pdf")

	pdfContent := generatePDFReport(revenue, occupancy, petDist, orderStats, startDate, endDate)
	c.Data(http.StatusOK, "application/pdf", pdfContent)
}

func generateExcelReport(revenue []map[string]interface{}, occupancy float64, petDist []map[string]interface{}, orderStats map[string]interface{}, startDate, endDate string) []byte {
	return []byte("Excel report placeholder")
}

func generatePDFReport(revenue []map[string]interface{}, occupancy float64, petDist []map[string]interface{}, orderStats map[string]interface{}, startDate, endDate string) []byte {
	return []byte("PDF report placeholder")
}
