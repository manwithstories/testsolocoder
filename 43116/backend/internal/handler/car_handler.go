package handler

import (
	"car-rental/internal/config"
	"car-rental/internal/model"
	"car-rental/internal/service"
	"car-rental/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CarHandler struct {
	carService *service.CarService
	cfg        *config.Config
}

func NewCarHandler(cfg *config.Config) *CarHandler {
	return &CarHandler{
		carService: service.NewCarService(cfg),
		cfg:        cfg,
	}
}

func (h *CarHandler) CreateCar(c *gin.Context) {
	var req service.CreateCarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	car, err := h.carService.CreateCar(&req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, car)
}

func (h *CarHandler) GetCarByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	car, err := h.carService.GetCarByID(uint(id))
	if err != nil {
		utils.NotFound(c, "车辆不存在")
		return
	}

	utils.Success(c, car)
}

func (h *CarHandler) GetAllCars(c *gin.Context) {
	page, pageSize, _, _ := utils.ParsePageParams(c)
	keyword := c.Query("keyword")
	status := c.Query("status")
	brand := c.Query("brand")
	storeID, _ := strconv.ParseUint(c.Query("store_id"), 10, 64)

	cars, total, err := h.carService.GetAllCars(page, pageSize, keyword, status, brand, uint(storeID))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, cars, total, page, pageSize)
}

func (h *CarHandler) UpdateCar(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.carService.UpdateCar(uint(id), updates)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *CarHandler) UpdateCarStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req struct {
		Status model.CarStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	err := h.carService.UpdateStatus(uint(id), req.Status)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *CarHandler) DeleteCar(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := h.carService.DeleteCar(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *CarHandler) UploadCarImage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请上传文件")
		return
	}

	if file.Size > h.cfg.Upload.MaxSize {
		utils.BadRequest(c, "文件大小超出限制")
		return
	}

	image, err := h.carService.UploadCarImage(uint(id), file)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, image)
}

func (h *CarHandler) BatchUploadCarImages(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	form, err := c.MultipartForm()
	if err != nil {
		utils.BadRequest(c, "请上传文件")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		utils.BadRequest(c, "请上传至少一个文件")
		return
	}

	for _, file := range files {
		if file.Size > h.cfg.Upload.MaxSize {
			utils.BadRequest(c, "文件大小超出限制: "+file.Filename)
			return
		}
	}

	images, err := h.carService.BatchUploadCarImages(uint(id), files)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, images)
}

func (h *CarHandler) DeleteCarImage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	imageID, _ := strconv.ParseUint(c.Param("imageId"), 10, 64)
	_ = id

	err := h.carService.DeleteCarImage(uint(imageID))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *CarHandler) GetCarImages(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	images, err := h.carService.GetCarImages(uint(id))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, images)
}

func (h *CarHandler) GetAvailableCars(c *gin.Context) {
	page, pageSize, _, _ := utils.ParsePageParams(c)
	storeID, _ := strconv.ParseUint(c.Query("store_id"), 10, 64)

	cars, total, err := h.carService.GetAvailableCars(uint(storeID), page, pageSize)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithPage(c, cars, total, page, pageSize)
}
