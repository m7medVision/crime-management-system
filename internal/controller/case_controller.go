package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/dto"
	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/service"
)

type CaseController struct {
	caseService *service.CaseService
}

func NewCaseController(caseService *service.CaseService) *CaseController {
	return &CaseController{caseService: caseService}
}

func (ctrl *CaseController) CreateCase(c *gin.Context) {
	var caseDTO dto.CaseDTO
	if err := c.ShouldBindJSON(&caseDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case data"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userID := user.(*model.User).ID
	caseData := &model.Case{
		Name:               caseDTO.Name,
		Description:        caseDTO.Description,
		Area:               caseDTO.Area,
		CaseType:           caseDTO.CaseType,
		AuthorizationLevel: caseDTO.AuthorizationLevel,
		CreatedByID:        userID,
		ReportedByID:       caseDTO.ReportedByID,
	}

	result, err := ctrl.caseService.CreateCase(caseData)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (ctrl *CaseController) UpdateCase(c *gin.Context) {
	var caseDTO dto.CaseDTO
	if err := c.ShouldBindJSON(&caseDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case data"})
		return
	}

	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userID := user.(*model.User).ID
	caseData, err := ctrl.caseService.GetCaseByID(uint(caseID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Case not found"})
		return
	}

	caseData.Name = caseDTO.Name
	caseData.Description = caseDTO.Description
	caseData.Area = caseDTO.Area
	caseData.CaseType = caseDTO.CaseType
	caseData.AuthorizationLevel = caseDTO.AuthorizationLevel
	caseData.CreatedByID = userID
	caseData.ReportedByID = caseDTO.ReportedByID

	result, err := ctrl.caseService.UpdateCase(caseData)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ctrl *CaseController) GetCaseByID(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	caseData, err := ctrl.caseService.GetCaseByID(uint(caseID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Case not found"})
		return
	}

	c.JSON(http.StatusOK, caseData)
}

func (ctrl *CaseController) ListCases(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")

	cases, total, err := ctrl.caseService.ListCases(offset, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cases": cases,
		"total": total,
	})
}

func (ctrl *CaseController) GetAssignees(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	assignees, err := ctrl.caseService.GetAssignees(uint(caseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assignees)
}

func (ctrl *CaseController) AddAssignee(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	var assigneeDTO dto.AssigneeDTO
	if err := c.ShouldBindJSON(&assigneeDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignee data"})
		return
	}

	err = ctrl.caseService.AddAssignee(uint(caseID), assigneeDTO.UserID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignee added successfully"})
}

func (ctrl *CaseController) RemoveAssignee(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	var assigneeDTO dto.AssigneeDTO
	if err := c.ShouldBindJSON(&assigneeDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignee data"})
		return
	}

	err = ctrl.caseService.RemoveAssignee(uint(caseID), assigneeDTO.UserID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignee removed successfully"})
}
