package database

import (
	"hotel-system/internal/model"
	"hotel-system/internal/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func InitSeedData(db *gorm.DB) {
	logger.Info("Starting seed data initialization...")

	createDefaultUsers(db)
	createMemberLevels(db)
	createTestRoomTypes(db)
	createTestRooms(db)

	logger.Info("Seed data initialization completed")
}

func createDefaultUsers(db *gorm.DB) {
	var adminCount int64
	db.Model(&model.User{}).Where("username = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			logger.Errorf("Failed to hash admin password: %v", err)
			return
		}

		admin := model.User{
			Username: "admin",
			Password: string(hashedPassword),
			RealName: "系统管理员",
			Role:     model.UserRoleAdmin,
			Status:   model.UserStatusActive,
		}
		if err := db.Create(&admin).Error; err != nil {
			logger.Errorf("Failed to create admin user: %v", err)
		} else {
			logger.Info("Default admin user created: admin/admin123")
		}
	} else {
		logger.Info("Admin user already exists, skipping")
	}

	var frontDeskCount int64
	db.Model(&model.User{}).Where("username = ?", "frontdesk").Count(&frontDeskCount)
	if frontDeskCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("frontdesk123"), bcrypt.DefaultCost)
		if err != nil {
			logger.Errorf("Failed to hash frontdesk password: %v", err)
			return
		}

		frontDesk := model.User{
			Username: "frontdesk",
			Password: string(hashedPassword),
			RealName: "前台接待",
			Role:     model.UserRoleFrontDesk,
			Status:   model.UserStatusActive,
		}
		if err := db.Create(&frontDesk).Error; err != nil {
			logger.Errorf("Failed to create frontdesk user: %v", err)
		} else {
			logger.Info("Default frontdesk user created: frontdesk/frontdesk123")
		}
	} else {
		logger.Info("Frontdesk user already exists, skipping")
	}
}

func createMemberLevels(db *gorm.DB) {
	levels := []struct {
		name         string
		discountRate float64
		pointsRate   float64
		minPoints    int
		maxPoints    int
	}{
		{name: "普通会员", discountRate: 1.0, pointsRate: 1.0, minPoints: 0, maxPoints: 999},
		{name: "银卡会员", discountRate: 0.95, pointsRate: 1.2, minPoints: 1000, maxPoints: 4999},
		{name: "金卡会员", discountRate: 0.9, pointsRate: 1.5, minPoints: 5000, maxPoints: 19999},
		{name: "钻石会员", discountRate: 0.85, pointsRate: 2.0, minPoints: 20000, maxPoints: 999999},
	}

	for _, level := range levels {
		var count int64
		db.Model(&model.MemberLevel{}).Where("name = ?", level.name).Count(&count)
		if count == 0 {
			memberLevel := model.MemberLevel{
				Name:         level.name,
				DiscountRate: level.discountRate,
				PointsRate:   level.pointsRate,
				MinPoints:    level.minPoints,
				MaxPoints:    level.maxPoints,
			}
			if err := db.Create(&memberLevel).Error; err != nil {
				logger.Errorf("Failed to create member level %s: %v", level.name, err)
			} else {
				logger.Infof("Member level created: %s", level.name)
			}
		} else {
			logger.Infof("Member level %s already exists, skipping", level.name)
		}
	}
}

