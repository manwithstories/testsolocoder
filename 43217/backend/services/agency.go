package services

import (
	"errors"
	"fmt"
	"time"

	"health-platform/models"
	"health-platform/repository"
	"health-platform/utils"
)

type AgencyService struct {
	agencyRepo    *repository.AgencyRepository
	packageRepo   *repository.PackageRepository
	timeSlotRepo  *repository.TimeSlotRepository
	userRepo      *repository.UserRepository
}

func NewAgencyService() *AgencyService {
	return &AgencyService{
		agencyRepo:   repository.NewAgencyRepository(),
		packageRepo:  repository.NewPackageRepository(),
		timeSlotRepo: repository.NewTimeSlotRepository(),
		userRepo:     repository.NewUserRepository(),
	}
}

type RegisterAgencyRequest struct {
	Name            string `json:"name" binding:"required,max=100"`
	UnifiedCode     string `json:"unified_code" binding:"required,max=50"`
	LegalPerson     string `json:"legal_person" binding:"required,max=50"`
	ContactPhone    string `json:"contact_phone" binding:"required,max=20"`
	ContactEmail    string `json:"contact_email" binding:"omitempty,email,max=100"`
	Address         string `json:"address" binding:"max=255"`
	Description     string `json:"description"`
	AdminUsername   string `json:"admin_username" binding:"required,min=3,max=50"`
	AdminPassword   string `json:"admin_password" binding:"required,min=6,max=50"`
	AdminRealName   string `json:"admin_real_name" binding:"required,max=50"`
	AdminPhone      string `json:"admin_phone" binding:"required,max=20"`
	AdminEmail      string `json:"admin_email" binding:"omitempty,email,max=100"`
}

type UpdateAgencyRequest struct {
	ContactPhone    string `json:"contact_phone" binding:"max=20"`
	ContactEmail    string `json:"contact_email" binding:"omitempty,email,max=100"`
	Address         string `json:"address" binding:"max=255"`
	Description     string `json:"description"`
}

type CreatePackageRequest struct {
	AgencyID      uint   `json:"agency_id" binding:"required"`
	Name          string `json:"name" binding:"required,max=100"`
	Description   string `json:"description"`
	OriginalPrice float64 `json:"original_price"`
	Price         float64 `json:"price" binding:"required"`
	SuitableFor   string `json:"suitable_for" binding:"max=50"`
	GenderLimit   int    `json:"gender_limit"`
	MinAge        int    `json:"min_age"`
	MaxAge        int    `json:"max_age"`
	Notes         string `json:"notes"`
	Items         []PackageItemRequest `json:"items"`
}

type PackageItemRequest struct {
	ItemName    string `json:"item_name" binding:"required,max=100"`
	ItemCode    string `json:"item_code" binding:"max=50"`
	Description string `json:"description"`
	Department  string `json:"department" binding:"max=50"`
	NormalRange string `json:"normal_range" binding:"max=255"`
	Unit        string `json:"unit" binding:"max=20"`
	IsRequired  bool   `json:"is_required"`
	SortOrder   int    `json:"sort_order"`
}

type UpdatePackageRequest struct {
	Name          string `json:"name" binding:"max=100"`
	Description   string `json:"description"`
	OriginalPrice float64 `json:"original_price"`
	Price         float64 `json:"price"`
	SuitableFor   string `json:"suitable_for" binding:"max=50"`
	GenderLimit   int    `json:"gender_limit"`
	MinAge        int    `json:"min_age"`
	MaxAge        int    `json:"max_age"`
	Notes         string `json:"notes"`
}

type UpdatePackageItemRequest struct {
	ItemName    string `json:"item_name" binding:"max=100"`
	Description string `json:"description"`
	NormalRange string `json:"normal_range" binding:"max=255"`
	Unit        string `json:"unit" binding:"max=20"`
	IsRequired  bool   `json:"is_required"`
	SortOrder   int    `json:"sort_order"`
}

