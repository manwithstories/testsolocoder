package handlers

import (
	"multishop/internal/database"
	"multishop/internal/dto"
	"multishop/internal/middleware"
	"multishop/internal/models"
	"multishop/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ShopHandler struct{}

func NewShopHandler() *ShopHandler {
	return &ShopHandler{}
}

func (h *ShopHandler) Apply(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req dto.ShopApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	var count int64
	database.DB.Model(&models.Shop{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		utils.Error(c, http.StatusBadRequest, "店铺名称已存在")
		return
	}

	database.DB.Model(&models.Shop{}).Where("user_id = ?", userID).Count(&count)
	if count > 0 {
		utils.Error(c, http.StatusBadRequest, "您已申请过店铺")
		return
	}

	shop := models.Shop{
		UserID:          userID,
		Name:            req.Name,
		Description:     req.Description,
		Logo:            req.Logo,
		ContactName:     req.ContactName,
		ContactPhone:    req.ContactPhone,
		Address:         req.Address,
		IDCardFront:     req.IDCardFront,
		IDCardBack:      req.IDCardBack,
		BusinessLicense: req.BusinessLicense,
		Status:          models.ShopStatusPending,
	}

	if err := database.DB.Create(&shop).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "申请失败")
		return
	}

	utils.Success(c, gin.H{"shop_id": shop.ID})
}

func (h *ShopHandler) GetMyShop(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var shop models.Shop
	if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "店铺不存在")
		return
	}

	var productCount int64
	database.DB.Model(&models.Product{}).Where("shop_id = ?", shop.ID).Count(&productCount)
	
	var soldCount int64
	database.DB.Model(&models.Product{}).Where("shop_id = ?", shop.ID).Select("COALESCE(SUM(sales), 0)").Scan(&soldCount)

	utils.Success(c, dto.ShopDetail{
		ShopInfo: dto.ShopInfo{
			ID:           shop.ID,
			UserID:       shop.UserID,
			Name:         shop.Name,
			Description:  shop.Description,
			Logo:         shop.Logo,
			Status:       shop.Status,
			Rating:       shop.Rating,
			CreatedAt:    shop.CreatedAt.Format(time.RFC3339),
			ContactName:  shop.ContactName,
			ContactPhone: shop.ContactPhone,
			Address:      shop.Address,
			ProductCount: productCount,
			SoldCount:    soldCount,
		},
		RejectReason:    shop.RejectReason,
		IDCardFront:     shop.IDCardFront,
		IDCardBack:      shop.IDCardBack,
		BusinessLicense: shop.BusinessLicense,
	})
}

func (h *ShopHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req dto.ShopUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var shop models.Shop
	if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "店铺不存在")
		return
	}

	if shop.Status != models.ShopStatusApproved {
		utils.Error(c, http.StatusBadRequest, "店铺未审核通过，无法修改")
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Logo != "" {
		updates["logo"] = req.Logo
	}
	if req.ContactName != "" {
		updates["contact_name"] = req.ContactName
	}
	if req.ContactPhone != "" {
		updates["contact_phone"] = req.ContactPhone
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}

	if err := database.DB.Model(&shop).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	utils.Success(c, nil)
}

func (h *ShopHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var shop models.Shop
	if err := database.DB.Preload("User").First(&shop, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "店铺不存在")
		return
	}

	utils.Success(c, dto.ShopInfo{
		ID:          shop.ID,
		UserID:      shop.UserID,
		Name:        shop.Name,
		Description: shop.Description,
		Logo:        shop.Logo,
		Status:      shop.Status,
		Rating:      shop.Rating,
		CreatedAt:   shop.CreatedAt.Format(time.RFC3339),
	})
}

func (h *ShopHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	query := database.DB.Model(&models.Shop{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var shops []models.Shop
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&shops)

	shopInfos := make([]dto.ShopInfo, 0, len(shops))
	for _, shop := range shops {
		var productCount int64
		database.DB.Model(&models.Product{}).Where("shop_id = ?", shop.ID).Count(&productCount)
		
		var soldCount int64
		database.DB.Model(&models.Product{}).Where("shop_id = ?", shop.ID).Select("COALESCE(SUM(sales), 0)").Scan(&soldCount)
		
		shopInfos = append(shopInfos, dto.ShopInfo{
			ID:           shop.ID,
			UserID:       shop.UserID,
			Name:         shop.Name,
			Description:  shop.Description,
			Logo:         shop.Logo,
			Status:       shop.Status,
			Rating:       shop.Rating,
			CreatedAt:    shop.CreatedAt.Format(time.RFC3339),
			ContactName:  shop.ContactName,
			ContactPhone: shop.ContactPhone,
			Address:      shop.Address,
			ProductCount: productCount,
			SoldCount:    soldCount,
		})
	}

	utils.Paginated(c, shopInfos, total, page, pageSize)
}

func (h *ShopHandler) Review(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req dto.ShopReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var shop models.Shop
	if err := database.DB.First(&shop, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "店铺不存在")
		return
	}

	updates := make(map[string]interface{})
	updates["status"] = req.Status
	if req.Status == models.ShopStatusApproved {
		now := time.Now()
		updates["approved_at"] = &now
	} else {
		updates["reject_reason"] = req.RejectReason
	}

	if err := database.DB.Model(&shop).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "审核失败")
		return
	}

	utils.Success(c, nil)
}
