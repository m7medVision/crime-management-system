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

func (ctrl *CaseController) SubmitCrimeReport(c *gin.Context) {
	var reportDTO dto.ReportDTO
	if err := c.ShouldBindJSON(&reportDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report data"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userID := user.(*model.User).ID
	report := &model.Report{
		Title:        reportDTO.Title,
		Description:  reportDTO.Description,
		Location:     reportDTO.Location,
		ReportedByID: userID,
	}

	result, err := ctrl.caseService.SubmitCrimeReport(report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"reportId": result.ID})
}

func (ctrl *CaseController) GetEvidence(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	evidence, err := ctrl.caseService.GetEvidence(uint(caseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, evidence)
}

func (ctrl *CaseController) GetSuspects(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	suspects, err := ctrl.caseService.GetSuspects(uint(caseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, suspects)
}

func (ctrl *CaseController) GetVictims(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	victims, err := ctrl.caseService.GetVictims(uint(caseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, victims)
}

func (ctrl *CaseController) GetWitnesses(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	witnesses, err := ctrl.caseService.GetWitnesses(uint(caseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, witnesses)
}

func (c *CaseController) ExtractLinks(ctx *gin.Context) {
	// Parse case ID from URL parameter
	caseIDStr := ctx.Param("id")
	caseID, err := strconv.ParseUint(caseIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	// Call service to extract links
	links, err := c.caseService.ExtractLinksFromCase(uint(caseID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the extracted links
	ctx.JSON(http.StatusOK, gin.H{"links": links})
}

// GetCaseStatusByReportID returns the case status for a given report ID
func (ctrl *CaseController) GetCaseStatusByReportID(c *gin.Context) {
	reportID := c.Param("reportId")
	if reportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID is required"})
		return
	}

	status, err := ctrl.caseService.GetCaseStatusByReportID(reportID)
	if err != nil {
		if err.Error() == "report not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reportId": reportID,
		"status":   status,
	})
}
