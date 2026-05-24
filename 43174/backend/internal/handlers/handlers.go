package handlers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"campus-trade-platform/internal/middleware"
	"campus-trade-platform/internal/models"
	"campus-trade-platform/internal/services"
	"campus-trade-platform/internal/utils"
	"campus-trade-platform/pkg/upload"
)

type UserHandler struct {
	userService  *services.UserService
	jwtSecret    string
	expireHours  int
	uploadPath   string
	maxUploadSize int64
}

func NewUserHandler(userService *services.UserService, jwtSecret string, expireHours int, uploadPath string, maxUploadSize int64) *UserHandler {
	return &UserHandler{
		userService:   userService,
		jwtSecret:     jwtSecret,
		expireHours:   expireHours,
		uploadPath:    uploadPath,
		maxUploadSize: maxUploadSize,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	response, err := h.userService.Login(&req, h.jwtSecret, h.expireHours)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.Success(c, response)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userCtx, exists := c.Get(string(middleware.UserContextKey))
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	user := userCtx.(middleware.UserContext)

	profile, err := h.userService.GetUserByID(user.UserID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, profile)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	role := c.Query("role")
	status := c.Query("status")

	users, total, err := h.userService.GetAllUsers(page, pageSize, role, status)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, users, total, page, pageSize)
}

func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status models.UserStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.userService.UpdateUserStatus(id, req.Status)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "User status updated successfully"})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userCtx, exists := c.Get(string(middleware.UserContextKey))
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	user := userCtx.(middleware.UserContext)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.userService.UpdateUserProfile(user.UserID, updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Profile updated successfully"})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userCtx, exists := c.Get(string(middleware.UserContextKey))
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	user := userCtx.(middleware.UserContext)

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.userService.ChangePassword(user.UserID, req.OldPassword, req.NewPassword)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Password changed successfully"})
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userCtx, exists := c.Get(string(middleware.UserContextKey))
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	user := userCtx.(middleware.UserContext)

	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "File is required")
		return
	}

	os.MkdirAll(h.uploadPath, 0755)

	filePath, err := upload.UploadFile(file, h.uploadPath, upload.AllowedImageTypes, h.maxUploadSize)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err = h.userService.UpdateUserProfile(user.UserID, map[string]interface{}{
		"avatar": filePath,
	})
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"avatar_url": filePath})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.DeleteUser(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) GetTopRatedUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, err := h.userService.GetTopRatedUsers(limit)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, users)
}

type TextbookHandler struct {
	textbookService *services.TextbookService
	uploadPath      string
	maxUploadSize   int64
}

func NewTextbookHandler(textbookService *services.TextbookService, uploadPath string, maxUploadSize int64) *TextbookHandler {
	return &TextbookHandler{
		textbookService: textbookService,
		uploadPath:      uploadPath,
		maxUploadSize:   maxUploadSize,
	}
}

func (h *TextbookHandler) CreateTextbook(c *gin.Context) {
	var req services.CreateTextbookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userCtx, exists := c.Get(string(middleware.UserContextKey))
	if exists {
		user := userCtx.(middleware.UserContext)
		req.SellerID = user.UserID
	}

	textbook, err := h.textbookService.CreateTextbook(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, textbook)
}

func (h *TextbookHandler) GetTextbookByID(c *gin.Context) {
	id := c.Param("id")

	textbook, err := h.textbookService.GetTextbookByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, textbook)
}

func (h *TextbookHandler) GetAllTextbooks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	categoryID := c.Query("category_id")
	status := c.Query("status")

	textbooks, total, err := h.textbookService.GetAllTextbooks(page, pageSize, keyword, categoryID, status)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, textbooks, total, page, pageSize)
}

func (h *TextbookHandler) SearchByISBN(c *gin.Context) {
	isbn := c.Query("isbn")
	if isbn == "" {
		utils.BadRequest(c, "ISBN is required")
		return
	}

	if !utils.IsValidISBN(isbn) {
		utils.BadRequest(c, "Invalid ISBN format")
		return
	}

	textbook, err := h.textbookService.SearchByISBN(isbn)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, textbook)
}

