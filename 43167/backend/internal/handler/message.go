package handler

import (
	"strconv"

	"watchplatform/internal/app"
	"watchplatform/internal/database"
	"watchplatform/internal/middleware"
	"watchplatform/internal/model"

	"github.com/gin-gonic/gin"
)

func pushMessage(userID uint, typ, title, content, refType string, refID uint) {
	if userID == 0 {
		return
	}
	m := model.Message{
		UserID:  userID,
		Type:    typ,
		Title:   title,
		Content: content,
		RefType: refType,
		RefID:   refID,
	}
	_ = database.DB.Create(&m).Error
}

func ListMessages(c *gin.Context) {
	u := middleware.CurrentUser(c)
	unread := c.Query("unread")
	db := database.DB.Model(&model.Message{}).Where("user_id = ?", u.ID)
	if unread == "1" {
		db = db.Where("`read` = ?", false)
	}
	var total int64
	db.Count(&total)
	var list []model.Message
	db.Order("created_at desc").Limit(100).Find(&list)
	app.OK(c, gin.H{"total": total, "list": list})
}

func MarkMessageRead(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	database.DB.Model(&model.Message{}).Where("id = ? AND user_id = ?", id, u.ID).Update("read", true)
	app.OK(c, nil)
}

func MarkAllRead(c *gin.Context) {
	u := middleware.CurrentUser(c)
	database.DB.Model(&model.Message{}).Where("user_id = ?", u.ID).Update("read", true)
	app.OK(c, nil)
}
