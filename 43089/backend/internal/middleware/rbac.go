package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"travel-planner/internal/database"
	"travel-planner/internal/models"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetCurrentUser(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userRoles := make(map[string]bool)
		for _, role := range user.Roles {
			userRoles[role.Name] = true
		}

		for _, role := range roles {
			if userRoles[role] {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetCurrentUser(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		for _, role := range user.Roles {
			var count int64
			database.DB.Table("role_permissions rp").
				Joins("JOIN permissions p ON rp.permission_id = p.id").
				Where("rp.role_id = ? AND p.resource = ? AND p.action = ?", role.ID, resource, action).
				Count(&count)
			if count > 0 {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

func CheckPlanAccess(requiredEdit bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetCurrentUserID(c)
		planID := c.Param("plan_id")
		if planID == "" {
			planID = c.Param("id")
		}

		planUUID, err := uuid.Parse(planID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
			c.Abort()
			return
		}

		var plan models.TravelPlan
		if err := database.DB.First(&plan, "id = ?", planUUID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			c.Abort()
			return
		}

		if plan.OwnerID == userID {
			c.Set("plan", plan)
			c.Set("plan_access", "owner")
			c.Next()
			return
		}

		var participant models.PlanParticipant
		if err := database.DB.Where("plan_id = ? AND user_id = ?", planUUID, userID).First(&participant).Error; err != nil {
			if plan.IsPublic && !requiredEdit {
				c.Set("plan", plan)
				c.Set("plan_access", "viewer")
				c.Next()
				return
			}
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this plan"})
			c.Abort()
			return
		}

		if requiredEdit && !participant.CanEdit {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have edit permission for this plan"})
			c.Abort()
			return
		}

		c.Set("plan", plan)
		c.Set("plan_access", participant.Role)
		c.Next()
	}
}
