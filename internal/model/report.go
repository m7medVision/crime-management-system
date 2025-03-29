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
	CivilID      string `gorm:"not null"`
	Email        string `gorm:"not null"`
	Name         string `gorm:"not null"`
	Role         string `gorm:"not null"`
	ReportStatus string `gorm:"not null"`
	Cases        []Case `gorm:"many2many:case_reports;"`
}

type ReportStatus string

const (
	ReportStatusOpen       ReportStatus = "open"
	ReportStatusInProgress ReportStatus = "in_progress"
	ReportStatusClosed     ReportStatus = "closed"
	ReportStatusResolved   ReportStatus = "resolved"
	ReportStatusRejected   ReportStatus = "rejected"
	ReportStatusPending    ReportStatus = "pending"
	ReportStatusEscalated  ReportStatus = "escalated"
	ReportStatusArchived   ReportStatus = "archived"
)
