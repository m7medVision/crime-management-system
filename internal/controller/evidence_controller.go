package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/dto"
	"github.com/m7medVision/crime-management-system/internal/service"
	"github.com/m7medVision/crime-management-system/internal/model"
)

type EvidenceController struct {
	evidenceService *service.EvidenceService
}

func NewEvidenceController(evidenceService *service.EvidenceService) *EvidenceController {
	return &EvidenceController{evidenceService: evidenceService}
}

func (ctrl *EvidenceController) RecordEvidence(c *gin.Context) {
	var evidenceDTO dto.EvidenceDTO
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
	evidenceData := &model.Evidence{
		CaseID:    evidenceDTO.CaseID,
		Type:      evidenceDTO.Type,
		Content:   evidenceDTO.Content,
		ImagePath: evidenceDTO.ImagePath,
		Remarks:   evidenceDTO.Remarks,
		AddedByID: userID,
	}

	result, err := ctrl.evidenceService.RecordEvidence(evidenceData)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (ctrl *EvidenceController) GetEvidenceByID(c *gin.Context) {
	evidenceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evidence ID"})
		return
	}

	evidenceData, err := ctrl.evidenceService.GetEvidenceByID(uint(evidenceID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Evidence not found"})
		return
	}

	c.JSON(http.StatusOK, evidenceData)
}

func (ctrl *EvidenceController) GetEvidenceImageByID(c *gin.Context) {
	evidenceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evidence ID"})
		return
	}

	imageData, err := ctrl.evidenceService.GetEvidenceImageByID(uint(evidenceID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Evidence image not found"})
		return
	}

	c.JSON(http.StatusOK, imageData)
}

func (ctrl *EvidenceController) UpdateEvidence(c *gin.Context) {
	var evidenceDTO dto.EvidenceDTO
	if err := c.ShouldBindJSON(&evidenceDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evidence data"})
		return
	}

	evidenceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evidence ID"})
		return
	}

	evidenceData, err := ctrl.evidenceService.GetEvidenceByID(uint(evidenceID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Evidence not found"})
		return
	}

	evidenceData.Content = evidenceDTO.Content
	evidenceData.Remarks = evidenceDTO.Remarks

	result, err := ctrl.evidenceService.UpdateEvidence(evidenceData)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ctrl *EvidenceController) SoftDeleteEvidence(c *gin.Context) {
	evidenceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid evidence ID"})
		return
	}

	err = ctrl.evidenceService.SoftDeleteEvidence(uint(evidenceID))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Evidence soft deleted successfully"})
}
