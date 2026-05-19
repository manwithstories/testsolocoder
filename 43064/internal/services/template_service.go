package services

import (
	"fmt"
	"regexp"

	"github.com/notification-center/internal/database"
	"github.com/notification-center/internal/errors"
	"github.com/notification-center/internal/logger"
	"github.com/notification-center/internal/models"
	"github.com/notification-center/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TemplateService struct {
	channelService *ChannelService
}

func NewTemplateService(channelService *ChannelService) *TemplateService {
	return &TemplateService{
		channelService: channelService,
	}
}

func (s *TemplateService) Create(template *models.Template) (*models.Template, error) {
	if _, err := s.channelService.GetByID(template.ChannelID); err != nil {
		return nil, err
	}

	template.Variables = s.ExtractVariables(template.Content)

	if template.IsDefault {
		database.DB.Model(&models.Template{}).
			Where("channel_id = ? AND language = ? AND id != ?", template.ChannelID, template.Language, 0).
			Update("is_default", false)
	}

	if err := database.DB.Create(template).Error; err != nil {
		logger.Error("create template failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}

	logger.Info("template created", zap.Uint("id", template.ID), zap.String("name", template.Name))
	return template, nil
}

func (s *TemplateService) ExtractVariables(content string) []models.TemplateVariable {
	varNames := utils.ParseTemplateVariables(content)
	variables := make([]models.TemplateVariable, 0, len(varNames))
	for _, name := range varNames {
		variables = append(variables, models.TemplateVariable{
			Name:     name,
			Type:     "string",
			Required: true,
		})
	}
	return variables
}

func (s *TemplateService) ValidateVariables(templateContent string, variables map[string]interface{}) error {
	requiredVars := utils.ParseTemplateVariables(templateContent)
	missing := make([]string, 0)

	for _, v := range requiredVars {
		if _, exists := variables[v]; !exists {
			missing = append(missing, v)
		}
	}

	if len(missing) > 0 {
		return errors.TemplateError(fmt.Sprintf("missing required variables: %v", missing), nil)
	}

	return nil
}

func (s *TemplateService) GetByID(id uint) (*models.Template, error) {
	var template models.Template
	if err := database.DB.Preload("Channel").First(&template, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.TemplateNotFound(id)
		}
		logger.Error("get template failed", zap.Error(err), zap.Uint("id", id))
		return nil, errors.DatabaseError(err)
	}
	return &template, nil
}

func (s *TemplateService) GetByChannelAndLanguage(channelID uint, language string) (*models.Template, error) {
	var template models.Template
	err := database.DB.Where("channel_id = ? AND language = ? AND is_default = ?", channelID, language, true).
		First(&template).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("template", err)
		}
		logger.Error("get template by channel and language failed", zap.Error(err))
		return nil, errors.DatabaseError(err)
	}
	return &template, nil
}

func (s *TemplateService) List(channelID uint, language string, page, pageSize int) ([]models.Template, int64, error) {
	var templates []models.Template
	var total int64

	query := database.DB.Model(&models.Template{}).Preload("Channel")

	if channelID > 0 {
		query = query.Where("channel_id = ?", channelID)
	}
	if language != "" {
		query = query.Where("language = ?", language)
	}

	if err := query.Count(&total).Error; err != nil {
		logger.Error("count templates failed", zap.Error(err))
		return nil, 0, errors.DatabaseError(err)
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&templates).Error; err != nil {
		logger.Error("list templates failed", zap.Error(err))
		return nil, 0, errors.DatabaseError(err)
	}

	return templates, total, nil
}

func (s *TemplateService) Update(id uint, updates map[string]interface{}) (*models.Template, error) {
	template, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if content, ok := updates["content"].(string); ok {
		updates["variables"] = s.ExtractVariables(content)
	}

	if isDefault, ok := updates["is_default"].(bool); ok && isDefault {
		database.DB.Model(&models.Template{}).
			Where("channel_id = ? AND language = ? AND id != ?", template.ChannelID, template.Language, id).
			Update("is_default", false)
	}

	if err := database.DB.Model(template).Updates(updates).Error; err != nil {
		logger.Error("update template failed", zap.Error(err), zap.Uint("id", id))
		return nil, errors.DatabaseError(err)
	}

	if err := database.DB.Preload("Channel").First(template, id).Error; err != nil {
		logger.Error("reload template failed", zap.Error(err), zap.Uint("id", id))
		return nil, errors.DatabaseError(err)
	}

	logger.Info("template updated", zap.Uint("id", id))
	return template, nil
}

func (s *TemplateService) Delete(id uint) error {
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err := database.DB.Delete(&models.Template{}, id).Error; err != nil {
		logger.Error("delete template failed", zap.Error(err), zap.Uint("id", id))
		return errors.DatabaseError(err)
	}

	logger.Info("template deleted", zap.Uint("id", id))
	return nil
}

func (s *TemplateService) Render(template *models.Template, variables map[string]interface{}) (string, string, error) {
	if err := s.ValidateVariables(template.Content, variables); err != nil {
		return "", "", err
	}

	subject, err := utils.RenderTemplate(template.Subject, variables)
	if err != nil {
		logger.Warn("render subject failed", zap.Error(err))
	}

	content, err := utils.RenderTemplate(template.Content, variables)
	if err != nil {
		return "", "", errors.TemplateError("render template failed", err)
	}

	return subject, content, nil
}

func (s *TemplateService) ValidateTemplateSyntax(content string) error {
	pattern := regexp.MustCompile(`\{\{[^}]+\}\}`)
	if !pattern.MatchString(content) {
		return nil
	}

	validPattern := regexp.MustCompile(`\{\{[a-zA-Z0-9_]+\}\}`)
	matches := pattern.FindAllString(content, -1)

	for _, match := range matches {
		if !validPattern.MatchString(match) {
			return errors.TemplateError(fmt.Sprintf("invalid variable syntax: %s", match), nil)
		}
	}

	return nil
}
