package utils

import (
	"log"
	"matchmaking-platform/config"
	"matchmaking-platform/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(config.Cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	err = DB.AutoMigrate(
		&model.User{},
		&model.Profile{},
		&model.MatchRecord{},
		&model.DateRecord{},
		&model.DateReview{},
		&model.MatchmakerMember{},
		&model.MatchmakerService{},
		&model.MatchmakerStats{},
		&model.ChatMessage{},
		&model.ChatSession{},
		&model.MemberBenefit{},
		&model.MemberOrder{},
		&model.InteractLog{},
		&model.SensitiveWord{},
		&model.SystemLog{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	initDefaultData()
	log.Println("Database connected and migrated successfully")
}

func initDefaultData() {
	var count int64
	DB.Model(&model.SensitiveWord{}).Count(&count)
	if count == 0 {
		words := []model.SensitiveWord{
			{Word: "赌博", Category: "illegal"},
			{Word: "色情", Category: "illegal"},
			{Word: "贷款", Category: "illegal"},
			{Word: "诈骗", Category: "illegal"},
			{Word: "毒品", Category: "illegal"},
			{Word: "枪支", Category: "illegal"},
			{Word: "暴力", Category: "illegal"},
			{Word: "色情内容", Category: "illegal"},
			{Word: "恐怖", Category: "illegal"},
			{Word: "传销", Category: "illegal"},
		}
		DB.Create(&words)
	}

	DB.Model(&model.MemberBenefit{}).Count(&count)
	if count == 0 {
		benefits := []model.MemberBenefit{
			{
				Level:           model.MemberFree,
				DailyInteract:   5,
				UnlimitedChat:  false,
				ViewWhoLiked:   false,
				PriorityMatch:  false,
				AdvancedFilter:  false,
				VideoChat:       false,
				HideOnline:      false,
				NoAds:           false,
				MatchmakerAssist: false,
				PricePerMonth:   0,
				Description:    "免费会员：每天5次互动",
			},
			{
				Level:           model.MemberSilver,
				DailyInteract:   20,
				UnlimitedChat:  true,
				ViewWhoLiked:   true,
				PriorityMatch:  false,
				AdvancedFilter:  false,
				VideoChat:       false,
				HideOnline:      false,
				NoAds:           false,
				MatchmakerAssist: false,
				PricePerMonth:   29.9,
				Description:    "白银会员：无限聊天、查看谁喜欢我",
			},
			{
				Level:           model.MemberGold,
				DailyInteract:   50,
				UnlimitedChat:  true,
				ViewWhoLiked:   true,
				PriorityMatch:  true,
				AdvancedFilter:  true,
				VideoChat:       true,
				HideOnline:      false,
				NoAds:           false,
				MatchmakerAssist: false,
				PricePerMonth:   99.9,
				Description:    "黄金会员：优先匹配、高级筛选、视频聊天",
			},
			{
				Level:           model.MemberDiamond,
				DailyInteract:   999,
				UnlimitedChat:  true,
				ViewWhoLiked:   true,
				PriorityMatch:  true,
				AdvancedFilter:  true,
				VideoChat:       true,
				HideOnline:      true,
				NoAds:           true,
				MatchmakerAssist: true,
				PricePerMonth:   299.9,
				Description:    "钻石会员：全部特权+红娘服务",
			},
		}
		DB.Create(&benefits)
	}
}
