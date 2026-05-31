package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"pet-board/internal/dto"
	"pet-board/internal/service"
	"pet-board/internal/utils"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.userService.Register(&req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.userService.Login(&req)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.Success(c, resp)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	profile, err := h.userService.GetProfile(uid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, profile)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.UpdateProfile(uid, &req); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) UpdateStoreInfo(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var req dto.UpdateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.UpdateStoreInfo(uid, &req); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) UpdateKeeperInfo(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var req dto.UpdateKeeperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.UpdateKeeperInfo(uid, &req); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.ChangePassword(uid, &req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page := utils.AtoiOrZero(c.DefaultQuery("page", "1"))
	pageSize := utils.AtoiOrZero(c.DefaultQuery("page_size", "10"))
	role := c.Query("role")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	users, total, err := h.userService.ListUsers(page, pageSize, role)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, dto.NewPagedResult(total, page, pageSize, users))
}

type PetHandler struct {
	petService *service.PetService
}

func NewPetHandler(petService *service.PetService) *PetHandler {
	return &PetHandler{petService: petService}
}

func (h *PetHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var req dto.PetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	pet, err := h.petService.Create(uid, &req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, pet)
}

func (h *PetHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	petID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid pet ID")
		return
	}

	pet, err := h.petService.GetByID(petID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, pet)
}

func (h *PetHandler) ListByOwner(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	page := utils.AtoiOrZero(c.DefaultQuery("page", "1"))
	pageSize := utils.AtoiOrZero(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	pets, total, err := h.petService.ListByOwner(uid, page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, dto.NewPagedResult(total, page, pageSize, pets))
}

func (h *PetHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	id := c.Param("id")
	petID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid pet ID")
		return
	}

	var req dto.PetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.petService.Update(petID, uid, &req); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *PetHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	id := c.Param("id")
	petID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid pet ID")
		return
	}

	if err := h.petService.Delete(petID, uid); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *PetHandler) AddVaccineRecord(c *gin.Context) {
	var req dto.VaccineRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	record, err := h.petService.AddVaccineRecord(&req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, record)
}

func (h *PetHandler) GetVaccineRecords(c *gin.Context) {
	id := c.Param("id")
	petID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid pet ID")
		return
	}

	records, err := h.petService.GetVaccineRecords(petID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, records)
}

func (h *PetHandler) AddDewormRecord(c *gin.Context) {
	var req dto.DewormRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	record, err := h.petService.AddDewormRecord(&req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, record)
}

func (h *PetHandler) GetDewormRecords(c *gin.Context) {
	id := c.Param("id")
	petID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid pet ID")
		return
	}

	records, err := h.petService.GetDewormRecords(petID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, records)
}

type PackageHandler struct {
	pkgService *service.PackageService
}

func NewPackageHandler(pkgService *service.PackageService) *PackageHandler {
	return &PackageHandler{pkgService: pkgService}
}

func (h *PackageHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var req dto.BoardingPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	pkg, err := h.pkgService.Create(uid, &req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, pkg)
}

func (h *PackageHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	pkgID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid package ID")
		return
	}

	pkg, err := h.pkgService.GetByID(pkgID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, pkg)
}

func (h *PackageHandler) ListByStore(c *gin.Context) {
	storeID := c.Query("store_id")
	uid, _ := uuid.Parse(storeID)

	pkgType := c.Query("type")
	page := utils.AtoiOrZero(c.DefaultQuery("page", "1"))
	pageSize := utils.AtoiOrZero(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	packages, total, err := h.pkgService.ListByStore(uid, pkgType, page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, dto.NewPagedResult(total, page, pageSize, packages))
}

func (h *PackageHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	id := c.Param("id")
	pkgID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid package ID")
		return
	}

	var req dto.BoardingPackageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.pkgService.Update(pkgID, uid, &req); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *PackageHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	id := c.Param("id")
	pkgID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequest(c, "invalid package ID")
		return
	}

	if err := h.pkgService.Delete(pkgID, uid); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}
