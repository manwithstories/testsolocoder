package handlers

import (
	"net/http"

	"booklibrary/internal/database"
	"booklibrary/internal/errors"
	"booklibrary/internal/logger"
	"booklibrary/internal/models"
	"booklibrary/internal/utils"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct{}

func NewNoteHandler() *NoteHandler {
	return &NoteHandler{}
}

type CreateNoteRequest struct {
	BookID  uint64 `json:"book_id" binding:"required"`
	Page    int    `json:"page"`
	Content string `json:"content" binding:"required"`
}

type UpdateNoteRequest struct {
	Page    int    `json:"page"`
	Content string `json:"content"`
}

func (h *NoteHandler) GetNotesByBook(c *gin.Context) {
	bookID, err := utils.GetUintParam(c, "bookId")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var notes []models.ReadingNote
	result := database.DB.Where("book_id = ?", bookID).Order("page asc, created_at desc").Find(&notes)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, notes)
}

func (h *NoteHandler) GetNote(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var note models.ReadingNote
	result := database.DB.First(&note, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	var req CreateNoteRequest
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

	note := models.ReadingNote{
		BookID:  req.BookID,
		Page:    req.Page,
		Content: req.Content,
	}

	result = database.DB.Create(&note)
	if result.Error != nil {
		logger.Errorf("Create note failed: %v", result.Error)
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}

	logger.Infof("Note created: %d for book %d", note.ID, note.BookID)
	c.JSON(http.StatusCreated, note)
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	var note models.ReadingNote
	result := database.DB.First(&note, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	var req UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.ErrValidation)
		return
	}

	if req.Page >= 0 {
		note.Page = req.Page
	}
	if req.Content != "" {
		note.Content = req.Content
	}

	database.DB.Save(&note)
	logger.Infof("Note updated: %d", note.ID)
	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id, err := utils.GetUintParam(c, "id")
	if err != nil {
		errors.HandleError(c, errors.ErrBadRequest)
		return
	}

	result := database.DB.Delete(&models.ReadingNote{}, id)
	if result.Error != nil {
		errors.HandleError(c, errors.ErrInternalServer)
		return
	}
	if result.RowsAffected == 0 {
		errors.HandleError(c, errors.ErrNotFound)
		return
	}

	logger.Infof("Note deleted: %d", id)
	c.JSON(http.StatusNoContent, nil)
}
