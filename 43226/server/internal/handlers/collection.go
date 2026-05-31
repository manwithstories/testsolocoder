package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"museum-server/internal/dto"
	"museum-server/internal/models"
	"museum-server/internal/services"
	"museum-server/pkg/response"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	result, err := h.userService.Register(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	result, err := h.userService.Login(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	user, err := h.userService.GetUser(uid)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.userService.UpdateUser(uid, &req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	var query dto.CollectionListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	users, total, err := h.userService.ListUsers(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.SuccessWithPage(c, users, total, query.Page, query.PageSize)
}

func (h *UserHandler) GetGuides(c *gin.Context) {
	guides, err := h.userService.GetGuides()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, guides)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	response.Success(c, user)
}

type CollectionHandler struct {
	collectionService *services.CollectionService
}

func NewCollectionHandler(collectionService *services.CollectionService) *CollectionHandler {
	return &CollectionHandler{collectionService: collectionService}
}

func (h *CollectionHandler) Create(c *gin.Context) {
	var req dto.CollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	museumID := c.GetUint("museum_id")
	if museumID == 0 {
		museumID = 1
	}

	collection, err := h.collectionService.Create(museumID, &req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, collection)
}

func (h *CollectionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid collection ID")
		return
	}

	collection, err := h.collectionService.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}

	response.Success(c, collection)
}

func (h *CollectionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid collection ID")
		return
	}

	var req dto.CollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.collectionService.Update(uint(id), &req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *CollectionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid collection ID")
		return
	}

	if err := h.collectionService.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *CollectionHandler) List(c *gin.Context) {
	var query dto.CollectionListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	collections, total, err := h.collectionService.List(&query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.SuccessWithPage(c, collections, total, query.Page, query.PageSize)
}

func (h *CollectionHandler) CreateCategory(c *gin.Context) {
	var req dto.CollectionCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	category, err := h.collectionService.CreateCategory(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, category)
}

func (h *CollectionHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid category ID")
		return
	}

	var req dto.CollectionCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	if err := h.collectionService.UpdateCategory(uint(id), &req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *CollectionHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid category ID")
		return
	}

	if err := h.collectionService.DeleteCategory(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *CollectionHandler) ListCategories(c *gin.Context) {
	categories, err := h.collectionService.ListCategories()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, categories)
}

func (h *CollectionHandler) ListTags(c *gin.Context) {
	tags, err := h.collectionService.ListTags()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, tags)
}

func (h *CollectionHandler) CreateTag(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required,max=32"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	tag, err := h.collectionService.CreateTag(req.Name)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, tag)
}

func (h *CollectionHandler) BatchImport(c *gin.Context) {
	var collections []models.Collection
	if err := c.ShouldBindJSON(&collections); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	museumID := c.GetUint("museum_id")
	if museumID == 0 {
		museumID = 1
	}

	if err := h.collectionService.BatchImport(museumID, collections); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	response.Success(c, gin.H{"imported_count": len(collections)})
}
