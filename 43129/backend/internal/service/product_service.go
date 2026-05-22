package service

import (
	"fmt"
	"time"

	"beauty-salon-system/internal/model"
	"beauty-salon-system/internal/repository"
	"beauty-salon-system/internal/repository/mysql"
)

type ProductService struct {
	productRepo       *repository.ProductRepository
	productRecordRepo *repository.ProductRecordRepository
	productSaleRepo   *repository.ProductSaleRepository
}

func NewProductService(
	productRepo *repository.ProductRepository,
	productRecordRepo *repository.ProductRecordRepository,
	productSaleRepo *repository.ProductSaleRepository,
) *ProductService {
	return &ProductService{
		productRepo:       productRepo,
		productRecordRepo: productRecordRepo,
		productSaleRepo:   productSaleRepo,
	}
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	Stock       int     `json:"stock"`
	Threshold   int     `json:"threshold"`
	Price       float64 `json:"price"`
	RetailPrice float64 `json:"retail_price"`
	Supplier    string  `json:"supplier"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	Threshold   int     `json:"threshold"`
	Price       float64 `json:"price"`
	RetailPrice float64 `json:"retail_price"`
	Supplier    string  `json:"supplier"`
	Status      int     `json:"status"`
}

type AddStockRequest struct {
	ProductID  uint   `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required"`
	OperatorID uint   `json:"operator_id"`
	Remark     string `json:"remark"`
}

type DeductStockRequest struct {
	ProductID     uint   `json:"product_id" binding:"required"`
	Quantity      int    `json:"quantity" binding:"required"`
	AppointmentID *uint  `json:"appointment_id"`
	OperatorID    uint   `json:"operator_id"`
	Remark        string `json:"remark"`
}

type SaleProductRequest struct {
	CustomerID uint   `json:"customer_id" binding:"required"`
	ProductID  uint   `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required"`
	PayMethod  string `json:"pay_method" binding:"required"`
	OperatorID uint   `json:"operator_id"`
}

func (s *ProductService) Create(req *CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		Name:        req.Name,
		Category:    req.Category,
		Description: req.Description,
		Unit:        req.Unit,
		Stock:       req.Stock,
		Threshold:   req.Threshold,
		Price:       req.Price,
		RetailPrice: req.RetailPrice,
		Supplier:    req.Supplier,
		Status:      1,
	}

	if product.Unit == "" {
		product.Unit = "件"
	}
	if product.Threshold <= 0 {
		product.Threshold = 10
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	return product, nil
}

func (s *ProductService) GetByID(id uint) (*model.Product, error) {
	return s.productRepo.GetByID(id)
}

