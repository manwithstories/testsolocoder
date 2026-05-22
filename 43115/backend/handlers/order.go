package handlers

import (
	"fmt"
	"strconv"
	"time"

	"housekeeping-platform/config"
	"housekeeping-platform/models"
	"housekeeping-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateOrderRequest struct {
	ServiceItemID   uint      `json:"service_item_id" binding:"required"`
	AddressID       uint      `json:"address_id" binding:"required"`
	AppointmentTime time.Time `json:"appointment_time" binding:"required"`
	Duration        int       `json:"duration" binding:"required"`
	Remark          string    `json:"remark"`
}

func CheckTimeConflict(providerID uint, appointmentTime time.Time, duration int) bool {
	endTime := appointmentTime.Add(time.Duration(time.Duration(duration) * time.Minute))

	var count int64
	config.DB.Model(&models.Order{}).Where(
		"provider_id = ? AND status IN ? AND appointment_time < ? AND appointment_time + (duration || ' minutes')::interval > ?",
		providerID,
		[]models.OrderStatus{models.OrderStatusConfirmed, models.OrderStatusInService},
		endTime,
		appointmentTime,
	).Count(&count)

	return count > 0
}

func CreateOrder(c *gin.Context) {
	customerID := c.GetUint("user_id")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if !utils.IsFutureTime(req.AppointmentTime) {
		utils.BadRequest(c, "预约时间必须在未来")
		return
	}

	if !utils.ValidateDuration(req.Duration) {
		utils.BadRequest(c, "服务时长必须在30-480分钟之间")
		return
	}

	var serviceItem models.ServiceItem
	if result := config.DB.First(&serviceItem, req.ServiceItemID); result.Error != nil {
		utils.NotFound(c, "服务不存在")
		return
	}

	if !serviceItem.IsActive {
		utils.BadRequest(c, "该服务已下架")
		return
	}

	if req.Duration < serviceItem.MinDuration || req.Duration > serviceItem.MaxDuration {
		utils.BadRequest(c, fmt.Sprintf("服务时长必须在%d-%d分钟之间", serviceItem.MinDuration, serviceItem.MaxDuration))
		return
	}

	var address models.Address
	if result := config.DB.Where("id = ? AND user_id = ?", req.AddressID, customerID).First(&address); result.Error != nil {
		utils.NotFound(c, "地址不存在")
		return
	}

	cfg := config.Load()
	totalAmount := utils.CalculateTotal(serviceItem.BasePrice, req.Duration)
	platformFee := utils.CalculatePlatformFee(totalAmount, cfg.App.PlatformCommissionRate)
	providerIncome := utils.CalculateProviderIncome(totalAmount, platformFee)

	order := models.Order{
		OrderNo:         utils.GenerateOrderNo(),
		CustomerID:      customerID,
		ServiceItemID:   req.ServiceItemID,
		Status:          models.OrderStatusPending,
		AddressID:       req.AddressID,
		ServiceAddress:  fmt.Sprintf("%s%s%s%s", address.Province, address.City, address.District, address.Address),
		ContactName:     address.ContactName,
		ContactPhone:    address.ContactPhone,
		Longitude:       address.Longitude,
		Latitude:        address.Latitude,
		AppointmentTime: req.AppointmentTime,
		Duration:        req.Duration,
		BasePrice:       serviceItem.BasePrice,
		TotalAmount:     totalAmount,
		PlatformFee:     platformFee,
		ProviderIncome:  providerIncome,
		Remark:          req.Remark,
	}

	tx := config.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "创建订单失败")
		return
	}

	var providers []models.User
	query := tx.Where("role = ? AND provider_status = ? AND is_active = ?",
		models.RoleServiceProvider,
		models.ProviderStatusApproved,
		true,
	)

	if len(serviceItem.ServiceAreas) > 0 {
		query = query.Joins("JOIN service_item_areas ON service_item_areas.service_area_id IN (?)",
			tx.Model(&models.ServiceArea{}).Where("district = ?", address.District).Select("id"))
	}

	query.Order("rating DESC, order_count ASC").Limit(5).Find(&providers)

	expiredAt := time.Now().Add(30 * time.Minute)
	for _, provider := range providers {
		invitation := models.OrderInvitation{
			OrderID:    order.ID,
			ProviderID: provider.ID,
			Status:     models.InvitationStatusPending,
			ExpiredAt:  expiredAt,
		}
		if err := tx.Create(&invitation).Error; err != nil {
			continue
		}
		utils.SendInvitationNotification(provider.ID, order.OrderNo)
	}

	tx.Commit()

	config.DB.Model(&models.ServiceItem{}).Where("id = ?", req.ServiceItemID).UpdateColumn("order_count", gorm.Expr("order_count + ?", 1))

	go utils.LogOperation(customerID, "customer", "order", "create", &order.ID, "order", "创建订单", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, order)
}

