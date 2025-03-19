package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/service"
)

type LinkExtractionController struct {
	linkExtractionService *service.LinkExtractionService
}

func NewLinkExtractionController(linkExtractionService *service.LinkExtractionService) *LinkExtractionController {
	return &LinkExtractionController{linkExtractionService: linkExtractionService}
}

func (ctrl *LinkExtractionController) ExtractLinks(c *gin.Context) {
	caseID := c.Param("id")

	links, err := ctrl.linkExtractionService.ExtractLinks(caseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"links": links})
}
