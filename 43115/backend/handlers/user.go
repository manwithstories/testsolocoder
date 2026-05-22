package handlers

import (
	"strconv"

	"housekeeping-platform/config"
	"housekeeping-platform/models"
	"housekeeping-platform/utils"

	"github.com/gin-gonic/gin"
)

func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	role := c.Query("role")
	status := c.Query("status")
	keyword := c.Query("keyword")

	query := config.DB.Model(&models.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if status != "" {
		query = query.Where("provider_status = ?", status)
	}
	if keyword != "" {
		query = query.Where("nickname LIKE ? OR phone LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var users []models.User
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users)

	utils.Success(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"list":      users,
	})
}

func GetUserDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if result := config.DB.Preload("Addresses").Preload("Certifications").First(&user, id); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func ToggleUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if result := config.DB.First(&user, id); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	user.IsActive = !user.IsActive
	config.DB.Save(&user)

	operatorID := c.GetUint("user_id")
	operatorRole := c.GetString("role")
	uid := uint(id)
	go utils.LogOperation(operatorID, operatorRole, "user", "toggle_status", &uid, "user", "切换用户状态", c.ClientIP(), c.Request.UserAgent())

	utils.SuccessWithMessage(c, "操作成功", gin.H{"is_active": user.IsActive})
}

func GetAddressList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var addresses []models.Address
	config.DB.Where("user_id = ?", userID).Order("is_default DESC, id DESC").Find(&addresses)

	utils.Success(c, addresses)
}

func CreateAddress(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		ContactName  string  `json:"contact_name" binding:"required"`
		ContactPhone string  `json:"contact_phone" binding:"required"`
		Province     string  `json:"province" binding:"required"`
		City         string  `json:"city" binding:"required"`
		District     string  `json:"district" binding:"required"`
		Address      string  `json:"address" binding:"required"`
		Longitude    float64 `json:"longitude"`
		Latitude     float64 `json:"latitude"`
		IsDefault    bool    `json:"is_default"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if !utils.IsValidPhone(req.ContactPhone) {
		utils.BadRequest(c, "联系电话格式错误")
		return
	}

	if req.IsDefault {
		config.DB.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", false)
	}

	address := models.Address{
		UserID:       userID,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		IsDefault:    req.IsDefault,
	}

	if result := config.DB.Create(&address); result.Error != nil {
		utils.InternalError(c, "创建失败")
		return
	}

	utils.Success(c, address)
}

func UpdateAddress(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var address models.Address
	if result := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&address); result.Error != nil {
		utils.NotFound(c, "地址不存在")
		return
	}

	var req struct {
		ContactName  *string  `json:"contact_name"`
		ContactPhone *string  `json:"contact_phone"`
		Province     *string  `json:"province"`
		City         *string  `json:"city"`
		District     *string  `json:"district"`
		Address      *string  `json:"address"`
		Longitude    *float64 `json:"longitude"`
		Latitude     *float64 `json:"latitude"`
		IsDefault    *bool    `json:"is_default"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if req.IsDefault != nil && *req.IsDefault {
		config.DB.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", false)
	}

	updates := map[string]interface{}{}
	if req.ContactName != nil {
		updates["contact_name"] = *req.ContactName
	}
	if req.ContactPhone != nil {
		updates["contact_phone"] = *req.ContactPhone
	}
	if req.Province != nil {
		updates["province"] = *req.Province
	}
	if req.City != nil {
		updates["city"] = *req.City
	}
	if req.District != nil {
		updates["district"] = *req.District
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.Longitude != nil {
		updates["longitude"] = *req.Longitude
	}
	if req.Latitude != nil {
		updates["latitude"] = *req.Latitude
	}
	if req.IsDefault != nil {
		updates["is_default"] = *req.IsDefault
	}

	config.DB.Model(&address).Updates(updates)

	utils.Success(c, nil)
}

func DeleteAddress(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if result := config.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Address{}); result.Error != nil {
		utils.InternalError(c, "删除失败")
		return
	}

	utils.Success(c, nil)
}

func SetDefaultAddress(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	config.DB.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", false)
	config.DB.Model(&models.Address{}).Where("id = ? AND user_id = ?", id, userID).Update("is_default", true)

	utils.Success(c, nil)
}

func SubmitCertification(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		RealName          string `json:"real_name" binding:"required"`
		IDCard            string `json:"id_card" binding:"required"`
		CertificationDesc string `json:"certification_desc"`
		Certifications    []struct {
			CertType   string `json:"cert_type" binding:"required"`
			CertName   string `json:"cert_name" binding:"required"`
			CertNumber string `json:"cert_number"`
			CertImage  string `json:"cert_image"`
		} `json:"certifications"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	if !utils.IsValidIDCard(req.IDCard) {
		utils.BadRequest(c, "身份证号格式错误")
		return
	}

	tx := config.DB.Begin()

	if err := tx.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"real_name":          req.RealName,
		"id_card":            req.IDCard,
		"certification_desc": req.CertificationDesc,
		"provider_status":    models.ProviderStatusPending,
	}).Error; err != nil {
		tx.Rollback()
		utils.InternalError(c, "提交失败")
		return
	}

	for _, cert := range req.Certifications {
		certification := models.ServiceProviderCert{
			ProviderID: userID,
			CertType:   cert.CertType,
			CertName:   cert.CertName,
			CertNumber: cert.CertNumber,
			CertImage:  cert.CertImage,
		}
		if err := tx.Create(&certification).Error; err != nil {
			tx.Rollback()
			utils.InternalError(c, "提交失败")
			return
		}
	}

	tx.Commit()

	go utils.LogOperation(userID, "service_provider", "provider", "submit_certification", &userID, "user", "提交认证申请", c.ClientIP(), c.Request.UserAgent())

	utils.SuccessWithMessage(c, "认证申请已提交，请等待审核", nil)
}

func ReviewCertification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Approved     bool   `json:"approved"`
		RejectReason string `json:"reject_reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	var user models.User
	if result := config.DB.First(&user, id); result.Error != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	status := models.ProviderStatusApproved
	rejectReason := ""
	if !req.Approved {
		status = models.ProviderStatusRejected
		rejectReason = req.RejectReason
	}

	updates := map[string]interface{}{
		"provider_status": status,
	}
	if rejectReason != "" {
		updates["reject_reason"] = rejectReason
	}

	config.DB.Model(&user).Updates(updates)

	if req.Approved {
		utils.SendSystemNotification(user.ID, "认证审核通过", "恭喜您的认证申请已通过审核，可以开始接单服务了")
	} else {
		utils.SendSystemNotification(user.ID, "认证审核未通过", "很遗憾您的认证申请未通过审核，原因："+rejectReason)
	}

	operatorID := c.GetUint("user_id")
	operatorRole := c.GetString("role")
	uid := uint(id)
	action := "approve_certification"
	if !req.Approved {
		action = "reject_certification"
	}
	go utils.LogOperation(operatorID, operatorRole, "provider", action, &uid, "user", "审核认证申请", c.ClientIP(), c.Request.UserAgent())

	utils.SuccessWithMessage(c, "操作成功", nil)
}
