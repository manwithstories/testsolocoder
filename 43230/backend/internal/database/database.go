package database

import (
	"fmt"
	"log"

	"print3d-platform/internal/config"
	"print3d-platform/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	log.Println("Database connection established successfully")

	return db, nil
}

func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.DesignerProfile{},
		&models.PrinterProfile{},
		&models.Model3D{},
		&models.ModelVersion{},
		&models.ModelPurchase{},
		&models.ModelFavorite{},
		&models.Material{},
		&models.ColorOption{},
		&models.PrintOrder{},
		&models.OrderHistory{},
		&models.Settlement{},
		&models.PrinterDevice{},
		&models.MaterialInventory{},
		&models.PrintSchedule{},
		&models.Review{},
		&models.FileUpload{},
		&models.FileAccessLog{},
		&models.DownloadRecord{},
		&models.Notification{},
		&models.Transaction{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func SeedData() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	seedMaterials()
	log.Println("Seed data completed")
	return nil
}

func seedMaterials() {
	materials := []models.Material{
		{
			Type:           models.MaterialPLA,
			Name:           "PLA 通用型",
			Description:    "聚乳酸，环保可降解，适合初学者，打印温度190-220°C",
			PricePerGram:   0.05,
			PrintSpeed:     60,
			Density:        1.24,
			Strength:       "中等",
			TemperatureResistance: "60°C",
			IsAvailable:    true,
		},
		{
			Type:           models.MaterialABS,
			Name:           "ABS 工程塑料",
			Description:    "丙烯腈-丁二烯-苯乙烯，强度高，耐冲击，打印温度220-250°C",
			PricePerGram:   0.08,
			PrintSpeed:     50,
			Density:        1.04,
			Strength:       "高",
			TemperatureResistance: "90°C",
			IsAvailable:    true,
		},
		{
			Type:           models.MaterialPETG,
			Name:           "PETG 透明级",
			Description:    "聚对苯二甲酸乙二醇酯，透明，耐化学腐蚀，打印温度220-250°C",
			PricePerGram:   0.07,
			PrintSpeed:     55,
			Density:        1.27,
			Strength:       "高",
			TemperatureResistance: "80°C",
			IsAvailable:    true,
		},
		{
			Type:           models.MaterialTPU,
			Name:           "TPU 柔性材料",
			Description:    "热塑性聚氨酯，柔性弹性，适合打印柔性件，打印温度210-240°C",
			PricePerGram:   0.12,
			PrintSpeed:     30,
			Density:        1.21,
			Strength:       "中等",
			TemperatureResistance: "80°C",
			IsAvailable:    true,
		},
		{
			Type:           models.MaterialResin,
			Name:           "光固化树脂",
			Description:    "光固化树脂，高精度，适合细节丰富的模型",
			PricePerGram:   0.15,
			PrintSpeed:     20,
			Density:        1.1,
			Strength:       "中等",
			TemperatureResistance: "60°C",
			IsAvailable:    true,
		},
	}

	for _, m := range materials {
		DB.FirstOrCreate(&m, models.Material{Type: m.Type, Name: m.Name})
	}
}
