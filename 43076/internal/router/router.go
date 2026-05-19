package router

import (
	"ticket-system/internal/handlers/attachment"
	"ticket-system/internal/handlers/comment"
	"ticket-system/internal/handlers/customer"
	"ticket-system/internal/handlers/report"
	"ticket-system/internal/handlers/sla"
	"ticket-system/internal/handlers/team"
	"ticket-system/internal/handlers/ticket"
	"ticket-system/internal/middleware"
	"ticket-system/internal/models"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	api.POST("/auth/login", team.Login)

	auth := api.Group("")
	auth.Use(middleware.Auth())
	{
		auth.GET("/auth/me", team.GetCurrentUser)

		admin := auth.Group("")
		admin.Use(middleware.RoleRequired(models.RoleAdmin, models.RoleManager))
		{
			admin.POST("/departments", team.CreateDepartment)
			admin.GET("/departments", team.ListDepartments)
			admin.GET("/departments/:id", team.GetDepartment)
			admin.PUT("/departments/:id", team.UpdateDepartment)
			admin.DELETE("/departments/:id", team.DeleteDepartment)

			admin.POST("/skill-groups", team.CreateSkillGroup)
			admin.GET("/skill-groups", team.ListSkillGroups)
			admin.GET("/skill-groups/:id", team.GetSkillGroup)
			admin.PUT("/skill-groups/:id", team.UpdateSkillGroup)
			admin.DELETE("/skill-groups/:id", team.DeleteSkillGroup)

			admin.POST("/users", team.CreateUser)
			admin.PUT("/users/:id", team.UpdateUser)
			admin.DELETE("/users/:id", team.DeleteUser)

			admin.POST("/assignment-rules", ticket.CreateAssignmentRule)
			admin.GET("/assignment-rules", ticket.ListAssignmentRules)
			admin.GET("/assignment-rules/:id", ticket.GetAssignmentRule)
			admin.PUT("/assignment-rules/:id", ticket.UpdateAssignmentRule)
			admin.DELETE("/assignment-rules/:id", ticket.DeleteAssignmentRule)
		}

		users := auth.Group("/users")
		users.Use(middleware.RoleRequired(models.RoleAdmin, models.RoleManager, models.RoleAgent))
		{
			users.GET("", team.ListUsers)
			users.GET("/:id", team.GetUser)
		}

		customers := auth.Group("/customers")
		customers.Use(middleware.RoleRequired(models.RoleAdmin, models.RoleManager, models.RoleAgent))
		{
			customers.POST("", customer.CreateCustomer)
			customers.GET("", customer.ListCustomers)
			customers.GET("/:id", customer.GetCustomer)
			customers.GET("/:id/tickets", customer.GetCustomerWithTickets)
			customers.PUT("/:id", customer.UpdateCustomer)
			customers.DELETE("/:id", customer.DeleteCustomer)
		}

		tickets := auth.Group("/tickets")
		tickets.Use(middleware.RoleRequired(models.RoleAdmin, models.RoleManager, models.RoleAgent))
		{
			tickets.POST("", ticket.CreateTicket)
			tickets.GET("", ticket.ListTickets)
			tickets.GET("/:id", ticket.GetTicket)
			tickets.GET("/no/:ticket_no", ticket.GetTicketByNo)
			tickets.PUT("/:id", ticket.UpdateTicket)
			tickets.PUT("/:id/status", ticket.UpdateTicketStatus)
			tickets.POST("/:id/assign", ticket.AssignTicket)
			tickets.GET("/:id/logs", ticket.GetTicketLogs)
			tickets.POST("/auto-assign", ticket.AutoAssign)
			tickets.GET("/agents/workload", ticket.GetAgentWorkload)
		}

		comments := auth.Group("/comments")
		comments.Use(middleware.RoleRequired(models.RoleAdmin, models.RoleManager, models.RoleAgent))
		{
			comments.POST("", comment.CreateComment)
			comments.GET("", comment.ListComments)
			comments.GET("/:id", comment.GetComment)
			comments.PUT("/:id", comment.UpdateComment)
			comments.DELETE("/:id", comment.DeleteComment)
		}

		attachments := auth.Group("/attachments")
		attachments.Use(middleware.RoleRequired(models.RoleAdmin, models.RoleManager, models.RoleAgent))
		{
			attachments.POST("", attachment.UploadAttachment)
			attachments.GET("", attachment.ListAttachments)
			attachments.GET("/:id", attachment.GetAttachment)
			attachments.DELETE("/:id", attachment.DeleteAttachment)
		}

		slaGroup := auth.Group("/sla")
		slaGroup.Use(middleware.RoleRequired(models.RoleAdmin, models.RoleManager, models.RoleAgent))
		{
			slaGroup.GET("/ticket/:ticket_id", sla.GetSLARecord)
			slaGroup.GET("/breached", sla.ListBreachedSLA)
			slaGroup.POST("/ticket/:ticket_id/escalate", sla.EscalateTicket)
			slaGroup.GET("/stats", sla.GetSLAStats)
		}

		reports := auth.Group("/reports")
		reports.Use(middleware.RoleRequired(models.RoleAdmin, models.RoleManager))
		{
			reports.GET("/overview", report.GetOverallStats)
			reports.GET("/daily", report.GetDailyStats)
			reports.GET("/weekly", report.GetWeeklyStats)
			reports.GET("/monthly", report.GetMonthlyStats)
			reports.GET("/priority-distribution", report.GetPriorityDistribution)
			reports.GET("/status-distribution", report.GetStatusDistribution)
			reports.GET("/type-distribution", report.GetTypeDistribution)
			reports.GET("/agent-performance", report.GetAgentPerformance)
		}
	}
}
