package services

import (
	"errors"
	"fmt"
	"time"

	"recruitment-platform/internal/models"
	"recruitment-platform/internal/repository"
)

type JobService struct {
	jobRepo *repository.JobRepository
}

func NewJobService(jobRepo *repository.JobRepository) *JobService {
	return &JobService{jobRepo: jobRepo}
}

type CreateJobRequest struct {
	Title        string     `json:"title" binding:"required"`
	Description  string     `json:"description" binding:"required"`
	SalaryMin    int        `json:"salary_min"`
	SalaryMax    int        `json:"salary_max"`
	SalaryType   string     `json:"salary_type"`
	Location     string     `json:"location" binding:"required"`
	JobType      string     `json:"job_type"`
	Experience   string     `json:"experience"`
	Education    string     `json:"education"`
	Skills       string     `json:"skills"`
	Requirements string     `json:"requirements"`
	Benefits     string     `json:"benefits"`
	Deadline     *time.Time `json:"deadline"`
}

type UpdateJobRequest struct {
	Title        *string     `json:"title"`
	Description  *string     `json:"description"`
	SalaryMin    *int        `json:"salary_min"`
	SalaryMax    *int        `json:"salary_max"`
	SalaryType   *string     `json:"salary_type"`
	Location     *string     `json:"location"`
	JobType      *string     `json:"job_type"`
	Experience   *string     `json:"experience"`
	Education    *string     `json:"education"`
	Skills       *string     `json:"skills"`
	Requirements *string     `json:"requirements"`
	Benefits     *string     `json:"benefits"`
	Deadline     *time.Time  `json:"deadline"`
	Status       *string     `json:"status"`
}

type BulkJobItem struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	SalaryMin    int    `json:"salary_min"`
	SalaryMax    int    `json:"salary_max"`
	Location     string `json:"location"`
	JobType      string `json:"job_type"`
	Experience   string `json:"experience"`
	Skills       string `json:"skills"`
	Requirements string `json:"requirements"`
}

func (s *JobService) CreateJob(companyID uint, req *CreateJobRequest) (*models.Job, error) {
	if req.SalaryMax > 0 && req.SalaryMin > req.SalaryMax {
		return nil, errors.New("最低薪资不能高于最高薪资")
	}

	if req.Deadline != nil && req.Deadline.Before(time.Now()) {
		return nil, errors.New("截止日期不能早于当前时间")
	}

	job := &models.Job{
		CompanyID:    companyID,
		Title:        req.Title,
		Description:  req.Description,
		SalaryMin:    req.SalaryMin,
		SalaryMax:    req.SalaryMax,
		SalaryType:   req.SalaryType,
		Location:     req.Location,
		JobType:      req.JobType,
		Experience:   req.Experience,
		Education:    req.Education,
		Skills:       req.Skills,
		Requirements: req.Requirements,
		Benefits:     req.Benefits,
		Deadline:     req.Deadline,
		Status:       models.JobStatusDraft,
	}

	if err := s.jobRepo.Create(job); err != nil {
		return nil, errors.New("创建职位失败")
	}

	return job, nil
}

func (s *JobService) GetJob(id uint) (*models.Job, error) {
	job, err := s.jobRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("职位不存在")
	}
	return job, nil
}

func (s *JobService) GetJobWithView(id uint, ipAddress string, userID *uint) (*models.Job, error) {
	job, err := s.jobRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("职位不存在")
	}

	s.jobRepo.IncrementViewCount(id)
	s.jobRepo.LogView(id, userID, ipAddress)

	return job, nil
}

func (s *JobService) UpdateJob(id, companyID uint, req *UpdateJobRequest) (*models.Job, error) {
	job, err := s.jobRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("职位不存在")
	}

	if job.CompanyID != companyID {
		return nil, errors.New("无权限修改此职位")
	}

	if req.Title != nil {
		job.Title = *req.Title
	}
	if req.Description != nil {
		job.Description = *req.Description
	}
	if req.SalaryMin != nil {
		job.SalaryMin = *req.SalaryMin
	}
	if req.SalaryMax != nil {
		job.SalaryMax = *req.SalaryMax
	}
	if req.SalaryType != nil {
		job.SalaryType = *req.SalaryType
	}
	if req.Location != nil {
		job.Location = *req.Location
	}
	if req.JobType != nil {
		job.JobType = *req.JobType
	}
	if req.Experience != nil {
		job.Experience = *req.Experience
	}
	if req.Education != nil {
		job.Education = *req.Education
	}
	if req.Skills != nil {
		job.Skills = *req.Skills
	}
	if req.Requirements != nil {
		job.Requirements = *req.Requirements
	}
	if req.Benefits != nil {
		job.Benefits = *req.Benefits
	}
	if req.Deadline != nil {
		job.Deadline = req.Deadline
	}
	if req.Status != nil {
		job.Status = models.JobStatus(*req.Status)
	}

	if job.SalaryMax > 0 && job.SalaryMin > job.SalaryMax {
		return nil, errors.New("最低薪资不能高于最高薪资")
	}

	if err := s.jobRepo.Update(job); err != nil {
		return nil, errors.New("更新职位失败")
	}

	return job, nil
}

