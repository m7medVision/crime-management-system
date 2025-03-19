package repository

import (
	"github.com/m7medVision/crime-management-system/internal/model"
	"gorm.io/gorm"
)

type TextAnalysisRepository struct {
	db *gorm.DB
}

func NewTextAnalysisRepository(db *gorm.DB) *TextAnalysisRepository {
	return &TextAnalysisRepository{db: db}
}

func (r *TextAnalysisRepository) GetAllTextEvidence() ([]model.Evidence, error) {
	var evidences []model.Evidence
	err := r.db.Where("type = ?", model.EvidenceTypeText).Find(&evidences).Error
	return evidences, err
}
