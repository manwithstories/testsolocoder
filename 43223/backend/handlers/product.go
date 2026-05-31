package handlers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"coffee-platform/config"
	"coffee-platform/database"
	"coffee-platform/models"
	"coffee-platform/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	cfg *config.Config
}

func NewProductHandler(cfg *config.Config) *ProductHandler {
	return &ProductHandler{cfg: cfg}
}

func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(h.cfg.App.PageSize)))
	status := c.Query("status")
	keyword := c.Query("keyword")
	origin := c.Query("origin")
	roastLevel := c.Query("roast_level")
	processMethod := c.Query("process_method")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	roasterID := c.Query("roaster_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = h.cfg.App.PageSize
	}

	query := database.DB.Model(&models.Product{}).Preload("Images")

	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status = ?", models.ProductStatusOnSale)
	}

	if keyword != "" {
		query = query.Where("name LIKE ? OR origin LIKE ? OR flavor_notes LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if origin != "" {
		query = query.Where("origin = ?", origin)
	}
	if roastLevel != "" {
		query = query.Where("roast_level = ?", roastLevel)
	}
	if processMethod != "" {
		query = query.Where("process_method = ?", processMethod)
	}
	if minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}
	if roasterID != "" {
		query = query.Where("roaster_id = ?", roasterID)
	}

	allowedSortFields := map[string]bool{
		"created_at": true, "price": true, "cupping_score": true, "name": true,
	}
	if allowedSortFields[sortBy] {
		order := sortBy + " " + sortOrder
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	var total int64
	query.Count(&total)

	var products []models.Product
	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)

	utils.PaginatedResponse(c, products, total, page, pageSize)
}

func (h *ProductHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var product models.Product
	if err := database.DB.Preload("Images").Preload("Roaster").First(&product, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "商品不存在")
		return
	}

	utils.Success(c, product)
}

func (h *ProductHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	product := models.Product{
		Name:          req.Name,
		Origin:        req.Origin,
		Farm:          req.Farm,
		Variety:       req.Variety,
		Altitude:      req.Altitude,
		ProcessMethod: req.ProcessMethod,
		RoastLevel:    req.RoastLevel,
		FlavorNotes:   req.FlavorNotes,
		CuppingScore:  req.CuppingScore,
		Description:   req.Description,
		Price:         req.Price,
		Weight:        req.Weight,
		Stock:         req.Stock,
		Status:        req.Status,
		RoasterID:     userID,
	}

	if userRole == string(models.RoleAdmin) && req.Status == "" {
		product.Status = models.ProductStatusOnSale
	} else if product.Status == "" {
		product.Status = models.ProductStatusDraft
	}

	if err := database.DB.Create(&product).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}

	utils.SuccessWithMessage(c, "创建成功", product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "商品不存在")
		return
	}

	if product.RoasterID != userID && userRole != string(models.RoleAdmin) {
		utils.Error(c, http.StatusForbidden, "无权修改此商品")
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Origin != nil {
		updates["origin"] = *req.Origin
	}
	if req.Farm != nil {
		updates["farm"] = *req.Farm
	}
	if req.Variety != nil {
		updates["variety"] = *req.Variety
	}
	if req.Altitude != nil {
		updates["altitude"] = *req.Altitude
	}
	if req.ProcessMethod != nil {
		updates["process_method"] = *req.ProcessMethod
	}
	if req.RoastLevel != nil {
		updates["roast_level"] = *req.RoastLevel
	}
	if req.FlavorNotes != nil {
		updates["flavor_notes"] = *req.FlavorNotes
	}
	if req.CuppingScore != nil {
		updates["cupping_score"] = *req.CuppingScore
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
	}

	if err := database.DB.Model(&product).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	database.DB.Preload("Images").First(&product, id)
	utils.SuccessWithMessage(c, "更新成功", product)
}

