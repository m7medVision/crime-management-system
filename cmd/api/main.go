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

	// Initialize MinIO client
	if err := util.InitMinio(); err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	caseRepo := repository.NewCaseRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.Auth.Secret, cfg.Auth.ExpiryTime)
	caseService := service.NewCaseService(caseRepo, userRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	caseController := controller.NewCaseController(caseService)

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
	}

	// Start server
	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
