package handler

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"watchplatform/internal/app"
	"watchplatform/internal/config"
	"watchplatform/internal/database"
	"watchplatform/internal/middleware"
	"watchplatform/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WatchCreateReq struct {
	Brand        string  `json:"brand" binding:"required,max=64"`
	Model        string  `json:"model" binding:"required,max=128"`
	ReferenceNo  string  `json:"reference_no"`
	Year         int     `json:"year"`
	Movement     string  `json:"movement"`
	CaseSizeMM   float64 `json:"case_size_mm"`
	CaseMaterial string  `json:"case_material"`
	DialColor    string  `json:"dial_color"`
	Bracelet     string  `json:"bracelet"`
	Condition    string  `json:"condition"`
	Description  string  `json:"description"`
	Price        float64 `json:"price" binding:"required,gt=0"`
}

type WatchUpdateReq struct {
	Brand        *string  `json:"brand"`
	Model        *string  `json:"model"`
	ReferenceNo  *string  `json:"reference_no"`
	Year         *int     `json:"year"`
	Movement     *string  `json:"movement"`
	CaseSizeMM   *float64 `json:"case_size_mm"`
	CaseMaterial *string  `json:"case_material"`
	DialColor    *string  `json:"dial_color"`
	Bracelet     *string  `json:"bracelet"`
	Condition    *string  `json:"condition"`
	Description  *string  `json:"description"`
	Price        *float64 `json:"price"`
	Status       *string  `json:"status"`
}

func CreateWatch(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var req WatchCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	w := model.Watch{
		SellerID:     u.ID,
		Brand:        req.Brand,
		Model:        req.Model,
		ReferenceNo:  req.ReferenceNo,
		Year:         req.Year,
		Movement:     req.Movement,
		CaseSizeMM:   req.CaseSizeMM,
		CaseMaterial: req.CaseMaterial,
		DialColor:    req.DialColor,
		Bracelet:     req.Bracelet,
		Condition:    req.Condition,
		Description:  req.Description,
		Price:        req.Price,
		Status:       "on_sale",
	}
	if err := database.DB.Create(&w).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	app.OK(c, w)
}

func UpdateWatch(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var w model.Watch
	if err := database.DB.First(&w, id).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "手表不存在")
		return
	}
	if w.SellerID != u.ID {
		app.Fail(c, http.StatusForbidden, "无权修改")
		return
	}
	var req WatchUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	updates := map[string]interface{}{}
	if req.Brand != nil {
		updates["brand"] = *req.Brand
	}
	if req.Model != nil {
		updates["model"] = *req.Model
	}
	if req.ReferenceNo != nil {
		updates["reference_no"] = *req.ReferenceNo
	}
	if req.Year != nil {
		updates["year"] = *req.Year
	}
	if req.Movement != nil {
		updates["movement"] = *req.Movement
	}
	if req.CaseSizeMM != nil {
		updates["case_size_mm"] = *req.CaseSizeMM
	}
	if req.CaseMaterial != nil {
		updates["case_material"] = *req.CaseMaterial
	}
	if req.DialColor != nil {
		updates["dial_color"] = *req.DialColor
	}
	if req.Bracelet != nil {
		updates["bracelet"] = *req.Bracelet
	}
	if req.Condition != nil {
		updates["condition"] = *req.Condition
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if err := database.DB.Model(&w).Updates(updates).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	database.DB.Preload("Photos").First(&w, w.ID)
	app.OK(c, w)
}

func DeleteWatch(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var w model.Watch
	if err := database.DB.First(&w, id).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "手表不存在")
		return
	}
	if w.SellerID != u.ID {
		app.Fail(c, http.StatusForbidden, "无权删除")
		return
	}
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("watch_id = ?", w.ID).Delete(&model.WatchPhoto{}).Error; err != nil {
			return err
		}
		return tx.Delete(&w).Error
	}); err != nil {
		app.BizFail(c, err)
		return
	}
	app.OK(c, nil)
}

func ListWatches(c *gin.Context) {
	brand := c.Query("brand")
	keyword := c.Query("keyword")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	db := database.DB.Model(&model.Watch{}).Where("status = ?", "on_sale")
	if brand != "" {
		db = db.Where("brand = ?", brand)
	}
	if keyword != "" {
		k := "%" + keyword + "%"
		db = db.Where("brand LIKE ? OR model LIKE ? OR description LIKE ?", k, k, k)
	}
	if minPrice != "" {
		if v, err := strconv.ParseFloat(minPrice, 64); err == nil {
			db = db.Where("price >= ?", v)
		}
	}
	if maxPrice != "" {
		if v, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			db = db.Where("price <= ?", v)
		}
	}
	var total int64
	db.Count(&total)
	var list []model.Watch
	db.Preload("Photos").Order("created_at desc").Offset((page - 1) * size).Limit(size).Find(&list)
	app.OK(c, gin.H{"total": total, "page": page, "size": size, "list": list})
}

func GetWatch(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var w model.Watch
	if err := database.DB.Preload("Photos").First(&w, id).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "手表不存在")
		return
	}
	app.OK(c, w)
}

func UploadWatchPhotos(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var w model.Watch
	if err := database.DB.First(&w, id).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "手表不存在")
		return
	}
	if w.SellerID != u.ID {
		app.Fail(c, http.StatusForbidden, "无权上传")
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		app.Fail(c, http.StatusBadRequest, "解析文件失败")
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		app.Fail(c, http.StatusBadRequest, "未提供文件")
		return
	}
	dir := filepath.Join(config.Cfg.UploadDir, "watches", strconv.FormatUint(uint64(w.ID), 10))
	var urls []string
	for _, f := range files {
		if f.Size > config.Cfg.MaxFileSize {
			continue
		}
		ext := strings.ToLower(filepath.Ext(f.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			continue
		}
		name := uuid.New().String() + ext
		dst := filepath.Join(dir, name)
		if err := c.SaveUploadedFile(f, dst); err != nil {
			continue
		}
		url := "/uploads/watches/" + strconv.FormatUint(uint64(w.ID), 10) + "/" + name
		urls = append(urls, url)
		_ = database.DB.Create(&model.WatchPhoto{WatchID: w.ID, URL: url})
	}
	app.OK(c, gin.H{"urls": urls})
}
