package service

import (
	"errors"
	"fmt"
	"hotel-system/internal/dto"
	"hotel-system/internal/model"
	"hotel-system/internal/pkg/logger"
	"hotel-system/internal/repository"
	"math"
	"time"

	"gorm.io/gorm"
)

type MemberService interface {
	RegisterMember(req *dto.MemberRegisterRequest) (*model.Member, error)
	GetMember(id uint) (*model.Member, error)
	GetMemberByMemberNo(memberNo string) (*model.Member, error)
	GetMemberByPhone(phone string) (*model.Member, error)
	UpdateMember(id uint, req *dto.MemberUpdateRequest) (*model.Member, error)
	DeleteMember(id uint) error
	ListMembers(req *dto.MemberListRequest) ([]model.Member, int64, error)

	CreateMemberLevel(req *dto.MemberLevelCreateRequest) (*model.MemberLevel, error)
	GetMemberLevel(id uint) (*model.MemberLevel, error)
	UpdateMemberLevel(id uint, req *dto.MemberLevelUpdateRequest) (*model.MemberLevel, error)
	DeleteMemberLevel(id uint) error
	ListMemberLevels() ([]model.MemberLevel, int64, error)

	CalculateDiscount(memberID uint) (float64, error)
	CalculatePoints(memberID uint, amount float64) (int, error)
	UsePoints(req *dto.PointsUseRequest) (float64, error)
	AddPoints(memberID uint, amount float64, description string, orderNo string) (int, error)
	CheckAndUpgradeLevel(memberID uint) (*model.MemberLevel, error)

	GetMemberDiscount(memberID uint) (*dto.MemberDiscountResponse, error)
	RechargePoints(req *dto.PointsRechargeRequest) error
	GetMemberConsumptionHistory(memberID uint, page, pageSize int) ([]map[string]interface{}, int64, error)
}

type memberService struct {
	memberRepo      repository.MemberRepository
	memberLevelRepo repository.MemberLevelRepository
	db              *gorm.DB
}

func NewMemberService(memberRepo repository.MemberRepository, memberLevelRepo repository.MemberLevelRepository, db *gorm.DB) MemberService {
	return &memberService{
		memberRepo:      memberRepo,
		memberLevelRepo: memberLevelRepo,
		db:              db,
	}
}

func (s *memberService) generateMemberNo() (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")

	count, err := s.memberRepo.GetTodayMemberCount(dateStr)
	if err != nil {
		logger.Errorf("获取今日会员数量失败: %v", err)
		return "", err
	}

	seq := count + 1
	memberNo := fmt.Sprintf("MB%s%04d", dateStr, seq)

	return memberNo, nil
}

func (s *memberService) RegisterMember(req *dto.MemberRegisterRequest) (*model.Member, error) {
	existingMember, _ := s.memberRepo.GetByPhone(req.Phone)
	if existingMember != nil {
		logger.Warnf("手机号已注册: %s", req.Phone)
		return nil, errors.New("该手机号已注册会员")
	}

	memberNo, err := s.generateMemberNo()
	if err != nil {
		logger.Errorf("生成会员号失败: %v", err)
		return nil, errors.New("生成会员号失败")
	}

	defaultLevel, err := s.memberLevelRepo.GetLevelByPoints(0)
	if err != nil {
		logger.Errorf("获取默认会员等级失败: %v", err)
		return nil, errors.New("获取默认会员等级失败")
	}

	member := &model.Member{
		MemberNo: memberNo,
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		IDCard:   req.IDCard,
		LevelID:  defaultLevel.ID,
		Points:   0,
		Balance:  0,
		Status:   model.MemberStatusActive,
	}

	err = s.memberRepo.Create(member)
	if err != nil {
		logger.Errorf("创建会员失败: %v", err)
		return nil, errors.New("创建会员失败")
	}

	logger.Infof("会员注册成功: memberNo=%s, name=%s", member.MemberNo, member.Name)
	return member, nil
}

func (s *memberService) GetMember(id uint) (*model.Member, error) {
	member, err := s.memberRepo.GetByID(id)
	if err != nil {
		logger.Errorf("获取会员信息失败: id=%d, err=%v", id, err)
		return nil, errors.New("会员不存在")
	}
	return member, nil
}

