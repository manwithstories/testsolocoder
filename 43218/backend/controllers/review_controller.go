package controllers

import (
	"secondhand-platform/models"
	"secondhand-platform/services"
	"secondhand-platform/utils"

	"github.com/gin-gonic/gin"
)

type ReviewController struct {
	reviewService *services.ReviewService
}

func NewReviewController() *ReviewController {
	return &ReviewController{
		reviewService: services.NewReviewService(),
	}
}

type CreateReviewRequest struct {
	OrderID       *uint  `json:"order_id"`
	RepairOrderID *uint  `json:"repair_order_id"`
	RevieweeID    uint   `json:"reviewee_id" binding:"required"`
	ReviewType    string `json:"review_type" binding:"required,oneof=product seller repair technician"`
	Rating        int    `json:"rating" binding:"required,min=1,max=5"`
	Content       string `json:"content"`
	Images        string `json:"images"`
	QualityScore  *int   `json:"quality_score"`
	ServiceScore  *int   `json:"service_score"`
}

func (ctrl *ReviewController) CreateReview(c *gin.Context) {
	reviewerID := c.GetUint("user_id")

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if req.OrderID == nil && req.RepairOrderID == nil {
		utils.Error(c, 400, "必须指定订单ID或维修订单ID")
		return
	}

	if ctrl.reviewService.CheckReviewed(req.OrderID, req.RepairOrderID, reviewerID) {
		utils.Error(c, 400, "您已经评价过此订单")
		return
	}

	review, err := ctrl.reviewService.CreateReview(
		reviewerID, req.RevieweeID, req.OrderID, req.RepairOrderID,
		req.ReviewType, req.Rating, req.Content, req.Images,
		req.QualityScore, req.ServiceScore,
	)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, review)
}

func (ctrl *ReviewController) GetReview(c *gin.Context) {
	id := parseIntParam(c, "id")

	review, err := ctrl.reviewService.GetReviewByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "评价不存在")
		return
	}

	utils.Success(c, review)
}

func (ctrl *ReviewController) ListReviews(c *gin.Context) {
	page := getPage(c)
	pageSize := getPageSize(c)
	revieweeID := parseInt(c.Query("reviewee_id"), 0)
	reviewType := c.Query("review_type")
	minRating := parseInt(c.Query("min_rating"), 0)

	reviews, total, err := ctrl.reviewService.ListReviews(page, pageSize, uint(revieweeID), reviewType, minRating)
	if err != nil {
		utils.Error(c, 500, "获取评价列表失败")
		return
	}

	utils.SuccessWithPagination(c, reviews, page, pageSize, total)
}

func (ctrl *ReviewController) GetAverageRating(c *gin.Context) {
	userID := parseIntParam(c, "user_id")
	reviewType := c.Query("review_type")

	avgRating, err := ctrl.reviewService.GetAverageRating(uint(userID), reviewType)
	if err != nil {
		utils.Error(c, 500, "获取评分失败")
		return
	}

	utils.Success(c, gin.H{
		"user_id":     userID,
		"avg_rating":  avgRating,
		"review_type": reviewType,
	})
}

