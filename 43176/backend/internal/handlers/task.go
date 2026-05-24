package handlers

import (
	"context"
	"errand-service/internal/models"
	"errand-service/internal/utils"
	"errand-service/pkg/logger"
	redisClient "errand-service/pkg/redis"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewTaskHandler(db *gorm.DB, redisClient *redis.Client) *TaskHandler {
	return &TaskHandler{db: db, redis: redisClient}
}

type CreateTaskRequest struct {
	Type        string    `json:"type" binding:"required,oneof=buy pickup deliver queue errand"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	StartAddr   string    `json:"start_addr" binding:"required"`
	StartLat    float64   `json:"start_lat" binding:"required"`
	StartLng    float64   `json:"start_lng" binding:"required"`
	EndAddr     string    `json:"end_addr" binding:"required"`
	EndLat      float64   `json:"end_lat" binding:"required"`
	EndLng      float64   `json:"end_lng" binding:"required"`
	Deadline    time.Time `json:"deadline" binding:"required"`
	Reward      float64   `json:"reward" binding:"required,gt=0"`
	Images      []string  `json:"images"`
}

type ListTasksQuery struct {
	Status   string `form:"status"`
	Type     string `form:"type"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
}

func (h *TaskHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	if userRole != "publisher" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Only publishers can create tasks"})
		return
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request parameters", "error": err.Error()})
		return
	}

	if req.Deadline.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Deadline must be in the future"})
		return
	}

	var publisher models.User
	if err := h.db.First(&publisher, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
		return
	}

	if publisher.Balance < req.Reward {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Insufficient balance, please deposit first"})
		return
	}

	task := models.Task{
		PublisherID: userID,
		Type:        models.TaskType(req.Type),
		Title:       req.Title,
		Description: req.Description,
		StartAddr:   req.StartAddr,
		StartLat:    req.StartLat,
		StartLng:    req.StartLng,
		EndAddr:     req.EndAddr,
		EndLat:      req.EndLat,
		EndLng:      req.EndLng,
		PublishTime: time.Now(),
		Deadline:    req.Deadline,
		Reward:      req.Reward,
		Status:      models.TaskStatusPending,
	}

	tx := h.db.Begin()

	if err := tx.Create(&task).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to create task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create task"})
		return
	}

	for i, imgURL := range req.Images {
		taskImage := models.TaskImage{
			TaskID:    task.ID,
			ImageURL:  imgURL,
			SortOrder: i,
		}
		if err := tx.Create(&taskImage).Error; err != nil {
			tx.Rollback()
			logger.Errorf("Failed to create task image: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to save task images"})
			return
		}
	}

	if err := tx.Model(&publisher).Update("balance", gorm.Expr("balance - ?", req.Reward)).Error; err != nil {
		tx.Rollback()
		logger.Errorf("Failed to deduct balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to process payment"})
		return
	}

	tx.Commit()

	ctx := context.Background()
	geoKey := "tasks:nearby"
	redisClient.GeoAdd(ctx, geoKey, task.StartLng, task.StartLat, fmt.Sprintf("task:%d", task.ID))
	redisClient.Set(ctx, fmt.Sprintf("task:%d", task.ID), task.Status, 24*time.Hour)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Task created successfully",
		"data":    task,
	})
}

func (h *TaskHandler) List(c *gin.Context) {
	var query ListTasksQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid query parameters"})
		return
	}

	db := h.db.Model(&models.Task{}).Preload("Publisher").Preload("Images")

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}

	var total int64
	db.Count(&total)

	var tasks []models.Task
	offset := (query.Page - 1) * query.PageSize
	db.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&tasks)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total":     total,
			"page":      query.Page,
			"page_size": query.PageSize,
			"items":     tasks,
		},
	})
}

func (h *TaskHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := h.db.Preload("Publisher").Preload("Courier").Preload("Images").
		First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": task})
}

