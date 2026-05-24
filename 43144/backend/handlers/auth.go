package handlers

import (
	"net/http"
	"strconv"

	"pet-adoption-platform/models"
	"pet-adoption-platform/services"
	"pet-adoption-platform/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := services.Register(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	response, err := services.Login(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, response)
}

func GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := services.GetUserByID(userID)
	if err != nil {
		utils.InternalError(c, "failed to get user profile")
		return
	}

	utils.Success(c, user)
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	allowedFields := []string{"name", "phone", "address"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := req[field]; ok {
			updates[field] = val
		}
	}

	user, err := services.UpdateUser(userID, updates)
	if err != nil {
		utils.InternalError(c, "failed to update profile")
		return
	}

	utils.Success(c, user)
}

func ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	role := c.Query("role")

	users, total, err := services.ListUsers(page, pageSize, role)
	if err != nil {
		utils.InternalError(c, "failed to list users")
		return
	}

	utils.PaginatedSuccess(c, users, total, page, pageSize)
}

func GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}
	user, err := services.GetUserByID(uint(id))
	if err != nil {
		utils.NotFound(c, "user not found")
		return
	}

	utils.Success(c, user)
}

func VerifyUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}
	user, err := services.UpdateUser(uint(id), map[string]interface{}{
		"is_verified": true,
	})
	if err != nil {
		utils.InternalError(c, "failed to verify user")
		return
	}

	utils.Success(c, user)
}

func GetHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
