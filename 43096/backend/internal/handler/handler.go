package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"ticket-system/internal/dto"
	"ticket-system/internal/middleware"
	"ticket-system/internal/service"
)

type Handler struct {
	userSvc      *service.UserService
	showSvc      *service.ShowService
	orderSvc     *service.OrderService
	checkinSvc   *service.CheckinService
	statisticsSvc *service.StatisticsService
}

func NewHandler() *Handler {
	return &Handler{
		userSvc:      service.NewUserService(),
		showSvc:      service.NewShowService(),
		orderSvc:     service.NewOrderService(),
		checkinSvc:   service.NewCheckinService(),
		statisticsSvc: service.NewStatisticsService(),
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误: " + err.Error()})
		return
	}

	user, err := h.userSvc.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "注册成功", Data: user})
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	resp, err := h.userSvc.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{Code: 401, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "登录成功", Data: resp})
}

func (h *Handler) GetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userSvc.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: user})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	user, err := h.userSvc.UpdateUser(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "更新成功", Data: user})
}

func (h *Handler) GetMemberLevels(c *gin.Context) {
	levels, err := h.userSvc.GetMemberLevels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: levels})
}

func (h *Handler) ExchangeCoupon(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.CouponExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	coupon, err := h.userSvc.ExchangeCoupon(userID, req.Points)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "兑换成功", Data: coupon})
}

func (h *Handler) GetUserCoupons(c *gin.Context) {
	userID := middleware.GetUserID(c)
	coupons, err := h.userSvc.GetUserCoupons(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: coupons})
}

func (h *Handler) CreateShow(c *gin.Context) {
	var req dto.ShowCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误: " + err.Error()})
		return
	}

	show, err := h.showSvc.CreateShow(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: show})
}

func (h *Handler) GetShow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	show, err := h.showSvc.GetShow(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "演出不存在"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: show})
}

func (h *Handler) ListShows(c *gin.Context) {
	var req dto.ShowListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	resp, err := h.showSvc.ListShows(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: resp})
}

func (h *Handler) UpdateShow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.ShowUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	show, err := h.showSvc.UpdateShow(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "更新成功", Data: show})
}

func (h *Handler) DeleteShow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err := h.showSvc.DeleteShow(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "删除失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "删除成功"})
}

func (h *Handler) CreateSession(c *gin.Context) {
	var req dto.SessionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误: " + err.Error()})
		return
	}

	session, err := h.showSvc.CreateSession(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: session})
}

func (h *Handler) GetSessions(c *gin.Context) {
	showID, _ := strconv.ParseUint(c.Param("showId"), 10, 64)
	sessions, err := h.showSvc.GetSessions(showID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: sessions})
}

func (h *Handler) CreateSeatArea(c *gin.Context) {
	var req dto.SeatAreaCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误: " + err.Error()})
		return
	}

	area, err := h.showSvc.CreateSeatArea(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: area})
}

func (h *Handler) GetSeatAreas(c *gin.Context) {
	sessionID, _ := strconv.ParseUint(c.Param("sessionId"), 10, 64)
	areas, err := h.showSvc.GetSeatAreas(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: areas})
}

func (h *Handler) CreateSeat(c *gin.Context) {
	var req dto.SeatCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误: " + err.Error()})
		return
	}

	seat, err := h.showSvc.CreateSeat(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: seat})
}

func (h *Handler) BatchCreateSeats(c *gin.Context) {
	var req dto.SeatBatchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误: " + err.Error()})
		return
	}

	count, err := h.showSvc.BatchCreateSeats(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "批量创建成功", Data: gin.H{"created_count": count}})
}

func (h *Handler) GetSeats(c *gin.Context) {
	sessionID, _ := strconv.ParseUint(c.Param("sessionId"), 10, 64)
	seats, err := h.showSvc.GetSeats(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: seats})
}

func (h *Handler) LockSeats(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.SeatLockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	failedSeats, err := h.showSvc.LockSeats(req.SessionID, req.SeatIDs, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Code:    400,
			Message: err.Error(),
			Data:    gin.H{"failed_seats": failedSeats},
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "锁定成功"})
}

