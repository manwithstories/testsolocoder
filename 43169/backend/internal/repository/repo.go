package repository

import (
	"matchmaking-platform/internal/model"
	"matchmaking-platform/internal/utils"
)

type MatchRepo struct{}

func NewMatchRepo() *MatchRepo {
	return &MatchRepo{}
}

func (r *MatchRepo) Create(record *model.MatchRecord) error {
	return utils.DB.Create(record).Error
}

func (r *MatchRepo) FindByUserAndTarget(userID, targetID uint) (*model.MatchRecord, error) {
	var record model.MatchRecord
	err := utils.DB.Where("user_id = ? AND target_id = ?", userID, targetID).First(&record).Error
	return &record, err
}

func (r *MatchRepo) Update(id uint, updates map[string]interface{}) error {
	return utils.DB.Model(&model.MatchRecord{}).Where("id = ?", id).Updates(updates).Error
}

func (r *MatchRepo) GetFavorites(userID uint, page, pageSize int) ([]model.MatchRecord, int64, error) {
	var records []model.MatchRecord
	var total int64
	db := utils.DB.Model(&model.MatchRecord{}).Where("user_id = ? AND is_favorited = ?", userID, true)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&records).Error
	return records, total, err
}

func (r *MatchRepo) GetBlocked(userID uint, page, pageSize int) ([]model.MatchRecord, int64, error) {
	var records []model.MatchRecord
	var total int64
	db := utils.DB.Model(&model.MatchRecord{}).Where("user_id = ? AND is_blocked = ?", userID, true)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&records).Error
	return records, total, err
}

func (r *MatchRepo) GetBlockedIDs(userID uint) ([]uint, error) {
	var ids []uint
	err := utils.DB.Model(&model.MatchRecord{}).Where("user_id = ? AND is_blocked = ?", userID, true).Pluck("target_id", &ids).Error
	return ids, err
}

func (r *MatchRepo) CountMatchSuccess() (int64, error) {
	var count int64
	err := utils.DB.Model(&model.MatchRecord{}).Where("is_favorited = ?", true).Count(&count).Error
	return count, err
}

type DateRepo struct{}

func NewDateRepo() *DateRepo {
	return &DateRepo{}
}

func (r *DateRepo) Create(record *model.DateRecord) error {
	return utils.DB.Create(record).Error
}

func (r *DateRepo) FindByID(id uint) (*model.DateRecord, error) {
	var record model.DateRecord
	err := utils.DB.First(&record, id).Error
	return &record, err
}

func (r *DateRepo) Update(id uint, updates map[string]interface{}) error {
	return utils.DB.Model(&model.DateRecord{}).Where("id = ?", id).Updates(updates).Error
}

func (r *DateRepo) ListByUser(userID uint, page, pageSize int) ([]model.DateRecord, int64, error) {
	var records []model.DateRecord
	var total int64
	db := utils.DB.Model(&model.DateRecord{}).Where("initiator_id = ? OR receiver_id = ?", userID, userID)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("date_at DESC").Find(&records).Error
	return records, total, err
}

func (r *DateRepo) ListByStatus(status model.DateStatus, page, pageSize int) ([]model.DateRecord, int64, error) {
	var records []model.DateRecord
	var total int64
	db := utils.DB.Model(&model.DateRecord{}).Where("status = ?", status)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("date_at DESC").Find(&records).Error
	return records, total, err
}

func (r *DateRepo) CountTotal() (int64, error) {
	var count int64
	err := utils.DB.Model(&model.DateRecord{}).Count(&count).Error
	return count, err
}

