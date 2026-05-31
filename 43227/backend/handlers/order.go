package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"beehive-platform/database"
	"beehive-platform/models"
	"beehive-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func generateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("ORD%s%d", now.Format("20060102150405"), now.UnixNano()%1000)
}

func (h *OrderHandler) Create(c *gin.Context) {
	buyerID, _ := c.Get("user_id")

	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var product models.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "product not found")
		return
	}

	if product.Status != "on_sale" {
		utils.Fail(c, http.StatusBadRequest, "product is not on sale")
	}

	if product.Stock < req.Quantity {
		utils.Fail(c, http.StatusBadRequest, "insufficient stock")
		return
	}

	totalAmount := product.Price * req.Quantity

	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to begin transaction")
		return
	}

	orderNo := generateOrderNo()
	order := models.Order{
		OrderNo:         orderNo,
		BuyerID:         buyerID.(uint),
		SellerID:        product.UserID,
		ProductID:       req.ProductID,
		Quantity:        req.Quantity,
		UnitPrice:       product.Price,
		TotalAmount:     totalAmount,
		Status:          "pending",
		PaymentStatus:   "unpaid",
		ShippingAddress: req.ShippingAddress,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		utils.FailWithError(c, http.StatusInternalServerError, "failed to create order", err)
		return
	}

	if err := tx.Model(&product).Update("stock", gorm.Expr("stock - ?", req.Quantity)).Error; err != nil {
		tx.Rollback()
		utils.FailWithError(c, http.StatusInternalServerError, "failed to update product stock", err)
		return
	}

	if product.Stock-req.Quantity <= 0 {
		tx.Model(&product).Update("status", "sold_out")
	}

	if err := tx.Model(&models.Inventory{}).Where("id = ?", product.InventoryID).
		Update("quantity", gorm.Expr("quantity - ?", req.Quantity)).Error; err != nil {
		tx.Rollback()
		utils.FailWithError(c, http.StatusInternalServerError, "failed to update inventory", err)
		return
	}

	sellerNotification := models.Notification{
		UserID:    product.UserID,
		Type:      "order_status",
		Title:     "新订单通知",
		Content:   fmt.Sprintf("您有一个新订单：%s，数量：%.2f%s，金额：￥%.2f", orderNo, req.Quantity, product.Unit, totalAmount),
		RelatedID: &order.ID,
	}
	tx.Create(&sellerNotification)

	if err := tx.Commit().Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to commit transaction")
		return
	}

	utils.Success(c, order)
}

func (h *OrderHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var pageParams utils.PageParams
	if err := c.ShouldBindQuery(&pageParams); err != nil {
		pageParams = utils.PageParams{Page: 1, PageSize: 10}
	}
	if pageParams.Page < 1 {
		pageParams.Page = 1
	}
	if pageParams.PageSize < 1 {
		pageParams.PageSize = 10
	}

	role, _ := c.Get("user_role")
	query := database.DB.Model(&models.Order{})

	if role == "buyer" {
		query = query.Where("buyer_id = ?", userID)
	} else if role == "beekeeper" {
		query = query.Where("seller_id = ?", userID)
	}

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	query = query.Order(sortBy + " " + sortOrder)

	var total int64
	query.Count(&total)

	var orders []models.Order
	query.Offset(pageParams.GetOffset()).Limit(pageParams.PageSize).
		Preload("Buyer").Preload("Seller").Preload("Product").Find(&orders)

	utils.SuccessWithTotal(c, orders, total)
}

func (h *OrderHandler) Get(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("user_role")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var order models.Order
	query := database.DB.Preload("Buyer").Preload("Seller").Preload("Product").Where("id = ?", id)

	if role == "buyer" {
		query = query.Where("buyer_id = ?", userID)
	} else if role == "beekeeper" {
		query = query.Where("seller_id = ?", userID)
	}

	if err := query.First(&order).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "order not found")
		return
	}

	utils.Success(c, order)
}

func (h *OrderHandler) Pay(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var order models.Order
	if err := database.DB.Where("id = ? AND buyer_id = ?", id, userID).First(&order).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "order not found")
		return
	}

	if order.PaymentStatus != "unpaid" {
		utils.Fail(c, http.StatusBadRequest, "order has already been paid")
		return
	}

	now := time.Now()
	tx := database.DB.Begin()
	tx.Model(&order).Updates(map[string]interface{}{
		"payment_status": "paid",
		"status":         "paid",
		"payment_time":   &now,
	})

	sellerNotification := models.Notification{
		UserID:    order.SellerID,
		Type:      "order_status",
		Title:     "订单支付成功",
		Content:   fmt.Sprintf("订单 %s 已支付成功，请及时发货", order.OrderNo),
		RelatedID: &order.ID,
	}
	tx.Create(&sellerNotification)

	tx.Commit()

	utils.Success(c, nil)
}

