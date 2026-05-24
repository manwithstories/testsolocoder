package handlers

import (
	"strings"

	"recruitment-platform/middleware"
	"recruitment-platform/models"
	"recruitment-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecommendationHandler struct {
	DB *gorm.DB
}

func NewRecommendationHandler(db *gorm.DB) *RecommendationHandler {
	return &RecommendationHandler{DB: db}
}

func (h *RecommendationHandler) GetRecommendedJobs(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var jobSeeker models.JobSeeker
	if err := h.DB.Where("user_id = ?", userID).First(&jobSeeker).Error; err != nil {
		utils.NotFound(c, "Job seeker not found")
		return
	}

	var resumes []models.Resume
	h.DB.Preload("Skills").Where("jobseeker_id = ?", jobSeeker.ID).Find(&resumes)

	var skillNames []string
	for _, resume := range resumes {
		for _, skill := range resume.Skills {
			skillNames = append(skillNames, strings.ToLower(skill.Name))
		}
	}

	if len(skillNames) == 0 {
		var jobs []models.Job
		h.DB.Preload("Company").Preload("Department").
			Where("status = ?", models.JobStatusOpen).
			Order("created_at desc").Limit(10).Find(&jobs)
		utils.Success(c, jobs)
		return
	}

	var appliedJobIDs []uint
	h.DB.Model(&models.Application{}).Where("jobseeker_id = ?", jobSeeker.ID).Pluck("job_id", &appliedJobIDs)

	var jobs []models.Job
	query := h.DB.Preload("Company").Preload("Department").
		Where("status = ?", models.JobStatusOpen)

	if len(appliedJobIDs) > 0 {
		query = query.Where("id NOT IN ?", appliedJobIDs)
	}

	for _, skill := range skillNames {
		query = query.Or("skills ILIKE ? AND status = ?", "%"+skill+"%", models.JobStatusOpen)
	}

	query.Order("created_at desc").Limit(20).Find(&jobs)

	type scoredJob struct {
		models.Job
		MatchScore    int      `json:"match_score"`
		MatchedSkills []string `json:"matched_skills"`
	}

	var scoredJobs []scoredJob
	for _, job := range jobs {
		jobSkills := strings.Split(strings.ToLower(job.Skills), ",")
		matched := 0
		var matchedList []string
		for _, s := range skillNames {
			for _, js := range jobSkills {
				if strings.Contains(strings.TrimSpace(js), s) {
					matched++
					matchedList = append(matchedList, s)
					break
				}
			}
		}
		scoredJobs = append(scoredJobs, scoredJob{
			Job:           job,
			MatchScore:    matched,
			MatchedSkills: matchedList,
		})
	}

	for i := 0; i < len(scoredJobs)-1; i++ {
		for j := i + 1; j < len(scoredJobs); j++ {
			if scoredJobs[j].MatchScore > scoredJobs[i].MatchScore {
				scoredJobs[i], scoredJobs[j] = scoredJobs[j], scoredJobs[i]
			}
		}
	}

	utils.Success(c, scoredJobs)
}

func (h *RecommendationHandler) GetSimilarJobs(c *gin.Context) {
	jobID := c.Param("id")

	var job models.Job
	if err := h.DB.Preload("Company").First(&job, jobID).Error; err != nil {
		utils.NotFound(c, "Job not found")
		return
	}

	var similarJobs []models.Job
	query := h.DB.Preload("Company").Preload("Department").
		Where("status = ? AND id != ?", models.JobStatusOpen, job.ID)

	if job.Location != "" {
		query = query.Or("location ILIKE ? AND status = ? AND id != ?", "%"+job.Location+"%", models.JobStatusOpen, job.ID)
	}

	if job.JobType != "" {
		query = query.Or("job_type = ? AND status = ? AND id != ?", job.JobType, models.JobStatusOpen, job.ID)
	}

	if job.DepartmentID != nil {
		query = query.Or("department_id = ? AND status = ? AND id != ?", *job.DepartmentID, models.JobStatusOpen, job.ID)
	}

	jobSkills := strings.Split(strings.ToLower(job.Skills), ",")
	for _, skill := range jobSkills {
		skill = strings.TrimSpace(skill)
		if skill != "" {
			query = query.Or("skills ILIKE ? AND status = ? AND id != ?", "%"+skill+"%", models.JobStatusOpen, job.ID)
		}
	}

	query.Order("created_at desc").Limit(10).Find(&similarJobs)

	utils.Success(c, similarJobs)
}
