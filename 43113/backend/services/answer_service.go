package services

import (
	"errors"
	"qa-platform/models"
	"qa-platform/repository"
	"qa-platform/utils"
	"time"

	"gorm.io/gorm"
)

type AnswerService struct{}

func NewAnswerService() *AnswerService {
	return &AnswerService{}
}

type CreateAnswerRequest struct {
	QuestionID uint   `json:"questionId" binding:"required"`
	Content    string `json:"content" binding:"required,min=10"`
}

type UpdateAnswerRequest struct {
	Content string `json:"content" binding:"required,min=10"`
}

func (s *AnswerService) CreateAnswer(userID uint, req *CreateAnswerRequest) (*models.Answer, error) {
	var question models.Question
	if err := repository.DB.First(&question, req.QuestionID).Error; err != nil {
		return nil, errors.New("问题不存在")
	}

	if question.Status != "published" {
		return nil, errors.New("问题不存在")
	}

	hasSensitive, sensitiveWords := utils.SensitiveFilter.Check(req.Content)
	if hasSensitive {
		return nil, errors.New("内容包含敏感词: " + sensitiveWords[0])
	}

	var answer models.Answer
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		answer = models.Answer{
			QuestionID:  req.QuestionID,
			UserID:      userID,
			Content:     req.Content,
			Status:      "published",
			AuditStatus: "pending",
		}

		if err := tx.Create(&answer).Error; err != nil {
			return err
		}

		tx.Model(&question).UpdateColumn("answer_count", question.AnswerCount+1)

		if err := AddPoints(userID, 5, "answer_created", "回答问题", "question", question.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &answer, nil
}

func (s *AnswerService) GetAnswerByID(id uint) (*models.Answer, error) {
	var answer models.Answer
	if err := repository.DB.Preload("User").Preload("Question").First(&answer, id).Error; err != nil {
		return nil, errors.New("回答不存在")
	}
	return &answer, nil
}

func (s *AnswerService) UpdateAnswer(id, userID uint, req *UpdateAnswerRequest) error {
	var answer models.Answer
	if err := repository.DB.First(&answer, id).Error; err != nil {
		return errors.New("回答不存在")
	}

	if answer.UserID != userID {
		return errors.New("无权修改")
	}

	hasSensitive, sensitiveWords := utils.SensitiveFilter.Check(req.Content)
	if hasSensitive {
		return errors.New("内容包含敏感词: " + sensitiveWords[0])
	}

	answer.Content = req.Content
	answer.UpdatedAt = time.Now()
	return repository.DB.Save(&answer).Error
}

func (s *AnswerService) DeleteAnswer(id, userID uint, role string) error {
	var answer models.Answer
	if err := repository.DB.First(&answer, id).Error; err != nil {
		return errors.New("回答不存在")
	}

	if answer.UserID != userID && role != "admin" {
		return errors.New("无权删除")
	}

	answer.Status = "deleted"
	repository.DB.Save(&answer)

	var question models.Question
	repository.DB.First(&question, answer.QuestionID)
	if question.AnswerCount > 0 {
		repository.DB.Model(&question).UpdateColumn("answer_count", question.AnswerCount-1)
	}

	return nil
}

func (s *AnswerService) LikeAnswer(id, userID uint) error {
	var answer models.Answer
	if err := repository.DB.First(&answer, id).Error; err != nil {
		return errors.New("回答不存在")
	}

	answer.LikeCount++
	return repository.DB.Save(&answer).Error
}

func (s *AnswerService) DislikeAnswer(id, userID uint) error {
	var answer models.Answer
	if err := repository.DB.First(&answer, id).Error; err != nil {
		return errors.New("回答不存在")
	}

	answer.DislikeCount++
	return repository.DB.Save(&answer).Error
}

type CreateCommentRequest struct {
	AnswerID uint   `json:"answerId" binding:"required"`
	Content  string `json:"content" binding:"required,min=1,max=500"`
}

type CommentService struct{}

func NewCommentService() *CommentService {
	return &CommentService{}
}

func (s *CommentService) CreateComment(userID uint, req *CreateCommentRequest) (*models.Comment, error) {
	var answer models.Answer
	if err := repository.DB.First(&answer, req.AnswerID).Error; err != nil {
		return nil, errors.New("回答不存在")
	}

	hasSensitive, sensitiveWords := utils.SensitiveFilter.Check(req.Content)
	if hasSensitive {
		return nil, errors.New("内容包含敏感词: " + sensitiveWords[0])
	}

	comment := models.Comment{
		UserID:      userID,
		AnswerID:    req.AnswerID,
		Content:     req.Content,
		Status:      "published",
		AuditStatus: "pending",
	}

	if err := repository.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (s *CommentService) DeleteComment(id, userID uint, role string) error {
	var comment models.Comment
	if err := repository.DB.First(&comment, id).Error; err != nil {
		return errors.New("评论不存在")
	}

	if comment.UserID != userID && role != "admin" {
		return errors.New("无权删除")
	}

	comment.Status = "deleted"
	return repository.DB.Save(&comment).Error
}

func (s *CommentService) LikeComment(id, userID uint) error {
	var comment models.Comment
	if err := repository.DB.First(&comment, id).Error; err != nil {
		return errors.New("评论不存在")
	}

	comment.LikeCount++
	return repository.DB.Save(&comment).Error
}
