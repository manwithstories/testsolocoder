package handlers

import (
	"garden-planner/database"
	"garden-planner/middleware"
	"garden-planner/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateDiseaseDiagnosisRequest struct {
	PlantName   string `json:"plant_name"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
	Symptoms    string `json:"symptoms"`
}

type UpdateDiseaseDiagnosisRequest struct {
	Status      string `json:"status"`
	Diagnosis   string `json:"diagnosis"`
	Severity    string `json:"severity"`
	Treatment   string `json:"treatment"`
	Confidence  float64 `json:"confidence"`
}

var diseaseDatabase = []map[string]string{
	{
		"name":        "白粉病",
		"symptoms":    "叶片出现白色粉末状物质,叶片变黄,生长受阻",
		"severity":    "中等",
		"treatment":   "使用硫磺粉或专门的杀菌剂,保持通风,避免过度浇水",
		"keywords":    "白色粉末,白粉,叶片变白",
	},
	{
		"name":        "黑斑病",
		"symptoms":    "叶片出现黑色或棕色斑点,斑点逐渐扩大,叶片脱落",
		"severity":    "中等",
		"treatment":   "及时清除病叶,使用多菌灵或甲基托布津防治,避免叶面积水",
		"keywords":    "黑斑,黑点,棕色斑点,叶片斑点",
	},
	{
		"name":        "蚜虫",
		"symptoms":    "叶片卷曲,出现粘稠分泌物(蜜露),叶片变黄,生长缓慢",
		"severity":    "轻微",
		"treatment":   "使用肥皂水喷洒,或使用吡虫啉等杀虫剂,引入天敌如瓢虫",
		"keywords":    "蚜虫,小虫,卷曲,蜜露",
	},
	{
		"name":        "霜霉病",
		"symptoms":    "叶片出现黄色斑点,叶背有灰色霉层,叶片枯死",
		"severity":    "严重",
		"treatment":   "使用甲霜灵或霜脲氰防治,加强通风,避免夜间浇水",
		"keywords":    "霜霉,灰色霉层,黄色斑点,霉层",
	},
	{
		"name":        "根腐病",
		"symptoms":    "根部腐烂,植株萎蔫,叶片变黄,生长停滞",
		"severity":    "严重",
		"treatment":   "改善排水,使用恶霉灵灌根,减少浇水,更换土壤",
		"keywords":    "根腐,萎蔫,烂根,根部腐烂",
	},
	{
		"name":        "潜叶虫",
		"symptoms":    "叶片出现弯曲的隧道状痕迹,叶片变形,生长受阻",
		"severity":    "轻微",
		"treatment":   "清除受害叶片,使用阿维菌素防治,设置黄板诱捕",
		"keywords":    "潜叶,隧道,虫道,叶片痕迹",
	},
	{
		"name":        "叶斑病",
		"symptoms":    "叶片出现圆形或不规则褐色斑点,斑点中心灰白色",
		"severity":    "中等",
		"treatment":   "清除病叶,使用百菌清或代森锰锌防治,避免叶面积水",
		"keywords":    "叶斑,褐色斑点,圆形斑点",
	},
	{
		"name":        "缺氮症",
		"symptoms":    "老叶变黄,植株生长缓慢,叶片小而薄",
		"severity":    "轻微",
		"treatment":   "施用氮肥如尿素或腐熟有机肥,注意适量施用",
		"keywords":    "缺氮,黄叶,生长缓慢,叶片变黄",
	},
	{
		"name":        "缺磷症",
		"symptoms":    "叶片变紫或变红,植株矮小,根系发育不良",
		"severity":    "轻微",
		"treatment":   "施用磷肥如过磷酸钙,或使用骨粉等有机磷肥",
		"keywords":    "缺磷,紫色叶,红色叶,植株矮小",
	},
	{
		"name":        "缺钾症",
		"symptoms":    "叶片边缘焦枯,叶片出现斑点,果实品质下降",
		"severity":    "中等",
		"treatment":   "施用钾肥如硫酸钾或草木灰,注意平衡施肥",
		"keywords":    "缺钾,焦枯,叶缘枯,斑点",
	},
}

func CreateDiseaseDiagnosis(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateDiseaseDiagnosisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ImageURL == "" && req.Description == "" && req.Symptoms == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide image, description, or symptoms"})
		return
	}

	diagnosis := diagnoseDisease(req.Symptoms + " " + req.Description)

	record := models.DiseaseDiagnosis{
		ID:          uuid.New(),
		UserID:      userID,
		PlantName:   req.PlantName,
		ImageURL:    req.ImageURL,
		Description: req.Description,
		Symptoms:    req.Symptoms,
		Diagnosis:   diagnosis["name"],
		Severity:    diagnosis["severity"],
		Treatment:   diagnosis["treatment"],
		Confidence:  diagnosis["confidence"].(float64),
		Status:      "completed",
	}

	if err := database.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create diagnosis record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Diagnosis completed",
		"diagnosis": record,
	})
}

func GetDiseaseDiagnoses(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var diagnoses []models.DiseaseDiagnosis
	database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&diagnoses)

	c.JSON(http.StatusOK, gin.H{"diagnoses": diagnoses})
}

func GetDiseaseDiagnosis(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var diagnosis models.DiseaseDiagnosis
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&diagnosis).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diagnosis not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"diagnosis": diagnosis})
}

func UpdateDiseaseDiagnosis(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var diagnosis models.DiseaseDiagnosis
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&diagnosis).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diagnosis not found"})
		return
	}

	var req UpdateDiseaseDiagnosisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status != "" {
		diagnosis.Status = req.Status
	}
	if req.Diagnosis != "" {
		diagnosis.Diagnosis = req.Diagnosis
	}
	if req.Severity != "" {
		diagnosis.Severity = req.Severity
	}
	if req.Treatment != "" {
		diagnosis.Treatment = req.Treatment
	}
	if req.Confidence > 0 {
		diagnosis.Confidence = req.Confidence
	}

	if err := database.DB.Save(&diagnosis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update diagnosis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Diagnosis updated successfully",
		"diagnosis": diagnosis,
	})
}

func DeleteDiseaseDiagnosis(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.DiseaseDiagnosis{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete diagnosis"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diagnosis not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diagnosis deleted successfully"})
}

func diagnoseDisease(input string) map[string]interface{} {
	bestMatch := map[string]interface{}{
		"name":       "未知病虫害",
		"symptoms":   "无法确定具体问题,建议咨询园艺专家",
		"severity":   "未知",
		"treatment":  "建议观察植物状况,保持良好的栽培管理,必要时咨询专业人士",
		"confidence": 0.0,
	}

	maxConfidence := 0.0

	for _, disease := range diseaseDatabase {
		confidence := calculateMatchConfidence(input, disease["keywords"])
		if confidence > maxConfidence && confidence > 0.2 {
			maxConfidence = confidence
			bestMatch = map[string]interface{}{
				"name":       disease["name"],
				"symptoms":   disease["symptoms"],
				"severity":   disease["severity"],
				"treatment":  disease["treatment"],
				"confidence": confidence,
			}
		}
	}

	return bestMatch
}

func calculateMatchConfidence(input, keywords string) float64 {
	if keywords == "" {
		return 0
	}

	keywordList := splitKeywords(keywords)
	matches := 0

	for _, keyword := range keywordList {
		if containsString(input, keyword) {
			matches++
		}
	}

	if len(keywordList) == 0 {
		return 0
	}

	return float64(matches) / float64(len(keywordList))
}

func splitKeywords(s string) []string {
	var keywords []string
	current := ""
	for _, c := range s {
		if c == ',' {
			if current != "" {
				keywords = append(keywords, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		keywords = append(keywords, current)
	}
	return keywords
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
