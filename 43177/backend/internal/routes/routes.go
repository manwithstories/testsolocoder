package routes

import (
	"github.com/gin-gonic/gin"
	"repair-platform/internal/handlers"
	"repair-platform/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	userHandler := handlers.NewUserHandler()
	categoryHandler := handlers.NewCategoryHandler()
	serviceItemHandler := handlers.NewServiceItemHandler()
	orderHandler := handlers.NewOrderHandler()
	partHandler := handlers.NewPartHandler()
	financeHandler := handlers.NewFinanceHandler()
	adminHandler := handlers.NewAdminHandler()

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		api.POST("/auth/register", userHandler.Register)
		api.POST("/auth/login", userHandler.Login)

		api.GET("/categories", categoryHandler.GetCategories)
		api.GET("/categories/:id", categoryHandler.GetCategoryDetail)
		api.GET("/service-items", serviceItemHandler.GetServiceItems)
		api.GET("/service-items/:id", serviceItemHandler.GetServiceItemDetail)

		api.GET("/technicians", userHandler.GetTechnicianList)
		api.GET("/technicians/:id", userHandler.GetTechnicianDetail)
		api.GET("/reviews", userHandler.GetReviewList)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/user/profile", userHandler.GetUserInfo)
			auth.PUT("/user/profile", userHandler.UpdateProfile)
			auth.POST("/user/certificate", userHandler.SubmitCertificate)
			auth.POST("/reviews", userHandler.ReviewTechnician)
			auth.POST("/reviews/:id/reply", userHandler.ReplyReview)

			auth.POST("/orders", orderHandler.CreateOrder)
			auth.GET("/orders", orderHandler.GetOrderList)
			auth.GET("/orders/:id", orderHandler.GetOrderDetail)
			auth.POST("/orders/:id/accept", orderHandler.AcceptOrder)
			auth.POST("/orders/:id/arrive", orderHandler.ArriveAtSite)
			auth.POST("/orders/:id/start", orderHandler.StartRepair)
			auth.POST("/orders/:id/complete", orderHandler.CompleteOrder)
			auth.POST("/orders/:id/cancel", orderHandler.CancelOrder)
			auth.POST("/orders/:id/refund", orderHandler.RequestRefund)

			auth.POST("/parts/use", partHandler.UsePart)
			auth.POST("/part-requests", partHandler.CreatePartRequest)
			auth.GET("/part-requests", partHandler.GetPartRequestList)
			auth.GET("/part-requests/:id", partHandler.GetPartRequestDetail)
			auth.POST("/part-requests/:id/receive", partHandler.ReceivePartRequest)

			auth.GET("/finance/balance", financeHandler.GetBalance)
			auth.POST("/finance/withdraw", financeHandler.CreateWithdrawRequest)
			auth.GET("/finance/withdraws", financeHandler.GetWithdrawRequestList)
			auth.GET("/finance/withdraws/:id", financeHandler.GetWithdrawRequestDetail)
			auth.GET("/finance/transactions", financeHandler.GetTransactionList)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
		{
			admin.GET("/dashboard", adminHandler.GetDashboardStats)

			admin.POST("/categories", categoryHandler.CreateCategory)
			admin.PUT("/categories/:id", categoryHandler.UpdateCategory)
			admin.DELETE("/categories/:id", categoryHandler.DeleteCategory)

			admin.POST("/service-items", serviceItemHandler.CreateServiceItem)
			admin.PUT("/service-items/:id", serviceItemHandler.UpdateServiceItem)
			admin.DELETE("/service-items/:id", serviceItemHandler.DeleteServiceItem)

			admin.GET("/users", adminHandler.GetUserList)
			admin.GET("/users/:id", adminHandler.GetUserDetail)
			admin.PUT("/users/:id/status", adminHandler.UpdateUserStatus)

			admin.GET("/technicians/verify", adminHandler.GetTechnicianVerifyList)
			admin.POST("/technicians/:id/verify", adminHandler.VerifyTechnician)

			admin.GET("/orders", adminHandler.GetAllOrders)
			admin.GET("/orders/status/:status", orderHandler.GetOrdersByStatus)
			admin.GET("/refunds", adminHandler.GetRefundList)
			admin.POST("/refunds/:id/approve", orderHandler.ApproveRefund)
			admin.POST("/refunds/:id/reject", orderHandler.RejectRefund)

			admin.POST("/parts", partHandler.CreatePart)
			admin.GET("/parts", partHandler.GetPartList)
			admin.GET("/parts/:id", partHandler.GetPartDetail)
			admin.PUT("/parts/:id", partHandler.UpdatePart)
			admin.DELETE("/parts/:id", partHandler.DeletePart)

			admin.GET("/part-requests", partHandler.GetPartRequestList)
			admin.GET("/part-requests/:id", partHandler.GetPartRequestDetail)
			admin.POST("/part-requests/:id/approve", partHandler.ApprovePartRequest)
			admin.POST("/part-requests/:id/reject", partHandler.RejectPartRequest)
			admin.POST("/part-requests/:id/ship", partHandler.ShipPartRequest)

			admin.GET("/withdraws", financeHandler.GetWithdrawRequestList)
			admin.POST("/withdraws/:id/approve", financeHandler.ApproveWithdraw)
			admin.POST("/withdraws/:id/reject", financeHandler.RejectWithdraw)
			admin.POST("/withdraws/:id/complete", financeHandler.CompleteWithdraw)

			admin.GET("/finance/report", financeHandler.GetMonthlyReport)
			admin.GET("/finance/performance", financeHandler.GetTechnicianPerformance)
			admin.POST("/finance/settle", financeHandler.SettleTechnicianIncome)

			admin.GET("/reviews/low-rating", adminHandler.GetLowRatingReviews)
			admin.POST("/reviews/:id/intervene", adminHandler.HandleLowRatingReview)

			admin.GET("/order-logs", adminHandler.GetOrderLogs)
		}
	}
}
