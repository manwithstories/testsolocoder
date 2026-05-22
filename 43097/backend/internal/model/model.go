package model

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&MemberLevel{},
		&Member{},
		&MemberPointsLog{},
		&RoomType{},
		&Room{},
		&Booking{},
		&CheckIn{},
		&Payment{},
	)
}
