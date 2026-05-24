package services

import (
	"errors"
	"mime/multipart"

	"recruitment-platform/internal/config"
	"recruitment-platform/internal/models"
	"recruitment-platform/internal/repository"
	"recruitment-platform/internal/utils"
)

type ResumeService struct {
	resumeRepo    *repository.ResumeRepository
	fileValidator *utils.FileValidator
}

func NewResumeService(resumeRepo *repository.ResumeRepository) *ResumeService {
	cfg := config.AppConfig
	fileValidator := utils.NewFileValidator(
		cfg.Upload.MaxSizeMB,
		cfg.Upload.AllowedTypes,
		cfg.Upload.ResumePath,
	)
	return &ResumeService{
		resumeRepo:    resumeRepo,
		fileValidator: fileValidator,
	}
}

type CreateResumeRequest struct {
	Title      string `json:"title" binding:"required"`
	FullName   string `json:"full_name" binding:"required"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Location   string `json:"location"`
	Education  string `json:"education"`
	Experience string `json:"experience"`
	Skills     string `json:"skills"`
	Summary    string `json:"summary"`
	Projects   string `json:"projects"`
	IsDefault  bool   `json:"is_default"`
}

func (s *ResumeService) CreateResume(userID uint, req *CreateResumeRequest) (*models.Resume, error) {
	resume := &models.Resume{
		UserID:     userID,
		Title:      req.Title,
		FullName:   req.FullName,
		Email:      req.Email,
		Phone:      req.Phone,
		Location:   req.Location,
		Education:  req.Education,
		Experience: req.Experience,
		Skills:     req.Skills,
		Summary:    req.Summary,
		Projects:   req.Projects,
		IsDefault:  req.IsDefault,
	}

	if req.IsDefault {
		resume.IsDefault = true
	}

	if err := s.resumeRepo.Create(resume); err != nil {
		return nil, errors.New("创建简历失败")
	}

	if req.IsDefault {
		s.resumeRepo.SetDefault(userID, resume.ID)
	}

	return resume, nil
}

func (s *ResumeService) GetResume(id, userID uint) (*models.Resume, error) {
	resume, err := s.resumeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("简历不存在")
	}

	if resume.UserID != userID {
		return nil, errors.New("无权限查看此简历")
	}

	return resume, nil
}

func (s *ResumeService) UpdateResume(id, userID uint, req *CreateResumeRequest) (*models.Resume, error) {
	resume, err := s.resumeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("简历不存在")
	}

	if resume.UserID != userID {
		return nil, errors.New("无权限修改此简历")
	}

	resume.Title = req.Title
	resume.FullName = req.FullName
	resume.Email = req.Email
	resume.Phone = req.Phone
	resume.Location = req.Location
	resume.Education = req.Education
	resume.Experience = req.Experience
	resume.Skills = req.Skills
	resume.Summary = req.Summary
	resume.Projects = req.Projects

	if req.IsDefault {
		s.resumeRepo.SetDefault(userID, id)
	}

	if err := s.resumeRepo.Update(resume); err != nil {
		return nil, errors.New("更新简历失败")
	}

	return resume, nil
}

func (s *ResumeService) DeleteResume(id, userID uint) error {
	resume, err := s.resumeRepo.FindByID(id)
	if err != nil {
		return errors.New("简历不存在")
	}

	if resume.UserID != userID {
		return errors.New("无权限删除此简历")
	}

	if resume.FileURL != "" {
		utils.DeleteFile(resume.FileURL)
	}

	return s.resumeRepo.Delete(id)
}

func (s *ResumeService) ListResumes(userID uint) ([]models.Resume, error) {
	return s.resumeRepo.ListByUserID(userID)
}

func (s *ResumeService) SetDefaultResume(userID, resumeID uint) error {
	return s.resumeRepo.SetDefault(userID, resumeID)
}

func (s *ResumeService) UploadResumeFile(userID, resumeID uint, file *multipart.FileHeader) (*models.Resume, error) {
	resume, err := s.resumeRepo.FindByID(resumeID)
	if err != nil {
		return nil, errors.New("简历不存在")
	}

	if resume.UserID != userID {
		return nil, errors.New("无权限操作此简历")
	}

	if err := s.fileValidator.ValidateFile(file); err != nil {
		return nil, err
	}

	if resume.FileURL != "" {
		utils.DeleteFile(resume.FileURL)
	}

	fullPath, relativePath, err := s.fileValidator.SaveFile(file, "resumes")
	if err != nil {
		return nil, errors.New("文件保存失败")
	}

	resume.FileURL = relativePath
	resume.FileType = file.Header.Get("Content-Type")
	resume.FileSize = file.Size

	if err := s.resumeRepo.Update(resume); err != nil {
		utils.DeleteFile(fullPath)
		return nil, errors.New("更新简历失败")
	}

	return resume, nil
}

func (s *ResumeService) GetDefaultResume(userID uint) (*models.Resume, error) {
	return s.resumeRepo.GetDefaultResume(userID)
}

func (s *ResumeService) SearchResumes(keyword string, skills []string, page, pageSize int) ([]models.Resume, int64, error) {
	return s.resumeRepo.Search(keyword, skills, page, pageSize)
}