func (ctrl *ReviewController) DeleteReview(c *gin.Context) {
	reviewerID := c.GetUint("user_id")
	reviewID := parseIntParam(c, "id")

	if err := ctrl.reviewService.DeleteReview(reviewerID, uint(reviewID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

type ReportController struct {
	reportService *services.ReportService
}

func NewReportController() *ReportController {
	return &ReportController{
		reportService: services.NewReportService(),
	}
}

type CreateReportRequest struct {
	TargetType  string `json:"target_type" binding:"required,oneof=product user service review"`
	TargetID    uint   `json:"target_id" binding:"required"`
	Reason      string `json:"reason" binding:"required"`
	Description string `json:"description"`
	Images      string `json:"images"`
}

type HandleReportRequest struct {
	Approved     bool   `json:"approved"`
	HandleResult string `json:"handle_result"`
}

func (ctrl *ReportController) CreateReport(c *gin.Context) {
	reporterID := c.GetUint("user_id")

	var req CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	report, err := ctrl.reportService.CreateReport(
		reporterID, req.TargetType, req.TargetID, req.Reason, req.Description, req.Images,
	)
	if err != nil {
		utils.Error(c, 500, "举报失败")
		return
	}

	utils.Success(c, report)
}

func (ctrl *ReportController) GetReport(c *gin.Context) {
	id := parseIntParam(c, "id")

	report, err := ctrl.reportService.GetReportByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "举报不存在")
		return
	}

	utils.Success(c, report)
}

func (ctrl *ReportController) ListReports(c *gin.Context) {
	page := getPage(c)
	pageSize := getPageSize(c)
	status := parseInt(c.Query("status"), 0)
	targetType := c.Query("target_type")

	reports, total, err := ctrl.reportService.ListReports(page, pageSize, status, targetType)
	if err != nil {
		utils.Error(c, 500, "获取举报列表失败")
		return
	}

	utils.SuccessWithPagination(c, reports, page, pageSize, total)
}

func (ctrl *ReportController) HandleReport(c *gin.Context) {
	reportID := parseIntParam(c, "id")
	handledBy := c.GetUint("user_id")

	var req HandleReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	if err := ctrl.reportService.HandleReport(uint(reportID), req.Approved, req.HandleResult, handledBy); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

type WarrantyController struct {
	warrantyService *services.WarrantyService
}

func NewWarrantyController() *WarrantyController {
	return &WarrantyController{
		warrantyService: services.NewWarrantyService(),
	}
}

type CreateWarrantyRequest struct {
	OrderID       *uint  `json:"order_id"`
	RepairOrderID *uint  `json:"repair_order_id"`
	Type          string `json:"type" binding:"required,oneof=product repair"`
	Description   string `json:"description" binding:"required"`
	Images        string `json:"images"`
}

type HandleWarrantyRequest struct {
	Status       int    `json:"status" binding:"required,oneof=1 2 3 4"`
	HandleResult string `json:"handle_result"`
}

func (ctrl *WarrantyController) CreateWarranty(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateWarrantyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if req.OrderID == nil && req.RepairOrderID == nil {
		utils.Error(c, 400, "必须指定订单ID或维修订单ID")
		return
	}

	warranty, err := ctrl.warrantyService.CreateWarrantyClaim(
		userID, req.OrderID, req.RepairOrderID, req.Type, req.Description, req.Images,
	)
	if err != nil {
		utils.Error(c, 500, "申请质保失败")
		return
	}

	utils.Success(c, warranty)
}

func (ctrl *WarrantyController) GetWarranty(c *gin.Context) {
	id := parseIntParam(c, "id")

	warranty, err := ctrl.warrantyService.GetWarrantyByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "质保申请不存在")
		return
	}

	utils.Success(c, warranty)
}

func (ctrl *WarrantyController) ListWarranties(c *gin.Context) {
	userID := c.GetUint("user_id")
	page := getPage(c)
	pageSize := getPageSize(c)
	status := parseInt(c.Query("status"), 0)

	warranties, total, err := ctrl.warrantyService.ListWarranties(page, pageSize, userID, status)
	if err != nil {
		utils.Error(c, 500, "获取质保列表失败")
		return
	}

	utils.SuccessWithPagination(c, warranties, page, pageSize, total)
}

func (ctrl *WarrantyController) HandleWarranty(c *gin.Context) {
	warrantyID := parseIntParam(c, "id")
	handledBy := c.GetUint("user_id")

	var req HandleWarrantyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	if err := ctrl.warrantyService.HandleWarranty(uint(warrantyID), req.Status, req.HandleResult, handledBy); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

type NotificationController struct {
	notificationService *services.NotificationService
}

func NewNotificationController() *NotificationController {
	return &NotificationController{
		notificationService: services.NewNotificationService(),
	}
}

func (ctrl *NotificationController) ListNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	page := getPage(c)
	pageSize := getPageSize(c)
	isRead := c.Query("is_read")

	var isReadPtr *bool
	if isRead != "" {
		val := isRead == "true"
		isReadPtr = &val
	}

	notifications, total, err := ctrl.notificationService.ListNotifications(userID, page, pageSize, isReadPtr)
	if err != nil {
		utils.Error(c, 500, "获取通知列表失败")
		return
	}

	utils.SuccessWithPagination(c, notifications, page, pageSize, total)
}

func (ctrl *NotificationController) GetNotification(c *gin.Context) {
	id := parseIntParam(c, "id")

	notification, err := ctrl.notificationService.GetNotificationByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "通知不存在")
		return
	}

	utils.Success(c, notification)
}

func (ctrl *NotificationController) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	notificationID := parseIntParam(c, "id")

	if err := ctrl.notificationService.MarkAsRead(userID, uint(notificationID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *NotificationController) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := ctrl.notificationService.MarkAllAsRead(userID); err != nil {
		utils.Error(c, 500, "操作失败")
		return
	}

	utils.Success(c, nil)
}

func (ctrl *NotificationController) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	count, err := ctrl.notificationService.GetUnreadCount(userID)
	if err != nil {
		utils.Error(c, 500, "获取未读数量失败")
		return
	}

	utils.Success(c, gin.H{"unread_count": count})
}

