package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/config"
	"github.com/m7medVision/crime-management-system/internal/controller"
	"github.com/m7medVision/crime-management-system/internal/middleware"
	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/repository"
	"github.com/m7medVision/crime-management-system/internal/service"
	"github.com/m7medVision/crime-management-system/internal/util"

	_ "github.com/m7medVision/crime-management-system/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title District Core Crime Management System API
// @version 1.0
// @description API service for the District Core Crime Management System

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.basic BasicAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

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

	// Ensure temp directory exists
	tmpDir := "./tmp"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	caseRepo := repository.NewCaseRepository(db)
	evidenceRepo := repository.NewEvidenceRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.Auth.Secret, cfg.Auth.ExpiryTime)
	caseService := service.NewCaseService(caseRepo, userRepo)
	evidenceService := service.NewEvidenceService(evidenceRepo, caseRepo, userRepo, cfg.Storage.Minio.Bucket)
	userService := service.NewUserService(userRepo)

	// Initialize report service
	templatePath := filepath.Join("templates", "case_report.tex")
	reportService, err := service.NewReportService(caseRepo, templatePath)
	if err != nil {
		log.Fatalf("Failed to initialize report service: %v", err)
	}

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	caseController := controller.NewCaseController(caseService)
	evidenceController := controller.NewEvidenceController(evidenceService)
	reportController := controller.NewReportController(reportService)
	userController := controller.NewUserController(userService)

	// Setup Gin router
	router := gin.Default()

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes
	router.POST("/login", authController.Login)

	// Public routes
	public := router.Group("/api/public")
	public.POST("/reports", caseController.SubmitCrimeReport)
	public.GET("/reports/:reportId/status", caseController.GetCaseStatusByReportID)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.BasicAuth(db))
	{
		// User management routes (admin only)
		userRoutes := protected.Group("/users")
		userRoutes.Use(middleware.RequireRole(model.RoleAdmin))
		{
			userRoutes.POST("", userController.CreateUser)
			userRoutes.GET("", userController.ListUsers)
			userRoutes.GET("/:id", userController.GetUser)
			userRoutes.PUT("/:id", userController.UpdateUser)
			userRoutes.DELETE("/:id", userController.DeleteUser)
		}

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

		// Report routes
		protected.GET("/cases/:id/report", reportController.GenerateCaseReport)
	}

	// Start server
	port := cfg.Server.Port
	addr := "0.0.0.0:" + port
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
