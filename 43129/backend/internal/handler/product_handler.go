package handler

import (
	"strconv"

	"beauty-salon-system/internal/service"
	"beauty-salon-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService      *service.ProductService
	notificationService *service.NotificationService
	auditService        *service.AuditService
}

func NewProductHandler(
	productService *service.ProductService,
	notificationService *service.NotificationService,
	auditService *service.AuditService,
) *ProductHandler {
	return &ProductHandler{
		productService:      productService,
		notificationService: notificationService,
		auditService:        auditService,
	}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.productService.Create(&req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	h.auditService.Log(userID.(uint), "create", "product", "Create product", c.ClientIP())

	response.Success(c, result)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	result, err := h.productService.GetByID(uint(id))
	if err != nil {
		response.Error(c, 404, "product not found")
		return
	}

	response.Success(c, result)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.productService.Update(uint(id), &req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.productService.Delete(uint(id)); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	lowStock, _ := strconv.ParseBool(c.DefaultQuery("low_stock", "false"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	result, total, err := h.productService.List(page, pageSize, category, lowStock)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  result,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

func (h *ProductHandler) ListAll(c *gin.Context) {
	result, err := h.productService.ListAll()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *ProductHandler) AddStock(c *gin.Context) {
	var req service.AddStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.productService.AddStock(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	h.auditService.Log(userID.(uint), "add_stock", "product", "Add stock", c.ClientIP())

	response.Success(c, nil)
}

func (h *ProductHandler) DeductStock(c *gin.Context) {
	var req service.DeductStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.productService.DeductStock(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ProductHandler) GetRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	productID, _ := strconv.ParseUint(c.Query("product_id"), 10, 32)
	changeType := c.Query("change_type")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	result, total, err := h.productService.GetProductRecords(page, pageSize, uint(productID), changeType)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  result,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

func (h *ProductHandler) GetLowStock(c *gin.Context) {
	result, err := h.productService.GetLowStockProducts()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	h.notificationService.SendLowStockAlert(result)

	response.Success(c, result)
}

func (h *ProductHandler) Sale(c *gin.Context) {
	var req service.SaleProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.productService.SaleProduct(&req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	h.auditService.Log(userID.(uint), "sale", "product", "Sale product", c.ClientIP())

	response.Success(c, result)
}

func (h *ProductHandler) GetSales(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	customerID, _ := strconv.ParseUint(c.Query("customer_id"), 10, 32)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	result, total, err := h.productService.GetProductSales(page, pageSize, uint(customerID))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  result,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

func (h *ProductHandler) StockTake(c *gin.Context) {
	var req map[uint]int
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.productService.RecordStockTake(req, userID.(uint)); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	h.auditService.Log(userID.(uint), "stock_take", "product", "Stock take", c.ClientIP())

	response.Success(c, nil)
}