func (s *memberService) GetMemberByMemberNo(memberNo string) (*model.Member, error) {
	member, err := s.memberRepo.GetByMemberNo(memberNo)
	if err != nil {
		logger.Errorf("获取会员信息失败: memberNo=%s, err=%v", memberNo, err)
		return nil, errors.New("会员不存在")
	}
	return member, nil
}

func (s *memberService) GetMemberByPhone(phone string) (*model.Member, error) {
	member, err := s.memberRepo.GetByPhone(phone)
	if err != nil {
		logger.Errorf("获取会员信息失败: phone=%s, err=%v", phone, err)
		return nil, errors.New("会员不存在")
	}
	return member, nil
}

func (s *memberService) UpdateMember(id uint, req *dto.MemberUpdateRequest) (*model.Member, error) {
	member, err := s.memberRepo.GetByID(id)
	if err != nil {
		logger.Errorf("更新会员失败，会员不存在: id=%d, err=%v", id, err)
		return nil, errors.New("会员不存在")
	}

	if req.Name != "" {
		member.Name = req.Name
	}
	if req.Phone != "" {
		member.Phone = req.Phone
	}
	if req.Email != "" {
		member.Email = req.Email
	}
	if req.IDCard != "" {
		member.IDCard = req.IDCard
	}
	if req.Status != "" {
		member.Status = req.Status
	}

	err = s.memberRepo.Update(member)
	if err != nil {
		logger.Errorf("更新会员失败: %v", err)
		return nil, errors.New("更新会员失败")
	}

	logger.Infof("会员更新成功: id=%d, memberNo=%s", member.ID, member.MemberNo)
	return member, nil
}

func (s *memberService) DeleteMember(id uint) error {
	_, err := s.memberRepo.GetByID(id)
	if err != nil {
		logger.Errorf("删除会员失败，会员不存在: id=%d, err=%v", id, err)
		return errors.New("会员不存在")
	}

	err = s.memberRepo.Delete(id)
	if err != nil {
		logger.Errorf("删除会员失败: %v", err)
		return errors.New("删除会员失败")
	}

	logger.Infof("会员删除成功: id=%d", id)
	return nil
}

func (s *memberService) ListMembers(req *dto.MemberListRequest) ([]model.Member, int64, error) {
	page := req.GetPage()
	pageSize := req.GetPageSize()

	members, total, err := s.memberRepo.List(page, pageSize, req.Name, req.Phone, req.LevelID, req.Status)
	if err != nil {
		logger.Errorf("获取会员列表失败: %v", err)
		return nil, 0, errors.New("获取会员列表失败")
	}

	return members, total, nil
}

func (s *memberService) CreateMemberLevel(req *dto.MemberLevelCreateRequest) (*model.MemberLevel, error) {
	level := &model.MemberLevel{
		Name:         req.Name,
		DiscountRate: req.DiscountRate,
		PointsRate:   req.PointsRate,
		MinPoints:    req.MinPoints,
		MaxPoints:    req.MaxPoints,
	}

	err := s.memberLevelRepo.Create(level)
	if err != nil {
		logger.Errorf("创建会员等级失败: %v", err)
		return nil, errors.New("创建会员等级失败")
	}

	logger.Infof("会员等级创建成功: name=%s", level.Name)
	return level, nil
}

func (s *memberService) GetMemberLevel(id uint) (*model.MemberLevel, error) {
	level, err := s.memberLevelRepo.GetByID(id)
	if err != nil {
		logger.Errorf("获取会员等级失败: id=%d, err=%v", id, err)
		return nil, errors.New("会员等级不存在")
	}
	return level, nil
}

func (s *memberService) UpdateMemberLevel(id uint, req *dto.MemberLevelUpdateRequest) (*model.MemberLevel, error) {
	level, err := s.memberLevelRepo.GetByID(id)
	if err != nil {
		logger.Errorf("更新会员等级失败，等级不存在: id=%d, err=%v", id, err)
		return nil, errors.New("会员等级不存在")
	}

	if req.Name != "" {
		level.Name = req.Name
	}
	if req.DiscountRate > 0 {
		level.DiscountRate = req.DiscountRate
	}
	if req.PointsRate > 0 {
		level.PointsRate = req.PointsRate
	}
	if req.MinPoints >= 0 {
		level.MinPoints = req.MinPoints
	}
	if req.MaxPoints >= 0 {
		level.MaxPoints = req.MaxPoints
	}

	err = s.memberLevelRepo.Update(level)
	if err != nil {
		logger.Errorf("更新会员等级失败: %v", err)
		return nil, errors.New("更新会员等级失败")
	}

	logger.Infof("会员等级更新成功: id=%d, name=%s", level.ID, level.Name)
	return level, nil
}

