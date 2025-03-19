package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/config"
	"github.com/m7medVision/crime-management-system/internal/controller"
	"github.com/m7medVision/crime-management-system/internal/middleware"
	"github.com/m7medVision/crime-management-system/internal/repository"
	"github.com/m7medVision/crime-management-system/internal/service"
	"github.com/m7medVision/crime-management-system/internal/util"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := util.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	caseRepo := repository.NewCaseRepository(db)
	evidenceRepo := repository.NewEvidenceRepository(db)
	reportRepo := repository.NewReportRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)
	textAnalysisRepo := repository.NewTextAnalysisRepository(db)
	linkExtractionRepo := repository.NewLinkExtractionRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.Auth.Secret, cfg.Auth.ExpiryTime)
	caseService := service.NewCaseService(caseRepo, userRepo)
	evidenceService := service.NewEvidenceService(evidenceRepo)
	reportService := service.NewReportService(reportRepo, caseRepo)
	textAnalysisService := service.NewTextAnalysisService(textAnalysisRepo)
	linkExtractionService := service.NewLinkExtractionService(caseRepo)
	auditLogService := service.NewAuditLogService(auditLogRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	caseController := controller.NewCaseController(caseService)
	evidenceController := controller.NewEvidenceController(evidenceService)
	reportController := controller.NewReportController(reportService)
	textAnalysisController := controller.NewTextAnalysisController(textAnalysisService)
	linkExtractionController := controller.NewLinkExtractionController(linkExtractionService)
	auditLogController := controller.NewAuditLogController(auditLogService)

	// Setup Gin router
	router := gin.Default()

	// Auth routes
	router.POST("/login", authController.Login)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.BasicAuth(db))
	{
		protected.POST("/cases", caseController.CreateCase)
		protected.PUT("/cases/:id", caseController.UpdateCase)
		protected.GET("/cases/:id", caseController.GetCaseByID)
		protected.GET("/cases", caseController.ListCases)

		protected.GET("/cases/:id/assignees", caseController.GetAssignees)
		protected.POST("/cases/:id/assignees", caseController.AddAssignee)
		protected.DELETE("/cases/:id/assignees", caseController.RemoveAssignee)

		protected.POST("/evidence", evidenceController.RecordEvidence)
		protected.GET("/evidence/:id", evidenceController.GetEvidenceByID)
		protected.GET("/evidence/:id/image", evidenceController.GetEvidenceImageByID)
		protected.PUT("/evidence/:id", evidenceController.UpdateEvidence)
		protected.DELETE("/evidence/:id", evidenceController.SoftDeleteEvidence)

		protected.GET("/audit-logs", auditLogController.GetAdminLogs)
	}

	// Public routes
	public := router.Group("/")
	{
		public.POST("/reports", reportController.GenerateReport)
		public.GET("/reports/:id/status", reportController.GetReportStatus)
		public.GET("/text-analysis", textAnalysisController.ExtractTopWords)
		public.GET("/cases/:id/links", linkExtractionController.ExtractLinks)
	}

	// Start server
	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
