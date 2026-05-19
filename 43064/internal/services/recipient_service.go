package services

import (
	"strings"

	"github.com/notification-center/internal/database"
	"github.com/notification-center/internal/errors"
	"github.com/notification-center/internal/logger"
	"github.com/notification-center/internal/models"
	"github.com/notification-center/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RecipientService struct{}

func NewRecipientService() *RecipientService {
	return &RecipientService{}
}

func (s *RecipientService) Create(recipient *models.Recipient) (*models.Recipient, error) {
	if recipient.Email == "" && recipient.Phone == "" && recipient.WeChatOpenID == "" && recipient.DingTalkUserID == "" {
		return nil, errors.InvalidParams("at least one contact method is required", nil)
	}

	existing, err := s.FindDuplicate(recipient)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	if err := database.DB.Create(recipient).Error; err != nil {
		logger.Error("create recipient failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	logger.Info("recipient created", zap.Uint("id", recipient.ID), zap.String("email", recipient.Email))
	return recipient, nil
}

func (s *RecipientService) FindDuplicate(recipient *models.Recipient) (*models.Recipient, error) {
	var existing models.Recipient
	query := database.DB

	conditions := make([]string, 0)
	args := make([]interface{}, 0)

	if recipient.Email != "" {
		conditions = append(conditions, "email = ?")
		args = append(args, recipient.Email)
	}
	if recipient.Phone != "" {
		conditions = append(conditions, "phone = ?")
		args = append(args, recipient.Phone)
	}
	if recipient.WeChatOpenID != "" {
		conditions = append(conditions, "we_chat_open_id = ?")
		args = append(args, recipient.WeChatOpenID)
	}
	if recipient.DingTalkUserID != "" {
		conditions = append(conditions, "ding_talk_user_id = ?")
		args = append(args, recipient.DingTalkUserID)
	}

	if len(conditions) == 0 {
		return nil, nil
	}

	query = query.Where(strings.Join(conditions, " OR "), args...)

	if err := query.First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.DatabaseError(err)
	}

	return &existing, nil
}

func (s *RecipientService) GetByID(id uint) (*models.Recipient, error) {
	var recipient models.Recipient
	if err := database.DB.Preload("Tags").Preload("Groups").First(&recipient, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.RecipientNotFound(id)
		}
		logger.Error("get recipient failed", zap.Error(err), zap.Uint("id", id))
		return nil, errors.DatabaseError(err)
	}
	return &recipient, nil
}

func (s *RecipientService) List(keyword string, tagIDs, groupIDs []uint, enabledOnly bool, page, pageSize int) ([]models.Recipient, int64, error) {
	var recipients []models.Recipient
	var total int64

	query := database.DB.Model(&models.Recipient{}).Preload("Tags").Preload("Groups")

	if keyword != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if enabledOnly {
		query = query.Where("enabled = ?", true)
	}
	if len(tagIDs) > 0 {
		query = query.Joins("JOIN recipient_tags ON recipient_tags.recipient_id = recipients.id").
			Where("recipient_tags.tag_id IN (?)", tagIDs)
	}
	if len(groupIDs) > 0 {
		query = query.Joins("JOIN recipient_group_members ON recipient_group_members.recipient_id = recipients.id").
			Where("recipient_group_members.recipient_group_id IN (?)", groupIDs)
	}

	if err := query.Distinct("recipients.*").Count(&total).Error; err != nil {
		logger.Error("count recipients failed", zap.Error(err))
		return nil, 0, errors.DatabaseError(err)
	}

	offset := (page - 1) * pageSize
	if err := query.Distinct("recipients.*").Order("id DESC").Offset(offset).Limit(pageSize).Find(&recipients).Error; err != nil {
		logger.Error("list recipients failed", zap.Error(err))
		return nil, 0, errors.DatabaseError(err)
	}

	return recipients, total, nil
}

func (s *RecipientService) Update(id uint, updates map[string]interface{}) (*models.Recipient, error) {
	recipient, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := database.DB.Model(recipient).Updates(updates).Error; err != nil {
		logger.Error("update recipient failed", zap.Error(err), zap.Uint("id", id))
		return nil, errors.DatabaseError(err)
	}

	if err := database.DB.Preload("Tags").Preload("Groups").First(recipient, id).Error; err != nil {
		logger.Error("reload recipient failed", zap.Error(err), zap.Uint("id", id))
		return nil, errors.DatabaseError(err)
	}

	logger.Info("recipient updated", zap.Uint("id", id))
	return recipient, nil
}

func (s *RecipientService) Delete(id uint) error {
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err := database.DB.Delete(&models.Recipient{}, id).Error; err != nil {
		logger.Error("delete recipient failed", zap.Error(err), zap.Uint("id", id))
		return errors.DatabaseError(err)
	}

	logger.Info("recipient deleted", zap.Uint("id", id))
	return nil
}

func (s *RecipientService) AddTags(recipientID uint, tagIDs []uint) error {
	recipient, err := s.GetByID(recipientID)
	if err != nil {
		return err
	}

	var tags []models.Tag
	if err := database.DB.Where("id IN (?)", tagIDs).Find(&tags).Error; err != nil {
		return errors.DatabaseError(err)
	}

	if err := database.DB.Model(recipient).Association("Tags").Append(&tags); err != nil {
		logger.Error("add tags to recipient failed", zap.Error(err))
		return errors.DatabaseError(err)
	}

	logger.Info("tags added to recipient", zap.Uint("recipient_id", recipientID), zap.Any("tag_ids", tagIDs))
	return nil
}

func (s *RecipientService) RemoveTags(recipientID uint, tagIDs []uint) error {
	recipient, err := s.GetByID(recipientID)
	if err != nil {
		return err
	}

	var tags []models.Tag
	if err := database.DB.Where("id IN (?)", tagIDs).Find(&tags).Error; err != nil {
		return errors.DatabaseError(err)
	}

	if err := database.DB.Model(recipient).Association("Tags").Delete(&tags); err != nil {
		logger.Error("remove tags from recipient failed", zap.Error(err))
		return errors.DatabaseError(err)
	}

	logger.Info("tags removed from recipient", zap.Uint("recipient_id", recipientID), zap.Any("tag_ids", tagIDs))
	return nil
}

func (s *RecipientService) AddToGroups(recipientID uint, groupIDs []uint) error {
	recipient, err := s.GetByID(recipientID)
	if err != nil {
		return err
	}

	var groups []models.RecipientGroup
	if err := database.DB.Where("id IN (?)", groupIDs).Find(&groups).Error; err != nil {
		return errors.DatabaseError(err)
	}

	if err := database.DB.Model(recipient).Association("Groups").Append(&groups); err != nil {
		logger.Error("add recipient to groups failed", zap.Error(err))
		return errors.DatabaseError(err)
	}

	logger.Info("recipient added to groups", zap.Uint("recipient_id", recipientID), zap.Any("group_ids", groupIDs))
	return nil
}

func (s *RecipientService) RemoveFromGroups(recipientID uint, groupIDs []uint) error {
	recipient, err := s.GetByID(recipientID)
	if err != nil {
		return err
	}

	var groups []models.RecipientGroup
	if err := database.DB.Where("id IN (?)", groupIDs).Find(&groups).Error; err != nil {
		return errors.DatabaseError(err)
	}

	if err := database.DB.Model(recipient).Association("Groups").Delete(&groups); err != nil {
		logger.Error("remove recipient from groups failed", zap.Error(err))
		return errors.DatabaseError(err)
	}

	logger.Info("recipient removed from groups", zap.Uint("recipient_id", recipientID), zap.Any("group_ids", groupIDs))
	return nil
}

func (s *RecipientService) BatchImport(data [][]string) (int, int, []string) {
	success := 0
	failed := 0
	errors := make([]string, 0)

	for _, row := range data {
		if len(row) < 1 {
			continue
		}

		recipient := &models.Recipient{
			Email:  getColumn(row, 0),
			Phone:  getColumn(row, 1),
			Name:   getColumn(row, 2),
			WeChatOpenID: getColumn(row, 3),
			DingTalkUserID: getColumn(row, 4),
		}

		_, err := s.Create(recipient)
		if err != nil {
			failed++
			errors = append(errors, utils.TruncateString(err.Error(), 100))
		} else {
			success++
		}
	}

	logger.Info("batch import completed", zap.Int("success", success), zap.Int("failed", failed))
	return success, failed, errors
}

func (s *RecipientService) Export(groupIDs, tagIDs []uint) ([][]string, error) {
	recipients, _, err := s.List("", tagIDs, groupIDs, false, 1, 100000)
	if err != nil {
		return nil, err
	}

	data := make([][]string, 0, len(recipients)+1)
	data = append(data, []string{"email", "phone", "name", "wechat_open_id", "dingtalk_user_id"})

	for _, r := range recipients {
		data = append(data, []string{
			r.Email, r.Phone, r.Name, r.WeChatOpenID, r.DingTalkUserID,
		})
	}

	return data, nil
}

func (s *RecipientService) CreateTag(tag *models.Tag) (*models.Tag, error) {
	var existing models.Tag
	if err := database.DB.Where("name = ?", tag.Name).First(&existing).Error; err == nil {
		return &existing, nil
	}

	if err := database.DB.Create(tag).Error; err != nil {
		logger.Error("create tag failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	logger.Info("tag created", zap.Uint("id", tag.ID), zap.String("name", tag.Name))
	return tag, nil
}

func (s *RecipientService) ListTags() ([]models.Tag, error) {
	var tags []models.Tag
	if err := database.DB.Order("name ASC").Find(&tags).Error; err != nil {
		return nil, errors.DatabaseError(err)
	}
	return tags, nil
}

func (s *RecipientService) DeleteTag(id uint) error {
	if err := database.DB.Delete(&models.Tag{}, id).Error; err != nil {
		return errors.DatabaseError(err)
	}
	return nil
}

func (s *RecipientService) CreateGroup(group *models.RecipientGroup) (*models.RecipientGroup, error) {
	if err := database.DB.Create(group).Error; err != nil {
		logger.Error("create group failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	logger.Info("group created", zap.Uint("id", group.ID), zap.String("name", group.Name))
	return group, nil
}

func (s *RecipientService) ListGroups() ([]models.RecipientGroup, error) {
	var groups []models.RecipientGroup
	if err := database.DB.Order("name ASC").Find(&groups).Error; err != nil {
		return nil, errors.DatabaseError(err)
	}

	for i := range groups {
		var count int64
		database.DB.Table("recipient_group_members").
			Where("recipient_group_id = ?", groups[i].ID).
			Count(&count)
		groups[i].Count = int(count)
	}

	return groups, nil
}

func (s *RecipientService) GetGroup(id uint) (*models.RecipientGroup, error) {
	var group models.RecipientGroup
	if err := database.DB.Preload("Recipients").First(&group, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("group", err)
		}
		return nil, errors.DatabaseError(err)
	}

	var count int64
	database.DB.Table("recipient_group_members").
		Where("recipient_group_id = ?", group.ID).
		Count(&count)
	group.Count = int(count)

	return &group, nil
}

func (s *RecipientService) DeleteGroup(id uint) error {
	if err := database.DB.Delete(&models.RecipientGroup{}, id).Error; err != nil {
		return errors.DatabaseError(err)
	}
	return nil
}

func getColumn(row []string, index int) string {
	if index < len(row) {
		return strings.TrimSpace(row[index])
	}
	return ""
}