func (h *TextbookHandler) GetSellerTextbooks(c *gin.Context) {
	sellerID := c.Param("seller_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	textbooks, total, err := h.textbookService.GetSellerTextbooks(sellerID, page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, textbooks, total, page, pageSize)
}

func (h *TextbookHandler) UpdateTextbook(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.textbookService.UpdateTextbook(id, updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Textbook updated successfully"})
}

func (h *TextbookHandler) DeleteTextbook(c *gin.Context) {
	id := c.Param("id")

	err := h.textbookService.DeleteTextbook(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Textbook deleted successfully"})
}

func (h *TextbookHandler) UpdateTextbookStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status models.TextbookStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.textbookService.UpdateTextbookStatus(id, req.Status)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Textbook status updated successfully"})
}

func (h *TextbookHandler) UploadCoverImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "File is required")
		return
	}

	os.MkdirAll(h.uploadPath, 0755)

	filePath, err := upload.UploadFile(file, h.uploadPath, upload.AllowedImageTypes, h.maxUploadSize)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"image_url": filePath})
}

func (h *TextbookHandler) GetPopularTextbooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	textbooks, err := h.textbookService.GetPopularTextbooks(limit)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, textbooks)
}

type NoteHandler struct {
	noteService   *services.NoteService
	uploadPath    string
	maxUploadSize int64
}

func NewNoteHandler(noteService *services.NoteService, uploadPath string, maxUploadSize int64) *NoteHandler {
	return &NoteHandler{
		noteService:   noteService,
		uploadPath:    uploadPath,
		maxUploadSize: maxUploadSize,
	}
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	var req services.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userCtx, exists := c.Get(string(middleware.UserContextKey))
	if exists {
		user := userCtx.(middleware.UserContext)
		req.UploaderID = user.UserID
	}

	note, err := h.noteService.CreateNote(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, note)
}

func (h *NoteHandler) GetNoteByID(c *gin.Context) {
	id := c.Param("id")

	note, err := h.noteService.GetNoteByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, note)
}

func (h *NoteHandler) GetAllNotes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	subject := c.Query("subject")
	categoryID := c.Query("category_id")
	isFeatured := c.Query("is_featured") == "true"

	notes, total, err := h.noteService.GetAllNotes(page, pageSize, keyword, subject, categoryID, isFeatured)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, notes, total, page, pageSize)
}

func (h *NoteHandler) GetUploaderNotes(c *gin.Context) {
	uploaderID := c.Param("uploader_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	notes, total, err := h.noteService.GetUploaderNotes(uploaderID, page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, notes, total, page, pageSize)
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.noteService.UpdateNote(id, updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Note updated successfully"})
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id := c.Param("id")

	err := h.noteService.DeleteNote(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Note deleted successfully"})
}

func (h *NoteHandler) UploadNoteFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "File is required")
		return
	}

	os.MkdirAll(h.uploadPath, 0755)

	filePath, err := upload.UploadFile(file, h.uploadPath, upload.AllowedFileTypes, h.maxUploadSize*5)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"file_url": filePath})
}

