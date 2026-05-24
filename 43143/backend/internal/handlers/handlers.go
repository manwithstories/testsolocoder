package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"skillshare/internal/models"
	"skillshare/internal/service"
	"skillshare/pkg/response"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	ip := c.ClientIP()
	user, token, err := h.userService.Register(&input, ip)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	ip := c.ClientIP()
	user, token, err := h.userService.Login(&input, ip)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	token, err := h.userService.RefreshToken(input.RefreshToken)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, gin.H{"token": token})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.userService.GetUser(userID.(uuid.UUID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	user, err := h.userService.UpdateProfile(userID.(uuid.UUID), input)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的用户ID", nil)
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	users, total, err := h.userService.ListUsers(page, pageSize, keyword)
	if err != nil {
		response.InternalServerError(c, "获取用户列表失败", err.Error())
		return
	}

	response.Paginated(c, users, total, page, pageSize)
}

func (h *UserHandler) AddSkillTags(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input struct {
		TagIDs []uuid.UUID `json:"tag_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := h.userService.AddSkillTags(userID.(uuid.UUID), input.TagIDs); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) RemoveSkillTag(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tagID, err := uuid.Parse(c.Param("tag_id"))
	if err != nil {
		response.BadRequest(c, "无效的标签ID", nil)
		return
	}

	if err := h.userService.RemoveSkillTag(userID.(uuid.UUID), tagID); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

type SkillHandler struct {
	skillService *service.SkillService
}

func NewSkillHandler(skillService *service.SkillService) *SkillHandler {
	return &SkillHandler{skillService: skillService}
}

func (h *SkillHandler) CreateCategory(c *gin.Context) {
	var input struct {
		Name      string `json:"name" binding:"required"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	category, err := h.skillService.CreateCategory(input.Name, input.Icon, input.SortOrder)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, category)
}

func (h *SkillHandler) GetCategories(c *gin.Context) {
	categories, err := h.skillService.GetCategories()
	if err != nil {
		response.InternalServerError(c, "获取分类失败", err.Error())
		return
	}

	response.Success(c, categories)
}

func (h *SkillHandler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的分类ID", nil)
		return
	}

	var input struct {
		Name      string `json:"name" binding:"required"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	category, err := h.skillService.UpdateCategory(id, input.Name, input.Icon, input.SortOrder)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, category)
}

func (h *SkillHandler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的分类ID", nil)
		return
	}

	if err := h.skillService.DeleteCategory(id); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

func (h *SkillHandler) CreateTag(c *gin.Context) {
	var input struct {
		Name       string     `json:"name" binding:"required"`
		CategoryID *uuid.UUID `json:"category_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	tag, err := h.skillService.CreateTag(input.Name, input.CategoryID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, tag)
}

func (h *SkillHandler) GetTags(c *gin.Context) {
	var categoryID *uuid.UUID
	if c.Query("category_id") != "" {
		id, err := uuid.Parse(c.Query("category_id"))
		if err == nil {
			categoryID = &id
		}
	}

	tags, err := h.skillService.GetTags(categoryID)
	if err != nil {
		response.InternalServerError(c, "获取标签失败", err.Error())
		return
	}

	response.Success(c, tags)
}

func (h *SkillHandler) UpdateTag(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的标签ID", nil)
		return
	}

	var input struct {
		Name       string     `json:"name" binding:"required"`
		CategoryID *uuid.UUID `json:"category_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	tag, err := h.skillService.UpdateTag(id, input.Name, input.CategoryID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, tag)
}

func (h *SkillHandler) DeleteTag(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的标签ID", nil)
		return
	}

	if err := h.skillService.DeleteTag(id); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

func (h *SkillHandler) CreateSkill(c *gin.Context) {
	var input service.CreateSkillInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	skill, err := h.skillService.CreateSkill(&input)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, skill)
}

func (h *SkillHandler) GetSkills(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var categoryID *uuid.UUID
	if c.Query("category_id") != "" {
		id, err := uuid.Parse(c.Query("category_id"))
		if err == nil {
			categoryID = &id
		}
	}

	skills, total, err := h.skillService.GetSkills(page, pageSize, categoryID, keyword)
	if err != nil {
		response.InternalServerError(c, "获取技能列表失败", err.Error())
		return
	}

	response.Paginated(c, skills, total, page, pageSize)
}

func (h *SkillHandler) GetSkill(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的技能ID", nil)
		return
	}

	skill, err := h.skillService.GetSkill(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, skill)
}

func (h *SkillHandler) UpdateSkill(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的技能ID", nil)
		return
	}

	var input service.CreateSkillInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	skill, err := h.skillService.UpdateSkill(id, &input)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, skill)
}

func (h *SkillHandler) DeleteSkill(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的技能ID", nil)
		return
	}

	if err := h.skillService.DeleteSkill(id); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

func (h *SkillHandler) GetPopularSkills(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	skills, err := h.skillService.GetPopularSkills(limit)
	if err != nil {
		response.InternalServerError(c, "获取热门技能失败", err.Error())
		return
	}

	response.Success(c, skills)
}

func (h *SkillHandler) CreatePosting(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input service.CreatePostingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	posting, err := h.skillService.CreatePosting(userID.(uuid.UUID), &input)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, posting)
}

func (h *SkillHandler) GetPostings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var skillID *uuid.UUID
	if c.Query("skill_id") != "" {
		id, err := uuid.Parse(c.Query("skill_id"))
		if err == nil {
			skillID = &id
		}
	}

	var teacherID *uuid.UUID
	if c.Query("teacher_id") != "" {
		id, err := uuid.Parse(c.Query("teacher_id"))
		if err == nil {
			teacherID = &id
		}
	}

	postings, total, err := h.skillService.GetPostings(page, pageSize, skillID, teacherID)
	if err != nil {
		response.InternalServerError(c, "获取课程列表失败", err.Error())
		return
	}

	response.Paginated(c, postings, total, page, pageSize)
}

func (h *SkillHandler) GetPosting(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的课程ID", nil)
		return
	}

	posting, err := h.skillService.GetPosting(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, posting)
}

func (h *SkillHandler) UpdatePosting(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的课程ID", nil)
		return
	}

	var input service.CreatePostingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	posting, err := h.skillService.UpdatePosting(id, &input)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, posting)
}

func (h *SkillHandler) DeletePosting(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的课程ID", nil)
		return
	}

	if err := h.skillService.DeletePosting(id); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

func (h *SkillHandler) MatchSkills(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filter := &service.MatchFilter{}
	if c.Query("skill_id") != "" {
		id, err := uuid.Parse(c.Query("skill_id"))
		if err == nil {
			filter.SkillID = &id
		}
	}
	if minRating, err := strconv.ParseFloat(c.Query("min_rating"), 64); err == nil {
		filter.MinRating = minRating
	}
	if method := c.Query("method"); method != "" {
		filter.Method = &method
	}
	if maxPrice, err := strconv.ParseFloat(c.Query("max_price"), 64); err == nil {
		filter.MaxPrice = &maxPrice
	}

	postings, total, err := h.skillService.MatchSkills(userID.(uuid.UUID), filter, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "匹配失败", err.Error())
		return
	}

	response.Paginated(c, postings, total, page, pageSize)
}

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input service.CreateBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	booking, err := h.bookingService.CreateBooking(userID.(uuid.UUID), &input)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, booking)
}

func (h *BookingHandler) GetBooking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的预约ID", nil)
		return
	}

	booking, err := h.bookingService.GetBooking(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, booking)
}

func (h *BookingHandler) ListBookings(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role := c.Query("role")
	if role == "" {
		role = "student"
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	status := c.Query("status")

	bookings, total, err := h.bookingService.ListBookings(userID.(uuid.UUID), role, page, pageSize, &status)
	if err != nil {
		response.InternalServerError(c, "获取预约列表失败", err.Error())
		return
	}

	response.Paginated(c, bookings, total, page, pageSize)
}

func (h *BookingHandler) ConfirmBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的预约ID", nil)
		return
	}

	if err := h.bookingService.ConfirmBooking(id, userID.(uuid.UUID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

func (h *BookingHandler) RejectBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的预约ID", nil)
		return
	}

	var input struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&input)

	if err := h.bookingService.RejectBooking(id, userID.(uuid.UUID), input.Reason); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的预约ID", nil)
		return
	}

	var input struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&input)

	if err := h.bookingService.CancelBooking(id, userID.(uuid.UUID), input.Reason); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

func (h *BookingHandler) CompleteBooking(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的预约ID", nil)
		return
	}

	if err := h.bookingService.CompleteBooking(id, userID.(uuid.UUID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

func (h *BookingHandler) CreateReview(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input service.CreateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	review, err := h.bookingService.CreateReview(userID.(uuid.UUID), &input)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, review)
}

func (h *BookingHandler) GetReviewsByPosting(c *gin.Context) {
	postingID, err := uuid.Parse(c.Param("posting_id"))
	if err != nil {
		response.BadRequest(c, "无效的课程ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	reviews, total, err := h.bookingService.GetReviewsByPosting(postingID, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "获取评价失败", err.Error())
		return
	}

	response.Paginated(c, reviews, total, page, pageSize)
}

func (h *BookingHandler) GetReviewsByUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.BadRequest(c, "无效的用户ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	reviews, total, err := h.bookingService.GetReviewsByUser(userID, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "获取评价失败", err.Error())
		return
	}

	response.Paginated(c, reviews, total, page, pageSize)
}

func (h *BookingHandler) CreateComplaint(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input service.CreateComplaintInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	complaint, err := h.bookingService.CreateComplaint(userID.(uuid.UUID), &input)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, complaint)
}

func (h *BookingHandler) GetComplaints(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	status := c.Query("status")

	complaints, total, err := h.bookingService.GetComplaints(page, pageSize, &status)
	if err != nil {
		response.InternalServerError(c, "获取投诉列表失败", err.Error())
		return
	}

	response.Paginated(c, complaints, total, page, pageSize)
}

func (h *BookingHandler) HandleComplaint(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的投诉ID", nil)
		return
	}

	var input struct {
		Result string `json:"result" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := h.bookingService.HandleComplaint(id, userID.(uuid.UUID), input.Result); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input service.SendMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	message, err := h.messageService.SendMessage(userID.(uuid.UUID), &input)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, message)
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID, _ := c.Get("user_id")
	receiverID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.BadRequest(c, "无效的用户ID", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	messages, total, err := h.messageService.GetMessages(userID.(uuid.UUID), receiverID, page, pageSize)
	if err != nil {
		response.InternalServerError(c, "获取消息失败", err.Error())
		return
	}

	response.Paginated(c, messages, total, page, pageSize)
}

func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID, _ := c.Get("user_id")

	conversations, err := h.messageService.GetConversations(userID.(uuid.UUID))
	if err != nil {
		response.InternalServerError(c, "获取会话列表失败", err.Error())
		return
	}

	response.Success(c, conversations)
}

