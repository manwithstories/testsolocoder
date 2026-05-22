package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"auction-system/internal/dto"
	"auction-system/internal/middleware"
	"auction-system/internal/services"
	"auction-system/pkg/response"
)

type BidController struct {
	bidService     *services.BidService
	autoBidService *services.AutoBidService
}

func NewBidController() *BidController {
	return &BidController{
		bidService:     services.NewBidService(),
		autoBidService: services.NewAutoBidService(),
	}
}

func (ctrl *BidController) PlaceBid(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	itemID, _ := strconv.Atoi(c.Param("id"))

	var req dto.BidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	bid, err := ctrl.bidService.PlaceBid(userID, uint(itemID), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, bid)
}

func (ctrl *BidController) GetBidHistory(c *gin.Context) {
	itemID, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	bids, total, err := ctrl.bidService.GetBidHistory(uint(itemID), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      bids,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ctrl *BidController) GetMyBids(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	bids, total, err := ctrl.bidService.GetUserBids(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      bids,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ctrl *BidController) GetCurrentBid(c *gin.Context) {
	itemID, _ := strconv.Atoi(c.Param("id"))

	bid, err := ctrl.bidService.GetCurrentBid(uint(itemID))
	if err != nil {
		response.NotFound(c, "暂无出价记录")
		return
	}

	response.Success(c, bid)
}

func (ctrl *BidController) SetAutoBid(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var req dto.SetAutoBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	autoBid, err := ctrl.autoBidService.SetAutoBid(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, autoBid)
}

func (ctrl *BidController) CancelAutoBid(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	autoBidID, _ := strconv.Atoi(c.Param("id"))

	if err := ctrl.autoBidService.CancelAutoBid(userID, uint(autoBidID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *BidController) GetMyAutoBids(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	autoBids, err := ctrl.autoBidService.GetUserAutoBids(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, autoBids)
}

func (ctrl *BidController) GetAutoBid(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	itemID, _ := strconv.Atoi(c.Param("id"))

	autoBid, err := ctrl.autoBidService.GetAutoBid(uint(itemID), userID)
	if err != nil {
		response.NotFound(c, "未设置自动出价")
		return
	}

	response.Success(c, autoBid)
}