func (h *NoteHandler) IncrementDownload(c *gin.Context) {
	id := c.Param("id")

	err := h.noteService.IncrementDownload(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Download count updated"})
}

func (h *NoteHandler) SetFeatured(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		IsFeatured bool `json:"is_featured"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.noteService.SetNoteFeatured(id, req.IsFeatured)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Note featured status updated"})
}

func (h *NoteHandler) GetFeaturedNotes(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	notes, err := h.noteService.GetFeaturedNotes(limit)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, notes)
}

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req services.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	transaction, err := h.transactionService.CreateTransaction(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, transaction)
}

func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	id := c.Param("id")

	transaction, err := h.transactionService.GetTransactionByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, transaction)
}

func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	textbookID := c.Query("textbook_id")

	transactions, total, err := h.transactionService.GetAllTransactions(page, pageSize, status, textbookID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, transactions, total, page, pageSize)
}

func (h *TransactionHandler) GetBuyerTransactions(c *gin.Context) {
	buyerID := c.Param("buyer_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	transactions, total, err := h.transactionService.GetBuyerTransactions(buyerID, page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, transactions, total, page, pageSize)
}

func (h *TransactionHandler) GetSellerTransactions(c *gin.Context) {
	sellerID := c.Param("seller_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	transactions, total, err := h.transactionService.GetSellerTransactions(sellerID, page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, transactions, total, page, pageSize)
}

func (h *TransactionHandler) ConfirmTransaction(c *gin.Context) {
	id := c.Param("id")

	err := h.transactionService.ConfirmTransaction(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Transaction confirmed successfully"})
}

func (h *TransactionHandler) CompleteTransaction(c *gin.Context) {
	id := c.Param("id")

	err := h.transactionService.CompleteTransaction(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Transaction completed successfully"})
}

func (h *TransactionHandler) CancelTransaction(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.transactionService.CancelTransaction(id, req.Reason)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Transaction cancelled successfully"})
}

func (h *TransactionHandler) StartNegotiation(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Price float64 `json:"price" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.transactionService.StartNegotiation(id, req.Price)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Negotiation started"})
}

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req services.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	order, err := h.orderService.CreateOrder(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, order)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")

	order, err := h.orderService.GetOrderByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, order)
}

func (h *OrderHandler) GetOrderByOrderNo(c *gin.Context) {
	orderNo := c.Param("order_no")

	order, err := h.orderService.GetOrderByOrderNo(orderNo)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, order)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	buyerID := c.Query("buyer_id")
	sellerID := c.Query("seller_id")

	orders, total, err := h.orderService.GetAllOrders(page, pageSize, status, buyerID, sellerID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, orders, total, page, pageSize)
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userCtx, exists := c.Get(string(middleware.UserContextKey))
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	user := userCtx.(middleware.UserContext)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orders, total, err := h.orderService.GetUserOrders(user.UserID, page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, orders, total, page, pageSize)
}

func (h *OrderHandler) PayOrder(c *gin.Context) {
	id := c.Param("id")

	err := h.orderService.PayOrder(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Order paid successfully"})
}

func (h *OrderHandler) ShipOrder(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		TrackingNumber string `json:"tracking_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.orderService.ShipOrder(id, req.TrackingNumber)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Order shipped successfully"})
}

func (h *OrderHandler) DeliverOrder(c *gin.Context) {
	id := c.Param("id")

	err := h.orderService.DeliverOrder(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Order delivered successfully"})
}

func (h *OrderHandler) CompleteOrder(c *gin.Context) {
	id := c.Param("id")

	err := h.orderService.CompleteOrder(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Order completed successfully"})
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.orderService.CancelOrder(id, req.Reason)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Order cancelled successfully"})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status models.OrderStatus `json:"status" binding:"required"`
		Remark string            `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.orderService.UpdateOrderStatus(id, req.Status, req.Remark)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Order status updated successfully"})
}

type MessageHandler struct {
	messageService *services.MessageService
}

func NewMessageHandler(messageService *services.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var req services.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	message, err := h.messageService.CreateMessage(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, message)
}

func (h *MessageHandler) GetConversation(c *gin.Context) {
	userID1 := c.Query("user_id_1")
	userID2 := c.Query("user_id_2")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if userID1 == "" || userID2 == "" {
		utils.BadRequest(c, "Both user_id_1 and user_id_2 are required")
		return
	}

	messages, total, err := h.messageService.GetConversation(userID1, userID2, page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, messages, total, page, pageSize)
}

func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	userCtx, exists := c.Get(string(middleware.UserContextKey))
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	user := userCtx.(middleware.UserContext)

	count, err := h.messageService.GetUnreadCount(user.UserID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"unread_count": count})
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userCtx, exists := c.Get(string(middleware.UserContextKey))
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	user := userCtx.(middleware.UserContext)

	err := h.messageService.MarkAsRead(user.UserID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Messages marked as read"})
}

