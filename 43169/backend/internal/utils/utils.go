package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
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

func RandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func RandomCode(length int) string {
	const digits = "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = digits[time.Now().UnixNano()%int64(len(digits))]
	}
	return string(b)
}

func MD5Hash(str string) string {
	h := md5.Sum([]byte(str))
	return hex.EncodeToString(h[:])
}

func MaskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}

func MaskIDCard(idCard string) string {
	if len(idCard) < 8 {
		return idCard
	}
	return idCard[:4] + "********" + idCard[len(idCard)-4:]
}

func ValidatePhone(phone string) bool {
	reg := `^1[3-9]\d{9}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}

func ValidateEmail(email string) bool {
	reg := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(email)
}

func ValidateIDCard(idCard string) bool {
	reg := `^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(idCard)
}

func FilterSensitiveWords(content string, words []string) string {
	result := content
	for _, w := range words {
		replacer := strings.Repeat("*", len(w))
		result = strings.ReplaceAll(result, w, replacer)
	}
	return result
}

func CalcAge(birthday *time.Time) int {
	if birthday == nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - birthday.Year()
	if now.YearDay() < birthday.YearDay() {
		age--
	}
	return age
}

func DistanceKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := (lat2 - lat1) * 3.141592653589793 / 180.0
	dLon := (lon2 - lon1) * 3.141592653589793 / 180.0
	lat1 = lat1 * 3.141592653589793 / 180.0
	lat2 = lat2 * 3.141592653589793 / 180.0
	a := dLat/2*dLat/2 + dLon/2*dLon/2*lat1*lat2
	c := 2 * atan2(sqrt(a), sqrt(1-a))
	return earthRadiusKm * c
}

func atan2(y, x float64) float64 {
	if x > 0 {
		return atan(y / x)
	} else if x < 0 {
		if y >= 0 {
			return atan(y/x) + 3.141592653589793
		}
		return atan(y/x) - 3.141592653589793
	} else {
		if y > 0 {
			return 3.141592653589793 / 2
		} else if y < 0 {
			return -3.141592653589793 / 2
		}
		return 0
	}
}

func atan(x float64) float64 {
	if x > 1.0 {
		return 3.141592653589793/2 - atan(1.0/x)
	} else if x < -1.0 {
		return -3.141592653589793/2 - atan(1.0/x)
	}
	x2 := x * x
	return x - x*x2/3 + x*x2*x2/5 - x*x2*x2*x2/7 + x*x2*x2*x2*x2/9
}

func sqrt(x float64) float64 {
	z := 1.0
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}

func Paginate(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func SafeDiv(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