func (h *TaskHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := h.db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Task not found"})
		return
	}

	if task.PublisherID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "You can only update your own tasks"})
		return
	}

	if task.Status != models.TaskStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Task can only be updated when pending"})
		return
	}

	var req struct {
		Title       *string    `json:"title"`
		Description *string    `json:"description"`
		Deadline    *time.Time `json:"deadline"`
		Reward      *float64   `json:"reward"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid request parameters"})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Deadline != nil {
		if req.Deadline.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Deadline must be in the future"})
			return
		}
		updates["deadline"] = *req.Deadline
	}
	if req.Reward != nil {
		if *req.Reward <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Reward must be positive"})
			return
		}
		updates["reward"] = *req.Reward
	}

	if len(updates) > 0 {
		h.db.Model(&task).Updates(updates)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Task updated successfully"})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := h.db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Task not found"})
		return
	}

	if task.PublisherID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "You can only delete your own tasks"})
		return
	}

	if task.Status != models.TaskStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Only pending tasks can be deleted"})
		return
	}

	tx := h.db.Begin()

	if err := tx.Delete(&task).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to delete task"})
		return
	}

	if err := tx.Model(&models.User{}).Where("id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", task.Reward)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to refund"})
		return
	}

	tx.Commit()

	ctx := context.Background()
	redisClient.GeoRemove(ctx, "tasks:nearby", fmt.Sprintf("task:%d", id))
	redisClient.Delete(ctx, fmt.Sprintf("task:%d", id))

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Task deleted successfully, refund processed"})
}

func (h *TaskHandler) Accept(c *gin.Context) {
	userID := c.GetUint("user_id")
	userRole := c.GetString("user_role")

	if userRole != "courier" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Only couriers can accept tasks"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid task ID"})
		return
	}

	ctx := context.Background()
	lockKey := fmt.Sprintf("task:lock:%d", id)
	locked, err := redisClient.SetNX(ctx, lockKey, userID, 30*time.Second)
	if err != nil || !locked {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "Task is being processed by another courier"})
		return
	}
	defer redisClient.Delete(ctx, lockKey)

	var courierProfile models.CourierProfile
	if err := h.db.Where("user_id = ?", userID).First(&courierProfile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Courier profile not found"})
		return
	}

	if courierProfile.Status != "approved" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Your account has not been approved yet"})
		return
	}

	if courierProfile.CurrentTaskID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "You already have an ongoing task"})
		return
	}

	var task models.Task
	if err := h.db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Task not found"})
		return
	}

	if task.Status != models.TaskStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Task is no longer available"})
		return
	}

	tx := h.db.Begin()

	task.Status = models.TaskStatusAccepted
	task.CourierID = &userID
	if err := tx.Save(&task).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to accept task"})
		return
	}

	courierProfile.CurrentTaskID = &task.ID
	if err := tx.Save(&courierProfile).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update courier"})
		return
	}

	order := models.Order{
		TaskID:       task.ID,
		PublisherID:  task.PublisherID,
		CourierID:    userID,
		Status:       models.OrderStatusAccepted,
		Reward:       task.Reward,
		ServiceFee:   utils.CalculateServiceFee(task.Reward),
		ActualPayment: task.Reward - utils.CalculateServiceFee(task.Reward),
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create order"})
		return
	}

	tx.Commit()

	redisClient.Set(ctx, fmt.Sprintf("task:%d", id), task.Status, 24*time.Hour)
	redisClient.GeoRemove(ctx, "tasks:nearby", fmt.Sprintf("task:%d", id))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Task accepted successfully",
		"data":    gin.H{"task": task, "order": order},
	})
}

func (h *TaskHandler) Complete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := h.db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Task not found"})
		return
	}

	if task.CourierID == nil || *task.CourierID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Only the assigned courier can complete this task"})
		return
	}

	if task.Status != models.TaskStatusAccepted && task.Status != models.TaskStatusInProgress {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Task cannot be completed"})
		return
	}

	var req struct {
		ProofImages []string `json:"proof_images"`
	}
	c.ShouldBindJSON(&req)

	tx := h.db.Begin()

	oldStatus := task.Status
	task.Status = models.TaskStatusCompleted
	if err := tx.Save(&task).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to complete task"})
		return
	}

	var order models.Order
	if err := tx.Where("task_id = ?", id).First(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Order not found"})
		return
	}

	now := time.Now()
	order.Status = models.OrderStatusCompleted
	order.EndTime = &now
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update order"})
		return
	}

	for _, imgURL := range req.ProofImages {
		proof := models.OrderProof{
			OrderID:   order.ID,
			ImageURL:  imgURL,
			ProofType: "completion",
		}
		tx.Create(&proof)
	}

	statusLog := models.OrderStatusLog{
		OrderID:   order.ID,
		OldStatus: string(oldStatus),
		NewStatus: string(models.OrderStatusCompleted),
		Reason:    "Task completed by courier",
	}
	tx.Create(&statusLog)

	serviceFee := utils.CalculateServiceFee(task.Reward)
	courierEarning := task.Reward - serviceFee

	if err := tx.Model(&models.User{}).Where("id = ?", task.CourierID).
		Update("balance", gorm.Expr("balance + ?", courierEarning)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to settle payment"})
		return
	}

	settlementTx := models.Transaction{
		OrderID:        &order.ID,
		UserID:         *task.CourierID,
		Type:           models.TxTypeSettlement,
		Amount:         courierEarning,
		Status:         models.TxStatusCompleted,
		Description:    fmt.Sprintf("Settlement for task #%d", task.ID),
		TransactionNo:  utils.GenerateOrderNo(*task.CourierID),
		CompletedAt:    &now,
	}
	tx.Create(&settlementTx)

	serviceFeeTx := models.Transaction{
		OrderID:        &order.ID,
		UserID:         0,
		Type:           models.TxTypeServiceFee,
		Amount:         serviceFee,
		Status:         models.TxStatusCompleted,
		Description:    fmt.Sprintf("Service fee for task #%d", task.ID),
		TransactionNo:  utils.GenerateOrderNo(0),
		CompletedAt:    &now,
	}
	tx.Create(&serviceFeeTx)

	var courierProfile models.CourierProfile
	tx.Where("user_id = ?", userID).First(&courierProfile)
	courierProfile.CurrentTaskID = nil
	courierProfile.CompletedOrders++
	courierProfile.TotalOrders++
	tx.Save(&courierProfile)

	var courier models.User
	tx.First(&courier, userID)
	courier.OrderCount++
	tx.Save(&courier)

	tx.Commit()

	ctx := context.Background()
	redisClient.Delete(ctx, fmt.Sprintf("task:%d", id))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Task completed successfully",
		"data":    gin.H{"task": task, "order": order},
	})
}

func (h *TaskHandler) Cancel(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := h.db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Task not found"})
		return
	}

	if task.PublisherID != userID {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Only the publisher can cancel this task"})
		return
	}

	if task.Status == models.TaskStatusCompleted || task.Status == models.TaskStatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Task cannot be cancelled"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	tx := h.db.Begin()

	oldStatus := task.Status
	task.Status = models.TaskStatusCancelled
	tx.Save(&task)

	if task.CourierID != nil {
		var courierProfile models.CourierProfile
		tx.Where("user_id = ?", *task.CourierID).First(&courierProfile)
		courierProfile.CurrentTaskID = nil
		courierProfile.CancelledOrders++
		tx.Save(&courierProfile)
	}

	tx.Model(&models.User{}).Where("id = ?", task.PublisherID).
		Update("balance", gorm.Expr("balance + ?", task.Reward))

	var order models.Order
	if err := tx.Where("task_id = ?", id).First(&order).Error; err == nil {
		order.Status = models.OrderStatusCancelled
		tx.Save(&order)

		statusLog := models.OrderStatusLog{
			OrderID:   order.ID,
			OldStatus: string(oldStatus),
			NewStatus: string(models.OrderStatusCancelled),
			Reason:    req.Reason,
		}
		tx.Create(&statusLog)
	}

	tx.Commit()

	ctx := context.Background()
	redisClient.Delete(ctx, fmt.Sprintf("task:%d", id))
	redisClient.GeoRemove(ctx, "tasks:nearby", fmt.Sprintf("task:%d", id))

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Task cancelled successfully"})
}

func (h *TaskHandler) GetNearby(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	radius, _ := strconv.ParseFloat(c.DefaultQuery("radius", "5"), 64)

	if lat == 0 || lng == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Latitude and longitude are required"})
		return
	}

	ctx := context.Background()
	locations, err := redisClient.GeoRadius(ctx, "tasks:nearby", lng, lat, radius, "km")
	if err != nil {
		logger.Errorf("GeoRadius error: %v", err)
		var tasks []models.Task
		h.db.Where("status = ?", models.TaskStatusPending).
			Preload("Publisher").Preload("Images").
			Find(&tasks)
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": tasks})
		return
	}

	var taskIDs []uint
	for _, loc := range locations {
		var id uint
		fmt.Sscanf(loc.Name, "task:%d", &id)
		if id > 0 {
			taskIDs = append(taskIDs, id)
		}
	}

	if len(taskIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": []models.Task{}})
		return
	}

	var tasks []models.Task
	h.db.Where("id IN ? AND status = ?", taskIDs, models.TaskStatusPending).
		Preload("Publisher").Preload("Images").
		Find(&tasks)

	type TaskWithDistance struct {
		models.Task
		Distance float64 `json:"distance"`
	}

	var result []TaskWithDistance
	for _, task := range tasks {
		distance := utils.CalculateDistance(lat, lng, task.StartLat, task.StartLng)
		result = append(result, TaskWithDistance{Task: task, Distance: distance})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result})
}
