package handler

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"watchplatform/internal/app"
	"watchplatform/internal/config"
	"watchplatform/internal/database"
	"watchplatform/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func StatsOverview(c *gin.Context) {
	var totalTrades int64
	var totalAmount float64
	var completed int64
	database.DB.Model(&model.Trade{}).Count(&totalTrades)
	database.DB.Model(&model.Trade{}).Where("status = ?", model.TradeCompleted).
		Select("COALESCE(SUM(final_price),0)").Scan(&totalAmount)
	database.DB.Model(&model.Trade{}).Where("status = ?", model.TradeCompleted).Count(&completed)
	var watches int64
	var users int64
	var authOrders int64
	database.DB.Model(&model.Watch{}).Count(&watches)
	database.DB.Model(&model.User{}).Count(&users)
	database.DB.Model(&model.AuthOrder{}).Count(&authOrders)
	app.OK(c, gin.H{
		"total_trades":   totalTrades,
		"completed":      completed,
		"total_amount":   totalAmount,
		"watches":        watches,
		"users":          users,
		"auth_orders":    authOrders,
	})
}

func StatsBrands(c *gin.Context) {
	type Row struct {
		Brand string  `json:"brand"`
		Count int64   `json:"count"`
		Total float64 `json:"total"`
	}
	var list []Row
	database.DB.Model(&model.Trade{}).
		Select("watches.brand as brand, count(trades.id) as count, COALESCE(SUM(trades.final_price),0) as total").
		Joins("JOIN watches ON watches.id = trades.watch_id").
		Where("trades.status = ?", model.TradeCompleted).
		Group("watches.brand").
		Order("count desc").
		Limit(20).
		Scan(&list)
	app.OK(c, list)
}

func StatsExport(c *gin.Context) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetCellValue(sheet, "A1", "交易ID")
	f.SetCellValue(sheet, "B1", "卖家ID")
	f.SetCellValue(sheet, "C1", "买家ID")
	f.SetCellValue(sheet, "D1", "手表品牌")
	f.SetCellValue(sheet, "E1", "手表型号")
	f.SetCellValue(sheet, "F1", "成交价")
	f.SetCellValue(sheet, "G1", "状态")
	f.SetCellValue(sheet, "H1", "创建时间")

	type Row struct {
		ID         uint
		SellerID   uint
		BuyerID    uint
		Brand      string
		Model      string
		FinalPrice float64
		Status     string
		CreatedAt  time.Time
	}
	var rows []Row
	database.DB.Table("trades").
		Select("trades.id, trades.seller_id, trades.buyer_id, watches.brand, watches.model, trades.final_price, trades.status, trades.created_at").
		Joins("JOIN watches ON watches.id = trades.watch_id").
		Order("trades.id desc").
		Scan(&rows)

	for i, r := range rows {
		row := strconv.Itoa(i + 2)
		f.SetCellValue(sheet, "A"+row, r.ID)
		f.SetCellValue(sheet, "B"+row, r.SellerID)
		f.SetCellValue(sheet, "C"+row, r.BuyerID)
		f.SetCellValue(sheet, "D"+row, r.Brand)
		f.SetCellValue(sheet, "E"+row, r.Model)
		f.SetCellValue(sheet, "F"+row, r.FinalPrice)
		f.SetCellValue(sheet, "G"+row, r.Status)
		f.SetCellValue(sheet, "H"+row, r.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	filename := "stats_" + time.Now().Format("20060102150405") + ".xlsx"
	path := filepath.Join(config.Cfg.UploadDir, filename)
	if err := f.SaveAs(path); err != nil {
		app.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.FileAttachment(path, filename)
}
