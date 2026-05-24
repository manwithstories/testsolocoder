package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"business-registration-platform/models"
	"business-registration-platform/services"
	"business-registration-platform/utils"
)

type ApplicationController struct {
	applicationService *services.ApplicationService
	notificationService *services.NotificationService
}

func NewApplicationController() *ApplicationController {
	return &ApplicationController{
		applicationService:  services.NewApplicationService(),
		notificationService: services.NewNotificationService(),
	}
}

func (ctrl *ApplicationController) CreateApplication(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var req services.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	req.EntrepreneurID = userID.(uint)

	application, err := ctrl.applicationService.CreateApplication(&req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, application)
}

func (ctrl *ApplicationController) SubmitApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid application ID")
		return
	}

	if err := ctrl.applicationService.SubmitApplication(uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *ApplicationController) GetApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid application ID")
		return
	}

	application, err := ctrl.applicationService.GetApplicationByID(uint(id))
	if err != nil {
		utils.NotFound(c, "Application not found")
		return
	}

	utils.Success(c, application)
}

func (ctrl *ApplicationController) GetApplicationList(c *gin.Context) {
	var query services.ApplicationListQuery

	if v := c.Query("page"); v != "" {
		query.Page, _ = strconv.Atoi(v)
	}
	if v := c.Query("pageSize"); v != "" {
		query.PageSize, _ = strconv.Atoi(v)
	}
	query.Status = models.ApplicationStatus(c.Query("status"))
	query.Keyword = c.Query("keyword")

	userRole, _ := c.Get("userRole")
	userID, _ := c.Get("userID")

	if userRole.(models.UserRole) == models.RoleEntrepreneur {
		uid := userID.(uint)
		query.EntrepreneurID = &uid
	} else if userRole.(models.UserRole) == models.RoleAgent {
		uid := userID.(uint)
		query.AgentID = &uid
	}

	applications, total, err := ctrl.applicationService.GetApplicationList(&query)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":  applications,
		"total": total,
		"page":  query.Page,
		"pageSize": query.PageSize,
	})
}

func (ctrl *ApplicationController) UpdateApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid application ID")
		return
	}

	var req services.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.applicationService.UpdateApplication(uint(id), &req); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *ApplicationController) UploadMaterials(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid application ID")
		return
	}

	field := c.PostForm("field")
	if field == "" {
		utils.BadRequest(c, "Field name is required")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "File is required")
		return
	}
	defer file.Close()

	if err := utils.ValidateFileSize(header.Size); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	filePath, err := utils.UploadFile(file, header.Filename, "materials")
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	if err := ctrl.applicationService.UploadMaterials(uint(id), field, filePath); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"filePath": filePath, "fileUrl": utils.GetFileURL(filePath)})
}

func (ctrl *ApplicationController) ReviewApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid application ID")
		return
	}

	var req struct {
		Approved bool   `json:"approved" binding:"required"`
		Comments string `json:"comments"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.applicationService.ReviewApplication(uint(id), req.Approved, req.Comments); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	go ctrl.notificationService.SendApplicationStatusNotification(uint(id), "processing")

	utils.Success(c, nil)
}

func (ctrl *ApplicationController) CancelApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid application ID")
		return
	}

	if err := ctrl.applicationService.CancelApplication(uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *ApplicationController) AssignAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid application ID")
		return
	}

	var req struct {
		AgentID uint `json:"agentId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.applicationService.AssignAgent(uint(id), req.AgentID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}
