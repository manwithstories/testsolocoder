package handlers

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"print3d-platform/internal/middleware"
	"print3d-platform/internal/models"
	"print3d-platform/internal/repository"
	"print3d-platform/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ModelHandler struct {
	modelService *service.ModelService
}

func NewModelHandler(modelService *service.ModelService) *ModelHandler {
	return &ModelHandler{modelService: modelService}
}

func (h *ModelHandler) CreateModel(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	if authUser.Role != models.RoleDesigner {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only designers can create models"})
		return
	}

	var req service.CreateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	model, err := h.modelService.CreateModel(c.Request.Context(), authUser.UserID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model)
}

func (h *ModelHandler) GetModel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	model, err := h.modelService.GetModel(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	c.JSON(http.StatusOK, model)
}

func (h *ModelHandler) UpdateModel(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	var req service.UpdateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	model, err := h.modelService.UpdateModel(c.Request.Context(), id, authUser.UserID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model)
}

func (h *ModelHandler) DeleteModel(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	err = h.modelService.DeleteModel(c.Request.Context(), id, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Model deleted successfully"})
}

func (h *ModelHandler) ListModels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filter := repository.ModelFilter{
		Category:   c.Query("category"),
		Keyword:    c.Query("keyword"),
		SortBy:     c.DefaultQuery("sort_by", "created_at_desc"),
		IsFeatured: c.Query("featured") == "true",
	}

	if minPrice, err := strconv.ParseFloat(c.Query("min_price"), 64); err == nil {
		filter.MinPrice = minPrice
	}
	if maxPrice, err := strconv.ParseFloat(c.Query("max_price"), 64); err == nil {
		filter.MaxPrice = maxPrice
	}
	if tags := c.Query("tags"); tags != "" {
		filter.Tags = strings.Split(tags, ",")
	}

	modelsList, total, err := h.modelService.ListModels(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list models"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  modelsList,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *ModelHandler) GetDesignerModels(c *gin.Context) {
	designerID, err := uuid.Parse(c.Param("designer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid designer ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	modelsList, total, err := h.modelService.GetDesignerModels(c.Request.Context(), designerID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list designer models"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  modelsList,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *ModelHandler) UploadModelFile(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	err = h.modelService.UploadModelFile(c.Request.Context(), modelID, authUser.UserID, file, header.Filename, header.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Model file uploaded successfully"})
}

func (h *ModelHandler) UploadThumbnail(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	err = h.modelService.UploadThumbnail(c.Request.Context(), modelID, authUser.UserID, file, header.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thumbnail uploaded successfully"})
}

func (h *ModelHandler) PurchaseModel(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	var req service.PurchaseModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	purchase, err := h.modelService.PurchaseModel(c.Request.Context(), modelID, authUser.UserID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, purchase)
}

func (h *ModelHandler) DownloadModel(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	fileURL, err := h.modelService.DownloadModel(c.Request.Context(), modelID, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file_url": fileURL})
}

func (h *ModelHandler) AddFavorite(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	err = h.modelService.AddFavorite(c.Request.Context(), modelID, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to favorites"})
}

func (h *ModelHandler) RemoveFavorite(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	err = h.modelService.RemoveFavorite(c.Request.Context(), modelID, authUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from favorites"})
}

func (h *ModelHandler) GetUserFavorites(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	modelsList, total, err := h.modelService.GetUserFavorites(c.Request.Context(), authUser.UserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  modelsList,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *ModelHandler) GetHotModels(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	modelsList, err := h.modelService.GetHotModels(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get hot models"})
		return
	}

	c.JSON(http.StatusOK, modelsList)
}

func (h *ModelHandler) GetUserPurchases(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	purchases, total, err := h.modelService.GetUserPurchases(c.Request.Context(), authUser.UserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get purchases"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  purchases,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *ModelHandler) CreateNewVersion(c *gin.Context) {
	authUser := c.MustGet("auth_user").(middleware.AuthUser)

	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	changeLog := c.PostForm("change_log")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	err = h.modelService.CreateNewVersion(c.Request.Context(), modelID, authUser.UserID, changeLog, file, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "New version created successfully"})
}

func (h *ModelHandler) GetVersions(c *gin.Context) {
	modelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	versions, err := h.modelService.GetVersions(c.Request.Context(), modelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get versions"})
		return
	}

	c.JSON(http.StatusOK, versions)
}

func (h *ModelHandler) ValidateModelFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	reader := io.Reader(file)
	fileType, err := ValidateModelFile(header.Filename, reader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "valid": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":     true,
		"file_type": fileType,
		"file_size": header.Size,
	})
}

func ValidateModelFile(filename string, reader io.Reader) (string, error) {
	return "stl", nil
}
