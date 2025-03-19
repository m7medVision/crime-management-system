package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/service"
	"github.com/m7medVision/crime-management-system/internal/dto"
)

type ReportController struct {
	reportService *service.ReportService
}

func NewReportController(reportService *service.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

func (ctrl *ReportController) GenerateReport(c *gin.Context) {
	var reportDTO dto.ReportDTO
	if err := c.ShouldBindJSON(&reportDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report data"})
		return
	}

	report, err := ctrl.reportService.GenerateReport(&reportDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, report)
}

func (ctrl *ReportController) GetReportStatus(c *gin.Context) {
	reportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	status, err := ctrl.reportService.GetReportStatus(uint(reportID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
