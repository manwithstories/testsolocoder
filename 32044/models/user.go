package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;uniqueIndex;not null" json:"username" binding:"required,min=3,max=50"`
	Password  string         `gorm:"size:255;not null" json:"-" binding:"required,min=6,max=50"`
	Email     string         `gorm:"size:100;uniqueIndex" json:"email" binding:"omitempty,email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Accounts    []Account    `gorm:"foreignKey:UserID" json:"-"`
	Categories  []Category   `gorm:"foreignKey:UserID" json:"-"`
	Transactions []Transaction `gorm:"foreignKey:UserID" json:"-"`
	Budgets     []Budget     `gorm:"foreignKey:UserID" json:"-"`
}

func (u *User) HashPassword() error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedBytes)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
