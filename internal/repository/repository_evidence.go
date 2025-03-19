package repository

import (
	"github.com/m7medVision/crime-management-system/internal/model"
	"gorm.io/gorm"
)

type EvidenceRepository struct {
	db *gorm.DB
}

func NewEvidenceRepository(db *gorm.DB) *EvidenceRepository {
	return &EvidenceRepository{db: db}
}

func (r *EvidenceRepository) Create(evidence *model.Evidence) error {
	return r.db.Create(evidence).Error
}

func (r *EvidenceRepository) GetByID(id uint) (*model.Evidence, error) {
	var evidence model.Evidence
	result := r.db.First(&evidence, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &evidence, nil
}

func (r *EvidenceRepository) Update(evidence *model.Evidence) error {
	return r.db.Save(evidence).Error
}

func (r *EvidenceRepository) SoftDelete(id uint) error {
	evidence, err := r.GetByID(id)
	if err != nil {
		return err
	}

	evidence.IsDeleted = true
	return r.Update(evidence)
}
