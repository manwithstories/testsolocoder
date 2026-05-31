package controllers

import (
	"secondhand-platform/models"
	"secondhand-platform/services"
	"secondhand-platform/utils"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService *services.ProductService
}

func NewProductController() *ProductController {
	return &ProductController{
		productService: services.NewProductService(),
	}
}

type CreateProductRequest struct {
	Title         string  `json:"title" binding:"required,min=2,max=200"`
	Category      string  `json:"category" binding:"required"`
	Brand         string  `json:"brand" binding:"required,max=100"`
	Model         string  `json:"model" binding:"required,max=100"`
	Condition     string  `json:"condition" binding:"required,oneof=全新 95新 9成新 8成新 7成新及以下"`
	Price         float64 `json:"price" binding:"required,min=0"`
	OriginalPrice float64 `json:"original_price"`
	Description   string  `json:"description" binding:"required"`
	WarrantyDays  int     `json:"warranty_days"`
	Images        string  `json:"images"`
}

type UpdateProductRequest struct {
	Title         string  `json:"title" binding:"omitempty,min=2,max=200"`
	Category      string  `json:"category" binding:"omitempty"`
	Brand         string  `json:"brand" binding:"omitempty,max=100"`
	Model         string  `json:"model" binding:"omitempty,max=100"`
	Condition     string  `json:"condition" binding:"omitempty,oneof=全新 95新 9成新 8成新 7成新及以下"`
	Price         float64 `json:"price" binding:"omitempty,min=0"`
	OriginalPrice float64 `json:"original_price"`
	Description   string  `json:"description" binding:"omitempty"`
	WarrantyDays  int     `json:"warranty_days"`
	Images        string  `json:"images"`
}

type ReviewProductRequest struct {
	Approved     bool   `json:"approved"`
	RejectReason string `json:"reject_reason"`
}

func (ctrl *ProductController) CreateProduct(c *gin.Context) {
	sellerID := c.GetUint("user_id")

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	warrantyDays := req.WarrantyDays
	if warrantyDays <= 0 {
		warrantyDays = 30
	}

	product, err := ctrl.productService.CreateProduct(
		sellerID, req.Title, req.Category, req.Brand, req.Model,
		req.Condition, req.Price, req.OriginalPrice, req.Description,
		warrantyDays, req.Images,
	)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, product)
}

func (ctrl *ProductController) GetProduct(c *gin.Context) {
	id := parseIntParam(c, "id")

	product, err := ctrl.productService.GetProductByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "商品不存在")
		return
	}

	utils.Success(c, product)
}

func (ctrl *ProductController) ListProducts(c *gin.Context) {
	page := getPage(c)
	pageSize := getPageSize(c)
	category := c.Query("category")
	condition := c.Query("condition")
	minPrice := parseFloat(c.Query("min_price"), 0)
	maxPrice := parseFloat(c.Query("max_price"), 0)
	keyword := c.Query("keyword")
	sortBy := c.Query("sort_by")
	sellerID := parseInt(c.Query("seller_id"), 0)

	products, total, err := ctrl.productService.ListProducts(page, pageSize, category, condition, minPrice, maxPrice, keyword, sortBy, uint(sellerID))
	if err != nil {
		utils.Error(c, 500, "获取商品列表失败")
		return
	}

	utils.SuccessWithPagination(c, products, page, pageSize, total)
}

