package handlers

import (
	"net/http"
	"strconv"
	"time"

	"beehive-platform/database"
	"beehive-platform/models"
	"beehive-platform/utils"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func getSeason(month time.Month) string {
	switch month {
	case time.March, time.April, time.May:
		return "spring"
	case time.June, time.July, time.August:
		return "summer"
	case time.September, time.October, time.November:
		return "autumn"
	default:
		return "winter"
	}
}

func getSeasonalRecommendations(season, healthStatus string) string {
	recommendations := map[string]map[string]string{
		"spring": {
			"healthy": "春季蜜蜂活跃，检查蜂王状态，准备扩大蜂群，注意防螨。",
			"warning": "春季需加强巡查，及时处理轻微病虫害，补充花粉饲料。",
			"critical": "紧急处理病虫害，考虑隔离病群，联系专业人员指导。",
		},
		"summer": {
			"healthy": "夏季高温注意蜂箱通风遮阳，保证水源供应，适时采收蜂蜜。",
			"warning": "检查是否受热应激，增加喂水频率，减少开箱次数。",
			"critical": "立即移至阴凉处，增加通风，紧急处理病虫害。",
		},
		"autumn": {
			"healthy": "秋季为越冬做准备，检查蜂群储蜜，进行秋繁，防治病虫害。",
			"warning": "补充越冬饲料，加强螨类防治，确保蜂群健康越冬。",
			"critical": "集中处理病害，考虑合并弱群，确保越冬前恢复健康。",
		},
		"winter": {
			"healthy": "冬季保温防冻，定期检查越冬状态，减少开箱干扰。",
			"warning": "加强保温措施，检查储蜜是否充足，防止饥饿。",
			"critical": "紧急救治，补充饲料，加强保温，必要时室内越冬。",
		},
	}

	if recs, ok := recommendations[season]; ok {
		if rec, ok := recs[healthStatus]; ok {
			return rec
		}
	}
	return ""
}

func (h *HealthHandler) CreateRecord(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateHealthRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	var beehive models.Beehive
	if err := database.DB.Where("id = ? AND user_id = ?", req.BeehiveID, userID).First(&beehive).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "beehive not found")
		return
	}

	recordDate, err := time.Parse("2006-01-02", req.RecordDate)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "invalid date format")
		return
	}

	season := getSeason(recordDate.Month())
	if req.Season != "" {
		season = req.Season
	}

	healthStatus := "healthy"
	if req.HasDisease {
		if req.DiseaseSeverity == "severe" {
			healthStatus = "critical"
		} else {
			healthStatus = "warning"
		}
	}

	recommendations := getSeasonalRecommendations(season, healthStatus)

	record := models.HealthRecord{
		BeehiveID:       req.BeehiveID,
		RecordDate:      recordDate,
		QueenStatus:     req.QueenStatus,
		WorkerCount:     req.WorkerCount,
		HasDisease:      req.HasDisease,
		DiseaseType:     req.DiseaseType,
		DiseaseSeverity: req.DiseaseSeverity,
		Treatment:       req.Treatment,
		Season:          season,
		Recommendations: recommendations,
		Notes:           req.Notes,
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to begin transaction")
		return
	}

	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		utils.FailWithError(c, http.StatusInternalServerError, "failed to create health record", err)
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"health_status":  healthStatus,
		"queen_status":   req.QueenStatus,
		"worker_count":   req.WorkerCount,
		"last_inspection": &now,
	}
	if err := tx.Model(&beehive).Updates(updates).Error; err != nil {
		tx.Rollback()
		utils.FailWithError(c, http.StatusInternalServerError, "failed to update beehive", err)
		return
	}

	if req.HasDisease {
		notification := models.Notification{
			UserID:    userID.(uint),
			Type:      "disease_warning",
			Title:     "蜂箱病虫害预警",
			Content:   "蜂箱 " + beehive.Name + "(" + beehive.Code + ") 出现" + req.DiseaseType + "，严重程度：" + req.DiseaseSeverity + "。建议：" + recommendations,
			RelatedID: &beehive.ID,
		}
		tx.Create(&notification)
	} else {
		notification := models.Notification{
			UserID:    userID.(uint),
			Type:      "seasonal_tip",
			Title:     "季节性管理建议",
			Content:   recommendations,
			RelatedID: &beehive.ID,
		}
		tx.Create(&notification)
	}

	if err := tx.Commit().Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, "failed to commit transaction")
		return
	}

	utils.Success(c, record)
}

