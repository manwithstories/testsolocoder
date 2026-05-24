package router

import (
	"drone-rental/internal/handler"
	"drone-rental/internal/middleware"
	"drone-rental/internal/model"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	userHandler := handler.NewUserHandler()
	droneHandler := handler.NewDroneHandler()
	orderHandler := handler.NewOrderHandler()
	serviceHandler := handler.NewServiceHandler()
	flightHandler := handler.NewFlightHandler()
	insuranceHandler := handler.NewInsuranceHandler()
	reviewHandler := handler.NewReviewHandler()
	statsHandler := handler.NewStatsHandler()
	uploadHandler := handler.NewUploadHandler()
	excelHandler := handler.NewExcelHandler()

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		api.POST("/auth/register", userHandler.Register)
		api.POST("/auth/login", userHandler.Login)

		api.GET("/drones/search", droneHandler.SearchAvailable)
		api.GET("/drones/:id", droneHandler.GetByID)
		api.GET("/drones", droneHandler.List)
	}

	auth := api.Group("")
	auth.Use(middleware.JWTAuth())
	{
		auth.GET("/user/profile", userHandler.Profile)
		auth.PUT("/user/profile", userHandler.UpdateProfile)
		auth.POST("/user/verify-pilot", userHandler.VerifyPilot)
		auth.POST("/user/verify-owner", userHandler.VerifyOwner)

		auth.GET("/users", middleware.RoleAuth(model.RoleOwner, model.RolePilot), userHandler.List)
		auth.PUT("/users/audit", middleware.RoleAuth(model.RoleOwner), userHandler.AuditVerify)

		auth.POST("/drones", middleware.RoleAuth(model.RoleOwner), droneHandler.Create)
		auth.GET("/my-drones", middleware.RoleAuth(model.RoleOwner), droneHandler.MyDrones)
		auth.PUT("/drones/:id", middleware.RoleAuth(model.RoleOwner), droneHandler.Update)
		auth.PUT("/drones/:id/status", middleware.RoleAuth(model.RoleOwner), droneHandler.UpdateStatus)
		auth.DELETE("/drones/:id", middleware.RoleAuth(model.RoleOwner), droneHandler.Delete)
		auth.POST("/drones/batch-import", middleware.RoleAuth(model.RoleOwner), droneHandler.BatchImport)
		auth.POST("/drones/batch-import-excel", middleware.RoleAuth(model.RoleOwner), droneHandler.BatchImportExcel)

		auth.POST("/orders", middleware.RoleAuth(model.RoleClient), orderHandler.Create)
		auth.GET("/orders", orderHandler.List)
		auth.GET("/my-orders", orderHandler.MyOrders)
		auth.GET("/orders/:id", orderHandler.GetByID)
		auth.POST("/orders/pay", orderHandler.Pay)
		auth.POST("/orders/cancel", orderHandler.Cancel)
		auth.POST("/orders/:id/pickup", orderHandler.Pickup)
		auth.POST("/orders/confirm-return", orderHandler.ConfirmReturn)
		auth.POST("/orders/:id/complete", orderHandler.Complete)

		auth.POST("/services", middleware.RoleAuth(model.RoleClient), serviceHandler.Create)
		auth.GET("/services", serviceHandler.List)
		auth.GET("/my-services", serviceHandler.MyServices)
		auth.GET("/services/:id", serviceHandler.GetByID)
		auth.POST("/services/bid", middleware.RoleAuth(model.RolePilot), serviceHandler.CreateBid)
		auth.GET("/services/:id/bids", serviceHandler.ListBids)
		auth.POST("/services/accept-bid", middleware.RoleAuth(model.RoleClient), serviceHandler.AcceptBid)
		auth.PUT("/services/status", serviceHandler.UpdateStatus)

		auth.POST("/flights", middleware.RoleAuth(model.RolePilot), flightHandler.Create)
		auth.GET("/flights", flightHandler.List)
		auth.GET("/flights/:id", flightHandler.GetByID)
		auth.PUT("/flights/:id", middleware.RoleAuth(model.RolePilot), flightHandler.Update)
		auth.DELETE("/flights/:id", middleware.RoleAuth(model.RolePilot), flightHandler.Delete)

		auth.POST("/insurance/claims", insuranceHandler.CreateClaim)
		auth.GET("/insurance/claims", insuranceHandler.ListClaims)
		auth.GET("/insurance/claims/:id", insuranceHandler.GetClaimByID)
		auth.PUT("/insurance/claims/review", middleware.RoleAuth(model.RoleOwner), insuranceHandler.ReviewClaim)

		auth.POST("/reviews", reviewHandler.Create)
		auth.GET("/reviews", reviewHandler.List)
		auth.GET("/reviews/:id", reviewHandler.GetByID)
		auth.POST("/reviews/reply", reviewHandler.Reply)

		auth.GET("/stats/revenue", statsHandler.Revenue)
		auth.GET("/stats/region", statsHandler.Region)
		auth.GET("/stats/drone", statsHandler.Drone)

		auth.GET("/export/revenue", excelHandler.ExportRevenue)
		auth.GET("/export/region", excelHandler.ExportRegion)
		auth.GET("/export/drone", excelHandler.ExportDrone)

		auth.POST("/upload/image", uploadHandler.UploadImage)
		auth.POST("/upload/license", uploadHandler.UploadLicense)
		auth.POST("/upload/damage", uploadHandler.UploadDamage)
	}
}