func GetOrderList(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	query := config.DB.Model(&models.Order{})

	switch role {
	case "customer":
		query = query.Where("customer_id = ?", userID)
	case "service_provider":
		query = query.Where("provider_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	offset := (page - 1) * pageSize
	query.Preload("Customer").Preload("Provider").Preload("ServiceItem").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&orders)

	utils.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"list":      orders,
	})
}

func GetOrderDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var order models.Order
	query := config.DB.Preload("Customer").Preload("Provider").Preload("ServiceItem").
		Preload("Invitations").Preload("Review")

	if role == "customer" {
		query = query.Where("customer_id = ?", userID)
	} else if role == "service_provider" {
		query = query.Where("provider_id = ?", userID)
	}

	if result := query.First(&order, id); result.Error != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	utils.Success(c, order)
}

func RespondInvitation(c *gin.Context) {
	providerID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Accepted     bool   `json:"accepted"`
		RejectReason string `json:"reject_reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var invitation models.OrderInvitation
	if result := config.DB.Where("id = ? AND provider_id = ?", id, providerID).First(&invitation); result.Error != nil {
		utils.NotFound(c, "邀请不存在")
		return
	}

	if invitation.Status != models.InvitationStatusPending {
		utils.BadRequest(c, "该邀请已处理")
		return
	}

	if time.Now().After(invitation.ExpiredAt) {
		invitation.Status = models.InvitationStatusExpired
		config.DB.Save(&invitation)
		utils.BadRequest(c, "邀请已过期")
		return
	}

	var order models.Order
	if result := config.DB.First(&order, invitation.OrderID); result.Error != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.Status != models.OrderStatusPending {
		utils.BadRequest(c, "订单已被处理")
		return
	}

	tx := config.DB.Begin()

	now := time.Now()
	if req.Accepted {
		invitation.Status = models.InvitationStatusAccepted
		invitation.RespondedAt = &now

		order.ProviderID = providerID
		order.Status = models.OrderStatusConfirmed

		if err := tx.Save(&invitation).Error; err != nil {
			tx.Rollback()
			utils.InternalError(c, "处理失败")
			return
		}

		if err := tx.Save(&order).Error; err != nil {
			tx.Rollback()
			utils.InternalError(c, "处理失败")
			return
		}

		tx.Model(&models.OrderInvitation{}).
			Where("order_id = ? AND id != ? AND status = ?", order.ID, invitation.ID, models.InvitationStatusPending).
			Updates(map[string]interface{}{
				"status":       models.InvitationStatusRejected,
				"responded_at": now,
			})

		utils.SendOrderNotification(order.CustomerID, order.OrderNo, string(models.OrderStatusConfirmed))
	} else {
		invitation.Status = models.InvitationStatusRejected
		invitation.RespondedAt = &now
		invitation.RejectReason = req.RejectReason

		if err := tx.Save(&invitation).Error; err != nil {
			tx.Rollback()
			utils.InternalError(c, "处理失败")
			return
		}
	}

	tx.Commit()

	go utils.LogOperation(providerID, "service_provider", "order", "respond_invitation", &invitation.ID, "order_invitation", "响应预约邀请", c.ClientIP(), c.Request.UserAgent())

	utils.SuccessWithMessage(c, "操作成功", nil)
}

func StartService(c *gin.Context) {
	providerID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		Location  string  `json:"location"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var order models.Order
	if result := config.DB.Where("id = ? AND provider_id = ?", id, providerID).First(&order); result.Error != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.Status != models.OrderStatusConfirmed {
		utils.BadRequest(c, "订单状态不正确")
		return
	}

	now := time.Now()
	order.Status = models.OrderStatusInService
	order.ActualStartTime = &now
	order.StartLocation = req.Location
	order.Longitude = req.Longitude
	order.Latitude = req.Latitude

	config.DB.Save(&order)

	utils.SendOrderNotification(order.CustomerID, order.OrderNo, string(models.OrderStatusInService))

	go utils.LogOperation(providerID, "service_provider", "order", "start_service", &order.ID, "order", "开始服务打卡", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, nil)
}

