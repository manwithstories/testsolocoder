package services

import (
	"errors"
	"time"

	"recruitment-platform/internal/models"
	"recruitment-platform/internal/repository"
)

type ApplicationService struct {
	applicationRepo *repository.ApplicationRepository
	jobRepo         *repository.JobRepository
	resumeRepo      *repository.ResumeRepository
}

func NewApplicationService(
	applicationRepo *repository.ApplicationRepository,
	jobRepo *repository.JobRepository,
	resumeRepo *repository.ResumeRepository,
) *ApplicationService {
	return &ApplicationService{
		applicationRepo: applicationRepo,
		jobRepo:         jobRepo,
		resumeRepo:      resumeRepo,
	}
}

type ApplyRequest struct {
	JobID       uint   `json:"job_id" binding:"required"`
	ResumeID    uint   `json:"resume_id" binding:"required"`
	CoverLetter string `json:"cover_letter"`
}

type UpdateStatusRequest struct {
	Status       string `json:"status" binding:"required"`
	ChangeReason string `json:"change_reason"`
}

type BulkUpdateRequest struct {
	IDs    []uint `json:"ids" binding:"required"`
	Status string `json:"status" binding:"required"`
	Reason string `json:"reason"`
}

func (s *ApplicationService) Apply(applicantID uint, req *ApplyRequest) (*models.Application, error) {
	job, err := s.jobRepo.FindByID(req.JobID)
	if err != nil {
		return nil, errors.New("职位不存在")
	}

	if job.Status != models.JobStatusOpen {
		return nil, errors.New("该职位暂不接受投递")
	}

	if job.Deadline != nil && job.Deadline.Before(time.Now()) {
		return nil, errors.New("该职位已截止")
	}

	existing, _ := s.applicationRepo.FindByJobAndApplicant(req.JobID, applicantID)
	if existing != nil {
		return nil, errors.New("您已投递过该职位")
	}

	resume, err := s.resumeRepo.FindByID(req.ResumeID)
	if err != nil || resume.UserID != applicantID {
		return nil, errors.New("简历不存在或无权限使用")
	}

	application := &models.Application{
		JobID:       req.JobID,
		ApplicantID: applicantID,
		ResumeID:    req.ResumeID,
		Status:      models.ApplicationStatusPending,
		CoverLetter: req.CoverLetter,
		AppliedAt:   time.Now(),
	}

	if err := s.applicationRepo.Create(application); err != nil {
		return nil, errors.New("投递失败")
	}

	s.jobRepo.IncrementApplyCount(req.JobID)

	return application, nil
}

func (s *ApplicationService) GetApplication(id, userID uint, role models.UserRole) (*models.Application, error) {
	application, err := s.applicationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("投递记录不存在")
	}

	if role == models.RoleApplicant && application.ApplicantID != userID {
		return nil, errors.New("无权限查看此投递记录")
	}

	if role == models.RoleCompany {
		if application.Job == nil || application.Job.CompanyID == 0 {
			return nil, errors.New("无权限查看此投递记录")
		}
	}

	return application, nil
}

func (s *ApplicationService) ListByApplicant(applicantID uint, page, pageSize int, status string) ([]models.Application, int64, error) {
	return s.applicationRepo.ListByApplicant(applicantID, page, pageSize, models.ApplicationStatus(status))
}

func (s *ApplicationService) ListByJob(jobID, companyID uint, page, pageSize int, status string) ([]models.Application, int64, error) {
	job, err := s.jobRepo.FindByID(jobID)
	if err != nil {
		return nil, 0, errors.New("职位不存在")
	}

	if job.CompanyID != companyID {
		return nil, 0, errors.New("无权限查看此职位的投递")
	}

	return s.applicationRepo.ListByJob(jobID, page, pageSize, models.ApplicationStatus(status))
}

func (s *ApplicationService) ListByCompany(companyID uint, page, pageSize int, status, keyword string) ([]models.Application, int64, error) {
	return s.applicationRepo.ListByCompany(companyID, page, pageSize, models.ApplicationStatus(status), keyword)
}

