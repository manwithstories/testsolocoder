package handlers

import (
	"encoding/json"
	"fmt"
	"multishop/internal/database"
	"multishop/internal/dto"
	"multishop/internal/middleware"
	"multishop/internal/models"
	"multishop/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) AddToCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req dto.CartAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var product models.Product
	if err := database.DB.Where("id = ? AND status = ?", req.ProductID, models.ProductStatusOnSale).First(&product).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "商品不存在或已下架")
		return
	}

	var existingCart models.CartItem
	query := database.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID)
	if req.SKU_ID != nil {
		query = query.Where("sku_id = ?", *req.SKU_ID)
	}

	if err := query.First(&existingCart).Error; err == nil {
		existingCart.Quantity += req.Quantity
		if err := database.DB.Save(&existingCart).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "添加失败")
			return
		}
	} else {
		cartItem := models.CartItem{
			UserID:    userID,
			ProductID: req.ProductID,
			SKU_ID:    req.SKU_ID,
			Quantity:  req.Quantity,
		}
		if err := database.DB.Create(&cartItem).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "添加失败")
			return
		}
	}

	utils.Success(c, nil)
}

func (h *OrderHandler) GetCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var cartItems []models.CartItem
	database.DB.Where("user_id = ?", userID).Preload("Product").Find(&cartItems)

	result := make([]dto.CartItemInfo, 0, len(cartItems))
	for _, item := range cartItems {
		cartInfo := dto.CartItemInfo{
			ID:        item.ID,
			ProductID: item.ProductID,
			SKU_ID:    item.SKU_ID,
			Quantity:  item.Quantity,
			Product: dto.ProductInfo{
				ID:        item.Product.ID,
				Name:      item.Product.Name,
				MainImage: item.Product.MainImage,
				Price:     item.Product.Price,
				Stock:     item.Product.Stock,
				ShopID:    item.Product.ShopID,
			},
		}

		if item.SKU_ID != nil {
			var sku models.SKU
			database.DB.First(&sku, *item.SKU_ID)
			var specMap map[string]string
			json.Unmarshal([]byte(sku.Specs), &specMap)
			cartInfo.SKU = &dto.SKUInfo{
				ID:      sku.ID,
				Specs:   specMap,
				Price:   sku.Price,
				Stock:   sku.Stock,
				SKUCode: sku.SKUCode,
			}
		}

		result = append(result, cartInfo)
	}

	utils.Success(c, result)
}

func (h *OrderHandler) UpdateCart(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req dto.CartUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var cartItem models.CartItem
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&cartItem).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "购物车项不存在")
		return
	}

	cartItem.Quantity = req.Quantity
	database.DB.Save(&cartItem)
	utils.Success(c, nil)
}

