package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"pet-board/internal/config"
	"pet-board/internal/dto"
	"pet-board/internal/models"
	"pet-board/internal/repository"
	"pet-board/internal/utils"
)

type UserService struct {
	userRepo  *repository.UserRepository
	jwtConfig config.JWTConfig
	logger    *logrus.Logger
}

func NewUserService(userRepo *repository.UserRepository, jwtConfig config.JWTConfig, logger *logrus.Logger) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
		logger:    logger,
	}
}

func (s *UserService) Register(req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	existingUser, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingEmail, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if existingEmail != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userID := uuid.New()
	user := &models.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Role:     models.Role(req.Role),
		RealName: req.RealName,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if req.Role == string(models.RoleStore) {
		storeInfo := &models.StoreInfo{
			ID:        uuid.New(),
			UserID:    userID,
			StoreName: req.Username + "'s Store",
		}
		if err := s.userRepo.CreateStoreInfo(storeInfo); err != nil {
			return nil, fmt.Errorf("failed to create store info: %w", err)
		}
	} else if req.Role == string(models.RoleKeeper) {
		keeperInfo := &models.KeeperInfo{
			ID:       uuid.New(),
			UserID:   userID,
			NickName: req.Username,
		}
		if err := s.userRepo.CreateKeeperInfo(keeperInfo); err != nil {
			return nil, fmt.Errorf("failed to create keeper info: %w", err)
		}
	}

	token, err := utils.GenerateToken(userID, user.Username, string(user.Role), s.jwtConfig.Secret, s.jwtConfig.ExpireHour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Infof("User registered: %s (role: %s)", user.Username, user.Role)

	return &dto.LoginResponse{
		Token:     token,
		UserID:    userID,
		Username:  user.Username,
		Role:      string(user.Role),
		ExpiresAt: time.Now().Add(time.Duration(s.jwtConfig.ExpireHour) * time.Hour),
	}, nil
}

func (s *UserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var user *models.User
	var err error

	user, err = s.userRepo.GetByUsername(req.Account)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if user == nil {
		user, err = s.userRepo.GetByEmail(req.Account)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	if user.Status != "active" {
		return nil, errors.New("account is not active")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, string(user.Role), s.jwtConfig.Secret, s.jwtConfig.ExpireHour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Infof("User logged in: %s (role: %s)", user.Username, user.Role)

	return &dto.LoginResponse{
		Token:     token,
		UserID:    user.ID,
		Username:  user.Username,
		Role:      string(user.Role),
		ExpiresAt: time.Now().Add(time.Duration(s.jwtConfig.ExpireHour) * time.Hour),
	}, nil
}

func (s *UserService) GetProfile(userID uuid.UUID) (*dto.UserProfile, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	profile := &dto.UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      string(user.Role),
		AvatarURL: user.AvatarURL,
		RealName:  user.RealName,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}

	if user.StoreInfo != nil {
		profile.StoreInfo = &dto.StoreInfoDTO{
			ID:            user.StoreInfo.ID,
			StoreName:     user.StoreInfo.StoreName,
			Address:       user.StoreInfo.Address,
			LicenseNo:     user.StoreInfo.LicenseNo,
			BusinessHours: user.StoreInfo.BusinessHours,
			Description:   user.StoreInfo.Description,
			Rating:        user.StoreInfo.Rating,
			ReviewCount:   user.StoreInfo.ReviewCount,
		}
	}

	if user.KeeperInfo != nil {
		profile.KeeperInfo = &dto.KeeperInfoDTO{
			ID:            user.KeeperInfo.ID,
			NickName:      user.KeeperInfo.NickName,
			Experience:    user.KeeperInfo.Experience,
			Specialty:     user.KeeperInfo.Specialty,
			Rating:        user.KeeperInfo.Rating,
			ReviewCount:   user.KeeperInfo.ReviewCount,
			Certification: user.KeeperInfo.Certification,
			StoreID:       user.KeeperInfo.StoreID,
		}
	}

	return profile, nil
}

func (s *UserService) UpdateProfile(userID uuid.UUID, req *dto.UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.Phone = req.Phone
	user.AvatarURL = req.AvatarURL
	user.RealName = req.RealName

	return s.userRepo.Update(user)
}

func (s *UserService) UpdateStoreInfo(userID uuid.UUID, req *dto.UpdateStoreRequest) error {
	storeInfo, err := s.userRepo.GetStoreInfoByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get store info: %w", err)
	}
	if storeInfo == nil {
		return errors.New("store info not found")
	}

	storeInfo.StoreName = req.StoreName
	storeInfo.Address = req.Address
	storeInfo.LicenseNo = req.LicenseNo
	storeInfo.BusinessHours = req.BusinessHours
	storeInfo.Description = req.Description

	return s.userRepo.UpdateStoreInfo(storeInfo)
}

func (s *UserService) UpdateKeeperInfo(userID uuid.UUID, req *dto.UpdateKeeperRequest) error {
	keeperInfo, err := s.userRepo.GetKeeperInfoByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get keeper info: %w", err)
	}
	if keeperInfo == nil {
		return errors.New("keeper info not found")
	}

	keeperInfo.NickName = req.NickName
	keeperInfo.Experience = req.Experience
	keeperInfo.Specialty = req.Specialty
	keeperInfo.Certification = req.Certification

	if req.StoreID != "" {
		storeID, err := uuid.Parse(req.StoreID)
		if err != nil {
			return errors.New("invalid store ID")
		}
		keeperInfo.StoreID = storeID
	}

	return s.userRepo.UpdateKeeperInfo(keeperInfo)
}

