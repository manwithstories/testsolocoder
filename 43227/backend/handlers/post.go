package handlers

import (
	"net/http"
	"strconv"

	"beehive-platform/database"
	"beehive-platform/models"
	"beehive-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostHandler struct{}

func NewPostHandler() *PostHandler {
	return &PostHandler{}
}

func (h *PostHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	post := models.Post{
		UserID:   userID.(uint),
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Tags:     req.Tags,
		Images:   req.Images,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		utils.FailWithError(c, http.StatusInternalServerError, "failed to create post", err)
		return
	}

	utils.Success(c, post)
}

func (h *PostHandler) List(c *gin.Context) {
	var pageParams utils.PageParams
	if err := c.ShouldBindQuery(&pageParams); err != nil {
		pageParams = utils.PageParams{Page: 1, PageSize: 10}
	}
	if pageParams.Page < 1 {
		pageParams.Page = 1
	}
	if pageParams.PageSize < 1 {
		pageParams.PageSize = 10
	}

	query := database.DB.Model(&models.Post{})

	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if tag := c.Query("tag"); tag != "" {
		query = query.Where("tags::text LIKE ?", "%"+tag+"%")
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	query = query.Order(sortBy + " " + sortOrder)

	var total int64
	query.Count(&total)

	var posts []models.Post
	query.Offset(pageParams.GetOffset()).Limit(pageParams.PageSize).
		Preload("User").Find(&posts)

	utils.SuccessWithTotal(c, posts, total)
}

func (h *PostHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var post models.Post
	if err := database.DB.Preload("User").First(&post, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "post not found")
		return
	}

	database.DB.Model(&post).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	utils.Success(c, post)
}

func (h *PostHandler) Update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Title    *string   `json:"title" binding:"omitempty,max=200"`
		Content  *string   `json:"content"`
		Category *string   `json:"category" binding:"omitempty,oneof=disease_control harvest_technique seasonal_management equipment market general"`
		Tags     *[]string `json:"tags"`
		Images   *[]string `json:"images"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var post models.Post
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&post).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "post not found")
		return
	}

	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Category != nil {
		post.Category = *req.Category
	}
	if req.Tags != nil {
		post.Tags = *req.Tags
	}
	if req.Images != nil {
		post.Images = *req.Images
	}

	database.DB.Save(&post)

	utils.Success(c, post)
}

func (h *PostHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Post{})
	if result.Error != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to delete post")
		return
	}
	if result.RowsAffected == 0 {
		utils.Fail(c, http.StatusNotFound, "post not found")
		return
	}

	utils.Success(c, nil)
}

func (h *PostHandler) Like(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "post not found")
		return
	}

	database.DB.Model(&post).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))

	utils.Success(c, nil)
}

func (h *PostHandler) CreateComment(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var post models.Post
	if err := database.DB.First(&post, req.PostID).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "post not found")
		return
	}

	comment := models.Comment{
		PostID:   req.PostID,
		UserID:   userID.(uint),
		Content:  req.Content,
		ParentID: req.ParentID,
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to begin transaction")
		return
	}

	if err := tx.Create(&comment).Error; err != nil {
		tx.Rollback()
		utils.FailWithError(c, http.StatusInternalServerError, "failed to create comment", err)
		return
	}

	tx.Model(&post).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))

	tx.Commit()

	utils.Success(c, comment)
}

func (h *PostHandler) ListComments(c *gin.Context) {
	postID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var comments []models.Comment
	database.DB.Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("User").
		Order("created_at DESC").
		Find(&comments)

	utils.Success(c, comments)
}

func (h *PostHandler) GetCategories(c *gin.Context) {
	categories := []map[string]string{
		{"key": "disease_control", "name": "病虫害防治"},
		{"key": "harvest_technique", "name": "采收技巧"},
		{"key": "seasonal_management", "name": "季节性管理"},
		{"key": "equipment", "name": "设备工具"},
		{"key": "market", "name": "市场行情"},
		{"key": "general", "name": "综合讨论"},
	}

	utils.Success(c, categories)
}
