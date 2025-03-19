package service

import (
	"regexp"
	"strings"

	"github.com/m7medVision/crime-management-system/internal/repository"
)

type LinkExtractionService struct {
	caseRepo *repository.CaseRepository
}

func NewLinkExtractionService(caseRepo *repository.CaseRepository) *LinkExtractionService {
	return &LinkExtractionService{caseRepo: caseRepo}
}

func (s *LinkExtractionService) ExtractLinks(caseID string) ([]string, error) {
	caseData, err := s.caseRepo.GetByCaseNumber(caseID)
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
