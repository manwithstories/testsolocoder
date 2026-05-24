package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/service"
	"furniture-platform/pkg/response"
)

// TicketHandler 售后工单 HTTP 处理器
type TicketHandler struct {
	service *service.TicketService
}

// NewTicketHandler 创建售后工单处理器
func NewTicketHandler(svc *service.TicketService) *TicketHandler {
	return &TicketHandler{service: svc}
}

// toTicketResponse 将模型转为响应 DTO
func toTicketResponse(t *model.Ticket) dto.TicketResponse {
	images, _ := t.ParseImages()
	if images == nil {
		images = []string{}
	}
	return dto.TicketResponse{
		ID:        t.ID,
		OrderID:   t.OrderID,
		OwnerID:   t.OwnerID,
		Type:      t.Type,
		Title:     t.Title,
		Content:   t.Content,
		Status:    t.Status,
		Images:    images,
		CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: t.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// parseTicketID 解析工单 ID 路径参数
func parseTicketID(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// CreateTicket 业主创建售后工单
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	userID, _, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	ticket, err := h.service.CreateTicket(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toTicketResponse(ticket))
}

// UpdateTicket 更新工单（状态流转/内容修改）
func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	userID, role, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := parseTicketID(c)
	if err != nil {
		response.BadRequest(c, "工单 ID 格式错误")
		return
	}

	var req dto.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	ticket, err := h.service.UpdateTicket(id, userID, role, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toTicketResponse(ticket))
}

// GetTicket 获取工单详情
func (h *TicketHandler) GetTicket(c *gin.Context) {
	userID, role, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := parseTicketID(c)
	if err != nil {
		response.BadRequest(c, "工单 ID 格式错误")
		return
	}

	ticket, err := h.service.GetByID(id, userID, role)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, toTicketResponse(ticket))
}

// ListTickets 分页查询工单列表
func (h *TicketHandler) ListTickets(c *gin.Context) {
	userID, role, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.TicketListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}
	req.UserID = userID
	req.Role = role

	list, total, err := h.service.List(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	resp := make([]dto.TicketResponse, 0, len(list))
	for _, t := range list {
		resp = append(resp, toTicketResponse(t))
	}
	response.Success(c, dto.TicketListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     resp,
	})
}
