package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/config"
	"github.com/m7medVision/crime-management-system/internal/controller"
	"github.com/m7medVision/crime-management-system/internal/middleware"
	"github.com/m7medVision/crime-management-system/internal/model"
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

	// Initialize MinIO client
	if err := util.InitMinio(cfg); err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	caseRepo := repository.NewCaseRepository(db)
	evidenceRepo := repository.NewEvidenceRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.Auth.Secret, cfg.Auth.ExpiryTime)
	caseService := service.NewCaseService(caseRepo, userRepo)
	evidenceService := service.NewEvidenceService(evidenceRepo, caseRepo, userRepo, cfg.Storage.Minio.Bucket)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	caseController := controller.NewCaseController(caseService)
	evidenceController := controller.NewEvidenceController(evidenceService)

	// Setup Gin router
	router := gin.Default()

	// Auth routes
	router.POST("/login", authController.Login)

	// Public routes
	router.POST("/reports", caseController.SubmitCrimeReport)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.BasicAuth(db))
	{
		// Case routes
		protected.POST("/cases", middleware.RequireRole(model.RoleInvestigator, model.RoleAdmin), caseController.CreateCase)
		protected.PUT("/cases/:id", middleware.RequireRole(model.RoleInvestigator, model.RoleAdmin), caseController.UpdateCase)
		protected.GET("/cases/:id", caseController.GetCaseByID)
		protected.GET("/cases", caseController.ListCases)

		protected.GET("/cases/:id/assignees", caseController.GetAssignees)
		protected.POST("/cases/:id/assignees", middleware.RequireRole(model.RoleInvestigator, model.RoleAdmin), caseController.AddAssignee)
		protected.DELETE("/cases/:id/assignees", middleware.RequireRole(model.RoleInvestigator, model.RoleAdmin), caseController.RemoveAssignee)

		// Evidence routes
		protected.POST("/evidence/text", evidenceController.CreateTextEvidence)
		protected.POST("/evidence/image", evidenceController.CreateImageEvidence)
		protected.GET("/evidence/:id", evidenceController.GetEvidenceByID)
		protected.GET("/evidence/:id/image", evidenceController.GetEvidenceImage)
		protected.PUT("/evidence/:id", evidenceController.UpdateEvidence)
		protected.DELETE("/evidence/:id", middleware.RequireRole(model.RoleInvestigator, model.RoleAdmin), evidenceController.SoftDeleteEvidence)
		protected.DELETE("/evidence/:id/permanent", middleware.RequireRole(model.RoleAdmin), evidenceController.HardDeleteEvidence)
		protected.GET("/evidence/:id/audit", middleware.RequireRole(model.RoleAdmin), evidenceController.GetEvidenceAuditLogs)
		protected.GET("/cases/:id/evidence", caseController.GetEvidence)
		protected.GET("/cases/:id/links", caseController.ExtractLinks)
	}

	// Start server
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