func (r *DateRepo) CountByStatus(status model.DateStatus) (int64, error) {
	var count int64
	err := utils.DB.Model(&model.DateRecord{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

type DateReviewRepo struct{}

func NewDateReviewRepo() *DateReviewRepo {
	return &DateReviewRepo{}
}

func (r *DateReviewRepo) Create(review *model.DateReview) error {
	return utils.DB.Create(review).Error
}

func (r *DateReviewRepo) FindByDateAndReviewer(dateID, reviewerID uint) (*model.DateReview, error) {
	var review model.DateReview
	err := utils.DB.Where("date_id = ? AND reviewer_id = ?", dateID, reviewerID).First(&review).Error
	return &review, err
}

func (r *DateReviewRepo) ListByUser(userID uint, page, pageSize int) ([]model.DateReview, int64, error) {
	var reviews []model.DateReview
	var total int64
	db := utils.DB.Model(&model.DateReview{}).Where("reviewer_id = ? OR target_id = ?", userID, userID)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&reviews).Error
	return reviews, total, err
}

type MatchmakerRepo struct{}

func NewMatchmakerRepo() *MatchmakerRepo {
	return &MatchmakerRepo{}
}

func (r *MatchmakerRepo) AddMember(member *model.MatchmakerMember) error {
	return utils.DB.Create(member).Error
}

func (r *MatchmakerRepo) RemoveMember(matchmakerID, memberID uint) error {
	return utils.DB.Where("matchmaker_id = ? AND member_id = ?", matchmakerID, memberID).Delete(&model.MatchmakerMember{}).Error
}

func (r *MatchmakerRepo) ListMembers(matchmakerID uint, page, pageSize int) ([]model.MatchmakerMember, int64, error) {
	var members []model.MatchmakerMember
	var total int64
	db := utils.DB.Model(&model.MatchmakerMember{}).Where("matchmaker_id = ?", matchmakerID)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&members).Error
	return members, total, err
}

func (r *MatchmakerRepo) IsMember(matchmakerID, memberID uint) (bool, error) {
	var count int64
	err := utils.DB.Model(&model.MatchmakerMember{}).Where("matchmaker_id = ? AND member_id = ?", matchmakerID, memberID).Count(&count).Error
	return count > 0, err
}

func (r *MatchmakerRepo) CreateService(service *model.MatchmakerService) error {
	return utils.DB.Create(service).Error
}

func (r *MatchmakerRepo) UpdateService(id uint, updates map[string]interface{}) error {
	return utils.DB.Model(&model.MatchmakerService{}).Where("id = ?", id).Updates(updates).Error
}

func (r *MatchmakerRepo) ListServices(matchmakerID uint, page, pageSize int) ([]model.MatchmakerService, int64, error) {
	var services []model.MatchmakerService
	var total int64
	db := utils.DB.Model(&model.MatchmakerService{}).Where("matchmaker_id = ?", matchmakerID)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&services).Error
	return services, total, err
}

func (r *MatchmakerRepo) GetOrCreateStats(matchmakerID uint) (*model.MatchmakerStats, error) {
	var stats model.MatchmakerStats
	err := utils.DB.Where("matchmaker_id = ?", matchmakerID).First(&stats).Error
	if err != nil {
		stats = model.MatchmakerStats{MatchmakerID: matchmakerID}
		if err := utils.DB.Create(&stats).Error; err != nil {
			return nil, err
		}
	}
	return &stats, nil
}

func (r *MatchmakerRepo) UpdateStats(matchmakerID uint, updates map[string]interface{}) error {
	return utils.DB.Model(&model.MatchmakerStats{}).Where("matchmaker_id = ?", matchmakerID).Updates(updates).Error
}

type ChatRepo struct{}

func NewChatRepo() *ChatRepo {
	return &ChatRepo{}
}

func (r *ChatRepo) SendMessage(msg *model.ChatMessage) error {
	return utils.DB.Create(msg).Error
}

func (r *ChatRepo) GetHistory(userAID, userBID uint, page, pageSize int) ([]model.ChatMessage, int64, error) {
	var messages []model.ChatMessage
	var total int64
	db := utils.DB.Model(&model.ChatMessage{}).Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userAID, userBID, userBID, userAID,
	)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&messages).Error
	return messages, total, err
}

