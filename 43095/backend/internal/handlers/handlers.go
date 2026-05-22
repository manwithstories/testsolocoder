package handlers

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup) {
	RegisterAuthRoutes(r)
	RegisterDepartmentRoutes(r)
	RegisterDoctorRoutes(r)
	RegisterPatientRoutes(r)
	RegisterAppointmentRoutes(r)
	RegisterConsultationRoutes(r)
	RegisterHealthRecordRoutes(r)
	RegisterNotificationRoutes(r)
	RegisterPaymentRoutes(r)
	RegisterReviewRoutes(r)
	RegisterUploadRoutes(r)
}
