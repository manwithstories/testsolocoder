package handlers

import (
	"time"

	"housekeeping-platform/config"
	"housekeeping-platform/models"
	"housekeeping-platform/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Role     string `json:"role" binding:"required"`
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string      `json:"token"`
	UserInfo interface{} `json:"user_info"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if !utils.IsValidPhone(req.Phone) {
		utils.BadRequest(c, "手机号格式错误")
		return
	}

	if !utils.IsValidPassword(req.Password) {
		utils.BadRequest(c, "密码格式错误，需6-20位且包含字母和数字")
		return
	}

	var existingUser models.User
	if result := config.DB.Where("phone = ?", req.Phone).First(&existingUser); result.Error == nil {
		utils.BadRequest(c, "手机号已注册")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "密码加密失败")
		return
	}

	user := models.User{
		Phone:    req.Phone,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Role:     models.UserRole(req.Role),
		IsActive: true,
	}

	if user.Role == models.RoleServiceProvider {
		user.ProviderStatus = models.ProviderStatusPending
	}

	if result := config.DB.Create(&user); result.Error != nil {
		utils.InternalError(c, "注册失败")
		return
	}

	cfg := config.Load()
	token, err := utils.GenerateToken(user.ID, user.Phone, string(user.Role), user.Nickname, cfg.JWT.Secret, cfg.JWT.ExpireHour)
	if err != nil {
		utils.InternalError(c, "生成令牌失败")
		return
	}

	go utils.LogOperation(user.ID, string(user.Role), "auth", "register", &user.ID, "user", "用户注册", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, LoginResponse{
		Token:    token,
		UserInfo: user,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if !utils.IsValidPhone(req.Phone) {
		utils.BadRequest(c, "手机号格式错误")
		return
	}

	var user models.User
	if result := config.DB.Where("phone = ?", req.Phone).First(&user); result.Error != nil {
		utils.Unauthorized(c, "手机号或密码错误")
		return
	}

	if !user.IsActive {
		utils.Forbidden(c, "账号已被禁用")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Unauthorized(c, "手机号或密码错误")
		return
	}

	now := time.Now()
	config.DB.Model(&user).Update("last_login_at", now)

	cfg := config.Load()
	token, err := utils.GenerateToken(user.ID, user.Phone, string(user.Role), user.Nickname, cfg.JWT.Secret, cfg.JWT.ExpireHour)
	if err != nil {
		utils.InternalError(c, "生成令牌失败")
		return
	}

	go utils.LogOperation(user.ID, string(user.Role), "auth", "login", &user.ID, "user", "用户登录", c.ClientIP(), c.Request.UserAgent())

	utils.Success(c, LoginResponse{
		Token:    token,
		UserInfo: user,
	})
}

func GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if result := config.DB.Preload("Addresses").Preload("Certifications").First(&user, userID); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Nickname   string `json:"nickname"`
		Avatar     string `json:"avatar"`
		Gender     string `json:"gender"`
		Age        int    `json:"age"`
		RealName   string `json:"real_name"`
		IDCard     string `json:"id_card"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if req.IDCard != "" && !utils.IsValidIDCard(req.IDCard) {
		utils.BadRequest(c, "身份证号格式错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.Age > 0 {
		updates["age"] = req.Age
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.IDCard != "" {
		updates["id_card"] = req.IDCard
	}

	var user models.User
	if result := config.DB.Model(&user).Where("id = ?", userID).Updates(updates); result.Error != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	go utils.LogOperation(userID, c.GetString("role"), "user", "update_profile", &userID, "user", "更新个人资料", c.ClientIP(), c.Request.UserAgent())

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if !utils.IsValidPassword(req.NewPassword) {
		utils.BadRequest(c, "新密码格式错误")
		return
	}

	var user models.User
	if result := config.DB.First(&user, userID); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		utils.BadRequest(c, "原密码错误")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "密码加密失败")
		return
	}

	config.DB.Model(&user).Update("password", string(hashedPassword))

	go utils.LogOperation(userID, c.GetString("role"), "user", "change_password", &userID, "user", "修改密码", c.ClientIP(), c.Request.UserAgent())

	utils.SuccessWithMessage(c, "密码修改成功", nil)
}

func Logout(c *gin.Context) {
	userID := c.GetUint("user_id")
	go utils.LogOperation(userID, c.GetString("role"), "auth", "logout", &userID, "user", "用户退出", c.ClientIP(), c.Request.UserAgent())

	utils.SuccessWithMessage(c, "退出成功", nil)
}
