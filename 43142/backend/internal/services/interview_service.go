package services

import (
	"errors"
	"time"

	"recruitment-platform/internal/models"
	"recruitment-platform/internal/repository"
	"recruitment-platform/internal/utils"
)

type InterviewService struct {
	interviewRepo   *repository.InterviewRepository
	applicationRepo *repository.ApplicationRepository
	jobRepo         *repository.JobRepository
	emailService    *utils.EmailService
}

func NewInterviewService(
	interviewRepo *repository.InterviewRepository,
	applicationRepo *repository.ApplicationRepository,
	jobRepo *repository.JobRepository,
	emailService *utils.EmailService,
) *InterviewService {
	return &InterviewService{
		interviewRepo:   interviewRepo,
		applicationRepo: applicationRepo,
		jobRepo:         jobRepo,
		emailService:    emailService,
	}
}

type ScheduleInterviewRequest struct {
	ApplicationID uint      `json:"application_id" binding:"required"`
	Interviewer   string    `json:"interviewer"`
	InterviewType string    `json:"interview_type"`
	Location      string    `json:"location"`
	MeetingLink   string    `json:"meeting_link"`
	ScheduledAt   time.Time `json:"scheduled_at" binding:"required"`
	Duration      int       `json:"duration"`
	Notes         string    `json:"notes"`
}

