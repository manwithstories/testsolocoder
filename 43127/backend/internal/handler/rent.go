package handler

import (
	"property-management/internal/database"
	"property-management/internal/model"
	"property-management/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RentHandler struct{}

func NewRentHandler() *RentHandler {
	return &RentHandler{}
}

type GenerateBillRequest struct {
	Month string `json:"month"`
}

func (h *RentHandler) GenerateMonthlyBills(c *gin.Context) {
	var req GenerateBillRequest
	c.ShouldBindJSON(&req)
	
	month := req.Month
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	var contracts []model.Contract
	database.DB.Where("status = ?", 1).Find(&contracts)

	generatedCount := 0
	for _, contract := range contracts {
		var count int64
		database.DB.Model(&model.RentRecord{}).
			Where("contract_id = ? AND month = ?", contract.ID, month).
			Count(&count)

		if count == 0 {
			dueDate, _ := time.Parse("2006-01-02", month+"-05")
			record := model.RentRecord{
				ContractID: contract.ID,
				TenantID:   contract.TenantID,
				PropertyID: contract.PropertyID,
				Month:      month,
				Amount:     contract.Rent,
				Status:     0,
				DueDate:    dueDate,
			}
			database.DB.Create(&record)
			generatedCount++
		}
	}

	utils.Success(c, gin.H{
		"generated": generatedCount,
		"month":     month,
	})
}

func (h *RentHandler) ListBills(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query := database.DB.Model(&model.RentRecord{})

	if month := c.Query("month"); month != "" {
		query = query.Where("month = ?", month)
	}
	if status := c.Query("status"); status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}
	if contractId := c.Query("contractId"); contractId != "" {
		id, _ := strconv.Atoi(contractId)
		query = query.Where("contract_id = ?", id)
	}

	var total int64
	query.Count(&total)

	var bills []model.RentRecord
	offset := (page - 1) * pageSize
	query.Preload("Contract").Preload("Contract.Tenant").Preload("Contract.Property").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&bills)

	var totalAmount float64
	var paidAmount float64
	for _, bill := range bills {
		totalAmount += bill.Amount
		if bill.Status == 1 {
			paidAmount += bill.Amount
		}
	}

	utils.Success(c, gin.H{
		"list":        bills,
		"total":       total,
		"page":        page,
		"pageSize":    pageSize,
		"totalAmount": totalAmount,
		"paidAmount":  paidAmount,
	})
}

type PayBillRequest struct {
	PaidAmount float64 `json:"paidAmount" binding:"required"`
	Remark     string  `json:"remark"`
}

func (h *RentHandler) PayBill(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var record model.RentRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		utils.Error(c, 404, "Bill not found")
		return
	}

	var req PayBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request parameters")
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":   1,
		"paid_at":  now,
		"remark":   req.Remark,
	}

	if now.After(record.DueDate) {
		daysLate := int(now.Sub(record.DueDate).Hours() / 24)
		lateFee := record.Amount * 0.005 * float64(daysLate)
		updates["late_fee"] = lateFee
	}

	database.DB.Model(&record).Updates(updates)
	utils.Success(c, nil)
}

func (h *RentHandler) CalculateLateFee(c *gin.Context) {
	var unpaidBills []model.RentRecord
	now := time.Now()
	database.DB.Where("status = ? AND due_date < ?", 0, now).Find(&unpaidBills)

	for _, bill := range unpaidBills {
		daysLate := int(now.Sub(bill.DueDate).Hours() / 24)
		if daysLate > 0 {
			lateFee := bill.Amount * 0.005 * float64(daysLate)
			database.DB.Model(&bill).Update("late_fee", lateFee)
		}
	}

	utils.Success(c, gin.H{
		"processed": len(unpaidBills),
	})
}

func (h *RentHandler) GetBillDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var record model.RentRecord
	if err := database.DB.Preload("Contract").Preload("Contract.Tenant").Preload("Contract.Property").
		First(&record, id).Error; err != nil {
		utils.Error(c, 404, "Bill not found")
		return
	}
	utils.Success(c, record)
}