func (h *OrderHandler) RemoveCartItem(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.CartItem{})
	utils.Success(c, nil)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req dto.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var cartItems []models.CartItem
	if len(req.CartIDs) > 0 {
		database.DB.Where("id IN ? AND user_id = ?", req.CartIDs, userID).Preload("Product").Find(&cartItems)
	} else {
		database.DB.Where("user_id = ?", userID).Preload("Product").Find(&cartItems)
	}

	if len(cartItems) == 0 {
		utils.Error(c, http.StatusBadRequest, "购物车为空")
		return
	}

	shopCartMap := make(map[uint][]models.CartItem)
	for _, item := range cartItems {
		shopCartMap[item.Product.ShopID] = append(shopCartMap[item.Product.ShopID], item)
	}

	tx := database.DB.Begin()
	var orderNos []string
	var totalAmount float64

	for shopID, items := range shopCartMap {
		shopOrderNo := fmt.Sprintf("ORD%d%s", shopID, time.Now().Format("20060102150405"))
		orderNos = append(orderNos, shopOrderNo)

		var orderTotal float64
		orderItems := make([]models.OrderItem, 0, len(items))

		for _, item := range items {
			var price float64
			var stock int

			if item.SKU_ID != nil {
				var sku models.SKU
				if err := tx.First(&sku, *item.SKU_ID).Error; err != nil {
					tx.Rollback()
					utils.Error(c, http.StatusBadRequest, "SKU不存在")
					return
				}

				result := tx.Model(&sku).Where("id = ? AND stock >= ?", sku.ID, item.Quantity).
					Update("stock", gorm.Expr("stock - ?", item.Quantity))
				if result.RowsAffected == 0 {
					tx.Rollback()
					utils.Error(c, http.StatusBadRequest, "库存不足: "+item.Product.Name)
					return
				}
				price = sku.Price
				stock = sku.Stock
			} else {
				result := tx.Model(&item.Product).Where("id = ? AND stock >= ?", item.Product.ID, item.Quantity).
					Update("stock", gorm.Expr("stock - ?", item.Quantity))
				if result.RowsAffected == 0 {
					tx.Rollback()
					utils.Error(c, http.StatusBadRequest, "库存不足: "+item.Product.Name)
					return
				}
				price = item.Product.Price
				stock = item.Product.Stock
			}

			subtotal := price * float64(item.Quantity)
			orderTotal += subtotal

			var specsJSON string
			if item.SKU_ID != nil {
				var sku models.SKU
				tx.First(&sku, *item.SKU_ID)
				specsJSON = sku.Specs
			}

			orderItem := models.OrderItem{
				ProductID:    item.ProductID,
				SKU_ID:       item.SKU_ID,
				ProductName:  item.Product.Name,
				ProductImage: item.Product.MainImage,
				Specs:        specsJSON,
				Price:        price,
				Quantity:     item.Quantity,
				Subtotal:     subtotal,
			}
			orderItems = append(orderItems, orderItem)
			_ = stock
		}

		order := models.Order{
			OrderNo:         shopOrderNo,
			UserID:          userID,
			ShopID:          shopID,
			TotalAmount:     orderTotal,
			Status:          models.OrderStatusPendingPayment,
			ReceiverName:    req.ReceiverName,
			ReceiverPhone:   req.ReceiverPhone,
			ReceiverAddress: req.ReceiverAddress,
			Remark:          req.Remark,
		}

		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			utils.Error(c, http.StatusInternalServerError, "创建订单失败")
			return
		}

		for i := range orderItems {
			orderItems[i].OrderID = order.ID
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			tx.Rollback()
			utils.Error(c, http.StatusInternalServerError, "创建订单失败")
			return
		}

		CreateNotification(tx, shopID, models.NotificationTypeOrder,
			"新订单提醒",
			fmt.Sprintf("您有新的订单，请及时处理: %s", shopOrderNo),
			nil)

		totalAmount += orderTotal
	}

	for _, item := range cartItems {
		tx.Delete(&item)
	}

	tx.Commit()

	utils.Success(c, dto.OrderCreateResponse{
		OrderNos:    orderNos,
		TotalAmount: totalAmount,
	})
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	query := database.DB.Model(&models.Order{})

	if userRole == models.RoleSeller {
		var shop models.Shop
		database.DB.Where("user_id = ?", userID).First(&shop)
		query = query.Where("shop_id = ?", shop.ID)
	} else {
		query = query.Where("user_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	offset := (page - 1) * pageSize
	query.Preload("Shop").Preload("Items").
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders)

	result := make([]dto.OrderInfo, 0, len(orders))
	for _, order := range orders {
		result = append(result, convertToOrderInfo(order))
	}

	utils.Paginated(c, result, total, page, pageSize)
}

func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var order models.Order
	query := database.DB.Preload("Shop").Preload("Items").Preload("Payment")

	if userRole == models.RoleSeller {
		var shop models.Shop
		database.DB.Where("user_id = ?", userID).First(&shop)
		query = query.Where("shop_id = ?", shop.ID)
	} else if userRole == models.RoleBuyer {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&order, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "订单不存在")
		return
	}

	items := make([]dto.OrderItemInfo, 0, len(order.Items))
	for _, item := range order.Items {
		var specMap map[string]string
		if item.Specs != "" {
			json.Unmarshal([]byte(item.Specs), &specMap)
		}
		items = append(items, dto.OrderItemInfo{
			ID:            item.ID,
			ProductID:     item.ProductID,
			SKU_ID:        item.SKU_ID,
			ProductName:   item.ProductName,
			ProductImage:  item.ProductImage,
			Specs:         specMap,
			Price:         item.Price,
			Quantity:      item.Quantity,
			Subtotal:      item.Subtotal,
			Reviewed:      item.Reviewed,
		})
	}

	utils.Success(c, dto.OrderDetail{
		OrderInfo: convertToOrderInfo(order),
		Items:     items,
	})
}

func (h *OrderHandler) PayOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ? AND status = ?", id, userID, models.OrderStatusPendingPayment).First(&order).Error; err != nil {
		utils.Error(c, http.StatusBadRequest, "订单不存在或状态不正确")
		return
	}

	tx := database.DB.Begin()

	now := time.Now()
	order.Status = models.OrderStatusPendingShip
	order.PaidAt = &now
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, "支付失败")
		return
	}

	paymentNo := fmt.Sprintf("PAY%s", time.Now().Format("20060102150405"))
	payment := models.Payment{
		OrderID:   order.ID,
		PaymentNo: paymentNo,
		Amount:    order.TotalAmount,
		Method:    "alipay",
		Status:    models.PaymentStatusPaid,
		PaidAt:    &now,
	}
	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, "支付失败")
		return
	}

	CreateNotification(tx, order.ShopID, models.NotificationTypeOrder,
		"订单已支付",
		fmt.Sprintf("订单 %s 已支付，请尽快发货", order.OrderNo),
		nil)

	tx.Commit()
	utils.Success(c, nil)
}

