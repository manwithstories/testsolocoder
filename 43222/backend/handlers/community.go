package handlers

import (
	"garden-planner/database"
	"garden-planner/middleware"
	"garden-planner/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreatePostRequest struct {
	Title     string `json:"title" binding:"required,max=200"`
	Content   string `json:"content" binding:"required"`
	ImageURLs string `json:"image_urls"`
	Category  string `json:"category"`
	Tags      string `json:"tags"`
}

type UpdatePostRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	ImageURLs string `json:"image_urls"`
	Category  string `json:"category"`
	Tags      string `json:"tags"`
}

type CreateCommentRequest struct {
	Content  string     `json:"content" binding:"required"`
	ParentID *uuid.UUID `json:"parent_id"`
}

func CreatePost(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{
		ID:         uuid.New(),
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		ImageURLs: req.ImageURLs,
		Category:   req.Category,
		Tags:       req.Tags,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post":    post,
	})
}

func GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	category := c.Query("category")
	userID := c.Query("user_id")
	search := c.Query("search")

	var posts []models.Post
	var total int64

	query := database.DB.Model(&models.Post{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	query.Preload("User").
		Preload("Comments.User").
		Offset(offset).Limit(pageSize).
		Order("CASE WHEN is_pinned THEN 0 ELSE 1 END, created_at DESC").
		Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"posts":     posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetPost(c *gin.Context) {
	id := c.Param("id")

	var post models.Post
	if err := database.DB.Preload("User").
		Preload("Comments.User").
		First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	post.ViewCount++
	database.DB.Model(&post).Update("view_count", post.ViewCount)

	c.JSON(http.StatusOK, gin.H{"post": post})
}

func UpdatePost(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var post models.Post
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.ImageURLs != "" {
		post.ImageURLs = req.ImageURLs
	}
	if req.Category != "" {
		post.Category = req.Category
	}
	if req.Tags != "" {
		post.Tags = req.Tags
	}

	if err := database.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"post":    post,
	})
}

func DeletePost(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Post{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func LikePost(c *gin.Context) {
	userID := middleware.GetUserID(c)
	postID := c.Param("id")

	postUUID, err := uuid.Parse(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var existingLike models.Like
	result := database.DB.Where("post_id = ? AND user_id = ?", postUUID, userID).First(&existingLike)

	if result.Error == nil {
		database.DB.Delete(&existingLike)
		database.DB.Model(&models.Post{}).Where("id = ?", postUUID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))
		c.JSON(http.StatusOK, gin.H{"message": "Post unliked", "liked": false})
		return
	}

	like := models.Like{
		ID:     uuid.New(),
		PostID: &postUUID,
		UserID: userID,
	}

	if err := database.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}

	database.DB.Model(&models.Post{}).Where("id = ?", postUUID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))

	c.JSON(http.StatusOK, gin.H{"message": "Post liked", "liked": true})
}

func CreateComment(c *gin.Context) {
	userID := middleware.GetUserID(c)
	postID := c.Param("id")

	postUUID, err := uuid.Parse(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := models.Comment{
		ID:       uuid.New(),
		PostID:   postUUID,
		UserID:   userID,
		Content:  req.Content,
		ParentID: req.ParentID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

func GetComments(c *gin.Context) {
	postID := c.Param("id")

	var comments []models.Comment
	database.DB.Where("post_id = ?", postID).
		Preload("User").
		Order("created_at ASC").
		Find(&comments)

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func DeleteComment(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Comment{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

func FollowUser(c *gin.Context) {
	followerID := middleware.GetUserID(c)
	followingID := c.Param("id")

	followingUUID, err := uuid.Parse(followingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if followerID == followingUUID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}

	var existingFollow models.Follow
	result := database.DB.Where("follower_id = ? AND following_id = ?", followerID, followingUUID).First(&existingFollow)

	if result.Error == nil {
		database.DB.Delete(&existingFollow)
		c.JSON(http.StatusOK, gin.H{"message": "Unfollowed", "following": false})
		return
	}

	follow := models.Follow{
		ID:          uuid.New(),
		FollowerID:  followerID,
		FollowingID: followingUUID,
	}

	if err := database.DB.Create(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Followed", "following": true})
}

func GetFollowers(c *gin.Context) {
	userID := c.Param("id")

	var followers []models.User
	database.DB.Joins("JOIN follows ON follows.follower_id = users.id").
		Where("follows.following_id = ?", userID).
		Find(&followers)

	c.JSON(http.StatusOK, gin.H{"followers": followers})
}

func GetFollowing(c *gin.Context) {
	userID := c.Param("id")

	var following []models.User
	database.DB.Joins("JOIN follows ON follows.following_id = users.id").
		Where("follows.follower_id = ?", userID).
		Find(&following)

	c.JSON(http.StatusOK, gin.H{"following": following})
}