type ReviewHandler struct {
	reviewService *services.ReviewService
}

func NewReviewHandler(reviewService *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req services.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	review, err := h.reviewService.CreateReview(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, review)
}

func (h *ReviewHandler) GetTextbookReviews(c *gin.Context) {
	textbookID := c.Param("textbook_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	reviews, total, err := h.reviewService.GetTextbookReviews(textbookID, page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, reviews, total, page, pageSize)
}

func (h *ReviewHandler) GetNoteReviews(c *gin.Context) {
	noteID := c.Param("note_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	reviews, total, err := h.reviewService.GetNoteReviews(noteID, page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, reviews, total, page, pageSize)
}

func (h *ReviewHandler) GetAllReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	isMalicious := c.Query("is_malicious") == "true"

	reviews, total, err := h.reviewService.GetAllReviews(page, pageSize, isMalicious)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, reviews, total, page, pageSize)
}

func (h *ReviewHandler) HideReview(c *gin.Context) {
	id := c.Param("id")

	err := h.reviewService.HideReview(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Review hidden successfully"})
}

func (h *ReviewHandler) MarkMalicious(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		IsMalicious bool `json:"is_malicious"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.reviewService.MarkMalicious(id, req.IsMalicious)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Review status updated"})
}

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req services.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	category, err := h.categoryService.CreateCategory(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, category)
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	category, err := h.categoryService.GetCategoryByID(id)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, category)
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, categories)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	err := h.categoryService.UpdateCategory(id, updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Category updated successfully"})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	err := h.categoryService.DeleteCategory(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Category deleted successfully"})
}

type StatisticsHandler struct {
	statisticsService *services.StatisticsService
}

func NewStatisticsHandler(statisticsService *services.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
	}
}

func (h *StatisticsHandler) GetTextbookStats(c *gin.Context) {
	stats, err := h.statisticsService.GetTextbookStats()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

func (h *StatisticsHandler) GetUserStats(c *gin.Context) {
	stats, err := h.statisticsService.GetUserStats()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

func (h *StatisticsHandler) GetOrderStats(c *gin.Context) {
	stats, err := h.statisticsService.GetOrderStats()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

func (h *StatisticsHandler) GetPopularTextbooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	textbooks, err := h.statisticsService.GetPopularTextbooks(limit)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, textbooks)
}

func (h *StatisticsHandler) GetTopUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, err := h.statisticsService.GetTopUsers(limit)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, users)
}

func (h *StatisticsHandler) GetMonthlyStats(c *gin.Context) {
	months, _ := strconv.Atoi(c.DefaultQuery("months", "6"))

	stats, err := h.statisticsService.GetMonthlyStats(months)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

func (h *StatisticsHandler) ExportMonthlyReport(c *gin.Context) {
	month := c.Query("month")
	if month == "" {
		utils.BadRequest(c, "Month parameter is required (format: YYYY-MM)")
		return
	}

	os.MkdirAll("./exports", 0755)

	filePath, err := h.statisticsService.ExportMonthlyReport(month)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	c.FileAttachment(filePath, "monthly_report_"+month+".xlsx")
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "campus-trade-platform",
		"version": "1.0.0",
	})
}

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	userCtx, exists := c.Get(string(middleware.UserContextKey))
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	user := userCtx.(middleware.UserContext)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	notifications, total, err := h.notificationService.GetUserNotifications(user.UserID, page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Paginated(c, notifications, total, page, pageSize)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	err := h.notificationService.MarkAsRead(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "Notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userCtx, exists := c.Get(string(middleware.UserContextKey))
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	user := userCtx.(middleware.UserContext)

	err := h.notificationService.MarkAllAsRead(user.UserID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "All notifications marked as read"})
}
