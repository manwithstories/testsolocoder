package handlers

import (
	"net/http"
	"time"

	"booklibrary/internal/database"
	"booklibrary/internal/errors"
	"booklibrary/internal/logger"
	"booklibrary/internal/models"
	"booklibrary/internal/utils"

	"github.com/gin-gonic/gin"
)

type GoalHandler struct{}

func NewGoalHandler() *GoalHandler {
	return &GoalHandler{}
}

type CreateGoalRequest struct {
	Year        int  `json:"year" binding:"required"`
	Month       *int `json:"month"`
	TargetBooks int  `json:"target_books"`
	TargetPages int  `json:"target_pages"`
}

type UpdateGoalRequest struct {
	TargetBooks int `json:"target_books"`
	TargetPages int `json:"target_pages"`
}

type GoalProgress struct {
	models.ReadingGoal
	CompletedBooks int `json:"completed_books"`
	CompletedPages int `json:"completed_pages"`
	BookProgress   float64 `json:"book_progress"`
	PageProgress   float64 `json:"page_progress"`
}

func (h *GoalHandler) GetGoals(c *gin.Context) {
	year := utils.GetIntQuery(c, "year", 0)

	var goals []models.ReadingGoal
	query := database.DB

	if year > 0 {
		query = query.Where("year = ?", year)
	}

	result := query.Order("year desc, month desc").Find(&goals)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, goals)
}

func (h *GoalHandler) GetGoal(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var goal models.ReadingGoal
	result := database.DB.First(&goal, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	progress := calculateGoalProgress(goal)
	c.JSON(http.StatusOK, progress)
}

func (h *GoalHandler) CreateGoal(c *gin.Context) {
	var req CreateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrValidation)
		return
	}

	var existing models.ReadingGoal
	query := database.DB.Where("year = ?", req.Year)
	if req.Month != nil {
		query = query.Where("month = ?", *req.Month)
	} else {
		query = query.Where("month IS NULL")
	}

	if result := query.First(&existing); result.Error == nil {
		errors.ErrorResponse(c, http.StatusConflict, "该目标已存在")
		return
	}

	goal := models.ReadingGoal{
		Year:        req.Year,
		Month:       req.Month,
		TargetBooks: req.TargetBooks,
		TargetPages: req.TargetPages,
	}

	result := database.DB.Create(&goal)
	if result.Error != nil {
		logger.Errorf("Create goal failed: %v", result.Error)
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}

	logger.Infof("Goal created: %d", goal.ID)
	c.JSON(http.StatusCreated, goal)
}

func (h *GoalHandler) UpdateGoal(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var goal models.ReadingGoal
	result := database.DB.First(&goal, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	var req UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrValidation)
		return
	}

	goal.TargetBooks = req.TargetBooks
	goal.TargetPages = req.TargetPages

	database.DB.Save(&goal)
	logger.Infof("Goal updated: %d", goal.ID)
	c.JSON(http.StatusOK, goal)
}

func (h *GoalHandler) DeleteGoal(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	result := database.DB.Delete(&models.ReadingGoal{}, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}
	if result.RowsAffected == 0 {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	logger.Infof("Goal deleted: %d", id)
	c.JSON(http.StatusNoContent, nil)
}

func (h *GoalHandler) GetYearlyGoalProgress(c *gin.Context) {
	year := utils.GetIntQuery(c, "year", time.Now().Year())

	var yearlyGoal models.ReadingGoal
	database.DB.Where("year = ? AND month IS NULL", year).First(&yearlyGoal)

	var monthlyGoals []models.ReadingGoal
	database.DB.Where("year = ? AND month IS NOT NULL", year).Order("month asc").Find(&monthlyGoals)

	yearlyProgress := calculateGoalProgress(yearlyGoal)

	var monthlyProgress []GoalProgress
	for _, mg := range monthlyGoals {
		monthlyProgress = append(monthlyProgress, calculateGoalProgress(mg))
	}

	c.JSON(http.StatusOK, gin.H{
		"yearly":  yearlyProgress,
		"monthly": monthlyProgress,
	})
}

func calculateGoalProgress(goal models.ReadingGoal) GoalProgress {
	progress := GoalProgress{
		ReadingGoal: goal,
	}

	var startDate, endDate time.Time
	if goal.Month != nil {
		startDate = utils.StartOfMonth(goal.Year, *goal.Month)
		endDate = utils.EndOfMonth(goal.Year, *goal.Month)
	} else {
		startDate = utils.StartOfYear(goal.Year)
		endDate = utils.EndOfYear(goal.Year)
	}

	var completedBooks []models.Book
	database.DB.Where("reading_status = ? AND end_date BETWEEN ? AND ?",
		models.StatusCompleted, startDate, endDate).
		Find(&completedBooks)

	progress.CompletedBooks = len(completedBooks)
	for _, book := range completedBooks {
		progress.CompletedPages += book.TotalPages
	}

	if goal.TargetBooks > 0 {
		progress.BookProgress = float64(progress.CompletedBooks) / float64(goal.TargetBooks) * 100
	}
	if goal.TargetPages > 0 {
		progress.PageProgress = float64(progress.CompletedPages) / float64(goal.TargetPages) * 100
	}

	return progress
}
