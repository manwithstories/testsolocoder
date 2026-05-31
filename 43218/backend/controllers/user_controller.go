package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"secondhand-platform/config"
	"secondhand-platform/models"
	"secondhand-platform/services"
	"secondhand-platform/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
	Role     string `json:"role" binding:"required,oneof=seller buyer technician"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type RealNameAuthRequest struct {
	RealName string `json:"real_name" binding:"required"`
	IDCard   string `json:"id_card" binding:"required,len=18"`
}

type TechnicianCertRequest struct {
	CertType   string `json:"cert_type" binding:"required"`
	CertNumber string `json:"cert_number" binding:"required"`
	CertImage  string `json:"cert_image"`
}

type RechargeRequest struct {
	Amount        float64 `json:"amount" binding:"required,min=0.01"`
	PaymentMethod string  `json:"payment_method" binding:"required,oneof=alipay wechat wallet"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,min=1"`
}

func (ctrl *UserController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	user, err := ctrl.userService.Register(req.Username, req.Password, req.Email, req.Phone, req.Role, req.Nickname)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

func (ctrl *UserController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	user, accessToken, refreshToken, err := ctrl.userService.Login(req.Username, req.Password)
	if err != nil {
		utils.Error(c, 401, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"user": gin.H{
			"id":               user.ID,
			"username":         user.Username,
			"role":             user.Role,
			"nickname":         user.Nickname,
			"avatar":           user.Avatar,
			"credit_score":     user.CreditScore,
			"is_authenticated": user.IsAuthenticated,
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    config.AppConfig.JWT.AccessTokenExpire * 3600,
	})
}

func (ctrl *UserController) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if err := ctrl.userService.Logout(token); err != nil {
		utils.Error(c, 500, "登出失败")
		return
	}

	utils.Success(c, nil)
}

func (ctrl *UserController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	newToken, err := ctrl.userService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.Error(c, 401, "刷新令牌失败")
		return
	}

	utils.Success(c, gin.H{
		"access_token": newToken,
		"token_type":   "Bearer",
		"expires_in":   config.AppConfig.JWT.AccessTokenExpire * 3600,
	})
}

func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := ctrl.userService.GetUserByID(userID)
	if err != nil {
		utils.Error(c, 404, "用户不存在")
		return
	}

	utils.Success(c, gin.H{
		"id":               user.ID,
		"username":         user.Username,
		"email":            user.Email,
		"phone":            user.Phone,
		"nickname":         user.Nickname,
		"avatar":           user.Avatar,
		"role":             user.Role,
		"real_name":        user.RealName,
		"id_card":          user.IDCard,
		"is_authenticated": user.IsAuthenticated,
		"credit_score":     user.CreditScore,
		"wallet_balance":   user.WalletBalance,
		"created_at":       user.CreatedAt,
	})
}

func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Nickname string `json:"nickname" binding:"omitempty,max=50"`
		Email    string `json:"email" binding:"omitempty,email"`
		Phone    string `json:"phone" binding:"omitempty,len=11"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if err := ctrl.userService.UpdateProfile(userID, updates); err != nil {
		utils.Error(c, 500, "更新失败")
		return
	}

	utils.Success(c, nil)
}

func (ctrl *UserController) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *UserController) SubmitRealNameAuth(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req RealNameAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.userService.SubmitRealNameAuth(userID, req.RealName, req.IDCard); err != nil {
		utils.Error(c, 500, "提交失败")
		return
	}

	utils.Success(c, nil)
}

func (ctrl *UserController) SubmitTechnicianCert(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req TechnicianCertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	cert, err := ctrl.userService.SubmitTechnicianCert(userID, req.CertType, req.CertNumber, req.CertImage)
	if err != nil {
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, cert)
}

func (ctrl *UserController) Recharge(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	transaction, err := ctrl.userService.Recharge(userID, req.Amount, req.PaymentMethod)
	if err != nil {
		utils.Error(c, 500, "充值失败")
		return
	}

	utils.Success(c, transaction)
}

func (ctrl *UserController) Withdraw(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.userService.Withdraw(userID, req.Amount); err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (ctrl *UserController) GetWalletBalance(c *gin.Context) {
	userID := c.GetUint("user_id")

	balance, err := ctrl.userService.GetWalletBalance(userID)
	if err != nil {
		utils.Error(c, 500, "获取余额失败")
		return
	}

	utils.Success(c, gin.H{"balance": balance})
}

func (ctrl *UserController) ListWalletLogs(c *gin.Context) {
	userID := c.GetUint("user_id")
	page := getPage(c)
	pageSize := getPageSize(c)

	logs, total, err := ctrl.userService.ListWalletLogs(userID, page, pageSize)
	if err != nil {
		utils.Error(c, 500, "获取记录失败")
		return
	}

	utils.SuccessWithPagination(c, logs, page, pageSize, total)
}

func (ctrl *UserController) GetUserStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := ctrl.userService.GetUserStats(userID)
	if err != nil {
		utils.Error(c, 500, "获取统计失败")
		return
	}

	utils.Success(c, stats)
}

func (ctrl *UserController) ListUsers(c *gin.Context) {
	page := getPage(c)
	pageSize := getPageSize(c)
	role := c.Query("role")
	status := parseInt(c.Query("status"), 0)

	users, total, err := ctrl.userService.ListUsers(page, pageSize, role, status)
	if err != nil {
		utils.Error(c, 500, "获取用户列表失败")
		return
	}

	utils.SuccessWithPagination(c, users, page, pageSize, total)
}

func (ctrl *UserController) UpdateUserStatus(c *gin.Context) {
	userID := parseIntParam(c, "id")
	var req struct {
		Status int `json:"status" binding:"required,oneof=1 2 3"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	if err := ctrl.userService.UpdateUserStatus(uint(userID), req.Status); err != nil {
		utils.Error(c, 500, "更新失败")
		return
	}

	utils.Success(c, nil)
}

