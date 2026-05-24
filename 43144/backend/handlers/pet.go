package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"pet-adoption-platform/config"
	"pet-adoption-platform/database"
	"pet-adoption-platform/models"
	"pet-adoption-platform/services"
	"pet-adoption-platform/utils"

	"github.com/gin-gonic/gin"
)

func CreatePet(c *gin.Context) {
	rescueID := c.GetUint("rescue_id")

	var req models.CreatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	pet, err := services.CreatePet(&req, rescueID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, pet)
}

func GetPet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid pet id")
		return
	}

	pet, err := services.GetPetByID(uint(id))
	if err != nil {
		utils.NotFound(c, "pet not found")
		return
	}

	utils.Success(c, pet)
}

func UpdatePet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid pet id")
		return
	}

	rescueID := c.GetUint("rescue_id")
	_, err = services.ValidatePetOwnership(uint(id), rescueID)
	if err != nil {
		utils.Forbidden(c, err.Error())
		return
	}

	var req models.UpdatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	pet, err := services.UpdatePet(uint(id), &req)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, pet)
}

func DeletePet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid pet id")
		return
	}

	rescueID := c.GetUint("rescue_id")
	_, err = services.ValidatePetOwnership(uint(id), rescueID)
	if err != nil {
		utils.Forbidden(c, err.Error())
		return
	}

	if err := services.DeletePet(uint(id)); err != nil {
		utils.InternalError(c, "failed to delete pet")
		return
	}

	utils.Success(c, gin.H{"message": "pet deleted successfully"})
}

func ListPets(c *gin.Context) {
	var query models.PetListQuery
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

	pets, total, err := services.ListPets(&query)
	if err != nil {
		utils.InternalError(c, "failed to list pets")
		return
	}

	utils.PaginatedSuccess(c, pets, total, query.Page, query.PageSize)
}

func UploadPetPhotos(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid pet id")
		return
	}

	rescueID := c.GetUint("rescue_id")
	_, err = services.ValidatePetOwnership(uint(id), rescueID)
	if err != nil {
		utils.Forbidden(c, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		utils.BadRequest(c, "failed to parse form")
		return
	}

	files := form.File["photos"]
	if len(files) == 0 {
		utils.BadRequest(c, "no files uploaded")
		return
	}

	cfg := config.Load()
	uploadDir := filepath.Join(cfg.UploadDir, "pets")
	os.MkdirAll(uploadDir, 0755)

	var paths []string
	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		if !strings.EqualFold(ext, ".jpg") && !strings.EqualFold(ext, ".jpeg") &&
			!strings.EqualFold(ext, ".png") && !strings.EqualFold(ext, ".gif") {
			continue
		}

		filename := "pet_" + strconv.FormatUint(uint64(id), 10) + "_" + strconv.FormatInt(int64(len(paths)), 10) + ext
		savePath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			continue
		}

		paths = append(paths, "/uploads/pets/"+filename)
	}

	photosStr := strings.Join(paths, ",")
	if err := services.UpdatePetPhotos(uint(id), photosStr); err != nil {
		utils.InternalError(c, "failed to update photos")
		return
	}

	utils.Success(c, gin.H{"photos": paths})
}

func UploadPetVideos(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid pet id")
		return
	}

	rescueID := c.GetUint("rescue_id")
	_, err = services.ValidatePetOwnership(uint(id), rescueID)
	if err != nil {
		utils.Forbidden(c, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		utils.BadRequest(c, "failed to parse form")
		return
	}

	files := form.File["videos"]
	if len(files) == 0 {
		utils.BadRequest(c, "no files uploaded")
		return
	}

	cfg := config.Load()
	uploadDir := filepath.Join(cfg.UploadDir, "pets")
	os.MkdirAll(uploadDir, 0755)

	var paths []string
	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		if !strings.EqualFold(ext, ".mp4") && !strings.EqualFold(ext, ".webm") &&
			!strings.EqualFold(ext, ".mov") {
			continue
		}

		filename := "pet_" + strconv.FormatUint(uint64(id), 10) + "_video_" + strconv.FormatInt(int64(len(paths)), 10) + ext
		savePath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			continue
		}

		paths = append(paths, "/uploads/pets/"+filename)
	}

	videosStr := strings.Join(paths, ",")
	if err := services.UpdatePetVideos(uint(id), videosStr); err != nil {
		utils.InternalError(c, "failed to update videos")
		return
	}

	utils.Success(c, gin.H{"videos": paths})
}

func GetPetAdoptionHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid pet id")
		return
	}

	history, err := services.GetPetAdoptionHistory(uint(id))
	if err != nil {
		utils.InternalError(c, "failed to get adoption history")
		return
	}

	utils.Success(c, history)
}

func UpdatePetStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid pet id")
		return
	}

	rescueID := c.GetUint("rescue_id")
	_, err = services.ValidatePetOwnership(uint(id), rescueID)
	if err != nil {
		utils.Forbidden(c, err.Error())
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=adoptable adopted treatment deceased"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := services.UpdatePetStatus(uint(id), models.PetStatus(req.Status)); err != nil {
		utils.InternalError(c, "failed to update status")
		return
	}

	utils.Success(c, gin.H{"message": "status updated"})
}

func GetMyPets(c *gin.Context) {
	role := c.GetString("role")

	var query models.PetListQuery
	c.ShouldBindQuery(&query)

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}

	if role == "rescue" {
		rescueID := c.GetUint("rescue_id")
		query.RescueID = rescueID
	} else if role == "adopter" {
		query.Search = ""
	}

	pets, total, err := services.ListPets(&query)
	if err != nil {
		utils.InternalError(c, "failed to list pets")
		return
	}

	utils.PaginatedSuccess(c, pets, total, query.Page, query.PageSize)
}

func GetMyAdoptedPets(c *gin.Context) {
	userID := c.GetUint("user_id")

	var pets []models.Pet
	if err := database.DB.Where("adopter_id = ?", userID).Find(&pets).Error; err != nil {
		utils.InternalError(c, "failed to get adopted pets")
		return
	}

	utils.Success(c, pets)
}

func GetUploadedFile(c *gin.Context) {
	filename := c.Param("filename")
	c.FileAttachment("./uploads/pets/"+filename, filename)
	c.Status(http.StatusOK)
}
