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

type RoastingHandler struct {
	cfg *config.Config
}

func NewRoastingHandler(cfg *config.Config) *RoastingHandler {
	return &RoastingHandler{cfg: cfg}
}

func (h *RoastingHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(h.cfg.App.PageSize)))
	productID := c.Query("product_id")
	roasterID := c.Query("roaster_id")
	batchNumber := c.Query("batch_number")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = h.cfg.App.PageSize
	}

	query := database.DB.Model(&models.RoastingRecord{}).Preload("Product").Preload("Roaster")

	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}
	if roasterID != "" {
		query = query.Where("roaster_id = ?", roasterID)
	}
	if batchNumber != "" {
		query = query.Where("batch_number LIKE ?", "%"+batchNumber+"%")
	}

	var total int64
	query.Count(&total)

	var records []models.RoastingRecord
	query.Order("roasted_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&records)

	utils.PaginatedResponse(c, records, total, page, pageSize)
}

func (h *RoastingHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var record models.RoastingRecord
	if err := database.DB.Preload("DataPoints").Preload("Product").Preload("Roaster").
		First(&record, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "烘焙记录不存在")
		return
	}

	utils.Success(c, record)
}

func (h *RoastingHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.CreateRoastingRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	var existingRecord models.RoastingRecord
	if database.DB.Where("batch_number = ?", req.BatchNumber).First(&existingRecord).Error == nil {
		utils.Error(c, http.StatusConflict, "批次号已存在")
		return
	}

	record := models.RoastingRecord{
		ProductID:       req.ProductID,
		RoasterID:       userID,
		BatchNumber:     req.BatchNumber,
		GreenBeanWeight: req.GreenBeanWeight,
		RoastedWeight:   req.RoastedWeight,
		InputTemp:       req.InputTemp,
		TurningPoint:    req.TurningPoint,
		TurningTime:     req.TurningTime,
		FirstCrackTemp:  req.FirstCrackTemp,
		FirstCrackTime:  req.FirstCrackTime,
		SecondCrackTemp: req.SecondCrackTemp,
		SecondCrackTime: req.SecondCrackTime,
		DropTemp:        req.DropTemp,
		TotalRoastTime:  req.TotalRoastTime,
		Notes:           req.Notes,
		RoastedAt:       req.RoastedAt,
	}

	if err := database.DB.Create(&record).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "创建失败")
		return
	}

	if len(req.DataPoints) > 0 {
		for _, dp := range req.DataPoints {
			dataPoint := models.RoastingDataPoint{
				RoastingRecordID: record.ID,
				TimeElapsed:      dp.TimeElapsed,
				BeanTemp:         dp.BeanTemp,
				EnvTemp:          dp.EnvTemp,
				RateOfRise:       dp.RateOfRise,
			}
			database.DB.Create(&dataPoint)
		}
	}

	database.DB.Preload("DataPoints").First(&record, record.ID)
	utils.SuccessWithMessage(c, "创建成功", record)
}

func (h *RoastingHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	var record models.RoastingRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "烘焙记录不存在")
		return
	}

	if record.RoasterID != userID && userRole != "admin" {
		utils.Error(c, http.StatusForbidden, "无权修改此记录")
		return
	}

	var req models.CreateRoastingRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	updates := map[string]interface{}{
		"product_id":       req.ProductID,
		"batch_number":     req.BatchNumber,
		"green_bean_weight": req.GreenBeanWeight,
		"roasted_weight":   req.RoastedWeight,
		"input_temp":       req.InputTemp,
		"turning_point":    req.TurningPoint,
		"turning_time":     req.TurningTime,
		"first_crack_temp": req.FirstCrackTemp,
		"first_crack_time": req.FirstCrackTime,
		"second_crack_temp": req.SecondCrackTemp,
		"second_crack_time": req.SecondCrackTime,
		"drop_temp":        req.DropTemp,
		"total_roast_time": req.TotalRoastTime,
		"notes":            req.Notes,
		"roasted_at":       req.RoastedAt,
	}

	if err := database.DB.Model(&record).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "更新失败")
		return
	}

	if len(req.DataPoints) > 0 {
		database.DB.Where("roasting_record_id = ?", id).Delete(&models.RoastingDataPoint{})
		for _, dp := range req.DataPoints {
			dataPoint := models.RoastingDataPoint{
				RoastingRecordID: record.ID,
				TimeElapsed:      dp.TimeElapsed,
				BeanTemp:         dp.BeanTemp,
				EnvTemp:          dp.EnvTemp,
				RateOfRise:       dp.RateOfRise,
			}
			database.DB.Create(&dataPoint)
		}
	}

	database.DB.Preload("DataPoints").First(&record, id)
	utils.SuccessWithMessage(c, "更新成功", record)
}

