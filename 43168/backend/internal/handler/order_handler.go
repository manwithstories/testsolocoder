package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/service"
	"furniture-platform/pkg/response"
)

// OrderHandler 订单 HTTP 处理器
type OrderHandler struct {
	service *service.OrderService
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{service: svc}
}

// toOrderItemDTO 将订单项模型转为 DTO
func toOrderItemDTO(it *model.OrderItem) dto.OrderItemDTO {
	opts, _ := it.ParseCustomOptions()
	optDTOs := make([]dto.CustomOptionDTO, 0, len(opts))
	for _, o := range opts {
		optDTOs = append(optDTOs, dto.CustomOptionDTO{
			OptionType:  o.OptionType,
			OptionValue: o.OptionValue,
			PriceAdjust: o.PriceAdjust,
		})
	}
	return dto.OrderItemDTO{
		ID:            it.ID,
		OrderID:       it.OrderID,
		ProductID:     it.ProductID,
		ProductName:   it.ProductName,
		BasePrice:     it.BasePrice,
		CustomOptions: optDTOs,
		Quantity:      it.Quantity,
		Subtotal:      it.Subtotal,
	}
}

// toOrderResponse 将订单模型转为响应 DTO
func toOrderResponse(o *model.Order) dto.OrderResponse {
	return dto.OrderResponse{
		ID:             o.ID,
		OrderNo:        o.OrderNo,
		OwnerID:        o.OwnerID,
		DesignerID:     o.DesignerID,
		ManufacturerID: o.ManufacturerID,
		TotalAmount:    o.TotalAmount,
		Discount:       o.Discount,
		FinalAmount:    o.FinalAmount,
		Status:         o.Status,
		Address:        o.Address,
		ContactName:    o.ContactName,
		ContactPhone:   o.ContactPhone,
		Remark:         o.Remark,
		CreatedAt:      o.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      o.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// toOrderDetailResponse 转为订单详情响应
func toOrderDetailResponse(o *model.Order, items []*model.OrderItem, histories []*model.OrderHistory) dto.OrderDetailResponse {
	itemDTOs := make([]dto.OrderItemDTO, 0, len(items))
	for _, it := range items {
		itemDTOs = append(itemDTOs, toOrderItemDTO(it))
	}
	hisDTOs := make([]dto.OrderHistoryDTO, 0, len(histories))
	for _, h := range histories {
		hisDTOs = append(hisDTOs, dto.OrderHistoryDTO{
			ID:           h.ID,
			OrderID:      h.OrderID,
			Status:       h.Status,
			OperatorID:   h.OperatorID,
			OperatorRole: h.OperatorRole,
			Remark:       h.Remark,
			CreatedAt:    h.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return dto.OrderDetailResponse{
		OrderResponse: toOrderResponse(o),
		Items:         itemDTOs,
		Histories:     hisDTOs,
	}
}

// parseOrderID 解析订单 ID 路径参数
func parseOrderID(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// CreateInquiry 业主发起询价
func (h *OrderHandler) CreateInquiry(c *gin.Context) {
	userID, _, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.CreateInquiryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	order, items, err := h.service.CreateInquiry(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toOrderDetailResponse(order, items, nil))
}

// Quote 厂商报价
func (h *OrderHandler) Quote(c *gin.Context) {
	userID, _, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := parseOrderID(c)
	if err != nil {
		response.BadRequest(c, "订单 ID 格式错误")
		return
	}

	var req dto.QuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	order, items, err := h.service.Quote(id, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	histories, _ := h.service.GetHistories(order.ID)
	response.Success(c, toOrderDetailResponse(order, items, histories))
}

// ConfirmOrder 业主确认下单
func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
	userID, _, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := parseOrderID(c)
	if err != nil {
		response.BadRequest(c, "订单 ID 格式错误")
		return
	}

	var body struct {
		Remark string `json:"remark"`
	}
	_ = c.ShouldBindJSON(&body)

	order, err := h.service.ConfirmOrder(id, userID, body.Remark)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toOrderResponse(order))
}

// GetOrder 获取订单详情
func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID, role, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := parseOrderID(c)
	if err != nil {
		response.BadRequest(c, "订单 ID 格式错误")
		return
	}

	order, items, histories, err := h.service.GetByID(id, userID, role)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, toOrderDetailResponse(order, items, histories))
}

// ListOrders 分页查询订单列表
func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID, role, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	var req dto.OrderListRequest
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

	resp := make([]dto.OrderResponse, 0, len(list))
	for _, o := range list {
		resp = append(resp, toOrderResponse(o))
	}
	response.Success(c, dto.OrderListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     resp,
	})
}

// CancelOrder 取消订单
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID, role, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := parseOrderID(c)
	if err != nil {
		response.BadRequest(c, "订单 ID 格式错误")
		return
	}

	var body struct {
		Remark string `json:"remark"`
	}
	_ = c.ShouldBindJSON(&body)

	order, err := h.service.CancelOrder(id, userID, role, body.Remark)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toOrderResponse(order))
}

// UpdateStatus 更新订单状态
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	userID, role, ok := currentUser(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := parseOrderID(c)
	if err != nil {
		response.BadRequest(c, "订单 ID 格式错误")
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	order, err := h.service.UpdateStatus(id, userID, role, req.Status, req.Remark)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toOrderResponse(order))
}
