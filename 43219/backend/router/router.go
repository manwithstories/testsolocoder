package router

import (
	"housekeeping/handlers"
	"housekeeping/middleware"
	"housekeeping/models"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.New()
	r.Use(middleware.CORS(), middleware.Logger(), middleware.Recovery())

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		svc := api.Group("/services")
		{
			svc.GET("", handlers.ListServices)
			svc.GET("/:id", handlers.GetService)
		}

		staff := api.Group("/staff")
		{
			staff.GET("", handlers.ListStaff)
			staff.GET("/:id", handlers.GetStaffDetail)
			staff.GET("/:id/reviews", handlers.ListStaffReviews)
		}

		authed := api.Group("")
		authed.Use(middleware.Auth())
		{
			authed.GET("/me", handlers.Me)
			authed.PUT("/me", handlers.UpdateMe)

			bookings := authed.Group("/bookings")
			{
				bookings.POST("", handlers.CreateBooking)
				bookings.GET("", handlers.ListBookings)
				bookings.GET("/:id", handlers.GetBooking)
				bookings.POST("/:id/confirm", handlers.ConfirmBooking)
				bookings.POST("/:id/reschedule", handlers.RescheduleBooking)
				bookings.POST("/:id/cancel", handlers.CancelBooking)
			}

			orders := authed.Group("/orders")
			{
				orders.GET("", handlers.ListOrders)
				orders.GET("/:id", handlers.GetOrder)
				orders.POST("/:id/report", middleware.RequireRole(models.RoleStaff), handlers.SubmitReport)
				orders.POST("/:id/confirm", middleware.RequireRole(models.RoleCustomer), handlers.ConfirmOrder)
				orders.POST("/:id/refund", middleware.RequireRole(models.RoleCustomer), handlers.RequestRefund)
			}

			reviews := authed.Group("/reviews")
			{
				reviews.POST("", middleware.RequireRole(models.RoleCustomer), handlers.CreateReview)
				reviews.GET("", handlers.ListReviews)
				reviews.GET("/mine", handlers.MyReviews)
			}

			tickets := authed.Group("/tickets")
			{
				tickets.POST("", handlers.CreateTicket)
				tickets.GET("", handlers.ListTickets)
			}

			customer := authed.Group("/customer")
			customer.Use(middleware.RequireRole(models.RoleCustomer))
			{
				customer.GET("/tickets", handlers.ListTickets)
			}

			staffGrp := authed.Group("/staff-area")
			staffGrp.Use(middleware.RequireRole(models.RoleStaff))
			{
				staffGrp.PUT("/certs", handlers.UpdateStaffCert)
				staffGrp.GET("/earnings", handlers.MyEarnings)
				staffGrp.POST("/withdraw", handlers.RequestWithdrawal)
				staffGrp.GET("/wallet", handlers.WalletInfo)
			}

			company := authed.Group("/company")
			company.Use(middleware.RequireRole(models.RoleCompany))
			{
				company.POST("/services", handlers.CreateService)
				company.GET("/my-services", handlers.MyServices)
				company.PUT("/services/:id", handlers.UpdateService)
				company.DELETE("/services/:id", handlers.DeleteService)
				company.GET("/dashboard", handlers.StatsCompanyDashboard)
				company.GET("/finance/monthly", handlers.CompanyMonthlyStats)
				company.GET("/finance/export", handlers.ExportFinanceCSV)
				company.POST("/bookings/:id/review-reschedule", handlers.ReviewNeedReschedule)
			}

			admin := authed.Group("/admin")
			admin.Use(middleware.RequireRole(models.RoleAdmin))
			{
				admin.GET("/stats/overview", handlers.StatsOverview)
				admin.GET("/stats/revenue", handlers.StatsRevenueTrend)
				admin.GET("/stats/category", handlers.StatsOrdersByCategory)
				admin.GET("/stats/staff", handlers.StatsStaffPerf)
				admin.POST("/tickets/:id/assign", handlers.AssignTicket)
				admin.POST("/tickets/:id/resolve", handlers.ResolveTicket)
				admin.POST("/tickets/:id/close", handlers.CloseTicket)
				admin.PUT("/staff/:id/suspend", handlers.AdminSuspendStaff)
				admin.GET("/finance/export", handlers.ExportFinanceCSV)
			}
		}
	}
	return r
}
