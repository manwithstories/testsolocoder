package service

import (
	"car-rental/internal/model"
	"car-rental/internal/repository"
	"errors"
)

type StoreService struct {
	storeRepo *repository.StoreRepository
}

func NewStoreService() *StoreService {
	return &StoreService{
		storeRepo: repository.NewStoreRepository(),
	}
}

type CreateCityRequest struct {
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code"`
	Province string `json:"province"`
}

type CreateStoreRequest struct {
	Name          string  `json:"name" binding:"required"`
	CityID        uint    `json:"city_id" binding:"required"`
	Address       string  `json:"address" binding:"required"`
	Phone         string  `json:"phone"`
	BusinessHours string  `json:"business_hours"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
}

func (s *StoreService) CreateCity(req *CreateCityRequest) (*model.City, error) {
	if s.storeRepo.ExistsCityByName(req.Name) {
		return nil, errors.New("城市名称已存在")
	}

	city := &model.City{
		Name:     req.Name,
		Code:     req.Code,
		Province: req.Province,
	}

	err := s.storeRepo.CreateCity(city)
	if err != nil {
		return nil, err
	}

	return city, nil
}

func (s *StoreService) GetCityByID(id uint) (*model.City, error) {
	return s.storeRepo.FindCityByID(id)
}

func (s *StoreService) GetAllCities() ([]model.City, error) {
	return s.storeRepo.FindAllCities()
}

func (s *StoreService) UpdateCity(id uint, updates map[string]interface{}) error {
	city, err := s.storeRepo.FindCityByID(id)
	if err != nil {
		return err
	}

	if name, ok := updates["name"]; ok {
		if name.(string) != city.Name && s.storeRepo.ExistsCityByName(name.(string)) {
			return errors.New("城市名称已存在")
		}
		city.Name = name.(string)
	}
	if code, ok := updates["code"]; ok {
		city.Code = code.(string)
	}
	if province, ok := updates["province"]; ok {
		city.Province = province.(string)
	}

	return s.storeRepo.UpdateCity(city)
}

func (s *StoreService) DeleteCity(id uint) error {
	return s.storeRepo.DeleteCity(id)
}

func (s *StoreService) CreateStore(req *CreateStoreRequest) (*model.Store, error) {
	_, err := s.storeRepo.FindCityByID(req.CityID)
	if err != nil {
		return nil, errors.New("城市不存在")
	}

	store := &model.Store{
		Name:          req.Name,
		CityID:        req.CityID,
		Address:       req.Address,
		Phone:         req.Phone,
		BusinessHours: req.BusinessHours,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		IsActive:      true,
	}

	err = s.storeRepo.CreateStore(store)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *StoreService) GetStoreByID(id uint) (*model.Store, error) {
	return s.storeRepo.FindStoreByID(id)
}

func (s *StoreService) GetAllStores(page, pageSize int, cityID uint, keyword string) ([]model.Store, int64, error) {
	return s.storeRepo.FindAllStores(page, pageSize, cityID, keyword)
}

func (s *StoreService) UpdateStore(id uint, updates map[string]interface{}) error {
	store, err := s.storeRepo.FindStoreByID(id)
	if err != nil {
		return err
	}

	if name, ok := updates["name"]; ok {
		store.Name = name.(string)
	}
	if cityID, ok := updates["city_id"]; ok {
		store.CityID = uint(cityID.(float64))
	}
	if address, ok := updates["address"]; ok {
		store.Address = address.(string)
	}
	if phone, ok := updates["phone"]; ok {
		store.Phone = phone.(string)
	}
	if businessHours, ok := updates["business_hours"]; ok {
		store.BusinessHours = businessHours.(string)
	}
	if latitude, ok := updates["latitude"]; ok {
		store.Latitude = latitude.(float64)
	}
	if longitude, ok := updates["longitude"]; ok {
		store.Longitude = longitude.(float64)
	}
	if isActive, ok := updates["is_active"]; ok {
		store.IsActive = isActive.(bool)
	}

	return s.storeRepo.UpdateStore(store)
}

func (s *StoreService) DeleteStore(id uint) error {
	return s.storeRepo.DeleteStore(id)
}

func (s *StoreService) GetStoresByCity(cityID uint) ([]model.Store, error) {
	return s.storeRepo.GetStoresByCity(cityID)
}