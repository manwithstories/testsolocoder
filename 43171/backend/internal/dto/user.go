package dto

import "drone-rental/internal/model"

type RegisterReq struct {
	Username string     `json:"username" binding:"required,min=3,max=64"`
	Password string     `json:"password" binding:"required,min=6,max=64"`
	Nickname string     `json:"nickname"`
	Phone    string     `json:"phone"`
	Email    string     `json:"email"`
	Role     model.Role `json:"role" binding:"required,oneof=client pilot owner"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token    string    `json:"token"`
	UserID   uint      `json:"user_id"`
	Username string    `json:"username"`
	Role     model.Role `json:"role"`
	Nickname string    `json:"nickname"`
	Avatar   string    `json:"avatar"`
	VerifyStatus model.VerifyStatus `json:"verify_status"`
}

type VerifyPilotReq struct {
	RealName string `json:"real_name" binding:"required"`
	IDCardNo string `json:"id_card_no" binding:"required"`
	LicenseNo string `json:"license_no" binding:"required"`
}

type VerifyOwnerReq struct {
	RealName string `json:"real_name" binding:"required"`
	IDCardNo string `json:"id_card_no" binding:"required"`
}

type UpdateUserReq struct {
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}
