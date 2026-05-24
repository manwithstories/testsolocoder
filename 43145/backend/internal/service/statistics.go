package service

import (
	"errors"
	"fmt"
	"sort"
	"survey-platform/internal/dto"
	"survey-platform/internal/model"
	"survey-platform/internal/repository"
	"survey-platform/internal/utils"
	"strings"
	"time"
)

type StatisticsService struct {
	responseRepo *repository.ResponseRepository
	surveyRepo   *repository.SurveyRepository
	questionRepo *repository.QuestionRepository
}

func NewStatisticsService(
	responseRepo *repository.ResponseRepository,
	surveyRepo *repository.SurveyRepository,
	questionRepo *repository.QuestionRepository,
) *StatisticsService {
	return &StatisticsService{
		responseRepo: responseRepo,
		surveyRepo:   surveyRepo,
		questionRepo: questionRepo,
	}
}

func (s *StatisticsService) GetStatistics(query *dto.StatisticsQuery) (*dto.StatisticsResponse, error) {
	cacheKey := fmt.Sprintf("stats:%d:%v:%v:%s", query.SurveyID, query.StartDate, query.EndDate, query.Channel)
	var cached dto.StatisticsResponse
	if err := utils.CacheGetJSON(cacheKey, &cached); err == nil && cached.SurveyID > 0 {
		return &cached, nil
	}

	survey, err := s.surveyRepo.FindByID(query.SurveyID)
	if err != nil {
		return nil, err
	}
	if survey == nil {
		return nil, errors.New("survey not found")
	}

	responses, _, err := s.responseRepo.List(query.SurveyID, 1, 10000, 0, query.StartDate, query.EndDate, "", "created_at", "desc")
	if err != nil {
		return nil, err
	}

	questions, err := s.questionRepo.FindBySurveyID(query.SurveyID)
	if err != nil {
		return nil, err
	}

	var completedCount, inProgressCount, abandonedCount int
	var totalDuration int64
	timeDistMap := make(map[string]int)

	for _, r := range responses {
		switch r.Status {
		case 2:
			completedCount++
			totalDuration += int64(r.Duration)
		case 1:
			inProgressCount++
		case 3:
			abandonedCount++
		}

		dateKey := r.CreatedAt.Format("2006-01-02")
		timeDistMap[dateKey]++
	}

	completionRate := float64(0)
	if len(responses) > 0 {
		completionRate = float64(completedCount) / float64(len(responses)) * 100
	}

	avgDuration := float64(0)
	if completedCount > 0 {
		avgDuration = float64(totalDuration) / float64(completedCount)
	}

	var timeDistribution []dto.TimeDataPoint
	for date, count := range timeDistMap {
		timeDistribution = append(timeDistribution, dto.TimeDataPoint{
			Date:  date,
			Count: count,
		})
	}
	sort.Slice(timeDistribution, func(i, j int) bool {
		return timeDistribution[i].Date < timeDistribution[j].Date
	})

	var questionStats []dto.QuestionStat
	for _, q := range questions {
		stat := s.calculateQuestionStat(q, responses)
		questionStats = append(questionStats, stat)
	}

	result := &dto.StatisticsResponse{
		SurveyID:         query.SurveyID,
		TotalResponses:   len(responses),
		CompletedCount:   completedCount,
		InProgressCount:  inProgressCount,
		AbandonedCount:   abandonedCount,
		CompletionRate:   completionRate,
		AvgDuration:      avgDuration,
		QuestionStats:    questionStats,
		TimeDistribution: timeDistribution,
	}

	utils.CacheSetJSON(cacheKey, result, 5*time.Minute)
	return result, nil
}

