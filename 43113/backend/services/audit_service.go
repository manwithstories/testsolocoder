package services

import (
	"errors"
	"qa-platform/models"
	"qa-platform/repository"
	"qa-platform/utils"
	"time"
)

type AuditService struct{}

func NewAuditService() *AuditService {
	return &AuditService{}
}

type AuditQuery struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
	TargetType string `form:"targetType"`
	Status     string `form:"status"`
}

type AuditRequest struct {
	ID     uint   `json:"id" binding:"required"`
	TargetType string `json:"targetType" binding:"required"`
	Action string `json:"action" binding:"required"`
	Reason string `json:"reason"`
}

type ReportQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Status   string `form:"status"`
}

type HandleReportRequest struct {
	Status  string `json:"status" binding:"required"`
	Result  string `json:"result"`
}

func (s *AuditService) GetAuditList(query AuditQuery) ([]models.AuditRecord, int64, error) {
	var records []models.AuditRecord
	var total int64

	dbQuery := repository.DB.Preload("Admin")
	if query.TargetType != "" {
		dbQuery = dbQuery.Where("target_type = ?", query.TargetType)
	}
	if query.Status != "" {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	dbQuery.Model(&models.AuditRecord{}).Count(&total)
	dbQuery.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize).
		Order("created_at DESC").Find(&records)

	return records, total, nil
}

func (s *AuditService) CreateAuditRecord(adminID uint, targetType string, targetID uint, action, reason, status, ip string) error {
	record := models.AuditRecord{
		AdminID:    adminID,
		TargetType: targetType,
		TargetID:   targetID,
		Action:     action,
		Reason:     reason,
		Status:     status,
		OperatorIP: ip,
	}

	return repository.DB.Create(&record).Error
}

func (s *AuditService) AuditContent(adminID uint, req *AuditRequest) error {
	switch req.TargetType {
	case "question":
		var question models.Question
		if err := repository.DB.First(&question, req.ID).Error; err != nil {
			return errors.New("问题不存在")
		}

		question.AuditStatus = req.Action
		if err := repository.DB.Save(&question).Error; err != nil {
			return err
		}

		s.CreateAuditRecord(adminID, "question", question.ID, req.Action, req.Reason, "approved", "")

	case "answer":
		var answer models.Answer
		if err := repository.DB.First(&answer, req.ID).Error; err != nil {
			return errors.New("回答不存在")
		}

		answer.AuditStatus = req.Action
		if err := repository.DB.Save(&answer).Error; err != nil {
			return err
		}

		s.CreateAuditRecord(adminID, "answer", answer.ID, req.Action, req.Reason, "approved", "")

	case "comment":
		var comment models.Comment
		if err := repository.DB.First(&comment, req.ID).Error; err != nil {
			return errors.New("评论不存在")
		}

		comment.AuditStatus = req.Action
		if err := repository.DB.Save(&comment).Error; err != nil {
			return err
		}

		s.CreateAuditRecord(adminID, "comment", comment.ID, req.Action, req.Reason, "approved", "")
	}

	return nil
}

func (s *AuditService) GetPendingAuditCount() (map[string]int64, error) {
	result := make(map[string]int64)

	var questionCount int64
	repository.DB.Model(&models.Question{}).Where("audit_status = ?", "pending").Count(&questionCount)
	result["questions"] = questionCount

	var answerCount int64
	repository.DB.Model(&models.Answer{}).Where("audit_status = ?", "pending").Count(&answerCount)
	result["answers"] = answerCount

	var commentCount int64
	repository.DB.Model(&models.Comment{}).Where("audit_status = ?", "pending").Count(&commentCount)
	result["comments"] = commentCount

	return result, nil
}

func (s *AuditService) GetReportList(query ReportQuery) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	dbQuery := repository.DB.Preload("Reporter").Preload("Handler")
	if query.Status != "" {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}

	dbQuery.Model(&models.Report{}).Count(&total)
	dbQuery.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize).
		Order("created_at DESC").Find(&reports)

	return reports, total, nil
}

func (s *AuditService) HandleReport(id, handlerID uint, req *HandleReportRequest) error {
	var report models.Report
	if err := repository.DB.First(&report, id).Error; err != nil {
		return errors.New("举报不存在")
	}

	now := time.Now()
	report.Status = req.Status
	report.HandlerID = &handlerID
	report.HandleResult = req.Result
	report.HandledAt = &now

	if err := repository.DB.Save(&report).Error; err != nil {
		return err
	}

	if req.Status == "resolved" && req.Result == "delete" {
		switch report.TargetType {
		case "question":
			repository.DB.Model(&models.Question{}).Where("id = ?", report.TargetID).Update("status", "deleted")
		case "answer":
			repository.DB.Model(&models.Answer{}).Where("id = ?", report.TargetID).Update("status", "deleted")
		case "comment":
			repository.DB.Model(&models.Comment{}).Where("id = ?", report.TargetID).Update("status", "deleted")
		}
	}

	return nil
}

func (s *AuditService) CreateReport(reporterID uint, targetType string, targetID uint, reason, description string) error {
	report := models.Report{
		ReporterID:  reporterID,
		TargetType:  targetType,
		TargetID:    targetID,
		Reason:      reason,
		Description: description,
		Status:      "pending",
	}

	return repository.DB.Create(&report).Error
}

type SensitiveWordQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Category string `form:"category"`
	Keyword  string `form:"keyword"`
}

func (s *AuditService) GetSensitiveWords(query SensitiveWordQuery) ([]models.SensitiveWord, int64, error) {
	var words []models.SensitiveWord
	var total int64

	dbQuery := repository.DB.Model(&models.SensitiveWord{})
	if query.Category != "" {
		dbQuery = dbQuery.Where("category = ?", query.Category)
	}
	if query.Keyword != "" {
		dbQuery = dbQuery.Where("word LIKE ?", "%"+query.Keyword+"%")
	}

	dbQuery.Count(&total)
	dbQuery.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize).
		Order("created_at DESC").Find(&words)

	return words, total, nil
}

func (s *AuditService) CreateSensitiveWord(word, category, replaceTo string, level int) error {
	existingWord := models.SensitiveWord{}
	result := repository.DB.Where("word = ?", word).First(&existingWord)
	if result.Error == nil {
		return errors.New("敏感词已存在")
	}

	sw := models.SensitiveWord{
		Word:      word,
		Category:  category,
		ReplaceTo: replaceTo,
		Level:     level,
	}

	if err := repository.DB.Create(&sw).Error; err != nil {
		return err
	}

	utils.SensitiveFilter.AddWord(word, replaceTo)
	return nil
}

func (s *AuditService) DeleteSensitiveWord(id uint) error {
	return repository.DB.Delete(&models.SensitiveWord{}, id).Error
}

func (s *AuditService) BatchCheckContent(content string) (string, []string) {
	return utils.SensitiveFilter.Filter(content)
}
