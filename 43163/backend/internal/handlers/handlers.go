package handlers

import (
	"net/http"
	"strconv"
	"time"

	"printshop/internal/auth"
	"printshop/internal/config"
	"printshop/internal/middleware"
	"printshop/internal/models"
	"printshop/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	db       *gorm.DB
	cfg      *config.Config
	services *Services
}

type Services struct {
	Pricing    *service.PricingService
	Production *service.ProductionService
	Customer   *service.CustomerService
	Dashboard  *service.DashboardService
	File       *service.FileService
}

func NewHandler(db *gorm.DB, cfg *config.Config) *Handler {
	return &Handler{
		db:  db,
		cfg: cfg,
		services: &Services{
			Pricing:    service.NewPricingService(db),
			Production: service.NewProductionService(db),
			Customer:   service.NewCustomerService(db),
			Dashboard:  service.NewDashboardService(db),
			File:       service.NewFileService(db),
		},
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/auth/login", h.Login)
	api.POST("/auth/register", h.Register)

	auth := api.Group("", middleware.AuthMiddleware(h.cfg))
	auth.Use(middleware.AuditMiddleware(h.db))

	auth.GET("/auth/me", h.Me)

	auth.GET("/roles", h.ListRoles)
	auth.POST("/roles", h.CreateRole)

	auth.GET("/templates", h.ListTemplates)
	auth.GET("/templates/:id", h.GetTemplate)
	auth.POST("/templates", h.CreateTemplate)
	auth.PUT("/templates/:id", h.UpdateTemplate)
	auth.DELETE("/templates/:id", h.DeleteTemplate)

	auth.GET("/orders", h.ListOrders)
	auth.GET("/orders/:id", h.GetOrder)
	auth.POST("/orders", h.CreateOrder)
	auth.PUT("/orders/:id/status", h.UpdateOrderStatus)
	auth.POST("/orders/:id/split", h.SplitOrder)

	auth.POST("/pricing/calculate", h.CalculatePrice)
	auth.GET("/pricing/rules", h.ListPriceRules)
	auth.POST("/pricing/rules", h.CreatePriceRule)
	auth.PUT("/pricing/rules/:id", h.UpdatePriceRule)
	auth.DELETE("/pricing/rules/:id", h.DeletePriceRule)

	auth.GET("/production/lines", h.ListProductionLines)
	auth.POST("/production/lines", h.CreateProductionLine)
	auth.POST("/production/schedules", h.CreateSchedule)
	auth.PUT("/production/schedules/:id/progress", h.UpdateScheduleProgress)
	auth.GET("/production/schedules", h.ListSchedules)

	auth.GET("/customers", h.ListCustomers)
	auth.GET("/customers/:id", h.GetCustomer)
	auth.POST("/customers", h.CreateCustomer)
	auth.PUT("/customers/:id", h.UpdateCustomer)
	auth.DELETE("/customers/:id", h.DeleteCustomer)
	auth.POST("/invoices", h.GenerateInvoice)
	auth.GET("/invoices", h.ListInvoices)

	auth.GET("/dashboard/stats", h.GetDashboardStats)

	auth.POST("/files/upload", h.UploadFile)
	auth.GET("/files/:id", h.GetFile)

	auth.GET("/audit-logs", h.ListAuditLogs)
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := h.db.Preload("Role").Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !user.Active {
		c.JSON(http.StatusForbidden, gin.H{"error": "account is disabled"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := auth.GenerateToken(h.cfg, user.ID, user.Username, user.RoleID, user.Role.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=64"`
		Password string `json:"password" binding:"required,min=6"`
		RealName string `json:"real_name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		RoleID   uint   `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	user := models.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		RealName:     req.RealName,
		Email:        req.Email,
		Phone:        req.Phone,
		RoleID:       req.RoleID,
		Active:       true,
	}
	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "username": user.Username})
}

func (h *Handler) Me(c *gin.Context) {
	userID := c.GetUint(string(middleware.CtxUserID))
	var user models.User
	if err := h.db.Preload("Role").Preload("Customer").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) ListRoles(c *gin.Context) {
	var roles []models.Role
	h.db.Find(&roles)
	c.JSON(http.StatusOK, roles)
}

func (h *Handler) CreateRole(c *gin.Context) {
	var req models.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.db.Create(&req)
	c.JSON(http.StatusCreated, req)
}

func (h *Handler) ListTemplates(c *gin.Context) {
	var templates []models.Template
	query := h.db.Preload("Materials").Preload("Processes").Preload("Options")
	if cat := c.Query("category"); cat != "" {
		query = query.Where("category = ?", cat)
	}
	query.Order("id DESC").Find(&templates)
	c.JSON(http.StatusOK, templates)
}

func (h *Handler) GetTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var t models.Template
	if err := h.db.Preload("Materials").Preload("Processes").Preload("Options").First(&t, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *Handler) CreateTemplate(c *gin.Context) {
	var req struct {
		Name        string                  `json:"name" binding:"required"`
		Category    string                  `json:"category" binding:"required"`
		WidthMM     float64                 `json:"width_mm" binding:"required"`
		HeightMM    float64                 `json:"height_mm" binding:"required"`
		Description string                  `json:"description"`
		Thumbnail   string                  `json:"thumbnail"`
		Materials   []models.TemplateMaterial `json:"materials"`
		Processes   []models.TemplateProcess  `json:"processes"`
		Options     []models.TemplateOption   `json:"options"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t := models.Template{
		Name:        req.Name,
		Category:    req.Category,
		WidthMM:     req.WidthMM,
		HeightMM:    req.HeightMM,
		Description: req.Description,
		Thumbnail:   req.Thumbnail,
		Materials:   req.Materials,
		Processes:   req.Processes,
		Options:     req.Options,
		Active:      true,
	}
	userID := c.GetUint(string(middleware.CtxUserID))
	t.CreatedBy = userID
	if err := h.db.Create(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *Handler) UpdateTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var t models.Template
	if err := h.db.First(&t, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.db.Save(&t)
	c.JSON(http.StatusOK, t)
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.db.Delete(&models.Template{}, id)
	h.db.Where("template_id = ?", id).Delete(&models.TemplateMaterial{})
	h.db.Where("template_id = ?", id).Delete(&models.TemplateProcess{})
	h.db.Where("template_id = ?", id).Delete(&models.TemplateOption{})
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) ListOrders(c *gin.Context) {
	var orders []models.Order
	query := h.db.Preload("Customer").Preload("Items").Preload("Items.Template").Preload("Items.FileAsset")
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if cid := c.Query("customer_id"); cid != "" {
		id, _ := strconv.ParseUint(cid, 10, 64)
		query = query.Where("customer_id = ?", id)
	}
	if urgent := c.Query("urgent"); urgent == "true" {
		query = query.Where("urgent = ?", true)
	}
	query.Order("created_at DESC").Find(&orders)
	c.JSON(http.StatusOK, orders)
}

func (h *Handler) GetOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var o models.Order
	if err := h.db.Preload("Customer").Preload("Items").Preload("Items.Template").Preload("Items.FileAsset").Preload("Schedules").First(&o, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req struct {
		CustomerID uint   `json:"customer_id" binding:"required"`
		Urgent     bool   `json:"urgent"`
		Remark     string `json:"remark"`
		Items      []struct {
			TemplateID    uint   `json:"template_id" binding:"required"`
			MaterialID    uint   `json:"material_id" binding:"required"`
			ProcessIDs    []uint `json:"process_ids"`
			Quantity      int    `json:"quantity" binding:"required"`
			FileAssetID   *uint  `json:"file_asset_id"`
			Specification string `json:"specification"`
		} `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customer models.Customer
	if err := h.db.First(&customer, req.CustomerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer not found"})
		return
	}

	order := models.Order{
		OrderNo:    generateOrderNo(),
		CustomerID: req.CustomerID,
		Status:     "created",
		Urgent:     req.Urgent,
		Remark:     req.Remark,
	}
	userID := c.GetUint(string(middleware.CtxUserID))
	order.CreatedBy = userID

	var total float64
	for _, item := range req.Items {
		processStr := uintSliceToString(item.ProcessIDs)
		unitPrice, _, err := h.services.Pricing.Calculate(service.PricingInput{
			TemplateID:    item.TemplateID,
			MaterialID:    item.MaterialID,
			ProcessIDs:    item.ProcessIDs,
			Quantity:      item.Quantity,
			CustomerLevel: customer.Level,
			Urgent:        req.Urgent,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		unitPerItem := unitPrice / float64(item.Quantity)
		order.Items = append(order.Items, models.OrderItem{
			TemplateID:    item.TemplateID,
			MaterialID:    item.MaterialID,
			ProcessIDs:    processStr,
			Quantity:      item.Quantity,
			FileAssetID:   item.FileAssetID,
			UnitPrice:     unitPerItem,
			SubTotal:      unitPrice,
			Specification: item.Specification,
		})
		total += unitPrice
	}
	order.TotalPrice = total
	order.FinalPrice = total

	availableCredit := customer.CreditLimit - customer.Balance
	if total > availableCredit {
		c.JSON(http.StatusBadRequest, gin.H{"error": "超出客户信用额度，可用额度：¥" + strconv.FormatFloat(availableCredit, 'f', 2, 64)})
		return
	}

	if err := h.db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	validStatuses := map[string]bool{
		"created": true, "reviewing": true, "reviewed": true,
		"producing": true, "produced": true, "shipped": true, "cancelled": true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	updates := map[string]interface{}{"status": req.Status}
	userID := c.GetUint(string(middleware.CtxUserID))
	now := time.Now()
	switch req.Status {
	case "reviewed":
		updates["reviewed_by"] = userID
		updates["reviewed_at"] = &now
	case "produced":
		updates["produced_at"] = &now
	case "shipped":
		updates["shipped_at"] = &now
	case "cancelled":
		updates["cancelled_at"] = &now
	}
	h.db.Model(&models.Order{}).Where("id = ?", id).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (h *Handler) SplitOrder(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		ItemIDs []uint `json:"item_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := h.db.Preload("Items").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	newOrder := models.Order{
		OrderNo:    generateOrderNo(),
		CustomerID: order.CustomerID,
		Status:     "created",
		Urgent:     order.Urgent,
		Remark:     order.Remark + " (split)",
		CreatedBy:  order.CreatedBy,
	}

	var remainingItems []models.OrderItem
	var splitItems []models.OrderItem
	for _, item := range order.Items {
		found := false
		for _, iid := range req.ItemIDs {
			if item.ID == iid {
				found = true
				splitItems = append(splitItems, item)
				break
			}
		}
		if !found {
			remainingItems = append(remainingItems, item)
		}
	}

	var newTotal float64
	for _, item := range splitItems {
		newOrder.Items = append(newOrder.Items, models.OrderItem{
			TemplateID:    item.TemplateID,
			MaterialID:    item.MaterialID,
			ProcessIDs:    item.ProcessIDs,
			Quantity:      item.Quantity,
			FileAssetID:   item.FileAssetID,
			UnitPrice:     item.UnitPrice,
			SubTotal:      item.SubTotal,
			Specification: item.Specification,
		})
		newTotal += item.SubTotal
	}
	newOrder.TotalPrice = newTotal
	newOrder.FinalPrice = newTotal

	h.db.Create(&newOrder)

	h.db.Where("order_id = ? AND id IN ?", id, req.ItemIDs).Delete(&models.OrderItem{})

	var remainTotal float64
	for _, item := range remainingItems {
		remainTotal += item.SubTotal
	}
	h.db.Model(&order).Updates(map[string]interface{}{"total_price": remainTotal, "final_price": remainTotal})

	c.JSON(http.StatusCreated, newOrder)
}

func (h *Handler) CalculatePrice(c *gin.Context) {
	var req service.PricingInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	price, processNames, err := h.services.Pricing.Calculate(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_price": price, "process_names": processNames})
}

func (h *Handler) ListPriceRules(c *gin.Context) {
	rules, err := h.services.Pricing.GetRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rules)
}

func (h *Handler) CreatePriceRule(c *gin.Context) {
	var rule models.PriceRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.services.Pricing.CreateRule(&rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rule)
}

func (h *Handler) UpdatePriceRule(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var rule models.PriceRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.services.Pricing.UpdateRule(uint(id), &rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rule)
}

func (h *Handler) DeletePriceRule(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.services.Pricing.DeleteRule(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) ListProductionLines(c *gin.Context) {
	lines, err := h.services.Production.GetLines()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lines)
}

func (h *Handler) CreateProductionLine(c *gin.Context) {
	var line models.ProductionLine
	if err := c.ShouldBindJSON(&line); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.services.Production.CreateLine(&line); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, line)
}

func (h *Handler) CreateSchedule(c *gin.Context) {
	var req struct {
		OrderID   uint      `json:"order_id" binding:"required"`
		LineID    uint      `json:"line_id" binding:"required"`
		PlannedQty int      `json:"planned_qty" binding:"required"`
		StartDate time.Time `json:"start_date" binding:"required"`
		EndDate   time.Time `json:"end_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sched, err := h.services.Production.ScheduleOrder(req.OrderID, req.LineID, req.PlannedQty, req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sched)
}

func (h *Handler) UpdateScheduleProgress(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		ProducedQty int `json:"produced_qty" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.services.Production.UpdateProgress(uint(id), req.ProducedQty); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "progress updated"})
}

func (h *Handler) ListSchedules(c *gin.Context) {
	schedules, err := h.services.Production.GetSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

func (h *Handler) ListCustomers(c *gin.Context) {
	customers, err := h.services.Customer.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func (h *Handler) GetCustomer(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	cust, err := h.services.Customer.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	c.JSON(http.StatusOK, cust)
}

func (h *Handler) CreateCustomer(c *gin.Context) {
	var cust models.Customer
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.services.Customer.Create(&cust); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cust)
}

func (h *Handler) UpdateCustomer(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var cust models.Customer
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.services.Customer.Update(uint(id), &cust); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) DeleteCustomer(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.services.Customer.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) GenerateInvoice(c *gin.Context) {
	var req struct {
		CustomerID  uint      `json:"customer_id" binding:"required"`
		PeriodStart time.Time `json:"period_start" binding:"required"`
		PeriodEnd   time.Time `json:"period_end" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inv, err := h.services.Customer.GenerateInvoice(req.CustomerID, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, inv)
}

func (h *Handler) ListInvoices(c *gin.Context) {
	customerID := c.Query("customer_id")
	var cid uint
	if customerID != "" {
		id, _ := strconv.ParseUint(customerID, 10, 64)
		cid = uint(id)
	}
	invoices, err := h.services.Customer.GetInvoices(cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, invoices)
}

func (h *Handler) GetDashboardStats(c *gin.Context) {
	stats, err := h.services.Dashboard.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	if err := h.services.File.ValidateFile(header.Filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := header.Filename
	asset := models.FileAsset{
		FileName:   filename,
		FilePath:   "/uploads/" + filename,
		FileSize:   header.Size,
		FileType:   header.Header.Get("Content-Type"),
		Status:     "uploaded",
		PreviewURL: "/uploads/" + filename,
	}
	userID := c.GetUint(string(middleware.CtxUserID))
	asset.UploaderID = userID

	if err := h.services.File.SaveFile(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, asset)
}

func (h *Handler) GetFile(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	f, err := h.services.File.GetFile(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	c.JSON(http.StatusOK, f)
}

func (h *Handler) ListAuditLogs(c *gin.Context) {
	var logs []models.AuditLog
	query := h.db.Order("created_at DESC")
	if action := c.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}
	if limit := c.Query("limit"); limit != "" {
		l, _ := strconv.Atoi(limit)
		query = query.Limit(l)
	}
	query.Find(&logs)
	c.JSON(http.StatusOK, logs)
}

func generateOrderNo() string {
	return "PO" + time.Now().Format("20060102150405")
}

func uintSliceToString(ids []uint) string {
	if len(ids) == 0 {
		return ""
	}
	result := strconv.Itoa(int(ids[0]))
	for _, id := range ids[1:] {
		result += "," + strconv.Itoa(int(id))
	}
	return result
}
