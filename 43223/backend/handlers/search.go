package handlers

import (
	"net/http"
	"strconv"

	"coffee-platform/config"
	"coffee-platform/database"
	"coffee-platform/models"
	"coffee-platform/utils"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	cfg *config.Config
}

func NewSearchHandler(cfg *config.Config) *SearchHandler {
	return &SearchHandler{cfg: cfg}
}

func (h *SearchHandler) SearchProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(h.cfg.App.PageSize)))

	keyword := c.Query("q")
	origin := c.Query("origin")
	roastLevel := c.Query("roast_level")
	processMethod := c.Query("process_method")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	minScore := c.Query("min_score")
	maxScore := c.Query("max_score")
	roasterID := c.Query("roaster_id")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = h.cfg.App.PageSize
	}

	query := database.DB.Model(&models.Product{}).Preload("Images").Preload("Roaster").
		Where("status = ?", models.ProductStatusOnSale)

	if keyword != "" {
		query = query.Where("name ILIKE ? OR origin ILIKE ? OR farm ILIKE ? OR variety ILIKE ? OR flavor_notes ILIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
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
	if minScore != "" {
		query = query.Where("cupping_score >= ?", minScore)
	}
	if maxScore != "" {
		query = query.Where("cupping_score <= ?", maxScore)
	}
	if roasterID != "" {
		query = query.Where("roaster_id = ?", roasterID)
	}

	allowedSortFields := map[string]bool{
		"created_at": true, "price": true, "cupping_score": true, "name": true, "origin": true,
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

	type SearchFilters struct {
		Origins       []string `json:"origins"`
		RoastLevels   []string `json:"roast_levels"`
		ProcessMethods []string `json:"process_methods"`
		PriceRange    struct {
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"price_range"`
		ScoreRange    struct {
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"score_range"`
	}

	var filters SearchFilters
	database.DB.Model(&models.Product{}).Where("status = ?", models.ProductStatusOnSale).
		Distinct("origin").Pluck("origin", &filters.Origins)

	filters.RoastLevels = []string{"light", "medium", "medium_dark", "dark"}
	filters.ProcessMethods = []string{"washed", "natural", "honey", "anaerobic", "wet_hulled"}

	var priceStats struct {
		MinPrice float64 `json:"min_price"`
		MaxPrice float64 `json:"max_price"`
	}
	database.DB.Model(&models.Product{}).Where("status = ?", models.ProductStatusOnSale).
		Select("COALESCE(MIN(price), 0) as min_price, COALESCE(MAX(price), 0) as max_price").
		Scan(&priceStats)
	filters.PriceRange.Min = priceStats.MinPrice
	filters.PriceRange.Max = priceStats.MaxPrice

	var scoreStats struct {
		MinScore float64 `json:"min_score"`
		MaxScore float64 `json:"max_score"`
	}
	database.DB.Model(&models.Product{}).Where("status = ?", models.ProductStatusOnSale).
		Select("COALESCE(MIN(cupping_score), 0) as min_score, COALESCE(MAX(cupping_score), 0) as max_score").
		Scan(&scoreStats)
	filters.ScoreRange.Min = scoreStats.MinScore
	filters.ScoreRange.Max = scoreStats.MaxScore

	utils.Success(c, gin.H{
		"items":     products,
		"total":     total,
		"page":      page,
		"pageSize":  pageSize,
		"filters":   filters,
	})
}

func (h *SearchHandler) Suggest(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		utils.Success(c, []interface{}{})
		return
	}

	type SuggestItem struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		Image string `json:"image"`
	}

	var suggestions []SuggestItem

	var products []models.Product
	database.DB.Preload("Images").
		Where("name ILIKE ? AND status = ?", "%"+keyword+"%", models.ProductStatusOnSale).
		Limit(5).Find(&products)

	for _, p := range products {
		item := SuggestItem{
			ID:   p.ID,
			Name: p.Name,
			Type: "product",
		}
		for _, img := range p.Images {
			if img.IsCover {
				item.Image = img.URL
				break
			}
		}
		suggestions = append(suggestions, item)
	}

	var origins []string
	database.DB.Model(&models.Product{}).
		Where("origin ILIKE ? AND status = ?", "%"+keyword+"%", models.ProductStatusOnSale).
		Distinct("origin").Limit(5).Pluck("origin", &origins)

	for _, o := range origins {
		suggestions = append(suggestions, SuggestItem{
			Name: o,
			Type: "origin",
		})
	}

	utils.Success(c, suggestions)
}

func (h *SearchHandler) GetSearchHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	type HistoryItem struct {
		Keyword   string `json:"keyword"`
		SearchCount int64 `json:"search_count"`
		LastSearch string `json:"last_search"`
	}

	var history []HistoryItem
	database.DB.Table("operation_logs").
		Select("path as keyword, COUNT(*) as search_count, MAX(created_at) as last_search").
		Where("user_id = ? AND path LIKE ?", userID, "%/search%").
		Group("path").
		Order("last_search DESC").
		Limit(10).
		Scan(&history)

	utils.Success(c, history)
}

func (h *SearchHandler) AdvancedSearch(c *gin.Context) {
	var req struct {
		Keywords      []string          `json:"keywords"`
		Filters       map[string]string `json:"filters"`
		PriceRange    [2]float64        `json:"price_range"`
		ScoreRange    [2]float64        `json:"score_range"`
		ExcludeFields []string          `json:"exclude_fields"`
		SortBy        string            `json:"sort_by"`
		SortOrder     string            `json:"sort_order"`
		Page          int               `json:"page"`
		PageSize      int               `json:"page_size"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = h.cfg.App.PageSize
	}

	query := database.DB.Model(&models.Product{}).Preload("Images").Preload("Roaster").
		Where("status = ?", models.ProductStatusOnSale)

	for _, kw := range req.Keywords {
		query = query.Where("name ILIKE ? OR origin ILIKE ? OR flavor_notes ILIKE ?",
			"%"+kw+"%", "%"+kw+"%", "%"+kw+"%")
	}

	for field, value := range req.Filters {
		query = query.Where(field+" = ?", value)
	}

	if req.PriceRange[0] > 0 {
		query = query.Where("price >= ?", req.PriceRange[0])
	}
	if req.PriceRange[1] > 0 {
		query = query.Where("price <= ?", req.PriceRange[1])
	}
	if req.ScoreRange[0] > 0 {
		query = query.Where("cupping_score >= ?", req.ScoreRange[0])
	}
	if req.ScoreRange[1] > 0 {
		query = query.Where("cupping_score <= ?", req.ScoreRange[1])
	}

	allowedSortFields := map[string]bool{
		"created_at": true, "price": true, "cupping_score": true, "name": true,
	}
	if allowedSortFields[req.SortBy] {
		order := req.SortBy + " " + req.SortOrder
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	var total int64
	query.Count(&total)

	var products []models.Product
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&products)

	utils.Success(c, gin.H{
		"items":    products,
		"total":    total,
		"page":     req.Page,
		"pageSize": req.PageSize,
	})
}
