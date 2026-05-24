package handler

import (
	"net/http"
	"strconv"

	"watchplatform/internal/app"
	"watchplatform/internal/database"
	"watchplatform/internal/middleware"
	"watchplatform/internal/model"

	"github.com/gin-gonic/gin"
)

type FavoriteGroupCreateReq struct {
	Brand string `json:"brand" binding:"required,max=64"`
	Name  string `json:"name"`
}

type FavoriteCreateReq struct {
	WatchID uint `json:"watch_id" binding:"required"`
	GroupID uint `json:"group_id"`
}

func ListFavoriteGroups(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var list []model.FavoriteGroup
	database.DB.Where("user_id = ?", u.ID).Find(&list)
	app.OK(c, list)
}

func CreateFavoriteGroup(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var req FavoriteGroupCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	g := model.FavoriteGroup{
		UserID: u.ID,
		Brand:  req.Brand,
		Name:   req.Name,
	}
	if g.Name == "" {
		g.Name = req.Brand
	}
	if err := database.DB.Create(&g).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	app.OK(c, g)
}

func DeleteFavoriteGroup(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := database.DB.Where("id = ? AND user_id = ?", id, u.ID).Delete(&model.FavoriteGroup{}).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	database.DB.Where("user_id = ? AND group_id = ?", u.ID, id).Update("group_id", nil)
	app.OK(c, nil)
}

func ListFavorites(c *gin.Context) {
	u := middleware.CurrentUser(c)
	groupID := c.Query("group_id")
	brand := c.Query("brand")
	db := database.DB.Model(&model.Favorite{}).Where("user_id = ?", u.ID)
	if groupID != "" {
		db = db.Where("group_id = ?", groupID)
	}
	if brand != "" {
		db = db.Joins("JOIN watches ON watches.id = favorites.watch_id").Where("watches.brand = ?", brand)
	}
	var list []model.Favorite
	db.Find(&list)
	app.OK(c, list)
}

func AddFavorite(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var req FavoriteCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	var existing model.Favorite
	database.DB.Where("user_id = ? AND watch_id = ?", u.ID, req.WatchID).First(&existing)
	if existing.ID != 0 {
		app.Fail(c, http.StatusConflict, "已收藏")
		return
	}
	f := model.Favorite{
		UserID:  u.ID,
		WatchID: req.WatchID,
		GroupID: req.GroupID,
	}
	if err := database.DB.Create(&f).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	app.OK(c, f)
}

func RemoveFavorite(c *gin.Context) {
	u := middleware.CurrentUser(c)
	watchID, _ := strconv.ParseUint(c.Param("watch_id"), 10, 64)
	if err := database.DB.Where("user_id = ? AND watch_id = ?", u.ID, watchID).Delete(&model.Favorite{}).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	app.OK(c, nil)
}