func (r *ChatRepo) MarkAsRead(senderID, receiverID uint) error {
	return utils.DB.Model(&model.ChatMessage{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = ?", senderID, receiverID, false).
		Update("is_read", true).Error
}

func (r *ChatRepo) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := utils.DB.Model(&model.ChatMessage{}).Where("receiver_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

func (r *ChatRepo) GetOrCreateSession(userAID, userBID uint) (*model.ChatSession, error) {
	var session model.ChatSession
	err := utils.DB.Where(
		"(user_a_id = ? AND user_b_id = ?) OR (user_a_id = ? AND user_b_id = ?)",
		userAID, userBID, userBID, userAID,
	).First(&session).Error
	if err != nil {
		session = model.ChatSession{
			UserAID: userAID,
			UserBID: userBID,
		}
		if err := utils.DB.Create(&session).Error; err != nil {
			return nil, err
		}
	}
	return &session, nil
}

func (r *ChatRepo) UpdateSession(session *model.ChatSession) error {
	return utils.DB.Save(session).Error
}

func (r *ChatRepo) ListSessions(userID uint, page, pageSize int) ([]model.ChatSession, int64, error) {
	var sessions []model.ChatSession
	var total int64
	db := utils.DB.Model(&model.ChatSession{}).Where("user_a_id = ? OR user_b_id = ?", userID, userID)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("last_time DESC").Find(&sessions).Error
	return sessions, total, err
}

type MemberRepo struct{}

func NewMemberRepo() *MemberRepo {
	return &MemberRepo{}
}

func (r *MemberRepo) GetBenefit(level string) (*model.MemberBenefit, error) {
	var benefit model.MemberBenefit
	err := utils.DB.Where("level = ?", level).First(&benefit).Error
	return &benefit, err
}

func (r *MemberRepo) ListBenefits() ([]model.MemberBenefit, error) {
	var benefits []model.MemberBenefit
	err := utils.DB.Find(&benefits).Error
	return benefits, err
}

func (r *MemberRepo) CreateOrder(order *model.MemberOrder) error {
	return utils.DB.Create(order).Error
}

func (r *MemberRepo) UpdateOrder(id uint, updates map[string]interface{}) error {
	return utils.DB.Model(&model.MemberOrder{}).Where("id = ?", id).Updates(updates).Error
}

func (r *MemberRepo) ListOrders(userID uint, page, pageSize int) ([]model.MemberOrder, int64, error) {
	var orders []model.MemberOrder
	var total int64
	db := utils.DB.Model(&model.MemberOrder{}).Where("user_id = ?", userID)
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&orders).Error
	return orders, total, err
}

func (r *MemberRepo) LogInteract(log *model.InteractLog) error {
	return utils.DB.Create(log).Error
}

func (r *MemberRepo) CountTodayInteracts(userID uint) (int64, error) {
	var count int64
	err := utils.DB.Model(&model.InteractLog{}).
		Where("user_id = ? AND DATE(created_at) = CURDATE()", userID).
		Count(&count).Error
	return count, err
}

func (r *MemberRepo) GetExpiredMembers() ([]model.User, error) {
	var users []model.User
	err := utils.DB.Model(&model.User{}).
		Where("member_level != ? AND member_expire IS NOT NULL AND member_expire < NOW()", model.MemberFree).
		Find(&users).Error
	return users, err
}

type SensitiveWordRepo struct{}

func NewSensitiveWordRepo() *SensitiveWordRepo {
	return &SensitiveWordRepo{}
}

func (r *SensitiveWordRepo) ListAll() ([]model.SensitiveWord, error) {
	var words []model.SensitiveWord
	err := utils.DB.Find(&words).Error
	return words, err
}

type SystemLogRepo struct{}

func NewSystemLogRepo() *SystemLogRepo {
	return &SystemLogRepo{}
}

func (r *SystemLogRepo) Create(log *model.SystemLog) error {
	return utils.DB.Create(log).Error
}

func (r *SystemLogRepo) List(page, pageSize int, filter map[string]interface{}) ([]model.SystemLog, int64, error) {
	var logs []model.SystemLog
	var total int64
	db := utils.DB.Model(&model.SystemLog{})
	for key, val := range filter {
		db = db.Where(key+" = ?", val)
	}
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&logs).Error
	return logs, total, err
}
