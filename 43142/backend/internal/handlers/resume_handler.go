package handlers

import (
	"github.com/gin-gonic/gin"

	"recruitment-platform/internal/services"
	"recruitment-platform/internal/utils"
	"recruitment-platform/pkg/logger"
)

type ResumeHandler struct {
	resumeService *services.ResumeService
}

func NewResumeHandler(resumeService *services.ResumeService) *ResumeHandler {
	return &ResumeHandler{resumeService: resumeService}
}

func (h *ResumeHandler) CreateResume(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req services.CreateResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resume, err := h.resumeService.CreateResume(userID, &req)
	if err != nil {
		logger.Error("创建简历失败: %v", err)
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, resume)
}

func (h *ResumeHandler) GetResume(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的简历ID")
		return
	}

	userID := c.GetUint("user_id")

	resume, err := h.resumeService.GetResume(id, userID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, resume)
}

func (h *ResumeHandler) UpdateResume(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的简历ID")
		return
	}

	userID := c.GetUint("user_id")

	var req services.CreateResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resume, err := h.resumeService.UpdateResume(id, userID, &req)
	if err != nil {
		logger.Error("更新简历失败: %v", err)
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, resume)
}

func (h *ResumeHandler) DeleteResume(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的简历ID")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.resumeService.DeleteResume(id, userID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *ResumeHandler) ListResumes(c *gin.Context) {
	userID := c.GetUint("user_id")

	resumes, err := h.resumeService.ListResumes(userID)
	if err != nil {
		logger.Error("获取简历列表失败: %v", err)
		utils.InternalError(c, "获取简历列表失败")
		return
	}

	utils.Success(c, resumes)
}

func (h *ResumeHandler) SetDefaultResume(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的简历ID")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.resumeService.SetDefaultResume(userID, id); err != nil {
		utils.InternalError(c, "设置默认简历失败")
		return
	}

	utils.SuccessWithMessage(c, "设置成功", nil)
}

func (h *ResumeHandler) UploadResumeFile(c *gin.Context) {
	id := parseUintParam(c.Param("id"))
	if id == 0 {
		utils.BadRequest(c, "无效的简历ID")
		return
	}

	userID := c.GetUint("user_id")

	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请上传文件")
		return
	}

	resume, err := h.resumeService.UploadResumeFile(userID, id, file)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, resume)
}

func (h *ResumeHandler) GetDefaultResume(c *gin.Context) {
	userID := c.GetUint("user_id")

	resume, err := h.resumeService.GetDefaultResume(userID)
	if err != nil {
		utils.NotFound(c, "默认简历不存在")
		return
	}

	utils.Success(c, resume)
}

func (h *ResumeHandler) SearchResumes(c *gin.Context) {
	page, pageSize := utils.GetPagination(c)
	keyword := c.Query("keyword")

	skills := c.QueryArray("skills")

	resumes, total, err := h.resumeService.SearchResumes(keyword, skills, page, pageSize)
	if err != nil {
		logger.Error("搜索简历失败: %v", err)
		utils.InternalError(c, "搜索简历失败")
		return
	}

	utils.Paginated(c, resumes, page, pageSize, total)
}