type UpdateInterviewRequest struct {
	Interviewer *string    `json:"interviewer"`
	Location    *string    `json:"location"`
	MeetingLink *string    `json:"meeting_link"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	Duration    *int       `json:"duration"`
	Notes       *string    `json:"notes"`
	Feedback    *string    `json:"feedback"`
	Rating      *int       `json:"rating"`
}

func (s *InterviewService) ScheduleInterview(companyID uint, req *ScheduleInterviewRequest) (*models.Interview, error) {
	application, err := s.applicationRepo.FindByID(req.ApplicationID)
	if err != nil {
		return nil, errors.New("投递记录不存在")
	}

	if application.Job == nil || application.Job.CompanyID != companyID {
		return nil, errors.New("无权限操作此投递记录")
	}

	if application.Status != models.ApplicationStatusInterested &&
		application.Status != models.ApplicationStatusInterview {
		return nil, errors.New("当前投递状态不允许安排面试")
	}

	if req.ScheduledAt.Before(time.Now()) {
		return nil, errors.New("面试时间不能早于当前时间")
	}

	duration := req.Duration
	if duration <= 0 {
		duration = 60
	}

	hasConflict, err := s.interviewRepo.CheckTimeConflict(application.ApplicantID, req.ScheduledAt, duration)
	if err != nil {
		return nil, errors.New("检查时间冲突失败")
	}
	if hasConflict {
		return nil, errors.New("该求职者在此时间段已有面试安排")
	}

	interview := &models.Interview{
		ApplicationID: req.ApplicationID,
		JobID:         application.JobID,
		ApplicantID:   application.ApplicantID,
		Interviewer:   req.Interviewer,
		InterviewType: req.InterviewType,
		Location:      req.Location,
		MeetingLink:   req.MeetingLink,
		ScheduledAt:   req.ScheduledAt,
		Duration:      duration,
		Status:        models.InterviewStatusPending,
		Notes:         req.Notes,
	}

	if err := s.interviewRepo.Create(interview); err != nil {
		return nil, errors.New("创建面试失败")
	}

	application.Status = models.ApplicationStatusInterview
	application.LastUpdateAt = time.Now()
	s.applicationRepo.Update(application)

	if application.Applicant != nil && application.Applicant.Profile != nil {
		s.emailService.SendInterviewInvitation(
			application.Applicant.Email,
			application.Applicant.Profile.FullName,
			application.Job.Title,
			req.Interviewer,
			req.ScheduledAt,
			req.Location,
		)
	}

	return interview, nil
}

func (s *InterviewService) GetInterview(id, userID uint, role models.UserRole) (*models.Interview, error) {
	interview, err := s.interviewRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("面试记录不存在")
	}

	if role == models.RoleApplicant && interview.ApplicantID != userID {
		return nil, errors.New("无权限查看此面试记录")
	}

	return interview, nil
}

func (s *InterviewService) ListByCompany(companyID uint, page, pageSize int, status string) ([]models.Interview, int64, error) {
	return s.interviewRepo.ListByCompany(companyID, page, pageSize, models.InterviewStatus(status))
}

func (s *InterviewService) ListByApplicant(applicantID uint, page, pageSize int, status string) ([]models.Interview, int64, error) {
	return s.interviewRepo.ListByApplicant(applicantID, page, pageSize, models.InterviewStatus(status))
}

func (s *InterviewService) UpdateInterview(id, companyID uint, req *UpdateInterviewRequest) (*models.Interview, error) {
	interview, err := s.interviewRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("面试记录不存在")
	}

	if interview.Application == nil || interview.Application.Job == nil ||
		interview.Application.Job.CompanyID != companyID {
		return nil, errors.New("无权限操作此面试记录")
	}

	if req.Interviewer != nil {
		interview.Interviewer = *req.Interviewer
	}
	if req.Location != nil {
		interview.Location = *req.Location
	}
	if req.MeetingLink != nil {
		interview.MeetingLink = *req.MeetingLink
	}
	if req.ScheduledAt != nil {
		if req.ScheduledAt.Before(time.Now()) {
			return nil, errors.New("面试时间不能早于当前时间")
		}

		hasConflict, _ := s.interviewRepo.CheckTimeConflict(interview.ApplicantID, *req.ScheduledAt, interview.Duration)
		if hasConflict {
			return nil, errors.New("该求职者在此时间段已有面试安排")
		}
		interview.ScheduledAt = *req.ScheduledAt
	}
	if req.Duration != nil && *req.Duration > 0 {
		interview.Duration = *req.Duration
	}
	if req.Notes != nil {
		interview.Notes = *req.Notes
	}
	if req.Feedback != nil {
		interview.Feedback = *req.Feedback
	}
	if req.Rating != nil {
		interview.Rating = *req.Rating
	}

	if err := s.interviewRepo.Update(interview); err != nil {
		return nil, errors.New("更新面试失败")
	}

	return interview, nil
}

func (s *InterviewService) AcceptInterview(id, applicantID uint) error {
	interview, err := s.interviewRepo.FindByID(id)
	if err != nil {
		return errors.New("面试记录不存在")
	}

	if interview.ApplicantID != applicantID {
		return errors.New("无权限操作此面试记录")
	}

	if interview.Status != models.InterviewStatusPending {
		return errors.New("当前状态不允许接受面试")
	}

	return s.interviewRepo.UpdateStatus(id, models.InterviewStatusAccepted)
}

func (s *InterviewService) RejectInterview(id, applicantID uint) error {
	interview, err := s.interviewRepo.FindByID(id)
	if err != nil {
		return errors.New("面试记录不存在")
	}

	if interview.ApplicantID != applicantID {
		return errors.New("无权限操作此面试记录")
	}

	if interview.Status != models.InterviewStatusPending {
		return errors.New("当前状态不允许拒绝面试")
	}

	return s.interviewRepo.UpdateStatus(id, models.InterviewStatusRejected)
}

func (s *InterviewService) CompleteInterview(id, companyID uint, feedback string, rating int) error {
	interview, err := s.interviewRepo.FindByID(id)
	if err != nil {
		return errors.New("面试记录不存在")
	}

	if interview.Application == nil || interview.Application.Job == nil ||
		interview.Application.Job.CompanyID != companyID {
		return errors.New("无权限操作此面试记录")
	}

	interview.Status = models.InterviewStatusCompleted
	interview.Feedback = feedback
	interview.Rating = rating

	return s.interviewRepo.Update(interview)
}

func (s *InterviewService) CancelInterview(id, companyID uint) error {
	interview, err := s.interviewRepo.FindByID(id)
	if err != nil {
		return errors.New("面试记录不存在")
	}

	if interview.Application == nil || interview.Application.Job == nil ||
		interview.Application.Job.CompanyID != companyID {
		return errors.New("无权限操作此面试记录")
	}

	return s.interviewRepo.UpdateStatus(id, models.InterviewStatusCancelled)
}
