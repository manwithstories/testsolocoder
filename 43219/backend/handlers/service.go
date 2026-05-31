package handlers

import (
	"housekeeping/database"
	"housekeeping/models"
	"housekeeping/utils"

	"github.com/gin-gonic/gin"
)

type ServiceInput struct {
	Name     string  `json:"name" binding:"required"`
	Category string  `json:"category"`
	Desc     string  `json:"desc"`
	MinPrice float64 `json:"min_price" binding:"required"`
	MaxPrice float64 `json:"max_price" binding:"required"`
	Duration int     `json:"duration" binding:"required"`
	Skills   string  `json:"skills"`
}

func CreateService(c *gin.Context) {
	uid, _ := c.Get("uid")
	var in ServiceInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	svc := models.Service{
		CompanyID: uid.(uint),
		Name:      in.Name,
		Category:  in.Category,
		Desc:      in.Desc,
		MinPrice:  in.MinPrice,
		MaxPrice:  in.MaxPrice,
		Duration:  in.Duration,
		Skills:    in.Skills,
		Status:    "active",
	}
	if err := database.DB.Create(&svc).Error; err != nil {
		utils.ServerError(c, "create service failed")
		return
	}
	utils.OK(c, svc)
}

func ListServices(c *gin.Context) {
	var list []models.Service
	q := database.DB.Where("status = ?", "active")
	if cat := c.Query("category"); cat != "" {
		q = q.Where("category = ?", cat)
	}
	if cid := c.Query("company_id"); cid != "" {
		q = q.Where("company_id = ?", cid)
	}
	if kw := c.Query("keyword"); kw != "" {
		q = q.Where("name LIKE ? OR `desc` LIKE ?", "%"+kw+"%", "%"+kw+"%")
	}
	if err := q.Order("id desc").Find(&list).Error; err != nil {
		utils.ServerError(c, "query failed")
		return
	}
	utils.OK(c, list)
}

func GetService(c *gin.Context) {
	id := c.Param("id")
	var svc models.Service
	if err := database.DB.First(&svc, id).Error; err != nil {
		utils.NotFound(c, "service not found")
		return
	}
	utils.OK(c, svc)
}

func UpdateService(c *gin.Context) {
	id := c.Param("id")
	uid, _ := c.Get("uid")
	var svc models.Service
	if err := database.DB.First(&svc, id).Error; err != nil {
		utils.NotFound(c, "service not found")
		return
	}
	if svc.CompanyID != uid.(uint) {
		utils.Forbidden(c, "not owner")
		return
	}
	var in ServiceInput
	if err := c.ShouldBindJSON(&in); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	svc.Name = in.Name
	svc.Category = in.Category
	svc.Desc = in.Desc
	svc.MinPrice = in.MinPrice
	svc.MaxPrice = in.MaxPrice
	svc.Duration = in.Duration
	svc.Skills = in.Skills
	if err := database.DB.Save(&svc).Error; err != nil {
		utils.ServerError(c, "update failed")
		return
	}
	utils.OK(c, svc)
}

func DeleteService(c *gin.Context) {
	id := c.Param("id")
	uid, _ := c.Get("uid")
	var svc models.Service
	if err := database.DB.First(&svc, id).Error; err != nil {
		utils.NotFound(c, "service not found")
		return
	}
	if svc.CompanyID != uid.(uint) {
		utils.Forbidden(c, "not owner")
		return
	}
	svc.Status = "inactive"
	database.DB.Save(&svc)
	utils.OK(c, "ok")
}

func MyServices(c *gin.Context) {
	uid, _ := c.Get("uid")
	var list []models.Service
	if err := database.DB.Where("company_id = ?", uid).Order("id desc").Find(&list).Error; err != nil {
		utils.ServerError(c, "query failed")
		return
	}
	utils.OK(c, list)
}