func (h *RoastingHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	var record models.RoastingRecord
	if err := database.DB.First(&record, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "烘焙记录不存在")
		return
	}

	if record.RoasterID != userID && userRole != "admin" {
		utils.Error(c, http.StatusForbidden, "无权删除此记录")
		return
	}

	database.DB.Where("roasting_record_id = ?", id).Delete(&models.RoastingDataPoint{})
	if err := database.DB.Delete(&record).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *RoastingHandler) Compare(c *gin.Context) {
	var req struct {
		RecordIDs []uint `json:"record_ids" binding:"required,min=2,max=5"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorWithDetails(c, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	var records []models.RoastingRecord
	database.DB.Preload("DataPoints").Preload("Product").
		Where("id IN ?", req.RecordIDs).Find(&records)

	if len(records) < 2 {
		utils.Error(c, http.StatusBadRequest, "需要至少2条记录进行对比")
		return
	}

	var totalInputTemp, totalDropTemp, totalTime, totalWeightLoss float64
	for _, r := range records {
		totalInputTemp += r.InputTemp
		totalDropTemp += r.DropTemp
		totalTime += float64(r.TotalRoastTime)
		if r.GreenBeanWeight > 0 {
			lossRate := (r.GreenBeanWeight - r.RoastedWeight) / r.GreenBeanWeight * 100
			totalWeightLoss += lossRate
		}
	}

	n := float64(len(records))
	result := models.RoastingComparisonResult{
		Records:        records,
		AvgInputTemp:   totalInputTemp / n,
		AvgDropTemp:    totalDropTemp / n,
		AvgTotalTime:   totalTime / n,
		WeightLossRate: totalWeightLoss / n,
	}

	utils.Success(c, result)
}

func (h *RoastingHandler) GetStats(c *gin.Context) {
	roasterID := c.Query("roaster_id")
	productID := c.Query("product_id")

	query := database.DB.Model(&models.RoastingRecord{})
	if roasterID != "" {
		query = query.Where("roaster_id = ?", roasterID)
	}
	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}

	var totalCount int64
	query.Count(&totalCount)

	var stats struct {
		TotalCount     int64   `json:"total_count"`
		AvgInputTemp   float64 `json:"avg_input_temp"`
		AvgDropTemp    float64 `json:"avg_drop_temp"`
		AvgTotalTime   float64 `json:"avg_total_time"`
		AvgWeightLoss  float64 `json:"avg_weight_loss"`
	}

	stats.TotalCount = totalCount

	var avgStats struct {
		AvgInputTemp  float64 `json:"avg_input_temp"`
		AvgDropTemp   float64 `json:"avg_drop_temp"`
		AvgTotalTime  float64 `json:"avg_total_time"`
	}
	query.Select("AVG(input_temp) as avg_input_temp, AVG(drop_temp) as avg_drop_temp, AVG(total_roast_time) as avg_total_time").
		Scan(&avgStats)

	stats.AvgInputTemp = avgStats.AvgInputTemp
	stats.AvgDropTemp = avgStats.AvgDropTemp
	stats.AvgTotalTime = avgStats.AvgTotalTime

	utils.Success(c, stats)
}