func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	count, err := h.messageService.GetUnreadCount(userID.(uuid.UUID))
	if err != nil {
		response.InternalServerError(c, "获取未读消息数失败", err.Error())
		return
	}

	response.Success(c, gin.H{"unread_count": count})
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	senderID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.BadRequest(c, "无效的用户ID", nil)
		return
	}

	if err := h.messageService.MarkAsRead(senderID, userID.(uuid.UUID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input service.CreatePaymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	payment, err := h.paymentService.CreatePayment(userID.(uuid.UUID), &input)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, payment)
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的支付ID", nil)
		return
	}

	payment, err := h.paymentService.GetPayment(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, payment)
}

func (h *PaymentHandler) GetUserPayments(c *gin.Context) {
	userID, _ := c.Get("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	payments, total, err := h.paymentService.GetUserPayments(userID.(uuid.UUID), page, pageSize)
	if err != nil {
		response.InternalServerError(c, "获取支付记录失败", err.Error())
		return
	}

	response.Paginated(c, payments, total, page, pageSize)
}

func (h *PaymentHandler) GetWallet(c *gin.Context) {
	userID, _ := c.Get("user_id")

	wallet, err := h.paymentService.GetWallet(userID.(uuid.UUID))
	if err != nil {
		response.InternalServerError(c, "获取钱包失败", err.Error())
		return
	}

	response.Success(c, wallet)
}

func (h *PaymentHandler) Withdraw(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	if err := h.paymentService.Withdraw(userID.(uuid.UUID), input.Amount); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

type ScheduleHandler struct {
	scheduleService *service.ScheduleService
}

func NewScheduleHandler(scheduleService *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{scheduleService: scheduleService}
}

func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input service.CreateScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	schedule, err := h.scheduleService.CreateSchedule(userID.(uuid.UUID), &input)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, schedule)
}