func (h *ProductHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req models.UpdateProductStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	if err := database.DB.Model(&models.Product{}).Where("id = ?", id).
		Update("status", req.Status).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "状态更新成功", nil)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "商品不存在")
		return
	}

	if product.RoasterID != userID && userRole != string(models.RoleAdmin) {
		utils.Error(c, http.StatusForbidden, "无权删除此商品")
		return
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *ProductHandler) UploadImage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := os.MkdirAll(h.cfg.App.UploadDir, 0755); err != nil {
		utils.Error(c, http.StatusInternalServerError, "创建上传目录失败")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "请上传文件")
		return
	}

	if file.Size > h.cfg.App.MaxUploadSize {
		utils.Error(c, http.StatusBadRequest, fmt.Sprintf("文件大小不能超过%dMB", h.cfg.App.MaxUploadSize/(1024*1024)))
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		utils.Error(c, http.StatusBadRequest, "不支持的文件格式")
		return
	}

	filename := fmt.Sprintf("%d-%s%s", id, uuid.New().String(), ext)
	filepath := filepath.Join(h.cfg.App.UploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		utils.Error(c, http.StatusInternalServerError, "文件保存失败")
		return
	}

	imageURL := "/uploads/" + filename

	isCover := c.DefaultPostForm("is_cover", "false") == "true"

	productImage := models.ProductImage{
		ProductID: uint(id),
		URL:       imageURL,
		IsCover:   isCover,
	}

	if isCover {
		database.DB.Model(&models.ProductImage{}).Where("product_id = ?", id).Update("is_cover", false)
	}

	if err := database.DB.Create(&productImage).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "保存图片信息失败")
		return
	}

	utils.Success(c, productImage)
}

func (h *ProductHandler) DeleteImage(c *gin.Context) {
	imageID, _ := strconv.ParseUint(c.Param("image_id"), 10, 64)

	var image models.ProductImage
	if err := database.DB.First(&image, imageID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "图片不存在")
		return
	}

	os.Remove(filepath.Join(".", image.URL))

	if err := database.DB.Delete(&image).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *ProductHandler) ImportCSV(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "请上传CSV文件")
		return
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
		utils.Error(c, http.StatusBadRequest, "请上传CSV格式文件")
		return
	}

	f, err := file.Open()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "文件读取失败")
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	headers, err := reader.Read()
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "CSV文件读取失败")
		return
	}

	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.TrimSpace(strings.ToLower(h))] = i
	}

	requiredFields := []string{"name", "origin", "process_method", "roast_level", "price", "weight"}
	for _, field := range requiredFields {
		if _, ok := headerMap[field]; !ok {
			utils.Error(c, http.StatusBadRequest, "CSV缺少必要字段: "+field)
			return
		}
	}

	var products []models.Product
	rowNum := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		rowNum++

		getField := func(name string) string {
			if idx, ok := headerMap[name]; ok && idx < len(record) {
				return strings.TrimSpace(record[idx])
			}
			return ""
		}

		price, _ := strconv.ParseFloat(getField("price"), 64)
		weight, _ := strconv.Atoi(getField("weight"))
		stock, _ := strconv.Atoi(getField("stock"))
		score, _ := strconv.ParseFloat(getField("cupping_score"), 64)

		if getField("name") == "" {
			continue
		}

		product := models.Product{
			Name:          getField("name"),
			Origin:        getField("origin"),
			Farm:          getField("farm"),
			Variety:       getField("variety"),
			Altitude:      getField("altitude"),
			ProcessMethod: models.ProcessMethod(getField("process_method")),
			RoastLevel:    models.RoastLevel(getField("roast_level")),
			FlavorNotes:   getField("flavor_notes"),
			CuppingScore:  score,
			Description:   getField("description"),
			Price:         price,
			Weight:        weight,
			Stock:         stock,
			Status:        models.ProductStatusDraft,
			RoasterID:     userID,
		}

		if product.Price <= 0 || product.Weight <= 0 {
			continue
		}

		products = append(products, product)

		if len(products) >= 100 {
			break
		}
	}

	if len(products) == 0 {
		utils.Error(c, http.StatusBadRequest, "没有有效的商品数据")
		return
	}

	var successCount int
	var failedRows []int
	for i, product := range products {
		if err := database.DB.Create(&product).Error; err != nil {
			failedRows = append(failedRows, i+2)
		} else {
			successCount++
		}
	}

	utils.Success(c, gin.H{
		"success_count": successCount,
		"total":         len(products),
		"failed_rows":   failedRows,
	})
}

func (h *ProductHandler) GetOrigins(c *gin.Context) {
	var origins []string
	database.DB.Model(&models.Product{}).
		Where("status = ?", models.ProductStatusOnSale).
		Distinct("origin").
		Pluck("origin", &origins)

	utils.Success(c, origins)
}

func (h *ProductHandler) GetRoastLevels(c *gin.Context) {
	utils.Success(c, []string{"light", "medium", "medium_dark", "dark"})
}

func (h *ProductHandler) GetProcessMethods(c *gin.Context) {
	utils.Success(c, []string{"washed", "natural", "honey", "anaerobic", "wet_hulled"})
}

var ErrProductNotFound = errors.New("product not found")
