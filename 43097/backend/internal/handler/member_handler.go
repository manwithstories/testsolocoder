package handler

import (
	"hotel-system/internal/dto"
	"hotel-system/internal/pkg/logger"
	"hotel-system/internal/service"
	"hotel-system/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberService service.MemberService
}

func NewMemberHandler(memberService service.MemberService) *MemberHandler {
	return &MemberHandler{memberService: memberService}
}

func (h *MemberHandler) RegisterMember(c *gin.Context) {
	var req dto.MemberRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("注册会员参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	member, err := h.memberService.RegisterMember(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	levelName := ""
	if member.Level != nil {
		levelName = member.Level.Name
	}

	memberResponse := dto.MemberDetailResponse{
		ID:        member.ID,
		MemberNo:  member.MemberNo,
		Name:      member.Name,
		Phone:     member.Phone,
		Email:     member.Email,
		IDCard:    member.IDCard,
		LevelID:   member.LevelID,
		LevelName: levelName,
		Points:    member.Points,
		Balance:   member.Balance,
		Status:    member.Status,
		CreatedAt: member.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: member.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.Success(c, memberResponse)
}

func (h *MemberHandler) GetMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的会员ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的会员ID")
		return
	}

	member, err := h.memberService.GetMember(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	levelName := ""
	if member.Level != nil {
		levelName = member.Level.Name
	}

	memberResponse := dto.MemberResponse{
		ID:        member.ID,
		MemberNo:  member.MemberNo,
		Name:      member.Name,
		Phone:     member.Phone,
		Email:     member.Email,
		IDCard:    member.IDCard,
		LevelID:   member.LevelID,
		LevelName: levelName,
		Points:    member.Points,
		Balance:   member.Balance,
		Status:    member.Status,
		CreatedAt: member.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: member.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.Success(c, memberResponse)
}

func (h *MemberHandler) GetMemberByPhone(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		logger.Warnf("手机号不能为空")
		utils.ErrorWithStatus(c, http.StatusBadRequest, "手机号不能为空")
		return
	}

	member, err := h.memberService.GetMemberByPhone(phone)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	levelName := ""
	if member.Level != nil {
		levelName = member.Level.Name
	}

	memberResponse := dto.MemberResponse{
		ID:        member.ID,
		MemberNo:  member.MemberNo,
		Name:      member.Name,
		Phone:     member.Phone,
		Email:     member.Email,
		IDCard:    member.IDCard,
		LevelID:   member.LevelID,
		LevelName: levelName,
		Points:    member.Points,
		Balance:   member.Balance,
		Status:    member.Status,
		CreatedAt: member.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: member.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.Success(c, memberResponse)
}

func (h *MemberHandler) UpdateMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的会员ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的会员ID")
		return
	}

	var req dto.MemberUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("更新会员参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	member, err := h.memberService.UpdateMember(uint(id), &req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	levelName := ""
	if member.Level != nil {
		levelName = member.Level.Name
	}

	memberResponse := dto.MemberResponse{
		ID:        member.ID,
		MemberNo:  member.MemberNo,
		Name:      member.Name,
		Phone:     member.Phone,
		Email:     member.Email,
		IDCard:    member.IDCard,
		LevelID:   member.LevelID,
		LevelName: levelName,
		Points:    member.Points,
		Balance:   member.Balance,
		Status:    member.Status,
		CreatedAt: member.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: member.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.Success(c, memberResponse)
}

func (h *MemberHandler) DeleteMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的会员ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的会员ID")
		return
	}

	err = h.memberService.DeleteMember(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *MemberHandler) ListMembers(c *gin.Context) {
	var req dto.MemberListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warnf("获取会员列表参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	members, total, err := h.memberService.ListMembers(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	var memberResponses []dto.MemberDetailResponse
	for _, member := range members {
		levelName := ""
		if member.Level != nil {
			levelName = member.Level.Name
		}
		memberResponses = append(memberResponses, dto.MemberDetailResponse{
			ID:        member.ID,
			MemberNo:  member.MemberNo,
			Name:      member.Name,
			Phone:     member.Phone,
			Email:     member.Email,
			IDCard:    member.IDCard,
			LevelID:   member.LevelID,
			LevelName: levelName,
			Points:    member.Points,
			Balance:   member.Balance,
			Status:    member.Status,
			CreatedAt: member.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: member.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	utils.PageResult(c, memberResponses, total, req.GetPage(), req.GetPageSize())
}

func (h *MemberHandler) CreateMemberLevel(c *gin.Context) {
	var req dto.MemberLevelCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("创建会员等级参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	level, err := h.memberService.CreateMemberLevel(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	levelResponse := dto.MemberLevelDetailResponse{
		ID:           level.ID,
		Name:         level.Name,
		DiscountRate: level.DiscountRate,
		PointsRate:   level.PointsRate,
		MinPoints:    level.MinPoints,
		MaxPoints:    level.MaxPoints,
		CreatedAt:    level.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    level.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.Success(c, levelResponse)
}

func (h *MemberHandler) GetMemberLevel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的会员等级ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的会员等级ID")
		return
	}

	level, err := h.memberService.GetMemberLevel(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	levelResponse := dto.MemberLevelResponse{
		ID:           level.ID,
		Name:         level.Name,
		DiscountRate: level.DiscountRate,
		PointsRate:   level.PointsRate,
		MinPoints:    level.MinPoints,
		MaxPoints:    level.MaxPoints,
		CreatedAt:    level.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    level.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.Success(c, levelResponse)
}

func (h *MemberHandler) UpdateMemberLevel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的会员等级ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的会员等级ID")
		return
	}

	var req dto.MemberLevelUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("更新会员等级参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	level, err := h.memberService.UpdateMemberLevel(uint(id), &req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	levelResponse := dto.MemberLevelResponse{
		ID:           level.ID,
		Name:         level.Name,
		DiscountRate: level.DiscountRate,
		PointsRate:   level.PointsRate,
		MinPoints:    level.MinPoints,
		MaxPoints:    level.MaxPoints,
		CreatedAt:    level.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    level.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	utils.Success(c, levelResponse)
}

func (h *MemberHandler) DeleteMemberLevel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的会员等级ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的会员等级ID")
		return
	}

	err = h.memberService.DeleteMemberLevel(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *MemberHandler) ListMemberLevels(c *gin.Context) {
	levels, total, err := h.memberService.ListMemberLevels()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	var levelResponses []dto.MemberLevelDetailResponse
	for _, level := range levels {
		levelResponses = append(levelResponses, dto.MemberLevelDetailResponse{
			ID:           level.ID,
			Name:         level.Name,
			DiscountRate: level.DiscountRate,
			PointsRate:   level.PointsRate,
			MinPoints:    level.MinPoints,
			MaxPoints:    level.MaxPoints,
			CreatedAt:    level.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    level.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	utils.Success(c, gin.H{
		"list":  levelResponses,
		"total": total,
	})
}

func (h *MemberHandler) UsePoints(c *gin.Context) {
	var req dto.PointsUseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("使用积分参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	cashValue, err := h.memberService.UsePoints(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"points":     req.Points,
		"cash_value": cashValue,
		"message":    "积分使用成功",
	})
}

func (h *MemberHandler) RechargePoints(c *gin.Context) {
	var req dto.PointsRechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("充值积分参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	err := h.memberService.RechargePoints(&req)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"points":  req.Points,
		"message": "积分充值成功",
	})
}

func (h *MemberHandler) GetMemberDiscount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的会员ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的会员ID")
		return
	}

	discount, err := h.memberService.GetMemberDiscount(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, discount)
}

func (h *MemberHandler) GetMemberConsumptionHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		logger.Warnf("无效的会员ID: %s", idStr)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "无效的会员ID")
		return
	}

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warnf("分页参数校验失败: %v", err)
		utils.ErrorWithStatus(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	records, total, err := h.memberService.GetMemberConsumptionHistory(uint(id), req.GetPage(), req.GetPageSize())
	if err != nil {
		logger.Errorf("获取会员消费历史失败: %v", err)
		utils.Error(c, "获取消费历史失败")
		return
	}

	utils.PageResult(c, records, total, req.GetPage(), req.GetPageSize())
}
