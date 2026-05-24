package repository

import (
	"errors"
	"survey-platform/internal/model"

	"gorm.io/gorm"
)

type SurveyRepository struct {
	db *gorm.DB
}

func NewSurveyRepository(db *gorm.DB) *SurveyRepository {
	return &SurveyRepository{db: db}
}

func (r *SurveyRepository) Create(survey *model.Survey) error {
	return r.db.Create(survey).Error
}

func (r *SurveyRepository) FindByID(id uint) (*model.Survey, error) {
	var survey model.Survey
	err := r.db.Preload("Questions.Options").Preload("Questions.LogicJumps").
		Preload("User").First(&survey, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &survey, nil
}

func (r *SurveyRepository) Update(survey *model.Survey) error {
	return r.db.Save(survey).Error
}

func (r *SurveyRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("survey_id = ?", id).Delete(&model.Question{}).Error; err != nil {
			return err
		}
		if err := tx.Where("survey_id = ?", id).Delete(&model.Response{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.Survey{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *SurveyRepository) List(userID uint, page, pageSize int, status int, category, keyword, sortBy, sortOrder string) ([]*model.Survey, int64, error) {
	var surveys []*model.Survey
	var total int64

	query := r.db.Model(&model.Survey{}).Where("user_id = ?", userID)
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := sortBy + " " + sortOrder
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order(order).Find(&surveys).Error; err != nil {
		return nil, 0, err
	}

	return surveys, total, nil
}

func (r *SurveyRepository) ListAll(page, pageSize int, status int, category, keyword, sortBy, sortOrder string) ([]*model.Survey, int64, error) {
	var surveys []*model.Survey
	var total int64

	query := r.db.Preload("User").Model(&model.Survey{})
	if status > 0 {
		query = query.Where("status = ?", status)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := sortBy + " " + sortOrder
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order(order).Find(&surveys).Error; err != nil {
		return nil, 0, err
	}

	return surveys, total, nil
}

func (r *SurveyRepository) UpdateStatus(id uint, status int) error {
	return r.db.Model(&model.Survey{}).Where("id = ?", id).Update("status", status).Error
}

func (r *SurveyRepository) IncrementResponseCount(surveyID uint) error {
	return r.db.Model(&model.Survey{}).Where("id = ?", surveyID).
		UpdateColumn("response_count", gorm.Expr("response_count + 1")).Error
}

func (r *SurveyRepository) Copy(originalID uint, newTitle string, newUserID uint) (*model.Survey, error) {
	original, err := r.FindByID(originalID)
	if err != nil || original == nil {
		return nil, errors.New("original survey not found")
	}

	newSurvey := &model.Survey{
		Title:         newTitle,
		Description:   original.Description,
		CoverImage:    original.CoverImage,
		UserID:        newUserID,
		Status:        1,
		Anonymous:     original.Anonymous,
		Password:      original.Password,
		MaxResponses:  original.MaxResponses,
		MaxPerUser:    original.MaxPerUser,
		RequiresLogin: original.RequiresLogin,
		AllowResume:   original.AllowResume,
		Category:      original.Category,
		Tags:          original.Tags,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newSurvey).Error; err != nil {
			return err
		}

		for _, q := range original.Questions {
			newQ := &model.Question{
				SurveyID:       newSurvey.ID,
				Title:          q.Title,
				Type:           q.Type,
				IsRequired:     q.IsRequired,
				OrderIndex:     q.OrderIndex,
				Description:    q.Description,
				Placeholder:    q.Placeholder,
				MinValue:       q.MinValue,
				MaxValue:       q.MaxValue,
				DefaultValue:   q.DefaultValue,
				ValidationRule: q.ValidationRule,
				DisplayLogic:   q.DisplayLogic,
			}
			if err := tx.Create(newQ).Error; err != nil {
				return err
			}

			for _, opt := range q.Options {
				newOpt := &model.Option{
					QuestionID: newQ.ID,
					Text:       opt.Text,
					OrderIndex: opt.OrderIndex,
					IsOther:    opt.IsOther,
					JumpTarget: opt.JumpTarget,
					Score:      opt.Score,
				}
				if err := tx.Create(newOpt).Error; err != nil {
					return err
				}
			}

			for _, lj := range q.LogicJumps {
				newLJ := &model.LogicJump{
					QuestionID: newQ.ID,
					Condition:  lj.Condition,
					Value:      lj.Value,
					JumpTo:     lj.JumpTo,
				}
				if err := tx.Create(newLJ).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return newSurvey, nil
}
