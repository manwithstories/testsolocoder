package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"printshop/internal/models"

	"gorm.io/gorm"
)

type PricingService struct {
	db *gorm.DB
}

func NewPricingService(db *gorm.DB) *PricingService {
	return &PricingService{db: db}
}

type PricingInput struct {
	TemplateID    uint
	MaterialID    uint
	ProcessIDs    []uint
	Quantity      int
	CustomerLevel string
	Urgent        bool
}

func (s *PricingService) Calculate(input PricingInput) (float64, []string, error) {
	var template models.Template
	if err := s.db.Preload("Materials").Preload("Processes").First(&template, input.TemplateID).Error; err != nil {
		return 0, nil, fmt.Errorf("template not found: %w", err)
	}
	if !template.Active {
		return 0, nil, errors.New("template is inactive")
	}

	var material models.TemplateMaterial
	if err := s.db.First(&material, input.MaterialID).Error; err != nil {
		return 0, nil, fmt.Errorf("material not found: %w", err)
	}
	if material.TemplateID != template.ID {
		return 0, nil, errors.New("material does not belong to template")
	}

	var rule models.PriceRule
	query := s.db.Where("category = ? AND min_qty <= ? AND max_qty >= ?", template.Category, input.Quantity, input.Quantity)
	query = query.Where("customer_level IN ?", []string{input.CustomerLevel, "all"})
	query = query.Order("discount DESC")
	if err := query.First(&rule).Error; err != nil {
		rule = models.PriceRule{UnitPrice: material.BasePrice, Discount: 0}
	}

	basePrice := rule.UnitPrice
	if basePrice == 0 {
		basePrice = material.BasePrice
	}
	unitPrice := basePrice

	var processExtra float64
	var processNames []string
	for _, pid := range input.ProcessIDs {
		var proc models.TemplateProcess
		if err := s.db.First(&proc, pid).Error; err != nil {
			continue
		}
		if proc.TemplateID != template.ID {
			continue
		}
		processExtra += proc.ExtraPrice
		processNames = append(processNames, proc.Name)
	}
	unitPrice += processExtra

	total := unitPrice * float64(input.Quantity)

	if rule.Discount > 0 {
		total = total * (1 - rule.Discount/100)
	}

	if input.Urgent {
		total = total * 1.3
	}

	return total, processNames, nil
}

