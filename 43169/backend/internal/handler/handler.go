package handler

import (
	"net/http"
	"strconv"
	"time"

	"matchmaking-platform/internal/dto"
	"matchmaking-platform/internal/service"
	"matchmaking-platform/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	resp, err := h.userService.Register(&req)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	resp, err := h.userService.Login(&req)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, resp)
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	info, err := h.userService.GetUserInfo(userID)
	if err != nil {
		utils.Fail(c, 404, "用户不存在")
		return
	}

	utils.Success(c, info)
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	info, err := h.userService.GetUserInfo(uint(id))
	if err != nil {
		utils.Fail(c, 404, "用户不存在")
		return
	}

	utils.Success(c, info)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}

	if err := h.userService.UpdateProfile(userID, &req); err != nil {
		utils.Fail(c, 500, "更新失败: "+err.Error())
		return
	}

	utils.SuccessMsg(c, "更新成功")
}

func (h *UserHandler) Verify(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := h.userService.Verify(userID, &req); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessMsg(c, "认证信息已提交，等待审核")
}

func (h *UserHandler) ApproveVerify(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.userService.ApproveVerify(uint(id)); err != nil {
		utils.Fail(c, 500, "操作失败")
		return
	}

	utils.SuccessMsg(c, "审核通过")
}

func (h *UserHandler) RejectVerify(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.userService.RejectVerify(uint(id)); err != nil {
		utils.Fail(c, 500, "操作失败")
		return
	}

	utils.SuccessMsg(c, "已拒绝")
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, 400, "请上传文件")
		return
	}

	filename := utils.RandomString(32) + "_" + file.Filename
	filePath := "uploads/avatars/" + filename

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.Fail(c, 500, "文件保存失败")
		return
	}

	avatarURL := "/uploads/avatars/" + filename
	h.userService.UploadAvatar(userID, avatarURL)

	utils.Success(c, gin.H{"avatar": avatarURL})
}

func (h *UserHandler) UploadPhotos(c *gin.Context) {
	userID := c.GetUint("user_id")

	form, err := c.MultipartForm()
	if err != nil {
		utils.Fail(c, 400, "请上传文件")
		return
	}

	files := form.File["files"]
	var urls []string

	for _, file := range files {
		filename := utils.RandomString(32) + "_" + file.Filename
		filePath := "uploads/photos/" + filename

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			continue
		}
		urls = append(urls, "/uploads/photos/"+filename)
	}

	h.userService.UploadPhotos(userID, urls)

	utils.Success(c, gin.H{"photos": urls})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.DefaultQuery("keyword", "")

	users, total, err := h.userService.ListUsers(page, pageSize, keyword)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Page(c, users, total, page, pageSize)
}

func (h *UserHandler) DisableUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.userService.DisableUser(uint(id)); err != nil {
		utils.Fail(c, 500, "操作失败")
		return
	}

	utils.SuccessMsg(c, "已禁用")
}

func (h *UserHandler) EnableUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.userService.EnableUser(uint(id)); err != nil {
		utils.Fail(c, 500, "操作失败")
		return
	}

	utils.SuccessMsg(c, "已启用")
}

func (h *UserHandler) SendSmsCode(c *gin.Context) {
	phone := c.Query("phone")
	if !utils.ValidatePhone(phone) {
		utils.Fail(c, 400, "手机号格式不正确")
		return
	}

	code := utils.RandomCode(6)

	utils.Success(c, gin.H{
		"phone": utils.MaskPhone(phone),
		"code":  code,
		"expires_in": 300,
	})
}

type MatchHandler struct {
	matchService *service.MatchService
}

func NewMatchHandler() *MatchHandler {
	return &MatchHandler{
		matchService: service.NewMatchService(),
	}
}

func (h *MatchHandler) SmartMatch(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	results, total, err := h.matchService.SmartMatch(userID, page, pageSize)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	utils.Page(c, results, total, page, pageSize)
}

func (h *MatchHandler) FilterMatch(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.MatchFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}

	results, total, err := h.matchService.FilterMatch(userID, &req)
	if err != nil {
		utils.Fail(c, 500, err.Error())
		return
	}

	page, pageSize := utils.Paginate(req.Page, req.PageSize)
	utils.Page(c, results, total, page, pageSize)
}

func (h *MatchHandler) Favorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	targetID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.matchService.Favorite(userID, uint(targetID)); err != nil {
		utils.Fail(c, 500, "操作失败")
		return
	}

	utils.SuccessMsg(c, "操作成功")
}

func (h *MatchHandler) Block(c *gin.Context) {
	userID := c.GetUint("user_id")
	targetID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.matchService.Block(userID, uint(targetID)); err != nil {
		utils.Fail(c, 500, "操作失败")
		return
	}

	utils.SuccessMsg(c, "操作成功")
}

