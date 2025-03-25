package main

import (
	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/controller"
	"github.com/m7medVision/crime-management-system/internal/repository"
	"github.com/m7medVision/crime-management-system/internal/service"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	caseRepo := &repository.CaseRepository{} // Initialize CaseRepository
	userRepo := &repository.UserRepository{} // Initialize UserRepository
	caseService := service.NewCaseService(caseRepo, userRepo)
	caseController := controller.NewCaseController(caseService)

	api := router.Group("/api")
	{
		api.POST("/reports", caseController.SubmitCrimeReport)
		api.GET("/cases", caseController.ListCases)
		api.POST("/cases", caseController.CreateCase)
		api.PUT("/cases/:id", caseController.UpdateCase)
		api.GET("/cases/:id", caseController.GetCaseByID)
		api.GET("/cases/:id/assignees", caseController.GetAssignees)
		api.POST("/cases/:id/assignees", caseController.AddAssignee)
		api.DELETE("/cases/:id/assignees", caseController.RemoveAssignee)
		api.GET("/cases/:id/evidence", caseController.GetEvidence)
		api.GET("/cases/:id/suspects", caseController.GetSuspects)
		api.GET("/cases/:id/victims", caseController.GetVictims)
		api.GET("/cases/:id/witnesses", caseController.GetWitnesses)
	}

	return router
}
