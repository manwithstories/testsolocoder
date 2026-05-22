package handler

import (
	"car-rental/internal/config"
	"car-rental/internal/middleware"
	"car-rental/internal/service"
	"car-rental/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(cfg *config.EmailConfig) *MessageHandler {
	return &MessageHandler{
		messageService: service.NewMessageService(cfg),
	}
}

func (h *MessageHandler) GetUserMessages(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	page, pageSize, _, _ := utils.ParsePageParams(c)

	messages, total, err := h.messageService.GetUserMessages(user.UserID, page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, messages, total, page, pageSize)
}

func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	count, err := h.messageService.GetUnreadCount(user.UserID)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"count": count})
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.messageService.MarkAsRead(uint(id), user.UserID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *MessageHandler) MarkAllAsRead(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	err := h.messageService.MarkAllAsRead(user.UserID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	user := middleware.GetUserContext(c)
	if user == nil {
		utils.Unauthorized(c, "未登录")
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.messageService.DeleteMessage(uint(id), user.UserID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{
		statsService: service.NewStatsService(),
	}
}

func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

func (h *StatsHandler) GetRevenueStats(c *gin.Context) {
	var startDate, endDate *int64
	if start := c.Query("start_date"); start != "" {
		if ts, err := strconv.ParseInt(start, 10, 64); err == nil {
			startDate = &ts
		}
	}
	if end := c.Query("end_date"); end != "" {
		if ts, err := strconv.ParseInt(end, 10, 64); err == nil {
			endDate = &ts
		}
	}

	_ = startDate
	_ = endDate

	stats, err := h.statsService.GetRevenueStats(nil, nil)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

type PromoHandler struct {
	promoService *service.PromoService
}

func NewPromoHandler() *PromoHandler {
	return &PromoHandler{
		promoService: service.NewPromoService(),
	}
}

func (h *PromoHandler) CreatePromo(c *gin.Context) {
	var req service.CreatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	promo, err := h.promoService.CreatePromo(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, promo)
}

func (h *PromoHandler) GetPromoByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	promo, err := h.promoService.GetPromoByID(uint(id))
	if err != nil {
		utils.NotFound(c, "优惠码不存在")
		return
	}

	utils.Success(c, promo)
}

func (h *PromoHandler) GetPromoByCode(c *gin.Context) {
	code := c.Param("code")

	promo, err := h.promoService.GetPromoByCode(code)
	if err != nil {
		utils.NotFound(c, "优惠码不存在")
		return
	}

	utils.Success(c, promo)
}

func (h *PromoHandler) GetAllPromos(c *gin.Context) {
	page, pageSize, _, _ := utils.ParsePageParams(c)
	keyword := c.Query("keyword")

	promos, total, err := h.promoService.GetAllPromos(page, pageSize, keyword)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, promos, total, page, pageSize)
}

func (h *PromoHandler) UpdatePromo(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.promoService.UpdatePromo(uint(id), updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *PromoHandler) DeletePromo(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.promoService.DeletePromo(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

type PricingHandler struct {
	pricingService *service.PricingService
}

func NewPricingHandler() *PricingHandler {
	return &PricingHandler{
		pricingService: service.NewPricingService(),
	}
}

func (h *PricingHandler) CreateRule(c *gin.Context) {
	var req service.CreatePricingRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	rule, err := h.pricingService.CreateRule(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, rule)
}

func (h *PricingHandler) GetRuleByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	rule, err := h.pricingService.GetRuleByID(uint(id))
	if err != nil {
		utils.NotFound(c, "规则不存在")
		return
	}

	utils.Success(c, rule)
}

func (h *PricingHandler) GetAllRules(c *gin.Context) {
	page, pageSize, _, _ := utils.ParsePageParams(c)
	ruleType := c.Query("type")

	rules, total, err := h.pricingService.GetAllRules(page, pageSize, ruleType)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, rules, total, page, pageSize)
}

func (h *PricingHandler) UpdateRule(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.pricingService.UpdateRule(uint(id), updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *PricingHandler) DeleteRule(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.pricingService.DeleteRule(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *PricingHandler) ToggleRuleActive(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.pricingService.ToggleActive(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}
