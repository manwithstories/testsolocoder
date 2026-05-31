package services

import (
	"errors"
	"fmt"
	"time"

	"health-platform/models"
	"health-platform/repository"

	"gorm.io/gorm"
)

type CompanyService struct {
	companyRepo    *repository.CompanyRepository
	departmentRepo *repository.BaseRepository
	employeeRepo   *repository.EmployeeRepository
	userRepo       *repository.UserRepository
	budgetRepo     *repository.CompanyBudgetRepository
	deptApptRepo   *repository.DepartmentAppointmentRepository
}

func NewCompanyService() *CompanyService {
	return &CompanyService{
		companyRepo:    repository.NewCompanyRepository(),
		departmentRepo: repository.NewBaseRepository(),
		employeeRepo:   repository.NewEmployeeRepository(),
		userRepo:       repository.NewUserRepository(),
		budgetRepo:     repository.NewCompanyBudgetRepository(),
		deptApptRepo:   repository.NewDepartmentAppointmentRepository(),
	}
}

type RegisterCompanyRequest struct {
	Name            string `json:"name" binding:"required,max=100"`
	UnifiedCode     string `json:"unified_code" binding:"required,max=50"`
	LegalPerson     string `json:"legal_person" binding:"required,max=50"`
	ContactPhone    string `json:"contact_phone" binding:"required,max=20"`
	ContactEmail    string `json:"contact_email" binding:"omitempty,email,max=100"`
	Address         string `json:"address" binding:"max=255"`
	HRUsername      string `json:"hr_username" binding:"required,min=3,max=50"`
	HRPassword      string `json:"hr_password" binding:"required,min=6,max=50"`
	HRRealName      string `json:"hr_real_name" binding:"required,max=50"`
	HRPhone         string `json:"hr_phone" binding:"required,max=20"`
	HREmail         string `json:"hr_email" binding:"omitempty,email,max=100"`
	AnnualBudget    float64 `json:"annual_budget"`
	PaymentType     int    `json:"payment_type"`
}

type UpdateCompanyRequest struct {
	ContactPhone    string  `json:"contact_phone" binding:"max=20"`
	ContactEmail    string  `json:"contact_email" binding:"omitempty,email,max=100"`
	Address         string  `json:"address" binding:"max=255"`
	AnnualBudget    float64 `json:"annual_budget"`
	PaymentType     int     `json:"payment_type"`
}

type AddDepartmentRequest struct {
	CompanyID    uint   `json:"company_id" binding:"required"`
	ParentID     *uint  `json:"parent_id"`
	Name         string `json:"name" binding:"required,max=100"`
	ManagerName  string `json:"manager_name" binding:"max=50"`
	ManagerPhone string `json:"manager_phone" binding:"max=20"`
}

type UpdateDepartmentRequest struct {
	Name         string `json:"name" binding:"max=100"`
	ManagerName  string `json:"manager_name" binding:"max=50"`
	ManagerPhone string `json:"manager_phone" binding:"max=20"`
	Status       int    `json:"status"`
}

type AddEmployeeRequest struct {
	CompanyID    uint   `json:"company_id" binding:"required"`
	DepartmentID uint   `json:"department_id" binding:"required"`
	EmployeeNo   string `json:"employee_no" binding:"max=50"`
	RealName     string `json:"real_name" binding:"required,max=50"`
	Gender       int    `json:"gender"`
	Birthday     *time.Time `json:"birthday"`
	IDCard       string `json:"id_card" binding:"max=20"`
	Phone        string `json:"phone" binding:"max=20"`
	Email        string `json:"email" binding:"omitempty,email,max=100"`
	Position     string `json:"position" binding:"max=50"`
	EntryDate    *time.Time `json:"entry_date"`
	Quota        int    `json:"quota"`
}

type UpdateEmployeeRequest struct {
	DepartmentID uint   `json:"department_id"`
	Position     string `json:"position" binding:"max=50"`
	Status       int    `json:"status"`
	Quota        int    `json:"quota"`
}

type SetBudgetRequest struct {
	CompanyID   uint    `json:"company_id" binding:"required"`
	Year        int     `json:"year" binding:"required"`
	TotalBudget float64 `json:"total_budget" binding:"required"`
	Frequency   int     `json:"frequency"`
}

