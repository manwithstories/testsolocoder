package services

import (
	"errors"
	"qa-platform/models"
	"qa-platform/repository"

	"gorm.io/gorm"
)

type RewardService struct{}

func NewRewardService() *RewardService {
	return &RewardService{}
}

func (s *RewardService) GetRewardList(page, pageSize int) ([]models.Reward, int64, error) {
	var rewards []models.Reward
	var total int64

	dbQuery := repository.DB.Where("status = ?", "active")
	dbQuery.Model(&models.Reward{}).Count(&total)
	dbQuery.Offset((page - 1) * pageSize).Limit(pageSize).Find(&rewards)

	return rewards, total, nil
}

func (s *RewardService) CreateReward(name, description, image string, pointsCost, stock int) (*models.Reward, error) {
	reward := models.Reward{
		Name:        name,
		Description: description,
		Image:       image,
		PointsCost:  pointsCost,
		Stock:       stock,
		Status:      "active",
	}

	if err := repository.DB.Create(&reward).Error; err != nil {
		return nil, err
	}

	return &reward, nil
}

func (s *RewardService) UpdateReward(id uint, updates map[string]interface{}) error {
	return repository.DB.Model(&models.Reward{}).Where("id = ?", id).Updates(updates).Error
}

func (s *RewardService) DeleteReward(id uint) error {
	return repository.DB.Model(&models.Reward{}).Where("id = ?", id).Update("status", "deleted").Error
}

func (s *RewardService) ExchangeReward(userID, rewardID uint) error {
	var reward models.Reward
	if err := repository.DB.First(&reward, rewardID).Error; err != nil {
		return errors.New("奖品不存在")
	}

	if reward.Status != "active" {
		return errors.New("奖品不可兑换")
	}

	if reward.Stock == 0 {
		return errors.New("奖品已兑换完")
	}

	var user models.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	if user.Points < reward.PointsCost {
		return errors.New("积分不足")
	}

	return repository.DB.Transaction(func(tx *gorm.DB) error {
		newPoints := user.Points - reward.PointsCost
		newLevel := CalculateLevel(newPoints)

		if err := tx.Model(&user).Updates(map[string]interface{}{
			"points": newPoints,
			"level":  newLevel,
		}).Error; err != nil {
			return err
		}

		if reward.Stock > 0 {
			if err := tx.Model(&reward).Update("stock", reward.Stock-1).Error; err != nil {
				return err
			}
		}

		exchange := models.RewardExchange{
			UserID:     userID,
			RewardID:   rewardID,
			PointsCost: reward.PointsCost,
			Status:     "completed",
		}
		if err := tx.Create(&exchange).Error; err != nil {
			return err
		}

		pointLog := models.PointLog{
			UserID:      userID,
			Type:        "reward_exchange",
			Points:      -reward.PointsCost,
			Balance:     newPoints,
			Description: "兑换: " + reward.Name,
			RefType:     "reward",
			RefID:       rewardID,
		}
		return tx.Create(&pointLog).Error
	})
}

func (s *RewardService) GetExchangeList(userID uint, page, pageSize int) ([]models.RewardExchange, int64, error) {
	var exchanges []models.RewardExchange
	var total int64

	dbQuery := repository.DB.Preload("Reward").Where("user_id = ?", userID)
	dbQuery.Model(&models.RewardExchange{}).Count(&total)
	dbQuery.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&exchanges)

	return exchanges, total, nil
}

func (s *RewardService) GetPointLogs(userID uint, page, pageSize int) ([]models.PointLog, int64, error) {
	var logs []models.PointLog
	var total int64

	dbQuery := repository.DB.Where("user_id = ?", userID)
	dbQuery.Model(&models.PointLog{}).Count(&total)
	dbQuery.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&logs)

	return logs, total, nil
}

type SearchService struct{}

func NewSearchService() *SearchService {
	return &SearchService{}
}

