package handlers

import (
	"housekeeping/database"
	"housekeeping/models"
	"housekeeping/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewInput struct {
	OrderID uint   `json:"order_id" binding:"required"`
	StaffID uint   `json:"staff_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Content string `json:"content"`
	Images  string `json:"images"`
}

func CreateReview(c *gin.Context) {
	uid, _ := c.Get("uid")
	var in ReviewInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	var order models.Order
	if err := database.DB.First(&order, in.OrderID).Error; err != nil {
		utils.NotFound(c, "order not found")
		return
	}
	if order.CustomerID != uid.(uint) {
		utils.Forbidden(c, "not your order")
		return
	}
	if order.Status != models.OrderPaid && order.Status != models.OrderConfirmed {
		utils.BadRequest(c, "order not confirmed")
		return
	}
	var existed models.Review
	if database.DB.Where("order_id = ?", in.OrderID).First(&existed).Error == nil {
		utils.BadRequest(c, "already reviewed")
		return
	}
	review := models.Review{
		OrderID:    in.OrderID,
		StaffID:    in.StaffID,
		CustomerID: uid.(uint),
		Rating:     in.Rating,
		Content:    in.Content,
		Images:     in.Images,
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&review).Error; err != nil {
			return err
		}
		var staff models.User
		if err := tx.First(&staff, in.StaffID).Error; err != nil {
			return err
		}
		newCount := staff.ReviewCount + 1
		newRating := (staff.Rating*float64(staff.ReviewCount) + float64(in.Rating)) / float64(newCount)
		if err := tx.Model(&staff).Updates(map[string]interface{}{
			"rating":       newRating,
			"review_count": newCount,
		}).Error; err != nil {
			return err
		}
		if in.Rating < 4 {
			if err := tx.Model(&staff).Update("suspended", true).Error; err != nil {
				return err
			}
		} else if newRating >= 4.5 && staff.Level < 5 {
			if err := tx.Model(&staff).Update("level", gorm.Expr("level + 1")).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		utils.ServerError(c, "create review failed")
		return
	}
	utils.Logger.Infow("review created", "order", in.OrderID, "rating", in.Rating)
	utils.OK(c, review)
}

func ListReviews(c *gin.Context) {
	var list []models.Review
	q := database.DB.Preload("Staff")
	if sid := c.Query("staff_id"); sid != "" {
		q = q.Where("staff_id = ?", sid)
	}
	if oid := c.Query("order_id"); oid != "" {
		q = q.Where("order_id = ?", oid)
	}
	if err := q.Order("id desc").Find(&list).Error; err != nil {
		utils.ServerError(c, "query failed")
		return
	}
	utils.OK(c, list)
}

func MyReviews(c *gin.Context) {
	uid, _ := c.Get("uid")
	role, _ := c.Get("role")
	var list []models.Review
	q := database.DB.Preload("Staff")
	switch role.(string) {
	case string(models.RoleCustomer):
		q = q.Where("customer_id = ?", uid)
	case string(models.RoleStaff):
		q = q.Where("staff_id = ?", uid)
	}
	if err := q.Order("id desc").Find(&list).Error; err != nil {
		utils.ServerError(c, "query failed")
		return
	}
	utils.OK(c, list)
}
