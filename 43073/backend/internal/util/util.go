package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateOrderNo() string {
	now := time.Now()
	return now.Format("20060102150405") + randomString(6)
}

func randomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateQRCodeContent(orderNo string, itemID uint) string {
	data := orderNo + "_" + fmt.Sprintf("%d", itemID)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func GenerateCouponCode() string {
	return "CPN" + randomString(8)
}
