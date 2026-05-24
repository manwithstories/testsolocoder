package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"campus-trade-platform/internal/models"
	"campus-trade-platform/internal/repository"
	jwtpkg "campus-trade-platform/pkg/jwt"
)

type UserService struct {
	userRepo *repository.UserRepository
	db       *gorm.DB
}

func NewUserService(userRepo *repository.UserRepository, db *gorm.DB) *UserService {
	return &UserService{userRepo: userRepo, db: db}
}

type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=50"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	Phone           string `json:"phone"`
	Role            string `json:"role" binding:"required,oneof=student merchant"`
	RealName        string `json:"real_name"`
	SchoolName      string `json:"school_name"`
	StudentID       string `json:"student_id"`
	StudentCardURL  string `json:"student_card_url"`
	BusinessLicense string `json:"business_license"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func (s *UserService) Register(req *RegisterRequest) (*models.User, error) {
	if s.userRepo.UsernameExists(req.Username) {
		return nil, errors.New("username already exists")
	}
	if s.userRepo.EmailExists(req.Email) {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		ID:              uuid.New(),
		Username:        req.Username,
		Email:           req.Email,
		Password:        string(hashedPassword),
		Phone:           req.Phone,
		Role:            models.UserRole(req.Role),
		Status:          models.UserStatusPending,
		RealName:        req.RealName,
		SchoolName:      req.SchoolName,
		StudentID:       req.StudentID,
		StudentCardURL:  req.StudentCardURL,
		BusinessLicense: req.BusinessLicense,
		Rating:          5.0,
		RatingCount:     0,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) Login(req *LoginRequest, secret string, expireHours int) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if user.Status == models.UserStatusBanned {
		return nil, errors.New("account is banned")
	}

	if user.Status == models.UserStatusRejected {
		return nil, errors.New("account registration was rejected")
	}

	token, err := jwtpkg.GenerateToken(user.ID, user.Username, string(user.Role), secret, expireHours)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) GetAllUsers(page, pageSize int, role, status string) ([]models.User, int64, error) {
	return s.userRepo.FindAll(page, pageSize, role, status)
}

func (s *UserService) UpdateUserStatus(id string, status models.UserStatus) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if status == models.UserStatusActive && user.Status != models.UserStatusPending {
		return errors.New("can only approve pending users")
	}

	return s.userRepo.UpdateStatus(id, status)
}

func (s *UserService) UpdateUserProfile(id string, updates map[string]interface{}) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if avatar, ok := updates["avatar"]; ok {
		user.Avatar = avatar.(string)
	}
	if phone, ok := updates["phone"]; ok {
		user.Phone = phone.(string)
	}
	if realName, ok := updates["real_name"]; ok {
		user.RealName = realName.(string)
	}

	return s.userRepo.Update(user)
}

func (s *UserService) ChangePassword(id, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

type UpdateRatingRequest struct {
	UserID    string
	NewRating int
}

func (s *UserService) UpdateUserRating(req *UpdateRatingRequest) error {
	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	newRating, newCount := func(currentRating float64, ratingCount int, newRating int) (float64, int) {
		total := currentRating * float64(ratingCount)
		total += float64(newRating)
		newCount := ratingCount + 1
		return total / float64(newCount), newCount
	}(user.Rating, user.RatingCount, req.NewRating)

	return s.userRepo.UpdateRating(req.UserID, newRating, newCount)
}

func (s *UserService) GetTopRatedUsers(limit int) ([]models.User, error) {
	var users []models.User
	err := s.db.Where("role IN ?", []string{"student", "merchant"}).
		Order("rating DESC").
		Limit(limit).
		Find(&users).Error
	return users, err
}

func (s *UserService) DeleteUser(id string) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.Delete(id)
}

type TextbookService struct {
	textbookRepo *repository.TextbookRepository
	categoryRepo *repository.CategoryRepository
	userRepo     *repository.UserRepository
	db           *gorm.DB
}

func NewTextbookService(textbookRepo *repository.TextbookRepository, categoryRepo *repository.CategoryRepository, userRepo *repository.UserRepository, db *gorm.DB) *TextbookService {
	return &TextbookService{
		textbookRepo: textbookRepo,
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
		db:           db,
	}
}

type CreateTextbookRequest struct {
	ISBN          string  `json:"isbn" binding:"required"`
	Title         string  `json:"title" binding:"required"`
	Author        string  `json:"author"`
	CourseName    string  `json:"course_name"`
	Edition       string  `json:"edition"`
	Publisher     string  `json:"publisher"`
	OriginalPrice float64 `json:"original_price"`
	Price         float64 `json:"price" binding:"required,gt=0"`
	Condition     string  `json:"condition" binding:"required"`
	Description   string  `json:"description"`
	CoverImage    string  `json:"cover_image"`
	SellerID      string  `json:"seller_id" binding:"required"`
	CategoryID    string  `json:"category_id"`
}

func (s *TextbookService) CreateTextbook(req *CreateTextbookRequest) (*models.Textbook, error) {
	sellerID, err := uuid.Parse(req.SellerID)
	if err != nil {
		return nil, errors.New("invalid seller ID")
	}

	_, err = s.userRepo.FindByID(req.SellerID)
	if err != nil {
		return nil, errors.New("seller not found")
	}

	var categoryID *uuid.UUID
	if req.CategoryID != "" {
		cid, err := uuid.Parse(req.CategoryID)
		if err == nil {
			categoryID = &cid
		}
	}

	textbook := &models.Textbook{
		ID:            uuid.New(),
		ISBN:          req.ISBN,
		Title:         req.Title,
		Author:        req.Author,
		CourseName:    req.CourseName,
		Edition:       req.Edition,
		Publisher:     req.Publisher,
		OriginalPrice: req.OriginalPrice,
		Price:         req.Price,
		Condition:     models.TextbookCondition(req.Condition),
		Description:   req.Description,
		CoverImage:    req.CoverImage,
		SellerID:      sellerID,
		CategoryID:   categoryID,
		Status:        models.TextbookStatusAvailable,
		ViewCount:     0,
	}

	err = s.textbookRepo.Create(textbook)
	if err != nil {
		return nil, fmt.Errorf("failed to create textbook: %w", err)
	}

	return textbook, nil
}

func (s *TextbookService) GetTextbookByID(id string) (*models.Textbook, error) {
	textbook, err := s.textbookRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("textbook not found")
	}

	go s.textbookRepo.IncrementViewCount(id)

	return textbook, nil
}

func (s *TextbookService) GetAllTextbooks(page, pageSize int, keyword, categoryID, status string) ([]models.Textbook, int64, error) {
	return s.textbookRepo.FindAll(page, pageSize, keyword, categoryID, status)
}

func (s *TextbookService) SearchByISBN(isbn string) (*models.Textbook, error) {
	textbook, err := s.textbookRepo.FindByISBN(isbn)
	if err != nil {
		return nil, errors.New("textbook not found with this ISBN")
	}
	return textbook, nil
}

func (s *TextbookService) GetSellerTextbooks(sellerID string, page, pageSize int) ([]models.Textbook, int64, error) {
	return s.textbookRepo.FindBySellerID(sellerID, page, pageSize)
}

func (s *TextbookService) UpdateTextbook(id string, updates map[string]interface{}) error {
	textbook, err := s.textbookRepo.FindByID(id)
	if err != nil {
		return errors.New("textbook not found")
	}

	if title, ok := updates["title"]; ok {
		textbook.Title = title.(string)
	}
	if price, ok := updates["price"]; ok {
		textbook.Price = price.(float64)
	}
	if description, ok := updates["description"]; ok {
		textbook.Description = description.(string)
	}
	if condition, ok := updates["condition"]; ok {
		textbook.Condition = models.TextbookCondition(condition.(string))
	}
	if coverImage, ok := updates["cover_image"]; ok {
		textbook.CoverImage = coverImage.(string)
	}

	textbook.UpdatedAt = time.Now()
	return s.textbookRepo.Update(textbook)
}

func (s *TextbookService) DeleteTextbook(id string) error {
	_, err := s.textbookRepo.FindByID(id)
	if err != nil {
		return errors.New("textbook not found")
	}
	return s.textbookRepo.Delete(id)
}

func (s *TextbookService) UpdateTextbookStatus(id string, status models.TextbookStatus) error {
	return s.textbookRepo.UpdateStatus(id, status)
}

func (s *TextbookService) GetPopularTextbooks(limit int) ([]models.Textbook, error) {
	return s.textbookRepo.GetPopular(limit)
}

type NoteService struct {
	noteRepo *repository.NoteRepository
	userRepo *repository.UserRepository
	db       *gorm.DB
}

func NewNoteService(noteRepo *repository.NoteRepository, userRepo *repository.UserRepository, db *gorm.DB) *NoteService {
	return &NoteService{
		noteRepo: noteRepo,
		userRepo: userRepo,
		db:       db,
	}
}

type CreateNoteRequest struct {
	Title       string `json:"title" binding:"required"`
	Subject     string `json:"subject"`
	CourseName  string `json:"course_name"`
	Description string `json:"description"`
	FileURL     string `json:"file_url" binding:"required"`
	FileType    string `json:"file_type"`
	FileSize    int64  `json:"file_size"`
	CoverImage  string `json:"cover_image"`
	UploaderID  string `json:"uploader_id" binding:"required"`
	CategoryID  string `json:"category_id"`
}

func (s *NoteService) CreateNote(req *CreateNoteRequest) (*models.Note, error) {
	uploaderID, err := uuid.Parse(req.UploaderID)
	if err != nil {
		return nil, errors.New("invalid uploader ID")
	}

	_, err = s.userRepo.FindByID(req.UploaderID)
	if err != nil {
		return nil, errors.New("uploader not found")
	}

	var categoryID *uuid.UUID
	if req.CategoryID != "" {
		cid, err := uuid.Parse(req.CategoryID)
		if err == nil {
			categoryID = &cid
		}
	}

	note := &models.Note{
		ID:            uuid.New(),
		Title:         req.Title,
		Subject:       req.Subject,
		CourseName:    req.CourseName,
		Description:   req.Description,
		FileURL:       req.FileURL,
		FileType:      req.FileType,
		FileSize:      req.FileSize,
		CoverImage:    req.CoverImage,
		UploaderID:   uploaderID,
		CategoryID:   categoryID,
		ViewCount:     0,
		DownloadCount: 0,
		Rating:        0,
		RatingCount:   0,
		IsFeatured:    false,
	}

	err = s.noteRepo.Create(note)
	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	return note, nil
}

func (s *NoteService) GetNoteByID(id string) (*models.Note, error) {
	note, err := s.noteRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("note not found")
	}

	go s.noteRepo.IncrementViewCount(id)

	return note, nil
}

func (s *NoteService) GetAllNotes(page, pageSize int, keyword, subject, categoryID string, isFeatured bool) ([]models.Note, int64, error) {
	return s.noteRepo.FindAll(page, pageSize, keyword, subject, categoryID, isFeatured)
}

func (s *NoteService) GetUploaderNotes(uploaderID string, page, pageSize int) ([]models.Note, int64, error) {
	return s.noteRepo.FindByUploaderID(uploaderID, page, pageSize)
}

func (s *NoteService) UpdateNote(id string, updates map[string]interface{}) error {
	note, err := s.noteRepo.FindByID(id)
	if err != nil {
		return errors.New("note not found")
	}

	if title, ok := updates["title"]; ok {
		note.Title = title.(string)
	}
	if description, ok := updates["description"]; ok {
		note.Description = description.(string)
	}
	if coverImage, ok := updates["cover_image"]; ok {
		note.CoverImage = coverImage.(string)
	}

	note.UpdatedAt = time.Now()
	return s.noteRepo.Update(note)
}

func (s *NoteService) DeleteNote(id string) error {
	_, err := s.noteRepo.FindByID(id)
	if err != nil {
		return errors.New("note not found")
	}
	return s.noteRepo.Delete(id)
}

func (s *NoteService) IncrementDownload(id string) error {
	return s.noteRepo.IncrementDownloadCount(id)
}

func (s *NoteService) SetNoteFeatured(id string, isFeatured bool) error {
	return s.noteRepo.SetFeatured(id, isFeatured)
}

func (s *NoteService) GetFeaturedNotes(limit int) ([]models.Note, error) {
	return s.noteRepo.GetFeatured(limit)
}

func (s *NoteService) UpdateNoteRating(noteID string, newRating int) error {
	note, err := s.noteRepo.FindByID(noteID)
	if err != nil {
		return errors.New("note not found")
	}

	updatedRating, updatedCount := func(currentRating float64, ratingCount int, newRating int) (float64, int) {
		total := currentRating * float64(ratingCount)
		total += float64(newRating)
		newCount := ratingCount + 1
		return total / float64(newCount), newCount
	}(note.Rating, note.RatingCount, newRating)

	return s.noteRepo.UpdateRating(noteID, updatedRating, updatedCount)
}
