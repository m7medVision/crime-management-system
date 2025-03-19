package repository

import (
	"github.com/m7medVision/crime-management-system/internal/model"
	"gorm.io/gorm"
)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(report *model.Report) error {
	return r.db.Create(report).Error
}

func (r *ReportRepository) GetByID(id uint) (*model.Report, error) {
	var report model.Report
	result := r.db.Preload("ReportedBy").First(&report, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &report, nil
}