func (ctrl *ProductController) UpdateProduct(c *gin.Context) {
	userID := c.GetUint("user_id")
	productID := parseIntParam(c, "id")

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Brand != "" {
		updates["brand"] = req.Brand
	}
	if req.Model != "" {
		updates["model"] = req.Model
	}
	if req.Condition != "" {
		updates["condition"] = req.Condition
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.WarrantyDays > 0 {
		updates["warranty_days"] = req.WarrantyDays
	}
	if req.Images != "" {
		updates["images"] = req.Images
	}
	updates["original_price"] = req.OriginalPrice

	if err := ctrl.productService.UpdateProduct(userID, uint(productID), updates); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *ProductController) DeleteProduct(c *gin.Context) {
	userID := c.GetUint("user_id")
	productID := parseIntParam(c, "id")

	if err := ctrl.productService.DeleteProduct(userID, uint(productID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *ProductController) OffShelfProduct(c *gin.Context) {
	userID := c.GetUint("user_id")
	productID := parseIntParam(c, "id")

	if err := ctrl.productService.OffShelfProduct(userID, uint(productID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *ProductController) ReviewProduct(c *gin.Context) {
	productID := parseIntParam(c, "id")
	reviewerID := c.GetUint("user_id")

	var req ReviewProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	if err := ctrl.productService.ReviewProduct(uint(productID), req.Approved, req.RejectReason, reviewerID); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *ProductController) ListPendingProducts(c *gin.Context) {
	page := getPage(c)
	pageSize := getPageSize(c)

	products, total, err := ctrl.productService.ListPendingProducts(page, pageSize)
	if err != nil {
		utils.Error(c, 500, "获取待审核商品失败")
		return
	}

	utils.SuccessWithPagination(c, products, page, pageSize, total)
}

func (ctrl *ProductController) ListMyProducts(c *gin.Context) {
	userID := c.GetUint("user_id")
	page := getPage(c)
	pageSize := getPageSize(c)
	status := parseInt(c.Query("status"), 0)

	products, total, err := ctrl.productService.ListMyProducts(userID, page, pageSize, status)
	if err != nil {
		utils.Error(c, 500, "获取我的商品失败")
		return
	}

	utils.SuccessWithPagination(c, products, page, pageSize, total)
}

func (ctrl *ProductController) ToggleFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	productID := parseIntParam(c, "id")

	isFavorited, err := ctrl.productService.ToggleFavorite(userID, uint(productID))
	if err != nil {
		utils.Error(c, 500, "操作失败")
		return
	}

	utils.Success(c, gin.H{"is_favorited": isFavorited})
}

func (ctrl *ProductController) ListFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")
	page := getPage(c)
	pageSize := getPageSize(c)

	favorites, total, err := ctrl.productService.ListFavorites(userID, page, pageSize)
	if err != nil {
		utils.Error(c, 500, "获取收藏列表失败")
		return
	}

	utils.SuccessWithPagination(c, favorites, page, pageSize, total)
}

func (ctrl *ProductController) GetCategories(c *gin.Context) {
	utils.Success(c, gin.H{
		"categories": ctrl.productService.GetCategories(),
		"conditions": ctrl.productService.GetConditions(),
	})
}

func (ctrl *ProductController) GetHotProducts(c *gin.Context) {
	limit := parseInt(c.Query("limit"), 10)

	products, err := ctrl.productService.GetHotProducts(limit)
	if err != nil {
		utils.Error(c, 500, "获取热门商品失败")
		return
	}

	utils.Success(c, products)
}

type RepairController struct {
	repairService *services.RepairService
}

func NewRepairController() *RepairController {
	return &RepairController{
		repairService: services.NewRepairService(),
	}
}

type CreateServiceRequest struct {
	ServiceType   string  `json:"service_type" binding:"required"`
	Title         string  `json:"title" binding:"required,min=2,max=200"`
	Description   string  `json:"description" binding:"required"`
	Price         float64 `json:"price" binding:"required,min=0"`
	MinPrice      float64 `json:"min_price"`
	MaxPrice      float64 `json:"max_price"`
	EstimatedDays int     `json:"estimated_days"`
	Images        string  `json:"images"`
}

type CreateRepairOrderRequest struct {
	TechnicianID    uint    `json:"technician_id" binding:"required"`
	ServiceID       uint    `json:"service_id" binding:"required"`
	DeviceType      string  `json:"device_type" binding:"required"`
	DeviceBrand     string  `json:"device_brand" binding:"required"`
	DeviceModel     string  `json:"device_model" binding:"required"`
	FaultDescription string `json:"fault_description" binding:"required"`
	ContactName     string  `json:"contact_name" binding:"required"`
	ContactPhone    string  `json:"contact_phone" binding:"required"`
	Address         string  `json:"address"`
	ServicePrice    float64 `json:"service_price" binding:"required,min=0"`
}

func (ctrl *RepairController) CreateService(c *gin.Context) {
	technicianID := c.GetUint("user_id")

	var req CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	service, err := ctrl.repairService.CreateService(
		technicianID, req.ServiceType, req.Title, req.Description,
		req.Price, req.MinPrice, req.MaxPrice, req.EstimatedDays, req.Images,
	)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, service)
}

func (ctrl *RepairController) GetService(c *gin.Context) {
	id := parseIntParam(c, "id")

	service, err := ctrl.repairService.GetServiceByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "服务不存在")
		return
	}

	utils.Success(c, service)
}

