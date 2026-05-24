package handler

import (
	"strconv"
	"time"

	"music-platform/internal/model"
	"music-platform/internal/service"
	apperrors "music-platform/pkg/errors"
	"music-platform/pkg/jwt"
	"music-platform/pkg/response"
	"music-platform/pkg/utils"

	"github.com/gin-gonic/gin"
)

type RevenueHandler struct {
	revenueService *service.RevenueService
}

func NewRevenueHandler() *RevenueHandler {
	return &RevenueHandler{
		revenueService: service.NewRevenueService(),
	}
}

func (h *RevenueHandler) GetRevenueRecords(c *gin.Context) {
	userID := jwt.GetUserID(c)

	page, pageSize := utils.GetPageAndPageSize(c)

	var startDate, endDate *time.Time
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = &t
		}
	}

	records, total, err := h.revenueService.GetRevenueRecords(userID, page, pageSize, startDate, endDate)
	if err != nil {
		response.InternalError(c, "获取收益记录失败")
		return
	}

	response.Page(c, records, total, page, pageSize)
}

func (h *RevenueHandler) GetArtistRevenueRecords(c *gin.Context) {
	artistID, err := strconv.ParseUint(c.Param("artist_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	page, pageSize := utils.GetPageAndPageSize(c)

	var startDate, endDate *time.Time
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = &t
		}
	}

	records, total, err := h.revenueService.GetArtistRevenueRecords(uint(artistID), page, pageSize, startDate, endDate)
	if err != nil {
		response.InternalError(c, "获取收益记录失败")
		return
	}

	response.Page(c, records, total, page, pageSize)
}

func (h *RevenueHandler) GetTotalRevenue(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var startDate, endDate *time.Time
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = &t
		}
	}

	total, err := h.revenueService.GetTotalRevenue(userID, startDate, endDate)
	if err != nil {
		response.InternalError(c, "获取总收益失败")
		return
	}

	response.Success(c, gin.H{"total_revenue": total})
}

func (h *RevenueHandler) GetRevenueSummary(c *gin.Context) {
	userID := jwt.GetUserID(c)

	summary, err := h.revenueService.GetRevenueSummary(userID)
	if err != nil {
		response.InternalError(c, "获取收益汇总失败")
		return
	}

	response.Success(c, summary)
}

func (h *RevenueHandler) RequestWithdraw(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var req service.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	withdrawRequest, err := h.revenueService.RequestWithdraw(userID, &req)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "申请提现失败")
		return
	}

	response.Success(c, withdrawRequest)
}

func (h *RevenueHandler) GetWithdrawRequests(c *gin.Context) {
	userID := jwt.GetUserID(c)

	page, pageSize := utils.GetPageAndPageSize(c)
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))

	requests, total, err := h.revenueService.GetWithdrawRequests(userID, page, pageSize, status)
	if err != nil {
		response.InternalError(c, "获取提现申请失败")
		return
	}

	response.Page(c, requests, total, page, pageSize)
}

func (h *RevenueHandler) GetAllWithdrawRequests(c *gin.Context) {
	page, pageSize := utils.GetPageAndPageSize(c)
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))

	requests, total, err := h.revenueService.GetAllWithdrawRequests(page, pageSize, status)
	if err != nil {
		response.InternalError(c, "获取提现申请失败")
		return
	}

	response.Page(c, requests, total, page, pageSize)
}

func (h *RevenueHandler) ApproveWithdraw(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	approvedBy := jwt.GetUserID(c)

	err = h.revenueService.ApproveWithdraw(uint(id), approvedBy)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "审批失败")
		return
	}

	response.Success(c, nil)
}

func (h *RevenueHandler) RejectWithdraw(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.revenueService.RejectWithdraw(uint(id), req.Reason)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "拒绝失败")
		return
	}

	response.Success(c, nil)
}

func (h *RevenueHandler) MarkWithdrawPaid(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		TransactionNo string `json:"transaction_no" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err = h.revenueService.MarkWithdrawPaid(uint(id), req.TransactionNo)
	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, 400, appErr.Code, appErr.Message)
			return
		}
		response.InternalError(c, "标记失败")
		return
	}

	response.Success(c, nil)
}

func (h *RevenueHandler) SettleRevenue(c *gin.Context) {
	var req struct {
		Period string `json:"period" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err := h.revenueService.SettleRevenue(req.Period)
	if err != nil {
		response.InternalError(c, "结算失败")
		return
	}

	response.Success(c, nil)
}

func (h *RevenueHandler) GetSubscriptions(c *gin.Context) {
	userID := jwt.GetUserID(c)

	page, pageSize := utils.GetPageAndPageSize(c)

	subscriptions, total, err := h.revenueService.GetSubscriptions(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取订阅列表失败")
		return
	}

	response.Page(c, subscriptions, total, page, pageSize)
}