func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的日程ID", nil)
		return
	}

	var input service.CreateScheduleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误", err.Error())
		return
	}

	schedule, err := h.scheduleService.UpdateSchedule(userID.(uuid.UUID), id, &input)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, schedule)
}

func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "无效的日程ID", nil)
		return
	}

	if err := h.scheduleService.DeleteSchedule(userID.(uuid.UUID), id); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, nil)
}

func (h *ScheduleHandler) GetUserSchedules(c *gin.Context) {
	userID, _ := c.Get("user_id")

	schedules, err := h.scheduleService.GetUserSchedules(userID.(uuid.UUID))
	if err != nil {
		response.InternalServerError(c, "获取日程失败", err.Error())
		return
	}

	response.Success(c, schedules)
}

func (h *ScheduleHandler) GetUserAvailability(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.BadRequest(c, "无效的用户ID", nil)
		return
	}

	dayOfWeek := c.Query("day_of_week")
	if dayOfWeek == "" {
		dayOfWeek = strings.ToLower(time.Now().Weekday().String())
	}

	schedules, err := h.scheduleService.GetUserAvailability(userID, models.DayOfWeek(dayOfWeek))
	if err != nil {
		response.InternalServerError(c, "获取可用时间失败", err.Error())
		return
	}

	response.Success(c, schedules)
}

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) GetTeacherStats(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var startDate, endDate time.Time
	var err error

	if c.Query("start_date") != "" {
		startDate, err = time.Parse("2006-01-02", c.Query("start_date"))
		if err != nil {
			startDate = time.Now().AddDate(0, -1, 0)
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0)
	}

	if c.Query("end_date") != "" {
		endDate, err = time.Parse("2006-01-02", c.Query("end_date"))
		if err != nil {
			endDate = time.Now()
		}
	} else {
		endDate = time.Now()
	}

	stats, err := h.statsService.GetTeacherStats(userID.(uuid.UUID), startDate, endDate)
	if err != nil {
		response.InternalServerError(c, "获取统计数据失败", err.Error())
		return
	}

	response.Success(c, stats)
}

func (h *StatsHandler) GetMonthlyReport(c *gin.Context) {
	userID, _ := c.Get("user_id")

	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month()))))

	report, err := h.statsService.GetMonthlyReport(userID.(uuid.UUID), year, month)
	if err != nil {
		response.InternalServerError(c, "获取月度报告失败", err.Error())
		return
	}

	response.Success(c, report)
}

func (h *StatsHandler) ExportReport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "报表导出功能",
		"format":  "pdf",
	})
}
