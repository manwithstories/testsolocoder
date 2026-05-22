package handler

import (
	"property-management/internal/database"
	"property-management/internal/model"
	redisclient "property-management/internal/redis"
	"property-management/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TenantHandler struct{}

func NewTenantHandler() *TenantHandler {
	return &TenantHandler{}
}

type TenantRequest struct {
	Name   string `json:"name" binding:"required"`
	Phone  string `json:"phone" binding:"required"`
	IDCard string `json:"idCard"`
	Email  string `json:"email"`
}

func (h *TenantHandler) Create(c *gin.Context) {
	var req TenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	tenant := model.Tenant{
		Name:   req.Name,
		Phone:  req.Phone,
		IDCard: req.IDCard,
		Email:  req.Email,
		Status: 1,
	}

	if err := database.DB.Create(&tenant).Error; err != nil {
		utils.Error(c, 500, "Failed to create tenant")
		return
	}

	utils.Success(c, tenant)
}

func (h *TenantHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var tenant model.Tenant
	if err := database.DB.First(&tenant, id).Error; err != nil {
		utils.Error(c, 404, "Tenant not found")
		return
	}

	var req TenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	database.DB.Model(&tenant).Updates(map[string]interface{}{
		"name":    req.Name,
		"phone":   req.Phone,
		"id_card": req.IDCard,
		"email":   req.Email,
	})

	utils.Success(c, nil)
}

func (h *TenantHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	database.DB.Delete(&model.Tenant{}, id)
	utils.Success(c, nil)
}

func (h *TenantHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query := database.DB.Model(&model.Tenant{})

	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var tenants []model.Tenant
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&tenants)

	utils.Success(c, gin.H{
		"list":     tenants,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *TenantHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var tenant model.Tenant
	if err := database.DB.First(&tenant, id).Error; err != nil {
		utils.Error(c, 404, "Tenant not found")
		return
	}
	utils.Success(c, tenant)
}

type AppointmentRequest struct {
	TenantID   uint   `json:"tenantId" binding:"required"`
	PropertyID uint   `json:"propertyId" binding:"required"`
	VisitTime  string `json:"visitTime" binding:"required"`
	Remark     string `json:"remark"`
}

func (h *TenantHandler) CreateAppointment(c *gin.Context) {
	var req AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	visitTime, err := time.Parse("2006-01-02 15:04:05", req.VisitTime)
	if err != nil {
		utils.Error(c, 400, "Invalid visit time format")
		return
	}

	appointment := model.Appointment{
		TenantID:   req.TenantID,
		PropertyID: req.PropertyID,
		VisitTime:  visitTime,
		Status:     1,
		Remark:     req.Remark,
	}

	if err := database.DB.Create(&appointment).Error; err != nil {
		utils.Error(c, 500, "Failed to create appointment")
		return
	}

	utils.Success(c, appointment)
}

func (h *TenantHandler) UpdateAppointmentStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	database.DB.Model(&model.Appointment{}).Where("id = ?", id).Update("status", req.Status)
	utils.Success(c, nil)
}

func (h *TenantHandler) ListAppointments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query := database.DB.Model(&model.Appointment{})

	if status := c.Query("status"); status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}

	var total int64
	query.Count(&total)

	var appointments []model.Appointment
	offset := (page - 1) * pageSize
	query.Preload("Tenant").Preload("Property").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&appointments)

	utils.Success(c, gin.H{
		"list":     appointments,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

type ContractRequest struct {
	TenantID    uint   `json:"tenantId" binding:"required"`
	PropertyID  uint   `json:"propertyId" binding:"required"`
	StartDate   string `json:"startDate" binding:"required"`
	EndDate     string `json:"endDate" binding:"required"`
	Rent        float64 `json:"rent" binding:"required"`
	Deposit     float64 `json:"deposit"`
	PaymentType string `json:"paymentType"`
	FileURL     string `json:"fileUrl"`
}

func (h *TenantHandler) CreateContract(c *gin.Context) {
	var req ContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	contract := model.Contract{
		TenantID:    req.TenantID,
		PropertyID:  req.PropertyID,
		StartDate:   startDate,
		EndDate:     endDate,
		Rent:        req.Rent,
		Deposit:     req.Deposit,
		PaymentType: req.PaymentType,
		Status:      1,
		FileURL:     req.FileURL,
	}

	if err := database.DB.Create(&contract).Error; err != nil {
		utils.Error(c, 500, "Failed to create contract")
		return
	}

	redisclient.AddContractReminder(contract.ID, contract.EndDate)

	var property model.Property
	database.DB.First(&property, req.PropertyID)
	property.Status = 2
	database.DB.Save(&property)
	redisclient.CachePropertyStatus(property.ID, 2)

	utils.Success(c, contract)
}

func (h *TenantHandler) UpdateContract(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var contract model.Contract
	if err := database.DB.First(&contract, id).Error; err != nil {
		utils.Error(c, 404, "Contract not found")
		return
	}

	var req ContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	redisclient.RemoveContractReminder(contract.ID)

	database.DB.Model(&contract).Updates(map[string]interface{}{
		"tenant_id":    req.TenantID,
		"property_id":  req.PropertyID,
		"start_date":   startDate,
		"end_date":     endDate,
		"rent":         req.Rent,
		"deposit":      req.Deposit,
		"payment_type": req.PaymentType,
		"file_url":     req.FileURL,
	})

	redisclient.AddContractReminder(contract.ID, endDate)

	utils.Success(c, nil)
}

func (h *TenantHandler) UpdateContractStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	database.DB.Model(&model.Contract{}).Where("id = ?", id).Update("status", req.Status)

	if req.Status == 2 {
		redisclient.RemoveContractReminder(uint(id))
	}

	utils.Success(c, nil)
}

func (h *TenantHandler) ListContracts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query := database.DB.Model(&model.Contract{})

	if status := c.Query("status"); status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}

	var total int64
	query.Count(&total)

	var contracts []model.Contract
	offset := (page - 1) * pageSize
	query.Preload("Tenant").Preload("Property").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&contracts)

	utils.Success(c, gin.H{
		"list":     contracts,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *TenantHandler) DetailContract(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var contract model.Contract
	if err := database.DB.Preload("Tenant").Preload("Property").First(&contract, id).Error; err != nil {
		utils.Error(c, 404, "Contract not found")
		return
	}
	utils.Success(c, contract)
}

func (h *TenantHandler) GetExpiringContracts(c *gin.Context) {
	now := time.Now()
	remindDate := now.AddDate(0, 1, 0)

	members := redisclient.GetExpiringContracts(now, remindDate)

	utils.Success(c, gin.H{
		"expiring": members,
		"count":   len(members),
	})
}