func (s *StatisticsService) calculateQuestionStat(question *model.Question, responses []*model.Response) dto.QuestionStat {
	stat := dto.QuestionStat{
		QuestionID:    question.ID,
		QuestionTitle: question.Title,
		QuestionType:  question.Type,
	}

	optionCounts := make(map[uint]int)
	textValues := make([]string, 0)
	numericValues := make([]float64, 0)
	totalAnswered := 0

	for _, response := range responses {
		for _, answer := range response.Answers {
			if answer.QuestionID == question.ID {
				totalAnswered++
				if answer.OptionID != nil {
					optionCounts[*answer.OptionID]++
				}
				if answer.TextValue != "" {
					textValues = append(textValues, answer.TextValue)
				}
				if answer.NumericValue != 0 {
					numericValues = append(numericValues, answer.NumericValue)
				}
			}
		}
	}

	stat.ResponseCount = totalAnswered

	if question.Type == "single_choice" || question.Type == "multi_choice" || question.Type == "ranking" {
		var optionStats []dto.OptionStat
		for _, opt := range question.Options {
			count := optionCounts[opt.ID]
			percent := float64(0)
			if totalAnswered > 0 {
				percent = float64(count) / float64(totalAnswered) * 100
			}
			optionStats = append(optionStats, dto.OptionStat{
				OptionID: opt.ID,
				Text:     opt.Text,
				Count:    count,
				Percent:  percent,
			})
		}
		stat.OptionStats = optionStats
	}

	if question.Type == "fill_in" && len(textValues) > 0 {
		wordCounts := s.generateWordCloud(textValues)
		stat.WordCloud = wordCounts

		totalLen := 0
		maxLen := 0
		minLen := 10000
		for _, tv := range textValues {
			l := len(tv)
			totalLen += l
			if l > maxLen {
				maxLen = l
			}
			if l < minLen {
				minLen = l
			}
		}
		stat.TextStats = &dto.TextStat{
			AvgLength: float64(totalLen) / float64(len(textValues)),
			MaxLength: maxLen,
			MinLength: minLen,
		}
	}

	if question.Type == "rating" && len(numericValues) > 0 {
		sort.Float64s(numericValues)
		stat.RatingStats = &dto.RatingStat{
			Average:      calculateAverage(numericValues),
			Median:       calculateMedian(numericValues),
			StdDev:       calculateStdDev(numericValues),
			Min:          numericValues[0],
			Max:          numericValues[len(numericValues)-1],
			Distribution: s.calculateRatingDistribution(question, numericValues),
		}
	}

	return stat
}

