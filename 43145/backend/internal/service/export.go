package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"survey-platform/internal/dto"
	"survey-platform/internal/model"
	"survey-platform/internal/repository"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExportService struct {
	surveyRepo     *repository.SurveyRepository
	responseRepo   *repository.ResponseRepository
	questionRepo   *repository.QuestionRepository
	statisticsSvc  *StatisticsService
}

func NewExportService(
	surveyRepo *repository.SurveyRepository,
	responseRepo *repository.ResponseRepository,
	questionRepo *repository.QuestionRepository,
	statisticsSvc *StatisticsService,
) *ExportService {
	return &ExportService{
		surveyRepo:    surveyRepo,
		responseRepo:  responseRepo,
		questionRepo:  questionRepo,
		statisticsSvc: statisticsSvc,
	}
}

func (s *ExportService) ExportExcel(surveyID uint) ([]byte, error) {
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return nil, err
	}
	if survey == nil {
		return nil, fmt.Errorf("survey not found")
	}

	questions, err := s.questionRepo.FindBySurveyID(surveyID)
	if err != nil {
		return nil, err
	}

	responses, _, err := s.responseRepo.List(surveyID, 1, 10000, 2, nil, nil, "", "created_at", "desc")
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()

	sheetName := "原始数据"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"答卷ID", "提交时间", "用时(秒)"}
	for _, q := range questions {
		headers = append(headers, fmt.Sprintf("Q%d: %s", q.OrderIndex+1, q.Title))
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#E0E0E0"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s1", columnLetter(len(headers))), style)

	for i, response := range responses {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), response.ID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), response.CompleteTime.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), response.Duration)

		for j, q := range questions {
			cell, _ := excelize.CoordinatesToCellName(j+4, row)
			var value string
			for _, answer := range response.Answers {
				if answer.QuestionID == q.ID {
					value = s.formatAnswerValue(q, &answer)
					break
				}
			}
			f.SetCellValue(sheetName, cell, value)
		}
	}

	statsSheet := "统计概览"
	f.NewSheet(statsSheet)

	statQuery := &dto.StatisticsQuery{
		SurveyID: surveyID,
	}
	stats, _ := s.statisticsSvc.GetStatistics(statQuery)

	statHeaders := []string{"指标", "数值"}
	f.SetCellValue(statsSheet, "A1", statHeaders[0])
	f.SetCellValue(statsSheet, "B1", statHeaders[1])
	f.SetCellStyle(statsSheet, "A1", "B1", style)

	statData := [][]interface{}{
		{"问卷标题", survey.Title},
		{"总答卷数", stats.TotalResponses},
		{"已完成", stats.CompletedCount},
		{"进行中", stats.InProgressCount},
		{"已放弃", stats.AbandonedCount},
		{"完成率(%)", fmt.Sprintf("%.2f", stats.CompletionRate)},
		{"平均用时(秒)", fmt.Sprintf("%.2f", stats.AvgDuration)},
	}

	for i, data := range statData {
		row := i + 2
		f.SetCellValue(statsSheet, fmt.Sprintf("A%d", row), data[0])
		f.SetCellValue(statsSheet, fmt.Sprintf("B%d", row), data[1])
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *ExportService) ExportPDF(surveyID uint) (map[string]interface{}, error) {
	survey, err := s.surveyRepo.FindByID(surveyID)
	if err != nil {
		return nil, err
	}
	if survey == nil {
		return nil, fmt.Errorf("survey not found")
	}

	statQuery := &dto.StatisticsQuery{
		SurveyID: surveyID,
	}
	stats, err := s.statisticsSvc.GetStatistics(statQuery)
	if err != nil {
		return nil, err
	}

	report := map[string]interface{}{
		"title":       survey.Title,
		"description": survey.Description,
		"generatedAt": time.Now().Format("2006-01-02 15:04:05"),
		"summary": map[string]interface{}{
			"totalResponses":  stats.TotalResponses,
			"completedCount":  stats.CompletedCount,
			"inProgressCount": stats.InProgressCount,
			"abandonedCount":  stats.AbandonedCount,
			"completionRate":  fmt.Sprintf("%.2f%%", stats.CompletionRate),
			"avgDuration":     fmt.Sprintf("%.2f秒", stats.AvgDuration),
		},
		"timeDistribution": stats.TimeDistribution,
		"questionStats":    stats.QuestionStats,
	}

	return report, nil
}

func (s *ExportService) ExportChartImages(surveyID uint) ([]map[string]interface{}, error) {
	statQuery := &dto.StatisticsQuery{
		SurveyID: surveyID,
	}
	stats, err := s.statisticsSvc.GetStatistics(statQuery)
	if err != nil {
		return nil, err
	}

	var charts []map[string]interface{}

	for _, qs := range stats.QuestionStats {
		if qs.QuestionType == "single_choice" || qs.QuestionType == "multi_choice" {
			charts = append(charts, map[string]interface{}{
				"type":     "bar",
				"title":    qs.QuestionTitle,
				"questionId": qs.QuestionID,
				"data":     qs.OptionStats,
			})
		} else if qs.QuestionType == "rating" && qs.RatingStats != nil {
			charts = append(charts, map[string]interface{}{
				"type":     "distribution",
				"title":    qs.QuestionTitle,
				"questionId": qs.QuestionID,
				"data":     qs.RatingStats.Distribution,
			})
		} else if qs.QuestionType == "fill_in" && len(qs.WordCloud) > 0 {
			charts = append(charts, map[string]interface{}{
				"type":     "wordcloud",
				"title":    qs.QuestionTitle,
				"questionId": qs.QuestionID,
				"data":     qs.WordCloud,
			})
		}
	}

	return charts, nil
}

func (s *ExportService) formatAnswerValue(question *model.Question, answer *model.Answer) string {
	switch question.Type {
	case "single_choice", "multi_choice", "ranking":
		if answer.OptionID != nil {
			for _, opt := range question.Options {
				if opt.ID == *answer.OptionID {
					if opt.IsOther && answer.TextValue != "" {
						return fmt.Sprintf("其他: %s", answer.TextValue)
					}
					return opt.Text
				}
			}
		}
		return answer.TextValue
	case "fill_in":
		return answer.TextValue
	case "rating":
		return strconv.FormatFloat(answer.NumericValue, 'f', 0, 64)
	case "matrix":
		if answer.MatrixValues != "" {
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(answer.MatrixValues), &result); err == nil {
				jsonBytes, _ := json.Marshal(result)
				return string(jsonBytes)
			}
		}
		return answer.MatrixValues
	default:
		return answer.TextValue
	}
}

func columnLetter(col int) string {
	result := ""
	for col > 0 {
		col--
		result = string(rune('A'+col%26)) + result
		col /= 26
	}
	return result
}
