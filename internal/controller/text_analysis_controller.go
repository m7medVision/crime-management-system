package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m7medVision/crime-management-system/internal/service"
)

type TextAnalysisController struct {
	textAnalysisService *service.TextAnalysisService
}

func NewTextAnalysisController(textAnalysisService *service.TextAnalysisService) *TextAnalysisController {
	return &TextAnalysisController{textAnalysisService: textAnalysisService}
}

func (ctrl *TextAnalysisController) ExtractTopWords(c *gin.Context) {
	topWords, err := ctrl.textAnalysisService.ExtractTopWords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"top_words": topWords})
}
