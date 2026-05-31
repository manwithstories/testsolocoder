package utils

import (
	"math"
	"translation-platform/internal/models"
)

const (
	BasePricePerWord = 0.15
)

var LanguageDifficulty = map[string]float64{
	"zh-en": 1.2,
	"en-zh": 1.2,
	"zh-ja": 1.3,
	"ja-zh": 1.3,
	"zh-ko": 1.3,
	"ko-zh": 1.3,
	"zh-fr": 1.5,
	"fr-zh": 1.5,
	"zh-de": 1.5,
	"de-zh": 1.5,
	"zh-es": 1.5,
	"es-zh": 1.5,
	"zh-ru": 1.6,
	"ru-zh": 1.6,
	"en-ja": 1.4,
	"ja-en": 1.4,
	"en-fr": 1.1,
	"fr-en": 1.1,
	"en-de": 1.1,
	"de-en": 1.1,
}

var UrgencyMultiplier = map[models.UrgencyLevel]float64{
	models.UrgencyNormal: 1.0,
	models.UrgencyFast:   1.3,
	models.UrgencyUrgent: 1.6,
}

var DomainDifficulty = map[string]float64{
	"general":    1.0,
	"technology": 1.2,
	"legal":      1.4,
	"medical":    1.5,
	"finance":    1.3,
	"academic":   1.2,
	"marketing":  1.1,
	"literary":   1.3,
}

func CalculateFee(wordCount int, sourceLang, targetLang string, urgency models.UrgencyLevel, domains []string) map[string]float64 {
	languagePair := sourceLang + "-" + targetLang
	languageDiff := LanguageDifficulty[languagePair]
	if languageDiff == 0 {
		languageDiff = 1.0
	}

	urgencyMult := UrgencyMultiplier[urgency]
	if urgencyMult == 0 {
		urgencyMult = 1.0
	}

	domainDiff := 1.0
	if len(domains) > 0 {
		totalDiff := 0.0
		for _, d := range domains {
			if diff, ok := DomainDifficulty[d]; ok {
				totalDiff += diff
			} else {
				totalDiff += 1.0
			}
		}
		domainDiff = totalDiff / float64(len(domains))
	}

	baseAmount := float64(wordCount) * BasePricePerWord
	baseWithLanguage := baseAmount * languageDiff
	urgencyFee := baseWithLanguage * (urgencyMult - 1.0)
	difficultyFee := baseWithLanguage * (domainDiff - 1.0)
	totalAmount := baseWithLanguage + urgencyFee + difficultyFee

	return map[string]float64{
		"base_amount":     roundFloat(baseAmount, 2),
		"urgency_fee":     roundFloat(urgencyFee, 2),
		"difficulty_fee":  roundFloat(difficultyFee, 2),
		"total_amount":    roundFloat(totalAmount, 2),
		"language_diff":   languageDiff,
		"urgency_mult":    urgencyMult,
		"domain_diff":     roundFloat(domainDiff, 2),
	}
}

func roundFloat(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func CalculateTranslatorScore(t *models.User, project *models.Project) float64 {
	if t == nil || project == nil {
		return 0
	}

	score := 0.0

	score += t.Rating * 20

	languageMatch := false
	for _, lp := range t.LanguagePairs {
		if lp.SourceLang == project.SourceLang && lp.TargetLang == project.TargetLang {
			languageMatch = true
			break
		}
	}
	if languageMatch {
		score += 30
	}

	expertiseMatch := 0
	for _, tag := range t.ExpertiseTags {
		for _, pt := range project.ExpertiseTags {
			if tag.ID == pt.ID {
				expertiseMatch++
			}
		}
	}
	if len(project.ExpertiseTags) > 0 {
		score += float64(expertiseMatch) / float64(len(project.ExpertiseTags)) * 20
	}

	completionRatio := 0.0
	if t.CurrentWorkload < t.DailyCapacity {
		completionRatio = 1.0 - float64(t.CurrentWorkload)/float64(t.DailyCapacity)
	}
	score += completionRatio * 15

	score += float64(t.CompletedCount) / 100.0 * 15

	return roundFloat(score, 2)
}