func (s *memberService) DeleteMemberLevel(id uint) error {
	_, err := s.memberLevelRepo.GetByID(id)
	if err != nil {
		logger.Errorf("删除会员等级失败，等级不存在: id=%d, err=%v", id, err)
		return errors.New("会员等级不存在")
	}

	err = s.memberLevelRepo.Delete(id)
	if err != nil {
		logger.Errorf("删除会员等级失败: %v", err)
		return errors.New("删除会员等级失败")
	}

	logger.Infof("会员等级删除成功: id=%d", id)
	return nil
}

func (s *memberService) ListMemberLevels() ([]model.MemberLevel, int64, error) {
	levels, total, err := s.memberLevelRepo.List()
	if err != nil {
		logger.Errorf("获取会员等级列表失败: %v", err)
		return nil, 0, errors.New("获取会员等级列表失败")
	}

	return levels, total, nil
}

func (s *memberService) CalculateDiscount(memberID uint) (float64, error) {
	member, err := s.memberRepo.GetByID(memberID)
	if err != nil {
		logger.Errorf("计算折扣失败，会员不存在: id=%d, err=%v", memberID, err)
		return 1.0, errors.New("会员不存在")
	}

	if member.Level == nil {
		return 1.0, nil
	}

	return member.Level.DiscountRate, nil
}

func (s *memberService) CalculatePoints(memberID uint, amount float64) (int, error) {
	member, err := s.memberRepo.GetByID(memberID)
	if err != nil {
		logger.Errorf("计算积分失败，会员不存在: id=%d, err=%v", memberID, err)
		return 0, errors.New("会员不存在")
	}

	pointsRate := 1.0
	if member.Level != nil {
		pointsRate = member.Level.PointsRate
	}

	points := int(math.Floor(amount * pointsRate))
	return points, nil
}

func (s *memberService) UsePoints(req *dto.PointsUseRequest) (float64, error) {
	member, err := s.memberRepo.GetByID(req.MemberID)
	if err != nil {
		logger.Errorf("使用积分失败，会员不存在: id=%d, err=%v", req.MemberID, err)
		return 0, errors.New("会员不存在")
	}

	if member.Points < req.Points {
		logger.Warnf("积分不足: memberID=%d, currentPoints=%d, needPoints=%d", req.MemberID, member.Points, req.Points)
		return 0, errors.New("积分不足")
	}

	cashValue := float64(req.Points) / 100.0

	description := req.Description
	if description == "" {
		description = fmt.Sprintf("使用%d积分抵扣%.2f元", req.Points, cashValue)
	}

	err = s.memberRepo.UpdatePoints(req.MemberID, -req.Points, description, model.PointsLogTypeUse, req.OrderNo)
	if err != nil {
		logger.Errorf("使用积分失败: %v", err)
		return 0, errors.New("使用积分失败")
	}

	logger.Infof("积分使用成功: memberID=%d, points=%d, cash=%.2f", req.MemberID, req.Points, cashValue)
	return cashValue, nil
}

func (s *memberService) AddPoints(memberID uint, amount float64, description string, orderNo string) (int, error) {
	member, err := s.memberRepo.GetByID(memberID)
	if err != nil {
		logger.Errorf("增加积分失败，会员不存在: id=%d, err=%v", memberID, err)
		return 0, errors.New("会员不存在")
	}

	pointsRate := 1.0
	if member.Level != nil {
		pointsRate = member.Level.PointsRate
	}

	points := int(math.Floor(amount * pointsRate))

	if points <= 0 {
		return 0, nil
	}

	desc := description
	if desc == "" {
		desc = fmt.Sprintf("消费%.2f元获得%d积分", amount, points)
	}

	err = s.memberRepo.UpdatePoints(memberID, points, desc, model.PointsLogTypeEarn, orderNo)
	if err != nil {
		logger.Errorf("增加积分失败: %v", err)
		return 0, errors.New("增加积分失败")
	}

	_, _ = s.CheckAndUpgradeLevel(memberID)

	logger.Infof("积分增加成功: memberID=%d, points=%d", memberID, points)
	return points, nil
}

