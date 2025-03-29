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

// CreateCase godoc
// @Summary Create a new case
// @Description Create a new criminal case record
// @Tags cases
// @Accept json
// @Produce json
// @Param case body dto.CaseDTO true "Case details"
// @Success 201 {object} model.Case
// @Failure 400 {object} map[string]string "Invalid case data"
// @Failure 403 {object} map[string]string "Permission denied"
// @Security ApiKeyAuth
// @Router /cases [post]
func (ctrl *CaseController) CreateCase(c *gin.Context) {
	var caseDTO dto.CaseDTO
	if err := c.ShouldBindJSON(&caseDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid case data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorDTO{
			Message: "Authentication required",
			Code:    http.StatusUnauthorized,
		})
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
		c.JSON(http.StatusForbidden, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusForbidden,
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// UpdateCase godoc
// @Summary Update an existing case
// @Description Update an existing criminal case record
// @Tags cases
// @Accept json
// @Produce json
// @Param id path int true "Case ID"
// @Param case body dto.CaseDTO true "Updated case details"
// @Success 200 {object} model.Case
// @Failure 400 {object} map[string]string "Invalid case data"
// @Failure 403 {object} map[string]string "Permission denied"
// @Failure 404 {object} map[string]string "Case not found"
// @Security ApiKeyAuth
// @Router /cases/{id} [put]
func (ctrl *CaseController) UpdateCase(c *gin.Context) {
	var caseDTO dto.CaseDTO
	if err := c.ShouldBindJSON(&caseDTO); err != nil {
		c.JSON(http.StatusBadRequest,
			dto.ErrorDTO{
				Message: "Invalid case data",
				Code:    http.StatusBadRequest,
			})
		return
	}

	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid case ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorDTO{
			Message: "Authentication required",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	userID := user.(*model.User).ID
	caseData, err := ctrl.caseService.GetCaseByID(uint(caseID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorDTO{
			Message: "Case not found",
			Code:    http.StatusNotFound,
		})
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
		c.JSON(http.StatusForbidden, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusForbidden,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCaseByID godoc
// @Summary Get case details
// @Description Retrieve detailed information about a specific case
// @Tags cases
// @Accept json
// @Produce json
// @Param id path int true "Case ID"
// @Success 200 {object} model.Case
// @Failure 400 {object} map[string]string "Invalid case ID"
// @Failure 404 {object} map[string]string "Case not found"
// @Security ApiKeyAuth
// @Router /cases/{id} [get]
func (ctrl *CaseController) GetCaseByID(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid case ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	caseData, err := ctrl.caseService.GetCaseByID(uint(caseID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorDTO{
			Message: "Case not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, caseData)
}

// ListCases godoc
// @Summary List all cases
// @Description Get a paginated list of cases with optional search
// @Tags cases
// @Accept json
// @Produce json
// @Param offset query int false "Pagination offset" default(0)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term for case name or description"
// @Success 200 {object} map[string]interface{} "cases and total count"
// @Failure 500 {object} map[string]string "Server error"
// @Security ApiKeyAuth
// @Router /cases [get]
func (ctrl *CaseController) ListCases(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")

	cases, total, err := ctrl.caseService.ListCases(offset, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cases": cases,
		"total": total,
	})
}

// GetAssignees godoc
// @Summary Get case assignees
// @Description Retrieve list of users assigned to a case
// @Tags cases,assignees
// @Accept json
// @Produce json
// @Param id path int true "Case ID"
// @Success 200 {array} model.User
// @Failure 400 {object} map[string]string "Invalid case ID"
// @Failure 500 {object} map[string]string "Server error"
// @Security ApiKeyAuth
// @Router /cases/{id}/assignees [get]
func (ctrl *CaseController) GetAssignees(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid case ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	assignees, err := ctrl.caseService.GetAssignees(uint(caseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, assignees)
}

func (ctrl *CaseController) AddAssignee(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid case ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var assigneeDTO dto.AssigneeDTO
	if err := c.ShouldBindJSON(&assigneeDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid assignee data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = ctrl.caseService.AddAssignee(uint(caseID), assigneeDTO.UserID)
	if err != nil {
		c.JSON(http.StatusForbidden, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusForbidden,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignee added successfully"})
}

func (ctrl *CaseController) RemoveAssignee(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid case ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var assigneeDTO dto.AssigneeDTO
	if err := c.ShouldBindJSON(&assigneeDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid assignee data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = ctrl.caseService.RemoveAssignee(uint(caseID), assigneeDTO.UserID)
	if err != nil {
		c.JSON(http.StatusForbidden, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusForbidden,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignee removed successfully"})
}

// SubmitCrimeReport godoc
// @Summary Submit a crime report
// @Description Public endpoint to submit a crime report
// @Tags public,reports
// @Accept json
// @Produce json
// @Param report body dto.ReportDTO true "Crime report details"
// @Success 201 {object} map[string]uint "reportId"
// @Failure 400 {object} map[string]string "Invalid report data"
// @Failure 500 {object} map[string]string "Server error"
// @Router /public/reports [post]
func (ctrl *CaseController) SubmitCrimeReport(c *gin.Context) {
	var reportDTO dto.ReportDTO
	if err := c.ShouldBindJSON(&reportDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid report data",
			Code:    http.StatusBadRequest,
		})
		return
	}

	report := &model.Report{
		Title:        reportDTO.Title,
		Description:  reportDTO.Description,
		Location:     reportDTO.Location,
		CivilID:      reportDTO.CivilID,
		Email:        reportDTO.Email,
		Name:         reportDTO.Name,
		ReportStatus: string(model.ReportStatusPending),
	}

	result, err := ctrl.caseService.SubmitCrimeReport(report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"reportId": result.ID})
}

func (ctrl *CaseController) GetEvidence(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid case ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	evidence, err := ctrl.caseService.GetEvidence(uint(caseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, evidence)
}

func (ctrl *CaseController) GetSuspects(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid case ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	suspects, err := ctrl.caseService.GetSuspects(uint(caseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, suspects)
}

func (ctrl *CaseController) GetVictims(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid case ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	victims, err := ctrl.caseService.GetVictims(uint(caseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, victims)
}

func (ctrl *CaseController) GetWitnesses(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid case ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	witnesses, err := ctrl.caseService.GetWitnesses(uint(caseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, witnesses)
}

func (c *CaseController) ExtractLinks(ctx *gin.Context) {
	// Parse case ID from URL parameter
	caseIDStr := ctx.Param("id")
	caseID, err := strconv.ParseUint(caseIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid case ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Call service to extract links
	links, err := c.caseService.ExtractLinksFromCase(uint(caseID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Return the extracted links
	ctx.JSON(http.StatusOK, gin.H{"links": links})
}

// GetCaseStatusByReportID godoc
// @Summary Get case status by report ID
// @Description Public endpoint to check a case status using report ID
// @Tags public,reports
// @Accept json
// @Produce json
// @Param reportId path string true "Report ID"
// @Success 200 {object} map[string]string "Report ID and status"
// @Failure 400 {object} map[string]string "Invalid Report ID"
// @Failure 404 {object} map[string]string "Report not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /public/reports/{reportId}/status [get]
func (ctrl *CaseController) GetCaseStatusByReportID(c *gin.Context) {
	reportID := c.Param("reportId")
	if reportID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Report ID is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	status, err := ctrl.caseService.GetCaseStatusByReportID(reportID)
	if err != nil {
		if err.Error() == "report not found" {
			c.JSON(http.StatusNotFound, dto.ErrorDTO{
				Message: "Report not found",
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reportId": reportID,
		"status":   status,
	})
}

func (ctrl *CaseController) UpdateCaseStatus(c *gin.Context) {
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	var statusUpdate struct {
		Status model.CaseStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status update data"})
		return
	}

	// Validate status value
	if statusUpdate.Status != model.StatusPending && 
	   statusUpdate.Status != model.StatusOngoing && 
	   statusUpdate.Status != model.StatusClosed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	// Get current user
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	userObj := user.(*model.User)

	// Get the case
	caseData, err := ctrl.caseService.GetCaseByID(uint(caseID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Case not found"})
		return
	}

	// Check if user is assigned to the case
	isAssigned := false
	assignees, err := ctrl.caseService.GetAssignees(uint(caseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check case assignment"})
		return
	}

	for _, assignee := range assignees {
		if assignee.ID == userObj.ID {
			isAssigned = true
			break
		}
	}

	if !isAssigned && userObj.Role == model.RoleOfficer {
		c.JSON(http.StatusForbidden, gin.H{"error": "Officer must be assigned to the case to update its status"})
		return
	}

	// Update the status
	caseData.Status = statusUpdate.Status
	result, err := ctrl.caseService.UpdateCase(caseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
