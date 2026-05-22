package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"event-platform/internal/database"
	applog "event-platform/internal/logger"
	"event-platform/internal/model"
	"event-platform/internal/queue"
	"event-platform/internal/repository"
	"event-platform/pkg/crypto"
	"event-platform/pkg/jwt"
)

type UserService struct {
	repo    *repository.UserRepo
	logRepo *repository.LogRepo
	jm      *jwt.Manager
}

func NewUserService(repo *repository.UserRepo, logRepo *repository.LogRepo, jm *jwt.Manager) *UserService {
	return &UserService{repo: repo, logRepo: logRepo, jm: jm}
}

func (s *UserService) Register(username, password, realName, idCard, phone, email string) (*model.User, error) {
	_, err := s.repo.GetByUsername(username)
	if err == nil {
		return nil, errors.New("username already exists")
	}
	hash, _ := crypto.HashPassword(password)
	u := &model.User{
		Username: username,
		Password: hash,
		RealName: realName,
		IdCard:   idCard,
		Phone:    phone,
		Email:    email,
		Role:     model.RoleUser,
		Status:   1,
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	applog.Infof("user registered: id=%d username=%s", u.ID, u.Username)
	return u, nil
}

func (s *UserService) Login(username, password string) (string, *model.User, error) {
	u, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid username or password")
	}
	if !crypto.CheckPassword(u.Password, password) {
		return "", nil, errors.New("invalid username or password")
	}
	if u.Status != 1 {
		return "", nil, errors.New("account disabled")
	}
	token, err := s.jm.Generate(u.ID, u.Username, string(u.Role))
	if err != nil {
		return "", nil, err
	}
	_ = s.repo.UpdateLastLogin(u.ID)
	return token, u, nil
}

func (s *UserService) Verify(userID uint, realName, idCard string) error {
	u, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}
	u.RealName = realName
	u.IdCard = idCard
	u.Verified = true
	return s.repo.Update(u)
}

func (s *UserService) GetByID(id uint) (*model.User, error) { return s.repo.GetByID(id) }

func (s *UserService) List(page, size int) ([]model.User, int64, error) { return s.repo.List(page, size) }

type EventService struct {
	repo    *repository.EventRepo
	logRepo *repository.LogRepo
	regRepo *repository.RegistrationRepo
}

func NewEventService(repo *repository.EventRepo, logRepo *repository.LogRepo, regRepo *repository.RegistrationRepo) *EventService {
	return &EventService{repo: repo, logRepo: logRepo, regRepo: regRepo}
}

func (s *EventService) Create(e *model.Event, items []model.EventItem) (*model.Event, error) {
	if err := s.repo.Create(e); err != nil {
		return nil, err
	}
	for i := range items {
		items[i].EventID = e.ID
		if err := s.repo.CreateItem(&items[i]); err != nil {
			return nil, err
		}
	}
	e.Items = items
	applog.Infof("event created: id=%d name=%s", e.ID, e.Name)
	return e, nil
}

func (s *EventService) Update(e *model.Event) error {
	return s.repo.Update(e)
}

func (s *EventService) Publish(id uint) error {
	e, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	e.IsPublished = true
	e.Status = "published"
	return s.repo.Update(e)
}

func (s *EventService) Unpublish(id uint) error {
	e, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	e.IsPublished = false
	e.Status = "draft"
	return s.repo.Update(e)
}

func (s *EventService) Get(id uint) (*model.Event, error) { return s.repo.GetByID(id) }
func (s *EventService) Item(id uint) (*model.EventItem, error) { return s.repo.GetItem(id) }
func (s *EventService) ListPublished(page, size int) ([]model.Event, int64, error) {
	return s.repo.ListPublished(page, size)
}
func (s *EventService) ListAll(page, size int) ([]model.Event, int64, error) {
	return s.repo.ListAll(page, size)
}

type RegistrationService struct {
	repo       *repository.RegistrationRepo
	eventRepo  *repository.EventRepo
	userRepo   *repository.UserRepo
	msgQueue   *queue.Queue
	logRepo    *repository.LogRepo
}

func NewRegistrationService(
	repo *repository.RegistrationRepo,
	eventRepo *repository.EventRepo,
	userRepo *repository.UserRepo,
	msgQueue *queue.Queue,
	logRepo *repository.LogRepo,
) *RegistrationService {
	return &RegistrationService{repo: repo, eventRepo: eventRepo, userRepo: userRepo, msgQueue: msgQueue, logRepo: logRepo}
}

func calcAge(idCard string) int {
	if len(idCard) < 10 {
		return 0
	}
	year := idCard[6:10]
	month := idCard[10:12]
	day := idCard[12:14]
	birth, err := time.Parse("20060102", year+month+day)
	if err != nil {
		return 0
	}
	age := time.Now().Year() - birth.Year()
	if time.Now().YearDay() < birth.YearDay() {
		age--
	}
	return age
}

