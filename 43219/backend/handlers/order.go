package handlers

import (
	"time"

	"housekeeping/config"
	"housekeeping/database"
	"housekeeping/models"
	"housekeeping/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderReportInput struct {
	ReportText   string `json:"report_text"`
	ReportImages string `json:"report_images"`
}

func SubmitReport(c *gin.Context) {
	uid, _ := c.Get("uid")
	id := c.Param("id")
	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.NotFound(c, "order not found")
		return
	}
	if order.StaffID != uid.(uint) {
		utils.Forbidden(c, "not your order")
		return
	}
	if order.Status != models.OrderCreated && order.Status != models.OrderReported {
		utils.BadRequest(c, "cannot report")
		return
	}
	var in OrderReportInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	now := timeNow()
	order.ReportText = in.ReportText
	order.ReportImages = in.ReportImages
	order.ReportedAt = &now
	order.Status = models.OrderReported
	database.DB.Save(&order)
	utils.OK(c, order)
}

func ConfirmOrder(c *gin.Context) {
	uid, _ := c.Get("uid")
	id := c.Param("id")
	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.NotFound(c, "order not found")
		return
	}
	if order.CustomerID != uid.(uint) {
		utils.Forbidden(c, "not your order")
		return
	}
	if order.Status != models.OrderReported && order.Status != models.OrderCreated {
		utils.BadRequest(c, "cannot confirm")
		return
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		now := timeNow()
		order.ConfirmedAt = &now
		order.Status = models.OrderPaid
		order.PaidAt = &now
		if err := tx.Save(&order).Error; err != nil {
			return err
		}
		ratio := config.C.Payment.StaffRatio
		if ratio <= 0 {
			ratio = 0.7
		}
		staffShare := order.TotalAmount * ratio
		companyShare := order.TotalAmount - staffShare
		settle := models.Settlement{
			OrderID:      order.ID,
			CompanyID:    order.CompanyID,
			StaffID:      order.StaffID,
			TotalAmount:  order.TotalAmount,
			StaffShare:   staffShare,
			CompanyShare: companyShare,
			Status:       "paid",
		}
		if err := tx.Create(&settle).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Wallet{}).Where("user_id = ?", order.StaffID).
			UpdateColumn("balance", gorm.Expr("balance + ?", staffShare)).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Wallet{}).Where("user_id = ?", order.CompanyID).
			UpdateColumn("balance", gorm.Expr("balance + ?", companyShare)).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		utils.ServerError(c, "confirm failed")
		return
	}
	utils.OK(c, order)
}

func AutoConfirmOrders() {
	limit := config.C.Order.AutoConfirmHours
	if limit <= 0 {
		return
	}
	cutoff := time.Now().Add(-time.Duration(limit) * time.Hour)
	var orders []models.Order
	database.DB.Where("status = ? AND reported_at IS NOT NULL AND reported_at < ?", models.OrderReported, cutoff).Find(&orders)
	for _, o := range orders {
		confirmOrderInternal(&o)
	}
}

func confirmOrderInternal(order *models.Order) {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		now := timeNow()
		order.ConfirmedAt = &now
		order.Status = models.OrderPaid
		order.PaidAt = &now
		if err := tx.Save(order).Error; err != nil {
			return err
		}
		ratio := config.C.Payment.StaffRatio
		if ratio <= 0 {
			ratio = 0.7
		}
		staffShare := order.TotalAmount * ratio
		companyShare := order.TotalAmount - staffShare
		settle := models.Settlement{
			OrderID:      order.ID,
			CompanyID:    order.CompanyID,
			StaffID:      order.StaffID,
			TotalAmount:  order.TotalAmount,
			StaffShare:   staffShare,
			CompanyShare: companyShare,
			Status:       "paid",
		}
		if err := tx.Create(&settle).Error; err != nil {
			return err
		}
		tx.Model(&models.Wallet{}).Where("user_id = ?", order.StaffID).UpdateColumn("balance", gorm.Expr("balance + ?", staffShare))
		tx.Model(&models.Wallet{}).Where("user_id = ?", order.CompanyID).UpdateColumn("balance", gorm.Expr("balance + ?", companyShare))
		return nil
	})
	if err != nil {
		utils.Logger.Errorw("auto confirm failed", "order", order.ID, "err", err)
	}
}

func RequestRefund(c *gin.Context) {
	uid, _ := c.Get("uid")
	id := c.Param("id")
	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.NotFound(c, "order not found")
		return
	}
	if order.CustomerID != uid.(uint) {
		utils.Forbidden(c, "not your order")
		return
	}
	if order.Status == models.OrderRefunding || order.Status == models.OrderRefunded {
		utils.BadRequest(c, "already requesting refund")
		return
	}
	order.Status = models.OrderRefunding
	database.DB.Save(&order)
	var body struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&body)
	ticket := models.Ticket{
		OrderID:      &order.ID,
		CustomerID:   order.CustomerID,
		StaffID:      &order.StaffID,
		Type:         models.TicketRefund,
		Title:        "退款申请 - 订单 #" + uintStr(order.ID),
		Content:      body.Reason,
		Status:       models.TicketOpen,
		LastActionAt: timeNow(),
	}
	database.DB.Create(&ticket)
	utils.OK(c, gin.H{"order": order, "ticket": ticket})
}

func ListOrders(c *gin.Context) {
	uid, _ := c.Get("uid")
	role, _ := c.Get("role")
	var list []models.Order
	q := database.DB.Preload("Booking").Preload("Booking.Service")
	switch role.(string) {
	case string(models.RoleCustomer):
		q = q.Where("customer_id = ?", uid)
	case string(models.RoleStaff):
		q = q.Where("staff_id = ?", uid)
	case string(models.RoleCompany):
		q = q.Where("company_id = ?", uid)
	}
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("id desc").Find(&list).Error; err != nil {
		utils.ServerError(c, "query failed")
		return
	}
	utils.OK(c, list)
}

func GetOrder(c *gin.Context) {
	id := c.Param("id")
	var o models.Order
	if err := database.DB.Preload("Booking").Preload("Booking.Service").First(&o, id).Error; err != nil {
		utils.NotFound(c, "order not found")
		return
	}
	utils.OK(c, o)
}
