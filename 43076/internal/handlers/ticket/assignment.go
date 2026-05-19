package ticket

import (
	"errors"
	"strconv"
	"sync"

	"ticket-system/internal/database"
	"ticket-system/internal/middleware"
	"ticket-system/internal/models"
	"ticket-system/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssignmentRuleRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	Mode        string `json:"mode" binding:"required,assignment_mode"`
	SkillGroupID *uint  `json:"skill_group_id"`
	IsDefault   bool   `json:"is_default"`
	Enabled     bool   `json:"enabled"`
}

type AutoAssignRequest struct {
	TicketID     uint   `json:"ticket_id" binding:"required"`
	RuleID       *uint  `json:"rule_id"`
	Mode         string `json:"mode" binding:"omitempty,assignment_mode"`
	SkillGroupID *uint  `json:"skill_group_id"`
}

var (
	roundRobinCounter = make(map[uint]int)
	roundRobinMutex   sync.Mutex
)

func CreateAssignmentRule(c *gin.Context) {
	var req AssignmentRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	tx := database.DB.Begin()

	if req.IsDefault {
		if err := tx.Model(&models.AssignmentRule{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			tx.Rollback()
			utils.InternalServerError(c, "Failed to reset default rule")
			return
		}
	}

	rule := &models.AssignmentRule{
		Name:         req.Name,
		Description:  req.Description,
		Mode:         req.Mode,
		SkillGroupID: req.SkillGroupID,
		IsDefault:    req.IsDefault,
		Enabled:      req.Enabled,
	}

	if err := tx.Create(rule).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to create assignment rule")
		return
	}

	tx.Commit()
	utils.Success(c, rule)
}

func GetAssignmentRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid rule ID")
		return
	}

	var rule models.AssignmentRule
	if err := database.DB.First(&rule, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "Assignment rule not found")
			return
		}
		utils.InternalServerError(c, "Failed to get assignment rule")
		return
	}

	utils.Success(c, rule)
}

func ListAssignmentRules(c *gin.Context) {
	var rules []models.AssignmentRule
	query := database.DB.Model(&models.AssignmentRule{})

	if mode := c.Query("mode"); mode != "" {
		query = query.Where("mode = ?", mode)
	}
	if enabled := c.Query("enabled"); enabled != "" {
		query = query.Where("enabled = ?", enabled == "true")
	}

	if err := query.Order("is_default DESC, created_at DESC").Find(&rules).Error; err != nil {
		utils.InternalServerError(c, "Failed to list assignment rules")
		return
	}

	utils.Success(c, rules)
}

func UpdateAssignmentRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid rule ID")
		return
	}

	var req AssignmentRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	tx := database.DB.Begin()

	if req.IsDefault {
		if err := tx.Model(&models.AssignmentRule{}).Where("id != ? AND is_default = ?", uint(id), true).Update("is_default", false).Error; err != nil {
			tx.Rollback()
			utils.InternalServerError(c, "Failed to reset default rule")
			return
		}
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Mode != "" {
		updates["mode"] = req.Mode
	}
	if req.SkillGroupID != nil {
		updates["skill_group_id"] = *req.SkillGroupID
	}
	updates["is_default"] = req.IsDefault
	updates["enabled"] = req.Enabled

	if err := tx.Model(&models.AssignmentRule{}).Where("id = ?", uint(id)).Updates(updates).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to update assignment rule")
		return
	}

	tx.Commit()

	var rule models.AssignmentRule
	database.DB.First(&rule, uint(id))
	utils.Success(c, rule)
}

func DeleteAssignmentRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid rule ID")
		return
	}

	var rule models.AssignmentRule
	if err := database.DB.First(&rule, uint(id)).Error; err != nil {
		utils.NotFound(c, "Assignment rule not found")
		return
	}

	if rule.IsDefault {
		utils.BadRequest(c, "Cannot delete default assignment rule")
		return
	}

	if err := database.DB.Delete(&rule).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete assignment rule")
		return
	}

	utils.Success(c, nil)
}