func (s *RegistrationService) checkEligibility(user *model.User, item *model.EventItem) error {
	if item.MinAge > 0 || item.MaxAge > 0 {
		age := calcAge(user.IdCard)
		if item.MinAge > 0 && age < item.MinAge {
			return fmt.Errorf("年龄不满足要求，最小%d岁", item.MinAge)
		}
		if item.MaxAge > 0 && age > item.MaxAge {
			return fmt.Errorf("年龄超过上限，最大%d岁", item.MaxAge)
		}
	}
	return nil
}

func (s *RegistrationService) Register(userID, eventItemID uint, regType string, teamName string, teamMembers string) (*model.Registration, error) {
	item, err := s.eventRepo.GetItem(eventItemID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}
	event, err := s.eventRepo.GetByID(item.EventID)
	if err != nil {
		return nil, errors.New("赛事不存在")
	}
	if time.Now().After(event.RegistrationDeadline) {
		return nil, errors.New("报名已截止")
	}
	if !event.IsPublished {
		return nil, errors.New("赛事未发布")
	}
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if err := s.checkEligibility(user, item); err != nil {
		return nil, err
	}

	ctx := context.Background()
	quotaKey := fmt.Sprintf("event:quota:%d", eventItemID)
	waitlistKey := fmt.Sprintf("event:waitlist:%d", eventItemID)

	confirmedCount, err := s.repo.CountByItem(eventItemID, []model.RegistrationStatus{model.RegStatusConfirmed})
	if err != nil {
		return nil, err
	}
	waitlistCount, err := s.repo.CountByItem(eventItemID, []model.RegistrationStatus{model.RegStatusWaitlist})
	if err != nil {
		return nil, err
	}

	status := model.RegStatusConfirmed
	queuePos := 0
	if confirmedCount >= int64(item.Quota) {
		if item.WaitlistQuota > 0 && waitlistCount >= int64(item.WaitlistQuota) {
			return nil, errors.New("名额已满，候补也已排满")
		}
		status = model.RegStatusWaitlist
		queuePos = int(waitlistCount) + 1
	}

	reg := &model.Registration{
		UserID:        userID,
		EventID:       event.ID,
		EventItemID:   eventItemID,
		RegType:       model.RegistrationType(regType),
		TeamName:      teamName,
		TeamMembers:   teamMembers,
		Status:        status,
		QueuePosition: queuePos,
		PaymentStatus: "unpaid",
		Amount:        item.Fee,
	}

	incrKey := quotaKey
	if status == model.RegStatusWaitlist {
		incrKey = waitlistKey
	}
	database.RDB.Incr(ctx, incrKey)
	database.RDB.Expire(ctx, incrKey, 24*time.Hour)

	if err := s.repo.Create(reg); err != nil {
		database.RDB.Decr(ctx, incrKey)
		return nil, err
	}

	title := "报名成功"
	content := fmt.Sprintf("您已成功报名【%s - %s】", event.Name, item.Name)
	if status == model.RegStatusWaitlist {
		title = "候补报名"
		content = fmt.Sprintf("【%s - %s】名额已满，您已进入候补队列第%d位", event.Name, item.Name, queuePos)
	}
	_ = s.msgQueue.PublishMessage(ctx, queue.MessagePayload{
		UserIDs: []uint{userID},
		Type:    model.MsgTypeRegSuccess,
		Title:   title,
		Content: content,
		Extra:   fmt.Sprintf(`{"registration_id":%d}`, reg.ID),
	})

	return reg, nil
}

func (s *RegistrationService) ConfirmWaitlist(regID uint) error {
	reg, err := s.repo.GetByID(regID)
	if err != nil {
		return err
	}
	if reg.Status != model.RegStatusWaitlist {
		return errors.New("非候补状态")
	}
	reg.Status = model.RegStatusConfirmed
	reg.QueuePosition = 0
	if err := s.repo.Update(reg); err != nil {
		return err
	}
	item, _ := s.eventRepo.GetItem(reg.EventItemID)
	event, _ := s.eventRepo.GetByID(reg.EventID)
	ctx := context.Background()
	_ = s.msgQueue.PublishMessage(ctx, queue.MessagePayload{
		UserIDs: []uint{reg.UserID},
		Type:    model.MsgTypeRegSuccess,
		Title:   "候补已确认",
		Content: fmt.Sprintf("您在【%s - %s】的候补名额已确认，可正常参赛", event.Name, item.Name),
	})
	return nil
}

func (s *RegistrationService) ListByUser(userID uint, page, size int) ([]model.Registration, int64, error) {
	return s.repo.ListByUser(userID, page, size)
}