func (ctrl *RepairController) ListServices(c *gin.Context) {
	page := getPage(c)
	pageSize := getPageSize(c)
	serviceType := c.Query("service_type")
	keyword := c.Query("keyword")
	sortBy := c.Query("sort_by")
	technicianID := parseInt(c.Query("technician_id"), 0)

	services, total, err := ctrl.repairService.ListServices(page, pageSize, serviceType, keyword, sortBy, uint(technicianID))
	if err != nil {
		utils.Error(c, 500, "获取服务列表失败")
		return
	}

	utils.SuccessWithPagination(c, services, page, pageSize, total)
}

func (ctrl *RepairController) UpdateService(c *gin.Context) {
	userID := c.GetUint("user_id")
	serviceID := parseIntParam(c, "id")

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	if err := ctrl.repairService.UpdateService(userID, uint(serviceID), req); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *RepairController) DeleteService(c *gin.Context) {
	userID := c.GetUint("user_id")
	serviceID := parseIntParam(c, "id")

	if err := ctrl.repairService.DeleteService(userID, uint(serviceID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *RepairController) GetServiceTypes(c *gin.Context) {
	utils.Success(c, gin.H{"service_types": ctrl.repairService.GetServiceTypes()})
}

func (ctrl *RepairController) CreateRepairOrder(c *gin.Context) {
	buyerID := c.GetUint("user_id")

	var req CreateRepairOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	order, err := ctrl.repairService.CreateRepairOrder(
		buyerID, req.TechnicianID, req.ServiceID,
		req.DeviceType, req.DeviceBrand, req.DeviceModel,
		req.FaultDescription, req.ContactName, req.ContactPhone,
		req.Address, req.ServicePrice,
	)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, order)
}

func (ctrl *RepairController) GetRepairOrder(c *gin.Context) {
	id := parseIntParam(c, "id")

	order, err := ctrl.repairService.GetRepairOrderByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "订单不存在")
		return
	}

	utils.Success(c, order)
}

func (ctrl *RepairController) ListRepairOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")
	page := getPage(c)
	pageSize := getPageSize(c)
	status := parseInt(c.Query("status"), 0)

	orders, total, err := ctrl.repairService.ListRepairOrders(page, pageSize, userID, userRole, status)
	if err != nil {
		utils.Error(c, 500, "获取维修订单失败")
		return
	}

	utils.SuccessWithPagination(c, orders, page, pageSize, total)
}

