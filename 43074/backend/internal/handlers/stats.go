package handlers

import (
	"net/http"
	"time"

	"booklibrary/internal/database"
	"booklibrary/internal/models"
	"booklibrary/internal/utils"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct{}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

type YearlyStats struct {
	Year      int `json:"year"`
	Completed int `json:"completed"`
	Reading   int `json:"reading"`
	Abandoned int `json:"abandoned"`
	Total     int `json:"total"`
}

type MonthlyStats struct {
	Month     string `json:"month"`
	Completed int    `json:"completed"`
	PagesRead int    `json:"pages_read"`
}

type HeatmapData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ReadingDurationStats struct {
	Range    string `json:"range"`
	Count    int    `json:"count"`
	Category string `json:"category"`
}

func (h *StatsHandler) GetOverview(c *gin.Context) {
	var totalBooks int64
	database.DB.Model(&models.Book{}).Count(&totalBooks)

	var readingBooks int64
	database.DB.Model(&models.Book{}).Where("reading_status = ?", models.StatusReading).Count(&readingBooks)

	var completedBooks int64
	database.DB.Model(&models.Book{}).Where("reading_status = ?", models.StatusCompleted).Count(&completedBooks)

	var totalPagesRead int64
	rows, _ := database.DB.Model(&models.Book{}).
		Where("reading_status = ?", models.StatusCompleted).
		Select("COALESCE(SUM(total_pages), 0)").
		Rows()
	if rows.Next() {
		rows.Scan(&totalPagesRead)
	}
	rows.Close()

	var totalReadingTime int64
	database.DB.Model(&models.Book{}).Select("COALESCE(SUM(total_read_time), 0)").Scan(&totalReadingTime)

	var currentlyBorrowed int64
	database.DB.Model(&models.BorrowRecord{}).Where("return_date IS NULL").Count(&currentlyBorrowed)

	c.JSON(http.StatusOK, gin.H{
		"total_books":       totalBooks,
		"reading_books":     readingBooks,
		"completed_books":   completedBooks,
		"total_pages_read":  totalPagesRead,
		"total_read_time":   totalReadingTime,
		"currently_borrowed": currentlyBorrowed,
	})
}

func (h *StatsHandler) GetYearlyTrend(c *gin.Context) {
	year := utils.GetIntQuery(c, "year", time.Now().Year())

	startDate := utils.StartOfYear(year)
	endDate := utils.EndOfYear(year)

	var completed []models.Book
	database.DB.Where("reading_status = ? AND end_date BETWEEN ? AND ?",
		models.StatusCompleted, startDate, endDate).
		Find(&completed)

	monthly := make([]MonthlyStats, 12)
	for i := 0; i < 12; i++ {
		monthly[i] = MonthlyStats{
			Month:     time.Month(i + 1).String(),
			Completed: 0,
			PagesRead: 0,
		}
	}

	for _, book := range completed {
		if book.EndDate != nil {
			month := int(book.EndDate.Month()) - 1
			monthly[month].Completed++
			monthly[month].PagesRead += book.TotalPages
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"year":    year,
		"monthly": monthly,
	})
}

func (h *StatsHandler) GetReadingHeatmap(c *gin.Context) {
	year := utils.GetIntQuery(c, "year", time.Now().Year())

	startDate := utils.StartOfYear(year)
	endDate := utils.EndOfYear(year)

	var completed []models.Book
	database.DB.Where("reading_status = ? AND end_date BETWEEN ? AND ?",
		models.StatusCompleted, startDate, endDate).
		Find(&completed)

	heatmap := make(map[string]int)

	for _, book := range completed {
		if book.EndDate != nil {
			dateStr := book.EndDate.Format("2006-01-02")
			heatmap[dateStr]++
		}
		if book.StartDate != nil {
			dateStr := book.StartDate.Format("2006-01-02")
			heatmap[dateStr]++
		}
	}

	var result []HeatmapData
	for date, count := range heatmap {
		result = append(result, HeatmapData{
			Date:  date,
			Count: count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"year": year,
		"data": result,
	})
}

func (h *StatsHandler) GetDurationDistribution(c *gin.Context) {
	var completed []models.Book
	database.DB.Where("reading_status = ? AND total_read_time > 0",
		models.StatusCompleted).
		Find(&completed)

	ranges := []struct {
		Name string
		Min  int
		Max  int
	}{
		{"< 5小时", 0, 5 * 60},
		{"5-10小时", 5 * 60, 10 * 60},
		{"10-20小时", 10 * 60, 20 * 60},
		{"20-50小时", 20 * 60, 50 * 60},
		{"> 50小时", 50 * 60, 999999},
	}

	result := make([]ReadingDurationStats, len(ranges))
	for i, r := range ranges {
		result[i] = ReadingDurationStats{
			Range:    r.Name,
			Count:    0,
			Category: "阅读时长",
		}
	}

	for _, book := range completed {
		for i, r := range ranges {
			if book.TotalReadTime >= r.Min && book.TotalReadTime < r.Max {
				result[i].Count++
				break
			}
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *StatsHandler) GetCategoryStats(c *gin.Context) {
	type CategoryCount struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	var results []CategoryCount
	database.DB.Table("categories").
		Select("categories.name, COUNT(book_categories.book_id) as count").
		Joins("LEFT JOIN book_categories ON book_categories.category_id = categories.id").
		Group("categories.id, categories.name").
		Having("count > 0").
		Order("count desc").
		Scan(&results)

	c.JSON(http.StatusOK, results)
}

func (h *StatsHandler) GetTagStats(c *gin.Context) {
	type TagCount struct {
		Name  string `json:"name"`
		Color string `json:"color"`
		Count int    `json:"count"`
	}

	var results []TagCount
	database.DB.Table("tags").
		Select("tags.name, tags.color, COUNT(book_tags.book_id) as count").
		Joins("LEFT JOIN book_tags ON book_tags.tag_id = tags.id").
		Group("tags.id, tags.name, tags.color").
		Having("count > 0").
		Order("count desc").
		Scan(&results)

	c.JSON(http.StatusOK, results)
}