func CompleteService(c *gin.Context) {
	providerID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Location string `json:"location"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var order models.Order
	if result := config.DB.Where("id = ? AND provider_id = ?", id, providerID).First(&order); result.Error != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.Status != models.OrderStatusInService {
		utils.BadRequest(c, "订单状态不正确")
		return
	}

	now := time.Now()
	order.Status = models.OrderStatusCompleted
	order.ActualEndTime = &now
	order.EndLocation = req.Location
	order.CompletedAt = &now

	if order.ActualStartTime != nil {
		actualDuration := int(now.Sub(*order.ActualStartTime).Minutes())
		if actualDuration > 0 {
			order.Duration = actualDuration
			order.TotalAmount = utils.CalculateTotal(order.BasePrice, actualDuration)
			cfg := config.Load()
			order.PlatformFee = utils.CalculatePlatformFee(order.TotalAmount, cfg.App.PlatformCommissionRate)
			order.ProviderIncome = utils.CalculateProviderIncome(order.TotalAmount, order.PlatformFee)
		}
	}

	tx := config.DB.Begin()

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "操作失败")
		return
	}

	bill := models.Bill{
		ProviderID:    order.ProviderID,
		OrderID:       &order.ID,
		BillType:      models.BillTypeIncome,
		Amount:        order.ProviderIncome,
		Balance:       order.ProviderIncome,
		Status:        models.BillStatusCompleted,
		Description:   fmt.Sprintf("订单%s服务收入", order.OrderNo),
		TransactionNo: utils.GenerateTransactionNo(),
		SettledAt:     &now,
	}

	if err := tx.Create(&bill).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "操作失败")
		return
	}

	tx.Model(&models.User{}).Where("id = ?", order.ProviderID).
		Updates(map[string]interface{}{
			"balance":      gorm.Expr("balance + ?", order.ProviderIncome),
			"total_income": gorm.Expr("total_income + ?", order.ProviderIncome),
			"order_count":  gorm.Expr("order_count + ?", 1),
		})

	tx.Commit()

	utils.SendOrderNotification(order.CustomerID, order.OrderNo, string(models.OrderStatusCompleted))
	utils.SendReviewNotification(order.CustomerID, order.OrderNo)
	utils.SendReviewNotification(order.ProviderID, order.OrderNo)

	go utils.LogOperation(providerID, "service_provider", "order", "complete_service", &order.ID, "order", "完成服务打卡", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, nil)
}

func CancelOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var order models.Order
	query := config.DB.Model(&models.Order{})
	if role == "customer" {
		query = query.Where("customer_id = ?", userID)
	} else if role == "service_provider" {
		query = query.Where("provider_id = ?", userID)
	}

	if result := query.First(&order, id); result.Error != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.Status != models.OrderStatusPending && order.Status != models.OrderStatusConfirmed {
		utils.BadRequest(c, "当前状态无法取消订单")
		return
	}

	cfg := config.Load()
	now := time.Now()
	cancelDeadline := order.AppointmentTime.Add(-time.Duration(cfg.App.MaxCancelTimeMinutes) * time.Minute)

	penaltyAmount := 0.0
	if now.After(cancelDeadline) && order.Status == models.OrderStatusConfirmed {
		penaltyAmount = utils.CalculatePenalty(order.TotalAmount, cfg.App.PenaltyRate)
	}

	tx := config.DB.Begin()

	order.Status = models.OrderStatusCancelled
	order.CancelReason = req.Reason
	order.CancelledBy = role
	order.CancelledAt = &now
	order.PenaltyAmount = penaltyAmount

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "取消失败")
		return
	}

	if penaltyAmount > 0 && order.ProviderID != 0 {
		bill := models.Bill{
			ProviderID:    order.ProviderID,
			OrderID:       &order.ID,
			BillType:      models.BillTypePenalty,
			Amount:        penaltyAmount,
			Balance:       penaltyAmount,
			Status:        models.BillStatusCompleted,
			Description:   fmt.Sprintf("订单%s取消违约金", order.OrderNo),
			TransactionNo: utils.GenerateTransactionNo(),
			SettledAt:     &now,
		}
		tx.Create(&bill)

		tx.Model(&models.User{}).Where("id = ?", order.ProviderID).
			Update("balance", gorm.Expr("balance + ?", penaltyAmount))
	}

	tx.Commit()

	if role == "customer" && order.ProviderID != 0 {
		utils.SendOrderNotification(order.ProviderID, order.OrderNo, string(models.OrderStatusCancelled))
	} else if role == "service_provider" {
		utils.SendOrderNotification(order.CustomerID, order.OrderNo, string(models.OrderStatusCancelled))
	}

	go utils.LogOperation(userID, role, "order", "cancel", &order.ID, "order", "取消订单", c.ClientIP(), c.Request.UserAgent())

	utils.SuccessWithMessage(c, "订单已取消", gin.H{"penalty_amount": penaltyAmount})
}

func GetMyInvitations(c *gin.Context) {
	providerID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	query := config.DB.Model(&models.OrderInvitation{}).Where("provider_id = ?", providerID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var invitations []models.OrderInvitation
	offset := (page - 1) * pageSize
	query.Preload("Provider").Offset(offset).Limit(pageSize).Order("id DESC").Find(&invitations)

	utils.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"list":      invitations,
	})
}
