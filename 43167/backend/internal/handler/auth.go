package handler

import (
	"net/http"
	"regexp"
	"time"

	"watchplatform/internal/app"
	"watchplatform/internal/database"
	"watchplatform/internal/middleware"
	"watchplatform/internal/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterReq struct {
	Username string     `json:"username" binding:"required,min=3,max=32"`
	Password string     `json:"password" binding:"required,min=6,max=64"`
	Role     model.Role `json:"role" binding:"required,oneof=buyer seller appraiser"`
	Email    string     `json:"email"`
	Phone    string     `json:"phone"`
	RealName string     `json:"real_name"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var emailRe = regexp.MustCompile(`^[\w.+-]+@[\w-]+\.[\w.-]+$`)

func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	if req.Email != "" && !emailRe.MatchString(req.Email) {
		app.Fail(c, http.StatusBadRequest, "邮箱格式不正确")
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		app.BizFail(c, err)
		return
	}
	u := model.User{
		Username:    req.Username,
		Password:    string(hashed),
		Role:        req.Role,
		Email:       req.Email,
		Phone:       req.Phone,
		RealName:    req.RealName,
		CreditScore: 100,
	}
	if err := database.DB.Create(&u).Error; err != nil {
		app.Fail(c, http.StatusConflict, "用户名已存在")
		return
	}
	app.OK(c, gin.H{"id": u.ID})
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	var u model.User
	if err := database.DB.Where("username = ?", req.Username).First(&u).Error; err != nil {
		app.Fail(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		app.Fail(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	token, err := middleware.GenerateToken(&u)
	if err != nil {
		app.BizFail(c, err)
		return
	}
	app.OK(c, gin.H{
		"token": token,
		"expires_in": int(time.Now().Add(time.Hour * 72).Unix()),
		"user": gin.H{
			"id":           u.ID,
			"username":     u.Username,
			"role":         u.Role,
			"email":        u.Email,
			"credit_score": u.CreditScore,
		},
	})
}

func Me(c *gin.Context) {
	u := middleware.CurrentUser(c)
	if u == nil {
		app.Fail(c, http.StatusUnauthorized, "未登录")
		return
	}
	app.OK(c, u)
}
