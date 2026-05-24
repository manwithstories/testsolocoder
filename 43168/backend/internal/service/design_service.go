package service

import (
	"errors"

	"gorm.io/gorm"

	"furniture-platform/internal/dto"
	"furniture-platform/internal/model"
	"furniture-platform/internal/repository"
)

type DesignService struct {
	repo *repository.DesignRepository
	db   *gorm.DB
}

func NewDesignService(repo *repository.DesignRepository, db *gorm.DB) *DesignService {
	return &DesignService{
		repo: repo,
		db:   db,
	}
}

func (s *DesignService) CreateProject(designerID uint, req *dto.CreateProjectRequest) (*model.DesignProject, error) {
	if req.OwnerID == 0 {
		return nil, errors.New("业主 ID 不能为空")
	}
	if req.OwnerID == designerID {
		return nil, errors.New("业主不能与设计师相同")
	}

	project := &model.DesignProject{
		DesignerID:  designerID,
		OwnerID:     req.OwnerID,
		Name:        req.Name,
		Description: req.Description,
		Status:      model.ProjectStatusDraft,
		CoverImage:  req.CoverImage,
		RoomType:    req.RoomType,
		Area:        req.Area,
		Budget:      req.Budget,
	}

	if err := s.repo.CreateProject(project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *DesignService) UpdateProject(id uint, currentUserID uint, currentRole string, req *dto.UpdateProjectRequest) (*model.DesignProject, error) {
	project, err := s.repo.GetProjectByID(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("方案不存在")
	}

	if currentRole == model.RoleDesigner && project.DesignerID != currentUserID {
		return nil, errors.New("无权编辑他人方案")
	}
	if currentRole == model.RoleOwner && project.OwnerID != currentUserID {
		return nil, errors.New("无权编辑他人方案")
	}

	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.CoverImage != "" {
		project.CoverImage = req.CoverImage
	}
	if req.RoomType != "" {
		project.RoomType = req.RoomType
	}
	if req.Area != nil {
		project.Area = *req.Area
	}
	if req.Budget != nil {
		project.Budget = *req.Budget
	}
	if req.Status != "" {
		if !model.ValidProjectStatus(req.Status) {
			return nil, errors.New("状态不合法")
		}
		project.Status = req.Status
	}

	if err := s.repo.UpdateProject(project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *DesignService) DeleteProject(id uint, currentUserID uint, currentRole string) error {
	project, err := s.repo.GetProjectByID(id)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("方案不存在")
	}
	if currentRole == model.RoleDesigner && project.DesignerID != currentUserID {
		return errors.New("无权删除他人方案")
	}
	return s.repo.DeleteProject(id)
}

func (s *DesignService) GetProject(id uint, currentUserID uint, currentRole string) (*model.DesignProject, error) {
	project, err := s.repo.GetProjectByID(id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("方案不存在")
	}
	if currentRole == model.RoleDesigner && project.DesignerID != currentUserID {
		return nil, errors.New("无权查看他人方案")
	}
	if currentRole == model.RoleOwner && project.OwnerID != currentUserID {
		return nil, errors.New("无权查看他人方案")
	}
	return project, nil
}

func (s *DesignService) ListProjects(params *dto.ProjectListRequest) ([]*model.DesignProject, int64, error) {
	p := &repository.ProjectListParams{
		Page:     params.Page,
		PageSize: params.PageSize,
		Keyword:  params.Keyword,
		Status:   params.Status,
		Role:     params.Role,
		UserID:   params.UserID,
	}
	return s.repo.ListProjects(p)
}

func (s *DesignService) UploadImage(projectID uint, currentUserID uint, req *dto.UploadImageRequest) (*model.DesignImage, error) {
	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("方案不存在")
	}
	if project.DesignerID != currentUserID {
		return nil, errors.New("仅设计师可上传图片")
	}

	img := &model.DesignImage{
		ProjectID:   projectID,
		ImageURL:    req.ImageURL,
		Description: req.Description,
		Sort:        req.Sort,
	}
	if err := s.repo.CreateImage(img); err != nil {
		return nil, err
	}
	return img, nil
}

func (s *DesignService) ListImages(projectID uint, currentUserID uint, currentRole string) ([]*model.DesignImage, error) {
	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("方案不存在")
	}
	if currentRole == model.RoleDesigner && project.DesignerID != currentUserID {
		return nil, errors.New("无权查看他人方案")
	}
	if currentRole == model.RoleOwner && project.OwnerID != currentUserID {
		return nil, errors.New("无权查看他人方案")
	}
	return s.repo.ListImages(projectID)
}

func (s *DesignService) AddComment(projectID uint, userID uint, userRole string, req *dto.AddCommentRequest) (*model.DesignComment, error) {
	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("方案不存在")
	}
	if userRole == model.RoleDesigner && project.DesignerID != userID && project.OwnerID != userID {
		return nil, errors.New("无权批注他人方案")
	}
	if userRole == model.RoleOwner && project.OwnerID != userID {
		return nil, errors.New("无权批注他人方案")
	}

	if req.ParentID != 0 {
		parent, err := s.repo.GetCommentByID(req.ParentID)
		if err != nil {
			return nil, err
		}
		if parent == nil || parent.ProjectID != projectID {
			return nil, errors.New("父级批注不存在")
		}
	}

	if !model.ValidCommentType(req.Type) {
		return nil, errors.New("批注类型不合法")
	}

	var posX, posY float64
	if req.PositionX != nil {
		posX = *req.PositionX
	}
	if req.PositionY != nil {
		posY = *req.PositionY
	}

	comment := &model.DesignComment{
		ProjectID: projectID,
		UserID:    userID,
		UserRole:  userRole,
		Content:   req.Content,
		Type:      req.Type,
		PositionX: posX,
		PositionY: posY,
		ParentID:  req.ParentID,
	}

	if err := s.repo.CreateComment(comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *DesignService) ListComments(projectID uint, currentUserID uint, currentRole string) ([]*model.DesignComment, error) {
	project, err := s.repo.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("方案不存在")
	}
	if currentRole == model.RoleDesigner && project.DesignerID != currentUserID {
		return nil, errors.New("无权查看他人方案")
	}
	if currentRole == model.RoleOwner && project.OwnerID != currentUserID {
		return nil, errors.New("无权查看他人方案")
	}
	return s.repo.ListComments(projectID)
}