type SearchQuery struct {
	Keyword    string `form:"keyword"`
	CategoryID uint   `form:"categoryId"`
	TagID      uint   `form:"tagId"`
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
	Sort       string `form:"sort"`
}

func (s *SearchService) SearchQuestions(query SearchQuery) ([]models.Question, int64, error) {
	var questions []models.Question
	var total int64

	dbQuery := repository.DB.Model(&models.Question{}).
		Preload("User").Preload("Category").Preload("Tags").
		Where("status = ? AND audit_status = ?", "published", "approved")

	if query.Keyword != "" {
		dbQuery = dbQuery.Where("title LIKE ? OR content LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	if query.CategoryID > 0 {
		dbQuery = dbQuery.Where("category_id = ?", query.CategoryID)
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
	case "most_answers":
		dbQuery = dbQuery.Order("answer_count DESC")
	default:
		dbQuery = dbQuery.Order("created_at DESC")
	}

	dbQuery.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize).Find(&questions)

	return questions, total, nil
}

func (s *SearchService) GetRecommendations(userID uint, limit int) ([]models.Question, error) {
	var questions []models.Question

	var user models.User
	repository.DB.First(&user, userID)

	var follows []models.Follow
	repository.DB.Where("follower_id = ? AND following_type = ?", userID, "tag").Find(&follows)

	var tagIDs []uint
	for _, f := range follows {
		tagIDs = append(tagIDs, f.FollowingID)
	}

	var favorites []models.Favorite
	repository.DB.Where("user_id = ? AND target_type = ?", userID, "question").Find(&favorites)

	var favQuestionIDs []uint
	for _, fav := range favorites {
		favQuestionIDs = append(favQuestionIDs, fav.TargetID)
	}

	var favTagIDs []uint
	if len(favQuestionIDs) > 0 {
		var favQuestions []models.Question
		repository.DB.Preload("Tags").Find(&favQuestions, favQuestionIDs)
		for _, q := range favQuestions {
			for _, tag := range q.Tags {
				favTagIDs = append(favTagIDs, tag.ID)
			}
		}
	}

	allTagIDs := append(tagIDs, favTagIDs...)

	dbQuery := repository.DB.Model(&models.Question{}).
		Preload("User").Preload("Category").Preload("Tags").
		Where("status = ? AND audit_status = ?", "published", "approved").
		Where("user_id != ?", userID)

	if len(allTagIDs) > 0 {
		dbQuery = dbQuery.Joins("JOIN question_tags ON question_tags.question_id = questions.id").
			Where("question_tags.tag_id IN ?", allTagIDs)
	}

	dbQuery.Order("hot_score DESC, created_at DESC").Limit(limit).Find(&questions)

	if len(questions) < limit {
		var extraQuestions []models.Question
		repository.DB.Model(&models.Question{}).
			Preload("User").Preload("Category").Preload("Tags").
			Where("status = ? AND audit_status = ?", "published", "approved").
			Where("user_id != ?", userID).
			Where("id NOT IN ?", func() []uint {
				var ids []uint
				for _, q := range questions {
					ids = append(ids, q.ID)
				}
				return append(ids, 0)
			}()).
			Order("hot_score DESC, created_at DESC").
			Limit(limit - len(questions)).
			Find(&extraQuestions)
		questions = append(questions, extraQuestions...)
	}

	return questions, nil
}

type StatsService struct{}

func NewStatsService() *StatsService {
	return &StatsService{}
}

type DashboardStats struct {
	TotalQuestions    int64 `json:"totalQuestions"`
	TotalAnswers      int64 `json:"totalAnswers"`
	TotalUsers        int64 `json:"totalUsers"`
	TotalComments     int64 `json:"totalComments"`
	PendingAuditCount int64 `json:"pendingAuditCount"`
	TodayNewQuestions int64 `json:"todayNewQuestions"`
	TodayNewAnswers   int64 `json:"todayNewAnswers"`
	TodayNewUsers     int64 `json:"todayNewUsers"`
}

