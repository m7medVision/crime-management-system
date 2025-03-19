package repository

import (
	"github.com/m7medVision/crime-management-system/internal/model"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(auditLog *model.AuditLog) error {
	return r.db.Create(auditLog).Error
}

func (r *AuditLogRepository) GetByUserID(userID uint) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	result := r.db.Where("user_id = ?", userID).Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}
	return logs, nil
}