func (h *MatchHandler) GetFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	results, total, err := h.matchService.GetFavorites(userID, page, pageSize)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Page(c, results, total, page, pageSize)
}

func (h *MatchHandler) GetBlocked(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	results, total, err := h.matchService.GetBlocked(userID, page, pageSize)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Page(c, results, total, page, pageSize)
}

type DateHandler struct {
	dateService *service.DateService
}

func NewDateHandler() *DateHandler {
	return &DateHandler{
		dateService: service.NewDateService(),
	}
}

func (h *DateHandler) CreateInvite(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.DateInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	record, err := h.dateService.CreateInvite(userID, &req)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, record)
}

func (h *DateHandler) Accept(c *gin.Context) {
	userID := c.GetUint("user_id")
	dateID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.dateService.Accept(uint(dateID), userID); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessMsg(c, "已接受邀请")
}

func (h *DateHandler) Reject(c *gin.Context) {
	userID := c.GetUint("user_id")
	dateID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.dateService.Reject(uint(dateID), userID); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessMsg(c, "已拒绝")
}

func (h *DateHandler) Cancel(c *gin.Context) {
	userID := c.GetUint("user_id")
	dateID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.dateService.Cancel(uint(dateID), userID); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessMsg(c, "已取消")
}

func (h *DateHandler) Complete(c *gin.Context) {
	userID := c.GetUint("user_id")
	dateID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.dateService.Complete(uint(dateID), userID); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessMsg(c, "已完成")
}

func (h *DateHandler) GetUserDates(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	dates, total, err := h.dateService.GetUserDates(userID, page, pageSize)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Page(c, dates, total, page, pageSize)
}

func (h *DateHandler) CreateReview(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.DateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := h.dateService.CreateReview(userID, &req); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessMsg(c, "评价成功")
}

func (h *DateHandler) GetReviews(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	reviews, total, err := h.dateService.GetReviews(userID, page, pageSize)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Page(c, reviews, total, page, pageSize)
}

type MatchmakerHandler struct {
	service *service.MatchmakerService
}

func NewMatchmakerHandler() *MatchmakerHandler {
	return &MatchmakerHandler{
		service: service.NewMatchmakerService(),
	}
}

func (h *MatchmakerHandler) AddMember(c *gin.Context) {
	matchmakerID := c.GetUint("user_id")

	var req dto.MemberManageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}

	if err := h.service.AddMember(matchmakerID, &req); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessMsg(c, "添加成功")
}

func (h *MatchmakerHandler) RemoveMember(c *gin.Context) {
	matchmakerID := c.GetUint("user_id")
	memberID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.service.RemoveMember(matchmakerID, uint(memberID)); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessMsg(c, "已移除")
}

func (h *MatchmakerHandler) ListMembers(c *gin.Context) {
	matchmakerID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	members, total, err := h.service.ListMembers(matchmakerID, page, pageSize)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Page(c, members, total, page, pageSize)
}

func (h *MatchmakerHandler) CreateService(c *gin.Context) {
	matchmakerID := c.GetUint("user_id")

	var req dto.MatchmakerServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := h.service.CreateService(matchmakerID, &req); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessMsg(c, "服务创建成功")
}

func (h *MatchmakerHandler) UpdateProgress(c *gin.Context) {
	matchmakerID := c.GetUint("user_id")
	serviceID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var body struct {
		Progress int `json:"progress"`
	}
	c.ShouldBindJSON(&body)

	if err := h.service.UpdateServiceProgress(matchmakerID, uint(serviceID), body.Progress); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessMsg(c, "更新成功")
}

func (h *MatchmakerHandler) ListServices(c *gin.Context) {
	matchmakerID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	services, total, err := h.service.ListServices(matchmakerID, page, pageSize)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Page(c, services, total, page, pageSize)
}

func (h *MatchmakerHandler) GetStats(c *gin.Context) {
	matchmakerID := c.GetUint("user_id")

	stats, err := h.service.GetStats(matchmakerID)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Success(c, stats)
}

func (h *MatchmakerHandler) ListAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	matchmakers, total, err := h.service.ListAllMatchmakers(page, pageSize)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Page(c, matchmakers, total, page, pageSize)
}

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		chatService: service.NewChatService(),
	}
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	senderID := c.GetUint("user_id")

	var req dto.ChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	msg, err := h.chatService.SendMessage(senderID, req.ReceiverID, req.Type, req.Content)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, msg)
}

func (h *ChatHandler) GetHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	otherID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	messages, total, err := h.chatService.GetHistory(userID, uint(otherID), page, pageSize)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Page(c, messages, total, page, pageSize)
}

func (h *ChatHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	count, err := h.chatService.GetUnreadCount(userID)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Success(c, gin.H{"unread_count": count})
}