func (h *OrderHandler) ShipOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req dto.OrderShipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var shop models.Shop
	if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "店铺不存在")
		return
	}

	var order models.Order
	if err := database.DB.Where("id = ? AND shop_id = ? AND status = ?", id, shop.ID, models.OrderStatusPendingShip).First(&order).Error; err != nil {
		utils.Error(c, http.StatusBadRequest, "订单不存在或状态不正确")
		return
	}

	now := time.Now()
	order.Status = models.OrderStatusShipped
	order.ShippedAt = &now
	order.TrackingNo = req.TrackingNo
	order.TrackingCompany = req.TrackingCompany
	database.DB.Save(&order)

	CreateNotification(database.DB, order.UserID, models.NotificationTypeShipment,
		"订单已发货",
		fmt.Sprintf("订单 %s 已通过 %s 发货，物流单号: %s", order.OrderNo, req.TrackingCompany, req.TrackingNo),
		nil)

	utils.Success(c, nil)
}

func (h *OrderHandler) ConfirmReceive(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ? AND status = ?", id, userID, models.OrderStatusShipped).First(&order).Error; err != nil {
		utils.Error(c, http.StatusBadRequest, "订单不存在或状态不正确")
		return
	}

	now := time.Now()
	order.Status = models.OrderStatusCompleted
	order.CompletedAt = &now
	database.DB.Save(&order)

	var orderItems []models.OrderItem
	database.DB.Where("order_id = ?", order.ID).Find(&orderItems)
	for _, item := range orderItems {
		database.DB.Model(&models.Product{}).Where("id = ?", item.ProductID).
			UpdateColumn("sales", gorm.Expr("sales + ?", item.Quantity))
	}

	utils.Success(c, nil)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "订单不存在")
		return
	}

	if order.Status != models.OrderStatusPendingPayment {
		utils.Error(c, http.StatusBadRequest, "当前状态无法取消")
		return
	}

	tx := database.DB.Begin()

	order.Status = models.OrderStatusCancelled
	tx.Save(&order)

	var orderItems []models.OrderItem
	tx.Where("order_id = ?", order.ID).Find(&orderItems)
	for _, item := range orderItems {
		if item.SKU_ID != nil {
			tx.Model(&models.SKU{}).Where("id = ?", *item.SKU_ID).
				UpdateColumn("stock", gorm.Expr("stock + ?", item.Quantity))
		} else {
			tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
				UpdateColumn("stock", gorm.Expr("stock + ?", item.Quantity))
		}
	}

	tx.Commit()
	utils.Success(c, nil)
}

func (h *OrderHandler) ApplyRefund(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req dto.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", req.OrderID, userID).First(&order).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "订单不存在")
		return
	}

	if order.Status != models.OrderStatusShipped && order.Status != models.OrderStatusCompleted {
		utils.Error(c, http.StatusBadRequest, "当前状态无法申请退款")
		return
	}

	var refundAmount float64
	if req.OrderItemID != nil {
		var orderItem models.OrderItem
		if err := database.DB.Where("id = ? AND order_id = ?", *req.OrderItemID, order.ID).First(&orderItem).Error; err != nil {
			utils.Error(c, http.StatusNotFound, "订单项不存在")
			return
		}
		refundAmount = orderItem.Subtotal
	} else {
		refundAmount = order.TotalAmount
	}

	refundNo := fmt.Sprintf("REF%s", time.Now().Format("20060102150405"))
	refund := models.Refund{
		OrderID:     order.ID,
		OrderItemID: req.OrderItemID,
		UserID:      userID,
		ShopID:      order.ShopID,
		RefundNo:    refundNo,
		Amount:      refundAmount,
		Reason:      req.Reason,
		Status:      models.RefundStatusPending,
		Type:        req.Type,
	}

	if err := database.DB.Create(&refund).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "申请失败")
		return
	}

	utils.Success(c, gin.H{"refund_no": refundNo})
}

