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

type CheckInHandler struct {
	checkInService service.CheckInService
}

func NewCheckInHandler(checkInService service.CheckInService) *CheckInHandler {
	return &CheckInHandler{checkInService: checkInService}
}

func (h *CheckInHandler) CreateCheckIn(c *gin.Context) {
	var req dto.CheckInCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("办理入住参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	checkIn, err := h.checkInService.CreateCheckIn(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	checkInResponse := convertToCheckInResponse(checkIn)
	utils.Success(c, checkInResponse)
}

func (h *CheckInHandler) GetCheckIn(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的入住ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的入住ID")
		return
	}

	checkIn, err := h.checkInService.GetCheckInDetail(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	checkInResponse := convertToCheckInResponse(checkIn)
	utils.Success(c, checkInResponse)
}

func (h *CheckInHandler) ListCheckIns(c *gin.Context) {
	var req dto.CheckInListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warnf("获取入住列表参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	checkIns, total, err := h.checkInService.ListCheckIns(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	var checkInResponses []dto.CheckInResponse
	for _, checkIn := range checkIns {
		checkInResponses = append(checkInResponses, convertToCheckInResponse(&checkIn))
	}

	utils.PageResult(c, checkInResponses, total, req.GetPage(), req.GetPageSize())
}

func (h *CheckInHandler) CheckOut(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的入住ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的入住ID")
		return
	}

	var req dto.CheckOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("退房参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	result, err := h.checkInService.CheckOut(uint(id), &req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	checkIn, err := h.checkInService.GetCheckInDetail(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	checkInResponse := convertToCheckInResponse(checkIn)
	result.CheckIn = &checkInResponse

	utils.Success(c, result)
}

func (h *CheckInHandler) ExtendStay(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的入住ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的入住ID")
		return
	}

	var req dto.ExtendStayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("续住参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	checkIn, err := h.checkInService.ExtendStay(uint(id), &req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	checkInResponse := convertToCheckInResponse(checkIn)
	utils.Success(c, checkInResponse)
}

func (h *CheckInHandler) AddExtraCharge(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的入住ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的入住ID")
		return
	}

	var req dto.ExtraChargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("添加额外消费参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	checkIn, err := h.checkInService.AddExtraCharge(uint(id), &req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	checkInResponse := convertToCheckInResponse(checkIn)
	utils.Success(c, checkInResponse)
}

func convertToCheckInResponse(checkIn *model.CheckIn) dto.CheckInResponse {
	resp := dto.CheckInResponse{
		ID:               checkIn.ID,
		CheckInNo:        checkIn.CheckInNo,
		BookingID:        checkIn.BookingID,
		RoomID:           checkIn.RoomID,
		GuestName:        checkIn.GuestName,
		GuestPhone:       checkIn.GuestPhone,
		GuestIDCard:      checkIn.GuestIDCard,
		CheckInTime:      checkIn.CheckInTime.Format("2006-01-02 15:04:05"),
		ExpectedCheckOut: checkIn.ExpectedCheckOut.Format("2006-01-02 15:04:05"),
		Status:           checkIn.Status,
		Deposit:          checkIn.Deposit,
		ExtraCharges:     checkIn.ExtraCharges,
		TotalAmount:      checkIn.TotalAmount,
		CreatedAt:        checkIn.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        checkIn.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if checkIn.ActualCheckOut != nil {
		actualCheckOutStr := checkIn.ActualCheckOut.Format("2006-01-02 15:04:05")
		resp.ActualCheckOut = &actualCheckOutStr
	}

	if checkIn.Room != nil {
		roomTypeResp := &dto.RoomTypeResponse{
			ID:          checkIn.Room.RoomType.ID,
			Name:        checkIn.Room.RoomType.Name,
			Description: checkIn.Room.RoomType.Description,
			BasePrice:   checkIn.Room.RoomType.BasePrice,
			BedCount:    checkIn.Room.RoomType.BedCount,
			MaxGuests:   checkIn.Room.RoomType.MaxGuests,
			Facilities:  []string(checkIn.Room.RoomType.Facilities),
			CreatedAt:   checkIn.Room.RoomType.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   checkIn.Room.RoomType.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		resp.Room = &dto.RoomResponse{
			ID:         checkIn.Room.ID,
			RoomNo:     checkIn.Room.RoomNo,
			Floor:      checkIn.Room.Floor,
			RoomTypeID: checkIn.Room.RoomTypeID,
			RoomType:   roomTypeResp,
			Status:     checkIn.Room.Status,
			Price:      checkIn.Room.Price,
			Facilities: []string(checkIn.Room.Facilities),
			CreatedAt:  checkIn.Room.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  checkIn.Room.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return resp
}
