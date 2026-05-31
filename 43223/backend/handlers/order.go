package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"coffee-platform/config"
	"coffee-platform/database"
	"coffee-platform/models"
	"coffee-platform/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderHandler struct {
	cfg *config.Config
}

func NewOrderHandler(cfg *config.Config) *OrderHandler {
	return &OrderHandler{cfg: cfg}
}

func generateOrderNo() string {
	return fmt.Sprintf("CO%s%s", time.Now().Format("20060102150405"), uuid.New().String()[:8])
}

func (h *OrderHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(h.cfg.App.PageSize)))
	status := c.Query("status")
	orderNo := c.Query("order_no")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = h.cfg.App.PageSize
	}

	query := database.DB.Model(&models.Order{}).Preload("Items.Product").Preload("User")

	if userRole != string(models.RoleAdmin) {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if orderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+orderNo+"%")
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders)

	utils.PaginatedResponse(c, orders, total, page, pageSize)
}

func (h *OrderHandler) Get(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var order models.Order
	if err := database.DB.Preload("Items.Product").Preload("User").First(&order, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "订单不存在")
		return
	}

	if order.UserID != userID && userRole != string(models.RoleAdmin) {
		utils.Error(c, http.StatusForbidden, "无权查看此订单")
		return
	}

	utils.Success(c, order)
}

func (h *OrderHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	if len(req.Items) == 0 {
		utils.Error(c, http.StatusBadRequest, "订单商品不能为空")
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var totalAmount float64
		var orderItems []models.OrderItem

		for _, item := range req.Items {
			var product models.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				return fmt.Errorf("商品不存在: ID=%d", item.ProductID)
			}

			if product.Status != models.ProductStatusOnSale {
				return fmt.Errorf("商品[%s]已下架", product.Name)
			}

			if product.Stock < item.Quantity {
				return fmt.Errorf("商品[%s]库存不足", product.Name)
			}

			subtotal := product.Price * float64(item.Quantity)
			totalAmount += subtotal

			orderItems = append(orderItems, models.OrderItem{
				ProductID:   product.ID,
				ProductName: product.Name,
				Price:       product.Price,
				Quantity:    item.Quantity,
				Subtotal:    subtotal,
			})

			if err := tx.Model(&product).UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return fmt.Errorf("更新库存失败")
			}
		}

		order := models.Order{
			OrderNo:       generateOrderNo(),
			UserID:        userID,
			TotalAmount:   totalAmount,
			Status:        models.OrderStatusPending,
			PaymentStatus: models.PaymentStatusPending,
			ReceiverName:  req.ReceiverName,
			ReceiverPhone: req.ReceiverPhone,
			Address:       req.Address,
			Remark:        req.Remark,
		}

		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("创建订单失败")
		}

		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			return fmt.Errorf("创建订单项失败")
		}

		if userID > 0 {
			tx.Where("user_id = ?", userID).Delete(&models.CartItem{})
		}

		order.Items = orderItems
		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "订单创建成功", nil)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "订单不存在")
		return
	}

	validTransitions := map[models.OrderStatus][]models.OrderStatus{
		models.OrderStatusPending:    {models.OrderStatusPaid, models.OrderStatusCancelled},
		models.OrderStatusPaid:       {models.OrderStatusProcessing, models.OrderStatusRefunded},
		models.OrderStatusProcessing: {models.OrderStatusShipped, models.OrderStatusRefunded},
		models.OrderStatusShipped:    {models.OrderStatusDelivered},
	}

	validNext, exists := validTransitions[order.Status]
	if !exists {
		utils.Error(c, http.StatusBadRequest, "当前订单状态不允许修改")
		return
	}

	isValid := false
	for _, s := range validNext {
		if s == req.Status {
			isValid = true
			break
		}
	}

	if !isValid {
		utils.Error(c, http.StatusBadRequest, "不允许的状态流转")
		return
	}

	updates := map[string]interface{}{
		"status": req.Status,
	}

	now := time.Now()
	switch req.Status {
	case models.OrderStatusPaid:
		updates["payment_status"] = models.PaymentStatusSuccess
		updates["paid_at"] = &now
	case models.OrderStatusShipped:
		updates["shipped_at"] = &now
	case models.OrderStatusDelivered:
		updates["delivered_at"] = &now
	case models.OrderStatusCancelled:
		updates["cancelled_at"] = &now
	}

	if err := database.DB.Model(&order).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	if req.Status == models.OrderStatusCancelled || req.Status == models.OrderStatusRefunded {
		var items []models.OrderItem
		database.DB.Where("order_id = ?", order.ID).Find(&items)
		for _, item := range items {
			database.DB.Model(&models.Product{}).Where("id = ?", item.ProductID).
				UpdateColumn("stock", gorm.Expr("stock + ?", item.Quantity))
		}
	}

	utils.SuccessWithMessage(c, "状态更新成功", nil)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "订单不存在")
		return
	}

	if order.UserID != userID {
		utils.Error(c, http.StatusForbidden, "无权取消此订单")
		return
	}

	if order.Status != models.OrderStatusPending {
		utils.Error(c, http.StatusBadRequest, "只有待支付订单可以取消")
		return
	}

	now := time.Now()
	database.DB.Model(&order).Updates(map[string]interface{}{
		"status":       models.OrderStatusCancelled,
		"cancelled_at": &now,
	})

	var items []models.OrderItem
	database.DB.Where("order_id = ?", order.ID).Find(&items)
	for _, item := range items {
		database.DB.Model(&models.Product{}).Where("id = ?", item.ProductID).
			UpdateColumn("stock", gorm.Expr("stock + ?", item.Quantity))
	}

	utils.SuccessWithMessage(c, "订单已取消", nil)
}

