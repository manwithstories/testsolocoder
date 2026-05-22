package handlers

import (
	"strconv"

	"meeting-room/internal/middleware"
	"meeting-room/internal/services"
	"meeting-room/internal/utils"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomService *services.RoomService
}

func NewRoomHandler() *RoomHandler {
	return &RoomHandler{
		roomService: services.NewRoomService(),
	}
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req services.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	room, err := h.roomService.CreateRoom(&req, userID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, room)
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "会议室ID错误")
		return
	}

	room, err := h.roomService.GetRoom(uint(id))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, room)
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "会议室ID错误")
		return
	}

	var req services.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请求参数错误")
		return
	}

	room, err := h.roomService.UpdateRoom(uint(id), &req)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, room)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "会议室ID错误")
		return
	}

	err = h.roomService.DeleteRoom(uint(id))
	if err != nil {
		utils.InternalError(c, "删除会议室失败")
		return
	}

	utils.Success(c, nil)
}

func (h *RoomHandler) ListRooms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	floor := c.Query("floor")
	equipment := c.Query("equipment")

	rooms, total, err := h.roomService.ListRooms(page, pageSize, floor, equipment)
	if err != nil {
		utils.InternalError(c, "获取会议室列表失败")
		return
	}

	utils.Success(c, gin.H{
		"rooms": rooms,
		"total": total,
	})
}

func (h *RoomHandler) ListAllRooms(c *gin.Context) {
	rooms, err := h.roomService.ListAllRooms()
	if err != nil {
		utils.InternalError(c, "获取会议室列表失败")
		return
	}

	utils.Success(c, rooms)
}

func (h *RoomHandler) UploadPhoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "会议室ID错误")
		return
	}

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		utils.BadRequest(c, "请上传照片文件")
		return
	}
	defer file.Close()

	photo, err := h.roomService.UploadPhoto(uint(id), file, header.Filename)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, photo)
}

func (h *RoomHandler) DeletePhoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "照片ID错误")
		return
	}

	err = h.roomService.DeletePhoto(uint(id))
	if err != nil {
		utils.InternalError(c, "删除照片失败")
		return
	}

	utils.Success(c, nil)
}

func (h *RoomHandler) GetFloors(c *gin.Context) {
	floors, err := h.roomService.GetAllFloors()
	if err != nil {
		utils.InternalError(c, "获取楼层列表失败")
		return
	}

	utils.Success(c, floors)
}