func (h *OrderHandler) ReviewRefund(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req dto.RefundReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var shop models.Shop
	database.DB.Where("user_id = ?", userID).First(&shop)

	var refund models.Refund
	if err := database.DB.Where("id = ? AND shop_id = ? AND status = ?", id, shop.ID, models.RefundStatusPending).First(&refund).Error; err != nil {
		utils.Error(c, http.StatusBadRequest, "退款申请不存在或状态不正确")
		return
	}

	tx := database.DB.Begin()

	if req.Status == models.RefundStatusApproved {
		now := time.Now()
		refund.Status = models.RefundStatusApproved
		refund.ApprovedAt = &now

		var order models.Order
		tx.First(&order, refund.OrderID)
		order.Status = models.OrderStatusRefunded
		tx.Save(&order)

		if refund.OrderItemID != nil {
			var orderItem models.OrderItem
			tx.First(&orderItem, *refund.OrderItemID)
			if orderItem.SKU_ID != nil {
				tx.Model(&models.SKU{}).Where("id = ?", *orderItem.SKU_ID).
					UpdateColumn("stock", gorm.Expr("stock + ?", orderItem.Quantity))
			} else {
				tx.Model(&models.Product{}).Where("id = ?", orderItem.ProductID).
					UpdateColumn("stock", gorm.Expr("stock + ?", orderItem.Quantity))
			}
		} else {
			var orderItems []models.OrderItem
			tx.Where("order_id = ?", refund.OrderID).Find(&orderItems)
			for _, item := range orderItems {
				if item.SKU_ID != nil {
					tx.Model(&models.SKU{}).Where("id = ?", *item.SKU_ID).
						UpdateColumn("stock", gorm.Expr("stock + ?", item.Quantity))
				} else {
					tx.Model(&models.Product{}).Where("id = ?", item.ProductID).
						UpdateColumn("stock", gorm.Expr("stock + ?", item.Quantity))
				}
			}
		}

		refund.Status = models.RefundStatusCompleted
		refund.CompletedAt = &now

		CreateNotification(tx, refund.UserID, models.NotificationTypeRefund,
			"退款成功",
			fmt.Sprintf("退款申请 %s 已通过，退款金额: %.2f元", refund.RefundNo, refund.Amount),
			nil)
	} else {
		refund.Status = models.RefundStatusRejected
		refund.RejectReason = req.RejectReason

		CreateNotification(tx, refund.UserID, models.NotificationTypeRefund,
			"退款被拒绝",
			fmt.Sprintf("退款申请 %s 被拒绝: %s", refund.RefundNo, req.RejectReason),
			nil)
	}

	tx.Save(&refund)
	tx.Commit()

	utils.Success(c, nil)
}

func convertToOrderInfo(order models.Order) dto.OrderInfo {
	statusText := map[string]string{
		models.OrderStatusPendingPayment: "待付款",
		models.OrderStatusPendingShip:    "待发货",
		models.OrderStatusShipped:        "已发货",
		models.OrderStatusCompleted:      "已完成",
		models.OrderStatusRefunded:       "已退款",
		models.OrderStatusCancelled:      "已取消",
	}[order.Status]

	items := make([]dto.OrderItemInfo, 0, len(order.Items))
	for _, item := range order.Items {
		var specMap map[string]string
		if item.Specs != "" {
			json.Unmarshal([]byte(item.Specs), &specMap)
		}
		items = append(items, dto.OrderItemInfo{
			ID:            item.ID,
			ProductID:     item.ProductID,
			SKU_ID:        item.SKU_ID,
			ProductName:   item.ProductName,
			ProductImage:  item.ProductImage,
			Specs:         specMap,
			Price:         item.Price,
			Quantity:      item.Quantity,
			Subtotal:      item.Subtotal,
			Reviewed:      item.Reviewed,
		})
	}

	info := dto.OrderInfo{
		ID:              order.ID,
		OrderNo:         order.OrderNo,
		ShopID:          order.ShopID,
		ShopName:        order.Shop.Name,
		TotalAmount:     order.TotalAmount,
		ShippingFee:     order.ShippingFee,
		Status:          order.Status,
		StatusText:      statusText,
		ReceiverName:    order.ReceiverName,
		ReceiverPhone:   order.ReceiverPhone,
		ReceiverAddress: order.ReceiverAddress,
		Remark:          order.Remark,
		CreatedAt:       order.CreatedAt.Format(time.RFC3339),
		TrackingNo:      order.TrackingNo,
		TrackingCompany: order.TrackingCompany,
		Items:           items,
	}
	if order.PaidAt != nil {
		info.PaidAt = order.PaidAt.Format(time.RFC3339)
	}
	if order.ShippedAt != nil {
		info.ShippedAt = order.ShippedAt.Format(time.RFC3339)
	}
	if order.CompletedAt != nil {
		info.CompletedAt = order.CompletedAt.Format(time.RFC3339)
	}
	return info
}

func CreateNotification(db *gorm.DB, userID uint, notifType, title, content string, data map[string]interface{}) {
	var dataJSON string
	if data != nil {
		bytes, _ := json.Marshal(data)
		dataJSON = string(bytes)
	}

	notif := models.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Content: content,
		Data:    dataJSON,
		IsRead:  false,
	}
	db.Create(&notif)
}
