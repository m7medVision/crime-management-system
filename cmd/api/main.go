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
	"github.com/m7medVision/crime-management-system/internal/model"
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

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.Auth.Secret, cfg.Auth.ExpiryTime)
	caseService := service.NewCaseService(caseRepo, userRepo)
	userService := service.NewUserService(userRepo) // Initialize userService

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	caseController := controller.NewCaseController(caseService)
	userController := controller.NewUserController(userService) // Initialize userController

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

		// User management routes
		protected.POST("/users", middleware.RequireRole(model.RoleAdmin), userController.CreateUser)
		protected.PUT("/users/:id", middleware.RequireRole(model.RoleAdmin), userController.UpdateUser)
		protected.DELETE("/users/:id", middleware.RequireRole(model.RoleAdmin), userController.DeleteUser)
	}

	// Start server
	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
