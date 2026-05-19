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
	"gorm.io/gorm"
)

type BorrowHandler struct{}

func NewBorrowHandler() *BorrowHandler {
	return &BorrowHandler{}
}

type CreateBorrowRequest struct {
	BookID             uint64 `json:"book_id" binding:"required"`
	BorrowerName       string `json:"borrower_name" binding:"required"`
	BorrowerPhone      string `json:"borrower_phone"`
	BorrowerEmail      string `json:"borrower_email"`
	BorrowDate         string `json:"borrow_date"`
	ExpectedReturnDate string `json:"expected_return_date"`
	Notes              string `json:"notes"`
}

type ReturnBookRequest struct {
	ReturnDate string `json:"return_date"`
}

func (h *BorrowHandler) GetBorrows(c *gin.Context) {
	status := utils.GetStringQuery(c, "status", "")

	var borrows []models.BorrowRecord
	query := database.DB.Preload("Book")

	if status == "borrowed" {
		query = query.Where("return_date IS NULL")
	} else if status == "returned" {
		query = query.Where("return_date IS NOT NULL")
	}

	result := query.Order("borrow_date desc").Find(&borrows)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, borrows)
}

func (h *BorrowHandler) GetBorrow(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var borrow models.BorrowRecord
	result := database.DB.Preload("Book").First(&borrow, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, borrow)
}

func (h *BorrowHandler) CreateBorrow(c *gin.Context) {
	var req CreateBorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrValidation)
		return
	}

	var book models.Book
	result := database.DB.First(&book, req.BookID)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	var existingBorrow models.BorrowRecord
	result = database.DB.Where("book_id = ? AND return_date IS NULL", req.BookID).First(&existingBorrow)
	if result.Error == nil {
		errors.HandleError(c, errors.ErrBookBorrowed)
		return
	}

	borrowDate := time.Now()
	if req.BorrowDate != "" {
		parsed, err := utils.ParseDate(req.BorrowDate)
		if err == nil {
			borrowDate = parsed
		}
	}

	var expectedReturnDate *time.Time
	if req.ExpectedReturnDate != "" {
		parsed, err := utils.ParseDate(req.ExpectedReturnDate)
		if err == nil {
			expectedReturnDate = &parsed
		}
	}

	borrow := models.BorrowRecord{
		BookID:             req.BookID,
		BorrowerName:       req.BorrowerName,
		BorrowerPhone:      req.BorrowerPhone,
		BorrowerEmail:      req.BorrowerEmail,
		BorrowDate:         borrowDate,
		ExpectedReturnDate: expectedReturnDate,
		Status:             "borrowed",
		Notes:              req.Notes,
	}

	tx := database.DB.Begin()
	if err := tx.Create(&borrow).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Create borrow failed: %v", err)
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}
	tx.Commit()

	logger.Infof("Book %d borrowed by %s", req.BookID, req.BorrowerName)
	c.JSON(http.StatusCreated, borrow)
}

func (h *BorrowHandler) ReturnBook(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var borrow models.BorrowRecord
	result := database.DB.First(&borrow, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	if borrow.ReturnDate != nil {
		errors.ErrorResponse(c, http.StatusBadRequest, "该书已归还")
		return
	}

	var req ReturnBookRequest
	c.ShouldBindJSON(&req)

	returnDate := time.Now()
	if req.ReturnDate != "" {
		parsed, err := utils.ParseDate(req.ReturnDate)
		if err == nil {
			returnDate = parsed
		}
	}

	borrow.ReturnDate = &returnDate
	borrow.Status = "returned"
	database.DB.Save(&borrow)

	logger.Infof("Book %d returned", borrow.BookID)
	c.JSON(http.StatusOK, borrow)
}

func (h *BorrowHandler) DeleteBorrow(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	result := database.DB.Delete(&models.BorrowRecord{}, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}
	if result.RowsAffected == 0 {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	logger.Infof("Borrow record deleted: %d", id)
	c.JSON(http.StatusNoContent, nil)
}

func (h *BorrowHandler) GetOverdue(c *gin.Context) {
	now := time.Now()
	var overdue []models.BorrowRecord
	database.DB.Preload("Book").
		Where("return_date IS NULL AND expected_return_date < ?", now).
		Order("expected_return_date asc").
		Find(&overdue)

	c.JSON(http.StatusOK, overdue)
}

func (h *BorrowHandler) GetBorrowByBook(c *gin.Context) {
	bookID, err := utils.GetUintParam(c, "bookId")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var borrow models.BorrowRecord
	result := database.DB.Where("book_id = ? AND return_date IS NULL", bookID).First(&borrow)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, nil)
			return
		}
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, borrow)
}
