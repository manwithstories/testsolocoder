package service

import (
	"errors"
	"math"
	"strings"
	"time"

	"matchmaking-platform/internal/dto"
	"matchmaking-platform/internal/model"
	"matchmaking-platform/internal/repository"
	"matchmaking-platform/internal/utils"
)

type MatchService struct {
	matchRepo   *repository.MatchRepo
	profileRepo *repository.ProfileRepo
	userRepo    *repository.UserRepo
	logRepo     *repository.SystemLogRepo
	memberRepo  *repository.MemberRepo
}

func NewMatchService() *MatchService {
	return &MatchService{
		matchRepo:   repository.NewMatchRepo(),
		profileRepo: repository.NewProfileRepo(),
		userRepo:    repository.NewUserRepo(),
		logRepo:     repository.NewSystemLogRepo(),
		memberRepo:  repository.NewMemberRepo(),
	}
}

func (s *MatchService) SmartMatch(userID uint, page, pageSize int) ([]dto.MatchResultItem, int64, error) {
	userProfile, err := s.profileRepo.FindByUserID(userID)
	if err != nil {
		return nil, 0, errors.New("请先完善个人资料")
	}

	blockedIDs, _ := s.matchRepo.GetBlockedIDs(userID)

	filter := map[string]interface{}{}
	if userProfile.PreferCity != "" {
		filter["city"] = userProfile.PreferCity
	}

	var targetGender string
	if userProfile.Gender == model.GenderMale {
		targetGender = string(model.GenderFemale)
	} else {
		targetGender = string(model.GenderMale)
	}
	filter["gender"] = targetGender

	profiles, total, err := s.profileRepo.ListByFilter(filter, blockedIDs, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var results []dto.MatchResultItem
	for _, p := range profiles {
		if p.UserID == userID {
			continue
		}
		if p.Age > 0 && userProfile.MinAge > 0 && p.Age < userProfile.MinAge {
			continue
		}
		if p.Age > 0 && userProfile.MaxAge > 0 && p.Age > userProfile.MaxAge {
			continue
		}
		if p.Height > 0 && userProfile.MinHeight > 0 && p.Height < userProfile.MinHeight {
			continue
		}
		if p.Height > 0 && userProfile.MaxHeight > 0 && p.Height > userProfile.MaxHeight {
			continue
		}

		score := s.calculateScore(userProfile, &p)
		reason := s.generateReason(userProfile, &p)

		record, _ := s.matchRepo.FindByUserAndTarget(userID, p.UserID)

		item := dto.MatchResultItem{
			UserID:      p.UserID,
			ProfileInfo: s.toProfileInfo(&p),
			MatchScore:  score,
			MatchReason: reason,
		}
		if record != nil {
			item.IsFavorited = record.IsFavorited
			item.IsBlocked = record.IsBlocked
		}
		results = append(results, item)

		s.matchRepo.Create(&model.MatchRecord{
			UserID:      userID,
			TargetID:    p.UserID,
			MatchScore:  score,
			MatchReason: reason,
		})
	}

	return results, total, nil
}

func (s *MatchService) FilterMatch(userID uint, req *dto.MatchFilterRequest) ([]dto.MatchResultItem, int64, error) {
	page, pageSize := utils.Paginate(req.Page, req.PageSize)

	filter := map[string]interface{}{}
	if req.Gender != "" {
		filter["gender"] = req.Gender
	}
	if req.Education != "" {
		filter["education"] = req.Education
	}
	if req.Income != "" {
		filter["income"] = req.Income
	}
	if req.City != "" {
		filter["city"] = req.City
	}

	blockedIDs, _ := s.matchRepo.GetBlockedIDs(userID)

	profiles, total, err := s.profileRepo.ListByFilter(filter, blockedIDs, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var results []dto.MatchResultItem
	for _, p := range profiles {
		if p.UserID == userID {
			continue
		}
		if req.MinAge > 0 && p.Age < req.MinAge {
			continue
		}
		if req.MaxAge > 0 && p.Age > req.MaxAge {
			continue
		}
		if req.MinHeight > 0 && p.Height < req.MinHeight {
			continue
		}
		if req.MaxHeight > 0 && p.Height > req.MaxHeight {
			continue
		}

		record, _ := s.matchRepo.FindByUserAndTarget(userID, p.UserID)

		item := dto.MatchResultItem{
			UserID:      p.UserID,
			ProfileInfo: s.toProfileInfo(&p),
		}
		if record != nil {
			item.IsFavorited = record.IsFavorited
			item.IsBlocked = record.IsBlocked
			item.MatchScore = record.MatchScore
		}
		results = append(results, item)
	}

	return results, total, nil
}

func (s *MatchService) Favorite(userID, targetID uint) error {
	record, err := s.matchRepo.FindByUserAndTarget(userID, targetID)
	if err != nil {
		return s.matchRepo.Create(&model.MatchRecord{
			UserID:      userID,
			TargetID:    targetID,
			IsFavorited: true,
		})
	}
	return s.matchRepo.Update(record.ID, map[string]interface{}{"is_favorited": !record.IsFavorited})
}

func (s *MatchService) Block(userID, targetID uint) error {
	record, err := s.matchRepo.FindByUserAndTarget(userID, targetID)
	if err != nil {
		return s.matchRepo.Create(&model.MatchRecord{
			UserID:    userID,
			TargetID:  targetID,
			IsBlocked: true,
		})
	}
	return s.matchRepo.Update(record.ID, map[string]interface{}{"is_blocked": !record.IsBlocked})
}

func (s *MatchService) GetFavorites(userID uint, page, pageSize int) ([]dto.MatchResultItem, int64, error) {
	records, total, err := s.matchRepo.GetFavorites(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var results []dto.MatchResultItem
	for _, r := range records {
		profile, err := s.profileRepo.FindByUserID(r.TargetID)
		if err != nil {
			continue
		}
		results = append(results, dto.MatchResultItem{
			UserID:      r.TargetID,
			ProfileInfo: s.toProfileInfo(profile),
			MatchScore:  r.MatchScore,
			IsFavorited: true,
		})
	}
	return results, total, nil
}

func (s *MatchService) GetBlocked(userID uint, page, pageSize int) ([]dto.MatchResultItem, int64, error) {
	records, total, err := s.matchRepo.GetBlocked(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var results []dto.MatchResultItem
	for _, r := range records {
		profile, err := s.profileRepo.FindByUserID(r.TargetID)
		if err != nil {
			continue
		}
		results = append(results, dto.MatchResultItem{
			UserID:      r.TargetID,
			ProfileInfo: s.toProfileInfo(profile),
			IsBlocked:   true,
		})
	}
	return results, total, nil
}

func (s *MatchService) calculateScore(user, target *model.Profile) float64 {
	var score float64 = 50

	if target.Age >= user.MinAge && target.Age <= user.MaxAge {
		score += 10
	}
	if target.Height >= user.MinHeight && target.Height <= user.MaxHeight {
		score += 10
	}
	if target.Education == user.PreferEducation {
		score += 10
	}
	if target.Income == user.PreferIncome {
		score += 10
	}
	if target.City == user.PreferCity {
		score += 10
	}

	userTags := strings.Split(user.Tags, ",")
	targetTags := strings.Split(target.Tags, ",")
	commonTags := 0
	for _, ut := range userTags {
		for _, tt := range targetTags {
			if strings.TrimSpace(ut) == strings.TrimSpace(tt) && ut != "" {
				commonTags++
			}
		}
	}
	if commonTags > 0 {
		score += float64(commonTags) * 2
	}

	userHobbies := strings.Split(user.Hobbies, ",")
	targetHobbies := strings.Split(target.Hobbies, ",")
	commonHobbies := 0
	for _, uh := range userHobbies {
		for _, th := range targetHobbies {
			if strings.TrimSpace(uh) == strings.TrimSpace(th) && uh != "" {
				commonHobbies++
			}
		}
	}
	if commonHobbies > 0 {
		score += float64(commonHobbies) * 1.5
	}

	score = math.Min(score, 100)
	return math.Round(score*10) / 10
}

func (s *MatchService) generateReason(user, target *model.Profile) string {
	var reasons []string

	if target.Age >= user.MinAge && target.Age <= user.MaxAge {
		reasons = append(reasons, "年龄符合您的要求")
	}
	if target.City == user.PreferCity {
		reasons = append(reasons, "同城匹配")
	}

	userTags := strings.Split(user.Tags, ",")
	targetTags := strings.Split(target.Tags, ",")
	var commonTags []string
	for _, ut := range userTags {
		for _, tt := range targetTags {
			if strings.TrimSpace(ut) == strings.TrimSpace(tt) && ut != "" {
				commonTags = append(commonTags, ut)
			}
		}
	}
	if len(commonTags) > 0 {
		reasons = append(reasons, "共同标签: "+strings.Join(commonTags, "、"))
	}

	userHobbies := strings.Split(user.Hobbies, ",")
	targetHobbies := strings.Split(target.Hobbies, ",")
	var commonHobbies []string
	for _, uh := range userHobbies {
		for _, th := range targetHobbies {
			if strings.TrimSpace(uh) == strings.TrimSpace(th) && uh != "" {
				commonHobbies = append(commonHobbies, uh)
			}
		}
	}
	if len(commonHobbies) > 0 {
		reasons = append(reasons, "共同爱好: "+strings.Join(commonHobbies, "、"))
	}

	if len(reasons) == 0 {
		return "系统智能匹配"
	}
	return strings.Join(reasons, "；")
}

func (s *MatchService) toProfileInfo(p *model.Profile) dto.ProfileInfo {
	return dto.ProfileInfo{
		ID:         p.ID,
		UserID:     p.UserID,
		Nickname:   p.Nickname,
		Gender:     string(p.Gender),
		Birthday:   p.Birthday,
		Age:        p.Age,
		Height:     p.Height,
		Weight:     p.Weight,
		Education:  string(p.Education),
		Occupation: p.Occupation,
		Income:     string(p.Income),
		City:       p.City,
		District:   p.District,
		Intro:      p.Intro,
		Hobbies:    p.Hobbies,
		Tags:       p.Tags,
		Photos:     strings.Split(p.Photos, ","),
	}
}

type DateService struct {
	dateRepo     *repository.DateRepo
	reviewRepo   *repository.DateReviewRepo
	memberRepo   *repository.MemberRepo
	logRepo      *repository.SystemLogRepo
	userService  *UserService
}

func NewDateService() *DateService {
	return &DateService{
		dateRepo:    repository.NewDateRepo(),
		reviewRepo:  repository.NewDateReviewRepo(),
		memberRepo:  repository.NewMemberRepo(),
		logRepo:     repository.NewSystemLogRepo(),
		userService: NewUserService(),
	}
}

func (s *DateService) CreateInvite(initiatorID uint, req *dto.DateInviteRequest) (*model.DateRecord, error) {
	benefit, _ := s.memberRepo.GetBenefit("free")
	user, _ := s.userService.userRepo.FindByID(initiatorID)
	if user.MemberLevel != "free" {
		benefit, _ = s.memberRepo.GetBenefit(user.MemberLevel)
	}

	todayCount, _ := s.memberRepo.CountTodayInteracts(initiatorID)
	if todayCount >= int64(benefit.DailyInteract) {
		return nil, errors.New("今日互动次数已达上限，请升级会员")
	}

	if req.DateAt.Before(time.Now()) {
		return nil, errors.New("约会时间不能早于当前时间")
	}

	record := &model.DateRecord{
		InitiatorID: initiatorID,
		ReceiverID:  req.ReceiverID,
		Title:       req.Title,
		Location:    req.Location,
		DateAt:      req.DateAt,
		Duration:    req.Duration,
		Status:      model.DateStatusPending,
		Note:        req.Note,
	}

	if err := s.dateRepo.Create(record); err != nil {
		return nil, err
	}

	s.memberRepo.LogInteract(&model.InteractLog{
		UserID:   initiatorID,
		TargetID: req.ReceiverID,
		Action:   "date_invite",
	})

	s.logRepo.Create(&model.SystemLog{
		UserID: initiatorID,
		Module: "date",
		Action: "invite",
	})

	return record, nil
}

func (s *DateService) Accept(dateID, userID uint) error {
	record, err := s.dateRepo.FindByID(dateID)
	if err != nil {
		return errors.New("约会记录不存在")
	}
	if record.ReceiverID != userID {
		return errors.New("无权操作")
	}
	if record.Status != model.DateStatusPending {
		return errors.New("当前状态无法接受")
	}
	return s.dateRepo.Update(dateID, map[string]interface{}{"status": model.DateStatusAccepted})
}

func (s *DateService) Reject(dateID, userID uint) error {
	record, err := s.dateRepo.FindByID(dateID)
	if err != nil {
		return errors.New("约会记录不存在")
	}
	if record.ReceiverID != userID && record.InitiatorID != userID {
		return errors.New("无权操作")
	}
	if record.Status != model.DateStatusPending {
		return errors.New("当前状态无法拒绝")
	}
	return s.dateRepo.Update(dateID, map[string]interface{}{"status": model.DateStatusRejected})
}

func (s *DateService) Cancel(dateID, userID uint) error {
	record, err := s.dateRepo.FindByID(dateID)
	if err != nil {
		return errors.New("约会记录不存在")
	}
	if record.InitiatorID != userID {
		return errors.New("只有发起人可以取消")
	}
	if record.Status != model.DateStatusAccepted {
		return errors.New("当前状态无法取消")
	}
	return s.dateRepo.Update(dateID, map[string]interface{}{"status": model.DateStatusCanceled})
}

func (s *DateService) Complete(dateID, userID uint) error {
	record, err := s.dateRepo.FindByID(dateID)
	if err != nil {
		return errors.New("约会记录不存在")
	}
	if record.InitiatorID != userID && record.ReceiverID != userID {
		return errors.New("无权操作")
	}
	if record.Status != model.DateStatusAccepted {
		return errors.New("当前状态无法完成")
	}
	return s.dateRepo.Update(dateID, map[string]interface{}{"status": model.DateStatusCompleted})
}

func (s *DateService) GetUserDates(userID uint, page, pageSize int) ([]model.DateRecord, int64, error) {
	return s.dateRepo.ListByUser(userID, page, pageSize)
}

func (s *DateService) CreateReview(userID uint, req *dto.DateReviewRequest) error {
	record, err := s.dateRepo.FindByID(req.DateID)
	if err != nil {
		return errors.New("约会记录不存在")
	}
	if record.Status != model.DateStatusCompleted {
		return errors.New("只有完成的约会可以评价")
	}
	if record.InitiatorID != userID && record.ReceiverID != userID {
		return errors.New("无权评价")
	}

	existing, _ := s.reviewRepo.FindByDateAndReviewer(req.DateID, userID)
	if existing != nil {
		return errors.New("您已经评价过了")
	}

	review := &model.DateReview{
		DateID:     req.DateID,
		ReviewerID: userID,
		TargetID:   req.TargetID,
		Rating:     req.Rating,
		Content:    req.Content,
	}

	return s.reviewRepo.Create(review)
}

func (s *DateService) GetReviews(userID uint, page, pageSize int) ([]model.DateReview, int64, error) {
	return s.reviewRepo.ListByUser(userID, page, pageSize)
}
