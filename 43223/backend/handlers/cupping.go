package handlers

import (
	"net/http"
	"strconv"
	"time"

	"coffee-platform/config"
	"coffee-platform/database"
	"coffee-platform/models"
	"coffee-platform/utils"

	"github.com/gin-gonic/gin"
)

type CuppingHandler struct {
	cfg *config.Config
}

func NewCuppingHandler(cfg *config.Config) *CuppingHandler {
	return &CuppingHandler{cfg: cfg}
}

func calculateOverallScore(req models.CreateCuppingScoreRequest) float64 {
	total := req.DryFragrance + req.WetAroma + req.Body + req.Acidity +
		req.Sweetness + req.Aftertaste + req.Balance
	return total / 7
}

func (h *CuppingHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(h.cfg.App.PageSize)))
	productID := c.Query("product_id")
	userID := c.Query("user_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = h.cfg.App.PageSize
	}

	query := database.DB.Model(&models.CuppingScore{}).Preload("Product").Preload("User")

	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	var total int64
	query.Count(&total)

	var scores []models.CuppingScore
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&scores)

	utils.PaginatedResponse(c, scores, total, page, pageSize)
}

func (h *CuppingHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var score models.CuppingScore
	if err := database.DB.Preload("CriteriaItems").Preload("Product").Preload("User").
		First(&score, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "评分记录不存在")
		return
	}

	utils.Success(c, score)
}

func (h *CuppingHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.CreateCuppingScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	var existingScore models.CuppingScore
	if database.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).
		First(&existingScore).Error == nil {
		utils.Error(c, http.StatusConflict, "您已对该商品进行过杯测评分")
		return
	}

	var product models.Product
	if err := database.DB.First(&product, req.ProductID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "商品不存在")
		return
	}

	overallScore := calculateOverallScore(req)

	score := models.CuppingScore{
		ProductID:    req.ProductID,
		UserID:       userID,
		DryFragrance: req.DryFragrance,
		WetAroma:     req.WetAroma,
		Body:         req.Body,
		Acidity:      req.Acidity,
		Sweetness:    req.Sweetness,
		Aftertaste:   req.Aftertaste,
		Balance:      req.Balance,
		OverallScore: overallScore,
		Notes:        req.Notes,
	}

	if err := database.DB.Create(&score).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "评分失败")
		return
	}

	var avgScore float64
	database.DB.Model(&models.CuppingScore{}).
		Where("product_id = ?", req.ProductID).
		Select("COALESCE(AVG(overall_score), 0)").
		Scan(&avgScore)

	database.DB.Model(&models.Product{}).Where("id = ?", req.ProductID).
		Update("cupping_score", avgScore)

	utils.SuccessWithMessage(c, "评分成功", score)
}

func (h *CuppingHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var score models.CuppingScore
	if err := database.DB.First(&score, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "评分记录不存在")
		return
	}

	if score.UserID != userID {
		utils.Error(c, http.StatusForbidden, "无权修改此评分")
		return
	}

	var req models.CreateCuppingScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	overallScore := calculateOverallScore(req)

	updates := map[string]interface{}{
		"dry_fragrance": req.DryFragrance,
		"wet_aroma":     req.WetAroma,
		"body":          req.Body,
		"acidity":       req.Acidity,
		"sweetness":     req.Sweetness,
		"aftertaste":    req.Aftertaste,
		"balance":       req.Balance,
		"overall_score": overallScore,
		"notes":         req.Notes,
	}

	if err := database.DB.Model(&score).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	var avgScore float64
	database.DB.Model(&models.CuppingScore{}).
		Where("product_id = ?", score.ProductID).
		Select("COALESCE(AVG(overall_score), 0)").
		Scan(&avgScore)

	database.DB.Model(&models.Product{}).Where("id = ?", score.ProductID).
		Update("cupping_score", avgScore)

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (h *CuppingHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var score models.CuppingScore
	if err := database.DB.First(&score, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "评分记录不存在")
		return
	}

	if score.UserID != userID && userRole != string(models.RoleAdmin) {
		utils.Error(c, http.StatusForbidden, "无权删除此评分")
		return
	}

	productID := score.ProductID

	if err := database.DB.Delete(&score).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	var avgScore float64
	database.DB.Model(&models.CuppingScore{}).
		Where("product_id = ?", productID).
		Select("COALESCE(AVG(overall_score), 0)").
		Scan(&avgScore)

	database.DB.Model(&models.Product{}).Where("id = ?", productID).
		Update("cupping_score", avgScore)

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *CuppingHandler) GetProductStats(c *gin.Context) {
	productID := c.Query("product_id")

	if productID == "" {
		utils.Error(c, http.StatusBadRequest, "请指定商品ID")
		return
	}

	var stats models.CuppingStats
	database.DB.Table("cupping_scores").
		Select(`product_id,
			COUNT(*) as total_count,
			COALESCE(AVG(overall_score), 0) as avg_score,
			COALESCE(AVG(dry_fragrance), 0) as avg_dry_fragrance,
			COALESCE(AVG(wet_aroma), 0) as avg_wet_aroma,
			COALESCE(AVG(body), 0) as avg_body,
			COALESCE(AVG(acidity), 0) as avg_acidity,
			COALESCE(AVG(sweetness), 0) as avg_sweetness,
			COALESCE(AVG(aftertaste), 0) as avg_aftertaste,
			COALESCE(AVG(balance), 0) as avg_balance`).
		Where("product_id = ?", productID).
		Group("product_id").
		Scan(&stats)

	var product models.Product
	database.DB.First(&product, productID)
	stats.ProductName = product.Name

	utils.Success(c, stats)
}

func (h *CuppingHandler) GetScoreTrend(c *gin.Context) {
	productID := c.Query("product_id")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	if days < 1 || days > 365 {
		days = 30
	}

	since := time.Now().AddDate(0, 0, -days)

	var results []models.ScoreTrendItem
	query := database.DB.Table("cupping_scores").
		Select("DATE(created_at) as date, COALESCE(AVG(overall_score), 0) as avg_score, COUNT(*) as count")

	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}

	query.Where("created_at >= ?", since).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results)

	utils.Success(c, results)
}

func (h *CuppingHandler) GetUserScores(c *gin.Context) {
	userID := c.GetUint("user_id")

	var scores []models.CuppingScore
	database.DB.Preload("Product").Where("user_id = ?", userID).
		Order("created_at DESC").Find(&scores)

	utils.Success(c, scores)
}