func (s *UserService) ChangePassword(userID uuid.UUID, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

func (s *UserService) ListUsers(page, pageSize int, role string) ([]models.User, int64, error) {
	return s.userRepo.List(page, pageSize, role)
}

type PetService struct {
	petRepo *repository.PetRepository
	logger  *logrus.Logger
}

func NewPetService(petRepo *repository.PetRepository, logger *logrus.Logger) *PetService {
	return &PetService{
		petRepo: petRepo,
		logger:  logger,
	}
}

func (s *PetService) Create(ownerID uuid.UUID, req *dto.PetRequest) (*models.Pet, error) {
	pet := &models.Pet{
		ID:          uuid.New(),
		OwnerID:     ownerID,
		Name:        req.Name,
		Species:     req.Species,
		Breed:       req.Breed,
		Gender:      req.Gender,
		BirthDate:   req.BirthDate,
		Weight:      req.Weight,
		Color:       req.Color,
		AvatarURL:   req.AvatarURL,
		Allergies:   req.Allergies,
		DietHabit:   req.DietHabit,
		Temperament: req.Temperament,
	}

	if err := s.petRepo.Create(pet); err != nil {
		return nil, fmt.Errorf("failed to create pet: %w", err)
	}

	s.logger.Infof("Pet created: %s (owner: %s)", pet.Name, ownerID)
	return pet, nil
}

func (s *PetService) GetByID(petID uuid.UUID) (*models.Pet, error) {
	pet, err := s.petRepo.GetByID(petID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pet: %w", err)
	}
	if pet == nil {
		return nil, errors.New("pet not found")
	}
	return pet, nil
}

func (s *PetService) ListByOwner(ownerID uuid.UUID, page, pageSize int) ([]models.Pet, int64, error) {
	return s.petRepo.ListByOwner(ownerID, page, pageSize)
}

func (s *PetService) Update(petID, ownerID uuid.UUID, req *dto.PetRequest) error {
	pet, err := s.petRepo.GetByID(petID)
	if err != nil {
		return fmt.Errorf("failed to get pet: %w", err)
	}
	if pet == nil {
		return errors.New("pet not found")
	}
	if pet.OwnerID != ownerID {
		return errors.New("not authorized to update this pet")
	}

	pet.Name = req.Name
	pet.Species = req.Species
	pet.Breed = req.Breed
	pet.Gender = req.Gender
	pet.BirthDate = req.BirthDate
	pet.Weight = req.Weight
	pet.Color = req.Color
	pet.AvatarURL = req.AvatarURL
	pet.Allergies = req.Allergies
	pet.DietHabit = req.DietHabit
	pet.Temperament = req.Temperament

	return s.petRepo.Update(pet)
}

func (s *PetService) Delete(petID, ownerID uuid.UUID) error {
	pet, err := s.petRepo.GetByID(petID)
	if err != nil {
		return fmt.Errorf("failed to get pet: %w", err)
	}
	if pet == nil {
		return errors.New("pet not found")
	}
	if pet.OwnerID != ownerID {
		return errors.New("not authorized to delete this pet")
	}

	return s.petRepo.Delete(petID)
}

func (s *PetService) AddVaccineRecord(req *dto.VaccineRecordRequest) (*models.VaccineRecord, error) {
	petID, err := uuid.Parse(req.PetID)
	if err != nil {
		return nil, errors.New("invalid pet ID")
	}

	record := &models.VaccineRecord{
		ID:           uuid.New(),
		PetID:        petID,
		VaccineName:  req.VaccineName,
		VaccinatedAt: req.VaccinatedAt,
		ExpireAt:     req.ExpireAt,
		Hospital:     req.Hospital,
		ProofURL:     req.ProofURL,
	}

	if err := s.petRepo.CreateVaccineRecord(record); err != nil {
		return nil, fmt.Errorf("failed to create vaccine record: %w", err)
	}

	s.logger.Infof("Vaccine record added for pet: %s", petID)
	return record, nil
}

func (s *PetService) GetVaccineRecords(petID uuid.UUID) ([]models.VaccineRecord, error) {
	return s.petRepo.GetVaccineRecords(petID)
}

func (s *PetService) AddDewormRecord(req *dto.DewormRecordRequest) (*models.DewormRecord, error) {
	petID, err := uuid.Parse(req.PetID)
	if err != nil {
		return nil, errors.New("invalid pet ID")
	}

	record := &models.DewormRecord{
		ID:         uuid.New(),
		PetID:      petID,
		DewormType: req.DewormType,
		DewormedAt: req.DewormedAt,
		ExpireAt:   req.ExpireAt,
		Medicine:   req.Medicine,
	}

	if err := s.petRepo.CreateDewormRecord(record); err != nil {
		return nil, fmt.Errorf("failed to create deworm record: %w", err)
	}

	s.logger.Infof("Deworm record added for pet: %s", petID)
	return record, nil
}

func (s *PetService) GetDewormRecords(petID uuid.UUID) ([]models.DewormRecord, error) {
	return s.petRepo.GetDewormRecords(petID)
}

func (s *PetService) HasValidVaccine(petID uuid.UUID) bool {
	return s.petRepo.HasValidVaccine(petID)
}

type PackageService struct {
	pkgRepo *repository.PackageRepository
	logger  *logrus.Logger
}

func NewPackageService(pkgRepo *repository.PackageRepository, logger *logrus.Logger) *PackageService {
	return &PackageService{
		pkgRepo: pkgRepo,
		logger:  logger,
	}
}

func (s *PackageService) Create(storeID uuid.UUID, req *dto.BoardingPackageRequest) (*models.BoardingPackage, error) {
	pkg := &models.BoardingPackage{
		ID:           uuid.New(),
		StoreID:      storeID,
		Name:         req.Name,
		Type:         req.Type,
		Description:  req.Description,
		PricePerDay:  req.PricePerDay,
		Capacity:     req.Capacity,
		Features:     req.Features,
		SortOrder:    req.SortOrder,
		IsAvailable:  true,
	}

	if err := s.pkgRepo.Create(pkg); err != nil {
		return nil, fmt.Errorf("failed to create package: %w", err)
	}

	s.logger.Infof("Package created: %s (store: %s)", pkg.Name, storeID)
	return pkg, nil
}

func (s *PackageService) GetByID(id uuid.UUID) (*models.BoardingPackage, error) {
	pkg, err := s.pkgRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get package: %w", err)
	}
	if pkg == nil {
		return nil, errors.New("package not found")
	}
	return pkg, nil
}

func (s *PackageService) ListByStore(storeID uuid.UUID, pkgType string, page, pageSize int) ([]models.BoardingPackage, int64, error) {
	return s.pkgRepo.ListByStore(storeID, pkgType, page, pageSize)
}

func (s *PackageService) Update(id, storeID uuid.UUID, req *dto.BoardingPackageRequest) error {
	pkg, err := s.pkgRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get package: %w", err)
	}
	if pkg == nil {
		return errors.New("package not found")
	}
	if pkg.StoreID != storeID {
		return errors.New("not authorized to update this package")
	}

	pkg.Name = req.Name
	pkg.Type = req.Type
	pkg.Description = req.Description
	pkg.PricePerDay = req.PricePerDay
	pkg.Capacity = req.Capacity
	pkg.Features = req.Features
	pkg.SortOrder = req.SortOrder

	return s.pkgRepo.Update(pkg)
}

func (s *PackageService) Delete(id, storeID uuid.UUID) error {
	pkg, err := s.pkgRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get package: %w", err)
	}
	if pkg == nil {
		return errors.New("package not found")
	}
	if pkg.StoreID != storeID {
		return errors.New("not authorized to delete this package")
	}

	return s.pkgRepo.Delete(id)
}