func (ctrl *RepairController) AcceptRepairOrder(c *gin.Context) {
	technicianID := c.GetUint("user_id")
	orderID := parseIntParam(c, "id")

	if err := ctrl.repairService.AcceptRepairOrder(technicianID, uint(orderID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *RepairController) StartRepair(c *gin.Context) {
	technicianID := c.GetUint("user_id")
	orderID := parseIntParam(c, "id")

	if err := ctrl.repairService.StartRepair(technicianID, uint(orderID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *RepairController) CompleteRepair(c *gin.Context) {
	technicianID := c.GetUint("user_id")
	orderID := parseIntParam(c, "id")

	var req struct {
		FinalPrice float64 `json:"final_price"`
	}
	c.ShouldBindJSON(&req)

	if err := ctrl.repairService.CompleteRepair(technicianID, uint(orderID), req.FinalPrice); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *RepairController) PickUpDevice(c *gin.Context) {
	userID := c.GetUint("user_id")
	orderID := parseIntParam(c, "id")

	if err := ctrl.repairService.PickUpDevice(userID, uint(orderID)); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *RepairController) CancelRepairOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")
	orderID := parseIntParam(c, "id")

	if err := ctrl.repairService.CancelRepairOrder(userID, uint(orderID), userRole); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

type OrderController struct {
	orderService *services.OrderService
}

func NewOrderController() *OrderController {
	return &OrderController{
		orderService: services.NewOrderService(),
	}
}

type CreateOrderRequest struct {
	ProductID       uint    `json:"product_id" binding:"required"`
	ReceiverName    string  `json:"receiver_name" binding:"required"`
	ReceiverPhone   string  `json:"receiver_phone" binding:"required"`
	ReceiverAddress string  `json:"receiver_address" binding:"required"`
	NegotiatedPrice float64 `json:"negotiated_price"`
}

type PayOrderRequest struct {
	OrderNo       string `json:"order_no" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=alipay wechat wallet"`
}

type ShipOrderRequest struct {
	OrderNo         string `json:"order_no" binding:"required"`
	TrackingNo      string `json:"tracking_no" binding:"required"`
	TrackingCompany string `json:"tracking_company" binding:"required"`
}

type NegotiatePriceRequest struct {
	OrderNo       string  `json:"order_no" binding:"required"`
	OfferedPrice  float64 `json:"offered_price" binding:"required,min=0"`
	Message       string  `json:"message"`
}

type HandleNegotiationRequest struct {
	OrderNo      string  `json:"order_no" binding:"required"`
	Accepted     bool    `json:"accepted"`
	CounterPrice float64 `json:"counter_price"`
	Message      string  `json:"message"`
}

type RefundRequest struct {
	OrderNo string `json:"order_no" binding:"required"`
	Reason  string `json:"reason"`
}

type HandleRefundRequest struct {
	OrderID  uint `json:"order_id" binding:"required"`
	Approved bool `json:"approved"`
}

func (ctrl *OrderController) CreateOrder(c *gin.Context) {
	buyerID := c.GetUint("user_id")

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	order, err := ctrl.orderService.CreateOrder(
		buyerID, req.ProductID, req.ReceiverName, req.ReceiverPhone,
		req.ReceiverAddress, req.NegotiatedPrice,
	)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, order)
}

func (ctrl *OrderController) GetOrder(c *gin.Context) {
	id := parseIntParam(c, "id")

	order, err := ctrl.orderService.GetOrderByID(uint(id))
	if err != nil {
		utils.Error(c, 404, "订单不存在")
		return
	}

	utils.Success(c, order)
}

func (ctrl *OrderController) GetOrderByNo(c *gin.Context) {
	orderNo := c.Param("order_no")

	order, err := ctrl.orderService.GetOrderByOrderNo(orderNo)
	if err != nil {
		utils.Error(c, 404, "订单不存在")
		return
	}

	utils.Success(c, order)
}

func (ctrl *OrderController) ListOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")
	page := getPage(c)
	pageSize := getPageSize(c)
	status := parseInt(c.Query("status"), 0)

	orders, total, err := ctrl.orderService.ListOrders(page, pageSize, userID, userRole, status)
	if err != nil {
		utils.Error(c, 500, "获取订单列表失败")
		return
	}

	utils.SuccessWithPagination(c, orders, page, pageSize, total)
}

func (ctrl *OrderController) PayOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req PayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.orderService.PayOrder(userID, req.OrderNo, req.PaymentMethod); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *OrderController) ShipOrder(c *gin.Context) {
	sellerID := c.GetUint("user_id")

	var req ShipOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.orderService.ShipOrder(sellerID, req.OrderNo, req.TrackingNo, req.TrackingCompany); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *OrderController) ConfirmDelivery(c *gin.Context) {
	userID := c.GetUint("user_id")
	orderNo := c.Param("order_no")

	if err := ctrl.orderService.ConfirmDelivery(userID, orderNo); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *OrderController) CancelOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")
	orderNo := c.Param("order_no")

	if err := ctrl.orderService.CancelOrder(userID, orderNo, userRole); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *OrderController) RefundOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.orderService.RefundOrder(userID, req.OrderNo, req.Reason); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *OrderController) HandleRefund(c *gin.Context) {
	var req HandleRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	if err := ctrl.orderService.HandleRefund(req.OrderID, req.Approved); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *OrderController) NegotiatePrice(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req NegotiatePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.orderService.NegotiatePrice(userID, req.OrderNo, req.OfferedPrice, req.Message); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *OrderController) HandleNegotiation(c *gin.Context) {
	sellerID := c.GetUint("user_id")

	var req HandleNegotiationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.orderService.HandleNegotiation(sellerID, req.OrderNo, req.Accepted, req.CounterPrice, req.Message); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *OrderController) GetOrderStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	stats, err := ctrl.orderService.GetOrderStats(userID, userRole)
	if err != nil {
		utils.Error(c, 500, "获取订单统计失败")
		return
	}

	utils.Success(c, stats)
}

func init() {
	_ = models.ProductStatusOnSale
}
