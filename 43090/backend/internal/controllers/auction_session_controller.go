package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"auction-system/internal/dto"
	"auction-system/internal/services"
	"auction-system/pkg/response"
)

type AuctionSessionController struct {
	sessionService *services.AuctionSessionService
}

func NewAuctionSessionController() *AuctionSessionController {
	return &AuctionSessionController{
		sessionService: services.NewAuctionSessionService(),
	}
}

func (ctrl *AuctionSessionController) Create(c *gin.Context) {
	var req dto.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	session, err := ctrl.sessionService.CreateSession(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, session)
}

func (ctrl *AuctionSessionController) Update(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))

	var req dto.UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.sessionService.UpdateSession(uint(sessionID), &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *AuctionSessionController) GetByID(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))

	session, err := ctrl.sessionService.GetSessionByID(uint(sessionID))
	if err != nil {
		response.NotFound(c, "拍卖会不存在")
		return
	}

	response.Success(c, session)
}

func (ctrl *AuctionSessionController) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.DefaultQuery("keyword", "")
	var status *int
	if s, err := strconv.Atoi(c.Query("status")); err == nil {
		status = &s
	}

	sessions, total, err := ctrl.sessionService.GetSessionList(page, pageSize, status, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      sessions,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ctrl *AuctionSessionController) AddItems(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))

	var req dto.AddItemToSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.sessionService.AddItemsToSession(uint(sessionID), req.AuctionItemIDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *AuctionSessionController) RemoveItem(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))
	itemID, _ := strconv.Atoi(c.Param("item_id"))

	if err := ctrl.sessionService.RemoveItemFromSession(uint(sessionID), uint(itemID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *AuctionSessionController) Start(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))

	if err := ctrl.sessionService.StartSession(uint(sessionID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *AuctionSessionController) End(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))

	if err := ctrl.sessionService.EndSession(uint(sessionID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *AuctionSessionController) Cancel(c *gin.Context) {
	sessionID, _ := strconv.Atoi(c.Param("id"))

	if err := ctrl.sessionService.CancelSession(uint(sessionID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *AuctionSessionController) GetActive(c *gin.Context) {
	sessions, err := ctrl.sessionService.GetActiveSessions()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, sessions)
}
