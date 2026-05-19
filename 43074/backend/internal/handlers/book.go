package handlers

import (
	"net/http"
	"time"

	"booklibrary/internal/database"
	"booklibrary/internal/errors"
	"booklibrary/internal/logger"
	"booklibrary/internal/models"
	"booklibrary/internal/utils"
	isbnutil "booklibrary/internal/utils/isbn"
	uploadutil "booklibrary/internal/utils/upload"

	"github.com/gin-gonic/gin"
)

type BookHandler struct{}

func NewBookHandler() *BookHandler {
	return &BookHandler{}
}

type CreateBookRequest struct {
	Title       string   `json:"title" binding:"required"`
	Author      string   `json:"author"`
	Publisher   string   `json:"publisher"`
	ISBN        string   `json:"isbn"`
	Summary     string   `json:"summary"`
	TotalPages  int      `json:"total_pages"`
	TagIDs      []uint64 `json:"tag_ids"`
	CategoryIDs []uint64 `json:"category_ids"`
}

type UpdateBookRequest struct {
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Publisher   string   `json:"publisher"`
	ISBN        string   `json:"isbn"`
	Summary     string   `json:"summary"`
	TotalPages  int      `json:"total_pages"`
	TagIDs      []uint64 `json:"tag_ids"`
	CategoryIDs []uint64 `json:"category_ids"`
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	page := utils.GetIntQuery(c, "page", 1)
	pageSize := utils.GetIntQuery(c, "page_size", 20)
	search := utils.GetStringQuery(c, "search", "")
	status := utils.GetStringQuery(c, "status", "")
	tagIDs := c.QueryArray("tag_ids")
	categoryIDs := c.QueryArray("category_ids")
	sortBy := utils.GetStringQuery(c, "sort_by", "created_at")
	sortOrder := utils.GetStringQuery(c, "sort_order", "desc")

	query := database.DB.Model(&models.Book{}).Preload("Tags").Preload("Categories")

	if search != "" {
		query = query.Where("title LIKE ? OR author LIKE ? OR isbn LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if status != "" {
		query = query.Where("reading_status = ?", status)
	}

	if len(tagIDs) > 0 {
		query = query.Joins("JOIN book_tags ON book_tags.book_id = books.id").
			Where("book_tags.tag_id IN ?", tagIDs)
	}

	if len(categoryIDs) > 0 {
		query = query.Joins("JOIN book_categories ON book_categories.book_id = books.id").
			Where("book_categories.category_id IN ?", categoryIDs)
	}

	var total int64
	query.Count(&total)

	var books []models.Book
	offset := (page - 1) * pageSize
	query.Order(sortBy + " " + sortOrder).Offset(offset).Limit(pageSize).Find(&books)

	c.JSON(http.StatusOK, gin.H{
		"data":  books,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *BookHandler) GetBook(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var book models.Book
	result := database.DB.Preload("Tags").Preload("Categories").Preload("ReadingNotes").
		Preload("BorrowRecord").First(&book, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid request: %v", err)
		errors.HandleError(c, errors.ErrValidation)
		return
	}

	var existingBook models.Book
	if req.ISBN != "" {
		normalizedISBN := isbnutil.Normalize(req.ISBN)
		if !isbnutil.Validate(normalizedISBN) {
			errors.HandleError(c, errors.ErrISBNInvalid)
			return
		}
		result := database.DB.Where("isbn = ?", normalizedISBN).First(&existingBook)
		if result.Error == nil {
			errors.ErrorResponse(c, http.StatusConflict, "该ISBN已存在")
			return
		}
		req.ISBN = normalizedISBN
	}

	book := models.Book{
		Title:         req.Title,
		Author:        req.Author,
		Publisher:     req.Publisher,
		ISBN:          req.ISBN,
		Summary:       req.Summary,
		TotalPages:    req.TotalPages,
		ReadingStatus: models.StatusToRead,
	}

	if len(req.TagIDs) > 0 {
		var tags []*models.Tag
		database.DB.Find(&tags, req.TagIDs)
		book.Tags = tags
	}

	if len(req.CategoryIDs) > 0 {
		var categories []*models.Category
		database.DB.Find(&categories, req.CategoryIDs)
		book.Categories = categories
	}

	result := database.DB.Create(&book)
	if result.Error != nil {
		logger.Errorf("Create book failed: %v", result.Error)
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}

	database.DB.Preload("Tags").Preload("Categories").First(&book, book.ID)

	logger.Infof("Book created: %d - %s", book.ID, book.Title)
	c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var book models.Book
	result := database.DB.First(&book, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("Invalid request: %v", err)
		errors.HandleError(c, errors.ErrValidation)
		return
	}

	if req.ISBN != "" && req.ISBN != book.ISBN {
		normalizedISBN := isbnutil.Normalize(req.ISBN)
		if !isbnutil.Validate(normalizedISBN) {
			errors.HandleError(c, errors.ErrISBNInvalid)
			return
		}
		var existingBook models.Book
		result := database.DB.Where("isbn = ? AND id != ?", normalizedISBN, id).First(&existingBook)
		if result.Error == nil {
			errors.ErrorResponse(c, http.StatusConflict, "该ISBN已被其他图书使用")
			return
		}
		book.ISBN = normalizedISBN
	}

	if req.Title != "" {
		book.Title = req.Title
	}
	book.Author = req.Author
	book.Publisher = req.Publisher
	book.Summary = req.Summary
	if req.TotalPages > 0 {
		book.TotalPages = req.TotalPages
		book.ReadingProgress = utils.CalculateProgress(book.CurrentPage, book.TotalPages)
	}

	if req.TagIDs != nil {
		var tags []*models.Tag
		database.DB.Find(&tags, req.TagIDs)
		database.DB.Model(&book).Association("Tags").Replace(tags)
	}

	if req.CategoryIDs != nil {
		var categories []*models.Category
		database.DB.Find(&categories, req.CategoryIDs)
		database.DB.Model(&book).Association("Categories").Replace(categories)
	}

	database.DB.Save(&book)
	database.DB.Preload("Tags").Preload("Categories").First(&book, id)

	logger.Infof("Book updated: %d - %s", book.ID, book.Title)
	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var book models.Book
	result := database.DB.First(&book, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	tx := database.DB.Begin()

	tx.Model(&book).Association("Tags").Clear()
	tx.Model(&book).Association("Categories").Clear()
	tx.Where("book_id = ?", id).Delete(&models.ReadingNote{})
	tx.Where("book_id = ?", id).Delete(&models.ReadSession{})
	tx.Where("book_id = ?", id).Delete(&models.BorrowRecord{})

	if book.CoverImage != "" {
		uploadutil.DeleteImage(book.CoverImage)
	}

	tx.Delete(&book)
	tx.Commit()

	logger.Infof("Book deleted: %d", id)
	c.JSON(http.StatusNoContent, nil)
}

func (h *BookHandler) UploadCover(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var book models.Book
	result := database.DB.First(&book, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	file, err := c.FormFile("cover")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	accessURL, err := uploadutil.SaveImage(file)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if book.CoverImage != "" {
		uploadutil.DeleteImage(book.CoverImage)
	}

	book.CoverImage = accessURL
	database.DB.Save(&book)

	logger.Infof("Cover uploaded for book: %d", id)
	c.JSON(http.StatusOK, gin.H{
		"cover_image": accessURL,
	})
}

func (h *BookHandler) FetchByISBN(c *gin.Context) {
	isbn := c.Param("isbn")
	if isbn == "" {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	info, err := isbnutil.FetchBookInfo(isbn)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, info)
}

type UpdateReadingProgressRequest struct {
	CurrentPage int `json:"current_page" binding:"required,min=0"`
}

func (h *BookHandler) UpdateReadingProgress(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var req UpdateReadingProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrValidation)
		return
	}

	var book models.Book
	result := database.DB.First(&book, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	if req.CurrentPage > book.TotalPages && book.TotalPages > 0 {
		errors.HandleError(c, errors.ErrInvalidPageNumber)
		return
	}

	now := time.Now()
	wasReading := book.ReadingStatus == models.StatusReading

	book.CurrentPage = req.CurrentPage
	book.ReadingProgress = utils.CalculateProgress(req.CurrentPage, book.TotalPages)

	if req.CurrentPage == 0 {
		book.ReadingStatus = models.StatusToRead
		book.StartDate = nil
		book.EndDate = nil
	} else if req.CurrentPage >= book.TotalPages && book.TotalPages > 0 {
		book.ReadingStatus = models.StatusCompleted
		if book.StartDate == nil {
			book.StartDate = &now
		}
		if book.EndDate == nil {
			book.EndDate = &now
		}
	} else {
		book.ReadingStatus = models.StatusReading
		if book.StartDate == nil {
			book.StartDate = &now
		}
	}

	if !wasReading && book.ReadingStatus == models.StatusReading {
		logger.Infof("Started reading book: %d - %s", book.ID, book.Title)
	}

	database.DB.Save(&book)

	c.JSON(http.StatusOK, gin.H{
		"current_page":     book.CurrentPage,
		"reading_progress": book.ReadingProgress,
		"reading_status":   book.ReadingStatus,
	})
}

func (h *BookHandler) UpdateStatus(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	status := models.ReadingStatus(c.Query("status"))
	validStatuses := map[models.ReadingStatus]bool{
		models.StatusToRead:    true,
		models.StatusReading:   true,
		models.StatusCompleted: true,
		models.StatusAbandoned: true,
	}

	if !validStatuses[status] {
		errors.ErrorResponse(c, http.StatusBadRequest, "无效的阅读状态")
		return
	}

	var book models.Book
	result := database.DB.First(&book, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	now := time.Now()
	book.ReadingStatus = status

	if status == models.StatusReading && book.StartDate == nil {
		book.StartDate = &now
	}
	if status == models.StatusCompleted && book.EndDate == nil {
		book.EndDate = &now
		book.ReadingProgress = 100
		if book.TotalPages > 0 {
			book.CurrentPage = book.TotalPages
		}
	}

	database.DB.Save(&book)

	logger.Infof("Book %d status updated to: %s", id, status)
	c.JSON(http.StatusOK, gin.H{
		"reading_status": book.ReadingStatus,
	})
}

func (h *BookHandler) GetCurrentlyReading(c *gin.Context) {
	var books []models.Book
	database.DB.Where("reading_status = ?", models.StatusReading).
		Preload("Tags").Preload("Categories").
		Order("updated_at desc").
		Find(&books)

	c.JSON(http.StatusOK, books)
}
