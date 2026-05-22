package handlers

import (
	"fmt"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"e-learning-platform/internal/config"
	"e-learning-platform/internal/database"
	"e-learning-platform/internal/models"
	"e-learning-platform/internal/utils"
)

type OrderHandler struct {
	cfg *config.Config
}

func NewOrderHandler(cfg *config.Config) *OrderHandler {
	return &OrderHandler{cfg: cfg}
}

type CreateOrderRequest struct {
	CourseID uuid.UUID  `json:"course_id" binding:"required"`
	CouponCode string   `json:"coupon_code"`
}

type RefundRequest struct {
	Reason string `json:"reason" binding:"required,max=500"`
}

type CreateCouponRequest struct {
	Code        string     `json:"code" binding:"required,max=50"`
	Type        string     `json:"type" binding:"required,oneof=fixed percent"`
	Value       float64    `json:"value" binding:"required"`
	MinAmount   float64    `json:"min_amount"`
	MaxDiscount float64    `json:"max_discount"`
	TotalCount  int        `json:"total_count"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=paid refunding refunded cancelled failed"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var course models.Course
	if err := database.DB.Where("id = ? AND status = ?", req.CourseID, models.CoursePublished).First(&course).Error; err != nil {
		utils.NotFound(c, "Course not found or not available")
		return
	}

	if course.IsFree {
		utils.BadRequest(c, "This course is free")
		return
	}

	var existingOrder models.Order
	if database.DB.Where("user_id = ? AND course_id = ? AND status IN ?",
		userID, req.CourseID, []models.OrderStatus{models.OrderPaid, models.OrderPending}).
		First(&existingOrder).Error == nil {
		utils.BadRequest(c, "You have already purchased or have a pending order for this course")
		return
	}

	originalPrice := course.Price
	discount := 0.0

	if req.CouponCode != "" {
		var coupon models.Coupon
		if err := database.DB.Where("code = ? AND is_active = true", req.CouponCode).First(&coupon).Error; err != nil {
			utils.BadRequest(c, "Invalid coupon code")
			return
		}

		if coupon.ExpiresAt != nil && time.Now().After(*coupon.ExpiresAt) {
			utils.BadRequest(c, "Coupon has expired")
			return
		}

		if coupon.TotalCount > 0 && coupon.UsedCount >= coupon.TotalCount {
			utils.BadRequest(c, "Coupon has been fully used")
			return
		}

		if originalPrice < coupon.MinAmount {
			utils.BadRequest(c, fmt.Sprintf("Minimum order amount for this coupon is %.2f", coupon.MinAmount))
			return
		}

		var couponUsed models.CouponUsed
		if database.DB.Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).First(&couponUsed).Error == nil {
			utils.BadRequest(c, "You have already used this coupon")
			return
		}

		if coupon.Type == models.CouponFixed {
			discount = math.Min(coupon.Value, originalPrice)
		} else {
			discount = originalPrice * coupon.Value / 100
			if coupon.MaxDiscount > 0 {
				discount = math.Min(discount, coupon.MaxDiscount)
			}
		}
	}

	amount := originalPrice - discount
	if amount < 0 {
		amount = 0
	}

	orderNo := generateOrderNo()

	order := models.Order{
		UserID:        userID.(uuid.UUID),
		CourseID:      req.CourseID,
		OrderNo:       orderNo,
		Amount:        amount,
		OriginalPrice: originalPrice,
		Discount:      discount,
		Status:        models.OrderPending,
	}

	if req.CouponCode != "" {
		var coupon models.Coupon
		database.DB.Where("code = ?", req.CouponCode).First(&coupon)
		order.CouponID = &coupon.ID
	}

	if err := database.DB.Create(&order).Error; err != nil {
		utils.InternalError(c, "Failed to create order")
		return
	}

	utils.Created(c, order)
}

func (h *OrderHandler) PayOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	orderID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid order ID")
		return
	}

	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		utils.NotFound(c, "Order not found")
		return
	}

	if order.UserID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not your order")
		return
	}

	if order.Status != models.OrderPending {
		utils.BadRequest(c, "Order cannot be paid")
		return
	}

	tx := database.DB.Begin()

	now := time.Now()
	order.Status = models.OrderPaid
	order.PaidAt = &now
	order.PayMethod = models.PayMethodBalance

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to update order")
		return
	}

	if err := tx.Model(&models.Course{}).
		Where("id = ?", order.CourseID).
		UpdateColumn("student_count", gorm.Expr("student_count + 1")).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "Failed to update course student count")
		return
	}

	if order.CouponID != nil {
		if err := tx.Model(&models.Coupon{}).
			Where("id = ?", *order.CouponID).
			UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error; err != nil {
			tx.Rollback()
			utils.InternalError(c, "Failed to update coupon usage")
			return
		}

		couponUsed := models.CouponUsed{
			CouponID: *order.CouponID,
			UserID:   userID.(uuid.UUID),
			OrderID:  order.ID,
			UsedAt:   now,
		}
		if err := tx.Create(&couponUsed).Error; err != nil {
			tx.Rollback()
			utils.InternalError(c, "Failed to record coupon usage")
			return
		}
	}

	tx.Commit()
	utils.Success(c, gin.H{"message": "Payment successful", "order": order})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	id := c.Param("id")
	orderID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid order ID")
		return
	}

	var order models.Order
	query := database.DB.Preload("Course").Preload("Course.Instructor")
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&order, orderID).Error; err != nil {
		utils.NotFound(c, "Order not found")
		return
	}

	utils.Success(c, order)
}

func (h *OrderHandler) ListMyOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	status := c.Query("status")

	query := database.DB.Model(&models.Order{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var orders []models.Order
	var total int64

	query.Count(&total)

	page, pageSize := getPagination(c)
	query.Preload("Course").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders)

	utils.Paginated(c, orders, total, page, pageSize)
}

func (h *OrderHandler) ListAllOrders(c *gin.Context) {
	status := c.Query("status")
	search := c.Query("search")

	query := database.DB.Model(&models.Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Joins("JOIN users ON users.id = orders.user_id").
			Where("users.username ILIKE ? OR orders.order_no ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var orders []models.Order
	var total int64

	query.Count(&total)

	page, pageSize := getPagination(c)
	query.Preload("User").Preload("Course").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders)

	utils.Paginated(c, orders, total, page, pageSize)
}

func (h *OrderHandler) ApplyRefund(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")
	orderID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid order ID")
		return
	}

	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		utils.NotFound(c, "Order not found")
		return
	}

	if order.UserID != userID.(uuid.UUID) {
		utils.Forbidden(c, "Not your order")
		return
	}

	if order.Status != models.OrderPaid {
		utils.BadRequest(c, "Only paid orders can be refunded")
		return
	}

	var req RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	order.Status = models.OrderRefunding
	order.RefundReason = req.Reason

	if err := database.DB.Save(&order).Error; err != nil {
		utils.InternalError(c, "Failed to apply refund")
		return
	}

	utils.Success(c, gin.H{"message": "Refund request submitted"})
}

func (h *OrderHandler) ProcessRefund(c *gin.Context) {
	id := c.Param("id")
	orderID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid order ID")
		return
	}

	var req struct {
		Approved bool   `json:"approved"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		utils.NotFound(c, "Order not found")
		return
	}

	if order.Status != models.OrderRefunding {
		utils.BadRequest(c, "Order is not in refunding status")
		return
	}

	tx := database.DB.Begin()

	if req.Approved {
		now := time.Now()
		order.Status = models.OrderRefunded
		order.RefundedAt = &now
		order.Remark = req.Remark

		if err := tx.Save(&order).Error; err != nil {
			tx.Rollback()
			utils.InternalError(c, "Failed to process refund")
			return
		}

		if err := tx.Model(&models.Course{}).
			Where("id = ?", order.CourseID).
			UpdateColumn("student_count", gorm.Expr("student_count - 1")).Error; err != nil {
			tx.Rollback()
			utils.InternalError(c, "Failed to update course student count")
			return
		}
	} else {
		order.Status = models.OrderPaid
		order.Remark = req.Remark
		if err := tx.Save(&order).Error; err != nil {
			tx.Rollback()
			utils.InternalError(c, "Failed to reject refund")
			return
		}
	}

	tx.Commit()
	utils.Success(c, gin.H{"message": "Refund processed"})
}

