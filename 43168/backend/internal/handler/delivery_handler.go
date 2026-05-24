package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/service"
	"furniture-platform/pkg/response"
)

// DeliveryHandler 配送安装 HTTP 处理器
type DeliveryHandler struct {
	service *service.DeliveryService
}

// NewDeliveryHandler 创建配送安装处理器
func NewDeliveryHandler(svc *service.DeliveryService) *DeliveryHandler {
	return &DeliveryHandler{service: svc}
}

// toDeliveryResponse 将模型转为响应 DTO
func toDeliveryResponse(d *model.Delivery) dto.DeliveryResponse {
	return dto.DeliveryResponse{
		ID:           d.ID,
		OrderID:      d.OrderID,
		OwnerID:      d.OwnerID,
		Type:         d.Type,
		TimeSlot:     d.TimeSlot,
		Address:      d.Address,
		ContactName:  d.ContactName,
		ContactPhone: d.ContactPhone,
		Status:       d.Status,
		InstallerID:  d.InstallerID,
		Remark:       d.Remark,
		CreatedAt:    d.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    d.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// parseDeliveryID 解析配送安装 ID 路径参数
func parseDeliveryID(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// CreateDelivery 业主创建配送安装预约
func (h *DeliveryHandler) CreateDelivery(c *gin.Context) {
	userID, _, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.CreateDeliveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	delivery, err := h.service.CreateDelivery(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toDeliveryResponse(delivery))
}

// UpdateDelivery 更新配送安装（预约信息/确认/完成/取消）
func (h *DeliveryHandler) UpdateDelivery(c *gin.Context) {
	userID, role, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := parseDeliveryID(c)
	if err != nil {
		response.BadRequest(c, "配送安装 ID 格式错误")
		return
	}

	var req dto.UpdateDeliveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	delivery, err := h.service.UpdateDelivery(id, userID, role, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toDeliveryResponse(delivery))
}

// GetDelivery 获取配送安装详情
func (h *DeliveryHandler) GetDelivery(c *gin.Context) {
	userID, role, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := parseDeliveryID(c)
	if err != nil {
		response.BadRequest(c, "配送安装 ID 格式错误")
		return
	}

	delivery, err := h.service.GetByID(id, userID, role)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, toDeliveryResponse(delivery))
}

// ListDeliveries 分页查询配送安装列表
func (h *DeliveryHandler) ListDeliveries(c *gin.Context) {
	userID, role, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.DeliveryListRequest
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

	resp := make([]dto.DeliveryResponse, 0, len(list))
	for _, d := range list {
		resp = append(resp, toDeliveryResponse(d))
	}
	response.Success(c, dto.DeliveryListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     resp,
	})
}
