package handler

import (
	"car-rental/internal/service"
	"car-rental/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeService *service.StoreService
}

func NewStoreHandler() *StoreHandler {
	return &StoreHandler{
		storeService: service.NewStoreService(),
	}
}

func (h *StoreHandler) CreateCity(c *gin.Context) {
	var req service.CreateCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	city, err := h.storeService.CreateCity(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, city)
}

func (h *StoreHandler) GetCityByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	city, err := h.storeService.GetCityByID(uint(id))
	if err != nil {
		utils.NotFound(c, "城市不存在")
		return
	}

	utils.Success(c, city)
}

func (h *StoreHandler) GetAllCities(c *gin.Context) {
	cities, err := h.storeService.GetAllCities()
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, cities)
}

func (h *StoreHandler) UpdateCity(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.storeService.UpdateCity(uint(id), updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *StoreHandler) DeleteCity(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.storeService.DeleteCity(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *StoreHandler) CreateStore(c *gin.Context) {
	var req service.CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	store, err := h.storeService.CreateStore(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, store)
}

func (h *StoreHandler) GetStoreByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	store, err := h.storeService.GetStoreByID(uint(id))
	if err != nil {
		utils.NotFound(c, "门店不存在")
		return
	}

	utils.Success(c, store)
}

func (h *StoreHandler) GetAllStores(c *gin.Context) {
	page, pageSize, _, _ := utils.ParsePageParams(c)
	cityID, _ := strconv.ParseUint(c.Query("city_id"), 10, 64)
	keyword := c.Query("keyword")

	stores, total, err := h.storeService.GetAllStores(page, pageSize, uint(cityID), keyword)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, stores, total, page, pageSize)
}

func (h *StoreHandler) UpdateStore(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.storeService.UpdateStore(uint(id), updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *StoreHandler) DeleteStore(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.storeService.DeleteStore(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *StoreHandler) GetStoresByCity(c *gin.Context) {
	cityID, _ := strconv.ParseUint(c.Param("cityId"), 10, 64)

	stores, err := h.storeService.GetStoresByCity(uint(cityID))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, stores)
}
