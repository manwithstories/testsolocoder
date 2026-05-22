package handlers

import (
	"strconv"
	"time"

	"housekeeping-platform/config"
	"housekeeping-platform/models"
	"housekeeping-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateReview(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	var req struct {
		OrderID            uint   `json:"order_id" binding:"required"`
		AttitudeRating     int    `json:"attitude_rating" binding:"required"`
		PunctualRating     int    `json:"punctual_rating" binding:"required"`
		ProfessionalRating int    `json:"professional_rating" binding:"required"`
		Content            string `json:"content"`
		Images             string `json:"images"`
		IsAnonymous        bool   `json:"is_anonymous"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if !utils.ValidateRating(req.AttitudeRating) || !utils.ValidateRating(req.PunctualRating) || !utils.ValidateRating(req.ProfessionalRating) {
		utils.BadRequest(c, "评分必须在1-5之间")
		return
	}

	var order models.Order
	if result := config.DB.First(&order, req.OrderID); result.Error != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	if order.Status != models.OrderStatusCompleted {
		utils.BadRequest(c, "订单未完成，无法评价")
		return
	}

	var existingReview models.Review
	config.DB.Where("order_id = ? AND reviewer_id = ?", req.OrderID, userID).First(&existingReview)
	if existingReview.ID != 0 {
		utils.BadRequest(c, "您已评价过此订单")
		return
	}

	if role == "customer" {
		if order.CustomerID != userID {
			utils.BadRequest(c, "无权评价此订单")
			return
		}
		if order.CustomerRated {
			utils.BadRequest(c, "您已评价过此订单")
			return
		}
	} else if role == "service_provider" {
		if order.ProviderID != userID {
			utils.BadRequest(c, "无权评价此订单")
			return
		}
		if order.ProviderRated {
			utils.BadRequest(c, "您已评价过此订单")
			return
		}
	}

	overallRating := float64(req.AttitudeRating+req.PunctualRating+req.ProfessionalRating) / 3.0
	overallRating = utils.RoundPrice(overallRating, 1)

	var revieweeID uint
	var reviewType models.ReviewType

	if role == "customer" {
		revieweeID = order.ProviderID
		reviewType = models.ReviewTypeCustomer
	} else {
		revieweeID = order.CustomerID
		reviewType = models.ReviewTypeProvider
	}

	review := models.Review{
		OrderID:            req.OrderID,
		ReviewerID:         userID,
		RevieweeID:         revieweeID,
		ReviewType:         reviewType,
		AttitudeRating:     req.AttitudeRating,
		PunctualRating:     req.PunctualRating,
		ProfessionalRating: req.ProfessionalRating,
		OverallRating:      overallRating,
		Content:            req.Content,
		Images:             req.Images,
		IsAnonymous:        req.IsAnonymous,
	}

	tx := config.DB.Begin()

	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "评价失败")
		return
	}

	if role == "customer" {
		tx.Model(&order).Update("customer_rated", true)
		var provider models.User
		tx.First(&provider, order.ProviderID)
		newRating := (provider.Rating*float64(provider.OrderCount) + overallRating) / float64(provider.OrderCount+1)
		tx.Model(&provider).Update("rating", utils.RoundPrice(newRating, 1))
	} else {
		tx.Model(&order).Update("provider_rated", true)
	}

	tx.Model(&models.ServiceItem{}).Where("id = ?", order.ServiceItemID).
		Updates(map[string]interface{}{
			"review_count": gorm.Expr("review_count + ?", 1),
			"rating":       overallRating,
		})

	tx.Commit()

	go utils.LogOperation(userID, role, "review", "create", &review.ID, "review", "创建评价", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, review)
}

func GetReviewList(c *gin.Context) {
	serviceItemID := c.Query("service_item_id")
	providerID := c.Query("provider_id")
	minRating := c.Query("min_rating")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := config.DB.Model(&models.Review{}).Where("review_type = ?", models.ReviewTypeCustomer)

	if serviceItemID != "" {
		var orderIDs []uint
		config.DB.Model(&models.Order{}).Where("service_item_id = ?", serviceItemID).Pluck("id", &orderIDs)
		query = query.Where("order_id IN ?", orderIDs)
	}

	if providerID != "" {
		query = query.Where("reviewee_id = ?", providerID)
	}

	if minRating != "" {
		query = query.Where("overall_rating >= ?", minRating)
	}

	var total int64
	query.Count(&total)

	var reviews []models.Review
	offset := (page - 1) * pageSize
	query.Preload("Reviewer").Preload("Reviewee").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&reviews)

	utils.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"list":      reviews,
	})
}

func GetReviewDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var review models.Review
	if result := config.DB.Preload("Reviewer").Preload("Reviewee").First(&review, id); result.Error != nil {
		utils.NotFound(c, "评价不存在")
		return
	}

	utils.Success(c, review)
}

func ReplyReview(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var review models.Review
	if result := config.DB.First(&review, id); result.Error != nil {
		utils.NotFound(c, "评价不存在")
		return
	}

	if review.RevieweeID != userID {
		utils.Forbidden(c, "无权回复此评价")
		return
	}

	now := time.Now()
	review.ReplyContent = req.Content
	review.RepliedAt = &now

	config.DB.Save(&review)

	go utils.LogOperation(userID, "service_provider", "review", "reply", &review.ID, "review", "回复评价", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, nil)
}

func CreateComplaint(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	var req struct {
		OrderID uint   `json:"order_id" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Images  string `json:"images"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var order models.Order
	if result := config.DB.First(&order, req.OrderID); result.Error != nil {
		utils.NotFound(c, "订单不存在")
		return
	}

	var respondentID uint
	if role == "customer" {
		if order.CustomerID != userID {
			utils.BadRequest(c, "无权投诉此订单")
			return
		}
		respondentID = order.ProviderID
	} else if role == "service_provider" {
		if order.ProviderID != userID {
			utils.BadRequest(c, "无权投诉此订单")
			return
		}
		respondentID = order.CustomerID
	}

	complaint := models.Complaint{
		OrderID:       req.OrderID,
		ComplainantID: userID,
		RespondentID:  respondentID,
		Title:         req.Title,
		Content:       req.Content,
		Images:        req.Images,
		Status:        models.ComplaintStatusPending,
	}

	if result := config.DB.Create(&complaint); result.Error != nil {
		utils.InternalError(c, "投诉失败")
		return
	}

	var admins []models.User
	config.DB.Where("role = ?", models.RoleAdmin).Find(&admins)
	for _, admin := range admins {
		utils.SendSystemNotification(admin.ID, "新投诉待处理", "订单"+order.OrderNo+"有新投诉待处理")
	}

	go utils.LogOperation(userID, role, "complaint", "create", &complaint.ID, "complaint", "创建投诉", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, complaint)
}

func GetComplaintList(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	query := config.DB.Model(&models.Complaint{})

	if role == "customer" || role == "service_provider" {
		query = query.Where("complainant_id = ? OR respondent_id = ?", userID, userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var complaints []models.Complaint
	offset := (page - 1) * pageSize
	query.Preload("Complainant").Preload("Respondent").Preload("Order").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&complaints)

	utils.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"list":      complaints,
	})
}

func HandleComplaint(c *gin.Context) {
	handlerID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
		Result string `json:"result" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var complaint models.Complaint
	if result := config.DB.First(&complaint, id); result.Error != nil {
		utils.NotFound(c, "投诉不存在")
		return
	}

	complaint.Status = models.ComplaintStatus(req.Status)
	complaint.HandlerID = &handlerID
	complaint.HandleResult = req.Result
	now := time.Now()
	complaint.HandledAt = &now

	config.DB.Save(&complaint)

	utils.SendSystemNotification(complaint.ComplainantID, "投诉处理结果", "您的投诉已处理，结果："+req.Result)
	utils.SendSystemNotification(complaint.RespondentID, "投诉处理结果", "针对您的投诉已处理，结果："+req.Result)

	go utils.LogOperation(handlerID, "admin", "complaint", "handle", &complaint.ID, "complaint", "处理投诉", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, nil)
}
