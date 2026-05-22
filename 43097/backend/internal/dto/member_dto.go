package dto

import "hotel-system/internal/model"

type MemberRegisterRequest struct {
	Name   string `json:"name" binding:"required,max=50"`
	Phone  string `json:"phone" binding:"required,max=20"`
	Email  string `json:"email" binding:"omitempty,email,max=100"`
	IDCard string `json:"idCard" binding:"omitempty,max=20"`
}

type MemberUpdateRequest struct {
	Name   string              `json:"name" binding:"omitempty,max=50"`
	Phone  string              `json:"phone" binding:"omitempty,max=20"`
	Email  string              `json:"email" binding:"omitempty,email,max=100"`
	IDCard string              `json:"idCard" binding:"omitempty,max=20"`
	Status model.MemberStatus `json:"status" binding:"omitempty,oneof=active inactive"`
}

type MemberDetailResponse struct {
	ID        uint                `json:"id"`
	MemberNo  string              `json:"memberNo"`
	Name      string              `json:"name"`
	Phone     string              `json:"phone"`
	Email     string              `json:"email"`
	IDCard    string              `json:"idCard"`
	LevelID   uint                `json:"levelId"`
	LevelName string              `json:"levelName"`
	Points    int                 `json:"points"`
	Balance   float64             `json:"balance"`
	Status    model.MemberStatus `json:"status"`
	CreatedAt string              `json:"createdAt"`
	UpdatedAt string              `json:"updatedAt"`
}

type MemberListRequest struct {
	PaginationRequest
	Name   string              `form:"name"`
	Phone  string              `form:"phone"`
	LevelID uint                `form:"levelId"`
	Status model.MemberStatus `form:"status"`
}

type MemberLevelCreateRequest struct {
	Name         string  `json:"name" binding:"required,max=50"`
	DiscountRate float64 `json:"discountRate" binding:"required,min=0,max=1"`
	PointsRate   float64 `json:"pointsRate" binding:"required,min=0"`
	MinPoints    int     `json:"minPoints" binding:"required,min=0"`
	MaxPoints    int     `json:"maxPoints" binding:"required,min=0"`
}

type MemberLevelUpdateRequest struct {
	Name         string  `json:"name" binding:"omitempty,max=50"`
	DiscountRate float64 `json:"discountRate" binding:"omitempty,min=0,max=1"`
	PointsRate   float64 `json:"pointsRate" binding:"omitempty,min=0"`
	MinPoints    int     `json:"minPoints" binding:"omitempty,min=0"`
	MaxPoints    int     `json:"maxPoints" binding:"omitempty,min=0"`
}

type MemberLevelDetailResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	DiscountRate float64 `json:"discountRate"`
	PointsRate   float64 `json:"pointsRate"`
	MinPoints    int     `json:"minPoints"`
	MaxPoints    int     `json:"maxPoints"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

type PointsUseRequest struct {
	MemberID    uint   `json:"memberId" binding:"required"`
	Points      int    `json:"points" binding:"required,min=1"`
	Description string `json:"description" binding:"omitempty,max=255"`
	OrderNo     string `json:"orderNo" binding:"omitempty,max=64"`
}

type PointsRechargeRequest struct {
	MemberID    uint   `json:"memberId" binding:"required"`
	Points      int    `json:"points" binding:"required,min=1"`
	Description string `json:"description" binding:"omitempty,max=255"`
	OrderNo     string `json:"orderNo" binding:"omitempty,max=64"`
}

type MemberDiscountResponse struct {
	MemberID     uint    `json:"memberId"`
	MemberNo     string  `json:"memberNo"`
	Name         string  `json:"name"`
	LevelID      uint    `json:"levelId"`
	LevelName    string  `json:"levelName"`
	DiscountRate float64 `json:"discountRate"`
	Points       int     `json:"points"`
	PointsRate   float64 `json:"pointsRate"`
}
