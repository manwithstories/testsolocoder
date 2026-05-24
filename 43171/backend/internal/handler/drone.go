package handler

import (
	"drone-rental/internal/dto"
	"drone-rental/internal/middleware"
	"drone-rental/internal/model"
	"drone-rental/internal/pkg/response"
	"drone-rental/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type DroneHandler struct {
	droneService *service.DroneService
}

func NewDroneHandler() *DroneHandler {
	return &DroneHandler{
		droneService: service.NewDroneService(),
	}
}

func (h *DroneHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.CreateDroneReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	drone, err := h.droneService.Create(userID, &req)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, drone)
}

func (h *DroneHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	drone, err := h.droneService.GetByID(uint(id))
	if err != nil {
		response.ErrNotFound(c, "设备不存在")
		return
	}
	response.Success(c, drone)
}

func (h *DroneHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.UpdateDroneReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.droneService.Update(uint(id), userID, &req); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *DroneHandler) UpdateStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.UpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	if err := h.droneService.UpdateStatus(uint(id), userID, req.Status); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *DroneHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.droneService.Delete(uint(id), userID); err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *DroneHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	ownerID, _ := strconv.ParseUint(c.Query("owner_id"), 10, 64)
	keyword := c.Query("keyword")
	region := c.Query("region")
	brand := c.Query("brand")
	status := model.DroneStatus(c.Query("status"))
	drones, total, err := h.droneService.List(page, pageSize, uint(ownerID), keyword, region, brand, status)
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Page(c, drones, total, page, pageSize)
}

func (h *DroneHandler) MyDrones(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")
	region := c.Query("region")
	status := model.DroneStatus(c.Query("status"))
	drones, total, err := h.droneService.List(page, pageSize, userID, keyword, region, "", status)
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Page(c, drones, total, page, pageSize)
}

func (h *DroneHandler) SearchAvailable(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	region := c.Query("region")
	keyword := c.Query("keyword")
	minPrice, _ := strconv.ParseFloat(c.DefaultQuery("min_price", "0"), 64)
	maxPrice, _ := strconv.ParseFloat(c.DefaultQuery("max_price", "0"), 64)
	drones, total, err := h.droneService.SearchAvailable(startDate, endDate, region, keyword, minPrice, maxPrice, page, pageSize)
	if err != nil {
		response.ErrServer(c, err.Error())
		return
	}
	response.Page(c, drones, total, page, pageSize)
}

func (h *DroneHandler) BatchImport(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req dto.BatchImportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrParam(c, err.Error())
		return
	}
	count, err := h.droneService.BatchImport(userID, req.Drones)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, gin.H{"imported": count})
}

func (h *DroneHandler) BatchImportExcel(c *gin.Context) {
	userID := middleware.GetUserID(c)
	
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrParam(c, "请上传Excel文件")
		return
	}
	
	f, err := file.Open()
	if err != nil {
		response.Fail(c, 500, "文件打开失败")
		return
	}
	defer f.Close()
	
	excelFile, err := excelize.OpenReader(f)
	if err != nil {
		response.Fail(c, 500, "Excel文件解析失败")
		return
	}
	defer excelFile.Close()
	
	sheetName := excelFile.GetSheetName(0)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		response.Fail(c, 500, "读取Excel行失败")
		return
	}
	
	if len(rows) < 2 {
		response.Fail(c, 400, "Excel文件至少需要包含表头和一行数据")
		return
	}
	
	var drones []dto.CreateDroneReq
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 6 {
			continue
		}
		
		drone := dto.CreateDroneReq{
			Name:        getCellValue(row, 0),
			Brand:       getCellValue(row, 1),
			Model:       getCellValue(row, 2),
			SerialNo:    getCellValue(row, 3),
			Region:      getCellValue(row, 4),
		}
		
		if price, err := strconv.ParseFloat(getCellValue(row, 5), 64); err == nil {
			drone.PricePerDay = price
		}
		if len(row) > 6 {
			if deposit, err := strconv.ParseFloat(getCellValue(row, 6), 64); err == nil {
				drone.Deposit = deposit
			}
		}
		if len(row) > 7 {
			drone.Description = getCellValue(row, 7)
		}
		
		if drone.Name != "" && drone.SerialNo != "" && drone.Region != "" && drone.PricePerDay > 0 {
			drones = append(drones, drone)
		}
	}
	
	count, err := h.droneService.BatchImport(userID, drones)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}
	response.Success(c, gin.H{"imported": count})
}

func getCellValue(row []string, index int) string {
	if index < len(row) {
		return row[index]
	}
	return ""
}

func (h *DroneHandler) Online(c *gin.Context) {
	h.UpdateStatus(c)
}
