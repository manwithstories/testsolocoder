package repository

import (
	"event-platform/internal/database"
	"event-platform/internal/model"
	"time"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo { return &UserRepo{} }

func (r *UserRepo) Create(u *model.User) error {
	return database.DB.Create(u).Error
}

func (r *UserRepo) GetByID(id uint) (*model.User, error) {
	var u model.User
	if err := database.DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByUsername(uname string) (*model.User, error) {
	var u model.User
	if err := database.DB.Where("username=?", uname).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Update(u *model.User) error {
	return database.DB.Save(u).Error
}

func (r *UserRepo) UpdateLastLogin(id uint) error {
	now := time.Now()
	return database.DB.Model(&model.User{}).Where("id=?", id).Update("last_login_at", now).Error
}

func (r *UserRepo) List(page, size int) ([]model.User, int64, error) {
	var list []model.User
	var total int64
	db := database.DB.Model(&model.User{})
	db.Count(&total)
	if err := db.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

type EventRepo struct{}

func NewEventRepo() *EventRepo { return &EventRepo{} }

func (r *EventRepo) Create(e *model.Event) error { return database.DB.Create(e).Error }
func (r *EventRepo) Update(e *model.Event) error { return database.DB.Save(e).Error }
func (r *EventRepo) GetByID(id uint) (*model.Event, error) {
	var e model.Event
	if err := database.DB.Preload("Items").First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}
func (r *EventRepo) ListPublished(page, size int) ([]model.Event, int64, error) {
	var list []model.Event
	var total int64
	db := database.DB.Model(&model.Event{}).Where("is_published=?", true)
	db.Count(&total)
	if err := db.Order("id desc").Preload("Items").Offset((page-1)*size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *EventRepo) ListAll(page, size int) ([]model.Event, int64, error) {
	var list []model.Event
	var total int64
	db := database.DB.Model(&model.Event{})
	db.Count(&total)
	if err := db.Order("id desc").Preload("Items").Offset((page-1)*size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *EventRepo) DeleteItem(id uint) error {
	return database.DB.Delete(&model.EventItem{}, id).Error
}
func (r *EventRepo) CreateItem(item *model.EventItem) error {
	return database.DB.Create(item).Error
}
func (r *EventRepo) GetItem(id uint) (*model.EventItem, error) {
	var it model.EventItem
	if err := database.DB.First(&it, id).Error; err != nil {
		return nil, err
	}
	return &it, nil
}

type RegistrationRepo struct{}

func NewRegistrationRepo() *RegistrationRepo { return &RegistrationRepo{} }

func (r *RegistrationRepo) Create(reg *model.Registration) error { return database.DB.Create(reg).Error }
func (r *RegistrationRepo) Update(reg *model.Registration) error { return database.DB.Save(reg).Error }
func (r *RegistrationRepo) GetByID(id uint) (*model.Registration, error) {
	var reg model.Registration
	if err := database.DB.First(&reg, id).Error; err != nil {
		return nil, err
	}
	return &reg, nil
}
func (r *RegistrationRepo) GetByUserAndItem(userID, itemID uint) (*model.Registration, error) {
	var reg model.Registration
	err := database.DB.Where("user_id=? AND event_item_id=? AND status NOT IN ?",
		userID, itemID, []model.RegistrationStatus{model.RegStatusCancelled, model.RegStatusRejected}).
		First(&reg).Error
	if err != nil {
		return nil, err
	}
	return &reg, nil
}
func (r *RegistrationRepo) CountByItem(itemID uint, statuses []model.RegistrationStatus) (int64, error) {
	var n int64
	q := database.DB.Model(&model.Registration{}).Where("event_item_id=?", itemID)
	if len(statuses) > 0 {
		q = q.Where("status IN ?", statuses)
	}
	if err := q.Count(&n).Error; err != nil {
		return 0, err
	}
	return n, nil
}
func (r *RegistrationRepo) ListByUser(userID uint, page, size int) ([]model.Registration, int64, error) {
	var list []model.Registration
	var total int64
	db := database.DB.Model(&model.Registration{}).Where("user_id=?", userID)
	db.Count(&total)
	if err := db.Preload("Event").Preload("EventItem").Order("id desc").
		Offset((page-1)*size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *RegistrationRepo) ListByEvent(eventID uint, page, size int) ([]model.Registration, int64, error) {
	var list []model.Registration
	var total int64
	db := database.DB.Model(&model.Registration{}).Where("event_id=?", eventID)
	db.Count(&total)
	if err := db.Preload("User").Preload("EventItem").Order("id desc").
		Offset((page-1)*size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *RegistrationRepo) ListWaitlistByItem(itemID uint) ([]model.Registration, error) {
	var list []model.Registration
	if err := database.DB.Where("event_item_id=? AND status=?", itemID, model.RegStatusWaitlist).
		Order("queue_position asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

type ScoreRepo struct{}

func NewScoreRepo() *ScoreRepo { return &ScoreRepo{} }

func (r *ScoreRepo) Create(s *model.Score) error { return database.DB.Create(s).Error }
func (r *ScoreRepo) Update(s *model.Score) error { return database.DB.Save(s).Error }
func (r *ScoreRepo) GetByUserAndItem(userID, itemID uint) (*model.Score, error) {
	var s model.Score
	err := database.DB.Where("user_id=? AND event_item_id=?", userID, itemID).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}
func (r *ScoreRepo) ListByItem(itemID uint) ([]model.Score, error) {
	var list []model.Score
	if err := database.DB.Where("event_item_id=? AND is_valid=?", itemID, true).
		Order("score desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
func (r *ScoreRepo) ListByUser(userID uint, page, size int) ([]model.Score, int64, error) {
	var list []model.Score
	var total int64
	db := database.DB.Model(&model.Score{}).Where("user_id=?", userID)
	db.Count(&total)
	if err := db.Order("id desc").Offset((page-1)*size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *ScoreRepo) BatchCreate(list []*model.Score) error {
	return database.DB.Create(&list).Error
}

type CertificateRepo struct{}

func NewCertificateRepo() *CertificateRepo { return &CertificateRepo{} }

func (r *CertificateRepo) Create(c *model.Certificate) error { return database.DB.Create(c).Error }
func (r *CertificateRepo) Update(c *model.Certificate) error { return database.DB.Save(c).Error }
func (r *CertificateRepo) GetByID(id uint) (*model.Certificate, error) {
	var c model.Certificate
	if err := database.DB.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *CertificateRepo) GetByUserAndItem(userID, itemID uint) (*model.Certificate, error) {
	var c model.Certificate
	err := database.DB.Where("user_id=? AND event_item_id=?", userID, itemID).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *CertificateRepo) ListByUser(userID uint) ([]model.Certificate, error) {
	var list []model.Certificate
	if err := database.DB.Where("user_id=?", userID).Order("id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
func (r *CertificateRepo) ListByStatus(status string) ([]model.Certificate, error) {
	var list []model.Certificate
	if err := database.DB.Where("status=?", status).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

type MessageRepo struct{}

func NewMessageRepo() *MessageRepo { return &MessageRepo{} }

func (r *MessageRepo) Create(m *model.Message) error { return database.DB.Create(m).Error }
func (r *MessageRepo) BatchCreate(list []*model.Message) error {
	return database.DB.Create(&list).Error
}
func (r *MessageRepo) ListByUser(userID uint, page, size int) ([]model.Message, int64, error) {
	var list []model.Message
	var total int64
	db := database.DB.Model(&model.Message{}).Where("user_id=?", userID)
	db.Count(&total)
	if err := db.Order("id desc").Offset((page-1)*size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func (r *MessageRepo) MarkRead(id, userID uint) error {
	return database.DB.Model(&model.Message{}).Where("id=? AND user_id=?", id, userID).Update("is_read", true).Error
}
func (r *MessageRepo) MarkAllRead(userID uint) error {
	return database.DB.Model(&model.Message{}).Where("user_id=? AND is_read=?", userID, false).Update("is_read", true).Error
}
func (r *MessageRepo) UnreadCount(userID uint) (int64, error) {
	var n int64
	if err := database.DB.Model(&model.Message{}).Where("user_id=? AND is_read=?", userID, false).Count(&n).Error; err != nil {
		return 0, err
	}
	return n, nil
}

type LogRepo struct{}

func NewLogRepo() *LogRepo { return &LogRepo{} }

func (r *LogRepo) Create(l *model.OperationLog) error { return database.DB.Create(l).Error }