func (s *PricingService) GetRules() ([]models.PriceRule, error) {
	var rules []models.PriceRule
	if err := s.db.Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

func (s *PricingService) CreateRule(rule *models.PriceRule) error {
	return s.db.Create(rule).Error
}

func (s *PricingService) UpdateRule(id uint, rule *models.PriceRule) error {
	return s.db.Model(&models.PriceRule{}).Where("id = ?", id).Updates(rule).Error
}

func (s *PricingService) DeleteRule(id uint) error {
	return s.db.Delete(&models.PriceRule{}, id).Error
}

type ProductionService struct {
	db *gorm.DB
}

func NewProductionService(db *gorm.DB) *ProductionService {
	return &ProductionService{db: db}
}

func (s *ProductionService) GetLines() ([]models.ProductionLine, error) {
	var lines []models.ProductionLine
	if err := s.db.Find(&lines).Error; err != nil {
		return nil, err
	}
	return lines, nil
}

func (s *ProductionService) CreateLine(line *models.ProductionLine) error {
	return s.db.Create(line).Error
}

func (s *ProductionService) ScheduleOrder(orderID uint, lineID uint, qty int, startDate time.Time, endDate time.Time) (*models.ProductionSchedule, error) {
	var line models.ProductionLine
	if err := s.db.First(&line, lineID).Error; err != nil {
		return nil, fmt.Errorf("production line not found: %w", err)
	}
	if !line.Active {
		return nil, errors.New("production line is inactive")
	}
	if line.Workload+qty > line.Capacity {
		return nil, errors.New("production line capacity exceeded")
	}
	sched := models.ProductionSchedule{
		OrderID:     orderID,
		LineID:      lineID,
		PlannedQty:  qty,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      "scheduled",
	}
	if err := s.db.Create(&sched).Error; err != nil {
		return nil, err
	}
	s.db.Model(&line).UpdateColumn("workload", gorm.Expr("workload + ?", qty))
	return &sched, nil
}

func (s *ProductionService) UpdateProgress(scheduleID uint, producedQty int) error {
	var sched models.ProductionSchedule
	if err := s.db.First(&sched, scheduleID).Error; err != nil {
		return err
	}
	s.db.Model(&sched).Update("produced_qty", producedQty)
	if producedQty >= sched.PlannedQty {
		s.db.Model(&sched).Update("status", "completed")
	}
	return nil
}

func (s *ProductionService) GetSchedules() ([]models.ProductionSchedule, error) {
	var schedules []models.ProductionSchedule
	if err := s.db.Preload("Line").Preload("Order").Order("start_date ASC").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

type CustomerService struct {
	db *gorm.DB
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{db: db}
}

func (s *CustomerService) GetAll() ([]models.Customer, error) {
	var customers []models.Customer
	if err := s.db.Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func (s *CustomerService) GetByID(id uint) (*models.Customer, error) {
	var c models.Customer
	if err := s.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *CustomerService) Create(c *models.Customer) error {
	return s.db.Create(c).Error
}

func (s *CustomerService) Update(id uint, c *models.Customer) error {
	return s.db.Model(&models.Customer{}).Where("id = ?", id).Updates(c).Error
}

func (s *CustomerService) Delete(id uint) error {
	return s.db.Delete(&models.Customer{}, id).Error
}

func (s *CustomerService) GenerateInvoice(customerID uint, periodStart, periodEnd time.Time) (*models.Invoice, error) {
	var orders []models.Order
	if err := s.db.Where("customer_id = ? AND created_at BETWEEN ? AND ? AND status != ?", customerID, periodStart, periodEnd, "cancelled").Find(&orders).Error; err != nil {
		return nil, err
	}

	invNo := generateInvoiceNo()
	inv := models.Invoice{
		InvoiceNo:   invNo,
		CustomerID:  customerID,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		Status:      "unpaid",
	}
	var total float64
	for _, o := range orders {
		inv.Items = append(inv.Items, models.InvoiceItem{
			OrderID: o.ID,
			Amount:  o.FinalPrice,
		})
		total += o.FinalPrice
	}
	inv.TotalAmount = total
	if err := s.db.Create(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (s *CustomerService) GetInvoices(customerID uint) ([]models.Invoice, error) {
	var invoices []models.Invoice
	query := s.db.Preload("Items").Preload("Customer")
	if customerID > 0 {
		query = query.Where("customer_id = ?", customerID)
	}
	if err := query.Order("created_at DESC").Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}

func generateInvoiceNo() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("INV%s%s", time.Now().Format("20060102"), hex.EncodeToString(b))
}

type DashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

type DashboardStats struct {
	TotalOrders    int64   `json:"total_orders"`
	TotalRevenue   float64 `json:"total_revenue"`
	PendingOrders  int64   `json:"pending_orders"`
	ProducingOrders int64 `json:"producing_orders"`
	CompletedOrders int64 `json:"completed_orders"`
	UrgentOrders   int64   `json:"urgent_orders"`
	TotalCustomers int64   `json:"total_customers"`
	CapacityUsage  float64 `json:"capacity_usage"`
	RecentOrders   []models.Order `json:"recent_orders"`
}

func (s *DashboardService) GetStats() (*DashboardStats, error) {
	var stats DashboardStats
	s.db.Model(&models.Order{}).Count(&stats.TotalOrders)
	s.db.Model(&models.Order{}).Where("status = ?", "created").Count(&stats.PendingOrders)
	s.db.Model(&models.Order{}).Where("status = ?", "producing").Count(&stats.ProducingOrders)
	s.db.Model(&models.Order{}).Where("status = ?", "shipped").Count(&stats.CompletedOrders)
	s.db.Model(&models.Order{}).Where("urgent = ? AND status NOT IN ?", true, []string{"shipped", "cancelled"}).Count(&stats.UrgentOrders)
	s.db.Model(&models.Customer{}).Count(&stats.TotalCustomers)

	s.db.Model(&models.Order{}).
		Where("status NOT IN ?", []string{"cancelled"}).
		Select("COALESCE(SUM(final_price), 0)").Scan(&stats.TotalRevenue)

	var lines []models.ProductionLine
	s.db.Find(&lines)
	var totalCap, totalLoad int
	for _, l := range lines {
		if l.Active {
			totalCap += l.Capacity
			totalLoad += l.Workload
		}
	}
	if totalCap > 0 {
		stats.CapacityUsage = float64(totalLoad) / float64(totalCap) * 100
	}

	s.db.Preload("Customer").Order("created_at DESC").Limit(10).Find(&stats.RecentOrders)

	return &stats, nil
}

type FileService struct {
	db *gorm.DB
}

func NewFileService(db *gorm.DB) *FileService {
	return &FileService{db: db}
}

var AllowedExtensions = map[string]bool{
	".pdf":  true,
	".ai":   true,
	".eps":  true,
	".psd":  true,
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".cdr":  true,
}

func (s *FileService) ValidateFile(filename string) error {
	ext := getExtension(filename)
	if !AllowedExtensions[ext] {
		return fmt.Errorf("unsupported file type: %s", ext)
	}
	return nil
}

func (s *FileService) SaveFile(file *models.FileAsset) error {
	return s.db.Create(file).Error
}

func (s *FileService) GetFile(id uint) (*models.FileAsset, error) {
	var f models.FileAsset
	if err := s.db.First(&f, id).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func getExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}
