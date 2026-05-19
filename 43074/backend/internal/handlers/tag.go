package handlers

import (
	"net/http"

	"booklibrary/internal/database"
	"booklibrary/internal/errors"
	"booklibrary/internal/logger"
	"booklibrary/internal/models"
	"booklibrary/internal/utils"

	"github.com/gin-gonic/gin"
)

type TagHandler struct{}

func NewTagHandler() *TagHandler {
	return &TagHandler{}
}

type CreateTagRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

func (h *TagHandler) GetTags(c *gin.Context) {
	var tags []models.Tag
	result := database.DB.Order("name asc").Find(&tags)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}
	c.JSON(http.StatusOK, tags)
}

func (h *TagHandler) GetTag(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var tag models.Tag
	result := database.DB.First(&tag, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrValidation)
		return
	}

	var existingTag models.Tag
	result := database.DB.Where("name = ?", req.Name).First(&existingTag)
	if result.Error == nil {
		errors.ErrorResponse(c, http.StatusConflict, "标签已存在")
		return
	}

	color := req.Color
	if color == "" {
		color = "#3b82f6"
	}

	tag := models.Tag{
		Name:  req.Name,
		Color: color,
	}

	result = database.DB.Create(&tag)
	if result.Error != nil {
		logger.Errorf("Create tag failed: %v", result.Error)
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}

	logger.Infof("Tag created: %d - %s", tag.ID, tag.Name)
	c.JSON(http.StatusCreated, tag)
}

func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var tag models.Tag
	result := database.DB.First(&tag, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrValidation)
		return
	}

	if req.Name != "" && req.Name != tag.Name {
		var existingTag models.Tag
		result := database.DB.Where("name = ? AND id != ?", req.Name, id).First(&existingTag)
		if result.Error == nil {
			errors.ErrorResponse(c, http.StatusConflict, "标签已存在")
			return
		}
		tag.Name = req.Name
	}

	if req.Color != "" {
		tag.Color = req.Color
	}

	database.DB.Save(&tag)
	logger.Infof("Tag updated: %d - %s", tag.ID, tag.Name)
	c.JSON(http.StatusOK, tag)
}

func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var tag models.Tag
	result := database.DB.First(&tag, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	tx := database.DB.Begin()
	tx.Model(&tag).Association("Books").Clear()
	tx.Delete(&tag)
	tx.Commit()

	logger.Infof("Tag deleted: %d", id)
	c.JSON(http.StatusNoContent, nil)
}
