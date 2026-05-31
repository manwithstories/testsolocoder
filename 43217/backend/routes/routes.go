package routes

import (
	"health-platform/controllers"
	"health-platform/middleware"
	"health-platform/models"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	authCtrl := controllers.NewAuthController()
	companyCtrl := controllers.NewCompanyController()
	agencyCtrl := controllers.NewAgencyController()
	appointmentCtrl := controllers.NewAppointmentController()
	reportCtrl := controllers.NewReportController()
	healthCtrl := controllers.NewHealthController()
	billingCtrl := controllers.NewBillingController()
	statisticsCtrl := controllers.NewStatisticsController()

	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", authCtrl.Register)
		api.POST("/auth/login", authCtrl.Login)

		api.POST("/companies/register", companyCtrl.RegisterCompany)
		api.POST("/agencies/register", agencyCtrl.RegisterAgency)

		api.GET("/agencies", agencyCtrl.ListAgencies)
		api.GET("/agencies/:id", agencyCtrl.GetAgency)
		api.GET("/packages", agencyCtrl.ListOnlinePackages)
		api.GET("/packages/:id", agencyCtrl.GetPackage)
		api.GET("/packages/hot", agencyCtrl.GetHotPackages)
		api.GET("/packages/:id/timeslots", agencyCtrl.GetPackageTimeSlots)

		auth := api.Group("")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/auth/userinfo", authCtrl.GetUserInfo)

			hr := auth.Group("")
			hr.Use(middleware.RoleAuth(string(models.RoleHR), string(models.RoleAdmin)))
			{
				hr.GET("/companies/:id", companyCtrl.GetCompany)
				hr.PUT("/companies", companyCtrl.UpdateCompany)
				
				hr.POST("/departments", companyCtrl.AddDepartment)
				hr.PUT("/departments/:id", companyCtrl.UpdateDepartment)
				hr.GET("/departments", companyCtrl.GetDepartments)

				hr.POST("/employees", companyCtrl.AddEmployee)
				hr.PUT("/employees/:id", companyCtrl.UpdateEmployee)
				hr.GET("/employees", companyCtrl.GetEmployees)
				hr.GET("/employees/:id", companyCtrl.GetEmployee)

				hr.POST("/budgets", companyCtrl.SetBudget)
				hr.GET("/budgets", companyCtrl.GetBudget)

				hr.POST("/department-appointments", companyCtrl.SetDepartmentAppointment)
				hr.GET("/departments/:department_id/appointments", companyCtrl.GetDepartmentAppointments)

				hr.GET("/company/appointments", appointmentCtrl.GetCompanyAppointments)
				hr.GET("/company/reports", reportCtrl.GetCompanyReports)

				hr.GET("/company/billings", billingCtrl.GetCompanyBillings)
				hr.GET("/company/balance", billingCtrl.GetCompanyBalance)
				hr.POST("/company/recharge", billingCtrl.Recharge)
				hr.GET("/company/transactions", billingCtrl.GetTransactions)

				hr.GET("/statistics/company", statisticsCtrl.GetCompanyStatistics)
				hr.GET("/statistics/department/:department_id", statisticsCtrl.GetDepartmentStatistics)
				hr.GET("/statistics/age-distribution", statisticsCtrl.GetAgeDistribution)
				hr.GET("/statistics/gender-distribution", statisticsCtrl.GetGenderDistribution)
				hr.GET("/statistics/abnormal-distribution", statisticsCtrl.GetAbnormalDistribution)
				hr.GET("/statistics/export", statisticsCtrl.ExportStatistics)
			}

			agencyAdmin := auth.Group("")
			agencyAdmin.Use(middleware.RoleAuth(string(models.RoleAgency), string(models.RoleAdmin)))
			{
				agencyAdmin.GET("/agencies/:id", agencyCtrl.GetAgency)
				agencyAdmin.PUT("/agencies", agencyCtrl.UpdateAgency)

				agencyAdmin.POST("/packages", agencyCtrl.CreatePackage)
				agencyAdmin.PUT("/packages/:id", agencyCtrl.UpdatePackage)
				agencyAdmin.PATCH("/packages/:id/price", agencyCtrl.UpdatePackagePrice)
				agencyAdmin.PATCH("/packages/:id/status", agencyCtrl.UpdatePackageStatus)
				agencyAdmin.GET("/agency/packages", agencyCtrl.GetAgencyPackages)

				agencyAdmin.POST("/timeslots", agencyCtrl.CreateTimeSlot)

				agencyAdmin.GET("/agency/appointments", appointmentCtrl.GetAgencyAppointments)
				agencyAdmin.PATCH("/appointments/:id/complete", appointmentCtrl.CompleteAppointment)

				agencyAdmin.POST("/reports", reportCtrl.CreateReport)
				agencyAdmin.POST("/reports/upload", reportCtrl.UploadReportFile)

				agencyAdmin.GET("/agency/billings", billingCtrl.GetAgencyBillings)
				agencyAdmin.POST("/billings/generate", billingCtrl.GenerateMonthlyBilling)

				agencyAdmin.GET("/statistics/agency-rating", statisticsCtrl.GetAgencyRating)
				agencyAdmin.GET("/statistics/package-ranking", statisticsCtrl.GetPackageRanking)
			}

			employee := auth.Group("")
			employee.Use(middleware.RoleAuth(string(models.RoleEmployee), string(models.RoleAdmin)))
			{
				employee.POST("/appointments", appointmentCtrl.CreateAppointment)
				employee.GET("/appointments/:id", appointmentCtrl.GetAppointment)
				employee.PUT("/appointments/:id/reschedule", appointmentCtrl.RescheduleAppointment)
				employee.PUT("/appointments/:id/cancel", appointmentCtrl.CancelAppointment)
				employee.GET("/employees/:employee_id/appointments", appointmentCtrl.GetEmployeeAppointments)
				employee.GET("/employees/:employee_id/appointment-status", appointmentCtrl.GetEmployeeAppointmentStatus)
				employee.GET("/employees/:employee_id/check-quota", appointmentCtrl.CheckQuota)

				employee.GET("/employees/:employee_id/reports", reportCtrl.GetEmployeeReports)
				employee.GET("/reports/:id", reportCtrl.GetReport)
				employee.GET("/appointments/:appointment_id/report", reportCtrl.GetReportByAppointment)
				employee.GET("/reports/:id/download", reportCtrl.DownloadReportFile)
				employee.GET("/employees/:employee_id/abnormal-reports", reportCtrl.GetAbnormalReports)
				employee.GET("/employees/:employee_id/abnormal-items", reportCtrl.GetAbnormalItems)

				employee.GET("/employees/:employee_id/health-records", healthCtrl.GetHealthRecords)
				employee.GET("/employees/:employee_id/health-records/:year", healthCtrl.GetHealthRecordByYear)
				employee.GET("/employees/:employee_id/health-records/all", healthCtrl.GetAllHealthRecords)
				employee.GET("/employees/:employee_id/health-trend", healthCtrl.GetTrendData)
				employee.GET("/employees/:employee_id/health-summary", healthCtrl.GetHealthSummary)

				employee.GET("/employees/:employee_id/abnormal-items", healthCtrl.GetAbnormalItems)
				employee.GET("/employees/:employee_id/abnormal-items/all", healthCtrl.GetAllAbnormalItems)
				employee.POST("/abnormal-items", healthCtrl.CreateAbnormalItem)
				employee.POST("/abnormal-items/set-recheck", healthCtrl.SetRecheckDate)
				employee.PUT("/abnormal-items/recheck-status", healthCtrl.UpdateRecheckStatus)
				employee.GET("/employees/:employee_id/need-recheck", healthCtrl.GetNeedRecheckItems)

				employee.GET("/employees/:employee_id/reminders", healthCtrl.GetReminders)
				employee.GET("/employees/:employee_id/reminders/unread", healthCtrl.GetUnreadReminders)
				employee.PUT("/reminders/:id/read", healthCtrl.MarkReminderAsRead)
			}

			admin := auth.Group("/admin")
			admin.Use(middleware.RoleAuth(string(models.RoleAdmin)))
			{
				admin.DELETE("/statistics/cache", statisticsCtrl.ClearStatisticsCache)
			}

			billing := auth.Group("")
			billing.Use(middleware.RoleAuth(string(models.RoleHR), string(models.RoleAdmin), string(models.RoleAgency)))
			{
				billing.GET("/billings/:id", billingCtrl.GetBilling)
				billing.POST("/billings/pay", billingCtrl.PayBilling)
			}
		}
	}
}
