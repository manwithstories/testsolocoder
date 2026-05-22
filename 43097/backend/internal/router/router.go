package router

import (
	"hotel-system/internal/database"
	"hotel-system/internal/handler"
	"hotel-system/internal/middleware"
	"hotel-system/internal/repository"
	"hotel-system/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	db := database.GetDB()

	userRepo := repository.NewUserRepository(db)
	roomTypeRepo := repository.NewRoomTypeRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	checkInRepo := repository.NewCheckInRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	memberLevelRepo := repository.NewMemberLevelRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	userService := service.NewUserService(userRepo)
	roomService := service.NewRoomService(roomTypeRepo, roomRepo, db)
	bookingService := service.NewBookingService(bookingRepo, roomRepo, db)
	checkInService := service.NewCheckInService(checkInRepo, bookingRepo, roomRepo, db)
	memberService := service.NewMemberService(memberRepo, memberLevelRepo, db)
	paymentService := service.NewPaymentService(paymentRepo, db)
	reportService := service.NewReportService(db)
	dashboardService := service.NewDashboardService(db)

	authHandler := handler.NewAuthHandler(userService)
	userHandler := handler.NewUserHandler(userService)
	roomHandler := handler.NewRoomHandler(roomService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	checkInHandler := handler.NewCheckInHandler(checkInService)
	memberHandler := handler.NewMemberHandler(memberService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	reportHandler := handler.NewReportHandler(reportService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.GET("/profile", middleware.JWTAuth(), authHandler.GetProfile)
		}

		authorized := api.Group("/")
		authorized.Use(middleware.JWTAuth())
		{
			users := authorized.Group("/users")
			{
				users.GET("", middleware.AdminRequired(), userHandler.ListUsers)
				users.GET("/:id", middleware.AdminRequired(), userHandler.GetUser)
				users.POST("", middleware.AdminRequired(), userHandler.CreateUser)
				users.PUT("/:id", middleware.AdminRequired(), userHandler.UpdateUser)
				users.DELETE("/:id", middleware.AdminRequired(), userHandler.DeleteUser)
			}

			roomTypes := authorized.Group("/room-types")
			{
				roomTypes.GET("", roomHandler.ListRoomTypes)
				roomTypes.GET("/:id", roomHandler.GetRoomType)
				roomTypes.POST("", middleware.AdminRequired(), roomHandler.CreateRoomType)
				roomTypes.PUT("/:id", middleware.AdminRequired(), roomHandler.UpdateRoomType)
				roomTypes.DELETE("/:id", middleware.AdminRequired(), roomHandler.DeleteRoomType)
			}

			rooms := authorized.Group("/rooms")
			{
				rooms.GET("", roomHandler.ListRooms)
				rooms.GET("/available", roomHandler.GetAvailableRooms)
				rooms.GET("/dashboard", roomHandler.GetRoomDashboard)
				rooms.GET("/:id", roomHandler.GetRoom)
				rooms.POST("", middleware.AdminRequired(), roomHandler.CreateRoom)
				rooms.POST("/batch-import", middleware.AdminRequired(), roomHandler.BatchImportRooms)
				rooms.PUT("/:id", middleware.FrontDeskRequired(), roomHandler.UpdateRoom)
				rooms.PUT("/:id/status", middleware.FrontDeskRequired(), roomHandler.UpdateRoomStatus)
				rooms.DELETE("/:id", middleware.AdminRequired(), roomHandler.DeleteRoom)
			}

			bookings := authorized.Group("/bookings")
			{
				bookings.GET("", bookingHandler.ListBookings)
				bookings.GET("/:id", bookingHandler.GetBooking)
				bookings.POST("", bookingHandler.CreateBooking)
				bookings.POST("/calculate-price", bookingHandler.CalculatePrice)
				bookings.PUT("/:id", bookingHandler.UpdateBooking)
				bookings.POST("/:id/cancel", bookingHandler.CancelBooking)
				bookings.POST("/:id/confirm", bookingHandler.ConfirmBooking)
			}

			checkins := authorized.Group("/checkins")
			{
				checkins.GET("", checkInHandler.ListCheckIns)
				checkins.GET("/:id", checkInHandler.GetCheckIn)
				checkins.POST("", middleware.FrontDeskRequired(), checkInHandler.CreateCheckIn)
				checkins.POST("/:id/checkout", middleware.FrontDeskRequired(), checkInHandler.CheckOut)
				checkins.POST("/:id/extend", middleware.FrontDeskRequired(), checkInHandler.ExtendStay)
				checkins.POST("/:id/extra-charge", middleware.FrontDeskRequired(), checkInHandler.AddExtraCharge)
			}

			payments := authorized.Group("/payments")
			{
				payments.GET("", paymentHandler.ListPayments)
				payments.GET("/:id", paymentHandler.GetPayment)
				payments.GET("/order", paymentHandler.GetOrderPayments)
				payments.POST("", middleware.FrontDeskRequired(), paymentHandler.CreatePayment)
				payments.POST("/refund", middleware.AdminRequired(), paymentHandler.RefundPayment)
				payments.GET("/:id/voucher", paymentHandler.GeneratePaymentVoucher)
			}

			members := authorized.Group("/members")
			{
				members.GET("", middleware.FrontDeskRequired(), memberHandler.ListMembers)
				members.GET("/by-phone", middleware.FrontDeskRequired(), memberHandler.GetMemberByPhone)
				members.GET("/:id", memberHandler.GetMember)
				members.POST("", middleware.FrontDeskRequired(), memberHandler.RegisterMember)
				members.PUT("/:id", middleware.FrontDeskRequired(), memberHandler.UpdateMember)
				members.DELETE("/:id", middleware.AdminRequired(), memberHandler.DeleteMember)
				members.GET("/:id/discount", memberHandler.GetMemberDiscount)
				members.GET("/:id/history", memberHandler.GetMemberConsumptionHistory)
				members.POST("/points/use", memberHandler.UsePoints)
				members.POST("/points/recharge", middleware.FrontDeskRequired(), memberHandler.RechargePoints)

				memberLevels := members.Group("/levels")
				{
					memberLevels.GET("", memberHandler.ListMemberLevels)
					memberLevels.GET("/:id", memberHandler.GetMemberLevel)
					memberLevels.POST("", middleware.AdminRequired(), memberHandler.CreateMemberLevel)
					memberLevels.PUT("/:id", middleware.AdminRequired(), memberHandler.UpdateMemberLevel)
					memberLevels.DELETE("/:id", middleware.AdminRequired(), memberHandler.DeleteMemberLevel)
				}
			}

			reports := authorized.Group("/reports")
			{
				reports.GET("/occupancy-rate", middleware.AdminRequired(), reportHandler.GetOccupancyRate)
				reports.GET("/revenue", middleware.AdminRequired(), reportHandler.GetRevenueReport)
				reports.GET("/export", middleware.AdminRequired(), reportHandler.ExportReport)
			}

			dashboard := authorized.Group("/dashboard")
			{
				dashboard.GET("/stats", middleware.FrontDeskRequired(), dashboardHandler.GetDashboardStats)
				dashboard.GET("/room-status", middleware.FrontDeskRequired(), dashboardHandler.GetRoomStatusBoard)
				dashboard.GET("/floor/:floor", middleware.FrontDeskRequired(), dashboardHandler.GetFloorRooms)
			}
		}
	}
}
