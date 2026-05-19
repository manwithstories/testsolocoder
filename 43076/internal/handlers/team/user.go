package team

import (
	"strconv"

	"ticket-system/internal/database"
	"ticket-system/internal/middleware"
	"ticket-system/internal/models"
	"ticket-system/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateUserRequest struct {
	Username     string `json:"username" binding:"required,max=50"`
	Email        string `json:"email" binding:"required,email,max=100"`
	Password     string `json:"password" binding:"required,min=6,max=50"`
	RealName     string `json:"real_name" binding:"max=50"`
	Phone        string `json:"phone" binding:"max=20"`
	Role         string `json:"role" binding:"required,user_role"`
	DepartmentID *uint  `json:"department_id"`
	SkillGroupIDs []uint `json:"skill_group_ids"`
}

type UpdateUserRequest struct {
	Username     string `json:"username" binding:"max=50"`
	Email        string `json:"email" binding:"omitempty,email,max=100"`
	Password     string `json:"password" binding:"omitempty,min=6,max=50"`
	RealName     string `json:"real_name" binding:"max=50"`
	Phone        string `json:"phone" binding:"max=20"`
	Role         string `json:"role" binding:"omitempty,user_role"`
	DepartmentID *uint  `json:"department_id"`
	SkillGroupIDs []uint `json:"skill_group_ids"`
	IsOnDuty     *bool  `json:"is_on_duty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  *models.User `json:"user"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Unauthorized(c, "Invalid username or password")
			return
		}
		utils.InternalServerError(c, "Failed to login")
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		utils.Unauthorized(c, "Invalid username or password")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token")
		return
	}

	utils.Success(c, LoginResponse{
		Token: token,
		User:  &user,
	})
}

func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	var existingCount int64
	database.DB.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&existingCount)
	if existingCount > 0 {
		utils.BadRequest(c, "Username or email already exists")
		return
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerError(c, "Failed to hash password")
		return
	}

	tx := database.DB.Begin()

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		RealName:     req.RealName,
		Phone:        req.Phone,
		Role:         req.Role,
		DepartmentID: req.DepartmentID,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to create user")
		return
	}

	if len(req.SkillGroupIDs) > 0 {
		var groups []models.SkillGroup
		tx.Find(&groups, req.SkillGroupIDs)
		if err := tx.Model(user).Association("SkillGroups").Append(&groups); err != nil {
			tx.Rollback()
			utils.InternalServerError(c, "Failed to add skill groups")
			return
		}
	}

	tx.Commit()

	database.DB.Preload("Department").Preload("SkillGroups").First(user, user.ID)
	utils.Success(c, user)
}

func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	var user models.User
	if err := database.DB.Preload("Department").Preload("SkillGroups").First(&user, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "User not found")
			return
		}
		utils.InternalServerError(c, "Failed to get user")
		return
	}

	utils.Success(c, user)
}

func GetCurrentUser(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)

	var user models.User
	if err := database.DB.Preload("Department").Preload("SkillGroups").First(&user, userID).Error; err != nil {
		utils.InternalServerError(c, "Failed to get user")
		return
	}

	utils.Success(c, user)
}

func ListUsers(c *gin.Context) {
	var users []models.User
	query := database.DB.Preload("Department").Preload("SkillGroups")

	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}
	if deptID := c.Query("department_id"); deptID != "" {
		query = query.Where("department_id = ?", deptID)
	}
	if skillGroupID := c.Query("skill_group_id"); skillGroupID != "" {
		query = query.Joins("JOIN user_skill_groups ON user_skill_groups.user_id = users.id").
			Where("user_skill_groups.skill_group_id = ?", skillGroupID)
	}

	if err := query.Find(&users).Error; err != nil {
		utils.InternalServerError(c, "Failed to list users")
		return
	}

	utils.Success(c, users)
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	tx := database.DB.Begin()

	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.DepartmentID != nil {
		updates["department_id"] = req.DepartmentID
	}
	if req.IsOnDuty != nil {
		updates["is_on_duty"] = *req.IsOnDuty
	}
	if req.Password != "" {
		passwordHash, err := utils.HashPassword(req.Password)
		if err != nil {
			tx.Rollback()
			utils.InternalServerError(c, "Failed to hash password")
			return
		}
		updates["password_hash"] = passwordHash
	}

	if len(updates) > 0 {
		if err := tx.Model(&models.User{}).Where("id = ?", uint(id)).Updates(updates).Error; err != nil {
			tx.Rollback()
			utils.InternalServerError(c, "Failed to update user")
			return
		}
	}

	if req.SkillGroupIDs != nil {
		var user models.User
		tx.First(&user, uint(id))
		tx.Model(&user).Association("SkillGroups").Clear()
		var groups []models.SkillGroup
		tx.Find(&groups, req.SkillGroupIDs)
		if err := tx.Model(&user).Association("SkillGroups").Append(&groups); err != nil {
			tx.Rollback()
			utils.InternalServerError(c, "Failed to update skill groups")
			return
		}
	}

	tx.Commit()

	var user models.User
	database.DB.Preload("Department").Preload("SkillGroups").First(&user, uint(id))
	utils.Success(c, user)
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	currentUserID, _ := middleware.GetCurrentUserID(c)
	if uint(id) == currentUserID {
		utils.BadRequest(c, "Cannot delete yourself")
		return
	}

	tx := database.DB.Begin()

	var user models.User
	if err := tx.First(&user, uint(id)).Error; err != nil {
		tx.Rollback()
		utils.NotFound(c, "User not found")
		return
	}

	tx.Model(&user).Association("SkillGroups").Clear()

	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to delete user")
		return
	}

	tx.Commit()
	utils.Success(c, nil)
}
