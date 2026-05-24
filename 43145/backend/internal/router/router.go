package router

import (
	"survey-platform/internal/handler"
	"survey-platform/internal/middleware"
	"survey-platform/internal/repository"
	"survey-platform/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	userRepo := repository.NewUserRepository(db)
	surveyRepo := repository.NewSurveyRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	responseRepo := repository.NewResponseRepository(db)
	distRepo := repository.NewDistributionRepository(db)

	emailSvc := service.NewEmailService()
	authSvc := service.NewAuthService(userRepo)
	surveySvc := service.NewSurveyService(surveyRepo)
	questionSvc := service.NewQuestionService(questionRepo, surveyRepo, responseRepo)
	responseSvc := service.NewResponseService(responseRepo, surveyRepo, distRepo)
	distSvc := service.NewDistributionService(distRepo, surveyRepo, emailSvc)
	statisticsSvc := service.NewStatisticsService(responseRepo, surveyRepo, questionRepo)
	exportSvc := service.NewExportService(surveyRepo, responseRepo, questionRepo, statisticsSvc)

	authHdlr := handler.NewAuthHandler(authSvc)
	surveyHdlr := handler.NewSurveyHandler(surveySvc)
	questionHdlr := handler.NewQuestionHandler(questionSvc)
	responseHdlr := handler.NewResponseHandler(responseSvc, distSvc)
	distHdlr := handler.NewDistributionHandler(distSvc)
	statisticsHdlr := handler.NewStatisticsHandler(statisticsSvc)
	exportHdlr := handler.NewExportHandler(exportSvc)
	userHdlr := handler.NewUserHandler(userRepo)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHdlr.Register)
			auth.POST("/login", authHdlr.Login)
		}

		authRequired := api.Group("")
		authRequired.Use(middleware.Auth())
		{
			authRequired.GET("/profile", authHdlr.GetProfile)
			authRequired.PUT("/profile", authHdlr.UpdateProfile)
			authRequired.PUT("/password", authHdlr.ChangePassword)

			surveys := authRequired.Group("/surveys")
			{
				surveys.POST("", surveyHdlr.Create)
				surveys.GET("", surveyHdlr.List)
				surveys.GET("/all", middleware.RequireRole("admin"), surveyHdlr.ListAll)
				surveys.GET("/:id", surveyHdlr.GetByID)
				surveys.PUT("/:id", surveyHdlr.Update)
				surveys.DELETE("/:id", surveyHdlr.Delete)
				surveys.POST("/:id/publish", surveyHdlr.Publish)
				surveys.POST("/:id/close", surveyHdlr.Close)
				surveys.POST("/:id/copy", surveyHdlr.Copy)
			}

			questions := authRequired.Group("/surveys/:survey_id/questions")
			{
				questions.POST("", questionHdlr.Create)
				questions.GET("", questionHdlr.GetBySurveyID)
				questions.POST("/batch", questionHdlr.BatchCreate)
				questions.PUT("/:id", questionHdlr.Update)
				questions.DELETE("/:id", questionHdlr.Delete)
				questions.PUT("/:id/reorder", questionHdlr.Reorder)
			}

			questionById := authRequired.Group("/questions")
			{
				questionById.GET("/:id", questionHdlr.GetByID)
			}

			distribution := authRequired.Group("/surveys/:survey_id/distribution")
			{
				distribution.POST("/link", distHdlr.CreateLink)
				distribution.GET("/links", distHdlr.ListBySurveyID)
				distribution.POST("/invitations", distHdlr.SendInvitations)
				distribution.GET("/invitations", distHdlr.ListInvitations)
			}

			distLinks := authRequired.Group("/distribution")
			{
				distLinks.DELETE("/:id", distHdlr.DeleteLink)
			}

			responses := authRequired.Group("/surveys/:survey_id/responses")
			{
				responses.GET("", responseHdlr.List)
				responses.POST("", responseHdlr.StartResponse)
				responses.POST("/save", responseHdlr.SaveProgress)
				responses.POST("/submit", responseHdlr.SubmitResponse)
			}

			responseById := authRequired.Group("/responses")
			{
				responseById.GET("/:id", responseHdlr.GetByID)
				responseById.DELETE("/:id", responseHdlr.Delete)
			}

			statistics := authRequired.Group("/statistics")
			{
				statistics.GET("", statisticsHdlr.GetStatistics)
				statistics.GET("/cross-analysis", statisticsHdlr.CrossAnalysis)
			}

			export := authRequired.Group("/export")
			{
				export.GET("/:survey_id/excel", exportHdlr.ExportExcel)
				export.GET("/:survey_id/pdf", exportHdlr.ExportPDF)
				export.GET("/:survey_id/charts", exportHdlr.ExportChartImages)
			}

			users := authRequired.Group("/users")
			{
				users.GET("", middleware.RequireRole("admin"), userHdlr.List)
				users.GET("/:id", middleware.RequireRole("admin"), userHdlr.GetByID)
				users.PUT("/:id/role", middleware.RequireRole("admin"), userHdlr.UpdateRole)
				users.PUT("/:id/status", middleware.RequireRole("admin"), userHdlr.UpdateStatus)
				users.DELETE("/:id", middleware.RequireRole("admin"), userHdlr.Delete)
			}
		}

		public := api.Group("/public")
		{
			public.GET("/survey/:token", distHdlr.GetByToken)
			public.POST("/survey/:survey_id/validate", responseHdlr.ValidateAccess)
			public.POST("/survey/:survey_id/start", responseHdlr.StartResponse)
			public.POST("/survey/:survey_id/save", responseHdlr.SaveProgress)
			public.POST("/survey/:survey_id/submit", responseHdlr.SubmitResponse)
		}
	}

	return r
}
