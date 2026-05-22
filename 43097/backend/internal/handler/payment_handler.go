package handler

import (
	"hotel-system/internal/dto"
	"hotel-system/internal/model"
	"hotel-system/internal/pkg/logger"
	"hotel-system/internal/service"
	"hotel-system/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req dto.PaymentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("创建支付记录参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	payment, err := h.paymentService.CreatePayment(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	paymentResponse := convertToPaymentResponse(payment)
	utils.Success(c, paymentResponse)
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的支付记录ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的支付记录ID")
		return
	}

	payment, err := h.paymentService.GetPaymentByID(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	paymentResponse := convertToPaymentResponse(payment)
	utils.Success(c, paymentResponse)
}

func (h *PaymentHandler) ListPayments(c *gin.Context) {
	var req dto.PaymentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warnf("获取支付记录列表参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	payments, total, err := h.paymentService.ListPayments(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	var paymentResponses []dto.PaymentResponse
	for _, payment := range payments {
		paymentResponses = append(paymentResponses, convertToPaymentResponse(&payment))
	}

	utils.PageResult(c, paymentResponses, total, req.GetPage(), req.GetPageSize())
}

func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	var req dto.PaymentRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("退款参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	refundPayment, err := h.paymentService.RefundPayment(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	paymentResponse := convertToPaymentResponse(refundPayment)
	utils.Success(c, paymentResponse)
}

func (h *PaymentHandler) GetOrderPayments(c *gin.Context) {
	var req dto.OrderPaymentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warnf("获取订单支付记录参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := h.paymentService.GetOrderPayments(req.OrderType, req.OrderID)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, result)
}

func (h *PaymentHandler) GeneratePaymentVoucher(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的支付记录ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的支付记录ID")
		return
	}

	voucher, err := h.paymentService.GeneratePaymentVoucher(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, voucher)
}

func convertToPaymentResponse(payment *model.Payment) dto.PaymentResponse {
	return dto.PaymentResponse{
		ID:            payment.ID,
		PaymentNo:     payment.PaymentNo,
		OrderType:     payment.OrderType,
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		PaymentType:   payment.PaymentType,
		Status:        payment.Status,
		TransactionID: payment.TransactionID,
		Remark:        payment.Remark,
		CreatedAt:     payment.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