func (h *OrderHandler) CreateCoupon(c *gin.Context) {
	var req CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	var existing models.Coupon
	if database.DB.Where("code = ?", req.Code).First(&existing).Error == nil {
		utils.BadRequest(c, "Coupon code already exists")
		return
	}

	coupon := models.Coupon{
		Code:        req.Code,
		Type:        models.CouponType(req.Type),
		Value:       req.Value,
		MinAmount:   req.MinAmount,
		MaxDiscount: req.MaxDiscount,
		TotalCount:  req.TotalCount,
		IsActive:    true,
		ExpiresAt:   req.ExpiresAt,
	}

	if err := database.DB.Create(&coupon).Error; err != nil {
		utils.InternalError(c, "Failed to create coupon")
		return
	}

	utils.Created(c, coupon)
}

func (h *OrderHandler) ListCoupons(c *gin.Context) {
	var coupons []models.Coupon
	database.DB.Order("created_at DESC").Find(&coupons)
	utils.Success(c, coupons)
}

func (h *OrderHandler) UpdateCoupon(c *gin.Context) {
	id := c.Param("id")
	couponID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid coupon ID")
		return
	}

	var req struct {
		IsActive *bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := database.DB.Model(&models.Coupon{}).Where("id = ?", couponID).Updates(updates).Error; err != nil {
		utils.InternalError(c, "Failed to update coupon")
		return
	}

	utils.Success(c, gin.H{"message": "Coupon updated"})
}

