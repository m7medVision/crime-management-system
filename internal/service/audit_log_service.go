package service

import (
	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/repository"
)

type AuditLogService struct {
	auditLogRepo *repository.AuditLogRepository
}

func NewAuditLogService(auditLogRepo *repository.AuditLogRepository) *AuditLogService {
	return &AuditLogService{auditLogRepo: auditLogRepo}
}

func (s *AuditLogService) GetAdminLogs(userID uint) ([]model.AuditLog, error) {
	return s.auditLogRepo.GetByUserID(userID)
}
