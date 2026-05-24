package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/service"
	"furniture-platform/pkg/response"
)

// ProductHandler 产品 HTTP 处理器
type ProductHandler struct {
	service *service.ProductService
}

// NewProductHandler 创建产品处理器
func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{service: svc}
}

// currentUser 从上下文中获取当前用户信息
func currentUser(c *gin.Context) (userID uint, role string, ok bool) {
	idVal, idExists := c.Get("user_id")
	roleVal, roleExists := c.Get("role")
	if !idExists || !roleExists {
		return 0, "", false
	}
	id, idOK := idVal.(uint)
	r, rOK := roleVal.(string)
	if !idOK || !rOK {
		return 0, "", false
	}
	return id, r, true
}

// toProductResponse 将模型转为响应 DTO
func toProductResponse(p *model.Product, options []*model.ProductOption, images []*model.ProductImage) dto.ProductResponse {
	opts := make([]dto.ProductOptionDTO, 0, len(options))
	for _, o := range options {
		opts = append(opts, dto.ProductOptionDTO{
			ID:          o.ID,
			OptionType:  o.OptionType,
			OptionValue: o.OptionValue,
			PriceAdjust: o.PriceAdjust,
			Sort:        o.Sort,
		})
	}
	imgs := make([]dto.ProductImageDTO, 0, len(images))
	for _, img := range images {
		imgs = append(imgs, dto.ProductImageDTO{
			ID:       img.ID,
			ImageURL: img.ImageURL,
			Sort:     img.Sort,
		})
	}
	return dto.ProductResponse{
		ID:             p.ID,
		ManufacturerID: p.ManufacturerID,
		Name:           p.Name,
		Category:       p.Category,
		Description:    p.Description,
		BasePrice:      p.BasePrice,
		Stock:          p.Stock,
		Status:         p.Status,
		IsHot:          p.IsHot,
		Options:        opts,
		Images:         imgs,
		CreatedAt:      p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// CreateProduct 创建产品
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	userID, _, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	product, options, images, err := h.service.Create(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toProductResponse(product, options, images))
}

// UpdateProduct 编辑产品
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	userID, _, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := parseProductID(c)
	if err != nil {
		response.BadRequest(c, "产品 ID 格式错误")
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	product, options, images, err := h.service.Update(id, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toProductResponse(product, options, images))
}

// DeleteProduct 删除产品
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	userID, _, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := parseProductID(c)
	if err != nil {
		response.BadRequest(c, "产品 ID 格式错误")
		return
	}

	if err := h.service.Delete(id, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetProduct 获取产品详情
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := parseProductID(c)
	if err != nil {
		response.BadRequest(c, "产品 ID 格式错误")
		return
	}

	product, options, images, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, toProductResponse(product, options, images))
}

// ListProducts 分页查询产品列表
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var req dto.ProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	products, optsMap, imgsMap, total, err := h.service.List(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	list := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		list = append(list, toProductResponse(p, optsMap[p.ID], imgsMap[p.ID]))
	}

	response.Success(c, dto.ProductListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	})
}

// GetHotProducts 获取热门产品列表（可匿名访问，走 Redis 缓存）
func (h *ProductHandler) GetHotProducts(c *gin.Context) {
	items, err := h.service.ListHot()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	resp := make([]dto.HotProductResponse, 0, len(items))
	for _, it := range items {
		resp = append(resp, dto.HotProductResponse{
			ID:             it.ID,
			ManufacturerID: it.ManufacturerID,
			Name:           it.Name,
			Category:       it.Category,
			BasePrice:      it.BasePrice,
			Cover:          it.Cover,
		})
	}
	response.Success(c, resp)
}

// parseProductID 解析产品 ID 路径参数
func parseProductID(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
