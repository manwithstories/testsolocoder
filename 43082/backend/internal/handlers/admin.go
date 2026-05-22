package handlers

import (
	"encoding/json"
	"multishop/internal/database"
	"multishop/internal/dto"
	"multishop/internal/middleware"
	"multishop/internal/models"
	"multishop/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	role := c.Query("role")
	keyword := c.Query("keyword")

	query := database.DB.Model(&models.User{})
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var users []models.User
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users)

	result := make([]dto.UserInfo, 0, len(users))
	for _, u := range users {
		result = append(result, dto.UserInfo{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Phone:    u.Phone,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Role:     u.Role,
		})
	}

	utils.Paginated(c, result, total, page, pageSize)
}

func (h *AdminHandler) GetDisputes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	query := database.DB.Model(&models.Dispute{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var disputes []models.Dispute
	offset := (page - 1) * pageSize
	query.Preload("User").Preload("Shop").
		Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&disputes)

	result := make([]dto.DisputeInfo, 0, len(disputes))
	for _, d := range disputes {
		var evidence []string
		json.Unmarshal([]byte(d.Evidence), &evidence)

		info := dto.DisputeInfo{
			ID:        d.ID,
			OrderID:   d.OrderID,
			UserID:    d.UserID,
			Username:  d.User.Username,
			ShopID:    d.ShopID,
			ShopName:  d.Shop.Name,
			Type:      d.Type,
			Reason:    d.Reason,
			Evidence:  evidence,
			Status:    d.Status,
			Result:    d.Result,
			AdminID:   d.AdminID,
			CreatedAt: d.CreatedAt.Format(time.RFC3339),
		}
		if d.ResolvedAt != nil {
			info.ResolvedAt = d.ResolvedAt.Format(time.RFC3339)
		}
		result = append(result, info)
	}

	utils.Paginated(c, result, total, page, pageSize)
}

func (h *AdminHandler) ResolveDispute(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req dto.DisputeResolveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var dispute models.Dispute
	if err := database.DB.First(&dispute, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "纠纷不存在")
		return
	}

	now := time.Now()
	dispute.Status = models.DisputeStatusResolved
	dispute.Result = req.Result
	dispute.AdminID = &adminID
	dispute.ResolvedAt = &now
	database.DB.Save(&dispute)

	utils.Success(c, nil)
}

func (h *AdminHandler) GetStatistics(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	var totalUsers int64
	database.DB.Model(&models.User{}).Count(&totalUsers)

	var totalShops int64
	database.DB.Model(&models.Shop{}).Count(&totalShops)

	var totalProducts int64
	database.DB.Model(&models.Product{}).Count(&totalProducts)

	var totalOrders int64
	database.DB.Model(&models.Order{}).Where("status != ?", models.OrderStatusCancelled).
		Where("created_at BETWEEN ? AND ?", startDate, endDate+" 23:59:59").
		Count(&totalOrders)

	var totalSales float64
	database.DB.Model(&models.Order{}).Where("status IN ?", []string{models.OrderStatusCompleted, models.OrderStatusRefunded}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate+" 23:59:59").
		Select("COALESCE(SUM(total_amount), 0)").Scan(&totalSales)

	var pendingShops int64
	database.DB.Model(&models.Shop{}).Where("status = ?", models.ShopStatusPending).Count(&pendingShops)

	var openDisputes int64
	database.DB.Model(&models.Dispute{}).Where("status IN ?", []string{models.DisputeStatusOpen, models.DisputeStatusProcessing}).Count(&openDisputes)

	type DailySalesResult struct {
		Date   string  `gorm:"column:date"`
		Amount float64 `gorm:"column:amount"`
		Orders int64   `gorm:"column:orders"`
	}
	var dailySalesResults []DailySalesResult
	database.DB.Model(&models.Order{}).
		Select("DATE(created_at) as date, COALESCE(SUM(total_amount), 0) as amount, COUNT(*) as orders").
		Where("status IN ?", []string{models.OrderStatusCompleted, models.OrderStatusRefunded}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate+" 23:59:59").
		Group("DATE(created_at)").Order("date ASC").
		Scan(&dailySalesResults)

	dailySales := make([]dto.DailySales, 0, len(dailySalesResults))
	for _, ds := range dailySalesResults {
		dailySales = append(dailySales, dto.DailySales{
			Date:   ds.Date,
			Amount: ds.Amount,
			Orders: ds.Orders,
		})
	}

	utils.Success(c, dto.AdminStatistics{
		TotalUsers:    totalUsers,
		TotalShops:    totalShops,
		TotalProducts: totalProducts,
		TotalOrders:   totalOrders,
		TotalSales:    totalSales,
		PendingShops:  pendingShops,
		OpenDisputes:  openDisputes,
		DailySales:    dailySales,
	})
}

func (h *AdminHandler) ExportOrders(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var orders []models.Order
	query := database.DB.Preload("Items").Preload("User").Preload("Shop")
	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate+" 23:59:59")
	}
	query.Order("created_at DESC").Find(&orders)

	f := excelize.NewFile()
	sheetName := "订单数据"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"订单号", "下单时间", "用户", "店铺", "商品", "金额", "状态", "收货人", "电话", "地址"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, header)
	}

	for i, order := range orders {
		row := i + 2
		items := ""
		for j, item := range order.Items {
			if j > 0 {
				items += "; "
			}
			items += item.ProductName + " x" + strconv.Itoa(item.Quantity)
		}

		statusText := map[string]string{
			models.OrderStatusPendingPayment: "待付款",
			models.OrderStatusPendingShip:    "待发货",
			models.OrderStatusShipped:        "已发货",
			models.OrderStatusCompleted:      "已完成",
			models.OrderStatusRefunded:       "已退款",
			models.OrderStatusCancelled:      "已取消",
		}[order.Status]

		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), order.OrderNo)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), order.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), order.User.Username)
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row), order.Shop.Name)
		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), items)
		f.SetCellValue(sheetName, "F"+strconv.Itoa(row), order.TotalAmount)
		f.SetCellValue(sheetName, "G"+strconv.Itoa(row), statusText)
		f.SetCellValue(sheetName, "H"+strconv.Itoa(row), order.ReceiverName)
		f.SetCellValue(sheetName, "I"+strconv.Itoa(row), order.ReceiverPhone)
		f.SetCellValue(sheetName, "J"+strconv.Itoa(row), order.ReceiverAddress)
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=orders.xlsx")
	f.Write(c.Writer)
}

func (h *AdminHandler) CreateDispute(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req dto.DisputeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", req.OrderID, userID).First(&order).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "订单不存在")
		return
	}

	evidenceJSON, _ := json.Marshal(req.Evidence)
	dispute := models.Dispute{
		OrderID:  req.OrderID,
		UserID:   userID,
		ShopID:   order.ShopID,
		Type:     req.Type,
		Reason:   req.Reason,
		Evidence: string(evidenceJSON),
		Status:   models.DisputeStatusOpen,
	}

	if err := database.DB.Create(&dispute).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "提交失败")
		return
	}

	utils.Success(c, gin.H{"dispute_id": dispute.ID})
}
