package handler

import (
	"net/http"
	"strconv"

	"watchplatform/internal/app"
	"watchplatform/internal/database"
	"watchplatform/internal/middleware"
	"watchplatform/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TradeCreateReq struct {
	WatchID    uint    `json:"watch_id" binding:"required"`
	StartPrice float64 `json:"start_price" binding:"required,gt=0"`
	Remark     string  `json:"remark"`
}

type BidReq struct {
	Price   float64 `json:"price" binding:"required,gt=0"`
	Message string  `json:"message"`
}

type AcceptBidReq struct {
	BidID uint `json:"bid_id" binding:"required"`
}

func CreateTrade(c *gin.Context) {
	u := middleware.CurrentUser(c)
	var req TradeCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	var w model.Watch
	if err := database.DB.First(&w, req.WatchID).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "手表不存在")
		return
	}
	if w.SellerID != u.ID {
		app.Fail(c, http.StatusForbidden, "非本人手表")
		return
	}
	t := model.Trade{
		SellerID:   u.ID,
		WatchID:    req.WatchID,
		StartPrice: req.StartPrice,
		FinalPrice: req.StartPrice,
		Status:     model.TradeOpen,
		Remark:     req.Remark,
	}
	if err := database.DB.Create(&t).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	app.OK(c, t)
}

func ListTrades(c *gin.Context) {
	u := middleware.CurrentUser(c)
	db := database.DB.Model(&model.Trade{})
	if u.Role == model.RoleSeller {
		db = db.Where("seller_id = ?", u.ID)
	} else if u.Role == model.RoleBuyer {
		db = db.Where("buyer_id = ? OR status in ?", u.ID, []string{string(model.TradeOpen), string(model.TradeBidding)})
	}
	status := c.Query("status")
	if status != "" {
		db = db.Where("status = ?", status)
	}
	var total int64
	db.Count(&total)
	var list []model.Trade
	db.Order("created_at desc").Find(&list)
	app.OK(c, gin.H{"total": total, "list": list})
}

func GetTrade(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var t model.Trade
	if err := database.DB.First(&t, id).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "交易不存在")
		return
	}
	var bids []model.TradeBid
	database.DB.Where("trade_id = ?", t.ID).Order("created_at desc").Find(&bids)
	app.OK(c, gin.H{"trade": t, "bids": bids})
}

func PlaceBid(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var t model.Trade
	if err := database.DB.First(&t, id).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "交易不存在")
		return
	}
	if t.Status != model.TradeOpen && t.Status != model.TradeBidding {
		app.Fail(c, http.StatusBadRequest, "当前状态不可出价")
		return
	}
	var req BidReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	bid := model.TradeBid{
		TradeID: t.ID,
		BuyerID: u.ID,
		Price:   req.Price,
		Message: req.Message,
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&bid).Error; err != nil {
			return err
		}
		return tx.Model(&t).Update("status", model.TradeBidding).Error
	})
	if err != nil {
		app.BizFail(c, err)
		return
	}
	pushMessage(t.SellerID, "new_bid", "收到新的出价", "您的交易收到一笔新出价", "trade", t.ID)
	app.OK(c, bid)
}

func AcceptBid(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var t model.Trade
	if err := database.DB.First(&t, id).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "交易不存在")
		return
	}
	if t.SellerID != u.ID {
		app.Fail(c, http.StatusForbidden, "无权操作")
		return
	}
	var req AcceptBidReq
	if err := c.ShouldBindJSON(&req); err != nil {
		app.BindFail(c, err)
		return
	}
	var bid model.TradeBid
	if err := database.DB.First(&bid, req.BidID).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "出价不存在")
		return
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&bid).Update("accepted", true).Error; err != nil {
			return err
		}
		if err := tx.Model(&t).Updates(map[string]interface{}{
			"buyer_id":    bid.BuyerID,
			"final_price": bid.Price,
			"status":      model.TradePendingDeal,
		}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Watch{}).Where("id = ?", t.WatchID).Update("status", "reserved").Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		app.BizFail(c, err)
		return
	}
	pushMessage(bid.BuyerID, "bid_accepted", "出价已被接受", "您的出价已被卖家接受", "trade", t.ID)
	app.OK(c, nil)
}

func ShipTrade(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var t model.Trade
	if err := database.DB.First(&t, id).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "交易不存在")
		return
	}
	if t.SellerID != u.ID {
		app.Fail(c, http.StatusForbidden, "无权操作")
		return
	}
	if t.Status != model.TradePendingDeal {
		app.Fail(c, http.StatusBadRequest, "当前状态不可发货")
		return
	}
	if err := database.DB.Model(&t).Update("status", model.TradeShipped).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	pushMessage(t.BuyerID, "trade_shipped", "商品已发货", "卖家已发货，请等待收货", "trade", t.ID)
	app.OK(c, nil)
}

func CompleteTrade(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var t model.Trade
	if err := database.DB.First(&t, id).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "交易不存在")
		return
	}
	if t.BuyerID != u.ID && t.SellerID != u.ID {
		app.Fail(c, http.StatusForbidden, "无权操作")
		return
	}
	if err := database.DB.Model(&t).Update("status", model.TradeCompleted).Error; err != nil {
		app.BizFail(c, err)
		return
	}
	pushMessage(t.SellerID, "trade_completed", "交易已完成", "交易已完成，可进行互评", "trade", t.ID)
	pushMessage(t.BuyerID, "trade_completed", "交易已完成", "交易已完成，可进行互评", "trade", t.ID)
	app.OK(c, nil)
}

func CancelTrade(c *gin.Context) {
	u := middleware.CurrentUser(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var t model.Trade
	if err := database.DB.First(&t, id).Error; err != nil {
		app.Fail(c, http.StatusNotFound, "交易不存在")
		return
	}
	if t.SellerID != u.ID && t.BuyerID != u.ID {
		app.Fail(c, http.StatusForbidden, "无权操作")
		return
	}
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&t).Update("status", model.TradeCanceled).Error; err != nil {
			return err
		}
		return tx.Model(&model.Watch{}).Where("id = ?", t.WatchID).Update("status", "on_sale").Error
	}); err != nil {
		app.BizFail(c, err)
		return
	}
	app.OK(c, nil)
}