func AutoAssign(c *gin.Context) {
	var req AutoAssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	operatorID, _ := middleware.GetCurrentUserID(c)

	tx := database.DB.Begin()

	var ticket models.Ticket
	if err := tx.First(&ticket, req.TicketID).Error; err != nil {
		tx.Rollback()
		utils.NotFound(c, "Ticket not found")
		return
	}

	if ticket.Status == models.TicketStatusClosed || ticket.Status == models.TicketStatusResolved {
		tx.Rollback()
		utils.BadRequest(c, "Cannot assign resolved or closed ticket")
		return
	}

	mode := req.Mode
	skillGroupID := req.SkillGroupID

	if mode == "" {
		var rule models.AssignmentRule
		if req.RuleID != nil {
			tx.First(&rule, *req.RuleID)
		} else {
			tx.Where("is_default = ? AND enabled = ?", true, true).First(&rule)
		}
		if rule.ID > 0 {
			mode = rule.Mode
			skillGroupID = rule.SkillGroupID
		} else {
			mode = models.AssignmentModeLoadBalance
		}
	}

	assigneeID, err := findAssignee(tx, mode, skillGroupID)
	if err != nil {
		tx.Rollback()
		utils.InternalServerError(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	updates["assignee_id"] = assigneeID
	updates["status"] = models.TicketStatusAssigned
	if skillGroupID != nil {
		updates["skill_group_id"] = *skillGroupID
	}

	if err := tx.Model(&ticket).Updates(updates).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to assign ticket")
		return
	}

	logEntry := &models.TicketLog{
		TicketID:   ticket.ID,
		OperatorID: operatorID,
		Action:     "auto_assign",
		NewValue:   strconv.FormatUint(uint64(assigneeID), 10),
		Remark:     "自动分配: " + mode,
	}
	if err := tx.Create(logEntry).Error; err != nil {
		tx.Rollback()
		utils.InternalServerError(c, "Failed to create ticket log")
		return
	}

	tx.Commit()

	database.DB.Preload("Assignee").First(&ticket, ticket.ID)
	utils.Success(c, ticket)
}

func findAssignee(tx *gorm.DB, mode string, skillGroupID *uint) (uint, error) {
	var users []models.User
	query := tx.Where("role IN ?", []string{models.RoleAdmin, models.RoleManager, models.RoleAgent})

	if skillGroupID != nil {
		query = query.Joins("JOIN user_skill_groups ON user_skill_groups.user_id = users.id").
			Where("user_skill_groups.skill_group_id = ?", *skillGroupID)
	}

	if err := query.Find(&users).Error; err != nil {
		return 0, err
	}

	if len(users) == 0 {
		return 0, errors.New("no available agents found")
	}

	switch mode {
	case models.AssignmentModeRoundRobin:
		return roundRobinAssign(users, skillGroupID)
	case models.AssignmentModeLoadBalance:
		return loadBalanceAssign(tx, users)
	case models.AssignmentModeSkillGroup:
		return loadBalanceAssign(tx, users)
	default:
		return loadBalanceAssign(tx, users)
	}
}

func roundRobinAssign(users []models.User, skillGroupID *uint) (uint, error) {
	roundRobinMutex.Lock()
	defer roundRobinMutex.Unlock()

	key := uint(0)
	if skillGroupID != nil {
		key = *skillGroupID
	}

	counter := roundRobinCounter[key]
	idx := counter % len(users)
	roundRobinCounter[key] = counter + 1

	return users[idx].ID, nil
}

func loadBalanceAssign(tx *gorm.DB, users []models.User) (uint, error) {
	type userLoad struct {
		UserID uint
		Count  int64
	}

	var loads []userLoad
	userIDs := make([]uint, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	tx.Model(&models.Ticket{}).
		Select("assignee_id as user_id, COUNT(*) as count").
		Where("assignee_id IN ? AND status NOT IN ?", userIDs, []string{models.TicketStatusClosed, models.TicketStatusResolved}).
		Group("assignee_id").
		Scan(&loads)

	loadMap := make(map[uint]int64)
	for _, l := range loads {
		loadMap[l.UserID] = l.Count
	}

	minLoad := int64(-1)
	var bestUserID uint
	for _, user := range users {
		load := loadMap[user.ID]
		if minLoad == -1 || load < minLoad {
			minLoad = load
			bestUserID = user.ID
		}
	}

	return bestUserID, nil
}

func GetAgentWorkload(c *gin.Context) {
	type Workload struct {
		UserID      uint   `json:"user_id"`
		Username    string `json:"username"`
		RealName    string `json:"real_name"`
		OpenCount   int64  `json:"open_count"`
		ActiveCount int64  `json:"active_count"`
		TotalCount  int64  `json:"total_count"`
	}

	var workloads []Workload

	rows, err := database.DB.Raw(`
		SELECT 
			u.id as user_id,
			u.username,
			u.real_name,
			COUNT(CASE WHEN t.status IN ('open', 'assigned') THEN 1 END) as open_count,
			COUNT(CASE WHEN t.status IN ('in_progress', 'pending', 'escalated') THEN 1 END) as active_count,
			COUNT(t.id) as total_count
		FROM users u
		LEFT JOIN tickets t ON t.assignee_id = u.id AND t.status NOT IN ('closed', 'resolved')
		WHERE u.role IN ('admin', 'manager', 'agent')
		GROUP BY u.id, u.username, u.real_name
		ORDER BY total_count DESC
	`).Rows()

	if err != nil {
		utils.InternalServerError(c, "Failed to get workload")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var w Workload
		rows.Scan(&w.UserID, &w.Username, &w.RealName, &w.OpenCount, &w.ActiveCount, &w.TotalCount)
		workloads = append(workloads, w)
	}

	utils.Success(c, workloads)
}
