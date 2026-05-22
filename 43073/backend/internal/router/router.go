package router

import (
	"ticket-system/config"
	"ticket-system/internal/controller"
	"ticket-system/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.Static("/uploads", config.App.UploadPath)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			userController := controller.NewUserController()
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}

		upload := api.Group("/upload")
		upload.Use(middleware.Auth())
		{
			uploadController := controller.NewUploadController()
			upload.POST("/image", uploadController.UploadImage)
		}

		activities := api.Group("/activities")
		{
			activityController := controller.NewActivityController()
			activities.GET("", activityController.GetList)
			activities.GET("/:id", activityController.Get)
			activities.Use(middleware.Auth(), middleware.Admin())
			{
				activities.POST("", activityController.Create)
				activities.PUT("/:id", activityController.Update)
				activities.PUT("/:id/status", activityController.UpdateStatus)
				activities.DELETE("/:id", activityController.Delete)
			}
		}

		ticketTypes := api.Group("/ticket-types")
		{
			ticketTypeController := controller.NewTicketTypeController()
			ticketTypes.GET("", ticketTypeController.GetList)
			ticketTypes.GET("/:id", ticketTypeController.Get)
			ticketTypes.Use(middleware.Auth(), middleware.Admin())
			{
				ticketTypes.POST("", ticketTypeController.Create)
				ticketTypes.PUT("/:id", ticketTypeController.Update)
				ticketTypes.DELETE("/:id", ticketTypeController.Delete)
			}
		}

		coupons := api.Group("/coupons")
		coupons.Use(middleware.Auth(), middleware.Admin())
		{
			couponController := controller.NewCouponController()
			coupons.POST("", couponController.Create)
			coupons.GET("", couponController.GetList)
			coupons.GET("/:id", couponController.Get)
			coupons.DELETE("/:id", couponController.Delete)
		}

		orders := api.Group("/orders")
		orders.Use(middleware.Auth())
		{
			orderController := controller.NewOrderController()
			orders.POST("", orderController.Create)
			orders.GET("", orderController.GetList)
			orders.GET("/:id", orderController.Get)
			orders.POST("/:id/pay", orderController.Pay)
			orders.POST("/:id/cancel", orderController.Cancel)
		}

		checkins := api.Group("/checkins")
		checkins.Use(middleware.Auth())
		{
			checkInController := controller.NewCheckInController()
			checkins.POST("", checkInController.CheckIn)
			checkins.GET("", checkInController.GetList)
			checkins.GET("/statistics", checkInController.GetStatistics)
		}

		statistics := api.Group("/statistics")
		statistics.Use(middleware.Auth(), middleware.Admin())
		{
			statisticsController := controller.NewStatisticsController()
			statistics.GET("/activities", statisticsController.GetActivityStatistics)
			statistics.GET("/ticket-types", statisticsController.GetTicketTypeStatistics)
			statistics.GET("/daily", statisticsController.GetDailyStatistics)
			statistics.GET("/export", statisticsController.ExportExcel)
		}

		users := api.Group("/users")
		users.Use(middleware.Auth())
		{
			userController := controller.NewUserController()
			users.GET("/me", userController.GetCurrentUser)
			users.Use(middleware.Admin())
			{
				users.GET("", userController.GetUserList)
				users.GET("/:id", userController.GetUser)
				users.PUT("/:id", userController.UpdateUser)
				users.DELETE("/:id", userController.DeleteUser)
			}
		}
	}
}
