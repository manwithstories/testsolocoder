package router

import (
	"property-management/internal/config"
	"property-management/internal/handler"
	"property-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Static("/uploads", cfg.UploadDir)

	authHandler := handler.NewAuthHandler(cfg.JWTSecret)
	propertyHandler := handler.NewPropertyHandler(cfg.UploadDir)
	tenantHandler := handler.NewTenantHandler()
	rentHandler := handler.NewRentHandler()
	repairHandler := handler.NewRepairHandler()
	feeHandler := handler.NewFeeHandler()
	noticeHandler := handler.NewNoticeHandler()
	statsHandler := handler.NewStatsHandler()

	api := r.Group("/api")
	{
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)

		auth := api.Group("", middleware.AuthMiddleware(cfg.JWTSecret))
		{
			auth.GET("/auth/profile", authHandler.GetProfile)
			auth.PUT("/auth/profile", authHandler.UpdateProfile)

			auth.GET("/properties", propertyHandler.List)
			auth.GET("/properties/:id", propertyHandler.Detail)
			auth.GET("/facilities", propertyHandler.ListFacilities)
			auth.POST("/properties/upload", propertyHandler.UploadImage)
			auth.GET("/properties/my/list", propertyHandler.MyProperties)

			auth.POST("/properties", middleware.RoleMiddleware("admin", "landlord"), propertyHandler.Create)
			auth.PUT("/properties/:id", middleware.RoleMiddleware("admin", "landlord"), propertyHandler.Update)
			auth.DELETE("/properties/:id", middleware.RoleMiddleware("admin"), propertyHandler.Delete)
			auth.PUT("/properties/:id/status", middleware.RoleMiddleware("admin", "landlord"), propertyHandler.UpdateStatus)

			auth.GET("/users", middleware.RoleMiddleware("admin"), authHandler.ListUsers)
			auth.PUT("/users/:id/status", middleware.RoleMiddleware("admin"), authHandler.UpdateUserStatus)
			auth.DELETE("/users/:id", middleware.RoleMiddleware("admin"), authHandler.DeleteUser)

			auth.GET("/tenants", tenantHandler.List)
			auth.GET("/tenants/:id", tenantHandler.Detail)
			auth.POST("/tenants", tenantHandler.Create)
			auth.PUT("/tenants/:id", tenantHandler.Update)
			auth.DELETE("/tenants/:id", middleware.RoleMiddleware("admin"), tenantHandler.Delete)

			auth.GET("/appointments", tenantHandler.ListAppointments)
			auth.POST("/appointments", tenantHandler.CreateAppointment)
			auth.PUT("/appointments/:id/status", tenantHandler.UpdateAppointmentStatus)

			auth.GET("/contracts", tenantHandler.ListContracts)
			auth.GET("/contracts/:id", tenantHandler.DetailContract)
			auth.POST("/contracts", middleware.RoleMiddleware("admin", "landlord"), tenantHandler.CreateContract)
			auth.PUT("/contracts/:id", middleware.RoleMiddleware("admin", "landlord"), tenantHandler.UpdateContract)
			auth.PUT("/contracts/:id/status", middleware.RoleMiddleware("admin", "landlord"), tenantHandler.UpdateContractStatus)
			auth.GET("/contracts/expiring", tenantHandler.GetExpiringContracts)

			auth.GET("/rent/bills", rentHandler.ListBills)
			auth.GET("/rent/bills/:id", rentHandler.GetBillDetail)
			auth.POST("/rent/generate", middleware.RoleMiddleware("admin", "landlord"), rentHandler.GenerateMonthlyBills)
			auth.POST("/rent/bills/:id/pay", rentHandler.PayBill)
			auth.POST("/rent/calculate-late-fee", middleware.RoleMiddleware("admin"), rentHandler.CalculateLateFee)

			auth.GET("/repairs", repairHandler.List)
			auth.GET("/repairs/my", repairHandler.MyOrders)
			auth.GET("/repairs/:id", repairHandler.Detail)
			auth.POST("/repairs", repairHandler.Create)
			auth.PUT("/repairs/:id/assign", middleware.RoleMiddleware("admin"), repairHandler.Assign)
			auth.PUT("/repairs/:id/status", repairHandler.UpdateStatus)

			auth.GET("/fees", feeHandler.List)
			auth.POST("/fees", middleware.RoleMiddleware("admin"), feeHandler.Create)
			auth.PUT("/fees/:id", middleware.RoleMiddleware("admin"), feeHandler.Update)
			auth.DELETE("/fees/:id", middleware.RoleMiddleware("admin"), feeHandler.Delete)
			auth.POST("/fees/:id/pay", feeHandler.Pay)
			auth.POST("/fees/batch", middleware.RoleMiddleware("admin"), feeHandler.BatchGenerate)

			auth.GET("/notices", noticeHandler.List)
			auth.GET("/notices/:id", noticeHandler.Detail)
			auth.POST("/notices", middleware.RoleMiddleware("admin", "landlord"), noticeHandler.Create)
			auth.PUT("/notices/:id", middleware.RoleMiddleware("admin"), noticeHandler.Update)
			auth.PATCH("/notices/:id", middleware.RoleMiddleware("admin"), noticeHandler.UpdateFields)
			auth.DELETE("/notices/:id", middleware.RoleMiddleware("admin"), noticeHandler.Delete)
			auth.PUT("/notices/:id/status", middleware.RoleMiddleware("admin"), noticeHandler.UpdateStatus)

			auth.GET("/stats/overview", statsHandler.GetOverview)
			auth.GET("/stats/occupancy-trend", statsHandler.GetOccupancyTrend)
			auth.GET("/stats/income-trend", statsHandler.GetIncomeTrend)
			auth.GET("/stats/repair-stats", statsHandler.GetRepairStats)
			auth.GET("/stats/export/rent", middleware.RoleMiddleware("admin"), statsHandler.ExportRentRecords)
			auth.GET("/stats/export/repairs", middleware.RoleMiddleware("admin"), statsHandler.ExportRepairOrders)
		}
	}

	return r
}
