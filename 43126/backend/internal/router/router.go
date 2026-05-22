package router

import (
	"meeting-room/internal/handlers"
	"meeting-room/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.Logger())

	r.Static("/uploads", "./uploads")

	userHandler := handlers.NewUserHandler()
	roomHandler := handlers.NewRoomHandler()
	bookingHandler := handlers.NewBookingHandler()
	calendarHandler := handlers.NewCalendarHandler()
	materialHandler := handlers.NewMaterialHandler()
	statsHandler := handlers.NewStatsHandler()

	api := r.Group("/api")
	{
		api.POST("/auth/login", userHandler.Login)
		api.POST("/auth/register", userHandler.Register)
	}

	auth := api.Group("")
	auth.Use(middleware.JWTAuth())
	{
		auth.GET("/user/profile", userHandler.GetProfile)
		auth.PUT("/user/profile", userHandler.UpdateProfile)

		auth.GET("/rooms/floors", roomHandler.GetFloors)
		auth.GET("/rooms/list", roomHandler.ListAllRooms)
		auth.GET("/rooms", roomHandler.ListRooms)
		auth.GET("/rooms/:id", roomHandler.GetRoom)

		auth.POST("/bookings", bookingHandler.CreateBooking)
		auth.GET("/bookings", bookingHandler.ListBookings)
		auth.GET("/bookings/:id", bookingHandler.GetBooking)
		auth.DELETE("/bookings/:id", bookingHandler.CancelBooking)
		auth.PUT("/bookings/:id/reschedule", bookingHandler.RescheduleBooking)

		auth.GET("/calendar/week", calendarHandler.GetWeekCalendar)
		auth.GET("/calendar/month", calendarHandler.GetMonthCalendar)
		auth.GET("/calendar/rooms/:id/availability", calendarHandler.GetRoomAvailability)

		auth.POST("/bookings/:id/materials", materialHandler.UploadMaterial)
		auth.GET("/bookings/:id/materials", materialHandler.GetMaterials)
		auth.GET("/materials/:id/download", materialHandler.DownloadMaterial)
		auth.DELETE("/materials/:id", materialHandler.DeleteMaterial)

		auth.GET("/stats", statsHandler.GetStats)
		auth.GET("/stats/export", statsHandler.ExportStats)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(), middleware.AdminOnly())
	{
		admin.GET("/users", userHandler.ListUsers)
		admin.PUT("/users/:id/role", userHandler.UpdateRole)
		admin.DELETE("/users/:id", userHandler.DeleteUser)

		admin.PUT("/bookings/:id/approve", bookingHandler.ApproveBooking)
		admin.PUT("/bookings/:id/complete", bookingHandler.CompleteBooking)
	}

	spaceAdmin := api.Group("/admin")
	spaceAdmin.Use(middleware.JWTAuth(), middleware.SpaceAdminOrAdmin())
	{
		spaceAdmin.POST("/rooms", roomHandler.CreateRoom)
		spaceAdmin.PUT("/rooms/:id", roomHandler.UpdateRoom)
		spaceAdmin.DELETE("/rooms/:id", roomHandler.DeleteRoom)
		spaceAdmin.POST("/rooms/:id/photos", roomHandler.UploadPhoto)
		spaceAdmin.DELETE("/rooms/photos/:id", roomHandler.DeletePhoto)
	}

	return r
}
