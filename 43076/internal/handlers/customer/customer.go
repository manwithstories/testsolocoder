package customer

import (
	"strconv"

	"ticket-system/internal/database"
	"ticket-system/internal/models"
	"ticket-system/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateCustomerRequest struct {
	Name    string `json:"name" binding:"required,max=100"`
	Email   string `json:"email" binding:"required,email,max=100"`
	Phone   string `json:"phone" binding:"max=20"`
	Company string `json:"company" binding:"max=100"`
	Address string `json:"address" binding:"max=500"`
}

type UpdateCustomerRequest struct {
	Name    string `json:"name" binding:"max=100"`
	Email   string `json:"email" binding:"omitempty,email,max=100"`
	Phone   string `json:"phone" binding:"max=20"`
	Company string `json:"company" binding:"max=100"`
	Address string `json:"address" binding:"max=500"`
}

func CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	var existingCount int64
	database.DB.Model(&models.Customer{}).Where("email = ?", req.Email).Count(&existingCount)
	if existingCount > 0 {
		utils.BadRequest(c, "Email already exists")
		return
	}

	customer := &models.Customer{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Company: req.Company,
		Address: req.Address,
	}

	if err := database.DB.Create(customer).Error; err != nil {
		utils.InternalServerError(c, "Failed to create customer")
		return
	}

	utils.Success(c, customer)
}

func GetCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid customer ID")
		return
	}

	var customer models.Customer
	if err := database.DB.First(&customer, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Customer not found")
			return
		}
		utils.InternalServerError(c, "Failed to get customer")
		return
	}

	utils.Success(c, customer)
}

func GetCustomerWithTickets(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid customer ID")
		return
	}

	var customer models.Customer
	if err := database.DB.Preload("Tickets", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(50)
	}).First(&customer, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Customer not found")
			return
		}
		utils.InternalServerError(c, "Failed to get customer")
		return
	}

	utils.Success(c, customer)
}

func ListCustomers(c *gin.Context) {
	var customers []models.Customer
	query := database.DB.Model(&models.Customer{})

	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if email := c.Query("email"); email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}
	if company := c.Query("company"); company != "" {
		query = query.Where("company LIKE ?", "%"+company+"%")
	}
	if phone := c.Query("phone"); phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}

	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if pn, err := strconv.Atoi(p); err == nil && pn > 0 {
			page = pn
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if psn, err := strconv.Atoi(ps); err == nil && psn > 0 {
			pageSize = psn
		}
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&customers).Error; err != nil {
		utils.InternalServerError(c, "Failed to list customers")
		return
	}

	utils.Success(c, gin.H{
		"items":     customers,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid customer ID")
		return
	}

	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Company != "" {
		updates["company"] = req.Company
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}

	if err := database.DB.Model(&models.Customer{}).Where("id = ?", uint(id)).Updates(updates).Error; err != nil {
		utils.InternalServerError(c, "Failed to update customer")
		return
	}

	var customer models.Customer
	database.DB.First(&customer, uint(id))
	utils.Success(c, customer)
}

func DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid customer ID")
		return
	}

	var count int64
	database.DB.Model(&models.Ticket{}).Where("customer_id = ?", uint(id)).Count(&count)
	if count > 0 {
		utils.BadRequest(c, "Cannot delete customer with existing tickets")
		return
	}

	if err := database.DB.Delete(&models.Customer{}, uint(id)).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete customer")
		return
	}

	utils.Success(c, nil)
}
