package handler

import (
	"hotel-system/internal/pkg/logger"
	"hotel-system/internal/service"
	"hotel-system/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService service.DashboardService
}

func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetRoomStatusBoard(c *gin.Context) {
	board, err := h.dashboardService.GetRoomStatusBoard()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, board)
}

func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.dashboardService.GetDashboardStats()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, stats)
}

func (h *DashboardHandler) GetFloorRooms(c *gin.Context) {
	floorStr := c.Param("floor")
	floor, err := strconv.Atoi(floorStr)
	if err != nil {
		logger.Warnf("无效的楼层参数: %s", floorStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的楼层参数")
		return
	}

	if floor <= 0 {
		logger.Warnf("楼层参数必须大于0: %d", floor)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "楼层参数必须大于0")
		return
	}

	result, err := h.dashboardService.GetFloorRooms(floor)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, result)
}
