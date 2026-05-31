package handlers

import (
	"net/http"
	"strconv"

	"coffee-platform/config"
	"coffee-platform/database"
	"coffee-platform/models"
	"coffee-platform/utils"

	"github.com/gin-gonic/gin"
)

type CertificationHandler struct {
	cfg *config.Config
}

func NewCertificationHandler(cfg *config.Config) *CertificationHandler {
	return &CertificationHandler{cfg: cfg}
}

func (h *CertificationHandler) Apply(c *gin.Context) {
	userID := c.GetUint("user_id")

	var existingCert models.RoasterCertification
	if database.DB.Where("user_id = ?", userID).First(&existingCert).Error == nil {
		utils.Error(c, http.StatusConflict, "您已提交过认证申请")
		return
	}

	var req models.ApplyCertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	cert := models.RoasterCertification{
		UserID:     userID,
		CertName:   req.CertName,
		CertNumber: req.CertNumber,
		OrgName:    req.OrgName,
		CertFile:   req.CertFile,
		Experience: req.Experience,
		Specialty:  req.Specialty,
		Status:     models.CertStatusPending,
	}

	if err := database.DB.Create(&cert).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "提交失败")
		return
	}

	utils.SuccessWithMessage(c, "认证申请已提交，等待审核", cert)
}

func (h *CertificationHandler) UpdateApplication(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cert models.RoasterCertification
	if err := database.DB.Where("user_id = ?", userID).First(&cert).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "未找到认证申请")
		return
	}

	if cert.Status != models.CertStatusPending {
		utils.Error(c, http.StatusBadRequest, "只有待审核状态可以修改")
		return
	}

	var req models.ApplyCertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	updates := map[string]interface{}{
		"cert_name":   req.CertName,
		"cert_number": req.CertNumber,
		"org_name":    req.OrgName,
		"cert_file":   req.CertFile,
		"experience":  req.Experience,
		"specialty":   req.Specialty,
	}

	if err := database.DB.Model(&cert).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (h *CertificationHandler) GetMyCertification(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cert models.RoasterCertification
	if err := database.DB.Preload("User").Where("user_id = ?", userID).First(&cert).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "未找到认证申请")
		return
	}

	utils.Success(c, cert)
}

func (h *CertificationHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(h.cfg.App.PageSize)))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = h.cfg.App.PageSize
	}

	query := database.DB.Model(&models.RoasterCertification{}).Preload("User")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var certs []models.RoasterCertification
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&certs)

	utils.PaginatedResponse(c, certs, total, page, pageSize)
}

func (h *CertificationHandler) Review(c *gin.Context) {
	reviewerID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var cert models.RoasterCertification
	if err := database.DB.First(&cert, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "认证申请不存在")
		return
	}

	if cert.Status != models.CertStatusPending {
		utils.Error(c, http.StatusBadRequest, "该申请已被审核")
		return
	}

	var req models.ReviewCertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	now := timeNow()
	updates := map[string]interface{}{
		"status":         req.Status,
		"reviewer_id":    reviewerID,
		"review_comment": req.ReviewComment,
		"reviewed_at":    now,
	}

	if err := database.DB.Model(&cert).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "审核失败")
		return
	}

	if req.Status == models.CertStatusApproved {
		database.DB.Model(&models.User{}).Where("id = ?", cert.UserID).
			Updates(map[string]interface{}{
				"role":              models.RoleRoaster,
				"is_certified":      true,
				"certification_id":  cert.ID,
			})
	}

	utils.SuccessWithMessage(c, "审核完成", nil)
}

func (h *CertificationHandler) GetRoasterProfile(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var user models.User
	if err := database.DB.Preload("Certification").First(&user, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	if user.Role != models.RoleRoaster && !user.IsCertified {
		utils.Error(c, http.StatusBadRequest, "该用户不是认证烘焙师")
		return
	}

	var products []models.Product
	database.DB.Preload("Images").Where("roaster_id = ? AND status = ?", id, models.ProductStatusOnSale).
		Order("created_at DESC").Limit(20).Find(&products)

	var records []models.RoastingRecord
	database.DB.Preload("Product").Where("roaster_id = ?", id).
		Order("roasted_at DESC").Limit(10).Find(&records)

	var totalProducts, totalRoasts int64
	database.DB.Model(&models.Product{}).Where("roaster_id = ?", id).Count(&totalProducts)
	database.DB.Model(&models.RoastingRecord{}).Where("roaster_id = ?", id).Count(&totalRoasts)

	var avgScore float64
	database.DB.Table("cupping_scores cs").
		Joins("JOIN products p ON cs.product_id = p.id").
		Where("p.roaster_id = ?", id).
		Select("COALESCE(AVG(cs.overall_score), 0)").
		Scan(&avgScore)

	profile := models.RoasterProfile{
		User:           &user,
		Certification:  user.Certification,
		Products:       products,
		RoastingRecords: records,
		TotalProducts:  totalProducts,
		TotalRoasts:    totalRoasts,
		AvgScore:       avgScore,
	}

	utils.Success(c, profile)
}

func (h *CertificationHandler) ListCertifiedRoasters(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(h.cfg.App.PageSize)))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = h.cfg.App.PageSize
	}

	query := database.DB.Model(&models.User{}).
		Where("role = ? AND is_certified = ?", models.RoleRoaster, true)

	if keyword != "" {
		query = query.Where("nickname LIKE ? OR username LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var roasters []models.User
	query.Preload("Certification").Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&roasters)

	utils.PaginatedResponse(c, roasters, total, page, pageSize)
}

func timeNow() *time.Time {
	t := time.Now()
	return &t
}