func (h *OrderHandler) ValidateCoupon(c *gin.Context) {
	code := c.Query("code")
	courseIDStr := c.Query("course_id")

	if code == "" || courseIDStr == "" {
		utils.BadRequest(c, "Code and course_id are required")
		return
	}

	var coupon models.Coupon
	if err := database.DB.Where("code = ? AND is_active = true", code).First(&coupon).Error; err != nil {
		utils.BadRequest(c, "Invalid coupon")
		return
	}

	if coupon.ExpiresAt != nil && time.Now().After(*coupon.ExpiresAt) {
		utils.BadRequest(c, "Coupon has expired")
		return
	}

	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		utils.BadRequest(c, "Invalid course ID")
		return
	}

	var course models.Course
	database.DB.First(&course, courseID)

	discount := 0.0
	if course.Price < coupon.MinAmount {
		utils.BadRequest(c, fmt.Sprintf("Minimum order amount is %.2f", coupon.MinAmount))
		return
	}

	if coupon.Type == models.CouponFixed {
		discount = math.Min(coupon.Value, course.Price)
	} else {
		discount = course.Price * coupon.Value / 100
		if coupon.MaxDiscount > 0 {
			discount = math.Min(discount, coupon.MaxDiscount)
		}
	}

	utils.Success(c, gin.H{
		"coupon":   coupon,
		"discount": discount,
		"final":    course.Price - discount,
	})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	orderID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid order ID")
		return
	}

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := database.DB.Model(&models.Order{}).
		Where("id = ?", orderID).
		Update("status", req.Status).Error; err != nil {
		utils.InternalError(c, "Failed to update order status")
		return
	}

	utils.Success(c, gin.H{"message": "Order status updated"})
}

func (h *OrderHandler) DeleteCoupon(c *gin.Context) {
	id := c.Param("id")
	couponID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "Invalid coupon ID")
		return
	}

	if err := database.DB.Delete(&models.Coupon{}, couponID).Error; err != nil {
		utils.InternalError(c, "Failed to delete coupon")
		return
	}

	utils.Success(c, gin.H{"message": "Coupon deleted"})
}

func generateOrderNo() string {
	return fmt.Sprintf("ORD%s%d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000)
}
