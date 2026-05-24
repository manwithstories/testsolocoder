package repository

import (
	"errors"
	"survey-platform/internal/model"

	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) Create(question *model.Question) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(question).Error; err != nil {
			return err
		}

		for i := range question.Options {
			question.Options[i].QuestionID = question.ID
			if err := tx.Create(&question.Options[i]).Error; err != nil {
				return err
			}
		}

		for i := range question.LogicJumps {
			question.LogicJumps[i].QuestionID = question.ID
			if err := tx.Create(&question.LogicJumps[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *QuestionRepository) FindByID(id uint) (*model.Question, error) {
	var question model.Question
	err := r.db.Preload("Options").Preload("LogicJumps").First(&question, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &question, nil
}

func (r *QuestionRepository) FindBySurveyID(surveyID uint) ([]*model.Question, error) {
	var questions []*model.Question
	err := r.db.Preload("Options").Preload("LogicJumps").
		Where("survey_id = ? AND status = 1", surveyID).
		Order("order_index ASC").Find(&questions).Error
	return questions, err
}

func (r *QuestionRepository) Update(question *model.Question) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(question).Error; err != nil {
			return err
		}

		if err := tx.Where("question_id = ?", question.ID).Delete(&model.Option{}).Error; err != nil {
			return err
		}
		for i := range question.Options {
			question.Options[i].QuestionID = question.ID
			if err := tx.Create(&question.Options[i]).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("question_id = ?", question.ID).Delete(&model.LogicJump{}).Error; err != nil {
			return err
		}
		for i := range question.LogicJumps {
			question.LogicJumps[i].QuestionID = question.ID
			if err := tx.Create(&question.LogicJumps[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *QuestionRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("question_id = ?", id).Delete(&model.Answer{}).Error; err != nil {
			return err
		}
		if err := tx.Where("question_id = ?", id).Delete(&model.Option{}).Error; err != nil {
			return err
		}
		if err := tx.Where("question_id = ?", id).Delete(&model.LogicJump{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Question{}, id).Error
	})
}

func (r *QuestionRepository) UpdateOrderIndex(id uint, orderIndex int) error {
	return r.db.Model(&model.Question{}).Where("id = ?", id).
		Update("order_index", orderIndex).Error
}

func (r *QuestionRepository) SoftDelete(id uint) error {
	return r.db.Model(&model.Question{}).Where("id = ?", id).
		Update("status", 2).Error
}

func (r *QuestionRepository) BatchCreate(surveyID uint, questions []*model.Question) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("survey_id = ?", surveyID).Delete(&model.Question{}).Error; err != nil {
			return err
		}

		for _, q := range questions {
			q.SurveyID = surveyID
			if err := tx.Create(q).Error; err != nil {
				return err
			}

			for i := range q.Options {
				q.Options[i].QuestionID = q.ID
				if err := tx.Create(&q.Options[i]).Error; err != nil {
					return err
				}
			}

			for i := range q.LogicJumps {
				q.LogicJumps[i].QuestionID = q.ID
				if err := tx.Create(&q.LogicJumps[i]).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
