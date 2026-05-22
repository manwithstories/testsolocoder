package controller

import (
	"strconv"
	"ticket-system/internal/common/response"
	"ticket-system/internal/dto"
	"ticket-system/internal/service"

	"github.com/gin-gonic/gin"
)

type CouponController struct {
	couponService *service.CouponService
}

func NewCouponController() *CouponController {
	return &CouponController{
		couponService: service.NewCouponService(),
	}
}

func (c *CouponController) Create(ctx *gin.Context) {
	var req dto.CouponCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	coupon, err := c.couponService.Create(&req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, coupon)
}

func (c *CouponController) GetList(ctx *gin.Context) {
	var req dto.CouponListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	coupons, total, err := c.couponService.GetList(&req)
	if err != nil {
		response.ServerError(ctx, "获取优惠券列表失败")
		return
	}

	response.Page(ctx, coupons, total, req.Page, req.PageSize)
}

func (c *CouponController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的优惠券ID")
		return
	}

	coupon, err := c.couponService.GetByID(uint(id))
	if err != nil {
		response.NotFound(ctx, "优惠券不存在")
		return
	}

	response.Success(ctx, coupon)
}

func (c *CouponController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的优惠券ID")
		return
	}

	if err := c.couponService.Delete(uint(id)); err != nil {
		response.ServerError(ctx, "删除优惠券失败")
		return
	}

	response.Success(ctx, nil)
}
