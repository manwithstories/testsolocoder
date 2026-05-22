package dto

import "hotel-system/internal/model"

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	RealName string `json:"realName" binding:"max=50"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
}

type UserResponse struct {
	ID        uint           `json:"id"`
	Username  string         `json:"username"`
	RealName  string         `json:"realName"`
	Phone     string         `json:"phone"`
	Email     string         `json:"email"`
	Role      model.UserRole `json:"role"`
	Status    model.UserStatus `json:"status"`
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
}

type UserUpdateRequest struct {
	RealName string           `json:"realName" binding:"max=50"`
	Phone    string           `json:"phone" binding:"omitempty,max=20"`
	Email    string           `json:"email" binding:"omitempty,email,max=100"`
	Role     model.UserRole   `json:"role" binding:"omitempty,oneof=admin frontdesk user"`
	Status   model.UserStatus `json:"status" binding:"omitempty,oneof=active inactive"`
	Password string           `json:"password" binding:"omitempty,min=6,max=50"`
}

type UserListRequest struct {
	PaginationRequest
	Username string         `form:"username"`
	Role     model.UserRole `form:"role"`
	Status   model.UserStatus `form:"status"`
}
