package services

import (
	"errors"
	"qa-platform/models"
	"qa-platform/repository"
	"qa-platform/utils"
	"time"

	"gorm.io/gorm"
)

type QuestionService struct{}

func NewQuestionService() *QuestionService {
	return &QuestionService{}
}

type CreateQuestionRequest struct {
	Title        string `json:"title" binding:"required,min=5,max=200"`
	Content      string `json:"content" binding:"required,min=10"`
	CategoryID   uint   `json:"categoryId" binding:"required"`
	TagIDs       []uint `json:"tagIds"`
	RewardPoints int    `json:"rewardPoints"`
}

type UpdateQuestionRequest struct {
	Title   string `json:"title" binding:"min=5,max=200"`
	Content string `json:"content" binding:"min=10"`
}

type QuestionListQuery struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
	CategoryID uint   `form:"categoryId"`
	TagID      uint   `form:"tagId"`
	Keyword    string `form:"keyword"`
	Sort       string `form:"sort"`
	Status     string `form:"status"`
	UserID     uint   `form:"userId"`
}

type QuestionDetailResponse struct {
	Question       models.Question   `json:"question"`
	Tags           []models.Tag      `json:"tags"`
	User           UserResponse      `json:"user"`
	Category       models.Category   `json:"category"`
	Answers        []AnswerResponse  `json:"answers"`
	AnswerCount    int64             `json:"answerCount"`
	IsFavorited    bool              `json:"isFavorited"`
}

type AnswerResponse struct {
	Answer      models.Answer `json:"answer"`
	User        UserResponse  `json:"user"`
	Comments    []CommentResponse `json:"comments"`
	IsLiked     bool          `json:"isLiked"`
	IsFavorited bool          `json:"isFavorited"`
}

type CommentResponse struct {
	Comment models.Comment `json:"comment"`
	User    UserResponse   `json:"user"`
	IsLiked bool           `json:"isLiked"`
}