func (s *JobService) DeleteJob(id, companyID uint) error {
	job, err := s.jobRepo.FindByID(id)
	if err != nil {
		return errors.New("职位不存在")
	}

	if job.CompanyID != companyID {
		return errors.New("无权限删除此职位")
	}

	return s.jobRepo.Delete(id)
}

func (s *JobService) ListJobs(page, pageSize int, companyID *uint, status, keyword string) ([]models.Job, int64, error) {
	return s.jobRepo.List(page, pageSize, companyID, models.JobStatus(status), keyword)
}

func (s *JobService) PublishJob(id, companyID uint) error {
	job, err := s.jobRepo.FindByID(id)
	if err != nil {
		return errors.New("职位不存在")
	}

	if job.CompanyID != companyID {
		return errors.New("无权限操作此职位")
	}

	if job.Title == "" || job.Description == "" || job.Location == "" {
		return errors.New("职位信息不完整，无法发布")
	}

	return s.jobRepo.UpdateStatus(id, models.JobStatusOpen)
}

func (s *JobService) CloseJob(id, companyID uint) error {
	job, err := s.jobRepo.FindByID(id)
	if err != nil {
		return errors.New("职位不存在")
	}

	if job.CompanyID != companyID {
		return errors.New("无权限操作此职位")
	}

	return s.jobRepo.UpdateStatus(id, models.JobStatusClosed)
}

func (s *JobService) BulkImport(companyID uint, items []BulkJobItem) (int, []string) {
	var jobs []models.Job
	var errors []string

	for i, item := range items {
		if item.Title == "" {
			errors = append(errors, fmt.Sprintf("第%d条数据：标题不能为空", i+1))
			continue
		}
		if item.Location == "" {
			errors = append(errors, fmt.Sprintf("第%d条数据：工作地点不能为空", i+1))
			continue
		}

		jobs = append(jobs, models.Job{
			CompanyID:    companyID,
			Title:        item.Title,
			Description:  item.Description,
			SalaryMin:    item.SalaryMin,
			SalaryMax:    item.SalaryMax,
			Location:     item.Location,
			JobType:      item.JobType,
			Experience:   item.Experience,
			Skills:       item.Skills,
			Requirements: item.Requirements,
			Status:       models.JobStatusDraft,
		})
	}

	if len(jobs) > 0 {
		if err := s.jobRepo.BulkCreate(jobs); err != nil {
			errors = append(errors, "批量导入失败："+err.Error())
			return 0, errors
		}
	}

	return len(jobs), errors
}

func (s *JobService) BulkDelete(ids []uint, companyID uint) error {
	for _, id := range ids {
		job, err := s.jobRepo.FindByID(id)
		if err != nil || job.CompanyID != companyID {
			continue
		}
	}
	return s.jobRepo.BulkDelete(ids)
}

func (s *JobService) ExportJobs(companyID uint) ([]map[string]interface{}, error) {
	jobs, err := s.jobRepo.GetJobsByCompanyID(companyID)
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	for _, job := range jobs {
		data = append(data, map[string]interface{}{
			"id":          job.ID,
			"title":       job.Title,
			"location":    job.Location,
			"salary":      "面议",
			"job_type":    job.JobType,
			"experience":  job.Experience,
			"status":      job.Status,
			"view_count":  job.ViewCount,
			"apply_count": job.ApplyCount,
			"created_at":  job.CreatedAt.Format("2006-01-02 15:04:05"),
		})
		if job.SalaryMin > 0 && job.SalaryMax > 0 {
			data[len(data)-1]["salary"] = fmt.Sprintf("%dK-%dK", job.SalaryMin/1000, job.SalaryMax/1000)
		}
	}

	return data, nil
}

func (s *JobService) GetViewStats(jobID uint, days int) (int64, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	return s.jobRepo.GetViewStats(jobID, startDate, time.Now())
}
