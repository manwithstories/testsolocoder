package handlers

import (
	"path/filepath"
	"strconv"

	"meeting-room/internal/middleware"
	"meeting-room/internal/services"
	"meeting-room/internal/utils"

	"github.com/gin-gonic/gin"
)

type MaterialHandler struct {
	materialService *services.MaterialService
}

func NewMaterialHandler() *MaterialHandler {
	return &MaterialHandler{
		materialService: services.NewMaterialService(),
	}
}

func (h *MaterialHandler) UploadMaterial(c *gin.Context) {
	userID := middleware.GetUserID(c)

	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "预订ID错误")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请上传文件")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	fileType := "unknown"
	switch ext {
	case ".pdf":
		fileType = "pdf"
	case ".doc", ".docx":
		fileType = "word"
	case ".xls", ".xlsx":
		fileType = "excel"
	case ".ppt", ".pptx":
		fileType = "ppt"
	case ".jpg", ".jpeg", ".png":
		fileType = "image"
	}

	material, err := h.materialService.UploadMaterial(&services.UploadMaterialRequest{
		BookingID: uint(bookingID),
		UserID:    userID,
		FileName:  header.Filename,
		FileSize:  header.Size,
		FileType:  fileType,
		FileData:  file,
	})
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, material)
}

func (h *MaterialHandler) GetMaterials(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "预订ID错误")
		return
	}

	materials, err := h.materialService.GetMaterialsByBooking(uint(bookingID))
	if err != nil {
		utils.InternalError(c, "获取材料列表失败")
		return
	}

	utils.Success(c, materials)
}

func (h *MaterialHandler) DownloadMaterial(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "材料ID错误")
		return
	}

	material, filePath, err := h.materialService.DownloadMaterial(uint(id))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	c.FileAttachment(filePath, material.FileName)
}

func (h *MaterialHandler) DeleteMaterial(c *gin.Context) {
	userID := middleware.GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "材料ID错误")
		return
	}

	err = h.materialService.DeleteMaterial(uint(id), userID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}
