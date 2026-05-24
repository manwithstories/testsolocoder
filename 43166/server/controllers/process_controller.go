package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"business-registration-platform/services"
	"business-registration-platform/utils"
)

type ProcessController struct {
	processService      *services.ProcessService
	notificationService *services.NotificationService
}

func NewProcessController() *ProcessController {
	return &ProcessController{
		processService:      services.NewProcessService(),
		notificationService: services.NewNotificationService(),
	}
}

func (ctrl *ProcessController) GetProcessSteps(c *gin.Context) {
	applicationID, err := strconv.ParseUint(c.Param("applicationId"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid application ID")
		return
	}

	steps, err := ctrl.processService.GetProcessSteps(uint(applicationID))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, steps)
}

func (ctrl *ProcessController) UpdateProcessStep(c *gin.Context) {
	stepID, err := strconv.ParseUint(c.Param("stepId"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid step ID")
		return
	}

	userID, _ := c.Get("userID")

	var req services.UpdateStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.processService.UpdateProcessStep(uint(stepID), userID.(uint), &req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *ProcessController) StartStep(c *gin.Context) {
	stepID, err := strconv.ParseUint(c.Param("stepId"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid step ID")
		return
	}

	userID, _ := c.Get("userID")

	if err := ctrl.processService.StartStep(uint(stepID), userID.(uint)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *ProcessController) CompleteStep(c *gin.Context) {
	stepID, err := strconv.ParseUint(c.Param("stepId"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid step ID")
		return
	}

	userID, _ := c.Get("userID")

	certificateFile := c.PostForm("certificateFile")
	remark := c.PostForm("remark")

	if err := ctrl.processService.CompleteStep(uint(stepID), userID.(uint), certificateFile, remark); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *ProcessController) SkipStep(c *gin.Context) {
	stepID, err := strconv.ParseUint(c.Param("stepId"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid step ID")
		return
	}

	userID, _ := c.Get("userID")

	var req struct {
		Remark string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err.Error())
		return
	}

	if err := ctrl.processService.SkipStep(uint(stepID), userID.(uint), req.Remark); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *ProcessController) UploadCertificate(c *gin.Context) {
	stepID, err := strconv.ParseUint(c.Param("stepId"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "Invalid step ID")
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

	filePath, err := utils.UploadFile(file, header.Filename, "certificates")
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"stepId":       stepID,
		"filePath":     filePath,
		"fileUrl":      utils.GetFileURL(filePath),
	})
}