func (s *ProductService) Update(id uint, req *UpdateProductRequest) (*model.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Unit != "" {
		product.Unit = req.Unit
	}
	if req.Threshold > 0 {
		product.Threshold = req.Threshold
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.RetailPrice > 0 {
		product.RetailPrice = req.RetailPrice
	}
	if req.Supplier != "" {
		product.Supplier = req.Supplier
	}
	if req.Status > 0 {
		product.Status = req.Status
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	return product, nil
}

func (s *ProductService) Delete(id uint) error {
	return s.productRepo.Delete(id)
}

func (s *ProductService) List(page, pageSize int, category string, lowStock bool) ([]model.Product, int64, error) {
	return s.productRepo.List(page, pageSize, category, lowStock)
}

func (s *ProductService) ListAll() ([]model.Product, error) {
	return s.productRepo.ListAll()
}

func (s *ProductService) AddStock(req *AddStockRequest) error {
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	record := &model.ProductRecord{
		ProductID:     req.ProductID,
		ChangeType:    "in",
		Quantity:      req.Quantity,
		BeforeStock:   product.Stock,
		AfterStock:    product.Stock + req.Quantity,
		OperatorID:    req.OperatorID,
		Remark:        req.Remark,
	}

	tx := mysql.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := s.productRepo.AddStock(req.ProductID, req.Quantity); err != nil {
		tx.Rollback()
		return fmt.Errorf("add stock: %w", err)
	}

	if err := s.productRecordRepo.Create(record, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("create record: %w", err)
	}

	return tx.Commit().Error
}

func (s *ProductService) DeductStock(req *DeductStockRequest) error {
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if product.Stock < req.Quantity {
		return fmt.Errorf("insufficient stock, current stock: %d", product.Stock)
	}

	record := &model.ProductRecord{
		ProductID:     req.ProductID,
		ChangeType:    "out",
		Quantity:      req.Quantity,
		BeforeStock:   product.Stock,
		AfterStock:    product.Stock - req.Quantity,
		AppointmentID: req.AppointmentID,
		OperatorID:    req.OperatorID,
		Remark:        req.Remark,
	}

	tx := mysql.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := s.productRepo.DeductStock(req.ProductID, req.Quantity, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("deduct stock: %w", err)
	}

	if err := s.productRecordRepo.Create(record, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("create record: %w", err)
	}

	return tx.Commit().Error
}

func (s *ProductService) GetProductRecords(page, pageSize int, productID uint, changeType string) ([]model.ProductRecord, int64, error) {
	return s.productRecordRepo.List(page, pageSize, productID, changeType)
}

func (s *ProductService) GetLowStockProducts() ([]model.Product, error) {
	return s.productRepo.GetLowStockProducts()
}

func (s *ProductService) SaleProduct(req *SaleProductRequest) (*model.ProductSale, error) {
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if product.Stock < req.Quantity {
		return nil, fmt.Errorf("insufficient stock")
	}

	totalPrice := product.RetailPrice * float64(req.Quantity)

	sale := &model.ProductSale{
		CustomerID: req.CustomerID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		UnitPrice:  product.RetailPrice,
		TotalPrice: totalPrice,
		PayMethod:  req.PayMethod,
		OperatorID: req.OperatorID,
	}

	record := &model.ProductRecord{
		ProductID:   req.ProductID,
		ChangeType:  "sale",
		Quantity:    req.Quantity,
		BeforeStock: product.Stock,
		AfterStock:  product.Stock - req.Quantity,
		OperatorID:  req.OperatorID,
		Remark:      fmt.Sprintf("Product sale, customer: %d", req.CustomerID),
	}

	tx := mysql.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := s.productRepo.DeductStock(req.ProductID, req.Quantity, tx); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("deduct stock: %w", err)
	}

	if err := s.productSaleRepo.Create(sale, tx); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("create sale: %w", err)
	}

	if err := s.productRecordRepo.Create(record, tx); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("create record: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return sale, nil
}

func (s *ProductService) GetProductSales(page, pageSize int, customerID uint) ([]model.ProductSale, int64, error) {
	return s.productSaleRepo.List(page, pageSize, customerID)
}

func (s *ProductService) StockCheck(products map[uint]int) error {
	for productID, quantity := range products {
		product, err := s.productRepo.GetByID(productID)
		if err != nil {
			return fmt.Errorf("product %d not found: %w", productID, err)
		}
		if product.Stock < quantity {
			return fmt.Errorf("product '%s' insufficient stock: need %d, have %d", product.Name, quantity, product.Stock)
		}
	}
	return nil
}

func (s *ProductService) BatchDeductStock(items map[uint]int, appointmentID *uint, operatorID uint) error {
	tx := mysql.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for productID, quantity := range items {
		product, err := s.productRepo.GetByID(productID)
		if err != nil {
			tx.Rollback()
			return err
		}

		if product.Stock < quantity {
			tx.Rollback()
			return fmt.Errorf("product '%s' insufficient stock", product.Name)
		}

		record := &model.ProductRecord{
			ProductID:     productID,
			ChangeType:    "out",
			Quantity:      quantity,
			BeforeStock:   product.Stock,
			AfterStock:    product.Stock - quantity,
			AppointmentID: appointmentID,
			OperatorID:    operatorID,
			Remark:        "Batch deduct for service",
		}

		if err := s.productRepo.DeductStock(productID, quantity, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := s.productRecordRepo.Create(record, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *ProductService) StockTake() ([]model.Product, error) {
	return s.productRepo.ListAll()
}

func (s *ProductService) RecordStockTake(results map[uint]int, operatorID uint) error {
	tx := mysql.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for productID, actualStock := range results {
		product, err := s.productRepo.GetByID(productID)
		if err != nil {
			tx.Rollback()
			return err
		}

		diff := actualStock - product.Stock
		if diff == 0 {
			continue
		}

		changeType := "adjust"
		if diff > 0 {
			changeType = "in"
		}

		record := &model.ProductRecord{
			ProductID:   productID,
			ChangeType:  changeType,
			Quantity:    diff,
			BeforeStock: product.Stock,
			AfterStock:  actualStock,
			OperatorID:  operatorID,
			Remark:      "Stock take adjustment",
		}

		if err := tx.Model(&model.Product{}).Where("id = ?", productID).
			Update("stock", actualStock).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := s.productRecordRepo.Create(record, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

type GetProductRecordsRequest struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	ProductID  uint   `json:"product_id"`
	ChangeType string `json:"change_type"`
}

func (s *ProductService) GetProductRecordsByRequest(req *GetProductRecordsRequest) ([]model.ProductRecord, int64, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.GetProductRecords(page, pageSize, req.ProductID, req.ChangeType)
}

func init() {
	_ = time.Now()
}