func (ctrl *UserController) ReviewRealNameAuth(c *gin.Context) {
	userID := parseIntParam(c, "id")
	var req struct {
		Approved     bool   `json:"approved"`
		RejectReason string `json:"reject_reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	if err := ctrl.userService.ReviewRealNameAuth(uint(userID), req.Approved, req.RejectReason); err != nil {
		utils.Error(c, 500, "审核失败")
		return
	}

	utils.Success(c, nil)
}

func (ctrl *UserController) ReviewTechnicianCert(c *gin.Context) {
	certID := parseIntParam(c, "id")
	adminID := c.GetUint("user_id")

	var req struct {
		Approved     bool   `json:"approved"`
		RejectReason string `json:"reject_reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "参数错误")
		return
	}

	if err := ctrl.userService.ReviewTechnicianCert(uint(certID), req.Approved, req.RejectReason, adminID); err != nil {
		utils.Error(c, 500, "审核失败")
		return
	}

	utils.Success(c, nil)
}

func UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.Error(c, 400, "上传文件错误")
		return
	}
	defer file.Close()

	if header.Size > config.AppConfig.Upload.MaxFileSize {
		utils.Error(c, 400, "文件大小超过限制")
		return
	}

	ext := filepath.Ext(header.Filename)
	ext = strings.ToLower(ext)

	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		utils.Error(c, 400, "不支持的文件类型")
		return
	}

	uploadPath := config.AppConfig.Upload.UploadPath
	os.MkdirAll(uploadPath, 0755)

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filepath := uploadPath + filename

	dst, err := os.Create(filepath)
	if err != nil {
		utils.Error(c, 500, "保存文件失败")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.Error(c, 500, "保存文件失败")
		return
	}

	utils.Success(c, gin.H{
		"url":      "/uploads/" + filename,
		"filename": filename,
	})
}

func UploadMultipleImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		utils.Error(c, 400, "上传文件错误")
		return
	}

	files := form.File["files"]
	if len(files) > 6 {
		utils.Error(c, 400, "最多上传6张图片")
		return
	}

	var urls []string
	uploadPath := config.AppConfig.Upload.UploadPath
	os.MkdirAll(uploadPath, 0755)

	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	for _, file := range files {
		if file.Size > config.AppConfig.Upload.MaxFileSize {
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExts[ext] {
			continue
		}

		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		filepath := uploadPath + filename

		if err := saveMultipartFile(file, filepath); err != nil {
			continue
		}

		urls = append(urls, "/uploads/"+filename)
	}

	utils.Success(c, gin.H{
		"urls": urls,
	})
}

func saveMultipartFile(file *multipart.FileHeader, filepath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func getPage(c *gin.Context) int {
	page := parseInt(c.Query("page"), 1)
	if page < 1 {
		page = 1
	}
	return page
}

func getPageSize(c *gin.Context) int {
	pageSize := parseInt(c.Query("page_size"), 10)
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageSize
}

func parseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

func parseIntParam(c *gin.Context, name string) int {
	return parseInt(c.Param(name), 0)
}

func parseFloat(s string, defaultVal float64) float64 {
	if s == "" {
		return defaultVal
	}
	var n float64
	fmt.Sscanf(s, "%f", &n)
	return n
}

func getUintFromContext(c *gin.Context, key string) uint {
	val, exists := c.Get(key)
	if !exists {
		return 0
	}
	return val.(uint)
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func init() {
	_ = models.UserStatusNormal
}