func (h *Handler) UnlockSeats(c *gin.Context) {
	var req dto.SeatUnlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	err := h.showSvc.UnlockSeats(req.SessionID, req.SeatIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "解锁成功"})
}

func (h *Handler) UpdateSeatChart(c *gin.Context) {
	var req dto.SeatChartUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	chart, err := h.showSvc.UpdateSeatChart(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "更新成功", Data: chart})
}

func (h *Handler) GetSeatChart(c *gin.Context) {
	sessionID, _ := strconv.ParseUint(c.Param("sessionId"), 10, 64)
	chart, err := h.showSvc.GetSeatChart(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: chart})
}

func (h *Handler) CreateOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误: " + err.Error()})
		return
	}

	order, err := h.orderSvc.CreateOrder(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "创建成功", Data: order})
}

func (h *Handler) GetOrder(c *gin.Context) {
	orderNo := c.Param("orderNo")
	order, err := h.orderSvc.GetOrder(orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Code: 404, Message: "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: order})
}

func (h *Handler) ListOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.OrderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	role, _ := c.Get("role")
	queryUserID := userID
	if role == "admin" {
		queryUserID = 0
	}

	resp, err := h.orderSvc.ListOrders(queryUserID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: resp})
}

func (h *Handler) PayOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.OrderPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	order, err := h.orderSvc.PayOrder(req.OrderNo, req.PayType, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "支付成功", Data: order})
}

func (h *Handler) CancelOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	orderNo := c.Param("orderNo")
	err := h.orderSvc.CancelOrder(orderNo, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "取消成功"})
}

func (h *Handler) RequestRefund(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.RefundCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	refund, err := h.orderSvc.RequestRefund(req.OrderNo, userID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "申请成功", Data: refund})
}

func (h *Handler) AuditRefund(c *gin.Context) {
	var req dto.RefundAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	err := h.orderSvc.AuditRefund(req.RefundNo, req.Status, req.AuditRemark)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "审核完成"})
}

func (h *Handler) ExportOrders(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))

	file, err := h.orderSvc.ExportOrders(startDate, endDate, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "导出失败"})
		return
	}

	filename := "orders_" + strings.ReplaceAll(startDate, "-", "") + "_" + strings.ReplaceAll(endDate, "-", "") + ".xlsx"
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	_ = file.Write(c.Writer)
}

func (h *Handler) Checkin(c *gin.Context) {
	operatorID := middleware.GetUserID(c)
	var req dto.CheckinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	result, err := h.checkinSvc.Checkin(req.TicketNo, operatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetSalesStatistics(c *gin.Context) {
	var req dto.StatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	stats, err := h.statisticsSvc.GetSalesStatistics(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: stats})
}

func (h *Handler) GetAreaSales(c *gin.Context) {
	sessionID, _ := strconv.ParseUint(c.Param("sessionId"), 10, 64)
	sales, err := h.statisticsSvc.GetAreaSales(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: sales})
}

func (h *Handler) GetSeatHeatmap(c *gin.Context) {
	sessionID, _ := strconv.ParseUint(c.Param("sessionId"), 10, 64)
	data, err := h.statisticsSvc.GetSeatHeatmap(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: data})
}

func (h *Handler) GetAudienceProfile(c *gin.Context) {
	var req dto.StatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	profile, err := h.statisticsSvc.GetAudienceProfile(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: profile})
}

func (h *Handler) GetDailySales(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := h.statisticsSvc.GetDailySales(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "获取失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: data})
}

func (h *Handler) ExportStatisticsPDF(c *gin.Context) {
	var req dto.StatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	pdfData, err := h.statisticsSvc.GenerateStatisticsPDF(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Code: 500, Message: "生成PDF失败: " + err.Error()})
		return
	}

	filename := "statistics_report_" + time.Now().Format("20060102150405") + ".pdf"
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: "参数错误"})
		return
	}

	err := h.userSvc.ChangePassword(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Code: 400, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "密码修改成功"})
}