type CreateTimeSlotRequest struct {
	PackageID uint   `json:"package_id" binding:"required"`
	Date      string `json:"date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Total     int    `json:"total" binding:"required"`
}

func (s *AgencyService) RegisterAgency(req *RegisterAgencyRequest) (*models.Agency, *models.User, error) {
	existingAgency, _ := s.agencyRepo.FindByUnifiedCode(req.UnifiedCode)
	if existingAgency != nil {
		return nil, nil, errors.New("该统一社会信用代码已注册")
	}

	existingName, _ := s.agencyRepo.FindByName(req.Name)
	if existingName != nil {
		return nil, nil, errors.New("该机构名称已注册")
	}

	agency := &models.Agency{
		Name:          req.Name,
		UnifiedCode:   req.UnifiedCode,
		LegalPerson:   req.LegalPerson,
		ContactPhone:  req.ContactPhone,
		ContactEmail:  req.ContactEmail,
		Address:       req.Address,
		Description:   req.Description,
		Status:        models.AgencyStatusActive,
		Rating:        5.0,
	}

	if err := s.agencyRepo.Create(agency); err != nil {
		return nil, nil, fmt.Errorf("创建机构失败: %w", err)
	}

	authService := NewAuthService()
	adminUser, err := authService.Register(&RegisterRequest{
		Username: req.AdminUsername,
		Password: req.AdminPassword,
		RealName: req.AdminRealName,
		Phone:    req.AdminPhone,
		Email:    req.AdminEmail,
		Role:     models.RoleAgency,
		AgencyID: &agency.ID,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("创建管理员账号失败: %w", err)
	}

	agency.AdminUserID = &adminUser.ID
	if err := s.agencyRepo.Save(agency); err != nil {
		return nil, nil, fmt.Errorf("更新机构信息失败: %w", err)
	}

	return agency, adminUser, nil
}

func (s *AgencyService) GetAgency(agencyID uint) (*models.Agency, error) {
	agency, err := s.agencyRepo.GetWithPackages(agencyID)
	if err != nil {
		return nil, err
	}
	return agency, nil
}

func (s *AgencyService) UpdateAgency(agencyID uint, req *UpdateAgencyRequest) error {
	var agency models.Agency
	if err := s.agencyRepo.FindByID(&agency, agencyID); err != nil {
		return err
	}

	if req.ContactPhone != "" {
		agency.ContactPhone = req.ContactPhone
	}
	if req.ContactEmail != "" {
		agency.ContactEmail = req.ContactEmail
	}
	if req.Address != "" {
		agency.Address = req.Address
	}
	agency.Description = req.Description

	return s.agencyRepo.Save(&agency)
}

func (s *AgencyService) ListAgencies() ([]models.Agency, error) {
	return s.agencyRepo.ListActiveAgencies()
}

func (s *AgencyService) CreatePackage(req *CreatePackageRequest) (*models.Package, error) {
	pkg := &models.Package{
		AgencyID:      req.AgencyID,
		Name:          req.Name,
		Description:   req.Description,
		OriginalPrice: req.OriginalPrice,
		Price:         req.Price,
		SuitableFor:   req.SuitableFor,
		GenderLimit:   req.GenderLimit,
		MinAge:        req.MinAge,
		MaxAge:        req.MaxAge,
		Notes:         req.Notes,
		Status:        models.PackageStatusOnline,
	}

	if err := s.packageRepo.Create(pkg); err != nil {
		return nil, fmt.Errorf("创建套餐失败: %w", err)
	}

	for i, item := range req.Items {
		pkgItem := &models.PackageItem{
			PackageID:   pkg.ID,
			ItemName:    item.ItemName,
			ItemCode:    item.ItemCode,
			Description: item.Description,
			Department:  item.Department,
			NormalRange: item.NormalRange,
			Unit:        item.Unit,
			IsRequired:  item.IsRequired,
			SortOrder:   i + 1,
		}
		if item.SortOrder > 0 {
			pkgItem.SortOrder = item.SortOrder
		}
		if err := s.packageRepo.Create(pkgItem); err != nil {
			return nil, fmt.Errorf("添加检查项失败: %w", err)
		}
	}

	return pkg, nil
}

func (s *AgencyService) GetPackage(packageID uint) (*models.Package, error) {
	s.packageRepo.IncrementViewCount(packageID)
	return s.packageRepo.GetWithItems(packageID)
}

func (s *AgencyService) UpdatePackage(packageID uint, req *UpdatePackageRequest) error {
	pkg, err := s.packageRepo.GetWithItems(packageID)
	if err != nil {
		return err
	}

	if req.Name != "" {
		pkg.Name = req.Name
	}
	pkg.Description = req.Description
	pkg.OriginalPrice = req.OriginalPrice
	pkg.Price = req.Price
	pkg.SuitableFor = req.SuitableFor
	pkg.GenderLimit = req.GenderLimit
	pkg.MinAge = req.MinAge
	pkg.MaxAge = req.MaxAge
	pkg.Notes = req.Notes

	return s.packageRepo.Save(pkg)
}

func (s *AgencyService) UpdatePackageItem(itemID uint, req *UpdatePackageItemRequest) error {
	var item models.PackageItem
	if err := s.packageRepo.FindByID(&item, itemID); err != nil {
		return err
	}

	if req.ItemName != "" {
		item.ItemName = req.ItemName
	}
	item.Description = req.Description
	item.NormalRange = req.NormalRange
	item.Unit = req.Unit
	item.IsRequired = req.IsRequired
	item.SortOrder = req.SortOrder

	return s.packageRepo.Save(&item)
}

func (s *AgencyService) UpdatePackagePrice(packageID uint, price float64) error {
	return s.packageRepo.UpdatePrice(packageID, price)
}

func (s *AgencyService) UpdatePackageStatus(packageID uint, status models.PackageStatus) error {
	return s.packageRepo.UpdateStatus(packageID, status)
}

func (s *AgencyService) GetAgencyPackages(agencyID uint, page, pageSize int) ([]models.Package, int64, error) {
	return s.packageRepo.FindByAgencyID(agencyID, page, pageSize)
}

func (s *AgencyService) ListOnlinePackages(page, pageSize int, keyword string) ([]models.Package, int64, error) {
	return s.packageRepo.ListOnlinePackages(page, pageSize, keyword)
}

func (s *AgencyService) GetHotPackages(limit int) ([]models.Package, error) {
	return s.packageRepo.GetHotPackages(limit)
}

func (s *AgencyService) CreateTimeSlot(req *CreateTimeSlotRequest) error {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return errors.New("日期格式错误")
	}

	timeSlot := &models.TimeSlot{
		PackageID: req.PackageID,
		Date:      date,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Total:     req.Total,
		Status:    1,
	}

	return s.timeSlotRepo.Create(timeSlot)
}

func (s *AgencyService) BatchCreateTimeSlots(packageID uint, dates []string, startTime, endTime string, total int) error {
	var timeSlots []models.TimeSlot
	for _, dateStr := range dates {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		timeSlots = append(timeSlots, models.TimeSlot{
			PackageID: packageID,
			Date:      date,
			StartTime: startTime,
			EndTime:   endTime,
			Total:     total,
			Status:    1,
		})
	}
	return s.timeSlotRepo.BatchCreate(timeSlots)
}

func (s *AgencyService) GetPackageTimeSlots(packageID uint, startDate, endDate *time.Time) ([]models.TimeSlot, error) {
	return s.timeSlotRepo.FindByPackageID(packageID, startDate, endDate)
}

func (s *AgencyService) GetAvailableTimeSlots(packageID uint) ([]models.TimeSlot, error) {
	return s.timeSlotRepo.FindAvailableByPackageID(packageID)
}
