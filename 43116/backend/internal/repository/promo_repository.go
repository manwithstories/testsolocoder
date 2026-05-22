package repository

import (
	"car-rental/internal/model"
	cachedb "car-rental/internal/config"

	"gorm.io/gorm"
)

type PromoCodeRepository struct {
	db *gorm.DB
}

func NewPromoCodeRepository() *PromoCodeRepository {
	return &PromoCodeRepository{db: cachedb.DB}
}

func (r *PromoCodeRepository) Create(promo *model.PromoCode) error {
	return r.db.Create(promo).Error
}

func (r *PromoCodeRepository) FindByID(id uint) (*model.PromoCode, error) {
	var promo model.PromoCode
	err := r.db.First(&promo, id).Error
	if err != nil {
		return nil, err
	}
	return &promo, nil
}

func (r *PromoCodeRepository) FindByCode(code string) (*model.PromoCode, error) {
	var promo model.PromoCode
	err := r.db.Where("code = ?", code).First(&promo).Error
	if err != nil {
		return nil, err
	}
	return &promo, nil
}

func (r *PromoCodeRepository) FindAll(page, pageSize int, keyword string, isActive *bool) ([]model.PromoCode, int64, error) {
	var promos []model.PromoCode
	var total int64

	query := r.db.Model(&model.PromoCode{})
	if keyword != "" {
		query = query.Where("code LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&promos).Error
	return promos, total, err
}

func (r *PromoCodeRepository) Update(promo *model.PromoCode) error {
	return r.db.Save(promo).Error
}

func (r *PromoCodeRepository) IncrementUsedCount(id uint) error {
	return r.db.Model(&model.PromoCode{}).Where("id = ?", id).UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}

func (r *PromoCodeRepository) Delete(id uint) error {
	return r.db.Delete(&model.PromoCode{}, id).Error
}

func (r *PromoCodeRepository) ExistsByCode(code string) bool {
	var count int64
	r.db.Model(&model.PromoCode{}).Where("code = ?", code).Count(&count)
	return count > 0
}

type PricingRuleRepository struct {
	db *gorm.DB
}

func NewPricingRuleRepository() *PricingRuleRepository {
	return &PricingRuleRepository{db: cachedb.DB}
}

func (r *PricingRuleRepository) Create(rule *model.PricingRule) error {
	return r.db.Create(rule).Error
}

func (r *PricingRuleRepository) FindByID(id uint) (*model.PricingRule, error) {
	var rule model.PricingRule
	err := r.db.First(&rule, id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *PricingRuleRepository) FindAll(page, pageSize int, ruleType string) ([]model.PricingRule, int64, error) {
	var rules []model.PricingRule
	var total int64

	query := r.db.Model(&model.PricingRule{})
	if ruleType != "" {
		query = query.Where("rule_type = ?", ruleType)
	}

	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("priority ASC, created_at DESC").Find(&rules).Error
	return rules, total, err
}

func (r *PricingRuleRepository) GetActiveRules() ([]model.PricingRule, error) {
	var rules []model.PricingRule
	err := r.db.Where("is_active = ?", true).Order("priority ASC").Find(&rules).Error
	return rules, err
}

func (r *PricingRuleRepository) Update(rule *model.PricingRule) error {
	return r.db.Save(rule).Error
}

func (r *PricingRuleRepository) Delete(id uint) error {
	return r.db.Delete(&model.PricingRule{}, id).Error
}

func (r *PricingRuleRepository) ToggleActive(id uint) error {
	var rule model.PricingRule
	if err := r.db.First(&rule, id).Error; err != nil {
		return err
	}
	rule.IsActive = !rule.IsActive
	return r.db.Save(&rule).Error
}