package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/dto"
	"github.com/m7medVision/crime-management-system/internal/service"
)

type ReportController struct {
	reportService *service.ReportService
}

func NewReportController(reportService *service.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

// GenerateCaseReport godoc
// @Summary Generate case report
// @Description Generate a PDF report for a case with all details
// @Tags reports
// @Accept json
// @Produce application/pdf
// @Param id path int true "Case ID"
// @Success 200 {file} binary "PDF report file"
// @Failure 400 {object} dto.ErrorDTO "Invalid case ID"
// @Failure 500 {object} dto.ErrorDTO "Server error"
// @Security ApiKeyAuth
// @Router /cases/{id}/report [get]
func (ctrl *ReportController) GenerateCaseReport(c *gin.Context) {
	// Parse case ID from URL parameter
	caseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorDTO{
			Message: "Invalid case ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Generate the PDF report
	pdf, err := ctrl.reportService.GenerateCaseReport(uint(caseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorDTO{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Set headers for PDF download
	fileName := "case_report_" + strconv.Itoa(caseID) + ".pdf"
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", strconv.Itoa(len(pdf)))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	// Write PDF data to response
	c.Data(http.StatusOK, "application/pdf", pdf)
}