func (h *RevenueHandler) GetArtistSubscribers(c *gin.Context) {
	artistID, err := strconv.ParseUint(c.Param("artist_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	page, pageSize := utils.GetPageAndPageSize(c)

	subscribers, total, err := h.revenueService.GetArtistSubscribers(uint(artistID), page, pageSize)
	if err != nil {
		response.InternalError(c, "获取订阅者列表失败")
		return
	}

	response.Page(c, subscribers, total, page, pageSize)
}

func (h *RevenueHandler) GetDailyStats(c *gin.Context) {
	userID := jwt.GetUserID(c)

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		startDate = utils.GetMonthStart()
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		endDate = time.Now()
	}

	stats, err := h.revenueService.GetDailyStats(userID, startDate, endDate)
	if err != nil {
		response.InternalError(c, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

func (h *RevenueHandler) GetArtistDailyStats(c *gin.Context) {
	artistID, err := strconv.ParseUint(c.Param("artist_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		startDate = utils.GetMonthStart()
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		endDate = time.Now()
	}

	stats, err := h.revenueService.GetArtistDailyStats(uint(artistID), startDate, endDate)
	if err != nil {
		response.InternalError(c, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

func (h *RevenueHandler) GetArtistStats(c *gin.Context) {
	artistID, err := strconv.ParseUint(c.Param("artist_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	stats, err := h.revenueService.GetArtistStats(uint(artistID))
	if err != nil {
		response.InternalError(c, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

func (h *RevenueHandler) GetOperationLogs(c *gin.Context) {
	page, pageSize := utils.GetPageAndPageSize(c)

	userID, _ := strconv.ParseUint(c.DefaultQuery("user_id", "0"), 10, 64)
	module := c.Query("module")
	keyword := c.Query("keyword")

	logs, total, err := h.revenueService.GetOperationLogs(page, pageSize, uint(userID), module, keyword)
	if err != nil {
		response.InternalError(c, "获取操作日志失败")
		return
	}

	response.Page(c, logs, total, page, pageSize)
}

func (h *RevenueHandler) ExportRevenueExcel(c *gin.Context) {
	userID := jwt.GetUserID(c)

	var startDate, endDate *time.Time
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = &t
		}
	}

	records, err := h.revenueService.ExportRevenueExcel(userID, startDate, endDate)
	if err != nil {
		response.InternalError(c, "导出失败")
		return
	}

	columns := []utils.ExcelColumn{
		{Field: "ID", Title: "ID", Width: 10},
		{Field: "Type", Title: "类型", Width: 15},
		{Field: "Amount", Title: "金额", Width: 15},
		{Field: "PlayCount", Title: "播放次数", Width: 15},
		{Field: "Rate", Title: "费率", Width: 10},
		{Field: "Status", Title: "状态", Width: 10},
		{Field: "Period", Title: "周期", Width: 15},
		{Field: "CreatedAt", Title: "创建时间", Width: 20},
	}

	_ = utils.ExportExcel(c, "revenue_records.xlsx", "收益记录", columns, records)
}

func (h *RevenueHandler) ExportWithdrawExcel(c *gin.Context) {
	userID := jwt.GetUserID(c)

	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))

	records, err := h.revenueService.ExportWithdrawExcel(userID, status)
	if err != nil {
		response.InternalError(c, "导出失败")
		return
	}

	columns := []utils.ExcelColumn{
		{Field: "ID", Title: "ID", Width: 10},
		{Field: "Amount", Title: "金额", Width: 15},
		{Field: "Fee", Title: "手续费", Width: 15},
		{Field: "ActualAmount", Title: "实际金额", Width: 15},
		{Field: "Method", Title: "提现方式", Width: 15},
		{Field: "Account", Title: "账号", Width: 25},
		{Field: "Status", Title: "状态", Width: 10},
		{Field: "CreatedAt", Title: "创建时间", Width: 20},
	}

	_ = utils.ExportExcel(c, "withdraw_requests.xlsx", "提现记录", columns, records)
}

func (h *RevenueHandler) GetWithdrawStatusList(c *gin.Context) {
	statusList := []map[string]interface{}{
		{"value": model.WithdrawStatusPending, "label": "待审核"},
		{"value": model.WithdrawStatusApproved, "label": "已通过"},
		{"value": model.WithdrawStatusRejected, "label": "已拒绝"},
		{"value": model.WithdrawStatusPaid, "label": "已打款"},
		{"value": model.WithdrawStatusFailed, "label": "失败"},
	}

	response.Success(c, statusList)
}
