package handler

import (
	"encoding/json"
	"strconv"
	"time"

	"event-platform/internal/dto"
	"event-platform/internal/logger"
	"event-platform/internal/model"
	"event-platform/internal/service"
	"event-platform/pkg/response"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	svc *service.EventService
}

func NewEventHandler(svc *service.EventService) *EventHandler { return &EventHandler{svc: svc} }

func (h *EventHandler) Create(c *gin.Context) {
	var req dto.EventCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("event create bind error: %v", err)
		response.BadRequest(c, "参数错误")
		return
	}
	start, _ := time.Parse(time.RFC3339, req.StartDate)
	end, _ := time.Parse(time.RFC3339, req.EndDate)
	dl, _ := time.Parse(time.RFC3339, req.RegistrationDeadline)
	uid, _ := c.Get("user_id")
	e := &model.Event{
		Name:                 req.Name,
		Description:          req.Description,
		Location:             req.Location,
		StartDate:            start,
		EndDate:              end,
		RegistrationDeadline: dl,
		Status:               "draft",
		Organizer:            req.Organizer,
		CoverImage:           req.CoverImage,
		IsPublished:          false,
		CreatedBy:            uid.(uint),
	}
	var items []model.EventItem
	for _, it := range req.Items {
		items = append(items, model.EventItem{
			Name:          it.Name,
			Category:      it.Category,
			Gender:        it.Gender,
			MinAge:        it.MinAge,
			MaxAge:        it.MaxAge,
			Quota:         it.Quota,
			WaitlistQuota: it.WaitlistQuota,
			Fee:           it.Fee,
			Requirements:  it.Requirements,
			Status:        "open",
		})
	}
	ev, err := h.svc.Create(e, items)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, ev)
}

func (h *EventHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ev, err := h.svc.Get(uint(id))
	if err != nil {
		response.NotFound(c, "赛事不存在")
		return
	}
	var req dto.EventCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	ev.Name = req.Name
	ev.Description = req.Description
	ev.Location = req.Location
	start, _ := time.Parse(time.RFC3339, req.StartDate)
	end, _ := time.Parse(time.RFC3339, req.EndDate)
	dl, _ := time.Parse(time.RFC3339, req.RegistrationDeadline)
	ev.StartDate = start
	ev.EndDate = end
	ev.RegistrationDeadline = dl
	ev.Organizer = req.Organizer
	ev.CoverImage = req.CoverImage
	if err := h.svc.Update(ev); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, ev)
}

func (h *EventHandler) Publish(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Publish(uint(id)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *EventHandler) Unpublish(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Unpublish(uint(id)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *EventHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ev, err := h.svc.Get(uint(id))
	if err != nil {
		response.NotFound(c, "赛事不存在")
		return
	}
	response.OK(c, ev)
}

func (h *EventHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	all := c.Query("all") == "1"
	var list []model.Event
	var total int64
	var err error
	if all {
		list, total, err = h.svc.ListAll(page, size)
	} else {
		list, total, err = h.svc.ListPublished(page, size)
	}
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, dto.PagedData{Total: total, Page: page, PageSize: size, List: list})
}

type RegistrationHandler struct {
	svc *service.RegistrationService
}

func NewRegistrationHandler(svc *service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{svc: svc}
}

func (h *RegistrationHandler) Register(c *gin.Context) {
	var req dto.RegistrationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	uid, _ := c.Get("user_id")
	var teamMembers string
	if req.TeamMembers != nil {
		b, _ := json.Marshal(req.TeamMembers)
		teamMembers = string(b)
	}
	reg, err := h.svc.Register(uid.(uint), req.EventItemID, req.RegType, req.TeamName, teamMembers)
	if err != nil {
		response.Fail(c, 400, 40001, err.Error())
		return
	}
	response.OK(c, reg)
}

func (h *RegistrationHandler) MyList(c *gin.Context) {
	uid, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	list, total, err := h.svc.ListByUser(uid.(uint), page, size)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, dto.PagedData{Total: total, Page: page, PageSize: size, List: list})
}

func (h *RegistrationHandler) ListByEvent(c *gin.Context) {
	eid, _ := strconv.ParseUint(c.Param("event_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	list, total, err := h.svc.ListByEvent(uint(eid), page, size)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, dto.PagedData{Total: total, Page: page, PageSize: size, List: list})
}

func (h *RegistrationHandler) ConfirmWaitlist(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.ConfirmWaitlist(uint(id)); err != nil {
		response.Fail(c, 400, 40001, err.Error())
		return
	}
	response.OK(c, nil)
}

type ScoreHandler struct {
	svc *service.ScoreService
}

func NewScoreHandler(svc *service.ScoreService) *ScoreHandler { return &ScoreHandler{svc: svc} }

func (h *ScoreHandler) Entry(c *gin.Context) {
	var req dto.ScoreEntryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	var list []model.Score
	for _, it := range req.Scores {
		list = append(list, model.Score{
			EventID:     req.EventID,
			EventItemID: req.EventItemID,
			UserID:      it.UserID,
			Score:       it.Score,
			TimeUsed:    it.TimeUsed,
			Remarks:     it.Remarks,
			IsValid:     true,
		})
	}
	if err := h.svc.Entry(list); err != nil {
		response.Fail(c, 400, 40001, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *ScoreHandler) ListByItem(c *gin.Context) {
	iid, _ := strconv.ParseUint(c.Param("item_id"), 10, 64)
	list, err := h.svc.ListByItem(uint(iid))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, list)
}

func (h *ScoreHandler) MyList(c *gin.Context) {
	uid, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	list, total, err := h.svc.ListByUser(uid.(uint), page, size)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, dto.PagedData{Total: total, Page: page, PageSize: size, List: list})
}

type MessageHandler struct {
	svc *service.MessageService
}

func NewMessageHandler(svc *service.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

func (h *MessageHandler) List(c *gin.Context) {
	uid, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	list, total, err := h.svc.List(uid.(uint), page, size)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, dto.PagedData{Total: total, Page: page, PageSize: size, List: list})
}

func (h *MessageHandler) MarkRead(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	uid, _ := c.Get("user_id")
	if err := h.svc.MarkRead(uint(id), uid.(uint)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *MessageHandler) MarkAllRead(c *gin.Context) {
	uid, _ := c.Get("user_id")
	if err := h.svc.MarkAllRead(uid.(uint)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *MessageHandler) UnreadCount(c *gin.Context) {
	uid, _ := c.Get("user_id")
	n, err := h.svc.UnreadCount(uid.(uint))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"count": n})
}

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler { return &UserHandler{svc: svc} }

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	user, err := h.svc.Register(req.Username, req.Password, req.RealName, req.IdCard, req.Phone, req.Email)
	if err != nil {
		response.Fail(c, 400, 40001, err.Error())
		return
	}
	response.OK(c, gin.H{"user_id": user.ID, "username": user.Username})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	token, user, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		response.Fail(c, 401, 40100, err.Error())
		return
	}
	response.OK(c, dto.LoginResp{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Role:     string(user.Role),
		RealName: user.RealName,
		Verified: user.Verified,
	})
}

func (h *UserHandler) Profile(c *gin.Context) {
	uid, _ := c.Get("user_id")
	user, err := h.svc.GetByID(uid.(uint))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.OK(c, user)
}

func (h *UserHandler) Verify(c *gin.Context) {
	var req dto.VerifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	uid, _ := c.Get("user_id")
	if err := h.svc.Verify(uid.(uint), req.RealName, req.IdCard); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, nil)
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	list, total, err := h.svc.List(page, size)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, dto.PagedData{Total: total, Page: page, PageSize: size, List: list})
}
