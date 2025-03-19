package service

import (
	"strings"

	"github.com/m7medVision/crime-management-system/internal/repository"
)

type TextAnalysisService struct {
	textAnalysisRepo *repository.TextAnalysisRepository
}

func NewTextAnalysisService(textAnalysisRepo *repository.TextAnalysisRepository) *TextAnalysisService {
	return &TextAnalysisService{textAnalysisRepo: textAnalysisRepo}
}

func (s *TextAnalysisService) ExtractTopWords() (map[string]int, error) {
	evidences, err := s.textAnalysisRepo.GetAllTextEvidence()
	if err != nil {
		return nil, err
	}

	wordCount := make(map[string]int)
	stopWords := map[string]bool{
		"and": true, "the": true, "to": true, "of": true, "a": true, "in": true, "that": true, "is": true, "was": true, "he": true,
		"for": true, "it": true, "with": true, "as": true, "his": true, "on": true, "be": true, "at": true, "by": true, "i": true,
		"this": true, "had": true, "not": true, "are": true, "but": true, "from": true, "or": true, "have": true, "an": true, "they": true,
		"which": true, "one": true, "you": true, "were": true, "her": true, "all": true, "she": true, "there": true, "would": true, "their": true,
		"we": true, "him": true, "been": true, "has": true, "when": true, "who": true, "will": true, "no": true, "more": true, "if": true,
		"out": true, "so": true, "said": true, "what": true, "up": true, "its": true, "about": true, "into": true, "than": true, "them": true,
		"can": true, "only": true, "other": true, "new": true, "some": true, "could": true, "time": true, "these": true, "two": true, "may": true,
		"then": true, "do": true, "first": true, "any": true, "my": true, "now": true, "such": true, "like": true, "our": true, "over": true,
		"man": true, "me": true, "even": true, "most": true, "made": true, "after": true, "also": true, "did": true, "many": true, "before": true,
		"must": true, "through": true, "back": true, "years": true, "where": true, "much": true, "your": true, "way": true, "well": true, "down": true,
		"should": true, "because": true, "each": true, "just": true, "those": true, "people": true, "mr": true, "how": true, "too": true, "little": true,
		"state": true, "good": true, "very": true, "make": true, "world": true, "still": true, "see": true, "own": true, "men": true, "work": true,
		"long": true, "get": true, "here": true, "between": true, "both": true, "life": true, "being": true, "under": true, "never": true, "day": true,
		"same": true, "another": true, "know": true, "while": true, "last": true, "might": true, "us": true, "great": true, "old": true, "year": true,
		"off": true, "come": true, "since": true, "against": true, "go": true, "came": true, "right": true, "used": true, "take": true, "three": true,
		"himself": true, "few": true, "house": true, "use": true, "during": true, "without": true, "again": true, "place": true, "around": true, "however": true,
		"small": true, "found": true, "thought": true, "went": true, "say": true, "part": true, "once": true, "general": true, "high": true, "upon": true,
		"school": true, "every": true, "does": true, "got": true, "though": true, "left": true, "until": true, "children": true, "always": true, "city": true,
		"set": true, "put": true, "war": true, "home": true, "read": true, "hand": true, "large": true, "end": true, "open": true, "seemed": true,
		"next": true, "example": true, "began": true, "took": true, "sometimes": true, "run": true, "number": true, "course": true, "yet": true, "among": true,
	}

	for _, evidence := range evidences {
		words := strings.Fields(evidence.Content)
		for _, word := range words {
			word = strings.ToLower(word)
			if !stopWords[word] {
				wordCount[word]++
			}
		}
	}

	return wordCount, nil
}