func (h *ChatHandler) GetSessions(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	sessions, total, err := h.chatService.GetSessions(userID, page, pageSize)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Page(c, sessions, total, page, pageSize)
}

func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	senderID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	h.chatService.MarkAsRead(uint(senderID), userID)
	utils.SuccessMsg(c, "已标记为已读")
}

func (h *ChatHandler) UploadChatFile(c *gin.Context) {
	fileType := c.DefaultQuery("type", "image")
	file, err := c.FormFile("file")
	if err != nil {
		utils.Fail(c, 400, "请上传文件")
		return
	}

	filename := utils.RandomString(32) + "_" + file.Filename
	dir := "uploads/chat/" + fileType + "/"
	filePath := dir + filename

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.Fail(c, 500, "文件保存失败")
		return
	}

	utils.Success(c, gin.H{
		"url":  "/uploads/chat/" + fileType + "/" + filename,
		"type": fileType,
	})
}

type MemberHandler struct {
	memberService *service.MemberService
}

func NewMemberHandler() *MemberHandler {
	return &MemberHandler{
		memberService: service.NewMemberService(),
	}
}

func (h *MemberHandler) GetBenefits(c *gin.Context) {
	benefits, err := h.memberService.GetBenefits()
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}
	utils.Success(c, benefits)
}

func (h *MemberHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.MemberUpgradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "参数错误")
		return
	}

	order, err := h.memberService.CreateOrder(userID, req.Level, req.Months)
	if err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.Success(c, order)
}

func (h *MemberHandler) PayOrder(c *gin.Context) {
	orderID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.memberService.PayOrder(uint(orderID)); err != nil {
		utils.Fail(c, 400, err.Error())
		return
	}

	utils.SuccessMsg(c, "支付成功")
}

func (h *MemberHandler) GetUserOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	orders, total, err := h.memberService.GetUserOrders(userID, page, pageSize)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Page(c, orders, total, page, pageSize)
}

func (h *MemberHandler) CheckInteractLimit(c *gin.Context) {
	userID := c.GetUint("user_id")

	count, _ := h.memberService.CountTodayInteracts(userID)
	utils.Success(c, gin.H{
		"today_count": count,
		"checked_at":  time.Now(),
	})
}

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{
		statsService: service.NewStatsService(),
	}
}

func (h *StatsHandler) GetPlatformStats(c *gin.Context) {
	stats, err := h.statsService.GetPlatformStats()
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}
	utils.Success(c, stats)
}

func (h *StatsHandler) GetDailyStats(c *gin.Context) {
	startDate, _ := time.Parse("2006-01-02", c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -14).Format("2006-01-02")))
	endDate, _ := time.Parse("2006-01-02", c.DefaultQuery("end_date", time.Now().Format("2006-01-02")))

	stats, err := h.statsService.GetDailyStats(startDate, endDate)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}
	utils.Success(c, stats)
}

func (h *StatsHandler) GetMatchmakerStats(c *gin.Context) {
	startDate, _ := time.Parse("2006-01-02", c.DefaultQuery("start_date", ""))
	endDate, _ := time.Parse("2006-01-02", c.DefaultQuery("end_date", ""))
	matchmakerID, _ := strconv.ParseUint(c.DefaultQuery("matchmaker_id", "0"), 10, 64)

	stats, err := h.statsService.GetMatchmakerStats(startDate, endDate, uint(matchmakerID))
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}
	utils.Success(c, stats)
}

func (h *StatsHandler) ExportExcel(c *gin.Context) {
	startDate, _ := time.Parse("2006-01-02", c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02")))
	endDate, _ := time.Parse("2006-01-02", c.DefaultQuery("end_date", time.Now().Format("2006-01-02")))

	data, err := h.statsService.ExportExcel(startDate, endDate)
	if err != nil {
		utils.Fail(c, 500, "导出失败")
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=report.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *StatsHandler) ExportPDF(c *gin.Context) {
	startDate, _ := time.Parse("2006-01-02", c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02")))
	endDate, _ := time.Parse("2006-01-02", c.DefaultQuery("end_date", time.Now().Format("2006-01-02")))

	data, err := h.statsService.ExportPDF(startDate, endDate)
	if err != nil {
		utils.Fail(c, 500, "导出失败")
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=report.pdf")
	c.Data(http.StatusOK, "application/pdf", data)
}

func (h *StatsHandler) GetSystemLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filter := map[string]interface{}{}
	if module := c.Query("module"); module != "" {
		filter["module"] = module
	}

	logs, total, err := h.statsService.GetSystemLogs(page, pageSize, filter)
	if err != nil {
		utils.Fail(c, 500, "查询失败")
		return
	}

	utils.Page(c, logs, total, page, pageSize)
}
