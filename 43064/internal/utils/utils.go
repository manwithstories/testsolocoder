package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var variablePattern = regexp.MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`)

func GenerateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func ParseTemplateVariables(content string) []string {
	matches := variablePattern.FindAllStringSubmatch(content, -1)
	variables := make([]string, 0, len(matches))
	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 && !seen[match[1]] {
			variables = append(variables, match[1])
			seen[match[1]] = true
		}
	}
	return variables
}

func RenderTemplate(content string, variables map[string]interface{}) (string, error) {
	result := content
	for name, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", name)
		strValue := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, strValue)
	}

	remaining := variablePattern.FindAllString(result, -1)
	if len(remaining) > 0 {
		return result, fmt.Errorf("missing variables: %v", remaining)
	}

	return result, nil
}

func ToJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(data)
}

func FromJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

func ComputeHMAC256(message, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func SignRequest(payload, secret string) string {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	message := fmt.Sprintf("%s.%s", timestamp, payload)
	signature := ComputeHMAC256(message, secret)
	return fmt.Sprintf("%s.%s", timestamp, signature)
}

func VerifySignature(payload, signed, secret string) bool {
	parts := strings.Split(signed, ".")
	if len(parts) != 2 {
		return false
	}
	timestamp := parts[0]
	signature := parts[1]
	message := fmt.Sprintf("%s.%s", timestamp, payload)
	expected := ComputeHMAC256(message, secret)
	return hmac.Equal([]byte(signature), []byte(expected))
}

func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func IsValidPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func DeduplicateStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

func CalculateBackoff(attempt int, initialBackoff, maxBackoff int, multiplier float64) time.Duration {
	backoff := float64(initialBackoff)
	for i := 0; i < attempt; i++ {
		backoff *= multiplier
		if backoff > float64(maxBackoff) {
			backoff = float64(maxBackoff)
			break
		}
	}
	return time.Duration(backoff) * time.Millisecond
}

func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

func GetStackTrace(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%+v", err)
}

func Pointer[T any](v T) *T {
	return &v
}

func ValueOrDefault[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}