func (h *OrderHandler) Ship(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		TrackingNumber string `json:"tracking_number" binding:"required"`
		TrackingStatus string `json:"tracking_status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var order models.Order
	if err := database.DB.Where("id = ? AND seller_id = ?", id, userID).First(&order).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "order not found")
		return
	}

	if order.Status != "paid" {
		utils.Fail(c, http.StatusBadRequest, "order cannot be shipped")
		return
	}

	tx := database.DB.Begin()
	tx.Model(&order).Updates(map[string]interface{}{
		"status":          "shipped",
		"tracking_number": req.TrackingNumber,
		"tracking_status": req.TrackingStatus,
	})

	buyerNotification := models.Notification{
		UserID:    order.BuyerID,
		Type:      "order_status",
		Title:     "订单已发货",
		Content:   fmt.Sprintf("订单 %s 已发货，物流单号：%s", order.OrderNo, req.TrackingNumber),
		RelatedID: &order.ID,
	}
	tx.Create(&buyerNotification)

	tx.Commit()

	utils.Success(c, nil)
}

func (h *OrderHandler) Deliver(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var order models.Order
	if err := database.DB.Where("id = ? AND buyer_id = ?", id, userID).First(&order).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "order not found")
		return
	}

	if order.Status != "shipped" {
		utils.Fail(c, http.StatusBadRequest, "order cannot be confirmed delivered")
		return
	}

	database.DB.Model(&order).Update("status", "delivered")

	utils.Success(c, nil)
}

func (h *OrderHandler) Complete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var order models.Order
	if err := database.DB.Where("id = ? AND buyer_id = ?", id, userID).First(&order).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "order not found")
		return
	}

	if order.Status != "delivered" {
		utils.Fail(c, http.StatusBadRequest, "order cannot be completed")
		return
	}

	database.DB.Model(&order).Update("status", "completed")

	utils.Success(c, nil)
}

func (h *OrderHandler) RateSeller(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req models.RateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var order models.Order
	if err := database.DB.Where("id = ? AND buyer_id = ?", id, userID).First(&order).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "order not found")
		return
	}

	if order.Status != "completed" {
		utils.Fail(c, http.StatusBadRequest, "can only rate completed orders")
		return
	}

	if order.BuyerRating != nil {
		utils.Fail(c, http.StatusBadRequest, "order has already been rated")
		return
	}

	tx := database.DB.Begin()

	tx.Model(&order).Updates(map[string]interface{}{
		"buyer_rating":  req.Rating,
		"buyer_comment": req.Comment,
	})

	var seller models.User
	tx.First(&seller, order.SellerID)

	var ratingCount int64
	var ratingSum float64
	tx.Model(&models.Order{}).
		Where("seller_id = ? AND buyer_rating IS NOT NULL", order.SellerID).
		Count(&ratingCount).
		Select("COALESCE(SUM(buyer_rating), 0)").Scan(&ratingSum)

	if ratingCount > 0 {
		newReputation := ratingSum / float64(ratingCount)
		tx.Model(&seller).Update("reputation", newReputation)
	}

	tx.Commit()

	utils.Success(c, nil)
}

func (h *OrderHandler) RateBuyer(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req models.RateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var order models.Order
	if err := database.DB.Where("id = ? AND seller_id = ?", id, userID).First(&order).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "order not found")
		return
	}

	if order.SellerRating != nil {
		utils.Fail(c, http.StatusBadRequest, "buyer has already been rated")
		return
	}

	database.DB.Model(&order).Updates(map[string]interface{}{
		"seller_rating":  req.Rating,
		"seller_comment": req.Comment,
	})

	utils.Success(c, nil)
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var order models.Order
	if err := database.DB.Where("id = ? AND buyer_id = ?", id, userID).First(&order).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "order not found")
		return
	}

	if order.Status != "pending" && order.Status != "paid" {
		utils.Fail(c, http.StatusBadRequest, "order cannot be cancelled")
		return
	}

	tx := database.DB.Begin()

	tx.Model(&order).Updates(map[string]interface{}{
		"status":         "cancelled",
		"payment_status": "refunded",
	})

	tx.Model(&models.Product{}).Where("id = ?", order.ProductID).
		Update("stock", gorm.Expr("stock + ?", order.Quantity))

	tx.Model(&models.Inventory{}).Where("id = ?", func() uint {
		var product models.Product
		database.DB.First(&product, order.ProductID)
		return product.InventoryID
	}()).Update("quantity", gorm.Expr("quantity + ?", order.Quantity))

	sellerNotification := models.Notification{
		UserID:    order.SellerID,
		Type:      "order_status",
		Title:     "订单已取消",
		Content:   fmt.Sprintf("订单 %s 已被买家取消", order.OrderNo),
		RelatedID: &order.ID,
	}
	tx.Create(&sellerNotification)

	tx.Commit()

	utils.Success(c, nil)
}
