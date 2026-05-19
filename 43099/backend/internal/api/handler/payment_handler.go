package handler

import (
	"net/http"
	"strconv"
	"venue-booking/internal/dto"
	"venue-booking/internal/service"
	"venue-booking/pkg/excel"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
	logService     *service.OperationLogService
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		paymentService: service.NewPaymentService(),
		logService:     service.NewOperationLogService(),
	}
}

func (h *PaymentHandler) ConfirmPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid payment ID"))
		return
	}

	var req dto.ConfirmPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	payment, err := h.paymentService.ConfirmPayment(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, err.Error()))
		return
	}

	userID, _ := c.Get("userID")
	h.logService.Log(c, userID.(uint), "confirm_payment", "payment", map[string]interface{}{
		"payment_id": payment.ID,
		"order_id":   payment.OrderID,
		"amount":     payment.Amount,
		"method":     payment.PaymentMethod,
	})

	c.JSON(http.StatusOK, dto.Success(payment))
}

func (h *PaymentHandler) List(c *gin.Context) {
	var req dto.PaymentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	payments, total, err := h.paymentService.List(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to get payments"))
		return
	}

	c.JSON(http.StatusOK, dto.Success(dto.PaginationResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     payments,
	}))
}

func (h *PaymentHandler) Export(c *gin.Context) {
	var req dto.PaymentExportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request parameters"))
		return
	}

	payments, err := h.paymentService.ExportForDateRange(req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to export payments"))
		return
	}

	f, err := excel.ExportPayments(payments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to generate Excel file"))
		return
	}

	filename := "payment_records_" + req.StartDate + "_" + req.EndDate + ".xlsx"
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	f.Write(c.Writer)
}