type SetDepartmentAppointmentRequest struct {
	CompanyID    uint       `json:"company_id" binding:"required"`
	DepartmentID uint       `json:"department_id" binding:"required"`
	AgencyID     uint       `json:"agency_id" binding:"required"`
	PackageID    uint       `json:"package_id" binding:"required"`
	Year         int        `json:"year" binding:"required"`
	TotalQuota   int        `json:"total_quota" binding:"required"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
}

func (s *CompanyService) RegisterCompany(req *RegisterCompanyRequest) (*models.Company, *models.User, error) {
	existingCompany, _ := s.companyRepo.FindByUnifiedCode(req.UnifiedCode)
	if existingCompany != nil {
		return nil, nil, errors.New("该统一社会信用代码已注册")
	}

	existingName, _ := s.companyRepo.FindByName(req.Name)
	if existingName != nil {
		return nil, nil, errors.New("该企业名称已注册")
	}

	company := &models.Company{
		Name:          req.Name,
		UnifiedCode:   req.UnifiedCode,
		LegalPerson:   req.LegalPerson,
		ContactPhone:  req.ContactPhone,
		ContactEmail:  req.ContactEmail,
		Address:       req.Address,
		Status:        models.CompanyStatusActive,
		AnnualBudget:  req.AnnualBudget,
		PaymentType:   req.PaymentType,
	}

	if err := s.companyRepo.Create(company); err != nil {
		return nil, nil, fmt.Errorf("创建企业失败: %w", err)
	}

	authService := NewAuthService()
	hrUser, err := authService.Register(&RegisterRequest{
		Username:  req.HRUsername,
		Password:  req.HRPassword,
		RealName:  req.HRRealName,
		Phone:     req.HRPhone,
		Email:     req.HREmail,
		Role:      models.RoleHR,
		CompanyID: &company.ID,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("创建HR账号失败: %w", err)
	}

	company.HRUserID = &hrUser.ID
	if err := s.companyRepo.Save(company); err != nil {
		return nil, nil, fmt.Errorf("更新企业信息失败: %w", err)
	}

	return company, hrUser, nil
}

func (s *CompanyService) GetCompany(companyID uint) (*models.Company, error) {
	company, err := s.companyRepo.GetWithDepartments(companyID)
	if err != nil {
		return nil, err
	}
	return company, nil
}

func (s *CompanyService) UpdateCompany(companyID uint, req *UpdateCompanyRequest) error {
	var company models.Company
	if err := s.companyRepo.FindByID(&company, companyID); err != nil {
		return err
	}

	if req.ContactPhone != "" {
		company.ContactPhone = req.ContactPhone
	}
	if req.ContactEmail != "" {
		company.ContactEmail = req.ContactEmail
	}
	if req.Address != "" {
		company.Address = req.Address
	}
	company.AnnualBudget = req.AnnualBudget
	company.PaymentType = req.PaymentType

	return s.companyRepo.Save(&company)
}

func (s *CompanyService) AddDepartment(req *AddDepartmentRequest) (*models.Department, error) {
	department := &models.Department{
		CompanyID:    req.CompanyID,
		ParentID:     req.ParentID,
		Name:         req.Name,
		ManagerName:  req.ManagerName,
		ManagerPhone: req.ManagerPhone,
		Status:       1,
	}

	if err := s.departmentRepo.Create(department); err != nil {
		return nil, fmt.Errorf("添加部门失败: %w", err)
	}

	return department, nil
}

func (s *CompanyService) UpdateDepartment(departmentID uint, req *UpdateDepartmentRequest) error {
	var department models.Department
	if err := s.departmentRepo.FindByID(&department, departmentID); err != nil {
		return err
	}

	if req.Name != "" {
		department.Name = req.Name
	}
	department.ManagerName = req.ManagerName
	department.ManagerPhone = req.ManagerPhone
	department.Status = req.Status

	return s.departmentRepo.Save(&department)
}

func (s *CompanyService) GetDepartments(companyID uint) ([]models.Department, error) {
	var departments []models.Department
	err := s.departmentRepo.FindWithConditions(&departments, map[string]interface{}{
		"company_id": companyID,
	})
	return departments, err
}

func (s *CompanyService) AddEmployee(req *AddEmployeeRequest) (*models.Employee, error) {
	existing, _ := s.employeeRepo.FindByEmployeeNo(req.CompanyID, req.EmployeeNo)
	if existing != nil {
		return nil, errors.New("该工号已存在")
	}

	employee := &models.Employee{
		CompanyID:    req.CompanyID,
		DepartmentID: req.DepartmentID,
		EmployeeNo:   req.EmployeeNo,
		RealName:     req.RealName,
		Gender:       req.Gender,
		Birthday:     req.Birthday,
		IDCard:       req.IDCard,
		Phone:        req.Phone,
		Email:        req.Email,
		Position:     req.Position,
		EntryDate:    req.EntryDate,
		Status:       1,
		Quota:        req.Quota,
	}

	if err := s.employeeRepo.Create(employee); err != nil {
		return nil, fmt.Errorf("添加员工失败: %w", err)
	}

	return employee, nil
}

func (s *CompanyService) UpdateEmployee(employeeID uint, req *UpdateEmployeeRequest) error {
	var employee models.Employee
	if err := s.employeeRepo.FindByID(&employee, employeeID); err != nil {
		return err
	}

	if req.DepartmentID > 0 {
		employee.DepartmentID = req.DepartmentID
	}
	if req.Position != "" {
		employee.Position = req.Position
	}
	employee.Status = req.Status
	employee.Quota = req.Quota

	return s.employeeRepo.Save(&employee)
}

func (s *CompanyService) GetEmployees(companyID uint, page, pageSize int) ([]models.Employee, int64, error) {
	return s.employeeRepo.FindByCompanyID(companyID, page, pageSize)
}

func (s *CompanyService) GetEmployee(employeeID uint) (*models.Employee, error) {
	var employee models.Employee
	if err := s.employeeRepo.FindByID(&employee, employeeID); err != nil {
		return nil, err
	}
	return &employee, nil
}

func (s *CompanyService) BatchImportEmployees(employees []AddEmployeeRequest) (int, error) {
	var count int
	for _, emp := range employees {
		employee := &models.Employee{
			CompanyID:    emp.CompanyID,
			DepartmentID: emp.DepartmentID,
			EmployeeNo:   emp.EmployeeNo,
			RealName:     emp.RealName,
			Gender:       emp.Gender,
			Birthday:     emp.Birthday,
			IDCard:       emp.IDCard,
			Phone:        emp.Phone,
			Email:        emp.Email,
			Position:     emp.Position,
			EntryDate:    emp.EntryDate,
			Status:       1,
			Quota:        emp.Quota,
		}
		if err := s.employeeRepo.Create(employee); err == nil {
			count++
		}
	}
	return count, nil
}

func (s *CompanyService) SetBudget(req *SetBudgetRequest) error {
	existing, _ := s.budgetRepo.FindByCompanyIDAndYear(req.CompanyID, req.Year)
	if existing != nil {
		existing.TotalBudget = req.TotalBudget
		existing.Frequency = req.Frequency
		return s.budgetRepo.Save(existing)
	}

	budget := &models.CompanyBudget{
		CompanyID:   req.CompanyID,
		Year:        req.Year,
		TotalBudget: req.TotalBudget,
		Frequency:   req.Frequency,
	}
	return s.budgetRepo.Create(budget)
}

func (s *CompanyService) GetBudget(companyID uint, year int) (*models.CompanyBudget, error) {
	return s.budgetRepo.FindByCompanyIDAndYear(companyID, year)
}

func (s *CompanyService) SetDepartmentAppointment(req *SetDepartmentAppointmentRequest) error {
	existing := s.deptApptRepo.FindByDepartmentID(req.DepartmentID, req.Year)
	
	for _, deptAppt := range existing {
		if deptAppt.AgencyID == req.AgencyID && deptAppt.PackageID == req.PackageID {
			deptAppt.TotalQuota = req.TotalQuota
			deptAppt.StartDate = req.StartDate
			deptAppt.EndDate = req.EndDate
			return s.deptApptRepo.Save(&deptAppt)
		}
	}

	deptAppt := &models.DepartmentAppointment{
		CompanyID:    req.CompanyID,
		DepartmentID: req.DepartmentID,
		AgencyID:     req.AgencyID,
		PackageID:    req.PackageID,
		Year:         req.Year,
		TotalQuota:   req.TotalQuota,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Status:       1,
	}
	return s.deptApptRepo.Create(deptAppt)
}

func (s *CompanyService) GetDepartmentAppointments(departmentID uint, year int) ([]models.DepartmentAppointment, error) {
	return s.deptApptRepo.FindByDepartmentID(departmentID, year)
}
