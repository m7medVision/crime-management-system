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

// CreateTextEvidence godoc
// @Summary Create text evidence
// @Description Add textual evidence to a case
// @Tags evidence
// @Accept json
// @Produce json
// @Param evidence body dto.CreateTextEvidenceDTO true "Text evidence details"
// @Success 201 {object} model.Evidence
// @Failure 400 {object} dto.ErrorDTO "Invalid evidence data"
// @Failure 401 {object} dto.ErrorDTO "Authentication required"
// @Failure 500 {object} dto.ErrorDTO "Server error"
// @Security BasicAuth
// @Router /evidence/text [post]
func (ctrl *EvidenceController) CreateTextEvidence(c *gin.Context) {
	var evidenceDTO dto.CreateTextEvidenceDTO
	if err := c.ShouldBindJSON(&evidenceDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid evidence data",
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
	evidence, err := ctrl.evidenceService.CreateTextEvidence(
		evidenceDTO.CaseID,
		userID,
		evidenceDTO.Content,
		evidenceDTO.Remarks,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, evidence)
}

// CreateImageEvidence godoc
// @Summary Upload image evidence
// @Description Upload an image as evidence for a case
// @Tags evidence
// @Accept multipart/form-data
// @Produce json
// @Param caseId formData int true "Case ID"
// @Param remarks formData string false "Optional remarks about the evidence"
// @Param image formData file true "Image file"
// @Success 201 {object} model.Evidence
// @Failure 400 {object} dto.ErrorDTO "Invalid request or not an image"
// @Failure 401 {object} dto.ErrorDTO "Authentication required"
// @Failure 500 {object} dto.ErrorDTO "Server error"
// @Security BasicAuth
// @Router /evidence/image [post]
func (ctrl *EvidenceController) CreateImageEvidence(c *gin.Context) {
	caseID, err := strconv.Atoi(c.PostForm("caseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid case ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	remarks := c.PostForm("remarks")
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Image file is required",
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
	evidence, err := ctrl.evidenceService.CreateImageEvidence(uint(caseID), userID, file, remarks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, evidence)
}

// GetEvidenceByID godoc
// @Summary Get evidence details
// @Description Retrieve details of a specific evidence item
// @Tags evidence
// @Accept json
// @Produce json
// @Param id path int true "Evidence ID"
// @Success 200 {object} model.Evidence
// @Failure 400 {object} dto.ErrorDTO "Invalid evidence ID"
// @Failure 404 {object} dto.ErrorDTO "Evidence not found"
// @Security BasicAuth
// @Router /evidence/{id} [get]
func (ctrl *EvidenceController) GetEvidenceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid evidence ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	evidence, err := ctrl.evidenceService.GetEvidenceByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorDTO{
			Message: "Evidence not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, evidence)
}

// GetEvidenceImage godoc
// @Summary Get evidence image
// @Description Stream an evidence image file
// @Tags evidence
// @Accept json
// @Produce image/*
// @Param id path int true "Evidence ID"
// @Success 200 {file} binary "Image file"
// @Failure 400 {object} dto.ErrorDTO "Invalid evidence ID"
// @Failure 500 {object} dto.ErrorDTO "Server error or not an image"
// @Security BasicAuth
// @Router /evidence/{id}/image [get]
func (ctrl *EvidenceController) GetEvidenceImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid evidence ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	object, size, contentType, err := ctrl.evidenceService.GetEvidenceImage(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
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

// UpdateEvidence godoc
// @Summary Update evidence
// @Description Update evidence remarks
// @Tags evidence
// @Accept json
// @Produce json
// @Param id path int true "Evidence ID"
// @Param evidence body dto.UpdateEvidenceDTO true "Updated evidence details"
// @Success 200 {object} model.Evidence
// @Failure 400 {object} dto.ErrorDTO "Invalid evidence ID or data"
// @Failure 401 {object} dto.ErrorDTO "Authentication required"
// @Failure 500 {object} dto.ErrorDTO "Server error"
// @Security BasicAuth
// @Router /evidence/{id} [put]
func (ctrl *EvidenceController) UpdateEvidence(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid evidence ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var updateDTO dto.UpdateEvidenceDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid update data",
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
	evidence, err := ctrl.evidenceService.UpdateEvidence(uint(id), userID, updateDTO.Remarks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, evidence)
}

// SoftDeleteEvidence godoc
// @Summary Soft delete evidence
// @Description Mark evidence as deleted (soft delete)
// @Tags evidence
// @Accept json
// @Produce json
// @Param id path int true "Evidence ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} dto.ErrorDTO "Invalid evidence ID"
// @Failure 401 {object} dto.ErrorDTO "Authentication required"
// @Failure 500 {object} dto.ErrorDTO "Server error"
// @Security BasicAuth
// @Router /evidence/{id} [delete]
func (ctrl *EvidenceController) SoftDeleteEvidence(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid evidence ID",
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
	err = ctrl.evidenceService.SoftDeleteEvidence(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Evidence soft deleted successfully"})
}

// HardDeleteEvidence godoc
// @Summary Permanently delete evidence
// @Description Permanently delete evidence (requires confirmation)
// @Tags evidence
// @Accept json
// @Produce json
// @Param id path int true "Evidence ID"
// @Param confirmation body dto.DeleteConfirmationDTO true "Deletion confirmation"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} dto.ErrorDTO "Invalid evidence ID or missing confirmation"
// @Failure 401 {object} dto.ErrorDTO "Authentication required"
// @Failure 500 {object} dto.ErrorDTO "Server error"
// @Security BasicAuth
// @Router /evidence/{id}/permanent [delete]
func (ctrl *EvidenceController) HardDeleteEvidence(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid evidence ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var confirmDTO dto.DeleteConfirmationDTO
	if err := c.ShouldBindJSON(&confirmDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Confirmation required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if confirmDTO.Confirmation != "CONFIRM_DELETE" {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid confirmation",
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
	err = ctrl.evidenceService.HardDeleteEvidence(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Evidence permanently deleted"})
}

// GetEvidenceAuditLogs godoc
// @Summary Get evidence audit logs
// @Description Retrieve audit logs for a specific evidence item
// @Tags evidence,audit
// @Accept json
// @Produce json
// @Param id path int true "Evidence ID"
// @Success 200 {array} model.AuditLog
// @Failure 400 {object} dto.ErrorDTO "Invalid evidence ID"
// @Failure 500 {object} dto.ErrorDTO "Server error"
// @Security BasicAuth
// @Router /evidence/{id}/audit [get]
func (ctrl *EvidenceController) GetEvidenceAuditLogs(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid evidence ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	logs, err := ctrl.evidenceService.GetEvidenceAuditLogs(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, logs)
}
