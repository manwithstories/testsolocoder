package handler

import (
	"net/http"
	"strconv"
	"strings"

	"luxury-trading-platform/internal/model"
	"luxury-trading-platform/internal/service"
	resp "luxury-trading-platform/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type AuthenticationHandler struct {
	authService *service.AuthenticationService
	pdfService  *service.PDFService
}

func NewAuthenticationHandler(authService *service.AuthenticationService, pdfService *service.PDFService) *AuthenticationHandler {
	return &AuthenticationHandler{
		authService: authService,
		pdfService:  pdfService,
	}
}

func (h *AuthenticationHandler) CreateAuthentication(c *gin.Context) {
	buyerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	var req struct {
		OrderID uint `json:"order_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	auth, err := h.authService.CreateAuthentication(c.Request.Context(), buyerID.(uint), req.OrderID)
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Created(c, auth)
}

func (h *AuthenticationHandler) GetAuthentication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	auth, err := h.authService.GetAuthentication(c.Request.Context(), uint(id))
	if err != nil {
		resp.NotFound(c, err.Error())
		return
	}

	resp.Success(c, auth)
}

func (h *AuthenticationHandler) GetAuthenticationByOrder(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("order_id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	auth, err := h.authService.GetAuthenticationByOrder(c.Request.Context(), uint(orderID))
	if err != nil {
		resp.NotFound(c, err.Error())
		return
	}

	resp.Success(c, auth)
}

func (h *AuthenticationHandler) ListAuthentications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	role, _ := c.Get("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	var authenticatorID, buyerID *uint
	if exists {
		uid := userID.(uint)
		if role == model.RoleAuthenticator {
			authenticatorID = &uid
		} else if role == model.RoleBuyer {
			buyerID = &uid
		}
	}

	auths, total, err := h.authService.ListAuthentications(page, pageSize, model.AuthenticationStatus(status), authenticatorID, buyerID)
	if err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.SuccessWithPage(c, auths, total, page, pageSize)
}

func (h *AuthenticationHandler) AcceptAuthentication(c *gin.Context) {
	authenticatorID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	auth, err := h.authService.AcceptAuthentication(c.Request.Context(), uint(id), authenticatorID.(uint))
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			resp.Forbidden(c, err.Error())
			return
		}
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, auth)
}

func (h *AuthenticationHandler) CompleteAuthentication(c *gin.Context) {
	authenticatorID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	var req service.CompleteAuthenticationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	reportFile := ""
	if req.ReportFile == "" {
		pdfData := &service.AuthenticationReportData{
			ReportID:          strconv.FormatUint(uint64(id), 10),
			ProductTitle:      "",
			ProductCategory:   "",
			BrandName:         "",
			AuthenticatorName: "",
			AuthenticatorID:   "",
			Result:            string(req.Result),
			ResultCN:          service.GetResultCN(req.Result),
			ReportContent:     req.ReportContent,
			Notes:             req.AuthenticatorNotes,
		}
		reportFile, _ = h.pdfService.GenerateAuthenticationReport(pdfData)
	}

	if reportFile != "" {
		req.ReportFile = reportFile
	}

	auth, err := h.authService.CompleteAuthentication(c.Request.Context(), uint(id), authenticatorID.(uint), &req)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			resp.Forbidden(c, err.Error())
			return
		}
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, auth)
}

func (h *AuthenticationHandler) RejectAuthentication(c *gin.Context) {
	authenticatorID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.BadRequest(c, err)
		return
	}

	auth, err := h.authService.RejectAuthentication(c.Request.Context(), uint(id), authenticatorID.(uint), req.Reason)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			resp.Forbidden(c, err.Error())
			return
		}
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, auth)
}

func (h *AuthenticationHandler) CancelAuthentication(c *gin.Context) {
	buyerID, exists := c.Get("user_id")
	if !exists {
		resp.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	if err := h.authService.CancelAuthentication(c.Request.Context(), uint(id), buyerID.(uint)); err != nil {
		resp.InternalError(c, err)
		return
	}

	resp.Success(c, gin.H{"message": "authentication cancelled successfully"})
}

func (h *AuthenticationHandler) DownloadReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		resp.BadRequest(c, err)
		return
	}

	auth, err := h.authService.GetAuthentication(c.Request.Context(), uint(id))
	if err != nil {
		resp.NotFound(c, err.Error())
		return
	}

	if auth.ReportFile == "" {
		resp.NotFound(c, "report file not found")
		return
	}

	filePath := "." + auth.ReportFile
	c.File(filePath)
	c.Status(http.StatusOK)
}
