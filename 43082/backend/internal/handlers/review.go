package handlers

import (
	"encoding/json"
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

type ReviewHandler struct{}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{}
}

func (h *ReviewHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req dto.ReviewCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var orderItem models.OrderItem
	if err := database.DB.Where("id = ? AND order_id = ? AND product_id = ?", req.OrderItemID, req.OrderID, req.ProductID).
		First(&orderItem).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "订单项不存在")
		return
	}

	if orderItem.Reviewed {
		utils.Error(c, http.StatusBadRequest, "该商品已评价")
		return
	}

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ? AND status = ?", req.OrderID, userID, models.OrderStatusCompleted).
		First(&order).Error; err != nil {
		utils.Error(c, http.StatusBadRequest, "订单不存在或未完成")
		return
	}

	imagesJSON, _ := json.Marshal(req.Images)
	review := models.Review{
		UserID:      userID,
		ProductID:   req.ProductID,
		OrderID:     req.OrderID,
		OrderItemID: req.OrderItemID,
		ShopID:      order.ShopID,
		Rating:      req.Rating,
		Content:     req.Content,
		Images:      string(imagesJSON),
	}

	tx := database.DB.Begin()
	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		utils.Error(c, http.StatusInternalServerError, "评价失败")
		return
	}

	tx.Model(&orderItem).Update("reviewed", true)

	var reviews []models.Review
	tx.Where("shop_id = ?", order.ShopID).Find(&reviews)
	if len(reviews) > 0 {
		var totalRating int
		for _, r := range reviews {
			totalRating += r.Rating
		}
		avgRating := float64(totalRating) / float64(len(reviews))
		tx.Model(&models.Shop{}).Where("id = ?", order.ShopID).Update("rating", avgRating)
	}

	tx.Commit()
	utils.Success(c, gin.H{"review_id": review.ID})
}

func (h *ReviewHandler) Reply(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req dto.ReviewReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var shop models.Shop
	if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "店铺不存在")
		return
	}

	var review models.Review
	if err := database.DB.Where("id = ? AND shop_id = ?", id, shop.ID).First(&review).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "评价不存在")
		return
	}

	now := time.Now()
	review.Reply = req.Reply
	review.ReplyAt = &now
	database.DB.Save(&review)

	utils.Success(c, nil)
}

func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	productID, _ := strconv.ParseUint(c.Param("product_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := database.DB.Model(&models.Review{}).Where("product_id = ?", productID)
	var total int64
	query.Count(&total)

	var reviews []models.Review
	offset := (page - 1) * pageSize
	query.Preload("User").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&reviews)

	result := make([]dto.ReviewInfo, 0, len(reviews))
	for _, r := range reviews {
		var images []string
		json.Unmarshal([]byte(r.Images), &images)

		reviewInfo := dto.ReviewInfo{
			ID:        r.ID,
			UserID:    r.UserID,
			Username:  r.User.Username,
			Avatar:    r.User.Avatar,
			ProductID: r.ProductID,
			Rating:    r.Rating,
			Content:   r.Content,
			Images:    images,
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
		}
		if r.Reply != "" {
			reviewInfo.Reply = r.Reply
			reviewInfo.ReplyAt = r.ReplyAt.Format(time.RFC3339)
		}
		result = append(result, reviewInfo)
	}

	utils.Paginated(c, result, total, page, pageSize)
}

func (h *ReviewHandler) GetShopReviews(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var shop models.Shop
	if err := database.DB.Where("user_id = ?", userID).First(&shop).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "店铺不存在")
		return
	}

	query := database.DB.Model(&models.Review{}).Where("shop_id = ?", shop.ID)
	var total int64
	query.Count(&total)

	var reviews []models.Review
	offset := (page - 1) * pageSize
	query.Preload("User").Preload("Product").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&reviews)

	result := make([]dto.ReviewInfo, 0, len(reviews))
	for _, r := range reviews {
		var images []string
		json.Unmarshal([]byte(r.Images), &images)

		reviewInfo := dto.ReviewInfo{
			ID:        r.ID,
			UserID:    r.UserID,
			Username:  r.User.Username,
			Avatar:    r.User.Avatar,
			ProductID: r.ProductID,
			Rating:    r.Rating,
			Content:   r.Content,
			Images:    images,
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
		}
		if r.Reply != "" {
			reviewInfo.Reply = r.Reply
			reviewInfo.ReplyAt = r.ReplyAt.Format(time.RFC3339)
		}
		result = append(result, reviewInfo)
	}

	utils.Paginated(c, result, total, page, pageSize)
}
