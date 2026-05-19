package dto

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"real_name" binding:"max=50"`
	Phone    string `json:"phone" binding:"max=20"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  interface{} `json:"user"`
}

type UpdateProfileRequest struct {
	RealName string `json:"real_name" binding:"max=50"`
	Phone    string `json:"phone" binding:"max=20"`
	Avatar   string `json:"avatar"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user admin super_admin"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type SendVerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