func (h *HealthHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	beehiveIDStr := c.Query("beehive_id")

	var pageParams utils.PageParams
	if err := c.ShouldBindQuery(&pageParams); err != nil {
		pageParams = utils.PageParams{Page: 1, PageSize: 10}
	}
	if pageParams.Page < 1 {
		pageParams.Page = 1
	}
	if pageParams.PageSize < 1 {
		pageParams.PageSize = 10
	}

	query := database.DB.Model(&models.HealthRecord{}).
		Joins("JOIN beehives ON health_records.beehive_id = beehives.id").
		Where("beehives.user_id = ?", userID)

	if beehiveIDStr != "" {
		beehiveID, _ := strconv.ParseUint(beehiveIDStr, 10, 64)
		query = query.Where("beehive_id = ?", beehiveID)
	}

	if hasDisease := c.Query("has_disease"); hasDisease != "" {
		query = query.Where("has_disease = ?", hasDisease == "true")
	}

	if season := c.Query("season"); season != "" {
		query = query.Where("season = ?", season)
	}

	sortBy := c.DefaultQuery("sort_by", "record_date")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	query = query.Order(sortBy + " " + sortOrder)

	var total int64
	query.Count(&total)

	var records []models.HealthRecord
	query.Offset(pageParams.GetOffset()).Limit(pageParams.PageSize).
		Preload("Beehive").Find(&records)

	utils.SuccessWithTotal(c, records, total)
}

func (h *HealthHandler) Get(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var record models.HealthRecord
	if err := database.DB.
		Joins("JOIN beehives ON health_records.beehive_id = beehives.id").
		Where("health_records.id = ? AND beehives.user_id = ?", id, userID).
		Preload("Beehive").
		First(&record).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "health record not found")
		return
	}

	utils.Success(c, record)
}

func (h *HealthHandler) GetDiseaseWarnings(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var warnings []models.DiseaseWarning
	database.DB.Model(&models.HealthRecord{}).
		Select("health_records.id, health_records.beehive_id, beehives.name as beehive_name, beehives.code as beehive_code, health_records.disease_type, health_records.disease_severity as severity, health_records.created_at as warning_time").
		Joins("JOIN beehives ON health_records.beehive_id = beehives.id").
		Where("beehives.user_id = ? AND health_records.has_disease = true", userID).
		Order("health_records.created_at DESC").
		Limit(20).
		Scan(&warnings)

	utils.Success(c, warnings)
}

func (h *HealthHandler) GetSeasonalTips(c *gin.Context) {
	season := c.Query("season")
	if season == "" {
		season = getSeason(time.Now().Month())
	}

	tips := map[string][]map[string]string{
		"spring": {
			{"title": "春季扩群", "content": "检查蜂王产卵情况，适时加脾扩巢"},
			{"title": "病害防治", "content": "早春重点防治蜂螨，可使用甲酸等药剂"},
			{"title": "饲料补充", "content": "蜜源不足时补充糖浆，刺激蜂王产卵"},
			{"title": "分蜂控制", "content": "注意分蜂热，及时加脾或人工分蜂"},
		},
		"summer": {
			{"title": "防暑降温", "content": "蜂箱遮阴，加强通风，设置喂水器"},
			{"title": "适时采蜜", "content": "大流蜜期及时取蜜，注意留足饲料"},
			{"title": "敌害防治", "content": "防止胡蜂、蚂蚁等敌害入侵"},
			{"title": "越夏管理", "content": "高温期减少开箱，保持安静"},
		},
		"autumn": {
			{"title": "越冬准备", "content": "培养越冬适龄蜂，紧缩蜂巢"},
			{"title": "饲料储备", "content": "留足越冬蜜，不足时补充糖浆"},
			{"title": "螨类防治", "content": "秋末断子期彻底治螨"},
			{"title": "蜂群调整", "content": "合并弱群，调整群势"},
		},
		"winter": {
			{"title": "保温防寒", "content": "做好箱内外保温，防止寒风侵入"},
			{"title": "定期检查", "content": "听测蜂群状态，防止闷死或冻死"},
			{"title": "饲料检查", "content": "确保越冬饲料充足，防止饥饿"},
			{"title": "湿度控制", "content": "防止箱内过湿，注意通风"},
		},
	}

	utils.Success(c, gin.H{
		"season": season,
		"tips":   tips[season],
	})
}
