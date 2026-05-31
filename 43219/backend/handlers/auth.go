package handlers

import (
	"errors"
	"time"

	"housekeeping/database"
	"housekeeping/models"
	"housekeeping/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username    string      `json:"username" binding:"required,min=3,max=64"`
	Password    string      `json:"password" binding:"required,min=6"`
	Role        models.Role `json:"role" binding:"required"`
	RealName    string      `json:"real_name"`
	Phone       string      `json:"phone"`
	CompanyID   *uint       `json:"company_id,omitempty"`
	CertFiles   string      `json:"cert_files"`
	HealthFiles string      `json:"health_files"`
	Skills      string      `json:"skills"`
	Intro       string      `json:"intro"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func hashPassword(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func verifyPassword(hashed, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pwd)) == nil
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if req.Role != models.RoleCompany && req.Role != models.RoleStaff && req.Role != models.RoleCustomer {
		utils.BadRequest(c, "invalid role")
		return
	}
	var existing models.User
	if database.DB.Where("username = ?", req.Username).First(&existing).Error == nil {
		utils.BadRequest(c, "username already exists")
		return
	}
	hp, err := hashPassword(req.Password)
	if err != nil {
		utils.ServerError(c, "hash password failed")
		return
	}
	user := models.User{
		Username: req.Username,
		Password: hp,
		RealName: req.RealName,
		Phone:    req.Phone,
		Role:     req.Role,
		Status:   "active",
		Rating:   5,
		Level:    1,
	}
	if req.Role == models.RoleStaff {
		user.CompanyID = req.CompanyID
	}
	tx := database.DB.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		utils.ServerError(c, "create user failed")
		return
	}
	if req.Role == models.RoleStaff {
		sp := models.StaffProfile{
			UserID:      user.ID,
			CertFiles:   req.CertFiles,
			HealthFiles: req.HealthFiles,
			Intro:       req.Intro,
		}
		if err := tx.Create(&sp).Error; err != nil {
			tx.Rollback()
			utils.ServerError(c, "create staff profile failed")
			return
		}
	}
	if req.Role == models.RoleStaff || req.Role == models.RoleCompany {
		if err := tx.Create(&models.Wallet{UserID: user.ID}).Error; err != nil {
			tx.Rollback()
			utils.ServerError(c, "create wallet failed")
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		utils.ServerError(c, "commit failed")
		return
	}
	utils.Logger.Infow("user registered", "id", user.ID, "role", user.Role)
	utils.OK(c, gin.H{"id": user.ID, "username": user.Username, "role": user.Role})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.BadRequest(c, "invalid username or password")
		return
	}
	if user.Status != "active" {
		utils.Forbidden(c, "account is disabled")
		return
	}
	if !verifyPassword(user.Password, req.Password) {
		utils.BadRequest(c, "invalid username or password")
		return
	}
	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		utils.ServerError(c, "generate token failed")
		return
	}
	utils.OK(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"role":       user.Role,
			"real_name":  user.RealName,
			"phone":      user.Phone,
			"expires_at": time.Now().Add(24 * time.Hour).Unix(),
		},
	})
}

func Me(c *gin.Context) {
	uid, _ := c.Get("uid")
	var user models.User
	if err := database.DB.First(&user, uid).Error; err != nil {
		utils.NotFound(c, "user not found")
		return
	}
	out := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"status":   user.Status,
	}
	if user.Role == models.RoleStaff {
		var sp models.StaffProfile
		database.DB.Where("user_id = ?", user.ID).First(&sp)
		out["profile"] = sp
		out["rating"] = user.Rating
		out["level"] = user.Level
	}
	utils.OK(c, out)
}

func UpdateMe(c *gin.Context) {
	uid, _ := c.Get("uid")
	var body struct {
		RealName string `json:"real_name"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
		Skills   string `json:"skills"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	updates := map[string]interface{}{}
	if body.RealName != "" {
		updates["real_name"] = body.RealName
	}
	if body.Phone != "" {
		updates["phone"] = body.Phone
	}
	if body.Avatar != "" {
		updates["avatar"] = body.Avatar
	}
	if body.Skills != "" {
		updates["skills"] = body.Skills
	}
	if err := database.DB.Model(&models.User{}).Where("id = ?", uid).Updates(updates).Error; err != nil {
		utils.ServerError(c, "update failed")
		return
	}
	utils.OK(c, "ok")
}

func UpdateStaffCert(c *gin.Context) {
	uid, _ := c.Get("uid")
	var body struct {
		CertFiles   string `json:"cert_files"`
		HealthFiles string `json:"health_files"`
		IDCard      string `json:"id_card"`
		Intro       string `json:"intro"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	var sp models.StaffProfile
	err := database.DB.Where("user_id = ?", uid).First(&sp).Error
	if errors.Is(err, database.DB.Dialector.Explain("", nil).Error) {
		// ignore
	}
	updates := map[string]interface{}{}
	if body.CertFiles != "" {
		updates["cert_files"] = body.CertFiles
	}
	if body.HealthFiles != "" {
		updates["health_files"] = body.HealthFiles
	}
	if body.IDCard != "" {
		updates["id_card"] = body.IDCard
	}
	if body.Intro != "" {
		updates["intro"] = body.Intro
	}
	if errors.Is(err, nil) {
		if err := database.DB.Model(&sp).Updates(updates).Error; err != nil {
			utils.ServerError(c, "update failed")
			return
		}
	} else {
		sp := models.StaffProfile{
			UserID:      uid.(uint),
			CertFiles:   body.CertFiles,
			HealthFiles: body.HealthFiles,
			IDCard:      body.IDCard,
			Intro:       body.Intro,
		}
		if err := database.DB.Create(&sp).Error; err != nil {
			utils.ServerError(c, "create profile failed")
			return
		}
	}
	utils.OK(c, "ok")
}

func ListStaff(c *gin.Context) {
	var users []models.User
	q := database.DB.Where("role = ?", models.RoleStaff)
	if cid := c.Query("company_id"); cid != "" {
		q = q.Where("company_id = ?", cid)
	}
	if err := q.Find(&users).Error; err != nil {
		utils.ServerError(c, "query failed")
		return
	}
	utils.OK(c, users)
}

func GetStaffDetail(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		utils.NotFound(c, "staff not found")
		return
	}
	if user.Role != models.RoleStaff {
		utils.NotFound(c, "staff not found")
		return
	}
	var sp models.StaffProfile
	database.DB.Where("user_id = ?", user.ID).First(&sp)
	var reviews []models.Review
	database.DB.Where("staff_id = ?", user.ID).Order("created_at desc").Limit(20).Find(&reviews)
	utils.OK(c, gin.H{"user": user, "profile": sp, "reviews": reviews})
}

func ListStaffReviews(c *gin.Context) {
	id := c.Param("id")
	var reviews []models.Review
	q := database.DB.Where("staff_id = ?", id).Order("created_at desc")
	if err := q.Find(&reviews).Error; err != nil {
		utils.ServerError(c, "query failed")
		return
	}
	utils.OK(c, reviews)
}

func AdminSuspendStaff(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Suspended bool `json:"suspended"`
	}
	c.ShouldBindJSON(&body)
	if err := database.DB.Model(&models.User{}).Where("id = ?", id).Update("suspended", body.Suspended).Error; err != nil {
		utils.ServerError(c, "update failed")
		return
	}
	utils.OK(c, "ok")
}
