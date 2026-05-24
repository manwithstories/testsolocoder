package handler

import (
	"mime/multipart"
	"strconv"
	"strings"

	"luxury-trading-platform/internal/model"
	"luxury-trading-platform/internal/service"
	resp "luxury-trading-platform/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	sellerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), sellerID.(uint), &req)
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Created(c, product)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), uint(id))
	if err != nil {
		resp.NotFound(c, err.Error())
		return
	}

	resp.Success(c, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	sellerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), uint(id), sellerID.(uint), &req)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			resp.Forbidden(c, err.Error())
			return
		}
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	sellerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), uint(id), sellerID.(uint)); err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			resp.Forbidden(c, err.Error())
			return
		}
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, gin.H{"message": "product deleted successfully"})
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	brandIDStr := c.Query("brand_id")
	status := c.Query("status")
	sortBy := c.DefaultQuery("sort_by", "newest")
	keyword := c.Query("keyword")

	var brandID *uint
	if brandIDStr != "" {
		id, err := strconv.ParseUint(brandIDStr, 10, 32)
		if err == nil {
			uid := uint(id)
			brandID = &uid
		}
	}

	products, total, err := h.productService.ListProducts(page, pageSize, model.ProductCategory(category), brandID, model.ProductStatus(status), sortBy, keyword)
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.SuccessWithPage(c, products, total, page, pageSize)
}

func (h *ProductHandler) ListSellerProducts(c *gin.Context) {
	sellerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	products, total, err := h.productService.ListSellerProducts(sellerID.(uint), page, pageSize, model.ProductStatus(status))
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.SuccessWithPage(c, products, total, page, pageSize)
}

func (h *ProductHandler) UpdateProductStatus(c *gin.Context) {
	sellerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	var req struct {
		Status model.ProductStatus `json:"status" binding:"required,oneof=draft on_sale sold removed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	if err := h.productService.UpdateProductStatus(c.Request.Context(), uint(id), sellerID.(uint), req.Status); err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, gin.H{"message": "product status updated successfully"})
}

func (h *ProductHandler) UploadImages(c *gin.Context) {
	sellerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		resp.BadRequest(c, err)
		return
	}

	images, err := h.productService.UploadImages(c.Request.Context(), uint(id), sellerID.(uint), files)
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, images)
}

func (h *ProductHandler) GetProductImages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	images, err := h.productService.GetProductImages(uint(id))
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, images)
}

func (h *ProductHandler) DeleteProductImages(c *gin.Context) {
	sellerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	if err := h.productService.DeleteProductImages(c.Request.Context(), uint(id), sellerID.(uint)); err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, gin.H{"message": "product images deleted successfully"})
}

func (h *ProductHandler) CreateBrand(c *gin.Context) {
	var brand model.Brand
	if err := c.ShouldBindJSON(&brand); err != nil {
		resp.BadRequest(c, err)
		return
	}

	created, err := h.productService.CreateBrand(c.Request.Context(), &brand)
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Created(c, created)
}

func (h *ProductHandler) ListBrands(c *gin.Context) {
	category := c.Query("category")

	brands, err := h.productService.ListBrands(model.ProductCategory(category))
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, brands)
}

func (h *ProductHandler) GetBrand(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	brand, err := h.productService.GetBrand(uint(id))
	if err != nil {
		resp.NotFound(c, err.Error())
		return
	}

	resp.Success(c, brand)
}

type fileHeaderSlice []*multipart.FileHeader
