package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"auction-system/internal/dto"
	"auction-system/internal/middleware"
	"auction-system/internal/services"
	"auction-system/pkg/response"
)

type AuctionItemController struct {
	itemService *services.AuctionItemService
}

func NewAuctionItemController() *AuctionItemController {
	return &AuctionItemController{
		itemService: services.NewAuctionItemService(),
	}
}

func (ctrl *AuctionItemController) Create(c *gin.Context) {
	sellerID := middleware.GetCurrentUserID(c)
	var req dto.CreateAuctionItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	item, err := ctrl.itemService.CreateItem(sellerID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, item)
}

func (ctrl *AuctionItemController) Update(c *gin.Context) {
	sellerID := middleware.GetCurrentUserID(c)
	itemID, _ := strconv.Atoi(c.Param("id"))

	var req dto.UpdateAuctionItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.itemService.UpdateItem(uint(itemID), sellerID, &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *AuctionItemController) GetByID(c *gin.Context) {
	itemID, _ := strconv.Atoi(c.Param("id"))

	item, err := ctrl.itemService.GetItemByID(uint(itemID))
	if err != nil {
		response.NotFound(c, "拍卖品不存在")
		return
	}

	response.Success(c, item)
}

func (ctrl *AuctionItemController) GetList(c *gin.Context) {
	var query dto.AuctionItemQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	items, total, err := ctrl.itemService.GetItemList(&query)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      items,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	})
}

func (ctrl *AuctionItemController) Online(c *gin.Context) {
	sellerID := middleware.GetCurrentUserID(c)
	itemID, _ := strconv.Atoi(c.Param("id"))

	if err := ctrl.itemService.OnlineItem(uint(itemID), sellerID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *AuctionItemController) Offline(c *gin.Context) {
	sellerID := middleware.GetCurrentUserID(c)
	itemID, _ := strconv.Atoi(c.Param("id"))

	if err := ctrl.itemService.OfflineItem(uint(itemID), sellerID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *AuctionItemController) Delete(c *gin.Context) {
	sellerID := middleware.GetCurrentUserID(c)
	itemID, _ := strconv.Atoi(c.Param("id"))

	if err := ctrl.itemService.DeleteItem(uint(itemID), sellerID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *AuctionItemController) UploadImages(c *gin.Context) {
	sellerID := middleware.GetCurrentUserID(c)
	itemID, _ := strconv.Atoi(c.Param("id"))

	images, err := ctrl.itemService.UploadImages(uint(itemID), sellerID, c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, images)
}

func (ctrl *AuctionItemController) GetMyItems(c *gin.Context) {
	sellerID := middleware.GetCurrentUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	var status *int
	if s, err := strconv.Atoi(c.Query("status")); err == nil {
		status = &s
	}

	items, total, err := ctrl.itemService.GetSellerItems(sellerID, page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