func (s *RegistrationService) ListByEvent(eventID uint, page, size int) ([]model.Registration, int64, error) {
	return s.repo.ListByEvent(eventID, page, size)
}

type ScoreService struct {
	repo      *repository.ScoreRepo
	regRepo   *repository.RegistrationRepo
	eventRepo *repository.EventRepo
	certRepo  *repository.CertificateRepo
	msgQueue  *queue.Queue
	logRepo   *repository.LogRepo
}

func NewScoreService(
	repo *repository.ScoreRepo,
	regRepo *repository.RegistrationRepo,
	eventRepo *repository.EventRepo,
	certRepo *repository.CertificateRepo,
	msgQueue *queue.Queue,
	logRepo *repository.LogRepo,
) *ScoreService {
	return &ScoreService{repo: repo, regRepo: regRepo, eventRepo: eventRepo, certRepo: certRepo, msgQueue: msgQueue, logRepo: logRepo}
}

func (s *ScoreService) calcRank(list []model.Score) {
	for i := range list {
		list[i].Rank = i + 1
		if list[i].Rank <= 3 {
			list[i].Points = float64(4 - list[i].Rank)
		} else if list[i].Rank <= 10 {
			list[i].Points = 1.0
		}
	}
}

func (s *ScoreService) Entry(list []model.Score) error {
	for i := range list {
		existing, _ := s.repo.GetByUserAndItem(list[i].UserID, list[i].EventItemID)
		if existing != nil {
			return fmt.Errorf("user %d already has score for item %d", list[i].UserID, list[i].EventItemID)
		}
	}
	var ptrs []*model.Score
	for i := range list {
		ptrs = append(ptrs, &list[i])
	}
	if err := s.repo.BatchCreate(ptrs); err != nil {
		return err
	}
	itemID := list[0].EventItemID
	all, _ := s.repo.ListByItem(itemID)
	s.calcRank(all)
	for _, sc := range all {
		_ = s.repo.Update(&sc)
	}
	ctx := context.Background()
	var userIDs []uint
	for _, sc := range list {
		userIDs = append(userIDs, sc.UserID)
	}
	item, _ := s.eventRepo.GetItem(itemID)
	_ = s.msgQueue.PublishMessage(ctx, queue.MessagePayload{
		UserIDs: userIDs,
		Type:    model.MsgTypeScorePublished,
		Title:   "成绩已发布",
		Content: fmt.Sprintf("【%s】成绩已发布，请查看", item.Name),
	})
	return nil
}

func (s *ScoreService) ListByItem(itemID uint) ([]model.Score, error) {
	return s.repo.ListByItem(itemID)
}

func (s *ScoreService) ListByUser(userID uint, page, size int) ([]model.Score, int64, error) {
	return s.repo.ListByUser(userID, page, size)
}

type CertificateService struct {
	repo       *repository.CertificateRepo
	scoreRepo  *repository.ScoreRepo
	eventRepo  *repository.EventRepo
	userRepo   *repository.UserRepo
	msgQueue   *queue.Queue
	certDir    string
}

func NewCertificateService(
	repo *repository.CertificateRepo,
	scoreRepo *repository.ScoreRepo,
	eventRepo *repository.EventRepo,
	userRepo *repository.UserRepo,
	msgQueue *queue.Queue,
	certDir string,
) *CertificateService {
	return &CertificateService{repo: repo, scoreRepo: scoreRepo, eventRepo: eventRepo, userRepo: userRepo, msgQueue: msgQueue, certDir: certDir}
}

func (s *CertificateService) Generate(certID uint) error {
	return nil
}

func (s *CertificateService) Get(id uint) (*model.Certificate, error) { return s.repo.GetByID(id) }
func (s *CertificateService) Update(c *model.Certificate) error       { return s.repo.Update(c) }

func (s *CertificateService) ListByUser(userID uint) ([]model.Certificate, error) {
	return s.repo.ListByUser(userID)
}

type MessageService struct {
	repo     *repository.MessageRepo
	msgQueue *queue.Queue
}

func NewMessageService(repo *repository.MessageRepo, msgQueue *queue.Queue) *MessageService {
	return &MessageService{repo: repo, msgQueue: msgQueue}
}

func (s *MessageService) List(userID uint, page, size int) ([]model.Message, int64, error) {
	return s.repo.ListByUser(userID, page, size)
}
func (s *MessageService) MarkRead(id, userID uint) error { return s.repo.MarkRead(id, userID) }
func (s *MessageService) MarkAllRead(userID uint) error { return s.repo.MarkAllRead(userID) }
func (s *MessageService) UnreadCount(userID uint) (int64, error) { return s.repo.UnreadCount(userID) }
