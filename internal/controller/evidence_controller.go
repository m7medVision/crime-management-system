package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/dto"
	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/service"
)

type EvidenceController struct {
	evidenceService *service.EvidenceService
}

func NewEvidenceController(evidenceService *service.EvidenceService) *EvidenceController {
	return &EvidenceController{evidenceService: evidenceService}
}

// CreateTextEvidence handles the creation of text evidence
func (ctrl *EvidenceController) CreateTextEvidence(c *gin.Context) {
	var evidenceDTO dto.CreateTextEvidenceDTO
	if err := c.ShouldBindJSON(&evidenceDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evidence data"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userID := user.(*model.User).ID
	evidence, err := ctrl.evidenceService.CreateTextEvidence(
		evidenceDTO.CaseID,
		userID,
		evidenceDTO.Content,
		evidenceDTO.Remarks,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, evidence)
}

// CreateImageEvidence handles the upload of image evidence
func (ctrl *EvidenceController) CreateImageEvidence(c *gin.Context) {
	caseID, err := strconv.Atoi(c.PostForm("caseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
		return
	}

	remarks := c.PostForm("remarks")
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userID := user.(*model.User).ID
	evidence, err := ctrl.evidenceService.CreateImageEvidence(uint(caseID), userID, file, remarks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, evidence)
}

// GetEvidenceByID retrieves evidence by ID
func (ctrl *EvidenceController) GetEvidenceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evidence ID"})
		return
	}

	evidence, err := ctrl.evidenceService.GetEvidenceByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Evidence not found"})
		return
	}

	c.JSON(http.StatusOK, evidence)
}

// GetEvidenceImage streams the evidence image
func (ctrl *EvidenceController) GetEvidenceImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evidence ID"})
		return
	}

	object, size, contentType, err := ctrl.evidenceService.GetEvidenceImage(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer object.Close()

	// Set necessary headers
	c.Header("Content-Disposition", "inline")
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", strconv.FormatInt(size, 10))

	// Stream the file to the client
	c.DataFromReader(http.StatusOK, size, contentType, object, nil)
}

// UpdateEvidence updates evidence remarks
func (ctrl *EvidenceController) UpdateEvidence(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evidence ID"})
		return
	}

	var updateDTO dto.UpdateEvidenceDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userID := user.(*model.User).ID
	evidence, err := ctrl.evidenceService.UpdateEvidence(uint(id), userID, updateDTO.Remarks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, evidence)
}

// SoftDeleteEvidence soft deletes evidence
func (ctrl *EvidenceController) SoftDeleteEvidence(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evidence ID"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userID := user.(*model.User).ID
	err = ctrl.evidenceService.SoftDeleteEvidence(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Evidence soft deleted successfully"})
}

// HardDeleteEvidence permanently deletes evidence with confirmation
func (ctrl *EvidenceController) HardDeleteEvidence(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evidence ID"})
		return
	}

	var confirmDTO dto.DeleteConfirmationDTO
	if err := c.ShouldBindJSON(&confirmDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Confirmation required"})
		return
	}

	if confirmDTO.Confirmation != "CONFIRM_DELETE" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid confirmation"})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	userID := user.(*model.User).ID
	err = ctrl.evidenceService.HardDeleteEvidence(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Evidence permanently deleted"})
}

// GetEvidenceAuditLogs retrieves audit logs for evidence
func (ctrl *EvidenceController) GetEvidenceAuditLogs(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evidence ID"})
		return
	}

	logs, err := ctrl.evidenceService.GetEvidenceAuditLogs(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
