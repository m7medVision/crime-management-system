package model

import (
	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	ReportID     string `gorm:"uniqueIndex;not null"`
	Title        string `gorm:"not null"`
	Description  string `gorm:"type:text;not null"`
	Location     string `gorm:"not null"`
	ReportedByID uint   `gorm:"not null"`
	ReportedBy   User   `gorm:"foreignKey:ReportedByID"`
	Cases        []Case `gorm:"many2many:case_reports;"`
}
