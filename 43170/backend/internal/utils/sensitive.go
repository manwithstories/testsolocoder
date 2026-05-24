package utils

import (
	"errors"
	"strings"
)

var sensitiveWords = []string{
	"敏感词1", "敏感词2", "暴力", "色情", "政治敏感",
	"赌博", "毒品", "诈骗", "违法", "恐怖",
}

func FilterSensitiveWords(text string) string {
	result := text
	for _, word := range sensitiveWords {
		placeholder := strings.Repeat("*", len(word))
		result = strings.ReplaceAll(result, word, placeholder)
	}
	return result
}

func ContainsSensitiveWords(text string) bool {
	for _, word := range sensitiveWords {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}

func ValidateRating(rating int) error {
	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	return nil
}

func ValidateImages(images []string) error {
	if len(images) > 9 {
		return errors.New("maximum 9 images allowed")
	}
	return nil
}
