package handler

import (
	"net/http"
	"strconv"
	"furniture-platform/internal/dto"
	"furniture-platform/internal/repository"
	"furniture-platform/internal/service"
	"furniture-platform/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductionHandler struct {
	service     *service.ProductionService
	orderRepo   *repository.OrderRepo
}

func NewProductionHandler(svc *service.ProductionService, orderRepo *repository.OrderRepo) *ProductionHandler {
	return &ProductionHandler{service: svc, orderRepo: orderRepo}
}

func (h *ProductionHandler) GetProduction(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id <= 0 {
		response.BadRequest(c, "无效的ID")
		return
	}
	p, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "生产记录不存在")
		return
	}
	order, _ := h.orderRepo.GetByID(c.Request.Context(), p.OrderID)
	orderNo := ""
	if order != nil {
		orderNo = order.OrderNo
	}
	response.Success(c, service.ToProductionResponse(p, orderNo))
}

func (h *ProductionHandler) ListProductions(c *gin.Context) {
	var req dto.ProductionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	filter := repository.ProductionListFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
		OrderID:  req.OrderID,
		Status:   req.Status,
	}
	list, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, "查询失败")
		return
	}
	result := make([]dto.ProductionResponse, 0, len(list))
	for _, p := range list {
		order, _ := h.orderRepo.GetByID(c.Request.Context(), p.OrderID)
		orderNo := ""
		if order != nil {
			orderNo = order.OrderNo
		}
		result = append(result, service.ToProductionResponse(&p, orderNo))
	}
	response.Success(c, dto.ProductionListResponse{List: result, Total: total})
}

func (h *ProductionHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if id <= 0 {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req dto.UpdateProductionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	operatorID, _ := c.Get("user_id")
	p, err := h.service.UpdateStatus(c.Request.Context(), id, &req, operatorID.(int64))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	order, _ := h.orderRepo.GetByID(c.Request.Context(), p.OrderID)
	orderNo := ""
	if order != nil {
		orderNo = order.OrderNo
	}
	response.Success(c, service.ToProductionResponse(p, orderNo))
}

func (h *ProductionHandler) CreateProduction(c *gin.Context) {
	orderID, _ := strconv.ParseInt(c.Query("order_id"), 10, 64)
	if orderID <= 0 {
		response.BadRequest(c, "无效的订单ID")
		return
	}
	operatorID, _ := c.Get("user_id")
	p, err := h.service.Create(c.Request.Context(), orderID, operatorID.(int64))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	order, _ := h.orderRepo.GetByID(c.Request.Context(), p.OrderID)
	orderNo := ""
	if order != nil {
		orderNo = order.OrderNo
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    service.ToProductionResponse(p, orderNo),
	})
}
