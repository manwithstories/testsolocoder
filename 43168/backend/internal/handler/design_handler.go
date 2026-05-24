package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/service"
	"furniture-platform/pkg/response"
)

type DesignHandler struct {
	service *service.DesignService
}

func NewDesignHandler(svc *service.DesignService) *DesignHandler {
	return &DesignHandler{service: svc}
}

func currentUserRole(c *gin.Context) (string, bool) {
	val, exists := c.Get("role")
	if !exists {
		return "", false
	}
	role, ok := val.(string)
	if !ok {
		return "", false
	}
	return role, true
}

func toProjectResponse(p *model.DesignProject) dto.ProjectResponse {
	return dto.ProjectResponse{
		ID:          p.ID,
		DesignerID:  p.DesignerID,
		OwnerID:     p.OwnerID,
		Name:        p.Name,
		Description: p.Description,
		Status:      p.Status,
		CoverImage:  p.CoverImage,
		RoomType:    p.RoomType,
		Area:        p.Area,
		Budget:      p.Budget,
		CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toImageResponse(img *model.DesignImage) dto.ImageResponse {
	return dto.ImageResponse{
		ID:          img.ID,
		ProjectID:   img.ProjectID,
		ImageURL:    img.ImageURL,
		Description: img.Description,
		Sort:        img.Sort,
		CreatedAt:   img.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toCommentResponse(c *model.DesignComment) dto.CommentResponse {
	return dto.CommentResponse{
		ID:        c.ID,
		ProjectID: c.ProjectID,
		UserID:    c.UserID,
		UserRole:  c.UserRole,
		Content:   c.Content,
		Type:      c.Type,
		PositionX: c.PositionX,
		PositionY: c.PositionY,
		ParentID:  c.ParentID,
		CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func parseID(c *gin.Context, name string) (uint, error) {
	idStr := c.Param(name)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (h *DesignHandler) CreateProject(c *gin.Context) {
	uid, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}
	role, ok := currentUserRole(c)
	if !ok || role != model.RoleDesigner {
		response.Forbidden(c, "仅设计师可创建设计方案")
		return
	}

	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	project, err := h.service.CreateProject(uid, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toProjectResponse(project))
}

func (h *DesignHandler) UpdateProject(c *gin.Context) {
	uid, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}
	role, ok := currentUserRole(c)
	if !ok {
		response.Forbidden(c, "权限不足")
		return
	}

	id, err := parseID(c, "id")
	if err != nil {
		response.BadRequest(c, "方案 ID 格式错误")
		return
	}

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	project, err := h.service.UpdateProject(id, uid, role, &req)
	if err != nil {
		if err.Error() == "方案不存在" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "无权编辑他人方案" {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toProjectResponse(project))
}

func (h *DesignHandler) DeleteProject(c *gin.Context) {
	uid, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}
	role, ok := currentUserRole(c)
	if !ok {
		response.Forbidden(c, "权限不足")
		return
	}

	id, err := parseID(c, "id")
	if err != nil {
		response.BadRequest(c, "方案 ID 格式错误")
		return
	}

	if err := h.service.DeleteProject(id, uid, role); err != nil {
		if err.Error() == "方案不存在" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "无权删除他人方案" {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *DesignHandler) GetProject(c *gin.Context) {
	uid, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}
	role, ok := currentUserRole(c)
	if !ok {
		response.Forbidden(c, "权限不足")
		return
	}

	id, err := parseID(c, "id")
	if err != nil {
		response.BadRequest(c, "方案 ID 格式错误")
		return
	}

	project, err := h.service.GetProject(id, uid, role)
	if err != nil {
		if err.Error() == "方案不存在" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "无权查看他人方案" {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	images, err := h.service.ListImages(id, uid, role)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	imgList := make([]dto.ImageResponse, 0, len(images))
	for _, img := range images {
		imgList = append(imgList, toImageResponse(img))
	}

	response.Success(c, dto.ProjectDetailResponse{
		ProjectResponse: toProjectResponse(project),
		Images:          imgList,
	})
}

func (h *DesignHandler) ListProjects(c *gin.Context) {
	uid, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}
	role, ok := currentUserRole(c)
	if !ok {
		response.Forbidden(c, "权限不足")
		return
	}

	var req dto.ProjectListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	req.UserID = uid
	req.Role = role

	projects, total, err := h.service.ListProjects(&req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	list := make([]dto.ProjectResponse, 0, len(projects))
	for _, p := range projects {
		list = append(list, toProjectResponse(p))
	}

	response.Success(c, dto.ProjectListResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	})
}

func (h *DesignHandler) UploadImage(c *gin.Context) {
	uid, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	id, err := parseID(c, "id")
	if err != nil {
		response.BadRequest(c, "方案 ID 格式错误")
		return
	}

	var req dto.UploadImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	img, err := h.service.UploadImage(id, uid, &req)
	if err != nil {
		if err.Error() == "方案不存在" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "仅设计师可上传图片" {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toImageResponse(img))
}

func (h *DesignHandler) AddComment(c *gin.Context) {
	uid, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}
	role, ok := currentUserRole(c)
	if !ok {
		response.Forbidden(c, "权限不足")
		return
	}

	id, err := parseID(c, "id")
	if err != nil {
		response.BadRequest(c, "方案 ID 格式错误")
		return
	}

	var req dto.AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	comment, err := h.service.AddComment(id, uid, role, &req)
	if err != nil {
		if err.Error() == "方案不存在" || err.Error() == "父级批注不存在" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "无权批注他人方案" {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, toCommentResponse(comment))
}

func (h *DesignHandler) GetComments(c *gin.Context) {
	uid, ok := currentUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}
	role, ok := currentUserRole(c)
	if !ok {
		response.Forbidden(c, "权限不足")
		return
	}

	id, err := parseID(c, "id")
	if err != nil {
		response.BadRequest(c, "方案 ID 格式错误")
		return
	}

	comments, err := h.service.ListComments(id, uid, role)
	if err != nil {
		if err.Error() == "方案不存在" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "无权查看他人方案" {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	list := make([]dto.CommentResponse, 0, len(comments))
	for _, cmt := range comments {
		list = append(list, toCommentResponse(cmt))
	}
	response.Success(c, dto.CommentListResponse{List: list})
}