func (ctrl *NotificationController) DeleteNotification(c *gin.Context) {
	userID := c.GetUint("user_id")
	notificationID := parseIntParam(c, "id")

	if err := ctrl.notificationService.DeleteNotification(userID, uint(notificationID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

type MessageController struct {
	messageService *services.MessageService
}

func NewMessageController() *MessageController {
	return &MessageController{
		messageService: services.NewMessageService(),
	}
}

type SendMessageRequest struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Type       string `json:"type" binding:"omitempty,oneof=text image"`
}

func (ctrl *MessageController) SendMessage(c *gin.Context) {
	senderID := c.GetUint("user_id")

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	msgType := req.Type
	if msgType == "" {
		msgType = "text"
	}

	message, err := ctrl.messageService.SendMessage(senderID, req.ReceiverID, req.Content, msgType)
	if err != nil {
		utils.Error(c, 500, "发送消息失败")
		return
	}

	utils.Success(c, message)
}

func (ctrl *MessageController) ListMessages(c *gin.Context) {
	userID := c.GetUint("user_id")
	otherUserID := parseIntParam(c, "user_id")
	page := getPage(c)
	pageSize := getPageSize(c)

	messages, total, err := ctrl.messageService.ListMessages(userID, uint(otherUserID), page, pageSize)
	if err != nil {
		utils.Error(c, 500, "获取消息列表失败")
		return
	}

	utils.SuccessWithPagination(c, messages, page, pageSize, total)
}

func (ctrl *MessageController) MarkMessagesAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	otherUserID := parseIntParam(c, "user_id")

	if err := ctrl.messageService.MarkMessagesAsRead(userID, uint(otherUserID)); err != nil {
		utils.Error(c, 500, "操作失败")
		return
	}

	utils.Success(c, nil)
}

func (ctrl *MessageController) GetUnreadMessageCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	count, err := ctrl.messageService.GetUnreadMessageCount(userID)
	if err != nil {
		utils.Error(c, 500, "获取未读消息数量失败")
		return
	}

	utils.Success(c, gin.H{"unread_count": count})
}

func (ctrl *MessageController) GetMessageContacts(c *gin.Context) {
	userID := c.GetUint("user_id")

	contacts, err := ctrl.messageService.GetMessageContacts(userID)
	if err != nil {
		utils.Error(c, 500, "获取联系人列表失败")
		return
	}

	utils.Success(c, contacts)
}

type AdminController struct {
	adminService *services.AdminService
	userService  *services.UserService
}

func NewAdminController() *AdminController {
	return &AdminController{
		adminService: services.NewAdminService(),
		userService:  services.NewUserService(),
	}
}

func (ctrl *AdminController) GetDashboardStats(c *gin.Context) {
	stats, err := ctrl.adminService.GetDashboardStats()
	if err != nil {
		utils.Error(c, 500, "获取统计数据失败")
		return
	}

	utils.Success(c, stats)
}

func (ctrl *AdminController) GetTransactionStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" {
		startDate = "1970-01-01"
	}
	if endDate == "" {
		endDate = "2099-12-31"
	}

	stats, err := ctrl.adminService.GetTransactionStats(startDate, endDate)
	if err != nil {
		utils.Error(c, 500, "获取交易统计失败")
		return
	}

	utils.Success(c, stats)
}

func (ctrl *AdminController) GetUserActivityStats(c *gin.Context) {
	days := parseInt(c.Query("days"), 30)

	stats, err := ctrl.adminService.GetUserActivityStats(days)
	if err != nil {
		utils.Error(c, 500, "获取用户活跃度统计失败")
		return
	}

	utils.Success(c, stats)
}

func (ctrl *AdminController) ListAllTransactions(c *gin.Context) {
	page := getPage(c)
	pageSize := getPageSize(c)
	transactionType := c.Query("type")
	status := parseInt(c.Query("status"), 0)

	transactions, total, err := ctrl.adminService.ListAllTransactions(page, pageSize, transactionType, status)
	if err != nil {
		utils.Error(c, 500, "获取交易列表失败")
		return
	}

	utils.SuccessWithPagination(c, transactions, page, pageSize, total)
}

func (ctrl *AdminController) ListAdminLogs(c *gin.Context) {
	page := getPage(c)
	pageSize := getPageSize(c)

	logs, total, err := ctrl.adminService.ListAdminLogs(page, pageSize)
	if err != nil {
		utils.Error(c, 500, "获取操作日志失败")
		return
	}

	utils.SuccessWithPagination(c, logs, page, pageSize, total)
}

func init() {
	_ = models.ReportStatusPending
}