func createTestRoomTypes(db *gorm.DB) {
	roomTypes := []struct {
		name        string
		description string
		basePrice   float64
		bedCount    int
		maxGuests   int
		facilities  model.StringArray
	}{
		{
			name:        "标准单人间",
			description: "经济实惠的单人间，配备1张单人床，适合商务出差或独自旅行者",
			basePrice:   198.00,
			bedCount:    1,
			maxGuests:   1,
			facilities:  model.StringArray{"免费WiFi", "独立卫浴", "空调", "电视"},
		},
		{
			name:        "标准双床间",
			description: "舒适的双床间，配备2张单人床，适合朋友或同事同住",
			basePrice:   258.00,
			bedCount:    2,
			maxGuests:   2,
			facilities:  model.StringArray{"免费WiFi", "独立卫浴", "空调", "电视", "迷你吧"},
		},
		{
			name:        "豪华大床房",
			description: "宽敞的大床房，配备1张双人床，适合情侣或家庭入住",
			basePrice:   388.00,
			bedCount:    1,
			maxGuests:   2,
			facilities:  model.StringArray{"免费WiFi", "独立卫浴", "空调", "智能电视", "迷你吧", "保险箱"},
		},
		{
			name:        "行政套房",
			description: "高端商务套房，独立卧室和客厅，配备高端设施",
			basePrice:   688.00,
			bedCount:    1,
			maxGuests:   2,
			facilities:  model.StringArray{"免费WiFi", "独立卫浴", "按摩浴缸", "空调", "智能电视", "迷你吧", "保险箱", "行政酒廊"},
		},
	}

	for _, rt := range roomTypes {
		var count int64
		db.Model(&model.RoomType{}).Where("name = ?", rt.name).Count(&count)
		if count == 0 {
			roomType := model.RoomType{
				Name:        rt.name,
				Description: rt.description,
				BasePrice:   rt.basePrice,
				BedCount:    rt.bedCount,
				MaxGuests:   rt.maxGuests,
				Facilities:  rt.facilities,
			}
			if err := db.Create(&roomType).Error; err != nil {
				logger.Errorf("Failed to create room type %s: %v", rt.name, err)
			} else {
				logger.Infof("Room type created: %s", rt.name)
			}
		} else {
			logger.Infof("Room type %s already exists, skipping", rt.name)
		}
	}
}

func createTestRooms(db *gorm.DB) {
	var roomCount int64
	db.Model(&model.Room{}).Count(&roomCount)
	if roomCount > 0 {
		logger.Info("Rooms already exist, skipping test room creation")
		return
	}

	var roomTypes []model.RoomType
	db.Find(&roomTypes)
	if len(roomTypes) == 0 {
		logger.Warn("No room types found, cannot create test rooms")
		return
	}

	roomTypeMap := make(map[string]uint)
	for _, rt := range roomTypes {
		roomTypeMap[rt.Name] = rt.ID
	}

	rooms := []struct {
		roomNo     string
		floor      int
		roomTypeID uint
		status     model.RoomStatus
		price      float64
	}{
		{roomNo: "101", floor: 1, roomTypeID: roomTypeMap["标准单人间"], status: model.RoomStatusAvailable, price: 198.00},
		{roomNo: "102", floor: 1, roomTypeID: roomTypeMap["标准单人间"], status: model.RoomStatusAvailable, price: 198.00},
		{roomNo: "103", floor: 1, roomTypeID: roomTypeMap["标准单人间"], status: model.RoomStatusMaintenance, price: 198.00},
		{roomNo: "201", floor: 2, roomTypeID: roomTypeMap["标准双床间"], status: model.RoomStatusAvailable, price: 258.00},
		{roomNo: "202", floor: 2, roomTypeID: roomTypeMap["标准双床间"], status: model.RoomStatusAvailable, price: 258.00},
		{roomNo: "203", floor: 2, roomTypeID: roomTypeMap["标准双床间"], status: model.RoomStatusAvailable, price: 258.00},
		{roomNo: "301", floor: 3, roomTypeID: roomTypeMap["豪华大床房"], status: model.RoomStatusAvailable, price: 388.00},
		{roomNo: "302", floor: 3, roomTypeID: roomTypeMap["豪华大床房"], status: model.RoomStatusAvailable, price: 388.00},
		{roomNo: "401", floor: 4, roomTypeID: roomTypeMap["行政套房"], status: model.RoomStatusAvailable, price: 688.00},
		{roomNo: "402", floor: 4, roomTypeID: roomTypeMap["行政套房"], status: model.RoomStatusAvailable, price: 688.00},
	}

	for _, r := range rooms {
		if r.roomTypeID == 0 {
			continue
		}
		room := model.Room{
			RoomNo:     r.roomNo,
			Floor:      r.floor,
			RoomTypeID: r.roomTypeID,
			Status:     r.status,
			Price:      r.price,
		}
		if err := db.Create(&room).Error; err != nil {
			logger.Errorf("Failed to create room %s: %v", r.roomNo, err)
		} else {
			logger.Infof("Test room created: %s", r.roomNo)
		}
	}
}