func (s *memberService) CheckAndUpgradeLevel(memberID uint) (*model.MemberLevel, error) {
	member, err := s.memberRepo.GetByID(memberID)
	if err != nil {
		logger.Errorf("检查会员等级失败，会员不存在: id=%d, err=%v", memberID, err)
		return nil, errors.New("会员不存在")
	}

	newLevel, err := s.memberLevelRepo.GetLevelByPoints(member.Points)
	if err != nil {
		logger.Errorf("获取对应等级失败: points=%d, err=%v", member.Points, err)
		return nil, errors.New("获取会员等级失败")
	}

	if newLevel.ID != member.LevelID {
		oldLevelID := member.LevelID
		member.LevelID = newLevel.ID
		err = s.memberRepo.Update(member)
		if err != nil {
			logger.Errorf("升级会员等级失败: %v", err)
			return nil, errors.New("升级会员等级失败")
		}
		logger.Infof("会员等级变更: memberID=%d, oldLevel=%d, newLevel=%d", memberID, oldLevelID, newLevel.ID)
	}

	return newLevel, nil
}

func (s *memberService) GetMemberDiscount(memberID uint) (*dto.MemberDiscountResponse, error) {
	member, err := s.memberRepo.GetByID(memberID)
	if err != nil {
		logger.Errorf("获取会员折扣信息失败: id=%d, err=%v", memberID, err)
		return nil, errors.New("会员不存在")
	}

	levelName := ""
	discountRate := 1.0
	pointsRate := 1.0

	if member.Level != nil {
		levelName = member.Level.Name
		discountRate = member.Level.DiscountRate
		pointsRate = member.Level.PointsRate
	}

	return &dto.MemberDiscountResponse{
		MemberID:     member.ID,
		MemberNo:     member.MemberNo,
		Name:         member.Name,
		LevelID:      member.LevelID,
		LevelName:    levelName,
		DiscountRate: discountRate,
		Points:       member.Points,
		PointsRate:   pointsRate,
	}, nil
}

func (s *memberService) RechargePoints(req *dto.PointsRechargeRequest) error {
	_, err := s.memberRepo.GetByID(req.MemberID)
	if err != nil {
		logger.Errorf("充值积分失败，会员不存在: id=%d, err=%v", req.MemberID, err)
		return errors.New("会员不存在")
	}

	description := req.Description
	if description == "" {
		description = fmt.Sprintf("充值%d积分", req.Points)
	}

	err = s.memberRepo.UpdatePoints(req.MemberID, req.Points, description, model.PointsLogTypeRecharge, req.OrderNo)
	if err != nil {
		logger.Errorf("充值积分失败: %v", err)
		return errors.New("充值积分失败")
	}

	_, _ = s.CheckAndUpgradeLevel(req.MemberID)

	logger.Infof("积分充值成功: memberID=%d, points=%d", req.MemberID, req.Points)
	return nil
}

func (s *memberService) GetMemberConsumptionHistory(memberID uint, page, pageSize int) ([]map[string]interface{}, int64, error) {
	type ConsumptionRecord struct {
		ID           uint   `json:"id"`
		BookingNo    string `json:"bookingNo"`
		CheckInDate  string `json:"checkInDate"`
		CheckOutDate string `json:"checkOutDate"`
		TotalPrice   float64 `json:"totalPrice"`
		Status       string `json:"status"`
		Type         string `json:"type"`
		CreatedAt    string `json:"createdAt"`
	}

	var records []ConsumptionRecord
	var total int64

	query := s.db.Table("bookings").
		Select("bookings.id, bookings.booking_no, bookings.check_in_date, bookings.check_out_date, bookings.total_price, bookings.status, 'booking' as type, bookings.created_at").
		Where("bookings.member_id = ?", memberID).
		Where("bookings.deleted_at IS NULL")

	if err := query.Count(&total).Error; err != nil {
		logger.Errorf("查询消费记录数量失败: %v", err)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("bookings.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&records).Error; err != nil {
		logger.Errorf("查询消费记录失败: %v", err)
		return nil, 0, err
	}

	result := make([]map[string]interface{}, len(records))
	for i, r := range records {
		result[i] = map[string]interface{}{
			"id":           r.ID,
			"bookingNo":    r.BookingNo,
			"checkInDate":  r.CheckInDate,
			"checkOutDate": r.CheckOutDate,
			"totalPrice":   r.TotalPrice,
			"status":       r.Status,
			"type":         r.Type,
			"createdAt":    r.CreatedAt,
		}
	}

	return result, total, nil
}
