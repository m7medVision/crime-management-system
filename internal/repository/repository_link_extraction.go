package repository

import (
	"regexp"
	"strings"

	"github.com/m7medVision/crime-management-system/internal/model"
	"gorm.io/gorm"
)

type LinkExtractionRepository struct {
	db *gorm.DB
}

func NewLinkExtractionRepository(db *gorm.DB) *LinkExtractionRepository {
	return &LinkExtractionRepository{db: db}
}

func (r *LinkExtractionRepository) ExtractLinks(caseID uint) ([]string, error) {
	var caseData model.Case
	err := r.db.First(&caseData, caseID).Error
	if err != nil {
		return nil, err
	}

	links := extractLinksFromText(caseData.Description)
	return links, nil
}

func extractLinksFromText(text string) []string {
	var links []string
	re := regexp.MustCompile(`https?://[^\s]+`)
	matches := re.FindAllString(text, -1)
	for _, match := range matches {
		links = append(links, strings.TrimSpace(match))
	}
	return links
}
