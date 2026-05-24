package handler

import (
	"net/http"
	"path/filepath"
	"strings"

	"watchplatform/internal/app"
	"watchplatform/internal/config"
	"watchplatform/internal/database"
	"watchplatform/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ProfileUpdateReq struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RealName string `json:"real_name"`
	Avatar   string `json:"avatar"`
}

func UpdateProfile(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var req ProfileUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	updates := map[string]interface{}{}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if len(updates) == 0 {
		app.Fail(c, http.StatusBadRequest, "无更新字段")
		return
	}
	if err := database.DB.Model(u).Updates(updates).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	database.DB.First(u, u.ID)
	app.OK(c, u)
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=64"`
}

func ChangePassword(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var req ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.OldPassword)); err != nil {
		app.Fail(c, http.StatusBadRequest, "旧密码错误")
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		app.BizFail(c, err)
		return
	}
	if err := database.DB.Model(u).Update("password", string(hashed)).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	app.OK(c, nil)
}

func UploadAvatar(c *gin.Context) {
	u := middleware.CurrentUser(c)
	file, err := c.FormFile("file")
	if err != nil {
		app.Fail(c, http.StatusBadRequest, "未找到文件")
		return
	}
	if file.Size > config.Cfg.MaxFileSize {
		app.Fail(c, http.StatusBadRequest, "文件过大")
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		app.Fail(c, http.StatusBadRequest, "仅支持图片文件")
		return
	}
	dir := filepath.Join(config.Cfg.UploadDir, "avatars")
	filename := uuid.New().String() + ext
	dst := filepath.Join(dir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		app.BizFail(c, err)
		return
	}
	url := "/uploads/avatars/" + filename
	if err := database.DB.Model(u).Update("avatar", url).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	app.OK(c, gin.H{"url": url})
}
