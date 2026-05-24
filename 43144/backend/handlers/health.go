package handlers

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"pet-adoption-platform/config"
	"pet-adoption-platform/models"
	"pet-adoption-platform/services"
	"pet-adoption-platform/utils"

	"github.com/gin-gonic/gin"
)

func CreateHealthRecord(c *gin.Context) {
	userID := c.GetUint("user_id")
	rescueID := c.GetUint("rescue_id")

	var req models.CreateHealthRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	record, err := services.CreateHealthRecord(&req, userID, rescueID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, record)
}

func GetHealthRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid record id")
		return
	}

	record, err := services.GetHealthRecordByID(uint(id))
	if err != nil {
		utils.NotFound(c, "record not found")
		return
	}

	utils.Success(c, record)
}

func ListHealthRecords(c *gin.Context) {
	var query models.HealthRecordListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}

	records, total, err := services.ListHealthRecords(&query)
	if err != nil {
		utils.InternalError(c, "failed to list records")
		return
	}

	utils.PaginatedSuccess(c, records, total, query.Page, query.PageSize)
}

func UpdateHealthRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid record id")
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	allowedFields := []string{"title", "description", "notes", "hospital", "vet_name", "next_date"}
	updates := make(map[string]interface{})
	for _, field := range allowedFields {
		if val, ok := req[field]; ok {
			updates[field] = val
		}
	}

	record, err := services.UpdateHealthRecord(uint(id), updates)
	if err != nil {
		utils.InternalError(c, "failed to update record")
		return
	}

	utils.Success(c, record)
}

func DeleteHealthRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid record id")
		return
	}

	if err := services.DeleteHealthRecord(uint(id)); err != nil {
		utils.InternalError(c, "failed to delete record")
		return
	}

	utils.Success(c, gin.H{"message": "record deleted"})
}

func UploadHealthReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid record id")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "no file uploaded")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" {
		utils.BadRequest(c, "only PDF files are allowed")
		return
	}

	cfg := config.Load()
	uploadDir := filepath.Join(cfg.UploadDir, "reports")
	os.MkdirAll(uploadDir, 0755)

	filename := "report_" + strconv.FormatUint(uint64(id), 10) + ext
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.InternalError(c, "failed to save file")
		return
	}

	if err := services.UpdateReportFile(uint(id), "/uploads/reports/"+filename); err != nil {
		utils.InternalError(c, "failed to update report")
		return
	}

	utils.Success(c, gin.H{"file": "/uploads/reports/" + filename})
}

func GetHealthReminders(c *gin.Context) {
	petID, err := strconv.ParseUint(c.Param("pet_id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid pet id")
		return
	}

	reminders, err := services.GetHealthReminders(uint(petID))
	if err != nil {
		utils.InternalError(c, "failed to get reminders")
		return
	}

	utils.Success(c, reminders)
}

func CompleteHealthReminder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid reminder id")
		return
	}

	if err := services.CompleteHealthReminder(uint(id)); err != nil {
		utils.InternalError(c, "failed to complete reminder")
		return
	}

	utils.Success(c, gin.H{"message": "reminder completed"})
}

func GetPetHealthSummary(c *gin.Context) {
	petID, err := strconv.ParseUint(c.Param("pet_id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid pet id")
		return
	}

	summary, err := services.GetPetHealthSummary(uint(petID))
	if err != nil {
		utils.InternalError(c, "failed to get health summary")
		return
	}

	utils.Success(c, summary)
}
