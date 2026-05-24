package models

import (
	"repair-platform/pkg/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	logger.Info("Auto-migrating database schemas...")

	err = DB.AutoMigrate(
		&User{},
		&TechnicianProfile{},
		&Category{},
		&ServiceItem{},
		&Order{},
		&OrderLog{},
		&Review{},
		&Part{},
		&PartRequest{},
		&PartRequestItem{},
		&PartUsage{},
		&WithdrawRequest{},
		&Transaction{},
		&MonthlyReport{},
	)
	if err != nil {
		return err
	}

	logger.Info("Database migration completed")
	return nil
}

func InitSeedData() {
	var categoryCount int64
	DB.Model(&Category{}).Count(&categoryCount)
	if categoryCount > 0 {
		return
	}

	categories := []Category{
		{Name: "家电维修", Code: "appliance", Icon: "🏠", Sort: 1},
		{Name: "数码维修", Code: "digital", Icon: "📱", Sort: 2},
		{Name: "汽车维修", Code: "auto", Icon: "🚗", Sort: 3},
	}

	for i := range categories {
		DB.Create(&categories[i])
	}

	serviceItems := []ServiceItem{
		{CategoryID: categories[0].ID, Name: "空调清洗", Description: "专业空调清洗保养服务", MinPrice: 100, MaxPrice: 300, EstimatedTime: 60},
		{CategoryID: categories[0].ID, Name: "冰箱维修", Description: "冰箱不制冷、漏水等故障维修", MinPrice: 150, MaxPrice: 500, EstimatedTime: 90},
		{CategoryID: categories[0].ID, Name: "洗衣机维修", Description: "洗衣机不排水、不转动等故障", MinPrice: 80, MaxPrice: 300, EstimatedTime: 60},
		{CategoryID: categories[0].ID, Name: "热水器维修", Description: "热水器不出热水、漏电等维修", MinPrice: 100, MaxPrice: 400, EstimatedTime: 90},
		{CategoryID: categories[1].ID, Name: "手机换屏", Description: "手机屏幕更换服务", MinPrice: 200, MaxPrice: 800, EstimatedTime: 60},
		{CategoryID: categories[1].ID, Name: "笔记本维修", Description: "笔记本电脑故障维修", MinPrice: 100, MaxPrice: 1000, EstimatedTime: 120},
		{CategoryID: categories[1].ID, Name: "平板维修", Description: "平板电脑故障维修服务", MinPrice: 150, MaxPrice: 600, EstimatedTime: 60},
		{CategoryID: categories[2].ID, Name: "汽车保养", Description: "常规保养、机油更换", MinPrice: 300, MaxPrice: 800, EstimatedTime: 90},
		{CategoryID: categories[2].ID, Name: "汽车维修", Description: "发动机、变速箱等故障维修", MinPrice: 500, MaxPrice: 5000, EstimatedTime: 240},
		{CategoryID: categories[2].ID, Name: "汽车美容", Description: "汽车美容、抛光打蜡", MinPrice: 200, MaxPrice: 500, EstimatedTime: 120},
	}

	for i := range serviceItems {
		DB.Create(&serviceItems[i])
	}

	parts := []Part{
		{Name: "空调压缩机", Code: "AC-COMP-001", Category: "家电配件", Price: 800, Stock: 50, MinStock: 10},
		{Name: "空调电容", Code: "AC-CAP-001", Category: "家电配件", Price: 50, Stock: 100, MinStock: 20},
		{Name: "洗衣机电机", Code: "WM-MOT-001", Category: "家电配件", Price: 300, Stock: 30, MinStock: 5},
		{Name: "iPhone屏幕总成", Code: "IPH-SCR-001", Category: "数码配件", Price: 500, Stock: 20, MinStock: 5},
		{Name: "华为屏幕总成", Code: "HW-SCR-001", Category: "数码配件", Price: 400, Stock: 15, MinStock: 5},
		{Name: "机油滤芯", Code: "AUTO-OIL-001", Category: "汽车配件", Price: 50, Stock: 100, MinStock: 30},
		{Name: "刹车片", Code: "AUTO-BRK-001", Category: "汽车配件", Price: 200, Stock: 40, MinStock: 10},
	}

	for i := range parts {
		DB.Create(&parts[i])
	}

	logger.Info("Seed data initialized")
}
