package database

import (
	"log"
	"property-management/internal/config"
	"property-management/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	log.Println("Database connected successfully")

	DB.AutoMigrate(
		&model.User{},
		&model.Property{},
		&model.PropertyImage{},
		&model.Facility{},
		&model.Tenant{},
		&model.Appointment{},
		&model.Contract{},
		&model.RentRecord{},
		&model.RepairOrder{},
		&model.UtilityFee{},
		&model.Notice{},
	)

	seedData()
}

func seedData() {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count == 0 {
		admin := model.User{
			Username: "admin",
			Password: hashPassword("admin123"),
			RealName: "系统管理员",
			Phone: "13800138000",
			Email: "admin@example.com",
			Role: "admin",
			Status: 1,
		}
		DB.Create(&admin)

		landlord := model.User{
			Username: "landlord",
			Password: hashPassword("123456"),
			RealName: "张房东",
			Phone: "13900139000",
			Email: "landlord@example.com",
			Role: "landlord",
			Status: 1,
		}
		DB.Create(&landlord)

		maintainer := model.User{
			Username: "maintainer",
			Password: hashPassword("123456"),
			RealName: "李维修",
			Phone: "13700137000",
			Email: "maintainer@example.com",
			Role: "maintainer",
			Status: 1,
		}
		DB.Create(&maintainer)

		facilities := []model.Facility{
			{Name: "空调", Icon: "❄️"},
			{Name: "洗衣机", Icon: "🧺"},
			{Name: "冰箱", Icon: "🧊"},
			{Name: "热水器", Icon: "🚿"},
			{Name: "宽带", Icon: "📶"},
			{Name: "电视", Icon: "📺"},
			{Name: "床", Icon: "🛏️"},
			{Name: "衣橱", Icon: "👔"},
			{Name: "沙发", Icon: "🛋️"},
			{Name: "独立卫浴", Icon: "🚽"},
		}
		DB.Create(&facilities)
	}
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