func (s *StatsService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	repository.DB.Model(&models.Question{}).Where("status = ?", "published").Count(&stats.TotalQuestions)
	repository.DB.Model(&models.Answer{}).Where("status = ?", "published").Count(&stats.TotalAnswers)
	repository.DB.Model(&models.User{}).Where("status = ?", "active").Count(&stats.TotalUsers)
	repository.DB.Model(&models.Comment{}).Where("status = ?", "published").Count(&stats.TotalComments)

	var pendingCount int64
	repository.DB.Model(&models.Question{}).Where("audit_status = ?", "pending").Count(&pendingCount)
	var pendingAnswers int64
	repository.DB.Model(&models.Answer{}).Where("audit_status = ?", "pending").Count(&pendingAnswers)
	var pendingComments int64
	repository.DB.Model(&models.Comment{}).Where("audit_status = ?", "pending").Count(&pendingComments)
	stats.PendingAuditCount = pendingCount + pendingAnswers + pendingComments

	repository.DB.Model(&models.Question{}).Where("DATE(created_at) = CURDATE()").Count(&stats.TodayNewQuestions)
	repository.DB.Model(&models.Answer{}).Where("DATE(created_at) = CURDATE()").Count(&stats.TodayNewAnswers)
	repository.DB.Model(&models.User{}).Where("DATE(created_at) = CURDATE()").Count(&stats.TodayNewUsers)

	return stats, nil
}

type ActivityReport struct {
	Date             string `json:"date"`
	NewQuestions     int64  `json:"newQuestions"`
	NewAnswers       int64  `json:"newAnswers"`
	NewUsers         int64  `json:"newUsers"`
	NewComments      int64  `json:"newComments"`
}

func (s *StatsService) GetActivityReport(startDate, endDate string) ([]ActivityReport, error) {
	var reports []ActivityReport

	rows, err := repository.DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(CASE WHEN target_type = 'question' THEN 1 END) as new_questions,
			COUNT(CASE WHEN target_type = 'answer' THEN 1 END) as new_answers,
			COUNT(CASE WHEN target_type = 'user' THEN 1 END) as new_users,
			COUNT(CASE WHEN target_type = 'comment' THEN 1 END) as new_comments
		FROM (
			SELECT 'question' as target_type, created_at FROM questions WHERE created_at BETWEEN ? AND ?
			UNION ALL
			SELECT 'answer' as target_type, created_at FROM answers WHERE created_at BETWEEN ? AND ?
			UNION ALL
			SELECT 'user' as target_type, created_at FROM users WHERE created_at BETWEEN ? AND ?
			UNION ALL
			SELECT 'comment' as target_type, created_at FROM comments WHERE created_at BETWEEN ? AND ?
		) as combined
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`, startDate, endDate, startDate, endDate, startDate, endDate, startDate, endDate).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var report ActivityReport
		rows.Scan(&report.Date, &report.NewQuestions, &report.NewAnswers, &report.NewUsers, &report.NewComments)
		reports = append(reports, report)
	}

	return reports, nil
}

type AuditReport struct {
	Date           string `json:"date"`
	ReviewedCount  int64  `json:"reviewedCount"`
	ApprovedCount  int64  `json:"approvedCount"`
	RejectedCount  int64  `json:"rejectedCount"`
}

func (s *StatsService) GetAuditReport(startDate, endDate string) ([]AuditReport, error) {
	var reports []AuditReport

	rows, err := repository.DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as reviewed_count,
			COUNT(CASE WHEN status = 'approved' THEN 1 END) as approved_count,
			COUNT(CASE WHEN status = 'rejected' THEN 1 END) as rejected_count
		FROM audit_records
		WHERE created_at BETWEEN ? AND ?
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`, startDate, endDate).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var report AuditReport
		rows.Scan(&report.Date, &report.ReviewedCount, &report.ApprovedCount, &report.RejectedCount)
		reports = append(reports, report)
	}

	return reports, nil
}