func (h *OrderHandler) Pay(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	var order models.Order
	if err := database.DB.First(&order, req.OrderID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "订单不存在")
		return
	}

	if order.UserID != userID {
		utils.Error(c, http.StatusForbidden, "无权支付此订单")
		return
	}

	if order.Status != models.OrderStatusPending {
		utils.Error(c, http.StatusBadRequest, "订单状态不允许支付")
		return
	}

	now := time.Now()
	database.DB.Model(&order).Updates(map[string]interface{}{
		"status":         models.OrderStatusPaid,
		"payment_status": models.PaymentStatusSuccess,
		"paid_at":        &now,
	})

	utils.SuccessWithMessage(c, "支付成功", nil)
}

func (h *OrderHandler) GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cartItems []models.CartItem
	database.DB.Preload("Product.Images").Where("user_id = ?", userID).Find(&cartItems)

	var totalAmount float64
	for _, item := range cartItems {
		if item.Product != nil {
			totalAmount += item.Product.Price * float64(item.Quantity)
		}
	}

	utils.Success(c, gin.H{
		"items":        cartItems,
		"total_amount": totalAmount,
		"total_count":  len(cartItems),
	})
}

func (h *OrderHandler) AddToCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	var product models.Product
	if err := database.DB.Where("id = ? AND status = ?", req.ProductID, models.ProductStatusOnSale).
		First(&product).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "商品不存在或已下架")
		return
	}

	if product.Stock < req.Quantity {
		utils.Error(c, http.StatusBadRequest, "库存不足")
		return
	}

	var cartItem models.CartItem
	if database.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).
		First(&cartItem).Error == nil {
		newQty := cartItem.Quantity + req.Quantity
		if product.Stock < newQty {
			utils.Error(c, http.StatusBadRequest, "库存不足")
			return
		}
		cartItem.Quantity = newQty
		database.DB.Save(&cartItem)
	} else {
		cartItem = models.CartItem{
			UserID:    userID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		database.DB.Create(&cartItem)
	}

	utils.SuccessWithMessage(c, "已添加到购物车", cartItem)
}

func (h *OrderHandler) UpdateCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	var cartItem models.CartItem
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&cartItem).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "购物车项不存在")
		return
	}

	var product models.Product
	database.DB.First(&product, cartItem.ProductID)
	if product.Stock < req.Quantity {
		utils.Error(c, http.StatusBadRequest, "库存不足")
		return
	}

	cartItem.Quantity = req.Quantity
	database.DB.Save(&cartItem)

	utils.SuccessWithMessage(c, "更新成功", cartItem)
}

func (h *OrderHandler) RemoveFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.CartItem{}).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "已从购物车移除", nil)
}

func (h *OrderHandler) ClearCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	database.DB.Where("user_id = ?", userID).Delete(&models.CartItem{})
	utils.SuccessWithMessage(c, "购物车已清空", nil)
}
