package controllers

import (
	"health-platform/services"
	"health-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CompanyController struct {
	companyService *services.CompanyService
}

func NewCompanyController() *CompanyController {
	return &CompanyController{
		companyService: services.NewCompanyService(),
	}
}

func (ctrl *CompanyController) RegisterCompany(c *gin.Context) {
	var req services.RegisterCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	company, hrUser, err := ctrl.companyService.RegisterCompany(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"company": company,
		"hr_user": hrUser,
	})
}

func (ctrl *CompanyController) GetCompany(c *gin.Context) {
	companyID := c.GetUint("company_id")
	if companyID == 0 {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		companyID = uint(id)
	}

	company, err := ctrl.companyService.GetCompany(companyID)
	if err != nil {
		utils.Error(c, 404, "企业不存在")
		return
	}

	utils.Success(c, company)
}

func (ctrl *CompanyController) UpdateCompany(c *gin.Context) {
	companyID := c.GetUint("company_id")
	
	var req services.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.companyService.UpdateCompany(companyID, &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *CompanyController) AddDepartment(c *gin.Context) {
	var req services.AddDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	department, err := ctrl.companyService.AddDepartment(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, department)
}

func (ctrl *CompanyController) UpdateDepartment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	var req services.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.companyService.UpdateDepartment(uint(id), &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *CompanyController) GetDepartments(c *gin.Context) {
	companyID := c.GetUint("company_id")
	
	departments, err := ctrl.companyService.GetDepartments(companyID)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, departments)
}

func (ctrl *CompanyController) AddEmployee(c *gin.Context) {
	var req services.AddEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	employee, err := ctrl.companyService.AddEmployee(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, employee)
}

func (ctrl *CompanyController) UpdateEmployee(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	var req services.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.companyService.UpdateEmployee(uint(id), &req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *CompanyController) GetEmployees(c *gin.Context) {
	companyID := c.GetUint("company_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	employees, total, err := ctrl.companyService.GetEmployees(companyID, page, pageSize)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.PaginatedResponse(c, total, page, pageSize, employees)
}

func (ctrl *CompanyController) GetEmployee(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	
	employee, err := ctrl.companyService.GetEmployee(uint(id))
	if err != nil {
		utils.Error(c, 404, "员工不存在")
		return
	}

	utils.Success(c, employee)
}

func (ctrl *CompanyController) SetBudget(c *gin.Context) {
	var req services.SetBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.companyService.SetBudget(&req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *CompanyController) GetBudget(c *gin.Context) {
	companyID := c.GetUint("company_id")
	year, _ := strconv.Atoi(c.DefaultQuery("year", "2024"))
	
	budget, err := ctrl.companyService.GetBudget(companyID, year)
	if err != nil {
		utils.Error(c, 404, "预算信息不存在")
		return
	}

	utils.Success(c, budget)
}

func (ctrl *CompanyController) SetDepartmentAppointment(c *gin.Context) {
	var req services.SetDepartmentAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.companyService.SetDepartmentAppointment(&req); err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *CompanyController) GetDepartmentAppointments(c *gin.Context) {
	departmentID, _ := strconv.ParseUint(c.Param("department_id"), 10, 64)
	year, _ := strconv.Atoi(c.DefaultQuery("year", "2024"))
	
	appointments, err := ctrl.companyService.GetDepartmentAppointments(uint(departmentID), year)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, appointments)
}