func (s *StatisticsService) generateWordCloud(texts []string) []dto.WordCloudItem {
	stopWords := map[string]bool{
		"的": true, "了": true, "是": true, "我": true, "在": true,
		"有": true, "和": true, "就": true, "不": true, "人": true,
		"都": true, "一": true, "一个": true, "上": true, "也": true,
		"很": true, "到": true, "说": true, "要": true, "去": true,
		"the": true, "a": true, "an": true, "is": true, "are": true,
		"was": true, "were": true, "be": true, "been": true, "being": true,
		"have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "will": true, "would": true, "shall": true, "should": true,
		"may": true, "might": true, "can": true, "could": true, "must": true,
		"of": true, "in": true, "to": true, "for": true, "with": true,
		"on": true, "at": true, "by": true, "from": true, "as": true,
		"and": true, "or": true, "but": true, "if": true, "then": true,
		"so": true, "that": true, "this": true, "it": true, "he": true,
		"she": true, "they": true, "we": true, "you": true, "i": true,
		"me": true, "him": true, "her": true, "us": true, "them": true,
	}

	wordCounts := make(map[string]int)
	for _, text := range texts {
		words := strings.FieldsFunc(text, func(r rune) bool {
			return r == ' ' || r == ',' || r == '.' || r == '!' || r == '?' || r == '，' || r == '。' || r == '！' || r == '？'
		})
		for _, word := range words {
			w := strings.ToLower(word)
			if len(w) >= 2 && !stopWords[w] {
				wordCounts[w]++
			}
		}
	}

	var items []dto.WordCloudItem
	for word, count := range wordCounts {
		items = append(items, dto.WordCloudItem{
			Word:  word,
			Count: count,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	if len(items) > 50 {
		items = items[:50]
	}

	return items
}

func (s *StatisticsService) calculateRatingDistribution(question *model.Question, values []float64) map[int]int {
	dist := make(map[int]int)
	minVal := question.MinValue
	maxVal := question.MaxValue
	if maxVal == 0 {
		maxVal = 5
	}
	if minVal == 0 {
		minVal = 1
	}

	for i := minVal; i <= maxVal; i++ {
		dist[i] = 0
	}
	for _, v := range values {
		dist[int(v)]++
	}
	return dist
}

func calculateAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := float64(0)
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateMedian(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	if n%2 == 0 {
		return (values[n/2-1] + values[n/2]) / 2
	}
	return values[n/2]
}

func calculateStdDev(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	mean := calculateAverage(values)
	variance := float64(0)
	for _, v := range values {
		variance += (v - mean) * (v - mean)
	}
	return variance / float64(n)
}

func (s *StatisticsService) CrossAnalysis(query *dto.CrossAnalysisQuery) (*dto.CrossAnalysisResponse, error) {
	rowQ, err := s.questionRepo.FindByID(query.RowQuestionID)
	if err != nil {
		return nil, err
	}
	if rowQ == nil {
		return nil, errors.New("row question not found")
	}

	colQ, err := s.questionRepo.FindByID(query.ColQuestionID)
	if err != nil {
		return nil, err
	}
	if colQ == nil {
		return nil, errors.New("column question not found")
	}

	if rowQ.Type != "single_choice" || colQ.Type != "single_choice" {
		return nil, errors.New("cross analysis only supports single choice questions")
	}

	responses, _, err := s.responseRepo.List(query.SurveyID, 1, 10000, 2, nil, nil, "", "created_at", "desc")
	if err != nil {
		return nil, err
	}

	rowOptions := rowQ.Options
	colOptions := colQ.Options

	rowCounts := make([]int, len(rowOptions))
	colCounts := make([]int, len(colOptions))
	matrix := make([][]int, len(rowOptions))
	for i := range matrix {
		matrix[i] = make([]int, len(colOptions))
	}

	for _, response := range responses {
		var rowVal, colVal *uint
		for _, answer := range response.Answers {
			if answer.QuestionID == query.RowQuestionID && answer.OptionID != nil {
				rowVal = answer.OptionID
			}
			if answer.QuestionID == query.ColQuestionID && answer.OptionID != nil {
				colVal = answer.OptionID
			}
		}

		if rowVal != nil && colVal != nil {
			for i, opt := range rowOptions {
				if opt.ID == *rowVal {
					rowCounts[i]++
					for j, colOpt := range colOptions {
						if colOpt.ID == *colVal {
							colCounts[j]++
							matrix[i][j]++
						}
					}
				}
			}
		}
	}

	percentMatrix := make([][]float64, len(rowOptions))
	for i := range percentMatrix {
		percentMatrix[i] = make([]float64, len(colOptions))
		for j := range percentMatrix[i] {
			if rowCounts[i] > 0 {
				percentMatrix[i][j] = float64(matrix[i][j]) / float64(rowCounts[i]) * 100
			}
		}
	}

	chiSquare := s.calculateChiSquare(matrix, rowCounts, colCounts)
	pValue := 0.0
	significant := chiSquare > 3.841

	rowOptionTexts := make([]string, len(rowOptions))
	for i, opt := range rowOptions {
		rowOptionTexts[i] = opt.Text
	}
	colOptionTexts := make([]string, len(colOptions))
	for i, opt := range colOptions {
		colOptionTexts[i] = opt.Text
	}

	return &dto.CrossAnalysisResponse{
		RowQuestionID: query.RowQuestionID,
		ColQuestionID: query.ColQuestionID,
		RowOptions:    rowOptionTexts,
		ColOptions:    colOptionTexts,
		Matrix:        matrix,
		PercentMatrix: percentMatrix,
		ChiSquare:     chiSquare,
		PValue:        pValue,
		Significant:   significant,
	}, nil
}

func (s *StatisticsService) calculateChiSquare(matrix [][]int, rowCounts, colCounts []int) float64 {
	total := 0
	for _, count := range rowCounts {
		total += count
	}
	if total == 0 {
		return 0
	}

	chiSquare := float64(0)
	for i := range matrix {
		for j := range matrix[i] {
			expected := float64(rowCounts[i]*colCounts[j]) / float64(total)
			if expected > 0 {
				diff := float64(matrix[i][j]) - expected
				chiSquare += (diff * diff) / expected
			}
		}
	}
	return chiSquare
}
