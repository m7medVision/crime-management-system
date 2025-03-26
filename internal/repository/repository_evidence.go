package repository

import (
	"errors"

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
	result := r.db.Preload("Case").Preload("AddedBy").First(&evidence, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &evidence, nil
}

func (r *EvidenceRepository) Update(evidence *model.Evidence) error {
	return r.db.Save(evidence).Error
}

func (r *EvidenceRepository) SoftDelete(id uint) error {
	return r.db.Model(&model.Evidence{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *EvidenceRepository) HardDelete(id uint) error {
	result := r.db.Delete(&model.Evidence{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("evidence not found")
	}
	return nil
}

func (r *EvidenceRepository) ListByCaseID(caseID uint) ([]model.Evidence, error) {
	var evidence []model.Evidence
	err := r.db.Where("case_id = ? AND is_deleted = ?", caseID, false).
		Preload("AddedBy").Find(&evidence).Error
	return evidence, err
}

func (r *EvidenceRepository) CountByCaseID(caseID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Evidence{}).
		Where("case_id = ? AND is_deleted = ?", caseID, false).
		Count(&count).Error
	return count, err
}

func (r *EvidenceRepository) CreateAuditLog(auditLog *model.AuditLog) error {
	return r.db.Create(auditLog).Error
}

func (r *EvidenceRepository) GetAuditLogsForEvidence(evidenceID uint) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	err := r.db.Where("entity_type = ? AND entity_id = ?", "evidence", evidenceID).
		Preload("User").Find(&logs).Error
	return logs, err
}
