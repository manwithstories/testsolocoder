package handlers

import (
	"fmt"
	"garden-planner/database"
	"garden-planner/middleware"
	"garden-planner/models"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateProductRequest struct {
	Name           string  `json:"name" binding:"required,max=200"`
	Description    string  `json:"description"`
	Category       string  `json:"category"`
	Price          float64 `json:"price" binding:"required"`
	Stock          int     `json:"stock"`
	ImageURLs      string  `json:"image_urls"`
	Specifications string  `json:"specifications"`
}

type UpdateProductRequest struct {
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Category       string  `json:"category"`
	Price          float64 `json:"price"`
	Stock          int     `json:"stock"`
	ImageURLs      string  `json:"image_urls"`
	Specifications string  `json:"specifications"`
	IsActive       *bool   `json:"is_active"`
}

type AddToCartRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity"`
}

type UpdateCartRequest struct {
	Quantity int `json:"quantity"`
}

type CreateOrderRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
	ShippingPhone   string `json:"shipping_phone" binding:"required"`
	ShippingName    string `json:"shipping_name" binding:"required"`
	PaymentMethod   string `json:"payment_method"`
	Remark          string `json:"remark"`
}

type UpdateOrderRequest struct {
	Status         string `json:"status"`
	TrackingNumber string `json:"tracking_number"`
	ShippingStatus string `json:"shipping_status"`
}

func CreateProduct(c *gin.Context) {
	userID := middleware.GetUserID(c)
	userType := middleware.GetUserType(c)

	if userType != "merchant" && userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only merchants and admins can create products"})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		ID:             uuid.New(),
		SellerID:       userID,
		Name:           req.Name,
		Description:    req.Description,
		Category:       req.Category,
		Price:          req.Price,
		Stock:          req.Stock,
		ImageURLs:      req.ImageURLs,
		Specifications: req.Specifications,
		IsActive:       true,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

func GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	category := c.Query("category")
	sellerID := c.Query("seller_id")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	var products []models.Product
	var total int64

	query := database.DB.Model(&models.Product{}).Where("is_active = ?", true)

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if sellerID != "" {
		query = query.Where("seller_id = ?", sellerID)
	}
	if minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			query = query.Where("price >= ?", price)
		}
	}
	if maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			query = query.Where("price <= ?", price)
		}
	}
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	orderClause := sortBy + " " + sortOrder
	query.Preload("Seller").
		Offset(offset).Limit(pageSize).
		Order(orderClause).
		Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"products":  products,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := database.DB.Preload("Seller").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func UpdateProduct(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var product models.Product
	if err := database.DB.Where("id = ? AND seller_id = ?", id, userID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.ImageURLs != "" {
		product.ImageURLs = req.ImageURLs
	}
	if req.Specifications != "" {
		product.Specifications = req.Specifications
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": product,
	})
}

func DeleteProduct(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	result := database.DB.Where("id = ? AND seller_id = ?", id, userID).Delete(&models.Product{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func AddToCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	quantity := req.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	var existingCart models.Cart
	result := database.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&existingCart)

	if result.Error == nil {
		existingCart.Quantity += quantity
		if err := database.DB.Save(&existingCart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Cart updated successfully",
			"cart":    existingCart,
		})
		return
	}

	cart := models.Cart{
		ID:        uuid.New(),
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  quantity,
	}

	if err := database.DB.Create(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Added to cart successfully",
		"cart":    cart,
	})
}

func GetCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var cartItems []models.Cart
	database.DB.Where("user_id = ?", userID).
		Preload("Product").
		Find(&cartItems)

	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += item.Product.Price * float64(item.Quantity)
	}

	c.JSON(http.StatusOK, gin.H{
		"cart_items":   cartItems,
		"total_amount": totalAmount,
		"total_count":  len(cartItems),
	})
}

func UpdateCart(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var cart models.Cart
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	var req UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Quantity <= 0 {
		database.DB.Delete(&cart)
		c.JSON(http.StatusOK, gin.H{"message": "Cart item removed"})
		return
	}

	cart.Quantity = req.Quantity
	if err := database.DB.Save(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart updated successfully",
		"cart":    cart,
	})
}

func RemoveFromCart(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Cart{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from cart"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from cart successfully"})
}

func ClearCart(c *gin.Context) {
	userID := middleware.GetUserID(c)

	database.DB.Where("user_id = ?", userID).Delete(&models.Cart{})

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}

func CreateOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cartItems []models.Cart
	if err := database.DB.Where("user_id = ?", userID).
		Preload("Product").
		Find(&cartItems).Error; err != nil || len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	for _, cartItem := range cartItems {
		if cartItem.Product.Stock < cartItem.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Product %s is out of stock", cartItem.Product.Name),
			})
			return
		}

		subtotal := cartItem.Product.Price * float64(cartItem.Quantity)
		totalAmount += subtotal

		orderItems = append(orderItems, models.OrderItem{
			ID:          uuid.New(),
			ProductID:   cartItem.ProductID,
			ProductName: cartItem.Product.Name,
			Price:       cartItem.Product.Price,
			Quantity:    cartItem.Quantity,
			Subtotal:    subtotal,
		})
	}

	orderNo := generateOrderNo()

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		order := models.Order{
			ID:              uuid.New(),
			UserID:          userID,
			OrderNo:         orderNo,
			TotalAmount:     totalAmount,
			Status:          "pending",
			PaymentMethod:   req.PaymentMethod,
			PaymentStatus:   "unpaid",
			ShippingAddress: req.ShippingAddress,
			ShippingPhone:   req.ShippingPhone,
			ShippingName:    req.ShippingName,
			Remark:          req.Remark,
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for i := range orderItems {
			orderItems[i].OrderID = order.ID
			if err := tx.Create(&orderItems[i]).Error; err != nil {
				return err
			}
		}

		for _, cartItem := range cartItems {
			if err := tx.Model(&models.Product{}).
				Where("id = ?", cartItem.ProductID).
				UpdateColumn("stock", gorm.Expr("stock - ?", cartItem.Quantity)).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Order created successfully",
		"order_no": orderNo,
		"total":    totalAmount,
	})
}

func GetOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	status := c.Query("status")

	var orders []models.Order
	var total int64

	query := database.DB.Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Model(&models.Order{}).Count(&total)
	query.Preload("OrderItems.Product").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"orders":    orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).
		Preload("OrderItems.Product").
		First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

func UpdateOrder(c *gin.Context) {
	userType := middleware.GetUserType(c)
	id := c.Param("id")

	if userType != "merchant" && userType != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only merchants and admins can update orders"})
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status != "" {
		order.Status = req.Status
	}
	if req.TrackingNumber != "" {
		order.TrackingNumber = req.TrackingNumber
	}
	if req.ShippingStatus != "" {
		order.ShippingStatus = req.ShippingStatus
	}

	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order updated successfully",
		"order":   order,
	})
}

func CancelOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only pending orders can be cancelled"})
		return
	}

	order.Status = "cancelled"
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

func generateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("GP%s%06d", now.Format("20060102150405"), rand.Intn(1000000))
}
