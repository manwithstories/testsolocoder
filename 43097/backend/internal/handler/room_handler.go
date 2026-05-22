package handler

import (
	"hotel-system/internal/dto"
	"hotel-system/internal/model"
	"hotel-system/internal/pkg/logger"
	"hotel-system/internal/service"
	"hotel-system/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomService service.RoomService
}

func NewRoomHandler(roomService service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

func (h *RoomHandler) CreateRoomType(c *gin.Context) {
	var req dto.RoomTypeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("创建房型参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	roomType, err := h.roomService.CreateRoomType(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	roomTypeResponse := dto.RoomTypeResponse{
		ID:          roomType.ID,
		Name:        roomType.Name,
		Description: roomType.Description,
		BasePrice:   roomType.BasePrice,
		BedCount:    roomType.BedCount,
		MaxGuests:   roomType.MaxGuests,
		Facilities:  []string(roomType.Facilities),
		CreatedAt:   roomType.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   roomType.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.Success(c, roomTypeResponse)
}

func (h *RoomHandler) GetRoomType(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的房型ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的房型ID")
		return
	}

	roomType, err := h.roomService.GetRoomType(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	roomTypeResponse := dto.RoomTypeResponse{
		ID:          roomType.ID,
		Name:        roomType.Name,
		Description: roomType.Description,
		BasePrice:   roomType.BasePrice,
		BedCount:    roomType.BedCount,
		MaxGuests:   roomType.MaxGuests,
		Facilities:  []string(roomType.Facilities),
		CreatedAt:   roomType.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   roomType.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.Success(c, roomTypeResponse)
}

func (h *RoomHandler) UpdateRoomType(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的房型ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的房型ID")
		return
	}

	var req dto.RoomTypeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("更新房型参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	roomType, err := h.roomService.UpdateRoomType(uint(id), &req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	roomTypeResponse := dto.RoomTypeResponse{
		ID:          roomType.ID,
		Name:        roomType.Name,
		Description: roomType.Description,
		BasePrice:   roomType.BasePrice,
		BedCount:    roomType.BedCount,
		MaxGuests:   roomType.MaxGuests,
		Facilities:  []string(roomType.Facilities),
		CreatedAt:   roomType.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   roomType.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.Success(c, roomTypeResponse)
}

func (h *RoomHandler) DeleteRoomType(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的房型ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的房型ID")
		return
	}

	err = h.roomService.DeleteRoomType(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *RoomHandler) ListRoomTypes(c *gin.Context) {
	roomTypes, err := h.roomService.ListRoomTypes()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	var roomTypeResponses []dto.RoomTypeResponse
	for _, rt := range roomTypes {
		roomTypeResponses = append(roomTypeResponses, dto.RoomTypeResponse{
			ID:          rt.ID,
			Name:        rt.Name,
			Description: rt.Description,
			BasePrice:   rt.BasePrice,
			BedCount:    rt.BedCount,
			MaxGuests:   rt.MaxGuests,
			Facilities:  []string(rt.Facilities),
			CreatedAt:   rt.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   rt.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	utils.Success(c, roomTypeResponses)
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req dto.RoomCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("创建房间参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	room, err := h.roomService.CreateRoom(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	roomResponse := convertToRoomResponse(room)
	utils.Success(c, roomResponse)
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的房间ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的房间ID")
		return
	}

	room, err := h.roomService.GetRoom(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	roomResponse := convertToRoomResponse(room)
	utils.Success(c, roomResponse)
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的房间ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的房间ID")
		return
	}

	var req dto.RoomUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("更新房间参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	room, err := h.roomService.UpdateRoom(uint(id), &req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	roomResponse := convertToRoomResponse(room)
	utils.Success(c, roomResponse)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的房间ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的房间ID")
		return
	}

	err = h.roomService.DeleteRoom(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *RoomHandler) ListRooms(c *gin.Context) {
	var req dto.RoomListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warnf("获取房间列表参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	rooms, total, err := h.roomService.ListRooms(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	var roomResponses []dto.RoomResponse
	for _, room := range rooms {
		roomResponses = append(roomResponses, convertToRoomResponse(&room))
	}

	utils.PageResult(c, roomResponses, total, req.GetPage(), req.GetPageSize())
}

func (h *RoomHandler) BatchImportRooms(c *gin.Context) {
	var req dto.RoomBatchImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("批量导入房间参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := h.roomService.BatchImportRooms(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, result)
}

func (h *RoomHandler) GetAvailableRooms(c *gin.Context) {
	roomTypeIDStr := c.Query("room_type_id")
	var roomTypeID uint
	if roomTypeIDStr != "" {
		id, err := strconv.ParseUint(roomTypeIDStr, 10, 32)
		if err != nil {
			logger.Warnf("无效的房型ID: %s", roomTypeIDStr)
			utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的房型ID")
			return
		}
		roomTypeID = uint(id)
	}

	rooms, err := h.roomService.GetAvailableRooms(roomTypeID)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	var roomResponses []dto.RoomResponse
	for _, room := range rooms {
		roomResponses = append(roomResponses, convertToRoomResponse(&room))
	}

	utils.Success(c, roomResponses)
}

func (h *RoomHandler) UpdateRoomStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的房间ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的房间ID")
		return
	}

	var req dto.RoomStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("更新房间状态参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	err = h.roomService.UpdateRoomStatus(uint(id), req.Status)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *RoomHandler) GetRoomDashboard(c *gin.Context) {
	dashboard, err := h.roomService.GetRoomDashboard()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, dashboard)
}

func convertToRoomResponse(room *model.Room) dto.RoomResponse {
	resp := dto.RoomResponse{
		ID:         room.ID,
		RoomNo:     room.RoomNo,
		Floor:      room.Floor,
		RoomTypeID: room.RoomTypeID,
		Status:     room.Status,
		Price:      room.Price,
		Facilities: []string(room.Facilities),
		CreatedAt:  room.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  room.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if room.RoomType != nil {
		resp.RoomType = &dto.RoomTypeResponse{
			ID:          room.RoomType.ID,
			Name:        room.RoomType.Name,
			Description: room.RoomType.Description,
			BasePrice:   room.RoomType.BasePrice,
			BedCount:    room.RoomType.BedCount,
			MaxGuests:   room.RoomType.MaxGuests,
			Facilities:  []string(room.RoomType.Facilities),
			CreatedAt:   room.RoomType.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   room.RoomType.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return resp
}
