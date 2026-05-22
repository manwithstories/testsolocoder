package service

import (
	"car-rental/internal/model"
	"car-rental/internal/repository"
	"car-rental/internal/config"
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"car-rental/internal/utils"
)

type CarService struct {
	carRepo *repository.CarRepository
	cfg     *config.Config
}

func NewCarService(cfg *config.Config) *CarService {
	return &CarService{
		carRepo: repository.NewCarRepository(),
		cfg:     cfg,
	}
}

type CreateCarRequest struct {
	Brand        string  `json:"brand" binding:"required"`
	Model        string  `json:"model" binding:"required"`
	Year         int     `json:"year" binding:"required"`
	Seats        int     `json:"seats" binding:"required"`
	Transmission string  `json:"transmission" binding:"required"`
	FuelType     string  `json:"fuel_type"`
	DailyRent    float64 `json:"daily_rent" binding:"required"`
	Deposit      float64 `json:"deposit"`
	LicensePlate string  `json:"license_plate"`
	Color        string  `json:"color"`
	Mileage      int     `json:"mileage"`
	StoreID      uint    `json:"store_id" binding:"required"`
	Features     string  `json:"features"`
	Description  string  `json:"description"`
}

func (s *CarService) CreateCar(req *CreateCarRequest) (*model.Car, error) {
	if req.LicensePlate != "" && s.carRepo.ExistsByLicensePlate(req.LicensePlate) {
		return nil, errors.New("车牌号已存在")
	}

	car := &model.Car{
		Brand:        req.Brand,
		Model:        req.Model,
		Year:         req.Year,
		Seats:        req.Seats,
		Transmission: req.Transmission,
		FuelType:     req.FuelType,
		DailyRent:    req.DailyRent,
		Deposit:      req.Deposit,
		LicensePlate: req.LicensePlate,
		Color:        req.Color,
		Mileage:      req.Mileage,
		StoreID:      req.StoreID,
		Features:     req.Features,
		Description:  req.Description,
		Status:       model.CarStatusAvailable,
	}

	err := s.carRepo.Create(car)
	if err != nil {
		return nil, err
	}

	return car, nil
}

func (s *CarService) GetCarByID(id uint) (*model.Car, error) {
	return s.carRepo.FindByID(id)
}

func (s *CarService) GetAllCars(page, pageSize int, keyword, status, brand string, storeID uint) ([]model.Car, int64, error) {
	return s.carRepo.FindAll(page, pageSize, keyword, status, brand, storeID)
}

func (s *CarService) UpdateCar(id uint, updates map[string]interface{}) error {
	car, err := s.carRepo.FindByID(id)
	if err != nil {
		return err
	}

	if brand, ok := updates["brand"]; ok {
		car.Brand = brand.(string)
	}
	if modelName, ok := updates["model"]; ok {
		car.Model = modelName.(string)
	}
	if year, ok := updates["year"]; ok {
		car.Year = int(year.(float64))
	}
	if seats, ok := updates["seats"]; ok {
		car.Seats = int(seats.(float64))
	}
	if transmission, ok := updates["transmission"]; ok {
		car.Transmission = transmission.(string)
	}
	if fuelType, ok := updates["fuel_type"]; ok {
		car.FuelType = fuelType.(string)
	}
	if dailyRent, ok := updates["daily_rent"]; ok {
		car.DailyRent = dailyRent.(float64)
	}
	if deposit, ok := updates["deposit"]; ok {
		car.Deposit = deposit.(float64)
	}
	if licensePlate, ok := updates["license_plate"]; ok {
		if licensePlate.(string) != car.LicensePlate && s.carRepo.ExistsByLicensePlate(licensePlate.(string)) {
			return errors.New("车牌号已存在")
		}
		car.LicensePlate = licensePlate.(string)
	}
	if color, ok := updates["color"]; ok {
		car.Color = color.(string)
	}
	if mileage, ok := updates["mileage"]; ok {
		car.Mileage = int(mileage.(float64))
	}
	if storeID, ok := updates["store_id"]; ok {
		car.StoreID = uint(storeID.(float64))
	}
	if features, ok := updates["features"]; ok {
		car.Features = features.(string)
	}
	if description, ok := updates["description"]; ok {
		car.Description = description.(string)
	}

	return s.carRepo.Update(car)
}

func (s *CarService) UpdateStatus(id uint, status model.CarStatus) error {
	return s.carRepo.UpdateStatus(id, status)
}

func (s *CarService) DeleteCar(id uint) error {
	return s.carRepo.Delete(id)
}

func (s *CarService) UploadCarImage(carID uint, file *multipart.FileHeader) (*model.CarImage, error) {
	car, err := s.carRepo.FindByID(carID)
	if err != nil {
		return nil, errors.New("车辆不存在")
	}

	uploadPath := filepath.Join(s.cfg.Upload.Path, "cars")
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return nil, err
	}

	fileName := utils.GenerateFileName(file.Filename)
	savePath := filepath.Join(uploadPath, fileName)

	if err := saveMultipartFile(file, savePath); err != nil {
		return nil, err
	}

	image := &model.CarImage{
		CarID: car.ID,
		URL:   "/uploads/cars/" + fileName,
	}

	err = s.carRepo.AddImage(image)
	if err != nil {
		os.Remove(savePath)
		return nil, err
	}

	return image, nil
}

func (s *CarService) BatchUploadCarImages(carID uint, files []*multipart.FileHeader) ([]model.CarImage, error) {
	var images []model.CarImage

	for _, file := range files {
		image, err := s.UploadCarImage(carID, file)
		if err != nil {
			continue
		}
		images = append(images, *image)
	}

	return images, nil
}

func (s *CarService) DeleteCarImage(id uint) error {
	images, err := s.carRepo.GetImages(0)
	if err != nil {
		return err
	}

	for _, img := range images {
		if img.ID == id {
			filePath := filepath.Join(s.cfg.Upload.Path, strings.TrimPrefix(img.URL, "/uploads/"))
			os.Remove(filePath)
			break
		}
	}

	return s.carRepo.DeleteImage(id)
}

func (s *CarService) GetCarImages(carID uint) ([]model.CarImage, error) {
	return s.carRepo.GetImages(carID)
}

func (s *CarService) GetAvailableCars(storeID uint, page, pageSize int) ([]model.Car, int64, error) {
	return s.carRepo.GetAvailableCars(storeID, page, pageSize)
}

func saveMultipartFile(fh *multipart.FileHeader, dst string) error {
	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(src)
	return err
}