package generator

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
)

const (
	LowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	UppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits           = "0123456789"
	SpecialChars     = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	AmbiguousChars   = "0O1lI"
)

type Options struct {
	Length          int
	IncludeLower    bool
	IncludeUpper    bool
	IncludeDigits   bool
	IncludeSpecial  bool
	ExcludeAmbiguous bool
}

func NewDefaultOptions() *Options {
	return &Options{
		Length:          16,
		IncludeLower:    true,
		IncludeUpper:    true,
		IncludeDigits:   true,
		IncludeSpecial:  true,
		ExcludeAmbiguous: true,
	}
}

func Generate(opts *Options) (string, error) {
	if opts == nil {
		opts = NewDefaultOptions()
	}

	if opts.Length < 8 {
		return "", errors.New("password length must be at least 8")
	}

	if opts.Length > 128 {
		return "", errors.New("password length must not exceed 128")
	}

	charset := ""
	if opts.IncludeLower {
		charset += LowercaseLetters
	}
	if opts.IncludeUpper {
		charset += UppercaseLetters
	}
	if opts.IncludeDigits {
		charset += Digits
	}
	if opts.IncludeSpecial {
		charset += SpecialChars
	}

	if opts.ExcludeAmbiguous {
		for _, c := range AmbiguousChars {
			charset = strings.ReplaceAll(charset, string(c), "")
		}
	}

	if len(charset) == 0 {
		return "", errors.New("no character types selected")
	}

	password := make([]byte, opts.Length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < opts.Length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		password[i] = charset[randomIndex.Int64()]
	}

	if !meetsRequirements(string(password), opts) {
		return Generate(opts)
	}

	return string(password), nil
}

func meetsRequirements(password string, opts *Options) bool {
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, c := range password {
		switch {
		case strings.ContainsRune(LowercaseLetters, c):
			hasLower = true
		case strings.ContainsRune(UppercaseLetters, c):
			hasUpper = true
		case strings.ContainsRune(Digits, c):
			hasDigit = true
		case strings.ContainsRune(SpecialChars, c):
			hasSpecial = true
		}
	}

	if opts.IncludeLower && !hasLower {
		return false
	}
	if opts.IncludeUpper && !hasUpper {
		return false
	}
	if opts.IncludeDigits && !hasDigit {
		return false
	}
	if opts.IncludeSpecial && !hasSpecial {
		return false
	}

	return true
}

func GeneratePhrase(wordCount int) (string, error) {
	if wordCount < 3 {
		return "", errors.New("word count must be at least 3")
	}
	if wordCount > 10 {
		return "", errors.New("word count must not exceed 10")
	}

	words := []string{
		"alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "golf", "hotel",
		"india", "juliett", "kilo", "lima", "mike", "november", "oscar", "papa",
		"quebec", "romeo", "sierra", "tango", "uniform", "victor", "whiskey", "xray",
		"yankee", "zulu", "apple", "banana", "cherry", "dragon", "eagle", "forest",
		"galaxy", "harbor", "island", "jungle", "knight", "lemon", "mountain", "nebula",
		"ocean", "penguin", "quantum", "river", "sunset", "thunder", "umbrella", "violet",
		"wizard", "xenon", "yellow", "zebra",
	}

	phrase := make([]string, wordCount)
	wordsLen := big.NewInt(int64(len(words)))

	for i := 0; i < wordCount; i++ {
		randomIndex, err := rand.Int(rand.Reader, wordsLen)
		if err != nil {
			return "", err
		}
		phrase[i] = words[randomIndex.Int64()]
	}

	return strings.Join(phrase, "-"), nil
}
