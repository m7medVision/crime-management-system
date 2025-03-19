package dto

type TextAnalysisDTO struct {
	TopWords map[string]int `json:"top_words"`
}