func (s *ApplicationService) UpdateStatus(id, userID uint, role models.UserRole, req *UpdateStatusRequest) (*models.Application, error) {
	application, err := s.applicationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("投递记录不存在")
	}

	oldStatus := application.Status
	newStatus := models.ApplicationStatus(req.Status)

	validTransitions := map[models.ApplicationStatus][]models.ApplicationStatus{
		models.ApplicationStatusPending:    {models.ApplicationStatusViewed, models.ApplicationStatusRejected, models.ApplicationStatusWithdrawn},
		models.ApplicationStatusViewed:     {models.ApplicationStatusInterested, models.ApplicationStatusRejected},
		models.ApplicationStatusInterested: {models.ApplicationStatusInterview, models.ApplicationStatusRejected},
		models.ApplicationStatusInterview:  {models.ApplicationStatusAccepted, models.ApplicationStatusRejected},
		models.ApplicationStatusAccepted:   {},
		models.ApplicationStatusRejected:   {},
		models.ApplicationStatusWithdrawn:  {},
	}

	validNext, exists := validTransitions[oldStatus]
	if !exists {
		return nil, errors.New("当前状态不允许变更")
	}

	canTransition := false
	for _, s := range validNext {
		if s == newStatus {
			canTransition = true
			break
		}
	}

	if !canTransition {
		return nil, errors.New("不允许从当前状态变更到目标状态")
	}

	if role == models.RoleApplicant && newStatus != models.ApplicationStatusWithdrawn {
		return nil, errors.New("求职者只能撤回投递")
	}

	if role == models.RoleApplicant && application.ApplicantID != userID {
		return nil, errors.New("无权限操作此投递记录")
	}

	application.Status = newStatus
	application.LastUpdateAt = time.Now()

	if err := s.applicationRepo.Update(application); err != nil {
		return nil, errors.New("更新状态失败")
	}

	history := &models.ApplicationHistory{
		ApplicationID: application.ID,
		OldStatus:     oldStatus,
		NewStatus:     newStatus,
		ChangedBy:     userID,
		ChangeReason:  req.ChangeReason,
	}
	s.applicationRepo.AddHistory(history)

	return application, nil
}

func (s *ApplicationService) BulkUpdateStatus(ids []uint, status, reason string, changedBy uint) (int, error) {
	newStatus := models.ApplicationStatus(status)

	if err := s.applicationRepo.BulkUpdateStatus(ids, newStatus); err != nil {
		return 0, errors.New("批量更新失败")
	}

	for _, id := range ids {
		history := &models.ApplicationHistory{
			ApplicationID: id,
			NewStatus:     newStatus,
			ChangedBy:     changedBy,
			ChangeReason:  reason,
		}
		s.applicationRepo.AddHistory(history)
	}

	return len(ids), nil
}

func (s *ApplicationService) Withdraw(id, applicantID uint) error {
	application, err := s.applicationRepo.FindByID(id)
	if err != nil {
		return errors.New("投递记录不存在")
	}

	if application.ApplicantID != applicantID {
		return errors.New("无权限操作此投递记录")
	}

	if application.Status != models.ApplicationStatusPending &&
		application.Status != models.ApplicationStatusViewed {
		return errors.New("当前状态不允许撤回")
	}

	oldStatus := application.Status
	application.Status = models.ApplicationStatusWithdrawn
	application.LastUpdateAt = time.Now()

	if err := s.applicationRepo.Update(application); err != nil {
		return errors.New("撤回失败")
	}

	history := &models.ApplicationHistory{
		ApplicationID: application.ID,
		OldStatus:     oldStatus,
		NewStatus:     models.ApplicationStatusWithdrawn,
		ChangedBy:     applicantID,
		ChangeReason:  "求职者主动撤回",
	}
	s.applicationRepo.AddHistory(history)

	return nil
}

func (s *ApplicationService) GetHistory(applicationID uint) ([]models.ApplicationHistory, error) {
	return s.applicationRepo.GetHistory(applicationID)
}

func (s *ApplicationService) GetStatusCountByJob(jobID uint) (map[string]int64, error) {
	return s.applicationRepo.GetStatusCountByJob(jobID)
}
