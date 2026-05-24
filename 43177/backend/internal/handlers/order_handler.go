package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"repair-platform/internal/models"
	"repair-platform/internal/utils"
	"repair-platform/pkg/logger"
	"gorm.io/gorm"
)

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

type CreateOrderRequest struct {
	ServiceItemID   uint       `json:"service_item_id" binding:"required"`
	Title           string     `json:"title" binding:"required"`
	Description     string     `json:"description"`
	Images          string     `json:"images"`
	Address         string     `json:"address" binding:"required"`
	Longitude       float64    `json:"longitude"`
	Latitude        float64    `json:"latitude"`
	ContactName     string     `json:"contact_name" binding:"required"`
	ContactPhone    string     `json:"contact_phone" binding:"required"`
	AppointmentTime *time.Time `json:"appointment_time"`
	UrgentLevel     int        `json:"urgent_level"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误: "+err.Error())
		return
	}

	var serviceItem models.ServiceItem
	if err := models.DB.First(&serviceItem, req.ServiceItemID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "服务项目不存在")
		return
	}

	order := models.Order{
		OrderNo:         utils.GenerateOrderNo("ORD"),
		CustomerID:      userID,
		ServiceItemID:   req.ServiceItemID,
		Title:           req.Title,
		Description:     req.Description,
		Images:          req.Images,
		Address:         req.Address,
		Longitude:       req.Longitude,
		Latitude:        req.Latitude,
		ContactName:     req.ContactName,
		ContactPhone:    req.ContactPhone,
		AppointmentTime: req.AppointmentTime,
		UrgentLevel:     req.UrgentLevel,
		Status:          models.OrderStatusPending,
		QuotedPrice:     (serviceItem.MinPrice + serviceItem.MaxPrice) / 2,
	}

	tx := models.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Create order error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "创建工单失败")
		return
	}

	orderLog := models.OrderLog{
		OrderID: order.ID,
		UserID:  userID,
		Action:  "create",
		Content: "客户创建工单",
	}
	if err := tx.Create(&orderLog).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Create order log error: %v", err)
		utils.Error(c, http.StatusInternalServerError, 500, "创建工单失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "创建工单失败")
		return
	}

	go h.autoMatchTechnician(order.ID)

	utils.Success(c, gin.H{
		"order_id": order.ID,
		"order_no": order.OrderNo,
		"message":  "工单创建成功，正在为您匹配技师",
	})
}

func (h *OrderHandler) autoMatchTechnician(orderID uint) {
	var order models.Order
	if err := models.DB.First(&order, orderID).Error; err != nil {
		logger.Errorf("Auto match: order not found %d", orderID)
		return
	}

	var serviceItem models.ServiceItem
	models.DB.First(&serviceItem, order.ServiceItemID)

	var techProfiles []models.TechnicianProfile
	query := models.DB.Preload("User").Where("is_verified = ?", true).
		Where("active_orders < max_active_orders").
		Where("specialty LIKE ?", "%"+serviceItem.Category.Code+"%").
		Order("rating DESC")

	if err := query.Find(&techProfiles).Error; err != nil {
		logger.Errorf("Auto match: find technicians error: %v", err)
		return
	}

	var bestMatch *models.TechnicianProfile
	minDistance := 0.0
	for i := range techProfiles {
		tech := techProfiles[i]
		distance := utils.CalculateDistance(order.Latitude, order.Longitude, tech.User.Latitude, tech.User.Longitude)
		if distance <= tech.ServiceRadius {
			if bestMatch == nil || distance < minDistance {
				bestMatch = &techProfiles[i]
				minDistance = distance
			}
		}
	}

	if bestMatch == nil && len(techProfiles) > 0 {
		bestMatch = &techProfiles[0]
	}

	if bestMatch != nil {
		now := time.Now()
		models.DB.Model(&order).Updates(map[string]interface{}{
			"technician_id": bestMatch.UserID,
			"status":        models.OrderStatusAssigned,
			"assigned_at":   now,
		})

		models.DB.Model(bestMatch).Update("active_orders", gorm.Expr("active_orders + 1"))

		orderLog := models.OrderLog{
			OrderID: order.ID,
			UserID:  bestMatch.UserID,
			Action:  "assign",
			Content: "系统自动匹配技师",
		}
		models.DB.Create(&orderLog)

		logger.Infof("Order %d matched to technician %d", order.ID, bestMatch.UserID)
	}
}

func (h *OrderHandler) GetOrderList(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	query := models.DB.Preload("ServiceItem").Preload("Customer").Preload("Technician")

	switch role {
	case "customer":
		query = query.Where("customer_id = ?", userID)
	case "technician":
		query = query.Where("technician_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.Order{}).Count(&total)

	var orders []models.Order
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders)

	utils.Success(c, gin.H{
		"list":     orders,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	id, _ := strconv.Atoi(c.Param("id"))

	var order models.Order
	if err := models.DB.Preload("ServiceItem").Preload("ServiceItem.Category").
		Preload("Customer").Preload("Technician").
		First(&order, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "工单不存在")
		return
	}

	if role == "customer" && order.CustomerID != userID {
		utils.Error(c, http.StatusForbidden, 403, "无权查看此工单")
		return
	}
	if role == "technician" && order.TechnicianID != nil && *order.TechnicianID != userID {
		utils.Error(c, http.StatusForbidden, 403, "无权查看此工单")
		return
	}

	var logs []models.OrderLog
	models.DB.Where("order_id = ?", order.ID).Order("created_at ASC").Find(&logs)

	var reviews []models.Review
	models.DB.Where("order_id = ?", order.ID).Find(&reviews)

	var partUsages []models.PartUsage
	models.DB.Preload("Part").Where("order_id = ?", order.ID).Find(&partUsages)

	utils.Success(c, gin.H{
		"order":       order,
		"logs":        logs,
		"reviews":     reviews,
		"part_usages": partUsages,
	})
}

type AcceptOrderRequest struct {
	QuotedPrice float64 `json:"quoted_price"`
}

func (h *OrderHandler) AcceptOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var req AcceptOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	var order models.Order
	if err := models.DB.First(&order, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "工单不存在")
		return
	}

	if order.TechnicianID == nil || *order.TechnicianID != userID {
		utils.Error(c, http.StatusForbidden, 403, "无权操作此工单")
		return
	}

	if order.Status != models.OrderStatusAssigned {
		utils.Error(c, http.StatusBadRequest, 400, "工单状态不允许接单")
		return
	}

	var serviceItem models.ServiceItem
	models.DB.First(&serviceItem, order.ServiceItemID)

	quotedPrice := req.QuotedPrice
	if quotedPrice == 0 {
		quotedPrice = order.QuotedPrice
	}

	now := time.Now()
	tx := models.DB.Begin()

	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":       models.OrderStatusAccepted,
		"accepted_at":  now,
		"quoted_price": quotedPrice,
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "接单失败")
		return
	}

	orderLog := models.OrderLog{
		OrderID: order.ID,
		UserID:  userID,
		Action:  "accept",
		Content: "技师接单，报价: " + strconv.FormatFloat(quotedPrice, 'f', 2, 64),
	}
	if err := tx.Create(&orderLog).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "接单失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "接单失败")
		return
	}

	utils.Success(c, gin.H{"message": "接单成功"})
}

func (h *OrderHandler) ArriveAtSite(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var order models.Order
	if err := models.DB.First(&order, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "工单不存在")
		return
	}

	if order.TechnicianID == nil || *order.TechnicianID != userID {
		utils.Error(c, http.StatusForbidden, 403, "无权操作此工单")
		return
	}

	if order.Status != models.OrderStatusAccepted {
		utils.Error(c, http.StatusBadRequest, 400, "工单状态不允许上门")
		return
	}

	now := time.Now()
	tx := models.DB.Begin()

	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":    models.OrderStatusOnSite,
		"arrived_at": now,
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	orderLog := models.OrderLog{
		OrderID:   order.ID,
		UserID:   userID,
		Action:   "arrive",
		Content:  "技师已到达现场",
		Longitude: order.Longitude,
		Latitude:  order.Latitude,
	}
	if err := tx.Create(&orderLog).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	utils.Success(c, gin.H{"message": "已标记到达现场"})
}

func (h *OrderHandler) StartRepair(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var order models.Order
	if err := models.DB.First(&order, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "工单不存在")
		return
	}

	if order.TechnicianID == nil || *order.TechnicianID != userID {
		utils.Error(c, http.StatusForbidden, 403, "无权操作此工单")
		return
	}

	if order.Status != models.OrderStatusOnSite {
		utils.Error(c, http.StatusBadRequest, 400, "工单状态不允许开始维修")
		return
	}

	now := time.Now()
	tx := models.DB.Begin()

	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":    models.OrderStatusRepairing,
		"started_at": now,
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	orderLog := models.OrderLog{
		OrderID: order.ID,
		UserID:  userID,
		Action:  "start_repair",
		Content: "开始维修",
	}
	if err := tx.Create(&orderLog).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	utils.Success(c, gin.H{"message": "已开始维修"})
}

type CompleteOrderRequest struct {
	FinalPrice float64 `json:"final_price"`
	Note       string  `json:"note"`
}

func (h *OrderHandler) CompleteOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var req CompleteOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	var order models.Order
	if err := models.DB.First(&order, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "工单不存在")
		return
	}

	if order.TechnicianID == nil || *order.TechnicianID != userID {
		utils.Error(c, http.StatusForbidden, 403, "无权操作此工单")
		return
	}

	if order.Status != models.OrderStatusRepairing {
		utils.Error(c, http.StatusBadRequest, 400, "工单状态不允许完工")
		return
	}

	now := time.Now()
	tx := models.DB.Begin()

	finalPrice := req.FinalPrice
	if finalPrice == 0 {
		finalPrice = order.QuotedPrice
	}

	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":       models.OrderStatusCompleted,
		"completed_at": now,
		"final_price":  finalPrice,
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	orderLog := models.OrderLog{
		OrderID: order.ID,
		UserID:  userID,
		Action:  "complete",
		Content: "维修完成，最终费用: " + strconv.FormatFloat(finalPrice, 'f', 2, 64),
	}
	if err := tx.Create(&orderLog).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	var profile models.TechnicianProfile
	models.DB.Where("user_id = ?", *order.TechnicianID).First(&profile)
	models.DB.Model(&profile).Updates(map[string]interface{}{
		"completed_orders": gorm.Expr("completed_orders + 1"),
		"active_orders":    gorm.Expr("active_orders - 1"),
	})

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	utils.Success(c, gin.H{"message": "工单已完工"})
}

type CancelOrderRequest struct {
	CancelReason string `json:"cancel_reason" binding:"required"`
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var req CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	var order models.Order
	if err := models.DB.First(&order, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "工单不存在")
		return
	}

	if order.CustomerID != userID {
		utils.Error(c, http.StatusForbidden, 403, "无权取消此工单")
		return
	}

	if order.Status != models.OrderStatusPending && order.Status != models.OrderStatusAssigned {
		utils.Error(c, http.StatusBadRequest, 400, "当前状态不允许取消")
		return
	}

	now := time.Now()
	tx := models.DB.Begin()

	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":        models.OrderStatusCancelled,
		"cancelled_at":  now,
		"cancel_reason": req.CancelReason,
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "取消失败")
		return
	}

	if order.TechnicianID != nil {
		var profile models.TechnicianProfile
		models.DB.Where("user_id = ?", *order.TechnicianID).First(&profile)
		models.DB.Model(&profile).Update("active_orders", gorm.Expr("active_orders - 1"))
	}

	orderLog := models.OrderLog{
		OrderID: order.ID,
		UserID:  userID,
		Action:  "cancel",
		Content: "客户取消工单，原因: " + req.CancelReason,
	}
	if err := tx.Create(&orderLog).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "取消失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "取消失败")
		return
	}

	utils.Success(c, gin.H{"message": "工单已取消"})
}

type RefundOrderRequest struct {
	RefundReason string `json:"refund_reason" binding:"required"`
	RefundAmount float64 `json:"refund_amount"`
}

func (h *OrderHandler) RequestRefund(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var req RefundOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	var order models.Order
	if err := models.DB.First(&order, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "工单不存在")
		return
	}

	if order.CustomerID != userID {
		utils.Error(c, http.StatusForbidden, 403, "无权申请退款")
		return
	}

	if order.Status != models.OrderStatusCompleted {
		utils.Error(c, http.StatusBadRequest, 400, "当前状态不允许申请退款")
		return
	}

	refundAmount := req.RefundAmount
	if refundAmount <= 0 {
		refundAmount = order.FinalPrice
	}

	if refundAmount > order.FinalPrice {
		refundAmount = order.FinalPrice
	}

	tx := models.DB.Begin()

	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":        models.OrderStatusRefunding,
		"refund_reason": req.RefundReason,
		"refund_amount": refundAmount,
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "退款申请失败")
		return
	}

	orderLog := models.OrderLog{
		OrderID: order.ID,
		UserID:  userID,
		Action:  "refund_request",
		Content: "客户申请退款，原因: " + req.RefundReason + "，金额: " + strconv.FormatFloat(refundAmount, 'f', 2, 64),
	}
	if err := tx.Create(&orderLog).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "退款申请失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "退款申请失败")
		return
	}

	utils.Success(c, gin.H{"message": "退款申请已提交，等待审核"})
}

func (h *OrderHandler) CheckOrderTimeout() {
	logger.Info("Checking order timeout...")

	var pendingOrders []models.Order
	models.DB.Where("status IN ? AND created_at < ?",
		[]models.OrderStatus{models.OrderStatusPending, models.OrderStatusAssigned},
		time.Now().Add(-24*time.Hour)).Find(&pendingOrders)

	for _, order := range pendingOrders {
		logger.Infof("Order %d timed out, cancelling", order.ID)

		models.DB.Model(&order).Updates(map[string]interface{}{
			"status":        models.OrderStatusCancelled,
			"cancelled_at":  time.Now(),
			"cancel_reason": "工单超时自动取消",
		})

		if order.TechnicianID != nil {
			var profile models.TechnicianProfile
			models.DB.Where("user_id = ?", *order.TechnicianID).First(&profile)
			models.DB.Model(&profile).Update("active_orders", gorm.Expr("active_orders - 1"))
		}
	}

	logger.Infof("Processed %d timeout orders", len(pendingOrders))
}

func (h *OrderHandler) ApproveRefund(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var order models.Order
	if err := models.DB.First(&order, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "工单不存在")
		return
	}

	if order.Status != models.OrderStatusRefunding {
		utils.Error(c, http.StatusBadRequest, 400, "工单状态不允许审核退款")
		return
	}

	tx := models.DB.Begin()

	now := time.Now()
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":     models.OrderStatusRefunded,
		"refunded_at": now,
	}).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "审核失败")
		return
	}

	var customer models.User
	models.DB.First(&customer, order.CustomerID)
	tx.Model(&customer).Update("balance", gorm.Expr("balance + ?", order.RefundAmount))

	transaction := models.Transaction{
		TransactionNo: utils.GenerateOrderNo("TXN"),
		UserID:        order.CustomerID,
		Type:          "refund",
		Amount:        order.RefundAmount,
		BalanceAfter:  customer.Balance + order.RefundAmount,
		OrderID:       &order.ID,
		Description:   "工单退款: " + order.OrderNo,
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "审核失败")
		return
	}

	orderLog := models.OrderLog{
		OrderID: order.ID,
		UserID:  c.GetUint("user_id"),
		Action:  "refund_approved",
		Content: "退款审核通过，退款金额: " + strconv.FormatFloat(order.RefundAmount, 'f', 2, 64),
	}
	if err := tx.Create(&orderLog).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "审核失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "审核失败")
		return
	}

	utils.Success(c, gin.H{"message": "退款审核通过"})
}

func (h *OrderHandler) RejectRefund(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	type RejectRequest struct {
		Reason string `json:"reason" binding:"required"`
	}

	var req RejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "请求参数错误")
		return
	}

	var order models.Order
	if err := models.DB.First(&order, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, 404, "工单不存在")
		return
	}

	if order.Status != models.OrderStatusRefunding {
		utils.Error(c, http.StatusBadRequest, 400, "工单状态不允许审核退款")
		return
	}

	tx := models.DB.Begin()

	if err := tx.Model(&order).Update("status", models.OrderStatusCompleted).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	orderLog := models.OrderLog{
		OrderID: order.ID,
		UserID:  c.GetUint("user_id"),
		Action:  "refund_rejected",
		Content: "退款申请被拒绝，原因: " + req.Reason,
	}
	if err := tx.Create(&orderLog).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "操作失败")
		return
	}

	utils.Success(c, gin.H{"message": "退款申请已拒绝"})
}

func (h *OrderHandler) GetOrdersByStatus(c *gin.Context) {
	status := c.Param("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var orders []models.Order
	query := models.DB.Preload("ServiceItem").Preload("Customer").Preload("Technician").
		Where("status = ?", status)

	var total int64
	query.Model(&models.Order{}).Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders)

	utils.Success(c, gin.H{
		"list":     orders,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

func CheckOrderExists(orderID uint) error {
	var order models.Order
	if err := models.DB.First(&order, orderID).Error; err != nil {
		return errors.New("工单不存在")
	}
	return nil
}
