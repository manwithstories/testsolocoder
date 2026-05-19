package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateVerifyCode() string {
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	return code
}

func SendVerifyEmail(email, code string) error {
	fmt.Printf("验证码已发送到 %s: %s (这是模拟邮件，实际项目中请接入邮件服务)\n", email, code)
	return nil
}
