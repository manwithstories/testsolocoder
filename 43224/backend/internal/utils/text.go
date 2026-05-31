package utils

import (
	"math"
	"regexp"
	"strings"
	"unicode"
)

func CountWords(text string) int {
	if text == "" {
		return 0
	}

	re := regexp.MustCompile(`[\p{Han}]+|\b[\w']+\b`)
	matches := re.FindAllString(text, -1)

	count := 0
	for _, m := range matches {
		if isChinese(m) {
			count += len([]rune(m))
		} else {
			count++
		}
	}

	return count
}

func isChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func CalculateSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	if s1 == "" || s2 == "" {
		return 0.0
	}

	r1 := []rune(s1)
	r2 := []rune(s2)

	len1 := len(r1)
	len2 := len(r2)

	if len1 == 0 || len2 == 0 {
		return 0.0
	}

	prev := make([]int, len2+1)
	curr := make([]int, len2+1)

	for i := 1; i <= len1; i++ {
		curr[0] = i
		for j := 1; j <= len2; j++ {
			cost := 1
			if r1[i-1] == r2[j-1] {
				cost = 0
			}
			curr[j] = minInt(prev[j]+1, curr[j-1]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}

	distance := prev[len2]
	maxLen := float64(maxInt(len1, len2))
	similarity := 1.0 - float64(distance)/maxLen

	return math.Round(similarity*100) / 100
}

func minInt(nums ...int) int {
	min := nums[0]
	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
	}
	return min
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ExtractSegments(text string, maxLen int) []string {
	if text == "" {
		return nil
	}

	var segments []string
	sentences := splitSentences(text)

	var current strings.Builder
	currentLen := 0

	for _, sentence := range sentences {
		sLen := len([]rune(sentence))
		if currentLen+sLen > maxLen && currentLen > 0 {
			segments = append(segments, strings.TrimSpace(current.String()))
			current.Reset()
			currentLen = 0
		}
		current.WriteString(sentence)
		currentLen += sLen
	}

	if currentLen > 0 {
		segments = append(segments, strings.TrimSpace(current.String()))
	}

	return segments
}

func splitSentences(text string) []string {
	re := regexp.MustCompile(`[^.!?。！？]+[.!?。！？]+|[^.!?。！？]+$`)
	matches := re.FindAllString(text, -1)
	if len(matches) == 0 {
		return []string{text}
	}
	return matches
}

func SanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	re := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
	name = re.ReplaceAllString(name, "_")

	if len(name) > 200 {
		name = name[:200]
	}

	return name
}