func (s *QuestionService) CreateQuestion(userID uint, req *CreateQuestionRequest) (*models.Question, error) {
	var user models.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	if req.RewardPoints > user.Points {
		return nil, errors.New("积分不足")
	}

	hasSensitive, sensitiveWords := utils.SensitiveFilter.Check(req.Title + req.Content)
	if hasSensitive {
		return nil, errors.New("内容包含敏感词: " + sensitiveWords[0])
	}

	var question models.Question
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		question = models.Question{
			UserID:       userID,
			Title:        req.Title,
			Content:      req.Content,
			CategoryID:   req.CategoryID,
			RewardPoints: req.RewardPoints,
			Status:       "published",
			AuditStatus:  "pending",
		}

		if err := tx.Create(&question).Error; err != nil {
			return err
		}

		if len(req.TagIDs) > 0 {
			var tags []models.Tag
			tx.Find(&tags, req.TagIDs)
			tx.Model(&question).Association("Tags").Append(tags)
			for _, tag := range tags {
				tx.Model(&tag).UpdateColumn("usage_count", tag.UsageCount+1)
			}
		}

		if req.RewardPoints > 0 {
			if err := AddPoints(userID, -req.RewardPoints, "reward_question", "设置悬赏", "question", question.ID); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (s *QuestionService) GetQuestionByID(id uint, currentUserID uint) (*QuestionDetailResponse, error) {
	var question models.Question
	if err := repository.DB.Preload("User").Preload("Category").Preload("Tags").First(&question, id).Error; err != nil {
		return nil, errors.New("问题不存在")
	}

	repository.DB.Model(&question).UpdateColumn("views", question.Views+1)

	var answers []models.Answer
	repository.DB.Where("question_id = ? AND status = ?", id, "published").
		Order("is_accepted DESC, like_count DESC, created_at ASC").
		Preload("User").
		Find(&answers)

	var answerResponses []AnswerResponse
	for _, answer := range answers {
		var comments []models.Comment
		repository.DB.Where("answer_id = ? AND status = ?", answer.ID, "published").
			Preload("User").
			Find(&comments)

		var commentResponses []CommentResponse
		for _, comment := range comments {
			commentResponses = append(commentResponses, CommentResponse{
				Comment: comment,
				User:    toUserResponse(comment.User),
				IsLiked: false,
			})
		}

		isFavorited := false
		if currentUserID > 0 {
			var fav models.Favorite
			result := repository.DB.Where("user_id = ? AND target_type = ? AND target_id = ?",
				currentUserID, "answer", answer.ID).First(&fav)
			isFavorited = result.Error == nil
		}

		answerResponses = append(answerResponses, AnswerResponse{
			Answer:      answer,
			User:        toUserResponse(answer.User),
			Comments:    commentResponses,
			IsFavorited: isFavorited,
		})
	}

	var answerCount int64
	repository.DB.Model(&models.Answer{}).Where("question_id = ? AND status = ?", id, "published").Count(&answerCount)

	isFavorited := false
	if currentUserID > 0 {
		var fav models.Favorite
		result := repository.DB.Where("user_id = ? AND target_type = ? AND target_id = ?",
			currentUserID, "question", id).First(&fav)
		isFavorited = result.Error == nil
	}

	return &QuestionDetailResponse{
		Question:    question,
		Tags:        question.Tags,
		User:        toUserResponse(question.User),
		Category:    question.Category,
		Answers:     answerResponses,
		AnswerCount: answerCount,
		IsFavorited: isFavorited,
	}, nil
}

func (s *QuestionService) GetQuestionList(query QuestionListQuery) ([]models.Question, int64, error) {
	var questions []models.Question
	var total int64

	dbQuery := repository.DB.Model(&models.Question{}).
		Preload("User").Preload("Category").Preload("Tags").
		Where("status = ?", "published").
		Where("audit_status = ?", "approved")

	if query.CategoryID > 0 {
		dbQuery = dbQuery.Where("category_id = ?", query.CategoryID)
	}

	if query.Keyword != "" {
		dbQuery = dbQuery.Where("title LIKE ? OR content LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	if query.UserID > 0 {
		dbQuery = dbQuery.Where("user_id = ?", query.UserID)
	}

	if query.TagID > 0 {
		dbQuery = dbQuery.Joins("JOIN question_tags ON question_tags.question_id = questions.id").
			Where("question_tags.tag_id = ?", query.TagID)
	}

	dbQuery.Count(&total)

	switch query.Sort {
	case "hot":
		dbQuery = dbQuery.Order("hot_score DESC")
	case "newest":
		dbQuery = dbQuery.Order("created_at DESC")
	case "unanswered":
		dbQuery = dbQuery.Where("answer_count = 0").Order("created_at DESC")
	default:
		dbQuery = dbQuery.Order("created_at DESC")
	}

	if query.Page > 0 && query.PageSize > 0 {
		dbQuery = dbQuery.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize)
	}

	dbQuery.Find(&questions)
	return questions, total, nil
}

func (s *QuestionService) UpdateQuestion(id, userID uint, req *UpdateQuestionRequest) error {
	var question models.Question
	if err := repository.DB.First(&question, id).Error; err != nil {
		return errors.New("问题不存在")
	}

	if question.UserID != userID {
		return errors.New("无权修改")
	}

	hasSensitive, sensitiveWords := utils.SensitiveFilter.Check(req.Title + req.Content)
	if hasSensitive {
		return errors.New("内容包含敏感词: " + sensitiveWords[0])
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}

	return repository.DB.Model(&question).Updates(updates).Error
}

func (s *QuestionService) DeleteQuestion(id, userID uint, role string) error {
	var question models.Question
	if err := repository.DB.First(&question, id).Error; err != nil {
		return errors.New("问题不存在")
	}

	if question.UserID != userID && role != "admin" {
		return errors.New("无权删除")
	}

	question.Status = "deleted"
	return repository.DB.Save(&question).Error
}

func (s *QuestionService) AcceptAnswer(questionID, answerID, userID uint) error {
	var question models.Question
	if err := repository.DB.First(&question, questionID).Error; err != nil {
		return errors.New("问题不存在")
	}

	if question.UserID != userID {
		return errors.New("只有提问者可以采纳答案")
	}

	if question.IsSolved {
		return errors.New("问题已解决")
	}

	var answer models.Answer
	if err := repository.DB.First(&answer, answerID).Error; err != nil {
		return errors.New("回答不存在")
	}

	return repository.DB.Transaction(func(tx *gorm.DB) error {
		answer.IsAccepted = true
		if err := tx.Save(&answer).Error; err != nil {
			return err
		}

		question.IsSolved = true
		question.AcceptedAnswerID = &answerID
		if err := tx.Save(&question).Error; err != nil {
			return err
		}

		if question.RewardPoints > 0 {
			if err := AddPoints(answer.UserID, question.RewardPoints, "reward_received", "获得悬赏", "question", question.ID); err != nil {
				return err
			}
		} else {
			if err := AddPoints(answer.UserID, 10, "answer_accepted", "回答被采纳", "question", question.ID); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *QuestionService) LikeQuestion(id, userID uint) error {
	var question models.Question
	if err := repository.DB.First(&question, id).Error; err != nil {
		return errors.New("问题不存在")
	}

	question.LikeCount++
	return repository.DB.Save(&question).Error
}

func (s *QuestionService) UpdateHotScore() {
	var questions []models.Question
	repository.DB.Where("status = ?", "published").Find(&questions)

	now := time.Now()
	for _, q := range questions {
		hours := now.Sub(q.CreatedAt).Hours()
		decay := 1.0
		if hours > 24 {
			decay = 1.0 / (hours / 24)
		}

		hotScore := float64(q.Views)*0.1 + float64(q.LikeCount)*2 + float64(q.AnswerCount)*5 + float64(q.CollectCount)*3
		hotScore *= decay

		repository.DB.Model(&q).Update("hot_score", hotScore)
	}
}
