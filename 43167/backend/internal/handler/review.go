package handler

import (
	"net/http"

	"watchplatform/internal/app"
	"watchplatform/internal/database"
	"watchplatform/internal/middleware"
	"watchplatform/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewCreateReq struct {
	TradeID uint   `json:"trade_id" binding:"required"`
	Role    string `json:"role" binding:"required,oneof=buyer seller"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Content string `json:"content"`
}

func CreateReview(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var req ReviewCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	var t model.Trade
	if err := database.DB.First(&t, req.TradeID).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "交易不存在")
		return
	}
	if t.Status != model.TradeCompleted {
		app.Fail(c, http.StatusBadRequest, "交易未完成，无法评价")
		return
	}
	var toUserID uint
	if req.Role == "seller" && u.ID == t.BuyerID {
		toUserID = t.SellerID
	} else if req.Role == "buyer" && u.ID == t.SellerID {
		toUserID = t.BuyerID
	} else {
		app.Fail(c, http.StatusBadRequest, "无权评价")
		return
	}
	var existing model.Review
	database.DB.Where("trade_id = ? AND from_user_id = ?", t.ID, u.ID).First(&existing)
	if existing.ID != 0 {
		app.Fail(c, http.StatusConflict, "已评价")
		return
	}
	r := model.Review{
		TradeID:    t.ID,
		FromUserID: u.ID,
		ToUserID:   toUserID,
		Role:       req.Role,
		Rating:     req.Rating,
		Content:    req.Content,
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&r).Error; err != nil {
			return err
		}
		var target model.User
		if err := tx.First(&target, toUserID).Error; err != nil {
			return err
		}
		delta := req.Rating - 3
		newScore := target.CreditScore + delta
		if newScore < 0 {
			newScore = 0
		}
		if newScore > 100 {
			newScore = 100
		}
		return tx.Model(&target).Updates(map[string]interface{}{
			"credit_score": newScore,
			"review_count": target.ReviewCount + 1,
		}).Error
	})
	if err != nil {
		app.BizFail(c, err)
		return
	}
	pushMessage(toUserID, "new_review", "收到新评价", "您收到了一笔新评价", "review", r.ID)
	app.OK(c, r)
}

func ListReviews(c *gin.Context) {
	userID := c.Query("user_id")
	db := database.DB.Model(&model.Review{})
	if userID != "" {
		db = db.Where("to_user_id = ?", userID)
	}
	var list []model.Review
	db.Order("created_at desc").Find(&list)
	app.OK(c, list)
}
